---
title: Promote standalone common behaviors into the shared live journey suite
status: ideation
source: "Desired live-test registry inventory, 2026-08-03"
started:
completed:
verdict:
score: 0.9
worktree:
issue:
sprint: live-test-truth
group: common-runner
sprint-readiness: ready
id: r4qk46605sjcphj44cvkcsk4
gates:
    version: 1
    records:
        - id: gate:r4qk46605sjcphj44cvkcsk4:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:r4qk46605sjcphj44cvkcsk4-backlog-1
              briefing:
                id: briefing:r4qk46605sjcphj44cvkcsk4:backlog:attempt-1:revision-1
                digest: sha256:88c06435c9d215a9101c54108777e7553ffbfd61eed6a8cc041d02829236d274
                request-digest: sha256:8f3ec87bf5943cdeed8428bfcc78f626e9f0b108612fe2037f11f506de0f30b6
                room-ref: ./promote-standalone-common-live-journeys/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:r4qk46605sjcphj44cvkcsk4:backlog:1
                briefing: briefing:r4qk46605sjcphj44cvkcsk4:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T11:20:28.996942Z"
                decision: approve
                reason: Captain directed the First Officer to continue the next risk-first ideation wave.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

## Problem

Full-cycle, default headless gate stop, withdrawn-gate recovery, zero discovery, and auto-continue behavior are common runtime promises but currently live in standalone or Claude-specific entry points. Their placement prevents uniform host coverage and hides missing implementations.

## Acceptance criteria

**AC-1 (VALUE)** Each named behavior is registered and runnable through its canonical TestLiveSharedScenarios/<journey-id> entry point with the registry fixture contract and host-neutral grading.
Verified by: focused shared-runner invocations for every promoted journey on each implemented adapter.

**AC-2** Auto-continue preserves both single-root and split-root fixture variants under one intention and one durable assertAutoContinue grader.
Verified by: both fixture variants fail under a stop-after-implementation negative control and pass under correct advancement.

**AC-3** Inline zero-discovery and pre-gate setup become named fixture builders, and the realistic lifecycle fixture has one named builder binding.
Verified by: focused fixture tests under default tags where no live host is required.

## Stage-specific test gates

- Ideation identifies overlaps with gate-guardrail and avoids duplicate assertions.
- Implementation changes no runtime-specific transport behavior.
- Validation runs promoted journeys through available live adapters plus full and race suites.

