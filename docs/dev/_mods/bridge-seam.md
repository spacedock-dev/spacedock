---
name: bridge-seam
description: Produce the Bridge `_bridge/` seam by direct file writes — a liveness heartbeat carrying the harness session id, captain-intent drain by cursor, terminal acks, and FO-authored cards/alerts — per Bridge's `docs/seam-contract.md`
version: 0.1.0
---

# Bridge Seam

[Bridge](https://github.com/spacedock-dev/bridge) is a read-only command-center UI over this fleet. It cannot push into a running FO session (a Claude Code / Codex / Pi session has no inbound API), so the seam between Bridge and the FO is a set of **plain files under `_bridge/`** (relative to the FO's working directory — the fleet root where the FO was launched). Bridge writes exactly one of them, the captain-intent queue `_bridge/inbox.jsonl`; every other file the FO **writes** and Bridge **reads**. This mod is the producer side of that contract: on boot and on every loop tick the FO writes a liveness heartbeat, drains the intent queue by a monotonic per-slug cursor, and appends terminal acks — all as direct file writes, no Spacedock CLI verb.

**Direct file writes, not a verb.** The mechanism is a file recipe a bare agent can follow: read a cursor integer, read newline-delimited JSONL, append single JSONL lines, replace a small JSON file. There are no packaged drain/ack/commit verbs — those are retired; the FO performs the file operations itself. The one binary touchpoint is the hook-driven egress producer (`spacedock bridge egress emit`, wired by the plugin hooks) and, on Claude, the synchronous Stop check (`spacedock bridge inbox check`) that keeps a parked FO alive long enough to drain. Neither is invoked by hand from this mod.

**The full contract is authoritative in Bridge's `docs/seam-contract.md`** (§2 per-file shapes, §3 the drain recipe); a Spacedock-local overview lives at `docs/dev/bridge-seam.md`. Every JSON shape below is quoted from that contract. Where a reader tolerates more than shown, the reader's tolerance is the contract.

**This is a pull, not a push.** Delivery latency is one FO loop cadence: queued captain intent is read whenever the FO next boots or idles, never instantly. A **parked** FO (stopped, waiting at the prompt) is nudged per host — on Codex, Bridge resumes the session with `spacedock bridge ingress wake`; on Claude, the packaged `Stop` hook runs `spacedock bridge inbox check` at every turn boundary and returns a `block` decision while intent is queued, so the FO drains in-session. Never resume a live Claude session out-of-band; its transcript has no write locking.

**Per-host session-id source (load-bearing).** The heartbeat MUST carry `session_id` set to the harness's own session id, because the Claude Stop check resolves *which* workflow slugs belong to the stopping session by matching the Stop payload's `session_id` against the `session_id` in each `_bridge/fo.<slug>.json` heartbeat. A heartbeat with a missing or wrong session id silently kills Claude intent delivery: the check resolves no slug, never blocks, and the queued intent sits forever. Read the id from the per-host variable:

| host | session-id source |
|---|---|
| claude | `$CLAUDE_CODE_SESSION_ID` |
| codex | `$CODEX_THREAD_ID` |
| pi | the pi runtime's session id |

**Mod placement (known limitation).** Lifecycle hooks run from `{workflow_dir}/_mods/`, so a workflow gets this seam only when `bridge-seam.md` is present in its `_mods/` dir. The dogfood copy at `docs/dev/_mods/bridge-seam.md` covers this repo's own workflows; a non-dogfood workflow must copy the mod into its `_mods/`. Automatic scaffolding at commission/refit time is a named follow-up, not solved here.

## Hook: startup

For each workflow slug `$SLUG` this FO owns (one per commissioned workflow — see **Your workflow slug(s)**), **before the greet**:

1. Write the liveness heartbeat `_bridge/fo.$SLUG.json` carrying the harness session id (see **Heartbeat**), with `state:"working"`. This makes Bridge show the workflow attached the moment the FO boots — even a greet-and-stop launch.
2. Run the **Drain procedure** so a freshly-booted FO picks up any intent the captain queued while no FO was attached, before its first dispatch.

## Hook: idle

At the top of each loop tick:

1. Run the **Drain procedure** for each `$SLUG` — it acts on any new intent and refreshes the heartbeat (`ts = now`).
2. When parking (turn end / awaiting the captain), fulfil the **On park** obligations.

## Agent Prompt

You are the producer side of the Bridge `_bridge/` seam. Your job is judgment plus a handful of exact file writes. Do not reimplement this with a retired CLI verb; do the file operations directly, using the shapes below verbatim.

### Your workflow slug(s)

`$SLUG` is this workflow's slug — the basename of its directory (`basename {dir}`). It names the per-slug cursor and heartbeat files. It must be a **safe slug**: a single path element, no `/`, no `.`/`..`. Several FOs (one per commissioned workflow) can share one fleet root and one `_bridge/` dir; every step below is scoped to a single `$SLUG`. In a **fleet** (this FO owns several members), perform every per-slug step — heartbeat, cursor, drain, acks — **once per member slug**. Bridge shows a member live only when *its own* `fo.<slug>.json` is fresh; a single shared heartbeat shows only one member attached.

### Heartbeat

Bridge reads per-workflow FO liveness from `_bridge/fo.$SLUG.json` (whole-file replace). It treats the workflow as live only when `ts` is **fresh** (within 30 minutes) and **not future-dated**. Write it on boot and refresh it every tick and every drain. Exact shape (contract §2.4):

```json
{"session_id":"sess_9f","ts":"2026-07-14T18:06:30Z","state":"idle","host":"claude"}
```

- `session_id` — the harness session id from the per-host source table above (`$CLAUDE_CODE_SESSION_ID` / `$CODEX_THREAD_ID` / pi's id). **This field is load-bearing for Claude wake delivery** — never omit it or write a placeholder.
- `ts` — present-time RFC3339 UTC. A zero/absent or future ts ⇒ Bridge reports not-attached (never fabricates `working`).
- `state` — `working` while acting; `idle` when parked awaiting the captain.
- `host` — `claude` \| `codex` \| `pi` (the harness stamping it).

### Drain procedure

Per member slug `$SLUG`, on boot and every loop tick:

1. **Read the cursor.** `C = int(contents of "_bridge/.inbox-cursor.$SLUG")`, or `0` if the file is absent, empty, non-integer, or negative. The cursor is the count of physical inbox lines this slug's FO has already drained (its high-water line number).
2. **Read the inbox by physical line number.** Read `_bridge/inbox.jsonl` line by line, counting **1-based physical newline-terminated lines** — every newline-terminated line consumes a number **even a malformed one** (skip it from acting but still count it); a trailing fragment with **no** newline at EOF is a torn write, skipped and **not** counted. This is the `wc -l` rule Bridge uses for read-back, so your cursor agrees with Bridge's line numbers. Each inbox line looks like (contract §2.1):
   ```json
   {"id":"bi_ab12","ts":"2026-07-14T18:03:00Z","kind":"tell","text":"ping","target":"all"}
   ```
   Fields: `id` (stable, the primary reply correlator), `ts`, `kind` (`tell`\|`conn`\|`decision`\|`permission-decision`), optional `text`, `granted` (`conn`), `target`, `target_set`, `entity`/`field`/`value`/`verdict`/`directives` (`decision`), `request_id`.
3. **Route.** For each line with physical number `N > C`, decide whether it is addressed to `$SLUG`:
   - If `target_set` is **present**, act **only** when `$SLUG` is in `target_set`. The frozen `target_set` is authoritative — ignore `target` entirely, including `target:"all"`.
   - If `target_set` is **absent**, act when `target == "$SLUG"`, `target == "all"`, or `target` is missing/empty.
   - A line not addressed to `$SLUG` is skipped, but still counts as processed so the cursor advances past it.
4. **Dedup, then act.** Before acting on an addressed intent with `id` X, scan `_bridge/fo-replies.jsonl` for a **terminal** ack whose `in_reply_to_id` is X (terminal = any status except the interim `acting`). If one exists, you already handled X — skip it (still counted). Otherwise act on it by kind:
   - `tell` — treat `text` as a captain directive for this tick; act on it, then append a `reply` ack with `status:"answered"`.
   - `conn` — a conn-handover change. `granted:true` → adopt the conn within the stated goal `text` (drive the covered entities to done without stopping at their gates; escalations stay non-delegable), then append a `conn-ack` with `status:"accepted"`. `granted:false` → take the conn back (stop at every gate again), then `conn-ack` `status:"released"`.
   - `decision` — the captain resolved a gate from Bridge (a captain decision, not FO self-approval; Bridge has not advanced the entity). Resolve `entity` in this workflow and apply the normal gate flow: self-described shape (`field`+`value`) → set the field and continue the gate as if decided in chat; plain shape (`verdict`+optional `directives`) → `approve` advances, `reject`/`redo` route to the gate's `feedback-to` stage with the directives; perform any external actions (GitHub, Linear) before terminal state. Then append a `decision-ack` with `status:"applied"` (finished/already-satisfied), `status:"blocked"` (valid but could not finish), or `status:"rejected"` (invalid/stale/unresolvable — including when `entity` does not resolve here). You may append an interim `acting` ack first.
   - `permission-decision` — the captain resolved a top-level FO permission alert. Match `request_id` to the open `_bridge/fo-alerts.jsonl` record you emitted. `value:"deny"` → do not retry; append `permission-ack` `status:"denied"`. `value:"approve-once"` → retry the exact blocked action once via the runtime's escalation path. `value:"approve-rule"` → retry with the alert's `prefix_rule` if present, else treat as `approve-once`. Append `permission-ack` `status:"accepted"` before retrying, or `status:"blocked"` if the retry could not start. A Bridge approval is FO intent, not a bypass of a host-native security prompt; honor any native dialog.
5. **Recount and advance the cursor.** After acting on every addressed line this tick, set `L` = the current **total physical line count** of `_bridge/inbox.jsonl` (`wc -l`, same physical-line rule as step 2). Write `L` to `_bridge/.inbox-cursor.$SLUG` as a single decimal integer (whole-file replace). This advances past everything read this tick, including other slugs' lines the routing filtered out.
6. **Refresh the heartbeat** (`ts = now`, `state:"working"` while still acting).
7. Report to the captain how many intents you drained for this workflow and what you did with each.

Missing/empty `_bridge/inbox.jsonl` ⇒ no Bridge attached; skip (write nothing but the heartbeat). A malformed line ⇒ skip acting on it but still count it toward the cursor; note the skip to the captain.

### Cursor safety

The cursor is **monotonic — never lower it.** Only ever raise `.inbox-cursor.$SLUG` to the current physical line count. If you are ever unsure of the count, re-read the inbox and recount from scratch before writing; do not guess a higher number. **An over-advanced cursor silently skips captain intent and nothing re-blocks it**: the Claude Stop check only ever sees intent *below* the cursor as still-pending, so any line you jumped over is invisible to the check and will never trigger another block. Under-counting is self-healing (the next tick re-drains); over-counting is permanent silent loss. When in doubt, count low.

### Replies / acks (`_bridge/fo-replies.jsonl`)

Append one JSONL line per addressed intent you handled or rejected (append-only). Bridge **drops** any line whose `schema` is not exactly `1`, or that lacks a correlator. Exact shape (contract §2.3):

```json
{"schema":1,"ts":"2026-07-14T18:05:00Z","kind":"reply","target":"my-wf","in_reply_to_id":"bi_ab12","intent_kind":"tell","status":"answered","text":"done"}
```

Required on every ack: `schema:1`; non-zero RFC3339 `ts`; `kind`; `target` = your `$SLUG` (safe slug); `in_reply_to_id` = the intent's `id` (the strong correlator — a record must carry either `in_reply_to_id`, or both `in_reply_to_line>0` and a non-zero `in_reply_to_ts`); `intent_kind` = the correlated intent's `kind`; a `status` valid for the `kind`.

`kind` is derived from the intent kind, and `status` must be valid for that `kind`:

| intent `kind` | reply `kind` | terminal statuses (done / failed) | interim |
|---|---|---|---|
| `tell` | `reply` | `answered` / `rejected`,`blocked` | `acting` |
| `conn` | `conn-ack` | `accepted`,`released` / `rejected`,`blocked` | `acting` |
| `decision` | `decision-ack` | `applied` / `rejected`,`blocked` | `acting` |
| `permission-decision` | `permission-ack` | `accepted` / `denied`,`rejected`,`blocked` | `acting` |

`acting` is the interim ack (received → acting → terminal) and is legal for all four kinds; a terminal ack never regresses to `acting`. Optional echo fields: `text`, `granted` (`conn-ack`), `entity`/`field`/`value` (`decision-ack`), `request_id`, `verdict`, `session_id`, `host`. A well-formed ack that correlates to no loaded intent is surfaced by Bridge (muted), never a hard error.

### Cards: status / reco / gate-review (`_bridge/fo-initiate.jsonl`)

FO judgment the captain sees. Append-only; Bridge **drops** a line unless `schema` is `1`, `id` is non-empty, and `kind` is known. Exact shape (contract §2.7):

```json
{"schema":1,"id":"init_7a","ts":"2026-07-14T18:07:00Z","kind":"gate-review","workflow":"linear-drc-ship","entity":"drc-3467","headline":"Ready to ship?","body":"...","status":"open"}
```

`kind` is `status` (ambient), `reco` (recommendation), or `gate-review` (decidable). Always write `status:"open"` — Bridge overlays resolution itself from the matching inbox `decision` (correlated by `request_id`, which defaults to `id`). Re-emitting the same `id` folds to one card (latest `ts` wins). An open `gate-review` is never evicted by the card cap, however old.

### Permission alerts (`_bridge/fo-alerts.jsonl`)

A high-priority FO→captain interrupt when a command is blocked. Append-only; a line is dropped unless `id` is non-empty and `kind == "permission-request"`. Exact shape (contract §2.8):

```json
{"schema":1,"id":"al_3c","ts":"2026-07-14T18:08:00Z","kind":"permission-request","workflow":"linear-drc-ship","reason":"rm outside repo","command":"rm -rf /tmp/x","prefix_rule":["rm -rf /tmp/"],"status":"open"}
```

`id` is the correlator for the captain's `permission-decision`. `status` empty defaults to `open`; Bridge overlays the decision. (The deferred permission-alert helper `fo-bridge.md` owns the exact emit prose.)

### On park

Do not silently park — a parked turn with no card is indistinguishable from finished work.

1. If you reached a gate, append a `gate-review` to `_bridge/fo-initiate.jsonl` with `status:"open"` and a `request_id` **before** parking.
2. If a command was blocked, append a `permission-request` to `_bridge/fo-alerts.jsonl`.
3. Refresh the heartbeat with `state:"idle"`.

### Lifecycle egress (harness-driven, not authored here)

The live working/idle badge and per-ship running signal come from `_bridge/events.jsonl` (turn-lifecycle lines) and `_bridge/sessions/<id>.json` (session→entity markers). Those are produced by the packaged plugin hooks (`spacedock bridge egress emit`) inside the harness turn lifecycle, not written by hand from this mod. Without them Bridge still renders durable state from git narration + cursors + heartbeats. `_bridge/fo-feed.jsonl` is **optional** enrichment only (contract §5): every signal it carries is already covered by git narration and the marker-derived feed, so this mod does not produce it.
