---
title: "Welcome"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-24 21:57:38"
---

# Spacedock

Spacedock runs your work as a series of stages. **Nothing crosses a gate without a decision you own.**

A gate is a checkpoint where the workflow pauses and puts the question to you: ship this, or not? You approve it, send it back, or escalate. You can also delegate the call to an agent. Either way, the decision is recorded with its evidence and its reason. That is the whole idea. Everything else is detail.

You are the captain. You set the bar and make the calls; the agents do the rest. The bar starts rough and sharpens every time you reject, so calls that once needed you become ones you can hand off with confidence. See [the operating model](concepts/operating-model/) for how the three roles divide the work.

## What's different

- **The agent doesn't get to judge its own work.** Review runs as a separate stage with fresh context, no access to the maker's reasoning. It pushes back on thin evidence and work that looks busy without proving its claim.
- **Every decision leaves a trail.** Each gate carries a stage report: findings, verdicts, artifacts, anomalies. You decide on evidence, not the transcript, and the record outlives the reviewer.
- **The bar sharpens as you use it.** Each stage declares what good means and the agent works to that line. When a standard turns out fuzzy in practice, the agent proposes an edit to the written criteria for your approval.
- **Batch the work; decide as it flows back.** Queue many work items at once. Agents advance each through its stages, and you handle gates as they surface, not one session at a time.
- **Work survives the context limit.** When an agent runs out of context, a successor carries forward what's in flight.

## Where to go next

- **[Get started](get-started/install/)**: [install](get-started/install/) Spacedock, then pick an entry. [Survey an existing project](get-started/survey/) to see where your agents burn your time and surface the workflow you are already running without naming it. Or [start a fresh workflow](get-started/first-workflow/) from a common shape like development or research.
- **[Concepts](concepts/operating-model/)** covers the operating model, workflows and entities, the stage lifecycle, and gates and decisions.
- **[Running workflows](running-workflows/commission/)** walks through commissioning a workflow, operating a running workflow, and debriefing and refitting between sessions.

## For agents using Spacedock

Agents read these docs too. Start from [`llms.txt`](/llms.txt), the curated index of these pages.

## Sitemap

- [Install](get-started/install/index.md)
- [Survey your project](get-started/survey/index.md)
- [Your first workflow](get-started/first-workflow/index.md)
