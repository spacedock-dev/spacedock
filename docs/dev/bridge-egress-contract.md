# Bridge egress contract (FO liveness/replies → Bridge)

The Bridge seam is two flows over the shared `_bridge/` dir (resolved from the repo root where the FO launched). They run in opposite directions and are documented apart:

- **Ingress** — captain intent → FO — is the `bridge-inbox` mod (`docs/dev/_mods/bridge-inbox.md`). It drains `_bridge/inbox.jsonl` on the FO's own loop ticks; it rides the portable Spacedock mod-hook loop and is already host-neutral. Bridge writes every new intent with an opaque `id` and writes a frozen `target_set` array when it can resolve the workflow slugs expected to drain and acknowledge the intent. For `target == "all"` with known fleet members, Bridge expands the current fleet member slugs into `target_set`; for a specific target, Bridge writes `[slug]`. If Bridge cannot resolve any current member slugs for a broadcast, it omits `target_set` so the FO preserves legacy `target == "all"` routing. When `target_set` is present it is authoritative, and the legacy `target` field is ignored for routing; `target` remains only for backward compatibility with records that predate `target_set` and unknown-recipient broadcasts.
- **Egress** — FO liveness/activity/replies → Bridge — is *this* contract: the `_bridge/` files a running FO (and its ensigns) write for the external Bridge command-center UI to tail. Bridge reads what already exists; it is not a source of truth.

**The schema is Spacedock-owned and harness-neutral. The *producer* of the host-event files (`events.jsonl` + the session marker) is per-host** — the heartbeat, feed, and FO reply stream are written host-neutrally by the `bridge-inbox` mod. The per-host producers are bound the same way runtime lifecycle capabilities are (see `fo-dispatch-core.md`'s "Claude: PRESENT / Codex: ABSENT / Pi: …" idiom): the shared contract names the file shape; the host adapter (`skills/first-officer/references/{claude,codex,pi}-first-officer-runtime.md`, under their `## Bridge egress` section) owns the concrete mechanism that emits it. The current implementation uses the shared `spacedock bridge egress emit --host <host>` command for host normalization. Claude, Codex, and Pi have packaged event producers; deterministic session→entity marker parity is proven only for Claude so far.

All files are gitignored session runtime (`.gitignore` carries `_bridge/`), append-only or last-write where noted, and strictly observe-only: a telemetry side-channel must never block, fail, or alter the FO.

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
- → **Codex:** PACKAGED/FIXTURE-COVERED — `.codex-plugin/plugin.json` points at `hooks/codex-hooks.json`, whose non-async command hooks call `spacedock bridge egress emit --host codex` directly via `SPACEDOCK_BIN` or `PATH`. The hook does not depend on plugin-root environment variables, so installed cache layout and dev checkout layout are both valid. This proves packaging and minimal lifecycle payload handling, not live marker parity.
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

## `_bridge/fo-replies.jsonl` — captain-intent acknowledgements

One JSON object per line, appended after the FO has interpreted, accepted, or applied an inbox intent. This is the explanatory reply/ack stream for Bridge's conversation loop; the per-workflow inbox cursor remains the delivery/read source of truth. Duplicate replies are allowed under at-least-once replay; Bridge folds them by intent id (or legacy line fallback), acknowledging target, and reply kind. The stream is best-effort explanatory content, not an exactly-once delivery ledger.

```
{"schema":1,"ts":"<rfc3339 UTC>","kind":"<reply|conn-ack|decision-ack|permission-ack>","target":"<acknowledging workflow slug>","in_reply_to_id":"<Bridge intent id>","in_reply_to_line":123,"in_reply_to_ts":"<original intent ts>","intent_kind":"<tell|conn|decision|permission-decision>","status":"<answered|accepted|released|applied|denied|rejected|blocked>","text":"optional one-line note","granted":true,"entity":"...","field":"...","value":"...","verdict":"...","request_id":"...","session_id":"optional","host":"optional"}
```

- `target` is the actual acknowledging workflow slug, never `all`.
- `in_reply_to_line` is the physical line number in `_bridge/inbox.jsonl` that produced the ack.
- `kind`: `reply` for `tell`, `conn-ack` for `conn`, `decision-ack` for `decision`, `permission-ack` for `permission-decision`.
- `status`: `answered` for a handled `tell`; `accepted` when the FO adopts a conn grant or accepts a permission retry; `released` when the FO gives the conn back; `applied` when a decision is present and gate resolution finished or was already satisfied; `denied` when the captain rejects a permission request; `blocked` when a valid intent could not finish; `rejected` when an intent is invalid or unresolvable.
- Echo `granted`, `entity`, `field`, `value`, `verdict`, and `request_id` when present and relevant to the intent. Keep `text` one line.
- Append-only and best-effort: write one complete newline-terminated JSON object in one append operation, and never rewrite `fo-replies.jsonl`.
- → **All hosts (host-neutral producer):** the `bridge-inbox` mod appends replies/acks while draining `_bridge/inbox.jsonl`.

## `_bridge/fo-alerts.jsonl` — top-level FO alerts

One JSON object per line, appended when the FO is blocked by a captain-owned host decision that should surface above ordinary chat replies. The first producer is the sandbox permission path: when a command cannot proceed because the host needs approval, the FO writes an open permission request before waiting for a decision.

```
{"schema":1,"id":"perm_<opaque>","ts":"<rfc3339 UTC>","kind":"permission-request","severity":"blocked","workflow":"<workflow slug>","entity":"<entity slug>","host":"<claude|codex|pi>","session_id":"<«session-id»>","reason":"<one-line reason>","command":"<command summary>","prefix_rule":["git","-C"],"status":"open"}
```

- `id` is the stable join key for the captain's response.
- `workflow` and `entity` scope the alert for fleet UI routing; `entity` may be empty when the block happens before an entity is selected.
- `command` is a concise command summary, not a secret-bearing shell transcript.
- `prefix_rule` is optional and present only when the FO has a narrowly scoped reusable approval to propose.
- Bridge overlays a later inbox `permission-decision` record for the same `request_id` to show approved/denied state; the alert file itself remains append-only.
- → **All hosts (host-neutral producer):** the FO writes this through `spacedock bridge alert permission ...`. The helper returns `{"id":"...","request_id":"...","queued":true}` and appends one line to `_bridge/fo-alerts.jsonl`. If the helper cannot queue the alert, it still exits without blocking the FO command and returns `{"queued":false,"error":"..."}`.

## `_bridge/fo-initiate.jsonl` — FO-authored feed lines and decidable gates

One JSON object per line, appended when the FO originates a message to the captain that is not a reply to a captain intent: a status note, a recommendation, or a decidable `gate-review`. Bridge reads it, folds records by `id` (keeping the latest per id), and renders each as an FO-authored feed line — a `gate-review` gets an inline Approve/Reject affordance.

```
{"schema":1,"id":"<stable fold key>","ts":"<rfc3339 UTC>","kind":"<status|reco|gate-review>","workflow":"<slug>","entity":"<slug>","ship_id":"<slug>/<entity>","host":"<claude|codex|pi>","session_id":"<«session-id»>","headline":"<one-line lede>","body":"optional supporting prose","request_id":"<gate loop-closure correlator>","status":"open"}
```

- `id` is REQUIRED and is the fold key. Unlike `fo-alerts.jsonl` there is NO random fallback: idempotency depends on a stable caller-supplied id. The FO re-emits the same record each drain tick; Bridge collapses them to one card by `id`. For `gate-review`, `present-gate` derives `id`/`request_id` deterministically from `(entity, stage)` so re-emit is byte-stable.
- `kind` ∈ `status | reco | gate-review`. `status`/`reco` are plain FO-authored lines; `gate-review` is a decidable card.
- `headline` is REQUIRED, bounded to 240 chars; `body` is optional, bounded to 2000 chars. Both are collapsed to one line (control chars and whitespace runs become single spaces).
- `request_id` is the gate loop-closure correlator; it defaults to `id` for `gate-review` and is omitted for `status`/`reco`. Approve/Reject write a decision intent to `inbox.jsonl` carrying this `request_id`; Bridge's fo-initiate reader overlays that decision to flip the card's status.
- `status`: the writer ALWAYS writes `open`. The READER overlays `resolved`/`approved`/`rejected` from decision intents — never trust a written non-open status.
- **Path anchoring:** the writer resolves `filepath.Abs(--repo-root or cwd)/_bridge` — EXACTLY like `bridgealert.AppendPermission`, NOT `bridgeegress.canonicalBridgeRoot`. It MUST be passed the same repo root Bridge resolves from, or the write lands in a divergent `_bridge/`.
- **Bounded, but open gates never evicted:** the file is capped (the writer truncates to a recent-tail window on each append), but it NEVER drops the latest record of a still-open `gate-review` id, so an open gate cannot scroll out of the read window.
- **Channel boundary:** a decidable gate lives in `fo-initiate.jsonl` ONLY — never `fo-feed.jsonl` (ambient git narration) and never `fo-replies.jsonl` (which requires an `in_reply_to` correlator and would silently drop an uncorrelated FO-originated push).
- → **All hosts (host-neutral producer):** the FO writes this through `spacedock bridge initiate --kind <kind> ...`. The helper returns `{"id":"...","request_id":"...","queued":true}` and appends one line. If it cannot queue, it exits without blocking the FO and returns `{"queued":false,"error":"..."}`. Gate emission is wired in the `present-gate` skill (loaded by every host), not in a per-host hook.

## `_bridge/sessions/<actor_id>.json` — session→entity marker (RUNNING-badge source)

First-write-wins per host actor, one file per live working actor. The filename is the normalized `actor_id`: Claude main/ensign markers currently use the session id, while hosts that provide child ids can use a host-scoped composite such as `session_id.agent_id`. The marker maps that actor to the ship it is driving, so Bridge can render the deterministic live FO-vs-ensign RUNNING badge.

```
{"host":"<claude|codex|pi>","session_id":"<«session-id»>","agent_id":"<child id, when present>","actor_id":"<host-scoped actor id>","entity":"<slug>","workflow":"<workflow dir name>"}
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
