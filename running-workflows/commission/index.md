---
title: "Commission a workflow"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-18 23:08:54"
---

# Commission a workflow

A workflow is designed in conversation: you make four decisions and commission derives the rest. [Your first workflow](../../get-started/first-workflow/) shows the whole flow, including the design summary and pilot run.

Invoke it from a session started with `spacedock claude`. You can pass the mission inline:

```
/spacedock:commission product idea to simulated customer interview
```

Text after the command name becomes the workflow mission; with no argument, commission greets you and asks.

## The four things you decide

1. **The mission and what each work item is.** The description you give becomes the label the workflow uses everywhere: "a design idea" makes it a workflow of ideas.
1. **The stages.** Commission matches your mission to a workflow archetype (development: ship code through review; experiment: test a hypothesis against evidence; refinement: iterate on an artifact until it is good enough) and proposes the stage list from it. Confirm, adjust, or trim; it pushes back on redundant names (`awaiting_validation` is just `validation`).
1. **The gated stages.** Where the workflow pauses for your decision: by default one gate before the terminal stage, each with a rejection target stated in plain language ("If you reject at `review`, it goes back to `draft` for revision").
1. **The per-stage quality bar.** What "good" means for each stage in the generated README: what to produce, the bar to meet, the anti-patterns to avoid. Starting prose, not commitments.

Everything else is derived or asked with a recommendation attached: the directory under `docs/` where everything lands as plain text, how entities are identified, how rejections route. Stages that write code give each entity an isolated worktree, so your main checkout stays clean; if you ship through PR review, commission offers the [pr-merge mod](../../advanced/mods-and-standing-teammates/), which manages the PR lifecycle so merging never needs to be a stage. The design summary shows it all; ask about the tradeoffs before you accept.

## Two ways to tighten the README

The generated per-stage rules are best-guesses; [tighten them before any work runs](../../concepts/workflows-and-entities/#the-readme-is-where-you-set-the-rules). Either:

- open the README and edit the bullets under each stage heading directly, or
- type `review stages` to have commission walk you through each stage, flag the rules that read as generic, and apply your amendments inline.

## After you accept

You can now tell the first officer to dispatch your entities. From there, [Operate a workflow](../operating/) covers the day-to-day loop.

## Sitemap

- [Operate a workflow](../operating/index.md)
- [Debrief a session](../debrief/index.md)
