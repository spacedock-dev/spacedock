---
title: First Officer surfaces workflow-declared obligation authority at gates
status: backlog
source: "Durable-decisions sprint audit of Subspace 198f762 and Spacedock 6y, 2026-07-24."
score: "0.8"
id: 6hg9nxd47yfnrjxhb4vrjdpd
---

A First Officer can review evidence at a gate without noticing that a worker or reviewer added obligations that the task, captain, or workflow never authorized. Codify the stage-neutral gate behavior that prevents an inferred constraint from silently becoming acceptance scope.

## Problem

The reusable First Officer contract cannot name development-specific stages such as ideation. It also must not turn reviewer feedback into authority. When a workflow declares an obligation-provenance or constraint-authority policy, the FO needs a stage-neutral way to compare the proposed result with the accepted scope and surface material added or widened obligations before deciding the gate.

The current 6y incident is the motivating failure: an implementation/review loop promoted public exact-child transport forensics and a route-by-host live matrix into mandatory proof even though the product value required only the recorded gate lifecycle, one successor dispatch, and a durable effect.

## Proposed approach

Ideation should define the smallest generic FO rule:

- At any gate, when the workflow declares an obligation-authority policy, surface material obligations added, widened, or removed since the last accepted Briefing and cite their stated authority.
- Treat worker and reviewer findings as evidence, never authority by themselves.
- Hold an unsupported obligation for captain decision and route it using the workflow's own gate/revision semantics.
- Do not name or assume ideation, implementation, validation, a particular task schema, or a new gate command.

## Obligation delta

- **Bearer:** the First Officer, only while assembling or deciding a gate for a workflow that declares this policy.
- **Burden:** one provenance comparison over material obligation changes; no new review round or ceremony.
- **Authority:** captain direction from the durable-decisions sprint plus the workflow-declared policy. Reviewer feedback supplies evidence only.
- **No inferred obligation:** workflows with no such declaration retain their existing gate behavior.

## Out of scope

- Defining the workflow-specific policy text or its stage-specific consumption.
- Transporting workflow sections into worker dispatches.
- Adding a new gate, recorder, resolution, or state-machine concept.
- Treating every wording change as an obligation change.

## Acceptance criteria

**AC-1 (VALUE) - A workflow-declared obligation policy prevents an unsupported material obligation from silently passing a gate.**
Verified by: a fixture or live FO journey with arbitrary stage names in which a worker widens one obligation without task, captain, or named-contract authority; the gate presentation identifies the delta and does not approve or advance it. A control with the same evidence but explicit captain authority may proceed.

**AC-2 - The reusable FO rule is workflow- and stage-neutral.**
Verified by: the same behavioral fixture using non-development stage names, plus inspection showing the generic contract does not require the strings `ideation`, `implementation`, or `validation`.

**AC-3 - Reviewer findings remain evidence rather than self-authorizing scope.**
Verified by: a fixture where a reviewer requests an extra constraint; without a cited authority the resulting gate is held for decision, while an ordinary finding that proves an already-authorized AC failure follows the workflow's normal revision route.

**AC-4 - Workflows that declare no obligation-authority policy do not acquire a new ceremony.**
Verified by: an existing FO gate lifecycle fixture with no declaration whose dispatch, presentation, decision, and transition trace remains unchanged.

## Test plan

Prefer augmenting the existing recorded-gate/FO lifecycle scenario rather than creating another runtime controller. Use one declared-policy fixture, one authority-present control, one no-policy regression, and a structural stage-neutrality check. A live host run is required only if the chosen contract wording claims agent behavior that fixture execution cannot establish.
