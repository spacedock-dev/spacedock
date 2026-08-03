---
title: Add Pi as a first-class adapter for every common live journey
status: ideation
source: "Desired live-test registry requires common journeys on every supported runtime, 2026-08-03"
started: 2026-08-03T10:48:52Z
completed:
verdict:
score: 0.9
worktree:
issue:
sprint: live-test-truth
group: pi-common-runner
sprint-readiness: ready
id: tj41e4f404mz7ast3yh9enwc
gates:
    version: 1
    records:
        - id: gate:tj41e4f404mz7ast3yh9enwc:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:tj41e4f404mz7ast3yh9enwc-backlog-1
              briefing:
                id: briefing:tj41e4f404mz7ast3yh9enwc:backlog:attempt-1:revision-1
                digest: sha256:75c779adf9f75759ab31035e7433523ee1a64a5a351eb725e3bc838061a5a7ca
                request-digest: sha256:b2e61d14956245e69ed188a2d46d2d53dd2bc2e1c2af97fa94def22578bc1346
                room-ref: ./add-pi-common-live-runner/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:tj41e4f404mz7ast3yh9enwc:backlog:1
                briefing: briefing:tj41e4f404mz7ast3yh9enwc:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T10:34:56.368971Z"
                decision: approve
                reason: Captain approved the prepared Sol ideation cohort with make it so.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
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

