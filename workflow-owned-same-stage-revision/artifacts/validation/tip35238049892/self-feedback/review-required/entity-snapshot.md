---
id: recorded-gate-task
title: Recorded Gate Task
status: validation
completed:
verdict:
worktree:
gates:
    version: 1
    records:
        - id: gate:recorded-gate-task:validation
          stage: validation
          attempts:
            - id: gate-attempt:recorded-gate-task-validation-1
              briefing:
                id: briefing:recorded-gate-task:validation:attempt-1:revision-1
                digest: sha256:7c076e9869ccfbc6adfcba59ddded9af53ba19d38d6c415fa1b69868c13315cc
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:recorded-gate-task:validation:1
                briefing: briefing:recorded-gate-task:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-17T15:12:38.410238125Z"
                decision: revise
                reason: Correct the plan against frozen input
started: 2026-09-17T15:15:12Z
---
# Recorded Gate Task

## Acceptance criteria

**AC-1** The selected plan matches the frozen input.

## Stage Report: validation

- DONE: Replayed retained evidence
  The real command fixture is green.

### Summary

Ready for the recorded decision gate.

## Stage Report: validation (revision 1)

- DONE: selected/plan.md corrected to match selected/frozen-input.txt exactly (AC-1)
  Corrected selected/plan.md from "KEEP message B; DELETE message A" to "KEEP message A; DELETE message B" to agree with the frozen, authoritative selected/frozen-input.txt. selected/frozen-input.txt was not altered.
- DONE: appended and committed this Stage Report: validation section documenting the correction

### Summary

Corrected selected/plan.md per resolution:spacedock:recorded-gate-task:validation:1 (decision: revise) so it matches the frozen input exactly. Ready for re-review.
