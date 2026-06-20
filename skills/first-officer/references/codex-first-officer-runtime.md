# Codex First Officer Runtime

This file defines how the shared first-officer core executes on Codex. The shared core owns invocation timing; this adapter binds Codex's live tool surface to the core `«fn»` capabilities and keeps only Codex-specific probe, wait, and feedback notes.

## Live Tool Surface Probe

At boot, bind Codex capabilities from the live tool surface, not from a runtime-version label. Probe the expected call shapes and report concrete tool-surface blockers instead of assuming historical labels:

- `«worker.spawn»`: `spawn_agent(task_name,message,fork_turns)`
- `«addressable-worker»`: `send_message(target,message)` plus `followup_task(target,message)`, or the compatibility `send_input` surface
- `«completion-signal»`: `wait_agent(timeout_ms)` for foreground waiting, while the async mailbox notification remains the signal
- `«roster-reconcile»`: `list_agents(path_prefix?)`
- `«worker.shutdown»`: unresolved unless a shutdown-specific probe proves a binding; Do not bless `interrupt_agent`

Do not infer capabilities from a Codex version name. Do not infer that the turn-starting addressable-worker route is absent from the absence of newer follow-up tools.

## Runtime implementation

- `«worker.spawn»`: When bound, call `spawn_agent(task_name,message,fork_turns)` with the helper-emitted prompt unchanged as `message`. sanitize the helper-emitted `name` to a lowercase digit/underscore `task_name`; retain the mapping from helper identity to sanitized task path, entity slug, stage, and cycle.
- `«addressable-worker»`: A spawned worker remains addressable by live task path/handle while active or completed and still listed, including a completed-but-still-addressable worker. `send_message(target,message)` sends non-triggering context or cooperative preservation text; `followup_task` starts a worker turn on the addressed worker; `send_input` is both the compatibility address and turn-starting route.
- `«async-dispatch»`: Codex dispatch is async; the spawn call returns a handle and the FO event loop remains available for mailbox notifications, gate work, and captain interaction.
- `«worker-identity»`: Record Codex task name, mailbox handle, sanitized path, helper-emitted identity, entity slug, stage, cycle, completion epoch, and thread-inherited model when the helper emits null.
- `«completion-signal»`: The single completion signal is the async final-status notification in the FO mailbox. Foreground wait is only an idle observation action and does not itself complete the worker.
- `«worker.shutdown»`: `«worker.shutdown»` remains unresolved until probed. Durable workflow state remains authoritative. `«addressable-worker»` may carry cooperative preservation text when present.
- `«context-budget»`: ABSENT; reuse-condition-0 is satisfied because Codex has no bound context-budget probe.
- `«roster-reconcile»`: `«roster-reconcile»` may provide active/completed task-path reads when bound. `list_agents(path_prefix?)` provides attribution, stale-cohort classification, debugging, and cleanup targeting.

## Codex wait notes

Use foreground wait only when no other dispatchable or gate-processing work is available and an unresolved worker completion is the next useful idle action. A wait timeout return is normal and retryable; it means no final-status mailbox update arrived before the deadline. Before foreground waiting, tell the captain that pressing Esc or otherwise causing an operator interruption only returns control; the worker is not failed, closed, or redispatched, and the next foreground wait is reinstalled when waiting is again the next useful idle action.

A wait return must be attributed by mailbox content, task path, `«roster-reconcile»` state when present, or durable workflow state, not by a handle argument. A captain message or shell-out during the wait is operator activity, not idle wake evidence.

### Foreground wait

The foreground-wait operator cue must continue to state that an interruption returns control, does not fail, close, or redispatch the worker, and allows the next foreground wait to retry the same unresolved worker when waiting is again the next useful idle action.

### Queued notification flushed by later activity

The FO ends the turn without foreground wait, and a later captain message, tool action, or shell-out causes Codex to deliver a worker final-status notification that had already been queued. This is queued/activity-driven delivery, not autonomous FO wake-up.

### Autonomous idle FO wake-up

The FO does not use foreground wait, performs no later captain message, tool action, shell-out, or terminal job, and Codex starts a new assistant turn from the worker final-status notification alone. This is the only observation that proves no-wait autonomous FO wake-up.

## Feedback reviewer reuse

Feedback rejection is the load-bearing exception to casual fresh dispatch: when `«addressable-worker»` is bound, keep the validation reviewer addressable after its REJECTED report, route the fix to the implementation worker, then re-run the kept-alive validation reviewer through the same turn-starting `«addressable-worker»` binding. A completed validation report does not by itself make that reviewer unavailable, even when the reviewer already sent a completion signal; first attempt reuse with the captured validation-worker handle. Fresh-dispatch the reviewer only when the existing reviewer is no longer addressable or reuse conditions fail.

## Captain Interaction

The captain is the user of the Codex session. Communicate gate results, clarifications, and status directly in the conversation.
