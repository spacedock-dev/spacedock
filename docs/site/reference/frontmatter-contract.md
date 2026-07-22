# Frontmatter contract

Spacedock reads YAML frontmatter as the machine-readable state of a workflow and its entities. The first officer and dispatched workers write it; you operate through them, not by hand-editing. The field-level contract (names, types, patterns, defaults, invariants) is two machine-checkable schemas.

## Workflow README

The workflow `README.md` frontmatter declares the entity type, the id style, and the stages with their per-stage defaults and gates. The contract is [`workflow-readme.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/workflow-readme.mdschema.yml), which also specifies the required per-stage body subsections.

## Entity

Each entity's frontmatter carries its id, current stage, outcome, and worktree state. The contract is [`entity.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/entity.mdschema.yml), which defines the fields, the custom-field policy, the recognized body headings, and the invariants.

A gate decision is recorded in the entity's versioned `gates` collection before any stage transition or worker dispatch. The `spacedock gate record` verb is the sole writer for the gate/attempt/briefing/resolution portion of that collection; callers provide a complete Briefing, a fully associated exact Result, or a chat decision, while the binary derives lifecycle and ids. It preserves application-owned and legacy extension subtrees without interpreting or bulk-rewriting them. Each attempt stores one briefing reference and digest: replaceable while Resolution is absent and frozen when present. `spacedock status --fields gate-state,gate-decision,gate-resolution` surfaces the current recorded result without changing the default status table.
