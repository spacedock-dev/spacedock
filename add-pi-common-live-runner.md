---
title: Add Pi as a first-class adapter for every common live journey
status: backlog
source: "Desired live-test registry requires common journeys on every supported runtime, 2026-08-03"
started:
completed:
verdict:
score: 0.9
worktree:
issue:
sprint: live-test-truth
group: pi-common-runner
sprint-readiness: ready
id: tj41e4f404mz7ast3yh9enwc
---

## Problem

Pi is a supported runtime target but currently has a coverage map, a front-door smoke, and one quarantined recorded-gate test rather than the common live journey runner. Consequently no common journey is proven across all supported targets.

## Acceptance criteria

**AC-1 (VALUE)** Every common registry journey is invoked through the canonical shared entry point on Pi and is graded by the same fixture contract and host-neutral assertion as Claude and Codex.
Verified by: the Pi live lane running the complete registered common suite with no coverage-map-only substitutions.

**AC-2** Pi-specific launch, package, child-agent, session, and artifact behavior stays within the Pi adapter; common scenarios contain no Pi branching.
Verified by: adapter-focused tests and source inspection paired with a real Pi run.

**AC-3** Missing or quarantined Pi journey support is visible as a failing or skipped required case rather than represented as live evidence by a coverage-map entry.
Verified by: a negative fixture demonstrating that an absent Pi runner cannot satisfy suite completeness.

## Stage-specific test gates

- Ideation must price live cost and sequence the cheapest non-mutating journeys first without weakening the all-journey end state.
- Spike one existing shared fixture through Pi before designing bulk migration.
- Validation requires the exact Pi live suite and full offline/race tests.

