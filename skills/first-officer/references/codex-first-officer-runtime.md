# Codex First Officer Runtime

This file defines how the shared first-officer core executes on Codex. The host-neutral dispatch and merge procedures are in `references/fo-dispatch-core.md` / `fo-merge-core.md` (named by the boot-resident core); this file is the Codex parts those defer to.

## Cleanup Capability Bindings

`«worker.shutdown»` remains unresolved until probed. `«roster-reconcile»` may provide active/completed task-path reads when bound. `«addressable-worker»` may carry cooperative preservation text when present. Durable workflow state remains authoritative.

`«context-budget»` is unavailable unless a future probe binds it.

## Live Tool Surface Probe

At boot, bind the Codex `«fn»`s from the actual tools exposed in this session, not from a runtime-version label. The expected live bindings are:

- `«worker.spawn»` when `spawn_agent(task_name,message,fork_turns)` is present.
- `«addressable-worker»` when either live binding is present: `send_message(target,message)` plus `followup_task(target,message)`, or the compatibility `send_input` surface. `send_message` is non-triggering context/preservation text. `followup_task` starts a worker turn on the addressed worker. `send_input` is both the compatibility address and turn-starting route. A turn-starting binding can address a completed-but-still-addressable worker; use it for critical-path reuse and re-review.
- `«completion-signal»` foreground waiting when `wait_agent(timeout_ms)` is present.
- `«roster-reconcile»` when `list_agents(path_prefix?)` is present.
- `«worker.shutdown»` remains unresolved. Do not bless `interrupt_agent` until a shutdown-specific probe proves it is a binding.

If a referenced tool is missing or rejects the call shape, report that concrete tool-surface blocker and fall back only to behavior the live surface actually supports. Do not infer capabilities from a Codex version name or from historical transcript labels.

## Dispatch

The `fo-dispatch-core.md` `## Dispatch Adapter` defers initial worker dispatch to `«worker.spawn»`. Assemble the dispatch input through the core's
mandatory `spacedock dispatch build` flow (write checklist, scope notes, and
feedback context into files, then the flag/file form); Codex host is normally
derived from `CODEX_THREAD_ID`, so pass `--host codex` only for deliberate tests or
cross-host tooling.

Forward the emitted prompt exactly as returned; for Codex it is a
read-dispatch-file prompt. The FO must never forward a prompt containing `Skill(skill="spacedock:ensign")` to a Codex worker.

When `«worker.spawn»` is bound, sanitize the helper-emitted `name` through `«worker-identity»` to a lowercase digit/underscore task name, pass the emitted `prompt` unchanged as the worker message, omit unsupported helper fields, and record the sanitized task path as the live worker handle. Helper-emitted `name` remains the host-neutral worker identity source; retain the mapping from sanitized task path back to entity slug, stage, and cycle so cohorts and teardown targets remain classifiable.

Use Codex task names and mailbox notifications as the worker handle. Codex dispatch prompts carry the file-pointer completion contract emitted by `spacedock dispatch build`.

## Reuse And Feedback Routing

When the live tool surface exposes addressable workers, a spawned worker remains addressable by its live task path/handle while it is active or completed and still listed by the runtime. Use `«addressable-worker»` for non-triggering context notes, cooperative preservation text, reuse, feedback, or rework that must start a worker turn according to the binding above. Reuse is appropriate when the next task depends heavily on the worker's retained context; otherwise dispatch a fresh worker from entity state.

Feedback rejection is the load-bearing exception to casual fresh dispatch: when the live surface provides `«addressable-worker»`, keep the validation reviewer addressable after its REJECTED report, route the fix to the implementation worker, then re-run the kept-alive validation reviewer through the same turn-starting `«addressable-worker»` binding. A completed validation report does not by itself make that reviewer unavailable; first attempt reuse with the captured validation-worker handle. Do not fresh-dispatch cycle-2 validation merely because the run is headless, scoped to one entity, running on Codex, or the reviewer already sent a completion signal. Do not infer that the turn-starting addressable-worker route is absent from the absence of newer follow-up tools, from the absence of earlier follow-up events, from a transcript that has only used spawn/wait so far, or from a historical wait event label. Attempt the probed `«addressable-worker»` binding; if the live tool surface rejects the call, report that concrete blocker instead of silently fresh-dispatching cycle-2 validation. Fresh-dispatch the reviewer only when the existing reviewer is no longer addressable or reuse conditions fail.

## Awaiting Completion

The Codex completion signal is the async final-status notification in the FO mailbox.
The live foreground wait action is not the completion signal itself.

Before using foreground wait, finish the available dispatch sequence: process
ready final-status notifications, make gate decisions, apply resulting state
transitions, and dispatch any newly dispatchable work. Use foreground wait only
when no other dispatchable or gate-processing work is available and an
unresolved worker completion is the next useful idle action. Non-critical
background work may still end the FO turn and rely on the mailbox.

### Foreground wait

When `«completion-signal»` is bound, it may use global foreground wait only after workflow and gate work is exhausted. A wait timeout return is normal and retryable. It means no final-status mailbox update arrived before the deadline; it is not a failure, a zombie signal, or a teardown trigger. Before foreground waiting, tell the captain that pressing Esc or otherwise causing an operator interruption only returns control; the worker is not failed, closed, or redispatched, and the next foreground wait is reinstalled when waiting is again the next useful idle action. A wait return must be attributed by mailbox content, task path, `«roster-reconcile»` state when present, or durable workflow state, not by a handle argument. A captain message or shell-out during the wait is operator activity, not idle wake evidence.

### Queued notification flushed by later activity

The FO does not use foreground wait, ends the turn, and a later captain message,
tool action, or shell-out causes Codex to deliver a worker final-status
notification that had already been queued. This proves mailbox ordering and is
classified as queued/activity-driven delivery, not autonomous FO wake-up.

### Autonomous idle FO wake-up

The FO does not use foreground wait, performs no later captain message, tool
action, shell-out, or terminal job, and Codex starts a new assistant turn from
the worker final-status notification alone. This is the only observation that
proves no-wait idle wake-up.

Between dispatch and the async final-status notification, do not poll, re-dispatch, or close the worker.
Continue ready workflow work, foreground-wait only under the rule above, or end
the turn and let Codex deliver the mailbox notification. Close workers only after
completion has been observed or when the captain explicitly requests shutdown.

## Completion Signal

An ensign completes by sending a concise final message in its Codex worker
thread. The FO observes that message through the Codex final-status notification
in the FO mailbox and then resumes the shared event loop.

## Captain Interaction

The captain is the user of the Codex session. Communicate gate results,
clarifications, and status directly in the conversation.

## Roster Reconcile Binding

When `«roster-reconcile»` is bound, Codex has active/completed task-path reads for attribution, stale-cohort classification, debugging, and cleanup targeting. Durable workflow state remains authoritative.
