---
title: "Workflows & entities"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-11 05:14:43"
---

# Workflows and entities

A workflow keeps its rules and work records as plain text in your repository. An **entity** is one work item.

## Workflow rules

The workflow README defines the stages, standards and [gate content](../gates-and-decisions/#gate-content).

Ask your agent to change these rules as you learn what the work requires. You can also edit the README directly.

## Work items

Each work item has a Markdown record with its goal, acceptance criteria, progress and stage reports. Related artifacts can live beside it.

Your agent maintains the records as work proceeds. You can read them to understand what happened or resume in another session.

## State storage

Workflow state can live in a [separate checkout](../../advanced/split-root-state/) to keep routine progress commits off your code branch.

**Related guides**

- [Workflow design](../../running-workflows/commission/)
- [Frontmatter contract](../../reference/frontmatter-contract/)

## Sitemap

- [The operating model](../operating-model/index.md)
- [The stage lifecycle](../stage-lifecycle/index.md)
