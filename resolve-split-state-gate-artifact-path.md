---
title: Resolve split-state gate artifact paths before preparation
status: backlog
source: "PR #679 run 31640122995, Codex job 94260562369, artifact 9158783630: TestLiveCommonGateGuardrail found committed review files in the split state checkout, then passed a state-relative path to gate prepare, which resolved it from the workflow root and failed."
score: 0.98
id: mvmpzgqxyb32t3b3vdw0x0h1
---

Make a split-root First Officer pass the exact committed artifact path that gate preparation can resolve.

## Problem

A valid committed gate review exists, but the First Officer can select the wrong path base. The required Codex lane then stops before later journeys.

## Proposed approach

Diagnose the path contract first. Prefer one canonical resolver-owned path over model-selected path reconstruction. Do not add global hooks or test-only product machinery.

## Out of scope

Do not change DVD or v8 behavior, XFAIL policy, or unrelated gate semantics.

## Acceptance criteria

**AC-1 - Split-root gate preparation accepts the committed selected artifact.**
Verified by: the exact gate-guardrail fixture reaches an open prepared gate; a path based at the wrong root fails.

**AC-2 - The solution uses one canonical path contract.**
Verified by: focused tests cover split-root and single-root resolution without host-specific parsing.

## Test plan

Add focused path-resolution tests, then run the exact local Codex gate-guardrail target once on final bytes. Run full and race suites once if Go behavior changes.
