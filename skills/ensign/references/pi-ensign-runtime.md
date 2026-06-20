# Pi Ensign Runtime

This file defines how the shared ensign core executes on Pi.

## Agent Surface

A Pi ensign receives a bounded assignment from a Pi-native substrate such as `pi-subagents` or, through an adapter, `pi-agent-teams`. The assignment content is authoritative: entity path, workflow directory, target stage, stage definition fetch command, worktree path when present, and completion checklist.

Do not assume Claude team tools exist in Pi. Completion is reported by the worker's final result in the Pi turn or by the active Pi adapter's task-completion notification.

## Pi-Specific Rules

- If the dispatch prompt is a read-dispatch-file pointer, read that file first and treat its content as the assignment.
- If no worktree path is provided, stay on the repo root/main-branch checkout named by the assignment.
- If a worktree path is provided, keep code reads, writes, tests, and code commits under that worktree.
- In split-root workflows, the entity body and stage report stay at the state-checkout entity path provided by the dispatch; do not invent a worktree copy of the entity.
- Do not modify YAML frontmatter in the entity file.
- Commit the entity body/stage report path-scoped in the state checkout when the assignment asks for a state commit.
- Treat each follow-up assignment as a fresh cycle; never assume a previous completion still satisfies the current assignment.

## Clarification

If requirements are unclear or ambiguous, ask for clarification via `contact_supervisor` with `reason: "need_decision"` rather than guessing. Describe what you understand and what's ambiguous so the FO can route a quick answer. The FO replies via `intercom({action:"reply", message:"..."})`; after receiving the reply, resume working. `need_decision` blocks with a 10-minute timeout — ask only when genuinely blocked, and surface enough context for a one-shot answer.

For non-blocking plan-changing discoveries (a riskier mechanism panning out, a scope boundary discovered, an unexpected dependency), send `contact_supervisor` with `reason: "progress_update"` — the FO reads and acknowledges; no reply is required and you continue working without waiting.

## Completion

When done, return one concise final result that names:

- the entity and stage completed;
- the entity file containing the stage report;
- the commit or durable evidence produced;
- any residual risk or blocker.

Either the final result (the subagent return) or an explicit done-message sent via `contact_supervisor` / `intercom send` is a valid completion signal — the FO verifies the entity-file stage report either way. After sending that completion result, stop. Do not idle waiting for another message unless the active Pi substrate explicitly delivers one or the FO routes follow-up through the back-channel.
