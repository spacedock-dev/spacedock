---
title: Align Claude break-glass recovery with the selected dispatch mode
status: backlog
score: 0.96
source: "PRs #627, #628, #629, and #631 fail TestLiveBreakGlassShimRecovery after PR #626 selected it for required CI. The worker completes through bare blocking dispatch, but the oracle requires a named background worker. History: named recovery template 8e66ead, blanket single-task bare rule ecffced, live selection 4cc0d8."
sprint: test-behavior-completeness
started:
completed:
verdict:
worktree:
pr:
issue:
mod-block:
id: 824ecawn5jttbykcgx82nbf4
---

Stable CI must evaluate the supported Claude break-glass behavior instead of rejecting a successful worker because two contracts disagree.

## Problem and value

Single-task Claude dispatch selects bare blocking mode. Break-glass recovery currently hard-codes a named background worker. The required live oracle enforces the older recovery shape.

This contradiction makes unrelated PRs red even when the recovery worker completes and commits its report. Release decisions then require repeated manual waivers.

## Acceptance criteria

**AC-1 (VALUE) — Required Claude live CI gives one correct result for successful break-glass recovery.**

The live proof passes when the selected dispatch mode completes the worker and durable report. It fails when no worker runs or no report lands.

**AC-2 — Break-glass recovery preserves the dispatch mode selected before helper failure.**

For a single-task bare dispatch, recovery remains bare and blocking. For a team dispatch, recovery uses the valid team worker shape.

**AC-3 — The contract, manual template, fixture, and live oracle describe the same behavior.**

Tests must reject a mode change during recovery and must not require team-only fields for a bare dispatch.

**AC-4 — The correction does not weaken worker-completion evidence.**

The proof still requires the ensign skill, stage definition, worker execution, committed stage report, and bounded stop.

## Test plan

First add focused fixtures for bare and team recovery modes. Then correct the smallest contract and oracle surfaces. Run focused tests, `go test ./...`, `go test ./... -race`, and the real Claude break-glass lanes.
