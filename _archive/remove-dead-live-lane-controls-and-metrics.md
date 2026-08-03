---
title: Remove or connect dead live-lane inputs and journey metrics
status: backlog
source: "Live CI inventory found decorative effort input and unconsumed metrics paths, 2026-08-03"
started:
completed:
verdict:
score: 0.75
worktree:
issue:
sprint:
group: lane-hygiene
sprint-readiness: absorbed
id: b8wfn1hm4p8208a99e5mh49z
---

## Problem

The runtime-live workflow exposes an effort input that reaches summaries but not host configuration, exports a Pi metrics directory with no emitter, and uploads Codex journey metrics that the delta-comment job does not consume. These surfaces imply control or evidence they do not provide.

## Acceptance criteria

**AC-1 (VALUE)** Every retained live-lane input changes an actual supported runtime setting, and every retained metrics artifact has a named producer and consumer.
Verified by: workflow-level tests or a one-off execution trace showing each retained input and artifact crossing its declared boundary.

**AC-2** Decorative effort and dead Pi metrics configuration are removed unless ideation demonstrates a cheaper real connection with user value.
Verified by: workflow dispatch/schema inspection and lane execution.

**AC-3** Codex metrics are either connected to the existing comparison consumer or removed from production and upload paths without affecting Claude metrics.
Verified by: artifact-path tests and a representative workflow run.

## Stage-specific test gates

- Ideation begins by proving the current producer/consumer graph and chooses deletion as the default for unconsumed surfaces.
- Validation checks workflow syntax, artifact behavior, and full/race suites.

