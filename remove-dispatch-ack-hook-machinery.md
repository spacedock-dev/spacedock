---
title: Remove dispatch acknowledgment hook machinery
status: validation
score: "1.0"
source: Captain emergency no-global-hook directive, 2026-08-10
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-product
id: ca7w23pffeynv53swt2b8zf3
started: 2026-08-10T22:49:48Z
worktree: .worktrees/spacedock-ensign-remove-dispatch-ack-hook-machinery
gates:
    version: 1
    records:
        - id: gate:ca7w23pffeynv53swt2b8zf3:validation
          stage: validation
          attempts:
            - id: gate-attempt:ca7w23pffeynv53swt2b8zf3-validation-1
              briefing:
                id: briefing:ca7w23pffeynv53swt2b8zf3:validation:attempt-1:revision-1
                digest: sha256:aeb61ccb352ed0228ea178896445894cd1a6c3a799afab6e39201953a4eb7d56
                request-digest: sha256:617bc292f332b01b14553404bb7a5aecdbdc2873805152713149e75aa86867b3
                room-ref: ./remove-dispatch-ack-hook-machinery/review/validation/briefing-1
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

## Stage Report: validation

- DONE: Check the immutable candidate and remote branch.
  Both resolve to `e2f4e90e604964ea15326f101b466edc9ed1127c`.
- DONE: Inspect all eight changed files against base `ff9bb4506be73787a684e5fd80b7b772ea7473a5`.
  The candidate changes 8 files with 67 additions and 469 deletions.
- DONE: Check the production hook removal.
  `hooks.json` contains no dispatch-ack `PreToolUse` or `SubagentStart` hook.
- DONE: Check the unrelated hook.
  The compact `SessionStart` hook is byte-equivalent to the base hook.
- DONE: Check the package and integration removal.
  The candidate deletes `internal/dispatchack` and its dispatch, envelope, status, hidden-ref, and temporary-state behavior.
- DONE: Check for replacement machinery.
  Repository search found acknowledgment mechanism strings only in the removal contract test.
- DONE: Check the restored lifecycle assertion.
  The default-headless journey again uses the existing native implementation-worker lifecycle assertion.
- DONE: Check later behavior and live-test ownership.
  The default-headless XFAIL bindings, assertions, reconciliation rows, and active owners match the base.
- DONE: Inspect the accepted implementation evidence.
  Private-cache full and race checks passed. The exact keep-moving check also passed.
- DONE: Inspect the focused and structural evidence.
  Removal, registry, owner, format, and diff checks passed.
- DONE: Record the validation-shell termination.
  The validation batch ended after its first passing focused command. This event is target-external infrastructure evidence.
- DONE: Run no live or Pi check.
  Validation used static inspection and the Captain-accepted local evidence.
- DONE: Classify Material findings.
  No Material product finding remains. Therefore, released workflow, harm, authority, and trigger fields do not apply.

### Recommendation

PASSED

### Summary

The candidate removes the global dispatch acknowledgment mechanism without a replacement.
The unrelated hook, later product behavior, default-headless bindings, and active owners remain unchanged.
