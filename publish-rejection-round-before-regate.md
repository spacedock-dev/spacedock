---
title: Publish the rejected validation round before correction and re-gating
status: backlog
source: "Replacement for archived rejected zbc; Runtime Live E2E rejection-flow evidence showed the FO claimed validation/1 was recorded without invoking gate record --round validation/1."
started:
completed:
verdict:
score: 0.9
sprint: test-behavior-completeness
sprint-readiness: ready
worktree:
issue:
pr:
mod-block:
id: zhcb4bcz1qgcn7ajx2ctxpxk
---

Restore the common rejection journey on Sonnet, Opus, Codex, and Pi by making the First Officer publish the rejected validation round before correction and re-gating.

## Problem

The desired `rejection-flow` journey requires durable evidence of the first rejected validation before rework. Live runs can complete two validation cycles but omit the actual `gate record --round validation/1` command. The final message can claim that the round was recorded even when the command log proves that it was not.

The journey is therefore TODO on Sonnet, Opus, Codex, and Pi. The active TODO
entries name this replacement task instead of archived rejected `zbc`.

## Value

After a rejection, operators can inspect one durable correction history. The first rejected review is recorded before rework, and the later gate concerns the corrected candidate.

## Proposed approach

Ideation must find the smallest change that makes a real First Officer invoke the existing `gate record --round validation/1` path at the correct lifecycle boundary. Reuse the existing recorder, rejection fixture, command log, and durable journey assertion.

When the behavior is proven, remove each TODO only after that runtime has passing
evidence from `TestLiveCommonRejectionFlow`.

## Out of scope

- Do not revive `zbc`'s new `correction-round` field or Reference-binding schema.
- Do not add another recorder, freshness protocol, or test harness.
- Do not weaken the rejection-flow assertion or accept claims from the final message as proof.
- Do not create a separate Opus evidence task. Opus is the pre-release lane of the release workflow.

## Acceptance criteria

**AC-1 (VALUE) - A supported First Officer publishes the first rejected validation round before it starts correction work.**
Verified by: the existing `TestLiveCommonRejectionFlow` command log contains a successful `gate record --round validation/1`, and the existing durable assertion finds the retained round. Removing the invocation makes the journey fail.

**AC-2 (VALUE) - Sonnet, Opus, Codex, and Pi have honest, actionable coverage state.**
Verified by: each runtime either passes the unchanged journey and has no TODO, or retains a TODO that names this active task. No entry names archived `zbc`.

**AC-3 - The fix reuses the current rejection lifecycle and recorder.**
Verified by: the implementation adds no gate storage field, recorder command, freshness schema, or parallel journey harness. A diff review can falsify this criterion by finding such a mechanism.

## Test plan

Run the smallest existing offline checks that cover rejection routing and round
recording. Then run the unchanged focused journey through Sonnet, Codex, and Pi
with local subscription-backed hosts where available. Run Opus during the
pre-release cadence. Remove a runtime TODO only after its real run passes.
