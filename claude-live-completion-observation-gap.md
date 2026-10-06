---
title: Claude live journeys cannot observe the completion path the agent reports
sprint: test-behavior-completeness
source: "Run 37421110550 claude-live: the first officer narrated advancing the entity, dispatching a fresh validator, and AC-1 passing with the stage report committed, then stopped at the gate; the harness recorded validation-worker-not-dispatched. Its message noted the validator's SendMessage to team-lead failed."
id: jfdrr2dbdtt952s4cvvy74k8
---

## Problem

At least one Claude journey grades a run as incomplete while the agent's own narration describes a
complete, correct lifecycle. The harness observes completion through a path this agent does not use,
so a correct run is graded red.

## Value

This is the largest source of Claude lane flakiness: the lane reddens on observation, not on
behaviour, which makes its failures uninformative.

## Acceptance criteria

**AC-1** For the affected journeys, identify how the agent signals completion and whether the harness reads that path.
**AC-2** Either make the harness observe it, or record the gap as known with an owner and a binding.
**AC-3** Record the per-journey verdict distribution across at least five runs.

## Verification

Read the journeys' final messages and the harness assertions side by side.
