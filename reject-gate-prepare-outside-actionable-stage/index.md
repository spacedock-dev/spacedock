---
title: Reject gate prepare outside an actionable gated stage
status: backlog
source: "Pre-0.27 gate-machinery necessity audit, 2026-08-01: the real binary prepared and persisted a room while the ticket was in ungated implementation, then gate record refused it."
started:
completed:
verdict:
score: "0.95"
worktree:
issue:
pr:
sprint: durable-decisions
id: hq3d00mewqrys3s0z9pf27df
---

`gate prepare` must fail before mutation unless the ticket's current workflow stage is a gate that can accept a new attempt. The command currently checks only that the stage exists. A real invocation at the ungated `implementation` stage exited zero, added a gate attempt, and wrote a room that the later recorder correctly refused as non-actionable.

## Problem

Preparation is the first durable operation in the gate journey. Allowing it at an ungated or terminal stage manufactures authoritative-looking metadata that no supported decision flow can spend. It also forced a false `hold` record in the live `1w6` journey merely to retire an attempt that should never have existed.

## Proposed approach

Make the existing workflow-stage definition the authority. Resolve the ticket's current stage before allocating an attempt or room, require `gate: true`, and run the existing attempt-lifecycle eligibility checks before any write. Refusal must name the current stage and why it is not actionable. Do not add a repair path for rooms created by this unreleased prototype; delete or one-off-transform pilot state as needed.

## Common journey effect

- At `backlog`, `ideation`, or `validation` in this workflow, an eligible gate can be prepared normally.
- At `implementation` or `done`, preparation exits nonzero and leaves the ticket and filesystem byte-identical.
- A legitimate stale open attempt is retired with the existing `gate withdraw` lifecycle owned by `0m6`; invalid-stage preparation is not another withdrawal case.

## Out of scope

Do not change decision recording, terminal consumption, withdrawal semantics, provider presentation, or workflow stage declarations. Do not add compatibility parsing for invalid pilot attempts.

## Acceptance criteria

**AC-1 - Gate preparation cannot create an attempt or room outside an actionable gated stage.**
Verified by: a focused real-CLI fixture that tries the command at one gated stage, one ungated stage, and the terminal stage; the gated case succeeds, while each refusal has a stable diagnostic and a before/after byte and tree comparison proving zero mutation.

**AC-2 - Existing valid re-entry remains possible.**
Verified by: the smallest existing prepare fixture covering a retired prior attempt followed by a valid successor at the same gated stage; changing the guard to reject all historical records makes it fail.

## Test plan

Add focused package and CLI tests around the pre-write guard. Reuse the existing gate-prepare fixture builder; do not add a new live lane or transcript parser.
