---
title: Replace reviewer-reuse inference with structured evidence
status: backlog
source: "Test-infrastructure audit 2026-07-14."
started:
completed:
verdict:
score:
worktree:
issue:
milestone: 0.26.0
id: 6vn56z3423xk3ej3wvk54z6r
---

## Problem

The reviewer-reuse oracle infers identity across transcript dialects, command strings, wait counts, free-form narration, and durable reports. Its own code admits that some Codex output lacks enough metadata to prove a distinct or reused reviewer, yet fallback paths still pass. More parsing cannot create the missing fact.

## Required outcome

Claim reviewer reuse only from structured native spawn/follow-up handle correlation. When a host does not expose that surface, report identity as unsupported; durable reports may prove that re-review occurred, but not who performed it.

## Acceptance criteria

- **AC-1 (VALUE):** A reuse claim passes only when structured runtime evidence ties both validation cycles to the same reviewer handle.
- **AC-2:** An unavailable or ambiguous identity surface produces an explicit unsupported result and never passes from commands, waits, reports, or narration alone.
- **AC-3:** Fixtures reject wrong-reviewer and ambiguous-transcript cases and remove self-source/prose assertions.
- **AC-4:** Existing supported-host routing remains behaviorally covered without adding another transcript dialect parser or identity state machine.

## Mechanism/value trace

Simplest route: correlate native structured handles. The simpler alternative cannot prove identity only when the runtime omits it; that case must remain honestly unproven rather than reconstructed.
