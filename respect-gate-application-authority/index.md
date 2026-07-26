---
title: Apply a recorded gate only within the First Officer's assigned authority
status: backlog
source: "Durable-decisions real-sprint correction, 2026-07-26: gate record correctly preserved ideation, but a Shaping FO invoked Commander-owned gate consume because the shipped lifecycle treated approval as authority to apply."
started:
completed:
verdict:
score: 1.0
worktree:
issue:
sprint: durable-decisions
id: mnea9vq3pv1rz1x1hdjbvdg9
gates:
    version: 1
    current:
        gate: gate:mnea9vq3pv1rz1x1hdjbvdg9:backlog
    records:
        - id: gate:mnea9vq3pv1rz1x1hdjbvdg9:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:mnea9vq3pv1rz1x1hdjbvdg9-backlog-1
              briefing:
                id: briefing:mnea9vq3pv1rz1x1hdjbvdg9:backlog:attempt-1:revision-1
                digest: sha256:064cebf9ee6699261c5213d4b8b9ff42350c64e13060bcd07876a995c7a8e8e7
                digest-domain: canonical-bytes
                request-digest: sha256:4fb4b00f0b0d4dd7dfe9be7b26b4dd904d00a9ed3f7361c2e73dd24871a1fcef
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:mnea9vq3pv1rz1x1hdjbvdg9:backlog:1
                briefing: briefing:mnea9vq3pv1rz1x1hdjbvdg9:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T11:42:35.622923Z"
                decision: approve
                reason: 'The sprint''s real role error proves a narrow general contract gap: recording approval must not imply transition authority, while an explicit conn still permits consume and dispatch. Shape the contract and existing live proof without changing gate mechanics.'
              application:
                action: advance
                target-stage: ideation
                state: pending
                blockers: []
---

Recording a binding approval and applying it are separate operations. The shipped First Officer lifecycle must preserve that separation at the role boundary: presenting or recording an approval never enlarges the current session's assigned transition scope. A Shaping First Officer explicitly assigned to hold members at a gate records the exact Resolution and stops with approved-awaiting-advance; a Commander with the conn consumes the same frozen attempt, advances, and dispatches the entered stage.

This is a general authority rule, not a development-workflow stage name. It must compose with ordinary single-FO operation: when the current First Officer's explicit assignment or conn includes applying the transition, approval may proceed directly through consume and successor dispatch.

## Acceptance criteria

**AC-1 (VALUE)** A First Officer explicitly assigned to record an approval but leave application to another role records the binding Result, commits it, reports approved-awaiting-advance, and does not call `gate consume`, mutate status, or dispatch the successor.

**AC-2 (VALUE)** A cold First Officer explicitly assigned the Commander/application role discovers that recorded approval, verifies eligibility, consumes it once, commits the transition, and dispatches the entered working stage without reconstructing or replacing the decision.

**AC-3 (GENERALITY)** The contract is phrased in terms of declared transition authority and assignment scope, not `ideation`, `implementation`, Shaping, or Commander names; a normal FO with explicit conn continues to record, consume, and drive without an artificial stop.

**AC-4 (PROOF)** Existing deterministic skill fixtures prove structural closure and absence only. An applicable existing live FO journey proves both natural-language assignments: record-and-handoff stops after recording, while authorized application consumes and dispatches. No new harness or compatibility layer is added.

## Scope

Shape the narrow change to the shared gate lifecycle/First Officer contract and its existing fixtures/live lane. Do not alter gate recorder, eligibility, consume semantics, gate schema, or provider integration.
