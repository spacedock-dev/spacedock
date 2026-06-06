---
id: 82rb21h6ynqqawbjbcgw78t3
title: Support Codex followup_task worker reuse
status: backlog
source: "captain request 2026-06-05: file backlog task after multi_agent_v2 followup_task probe"
worktree: ""
---
# Support Codex followup_task worker reuse

Codex multi-agent support now exposes a generic mailbox surface in this session:
`spawn_agent`, `followup_task`, `send_message`, `wait_agent`, `close_agent`, and
`list_agents`. The important new distinction is that `send_message` delivers a
message without triggering a new worker turn, while `followup_task` sends work to
an existing non-root worker and triggers a follow-up turn.

The current Spacedock Codex runtime text and tests still describe reusable worker
routing in terms of the older `send_input` concept. When the Codex agent surface
matures and the final API naming/semantics are stable, update Spacedock's Codex
adapter and runtime tests so feedback routing can intentionally reuse the prior
worker handle with the active follow-up primitive.

The motivating workflow is validation feedback routed back to implementation:
after a validation rejection with `feedback-to: implementation`, the first
officer should send concrete fix work to the existing implementation worker when
reuse conditions allow, wait for that worker's new completion, then re-engage the
kept-alive validation worker for re-review when appropriate. A passive
`send_message` must not be treated as a worker continuation.

## Acceptance criteria

**AC-1** Codex runtime reuse uses the active follow-up primitive.
Verified by a Codex runtime fixture or live transcript assertion that observes a
completed worker receiving follow-up work through the active continuation tool
(`followup_task` or its stable successor), not through passive `send_message`.

**AC-2** Feedback-to-implementation reuse waits for fresh completion evidence.
Verified by a rejection-flow fixture that fails if the first officer accepts the
implementation worker's stale prior completion after routing follow-up work.

**AC-3** Validation re-review reuse keeps the reviewer handle correlated.
Verified by a Codex transcript assertion that permits exactly one validation
spawn for a two-cycle rejection flow and requires the cycle-2 re-review to target
that original validation worker handle; the assertion must fail on a fresh
cycle-2 validation spawn or on follow-up routed only to the implementation
worker.

**AC-4** Runtime compatibility is explicit while the API is unstable.
Verified by tests that cover the available continuation surface: old traces using
`send_input`, current sessions exposing `followup_task`, and sessions where no
active continuation primitive exists. The no-primitive case must fresh-dispatch
or block loudly rather than silently using passive messaging as reuse.

## Stage test gates

- Ideation should begin with a small no-write Codex probe that proves the final
  continuation primitive actually starts a new worker turn after completion.
- Implementation should update the Codex first-officer runtime instructions,
  feedback-rejection routing text, and reuse assertions together so the wording
  and behavioral tests agree.
- Validation should run focused reuse tests plus the relevant shared rejection
  flow scenario or its cheapest fixture-backed equivalent.
