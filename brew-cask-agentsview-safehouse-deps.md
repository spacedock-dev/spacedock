---
id: 8pwjdj4ngx9dbnbynxgsagq0
title: Brew cask depends on agentsview and safehouse
status: backlog
source: captain (2026-06-12)
started:
completed:
verdict:
score:
worktree:
issue:
sprint: 0201-post-flip-release-model
group: ux-cleanup
sprint-readiness: ready
---

The Spacedock Homebrew cask should declare `agentsview` and `safehouse` as dependencies, so a brew install brings them in alongside the launcher.

## Problem

{What is broken or missing, and why it matters. Ideation fills this in. Seed framing: ideation should confirm where the cask/formula lives (the homebrew tap), whether agentsview and safehouse are casks or formulae and in which tap, and the right dependency declaration (`depends_on`).}

## Proposed approach

{How the task intends to solve the problem. Ideation fills this in.}

## Out of scope

{What this task deliberately does not cover, so the boundary is explicit.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - {End-state property.}**
Verified by: {test name / command output or exit code / file the change produces / resulting on-disk state — something outside this task body that a future reader can reproduce and that can fail.}

## Test plan

{What verifies the implementation, estimated cost/complexity, and whether fixture, CLI, or live workflow tests are needed. Note: install behavior is user-visible — ideation owes a doc diff (install docs) per the workflow's ideation output rules.}
