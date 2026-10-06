---
title: Claude live break-glass journey fails on uncommitted entity state
sprint: test-behavior-completeness
source: "Run 37421110550 claude-live: FAIL TestLiveBreakGlassShimRecovery/selected-team, dispatch_recovery_live_test.go:119: entity has uncommitted changes: M widget-task.md."
id: y5hde203rm3ga4s6mm30yv0r
---

## Problem

`TestLiveBreakGlassShimRecovery` fails because a worker left `widget-task.md` modified without
committing. The assertion is about durable state, and the run ends with uncommitted changes.

## Value

Break-glass dispatch recovery is a safety path; its journey cannot pass while state is left behind.

## Acceptance criteria

**AC-1** Determine whether the uncommitted state is a product fault, a fixture gap, or a worker-behaviour difference between hosts.
**AC-2** The journey passes, or the expectation is corrected with the reason recorded.

## Verification

Re-run the focused live target and read the worktree state at failure.
