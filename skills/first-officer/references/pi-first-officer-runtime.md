# Pi First Officer Runtime

This file defines how the shared first-officer core executes on Pi.

## Runtime Shape

Pi is a first-class runtime target, but it does not expose Claude Code team-tool signatures. Do not call or ask workers to call Claude team tools. Pi dispatch uses a Pi-native substrate selected by the launch/test harness:

- default: `pi-subagents`, where the parent first officer uses the `subagent(...)` tool to run a bounded ensign assignment and observes the returned result as completion evidence;
- optional: `pi-agent-teams`, where an adapter maps Spacedock lifecycle intents to the `teams` tool actions (`member_spawn` or `delegate`, `message_dm`, `member_shutdown`, `team_done`).

The Spacedock contract talks in terms of dispatch, completion, follow-up, and shutdown. The Pi adapter owns how those lifecycle events map to the active substrate.

## Dispatch

Use `spacedock dispatch build` with `host: "pi"` in the input JSON. Forward the emitted dispatch file prompt to the Pi worker without rewriting it into Claude syntax. For the first slice, dispatch should normally run with `bare_mode: true`: the Pi subagent result is the completion signal, so no `team_name` is required.

A Pi first officer may dispatch via `subagent(...)` when that tool is available. The task must include the emitted dispatch-file prompt or the dispatch file content, the workflow directory, the entity path, and the completion checklist. The child must load the Spacedock ensign skill and Pi ensign runtime adapter before working.

## Awaiting Completion

For `pi-subagents`, the completion signal is the subagent result returned to the parent. After the result arrives, read the entity file and verify the stage report exactly as the shared core requires. Do not advance state based only on a cheerful worker summary.

For `pi-agent-teams`, completion is observed through the adapter's task/member notification and then verified against the entity file. The adapter should expose a clear completed/failed result to the first officer; the entity stage report remains the source of truth.

## Follow-up and Reuse

Fresh redispatch is the default safe behavior for the first Pi slice. If a Pi substrate exposes a resumable worker handle, record only the minimum metadata needed to prevent stale reuse mistakes: worker label, substrate, run/session handle, entity slug, stage, state, and completion epoch. A follow-up assignment must increment the epoch, and a previous completion must never satisfy the new epoch.

## Shutdown

For `pi-subagents`, a completed child invocation needs no mailbox shutdown. Mark the worker complete/closed in first-officer memory and continue.

For `pi-agent-teams`, use the adapter's lifecycle mapping to request member shutdown or end the team run. Do not emulate Claude team deletion.

## Live Harness Isolation

Live Pi tests should run with an isolated Pi config directory and an isolated session directory. The harness may copy the operator's existing Pi auth file into the isolated config directory so OAuth/subscription credentials are reused without sharing global sessions, packages, or settings. This mirrors the Codex live runner pattern: isolate runtime state, reuse credentials only.

The durable proof for Pi support is not transcript phrasing. A valid live proof dispatches a Pi ensign against a temp split-root workflow and verifies exit code, state checkout file changes, git log, and stage report content.
