---
title: Replace Codex wait watchdog with single-run evidence
status: ideation
source: "Test-infrastructure audit 2026-07-14."
started: 2026-07-15T06:16:10Z
completed:
verdict:
score:
worktree:
issue:
milestone: 0.26.0
id: 15sdwn85ekhwkf0jxnh38h8j
---

## Problem

The Codex rejection-flow harness owns a parallel wait protocol: it parses runtime events, tracks wait epochs, fingerprints Markdown and Git state, kills the process, and retries the entire journey once. Archived evidence includes a green result on attempt 2, so a failed first run can disappear behind a false green.

## Required outcome

Make one supported Codex run authoritative. Bound the real process, preserve its artifacts, and grade its exit plus durable workflow state. Do not recreate Codex wait semantics or add another controller, recovery protocol, daemon, lease, or lifecycle layer.

## Acceptance criteria

- **AC-1 (VALUE):** A failed or stalled first rejection-flow run fails the test and remains diagnosable; no whole-scenario automatic retry can turn it green.
- **AC-2:** Evidence comes from the supported process boundary and durable entity/report/Git state, with no test-owned wait-state protocol beyond a simple time bound.
- **AC-3:** A fault-injected failure is caught, and an exact-head real Codex rejection-flow run passes once without retry.
- **AC-4:** The watchdog/retry implementation is removed or materially reduced, and the replacement has no new controller or lifecycle subsystem.

## Mechanism/value trace

Simplest route: one `codex exec`, a quiet/hard process bound, preserved artifacts, and durable-state grading. If that cannot prove the value, stop for architecture review instead of extending the watchdog.
