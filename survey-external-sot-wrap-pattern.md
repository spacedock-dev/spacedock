---
id: ye2etq07y7xyr3hfh1w8xp7d
title: The survey detects an external staged source of truth and proposes the wrap pattern
status: backlog
source: "Field case 2026-08-26 (Worca) + strategy diagnosis /tmp/survey-external-sot-diagnosis.md: the user's agent rejected adoption because two systems claimed staged state; the survey produced an accurate inventory and a generic offering, never connecting the user's named friction (skipped steps) to the mechanism that answers it; corroborated by the internal health-check verdict"
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
---

When the surveyed project's staged truth lives in an external system, the survey degrades to a generic report and the adoption gatekeeper — the user's own agent — correctly rejects Spacedock as a second record home. The survey must detect that case and propose the wrap pattern instead: records stay put; the entity is a decision shell (external ref, evidence snapshot in the gate room at decision time, write-back as a mod); Spacedock owns the decision lifecycle, never the domain lifecycle.

## Problem

{Ideation fills this in. Seeded: the detection tell is a triple — staged vocabulary in conversation and skill names with no staged artifact in the repo, source-of-truth client code or credentials, and approvals traveling by mail. The recommendation must lead with the user's named friction mapped to the enforcing mechanism (a stage machine that refuses to advance; decidable approval queues), not with an inventory. The wrap pattern needs no new core machinery: the issue field already carries external refs, the Briefing already digest-freezes snapshot artifacts in the room, and the mod system already models write-back (pr-merge precedent). The beachhead pick matters: usually the dev workflow — the one staged thing the external system does not own. A live experiment (decision-point-table reframing prompt) is in flight; its returned table is test data for this design.}

## Proposed approach

{Ideation fills this in. Scope: the survey skill's detection questions, the wrap-pattern proposal section with the decision-shell framing, and the friction-to-mechanism recommendation ordering. Naming and documenting the wrap pattern (the running GTM-against-HubSpot reference) may be a deliverable here or a sibling doc task — ideation rules with the captain.}

## Risk evidence

{Backlog: one live rejection with the agent's reasoning recorded, plus the independent internal health-check verdict, decide design should start. Evidence homes cited in the source diagnosis.}

## Out of scope

Any core binding feature, live stage mirroring (refused: a sync daemon is the second-implementation class), write-back mods themselves.

## Expected surface and tolerance

{Backlog seed; ideation estimates with the production/proof split; the surface is skill and doc prose.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A survey run against a project with an external staged source of truth proposes the wrap pattern with a beachhead and a friction-to-mechanism lead, instead of a generic inventory report.**
Verified by: {ideation refines — seed: a fixture project with the detection triple; failing-today baseline: the current survey emits the generic report on the same fixture. The in-flight field experiment's table, if returned, is the live falsification.}

## Test plan

{Ideation fills this in. skills/** is shipped scaffolding: the path-to-lane rule applies.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
