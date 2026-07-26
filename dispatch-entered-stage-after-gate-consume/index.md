---
title: Dispatch the entered working stage after gate consumption
status: backlog
source: "Codex live run 30197794474 on 2026-07-26: rejection-flow began at implementation, but the FO advanced directly to validation without an implementation worker/report; the strict two-cycle provenance assertion correctly failed."
started:
completed:
verdict:
score: 1.0
worktree:
issue:
sprint: durable-decisions
id: gqsw81ghf48hr2n3jg6k7nx8
---

A gate application can atomically move a ticket into a non-gated working stage, but the First Officer must still dispatch that entered stage before advancing again. In the observed Codex journey, the ticket began at `implementation`; the FO skipped its first worker, moved directly to `validation`, and later produced only the rework implementation report. The final functionality looked repaired, but the durable history lost the original implementation round.

This is not a request to weaken the live assertion. The sprint walking skeleton depends on the same consume-to-successor handoff, so the missing dispatch ownership must be made explicit and falsifiable.

## Boundary to shape

Determine the smallest host-neutral correction that makes an entered, unstarted working stage dispatchable exactly once. First exercise whether the existing `gate consume` target-stage output plus lifecycle prose is sufficient, or whether the status/boot projection lacks a necessary current-stage readiness signal. Do not create a second scheduler, infer completion from status alone, or special-case the rejection fixture.

## Acceptance criteria

**AC-1 (VALUE)** In a real two-cycle rejection journey, the durable ticket contains the original implementation report, first REJECTED validation, rework implementation report, and second PASSED validation in order; no stage is credited from narration or status mutation alone.

**AC-2 (DISPATCH)** After a gate consume enters a non-gated working stage, an observed worker spawn and that worker's valid Stage Report occur before any mutation to the following stage. A dispatch package build, `started` timestamp, or FO-authored report cannot substitute.

**AC-3 (RECOVERY)** A cold boot after consume but before spawn surfaces one unambiguous dispatch action; a boot after the durable report does not duplicate that stage's worker.

**AC-4 (SCOPE)** The correction composes with existing consume, `status --next`, standing dispatch, fresh/reuse, and feedback routing surfaces without new scheduler state or weakening current negative controls.

## Test plan

Begin with the exact archived Codex transcript and current rejection fixture. Add the cheapest offline negative that fails when implementation is skipped, then run the focused real Codex journey locally. Include a consume-then-crash cold-boot fixture and a report-present no-duplicate control. Ideation must declare exact instruction/code/test files and LOC before implementation.
