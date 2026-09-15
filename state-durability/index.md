---
title: Make state transitions durable
status: ideation
source: Captain-approved combined durability work, 2026-09-15
started:
completed:
verdict:
score: 0.95
worktree:
issue: spacedock-dev/spacedock#689
pr:
mod-block:
id: 3tzfv0rrctdb066xt942zjbd
---

Prevent false completion and silent loss of workflow state. One task covers GitHub issues #689, #790 and #630 with separate behavioral acceptance criteria.

## Problem

#689: a local-merge workflow without a registered merge hook can archive before manual Git delivery succeeds. A subsequent conflict leaves false done state.
#790: status --archive moves a retired task, but state commit refuses its dirty archived scope and can exit zero. The retirement remains uncommitted.
#630: a missing split-root state checkout looks like a healthy empty workflow; new can create tasks in an untracked location before state commit fails.

## Proposed approach

Fix the three transaction boundaries in existing command owners. Delivery failure must preserve pending approval and active state. Retirement must have a supported durable commit path. Missing state storage must be distinguishable and unsafe filing refused. Reuse existing state publication, guards and archive helpers; do not introduce a new transaction framework or workflow stage. Keep retirement distinct from successful delivery.

The captain explicitly chose one task with three acceptance criteria rather than three separate workflow cycles. Split only on demonstrated incompatible designs or dependencies, and return that evidence for decision.

## Risk evidence

Issues: https://github.com/spacedock-dev/spacedock/issues/689 ; https://github.com/spacedock-dev/spacedock/issues/790 ; https://github.com/spacedock-dev/spacedock/issues/630 . Reproduce each against current main before proposing fixes. A stale historical report alone does not establish a current defect.

## Out of scope

External entity path repair, worker fencing, evidence scanner changes, bulk workflow refits, new state services, and unrelated cleanup. No code push or CI during ideation.

## Expected surface and tolerance

Ideation must name exact owners, net-line/file estimate and tolerance after the three reproductions. Use current main 438053493838dc70c9478b3d991309d566783e85, not the stale root checkout.

## Acceptance criteria

**AC-1 — Failed local delivery leaves the task active with terminal approval unspent.**
Verified by: real temporary Git repositories with a conflicting merge and a successful merge. Assert pending/consumed approval, active/archive location and delivered commits. Premature archive must fail the conflict case; successful delivery must finalize exactly once.

**AC-2 — Retirement has a documented command path that commits the full archive move.**
Verified by: flat and folder-form entities, including sibling artifacts, in split-root local-only and remote-backed state. Assert tracked archive, removed active paths, publication result and nonzero failure behavior. An uncommitted move reported as durable success must fail.

**AC-3 — Missing state storage cannot masquerade as an empty initialized workflow or accept unsafe filing.**
Verified by: missing, invalid and valid-empty state checkouts. Assert distinct status diagnostics and machine output, byte-clean new refusal, and successful filing after supported initialization. A silently created ignored directory must fail.

## Test plan

Start with the three cheapest real CLI/Git reproductions. Add focused regression tests in existing owners before implementation. Exercise failure and recovery without adding custom orchestration. Run normal/race suites and formatting; use a detached adversarial audit for status/guard paths. Determine live lane requirements from the actual diff; do not trigger costly CI before local validation and stack readiness.

### Feedback Cycles

