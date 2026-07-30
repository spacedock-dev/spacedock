---
title: Make dispatch prompt ownership explicit in the First Officer contract
status: backlog
sprint:
source: "Captain follow-up to self-contained-ensign-dispatch, 2026-07-27"
started:
completed:
verdict:
score: "0.8"
worktree:
issue:
id: 5f3ndvve613vqkpgmesd6ydm
---

Prevent future First Officers or runtime adapters from reconstructing full assignment prompts after `dispatch build` has produced the authoritative dispatch artifact.

## Problem

The pointer-only decision can regress if it remains only in task history or one runtime adapter. The authoritative cross-runtime dispatch contract must make prompt ownership unmistakable: `dispatch build` constructs the assignment, and the First Officer transports the emitted prompt unchanged. Full-text manual construction is legal only through the existing break-glass recovery path after a real helper failure.

## Boundary

The primary contract belongs beside `«dispatch.build»` in `skills/first-officer/references/fo-dispatch-core.md`. `internal/dispatch/build.go`, dispatch behavior tests, and the ensign contract may mirror or mechanically prove it. Do not create an ADR, new protocol, lint, output field, or workflow-specific rule.

## Acceptance criteria

**AC-1 — Cross-runtime ownership is explicit.** The First Officer dispatch core states the named invariant **Self-contained artifact; pointer-only transport**, requires fresh and reuse-advance helper prompts to be forwarded unchanged, and forbids assignment-payload construction outside the existing break-glass path.

Verified by: contract review against all runtime adapters; deleting or weakening the cross-runtime rule makes the review fixture fail.

**AC-2 — The ownership boundary is mechanically protected.** Relational dispatch tests grow stage/context/checklist/standing/scope/feedback payloads and prove the dispatch file changes while fresh and reuse-advance outer prompts remain byte-identical and payload-free across Claude, Codex, and Pi.

Verified by: the focused dispatch fixture; moving any payload sentinel into an outer prompt makes it fail.

**AC-3 — The rule is not duplicated into the wrong authority layer.** Runtime adapters only map emitted fields to host tools, the ensign contract only consumes the artifact, and workflow-specific documentation does not redefine prompt construction.

Verified by: targeted review of the dispatch core, runtime adapters, ensign core, and dev workflow.

## Dependency

Reconcile with `self-contained-ensign-dispatch` (`kd`): if that implementation already delivers these exact properties, close this task from its evidence rather than duplicating product changes.
