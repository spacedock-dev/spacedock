---
title: Pi stage dispatches should force fresh subagent context
status: backlog
source: captain (2026-06-04) — FO mistakenly dispatched a Spacedock implementation worker with pi-subagents context=fork; stage workers should be fresh and independent
score: "0.29"
started:
completed:
verdict:
worktree:
issue:
id: d2w8z614c0q1yssmyr33a38y
---

Spacedock's Pi first-officer runtime should make stage dispatch context explicit. During `launcher-binary-path-passthrough`, the FO dispatched an implementation worker with `context: "fork"` because the builtin worker agent defaults to forked context. That preserved parent context but conflicts with the intended Spacedock stage model: implementation and validation workers should receive the assignment/context from the dispatch prompt and entity state, not inherit the FO session transcript by default.

## Problem

The current operating contract forbids `pi-subagents` acceptance contracts for Spacedock stages and recommends fresh redispatch for Pi feedback loops, but it does not yet guard the `context` parameter used when the FO calls the subagent harness. A forked implementation worker can accidentally inherit prior FO reasoning, validation expectations, or captain-side discussion that should not be part of an independent stage assignment.

## Desired outcome

Make Spacedock stage dispatches through Pi/subagents explicitly fresh-context by default, and make violations visible in tests or runtime checks.

## Acceptance criteria

**AC-1 - Pi/Spacedock stage dispatch instructions require fresh subagent context.**
Verified by: an instruction invariant test or runtime contract test that fails if the Pi first-officer stage dispatch guidance omits `context: "fresh"` or allows relying on the worker agent's default context.

**AC-2 - Spacedock stage dispatches still avoid `pi-subagents` acceptance contracts.**
Verified by: existing or extended integration tests over Pi first-officer runtime docs that require `subagent(...)` stage dispatches to put acceptance requirements in the task prompt and forbid `subagent(... acceptance: ...)`.

**AC-3 - Feedback/retry dispatch remains a fresh assignment cycle unless explicitly marked as manual/debug resume.**
Verified by: tests or docs invariants that distinguish normal stage dispatch from opt-in resume/debug tooling and require durable entity-stage evidence for every cycle.

## Notes

This is separate from `launcher-binary-path-passthrough`. That task should be redispatched or continued under the corrected policy, but the policy/guard itself belongs in this separate entity.
