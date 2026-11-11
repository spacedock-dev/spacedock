---
title: "Workflow operation"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-11 05:34:48"
---

# Workflow operation

Talk to your agent to add work, check progress and make decisions.

## Start a session

```
spacedock claude
```

The agent reads the workflow records and continues work that is ready. It coordinates workers and brings you decisions or blockers.

## Check progress

Ask questions such as:

- What needs my decision?
- What is blocked?
- What is ready next?

You can add work or change priorities in the same conversation.

## Decide at a gate

The agent presents the evidence specified by the workflow's [gate content](../../concepts/gates-and-decisions/#gate-content). You can approve, request changes or hold the work.

Your decision and feedback become part of the work item's record. The agent continues or routes the requested changes according to the workflow.

## End a session

Run [`/spacedock:debrief`](../debrief/) to record the completed work, decisions and follow-up items.

The next session reads the workflow records to resume the work.

## Sitemap

- [Overview](../../workflow/index.md)
- [Workflow design](../commission/index.md)
- [Gates and decisions](../../concepts/gates-and-decisions/index.md)
- [Review standards](../../workflow/standards/index.md)
