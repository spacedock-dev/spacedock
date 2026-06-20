---
title: Convert Codex and Pi runtime adapters to capability binding blocks
status: backlog
source: "Captain direction (2026-06-20): move toward per-host runtime files as bindings blocks keyed by core «fn» capability names, starting with Codex and Pi; recommend sequencing with ad/trim-dispatch-adapter-prose."
started:
completed:
verdict:
score:
worktree:
issue:
sprint: 0221-layered-fo
sprint-readiness:
id: t0gk2fatt18tj28xm6sr1xd1
---

Codex and Pi first-officer runtime references should become compact runtime implementation maps keyed by the shared core's `«fn»` capability names, rather than prose sections that re-narrate the lifecycle. The target shape is a short runtime intro, a `## Runtime implementation` binding list, and only short residual sections for probe/harness notes that do not fit a capability name.

## Problem

`skills/first-officer/references/pi-first-officer-runtime.md` and `skills/first-officer/references/codex-first-officer-runtime.md` already contain most of the required binding facts, but they are spread across lifecycle prose sections (`Dispatch`, `Awaiting Completion`, `Follow-up and Reuse`, `Shutdown`, etc.). Pi still carries negative host contrast and step-number coupling. Codex is probe-based after PR #414 but still mixes binding facts with lifecycle narration. This duplicates responsibilities the shared core should own: when each capability is invoked.

## Proposed approach

Create a focused Codex+Pi cleanup that converts each adapter to a binding-first shape:

- short file intro: the shared core owns invocation timing; this file binds host realization
- `## Runtime implementation` with bullets for `«worker.spawn»`, `«addressable-worker»`, `«async-dispatch»`, `«worker-identity»`, `«completion-signal»`, `«worker.shutdown»`, `«context-budget»`, and `«roster-reconcile»`
- short residual sections only for live-tool probes, harness isolation, or compatibility notes that are not lifecycle bindings
- no negative host contrast except deliberate compatibility hazards
- no mutable shared-procedure step-number coupling
- no dispatch/await/reuse lifecycle re-narration unless it carries a host-specific guardrail not expressible in a binding bullet

This task should add or adjust contractlint guards so Codex/Pi adapter tool names stay in their runtime binding/probe sections and so the adapters do not reintroduce `## Dispatch`, `## Awaiting Completion`, `## Follow-up and Reuse`, or `## Shutdown` as lifecycle narration after the conversion.

## Sequencing recommendation with `ad`

Recommended order: run this task before `ad`'s broad adapter trim.

Reason: this task is a narrow structural normalization for Codex+Pi only. It creates the binding-block shape that `ad` can then use as the target when trimming remaining per-adapter operational prose. Running `ad` first would force it to decide the same binding-placement question while also preserving Claude await/reuse/guardrails, increasing risk and review surface.

Boundary: this task should not touch Claude except for shared test utilities if unavoidable. `ad` remains owner of Claude adapter trim, await/reuse/guardrail preservation, broad per-adapter prose reduction, and any cross-runtime final pass after Codex+Pi have binding blocks. If `ad` adds first-class `«worker.spawn»` / `«worker.shutdown»` to the core first, this task should consume that; otherwise this task may add the minimal core capability headings required for Codex/Pi binding bullets and leave broad Claude migration to `ad`.

## Out of scope

Do not change runtime behavior, live runners, launch/install UX, or Claude guardrails. Do not remove Codex wait-interruption semantics or feedback-reviewer reuse semantics; express them as binding/probe notes if they remain load-bearing.

## Acceptance criteria

**AC-1 - Codex and Pi FO runtime adapters use a binding-block structure.**
Verified by: diff review and contractlint showing each file has a `## Runtime implementation` section keyed by the shared `«fn»` capability names and no longer uses lifecycle narration headings for dispatch/await/reuse/shutdown.

**AC-2 - Runtime support principles are enforced for Codex and Pi.**
Verified by: contractlint or targeted tests rejecting negative host contrast, mutable shared-procedure step-number coupling, and concrete host tool names outside binding/probe sections in the Codex/Pi FO runtime adapters.

**AC-3 - Load-bearing host-specific semantics are preserved.**
Verified by: focused tests and review showing Codex probe-based tool binding, wait reinstallation/interruption semantics, feedback reviewer reuse, Pi model stamping, Pi file-verification gate, and Pi harness isolation notes remain present in compact form.

**AC-4 - The task composes cleanly with `ad`.**
Verified by: implementation report naming what remains for `ad`, with no Claude adapter rewrite beyond unavoidable shared/core/test edits.

## Test plan

Run `go test ./internal/contractlint` and any focused Codex/Pi runtime contract tests touched by the change. If shared `fo-dispatch-core.md` capability headings are added or moved, run the existing capability-binding tests and `go test ./...`.
