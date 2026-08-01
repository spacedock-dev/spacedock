---
title: Minimize the unreleased v1 gate application schema
status: ideation
source: "Pre-0.27 gate-machinery necessity audit, 2026-08-01: production emits empty blockers but has no producer or demonstrated consumer for blockers, execution holds, or feedback payloads."
started:
completed:
verdict:
score: "0.9"
worktree:
issue:
pr:
sprint: durable-decisions
id: nthcevf1snz7hm75gny3kd2e
gates:
    version: 1
    current:
        gate: gate:nthcevf1snz7hm75gny3kd2e:backlog
    records:
        - id: gate:nthcevf1snz7hm75gny3kd2e:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:nthcevf1snz7hm75gny3kd2e-backlog-1
              briefing:
                id: briefing:nthcevf1snz7hm75gny3kd2e:backlog:attempt-1:revision-1
                digest: sha256:a969f6243a82f97def080d4c51f69732932877a9a96392b3df5ef48b28a7ba8f
                digest-domain: canonical-bytes
                request-digest: sha256:1c1c9a1b2ce8ff70509a08507db6e35cbea319b2411678da6e2a7d9aa3af160d
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:nthcevf1snz7hm75gny3kd2e:backlog:1
                briefing: briefing:nthcevf1snz7hm75gny3kd2e:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-01T13:53:46.825149Z"
                decision: approve
                reason: Captain approved dispatching this durable-decisions ideation lane in parallel with wj, hq, and jc.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

Keep only application state that a supported command creates and a supported consumer spends. The unreleased v1 model currently carries blockers, execution holds, feedback payloads, and leaf fields that were designed speculatively but have no live producer or consumer.

## Problem

The active state audit found 16 gate-bearing tickets. Fourteen contain only the mechanically emitted `blockers: []`; none contains a nonempty blocker, execution hold, or application feedback payload. Strict decoding makes every unused field a permanent public-format commitment. `hold` also receives a `none/not-applicable` application and `revise` receives a feedback application even though the Resolution already carries the decision and the workflow router owns rework.

## Proposed approach

Ideation must derive the minimum application shape from the actual consumers after `1w6` lands. The expected baseline is that only an approval which authorizes an advance owns an application; its target and one-use state remain. A hold or revise remains a durable Resolution without a dead application wrapper. Remove blocker, execution-hold, and feedback-only fields and their unreachable validators. Because v1 is unreleased, update pilot state and fixtures once and make the clean schema the only accepted form; do not ship dual decoding or migration code.

## Streamlined common scenarios

- Approve a nonterminal gate: record one pending advance application, then `gate consume` atomically spends it with the stage transition.
- Approve a terminal gate: record one pending application; `gate consume` routes toward delivery without spending, and `1w6` makes the terminal delivery transaction the sole spender.
- Hold: record the Resolution and stop. No application exists because there is nothing to apply.
- Revise: record the Resolution and let the workflow feedback route select the declared stage. No application subtree duplicates that route.

## Out of scope

Do not redesign terminal delivery, current-gate selection, review-round storage, or provider presentation. Do not retain prototype fields for compatibility.

## Acceptance criteria

**AC-1 - Every retained application field has both a supported producer and a supported consumer.**
Verified by: a package-level producer/consumer matrix backed by focused tests for approve, hold, and revise; a hold or revise that emits an application, or an application field with no exercised producer and consumer, fails the matrix.

**AC-2 - The clean v1 schema is the only schema accepted after the cut.**
Verified by: strict-decode negative fixtures for each removed field, positive fixtures for pending/consumed/superseded approval applications, and `spacedock status --workflow-dir docs/dev --validate` against the one-off-normalized active and archived pilot state.

**AC-3 - The ordinary chat gate journey remains shorter without losing authority evidence.**
Verified by: a real-CLI fixture covering approve and consume plus separate hold and revise records, asserting the exact minimal YAML and unchanged Resolution identity.

## Test plan

Change the model, validation, encoder, fixtures, and pilot state in one coordinated unreleased-v1 cut. Use focused Go tests, the full Go and race suites, and state validation. Do not create compatibility branches.
