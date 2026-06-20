# Codex First Officer Runtime

This file defines how the shared first-officer core executes on Codex. The host-neutral dispatch and merge procedures are in `references/fo-dispatch-core.md` / `fo-merge-core.md` (named by the boot-resident core); this file is the Codex parts those defer to.

## Terminal Teardown (Merge-and-Cleanup step 10)

Codex has no team registry and no `TeamDelete`, so there is no bounded team-teardown loop, settle interval, attempt cap, or terminal-status marker — those are Claude specifics, absent here. The obligation is still mandatory at the terminal boundary.

Codex multi_agent_v2 leaves `«worker.shutdown»` unresolved until probed. Do not bless `interrupt_agent` as a shutdown binding yet; first prove whether it terminates, pauses, or leaves a worker addressable. Keep cooperative preservation text separate from any hard interruption.

10. **Teardown agents at terminal.** Derive the entity's worker cohort — every Codex worker whose task name decomposes to this entity's slug. On Codex multi_agent_v2, use `«roster-reconcile»` through `list_agents(path_prefix?)` to inspect active/completed task paths for stale-cohort classification and terminal cleanup targeting; durable workflow state remains authoritative for advancement. Until `«worker.shutdown»` has a probed binding, send only cooperative preservation-aware shutdown text over `send_message(target,message)` and drop the worker from session memory after completion or explicit captain direction. There is no team object to delete and no roster to settle, so teardown is a single cooperative pass, not a retry loop. Teardown is mandatory whether the merge ran locally or via a PR host.

## Team Creation

Codex does not expose Claude Code's `TeamCreate` registry. Do not create or
delete teams, and do not inspect Claude team directories. Codex workers are
spawned and observed through the multi-agent mailbox surface.

Codex declares none for the context-budget probe. When the shared core asks for a
reuse budget signal, treat the probe as unavailable and prefer a fresh dispatch
unless the current task explicitly depends on the previous worker's context.

## Dispatch

The spawn call `fo-dispatch-core.md` `## Dispatch Adapter` defers to is `spawn_agent`
for initial worker dispatch. Assemble the dispatch input through the core's
mandatory `spacedock dispatch build` flow (write checklist, scope notes, and
feedback context into files, then the flag/file form); Codex host is normally
derived from `CODEX_THREAD_ID`, so pass `--host codex` only for deliberate tests or
cross-host tooling.

Forward the emitted prompt exactly as returned; for Codex it is a
read-dispatch-file prompt. The FO must never forward a prompt containing `Skill(skill="spacedock:ensign")` to a Codex worker.

For Codex multi_agent_v2, bind `«worker.spawn»` to `spawn_agent(task_name,message,fork_turns)`: sanitize the helper-emitted `name` through `«worker-identity»` to a lowercase digit/underscore `task_name`, pass the emitted `prompt` unchanged as `message`, omit unsupported `description`, `subagent_type`, and `model`, and record the sanitized task path as the live worker handle. Helper-emitted `name` remains the host-neutral worker identity source; retain the mapping from sanitized task path back to entity slug, stage, and cycle so cohorts and teardown targets remain classifiable.

Codex has no team registry, so there is no `team_name` lifecycle to create,
recover, or tear down. Use Codex task names and mailbox notifications as the
worker handle. Do not use Claude `SendMessage` completion syntax in Codex
dispatch prompts.

Codex has no shared standing-teammate surface; the core's standing-injection call is a no-op.

## Reuse And Feedback Routing

On Codex multi_agent_v2, `«addressable-worker»` is realized by `send_message(target,message)` for non-triggering context notes and cooperative preservation text, and by `followup_task(target,message)` for reuse, feedback, or rework that must start a worker turn. Reuse is appropriate when the next task depends heavily on the worker's retained context; otherwise dispatch a fresh worker from entity state. Legacy pre-v2 input-routing terms belong only to explicitly versioned legacy fixtures or references.

## Awaiting Completion

The Codex completion signal is the async final-status notification in the FO mailbox.
`wait_agent` is an explicit foreground wait action, not the completion signal
itself.

Before calling `wait_agent`, finish the available dispatch sequence: process
ready final-status notifications, make gate decisions, apply resulting state
transitions, and dispatch any newly dispatchable work. Call `wait_agent` only
when no other dispatchable or gate-processing work is available and an
unresolved worker completion is the next useful idle action. Non-critical
background work may still end the FO turn and rely on the mailbox.

### Foreground wait

On Codex multi_agent_v2, `«completion-signal»` may use global `wait_agent(timeout_ms)` only after workflow and gate work is exhausted. A wait timeout return is normal and retryable. It means no final-status mailbox update arrived before the deadline; it is not a failure, a zombie signal, or a teardown trigger. Before calling `wait_agent`, tell the captain that pressing Esc or otherwise causing an operator interruption only returns control; the worker is not failed, closed, or redispatched, and the next foreground wait is reinstalled when waiting is again the next useful idle action. A wait return must be attributed by mailbox content, task path, `list_agents(path_prefix?)` roster state, or durable workflow state, not by a handle argument. A captain message or shell-out during the wait is operator activity, not idle wake evidence.

### Queued notification flushed by later activity

The FO does not call `wait_agent`, ends the turn, and a later captain message,
tool action, or shell-out causes Codex to deliver a worker final-status
notification that had already been queued. This proves mailbox ordering and is
classified as queued/activity-driven delivery, not autonomous FO wake-up.

### Autonomous idle FO wake-up

The FO does not call `wait_agent`, performs no later captain message, tool
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
clarifications, and status directly in the conversation. Do not invent a
team-lead mailbox on Codex.

## Backstop (Codex)

Codex multi_agent_v2 has `«roster-reconcile»` through `list_agents(path_prefix?)` for active/completed task-path reads. It can support stale-cohort classification, debugging, and terminal cleanup targeting, but it is not a durable workflow-state source and does not replace state-checkout evidence. A missed terminal teardown or supersede shutdown can still persist until session end if the boundary step is skipped; the shared-core boundary steps remain the enforcement point.
