---
name: bridge-inbox
description: Drain captain intent queued by the Bridge command-center UI (_bridge/inbox.jsonl), routing per-workflow, and write the FO liveness heartbeat each tick
version: 0.2.0
fo-realm: "FO realm — the FO maintains this file directly; it is FO process (the seam to the Bridge UI), not product built under the dev workflow."
---

# Bridge Inbox

[Bridge](https://github.com/spacedock-dev/bridge) is a read-only command-center UI over this fleet. It cannot push into a running FO session (a Claude Code session has no inbound API), so it writes captain intent to a durable, **append-only** inbox at `_bridge/inbox.jsonl` (relative to the FO's working directory — the repo root where `spacedock claude` was launched). This mod drains that inbox on the FO's own loop ticks, acting only on intent **addressed to this workflow**, and writes a per-workflow liveness heartbeat Bridge reads.

**Path alignment (load-bearing).** Bridge anchors the inbox, heartbeat, and feed on its `--repo-root` flag, falling back to `--fleet` only when `--repo-root` is unset. So intent reaches this FO **only when the captain launches Bridge with `--repo-root` pointing at this FO's cwd** (the repo root). Under a multi-workflow layout (`--fleet <repo>/docs/spacedock` with no `--repo-root`), Bridge would write to `<repo>/docs/spacedock/_bridge/` while the FO reads `<repo>/_bridge/` — and intent silently never arrives. The dir is the same `_bridge/` Bridge resolves *iff* the two roots agree.

**This is a pull, not a push.** Delivery latency is one FO loop cadence: the captain's intent is read whenever the FO next idles or boots, never instantly. Bridge's UI says as much ("queued — read on the FO's next tick"); do not promise synchronous delivery.

**Your workflow slug.** Several FOs (one per commissioned workflow) can share one repo root and one `_bridge/` dir; everything below is scoped to THIS FO's workflow by its slug. Derive the slug once and validate it — it names the per-workflow cursor and heartbeat files, and an unsafe value must never escape `_bridge/`:

```
SLUG=$(basename "{dir}")
case "$SLUG" in
  ""|.|..)            echo "bridge-inbox: empty/relative slug — skipping" >&2; exit 0 ;;
  *[!A-Za-z0-9._-]*)  echo "bridge-inbox: unsafe slug '$SLUG' — skipping" >&2; exit 0 ;;
esac
```

**Inbox record schema** (one JSON object per line, written by Bridge):

```
{"ts": "<rfc3339>", "kind": "tell" | "conn" | "decision", "text": "<string>", "granted": <bool, conn only>, "target": "<workflow-slug>" | "all", "entity": "<slug, decision only>", "field": "<frontmatter field, decision only>", "value": "<value, decision only>"}
```

`target` routes the intent. Act on a record only when `target == "$SLUG"` **or** `target == "all"`; a **missing/empty `target` means `all`** (backward-compatible with older Bridge records, which carried no target). A record targeted at another workflow is skipped — but still counts as processed, so this workflow's cursor advances past it.

**Consume by a per-workflow cursor, never by rewrite.** Bridge only appends to the one shared `inbox.jsonl`; this mod advances only `_bridge/.inbox-cursor.$SLUG` (the count of inbox lines this workflow's FO has processed). Each workflow's FO owns its own cursor, so several FOs draining the same inbox never clobber each other, and re-firing with no new lines is a no-op.

**Run each hook's steps in one shell.** Every Bash invocation is a fresh shell, so `$SLUG` (and `$CURSOR`) only persist within a single invocation. Begin each hook below by deriving and validating `$SLUG` (the block above), then run that hook's heartbeat and drain steps in the **same** shell — do not split the heartbeat and drain into separate invocations that each expect `$SLUG` to already be set, or the second runs with an empty slug and writes a stray `_bridge/fo..json` Bridge never reads.

## Hook: startup

1. Derive and validate `$SLUG` (above), in the shell you run the rest of this tick in.
2. Write the heartbeat (see **Heartbeat** below) so Bridge shows this workflow attached as soon as the FO boots.
3. Drain any intent the captain queued while no FO was attached, so a freshly-booted FO picks up standing instructions before its first dispatch — run the **Drain** procedure below.

## Hook: idle

Derive and validate `$SLUG` (above) in this tick's shell, then refresh the heartbeat and drain — all in the same shell.

### Heartbeat

Bridge shows per-workflow FO liveness by reading `_bridge/fo.$SLUG.json`; it treats this workflow as live only when the `ts` is **fresh** (within 30 minutes) and **not in the future**. Stamp a present-time UTC timestamp on every tick so an attached FO keeps showing live:

```
mkdir -p _bridge
printf '{"session_id":"%s","ts":"%s","state":"idle"}\n' \
  "${SD_SESSION_ID:-${CLAUDE_CODE_SESSION_ID:-${CODEX_THREAD_ID:-}}}" "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > _bridge/fo.$SLUG.json
```

`SD_SESSION_ID` is the host-neutral session-id token owned by your runtime adapter's **Bridge egress** binding; the snippet reads it first and falls back to the host's own session var (`$CLAUDE_CODE_SESSION_ID` on Claude, `$CODEX_THREAD_ID` on Codex) so the heartbeat never silently blanks, with no per-tick `export` needed. An empty value (a host that exposes none) is still a valid liveness tick — Bridge reads freshness from `ts`. The full egress-surface contract — `events.jsonl`, this heartbeat, `fo-feed.jsonl`, the session→entity marker, and which host produces each — is `docs/dev/bridge-egress-contract.md`.

`state` is `idle`: this mod runs at startup/idle boundaries (you are between dispatches when it fires), so it cannot honestly claim `working` — the finer working/idle signal already lives in `_bridge/events.jsonl`. Writing the heartbeat is observe-only; never let it block the loop. (A session left idle with no captain interaction for over 30 minutes stops ticking and goes honestly not-attached in Bridge — that is intended, not a bug.)

### Drain

Drain newly-queued captain intent addressed to this workflow, if any:

1. If `_bridge/inbox.jsonl` does not exist, skip — no Bridge is attached.
2. Adopt the old shared cursor on first run (one-time migration: before this version the mod used a single shared `_bridge/.inbox-cursor`; seed from it so you do not re-drain — and re-apply — intent already processed, e.g. an old `conn` grant), then read this workflow's cursor:
   ```
   if [ ! -f _bridge/.inbox-cursor.$SLUG ] && [ -f _bridge/.inbox-cursor ]; then
     cp _bridge/.inbox-cursor _bridge/.inbox-cursor.$SLUG
   fi
   CURSOR=$(cat _bridge/.inbox-cursor.$SLUG 2>/dev/null || echo 0)
   ```
3. Snapshot the current line count and read exactly the new records (bounding the read so a concurrent Bridge append can't make the cursor skip a line):
   ```
   NEW=$(wc -l < _bridge/inbox.jsonl | tr -d ' ')
   sed -n "$((CURSOR + 1)),${NEW}p" _bridge/inbox.jsonl
   ```
   If `NEW` is not greater than `CURSOR`, there is nothing new — skip (idempotent).
4. For each new record, in order, parse `kind` / `text` / `granted` / `target` (and `entity` / `field` / `value` on a `decision` record). **Check the target first:** if `target` is present and is neither `"$SLUG"` nor `"all"`, this record is for another workflow's FO — skip it (it is not yours to act on); it still counts as processed (the cursor advances past it in step 5). Otherwise (target is `"$SLUG"`, `"all"`, or absent) act:
   - **`kind == "tell"`** — the captain sent you a message. Treat `text` as a directive or clarification for this tick: act on it as you would a captain instruction (commission or clear work, answer the implied question, adjust course), and acknowledge it to the captain.
   - **`kind == "conn"`** — a conn-handover change. `granted: true` → adopt the conn within the stated goal `text`: drive the entities the conn covers to done without stopping at their gates, per the conn rules in `first-officer-shared-core` (escalations remain non-delegable and still surface to the captain). `granted: false` → take the conn back: stop at every gate for the captain's call again.
   - **`kind == "decision"`** — the captain resolved a self-described decision gate from Bridge (Bridge cannot perform the gate's external side-effects — a Linear write, a label — so it queues the decision here instead of advancing the entity). Resolve `entity` (its slug) in THIS workflow, then treat `field`/`value` as the captain's gate verdict: set the field with `${SPACEDOCK_BIN:-spacedock} status --set --workflow-dir {dir} <entity> <field>=<value>`, then drive that entity through its current gate exactly as if the captain had decided it at the gate — your normal gate-resolution runs the workflow's own stage actions (including any external writes the stage prose defines) and advances it. Idempotent: if the entity is already resolved/terminal with that value, it is a no-op. If `entity` does not resolve in this workflow (it belongs to another member's slug), skip it like a mismatched target. Acknowledge to the captain which entity you resolved and how.
5. Advance this workflow's cursor to the snapshot you read: `echo "$NEW" > _bridge/.inbox-cursor.$SLUG`.
6. Report to the captain: how many intents you drained (for this workflow) and what you did with each.

If a record is malformed (not valid JSON, or an unknown `kind`), skip it but still advance the cursor past it, and note the skip to the captain — never block the loop on a bad record.

**Delivery is at-least-once, not exactly-once.** The cursor advances only *after* you act (step 5 follows step 4), so nothing is lost if the loop dies mid-drain. The trade-off: a crash between acting and writing the cursor re-delivers that batch on the next tick. Treat `conn`/`tell` handling as idempotent — re-adopting a `conn` you already hold (or re-relinquishing one you already gave back) is a no-op, and a repeated `tell` is at worst a duplicate acknowledgement. (This is distinct from the first-run migration seed above, which guards against re-applying the *entire* pre-versioning history.)

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
