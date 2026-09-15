---
title: Dispatch a gated successor before presenting its gate
status: backlog
source: Captain approved separate headless contract fix on 2026-09-15
started:
completed:
verdict:
score: 0.95
worktree:
issue:
pr:
mod-block:
id: vcfwr5ptn6pedv32wq7cbsy6
---

Remove the contradictory completion shortcut so a headless first officer runs a gated successor before preparing its review.

## Problem

The corrected pre-gate fixture dispatched and completed implementation, then attempted gate preparation while still at implementation. The loaded shared dispatch contract says to stop when the next stage is gated, conflicting with its Gate successor guard requiring that stage's dispatch and completed report first.

## Proposed approach

Remove the next-stage gate stop exception in skills/first-officer/references/fo-dispatch-core.md. Distinguish completed-stage gate routing from entering and dispatching a gated successor. Preserve the existing Gate successor guard, authority, freshness, and approval rules. Captain approved this separate stack layer; no new CLI behavior.

## Risk evidence

/tmp/spacedock-stack-headless-live/finding.md records exact failed candidate da50d61d6 and loaded conflicting text. TestLiveCommonDefaultHeadlessGateStop failed in 122.09 seconds after one implementation worker completed; no validation transition/worker occurred. GateGuardrail control passed. Wording causality is a hypothesis that requires the same targeted live journey after the correction.

## Out of scope

New state metadata, scheduler or CLI changes, changing assertions to accept an unfinished validation stage, and same-stage revision issue 792.

## Expected surface and tolerance

Estimate net LOC change: +4, across 2 files, approximately 8 insertions and 4 deletions. Tolerance net +35 and 3 files. Shared dispatch contract plus existing skill smoke or behavioral fixture only; no new runner. Layer branches from committed stack tip e52ece67abdb282ab061c41ee897a8a341b22954.

## Acceptance criteria

**AC-1 — Headless completion enters and completes the gated successor before presenting its gate.**
Verified by: targeted local Codex TestLiveCommonDefaultHeadlessGateStop observes implementation completion, validation transition/worker report, then one open gate, with no approval or consumption. Restoring the shortcut is the falsifying intervention; do not assert causal certainty from prose presence.

**AC-2 — A gate-ready workflow still stops without spending captain authority.**
Verified by: targeted local Codex TestLiveCommonGateGuardrail remains passing; no successor dispatch or consumption occurs without approval.

**AC-3 — The correction preserves existing worker authority and delivery semantics.**
Verified by: focused existing dispatch/skill smoke and durable lifecycle controls; code diff contains no CLI/state schema or grader weakening. Current tip full normal/race results are recorded before completion, distinguishing any reproduced baseline defect.

## Test plan

Use the existing default-headless and gate-guardrail live proof owners, existing skill smoke before editing, and one serialized local Codex run per affected journey at exact final tip. Keep lower four layers intact. Independent validation verifies the full sequence and frozen authority, including source fixture layer's previously failed AC-2. Run required normal/race and gofmt at final tip; no CI until every layer is locally verified and stack ready.

### Feedback Cycles
