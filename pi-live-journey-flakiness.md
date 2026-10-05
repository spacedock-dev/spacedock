---
title: Measure and stabilise Pi common-journey flakiness
sprint: pi-live-completeness
source: "Cross-run evidence: owned-conflict-owner-handoff XFAIL then XPASS; keep-moving-posture XPASS then FAIL, its binding removed on that single XPASS per its own recorded condition; ac-value-reanchor pass then FAIL with no binding."
id: psvqjf0w8xh2txp9604gsvmz
status: backlog
---

## Problem

Several Pi common journeys return different verdicts across runs with no code change, so bindings have been
adjusted against a signal that is not stable.

## Value

A binding is only meaningful if pass and fail are evidence. Today neither is.

## Acceptance criteria

**AC-1** Each affected journey's verdict distribution is measured across at least five runs.
**AC-2** Every unstable journey carries an owner and a binding, so the lane is reliably green.
**AC-3** A binding removal requires five consecutive passes, recorded on the owning task.

## Verification

Read the register from consecutive lane runs and tabulate per-journey outcomes.
