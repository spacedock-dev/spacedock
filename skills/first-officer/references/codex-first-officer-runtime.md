# Codex First Officer Runtime

This file defines how the shared first-officer core executes on Codex.

## Dispatch seam

The host-neutral dispatch machinery — the per-entity dispatch procedure, worker resolution, the dispatch-adapter assembly (`spacedock dispatch build` → spawn call), the reuse contract, worktree ownership, and the event-loop skeleton — lives in `references/fo-dispatch-core.md`. Read it at the first dispatch. This file supplies the Codex specifics the core defers to: the spawn call is `spawn_agent`, the reuse-advance handle is `send_input`, the completion signal is the mailbox final-status notification, and Codex declares no context-budget probe (reuse-condition-0 is satisfied as unavailable → prefer fresh).

## Merge seam

The host-neutral merge-and-cleanup ceremony — the set→invoke→clear mod-block sequence, the Ship-Local ceremony, worktree-removal safety, and Mod-Block Enforcement — lives in `references/fo-merge-core.md`. Read it at the terminal boundary. This file supplies the Codex half of Merge-and-Cleanup step 10 below.

### Terminal Teardown (Codex — Merge-and-Cleanup step 10)

This realizes Merge-and-Cleanup step 10 (`fo-merge-core.md`) on Codex. Codex has no team registry and no `TeamDelete`, so there is no bounded team-teardown loop, settle interval, attempt cap, or terminal-status marker — those are Claude specifics, absent here. The obligation is still mandatory at the terminal boundary:

10. **Teardown agents at terminal.** Derive the entity's worker cohort — every Codex worker whose task name decomposes to this entity's slug. Issue the cooperative-shutdown call over the mailbox surface (best-effort, fire-and-forget) and drop them from session memory. There is no team object to delete and no roster to settle, so teardown is a single cooperative-shutdown pass, not a retry loop. Teardown is mandatory whether the merge ran locally or via a PR host. Residual risk: per `## Backstop (Codex)`, there is no reconcile sweep, so a skipped teardown stays missed until session end — the boundary step is the only enforcement, so do not skip it.

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

Codex has no team registry, so there is no `team_name` lifecycle to create,
recover, or tear down. Use Codex task names and mailbox notifications as the
worker handle. Do not use Claude `SendMessage` completion syntax in Codex
dispatch prompts.

## Reuse And Feedback Routing

Route feedback or continuation to an existing Codex worker with `send_input`.
Set `interrupt=true` only when the worker must immediately change direction.
Reuse is appropriate when the next task depends heavily on the worker's retained
context; otherwise dispatch a fresh worker from entity state.

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

The FO calls `wait_agent(handle)` as the next useful idle action after the
scheduling priority above is exhausted. A wait_agent timeout return is normal and
retryable with the same handle. It means no final-status mailbox update arrived
before the deadline; it is not a failure, a zombie signal, or a teardown trigger.
Before calling `wait_agent`, tell the captain that pressing Esc or otherwise
causing an operator interruption only returns control; the worker is not failed,
closed, or redispatched, and the next foreground wait retries the same handle. A
captain message or shell-out during the wait is operator activity, not idle wake
evidence.

### Queued notification flushed by later activity

The FO does not call `wait_agent`, ends the turn, and a later captain message,
tool action, or shell-out causes Codex to deliver a worker final-status
notification that had already been queued. This proves mailbox ordering, not
autonomous FO wake-up.

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

None. There is no `spacedock dispatch reconcile` analog on Codex — no roster source, no drift classifier, no periodic sweep. A missed terminal teardown or supersede shutdown stays missed until session end; the shared-core boundary steps are the only enforcement.
