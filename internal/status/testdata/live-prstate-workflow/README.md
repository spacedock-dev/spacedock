---
entity-type: task
entity-label: task
entity-label-plural: tasks
id-style: sequential
stages:
  defaults:
    worktree: false
    concurrency: 1
  states:
    - name: backlog
      initial: true
    - name: implementation
      worktree: true
    - name: done
      terminal: true
---

# Live-PR-State Fixture Workflow

Pins the shallow-boot accuracy dependency: a PR-bearing non-terminal entity's
`pr_state` entry in `status --boot --json` must reflect the LIVE merge state (from
`gh pr view`), not just the stored `pr:` field. The shallow-boot greet's accurate
local state summary rests on this. A stubbed `gh` on PATH supplies the live state
in the test.
