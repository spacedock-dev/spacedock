---
title: Cut gate-guardrail turn and tool bloat
status: ideation
source: "Journey-metrics audit of PR #643, Runtime Live E2E run 31297186020, Claude job 93204212216, artifact 9033837253. The gate-guardrail journey used 22 assistant turns and 24 tool calls, up from 11 and 11 in the v0.26 Sonnet observation. Captain directed a separate filing on 2026-08-09."
started:
completed:
verdict:
score: 0.9
sprint: durable-decisions
sprint-readiness: ready
group: gate-lifecycle-ux
worktree:
issue:
pr:
mod-block:
id: 5k704rrfk5r75vqv3bwn1yhf
gates:
    version: 1
    records:
        - id: gate:5k704rrfk5r75vqv3bwn1yhf:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:5k704rrfk5r75vqv3bwn1yhf-backlog-1
              briefing:
                id: briefing:5k704rrfk5r75vqv3bwn1yhf:backlog:attempt-1:revision-1
                digest: sha256:bdd851d207fe8f84300b59cbd34359f1cc613df5fee26b869f323db6e63f7956
                request-digest: sha256:ef783c09a07b51873e1f7c8d3fef81b226bfa69901afc1a652389886a6d7b2d2
                room-ref: ./cut-gate-guardrail-turn-and-tool-bloat/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:5k704rrfk5r75vqv3bwn1yhf:backlog:1
                briefing: briefing:5k704rrfk5r75vqv3bwn1yhf:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T14:50:42.1944Z"
                decision: approve
                reason: Captain directed ideation dispatch; the measured-call classification constrains the smallest safe change.
              application:
                target-stage: ideation
                state: consumed
---

Make the supported gate-guardrail journey reach one committed open gate with materially fewer turns and tool calls.

## Problem

The latest successful journey used 22 assistant turns and 24 tool calls. The v0.26 Sonnet observation used 11 turns and 11 calls.

The current log includes filesystem searches for known paths, repeated Git and state reads, a gate help probe, and two boot projections after preparation.

The evidence identifies the waste, but it does not prove one mechanism. Some calls can come from the contract, missing command output, or agent error.

## Required outcome

During ideation, classify each call in the retained run as required, contract-caused, missing CLI information, or agent error.

Then ship the smallest product or contract change that removes the avoidable calls. Do not prescribe a new projection before this classification.

Preserve the journey result: one prepared and committed open gate is presented. No decision, consume, successor dispatch, or unauthorized state change occurs.

Do not add a scheduler, lifecycle state, standing check, simulator, or test that tests the journey harness.

## Acceptance criteria

**AC-1 (VALUE) - The real gate journey uses materially less operator effort.**
Verified by: one candidate run with the same Claude model and journey uses at most 16 assistant turns and 18 tool calls. Compare it with run `31297186020`.

**AC-2 - Known search and retry waste is absent.**
Verified by: the candidate log contains no filesystem search for a known workflow or contract path, no failed command-shape probe, and no repeated boot projection after gate preparation.

**AC-3 - Gate authority and final state stay correct.**
Verified by: durable state contains exactly one committed open gate attempt and no decision, consumption, successor dispatch, or unauthorized mutation.

**AC-4 - The fix changes the owning surface only.**
Verified by: ideation maps each removed call to its owner and the implementation changes only the smallest owning CLI or contract surface. No test-harness controller is added.

## Test plan

Use the retained run as the baseline. Classify its 24 calls before selecting a mechanism.

Run focused checks for the selected product or contract behavior. Then run the real gate-guardrail journey once and compare turns, calls, command failures, and durable state.

Do not add a permanent metrics simulator or a static prose-presence test.

