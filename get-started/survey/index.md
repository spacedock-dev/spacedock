---
title: "Project survey"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-11 05:14:43"
---

# Project survey

Survey reads your local agent history to identify recurring work and decisions that need your attention.

## Survey the project

In Claude Code, run:

```
/spacedock:survey
```

Survey uses [agentsview](https://agentsview.io/), a tool for reading agent session history. It asks before installing a missing dependency.

The survey is read-only. It reports:

- Repeated work and manual interventions.
- Current workstreams.
- Unresolved decisions.

## Use the findings

Your agent can propose a workflow from the findings. You decide whether to create it and which decisions remain yours.

To describe a workflow yourself, see [Your first workflow](../first-workflow/).

## Sitemap

- [Your first workflow](../first-workflow/index.md)
- [Session debrief](../../running-workflows/debrief/index.md)
