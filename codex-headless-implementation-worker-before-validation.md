---
title: Make Codex spawn the headless implementation worker before validation
status: backlog
source: "PR #583 run 31320596435, Codex job 93262943132, 2026-08-09"
started:
completed:
verdict:
score: 0.95
worktree:
issue:
pr:
mod-block:
sprint: test-behavior-completeness
group: common-evidence
sprint-readiness: ready
id: 98aa776adg66gn823a8gamdq
---

## Problem

The Codex `default-headless-gate-stop` journey created an implementation dispatch package but did not spawn the worker.

Codex then changed the task status to validation. It created a validation dispatch package and prepared the validation gate.

The final task had no worker-produced implementation report. The live oracle rejected the run with this error:

`gate hold crossed its committed no-authority boundary: implementation was not dispatched before validation`

The failure occurred on candidate `7c8708c8537fc73761e56813ddbd6a498959ef19`. The test ran for 166.56 seconds.

## Value

A headless Codex First Officer must run the required work before it presents the next gate. A dispatch package is not a worker dispatch.

## Acceptance criteria

**AC-1 (VALUE) — Codex runs the implementation worker before validation.**

The `default-headless-gate-stop` journey observes one implementation worker spawn. It also observes the worker completion before the validation transition.

Verified by: a passing Codex run of `TestLiveCommonDefaultHeadlessGateStop` with the current fixture and durable oracle.

**AC-2 — A status change cannot replace the worker dispatch.**

The First Officer must not credit implementation from `dispatch build`, narration, or a direct status change.

Verified by: the existing provider-neutral command log and durable entity reject each missing-worker form.

**AC-3 — The Codex TODO is removed only after passing live evidence.**

The journey keeps its Codex TODO until the exact candidate passes the required Codex lane.

Verified by: registry reconciliation and the passing live artifact use the same commit.

## Scope

- Correct the Codex First Officer behavior that skips the implementation worker.
- Preserve the coherent 26n fixture and its no-authority oracle.
- Use the shipped Codex worker-dispatch mechanism.
- Keep product behavior portable. Do not add a test-only host switch.

## Out of scope

- Weakening the live oracle.
- Restoring the contradictory fixture.
- Treating `dispatch build` as worker-spawn evidence.
- Adding a simulator or a test for test infrastructure.
- Repairing Pi or Claude behavior.

## Evidence

Run `31320596435`, job `93262943132`, records this sequence:

1. Codex changed `queued` to `implementation`.
2. Codex ran `dispatch build --stage implementation`.
3. Codex did not call a worker-spawn tool.
4. Codex changed `implementation` to `validation`.
5. Codex prepared and committed the validation gate.

The run proves a product behavior gap. The repaired fixture made the gap visible.
