---
id: x02375wsg6q61xek7p0t36j2
title: Repair the Pi keep-moving posture journey
status: backlog
source: "Deferred Pi follow-up from the test-behavior-completeness priority recarve, 2026-08-10"
started:
completed:
verdict:
score: 0.8
group: pi-live-followup
worktree:
issue:
pr:
mod-block:
sprint: pi-live-completeness
---

Pi still needs a product repair for the `keep-moving-posture` journey. Sonnet and Codex are complete. This task owns only the deferred Pi result.

## Problem

The Pi target remains executable but does not have a current passing result. Its former owner will complete after the Sonnet and Codex repair lands.

## Visible value

A Pi operator consumes the gate, dispatches the entered stage with durable evidence, and reaches the safe terminal result without a forced status write.

## Out of scope

- Sonnet or Codex behavior.
- Shared XFAIL policy.
- Evidence-only xp6 work.
- A new gate command, stored format, authority source, or CI lane.

## Acceptance criteria

**AC-1 — The exact Pi keep-moving target passes the repaired product journey.**

Verified by: the focused live Pi target exits successfully and retains a passed journey metric with the required dispatch and terminal state.

**AC-2 — The Pi binding stays honest until the repair passes.**

Verified by: the Pi XFAIL names this active task before repair and is removed only after exact XPASS and normal PASS evidence.

## Test plan

Use focused offline gate and terminalization controls first. Use one exact Pi target sequence only when Pi work is authorized. Preserve all Sonnet and Codex behavior.
