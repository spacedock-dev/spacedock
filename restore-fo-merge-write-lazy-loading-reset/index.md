---
title: Restore lazy loading for first-officer merge and write cores (clean reset)
status: implementation
source: Clean reset from rejected 1k implementation, captain direction 2026-07-14
started: 2026-07-14T14:12:36Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-restore-fo-merge-write-lazy-loading-reset
issue:
milestone: 0.25.0
id: gk7ceyrs4496jgp535w3awfp
---

Restore the intended first-officer loading boundary: boot reads the shared core and active runtime adapter, while write authority loads at the first FO-authored mutation and merge handling loads only at a terminal or merge-mod recovery boundary.

## Problem

The earlier fix for delayed-reference filesystem hunting overcorrected by eagerly importing `fo-write-core.md` and `fo-merge-core.md` from the first-officer entry skill. A mutation-free interactive run that stops at a gate now pays for mutation and terminal ceremony it never uses. The defect was nondeterministic discovery, not lazy loading: deferred reads must use the loader-supplied first-officer base plus one literal canonical suffix, never cwd, wrapper-skill discovery, alternate paths, or search.

The prior implementation branch and its proof harness were rejected. This task intentionally carries no prior stage reports, feedback cycles, parser design, or implementation plan.

## Required outcome

- `skills/first-officer/SKILL.md` eagerly imports only the shared core; the active runtime adapter remains a boot read.
- A mutation-free gate hold reads neither write nor merge core.
- The first FO mutation reads the exact write core before the mutation. The first terminal or `mod-block=merge:*` recovery reads the exact merge core before merge-owned action; when both trigger together, write precedes merge.
- Merged-PR discovery is not required on boot; it may happen on `engage`.
- Do not change mutation resolution, merge resolution, routing, or terminal semantics. If a resolution is broken, that is a separate defect shared by eager and lazy loading.

## Mechanism/value trace

- Served value: avoid at least 8,000 bytes of unused cold contract while preserving supported-host behavior.
- Simplest route: change the import/load cues and observe real host tool-call order plus durable workflow outcomes through the existing scenario runners.
- Rejected heavier route: no new command parser, operation-language interpreter, PTY/session controller, lifecycle state, daemon, lease, or recovery protocol. If existing host events cannot prove a claim, stop for architecture review instead of implementing a second runtime.

## Acceptance criteria

- **AC-1 (VALUE):** Fresh Claude and Codex mutation-free gate journeys load shared core plus the selected adapter, read neither deferred core through the gate, and save at least 8,000 cold bytes against the eager baseline.
- **AC-2:** Existing filing and terminal/recovery journeys observe exact write/merge reads before their real supported commands, with no broad search, wrapper invocation, alternate-path retry, or earlier owned action; their durable success/refusal/archive outcomes remain unchanged.
- **AC-3:** Structural coverage proves one eager canonical import and two resolvable delayed files only. It does not infer runtime order from instruction prose.
- **AC-4:** Focused tests, `go test ./...`, `go test ./... -race`, and the relevant exact-head local Codex live journeys are green before Roborev is requested. The harness orchestrates and observes the supported runtime; it does not model arbitrary shell semantics.
