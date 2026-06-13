---
id: gja5htstcgjxydcz5h2051wc
title: Show sandbox state on startup
status: ideation
source: captain (2026-06-12, UX improvement)
started: 2026-06-13T04:08:48Z
completed:
verdict:
score:
worktree:
issue:
---

UX improvement: on startup, Spacedock should also show whether sandboxing is enabled — distinguishing enabled / available-but-not-enabled / unavailable — so the operator knows the execution-isolation posture before dispatching work.

Captain amendment (2026-06-12): the same sandbox info should also appear in `--version`. In addition, `--version` should show each detected runtime (claude / codex / pi) with its install and spacedock-enablement status, e.g.:

```
codex: installed
claude: installed, spacedock enabled
pi: not installed
```

## Problem

{What is broken or missing, and why it matters. Ideation fills this in. Seed framing: the startup surface reports mods, ID style, orphans, PR state, and dispatchables, but says nothing about whether the session is sandboxed. The operator has no at-a-glance signal of the isolation posture.}

## Proposed approach

{How the task intends to solve the problem. Ideation fills this in. Ideation should pin down which startup surface(s) this lands on (launcher banner, `status --boot`, or both), how sandbox state is detected per host, and how `--version` detects each runtime (claude / codex / pi) plus whether spacedock is enabled in it.}

## Out of scope

{What this task deliberately does not cover, so the boundary is explicit.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - {End-state property.}**
Verified by: {test name / command output or exit code / file the change produces / resulting on-disk state — something outside this task body that a future reader can reproduce and that can fail.}

## Test plan

{What verifies the implementation, estimated cost/complexity, and whether fixture, CLI, or live workflow tests are needed.}
