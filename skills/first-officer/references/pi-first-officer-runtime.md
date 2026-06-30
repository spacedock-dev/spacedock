# Pi First Officer Runtime

This file defines how the shared first-officer core executes on Pi. The shared core owns invocation timing; this adapter owns the Pi realization. Pi-native substrate selected by the launch/test harness determines whether the active transport is `pi-subagents` or `pi-agent-teams`.

## Runtime implementation

- `«worker.spawn»`: Use `spacedock dispatch build` with `host: "pi"` and forward the emitted assignment source of truth to the worker. For `pi-subagents`, call `subagent(...)` with explicit `context: "fresh"` and `cwd: <resolved repo root>`; add only Pi transport metadata such as optional human-facing phase or label values. For `pi-agent-teams`, map initial creation to `member_spawn` or `delegate`. For standing injection, `pi-subagents` is no-op and `pi-agent-teams` may map to `member_spawn`.
- `«addressable-worker»`: PRESENT. Fresh redispatch remains the default first Pi slice; normal follow-up and retry dispatches are fresh assignment cycles, not context resumes. Reuse-advance over a kept-alive handle is deferred; a non-fresh resume is only an explicit manual/debug exception tied to durable metadata and the full `«worker-identity»` schema. On `pi-agent-teams`, follow-up steering maps to `message_dm`.
- `«async-dispatch»`: `pi-subagents` runs `subagent(... async: true)` and polls the run id; `pi-agent-teams` maps async lifecycle to its team/member task substrate. do not use the `subagent(... acceptance: ...)` contract; Spacedock acceptance requirements stay in the dispatch content and are verified by state/report evidence.
- `«worker-identity»`: Record worker address, substrate, run/session handle, entity slug, stage, epoch, completion state, and stamped model. When the helper emits null model, stamp the parent's live model via `intercom({action:"list"})`; stage-declared model overrides are passed unchanged. Pi's model-space binding is provider/model strings, including provider-qualified and `~`-prefixed Pi-valid model names.
- `«completion-signal»`: For `pi-subagents`, the primary completion signal is the child return/status; optional advisory is only a heads-up. file verification remains the completion gate: the FO reads the entity file and verifies the stage report before advancing state. For `pi-agent-teams`, task/member completion is likewise verified against the entity file.
- `«worker.shutdown»`: For `pi-subagents`, a completed child invocation needs no mailbox shutdown; mark the worker complete/closed in first-officer memory. For `pi-agent-teams`, map teardown to `member_shutdown` or `team_done` according to the active adapter lifecycle.
- `«context-budget»`: ABSENT; reuse-condition-0 is satisfied for Pi.
- `«roster-reconcile»`: ABSENT; Pi relies on durable entity state and adapter-held worker identity, not a shared roster sweep.

The build artifact carries the entity slug/name, entity path, workflow directory, target stage, stage definition fetch command, worktree path when applicable, completion checklist, and completion-signal wording. It must not be replaced by a locally composed assignment. The model stamped through `«worker-identity»` is a Pi-native value used by reuse-condition-4.

## Live Harness Isolation

Live Pi tests should run with an isolated Pi config directory and an isolated session directory. The harness may copy the operator's existing Pi auth file into the isolated config directory so OAuth/subscription credentials are reused without sharing global sessions, packages, or settings.

The durable proof for Pi support is not transcript phrasing. A valid live proof dispatches a Pi ensign against a temp split-root workflow and verifies process exit, state checkout file changes, git log, and stage report content.

## Bridge egress (FO liveness/activity → Bridge)

Bridge reads FO liveness and activity from `_bridge/events.jsonl` and `_bridge/sessions/<session_id>.json`. The events.jsonl line schema is the Spacedock-owned contract `{"ts","event","session_id","agent_id","agent_type","detail":{"tool","source"}}` (full surface: `docs/dev/bridge-egress-contract.md`); each host binds its own producer for it. On Pi this producer is not yet established.

- **FO event emission** — ABSENT/TODO on Pi: there is no established producer (a Pi event API for emitting the events.jsonl contract line is unconfirmed). Consequence: Bridge gets no `events.jsonl` and no `_bridge/sessions/` markers on a Pi host, so it shows neither live tool activity nor FO-vs-ensign attribution; liveness falls back to the bridge-inbox heartbeat (`_bridge/fo.$SLUG.json`) the FO writes each tick. Binding a Pi producer that emits the contract line is the open work.
- **«session-id» binding** — the bridge-inbox heartbeat resolves the neutral `SD_SESSION_ID`; on Pi this stays empty (no stable per-session id is exposed), so the heartbeat carries an empty session id — still a valid liveness tick (Bridge reads freshness from `ts`). Bind it once Pi exposes a stable per-session id.
