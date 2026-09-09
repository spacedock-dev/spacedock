---
title: Make stage occupancy and ungated completion explicit
status: backlog
source: Captain discussion of email triage seed skipping intake, 2026-09-09
started:
completed:
verdict:
score: 0.8
worktree:
issue:
pr:
id: 1ytadmcakh5s27r8wjn7qa6f
---

Design the smallest lifecycle change that makes stage occupancy distinct from completed stage work. The captain requested filing and ideation only; implementation requires a later design approval.

## Problem

A fresh email batch has status intake but no intake artifact. In intake -> triage -> done, the scheduler selects triage directly when intake is initial. Initial stages are exempt from completion proof. Non-initial ungated stages infer completion from committed reports. Gated work has a durable gate-attempt boundary. These meanings conflict.

The intended email flow is intake fetching a closed batch, triage classifying it and preparing human approval, then an execution hook reconciling approved changes before archive. Adding a queued stage hides the lifecycle inconsistency.

## Proposed direction for ideation

Prefer one ungated completion field, provisionally stage-complete: false|true. Status identifies the occupied stage. Initial selects the starting stage, not an exemption from work. An incomplete ungated stage dispatches itself. The FO validates the worker report and records completion through a guarded binary operation. Workers retain report ownership and never mutate lifecycle state.

Gate attempts remain authoritative for gated work. Do not add a competing gated completion flag or a new attempt registry. Determine the smallest representation and command surface; the field name and grammar are proposals, not approved implementation.

Actual entry or revision resets ungated completion atomically. Idempotent status stamps and dispatch retries must not reset completion. Specify how a new execution avoids accepting an old report. Decide absent-field semantics and migration explicitly so legacy backlog workflows do not suddenly execute their seed stages. Include initial gated seed behavior in that compatibility decision.

## Risk evidence

Current code: internal/status/entered_stage.go exempts initial, gated, and terminal stages from entered-stage completion checks. internal/status/format.go dispatchAnalysis projects the successor unless current-stage work is incomplete. internal/status/entered_stage_test.go explicitly tests initial-stage successor projection.

The email pilot observed missing intake artifacts with triage selected. Its no-op fixture covered approval-to-archive, not fresh-seed dispatch. Reproduce the initial-stage selection in a small isolated fixture before designing the change. Do not mutate live email workflows or call providers.

Relevant contracts: skills/first-officer/references/fo-dispatch-core.md, skills/commission/SKILL.md, docs/specs/gate-resolution-frontmatter-contract.md, and the state creation/mutation/dispatch implementations.

## Acceptance criteria

**AC-1 — Fresh batches execute intake before triage.**
Verified by: a fixture creates an incomplete intake seed and observes intake as the first dispatch target; triage becomes eligible only after validated intake completion. Restoring initial-stage successor projection must fail the check.

**AC-2 — Completion is durable and cannot leak across stage entry or revision.**
Verified by: fixture transitions exercise completion, restart, actual re-entry, and idempotent dispatch. A stale report or retained prior completion must not skip new work.

**AC-3 — Gated decision authority and legacy workflow behavior remain explicit.**
Verified by: fixtures cover gate approval/consume, ungated destinations, initial gated seeds, and the chosen absent-field migration rule. A duplicate completion authority or silent legacy seed dispatch must fail a check.

## Ideation deliverables

Recommend one minimal design, exact command and stored-field semantics, transition table, compatibility decision, affected contract wording, expected files/net LOC/tolerance, and a focused test plan. Name existing proof owners and the distinct failure each added test catches. Challenge whether the proposed field is necessary before expanding it. Stop at the ideation approval gate.

## Out of scope

Template-backed new command parameters, seed typo validation, email classification policy, generic external-proof registries, provider operations, and implementation during ideation.
