---
title: Resolve split-state gate artifact paths before preparation
status: backlog
source: "PR #679 run 31640122995, Codex job 94260562369, artifact 9158783630: TestLiveCommonGateGuardrail found committed review files in the split state checkout, then passed a state-relative path to gate prepare, which resolved it from the workflow root and failed."
score: 0.98
id: mvmpzgqxyb32t3b3vdw0x0h1
gates:
    version: 1
    records:
        - id: gate:mvmpzgqxyb32t3b3vdw0x0h1:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:mvmpzgqxyb32t3b3vdw0x0h1-backlog-1
              briefing:
                id: briefing:mvmpzgqxyb32t3b3vdw0x0h1:backlog:attempt-1:revision-1
                digest: sha256:449b528bc4625bcc07c6f1463e438d550123d3e19fd4ba1c0821e80c28d13d7a
                request-digest: sha256:d855623326f6ca92f7fffb5b906151ab29b31b736198878fa45a77ae514fe3a6
                room-ref: ./resolve-split-state-gate-artifact-path/review/backlog/briefing-1
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
