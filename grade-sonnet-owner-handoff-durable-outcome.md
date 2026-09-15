---
title: Grade completed Sonnet owner handoff from durable outcomes
status: backlog
source: "PR #679 run 31640122995, Sonnet job 94260562364, artifact 9159847220: TestLiveCommonOwnedConflictOwnerHandoff completed the same-owner handoff, committed owner-handoff.marker, and filed a DONE report, but emitted conflict-owner-handoff-violation. The same class was visible in PR #673.\n2026-09-15 captain-approved scope/evidence: pre3 Claude run https://github.com/spacedock-dev/spacedock/actions/runs/34923154318, job 104235689031, artifact 10379211788, owned-conflict-owner-handoff. Marker commit 82649f4 has Captain <captain@example.test>; later worktree report commit a76e382 has the test identity. assertConflictOwnerHandoff uses git log -1 without binding to the marker-changing commit, so its wrong-author diagnosis is false. Bind author verification to the marker-changing commit. Clarify fixture assignment to permit the required worktree report while prohibiting authoritative entity edits; assert both boundaries separately. No owner-handoff product authority expansion. Acceptance: valid marker then differently authored report passes; wrong marker author and authority mutation fail; targeted local Claude owned-conflict-owner-handoff passes. Test/assertion files are unchanged pre2 to pre3; this is not established as a runtime regression. This evidence supplements the existing task, not a duplicate task."
score: 0.94
id: zvafbvhkdgjtbvay5sk3jzcc
---

Make the Sonnet owner-handoff journey grade the completed durable result instead of rejecting a valid host event shape.

## Problem

The worker handoff completed correctly, but the oracle rejected it. This makes required CI red after user value is already present.

## Proposed approach

Compare the exact artifact with the current oracle. Use durable branch, marker, report, and lifecycle outcomes. Do not add instrumentation or product behavior for observability.

## Out of scope

Do not change Opus or Pi bindings, dispatch behavior, or the owner-handoff product contract.

## Acceptance criteria

**AC-1 - The retained PR #679 Sonnet artifact passes the owner-handoff oracle.**
Verified by: replaying artifact 9159847220 accepts its exact completed durable outcome.

**AC-2 - Incomplete and wrong-owner handoffs still fail.**
Verified by: focused negative fixtures omit the marker, report, completion, or same-owner correlation.

## Test plan

Use retained artifact replay and focused negative controls. Run no paid live test while the evidence-only correction is under development.
