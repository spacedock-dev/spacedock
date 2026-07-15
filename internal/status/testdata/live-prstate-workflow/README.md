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

Pins ordinary, non-identify boot compatibility: a PR-bearing non-terminal
entity's `pr_state` entry in `status --boot --json` must reflect the LIVE merge
state (from `gh pr view`), not just the stored `pr:` field. Identify boot is
covered separately and deliberately uses the local PR mirror. A stubbed `gh` on
PATH supplies the live state in this test.
