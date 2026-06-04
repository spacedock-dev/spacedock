---
id: p23sxe8ec3mmwgekvz9041a9
title: "spacedock pr complete — collapse the Merge-and-Cleanup ceremony into one binary command"
status: ideation
source: "captain (2026-06-04) — token-efficiency decomposition of first-officer-shared-core.md + the binary-simplification roadmap #3. Merge-and-Cleanup is MECHANICAL ceremony, so it routes to a binary command (not a lazy skill): collapsing it removes the prose outright. The FO ran this ceremony ~7x BY HAND this session (at/n3/2a/am/6b/ep/zd/p4 merges) — the heaviest, most error-prone sequence — the strongest single binary-command ROI."
score: "0.36"
worktree:
started: 2026-06-04T07:20:48Z
completed:
verdict:
issue:
---

The post-merge half of Merge-and-Cleanup (+ the Ship-Local Ceremony) is pure mechanical ceremony: mod-block clear (separate commit per the standalone rule) → terminalize (`--set completed verdict= worktree=`) → archive → worktree-remove → local-branch-delete → state push. It is the heaviest block in `first-officer-shared-core.md` (lines 207–242 today) and the most error-prone to run by hand — the FO ran it ~7× manually this session. Ceremony → BINARY command (the roadmap's lever 1), which removes the prose from the boot read entirely rather than deferring it to a skill. The team teardown is NOT part of this command — it is a Claude team-tool call (not shell-able), so it stays FO/runtime-side per the using-claude-team boundary.

## Command-shape options (the captain wants OPTIONS, not the roadmap name taken as given)

The roadmap names this `spacedock pr complete {slug}`. The captain is explicitly unsure about the shape. Below are the live axes, each weighed, with a recommendation.

### Axis A — naming / namespace

| Option | For | Against |
|---|---|---|
| **A1. `spacedock pr complete {slug}`** (roadmap name) | Matches the roadmap and the `pr-merge` mod vocabulary; "complete" reads as "finish the PR's lifecycle." | The ceremony also runs on the `merge: local` / ship-local path where there is NO PR. A `pr` namespace actively MISLEADS there — a captain on a local-merge workflow would not think to reach for `pr complete`. The post-merge ceremony is about ENTITY lifecycle, not PR lifecycle. |
| **A2. `spacedock complete {slug}`** | Short, host-frontdoor-adjacent. | A bare top-level `complete` is ambiguous against the launch verbs (`claude`/`codex`); reads like "complete the session." No namespace to group siblings under. |
| **A3. `spacedock entity complete {slug}`** (or `lifecycle`) | An `entity`/`lifecycle` namespace is the honest home: this acts on the ENTITY's terminal transition, PR or local. Gives siblings #1/#2 a shared roof (`entity advance`, `entity complete`). | No `entity` namespace exists yet; introduces a new top-level noun. `status`/`dispatch`/`state` are the current nouns — `entity` overlaps `status` conceptually. |
| **A4. Fold into `spacedock dispatch advance {slug} --to done`** (sibling #2) | One command for "move this entity to the next stage," terminal stages included; no new surface. | Conflates two very different operations: a mid-flow stage advance (frontmatter + state commit) vs. the terminal ceremony (mod-block guard dance + archive + worktree/branch teardown + the merge-already-happened precondition). Overloading `advance --to {terminal}` hides a 6-step destructive ceremony behind a generic verb. Couples #3 to #2's design and serializes their delivery. |

**Recommended: A3 — `spacedock dispatch complete {slug}`** (reuse the existing `dispatch` namespace rather than mint a new `entity` noun). Rationale: the ceremony is entity-lifecycle, not PR-bound, so a `pr` namespace (A1) is wrong on the ship-local path — the very path the captain flagged. `dispatch` already houses lifecycle-adjacent verbs (`dispatch build`, `dispatch show-stage-def`) and is the natural sibling home for #2 `dispatch advance`; `dispatch complete` reads as "complete this entity's dispatch lifecycle" and is true for both PR and local merges. It avoids a brand-new top-level noun (A3-`entity`) while keeping the grouping benefit. A4 is rejected: terminal ceremony is structurally different from a stage advance and must not be hidden behind a generic verb. A1's name survives only as a possible alias if the captain wants the roadmap spelling discoverable.

*(If the captain prefers a dedicated noun for the whole lifecycle family, A3-`entity complete` is the fallback — same behavior, different roof.)*

### Axis B — scope (what the command owns)

| Option | Owns | Notes |
|---|---|---|
| **B1. Post-merge ceremony ONLY** | mod-block clear (own commit) → terminalize → archive → worktree-remove → branch-delete → state push. | The mechanical, deterministic, idempotent core. The merge MUST already have landed (PR merged, or local `--no-ff` done + `pr=local-merge:{sha}` sentinel set). |
| **B2. B1 + detect-the-merge** | Also runs `gh pr view` / inspects the merge state before ceremony. | Detection is already the `pr-merge` mod's startup/idle hook job and the FO event-loop PR scan. Duplicating it here splits ownership and drags `gh` availability into a state-mutation command. Keep detection in the mod; this command's precondition is "merge already recorded." |
| **B3. B1 + team teardown** | Also tears down the agent cohort + team. | The team-tool teardown is a Claude tool, NOT shell-able. A Go binary cannot issue it. Hard out of scope — stays FO/runtime-side. |

**Recommended: B1.** The command owns the deterministic post-merge ceremony and nothing else. Its precondition is an already-recorded merge (PR `MERGED`, or the `pr=local-merge:{sha}` sentinel present). Detection (B2) stays in the `pr-merge` mod / FO scan; team teardown (B3) is non-shell-able and stays FO-side. This keeps the command pure-mechanical and host-tool-free, which is exactly what makes it idempotent and testable.

### Axis C — relationship to siblings #1 `state sync` / #2 `dispatch advance`

Neither sibling is shipped yet (confirmed: `state_sync_test.go` tests the git DISCIPLINE the FO runs by hand; no `runStateSync` command exists). Options: (C1) three independent commands under whatever namespace each lands in; (C2) a unified `entity`/lifecycle namespace housing advance + complete, with `state sync` separate (it is universal, not lifecycle). **Recommended: C2-lite** — put `complete` and (future) `advance` under the existing `dispatch` namespace; leave `state sync` under `state` (it is a universal git-sync op, not an entity-lifecycle op). This is consistent with the roadmap's own note that `state sync` is universal while `pr complete` is dev-specific.

### Dev-specificity flag

This command is **dev-workflow-specific**, not universal-binary behavior. The README top-level `merge:` key (default `pr`; `local` for host-less workflows) is the established precedent for opt-in workflow-bound behavior. The command must read the workflow's `merge:` policy and the entity's recorded merge state rather than assume a PR — that is precisely why the `pr`-namespace name (A1) is wrong.

## Recommended approach (firmed)

`spacedock dispatch complete {slug}` orchestrates the existing, already-proven binary operations transactionally, ordered so a partial failure leaves a safe, re-runnable state:

1. **Precondition check** — resolve `{slug}`; if already archived (file in `_archive/`), exit 0 as a no-op (idempotency). Else confirm the merge is recorded: `pr` is set (PR or `local-merge:{sha}` sentinel) OR the workflow is `merge: local`. If not, exit non-zero with "merge not recorded — run the merge hook first" (NO `--force` on the happy path).
2. **Clear mod-block in its own commit** — `status --set {slug} mod-block=` then a path-scoped state commit. Standalone, never combined with terminal fields (the existing `status --set` guard enforces this; the command does not need `--force`).
3. **Terminalize** — `status --set {slug} completed verdict={verdict} worktree=` (default `verdict=PASSED`; `--verdict` overrides), path-scoped commit.
4. **Archive** — `status --archive {slug}`, path-scoped commit.
5. **State push** — push the state branch; on non-fast-forward, `pull --rebase` then re-push; on rebase CONFLICT, HALT and surface (the existing state-sync discipline — this is where future #1 `state sync` would be reused as a helper).
6. **Worktree + branch teardown (main repo)** — `git worktree remove {path}` (no `--force`; the untracked-files refusal is the safety net) and `git branch -d {branch}`. The `worktree:` path comes from frontmatter read BEFORE step 3 clears it; the branch is `{worker_key}/{slug}`.

The FO prose (shared-core lines 207–242) shrinks to roughly: "on merged-PR detected (or after the local merge hook lands), run `spacedock dispatch complete {slug}`, then run the terminal team teardown (FO/runtime-side)." The guard discipline, separate-commit rule, archive, worktree/branch teardown, and state-push/rebase prose all move into the binary.

**Constraints:**
- **Idempotency** — re-running on an already-archived entity is a clean exit-0 no-op (step 1). Re-running after a partial failure converges (each step is itself idempotent: clearing an already-clear mod-block, terminalizing an already-terminal entity, archiving an already-archived entity, removing an already-removed worktree all no-op or are detected).
- **Guard-honoring** — the mod-block clear is its own commit (the `status --set` guard already enforces the no-combine rule); NO `--force` anywhere on the happy path. If a guard refuses, a precondition was unmet — the command surfaces it, it does not paper over with `--force`.
- HIGH-STAKES (status mutation/guard + CI/release machinery) → detached audit at validation.

### Riskiest unknown — no spike needed, with the one integration risk flagged

Every individual step is already-proven, composable binary behavior the FO runs by hand today: `status --set` (clear/terminalize, with its standalone-mod-block guard), `status --archive` (with its verdict + merge-hook + mod-block guards), `git worktree remove`, `git branch -d`, and the state push / pull-rebase / conflict-halt discipline (proven in `state_sync_test.go`). The design only COMPOSES proven mechanisms — **no spike needed for the individual operations.**

The ONE genuinely-new integration risk is **transactional ordering + partial-failure recovery across the split-root boundary**: steps 2–5 mutate the STATE checkout (commits + push) while step 6 mutates the MAIN repo (worktree/branch). A failure between step 5 and step 6 leaves the entity archived-and-pushed but the worktree/branch still present. The chosen ordering makes this safe — re-running detects "already archived" at step 1 and proceeds to a still-needed step 6 teardown (so the no-op path must still attempt teardown of any lingering worktree/branch, not bail at step 1). **This ordering-and-recovery claim is what the validation end-to-end fixture must exercise first** (the smallest run that would invalidate the rest), and is the detached-audit focus given the high-stakes surface. Recording it here so the determination is on the record.

## Acceptance criteria

- **AC-1:** `spacedock dispatch complete {slug}` takes a merged-state entity (PR merged or `pr=local-merge:{sha}` recorded, or `merge: local`) to fully archived — mod-block cleared in its OWN commit, terminalized with `verdict=PASSED` (or `--verdict`), worktree removed, local branch deleted, state branch pushed — with NO `--force` anywhere.
  *Verified by:* an end-to-end fixture (real git, per the `state_sync_test.go` / `reconcile_e_test.go` pattern) driving the binary from merged-state to archived and asserting on-disk: entity file in `_archive/`, `verdict: PASSED`, `worktree:` empty, two distinct commits for mod-block-clear vs terminalize, worktree dir gone, branch absent, state branch advanced. The fixture must pass through every guard (a `--force` in the happy path fails the test).
- **AC-2:** Idempotent and partial-failure-safe — re-running on an already-archived entity exits 0; the second run still tears down a lingering worktree/branch if present (no-op step 1 does not skip step 6).
  *Verified by:* a fixture that runs the command twice (second run asserts exit 0 + no error), and a fixture that simulates a crash after archive-push but before worktree teardown, then re-runs and asserts the worktree/branch are cleaned and exit is 0.
- **AC-3:** The merge precondition is enforced — running `complete` on an entity with no recorded merge (no `pr`, not `merge: local`) exits non-zero with a diagnostic, without mutating state.
  *Verified by:* a fixture asserting non-zero exit + unchanged entity frontmatter (still non-terminal, file still active).
- **AC-4:** The FO contract prose for the post-merge ceremony (`first-officer-shared-core.md` Merge-and-Cleanup steps 5/7/8/9 + the Ship-Local Ceremony steps 3–7) is replaced by the single command invocation plus the FO-side team-teardown step, measurably shrinking those blocks.
  *Verified by:* the post-implementation `first-officer-shared-core.md` Merge-and-Cleanup + Ship-Local Ceremony region is measurably shorter (state the measured before/after line + char count at implementation; target ≥30% reduction of those two blocks), AND a presence check that the surviving prose names `spacedock dispatch complete` and retains the FO-side team-teardown step.

## Test plan

- **Primary proof — Go end-to-end fixtures (real git).** Mirror `internal/cli/state_sync_test.go` and `internal/dispatch/reconcile_e_test.go`: a tmp repo + bare origin + state branch, seeded with a merged-state entity. Drive the compiled command path. Estimated complexity: MEDIUM — the git-fixture scaffolding (`cloneOnStateBranch`, `commitEntity`, `seedStateBranch`) already exists and is the template; the new code is the orchestration + precondition/idempotency logic (~200 LOC per roadmap est).
- **Guard-composition unit coverage.** Assert the command issues the mod-block clear as a SEPARATE `status --set` (the guard rejects a combined call) and never passes `--force` — drive the real `status` runner, not a mock, so the existing guards do the rejecting.
- **Prose-shrink check (AC-4).** A presence/measurement check over `first-officer-shared-core.md`: the Merge-and-Cleanup + Ship-Local region names the new command and is shorter by the measured target. This is a property-of-the-text check (legitimate per the ideation stage def — the claim IS about the text), not a substring stand-in for behavior.
- **Detached audit at validation.** HIGH-STAKES surface (status mutation/guard + the merge/teardown machinery). The validation stage gets a detached audit whose lens is the transactional-ordering / partial-failure-recovery claim (the riskiest unknown above) and the no-`--force` guarantee.
- **Live workflow smoke (optional, stretch).** A single live run of `dispatch complete` on a real merged entity in the dev workflow, only if fixture coverage leaves runtime doubt.

## Implementation sequencing (FO must serialize the shared-core edit)

This entity and its two siblings — `gate-presentation-skill-extraction` and `feedback-rejection-flow-skill-extraction` — all edit the SAME `skills/first-officer/references/first-officer-shared-core.md` at implementation. Their implementations MUST be serialized (the FO coordinates). Sections THIS entity touches:
- **`## Merge and Cleanup`** (lines ~207–226) — steps 5 (mod-block clear), 7 (terminalize), 8 (archive), 9 (worktree/branch removal) collapse into the command invocation; step 10 (team teardown) stays as FO-side prose.
- **`### Ship-Local Ceremony`** (lines ~228–242) — steps 3–7 collapse likewise; the `merge: local` policy read stays as the precondition the command honors.
- Untouched by this entity: `## Gate Presentation` (sibling #1's target) and `## Feedback Rejection Flow` (sibling #2's target) — so the conflict is only over the file, not the regions, and a clean serialization avoids merge conflicts entirely.

## Out of scope

- The team-tool teardown (Claude tool, not shell-able — stays in the using-claude-team skill / runtime; survives as FO-side prose in Merge-and-Cleanup step 10).
- Merge DETECTION (B2) — stays in the `pr-merge` mod's startup/idle hooks + the FO event-loop PR scan.
- #1 `state sync` and #2 `dispatch advance` themselves — sibling roadmap binary commands; file separately if/when Phase 1 is greenlit. This entity REUSES the state push/pull-rebase discipline #1 would wrap, but does not depend on #1 landing first.

## Notes

Roadmap #3 (binary-simplification-roadmap.md, refreshed 2026-06-04 — named the strongest single binary ROI; the FO ran it ~7× manually this session). The qs reconcile helper + `state_sync_test.go` are the test-coverage templates (real-git fixtures). Sibling lever-2 extractions: `gate-presentation-skill-extraction`, `feedback-rejection-flow-skill-extraction`. Command-shape recommendation departs from the roadmap's `pr complete` name to `dispatch complete` because the ceremony is entity-lifecycle, not PR-bound (the ship-local path has no PR).

## Stage Report: ideation

- DONE: The captain is NOT sure about the `pr` subcommand shape — enumerate and weigh a FEW concrete options before recommending.
  Added `## Command-shape options` weighing naming (A1–A4), scope (B1–B3), sibling relationship (C1/C2); recommended `dispatch complete` over the roadmap `pr complete` (entity-lifecycle, not PR-bound — the ship-local path has no PR), with dev-specificity flagged against the README `merge:` precedent.
- DONE: Firm the riskiest-unknown + ACs for the recommended option (idempotency, guard-honoring, no `--force`, end-to-end fixture through every guard) + measured prose-shrink target.
  Recorded "no spike needed: composes proven `status --set`/`--archive` guards + `git worktree remove`/`branch -d` + state push/pull-rebase (state_sync_test.go)"; flagged the ONE new integration risk (transactional ordering / partial-failure recovery across the split-root↔main boundary) as the validation fixture's first exercise + detached-audit lens. AC-1..AC-4 each carry a "Verified by" naming a real-git fixture or a property-of-text check; AC-4 sets a ≥30% prose-shrink target on the Merge-and-Cleanup + Ship-Local blocks.

### Summary

Firmed the design for `spacedock dispatch complete {slug}` (departing from the roadmap's `pr complete` name — the ceremony is entity-lifecycle, and the ship-local `merge: local` path has no PR, so a `pr` namespace would mislead). Scope is the deterministic post-merge ceremony ONLY (mod-block clear → terminalize → archive → worktree/branch teardown → state push); merge-detection stays in the pr-merge mod and team teardown stays FO-side (non-shell-able Claude tool). No spike needed — the design composes already-proven binary operations; the sole new integration risk (transactional ordering + partial-failure recovery across the split-root↔main boundary) is named as the validation fixture's first exercise and the detached-audit focus. Flagged the shared-core serialization: this entity edits Merge-and-Cleanup + Ship-Local Ceremony, disjoint from the two siblings' regions (Gate Presentation, Feedback Rejection Flow), so the FO can serialize the file edits conflict-free.
