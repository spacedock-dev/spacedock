# Frontmatter contract

Spacedock reads YAML frontmatter as the machine-readable state of a workflow and its entities. The first officer and dispatched workers write it; you operate through them, not by hand-editing. The field-level contract (names, types, patterns, defaults, invariants) is two machine-checkable schemas.

## Workflow README

The workflow `README.md` frontmatter declares the entity type, the id style, and the stages with their per-stage defaults and gates. The contract is [`workflow-readme.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/workflow-readme.mdschema.yml), which also specifies the required per-stage body subsections.

## Entity

Each entity's frontmatter carries its id, current stage, outcome, and worktree state. The contract is [`entity.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/entity.mdschema.yml), which defines the fields, the custom-field policy, the recognized body headings, and the invariants.

External gate apply paths may also write `gate-id` and `gate-verdict`. `gate-id` is caller-supplied by the surface that captured the human gate decision, such as a gate UI packet, chat approval recovery event, handoff packet, or other harness id; `spacedock status --apply-gate` records it but does not mint it. These fields record the last gate decision handed to `spacedock status --apply-gate`; they do not replace the terminal `verdict` field used when an entity closes.
