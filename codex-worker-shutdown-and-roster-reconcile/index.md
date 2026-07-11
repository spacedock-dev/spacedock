---
title: Reclaim Codex worker slots with real shutdown and roster reconciliation
status: backlog
source: "Captain question 2026-07-11 while 88t context degradation coincided with the live four-thread collaboration ceiling."
started:
completed:
verdict:
score: 0.9
worktree:
issue:
id: jvw9c1z08b2mw4k3myvae9kw
---

## Problem

The current Codex adapter binds roster reads through `list_agents` but leaves `«worker.shutdown»` unresolved. The live collaboration surface exposes `interrupt_agent`, which stops a turn but does not close or archive its thread, so it must not be treated as shutdown.

This session's collaboration ceiling is four active agents including the first officer. A standing communication officer therefore leaves two task-worker slots. The 88t implementation and validation pair filled those slots while the implementation worker reached about 87% of its 353,400-token context window. Thread pressure did not consume that context, but it reduced the first officer's ability to fresh-dispatch a replacement when budget-aware reuse should have failed.

Codex configuration documents `agents.max_threads` with a default of six when unset, yet the embedded collaboration surface in this session enforces four. Capacity must therefore be discovered from the live surface rather than assumed from user config or client version.

## Proposed direction

Keep shutdown/reconcile separate from the context-budget probe:

1. Spike the actual Codex/app-server operation, if any, that permanently closes or archives a completed child thread and demonstrably frees an active collaboration slot.
2. Bind `«worker.shutdown»` only to that proven operation. Explicitly reject `interrupt_agent` as a substitute.
3. Extend roster reconciliation to classify workers as active, completed-but-required-for-feedback, terminal, superseded, or stale.
4. Reclaim only terminal, superseded, and otherwise no-longer-routable cohorts. Preserve implementation/reviewer workers required by feedback routing and completed-but-addressable reuse.
5. Reconcile at bounded lifecycle points: before fresh dispatch when capacity is exhausted, after superseding a worker, and after terminal merge/cleanup. Do not poll or eagerly close useful workers after every completion.
6. Report discovered live capacity and unreclaimable occupants when no safe shutdown binding exists.
7. Evaluate lazy standing-teammate injection separately: it can recover one slot, but must not masquerade as worker lifecycle correctness.

The sibling task `codex-context-budget-probe` decides when reuse is unsafe. This task makes the fresh-dispatch path available by reclaiming slots. Neither substitutes for the other.

## Acceptance criteria

- **AC-1:** A live Codex probe proves a shutdown/close operation permanently removes an unneeded child from the active roster and permits one additional spawn at the same live thread cap.
- **AC-2:** `interrupt_agent` is never bound as `«worker.shutdown»`, and tests distinguish interrupted, completed-addressable, and actually closed workers.
- **AC-3:** Reconciliation preserves any implementation or validation worker still required by an entity's feedback/re-review contract.
- **AC-4:** Terminal, superseded, and stale cohorts are closed at their bounded lifecycle points, and repeated reconciliation is idempotent.
- **AC-5:** When capacity is exhausted, the dispatcher either reclaims a provably safe slot and fresh-dispatches or reports the exact occupants and missing shutdown capability; it never silently degrades to unsafe reuse.
- **AC-6 (VALUE):** At a four-thread live cap with a standing teammate and a completed disposable worker, reconciliation increases available task-worker capacity from zero to one without losing any routable feedback worker or durable stage output.
- **AC-7:** Runtime capacity is derived from live spawn/roster behavior rather than assuming the documented default or a configured `agents.max_threads`.
- **AC-8:** The design records whether lazy standing-teammate injection remains necessary after safe shutdown exists, without coupling that policy change to the lifecycle implementation.

## Test plan

First run a throwaway app-server/collaboration spike at a known low thread cap: spawn to capacity, complete one disposable child, invoke each candidate close/archive operation, reconcile, and prove a replacement spawn succeeds. Also prove that interruption alone does not free the slot.

Then add adapter/state-machine tests covering active, completed-addressable, feedback-required, terminal, superseded, stale, interrupted, absent-shutdown, idempotent cleanup, and exhausted-capacity diagnostics. Finish with a live workflow replay at cap that preserves the feedback reviewer, retires only the disposable cohort, and fresh-dispatches the over-budget implementation replacement.
