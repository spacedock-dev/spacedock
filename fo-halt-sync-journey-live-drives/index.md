---
id: ev3e0nmknh98sn365ky163en
title: Behavioral live drives for the FO halt/sync/journey contract behaviors (the Bucket-A cluster with no existing drive)
status: ideation
source: "captain (2026-06-05) — split from hwk (tautological-test-remediation) per the (b) phasing call. hwk demotes the halt/sync/journey presence checks to honest non-AC lints; this task supplies their real behavioral proof — a new cross-host live scenario (split-root-halt) plus the sync/journey behaviors that today have NO live drive."
score: "0.33"
started: 2026-06-05T04:40:09Z
completed:
verdict:
worktree:
issue:
---

The tautological-test remediation (hwk) found that most Bucket-A behavioral claims already have live drives (present-gate → gate-guardrail scenario; feedback-rejection-flow → rejection-flow scenario; using-claude-team → the team scenarios), so those demote-and-bind. But four Bucket-A tests assert FO behaviors that have **no live drive at all** — they are pure presence checks (banned tautologies) and hwk demotes them to honest non-AC lints, leaving the behaviors unproven. This task supplies the owed behavioral proof.

## Problem — the behaviors with no drive

- **`TestFOHaltGateProse`** (`internal/hostneutrality/split_root_sync_contract_test.go`) — the FO must HALT on an uninitialized split-root state dir (orphan branch on origin, no linked worktree) and point at `spacedock state init`, NOT dispatch against the EMPTY-rendering boot table. SAFETY behavior: a silent failure if it doesn't halt.
- **`TestFOSyncProse` / `TestEnsignSyncProse`** (same file) — the state-checkout sync rules (path-scoped commits, `pull --rebase` on push-rejection, `rebase --abort` + no `--force` on conflict).
- **`TestCommissionJourneyProse`** (same file) — the commission journey performs the orphan/worktree/`state init` sequence.

Today all four are `strings.Contains` presence checks over the contract prose (`skills/first-officer/references/first-officer-shared-core.md`, `skills/ensign/references/ensign-shared-core.md`, `skills/commission/SKILL.md`) — they prove the words are present, never the behavior. hwk honestly demotes them (Oracle-honesty non-AC lint naming THIS task as the owed oracle).

## Riskiest unknown — SPIKE RUN (result recorded)

**Question:** grading a HALT is grading the ABSENCE of action — can a durable-state assertion cleanly distinguish "FO correctly halted" from "FO failed to halt and dispatched against the empty table"?

**Spike (throwaway, offline + real-binary, this session):** staged an uninitialized split-root via the `state_init_test.go` `commissionSplitWorkflow` mechanic (orphan state branch pushed to a bare origin, then a fresh clone with the state path ABSENT), ran the real binary's `status --boot` and `--validate`, and prototyped the halt-grading assertion across all four end-states.

**Findings (empirical, on the built `cmd/spacedock`, contract 1):**
1. **The silent-failure trap is real.** On the uninitialized split-root, `status --boot` renders an EMPTY `DISPATCHABLE` table and `status --validate` prints `VALID` (exit 0). The ONLY signal of trouble is `STATE_BACKEND: split-root (entity_dir: …, present: false)`. A naive FO reads "nothing to do, all valid" and goes idle — the exact silent failure the halt-gate exists to catch.
2. **Pure entity before/after CANNOT grade this.** In the uninitialized case there is no entity file to read (the state checkout is absent), and BOTH a correct halt and a silent "no work" failure leave the same on-disk state (nothing dispatched, nothing mutated). So the existing `livescenario.Run`'s entity-body before/after triple is insufficient on its own — the assertion needs additional durable signals + the observed output.
3. **Observed-substring grading alone is WEAK (false-pass).** An FO that mentions "state init" in passing but dispatches anyway PASSES an observed-only `Contains("state init")` check. The spike reproduced this false pass.
4. **Two DURABLE signals make the grade sharp.** Empirically: a correct halt leaves code-branch porcelain CLEAN and creates NO `.worktrees/` dispatch dir; the correct `state init` remedy leaves porcelain CLEAN and creates only the (gitignored) state path, with the seeded entity untouched at its initial stage. A dispatch-against-empty is the ONLY path that creates a `.worktrees/` dir or dirties code porcelain or mutates the seeded entity. The strengthened assertion (no dispatch worktree AND clean porcelain AND [`state init` ran OR observed named the uninitialized state]) PASSES the two correct remedies and REDS all three broken modes (silent-no-work, dispatched-but-mentioned, dispatched-with-words). The 5-case offline prototype was green.

**Conclusion:** the halt is gradeable. The grade is a **two-track triple**: (durable) no dispatch worktree created + code porcelain clean + seeded entity unmutated; (observed) the FO ran `state init` OR surfaced "state not initialized"/`state init`. The durable track is the independent source of truth (it cannot be faked by phrasing); the observed track only disambiguates the silent-no-work failure where nothing moved on disk.

## Proposed approach

**Primary: a new `split-root-halt` shared runtime scenario** added to the existing host-neutral table + both runner maps (the structure the three current scenarios already use):

1. **Fixture** (`internal/ensigncycle/shared_fixtures_test.go`, default tag so the offline negative reuses it): a `writeSplitRootHaltWorkflow` helper that stages an uninitialized split-root — a bare origin with an orphan state branch carrying one seeded entity, and a fresh clone whose `.spacedock-state` path is ABSENT — reusing the proven `commissionSplitWorkflow` mechanic from `internal/cli/state_init_test.go`. The neutral prompt says `Use $spacedock:first-officer` and instructs the FO to inspect startup/boot and proceed; it does NOT tell the FO to halt (the halt must be the FO's own contract-driven decision, else the drive is rigged).
2. **Assertion** `assertSplitRootHaltHeld` (`internal/ensigncycle/shared_assertions_impl_test.go`): the two-track triple the spike proved. Because the entity lives behind the absent state checkout, the durable signal is checked against the workflow root (no `.worktrees/` created, code porcelain clean) and the seeded entity (if the state path got created, the entity is still at its initial stage), NOT the absent-entity before/after. This is the one place the new scenario diverges from the existing three (whose entity is always present) — call it out in the design so the runner adapters read the right paths.
3. **Runners**: `runClaudeSplitRootHaltScenario` + `runCodexSplitRootHaltScenario`, registered in `claudeScenarioRunners()` / `codexScenarioRunners()`, driving the same fixture/prompt/assertion — the host-specific surface stays auth/launch/observed-extract only.
4. **Table + meta**: add `split-root-halt` to `sharedRuntimeScenarios()` (with `oldPythonTest`, `intent`, `timeout`) and to the `want` slice in `TestSharedRuntimeScenarioDefinitions`; `TestSharedScenarioRunnerCoverage` then enforces both hosts carry it.
5. **Offline negative** (`internal/ensigncycle/shared_scenarios_negative_test.go`): build the SPECIFIC broken end-state (a `.worktrees/` dispatch dir present, or the seeded entity advanced) from the real fixture and prove `assertSplitRootHaltHeld` REDS — and that a tautological observed-only check would have stayed green.

**Sync + commission-journey assessment (done in ideation — DO NOT add a sync or commission live drive):**

The sync and commission behaviors are NOT in the same "no oracle at all" position the halt is. Three layers, two already covered by real (non-tautological) oracles:

- **Git sync MECHANICS** (does `pull --rebase` replay disjoint-path commits with zero conflict; does a same-entity edit CONFLICT and halt; does a non-force re-push stay rejected so the peer's edit survives) — **already proven by real-git two-writer e2e** in `internal/cli/state_sync_test.go` (`TestTwoWriterSyncHappyPath`, `TestTwoWriterSameEntityConflictHalts`). Real origin, two clones, real git. Not a tautology.
- **The binary EMITS the correct path-scoped commit + push + `pull --rebase` guidance** into the worker prompt — proven by `internal/dispatch/build_statecommit_test.go` (`TestStateCommitGuidanceResolvesPaths`, and the bare-`git add -A` warning guard). This is a legitimate code gate: the binary is the independent source, the test parses its real emitted command.
- **The orphan/worktree/`state init` MECHANICS** the commission journey performs — already proven by `internal/cli/state_init_test.go` (`TestCommissionOrphanBranchScaffolding`, `TestStateInitResumesFreshClone`). Real git, no mocks.

What the demoted `TestFOSyncProse`/`TestEnsignSyncProse`/`TestCommissionJourneyProse` tautologies were *trying* to assert is the residual third layer: that the FO/ensign AGENT actually ISSUES those git commands at the right moment, and that the commission SKILL agent actually runs the orphan/init sequence. That residual is genuine, but: (a) there is NO `spacedock state sync/push/pull` subcommand — sync is the agent issuing raw `git`, so a live drive would re-prove `git`'s own rebase semantics that `state_sync_test.go` already proves with full determinism; (b) commission is captain-INTERACTIVE (the orphan birth happens inside an interactive design session), which the deterministic shared-scenario runner is a poor fit for. **Decision: the sync residual FOLDS INTO the split-root-halt scenario for free** — the correct halt remedy is the FO running `state init` (the boot pull/sync entry point), so the same live drive already exercises the FO doing the right git-state thing on an uninitialized checkout. The commission-journey residual is recorded as covered-by-mechanics (`state_init_test.go` + `build_statecommit_test.go`); no commission live drive is warranted by this task. The three demoted lints all bind to THIS scenario as their named owed oracle, with the explicit note that sync/journey are oracle-covered at the mechanics+emission layer and the agent-issues-them layer rides the halt drive.

## Acceptance criteria (entity-level; per the proof-policy each AC's proof is RUNNING the behavior, from a source OTHER than the contract files under test)

**AC-1 (behavioral) — A `split-root-halt` shared runtime scenario exists and a real FO driven against an uninitialized split-root holds the halt: no `.worktrees/` dispatch dir is created, code-branch porcelain stays clean, the seeded entity is unmutated, and the FO either ran `state init` or surfaced "state not initialized" / `spacedock state init`.**
Verified by: `TestLiveClaudeSharedScenarios/split-root-halt` and `TestLiveCodexSharedScenarios/split-root-halt` (cited live runs, `-tags live -count=1`) PASS; the durable+observed triple is the assertion.

**AC-2 (negative, offline, mutation-controlled) — `assertSplitRootHaltHeld` REDS on the specific broken end-states (a dispatch worktree was created, or the seeded entity was advanced) built from the real fixture, and a tautological observed-only check would have stayed green on the same broken state.**
Verified by: a new case in `internal/ensigncycle/shared_scenarios_negative_test.go` (default tag, no model spend) that fails the assertion on the broken state and asserts the contrast with an observed-only check.

**AC-3 (parity) — both host runners cover the scenario and the host-neutral table defines it.**
Verified by: `TestSharedScenarioRunnerCoverage` (both runner maps key `split-root-halt`) and `TestSharedRuntimeScenarioDefinitions` (`want` includes `split-root-halt` with non-empty `intent`/`oldPythonTest` and positive `timeout`).

**AC-4 (gap closure) — the hwk demotions for the halt/sync/journey cluster name this scenario as their owed oracle, and the sync/journey residual is recorded as covered (sync mechanics by `state_sync_test.go`, commit-emission by `build_statecommit_test.go`, orphan/init mechanics by `state_init_test.go`; the agent-issues-state-git layer riding the halt drive).**
Verified by: the demoted non-AC lints in hwk reference `split-root-halt`, and `internal/cli/state_sync_test.go` + `internal/dispatch/build_statecommit_test.go` + `internal/cli/state_init_test.go` are green (the named mechanics oracles exist and pass).

## Test plan

- **Spike: DONE this session** (offline + real-binary, throwaway). Result recorded above: halt is gradeable via the two-track triple; observed-only is a confirmed false-pass; sync/journey mechanics already have real oracles.
- **Offline (default tag, no model spend):** the new negative case (AC-2) + the meta/coverage guards (AC-3) run in the normal `go test ./...` sweep. `assertSplitRootHaltHeld` is a pure function over (workflow-root durable signals, seeded-entity state, observed) so its negative cases are free.
- **Live (gated `-tags live -count=1`, ~per-host minutes):** `split-root-halt` under both `TestLiveClaudeSharedScenarios` and `TestLiveCodexSharedScenarios` (AC-1). Timeout class ~90s–2m (a single boot-read-and-halt, no multi-stage cycle — cheaper than rejection-flow). `-count=1` per the docs/dev README clause defeats a stale-cache false-green.
- **No new sync or commission live drive** (assessment decision above). The sync residual rides AC-1; the journey residual is mechanics-covered.
- **Cost/complexity:** one new fixture + one assertion + two thin runner adapters + one table entry + one negative case. Medium. High-stakes (touches the shared-scenario CI machinery + the parity guards) → an independent staff review of the scenario design and the assertion's distinguishing power is warranted before the gate (the ideation stage def calls for this on skill-integration / split-root work).
- Pairs with hwk (the demotions that name this scenario) + the proof-policy entity (`eykb` / `f8b257cf`).

## Notes

Provenance: hwk implementation (this session) surfaced that the halt/sync/journey Bucket-A cluster has no existing live drive; the captain chose (b) — split the new cross-host scenario into this task rather than bloat hwk. Precedent: gq's `feedback-3-cycle-escalation` (a new cross-host live scenario filed as its own entity).

The one structural novelty vs. the three existing shared scenarios: the graded entity lives behind an ABSENT state checkout, so the assertion reads durable signals at the workflow root (no `.worktrees/`, clean porcelain) rather than the entity-body before/after the present-entity scenarios use. Implementation must wire the runner adapters to the workflow root, not an entity path that doesn't exist yet at run start. This is the riskiest implementation seam and the spike's main carry-forward.

## Stage Report: ideation

- DONE: Run the riskiest-unknown spike (grading a HALT = grading the absence of action; can a durable-state assertion distinguish "FO halted" from "FO dispatched against the empty boot")
  Spike RUN this session (throwaway, offline + real binary built from `cmd/spacedock`, contract 1). Result in "Riskiest unknown — SPIKE RUN": the silent-failure trap is real (uninitialized split-root boots EMPTY + `--validate` VALID); pure entity before/after cannot grade it; observed-only is a confirmed false-pass; the two-track triple (durable: no `.worktrees/` + clean porcelain + entity unmutated; observed: ran `state init` OR named the uninitialized state) PASSES both correct remedies and REDS all three broken modes — 5-case offline prototype green. Spike artifacts removed after recording.
- DONE: Design the split-root-halt live scenario (uninitialized split-root fixture + neutral `Use $spacedock:first-officer` prompt + durable-state halt assertion), Claude + Codex runners + offline negative
  See "Proposed approach": `writeSplitRootHaltWorkflow` fixture (reusing the proven `commissionSplitWorkflow` mechanic), `assertSplitRootHaltHeld` (the spiked triple), `run{Claude,Codex}SplitRootHaltScenario` runners, table + meta entry, and the offline negative. The one structural novelty (entity behind an absent state checkout → assert at workflow root) is called out as the riskiest implementation seam.
- DONE: Assess whether the sync (rebase/push/no-force) + commission-journey behaviors need their own drives, fold in, or are already covered
  See "Sync + commission-journey assessment": sync MECHANICS already proven by `internal/cli/state_sync_test.go` (real-git two-writer e2e), commit-emission by `internal/dispatch/build_statecommit_test.go` (code gate), orphan/init mechanics by `internal/cli/state_init_test.go`. There is NO `state sync/push/pull` subcommand — sync is the agent issuing raw `git`. DECISION: no separate sync or commission live drive; the agent-issues-state-git residual FOLDS INTO the halt drive (the correct halt remedy is the FO running `state init`); commission is captain-interactive and mechanics-covered.
- DONE: Write entity-level ACs where each behavior's proof is a cited live run (mutation-controlled), and name how the hwk demotions bind to this scenario as their owed oracle
  AC-1 (cited live runs both hosts), AC-2 (offline negative, mutation-controlled), AC-3 (parity guards), AC-4 (gap closure: demotions name `split-root-halt`; sync/journey residual recorded as mechanics-covered). Each AC's "Verified by" names a test/run outside the contract files under test.

### Summary

The riskiest unknown is resolved empirically: a HALT on an uninitialized split-root IS gradeable, but only via a two-track triple — a durable signal (no dispatch worktree, clean code porcelain, seeded entity unmutated) that is the independent source of truth, plus an observed signal (ran `state init` or named the uninitialized state) that disambiguates the silent-no-work failure. Observed-only grading was proven to false-pass, so the durable track is load-bearing. The sync/commission assessment found those behaviors are NOT oracle-less like the halt: their git mechanics and the binary's command-emission already have real (non-tautological) oracles, so this task adds ONE scenario (split-root-halt) and folds the sync residual into it rather than spawning two more live drives. The single implementation risk to flag at the gate is the structural novelty: the graded entity sits behind an absent state checkout, so the runner adapters and assertion read the workflow root, not a not-yet-existing entity path.
