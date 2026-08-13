---
id: p17swb3375rt525fn7f8xt7e
title: Finish the Pi rejection-flow journey
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

Pi still needs a product repair for the `rejection-flow` journey. Sonnet is complete, and Codex has a separate active owner. This task owns only the deferred Pi result.

## Problem

The Pi target reaches the second validation gate but can time out before it records and completes the expected stop. Its former owner is complete and archived.

## Visible value

A Pi operator sees the rejected round, the corrected candidate, the second passed validation, and one fresh unresolved validation gate without a timeout.

## Out of scope

- Sonnet or Codex behavior.
- The Codex dvd repair.
- Shared XFAIL policy.
- Evidence-only xp6 work.
- A new gate command, stored format, authority source, or CI lane.

## Acceptance criteria

**AC-1 — The exact Pi rejection-flow target completes normally.**

Verified by: the focused live Pi target exits successfully and retains the complete two-validation state plus one fresh unresolved validation gate.

**AC-2 — The Pi binding stays honest until the repair passes.**

Verified by: the Pi XFAIL names this active task before repair and is removed only after exact XPASS and normal PASS evidence.

## Test plan

Use focused offline lifecycle controls first. Use one exact Pi target sequence only when Pi work is authorized. Preserve the Codex dvd binding and all completed Sonnet behavior.
