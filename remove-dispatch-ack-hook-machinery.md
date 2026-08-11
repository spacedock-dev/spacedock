---
title: Remove dispatch acknowledgment hook machinery
status: implementation
score: "1.0"
source: Captain emergency no-global-hook directive, 2026-08-10
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-product
id: ca7w23pffeynv53swt2b8zf3
started: 2026-08-10T22:49:48Z
worktree: .worktrees/spacedock-ensign-remove-dispatch-ack-hook-machinery
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

## Stage Report: implementation

- DONE: Remove only the production dispatch acknowledgment hooks.
  `hooks.json` now contains only the existing compact `SessionStart` hook.
- DONE: Remove the dispatch acknowledgment package and all product integration.
  The change deletes `internal/dispatchack` and removes its dispatch, envelope, and status code.
- DONE: Remove hidden Git refs, temporary acknowledgment state, and dependent live assertions.
  The prior native worker-lifecycle assertion now checks the default-headless journey again.
- DONE: Preserve later product changes and current live-test ownership.
  The current default-headless bindings, reconciliation rows, assertions, and active owners are unchanged.
- DONE: Add no replacement dispatch mechanism.
  A focused contract test rejects the removed hooks, package, commands, fields, refs, tokens, and temporary state.
- DONE: Run the required full and race checks with a private Go cache.
  `go test ./...` passed. `go test ./... -race` passed.
- DONE: Resolve the known load-sensitive check without changing the candidate.
  The initial 250ms load-sensitive check passed alone. The exact full rerun then passed.
- DONE: Complete the package-timeout triage.
  `TestDurableKeepMovingRequiresOverlappingJourneys` passed in 58.863 seconds.
- DONE: Run focused contract checks.
  The removal control, registry reconciliation, and active-owner checks passed.
- DONE: Check format and the final diff.
  The gofmt diff was empty. `git diff --check` passed.
- DONE: Check the remaining acknowledgment text.
  Acknowledgment mechanism strings remain only in the removal contract test.
- DONE: Commit and push the exact candidate.
  Candidate `e2f4e90e604964ea15326f101b466edc9ed1127c` is on `spacedock-ensign/remove-dispatch-ack-hook-machinery`.
- DONE: Do not run live or Pi checks.
  This implementation used only accepted local evidence.

### Changed files

- Changed `hooks.json`.
- Added `internal/contractlint/dispatch_ack_removal_test.go`.
- Changed `internal/dispatch/build.go`.
- Changed `internal/dispatch/dispatch.go`.
- Deleted `internal/dispatchack/ack.go`.
- Deleted `internal/dispatchack/ack_test.go`.
- Changed `internal/ensigncycle/claude_live_runner_test.go`.
- Changed `internal/status/handlers.go`.

### Summary

Spacedock no longer installs or runs the dispatch acknowledgment mechanism.
The compact session hook and all current live-test bindings and owners remain unchanged.
