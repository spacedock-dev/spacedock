# Frontmatter contract

An entity's YAML frontmatter is the machine-readable state Spacedock reads to track and move it. The field-level contract (names, types, patterns, defaults, invariants) is the two machine-checkable schemas:

- [`workflow-readme.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/workflow-readme.mdschema.yml) — the workflow `README.md` frontmatter and the required per-stage body subsections.
- [`entity.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/entity.mdschema.yml) — entity frontmatter fields, the custom-field policy, recognized body headings, and the invariants.

## A development-workflow entity

The fields you see most often are the ones the development workflow declares for a `task`. A new entity ships with these blank, and the runtime fills the rest as the entity moves:

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

You fill `title`, `status`, and `source` at creation. `status` names the current stage (`backlog`, `ideation`, `implementation`, `validation`, `done` for this workflow), and the first officer advances it as work moves. The runtime owns `started`, `completed`, `verdict` (`PASSED` or `REJECTED` at the [final gate](../concepts/gates-and-decisions.md#the-three-calls)), and `worktree`; don't hand-edit those while a dispatched agent is active.
