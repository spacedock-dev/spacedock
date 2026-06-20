---
title: Codex multi_agent_v2 runtime dispatch support
status: ideation
source: "FO live dispatch observation (2026-06-20): while dispatching status-validate-determinism under Codex multi_agent_v2, the current Codex runtime adapter no longer matches the live tool surface. v2 exposes spawn_agent(task_name,message,fork_turns), list_agents, wait_agent(timeout_ms), send_message, followup_task, and interrupt_agent. The shipped adapter still names send_input, wait_agent(handle), and hyphenated dispatch-build names as direct handles."
score:
started: 2026-06-20T03:58:46Z
completed:
verdict:
worktree:
issue:
sprint:
sprint-readiness:
id: f2r8cnyxj9pf24xrsf71szb0
---

# Codex multi_agent_v2 runtime dispatch support

## Problem

The Codex first-officer and ensign runtime references describe the older collab surface:

- initial dispatch uses `spawn_agent`, but `dispatch build` emits hyphenated `name` values such as `spacedock-ensign-status-validate-determinism-implementation`; live v2 `spawn_agent.task_name` rejects hyphens and accepts only lowercase letters, digits, and underscores
- feedback/reuse/shutdown text names `send_input`, but the live surface exposes `send_message` and `followup_task`
- foreground waiting text says `wait_agent(handle)`, but the live surface exposes global `wait_agent(timeout_ms)` and captain messages interrupt the wait without killing the worker
- the no-wait probe still did not prove autonomous idle wake-up; completion surfaced after later root activity, so the correct classification is queued/activity-driven delivery
- `list_agents` is now a usable roster-like read for active and completed task paths

The runtime contract needs a forward-compatible v2 variant or adapter section that captures these semantics without breaking existing v1/live-test assumptions before they are intentionally retired.

## Acceptance criteria

**AC-1 - Initial Codex v2 dispatch maps helper output to the live spawn surface.**  
Verified by: a test or live fixture showing `dispatch build` output with a hyphenated `name` is converted to a valid v2 `task_name` while forwarding the helper `prompt` unchanged as `message`, omitting unsupported `description` and `model` fields.

**AC-2 - Codex v2 continuation and steering use the live message tools.**  
Verified by: adapter text and tests distinguish `send_message` context delivery from `followup_task` turn-triggering reuse/feedback, and remove or explicitly version-gate `send_input` from v2 instructions.

**AC-3 - Codex v2 waiting semantics match observed behavior.**  
Verified by: adapter text and live/fixture evidence state that `wait_agent(timeout_ms)` is global, captain input interrupts wait without failing the worker, and no-wait completion remains queued/activity-driven unless an autonomous wake probe proves otherwise.

**AC-4 - Terminal/supersede shutdown has a concrete v2 call.**  
Verified by: adapter text maps cooperative shutdown to the live v2 surface, likely `interrupt_agent` or a `send_message`/`followup_task` shutdown protocol depending on what probes prove, and names residual risk.

**AC-5 - Old Codex runtime assumptions are preserved or deliberately retired.**  
Verified by: existing Codex live/shared tests either remain green under a versioned v1 path or are updated with v2-specific expectations; stale assertions around `send_input` and `wait_agent(handle)` are not left ambiguous.

## Evidence to seed ideation

- `spawn_agent` rejected helper name `spacedock-ensign-status-validate-determinism-ideation` with `agent_name must use only lowercase letters, digits, and underscores`.
- Retrying as `spacedock_ensign_status_validate_determinism_ideation` succeeded with the same helper prompt.
- `list_agents` showed `/root`, completed probe workers, and the running Spacedock ensign task.
- Captain messages interrupted `wait_agent`, but the worker remained running and could later complete.
- A no-wait window did not wake root autonomously; completion became visible after a later `list_agents` call and mailbox notification.
- `send_message` successfully delivered a non-interrupting note to the running worker, which then committed and finalized.

## Test plan

- Add fixture tests for Codex v2 spawn mapping and task-name sanitization.
- Add fixture tests or contractlint for v2 wording: no unversioned `send_input`, no `wait_agent(handle)`, and explicit queued-vs-autonomous completion classification.
- Add or update live Codex shared-scenario assertions so v2 collab events are graded by durable state plus the new tool names.
- Run `go test ./...`, `go test ./... -race`, and the relevant Codex live lane if runtime behavior is changed.
