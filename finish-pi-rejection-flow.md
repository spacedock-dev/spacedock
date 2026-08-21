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
gates:
    version: 1
    records:
        - id: gate:p17swb3375rt525fn7f8xt7e:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:p17swb3375rt525fn7f8xt7e-backlog-1
              briefing:
                id: briefing:p17swb3375rt525fn7f8xt7e:backlog:attempt-1:revision-1
                digest: sha256:04983d0e5394ec4c8457116492684ab31a3f0807fcc738d92de84662b147d4b9
                request-digest: sha256:923662fa5e0b785a0b9c57ef58cbaa884166aa1ff9756dbb3fee1e7f28eaf6b6
                room-ref: ./finish-pi-rejection-flow/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:p17swb3375rt525fn7f8xt7e:backlog:1
                briefing: briefing:p17swb3375rt525fn7f8xt7e:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-21T08:16:33.370156173Z"
                decision: approve
                reason: Captain conn granted for pi-related fixes including gates (2026-08-21 chat). Seed clearly scopes the deferred Pi rejection-flow XFAIL (registered owner); AC-1 is the focused live Pi target completing normally. Advance to ideation.
              application:
                target-stage: ideation
                state: pending
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
