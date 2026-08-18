---
title: "Your first workflow"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-18 04:17:46"
---

# Your first workflow

Your first workflow comes from one of two places: [survey](../survey/) offers to build one from what it found in your project, or you describe one yourself to the `/spacedock:commission` skill. A workflow is your work moving through stages you define, pausing at the gates where you decide.

## Commission a workflow

Describe the work in the launch line:

```
spacedock claude "/spacedock:commission ship features through design, implementation, and review"
```

Commission asks a few questions, one at a time: what each work item is, the stages it moves through, which stages pause for your decision, and the quality bar for each. You confirm or adjust each proposal, and it presents the design:

> I'll call you Captain; let me know if you prefer something else.
>
> For each run, we process tasks going through the following stages:
>
> a. backlog\
> new tasks wait here for triage
>
> b. design\
> the approach is worked out and written down
>
> c. implementation\
> the change is built on an isolated branch
>
> d. review\
> the result is checked against the design
>
> e. done\
> accepted and closed
>
> If you reject at review, it goes back to implementation for revision.
>
> Our pilot run will be with:
>
> - Add rate limiting to the API
> - Fix the flaky login test
>
> Entity identity will use `sd-b32`.
>
> All files will be created in `docs/ship-features/` for you to review.
>
> Accept this design, or tell me what to change.

## What gets generated

Everything is plain text in your repo: a README that holds the workflow's rules (what each stage expects) and one file per work item. The generated rules are a starting point; tighten them before any work runs, because an agent working to a vague bar is expensive to correct. Learn more about every design decision in [Commission a workflow](../../running-workflows/commission/).

## The pilot run

On accept, commission dispatches your seed items in parallel, each moving through the stages until everything is idle or waiting on you. When a work item reaches the `review` gate, you get a gate review:

```
Capability/change: token-bucket limiter at the API middleware layer.
Test and evidence: the sixth failed login returns 429, and a successful login resets the counter.
Reviewed snapshot: Briefing `...` at compact digest `sha256:1a2b3c4d…`.
Recommendation: approve the implementation.
Decision ask: approve to close, revise with feedback, or hold at review.
```

You approve, send it back with feedback, or reject. Details: [gates and decisions](../../concepts/gates-and-decisions/).

## Operate the workflow

Every work item's state is stored and serialized in the [workflow](../../concepts/workflows-and-entities/), so you do not need to worry about context limits, resuming, or clearing. A typical session:

```
spacedock claude
```

It picks up whatever is ready to dispatch. Keep dispatching, approving, rejecting, steering. Before you stop, run [`/spacedock:debrief`](../../running-workflows/debrief/) to record what happened, update the learnings into the workflow, and file the follow-up items. The workflow self-improves. [Operate a workflow](../../running-workflows/operating/) covers the details.

A project can hold [multiple workflows](../../advanced/multi-workflow/); Spacedock finds them and drives them together. Since this runs in your existing coding agent, you can just ask the agent if anything is unclear.

Now you have the first Spacedock-powered workflow: dispatch and let the agents work and hum when you stay calm!

## Sitemap

- [Install](../install/index.md)
- [Survey your project](../survey/index.md)
