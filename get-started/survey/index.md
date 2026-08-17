---
title: "Survey your project"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-17 14:56:14"
---

# Survey your project

When you use the `spacedock:survey` skill, it looks at your existing agent conversation logs on local disk (through [agentsview](https://agentsview.io/), an open source session-history tool). It is read-only; if `agentsview` is missing, it asks before installing it.

## What it reports

- **The repeated manual behavior**: the loop you have been driving by hand, run after run, without naming it.
- **The steering and interruptions observed**: how often your agents needed you to step in, and for what.
- **Your current workstreams**: what is in flight, clustered into tracks.
- **What is still undecided**: forks raised but never resolved.

## Turn the report into a workflow

Survey ends with an offer, not an action. It can turn what it found into a Spacedock [workflow](../../concepts/workflows-and-entities/): the repeated loop becomes the stages, the workstreams become the work items, and the undecided forks become [approval gates](../../concepts/gates-and-decisions/). Nothing changes in your project until you say yes; on a no, the survey stands on its own as orientation.

The offer matches your work:

- **Routine loops** (issue → worktree → PR) get automation: a workflow that gates the crucial decisions and lets the agent drive between gates.
- **Exploration** (creative or design work where your steering is the point) gets book-keeping: structure for the parallel threads and their state. There is no automate-the-human-out pitch; the involvement is the work.

To define a workflow yourself instead, see [your first workflow](../first-workflow/).

## Sitemap

- [Install](../install/index.md)
- [Your first workflow](../first-workflow/index.md)
