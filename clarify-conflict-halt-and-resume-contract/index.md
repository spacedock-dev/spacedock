---
title: Clarify safe conflict halt and resume in the First Officer contract
status: backlog
source: "Captain follow-up, 2026-08-03: the live contract and observed merge friction did not make conflict ownership, affected scope, and unrelated-work continuation obvious; Science Officer advisory confirmed the generic boundary."
started:
completed:
verdict:
score:
worktree:
issue:
pr:
id: q9yxrwpccbt8c410f19d0803
---

Make the generic First Officer contract state the safe behavior after a state-sync or merge conflict. Keep this follow-up outside the durable-decisions sprint because it changes the shared contract, not the development workflow deliverable.

## Problem

The contract names a halt for rebase conflicts, but it does not clearly state the affected-operation boundary, the evidence to preserve, or whether unrelated ready work can continue. This caused confusion between safe halt behavior and a proposed automatic resolver route.

## Proposed approach

Update the generic skills/first-officer contract through a dispatched worker.

- Halt only the affected operation and entity.
- Abort the in-progress rebase or merge. Surface the peer, conflict paths, and the exact next recovery condition.
- Never force-push, auto-resolve, or silently discard either side.
- Keep unrelated entities moving when they are already genuinely ready and owned.
- Keep terminal delivery retry behavior and explicit product rework separate from conflict handling.
- Do not add an automatic conflict classifier or resolver-worker dispatch. A development-only finding taxonomy belongs in docs/dev/README.md or a separate approved task.

## Out of scope

Do not change the development workflow's semantic labels. Do not create a generic resolver worker. Do not change merge policy, terminal approval semantics, or force-push rules.

## Acceptance criteria

**AC-1 - The generic First Officer contract defines the affected-operation halt, abort, evidence, and no-force/no-auto-resolve rules.**
Verified by: the changed skills/first-officer references and the existing contract-oracle or skill smoke tests.

**AC-2 - The contract states when unrelated ready entities may continue and names the recovery condition for the halted entity.**
Verified by: a contract-oracle fixture or equivalent event-loop test that observes the halted entity and an independent ready entity.

**AC-3 - The generic contract does not introduce automatic conflict classification or resolver-worker dispatch.**
Verified by: contract diff review plus the existing merge/state-sync tests showing the current safe halt behavior remains intact.

## Test plan

Run the focused First Officer contract and skill smoke tests, the state-sync and merge-guard conflict fixtures, go test ./..., go test ./... -race, and gofmt -w ./cmd ./internal. Use one live or fixture-backed conflict case to prove durable peer/path evidence and recovery re-entry.
