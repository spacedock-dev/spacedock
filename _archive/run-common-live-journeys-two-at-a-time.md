---
title: Run Claude common live journeys two at a time
status: done
source: "Captain request after PR #643 CI timing review, 2026-08-09. Run 31295813569 measured Claude common plus substrate at about 31 minutes. Codex was explicitly deferred until later."
started: 2026-08-09T05:56:00Z
completed: 2026-08-09T15:57:26Z
verdict: passed
score: 0.75
sprint:
sprint-readiness: ready
group: live-ci-performance
worktree: .worktrees/spacedock-ensign-run-common-live-journeys-two-at-a-time
issue:
pr: pr-merge:647
mod-block:
archived: 2026-08-09T15:57:26Z
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
            - id: gate-attempt:c1nbf39akyh263a4p1ts593m-ideation-2
              briefing:
                id: briefing:c1nbf39akyh263a4p1ts593m:ideation:attempt-2:revision-1
                digest: sha256:a0b9c66b28d091b6520619b0c14b711e5d18a6655c1993dcd36a314f261a64cd
                request-digest: sha256:73f62dd5b6e88d7d9af69dd26498b591d44909e7aad050cdd65f1d794b1c5db0
                room-ref: ./run-common-live-journeys-two-at-a-time/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:c1nbf39akyh263a4p1ts593m:ideation:2
                briefing: briefing:c1nbf39akyh263a4p1ts593m:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-09T15:03:55.961385Z"
                decision: approve
                reason: Captain accepts Go's queued-work cost after a Claude journey failure in exchange for faster successful runs, with at most two active Claude journeys.
              application:
                target-stage: implementation
                state: consumed
---

Reduce pull-request wait time by running at most two independent Claude common journeys at one time.

## Problem

The common live suite has 16 top-level journey entry points. Go runs them in sequence because no test calls `t.Parallel()`.

Sonnet currently runs 12 non-TODO common journeys and three substrate proofs. Run `31295813569` spent approximately 26 minutes on the common step. The substrate proofs used five more minutes.

The existing runner gives each journey a separate workflow root, artifact directory, and Claude configuration directory. This isolation permits parallel runs.

Candidate `25fe7f42e` uses Go's parallel-test semaphore. A live failure run showed the cost of this choice.

After one journey failed, Go started a queued journey while another journey remained active. Thus, `-failfast` does not stop all queued parallel work.

## Value

Developers receive the same named live evidence sooner after a successful run. The target is 20 minutes or less for the Sonnet common step.

The independent baseline is pull-request run `31295813569`, which used approximately 26 minutes. This task follows the timing review for PR `#643`.

The task accepts Go's queued-work cost after a failure. The concurrency cap limits active Claude work, but it does not limit total post-failure work.

## Proposed approach

Keep candidate `25fe7f42e` unchanged. Resolve the runtime target before the live driver is constructed. Apply target-scoped TODOs before the parallel pause.

For Claude common journeys only, call `t.Parallel()` before runner setup. Add `-parallel 2` to the Claude CI and local commands.

Keep the current artifacts, metrics files, workflow roots, and one-process watchdogs. Codex, Pi, and Claude substrate proofs stay in sequence.

Use the Go parallel-test semaphore. Go limits active parallel tests with `-parallel`.

Keep `-failfast`, but do not claim that it stops queued Claude journeys. Go can start queued work after the first failure.

The smallest alternative is to abandon concurrency. That choice avoids queued-work cost, but it does not reduce the successful-run duration.

A scheduler can stop queued work after a failure. It is not necessary for the successful-run value, and it exceeds this task's scope.

The Go semaphore serves AC-1 and AC-2. Existing isolation serves AC-4 without a new harness.

### Expected surface

The existing candidate changes five files. It has 22 gross insertions and a net increase of five lines.

One follow-up commit will correct the failure description in `docs/runtime-live-ci.md`. The cumulative expected surface stays at five files.

The cumulative estimate is 23 to 27 gross insertions and a net increase of 5 to 9 lines. The tolerance is 30 insertions and 10 net lines.

| File | Gross insertions | Net change | Purpose |
|---|---:|---:|---|
| `internal/ensigncycle/shared_live_runner_test.go` | 9 | +5 | Resolve the target first and parallelize Claude common journeys. |
| `.github/workflows/runtime-live-e2e.yml` | 3 | 0 | Add the Claude cap and replace the step comment. |
| `internal/contractlint/live_registry_reconciliation_test.go` | 2 | 0 | Replace the existing exact Claude command expectation. |
| `internal/release/journey_workflow_test.go` | 2 | 0 | Replace the existing exact Claude command expectation. |
| `docs/runtime-live-ci.md` | 7-11 | +1 to +4 | Add the command and describe the accepted failure cost. |

Command grammar changes only for the Claude common-suite command. It adds `-parallel 2`.

Runtime behavior changes only for Claude common journeys. At most two journey processes run at one time.

After a failure, Go can start queued Claude journeys. The task does not set a limit on total post-failure journey starts.

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
- A guarantee that `-failfast` stops queued parallel journeys.

## Acceptance criteria

**AC-1 (VALUE) - A successful Sonnet common run finishes sooner without losing named evidence.**
Verified by: one exact-candidate Sonnet run executes the unchanged selected set. The common step completes in 20 minutes or less. Every non-TODO journey emits its normal result and metrics artifact. Compare this duration with the 26-minute baseline from run `31295813569`.

**AC-2 - The active Claude workload has a fixed cap.**
Verified by: process evidence from the controlled three-journey failure shows no more than two Claude journey processes at one time. It also shows the queued start.

**AC-3 - The canonical registry and CI selection stay unchanged.**
Verified by: `TestRuntimeLiveRegistryReconciliation` passes with the same 16 journey IDs, fixtures, TODO bindings, and `^TestLiveCommon` selectors. The workflow adds no job or approval.

**AC-4 - Concurrent Claude journeys keep separate state and evidence.**
Verified by: a focused two-journey run produces two workflow roots, two configuration directories, and two artifact directories. Each directory contains only its journey data.

## Test plan

1. Keep the existing exact Claude command replacements. Do not add a contract-lint assertion or release workflow assertion.
2. Run `TestRuntimeLiveRegistryReconciliation`. This test fails if selectors, IDs, fixtures, or TODO bindings change.
3. Run `go test ./...` and `go test ./... -race`. These commands make sure that shared test state stays race-free.
4. Run `gofmt -w ./cmd ./internal`. Then make sure that this command changes no unrelated file.
5. Run two focused Claude journeys with `-parallel 2`. Record wall time, active process count, roots, configuration directories, and artifacts.
6. Keep the controlled failure evidence. It shows that queued work can start after a failure and that no more than two processes run.
7. Run the exact Sonnet candidate command. Compare its common-step duration with run `31295813569`.

Live Claude runs prove concurrency, isolation, accepted failure behavior, and value. Offline tests do not simulate these properties.

The local spike on 2026-08-08 used real Claude Code `2.1.220` and `claude-sonnet-5`. It started `filing` and `zero-discovery` together.

The spike reached two active journeys in 12 seconds. It created separate roots, configuration directories, and artifact directories.

Both launches then failed before first-officer work. The stored OAuth subscription token was expired, and no API key was available.

This failure proves the launch cap and isolation paths. It does not prove successful journey behavior or the wall-time target.

### Documentation change

Before: `Each command stops at the first non-TODO failure.`

After: `The Claude command keeps the -failfast flag, but Go can start queued parallel journeys after a failure. At most two Claude journeys run at one time. The sequential Codex and Pi commands stop at the first non-TODO failure.`

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
- Cycle 2: Captain authorized acceptance of Go's queued-work cost or abandonment. Ideation recommends the accepted cost and a narrower value promise.

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

## Stage Report: ideation

- DONE: Preserve candidate `25fe7f42e` and use its live failure as design evidence.
  The code candidate is unchanged. The plan records that Go can start queued Claude work after a failure.
- DONE: Choose the smallest honest outcome.
  The plan accepts Go's queued-work cost and keeps the two-process cap. It adds no scheduler, shard system, simulator, or infrastructure test.
- DONE: Rewrite the value, approach, acceptance criteria, expected surface, test plan, and documentation change in Simple English.
  AC-1 measures successful-run duration and named evidence. AC-2 promises only the observed active-process cap.

### Summary

The revised design keeps the five-file candidate and accepts Go's queued-work behavior. A small documentation follow-up will remove the false failure-stop promise.
