---
title: Repair Opus recorded gate lifecycle
status: backlog
score: "0.90"
source: "xp6 recorded-gate retained evidence"
sprint: opus-live-completeness
sprint-readiness: defer
group: common-evidence
id: 66dpwxgvsxt7cbxhmgvt3qp4
---
## Problem

The Opus recorded-gate lifecycle remains XFAIL after the evidence-only xp6 task archived.

## Value

A Claude Opus user completes the recorded-gate lifecycle with durable prepare, decision, and consume evidence.

## Acceptance criteria

- AC-1: The exact Opus recorded-gate target runs with its XFAIL and reports XPASS-green after the product repair.
- AC-2: Remove only this Opus binding, then prove an exact normal PASS.
- AC-3: Preserve Sonnet, Codex, Pi, other bindings, and recorded-gate assertions.
- AC-4: Full, race, registry, active-owner, and required PR checks pass.

## Priority hold

Do not dispatch this task while the Captain priority is Sonnet and Codex.
