# Bridge egress contract (FO liveness → Bridge)

The Bridge seam is two flows over the shared `_bridge/` dir (resolved from the repo root where the FO launched). They run in opposite directions and are documented apart:

- **Ingress** — captain intent → FO — is the `bridge-inbox` mod (`docs/dev/_mods/bridge-inbox.md`). It drains `_bridge/inbox.jsonl` on the FO's own loop ticks; it rides the portable Spacedock mod-hook loop and is already host-neutral.
- **Egress** — FO liveness/activity → Bridge — is *this* contract: the four `_bridge/` files a running FO (and its ensigns) write for the external Bridge command-center UI to tail. Bridge reads what already exists; it is not a source of truth.

**The schema is Spacedock-owned and harness-neutral. The *producer* of the host-event files (`events.jsonl` + the session marker) is per-host** — the heartbeat and feed are written host-neutrally by the `bridge-inbox` mod. The per-host producers are bound the same way runtime lifecycle capabilities are (see `fo-dispatch-core.md`'s "Claude: PRESENT / Codex: ABSENT / Pi: …" idiom): the shared contract names the file shape; the host adapter (`skills/first-officer/references/{claude,codex,pi}-first-officer-runtime.md`, under their `## Bridge egress` section) owns the concrete mechanism that emits it. The current implementation uses the shared `spacedock bridge egress emit --host <host>` command for host normalization. Claude, Codex, and Pi have packaged event producers; deterministic session→entity marker parity is proven only for Claude so far.

All four files are gitignored session runtime (`.gitignore` carries `_bridge/`), append-only or last-write where noted, and strictly observe-only: a telemetry side-channel must never block, fail, or alter the FO.

## Session-id binding `«session-id»`

Two of the surfaces stamp the FO/ensign session id. The neutral shell token is **`SD_SESSION_ID`**, owned by the host adapter's `## Bridge egress` binding. The bridge-inbox heartbeat reads it neutral-first with a built-in per-host fallback, so it never silently blanks and needs no per-tick `export`:

```
"${SD_SESSION_ID:-${CLAUDE_CODE_SESSION_ID:-${CODEX_THREAD_ID:-}}}"
```

- → **Claude:** resolves to `$CLAUDE_CODE_SESSION_ID` (set by Claude Code in the session). · **Codex:** resolves to `$CODEX_THREAD_ID` when the host exposes it to the FO shell. · **Pi:** ABSENT → empty.

`SD_SESSION_ID` is honored by the **heartbeat** producer only. The Claude event producer stamps `events.jsonl` and the session marker from its hook payload's `.session_id`, *not* from `SD_SESSION_ID` — so on Claude leave `SD_SESSION_ID` unset (the fallback already yields `$CLAUDE_CODE_SESSION_ID`, the same id the payload carries). Setting `SD_SESSION_ID` to a *different* value would desync the heartbeat from the event stream and break the join the contract exists to enable; an explicit override is only safe when it equals the host's own session id.

An empty value is honest, not a bug: the heartbeat is still a valid liveness tick (Bridge reads freshness from `ts`); only the join between this heartbeat and the event stream is unavailable.

## `_bridge/events.jsonl` — FO/ensign activity stream

One JSON object per line, appended on each lifecycle event. Liveness, not content: no tool inputs/outputs, no prompt text.

```
{"timestamp":"<rfc3339>","ts":"<rfc3339>","host":"<claude|codex|pi>","event":"<canonical Bridge lifecycle event>","session_id":"<«session-id»>","agent_id":"<id, ensign subagents only>","agent_type":"<e.g. spacedock:ensign, empty for the main FO>","actor_id":"<host-scoped actor id>","detail":{"tool":"<tool name, when available>","source":"<event source, when available>"}}
```

- `agent_id`/`agent_type` are empty for the main FO session and set for ensign subagents, so Bridge tells FO-vs-ensign activity apart.
- `event` is normalized by `spacedock bridge egress emit` into Bridge's canonical grammar before it reaches `events.jsonl`: `SessionStart`, `UserPromptSubmit`, `PostToolUse`, `Notification`, `Stop`, `SubagentStart`, or `SubagentStop`. Host-native names must not leak into this file; for example Pi `session_shutdown`/`turn_end` normalize to `Stop`, `agent_end` normalizes to `SubagentStop`, and `tool_execution_end` normalizes to `PostToolUse`.
- Best-effort size cap: the producer trims to the most recent lines past a bound (it is liveness, not a ledger; a few lost lines to a concurrent trim are acceptable).
- → **Claude:** PRESENT — `hooks/hooks.json` registers `scripts/spacedock-bridge-events.sh` for Claude lifecycle events (all async). The wrapper delegates to the shared `spacedock bridge egress emit --host claude` command so Bridge never couples to Claude's internal transcript format.
- → **Codex:** PACKAGED/FIXTURE-COVERED — `.codex-plugin/plugin.json` points at `hooks/codex-hooks.json`, whose non-async command hooks invoke `scripts/codex-bridge-events.sh`; the wrapper delegates to `spacedock bridge egress emit --host codex`. This proves packaging and minimal lifecycle payload handling, not live marker parity.
- → **Pi:** PACKAGED/SOURCE-COVERED — `package.json` advertises `.pi/extensions/spacedock.ts`, and the Pi extension forwards lifecycle event payloads to `spacedock bridge egress emit --host pi`. Local tests cover registration/source wiring and package discovery, not a live Pi run.

## `_bridge/fo.$SLUG.json` — per-workflow liveness heartbeat

Last-write (one object, not a stream), one file per workflow slug. Bridge treats the workflow as live only when `ts` is fresh (within 30 minutes) and not in the future.

```
{"session_id":"<«session-id»>","ts":"<rfc3339 UTC>","state":"idle"}
```

- `state` is `idle`: the producer runs at startup/idle boundaries (between dispatches), so it cannot honestly claim `working` — the finer working/idle signal lives in `events.jsonl`.
- → **All hosts (host-neutral producer):** the `bridge-inbox` mod writes this each startup/idle tick (it rides the portable hook loop). The only per-host part is `«session-id»` — see the binding above. The heartbeat still attaches the FO on every host.

## `_bridge/fo-feed.jsonl` — fleet-history narration

One JSON object per line, appended when the FO dispatches, advances, or completes an entity. Drives Bridge's fleet-history rail for a local-only workflow (entities gitignored, so the `dispatch:`/`advance:` git narration is empty).

```
{"ts":"<rfc3339 UTC>","verb":"<dispatch|advance|complete>","entity":"<slug>","workflow":"<$SLUG>","stage":"<stage entered>","text":"<one-line summary, ≤120 chars, no newlines or quotes>"}
```

- `verb`: `dispatch` (sent an ensign to a stage), `advance` (moved an entity to its next stage), `complete` (entity reached terminal).
- → **All hosts (host-neutral producer):** the `bridge-inbox` mod's feed step appends it. No per-host binding — it is FO-loop narration, not a host event.

## `_bridge/sessions/<session_id>.json` — session→entity marker (RUNNING-badge source)

Last-write (first-write-wins per session), one file per live working session. Maps a session id to the ship it is driving, so Bridge can render the deterministic live FO-vs-ensign RUNNING badge.

```
{"session_id":"<«session-id»>","entity":"<slug>","workflow":"<workflow dir name>"}
```

- The `entity`/`workflow` pair is derived from the ensign's first Read of its entity file under `docs/spacedock/<workflow>/<slug>.md` (flat or `<slug>/index.md`); the path carries both the workflow (so Bridge's join is collision-free across workflows reusing a ticket id) and the slug. First-write-wins records the ensign's own entity, read before any duplicate-check sibling read. Archived (`_archive/`) entities are never marked.
- → **Claude:** PRESENT — derived by the shared emitter from the Claude `PostToolUse`/`Read` hook payload (which fires on every tool call).
- → **Codex/Pi:** NOT CLAIMED — the shared emitter can write a marker from an explicit normalized `entity_path`, but the packaged Codex/Pi producers do not yet have durable live proof that they can supply the child actor and entity path reliably. Until that proof exists, Bridge must treat their event streams as activity-only and omit deterministic per-ship RUNNING markers.

## Decision: deterministic RUNNING badge is Claude-proven only (for now)

The session→entity marker is the **only** deterministic source for the live FO-vs-ensign RUNNING badge, and the live derivation is **Claude-proven only for now**. It relies on a hook that fires on *every* tool call (Claude Code `PostToolUse`) to record the marker on the ensign's first entity Read. The prose-driven alternative — having the ensign run a first-action shell to write its own marker — is unreliable: the ensign skips that step roughly three times out of four. Deriving it in the hook is what makes the badge dependable. Codex and Pi now have event producers, but their marker path stays unclaimed until their adapters prove durable child identity plus entity path evidence.

**On non-Claude hosts Bridge degrades gracefully, it does not break:**

- The heartbeat (`fo.$SLUG.json`) is host-neutral, so the FO still shows **attached**.
- Git `dispatch:`/`advance:` narration and `fo-feed.jsonl` (both host-neutral) still drive the **fleet-history** rail.
- Only the live per-ship **RUNNING badge** (FO-vs-ensign, this very moment) is absent unless the host has written a session marker — Bridge simply does not render it, rather than showing a wrong one.

When Codex or Pi can prove durable child identity plus entity path evidence, bind the session marker in that host's adapter (`## Bridge egress`) against this same schema; the badge lights up with no Bridge change.
