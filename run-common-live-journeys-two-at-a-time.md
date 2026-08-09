---
title: Run Claude and Codex common live journeys two at a time
status: backlog
source: "Captain request after PR #643 CI timing review, 2026-08-09. Run 31295813569 measured Codex common at 8.5 minutes and Claude common plus substrate at about 31 minutes."
started:
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
---

Reduce pull-request wait time by running at most two independent common live journeys concurrently on Claude and Codex.

## Problem

The common live suite has 16 separate top-level journey entry points, but Go runs them sequentially because none calls `t.Parallel()`. Sonnet currently runs 12 common journeys and then three substrate proofs. Codex runs nine common journeys. A recent exact run spent about 26 minutes in the Sonnet common step, five more minutes in Sonnet substrate proofs, and 8.5 minutes in the Codex common step.

The entry-point cleanup already gives each journey a separate workflow root and artifact directory. Claude also gives each scenario a separate config directory. Codex gives each runner an isolated `CODEX_HOME`, but its concurrent runners would write the same `_setup` artifact filenames. That collision must be removed before Codex concurrency is enabled.

## Value

Developers receive the same named live evidence sooner. One slow runtime lane no longer holds a pull request for avoidable sequential time.

## Proposed approach

Resolve the runtime target and apply target-scoped TODOs before constructing a live driver. For Claude and Codex only, call `t.Parallel()` before runner setup. Keep the CI commands on the canonical `^TestLiveCommon` selector and add `-parallel 2` to both commands.

Make Codex setup diagnostics use a per-test directory under `_setup`; keep the existing isolated `CODEX_HOME`, per-journey artifacts, metrics filenames, and one-process-per-journey watchdog. Do not parallelize Pi or the Claude substrate proofs in this task.

Use Go's built-in parallel-test semaphore. Do not add a scheduler, worker pool, retry loop, shard registry, or new test harness. With `-failfast`, at most one additional in-flight journey can finish after the first failure.

## Out of scope

- More than two concurrent live processes.
- Parallel Pi or Claude substrate proofs.
- CI matrix sharding or additional environment approvals.
- Changes to desired journeys, TODO ownership, fixtures, assertions, models, effort, or runtime timeouts.
- A committed simulator for concurrency behavior.

## Acceptance criteria

**AC-1 (VALUE) - Sonnet common-journey wall time falls without losing named evidence.**
Verified by: one exact-candidate Sonnet run executes the unchanged selected journey set with a maximum concurrency of two and completes the common step in 20 minutes or less. Every non-TODO journey still emits its normal result and metrics artifact. Removing `-parallel 2` or serializing the tests makes the measured step exceed the target against the recorded baseline.

**AC-2 (VALUE) - Codex receives the same bounded concurrency benefit.**
Verified by: one exact-candidate Codex run executes the unchanged selected journey set with a maximum concurrency of two and completes the common step in 6.5 minutes or less. Each runner uses a distinct `CODEX_HOME`, setup directory, workflow root, and scenario artifact directory. Reusing the shared `_setup` files or forcing `-parallel 1` fails the isolation or timing evidence.

**AC-3 - Failure cost and API pressure remain bounded.**
Verified by: process start and finish evidence shows no more than two Claude or two Codex journey processes active at once. A controlled failing focused run starts no new journey after the failure and has at most one already-running sibling. Increasing the cap above two fails this criterion.

**AC-4 - The canonical registry and CI selection stay unchanged.**
Verified by: `TestRuntimeLiveRegistryReconciliation` passes with the same 16 journey IDs, fixtures, TODO bindings, and exact `^TestLiveCommon` selectors. The workflow adds no live job or environment approval.

## Test plan

Keep implementation small: target resolution before driver construction, conditional `t.Parallel()`, per-test Codex setup artifacts, and `-parallel 2` in the two existing CI commands. Use existing offline, race, formatting, and registry checks. Prove the user-visible result with one local subscription-backed two-journey run for Claude and Codex before full exact-candidate runs. Record wall time and maximum observed concurrency. Do not commit a parallel-runtime simulator.
