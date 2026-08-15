---
id: 3n6qcjttj261wjcdkvc71k6z
title: User-facing journey for the 0.27 gate/resolution benefit
status: backlog
source: "Captain directive, 2026-08-14: for the 0.27 gate/resolution rollout, publish a user-facing journey that describes the benefit of this major release"
started:
completed:
verdict:
score:
worktree:
issue:
---

Ship an executable gate/resolution journey plus the user-facing story it backs, for the 0.27 rollout. The journey exercises the full flow — `gate prepare` briefing, captain decision (`gate record` approve/revise/hold), a feedback round routed back, `gate consume` into the durable Resolution record, `merge guard` to done — and captures its artifacts. The user-facing narrative is derived from that run: every benefit claim cites a captured artifact, so a separate check can confirm the story. The end-value framing stays: the captain keeps decision-quality control over agent work, with a durable audit trail of every decision.

Captain ruling (2026-08-14, verbatim): "DOC-ONLY TASK IS BANNED." The narrative must ride on a real, checkable change (README ideation criteria and Proof policy). This seed supersedes the original doc-only shape, which was non-qualifying.

Existing machinery for ideation to build on: the live journeys under internal/ensigncycle, internal/journeymetrics, the journey-delta release tooling, and internal/ensigncycle/recorded_gate_lifecycle_test.go. Dogfood proof points from the 2026-08-14 audit remain available as material: 392 decisions recorded, 282 resolutions consumed, 212 feedback-cycle entries, 319 entities driven to done.

## Problem

{What is broken or missing, and why it matters. Ideation fills this in.}

## Proposed approach

{Ideation fills this in: extend an existing live/fixture journey or add a scripted demo journey; where the narrative lands, following the README rule that the doc diff rides in the behavior task; and the walkthrough scenario.}

## Out of scope

{What this task deliberately does not cover, so the boundary is explicit.}

## Expected surface and tolerance

Estimate net LOC change: {+NNN}, across {M} files. Declare observable semantics touched; none expected beyond a new journey fixture/test plus its docs.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. Ideation refines; the seed anchors:

**AC-1 - An executable journey exercises prepare -> record (approve and revise paths) -> consume -> merge guard and leaves the durable decision record as resulting on-disk state.**
Verified by: the journey test/fixture exit code plus the resulting gate record files — behavior exercised, not described.

**AC-2 - The user-facing story exists and every benefit claim in it cites an artifact captured from that run.**
Verified by: docs build green plus a one-off citation-to-artifact check recorded in the validation report (one-off evidence, never a committed prose-grep, per Proof policy).

## Test plan

{Ideation details: journey fixture vs live smoke, estimated cost, and the docs check. No standalone prose-review-only verification — banned.}
