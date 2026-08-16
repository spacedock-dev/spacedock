---
id: 18963egcskzxaje6b5vnas3q
title: Simplify the feedback-rejection flow to five steps or fewer
status: ideation
source: "Captain principle 2026-08-16: anything more than 4 or 5 steps is a footgun. Live evidence: both rejection-flow failure mechanisms are step-count casualties - codex drops step 8's bundled tail (completion gap, ErrNoGateRecord); claude's oracle splits across the step-6/step-8 two-call publish (entries=2 vs 4, wrong round id)."
started:
completed:
verdict:
score: "0.85"
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:18963egcskzxaje6b5vnas3q:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:18963egcskzxaje6b5vnas3q-backlog-1
              briefing:
                id: briefing:18963egcskzxaje6b5vnas3q:backlog:attempt-1:revision-1
                digest: sha256:db36beb0817ac7684245ed81ade8d77f4e20fb9a3e8aefb91acfa479f41497e0
                request-digest: sha256:a8a79d493bc4736f5cae81387d4165fabc076cd3e6487a96d03ea5770a99db73
                room-ref: ./simplify-feedback-rejection-flow/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:18963egcskzxaje6b5vnas3q:backlog:1
                briefing: briefing:18963egcskzxaje6b5vnas3q:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-16T04:00:54.684905Z"
                decision: approve
                reason: 'Captain directive 2026-08-16: dispatch 18'
              application:
                target-stage: ideation
                state: consumed
---

Redesign skills/feedback-rejection-flow/SKILL.md from 8 steps to at most 5, with the single-publish shape the diagnostic evidence recommends: accumulate the round, publish ONCE under one round id, and make the gate re-entry a single unbundled step. Every step becomes one action with one completion condition; no step may bundle a tail an FO can drop.

Coordination, binding: hz (repair-codex-rejection-round-recording, in implementation) owns the binary side (inline state commit) and grading honesty; its reconciliation of the two diagnosis accounts determines whether the oracle repairs land there or the oracle re-anchors to this redesign's single-publish shape. Whichever lands second reconciles. The live rejection-flow journey and its oracles are the proof surface for both.

## Problem

{Ideation fills this in, from both diagnosis reports' quoted streams.}

## Proposed approach

{Ideation fills this in: the <=5-step flow, the single-publish round shape, and the oracle alignment.}

## Out of scope

The binary state-commit fix and grading-honesty repairs (hz owns them). The workflow-defined cycle-line grammar.

## Expected surface and tolerance

Estimate net LOC change: near zero or negative on the skill; oracle-side test adjustments in step with the new shape.

## Acceptance criteria

**AC-1 - The shipped flow has at most 5 steps, each one action with one completion condition.**
Verified by: the skill's numbered-step count and a per-step single-action review at the gate; falsifying edit - re-bundling any tail.

**AC-2 - One round publish per rejection cycle: the recorder is invoked exactly once per round with a stable round id, and the live oracle asserts that single call.**
Verified by: the updated oracle plus a targeted live run where the recorder invocation count is read from the stream.

**AC-3 - The live rejection-flow journey passes on both hosts under the new flow.**
Verified by: targeted local runs both runtimes, then the full-matrix stack run.

## Test plan

Skill rewrite plus oracle alignment; targeted local rejection-flow runs as the repair loop; contractlint reference closure; no new standing harness.
