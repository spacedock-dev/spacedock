# Pi Ensign Runtime

How the shared ensign core executes on Pi. The shared core owns the ensign discipline (assignment reading, worktree, split-root commit, frontmatter, proof, stage report); this adapter binds only the Pi-specific concerns.

## Runtime implementation

- `Clarification` -> `contact_supervisor` with `reason: "need_decision"` when genuinely blocked, naming what you understand and what is ambiguous for a one-shot answer; the FO replies via `intercom({action:"reply", message:"..."})` and you resume. `need_decision` blocks with a 10-minute timeout. For non-blocking plan-changing discoveries, use `reason: "progress_update"` — the FO acknowledges, no reply required, and you continue.
- `Completion signal` -> return one concise final result naming the entity and stage completed, the entity file containing the stage report, the commit or durable evidence produced, and any residual risk. Completion is reported by the worker's final result in the Pi turn or by the active Pi adapter's task-completion notification. The subagent return is the primary signal and the FO file-verifies the stage report; an optional `intercom send` done-advisory may precede the return as a heads-up. After returning, stop unless the active Pi substrate delivers another message or the FO routes follow-up through inter-agent communication.
