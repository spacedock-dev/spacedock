---
id: 3n6qcjttj261wjcdkvc71k6z
title: User-facing journey for the 0.27 gate/resolution benefit
status: backlog
source: "Captain directive, 2026-08-14: for the 0.27 gate/resolution rollout, publish a user-facing journey that describes the benefit of this major release"
started:
completed:
verdict:
score:
worktree:
issue:
---

Author a user-facing journey that shows the benefit of the 0.27 gate/resolution feature from the user's side: work is dispatched, a worker delivers, a gate briefing is prepared (`gate prepare`), the captain decides approve/revise/hold (`gate record`), feedback rounds route back to implementation, and the resolution is consumed (`gate consume`) into a durable decision record, ending in `merge guard` and done. The narrative must sell the end value — the captain keeps decision-quality control over agent work, with a durable audit trail of every decision — not the mechanism.

Dogfood proof points available from the 2026-08-14 0.27 audit: 347 gate briefings prepared, 392 decisions recorded, 282 resolutions consumed, 212 feedback-cycle entries across 53 entities, 319 entities driven to done.

Placement candidates for ideation to decide: extend docs/site/concepts/gates-and-decisions.md, a new page under docs/site/running-workflows/, a get-started/first-workflow.md tie-in, and/or the 0.27 release notes.

## Problem

{What is broken or missing, and why it matters. Ideation fills this in.}

## Proposed approach

{How the task intends to solve the problem. Ideation fills this in — including final placement and the walkthrough scenario.}

## Out of scope

{What this task deliberately does not cover, so the boundary is explicit.}

## Expected surface and tolerance

Estimate net LOC change: {+NNN}, across {M} files.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - {Ideation refines. End value to anchor on: a reader new to gate/resolution can state its benefit after reading the journey.}**
Verified by: {ideation names the concrete check}

## Test plan

{What verifies the implementation; docs-only change expected — link/build checks plus prose review.}
