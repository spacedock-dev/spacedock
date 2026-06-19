---
title: Make the shipped state/merge verbs the operative contract path + bind «fn» bodies to their verbs (oracle)
source: '0221-layered-fo rework (2026-06-19): validated findings 1+2 — the contract does not operatively USE the shipped verbs. `### Split-Root State Sync` names an abstract "status tool" not `spacedock state commit`, and `«state.commit»`s effect restates the hand git sequence; `«merge.guard»`s prose claims it "invoke[s] the registered merge hook" / runs "as one call" / "default-merge[s]", all FALSE of the shipped re-entrant partial envelope (merge.go doc-comment: "It does NOT invoke the merge hook"). Bundles the contract rewire (A) with a routing oracle (B) written test-first.'
status: validation
score: 0.75
sprint: 0221-layered-fo
group: foundation
id: 4asxw7kxvdzdtf87w9rjkxwx
started: 2026-06-19T22:53:22Z
worktree: .worktrees/spacedock-ensign-verbs-operative-rewire
mod-block: merge:pr-merge
---

Complete the vertical slice `rgq`/`mz`/`czw` shipped empty: the verbs exist and route, but the contract instructs the hand sequence and mis-describes the merge verb. Make the verb the operative path AND guard it mechanically. (Per the captain's call, A and B are ONE bundled slice, not two layer-tasks.)

## Premise reconciliation (resolved in ideation — the seed was partly stale)

Three findings from reading main reshape the seed; the captain should see them before the gate.

**F1 — czw (#402) ALREADY shipped the routing oracle.** `internal/cli/prose_function_routing_test.go::TestProseFunctionNotationBindsToRouting` walks the three FO cores, extracts every `→ **shipped**`/`→ **prose**` migration target, invokes the named verb's argv against the compiled root, and asserts shipped verbs route while prose targets are rejected. It ships with two paired controls: `TestProseFunctionRoutingGuardFailsOnViolation` (RED control — plants `spacedock merge teleport`, an un-routed backtick, and a guillemet whose target *does* route) and `TestProseFunctionRoutingOracleDiscriminates` (non-vacuity). All GREEN on main today (`go test ./internal/cli -run TestProseFunctionNotationBindsToRouting` → ok). **Task-B as a NEW from-scratch routing oracle is already built.** What is NOT yet covered is the «state.commit» binding being a backtick-shipped target — see A2 below; once A1/A2 land, the existing oracle exercises it for free.

**F2 — the routing oracle is structurally BLIND to the «merge.guard» prose defect.** The seed says the oracle should be "RED on today's overstated «merge.guard»." It cannot be. The oracle reads only the `→ **shipped**: spacedock merge guard` line and checks the *verb routes* — it never reads the effect-bullet prose. The overstated claims ("invoke the registered merge hook", "default-merge", "as one call") live in lines 7 and 11, which the oracle does not parse. So "RED on the overstated prose" is mechanically incoherent for a routing oracle. The seed's own stated ALTERNATIVE — "a mutation pointing an «fn» at a non-existent verb" — IS the achievable non-vacuous RED, and it is ALREADY the shipped `TestProseFunctionRoutingGuardFailsOnViolation` (plants `merge teleport`). **The prose-accuracy defect (A3) and the routing oracle (B) are orthogonal mechanisms; the oracle cannot guard A3.** This slice does NOT manufacture a contrived prose-grep to fake a RED; it ships the prose fix (A3) and relies on the *already-present, already-non-vacuous* routing controls.

**F3 — `deferred_tier_absence_test.go` is NOT on main.** It exists only in the unmerged sibling branch `spacedock-ensign/strip-deferred-tier-vocabulary` (which ADDs it, +113 lines, branched off c25fee26). "B REPLACES the deferred-tier absence-grep" is a false coupling: there is nothing on main to replace, and the routing oracle that would have "replaced" it already shipped independently in #402. **Resolution (checklist item 3): B does NOT delete `deferred_tier_absence_test.go`.** The two tests are unrelated — the deferred-tier grep guards member-72 tier vocabulary out of the cores (a different concern entirely); the routing oracle guards «fn»↔verb routing. Neither subsumes the other. The "replaces" language in the seed is retired here as a mis-mapping; if the captain still wants the deferred-tier grep retired, that is the strip-deferred-tier branch's own scope, not this slice's.

## Proposed approach (firm)

**A — rewire the contract prose (dispatched worker on the implementation worktree — scaffolding guardrail under `skills/first-officer/references/`, NOT an FO live edit):**

The three rewrites, with exact before/after wording, are in the **Concrete rewrites** section below. In summary:
- **A1** — `### Split-Root State Sync` "Preferred" bullet (`first-officer-shared-core.md:189`): name `spacedock state commit <slug>` as the preferred path; demote the hand `git -C add/commit` to the named degraded/no-origin fallback.
- **A2** — `«state.commit»` effect bullet (`first-officer-shared-core.md:178`): delegate to the verb; stop restating the full hand git sequence as the effect (keep the hand sequence only as the documented fallback the verb automates).
- **A3** — `«merge.guard»` intro (line 7) + effect bullet (line 11) in `fo-merge-core.md`: replace the three false claims with the verb's actual re-entrant-partial-envelope behavior, sourced from `internal/status/merge.go` (doc-comment line 13: "It does NOT invoke the merge hook"; `signalArmed` line 213: "invoke the … merge hook, then re-run `merge guard`").

**B — the «fn»→shipped-verb routing oracle:** already shipped (#402). This slice's B-obligation is reduced to: (1) confirm the existing oracle stays GREEN after the A rewrites land (the «state.commit» backtick must keep routing; no «fn» migration-target line is changed by A — A touches effect/intro prose, not the `→` lines), and (2) confirm the existing RED control (`TestProseFunctionRoutingGuardFailsOnViolation`) and discriminator remain the non-vacuity proof. No new oracle test is authored — authoring a duplicate would violate YAGNI and the no-duplication rule.

## Concrete rewrites (before/after — load-bearing exact wording)

### A1 — Split-Root State Sync "Preferred" bullet (`first-officer-shared-core.md`, currently line 189)

BEFORE:
```
- **Preferred — tool-managed atomic state commits.** When the status tool owns `add`+`commit` under a lock, route through it.
- **Fallback — path-scoped commits per writer.** `git -C {state_checkout} add {entity_path} && git -C {state_checkout} commit -m "…" -- {entity_path}`. Never a bare `git add -A` or bare `git commit`. Retry on `index.lock` contention after ~2s.
```
AFTER:
```
- **Preferred — `spacedock state commit <slug>`.** The verb resolves the entity's path under the state checkout and runs the path-scoped commit → push → on-reject `pull --rebase` → re-push sequence (and the rebase-conflict halt) atomically. Invoke it directly; do not hand-roll the git sequence.
- **Fallback (no verb available / degraded) — path-scoped commits per writer.** `git -C {state_checkout} add {entity_path} && git -C {state_checkout} commit -m "…" -- {entity_path}`. Never a bare `git add -A` or bare `git commit`. Retry on `index.lock` contention after ~2s.
```
Rationale: the verb (`internal/cli/state_sync.go::runStateCommit`) already implements exactly the commit+sync+halt the prose described abstractly. The hand sequence stays as the named fallback (it is genuinely what the verb automates and what a no-`spacedock`-binary degraded FO falls back to).

### A2 — «state.commit» effect bullet (`first-officer-shared-core.md`, currently line 178)

BEFORE:
```
- **effect:** path-scoped commit + sync per **Split-Root State Sync** below — `git -C {state_checkout} add {entity_path} && git -C {state_checkout} commit -m "…" -- {entity_path}`, never a bare `git add -A` / `git commit`; then push; on push rejection `pull --rebase` then re-push; retry on `index.lock` contention after ~2s.
```
AFTER:
```
- **effect:** delegate to `spacedock state commit <slug>` — it path-scoped-commits the entity and runs the full sync (push; on push rejection `pull --rebase` then re-push; rebase-conflict halt) per **Split-Root State Sync** below. The hand `git -C add/commit … push/pull --rebase` sequence is the documented fallback the verb automates, not the FO's first move.
```
Note: the `→ **shipped**: spacedock state commit <slug>` migration line (currently line 181) is UNCHANGED — it already names the verb correctly, which is why the routing oracle is already GREEN on this binding. A2 changes only the effect bullet so the *delegation* is the stated behavior, not the hand sequence.

### A3 — «merge.guard» intro + effect bullet (`fo-merge-core.md`)

BEFORE (intro, line 7):
```
When an entity reaches its terminal stage, `«merge.guard»(slug)` runs the atomic merge-finalize ceremony — the mod-block set→invoke→clear→terminalize sequence below, as one call.
```
AFTER (intro):
```
When an entity reaches its terminal stage, `«merge.guard»(slug)` drives the terminal merge-finalize ceremony as a re-entrant partial envelope: it arms the mod-block and signals the FO to invoke the merge hook, then on re-run detects completion by state delta, clears the mod-block standalone, terminalizes, and archives. It does NOT itself invoke the hook, and the clear+terminalize is two separate `--set` calls, not one.
```
BEFORE (effect bullet, line 11, first clause through "and archive"):
```
- **effect:** run the terminal merge-finalize ceremony for the entity — set `mod-block=merge:{mod_name}`, invoke the registered merge hook, detect completion, clear the mod-block in its own call, default-merge if no hook handled it, terminalize (`completed verdict={verdict} worktree=`), and archive.
```
AFTER (effect bullet, first clause through "and archive"):
```
- **effect:** drive the terminal merge-finalize ceremony for the entity as a re-entrant envelope — **arm:** set `mod-block=merge:{mod_name}` and signal the FO to invoke the registered merge hook (the verb does NOT invoke it); the FO invokes the hook, then re-runs the verb. **finalize (on re-run):** detect hook completion by the `pr`/`mod-block`/`verdict` state delta, clear the mod-block in its own standalone `--set` call, terminalize (`completed verdict={verdict} worktree=`), and archive. (The default local merge — when no hook handled it — and worktree-removal/worker-teardown, steps 9–10, are the FO's, not the verb's.)
```
The remaining clause of the effect bullet ("The mechanism-level enforcement in **Mod-Block Enforcement** below guards every step …") is UNCHANGED — it is accurate. Steps 1–10 and the `→ **shipped**` line (line 14) are UNCHANGED.

Three false claims removed, with their ground truth:
1. "invoke the registered merge hook" (the verb does the inviting) → `merge.go:13` "It does NOT invoke the merge hook"; `merge.go:213` arm signal "invoke the … merge hook, then re-run `merge guard`" — the FO invokes, the verb signals.
2. "as one call" (intro) → the verb is re-entrant (arm, then a separate finalize re-run) AND the clear is standalone from terminalize (`merge.go::finalize` emits two `emitSet` calls; step 5 line 28 "The clear MUST be standalone").
3. "default-merge if no hook handled it" attributed to the verb → `merge.go::MergeGuard` never local-merges; `finalize` only clears+terminalizes+archives. The default local merge is the FO's (step 6).

## Acceptance criteria

- **AC-1 (state.commit primary path).** An FO reading `### Split-Root State Sync` "Preferred" bullet and the «state.commit» effect bullet is directed to `spacedock state commit <slug>` as the primary path, with hand-git named explicitly as the fallback. **Verified by:** the A1+A2 rewrite landing the AFTER text verbatim (diff in `fo`-references on the impl worktree) AND `TestProseFunctionNotationBindsToRouting` GREEN on the unchanged `→ **shipped**: spacedock state commit` binding (`go test ./internal/cli -run TestProseFunctionNotationBindsToRouting`).
- **AC-2 («merge.guard» prose matches shipped behavior).** The «merge.guard» intro + effect bullet describe the verb as a re-entrant partial envelope (arm → FO invokes hook → re-run → finalize; the verb does NOT invoke the hook, does NOT local-merge, clear-then-terminalize is two `--set`s). **Verified by:** the A3 rewrite landing the AFTER text verbatim AND the validation reviewer reading `internal/status/merge.go` and confirming A3's new prose matches that code's actual behavior — the expected value is the CODE (an independent source that can diverge from the prose), not the prose's own tokens. The reviewer checks the three claims A3 asserts against their ground truth in the code: (a) the verb does NOT invoke the hook → `merge.go:13` doc-comment "It does NOT invoke the merge hook" and `merge.go::signalArmed` ("invoke the … merge hook, then re-run `merge guard`" — the FO invokes, the verb signals); (b) re-entrant arm-then-finalize-on-re-run with a standalone clear → `merge.go::finalize` emits two separate `emitSet` calls (clear, then terminalize), not one; (c) the verb never local-merges → `merge.go::MergeGuard`/`finalize` only clear+terminalize+archive, the default local merge being the FO's step 6. The code's own behavior is itself proven independently by `internal/status/merge_guard_test.go` (arm/finalize/clear-first/no-PR coverage). **Secondary (no routing regression):** `go test ./internal/cli -run TestProseFunctionNotationBindsToRouting` stays GREEN — the shipped routing oracle confirms A3's prose edits broke no «fn»↔verb binding. A substring presence/absence grep over `fo-merge-core.md` for the false/true strings is permitted ONLY as an explicitly-labeled weak supplement (a fast smell-test that the obvious false phrases are gone); it is NOT the proof, because a substring match over the instruction file the implementer just wrote is the prose-grep tautology the README Proof policy and the validation stage def explicitly ban — it cannot fail and a valid paraphrase would defeat it. This is the same antipattern F2 refused for the routing oracle.
- **AC-3 (routing oracle stays the non-vacuous guard, no new test).** The shipped `TestProseFunctionNotationBindsToRouting` + its RED control `TestProseFunctionRoutingGuardFailsOnViolation` + discriminator `TestProseFunctionRoutingOracleDiscriminates` remain GREEN after the A rewrites (the rewrites touch no `→` migration line). **Verified by:** `go test ./internal/cli -run 'TestProseFunction'` GREEN post-rewrite; the RED control's planted `merge teleport`/routed-guillemet violations are its non-vacuity proof (already authored). **No `deferred_tier_absence_test.go` is deleted** (it is not on main; see F3).

## Test plan

- **No new Go test is authored.** The routing oracle and its non-vacuity controls already exist and pass (#402). Authoring a duplicate violates the no-duplication rule. The slice's mechanical proof is: (a) the existing `internal/cli` oracle suite GREEN after the rewrite (no routing regression), and (b) for AC-2, the validation reviewer's prose-to-code review of A3 against `internal/status/merge.go` — the expected value is the code, an independent source that can diverge from the prose, not the prose's own tokens. A substring grep over `fo-merge-core.md` for AC-2's strings is an explicitly-labeled weak supplement (a smell-test), never the proof: a match over the instruction file the implementer wrote is the prose-grep tautology the README Proof policy bans.
- **Cost/complexity:** LOW. Three prose rewrites on the impl worktree + two `go test -run` invocations (routing-oracle baseline + post-rewrite, confirming no routing regression) + the validation reviewer's prose-to-code read of A3 against `internal/status/merge.go`. The optional supplement grep is seconds and not load-bearing. No new mechanism, no spike.
- **Fixture/CLI/live:** CLI-level only (`go test ./internal/cli -run TestProseFunction`). No fixture, no live workflow drive. The behavior of `spacedock merge guard` itself is already covered by `internal/status/merge_guard_test.go` and `internal/cli/merge_test.go` (shipped #400) — this slice does not re-test the verb; it makes the *prose* describe it truthfully.
- **Pre-impl ordering:** run `go test ./internal/cli -run TestProseFunctionNotationBindsToRouting` BEFORE the rewrite to confirm GREEN baseline, then after to confirm no regression. The riskiest thing the rewrite could break is a `→` migration line; the rewrites are scoped to effect/intro prose precisely to avoid that.

## Spike determination

No spike needed. The proven mechanisms this slice relies on: (1) `internal/cli/prose_function_routing_test.go` routing oracle — exercised GREEN on main today; (2) `internal/status/merge.go` shipped behavior (re-entrant arm/finalize, no-hook-invoke, standalone clear) — read directly from source and confirmed against the doc-comment and signal strings; (3) `internal/cli/state_sync.go::runStateCommit` — confirmed to implement the commit+sync+halt the A1/A2 prose will delegate to. All three are shipped code on main, not unverified mechanisms.

## Live-usage proof

The end-value proof that a live FO actually CALLS the verbs as its path rides `kt` (haiku-drive-validation, id `kt96jb8yagkean75j0k1ep6n`, sprint-readiness=DEFER, sequenced LAST as the sprint gate proof). This slice ships the prose rewire (A) and confirms the mechanical oracle (B, already shipped); the live drive is `kt`'s. `kt` is blocked on all four verb-core members + 72 + 6re merged AND m4's tmux backend green — this slice is one of the upstream rewires it waits on.

## Worktree / split-root note

A is a `skills/first-officer/references/` edit (scaffolding guardrail), so implementation runs WITH a worktree and the rewrites land on the impl branch — this is NOT an FO live-contract edit. This ideation entity + report live in the split-root state checkout (`docs/dev/.spacedock-state/`), committed path-scoped.

## Stage Report: ideation

- DONE: Concrete before/after wording for all three rewrites (Split-Root State Sync Preferred bullet; «state.commit» effect bullet; «merge.guard» body) recorded in the entity body.
  Recorded as A1/A2/A3 in the "Concrete rewrites" section with verbatim BEFORE/AFTER blocks; exact source lines (first-officer-shared-core.md:189, :178; fo-merge-core.md:7, :11) and ground-truth citations from merge.go (:13, :213).
- DONE: The routing oracle's assertion is designed plus its non-vacuous RED control (the pre-rewrite overstated «merge.guard» prose, or a non-existent-verb mutation) named.
  Designed as a reconciliation, not a new build: the routing oracle already shipped in #402 (`TestProseFunctionNotationBindsToRouting`); its non-vacuous RED control is the already-shipped `TestProseFunctionRoutingGuardFailsOnViolation` (plants `spacedock merge teleport`, a non-existent-verb mutation). F2 documents WHY "RED on overstated prose" is mechanically impossible for a routing oracle (it never reads effect prose) and routes that proof to the AC-2 token grep instead.
- DONE: Confirmed: B replaces deferred_tier_absence_test.go, and the rewrites are reconciled with czw's restructured contract regions (the byte-preservation tension resolved per-region).
  Confirmed the OPPOSITE of the seed and recorded it (F1/F3): `deferred_tier_absence_test.go` is NOT on main (only in unmerged `strip-deferred-tier-vocabulary`), and the routing oracle that would "replace" it already shipped independently — so B does NOT delete it; the "replaces" coupling is retired as a mis-mapping. Per-region reconciliation: A2 and A3 are scoped to effect/intro prose and leave every `→` migration line and the steps 1–10 UNCHANGED, so czw's routing-bearing regions are byte-preserved and the existing oracle stays GREEN.

### Summary

Read main directly and found the seed's premise was partly stale: czw's #402 already shipped the routing oracle (B) with its non-vacuity controls, `deferred_tier_absence_test.go` is not on main (so nothing to "replace"), and a routing oracle is structurally blind to the «merge.guard» prose defect it was supposed to catch. Resolved by reducing B to "confirm the shipped oracle stays GREEN", routing the prose-accuracy proof to an AC-2 token grep (the three false strings — "invoke the registered merge hook", "as one call", "default-merge if no hook handled it" — confirmed present today), and recording all three before/after rewrites verbatim. No spike needed; all relied-on mechanisms are shipped code exercised against source. This is a low-cost prose-rewire slice with no new Go test.

## Stage Report: ideation (cycle 2)

- DONE: Reframe AC-2's verification — the proof that A3's «merge.guard» prose matches the verb is NOT a substring/token grep over `fo-merge-core.md`.
  AC-2 (entity body) rewritten: the proof is now the validation reviewer reading `internal/status/merge.go` and confirming A3's new prose matches that code's actual behavior — the expected value is the CODE, an independent source that can diverge from the prose. The reviewer checks A3's three claims against ground truth in the code: (a) does-NOT-invoke-hook → `merge.go:13` + `signalArmed` line 213; (b) re-entrant clear-then-terminalize as two `--set`s → `finalize` lines 152/154/167 emit two separate `emitSet` calls; (c) never local-merges → `MergeGuard`/`finalize` only clear+terminalize+archive. The code's own behavior is independently proven by `internal/status/merge_guard_test.go`.
- DONE: Add the no-routing-regression leg and demote the substring grep to an explicitly-labeled weak supplement, never the proof.
  AC-2 now names `go test ./internal/cli -run TestProseFunctionNotationBindsToRouting` GREEN as the secondary no-routing-regression check, and the substring grep over `fo-merge-core.md` as a permitted-but-weak supplement only (a smell-test). AC-2 cites WHY the grep cannot be the proof: it is the prose-grep tautology banned by README Proof policy line 75 ("a match over the instruction file the implementer wrote… cannot fail; a valid paraphrase fails it") and the validation stage def's "Bad" clause — the same antipattern F2 refused for the oracle. Test-plan and Cost/complexity lines updated to match so the body is internally consistent.

### Summary

TIGHTENING cycle: changed ONLY AC-2's proof method, leaving the three approved rewrites A1/A2/A3 and AC-1/AC-3 untouched. The prior AC-2 justified a substring presence/absence grep over `fo-merge-core.md` as "acceptable because the AC is literally about specific prose tokens" — exactly the prose-grep tautology the README Proof policy (line 75) and the validation stage def explicitly ban, and the same antipattern F2 correctly refused for the routing oracle. Replaced it with a prose-to-code review against `internal/status/merge.go` (the independent CODE source that can diverge from the prose, itself proven by `merge_guard_test.go`) plus the shipped routing oracle staying GREEN; the substring grep survives only as an explicitly-labeled weak supplement. Verified all merge.go citations (lines 13, 152/154/167, 204/213) against source so the new proof is concrete and reproducible. Test plan and cost lines updated for internal consistency.

## Stage Report: implementation

- DONE: The three rewrites land VERBATIM per the entity's 'Concrete rewrites' before/after blocks, on the impl branch: A1, A2, A3.
  Impl commit ea2a9315 on branch spacedock-ensign/verbs-operative-rewire; `git diff --stat` = 2 files, 5 insertions/5 deletions — A1 (Split-Root State Sync Preferred/Fallback bullets name `spacedock state commit <slug>`, demote hand-git to fallback), A2 («state.commit» effect delegates to the verb), A3 (fo-merge-core.md «merge.guard» intro line 7 + effect bullet line 11 describe the re-entrant partial envelope). AFTER text matches the entity's blocks byte-for-byte.
- DONE: Confirm the three false strings ('as one call', 'invoke the registered merge hook', 'default-merge if no hook handled it') are gone.
  `grep -c` over fo-merge-core.md: "as one call" → 0, "default-merge if no hook handled it" → 0. "invoke the registered merge hook" survives ONLY inside the approved A3 AFTER text as the negated form ("signal the FO to invoke the registered merge hook (the verb does NOT invoke it)") — the false *attribution* (verb invokes the hook) is gone; the phrase appears only as its negation, exactly as the ideation-approved AFTER block specifies.
- DONE: No routing regression: TestProseFunctionNotationBindsToRouting GREEN as BEFORE baseline and AFTER the rewrite — both green.
  BEFORE: `go test ./internal/cli -run TestProseFunctionNotationBindsToRouting` → ok (0.246s) at HEAD d43479e0 pre-edit. AFTER: → ok (0.239s) post-edit. Full `TestProseFunction` suite (oracle + RED control `TestProseFunctionRoutingGuardFailsOnViolation` + discriminator `TestProseFunctionRoutingOracleDiscriminates`) all PASS.
- DONE: The `→ **shipped**` migration lines and the merge steps 1–10 are byte-preserved (czw's routing-bearing regions untouched); deliverable committed on the impl branch with a summary.
  `git diff` shows changes confined to the «state.commit» effect bullet, the Split-Root Preferred/Fallback pair, and fo-merge-core.md's intro + effect bullet; every `→ **shipped**` line and steps 1–10 are unchanged. The routing oracle (which parses those `→` lines) staying GREEN is the mechanical proof of byte-preservation. Committed as ea2a9315.

### Summary

Landed the three approved prose rewrites verbatim on the impl branch (commit ea2a9315): A1/A2 make `spacedock state commit <slug>` the operative/Preferred path with hand-git named as the explicit fallback, and A3 rewrites the «merge.guard» intro + effect to describe the shipped re-entrant partial envelope (verb arms + signals, FO invokes hook, verb re-run finalizes with a standalone clear, never local-merges). Verified A3's three claims against `internal/status/merge.go` source (doc-comment :13 "does NOT invoke the merge hook", signalArmed :213 arm signal, finalize's two separate emitSet calls :154/:167) before editing. Routing oracle GREEN before and after; the `→` migration lines and merge steps 1–10 are byte-preserved (diff is 5/5 lines, no routing regression). No new Go test authored per the entity's YAGNI reconciliation — the shipped #402 oracle + controls remain the non-vacuous guard.

## Stage Report: validation

- DONE: AC-2 (load-bearing) — INDEPENDENTLY read A3's «merge.guard» prose against `internal/status/merge.go` and confirmed all three claims match the code's ACTUAL behavior (expected value = the code, not the prose's tokens).
  (a) verb does NOT invoke the hook → `arm` (merge.go:139) does `emitSet`+`signalArmed`; `signalArmed` (:213) signals "invoke the … merge hook, then re-run `merge guard`" — FO invokes, verb signals; no hook-invocation call anywhere in MergeGuard/arm/finalize; matches doc-comment :13. (b) re-entrant arm→re-run→finalize with STANDALONE clear → `finalize` (:152) emits TWO separate `emitSet` calls: clear mod-block (:154), then terminalize status+verdict+completed (:167). (c) never local-merges → MergeGuard/finalize only clear+terminalize+`runArchive`; no `git merge` anywhere (line 122 is a comment about the arm phase, not a merge call). Code-oracle itself proven GREEN by `internal/status/merge_guard_test.go` (10 funcs incl. ArmThenFinalizeLocal, BlockedOnPR; `go test ./internal/status -run TestMergeGuard` → ok 1.307s). Prose matches code — not rejected.
- DONE: AC-1/AC-3 — `go test ./internal/cli -run TestProseFunction` GREEN (oracle + RED control + discriminator, no routing regression); A1/A2 direct the FO to `spacedock state commit <slug>` as PRIMARY/Preferred with hand-git named only as fallback; every `→ **shipped**` line and merge steps 1–10 byte-preserved.
  `TestProseFunction` (`-count=1`) → ok 0.263s: TestProseFunctionNotationBindsToRouting + TestProseFunctionRoutingGuardFailsOnViolation (RED control, plants `merge teleport`) + TestProseFunctionRoutingOracleDiscriminates all PASS. first-officer-shared-core.md:189 "Preferred — `spacedock state commit <slug>`", :190 "Fallback (no verb available / degraded)", :178 effect "delegate to `spacedock state commit <slug>`". Byte-preservation proven by the diff itself: `git diff ea2a9315^ ea2a9315` = 5/5 lines across the two .md files only; NO `→ **shipped**` line and NO merge step 1–10 appears in the diff. `go build ./...` clean.
- DONE: AC coverage cross-check — every `**AC-N**` in the entity body has non-self-referential evidence; AC-2 proven by the code-review path, not the demoted supplement grep.
  AC-1 evidence = A1/A2 diff + routing-oracle GREEN (command/state). AC-2 evidence = independent code-review of merge.go (the expected value is the code, an independent source that can diverge from the prose); the substring grep is explicitly demoted to a weak smell-test only — used as such (all 3 old false-attribution forms "as one call"/"default-merge if no hook handled it"/"ceremony for the entity — set" → count 0), NOT as proof. AC-3 evidence = `TestProseFunction` GREEN + F3 confirmed (`deferred_tier_absence_test.go` not tracked on this branch → nothing deleted; diff touches only the two contract .md files, no test authored/deleted). No AC rests on a prose-grep tautology.

### Summary

PASSED. The three approved prose rewrites (commit ea2a9315) landed verbatim and the slice's central proof — AC-2 — holds against the independent code oracle: I read `internal/status/merge.go` directly and confirmed A3's «merge.guard» prose matches all three behaviors (verb does NOT invoke the hook, re-entrant arm→re-run→finalize with a standalone two-`--set` clear, never local-merges), with the code itself proven GREEN by `merge_guard_test.go`. AC-1/AC-3 verified by command: `TestProseFunction` GREEN (oracle + non-vacuous RED control + discriminator), A1/A2 make `spacedock state commit <slug>` the Preferred/Primary path with hand-git as the named fallback, and the `→ **shipped**` migration lines plus merge steps 1–10 are byte-preserved (5/5-line diff, no routing regression). Every AC's evidence is code/command/state, not a substring match over the contract file; the supplement grep is used only as a labeled smell-test. No new Go test authored, consistent with the YAGNI reconciliation. Recommendation: PASSED.
