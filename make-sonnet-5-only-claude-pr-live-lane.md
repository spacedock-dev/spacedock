---
id: 0ytmjwn4ppg5en25z7vmna0p
title: Make Sonnet 5 the only Claude live lane on pull requests
status: backlog
source: Captain decision on 2026-08-07 after review of PR 626 and current Opus cost and failure evidence.
started:
completed:
verdict:
score: 0.85
worktree:
issue:
pr:
mod-block:
sprint:
---

Make ordinary pull-request CI run one Claude lane: Sonnet 5 with maximum effort. Keep Opus as a pre-release lane.

## Outcome

A pull request requires one Claude live approval and run. Pre-release validation retains explicit Opus coverage.

## Problem

The Claude live matrix runs Sonnet and Opus on ordinary pull requests. This adds approval delay, cost, and duplicate routine evidence.

PR 626 retained and documented this policy. The Opus matrix existed before that PR, but PR 626 treated it as the current normal lane.

## Proposed approach

Change only the Claude lane cadence. Select Sonnet 5 with maximum effort for ordinary pull requests. Select Opus only through the pre-release path.

Update the existing workflow checks and operating guide. Keep model-version strings in workflow configuration, not in the desired journey registry.

## Out of scope

- Do not change common journeys, fixtures, assertions, or target-specific TODO bindings.
- Do not change Codex or Pi cadence.
- Do not repair the Opus product gap owned by `a7`.
- Do not remove Opus from the desired-state registry.
- Do not add a new workflow when the existing workflow can express both cadences.

## Acceptance criteria

**AC-1 - Ordinary pull-request CI runs exactly one Claude live lane: Sonnet 5 with maximum effort.**
Verified by: the Runtime Live E2E workflow has one Claude pull-request selection, and its existing workflow contract check fails if Opus returns to that selection.

**AC-2 - Opus remains an explicit pre-release live lane.**
Verified by: the pre-release invocation selects Opus directly, while a pull-request event does not create an Opus approval or job.

**AC-3 - The desired journey and target registry remains unchanged.**
Verified by: the task changes no journey row, fixture binding, target requirement, or TODO owner in `docs/runtime-live-ci-registry.md`.

**AC-4 - Operators can identify the two cadences without reading workflow YAML.**
Verified by: `docs/runtime-live-ci.md` names Sonnet 5/max as the pull-request lane and Opus as the pre-release lane.

**AC-5 - The change is small and does not add another enforcement layer.**
Verified by: the implementation reuses the existing workflow and its existing contract checks. The expected surface is at most five files and approximately 40 additions and 40 deletions.

## Test plan

Run the focused existing Runtime Live E2E workflow contract checks. Then run `go test ./...` and `go test ./... -race`.

Use one workflow-dispatch dry inspection or equivalent event expansion to show that pull requests select Sonnet only and pre-release selects Opus. A paid Opus run is not required for this cadence change.
