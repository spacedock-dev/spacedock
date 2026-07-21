---
title: Two-gate multi-attempt contract replay
gates:
  version: 1
  current:
    gate: gate:docs:dev:falsifiability-ladder:validation
    attempt: gate-attempt:z7cvbvdv-validation-3
  records:
    - id: gate:docs:dev:falsifiability-ladder:ideation
      stage: ideation
      current-attempt: gate-attempt:z7cvbvdv-ideation-5
      attempts:
        - id: gate-attempt:z7cvbvdv-ideation-1
          sequence: 1
          state: closed
          briefing:
            id: briefing:z7cvbvdv-ideation-1
            digest: sha256:1111111111111111111111111111111111111111111111111111111111111111
            note: RAW-FILE PIN retained production history
          resolution:
            type: Resolution
            id: resolution:z7cvbvdv-ideation-1
            briefing: briefing:z7cvbvdv-ideation-1
            by: person:captain
            at: 2026-07-16T10:00:00Z
            decision: revise
            reason: Separate the portable decision from workflow application.
          application:
            action: feedback
            target-stage: ideation
            state: consumed
            feedback:
              correction-round: 1
        - id: gate-attempt:z7cvbvdv-ideation-2
          sequence: 2
          previous-attempt: gate-attempt:z7cvbvdv-ideation-1
          state: closed
          briefing:
            id: briefing:z7cvbvdv-ideation-2
            digest: sha256:2222222222222222222222222222222222222222222222222222222222222222
          resolution:
            type: Resolution
            id: resolution:z7cvbvdv-ideation-2
            briefing: briefing:z7cvbvdv-ideation-2
            by: person:captain
            at: 2026-07-17T10:00:00Z
            decision: hold
            reason: Preserve the provider portability boundary.
          application:
            action: none
            target-stage: ideation
            state: not-applicable
        - id: gate-attempt:z7cvbvdv-ideation-3
          sequence: 3
          previous-attempt: gate-attempt:z7cvbvdv-ideation-2
          state: closed
          briefing:
            id: briefing:z7cvbvdv-ideation-3
            digest: sha256:3333333333333333333333333333333333333333333333333333333333333333
            room-ref: ./review/ideation
          resolution:
            type: Resolution
            id: resolution:z7cvbvdv-ideation-3
            briefing: briefing:z7cvbvdv-ideation-3
            by: person:captain
            at: 2026-07-18T10:00:00Z
            decision: approve
          application:
            action: advance
            target-stage: implementation
            state: superseded
            blockers:
              - kind: evidence
                state: satisfied
        - id: gate-attempt:z7cvbvdv-ideation-4
          sequence: 4
          previous-attempt: gate-attempt:z7cvbvdv-ideation-3
          state: closed
          briefing:
            id: briefing:z7cvbvdv-ideation-4
            digest: sha256:4444444444444444444444444444444444444444444444444444444444444444
          resolution:
            type: Resolution
            id: resolution:z7cvbvdv-ideation-4
            briefing: briefing:z7cvbvdv-ideation-4
            by: person:captain
            at: 2026-07-19T10:00:00Z
            decision: approve
          application:
            action: advance
            target-stage: implementation
            state: consumed
          scope-amendment:
            decision: preserve the approved recorder-only scope
        - id: gate-attempt:z7cvbvdv-ideation-5
          sequence: 5
          previous-attempt: gate-attempt:z7cvbvdv-ideation-4
          state: closed
          briefing:
            id: briefing:z7cvbvdv-ideation-5
            digest: sha256:8888888888888888888888888888888888888888888888888888888888888888
          resolution:
            type: Resolution
            id: resolution:z7cvbvdv-ideation-5
            briefing: briefing:z7cvbvdv-ideation-5
            by: person:captain
            at: 2026-07-19T12:00:00Z
            decision: approve
          application:
            action: advance
            target-stage: implementation
            state: consumed
    - id: gate:docs:dev:falsifiability-ladder:validation
      stage: validation
      current-attempt: gate-attempt:z7cvbvdv-validation-3
      attempts:
        - id: gate-attempt:z7cvbvdv-validation-1
          sequence: 1
          state: closed
          briefing:
            id: briefing:z7cvbvdv-validation-1
            digest: sha256:5555555555555555555555555555555555555555555555555555555555555555
          resolution:
            type: Resolution
            id: resolution:z7cvbvdv-validation-1
            briefing: briefing:z7cvbvdv-validation-1
            by: person:captain
            at: 2026-07-20T10:00:00Z
            decision: revise
            reason: Strengthen conflict evidence.
          application:
            action: feedback
            target-stage: implementation
            state: consumed
        - id: gate-attempt:z7cvbvdv-validation-2
          sequence: 2
          previous-attempt: gate-attempt:z7cvbvdv-validation-1
          state: closed
          briefing:
            id: briefing:z7cvbvdv-validation-2
            digest: sha256:6666666666666666666666666666666666666666666666666666666666666666
          resolution:
            type: Resolution
            id: resolution:z7cvbvdv-validation-2
            briefing: briefing:z7cvbvdv-validation-2
            by: person:captain
            at: 2026-07-21T10:00:00Z
            decision: approve
          application:
            action: advance
            target-stage: done
            state: superseded
            execution-hold:
              active: true
              reason: await final evidence
        - id: gate-attempt:z7cvbvdv-validation-3
          sequence: 3
          previous-attempt: gate-attempt:z7cvbvdv-validation-2
          state: closed
          briefing:
            id: briefing:z7cvbvdv-validation-3
            digest: sha256:7777777777777777777777777777777777777777777777777777777777777777
            provider-snapshot: retained
          resolution:
            type: Resolution
            id: resolution:z7cvbvdv-validation-3
            briefing: briefing:z7cvbvdv-validation-3
            by: person:captain
            at: 2026-07-22T10:00:00Z
            decision: approve
            provider-audit:
              result: verified
          application:
            action: advance
            target-stage: done
            state: pending
            blockers:
              - kind: release
                state: satisfied
  fixture-purpose: exercise two logical gates with stable re-entry histories
---
# Two-gate multi-attempt contract replay

This fixture exercises the contract's multi-gate, re-entry, application, and extension
shapes. Separate frozen fixtures replay all eight 0260 production entities.
