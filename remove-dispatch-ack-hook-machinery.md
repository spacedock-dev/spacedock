---
title: Remove dispatch acknowledgment hook machinery
status: backlog
score: "1.0"
source: Captain emergency no-global-hook directive, 2026-08-10
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-product
id: ca7w23pffeynv53swt2b8zf3
---
## Problem

The merged dispatch acknowledgment mechanism installs production PreToolUse and SubagentStart hooks and durable hidden-ref state primarily for test observability. It does not initiate the dispatched worker, so it adds global production behavior without delivering the required dispatch value.

## Value

Spacedock no longer installs or runs the dispatch acknowledgment hook mechanism. Existing unrelated hooks, later product fixes, and current default-headless XFAIL bindings and active owners remain unchanged.

## Scope

- Remove only production PreToolUse and SubagentStart dispatch acknowledgment hook entries.
- Remove `internal/dispatchack` and its dispatch/status integration.
- Remove hidden-ref and temporary acknowledgment behavior and tests that depend only on that mechanism.
- Preserve unrelated hook configuration and behavior.
- Preserve later product fixes and every current default-headless binding, assertion, reconciliation row, and active owner.
- Add no replacement mechanism.
- Keep kky, 272, v8, and s9 held.
- Use local focused, full, and race checks only. Do not run paid live or Pi.

## Acceptance criteria

- AC-1: Generated Claude and Codex production hook configuration contains no dispatch-ack PreToolUse or SubagentStart entry.
- AC-2: The `internal/dispatchack` package, hidden-ref/temp acknowledgment behavior, and dispatch/status integration are absent.
- AC-3: Unrelated hooks and later product behavior remain byte-equivalent except for mechanically required references to the removed mechanism.
- AC-4: Current default-headless XFAIL bindings, assertions, reconciliation rows, and active owners are unchanged.
- AC-5: Focused removal controls, `go test ./...`, `go test ./... -race`, format, reconciliation, and active-owner checks pass locally.
