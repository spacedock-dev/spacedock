---
title: "The stage lifecycle"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-25 04:42:39"
---

# The stage lifecycle

An entity moves through an ordered chain of stages that the workflow defines, and each stage declares the work it owns and the proof it must produce. The first officer advances an entity stage by stage, pausing at the gates you declared.

## A typical dev workflow

```
flowchart LR
  backlog --> ideation --> implementation --> validation --> done
  validation -. rejected .-> implementation
  classDef gate stroke:#e3b04b,stroke-width:2.5px
  class ideation,validation gate
```

The gold-bordered stages are gates, your calls: before code is written, and before the result ships. Read the chain as a pipeline: each stage takes the prior stage's output as its input, and the bar rises from "is this clear?" to "is this proven?".

The property that matters most is `feedback-to`: rejected work bounces back to the stage that owns the fix (a rejected validation returns to implementation), not to the reviewer that flagged it.

## What a stage declares

A stage can pause at a [gate](../gates-and-decisions/) for your decision, run its work in an isolated worktree, demand a reviewer with no access to the maker's reasoning, route rejected work back to an earlier stage, cap how many items it holds at once, and hand its work to a specialist worker. All of it is declared in the workflow README; ask the first officer to set up or change any of it.

Beyond the declarations, the prose of each stage's section in the README is the stage definition. What you write there is exactly what the worker receives as its assignment.

## Fresh context at validation

Validation declares `fresh: true`, and `fresh: true` means adversarial review: the reviewer is never the maker, arrives without the implementer's reasoning in its context, sees only the entity body and the deliverable, and pushes back on thin evidence. This is the mechanism behind the README's claim that "the agent doesn't get to judge its own work." High-stakes changes can also get a detached, out-of-workflow pass; [Reviews beyond validation](../gates-and-decisions/#reviews-beyond-validation) covers the difference.

When validation recommends `REJECTED`, `feedback-to: implementation` routes the concrete finding back to the implementation stage for rework rather than closing the entity. The entity re-enters implementation, the finding is addressed, and a fresh validator checks it again. A hard cap on feedback cycles prevents an endless bounce; on the third cycle the first officer escalates to you.

## Isolated worktrees

When a stage declares `worktree: true`, everything from that stage onward happens in an isolated worktree: the work and its commits stay there, concurrent items never collide with each other or with you, and at the terminal stage the branch is merged back and the worktree cleaned up.

## Where to go next

- [The operating model](../operating-model/) for who does what: you, the orchestrator, the workers.
- [Gates and decisions](../gates-and-decisions/) to see exactly what you decide at a stage boundary and on what evidence.
- The [frontmatter contract](../../reference/frontmatter-contract/) for the fields these stages write.

## Sitemap

- [The operating model](../operating-model/index.md)
- [Workflows & entities](../workflows-and-entities/index.md)
- [Gates & decisions](../gates-and-decisions/index.md)
