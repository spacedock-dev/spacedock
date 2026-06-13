# Workflows & entities

**Your work lives as plain text in your repo: readable, editable, diffable, and nothing is lost between sessions.** A workflow is a directory plus a README; an entity is one markdown file inside it, one work item. Everything about a work item is in its file: the problem, the design notes, the bar for done, the reports filed as it advances.

## The README is where you set the rules

The README is the single source of truth: it declares the stages and defines what each stage means, what counts as good, and what a worker must produce. Commission generates it; you edit it like any file, at any time. When a bar turns out fuzzy in practice, tighten the prose and the next dispatch works to the new line.

A generated README's frontmatter looks like this:

```yaml
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

Each `states` entry is one stage; its flags decide where work starts and ends, which stages pause for your call, and where rejected work bounces back to. [What a stage declares](stage-lifecycle.md#what-a-stage-declares) defines every property.

## Each work item is one file

An entity lives as a flat file `{slug}.md`, or a folder `{slug}/index.md` when reports and artifacts accumulate beside it. The body is the human-readable record: the problem, the approach, the acceptance criteria, and the stage reports. On top sits YAML frontmatter, the machine-readable state: the item's id, its current stage, its outcome. The [frontmatter contract](../reference/frontmatter-contract.md) lists every field.

## Keep workflow state off your code branch

A workflow can keep its mutable state in a separate state checkout, so routine stage transitions never churn your code branch or collide with a feature PR. The README opts in with one field:

```yaml
state: .spacedock-state
```

The README stays on your code branch as the living spec; the work items and their reports live on a separate branch your code branch never sees. [Multi-workflow & split-root state](../advanced/split-root-state.md) covers the mechanics and the fresh-clone setup.
