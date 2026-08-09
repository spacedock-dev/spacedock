---
title: Make headless Sonnet and Codex spawn implementation before validation
status: backlog
source: "PR #583 run 31320596435, Codex job 93262943132 and Sonnet job 93262943118, 2026-08-09"
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

The Sonnet and Codex `default-headless-gate-stop` journeys created an implementation dispatch package but did not spawn the worker.

Each First Officer then changed the task status to validation. Each one prepared the validation gate without an implementation worker report.

The final task had no worker-produced implementation report. The live oracle rejected the run with this error:

`gate hold crossed its committed no-authority boundary: implementation was not dispatched before validation`

The failures occurred on candidate `7c8708c8537fc73761e56813ddbd6a498959ef19`. Codex ran for 166.56 seconds. Sonnet ran for 414.17 seconds.

## Value

A headless First Officer must run the required work before it presents the next gate. A dispatch package is not a worker dispatch.

## Required order

Task `ts7gq0mr9s3chx2w4wppd1kt` must land before this product repair starts.
It must run the Sonnet and Codex cells as strict XFAIL evidence.

Do not change the First Officer behavior before both cells report the expected
`implementation-worker-not-dispatched` failure code. After the repair, XPASS
must force removal of each repaired XFAIL binding.

## Acceptance criteria

**AC-1 (VALUE) — Sonnet and Codex run the implementation worker before validation.**

The `default-headless-gate-stop` journey observes one implementation worker spawn. It also observes the worker completion before the validation transition.

Verified by: passing Sonnet and Codex runs of `TestLiveCommonDefaultHeadlessGateStop` with the current fixture and durable oracle.

**AC-2 — A status change cannot replace the worker dispatch.**

The First Officer must not credit implementation from `dispatch build`, narration, or a direct status change.

Verified by: the existing provider-neutral command log and durable entity reject each missing-worker form.

**AC-3 — Each runtime TODO is removed only after passing live evidence.**

The journey keeps each TODO until the exact candidate passes that runtime lane.

Verified by: registry reconciliation and the passing live artifact use the same commit.

## Scope

- Correct the shared First Officer behavior that lets Sonnet and Codex skip the implementation worker.
- Preserve the coherent 26n fixture and its no-authority oracle.
- Use each shipped runtime worker-dispatch mechanism.
- Keep product behavior portable. Do not add a test-only host switch.

## Out of scope

- Weakening the live oracle.
- Restoring the contradictory fixture.
- Treating `dispatch build` as worker-spawn evidence.
- Adding a simulator or a test for test infrastructure.
- Repairing Pi behavior. Pi has no new failure evidence from this run.

## Evidence

Run `31320596435`, jobs `93262943132` and `93262943118`, record this sequence:

1. The First Officer changed `queued` to `implementation`.
2. The First Officer created the implementation dispatch package.
3. The First Officer did not spawn the implementation worker.
4. The First Officer changed `implementation` to `validation`.
5. The First Officer prepared and committed the validation gate.

The run proves a product behavior gap. The repaired fixture made the gap visible.
