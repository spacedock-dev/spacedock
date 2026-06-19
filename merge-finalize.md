---
title: Merge-guard verb — spacedock merge guard <slug> (atomic mod-block set→invoke→clear→terminalize)
status: implementation
source: 0205 carve (2026-06-17, captain "stamp them") — index DoD candidate "merge-finalize"; 2y MERGED unblocks it.
score: 0.5
sprint: 0221-layered-fo
group: verb-core
sprint-readiness: ready
id: mzmc0dgkq1nbazxx8j0mnfn6
started: 2026-06-17T20:55:24Z
worktree: .worktrees/spacedock-ensign-merge-finalize
---

`spacedock merge guard <slug>` enforces the mod-block set→invoke→clear→terminalize-only-after sequence atomically (set before the hook, detect completion by state delta, clear in a standalone `--set`, terminalize only after) — so the single highest Haiku merge risk (combining / skipping / reordering those steps) is owned by the binary. The status tool's existing refusal of terminal-with-mod-block-set is the backstop.

## Problem

The terminal merge ceremony (`fo-merge-core.md` Merge-and-Cleanup steps 1–9, and the Ship-Local variant) is a fixed ordered sequence of `status --set` / `--archive` calls the FO follows by hand:

1. set `mod-block=merge:{mod_name}` BEFORE invoking the merge hook,
2. invoke the hook (FO prose / host-specific; carries captain-approval judgment for `pr-merge`),
3. detect hook completion by the state delta (`pr` now set ⇒ blocked; mod-block still set + no external wait ⇒ completed),
4. if completed, clear `mod-block` in a STANDALONE `--set` (the guard refuses combining `mod-block=` with terminal fields),
5. terminalize (`status=done verdict={v} completed`) only AFTER the clear,
6. archive.

Every step is mechanical and judgment-free; the load is in *how many ordering rules*, not in *whether any rule is right* — exactly what a binary absorbs. The w4 haiku-loop spike (PR #393, PASSED) held `«merge»` mechanically across N=3 on the LINEAR happy path, but the 0205 carve flags it must-build because the happy path exercises none of the failure surface: a weak FO that hits a blocked hook, a rejected verdict, or `merge: local` can combine the clear with the terminalize, skip the standalone clear, or terminalize before the mod-block lands. The `status --set`/`--archive` guards (`runSet`/`runArchive` in `internal/status`) REFUSE each of those — proven by `merge_policy_guard_test.go` — but they are a backstop the FO trips into, not a verb that drives the sequence right the first time. There is no single call an FO can make that performs the whole ordered ceremony so it CANNOT be reordered or collapsed.

A related defect (`team-mode-verdict-omission`, backlog) shows the cost concretely: a team-mode headless `-p` FO non-deterministically sets `status: done` + archives + commits but OMITS `verdict:` ~1/3 of runs, because the multi-turn ceremony's verdict write is sometimes not reached before the cycle ends. A single binary call that writes status + verdict + completed in ONE `--set` (which the verb already must do, since the guard refuses a verdict-less finalize) makes the verdict write atomic with the status write — it cannot land one without the other.

## Proposed approach

Ship `spacedock merge guard <slug>` as a new cobra subcommand under the `workflow` group (a `merge` command with a `guard` subcommand, alongside `status` / `new` / `state` in `internal/cli/cli.go`), routing to a handler in `internal/status` that orchestrates the existing proven `--set` / `--archive` paths. The verb OWNS the ordered state-transition sequence; it does NOT invoke the merge hook, make the merge verdict, or run the host agent-teardown. It is the L1 mechanical envelope around the FO's L3 hook-invocation and verdict decision.

**Behavior (the state machine the verb drives, computed from the entity's current frontmatter + workflow merge policy):**

- **Inputs:** `<slug>` (positional), `--workflow-dir DIR`, `--verdict {passed|rejected}` (the FO/captain's decision, passed in — NOT computed), `--json`/`--quiet`. The merge policy (`merge:` README key, default `pr`) and the registered merge hooks are read from the workflow, not flagged.
- **Phase A — arm (idempotent), `merge: local` only:** under `merge: local`, if no `mod-block` is set and a merge hook is registered, set `mod-block=merge:{mod_name}` (its own `--set`, committed by the caller's state-commit step) and EXIT signalling `armed` (action=`invoke-hook`, naming the hook the FO must now invoke). The verb does not invoke the hook itself — that step carries the `pr-merge` captain-approval guardrail (outward-facing push) and is host/FO prose. Re-running `guard` after the hook is the resume path. Under `merge: pr` the verb does NOT auto-arm: arming opens a PR (an outward-facing, captain-approval-gated step the FO owns), and an empty-`mod-block` / empty-`pr` entity is indistinguishable by frontmatter from one that SKIPPED the ceremony — so a policy-blind arm could not honor AC-5. Under `merge: pr`, `guard --verdict passed` on such an entity goes straight to a finalize attempt that the merge-hook guard refuses (`cannot advance to terminal`, exit 1), and the verb resumes a correctly-armed `merge: pr` entity at Phase B (the hook has set `pr`). The Phase-C finalize atomicity is identical under both policies; only the auto-arm is gated.
- **Phase B — detect completion by state delta (re-run after the hook):** read the post-hook frontmatter.
  - `pr` is now set ⇒ the hook BLOCKED (PR pending). Leave `mod-block` set, EXIT signalling `blocked` (action=`await-pr`), do NOT clear/terminalize. (Matches fo-merge-core step 3a/4.)
  - `mod-block` still set, `pr` empty, no external-wait ⇒ COMPLETED. Proceed to Phase C.
  - `verdict=rejected` (passed via `--verdict rejected`) ⇒ the ceremony is vacuous (a rejected entity never merges). Skip the pr-requirement; still clear an in-flight mod-block standalone, then terminalize+archive.
- **Phase C — finalize (ordered, atomic per step):**
  1. clear `mod-block` in a STANDALONE `--set` (never combined with terminal fields),
  2. terminalize in ONE `--set`: `status={terminal} verdict={--verdict} completed` (status+verdict+completed together — the verdict-omission fix),
  3. archive.
  Each underlying `--set`/`--archive` is the already-proven guarded path; the verb's contribution is EMITTING them in this order and refusing to proceed to a later step if an earlier one's guard refuses (it propagates the guard's exit 1 + stderr verbatim and stops — it does not `--force` past a refusal).
- **Ship-Local vs PR:** the verb reads the `merge:` policy. Under `merge: local` Phase B treats a cleared/clearable mod-block as completion (no `pr` sentinel required — the policy exempts the pr-requirement, matching the guard). Under default `merge: pr` with a registered hook, Phase B requires either `pr` set (blocked) or `verdict=rejected` before it will finalize; otherwise it EXITS signalling the FO must record the merge (the `local-merge:{sha}` sentinel path stays the hook's/FO's, set before re-running `guard`). The verb does NOT compute or set the `local-merge` sentinel SHA (that requires the landed merge commit, an FO/git step).
- **State-delta completion detection** is the one genuinely-new logic: a pure read of `pr`/`mod-block`/`verdict` from `ParseFrontmatter` plus the policy, deciding armed/blocked/completed. No new mutation primitive.

**Doc diff (user-visible — new command surface):** add `merge` to the `verbs` list in both completion scripts (`bashCompletion`, `zshCompletion` in `cli.go`) and a `merge) compadd/compgen -- guard` case; add the `merge guard` line to the grouped help (`help.go` workflow group). Concrete before/after wording recorded in the Test plan's doc-diff item below. The fo-merge-core prose-function `«merge.guard»` flips from guillemets to backticks pointing at this verb (owned by `prose-function-restructure`, NOT this task — this task ships the verb; the restructure flips the notation).

## Out of scope

- **The merge VERDICT and any judgment.** The verb enforces the SEQUENCE; the PASS/REJECT decision is passed in via `--verdict`. Routing the verdict to L3 is `fo-tier-delegation`.
- **Invoking the merge hook.** Phase A stops at `armed`; the FO/hook invokes (the `pr-merge` captain-approval push guardrail is judgment + outward-facing). The verb resumes at Phase B.
- **Host agent-teardown (Merge-and-Cleanup step 10).** The `TERMINAL_TEARDOWN_BOUNDED` settle/cap/marker ceremony (`claude-fo-merge.md`) is multi-turn, team-tool-driven, host-specific, and runs AFTER archive. It stays the FO's runtime-adapter responsibility. The verb finalizes on-disk state (through archive) only.
- **Worktree removal / `worktree safe-remove` companion.** `git worktree remove` (no `--force`) + the untracked-files audit (`fo-merge-core` Worktree-removal safety) is a git operation against the CODE worktree, separable from the state-finalize sequence. AC-BOUNDARY DECISION: ship `merge guard` (state sequence) THIS task; a `merge worktree-remove` / `worktree safe-remove` companion is a follow-up only if the spike on THIS verb shows the FO still botches the git removal. The w4 spike did not exercise worktree removal (single-root fixture), so there is no must-build signal yet — recorded as a deferred companion, not folded in. The `merge` command group leaves room for the companion subcommand without re-shaping the surface.
- **The `local-merge:{sha}` sentinel write.** Requires the landed merge commit SHA (an FO/git step); the verb signals the FO must record it, but does not compute it.
- **Setting the `pr:` field.** That is the hook's blocking signal; the verb reads it, never writes it.

## Acceptance criteria

Each AC names a property of the finished entity and how it is verified. All proofs drive the real binary against a fixture workflow and assert exit code + stdout/stderr (the `merge_policy_guard_test.go` golden-envelope idiom); none is a prose-grep over an instruction file.

**AC-1 — The happy-path ceremony, driven by one verb invocation per phase, lands the entity terminal + archived with `verdict` recorded.**
On a `merge: local` fixture with a registered merge hook, `merge guard <slug> --verdict passed` armed (Phase A sets mod-block, exits `armed`), then re-run after the hook completes finalizes: the entity ends `status={terminal}`, `verdict=passed`, archived under `_archive/`. Verified by: a Go test (`internal/status`) driving the verb via `runNative`, asserting exit 0 and reading the archived entity's frontmatter for `status`/`verdict`. Seeded by the ideation spike below (which proved the underlying `--set`/`--archive` sequence on the real binary).

**AC-2 — The verb refuses to terminalize while `mod-block` is set, and never combines the clear with the terminalize.**
The verb's finalize emits the standalone-clear `--set` and the terminalize `--set` as SEPARATE calls; an adversarial attempt to make it collapse them (or terminalize with mod-block live) propagates the underlying guard's exit 1 and the entity is unchanged. Verified by: a Go test asserting that after Phase A (`armed`, mod-block set), a `guard` run that should finalize first clears standalone — proven by the archived/terminal entity's git-or-frontmatter history showing the clear resolved before terminal fields — and by a negative test that an injected combined call (the guard the verb relies on) exits 1 with the `combined mod-block clear with terminal transition` stderr. Backstop pinned by the existing `TestMergeLocalCombinedClearAndTerminalizeRefused`.

**AC-3 — A blocked hook (`pr` set) leaves the entity un-terminalized with mod-block intact, and the verb signals `blocked`.**
On a `merge: pr` fixture, after the hook sets `pr=#N` (PR pending), `merge guard <slug>` detects the state delta, exits signalling `blocked`/`await-pr`, leaves `mod-block` set, and does NOT terminalize or archive. Verified by: a Go test on the `merge-pr-workflow` fixture asserting the verb's exit/signal and that the entity frontmatter still shows the non-terminal status + intact mod-block + set `pr`.

**AC-4 — `status` + `verdict` + `completed` are written in ONE `--set` (the verdict-omission atomicity fix).**
The terminalize step is a single `--set status={terminal} verdict={--verdict} completed`; there is no reachable verb path that sets `status=done` without `verdict` in the same mutation. Verified by: a Go test inspecting the verb's emitted mutation (or the post-finalize frontmatter where both land together), plus a negative test that `merge guard` WITHOUT `--verdict` exits non-zero before any terminal write (the verb requires the verdict up front), so a verdict-less finalize is unreachable. This is the structural fix for `team-mode-verdict-omission`.

**AC-5 — Under default `merge: pr` with a registered hook and no `pr`/no `--verdict rejected`, the verb refuses to finalize (the backstop is honored, not bypassed).**
The verb never passes `--force` to route around the merge-hook guard; if asked to finalize an entity that skipped the ceremony, it propagates the guard's `cannot advance to terminal` refusal (exit 1) and stops. Verified by: a Go test on `merge-pr-workflow` 020-no-sentinel asserting `merge guard … --verdict passed` (with no `pr`, no arm) exits 1 with the merge-hook refusal in stderr and the entity unchanged. Backstop pinned by `TestMergePrDefaultNoSentinelStillRefuses`.

**AC-6 — `verdict=rejected` finalizes without a `pr` (the rejected entity never merged), but still clears an in-flight mod-block standalone.**
`merge guard <slug> --verdict rejected` on an entity with empty `pr` finalizes to terminal+rejected+archived without `--force`; if a mod-block was in flight it is cleared in its own `--set` first. Verified by: a Go test on `merge-pr-workflow` (the `040-rejected` / `050-rejected-pending` fixtures) asserting exit 0 for the clean-rejected case and that the rejected-with-mod-block case clears the block before terminalizing. Backstop pinned by `TestRejectedVerdictArchiveMatchesSet` + `TestRejectedVerdictModBlockPendingArchiveRefuses`.

**AC-7 — The new command surface is discoverable.**
`spacedock merge guard --help` renders usage; `merge` appears in the grouped help workflow group and in the bash/zsh completion verb lists with a `guard` subcommand completion. Verified by: a Go test asserting `merge --help` / `merge guard --help` exits 0 with usage text, and the completion-script golden includes `merge` + the `merge) … guard` case.

## Test plan

All tests are Go unit/behavior tests in `internal/status` (verb handler) and `internal/cli` (routing/help/completion), driving the real binary against fixtures — the smallest surface that proves each claim, no live workflow needed (the live Haiku drive is `haiku-drive-validation`'s job, not this verb's).

- **Fixtures:** reuse the existing `merge-local-workflow` (AC-1/2/4) and `merge-pr-workflow` (AC-3/5/6) testdata — they already declare a registered `## Hook: merge`, the `merge: local` vs default `merge: pr` policies, and the no-sentinel / pending / rejected / rejected-pending entity states the ACs need. A small fixture addition: an entity carrying `pr=#N` to exercise the AC-3 blocked path (the existing fixtures have no PR-pending entity).
- **Verb handler tests (`internal/status`):** golden-envelope tests (the `assertMergeGolden` / `runNative` idiom) for each phase — armed, blocked, completed-happy, rejected, backstop-refusal-propagated. Assert exit code, stdout signal (`armed`/`blocked`/finalized), stderr (refusal text propagated verbatim from the underlying guard), and the resulting on-disk frontmatter / archive location.
- **Routing/help/completion tests (`internal/cli`):** `merge`/`merge guard` resolve to the handler; `--help` renders; unknown `merge bogus` exits 2; completion golden includes `merge`.
- **Negative/adversarial:** verdict-less finalize unreachable (AC-4); the verb never emits `--force` (assert the emitted call set, or that a guard refusal is propagated rather than bypassed) (AC-5).
- **Cost/complexity:** LOW. No new mutation primitive — the verb composes proven `--set`/`--archive` paths; the only new logic is the frontmatter-read state-delta classifier. ~6–8 new Go tests + one fixture entity. No live/credentialed lane.
- **Doc-diff item (concrete before/after):**
  - `cli.go` `bashCompletion` verbs string: `verbs="claude codex pi install doctor status new state completion dispatch --version --help"` → add `merge`: `verbs="claude codex pi install doctor status new state merge completion dispatch --version --help"`, and add a case `merge) COMPREPLY=( $(compgen -W "guard" -- "$cur") ) ;;`.
  - `cli.go` `zshCompletion` verbs: same `merge` insertion; add `merge) compadd -- guard ;;`.
  - `help.go` workflow group: add a `merge   Run the terminal merge-finalize ceremony for an entity` line alongside `status`/`new`/`state` (exact wording finalized against the existing help table during implementation).

## Spike — riskiest mechanism (DONE in ideation, seeds the first test)

**Riskiest mechanism:** the atomic set→invoke→clear→terminalize sequence with completion-detection-by-state-delta, AND that the existing guards are a real backstop catching every skip / reorder / collapse.

**Exercised (2026-06-17, this ideation, on the real `${SPACEDOCK_BIN}` against a copy of the `merge-local-workflow` + `merge-pr-workflow` fixtures):** hand-drove the full sequence and the three backstop refusals.

- Backstop A — terminalize on a `merge: pr` workflow with no ceremony: REFUSED, exit 1, `cannot advance to terminal — workflow has merge hook(s) [local-merge] that have not run`.
- Backstop B — terminalize with `mod-block` set: REFUSED, exit 1, `has pending mod-block (merge:local-merge). Clear mod-block in a separate --set call`.
- Backstop C — combined `mod-block=` + `status=done verdict=passed` in one call: REFUSED, exit 1, `combined mod-block clear with terminal transition`.
- Happy path — `set mod-block` → standalone `mod-block=` clear → `status=done verdict=passed completed` → `--archive`, in order: all exit 0; entity archived under `_archive/` with `status: done` + `verdict: passed`.

**Conclusion:** the verb's whole job is mechanizable on the EXISTING surface — it emits these exact ordered `--set`/`--archive` calls and propagates (never bypasses) the guard's refusal. The only new code is the frontmatter-read state-delta classifier (armed/blocked/completed over `ParseFrontmatter` + policy), which is low-risk. The spike output IS the AC-1/AC-2/AC-5 first test. No further spike needed: the verb relies on the proven `updateFrontmatter`/`runSet`/`runArchive` guards (`merge_policy_guard_test.go` green) + `ParseFrontmatter` + `resolveMergePolicy`/`scanMods` + the cobra routing pattern (`internal/cli`); all are shipped and tested.

## Stage Report: ideation

- DONE: Design `spacedock merge guard <slug>` — the atomic mod-block set→invoke→clear→terminalize-only-after sequence — against the REAL post-2y host-neutral merge core and the pr-merge / Ship-Local ceremonies; verb owns the sequence so a weak FO cannot combine/skip/reorder; behavior-first inputs/exit-codes/state-delta/on-disk effects.
  Body `## Proposed approach` — three-phase state machine (arm / detect-by-state-delta / finalize) against `fo-merge-core.md` steps 1-9 + Ship-Local; inputs/signals/exits specified.
- DONE: Resolve the AC BOUNDARY: what `merge guard` owns vs an optional `merge ceremony` / `worktree safe-remove` companion; whether it folds the `team-mode-verdict-omission` (re) atomicity.
  `## Out of scope` — verb owns state-sequence-through-archive; hook-invocation/verdict/step-10 teardown/worktree-removal/local-merge-sentinel stay OUT. Worktree-safe-remove deferred (no must-build signal — w4 single-root, did not exercise it). Verdict-omission IS folded (AC-4: status+verdict+completed in one `--set`).
- DONE: State the riskiest mechanism explicitly and SPIKE it (smallest exercise firing the sequence + backstop refusal) OR record "no spike needed".
  `## Spike` — hand-drove the full sequence + 3 backstop refusals on the real `${SPACEDOCK_BIN}` against the merge-local/merge-pr fixtures; all behaved as designed. Spike DONE, seeds AC-1/2/5 first test.
- DONE: Behavior-first oracle-based ACs + test plan (never a prose-grep); map «merge.guard» to a backtick verb; propose the doc-diff if user-visible. Verdict/judgment stay OUT.
  AC-1..7 each cite a Go test driving the real binary (exit code + on-disk frontmatter/archive), backstops pinned to existing `merge_policy_guard_test.go` tests; `## Test plan` names fixtures + the concrete completion/help doc-diff. `«merge.guard»` notation flip assigned to `prose-function-restructure` (this task ships the verb).

### Summary

Designed `spacedock merge guard <slug>` as an L1 mechanical envelope: a new `merge` cobra subcommand routing to an `internal/status` handler that EMITS the proven `--set`/`--archive` ceremony calls in the mandatory order (arm mod-block → detect hook completion by state delta → standalone clear → atomic terminalize → archive), propagating — never bypassing — the existing guard refusals. Key decisions: the merge VERDICT and hook-invocation stay OUT (passed in via `--verdict`; the captain-approval push guardrail is the FO's), step-10 agent teardown and worktree-removal are deferred companions (no w4 must-build signal yet), and the `team-mode-verdict-omission` fix is FOLDED by writing status+verdict+completed in one `--set` (AC-4). An ideation spike hand-drove the full sequence + all three backstop refusals on the real binary, confirming the verb is fully mechanizable on the existing surface (only new logic: a low-risk frontmatter-read state classifier) — no further spike needed. COMPLEX (atomicity + ceremony boundary) — expect a staff review before the gate.

## Stage Report: implementation

- DONE: Finalize is ordered and atomic: a STANDALONE `mod-block` clear `--set`, then ONE `--set status={terminal} verdict={--verdict} completed` (never combined) — the verb propagates the underlying guard's exit-1 refusal verbatim and never `--force`s past it (AC-2, AC-5).
  `internal/status/merge.go` `finalize()` emits the standalone clear then the single terminalize `--set` via `emitSet` (force always false); `TestMergeGuardPrNoSentinelRefuses` asserts the `cannot advance to terminal` refusal reaches exit 1 with the entity unchanged; `TestMergeLocalCombinedClearAndTerminalizeRefused` (existing) is the backstop.
- DONE: The state-delta classifier routes armed/blocked/completed/rejected over `pr`/`mod-block`/`verdict` + `merge:` policy: blocked (`pr` set) leaves mod-block intact and un-terminalized (AC-3); `--verdict rejected` finalizes without a `pr` but still clears an in-flight mod-block standalone first (AC-6).
  `MergeGuard()` switch in `merge.go` (pr→blocked, rejected→finalize, merge:local+no-mod-block→arm, else→finalize); `TestMergeGuardBlockedOnPR`, `TestMergeGuardRejectedFinalizesNoPR`, `TestMergeGuardRejectedClearsModBlockFirst` all green; AC-3 fixture `070-pr-pending.md` added.
- DONE: A verdict-less finalize is structurally unreachable: `merge guard` without `--verdict` exits non-zero before any terminal write (AC-4, the team-mode-verdict-omission fix).
  The `--verdict` gate in `MergeGuard()` fires before workflow resolution; terminalize writes status+verdict+completed in one `fieldUpdate` slice; `TestMergeGuardVerdictOmissionUnreachable` asserts exit 1 + entity untouched.
- DONE: Resolved a policy-blindness contradiction between Phase A (prose) and AC-5 in favor of the AC, confirmed with the FO. An empty-`mod-block` / empty-`pr` entity is indistinguishable as "armable" vs "ceremony-skipping" by frontmatter, so a policy-blind arm cannot satisfy AC-5. The auto-arm is gated to `merge: local` (fully in-process); under `merge: pr` the verb routes an unarmed `--verdict passed` to the merge-hook guard's refusal. Phase A body prose rewritten to state the gating explicitly.
  `MergeGuard()` arm case is `policy == mergeLocal && modBlock == "" && hookRegistered`; pinned by `TestMergeGuardArmIsPolicyGated` (committed in `b882a058`), which asserts the inversion on the IDENTICAL precondition (merge:local→exit 0 armed; merge:pr→exit 1 refused, no arm). Mutation self-check (re-run against the committed test): dropping the `merge: local` gate to a policy-blind arm turns the merge:pr leg RED (`must REFUSE (exit 1, not arm), got 0 ... armed`); restoring `merge.go` from git returns it green and `git status` clean — so the pin is load-bearing, not tautological.

### Summary

Shipped `spacedock merge guard <slug>` as `internal/status/merge.go` (handler) + a `merge`/`guard` cobra command, grouped-help row, and bash/zsh completion entries in `internal/cli`. The verb composes the existing proven `runSet`/`runArchive` guarded paths through `emitSet` (never `--force`), so the only new logic is the frontmatter-read state classifier. One load-bearing reading was escalated to and confirmed by the FO, and the Phase A body prose was rewritten to remove the contradiction: the auto-arm fires ONLY under `merge: local` (under `merge: pr` the PR-opening arm+hook stays the FO's, and a no-pr `--verdict passed` finalize is left to hit the merge-hook guard's refusal — AC-5); Phase-C finalize atomicity is identical under both policies. All 7 ACs are exercised by 10 status-handler tests (including the `TestMergeGuardArmIsPolicyGated` inversion pin) + 6 CLI routing/help/completion/e2e tests driving the real binary against the fixtures; full `go test ./...` is green, gofmt + vet clean.
