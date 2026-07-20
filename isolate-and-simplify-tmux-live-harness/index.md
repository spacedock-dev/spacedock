---
title: Isolate and simplify the real-tmux live harness
status: backlog
source: "Test-infrastructure audit 2026-07-14."
started:
completed:
verdict:
score:
worktree:
issue:
milestone: 0.26.0
id: 3a9acywe2trdxswhjtszh23z
sprint:
group:
sprint-readiness:
---

## Problem

The interactive live harness correctly uses real tmux instead of a custom PTY, but it shares the operator's default tmux server, compensates for inherited server environment, resolves private transcript layout, and may send the captain command three times. Those recovery layers can mask first-input failure or an improper re-park.

## Required outcome

Keep tmux as the real terminal boundary while making the proof isolated and single-attempt: a disposable server/socket, one supported readiness observation, one literal captain command, native runtime state, and durable workflow outcomes.

## Acceptance criteria

- **AC-1 (VALUE):** Interactive Claude journeys prove the first captain command is accepted once and produces the expected native roster and durable workflow state.
- **AC-2:** Every run owns a unique disposable tmux server/socket and verifies cleanup without relying on the operator's server environment.
- **AC-3:** Multi-nudge recovery is removed; a dropped or mishandled first command fails visibly rather than being retried into green.
- **AC-4:** Private transcript/session heuristics are removed or narrowly justified by an explicit architecture review, and the live journeys remain green.

## Mechanism/value trace

Simplest route: disposable tmux, one send-keys, capture/native state, durable state, cleanup. Do not introduce a PTY driver, process-group supervisor, lease, or second terminal lifecycle.
