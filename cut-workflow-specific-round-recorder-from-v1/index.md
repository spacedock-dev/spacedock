---
title: Cut workflow-specific advisory round recording from v1
status: ideation
source: "Pre-0.27 gate-machinery necessity audit, 2026-08-01: generic gate code embeds the development workflow's Material/fixed/declined taxonomy, LOC estimate grammar, ensign role, and Feedback Cycles projection."
started:
completed:
verdict:
score: "0.95"
worktree:
issue:
pr:
sprint: durable-decisions
id: wjkhq0sktbbe3txx6jhnvcv2
gates:
    version: 1
    current:
        gate: gate:wjkhq0sktbbe3txx6jhnvcv2:backlog
    records:
        - id: gate:wjkhq0sktbbe3txx6jhnvcv2:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:wjkhq0sktbbe3txx6jhnvcv2-backlog-1
              briefing:
                id: briefing:wjkhq0sktbbe3txx6jhnvcv2:backlog:attempt-1:revision-1
                digest: sha256:45e247730eceb6bb9fb0e68f806630bab3cdd2df47830efc5f4f9e40792351f8
                digest-domain: canonical-bytes
                request-digest: sha256:13584896ffe679ba3def35d84099ea4402ac80d53f7836c89be8fa535efc32af
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:wjkhq0sktbbe3txx6jhnvcv2:backlog:1
                briefing: briefing:wjkhq0sktbbe3txx6jhnvcv2:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-01T13:50:02.905125Z"
                decision: approve
                reason: Captain approved dispatching this durable-decisions ideation lane in parallel with hq, nth, and jc.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

Remove the public `gate record --round` and `gate validate --round` surface from the stable v1 cut. The present implementation is coherent only for this repository's development review policy, so retaining it would make workflow-specific semantics part of the generic gate storage contract.

## Problem

`internal/gates` recognizes development-only classes such as `correct-but-disproportionate`, parses a files/LOC/estimate/percentage/AC sentence, distinguishes a fixed or declined material finding, assumes an `actor:ensign` worker, and edits a `### Feedback Cycles` body section. Those are workflow rules, not durable-decision primitives. A prose-only rename cannot repair the jurisdiction error.

## Proposed approach

For 0.27 v1, delete the public round command paths, their development-only model and parser surface, and any spec claim that generic recording decides whether a Material finding is fixed. Preserve ordinary binding gate preparation, Resolution recording, validation, withdrawal, eligibility, and consumption. Keep the existing `workflow-neutral-advisory-round-recorder` ticket as the later reintroduction path: it may return only after its storage graph is structural and the active workflow owns classification and projection.

## Streamlined common scenario

An in-stage review remains a workflow operation. The reviewer supplies findings, the worker proposes dispositions, the First Officer authorizes work ownership, and the workflow retains whatever report it requires. The v1 gate binary does not parse that policy. When a later generic recorder returns, it may retain immutable Briefing, Annotation, and advisory Resolution relationships without interpreting their bodies.

## Out of scope

Do not design the replacement recorder in this ticket. Do not add a compatibility alias, hidden legacy mode, or schema version. Do not remove binding gate decisions.

## Acceptance criteria

**AC-1 - The stable command surface exposes no workflow-specific advisory round operation.**
Verified by: CLI help and routing tests that reject the removed round grammar and keep every binding gate verb green.

**AC-2 - Generic gate packages contain no development finding-policy parser or body projector.**
Verified by: focused package tests plus a source inventory tied to the removed exported behavior; retaining the Material/fixed/declined classifier, LOC grammar, ensign requirement, or Feedback Cycles mutation fails the criterion.

**AC-3 - The development workflow still states who classifies and routes findings without claiming the gate binary records them.**
Verified by: workflow contract review and existing behavioral workflow tests; a stage instruction that tells an ensign to invoke the removed round command fails.

## Test plan

Delete tests whose only value is the removed feature; preserve or move generic graph fixtures only if a currently shipped gate behavior consumes them. Run the full Go and race suites and the relevant existing workflow live scenario. Do not build a replacement harness.
