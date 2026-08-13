---
id: h30c9jrfcf21fdh2qs5z58sd
title: Repair the Pi smallest-sufficient mechanism journey
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

Pi still needs a product repair for the `smallest-sufficient-mechanism` journey. Sonnet and Codex are complete. This task owns only the deferred Pi result.

## Problem

The Pi target remains executable but does not have a current passing result. Its former owner is complete and archived, so the Pi XFAIL needs an active product owner.

## Visible value

A Pi operator completes the smallest-sufficient dispatch journey with the intended ready-only task selection and no extra dispatch.

## Out of scope

- Sonnet or Codex behavior.
- Shared XFAIL policy.
- Evidence-only xp6 work.
- A new runtime, fixture, result format, or CI lane.

## Acceptance criteria

**AC-1 — The exact Pi target passes the repaired product journey.**

Verified by: the focused live Pi `smallest-sufficient-mechanism` target exits successfully and retains a passed journey metric.

**AC-2 — The Pi binding stays honest until the repair passes.**

Verified by: the Pi XFAIL names this active task before repair and is removed only after exact XPASS and normal PASS evidence.

## Test plan

Use focused offline controls first. Use one exact Pi target sequence only when Pi work is authorized. Preserve all Sonnet and Codex bindings and bytes.
