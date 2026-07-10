---
title: Replace mutable step-number references with named FO functions
status: backlog
score: 0.95
source: "Captain direction 2026-07-11: widen the step-number sweep; hooks and references use the «fn» notation."
id: 88tq5zyg9jvx13f33zz3eq28
---

## Problem

The active first-officer contract still addresses shared procedures by mutable numbers such as Startup step 3, Dispatch step 2, Event Loop step 0/3, and Merge-and-Cleanup step 10. The Claude adapter and shared cores therefore couple to list positions even though docs/runtime-support.md requires adapters to bind named capabilities.

## Recommended design

Keep numbered lists only where local ordering aids execution. Any cross-line, cross-section, cross-file, hook, or runtime-override reference must target a named `«fn»` capability instead of a step number.

- Boot and interaction references bind `«state.boot»()` and a named interaction-boundary function.
- Dispatch checklist references bind a named checklist function.
- Event-loop references bind `«roster-reconcile»()` and `«dispatch.next-action»()`.
- Terminal teardown references bind `«worker.shutdown»()`.
- Lifecycle-hook references bind a named hook-runner function parameterized by hook point.
- Runtime adapters override or bind these functions; they do not copy numbered shared procedures.

## Scope

- Sweep the active first-officer skill, its shared/Claude runtime references, and directly linked gate/legacy helpers.
- Extend contract lint across Claude and shared cores so mutable numeric procedure references cannot return.
- Preserve the behavior and local ordering of the underlying procedures; this is reference normalization, not a new boot or merge sequence.

## Acceptance criteria

- **AC-1:** No active first-officer runtime adapter or shared core refers to a procedure as `step N`, `step-N`, or equivalent mutable numeric address.
- **AC-2:** Cross-file boot, dispatch, event-loop, teardown, and lifecycle-hook references use named `«fn»` notation with an explicit binding or override.
- **AC-3:** A contract-lint regression test fails on representative numeric references in Claude, Codex, Pi, and shared-core surfaces.
- **AC-4:** Existing contract, offline, race, and live-runtime gates show unchanged workflow behavior after the reference normalization.

## Test plan

- Add the failing contract-lint coverage first and confirm it catches the current Claude/shared-core residues.
- Replace numeric references with named function bindings and run focused contractlint tests.
- Run formatting, full tests, race tests, and the protected runtime live matrix because shipped first-officer prompt text changed.
