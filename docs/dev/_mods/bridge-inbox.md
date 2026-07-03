---
name: bridge-inbox
description: Drain captain intent queued by the Bridge command-center UI (_bridge/inbox.jsonl), route per-workflow, acknowledge in _bridge/fo-replies.jsonl, and write the FO liveness heartbeat each tick
version: 0.3.0
fo-realm: "FO realm — the FO maintains this file directly; it is FO process (the seam to the Bridge UI), not product built under the dev workflow."
---

# Bridge Inbox

[Bridge](https://github.com/spacedock-dev/bridge) is a read-only command-center UI over this fleet. It cannot push into a running FO session (a Claude Code session has no inbound API), so it writes captain intent to a durable, **append-only** inbox at `_bridge/inbox.jsonl` (relative to the FO's working directory — the repo root where `spacedock claude` was launched). This mod drains that inbox on the FO's own loop ticks, acting only on intent **addressed to this workflow**, writes a per-workflow liveness heartbeat Bridge reads, and appends best-effort explanatory replies/acks to `_bridge/fo-replies.jsonl`.

**The mechanism is packaged — do not hand-write it.** Cursor math, per-line routing, the liveness heartbeat, and the JSONL reply/ack serialization all live in the binary, reached through four verbs:

- `spacedock bridge inbox drain --host «host» --slug «slug»` — stamps the heartbeat and returns the new records addressed to this workflow as JSON. It does **not** advance the cursor.
- `spacedock bridge inbox ack --host «host» --slug «slug» --line «n» --id «id» --ts «ts» --kind «intent-kind» --status «status» [--text … --granted … --entity … --field … --value … --verdict … --request-id …]` — appends one compact, newline-terminated reply/ack line (it derives the ack `kind` from the intent kind).
- `spacedock bridge inbox commit --slug «slug» --cursor «high_water»` — advances this workflow's cursor to the drained high-water mark (monotonic; a lower value is ignored).
- `spacedock bridge inbox check --host «host»` — the Claude `Stop`-hook decision helper (see **Wake**), not part of the manual loop.

`«host»` is your runtime adapter's host token (`claude` / `codex` / `pi`); `«slug»` is this workflow's slug (below). The commands resolve the session id from the adapter's binding themselves, so no `$SLUG`/cursor/JSON shell is authored by hand. **Your only job is judgment:** read the records the drain hands back, *act* on each (interpret a `tell`, resolve a gate), then ack and commit. Do not reimplement drain/ack/commit with `wc`/`sed`/`cat`/`echo`/`jq` — hand-rolled shell across separate tool calls is exactly what corrupted cursors and wrote non-compact JSONL in the past.

**Path alignment (load-bearing).** Bridge anchors the inbox, heartbeat, feed, and reply stream on its `--repo-root` flag, falling back to `--fleet` only when `--repo-root` is unset. So intent reaches this FO **only when the captain launches Bridge with `--repo-root` pointing at this FO's cwd** (the repo root). Under a multi-workflow layout (`--fleet <repo>/docs/spacedock` with no `--repo-root`), Bridge would write to `<repo>/docs/spacedock/_bridge/` while the FO reads `<repo>/_bridge/` — and intent silently never arrives. The dir is the same `_bridge/` Bridge resolves *iff* the two roots agree.

**This is a pull, not a push.** Delivery latency is one FO loop cadence: the captain's intent is read whenever the FO next idles or boots, never instantly. Bridge's UI says as much ("queued — read on the FO's next tick"); do not promise synchronous delivery.

**Wake (getting a parked FO to drain).** A busy FO drains eagerly each loop iteration and at idle, so intent reaches it within a tick. A **parked** FO — one that has stopped and is waiting at the prompt — is nudged differently per host. On Codex, Bridge resumes the session with `spacedock bridge ingress wake`. On Claude, external session resume is unsafe (the transcript has no write locking), so instead the packaged `Stop` hook runs `spacedock bridge inbox check` at every turn boundary: when intent is queued for this session's workflow it returns a `block` decision that keeps the FO going long enough to drain, entirely in-session. A fully idle/closed Claude session cannot be woken remotely; Bridge surfaces its queued count honestly and the captain nudges it. Never resume a live Claude session out-of-band to force a drain.

**Your workflow slug.** Several FOs (one per commissioned workflow) can share one repo root and one `_bridge/` dir; everything below is scoped to THIS FO's workflow by its slug. The slug is the basename of this workflow's dir (`basename {dir}`) — it names the per-workflow cursor and heartbeat files. Pass it verbatim as `--slug` to each `spacedock bridge inbox` verb; the command validates it and refuses an unsafe value (one that could escape `_bridge/`) rather than acting on it, so you never guard it by hand. Below, `$SLUG` refers to that slug.

**Inbox record schema** (one JSON object per line, written by Bridge):

```
{"id":"<opaque unique string>","ts":"<rfc3339>","kind":"tell"|"conn"|"decision"|"permission-decision","text":"<string>","granted":<bool, conn only>,"target":"<workflow-slug>"|"all","target_set":["<workflow-slug>", "..."],"entity":"<slug, decision only>","field":"<frontmatter field, self-described decision only>","value":"<value>","verdict":"approve|reject|redo, plain decision only","directives":["plain decision route-back notes"],"request_id":"<permission alert id>"}
```

`id` is Bridge's opaque unique id for this intent. Current Bridge records carry `target_set` when Bridge can freeze a recipient set: for `target == "all"` with known fleet members, Bridge writes the current fleet member slugs into `target_set`; for a specific target, Bridge writes `[slug]`. If Bridge cannot resolve any current member slugs for a broadcast, it omits `target_set` so the FO preserves legacy `target == "all"` routing. `target` remains for backward compatibility with older records and unknown-recipient broadcasts.

Routing is:

- If `target_set` is present, act only when `"$SLUG"` is in `target_set`. Ignore `target` entirely for routing in that case, including `target == "all"`; the frozen `target_set` is authoritative.
- If `target_set` is absent, preserve old target behavior: act when `target == "$SLUG"`, `target == "all"`, or `target` is missing/empty.
- A record not addressed to this workflow is skipped — but still counts as processed, so this workflow's cursor advances past it.

**Consume by a per-workflow cursor, never by rewrite (the command owns it).** Bridge only appends to the one shared `inbox.jsonl`; the drain/commit verbs advance only `_bridge/.inbox-cursor.$SLUG` (the count of inbox lines this workflow's FO has processed), performing the one-time migration from the pre-versioning shared `_bridge/.inbox-cursor` for you. Each workflow's FO owns its own cursor, so several FOs draining the same inbox never clobber each other, and a drain with no new lines is a no-op. Because the whole read/route/heartbeat step is one atomic command, there is no `$SLUG`/`$CURSOR` shell state to lose across tool calls.

## Hook: startup

1. Run the **Drain procedure** below with `--host «host» --slug $SLUG`. `drain` stamps the liveness heartbeat first (so Bridge shows this workflow attached as soon as the FO boots), then returns any intent the captain queued while no FO was attached, so a freshly-booted FO picks up standing instructions before its first dispatch.

## Hook: idle

Run the **Drain procedure** below with `--host «host» --slug $SLUG`. A single `drain` call both refreshes the heartbeat and returns new intent.

### Heartbeat

Bridge shows per-workflow FO liveness by reading `_bridge/fo.$SLUG.json`; it treats this workflow as live only when the `ts` is **fresh** (within 30 minutes) and **not in the future**. `spacedock bridge inbox drain` stamps a present-time UTC timestamp (and the `host` and `session_id` from your adapter's binding) on every call, so an attached FO that drains each tick keeps showing live — you never write `fo.$SLUG.json` by hand. An empty session id (a host that exposes none) is still a valid liveness tick — Bridge reads freshness from `ts`. The full egress-surface contract — `events.jsonl`, this heartbeat, `fo-feed.jsonl`, `fo-replies.jsonl`, the session→entity marker, and which host produces each — is `docs/dev/bridge-egress-contract.md`.

`state` is `idle`: the drain runs at startup/idle boundaries (you are between dispatches when it fires), so it cannot honestly claim `working` — the finer working/idle signal already lives in `_bridge/events.jsonl`. The heartbeat is observe-only; a drain that cannot write it degrades to no-op and never blocks the loop. (A session left idle with no captain interaction for over 30 minutes stops ticking and goes honestly not-attached in Bridge — that is intended, not a bug.)

### Drain procedure

Drain newly-queued captain intent addressed to this workflow, if any:

1. **Drain.** Run `spacedock bridge inbox drain --host «host» --slug $SLUG`. It stamps the heartbeat and returns JSON: `{"cursor":N,"high_water":M,"count":K,"records":[…]}`. The `records` are exactly the new lines addressed to this workflow (routing already applied), each carrying its physical `LINE` number, and already excluding any record you have acked before (idempotent replay-safe). If `count` is `0`, there is nothing to act on — you are done (the heartbeat is already refreshed). No `_bridge/inbox.jsonl` yet means no Bridge is attached; the command reports `no-inbox` and you skip.
2. **Act + ack, in order.** For each returned record, act on it (below), then acknowledge it with `spacedock bridge inbox ack --host «host» --slug $SLUG --line «line» --id «id» --ts «ts» --kind «kind» --status «status»` plus any relevant `--text`/`--granted`/`--entity`/`--field`/`--value`/`--verdict`/`--request-id`. The `ack` verb derives the reply `kind` from the intent `kind` and writes one compact line for you. Ack **after** you have interpreted, accepted, or applied the intent — not merely after shell-reading the line. The per-kind action and the `status` to use:
   - **`kind == "tell"`** — the captain sent you a message. Treat `text` as a directive or clarification for this tick: act on it as you would a captain instruction (commission or clear work, answer the implied question, adjust course), and append a `reply` record with `status:"answered"`.
   - **`kind == "conn"`** — a conn-handover change. `granted: true` → adopt the conn within the stated goal `text`: drive the entities the conn covers to done without stopping at their gates, per the conn rules in `first-officer-shared-core` (escalations remain non-delegable and still surface to the captain). Then append a `conn-ack` record with `status:"accepted"`. `granted: false` → take the conn back: stop at every gate for the captain's call again. Then append a `conn-ack` record with `status:"released"`.
   - **`kind == "decision"`** — the captain resolved a gate from Bridge. This is a captain decision, not FO self-approval, and Bridge has not advanced the entity. Resolve `entity` (its slug) in THIS workflow, verify the entity is still at a compatible current gate, and apply the workflow's normal gate-resolution flow:
     - Self-described decision shape: `field` plus `value`. Treat it as the captain's selected gate value: set the field with `${SPACEDOCK_BIN:-spacedock} status --set --workflow-dir {dir} <entity> <field>=<value>`, then continue the current gate exactly as if the captain had decided it in chat.
     - Plain decision shape: `verdict` plus optional `directives`. `approve` advances through the gate's own approve path and side effects. `reject` / `redo` route to the gate's `feedback-to` stage with the supplied directives. If the workflow has external actions (GitHub review, Linear update, labels, etc.), perform them before terminal state or final acknowledgement.
     - Append a `decision-ack` record only after the decision is applied, blocked, or rejected. Use `status:"applied"` when gate resolution finished or was already satisfied; `status:"blocked"` when the intent is valid but execution could not finish; `status:"rejected"` when the intent is invalid, stale, or unresolvable. If `entity` does not resolve in this workflow, acknowledge it as `rejected` rather than silently treating it like a mismatched target; this record was addressed to this workflow and Bridge needs to close the loop.
   - **`kind == "permission-decision"`** — the captain resolved a top-level FO permission alert. Match `request_id` to the open `_bridge/fo-alerts.jsonl` record you emitted. `value: "deny"` means do not retry the blocked action; append a `permission-ack` with `status:"denied"` and report the block remains. `value: "approve-once"` means retry the exact blocked host action once using the runtime's escalation path. `value: "approve-rule"` means retry with the alert's proposed `prefix_rule` when one was present; otherwise treat it as `approve-once`. Append `permission-ack` with `status:"accepted"` before retrying; append `status:"blocked"` if the retry could not be started. This Bridge approval is FO intent, not a bypass of a host-native security prompt; if the host still presents a native approval dialog, honor it normally.
3. **Commit.** After acting on and acking every returned record, advance the cursor once with `spacedock bridge inbox commit --slug $SLUG --cursor «high_water»` (the `high_water` from step 1). This moves the cursor past everything read this tick — including records for other workflows the drain filtered out — so they are never reconsidered. Commit is monotonic and is the last step, so nothing is lost if the loop dies mid-drain.
4. Report to the captain: how many intents you drained (for this workflow) and what you did with each.

If a record is malformed (not valid JSON, missing required fields for its `kind`, or an unknown `kind`), skip it but still advance the cursor past it, and note the skip to the captain — never block the loop on a bad record. Append a rejected ack only when the record is addressed to `"$SLUG"` and has enough valid metadata to produce a valid reply shape: an `id`, `ts`, and a known `kind` that determines `reply`/`conn-ack`/`decision-ack` plus `intent_kind`. Unknown-kind or unrouteable records cannot be represented by the reply schema; report the skip to the captain, advance the cursor, and do not invent an id or invalid ack kind.

### Replies / acks

`spacedock bridge inbox ack` appends one reply/ack to `_bridge/fo-replies.jsonl` for each addressed inbox record you handled or rejected. You never format this JSON by hand — you pass flat flags and the command serializes the record below:

```
{"schema":1,"ts":"<rfc3339 UTC>","kind":"reply"|"conn-ack"|"decision-ack"|"permission-ack","target":"<actual $SLUG>","in_reply_to_id":"<Bridge intent id>","in_reply_to_line":123,"in_reply_to_ts":"<original intent ts>","intent_kind":"tell"|"conn"|"decision"|"permission-decision","status":"answered"|"accepted"|"released"|"applied"|"denied"|"rejected"|"blocked","text":"optional one-line note","granted":true|false,"entity":"...","field":"...","value":"...","verdict":"...","request_id":"...","session_id":"optional","host":"optional"}
```

What the command guarantees (so you don't have to):

- `target` is the actual acknowledging workflow slug (`$SLUG`), never `"all"` — pass `--slug $SLUG`.
- `in_reply_to_line` is the physical inbox line number you processed (`--line`), not the count of addressed records.
- `in_reply_to_id` and `in_reply_to_ts` echo the original inbox `id` and `ts` (`--id`, `--ts`).
- `intent_kind` echoes the inbox `kind` (`--kind`); the command maps it to the reply `kind`.
- `text` is optional (`--text`) and is flattened to a single line for you.
- Pass `--granted`, `--entity`, `--field`, `--value`, `--verdict`, and `--request-id` when present and relevant; omitted flags are omitted from the record.
- The command resolves `session_id` from the adapter's binding and stamps `host`; you do not supply them.
- It writes one complete newline-terminated JSON object in one append operation. Do not rewrite, truncate, sort, or compact `fo-replies.jsonl` yourself.
- Cursor remains the delivery/read source of truth; `fo-replies.jsonl` is best-effort explanatory ack content. A failed ack must never block the FO from completing the drained intent or committing the cursor after action.
- Duplicate replies are allowed under at-least-once replay; Bridge folds them by intent id (or legacy line fallback), acknowledging target, and reply kind. Still prefer idempotent behavior.

**Delivery is at-least-once, not exactly-once.** The cursor advances only *after* you act (commit is the last step, after act + ack), so nothing is lost if the loop dies mid-drain. The trade-off: a crash between acking and committing re-surfaces at worst nothing (the drain filters already-acked records) and between acting and acking re-delivers that record on the next tick. Treat `conn`/`tell` handling as idempotent — re-adopting a `conn` you already hold (or re-relinquishing one you already gave back) is a no-op, and a repeated `tell` is at worst a duplicate acknowledgement. (This is distinct from the first-run migration seed the command applies, which guards against re-applying the *entire* pre-versioning history.)

## Feed

Bridge's fleet-history rail shows the FO's narration. For a workflow whose entities are committed it can read the `dispatch:`/`advance:` git narration — but a local-only workflow (entities gitignored; no such commits) leaves that history empty even while you drive. So append a narration line to `_bridge/fo-feed.jsonl` (relative to the repo root, the same `_bridge/` the heartbeat and inbox use) each time you **dispatch**, **advance**, or **complete** an entity:

```
mkdir -p _bridge
printf '{"ts":"%s","verb":"%s","entity":"%s","workflow":"%s","stage":"%s","text":"%s"}\n' \
  "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "dispatch" "{slug}" "$SLUG" "{stage}" "{short note}" \
  >> _bridge/fo-feed.jsonl
```

- `verb` is `dispatch` (you sent an ensign to a stage), `advance` (you moved an entity to its next stage), or `complete` (an entity reached terminal).
- `entity` is the entity slug; `workflow` is `$SLUG` (this member's workflow); `stage` is the stage entered.
- `text` is a one-line human summary (≤120 chars, no newlines or `"`). Keep it factual — the captain reads this stream to follow the drive.
- Append-only and best-effort, exactly like the event stream: never let it block or fail the loop, and never rewrite the file (Bridge tails it; a concurrent append is fine). It is gitignored session runtime, like `events.jsonl`.
