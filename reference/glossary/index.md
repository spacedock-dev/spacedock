---
title: "Glossary"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-26 08:45:55"
---

# Glossary

The core Spacedock terms, in one place. Each links to where the concept is covered in depth.

## Workflow

A directory plus a `README.md` that declares the stages a work item moves through and what each stage expects. A project can hold more than one. See [Workflows and entities](../../concepts/workflows-and-entities/).

## Entity

One work item, stored as a single markdown file inside a workflow. Its frontmatter holds the machine-readable state (id, current stage, outcome); its body holds the problem, design notes, the bar for done, and the reports filed as it advances. See [Workflows and entities](../../concepts/workflows-and-entities/).

## Stage

A named step in a workflow (for example `design`, `implementation`, `review`). The README defines what each stage means, what counts as good, and what a worker must produce. See [The stage lifecycle](../../concepts/stage-lifecycle/).

## Gate

The decision point at the end of a stage where nothing advances without your vote. When a stage declares a gate, the first officer stops, presents a review, and waits for you to approve, redo with feedback, or reject. See [Gates and decisions](../../concepts/gates-and-decisions/).

## First officer

The orchestrator agent. It runs the workflow for you: finds what is ready, dispatches workers, checks results against the bar you set, and brings each gate decision to you with evidence. There is one per session. See [The operating model](../../concepts/operating-model/).

## Ensign

The worker agent. It moves one work item through one stage and proves the work, then signals completion. Ensigns come and go with the work. See [The operating model](../../concepts/operating-model/).

## Mod

An opt-in extension that adds behavior to a workflow without changing its stages, such as the [pr-merge mod](../../advanced/mods-and-standing-teammates/) that manages the pull-request lifecycle so merging never has to be a stage. See [Mods and standing teammates](../../advanced/mods-and-standing-teammates/).

## Standing teammate

A long-lived agent the first officer keeps around across dispatches for a recurring role (for example a reviewer), rather than spawning a fresh worker each time. See [Mods and standing teammates](../../advanced/mods-and-standing-teammates/).

## Split-root state

A layout that keeps a workflow's mutable state (entity bodies, stage reports) in a separate state checkout from the code, so concurrent workers commit state without colliding on the code branch. See [Split-root state](../../advanced/split-root-state/).

## Sitemap

- [Command reference](../command-reference/index.md)
- [Frontmatter contract](../frontmatter-contract/index.md)
- [Supported sandboxes](../sandbox/index.md)
