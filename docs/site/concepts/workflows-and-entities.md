# Workflows & entities

A workflow is a directory of plain-text work items plus a README that defines the stages, the schema, and the gates. Everything about a work item lives in the file itself.

## A workflow is a directory plus a README

The README declares the workflow: its stages, its frontmatter schema, and which stages are gated. The work items live alongside it. The README is the single definition the first officer reads to know how to drive the workflow.

## An entity is one work item

An **entity** (also called a *work item*) is one markdown file — or a folder, when reports and artifacts accumulate beside it. Each entity is named `{slug}` or `{slug}.md` (lowercase, hyphens, no spaces). Use folder-form entities when reports may pile up; for example `native-go-status/index.md`.

Everything about the work item lives in the file: the problem, the design notes, the bar for done, and the stage reports. State survives a session; the next one picks up where you left off.

## Entity frontmatter

Every entity carries YAML frontmatter. The field set is declared by the workflow README. The development workflow's schema is the canonical example:

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique entity identifier |
| `title` | string | Human-readable name |
| `status` | enum | The current stage |
| `source` | string | Where the work item came from |
| `started` | ISO 8601 | When active work began |
| `completed` | ISO 8601 | When the entity reached its terminal stage |
| `verdict` | enum | PASSED or REJECTED — set at the final stage |
| `worktree` | string | Worktree path while a dispatched agent is active, empty otherwise |
| `issue` | string | Optional external ticket reference (e.g. `ENG-123`, `owner/repo#42`) |

The full field reference for the development workflow lives on the [Frontmatter contract](../reference/frontmatter-contract.md) page.

## The state checkout

Some workflows separate the workflow *definition* from the workflow *runtime state*. With the `state-branch` profile, the README stays in the main repo while the mutable entities live in a per-workflow `.spacedock-state` checkout. This keeps state transitions out of the code branch. See [Multi-workflow & split-root state](../advanced/split-root-state.md) for how the two roots compose.
