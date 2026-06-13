# Frontmatter contract

Every entity is a markdown file (or a folder with an `index.md`) whose YAML frontmatter carries the fields Spacedock reads to track and move it. This page is a quick lookup for the fields a development-workflow entity uses. The always-current schema is owned by the [development workflow README](https://github.com/spacedock-dev/spacedock/blob/main/docs/dev/README.md#field-reference); when this table and that README disagree, the README wins.

Keep fields flat and top-level; add more flat custom fields rather than nested YAML. Spacedock reads frontmatter line by line, so nested values are ignored.

## Entity fields

These are the fields the development workflow declares for a `task` entity. Other workflows may rename the entity type and adjust fields.

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique 24-character Spacedock Base32 ID (this workflow uses `id-style: sd-b32`). |
| `title` | string | Human-readable entity name. |
| `status` | enum | The current stage: `backlog`, `ideation`, `implementation`, `validation`, or `done`. |
| `source` | string | Where the entity came from. Also used by the external-tracker bridge. |
| `started` | ISO 8601 | When active work began. |
| `completed` | ISO 8601 | When the entity reached terminal status. |
| `verdict` | enum | `PASSED` or `REJECTED`, set at the final stage. A terminal close requires it; see [gates and decisions](../concepts/gates-and-decisions.md#the-three-calls). |
| `score` | number | Priority score, `0.0`–`1.0` (optional). |
| `worktree` | string | Worktree path while a dispatched agent is active; empty otherwise. |
| `issue` | string | Optional external ticket reference, such as `ENG-123` or `owner/repo#42`. |

`status` is the execution state: `spacedock status` reports each entity's `status` against the README's stage declarations, and `--set status=<stage>` advances it. The fields the runtime writes (`started`, `completed`, `verdict`, `worktree`) should not be hand-edited while a dispatched agent is active.

## Copy-paste starter

The development workflow's task template ships these fields blank for a new entity:

```yaml
---
id:
title: Task name here
status: backlog
source:
started:
completed:
verdict:
score:
worktree:
issue:
---
```

Fill `title`, `status`, and `source` at creation; the runtime writes the rest as the entity moves.

## External-tracker fields

`issue` and `source` are the bridge to an external ledger (Linear, GitHub Issues). `issue` is the human-facing ticket reference; `source` records the origin. Spacedock `status` stays the execution status, and the tracker does not redefine stage semantics inside the entity. See [Bridge an external tracker](../advanced/external-tracker.md) for the full bridge model.

## Validating an entity

Check the workflow and its entities against the contract:

```bash
spacedock status --workflow-dir docs/dev --validate
```

It exits 0 when valid, 1 with errors on stderr otherwise. Validation reads stages from the README and entities from the state checkout, so it enforces the schema the workflow declares.
