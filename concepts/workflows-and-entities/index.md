---
title: "Workflows & entities"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-18 00:31:29"
---

# Workflows & entities

**Your work lives as plain text in your repo: readable, editable, diffable, and nothing is lost between sessions.** A workflow is a directory plus a README; an entity is one markdown file inside it, one work item. Everything about a work item is in its file: the problem, the design notes, the bar for done, the reports filed as it advances.

## The README is where you set the rules

The README is the single source of truth: it declares the stages and defines what each stage means, what counts as good, and what a worker must produce. Commission generates it; you edit it like any file, at any time, or more likely ask the first officer to edit it to improve the workflow. When a bar turns out fuzzy in practice, tighten the prose and the next dispatch works to the new line. Spacedock's own [dev workflow README](https://github.com/spacedock-dev/spacedock/blob/main/docs/dev/README.md) is a live example.

A generated README's frontmatter looks like this:

```
---
commissioned-by: spacedock@0.20.0
entity-type: task
entity-label: task
entity-label-plural: tasks
id-style: sd-b32
stages:
  defaults:
    worktree: false
    concurrency: 3
  states:
    - name: backlog
      initial: true
      gate: true
    - name: implementation
      worktree: true
    - name: validation
      fresh: true
      feedback-to: implementation
      gate: true
    - name: done
      terminal: true
---
```

A stage can also select ordered context from other README sections:

```
stages:
  defaults:
    context-sections:
      - Constraint Authority
```

`context-sections` contains exact README heading text. A stage-specific list replaces the default, and `[]` clears it. `dispatch build` validates every selection before spawn and writes an exact `dispatch show-stage-def` command into the dispatch file. The command uses the resolved absolute launcher selected for that build, so a later `SPACEDOCK_BIN` or PATH change cannot redirect stage loading. Its output contains the stage definition followed by the selected sections in declaration order. Selected sections normalize newline boundaries to LF, drop trailing blank lines, use one blank line between sections, and end stdout with one LF. The outer fresh or reuse-advance prompt remains a file pointer plus fixed routing metadata.

## Each work item is one file

An entity lives as a flat file `{slug}.md`, or a folder `{slug}/index.md` when reports and artifacts accumulate beside it. The body is the human-readable record: the problem, the approach, the acceptance criteria, and the stage reports. On top sits YAML frontmatter, the machine-readable state: the item's id, its current stage, its outcome. The [frontmatter contract](../../reference/frontmatter-contract/) has the fields and the schemas that define them.

## Keep workflow state off your code branch

A workflow can keep its mutable state in a separate state checkout, so routine stage transitions never churn your code branch or collide with a feature PR. The README opts in with one field, and you can ask the first officer to change the setting:

```
state: .spacedock-state
```

[Split-root state](../../advanced/split-root-state/) covers the mechanics and the fresh-clone setup.

## Sitemap

- [The operating model](../operating-model/index.md)
- [The stage lifecycle](../stage-lifecycle/index.md)
- [Gates & decisions](../gates-and-decisions/index.md)
