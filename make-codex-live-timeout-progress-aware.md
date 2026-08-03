---
title: Make Codex live scenario liveness progress-aware
status: backlog
source: "Desired live coverage must not false-red a healthy progressing Codex run, 2026-08-03"
started:
completed:
verdict:
score: 0.85
worktree:
issue:
sprint:
group: lane-reliability
sprint-readiness: absorbed
id: wp4b3zdv05fxxfz2yyhf002s
---

## Problem

Codex shared scenarios use a fixed 15-minute process deadline that does not reset on JSONL activity, worker events, or durable progress. A healthy but slow run can therefore be reported as a timeout while still advancing.

## Acceptance criteria

**AC-1 (VALUE)** A progressing Codex journey remains alive while a stalled journey fails within the documented quiet budget.
Verified by: deterministic stream-watcher tests that advance beyond the old wall-clock boundary under activity and fail after sustained silence, plus one focused live run.

**AC-2** The suite-wide timeout remains a loose runaway backstop and does not become the primary liveness signal.
Verified by: runner configuration and focused timeout controls.

**AC-3** Failure diagnostics retain the last observed progress and artifact locations.
Verified by: a forced-stall test over the runner observation boundary.

## Stage-specific test gates

- Spike the Codex JSONL activity reset path before choosing a timer design.
- Validation runs focused runner tests, a Codex live journey, and full/race suites.

