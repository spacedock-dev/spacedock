---
title: Make FO event-loop ordering and idle wait explicit
status: backlog
source: "Captain follow-up after the 2026-08-03 durable-decisions execution-gap diagnosis."
started:
completed:
verdict:
score: 0.98
worktree:
issue:
sprint: durable-decisions
group: fo-contract
milestone: 0.27.0
id: ej9kwkvw94w6rh6n5ek7qrbf
---

Make the First Officer event loop mechanically explicit so a dispatch-only empty result cannot hide merge recovery, ready gates, or a required idle/reconcile pass. The task preserves the existing state and runtime boundaries; it makes the ordering observable and testable.

## Problem

The current contract has the required branches, but an FO can run a filtered `status --next` query, see `dispatchable: []`, and report idle without processing `mod-block`/PR recovery, the separate `ready_gates` surface, the one-shot idle hook and reconcile retry, or the unresolved-worker wait rule. This loses the distinction between no dispatchable worker, a pending merge action, a gate awaiting Captain input, and a genuinely idle session.

## Proposed approach

Specify one event-loop trace and bind every runtime adapter to it. Each iteration drains worker messages when supported, processes every `mod-block` and PR action before querying dispatchables, handles `ready_gates` as gate or merge actions rather than worker dispatches, and then runs `status --next`. After the first empty result, run the idle hook exactly once, reconcile the roster, and query `status --next` once more. Install async wait only when the worker set is unresolved and there is no dispatchable entity, presentable gate, mod action, or other state work; a completed, errored, or absent worker set never qualifies. Preserve keep-moving for unrelated ready tasks and do not add a new state field or scheduler command.

## Out of scope

- Changing the commissioned workflow definition or stage taxonomy.
- Adding a new status command, daemon, watchdog, or generic resolver worker.
- Auto-resolving PR conflicts, force-pushing, or changing merge/approval authority.
- Rewriting task-specific findings or changing the Codex, Claude, or Pi transport APIs.

## Acceptance criteria

**AC-1 (VALUE) — A ready action cannot be hidden by a dispatch-only empty result.**
Verified by: a fixture-backed event-loop trace with one `mod-block`/PR row, one `awaiting-captain` or `approved-awaiting-merge` gate, and one genuinely dispatchable task; the recorded order shows mod/PR and gate handling before `status --next`, and the unrelated ready task still dispatches.

**AC-2 — Empty dispatch results trigger one idle/reconcile retry and a truthful stop.**
Verified by: a command-log fixture that returns an empty first `status --next`, releases a task from the idle hook or reconcile, and asserts one idle hook, one reconcile, one second `status --next`, and dispatch of the newly unblocked task; the unchanged-empty variant asserts the explicit no-dispatchable stop reason.

**AC-3 — Async wait is reserved for unresolved workers.**
Verified by: Codex runtime fixtures covering active unresolved, completed, errored, and absent worker sets; only the active unresolved case with no other work may emit `wait_agent(timeout_ms: 300000)`, and no other case emits a wait.

**AC-4 — Gate readiness is not conflated with worker dispatchability.**
Verified by: boot/status fixtures that expose `awaiting-captain`, `approved-awaiting-advance`, and `approved-awaiting-merge` rows, and assert presentation, consume/advance, and merge routing respectively without a worker spawn or a false idle report.

**AC-5 — One blocked task does not park unrelated work.**
Verified by: a multi-entity fixture with one conflict/mod-blocked task and two independent ready tasks; the trace shows the blocked task held with its evidence while both independent tasks reach their declared dispatch or gate action.

## Expected surface and semantic boundaries

Expected change is limited to the First Officer dispatch/event-loop contract and its behavior fixtures: `skills/first-officer/references/first-officer-shared-core.md`, `skills/first-officer/references/fo-dispatch-core.md`, the host runtime adapter notes needed to bind the wait predicate, and focused event-loop/recorded-workflow tests. Estimate 4–8 files and +120/-20 lines, with a 2x tolerance. Allowed semantic change: ordering and stop-reason observability only. Stored state, command grammar, gate authority, worker transport, and stage behavior remain unchanged.

## Test plan

Start with the smallest falsifying fixture: replay the current false-stop trace and assert that a filtered `status --next=[]` cannot terminate the iteration while a mod-block, ready gate, or unresolved worker remains. Add deterministic command-log tests for each branch and one runtime-backed Codex wait test. Run the focused status/dispatch/ensigncycle suites, then `go test ./...`, `go test ./... -race`, formatting, and the required host live lanes when contract/runtime adapter files change.

## Stage-specific test gates

- Ideation must replay the falsifying fixture and record the exact ordered trace before selecting a landing shape.
- Implementation must keep the trace assertions independent of contract prose and preserve the existing state and command bytes.
- Validation must cover the empty-first-query retry, gate/merge separation, unresolved-worker wait predicate, and keep-moving matrix, then run full, race, formatting, and relevant live evidence.
