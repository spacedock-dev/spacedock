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

The common live suite has 16 top-level journey entry points. Go runs them in sequence because none calls `t.Parallel()`.

Sonnet currently runs 12 non-TODO common journeys and three substrate proofs. Run `31295813569` spent approximately 26 minutes on the common step. The substrate proofs used five more minutes.

The existing runner gives each journey a separate workflow root, artifact directory, and Claude configuration directory. This isolation makes bounded parallel runs possible.

## Value

Developers receive the same named live evidence sooner. The target is 20 minutes or less for the Sonnet common step.

The independent baseline is pull-request run `31295813569`, which used approximately 26 minutes. This task follows the timing review for PR `#643`.

## Proposed approach

Resolve the runtime target before the live driver is constructed. Apply target-scoped TODOs before the parallel pause.

For Claude common journeys only, call `t.Parallel()` before runner setup. Add `-parallel 2` to the Claude CI and local commands.

Keep the current artifacts, metrics files, workflow roots, and one-process watchdogs. Codex, Pi, and Claude substrate proofs stay in sequence.

Use the Go parallel-test semaphore. A scheduler is unnecessary because Go already limits active parallel tests with `-parallel`.

The simpler alternative is CI matrix sharding. It needs more jobs and approvals, so it does not serve the value at lower cost.

The CI command keeps the canonical `^TestLiveCommon` selector and `-failfast`. After a failure, Go starts no queued journey. One active sibling can finish.

This mechanism serves AC-1 and AC-2. Existing isolation serves AC-4 without a new harness.

### Expected surface

The expected change has five files and approximately 35 to 70 inserted lines. The tolerance is one additional test file and 20 lines.

| File | Estimated insertions | Purpose |
|---|---:|---|
| `internal/ensigncycle/shared_live_runner_test.go` | 10-20 | Resolve the target first and parallelize Claude common journeys. |
| `.github/workflows/runtime-live-e2e.yml` | 2-5 | Add the Claude concurrency cap and update the step comment. |
| `internal/contractlint/live_registry_reconciliation_test.go` | 8-15 | Require the Claude cap and preserve the other selectors. |
| `internal/release/journey_workflow_test.go` | 4-10 | Update the executable workflow guard. |
| `docs/runtime-live-ci.md` | 10-20 | Document Claude concurrency and the local command. |

Command grammar changes only for the Claude common-suite command. It adds `-parallel 2`.

Runtime behavior changes only for Claude common journeys. At most two journey processes run at the same time.

Stored formats and authority do not change. Journey IDs, TODOs, fixtures, models, effort, timeouts, and artifact names do not change.

## Out of scope

- More than two concurrent live processes.
- Parallel Codex, Pi, or Claude substrate proofs. Reconsider Codex only after Claude produces measured evidence.
- CI matrix sharding or additional environment approvals.
- Changes to desired journeys, TODO ownership, fixtures, assertions, models, effort, or runtime timeouts.
- A concurrency simulator, including a committed simulator.
- A scheduler, worker pool, retry loop, or shard registry.

## Acceptance criteria

**AC-1 (VALUE) - Sonnet common-journey wall time falls without losing named evidence.**
Verified by: one exact-candidate Sonnet run executes the unchanged selected set with a concurrency cap of two. The common step completes in 20 minutes or less. Every non-TODO journey emits its normal result and metrics artifact. Compare this duration with the 26-minute baseline from run `31295813569`.

**AC-2 - Failure cost and API pressure remain bounded.**
Verified by: process start and finish evidence shows no more than two Claude journey processes at one time. A controlled failing run starts no queued journey after the first failure. At most one active sibling finishes.

**AC-3 - The canonical registry and CI selection stay unchanged.**
Verified by: `TestRuntimeLiveRegistryReconciliation` passes with the same 16 journey IDs, fixtures, TODO bindings, and `^TestLiveCommon` selectors. The workflow adds no job or approval.

**AC-4 - Concurrent Claude journeys keep separate state and evidence.**
Verified by: a focused two-journey run produces two workflow roots, two configuration directories, and two artifact directories. Each directory contains only its journey data.

## Test plan

1. Add focused tests for target resolution and the Claude-only command change. These tests serve AC-2 and AC-3.
2. Run `TestRuntimeLiveRegistryReconciliation`. This test fails if selectors, IDs, fixtures, or TODO bindings change.
3. Run `go test ./...` and `go test ./... -race`. These commands make sure that shared test state stays race-free.
4. Run `gofmt -w ./cmd ./internal`. Then make sure that this command changes no unrelated file.
5. Run two focused Claude journeys with `-parallel 2`. Record wall time, active process count, roots, configuration directories, and artifacts.
6. Run a controlled failing Claude pair with `-failfast -parallel 2`. Record the active sibling and make sure that no queued journey starts.
7. Run the exact Sonnet candidate command. Compare its common-step duration with run `31295813569`.

The local spike on 2026-08-08 used real Claude Code `2.1.220` and `claude-sonnet-5`. It started `filing` and `zero-discovery` together.

The spike reached two active journeys in 12 seconds. It created separate roots, configuration directories, and artifact directories.

Both launches then failed before first-officer work. The stored OAuth subscription token was expired, and no API key was available.

This failure proves the launch cap and isolation paths. It does not prove successful journey behavior or the wall-time target.

### Documentation change

Before: `Each command uses the same 16 exported sequential tests and stops at the first non-TODO failure.`

After: `The Claude command runs at most two common journeys at one time. The Codex and Pi commands run the same journeys in sequence.`

Before: `SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 90m -run '^TestLiveCommon' -failfast ./internal/ensigncycle -v`

After: `SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 90m -run '^TestLiveCommon' -failfast -parallel 2 ./internal/ensigncycle -v`

## Stage Report: ideation

- DONE: Prove that two Claude journeys use separate roots, configuration, and artifacts, with no more than two live processes.
  The live spike created distinct `filing` and `zero-discovery` paths and observed a maximum of two journey processes.
- FAILED: Run one local subscription-backed two-journey exercise, and record the wall time, concurrency, and failure behavior.
  The 12-second run reached concurrency two, but both real launches failed because the stored OAuth token was expired.
- DONE: Write a lean Simple English plan with exact files, LOC estimates, value links, and no Codex, Pi, substrate, scheduler, or simulator scope.
  The plan names five files, a 35-to-70-line estimate, the PR `#643` timing value, and explicit scope limits.

### Summary

The plan uses Go's existing semaphore and changes only Claude common journeys. A real two-journey spike proved isolation and the cap, but expired subscription authentication blocked successful behavior evidence.
