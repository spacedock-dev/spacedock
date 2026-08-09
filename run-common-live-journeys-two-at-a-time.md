---
title: Run Claude common live journeys two at a time
status: ideation
source: "Captain request after PR #643 CI timing review, 2026-08-09. Run 31295813569 measured Claude common plus substrate at about 31 minutes. Codex was explicitly deferred until later."
started: 2026-08-09T05:56:00Z
completed:
verdict:
score: 0.75
sprint:
sprint-readiness: ready
group: live-ci-performance
worktree:
issue:
pr:
mod-block:
id: c1nbf39akyh263a4p1ts593m
gates:
    version: 1
    records:
        - id: gate:c1nbf39akyh263a4p1ts593m:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:c1nbf39akyh263a4p1ts593m-backlog-1
              briefing:
                id: briefing:c1nbf39akyh263a4p1ts593m:backlog:attempt-1:revision-1
                digest: sha256:da85889f12776fd63a653ff53047fa42a77f2bba3a911f7daf4e688242ae2dbf
                request-digest: sha256:d57c94fc26cc64f0a8030d97b47088f5998ae038480845b41f3d5e141b3c3026
                room-ref: ./run-common-live-journeys-two-at-a-time/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:c1nbf39akyh263a4p1ts593m:backlog:1
                briefing: briefing:c1nbf39akyh263a4p1ts593m:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-09T05:55:31.272689Z"
                decision: approve
                reason: Captain directed dispatch of the narrowed Claude-only concurrency task.
              application:
                target-stage: ideation
                state: consumed
---

Reduce pull-request wait time by running at most two independent common live journeys concurrently on Claude.

## Problem

The common live suite has 16 separate top-level journey entry points, but Go runs them sequentially because none calls `t.Parallel()`. Sonnet currently runs 12 common journeys and then three substrate proofs. A recent exact run spent about 26 minutes in the Sonnet common step and five more minutes in Sonnet substrate proofs.

The entry-point cleanup already gives each journey a separate workflow root, artifact directory, and Claude config directory.

## Value

Developers receive the same named live evidence sooner. One slow runtime lane no longer holds a pull request for avoidable sequential time.

## Proposed approach

Resolve the runtime target and apply target-scoped TODOs before constructing a live driver. For Claude only, call `t.Parallel()` before runner setup. Keep the CI command on the canonical `^TestLiveCommon` selector and add `-parallel 2` to the Claude command.

Keep the existing per-journey artifacts, metrics filenames, and one-process-per-journey watchdog. Do not parallelize Codex, Pi, or the Claude substrate proofs in this task.

Use Go's built-in parallel-test semaphore. Do not add a scheduler, worker pool, retry loop, shard registry, or new test harness. With `-failfast`, at most one additional in-flight journey can finish after the first failure.

## Out of scope

- More than two concurrent live processes.
- Parallel Codex, Pi, or Claude substrate proofs. Reconsider Codex only after Claude produces measured evidence.
- CI matrix sharding or additional environment approvals.
- Changes to desired journeys, TODO ownership, fixtures, assertions, models, effort, or runtime timeouts.
- A committed simulator for concurrency behavior.

## Acceptance criteria

**AC-1 (VALUE) - Sonnet common-journey wall time falls without losing named evidence.**
Verified by: one exact-candidate Sonnet run executes the unchanged selected journey set with a maximum concurrency of two and completes the common step in 20 minutes or less. Every non-TODO journey still emits its normal result and metrics artifact. Removing `-parallel 2` or serializing the tests makes the measured step exceed the target against the recorded baseline.

**AC-2 - Failure cost and API pressure remain bounded.**
Verified by: process start and finish evidence shows no more than two Claude journey processes active at once. A controlled failing focused run starts no new journey after the failure and has at most one already-running sibling. Increasing the cap above two fails this criterion.

**AC-3 - The canonical registry and CI selection stay unchanged.**
Verified by: `TestRuntimeLiveRegistryReconciliation` passes with the same 16 journey IDs, fixtures, TODO bindings, and exact `^TestLiveCommon` selectors. The workflow adds no live job or environment approval.

## Test plan

Keep implementation small: target resolution before driver construction, conditional `t.Parallel()`, and `-parallel 2` in the existing Claude CI command. Use existing offline, race, formatting, and registry checks. Prove the user-visible result with one local subscription-backed two-journey Claude run before the full exact-candidate run. Record wall time and maximum observed concurrency. Do not commit a parallel-runtime simulator.
