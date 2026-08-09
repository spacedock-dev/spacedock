---
id: qz0ap96nt5k93tgbsphq9ahy
title: Keep journey-delta reporting green across failed-job reruns
status: backlog
source: PR #639 Runtime Live E2E attempts 2 and 3 on 2026-08-08.
started:
completed:
verdict:
score: 0.75
worktree:
issue:
pr:
mod-block:
---

A failed-job rerun can reuse successful live jobs from an earlier attempt. The metrics job then finds Codex artifact metadata, but artifact extraction fails.

All test jobs pass in this state. However, the optional metrics job keeps the complete workflow red.

## Problem

`journey-delta-comment` treats each current-run artifact as required. GitHub failed the Codex artifact download after five retries in two consecutive rerun attempts.

This failure hides the correct test result. It also prevents a clean merge for a change that passed all required test lanes.

## Proposed approach

Make the journey-delta comment best-effort at the artifact boundary. Preserve the comment when both metric artifacts are available.

If an artifact is unavailable, emit a clear warning and complete the reporting job successfully. Do not publish a partial or misleading comment.

Ideation must examine whether GitHub can select artifacts from a prior successful attempt. Prefer that recovery when it is reliable and small.

## Out of scope

- Do not change the live journey assertions.
- Do not add a live lane.
- Do not rerun a passed live lane only to produce metrics.
- Do not change the journey-cost format or comment layout.

## Acceptance criteria

**AC-1 (VALUE) - An unavailable metrics artifact cannot make a pull-request workflow red after all required test jobs pass.**
Verified by: a deterministic workflow exercise makes one artifact unavailable and observes a successful reporting job with a warning. Removing the nonfatal path makes this exercise fail.

**AC-2 - A normal first attempt still posts or updates one journey-delta comment when both artifacts are available.**
Verified by: the existing comment path receives both fixture artifacts and produces one complete comment. Removing either artifact input makes this exercise fail.

**AC-3 - A rerun never publishes a partial journey-delta comment.**
Verified by: fixtures cover one missing Claude artifact and one missing Codex artifact. Each case emits a warning and produces no comment.

**AC-4 - The fix does not add or rerun a live test job.**
Verified by: the workflow job graph remains offline, one Claude lane, one Codex lane, and the reporting job. Adding another producer makes the graph check fail.

**AC-5 - Operators can distinguish a test failure from an optional metrics failure.**
Verified by: the missing-artifact exercise exits successfully and emits the artifact name in a GitHub warning. Removing the artifact name makes the exercise fail.

## Test plan

Use deterministic fixtures for successful and failed artifact downloads. Exercise the reporting step without a live runtime subscription.

Run the existing workflow checks, the affected Go package, `go test ./...`, and `go test ./... -race`.

Because this task changes CI machinery, validation must run a detached adversarial audit. The audit must make one artifact unavailable and observe the value in AC-1.
