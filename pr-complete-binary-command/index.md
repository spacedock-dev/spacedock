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

## Architecture boundary — command vs. mod vs. FO prose (cycle-2 resolution)

The captain's question: *"is this cmd part of the mod? or part of the cmd is part of the mod?"* The post-merge ceremony is **TRIPLE-HOMED** today — the same logic is spelled out in three places:

1. **FO shared-core `## Merge and Cleanup` + `### Ship-Local Ceremony` prose** (lines 207–242) — the FO follows it by hand.
2. **`_mods/pr-merge.md` startup + idle hooks** — these ALREADY re-list the full ceremony. Quoting the startup hook (pr-merge.md line 15): on `MERGED`, "advance the entity to its terminal stage: set `status` to the terminal stage, `completed` to ISO 8601 now, `verdict: PASSED`, clear `worktree`, archive the file, and clean up any worktree/branch." That is the ceremony, prose-duplicated a third time. The idle hook (line 25) says "same logic as the startup hook."
3. **The proposed binary command** — a fourth would-be home.

**Resolution — the command is the single ceremony IMPLEMENTATION; the mod and the FO/ship-local path are its two CALLERS.** Mapped by responsibility:

| Layer | Owns | After this entity |
|---|---|---|
| **COMMAND** (`spacedock dispatch complete {slug}`) | The deterministic mechanical ceremony — ONE implementation: precondition → mod-block clear (own commit) → terminalize → archive → state push → worktree/branch teardown. Idempotent, guard-honoring, no `--force`. | The single source of truth for the ceremony. |
| **MOD** (`_mods/pr-merge.md`) | PR-lifecycle **POLICY**: build/open the PR, then DETECT the merge (`gh pr view`). Its startup/idle hooks' re-listed ceremony steps **collapse to one line: invoke `spacedock dispatch complete {slug}`.** The mod decides WHEN (merge detected); the command does WHAT (the ceremony). | The PR-path **caller** of the command. |
| **FO prose** (shared-core) | Minimal: "at the terminal boundary, the mod (PR path) or the FO directly (ship-local path) runs `spacedock dispatch complete {slug}`; then run the team teardown (FO/runtime-side, non-shell-able)." | Points at the command; team teardown is the only ceremony step that stays prose (can't shell a Claude tool). |

So the captain's *"part of the cmd is part of the mod"* resolves to: **the mod CALLS the command.** The command is shared infrastructure; the `pr-merge` mod is the PR-path caller; the FO (ship-local `merge: local`) is the other caller. Because the command serves BOTH callers — PR-via-mod and direct-local-merge — a `pr` namespace (A1) is doubly wrong: it is not even mod-specific, let alone PR-specific. This confirms and strengthens the cycle-1 naming recommendation.

**Who invokes the command:**
- **PR path:** the `pr-merge` mod's startup/idle hook invokes it on detecting `MERGED` (replacing its re-listed steps). The FO, running the mod hook, is the proximate runner, but the INVOCATION text lives in the mod, not FO prose.
- **Ship-local path** (`merge: local`, or pr fallback): the FO invokes it directly after the local merge hook lands, because there is no PR-detection mod step to host the call.

This is a one-way migration (the roadmap's stated rule): once the command exists, neither the mod nor the FO prose re-states the ceremony — they name the command. The cycle-2 work removes the duplication, not just the FO-prose copy.

**Touched at implementation (cycle-2 addition to the file map):** `_mods/pr-merge.md` startup hook (the `MERGED` paragraph, ~line 15) and idle hook (~line 25) collapse to a command invocation. This is a NEW file in the change set beyond `first-officer-shared-core.md` — but it is a mod file, written by a dispatched worker (the FO does not edit mods directly per FO Write Scope), and it does NOT collide with the two siblings (which touch only shared-core), so it does not add to the serialization constraint.

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
- **AC-4:** The post-merge ceremony's triple-homed prose collapses to the command invocation in BOTH non-implementation homes: (a) `first-officer-shared-core.md` Merge-and-Cleanup steps 5/7/8/9 + Ship-Local Ceremony steps 3–7, and (b) the `_mods/pr-merge.md` startup-hook `MERGED` paragraph + idle hook. The surviving prose names the command and keeps the FO-side team-teardown step.
  *Verified by:* the post-implementation `first-officer-shared-core.md` Merge-and-Cleanup + Ship-Local Ceremony region is measurably shorter (state the measured before/after line + char count at implementation; target ≥30% reduction of those two blocks); the `pr-merge` mod's `MERGED` ceremony steps are replaced by the command invocation; AND a presence check that both surviving texts name `spacedock dispatch complete` and shared-core retains the FO-side team-teardown step.
- **AC-5 (live safety-net — HARD requirement, captain: "we gotta have test to make that change"):** The change to the post-merge ceremony is covered by the `pr-lifecycle-from-boot` scenario, codified as a host-neutral scenario AND run live: boot a workflow → observe PR-pending + dispatchable startup state → let the PR lifecycle advance (the mod detects the merge and invokes the command) → the entity reaches the correct DURABLE terminal state (archived, `verdict: PASSED`, `worktree:` empty, terminal status) WITHOUT bypassing merge hooks or archival rules. This promotes `pr-lifecycle-from-boot` from prose-only (a prioritized-but-unbuilt scenario in `docs/specs/scenario-testing-principles.md`) into a real codified + live scenario.
  *Verified by:* (i) a new `pr-lifecycle-from-boot` entry in the `sharedRuntimeScenarios()` table (`internal/ensigncycle/shared_scenarios_test.go`) with a per-host runner each side, so the existing shared-coverage meta-tests red if either host lacks it; (ii) the seed-scenario lock test binding `docs/specs/scenario-testing-principles.md` to the table updated to include it (doc↔code lock stays green); (iii) an authored `livescenario.Scenario` ({Runbook, Setup, Assert}) whose `Assert` grades the durable BEFORE→AFTER entity state (archived + terminal + verdict + empty worktree), NOT transcript phrasing, runnable via the codified executor in CI and the LLM executor live. The live LLM-executor run is the producer proof the citation gate requires for the runtime-observable claim "the FO+mod actually complete the lifecycle from boot."

## Test plan

- **Primary proof — Go end-to-end fixtures (real git).** Mirror `internal/cli/state_sync_test.go` and `internal/dispatch/reconcile_e_test.go`: a tmp repo + bare origin + state branch, seeded with a merged-state entity. Drive the compiled command path. Estimated complexity: MEDIUM — the git-fixture scaffolding (`cloneOnStateBranch`, `commitEntity`, `seedStateBranch`) already exists and is the template; the new code is the orchestration + precondition/idempotency logic (~200 LOC per roadmap est).
- **Guard-composition unit coverage.** Assert the command issues the mod-block clear as a SEPARATE `status --set` (the guard rejects a combined call) and never passes `--force` — drive the real `status` runner, not a mock, so the existing guards do the rejecting.
- **Prose-shrink check (AC-4).** A presence/measurement check over `first-officer-shared-core.md`: the Merge-and-Cleanup + Ship-Local region names the new command and is shorter by the measured target. This is a property-of-the-text check (legitimate per the ideation stage def — the claim IS about the text), not a substring stand-in for behavior.
- **`pr-lifecycle-from-boot` scenario — codified + live (AC-5, MANDATORY).** Author the scenario once via the `internal/livescenario` primitive ({Runbook, Setup, Assert}); register it in `sharedRuntimeScenarios()` with a per-host runner each side. Run two executors of the one scenario:
  - *Codified executor* (every CI, cheap): the deterministic Go path proves the modeled consumer logic — boot state observed, ceremony runs, durable terminal state reached.
  - *LLM executor* (live, on the producer claim): a real Claude/Codex agent boots the workflow and drives the lifecycle; graded on the durable BEFORE→AFTER entity outcome, never transcript wording. This is the producer proof the citation gate requires.
  Estimated complexity: MEDIUM-HIGH — the primitive + the shared-table + lock-test machinery already exist (p4 landed them); the new work is the scenario's Setup (stage a PR-pending merged entity from boot) and Assert (durable terminal-state grade), plus the two host runners. The `Setup` must stage the SAME merged-PR-from-boot state the AC-1 fixture uses, so the two share a fixture builder.
- **Detached audit at validation.** HIGH-STAKES surface (status mutation/guard + CI machinery + the mod contract). The validation stage gets a detached audit whose lens is the transactional-ordering / partial-failure-recovery claim (the riskiest unknown above), the no-`--force` guarantee, and the command↔mod boundary (that the mod truly delegates rather than re-implementing). Staff review may precede the gate given the surface.

## Implementation sequencing (FO must serialize the shared-core edit)

This entity and its two siblings — `gate-presentation-skill-extraction` and `feedback-rejection-flow-skill-extraction` — all edit the SAME `skills/first-officer/references/first-officer-shared-core.md` at implementation. Their implementations MUST be serialized (the FO coordinates). Sections THIS entity touches:
- **`## Merge and Cleanup`** (lines ~207–226) — steps 5 (mod-block clear), 7 (terminalize), 8 (archive), 9 (worktree/branch removal) collapse into the command invocation; step 10 (team teardown) stays as FO-side prose.
- **`### Ship-Local Ceremony`** (lines ~228–242) — steps 3–7 collapse likewise; the `merge: local` policy read stays as the precondition the command honors.
- Untouched by this entity: `## Gate Presentation` (sibling #1's target) and `## Feedback Rejection Flow` (sibling #2's target) — so the conflict is only over the file, not the regions, and a clean serialization avoids merge conflicts entirely.

Additional files THIS entity touches (no sibling collision):
- **`_mods/pr-merge.md`** — startup-hook `MERGED` paragraph + idle hook collapse to the command invocation (the boundary resolution above). A mod file, written by a dispatched worker, not the FO.
- **`internal/ensigncycle/shared_scenarios_test.go`** + **`docs/specs/scenario-testing-principles.md`** — register `pr-lifecycle-from-boot` in the shared table and update the doc↔code seed lock (AC-5).
- **`internal/livescenario`** consumer + per-host runners — the authored scenario (no change to the primitive itself; p4 already shipped it).

## Out of scope

- The team-tool teardown (Claude tool, not shell-able — stays in the using-claude-team skill / runtime; survives as FO-side prose in Merge-and-Cleanup step 10).
- Merge DETECTION (B2) — stays in the `pr-merge` mod's startup/idle hooks + the FO event-loop PR scan. The cycle-2 boundary resolution keeps detection in the mod but has the mod's detection step now CALL the command instead of re-listing the ceremony. The command's own precondition (merge already recorded) is the guard, not a second detector.
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

## Stage Report: ideation (cycle 2)

- DONE: Resolve the cmd↔mod boundary ("is this cmd part of the mod? or part of the cmd is part of the mod?").
  Added `## Architecture boundary` — confirmed by reading `_mods/pr-merge.md`: the ceremony is TRIPLE-HOMED (FO prose + the mod's startup/idle hooks, which already re-list set-terminal→verdict→clear-worktree→archive→clean-worktree/branch + the proposed command). Resolution: the COMMAND is the single ceremony implementation; the MOD detects the merge then CALLS the command (its re-listed steps collapse to one invocation line); FO prose shrinks to "the mod (PR path) or FO (ship-local) runs the command, then team teardown stays FO-side." So "part of the cmd is part of the mod" = the mod calls the command. Serving BOTH callers (PR-via-mod + direct local-merge) re-confirms a `pr` namespace is wrong; `dispatch complete` stands.
- DONE: Fold in `pr-lifecycle-from-boot` as a HARD test requirement ("we gotta have test to make that change").
  Added AC-5: codify `pr-lifecycle-from-boot` in `sharedRuntimeScenarios()` + the doc↔code seed lock, authored via p4's `internal/livescenario` primitive ({Runbook, Setup, Assert}), graded on durable terminal state (archived/verdict/empty-worktree), run codified in CI AND live (LLM executor = the citation-gate producer proof). Promotes the scenario from prose-only (a prioritized-but-unbuilt entry in scenario-testing-principles.md) into a real codified+live scenario. Test plan + file map updated (pr-merge.md, shared_scenarios_test.go, scenario-testing-principles.md added as touched, no sibling collision).
- DONE: Keep the sound cycle-1 parts (scope B1, offline guard ACs, prose-shrink target, transactional/partial-failure integration risk).
  Untouched: AC-1/2/3, the no-spike determination, the riskiest-unknown (split-root↔main transactional ordering). AC-4 extended to also require the mod's ceremony prose collapse (not just FO prose). Detached-audit lens widened to include the command↔mod boundary; staff review may precede the gate (HIGH-STAKES: status guard + CI machinery + mod contract).

### Summary

Cycle 2 resolves the captain's boundary reframing: the command is the ONE ceremony implementation, the `pr-merge` mod becomes a thin caller (its startup/idle hooks already duplicate the ceremony — they collapse to a single `spacedock dispatch complete {slug}` invocation), and the FO contract shrinks to naming the command; team teardown alone stays FO-side prose (non-shell-able). This kills a TRIPLE duplication, not just the FO copy, and re-confirms the non-`pr` name (the command serves both the PR-via-mod and direct local-merge callers). Added AC-5 making `pr-lifecycle-from-boot` a hard codified+live scenario via p4's livescenario primitive, promoting it from a prose-only roadmap entry into a real safety net graded on durable terminal state. The HIGH-STAKES surface now spans the status guard, CI/scenario machinery, AND the mod contract — detached audit at validation, staff review may precede the gate.
