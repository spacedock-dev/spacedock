---
title: Make review-finding disposition workflow-owned and separate in-stage review from rejection routing
status: backlog
source: "Captain boundary audit, 2026-07-30: an end-of-implementation Roborev round never enters feedback-rejection-flow, and finding classification must happen when findings arrive rather than only after rejection."
started:
completed:
verdict:
score: "0.90"
worktree:
issue:
pr:
sprint: durable-decisions
id: rhx820qrkn6vxpday10nch36
---

Separate three currently conflated paths: an implementation ensign consuming an in-stage review, a validator forming a gate recommendation, and the First Officer routing an already-recorded revise decision through `feedback-to`.

The shared rejection flow should route an already-decided correction and must not define Material, Deferred risk, Polish, Needs decision, correct-but-disproportionate, value-AC, ideation-estimate, or 2× policy. The development workflow and reusable development template should own Roborev classification, task ownership routing, expected-surface calibration, and advisory round recording. Shared gate preparation and presentation should preserve workflow-declared finding categories rather than force development tiers.

## Problem

The shared `feedback-rejection-flow` currently carries development-specific finding classification even though it triggers only after a feedback gate rejects. End-of-implementation Roborev findings arrive earlier, inside implementation, so the classification is both globally over-scoped and attached to the wrong trigger. The newer rule that Material does not imply task ownership exists only in this repo's development README and is absent from the shared skill and reusable development template.

## Proposed approach

Ideation should specify a workflow-owned policy boundary: in-stage review consumers classify before editing; validators classify before forming PASSED/REJECTED; cross-stage rejection routing forwards an already-dispositioned revise decision. Generic skills and commission scaffolding refer to the active stage's declared review policy without naming development classes. Development README/template carry the full Roborev taxonomy, four-field evidence, ownership rule, expected-surface/semantic boundary, tolerance, and AC-drift behavior.

## Out of scope

Do not redesign the advisory-round recorder's storage validation in this task; that is owned by the sibling workflow-neutral recorder task. Do not redesign the generic gate criteria-source interface or dispatch stage identity.
