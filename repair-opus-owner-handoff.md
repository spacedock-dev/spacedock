---
title: Repair Opus owner handoff
status: backlog
score: "0.90"
source: "xp6 owner-handoff retained evidence"
sprint: opus-live-completeness
sprint-readiness: defer
group: common-evidence
id: bqy97b9npym3zs62pagjchpk
---
## Problem

The Opus owner-handoff journey remains XFAIL after the evidence-only xp6 task archived.

## Value

A Claude Opus user sees one correct conflict owner handoff without stale ownership or hidden failure.

## Acceptance criteria

- AC-1: The exact Opus owner-handoff target runs with its XFAIL and reports XPASS-green after the product repair.
- AC-2: Remove only this Opus binding, then prove an exact normal PASS.
- AC-3: Preserve Sonnet, Codex, Pi, other bindings, and the owner-handoff assertion.
- AC-4: Full, race, registry, active-owner, and required PR checks pass.

## Priority hold

Do not dispatch this task while the Captain priority is Sonnet and Codex.
