---
id: zrcmxzx5c7arew7r8afmxfbw
title: Non-sandboxed launch uses Claude auto-mode (and Codex equivalent)
status: ideation
source: captain (2026-06-12)
started: 2026-06-13T04:31:09Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0201-post-flip-release-model
group: ux-cleanup
sprint-readiness: ready
---

When the launch is not sandboxed, `spacedock claude` should launch with Claude Code's auto-mode, and `spacedock codex` should use the equivalent Codex setting. Captain ask, filed verbatim; ideation should confirm the intended rationale and the exact mode mapping with the captain.

Couples to `startup-sandbox-status` (gja5htstcgjxydcz5h2051wc): both need the same sandbox-state detection; this task consumes it to choose the launch mode, that task surfaces it to the operator.

## Problem

{What is broken or missing, and why it matters. Ideation fills this in. Seed framing: launch permission posture is currently independent of sandbox state; ideation should confirm current launcher flag behavior and define the auto-mode mapping per host (Claude / Codex), plus whether the operator can override.}

## Proposed approach

{How the task intends to solve the problem. Ideation fills this in.}

## Out of scope

{What this task deliberately does not cover, so the boundary is explicit.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - {End-state property.}**
Verified by: {test name / command output or exit code / file the change produces / resulting on-disk state — something outside this task body that a future reader can reproduce and that can fail.}

## Test plan

{What verifies the implementation, estimated cost/complexity, and whether fixture, CLI, or live workflow tests are needed. Note: this changes user-visible launch behavior, so ideation owes a doc diff per the workflow's ideation output rules.}
