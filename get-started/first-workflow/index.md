---
title: "Your first workflow"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-11 08:30:05"
---

# Your first workflow

[Install Spacedock](../install/), then describe a repeatable piece of work and decide where it needs your judgment.

## Describe the work

```
spacedock claude "/spacedock:commission ship features through design, implementation, and review"
```

The agent asks about your work and proposes a workflow. For example:

> Each task moves through design, implementation and review.
>
> You approve the design before implementation. A separate agent reviews the result before your final decision.
>
> Requested changes return to implementation.

## Set the standards

Adjust the proposed stages, decision points and quality standards. Type `review stages` to walk through the standards with your agent.

For each decision point, specify the context and evidence you need. This is the workflow's [gate content](../../concepts/gates-and-decisions/#gate-content).

## Run the first tasks

Choose the first work items with your agent. After you accept the design and standards, the agent starts the work.

At a gate, the agent presents the evidence and asks for your decision. Approve, request changes or hold the work.

## Resume later

Run [`/spacedock:debrief`](../../running-workflows/debrief/) before you end the session.

To resume:

```
spacedock claude
```

The agent reads the workflow records and picks up the work.

**Related guides**

- [Workflow operation](../../running-workflows/operating/)
- [Project survey](../survey/): start from your existing agent history.

## Sitemap

- [Project survey](../survey/index.md)
- [Session debrief](../../running-workflows/debrief/index.md)
