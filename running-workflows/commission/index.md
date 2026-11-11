---
title: "Workflow design"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-11 05:18:33"
---

# Workflow design

Describe the work to your agent. It proposes stages, standards and decision points for you to adjust.

## Describe the work

In a session started with `spacedock claude`, use:

```
/spacedock:commission ship features through design, implementation, and review
```

## Review the design

You decide:

- **The work items:** what each task produces.
- **The stages:** how each task moves from an idea to completion.
- **The decision points:** where work pauses for you and where requested changes go.
- **The standards:** what agents must produce and check at each stage.

For each decision point, [gate content](../../concepts/gates-and-decisions/#gate-content) defines the context and evidence the agent presents to you.

## Refine the standards

The generated rules are a starting point. Ask your agent to refine them, or type `review stages` for a walkthrough of each stage.

## Run the workflow

After you accept the design and standards, the agent starts the first work items. [Your first workflow](../../get-started/first-workflow/) shows the full interaction.

## Sitemap

- [Overview](../../workflow/index.md)
- [Workflow operation](../operating/index.md)
- [Gates and decisions](../../concepts/gates-and-decisions/index.md)
- [Review standards](../../workflow/standards/index.md)
