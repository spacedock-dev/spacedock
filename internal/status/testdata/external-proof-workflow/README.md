---
entity-type: task
entity-label: task
entity-label-plural: tasks
id-style: sequential
require-external-proof: true
stages:
  defaults:
    worktree: false
    concurrency: 1
  states:
    - name: backlog
      initial: true
      gate: true
    - name: ideation
      gate: true
    - name: implementation
    - name: done
      terminal: true
---

# External-Proof Fixture Workflow

A fixture workflow opting into `require-external-proof: true` so the guard's
terminal-set + --validate sub-check are exercised end to end.
