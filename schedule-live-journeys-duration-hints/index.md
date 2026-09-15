---
title: Schedule live journeys with committed duration hints
status: ideation
source: Captain request 2026-09-15; CI run 34996910090
started: 2026-09-15T18:21:56Z
completed:
verdict:
score: 0.8
worktree:
issue:
pr:
mod-block:
id: pytyzge5v85mcy8c02t1khqv
---

Reduce live CI completion time with committed per-runtime duration hints and one bounded queue for common journeys and Claude substrate proofs.

## Problem

CI run 34996910090 tested stack tip PR #798 at 380bd7bf34b3dd1de07a125d75035c593f5ff930. Its common journeys use three slots without duration ordering. Claude then runs substrate proofs serially. Observed common test times: Codex 945.863 seconds; Claude 1334.022 seconds. Claude substrate adds 520.586 seconds.

## Proposed approach

Keep rounded duration hints in the repository, keyed by runtime and scheduled test. Start longer tests first with deterministic ties. Put Claude substrate proofs in the same three-slot queue as common journeys. Preserve existing test assertions, failure propagation, coverage and artifact retention. Metrics inform occasional manual hint updates when sustained changes affect ordering; scheduling reads no remote history or external state. Do not introduce automatic hint rewriting or timing failure gates.

Keep this as a separate task and eventual stack PR above the current tip or forthcoming Claude fix. Do not modify already-validated stack branches. No push or CI until the stack is ready and individual changes have targeted local evidence.

## Risk evidence

Downloaded Go test JSON and journey records are under /tmp/ci-schedule-34996910090. Calculation: /tmp/analyze-ci-schedule.py and /tmp/ci-schedule-analysis.json. Longest-first predicts common Codex 817.88 seconds and Claude 1186.59 seconds. Pooling Claude common and three substrate tests predicts 1363.72 seconds, versus 1854.608 seconds sequentially. These are retrospective fixed-duration estimates, not live guarantees.

Substrate tests use separate temporary workflows but CI shares CLAUDE_CONFIG_DIR. Prove isolation before overlapping them. Go t.Parallel registration order alone does not guarantee a priority queue. Spike the smallest supported scheduling mechanism. Journey metrics cover bare dispatch and both break-glass modes; merged dispatch currently lacks a journey record. Use whole-test Go elapsed durations for initial hints, aggregating sequential variants correctly.

## Out of scope

Provider behavior fixes, increased concurrency, exact optimal scheduling, remote metrics dependencies, automatic hint updates, changed acceptance assertions, and automatic CI runs on intermediate stack PRs.

## Expected surface and tolerance

Ideation must establish the smallest concrete file and net-line estimate with explicit tolerance before implementation. Likely owners: internal/ensigncycle scheduling and substrate tests, .github/workflows/runtime-live-e2e.yml, and docs/runtime-live-ci.md. Reuse existing registries and runners.

## Acceptance criteria

**AC-1 — Committed duration hints reduce scheduling idle time without external state.**
Verified by: deterministic scheduler exercises using independent fixed durations, including the retained run; longer-ready work starts first, ties are stable, and all selected tests run once. An inverted ordering must fail the exercise.

**AC-2 — Claude substrate and common journeys share at most three running tests with isolated state.**
Verified by: targeted local overlap of actual substrate and common journeys, plus a deterministic concurrency exercise. Preserve workflow, config, artifact and process isolation; the fourth simultaneous test must fail the bounded-concurrency exercise.

**AC-3 — CI keeps complete evidence and propagates failures.**
Verified by: targeted runner failure and cancellation exercises plus existing registry checks; no omitted test, false pass, or lost detail/metric artifact. A missing scheduled test or swallowed failure must fail verification.

## Test plan

Before implementation, exercise the actual proposed scheduler with short deterministic jobs, then perform the smallest local live overlap that can expose shared-state collisions. Run focused tests before changes, required normal/race suites and formatting after changes, and a detached adversarial audit for CI machinery. Reserve full Claude/Codex stack-tip CI until individual changes and the final stack are ready.

### Feedback Cycles

