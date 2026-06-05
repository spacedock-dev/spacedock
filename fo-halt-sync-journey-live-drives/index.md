---
id: ev3e0nmknh98sn365ky163en
title: Behavioral live drives for the FO halt/sync/journey contract behaviors (the Bucket-A cluster with no existing drive)
status: backlog
source: "captain (2026-06-05) — split from hwk (tautological-test-remediation) per the (b) phasing call. hwk demotes the halt/sync/journey presence checks to honest non-AC lints; this task supplies their real behavioral proof — a new cross-host live scenario (split-root-halt) plus the sync/journey behaviors that today have NO live drive."
score: "0.33"
started:
completed:
verdict:
worktree:
issue:
---

The tautological-test remediation (hwk) found that most Bucket-A behavioral claims already have live drives (present-gate → gate-guardrail scenario; feedback-rejection-flow → rejection-flow scenario; using-claude-team → the team scenarios), so those demote-and-bind. But four Bucket-A tests assert FO behaviors that have **no live drive at all** — they are pure presence checks (banned tautologies) and hwk demotes them to honest non-AC lints, leaving the behaviors unproven. This task supplies the owed behavioral proof.

## Problem — the behaviors with no drive

- **`TestFOHaltGateProse`** — the FO must HALT on an uninitialized split-root state dir (orphan branch, no linked worktree) and point at `spacedock state init`, NOT dispatch against an empty/EMPTY-rendering boot table. This is a SAFETY behavior (a silent failure if it doesn't halt).
- **`TestFOSyncProse` / `TestEnsignSyncProse`** — the state-checkout sync rules (path-scoped commits, `pull --rebase` on push-rejection, no `--force`).
- **`TestCommissionJourneyProse`** — the commission journey performs the orphan/worktree/init sequence.

Today all four are presence checks over the contract prose — they prove the words are present, never the behavior. hwk honestly demotes them (Oracle-honesty non-AC lint naming THIS task as the owed oracle).

## Proposed approach (ideation formalizes)

Primary: a new **split-root-halt live scenario** — a fixture with an uninitialized split-root state checkout, the neutral `Use $spacedock:first-officer` prompt, and a durable-state assertion that the FO HALTED (did NOT dispatch, did NOT mutate entity state; surfaced the "state not initialized" report / `state init` pointer). Claude + Codex runners (the `shared_coverage_meta` parity guard enforces both hosts) + an offline negative. Structurally the same as gq's `feedback-3-cycle-escalation` scenario.

Then assess the sync + commission-journey behaviors: which need their own live drive vs. are covered by an existing cycle (the live-cycle smoke already exercises state push/pull) vs. fold into the split-root scenario.

## Riskiest unknown (spike in ideation)

Grading a HALT is grading the ABSENCE of action — can a durable-state assertion cleanly distinguish "FO correctly halted" (no dispatch, no entity mutation, the init pointer surfaced) from "FO failed to halt and dispatched against the empty table"? Exercise the smallest version before committing — likely: seed the uninitialized state, drive, and assert the entity set is untouched + no worktree created + the halt report present. Record the result (this is the novel piece; the cross-host runner machinery is already proven by gq/the shared table).

## Acceptance criteria (ideation formalizes; per the proof-policy each AC's proof is RUNNING the behavior)

**AC-1 (behavioral) — a split-root-halt live scenario drives a real FO and reds on the broken behavior (FO dispatches instead of halting).** Verified by: `TestLive{Claude,Codex}SharedScenarios/split-root-halt` (cited live runs) + an offline negative that reds on a dispatched-against-uninitialized end-state, mutation-controlled.
**AC-2 (parity) — both host runners cover it.** Verified by `TestSharedScenarioRunnerCoverage` / `TestSharedRuntimeScenarioDefinitions`.
**AC-3 — the hwk demotions for the halt/sync/journey cluster now reference a real behavioral oracle (this scenario), closing the gap hwk tracked.** Verified by: the demoted lints in hwk name this scenario, and the behavior is driven here.

## Test plan

- Spike: the halt-grading question (offline, cheap). Then the live scenario + runners + offline negative (live-gated, ~per-host minutes). Sync/journey drives assessed in ideation. High-stakes (CI/scenario machinery) → detached audit.
- Pairs with hwk (the demotions) + `eykb`/`f8b257cf` (the proof-policy).

## Notes

Provenance: hwk implementation (this session) surfaced that the halt/sync/journey Bucket-A cluster has no existing live drive; the captain chose (b) — split the new cross-host scenario into this task rather than bloat hwk. Precedent: gq's `feedback-3-cycle-escalation` (a new cross-host live scenario filed as its own entity).
