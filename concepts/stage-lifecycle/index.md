---
title: "The stage lifecycle"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-06 02:11:11"
---

# The stage lifecycle

Each work item moves through the stages in your workflow. Each stage defines the work to do and the evidence required.

## A development workflow

```
flowchart LR
  design --> implementation --> review --> done
  review -. changes .-> implementation
  classDef gate stroke:#e3b04b,stroke-width:2.5px
  class design,review gate
```

In this example, the outlined stages pause for your decision: before implementation and after review.

Requested changes return to implementation. The reviewer checks the revised work before it returns to you.

## Stage settings

You can ask your agent to configure a stage to:

- Pause at a [gate](../gates-and-decisions/).
- Use a separate reviewer with fresh context.
- Return findings to the stage responsible for corrections.
- Limit concurrent work or use a specialist agent.

## Isolated work

Code stages can use separate Git worktrees so work items have their own branches and files.

**Related guides**

- [Workflow design](../../running-workflows/commission/)
- [Gates and decisions](../gates-and-decisions/)

## Sitemap

- [The operating model](../operating-model/index.md)
- [Workflows & entities](../workflows-and-entities/index.md)
