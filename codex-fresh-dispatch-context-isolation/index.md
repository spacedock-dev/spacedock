---
id: rt8eywbyf3d7nyc9bsmrymnq
title: Codex fresh dispatch must isolate parent turns
status: backlog
source: "/tmp/spacedock-codex-fresh-dispatch-incident.md; captain-requested legitimacy check and filing, 2026-07-19"
started:
completed:
verdict:
score: "0.92"
worktree:
issue:
sprint:
---

A Codex fresh dispatch currently guarantees a new worker handle but not a clean
context. CodexMultiAgentV2Spawn.ToolArgs omits fork_turns when the field is
empty, its unit test requires that omission, and Codex interprets omission as
fork_turns: "all". The live incident also recorded an independent validator
spawned explicitly with "all". That validator found a real defect, but its
independence evidence must be rerun.

Spacedock already has a simpler boundary: every spawn_agent call is a fresh
dispatch, while deliberate continuity uses followup_task. Therefore every
Spacedock Codex spawn should pass fork_turns: "none" unconditionally. The
adapter should not expose "all" or a numeric turn count.

The existing per-host-stage-model-override task is not a duplicate. It forces
"none" only when model or effort overrides require it and explicitly records
the all-spawn default as a follow-up candidate.

## Acceptance criteria

- **AC-1 (VALUE):** Every Spacedock Codex fresh dispatch reaches spawn_agent
  with fork_turns: "none", so a fresh validation worker receives only its
  generated dispatch artifact and no parent turns. Verified by adapter tests for
  plain and override-bearing envelopes plus a Codex live scenario that records
  the actual spawn arguments.
- **AC-2:** The adapter cannot emit an omitted, "all", or numeric fork_turns
  value. Verified by tests that fail each unsafe shape and confirm ToolArgs
  always contains exactly fork_turns: "none".
- **AC-3:** The Codex First Officer contract defines fresh dispatch as a new
  handle with no inherited parent turns and binds the spawn call to
  fork_turns="none". Verified with the behavioral live evidence from AC-1;
  wording-presence alone does not satisfy this criterion.
- **AC-4:** Deliberate reuse remains unchanged: eligible stage advancement and
  feedback cycles use followup_task with the existing handle, and fresh: true
  selects a new isolated spawn. Verified by adapter or scenario coverage for one
  reuse path and one fresh: true path.
- **AC-5:** The affected release validation is rerun with a new worker,
  fork_turns: "none", and only the generated dispatch artifact. Its durable
  report records the replacement evidence without invalidating findings from
  the original run merely because their independence was weakened.

## Scope

Do not force every implementation stage to be fresh. A workflow that requires
role separation should declare fresh: true; ideation-to-implementation reuse
remains valid when it does not. This task fixes the meaning of every actual
fresh spawn.
