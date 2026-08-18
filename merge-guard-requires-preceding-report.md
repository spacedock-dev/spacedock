---
id: vkatjs25g9a9gmk3jtvx5ce0
title: merge guard refuses a terminal transition with no preceding worker report
status: backlog
source: "Captain CL, 2026-08-18, from the live-lane inventory reframe. Failing assertion: internal/ensigncycle/shared_keep_moving_durable_test.go:103, 'first terminal transition must follow worker report', red in two consecutive claude-live runs (32092321763 attempt 2 and 32105482382) while the FO's own final messages claimed the reports had completed."
started:
completed:
verdict:
score:
worktree:
issue:
---

The terminal ceremony checks that the merge hook ran. It does not check that the work it is finalizing was ever reported. Make it refuse when no report commit precedes the terminal transition.

## Problem

`merge guard` drives the terminal ceremony as one ordered envelope — arm the mod-block, detect hook completion by state delta, terminalize, archive with a path-scoped commit, publish state. It owns the sequence so steps cannot be combined, skipped, or reordered.

It verifies the hook. It never verifies the report.

So an FO that trusts a worker's completion message over durable state can finalize an entity whose report never landed, and the ceremony's own atomicity hides it: one verb takes the entity from validation to archived-and-pushed with no intermediate state where the absence would show.

Observed twice in consecutive claude-live runs. The durable check that caught it walks the entity's git history and requires the terminal-transition commit to follow a commit carrying the completing stage's report.

The correct shape is already what a healthy run produces, verified on this session's own terminalizations: the worker commits its report touching no frontmatter (`+15` and `+43` lines, zero status change), then `merge guard` flips `status`, `verdict`, and `completed` in a later commit. Different writers, naturally ordered. Nothing has to change about how a good run behaves.

Two sub-shapes are possible and the evidence does not separate them: the report and the flip batched into one commit, or the report never landing at all. The second is the likelier one, because our deliberate coalescing batches FO-side state writes together and never mixes them with a worker's report.

## Proposed approach

{Ideation fills this in, and should settle the sub-shape first from the retained streams of the two failing runs — the fix differs. The check itself: before terminalizing, require a commit in the entity's own history that carries the completing stage's report and precedes the transition. All of it is state `merge guard` already reads.}

## Out of scope

Changing what the merge hook does, the archive step, or the ordering the envelope already enforces. The `--rework` path unless the same gap applies there.

## Expected surface and tolerance

Estimate net LOC change: +60 across 2 files. Report insertions and deletions separately. Do not declare a gross tolerance. Semantics changed: `merge guard` gains a refusal condition, so a previously-accepted call can now fail.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - No entity reaches a terminal stage without its completing stage's report preceding the transition in durable history.**
This is the measuring AC: the count of terminal transitions with no preceding report commit must be ZERO. Verified by driving `merge guard` against a fixture entity whose report was never committed and observing a refusal, then committing the report and observing it proceed. Fails on today's binary, which finalizes either way.

**AC-2 - A healthy run is unaffected.**
Verified by replaying this session's real terminalization shape — worker report commit touching no frontmatter, then the terminal flip — and observing `merge guard` proceed unchanged. Fails if the guard refuses a correctly-ordered run, which would block every normal merge.
