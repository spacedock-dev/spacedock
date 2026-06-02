# Codex First Officer Runtime

This file defines how the shared first-officer core executes on Codex.

## Team Creation

Codex does not expose Claude Code's `TeamCreate` registry. Do not create or
delete teams, and do not inspect Claude team directories. Codex workers are
spawned and observed through the multi-agent mailbox surface.

Codex declares none for the context-budget probe. When the shared core asks for a
reuse budget signal, treat the probe as unavailable and prefer a fresh dispatch
unless the current task explicitly depends on the previous worker's context.

## Dispatch

Use `spawn_agent` for initial worker dispatch. Assemble the dispatch input with
`spacedock dispatch build` and pass `host: "codex"` in the JSON envelope. Forward
the emitted prompt exactly as returned; for Codex it is a read-dispatch-file
prompt and must not be rewritten into `Skill(skill=...)`.

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

The Codex completion signal is the async final-status notification in the FO mailbox. `wait_agent` is only an optional accelerator when a critical-path step is blocked on the worker result.

A wait_agent timeout return is normal. It means no final-status mailbox update
arrived before the deadline; it is not a failure, a zombie signal, or a teardown
trigger.

Between dispatch and the async final-status notification, do not poll, re-dispatch, or close the worker. End the turn and let Codex deliver the mailbox notification. Close workers only after completion has been observed or when the captain explicitly requests shutdown.

## Completion Signal

An ensign completes by sending a concise final message in its Codex worker
thread. The FO observes that message through the Codex final-status notification
in the FO mailbox and then resumes the shared event loop.

## Captain Interaction

The captain is the user of the Codex session. Communicate gate results,
clarifications, and status directly in the conversation. Do not invent a
team-lead mailbox on Codex.
