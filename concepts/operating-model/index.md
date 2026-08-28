---
title: "The operating model"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-28 17:10:44"
---

# The operating model

Spacedock runs on three roles and one division of labor: you shape the work and make the calls; the agents drive each item through its stages and bring decisions back to you with evidence.

## Roles

| Role              | Who                    | Ownership                                                               |
| ----------------- | ---------------------- | ----------------------------------------------------------------------- |
| **Captain**       | You                    | The mission, and the call at every approval gate unless delegated       |
| **First Officer** | The orchestrator agent | Runs the workflow for you and brings each decision to you with evidence |
| **Ensign**        | The worker agent       | Moves one work item through one stage                                   |

Each session has one captain and one first officer; ensigns come and go with the work.

The first officer keeps the work moving so you do not have to: it finds what is ready, hands each item to a worker, and checks the result against the bar you set. Gated stages pause for your call; everything else flows forward without you.

## Shaping versus driving

The captain shapes the product and owns the workflow; the agents drive. These are different jobs, and the split is what keeps you out of the per-step loop.

**Shaping is the product judgment: the goal, the taste, the steering.** What to build, what good looks like, which direction survives a gate. You make the calls at gates (approve, redo with feedback, or reject); some you answer yourself, others resolve through a delegated agent review. This judgment is the part the agents cannot supply.

**Owning the workflow is the structural half.** You set the stages and the bar each stage must clear, declared and serialized in the workflow README, starting from [`/spacedock:commission`](../../running-workflows/commission/). The declaration is shapable mid-task: you do not have to get it right the first time. When a bar turns out fuzzy in practice, tighten the README and the next dispatch works to the new line.

**Driving is moving work items through the declared stages.** The first officer schedules and dispatches; the ensign does the stage work and proves it. The first officer acts on its own for routine, reversible steps and asks you only when something is genuinely ambiguous: unclear requirements, a design choice that would change the output, scope too vague to turn into criteria. Everything else happens without a prompt to you.

The line holds because of one rule: the maker does not judge its own work. Review runs as a separate stage with fresh context and no access to the maker's reasoning (see [Gates and decisions](../gates-and-decisions/)). The first officer never self-approves a gated stage.

## Batched, evidenced decisions

Decisions reach you batched and backed by evidence, not as a stream of interruptions. Your attention is the bottleneck; the agents queue work and surface only the calls that need a human.

**Batch the work; decide as it flows back.** Queue many work items at once. Ensigns advance each through its stages in parallel. You handle gates as they surface, not one session at a time, and not on the agent's schedule. While one item waits on a clarification, the first officer keeps dispatching the others.

**Every gate carries evidence.** The first officer does not hand you the transcript. It presents a short review: what was chosen, the evidence for it, and one recommendation you can approve with a single yes. [Gates and decisions](../gates-and-decisions/) shows the format.

**The decision leaves a trail.** Each gate records the verdict and its reason alongside the stage report in the work item's file. The record outlives the reviewer, so a bad result traces back to the call that caused it. When you end a session, [`/spacedock:debrief`](../../running-workflows/debrief/) captures what happened, and the next session picks up from it.

## Where to go next

- [Workflows and entities](../workflows-and-entities/) to see where your work lives as plain files.
- [Stage lifecycle](../stage-lifecycle/) to follow one item end to end.
- [Gates and decisions](../gates-and-decisions/) to see exactly what you decide and on what evidence.

## Sitemap

- [Workflows & entities](../workflows-and-entities/index.md)
- [The stage lifecycle](../stage-lifecycle/index.md)
- [Gates & decisions](../gates-and-decisions/index.md)
