---
title: Replace remaining development-shaped gate and dispatch assumptions with workflow-declared evidence interfaces
status: backlog
source: "Captain shared-contract audit, 2026-07-30: experiment and refinement templates expose criteria, review ledgers, and stage names that the shared FO/runtime contract does not understand."
started:
completed:
verdict:
score: "0.70"
worktree:
issue:
pr:
sprint: durable-decisions
id: yd1cd7jdywby0c3r46bqtm25
---

Define the smallest workflow-declared interface that lets generic First Officer and dispatch machinery locate gate criteria, review history, feedback-stage identity, and worker stage identity without assuming the development workflow's headings or five stage names.

## Problem

The shared FO gate cross-check requires `## Acceptance criteria` with `**AC-N**` mechanism/value pairings, while the experiment template uses `## Success criteria` and refinement uses `## Review notes`. Runtime prose names a `validation` reviewer, the FO write contract privileges `### Feedback Cycles`, and dispatch reconciliation falls back to `backlog / ideation / implementation / validation / done`. These are development-shaped assumptions in paths intended to operate commissioned workflows generally.

## Proposed approach

Ideation should first exercise one experiment gate and one refinement rejection end to end to identify the minimum missing pointers. Prefer existing workflow/stage metadata and dispatch-emitted identity over new schema. If a declaration is required, define one coherent interface consumed by gate evidence extraction, feedback routing, write scope, and roster reconciliation; do not add independent ad-hoc flags for each reader. Preserve development's current behavior as one template realization, not a global default.

## Out of scope

Do not revisit development's Roborev taxonomy or advisory-round body validation; sibling tasks own those corrections. Do not add compatibility handling for unreleased prototype formats. Do not introduce a second criteria parser or duplicate worker-identity registry when existing metadata can be extended.
