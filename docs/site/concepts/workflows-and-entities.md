# Workflows & entities

**A workflow is a directory plus a README, and an entity is one markdown file inside it.** The README defines the stages, the schema, and the gates; each entity is a work item that moves through those stages. Everything about a work item lives in the file itself: the problem, the design notes, the bar for done, the stage reports. State survives a session, so the next one picks up where you left off.

## The workflow: a directory and its README

The README is the single source of truth. Its frontmatter declares the stages, the entity type, and the ID style; its prose body defines what each stage means, what counts as good, and what a worker must produce. You commission a workflow with `/spacedock:commission`, which generates the directory, the README, and a few seed entities for you.

A minimal README frontmatter looks like this:

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

Each `states` entry is one stage. The per-stage flags decide behavior:

- **`initial: true`** marks where a new entity starts; **`terminal: true`** marks where it ends.
- **`gate: true`** makes the stage end at a gate: the first officer pauses and presents a decision instead of advancing on its own.
- **`worktree: true`** runs the stage in its own git worktree, so stages that touch shared state stay isolated. Lighter stages run inline.
- **`fresh: true`** dispatches the stage with no access to the prior worker's reasoning. Use it for review stages so the maker doesn't judge its own work.
- **`feedback-to: <stage>`** names where rejected work bounces back to for revision.

The README body documents each stage with `Inputs`, `Outputs`, `Good`, and `Bad`. That prose is the living spec; it is what every dispatched ensign works to. Tighten it to your actual bar before the first dispatch; editing it after agents have run against vague prose costs more.

## The entity: one work item

**An entity lives as a flat file `{slug}.md` or a folder `{slug}/index.md`.** Use the folder form when reports or artifacts accumulate beside the work item. Slugs are lowercase with hyphens, no spaces (`add-login.md` or `add-login/index.md`). `spacedock status` reads both forms.

The body holds the human-readable record: a description, a problem statement, the proposed approach, acceptance criteria, and the stage reports filed as the entity advances.

## Entity frontmatter

The YAML frontmatter at the top of the file is the machine-readable state. The full schema lives in the workflow's own README; the [frontmatter contract](../reference/frontmatter-contract.md) is the field reference across workflows. The fields you set and read most often:

| Field | What it holds |
|-------|---------------|
| `id` | The unique identifier; format set by `id-style` in the README. |
| `title` | Human-readable name. The filename slug is derived from it. |
| `status` | The current stage, one of the stage names declared in the README. |
| `source` | Where the entity came from (e.g. `commission seed`, `linear`). |
| `started` / `completed` | ISO 8601 timestamps for when work began and when the entity reached terminal. |
| `verdict` | `PASSED` or `REJECTED`, set at the final stage. |
| `score` | Optional priority, 0.0–1.0. |
| `worktree` | The worktree path while a dispatched agent is active; empty otherwise. |
| `issue` | Optional external ticket reference, e.g. `ENG-123`, `kata:task-abc123`, or `owner/repo#42`. |

The frontmatter parser is line-oriented: keep fields flat and top-level. If a workflow needs more metadata, add flat custom fields rather than nested YAML.

`status` is the field that drives everything: the first officer reads it to decide which entities are ready to advance. To see the queue, run the status viewer against the workflow directory:

```bash
spacedock status --workflow-dir docs/dev
```

To list only the entities ready for dispatch, run the query the first officer runs each loop:

```bash
spacedock status --workflow-dir docs/dev --next
```

## Where entities live: the state checkout

**A workflow can keep its mutable entity state separate from the README using a state checkout.** The README frontmatter declares it with one field:

```yaml
state: .spacedock-state
```

With this set, the README (the living spec) stays on your code branch, while the entity files, their reports, and the archive live under `.spacedock-state`. Routine stage transitions then never churn the code branch or collide with a feature PR. The path is resolved relative to the README's directory; `spacedock status` reads stages from the README and entities from the state checkout, and writes frontmatter updates and archive moves into the state checkout.

The state checkout is a linked git worktree on an orphan branch in the same repo, so state commits land on that branch and the code branch never sees them. On a fresh clone the state worktree is absent. Run `spacedock state init` to fetch the orphan branch and re-add the linked worktree before working the workflow.

Omit `state:` (or set `state: $inline`) for a standalone workflow that isn't embedded in a code repo you ship from. Then the entities live beside the README in the same directory, with no orphan branch and no extra checkout.
