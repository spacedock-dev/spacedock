---
title: Restore live evidence whose original repair owners are complete
status: backlog
source: "Live-test-truth close reconciliation, 2026-08-09"
started:
completed:
verdict:
score: 0.9
worktree:
issue:
pr:
mod-block:
sprint: test-behavior-completeness
group: common-evidence
sprint-readiness: ready
id: xp6c9qfe7y4wwp46enc3f85n
---

Three live TODO groups still name completed and archived repair tasks. The code
changes landed, but the required runtime evidence did not remove these TODOs.

## Value

Each remaining TODO points to active work. Passing evidence removes the TODO for
that exact runtime and journey.

## Required order

Use the strict XFAIL behavior from `ts7gq0mr9s3chx2w4wppd1kt` before a new
product repair. Convert a TODO only when the complete journey runs and produces
one stable semantic failure code.

Remove a TODO immediately when the journey passes. Remove an XFAIL binding when
the repaired journey reports XPASS.

## Acceptance criteria

**AC-1 (VALUE) — Restore `gate-guardrail` evidence on Codex and Pi.**

Run the current common journey on each target. Remove each TODO only after its
exact candidate passes the durable assertion.

**AC-2 (VALUE) — Restore Opus `recorded-gate-lifecycle` evidence.**

Run the pre-release Opus lane. Remove the TODO only after the exact candidate
passes the durable assertion.

**AC-3 (VALUE) — Restore `owned-conflict-owner-handoff` evidence on Sonnet,
Opus, and Pi.**

Run the current common journey on each target. Remove each TODO only after its
exact candidate passes the durable assertion.

**AC-4 (VALUE) — Restore Pi `default-headless-gate-stop` evidence.**

Run the coherent fixture and durable oracle on Pi. Remove the TODO only after
the exact candidate passes.

**AC-5 — Preserve honest failures.**

If a target fails, diagnose the product behavior. Keep the TODO until a product
repair and exact-target evidence pass.

## Scope

- Reuse the current common journey entry points, fixtures, and assertions.
- Use local subscription-backed live runs before paid CI where possible.
- Keep product repairs in this task only when they directly restore one named
  journey result.

## Out of scope

- Weakening a durable assertion.
- Adding a simulator.
- Adding tests for the reconciliation code.
- Changing the desired journey registry without a desired behavior change.
