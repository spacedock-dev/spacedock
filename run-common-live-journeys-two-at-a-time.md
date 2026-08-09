---
title: Run Claude common live journeys two at a time
status: validation
source: "Captain request after PR #643 CI timing review, 2026-08-09. Run 31295813569 measured Claude common plus substrate at about 31 minutes. Codex was explicitly deferred until later."
started: 2026-08-09T05:56:00Z
completed:
verdict:
score: 0.75
sprint:
sprint-readiness: ready
group: live-ci-performance
worktree: .worktrees/spacedock-ensign-run-common-live-journeys-two-at-a-time
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
        - id: gate:c1nbf39akyh263a4p1ts593m:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:c1nbf39akyh263a4p1ts593m-ideation-1
              briefing:
                id: briefing:c1nbf39akyh263a4p1ts593m:ideation:attempt-1:revision-1
                digest: sha256:a652a2335c0572d93e39271d332e1e0c9e25d9c157318e60deccafb69265f7f1
                request-digest: sha256:201d5101bf0d12f55269979cc995b7944f3b4474325e880b768ccb07864f7550
                room-ref: ./run-common-live-journeys-two-at-a-time/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:c1nbf39akyh263a4p1ts593m:ideation:1
                briefing: briefing:c1nbf39akyh263a4p1ts593m:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-09T13:56:08.747199Z"
                decision: approve
                reason: Captain approved proceeding after removal of tests for test infrastructure.
              application:
                target-stage: implementation
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

The expected change has five files. It has 15 to 30 gross insertions and a net increase of 3 to 12 lines.

The tolerance is five files, 40 gross insertions, and a net increase of 17 lines. A larger change returns to ideation.

| File | Gross insertions | Net change | Purpose |
|---|---:|---:|---|
| `internal/ensigncycle/shared_live_runner_test.go` | 8-15 | +5 to +10 | Resolve the target first and parallelize Claude common journeys. |
| `.github/workflows/runtime-live-e2e.yml` | 2-4 | 0 to +2 | Add the Claude cap and replace the step comment. |
| `internal/contractlint/live_registry_reconciliation_test.go` | 1-3 | 0 | Replace the existing exact Claude command expectation. |
| `internal/release/journey_workflow_test.go` | 1-3 | 0 | Replace the existing exact Claude command expectation. |
| `docs/runtime-live-ci.md` | 3-8 | -2 to +2 | Replace the sequence text and local Claude command. |

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
- New tests for test infrastructure.
- A new contract-lint assertion or release workflow assertion.

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

1. Replace the existing exact Claude command expectations. Do not add a contract-lint assertion or release workflow assertion.
2. Run `TestRuntimeLiveRegistryReconciliation`. This test fails if selectors, IDs, fixtures, or TODO bindings change.
3. Run `go test ./...` and `go test ./... -race`. These commands make sure that shared test state stays race-free.
4. Run `gofmt -w ./cmd ./internal`. Then make sure that this command changes no unrelated file.
5. Run two focused Claude journeys with `-parallel 2`. Record wall time, active process count, roots, configuration directories, and artifacts.
6. Run a controlled failing Claude pair with `-failfast -parallel 2`. Record the active sibling and make sure that no queued journey starts.
7. Run the exact Sonnet candidate command. Compare its common-step duration with run `31295813569`.

Live Claude runs prove concurrency, isolation, failure cost, and value. Offline tests do not simulate these properties.

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
- DONE: Run one local subscription-backed two-journey exercise, and record the wall time, concurrency, and failure behavior.
  The 12-second run reached concurrency two, but both real launches failed because the stored OAuth token was expired.
- DONE: Write a lean Simple English plan with exact files, LOC estimates, value links, and no Codex, Pi, substrate, scheduler, or simulator scope.
  The plan names five files, 15-to-30 gross insertions, 3-to-12 net lines, the PR `#643` timing value, and explicit scope limits.

### Summary

The lean plan uses Go's existing semaphore and changes only Claude common journeys. It adds no infrastructure test or concurrency simulator.

A real two-journey spike proved isolation and the cap. Expired subscription authentication blocked successful behavior evidence.

## Stage Report: implementation

- DONE: Make Claude common journeys run with Go concurrency two, while Codex, Pi, and substrate proofs stay sequential.
  Commit `25fe7f42e` lazily constructs live drivers after Claude calls `t.Parallel()`; only the Claude common command adds `-parallel 2`.
- DONE: Keep the change within 40 gross insertions and 17 net lines, with no new infrastructure tests or simulator.
  The five-file commit has 22 gross insertions and +5 net lines; it only replaces existing exact command expectations.
- DONE: Run the existing offline checks and the smallest available live Claude exercise, and record exact evidence and authentication limits.
  `go test ./...`, `go test ./... -race`, gofmt, and registry reconciliation passed; a 2.75-second Claude pair entered both tests and made separate artifact/config trees, then expired OAuth stopped both before FO work.

### Summary

Claude common journeys now use Go's built-in test semaphore with a concurrency cap of two. Codex, Pi, and Claude substrate commands remain sequential, and live success/value timing remains authentication-blocked for independent validation.

## Review-finding disposition

### Validation finding 1 - queued Claude journey starts after the first failure

- Reviewer observation: a real three-journey Claude run resumed queued `filing` after `merge-hook-guardrail` failed, while `zero-discovery` was still active.
- Released user and normal workflow: PR CI runs the canonical Claude common suite with `-failfast -parallel 2`; a journey regression is a normal supported trigger.
- Observable harm: another paid journey starts after the first failure, so failure cost and API pressure exceed the promised bound.
- Affected authority: `value-ac[AC-2]` Failure cost and API pressure must remain bounded; no queued journey may start after the first failure.
- Trigger evidence: Claude Code `2.1.220` emitted `CONT merge-hook`, `CONT zero-discovery`, the first failure, then `CONT filing` under the exact candidate flags.
- Defect kind: outcome defect. The observable live behavior contradicts AC-2, rather than merely lacking evidence.
- Release scope: Material. The trigger is any real common-journey failure in the supported PR lane, not an unsupported edge case.
- Proposed ownership and disposition: mechanism/design reset; preventing Go from releasing a queued parallel test needs control outside the approved semaphore-only scope.
- Validation recommendation: REJECTED. Candidate commit `25fe7f42e` stayed unchanged pending First Officer and captain direction.

### Feedback Cycles

- Cycle 1: REJECTED — validation reviewer; surface 5 files at +22/-17 vs estimate 5 files, 15-30 gross insertions, +3 to +12 net; AC unchanged; Captain authorized a design reset to ideation because Go's semaphore violates AC-2.

## Stage Report: validation

- DONE: Verify that the five-file diff adds no infrastructure test, simulator, scheduler, or concurrency outside Claude common journeys.
  Commit `25fe7f42e` changes exactly five existing files (+22/-17); test edits only replace existing exact command expectations, and only Claude common tests call `t.Parallel()`.
- DONE: Reproduce the existing offline, race, formatting, and registry checks against commit 25fe7f42e.
  `go test ./...`, `go test ./... -race`, focused registry/workflow guards, `git diff --check`, and gofmt passed; gofmt left the worktree clean.
- DONE: Inspect or run the smallest real Claude exercise, and separate proven structure from live evidence that must come from PR CI.
  A real Claude pair entered both tests and wrote separate journey artifact/config paths, but expired OAuth stopped both before FO work; successful behavior and Sonnet timing remain unproved PR-CI evidence.
- SKIPPED: AC-1 successful Sonnet common-step timing and complete named evidence.
  Local OAuth is expired, so validation does not claim the <=20-minute target or successful result/metrics artifacts; PR CI must supply both.
- FAILED: AC-2 failure cost and API pressure remain bounded.
  A real three-journey `-failfast -parallel 2` run started queued `filing` after the first failure while one sibling remained active, falsifying the no-new-start promise.
- DONE: AC-3 canonical registry and CI selection stay unchanged.
  `TestRuntimeLiveRegistryReconciliation` and exact command guards passed; the diff retains all 16 IDs and selectors and adds no job or approval.
- DONE: AC-4 concurrent Claude journeys keep separate state and evidence paths.
  The focused pair produced distinct `filing` and `zero-discovery` artifact directories, temp workflow roots, and `.claude/<journey>` configuration paths; authentication prevented successful artifact contents.

### Summary

Offline, race, formatting, registry, scope, and structural isolation checks passed. Validation recommends REJECTED because the real controlled failure proved Go's parallel-test semaphore starts a queued journey after the first failure; the successful Sonnet wall-time proof remains explicitly unclaimed until PR CI.
