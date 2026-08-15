---
title: Repair Sonnet keep-moving and selected-bare live flakes
status: backlog
score: "0.95"
source: "Exact-main Runtime Live E2E run 31894934349 attempts 1 and 2"
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
id: 060xp69y61yhrww23g3wvwqy
gates:
  version: 1
  records: []
started:
worktree:
mod-block:
pr:
---
## Problem

Exact main `c6a336a33f866536d2b1ea9ef313bbbc6f244dd1` produced opposite outcomes across two Sonnet attempts:

- `TestLiveCommonKeepMovingPosture` passed attempt 1 and failed attempt 2 with `keep-moving-violation` after the First Officer bypassed the formal revise/feedback route.
- `TestLiveBreakGlassShimRecovery/selected-bare` failed attempt 1 because the Agent call did not preserve selected bare mode, then passed attempt 2.

Both behaviors are flaky on unchanged bytes. The live registry must retain active ownership until each journey is repaired and passes without an XFAIL binding.

## Acceptance criteria

- **AC-1:** Sonnet keep-moving behavior uses the formal revise/feedback route and passes repeatedly without an XFAIL binding.
- **AC-2:** Sonnet selected-bare break-glass recovery preserves blocking bare mode and passes repeatedly without owned expected-failure handling.
- **AC-3:** Each temporary expected-failure marker names this active owner and is removed only after bound failure/XPASS and unchanged-byte unbound PASS evidence.
- **AC-4:** Registry reconciliation, active-owner join, full, race, and exact Sonnet live checks pass.

## Evidence

- Run `31894934349`, attempt 1: common keep-moving PASS; selected-bare FAIL.
- Run `31894934349`, attempt 2: keep-moving FAIL; selected-bare PASS.
