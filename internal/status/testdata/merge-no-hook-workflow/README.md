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
      gate: true
    - name: ideation
      gate: true
    - name: implementation
      worktree: true
    - name: done
      terminal: true
---

# No-Merge-Hook Fixture Workflow

The default-policy (`merge:` key absent) sibling with NO `_mods/` directory at
all — no merge hook is registered. `merge guard`'s Phase C default case reaches
these entities directly: the terminalize `--set` is unguarded (no merge hook to
require) and finalize proceeds without ever having invoked a hook or recorded a
merge sentinel, which is exactly the state the D3 finalized-line "no merge hook
registered" clause exists to name.
