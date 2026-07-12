---
id: w8rrjv6jsmgahc269arx9k5v
title: Restore Codex handoff bootstrap with post-fence options
status: backlog
source: Captain-reported regression, 2026-07-13
started:
completed:
verdict:
score:
worktree:
issue:
---

Restore a fresh Codex launch when a Spacedock handoff is supplied before the
front-door fence and Codex options follow it after the fence.

## Problem

`spacedock codex @/tmp/handoff-file.md -- --model gpt-5.6-sol` forwards the
Codex model option but drops the Spacedock first-officer bootstrap prompt and
the handoff text. The resume repair made every nonempty post-fence argv suppress
the bootstrap, which is broader than the resume case.

## Proposed approach

Keep post-fence argv opaque. Use the already-parsed presence of a pre-fence task
to preserve the explicitly requested fresh-launch bootstrap and its normal fresh
launch posture, without recognizing or reconstructing Codex options or
subcommands. Update the command reference to state this precedence.

## Out of scope

No Codex argv parser, new command surface, approval-policy change for no-task
post-fence launches, or changes to resume invocations without a pre-fence task.

## Acceptance criteria

**AC-1 (VALUE) - A fresh handoff launch retains both the handoff and the first-officer bootstrap when Codex options are supplied after `--`.**
Verified by: a focused `runCodex` argv assertion for `@/tmp/handoff-file.md -- --model gpt-5.6-sol`.

**AC-2 - Resume-family post-fence argv without a pre-fence task remains prompt-free.**
Verified by: the existing focused resume argv assertion continues to pass.

**AC-3 - The repair leaves post-fence Codex tokens in operator order without a Codex grammar table.**
Verified by: the focused argv assertion and the existing opaque-passthrough test run.

**AC-4 - The command reference distinguishes an explicit pre-fence task from a no-task opaque Codex launch.**
Verified by: review against the focused argv behavior and the updated launch reference.

## Test plan

Add one focused Go unit test at the existing fake-host seam, then run the focused
package test and `go test ./...`. The behavior is pure argv construction, so no
live host or timeout test is justified.
