# Frontmatter contract

Every entity carries YAML frontmatter. The field set a workflow uses is declared in its README. This page surfaces the development workflow's schema as the canonical reference; a standalone, workflow-neutral contract spec under `docs/specs/` is a planned follow-up.

## Field reference (development workflow)

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique 24-character Spacedock Base32 ID (this workflow uses `sd-b32`) |
| `title` | string | Human-readable task name |
| `status` | enum | One of: backlog, ideation, implementation, validation, done |
| `source` | string | Where this task came from |
| `started` | ISO 8601 | When active work began |
| `completed` | ISO 8601 | When the task reached terminal status |
| `verdict` | enum | PASSED or REJECTED — set at the final stage |
| `score` | number | Priority score, 0.0–1.0 (optional). Workflows can upgrade to a multi-dimension rubric in their README. |
| `worktree` | string | Worktree path while a dispatched agent is active, empty otherwise |
| `issue` | string | Optional external ticket reference, such as `ENG-123`, `kata:task-abc123`, or `owner/repo#42` |

The authoritative, always-current copy of this schema is the development workflow's README — see [The development workflow](../contributing/development-workflow.md), under "Schema / Field Reference".

## External tracker fields

When an entity originates from an external ticket system, two flat fields bridge it:

- `issue` — the human-facing external reference.
- `source` — where the entity came from (e.g. `linear`, `kata`).

Spacedock status remains the execution status unless a future bridge explicitly declares bidirectional ownership. For the full bridge contract, see [External-tracker bridge](../advanced/external-tracker.md).
