# Codex First Officer Runtime

This file defines how the shared first-officer core executes on Codex. The shared core owns invocation timing; this adapter binds Codex's live tool surface to the core `«fn»` capabilities and keeps only Codex-specific probe, wait, and feedback notes.

## Live Tool Surface Probe

At boot, bind Codex capabilities from the live tool surface, not from a runtime-version label. Probe the expected call shapes and report concrete tool-surface blockers instead of assuming historical labels:

- `«worker.spawn»`: `spawn_agent(task_name,message,fork_turns)`
- `«addressable-worker»`: PRESENT only when a turn-starting route is live: `followup_task(target,message)` or a proven equivalent such as legacy `send_input`; ABSENT when the live surface exposes only `spawn_agent` and `wait_agent`. `send_message(target,message)` is non-triggering only and never makes `«addressable-worker»` PRESENT by itself.
- `«completion-signal»`: `wait_agent(timeout_ms)` for async idle monitoring, while the async mailbox final-status notification remains the signal
- `«roster-reconcile»`: `list_agents(path_prefix?)`
- `«worker.shutdown»`: unresolved unless a shutdown-specific probe proves a binding; Do not bless `interrupt_agent`

Do not infer capabilities from a Codex version name. Do not infer that the turn-starting addressable-worker route is absent from the absence of newer follow-up tools.

## Runtime implementation

- `«worker.spawn»`: A successful `«dispatch.build»` is not a dispatch by itself. For every ready entity, including one just advanced by gate approval, call `spawn_agent(task_name,message,fork_turns="none")` with the helper-emitted prompt unchanged as `message`. Every spawn is a fresh dispatch; deliberate continuity uses `followup_task` with the existing handle. sanitize the helper-emitted `name` to a lowercase digit/underscore `task_name`; retain the mapping from helper identity to sanitized task path, entity slug, stage, and cycle.
  Build fresh Codex dispatches in the named shape: omit `--bare-mode`. A zero-exit build must carry `name` and `prompt`; map them through `CodexMultiAgentV2SpawnInput` to `spawn_agent(task_name,message,fork_turns="none")`.
- `«addressable-worker»`: When PRESENT, a spawned worker remains addressable by live task path/handle while active or completed and still listed, including a completed-but-still-addressable worker. `followup_task(target,message)` is the current turn-starting reuse/advance route — its advance payload is `${SPACEDOCK_BIN:-spacedock} dispatch build --advance`'s emitted `output.prompt`, forwarded verbatim. `send_message(target,message)` is non-triggering context/preservation only. Legacy `send_input` is a fallback only when that surface is actually present. When absent, the `«addressable-worker»` reuse condition fails and feedback re-review fresh-dispatches a separate validation reviewer.
- `«async-dispatch»`: Codex dispatch is async; the spawn call returns a handle. `wait_agent` is asynchronous monitoring: captain input resumes the FO's active loop while workers continue unchanged.
  A fresh spawn does not satisfy a single-entity stop.
  When this worker is the only unresolved work, continue async monitoring until its completion signal, then verify the report and run the immediate gated next action with the reported active entity path as `--artifact` and the workflow's committed `README.md` as the distinct `--reference`.
- `«worker-identity»`: Record Codex task name, mailbox handle, sanitized path, helper-emitted identity, entity slug, stage, cycle, completion epoch, and thread-inherited model when the helper emits null.
- `«completion-signal»`: The single completion signal is the async final-status notification in the FO mailbox. A `wait_agent` return is only an idle observation and never completes the worker: verify the durable report/state first, then immediately run `«dispatch.next-action»()`.
- `«worker.shutdown»`: `«worker.shutdown»` remains unresolved until probed. Durable workflow state remains authoritative. `«addressable-worker»` may carry cooperative preservation text when present.
- `«context-budget»`: ABSENT; its reuse condition is satisfied because Codex has no bound context-budget probe.
- `«roster-reconcile»`: `«roster-reconcile»` may provide active/completed task-path reads when bound. `list_agents(path_prefix?)` provides attribution, stale-cohort classification, debugging, and cleanup targeting.
When `«async-dispatch»` is live, async idle monitoring MUST call `wait_agent(timeout_ms: 300000)` as the long-task default: dispatched work can outlast Codex's short default. Five minutes reduces normal timeout churn; it cannot correct conceptual misuse of idle monitoring when captain-authorized active-scope work is ready. This is Spacedock's per-call policy and does not change Codex's global configuration.

## Codex wait notes

Async idle monitoring applies only after `«dispatch.next-action»()` has performed its one idle/reconcile/retry and finds an active unresolved Codex worker with no dispatchable, gate, mod/PR, or other state work in the captain-authorized scope. In that condition, the FO MUST call `wait_agent(timeout_ms: 300000)` before ending the turn or reporting idle/status. A completed, errored, or absent worker set never qualifies for `wait_agent`; report the loop's explicit stop instead. Captain input resumes the FO's active loop while workers continue unchanged. When the FO becomes idle again, it MUST resume monitoring unresolved workers. Do not preface each monitoring epoch with an interruption disclaimer.

Consecutive waits over the same unresolved worker set and unchanged workflow state are one async idle-monitoring epoch. A wait timeout return is normal and retryable; it means no final-status mailbox update arrived before the deadline. On an ordinary timeout, silently reinstall `wait_agent` only while that worker set, workflow state, and captain-authorized active scope remain unchanged; speak again only on worker completion, a blocker, a workflow-state transition, or a heartbeat the captain explicitly requested.

A wait result must be attributed by mailbox content, task path, `«roster-reconcile»` state when present, or durable workflow state, not by a handle argument; no wait result alone completes a worker. An async final-status notification starts completion verification, not completion itself: verify the durable entity report/state, then immediately run `«dispatch.next-action»()` rather than ending on a completion-only turn or reinstalling idle monitoring. A captain message or shell-out during the wait is operator activity, not idle wake evidence.

### Harness wait-return label

“Wait interrupted by new input” is a harness implementation label only. Do not repeat it as FO language or derive worker cancellation, closure, redispatch, or completion from it.

### Queued notification flushed by later activity

The FO ends the turn without async idle monitoring, and a later captain message, tool action, or shell-out causes Codex to deliver a worker final-status notification that had already been queued. This is queued/activity-driven delivery, not autonomous FO wake-up.

### Autonomous idle FO wake-up

The FO does not use async idle monitoring, performs no later captain message, tool action, shell-out, or terminal job, and Codex starts a new assistant turn from the worker final-status notification alone. This is the only observation that proves no-wait autonomous FO wake-up.

## Feedback reviewer reuse

Feedback rejection is the load-bearing exception to casual fresh dispatch. When `«addressable-worker»` is PRESENT, keep the validation reviewer addressable after its REJECTED report and MUST first re-run the kept-alive validation reviewer through `«addressable-worker»`. Do not fresh-dispatch merely because the reviewer completed or already sent a completion signal. When `«addressable-worker»` is ABSENT, fresh-dispatch a separate validation reviewer for cycle 2. After a PASSED re-review, re-enter the normal gate flow and advance or terminalize from durable state.

## Captain Interaction

The captain is the user of the Codex session. Communicate gate results, clarifications, and status directly in the conversation.
