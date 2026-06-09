# Frontmatter contract

Every entity is a markdown file (or a folder with an `index.md`) whose YAML frontmatter carries the fields Spacedock reads to track and move it. The always-current schema lives in the development workflow's [Schema / Field Reference](../contributing/development-workflow.md#field-reference); this page surfaces that table and the external-tracker bridge fields in one place for reference. A standalone `docs/specs/frontmatter-contract.md` is a planned follow-up; until it lands, the development workflow README is the source of truth.

The frontmatter parser is line-oriented. Keep fields flat and top-level. If richer metadata becomes necessary, add more flat custom fields rather than nested YAML, because the v1 parser preserves lines, not arbitrary structure.

## Entity fields

These are the fields the development workflow declares for a `task` entity. Other workflows may rename the entity type and adjust fields, but these are the contract a dev-workflow entity must satisfy.

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique 24-character Spacedock Base32 ID, because this workflow uses `id-style: sd-b32`. |
| `title` | string | Human-readable entity name. |
| `status` | enum | One of: `backlog`, `ideation`, `implementation`, `validation`, `done`. The current stage. |
| `source` | string | Where the entity came from. Also used by the external-tracker bridge (see below). |
| `started` | ISO 8601 | When active work began. |
| `completed` | ISO 8601 | When the entity reached terminal status. |
| `verdict` | enum | `PASSED` or `REJECTED`. Set at the final stage. |
| `score` | number | Priority score, `0.0`–`1.0` (optional). A workflow can upgrade to a multi-dimension rubric in its README. |
| `worktree` | string | Worktree path while a dispatched agent is active; empty otherwise. |
| `issue` | string | Optional external ticket reference, such as `ENG-123`, `kata:task-abc123`, or `owner/repo#42`. |

The `status` field is the execution state. `spacedock status` reads stage declarations from the workflow README and reports each entity's `status` against them; `--set status=<stage>` is the mutation that advances an entity. The status read path does not invent stages — if the README declares no stages block, membership cannot be validated.

The `verdict` field is guarded on the finalize action, not on reaching a terminal stage. The guard keys on the finalize shape — a `--set` that writes `completed`, or an `--archive` of a terminal entity — and refuses it (exit 1, entity unmutated) when the post-state `verdict` is empty (see `internal/status/verdict_guard_test.go`). A bare dispatch into a terminal stage that does not write `completed` passes without a verdict, because the verdict is the outcome of work that has not happened yet. A finalize on an entity that already carries a verdict also passes, even when that `--set` does not re-name it. `--force` bypasses the guard. The failure mode the guard catches is finalizing or archiving a terminal entity with no verdict on record.

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

Fill `title`, `status`, and `source` at creation. `started`, `completed`, `verdict`, and `worktree` are written by the runtime as the entity moves; do not edit them by hand while a dispatched agent is active.

## External-tracker fields

The `issue` and `source` fields are the v0 bridge to an external ledger such as kata, Linear, or GitHub Issues. They are flat top-level fields the current parser preserves; the bridge adds no tracker-specific stage rules. See [Tracking work in an external system](../advanced/external-tracker.md) for the full integration model.

```yaml
issue: ENG-123
source: linear
```

or:

```yaml
issue: kata:task-abc123
source: kata
```

The contract for these two fields:

- **`issue` is the human-facing external reference.** It points at the ticket the entity mirrors; Spacedock does not parse its internals.
- **`source` records where the entity came from** when useful — the tracker name, or any origin marker.
- **Spacedock `status` remains the execution status.** The external tracker does not redefine Spacedock stage semantics inside the entity, and ownership stays one-way unless a future bridge explicitly declares bidirectional sync.

## Validating an entity

Check an entity against the contract with the status command's `--validate` flag:

```bash
spacedock status --workflow-dir docs/dev --validate
```

It exits 0 when the workflow is valid and 1 when it is not, printing the errors to stderr; with `--json` it also emits a `{"command":"validate","valid":"true"}` (or `"false"`) envelope. `--validate` cannot be combined with the other status flags — `--next`, `--next-id`, `--boot`, `--where`, `--fields`/`--all-fields`, `--archived`, `--archive`, or `--set`; the command rejects the combination. Validation reads stages from the workflow README and entities from the state checkout, so it enforces the contract against the same schema the workflow declares — not an assumed one.
