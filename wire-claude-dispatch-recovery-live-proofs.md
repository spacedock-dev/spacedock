---
title: Select the required Claude dispatch-recovery live proofs
status: backlog
source: "Desired runtime-specific registry includes bare and break-glass Claude behavior, 2026-08-03"
started:
completed:
verdict:
score: 0.8
worktree:
issue:
sprint:
group: runtime-specific
sprint-readiness: absorbed
id: v5w042j910mrpmffhzqrajn4
---

## Problem

TestLiveBareReachable and TestLiveBreakGlassShimRecovery prove supported Claude-specific dispatch boundaries but no CI live selector invokes them. They are release assertions, not manual experiments.

## Acceptance criteria

**AC-1 (VALUE)** The Claude live lane fails when supported bare dispatch or break-glass recovery regresses.
Verified by: both focused tests running through the current Claude front door in CI and by a negative control for each defining assertion.

**AC-2** Both proofs continue to share writeDispatchRecoveryWorkflow, while only the break-glass case adds the failing dispatch-build shim.
Verified by: fixture bindings and focused tests.

**AC-3** The added lane cost and timeout budget are measured and remain within the ideation-approved live budget.
Verified by: recorded focused durations and the complete Claude lane result.

## Stage-specific test gates

- Ideation decides whether these join the existing Claude runtime-specific step or require a separate independently reported step.
- Validation runs both tests serially and then the full/race suites.

