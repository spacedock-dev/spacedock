---
title: Make live test results tell the truth
status: backlog
source: "Captain recarve of live-test-truth, 2026-08-03. Absorbs 1a and wp as design inputs."
score: 1.0
sprint: live-test-truth
group: truthful-results
sprint-readiness: ready
id: 3d2rqxrgvqky085mn170x3zp
---

## Outcome

A bad first-officer decision produces a red result. A healthy progressing Codex run does not produce a false timeout.

The task owns truthful classification at both edges. Oracle storage and liveness timers are implementation steps, not separate delivery units.

## Design inputs

- `ac2-reanchor-live-scenario-repair` (`1a`): reproduced the false green and proved a durable two-way gate oracle.
- `make-codex-live-timeout-progress-aware` (`wp`): records the current fixed-deadline false-red risk.

Preserve the AC2 spike. Do not repeat it unless a source change invalidates it.

## Acceptance criteria

**AC-1 (VALUE) — The AC2 scenario rejects the wrong durable decision.**
Verified by: revise reaches `rework`, approve reaches `accepted`, and the same oracle rejects approve and unchanged state.

**AC-2 (VALUE) — A progressing Codex run stays alive.**
Verified by: deterministic activity extends the run beyond the old wall-clock limit, while one focused live journey remains green.

**AC-3 — A stalled Codex run fails within the quiet budget.**
Verified by: a forced-stall test reports the last progress event and artifact locations.

**AC-4 — The suite timeout stays a runaway backstop.**
Verified by: quiet-progress controls, suite configuration, full tests, and race tests.

## Ideation requirements

Use the AC2 report as completed design evidence. Spike only the unproved Codex activity-reset path. Produce one plan that shows the two visible truthfulness outcomes.

Use `$simple-english` in pragmatic mode for the complete plan.

