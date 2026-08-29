---
title: "Operate a workflow"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-29 00:32:28"
---

# Operate a workflow

You operate a workflow by talking to the first officer: launch a session, let it move everything that is ready, and decide when work reaches a gate.

## The session loop

Start a session. Name a task if you have one in mind, or say nothing and it picks up where the last session left off:

```
spacedock claude
```

The first officer reads the workflow state and dispatches an ensign for every entity ready for its next stage. Completed stages flow forward on their own: when the next stage has no gate, the first officer advances the entity and dispatches again without waiting for you. It returns to you for four things only: a gate needs your decision, a finished entity is ready to close, something is blocked, or nothing is left to dispatch.

In between, ask it whatever you want to know: "what's ready?", "where is the rate-limiting task?", "what's still in review?". The first officer queries the workflow state and shows you; there are no status commands to learn.

## When a gate arrives

A gate stops the loop and hands you the call: the first officer presents the review and waits. You make one of [the three calls](../../concepts/gates-and-decisions/#the-three-calls): approve, redo with feedback, or reject.

Approving the last stage closes the entity: the first officer records the merge and the verdict, archives the entity file, removes the worktree, and stands the workers down. The loop then continues with whatever is ready next.

## Delegating the loop

The loop needs you less as the workflow matures. Rejections already run without you: findings bounce back to the stage that owns the fix and the reviewer re-runs; [a rejection reaches your desk](../../concepts/gates-and-decisions/#rejections) only after three failed rounds. When the bar you set is sharp enough to trust, [hand over the conn](../../concepts/gates-and-decisions/) and let the first officer drive whole tasks to done with auto-approval.

## Ending a session

Stop whenever you want; every entity's state is in the workflow files, so the next session resumes where this one ended. Before you stop, run [`/spacedock:debrief`](../debrief/) to record the session for the next one.

## Sitemap

- [Commission a workflow](../commission/index.md)
- [Debrief a session](../debrief/index.md)
