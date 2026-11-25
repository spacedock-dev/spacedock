---
id: mcsp2zneb00z5zw5rh7rh6j3
title: A first officer goes idle while dispatchable work waits
status: backlog
source: "MetaFO friction list 2026-11-24, sent on the word of CL, /tmp/spacedock-frictions-2026-11-24.md, group B item 7"
group: B-paseo-host-backend
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
---

A first officer keeps driving while work is ready, so the captain does not have to nudge it.

## Problem

On the Studfleet VM (Paseo host backend), a first officer ended a turn with work to do and stayed idle until a message came. The cause is not known yet: a missing wake-up when a worker stops, or the first officer's own choice. MetaFO asked and will send the answer.

## Proposed approach

{Ideation fills this in, after the cause is known.}

Related backlog: `completion-delivery-watchdog`, `fo-external-wait-watcher`, `helper-dispatch-completion-signal`.

## Risk evidence

Backlog: one observation, cause unknown. The backlog gate decides whether to wait for MetaFO's answer before ideation.

## Out of scope

Paseo's own notification routing.

## Expected surface and tolerance

Estimate net LOC change: {ideation}, across {M} files.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - When a worker stops and another entity is dispatchable, the first officer dispatches it without an outside message.**
Verified by: {ideation names the fixture or live journey once the cause is known}.

## Test plan

{Ideation fills this in.}

### Feedback Cycles
