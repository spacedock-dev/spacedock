---
name: bridge-inbox
description: Drain captain intent queued by the Bridge command-center UI (_bridge/inbox.jsonl), routing per-workflow, and write the FO liveness heartbeat each tick
version: 0.2.0
fo-realm: "FO realm — the FO maintains this file directly; it is FO process (the seam to the Bridge UI), not product built under the dev workflow."
---

# Bridge Inbox

[Bridge](https://github.com/spacedock-dev/bridge) is a read-only command-center UI over this fleet. It cannot push into a running FO session (a Claude Code session has no inbound API), so it writes captain intent to a durable, **append-only** inbox at `_bridge/inbox.jsonl` (relative to the FO's working directory — the repo root where `spacedock claude` was launched, the same root Bridge resolves). This mod drains that inbox on the FO's own loop ticks, acting only on intent **addressed to this workflow**, and writes a per-workflow liveness heartbeat Bridge reads.

**This is a pull, not a push.** Delivery latency is one FO loop cadence: the captain's intent is read whenever the FO next idles or boots, never instantly. Bridge's UI says as much ("queued — read on the FO's next tick"); do not promise synchronous delivery.

**Your workflow slug.** Several FOs (one per commissioned workflow) can share one repo root and one `_bridge/` dir; everything below is scoped to THIS FO's workflow by its slug. Derive the slug once and validate it — it names the per-workflow cursor and heartbeat files, and an unsafe value must never escape `_bridge/`:

```
SLUG=$(basename "{workflow_dir}")
case "$SLUG" in
  ""|.|..)            echo "bridge-inbox: empty/relative slug — skipping" >&2; exit 0 ;;
  *[!A-Za-z0-9._-]*)  echo "bridge-inbox: unsafe slug '$SLUG' — skipping" >&2; exit 0 ;;
esac
```

**Inbox record schema** (one JSON object per line, written by Bridge):

```
{"ts": "<rfc3339>", "kind": "tell" | "conn", "text": "<string>", "granted": <bool, conn only>, "target": "<workflow-slug>" | "all"}
```

`target` routes the intent. Act on a record only when `target == "$SLUG"` **or** `target == "all"`; a **missing/empty `target` means `all`** (backward-compatible with older Bridge records, which carried no target). A record targeted at another workflow is skipped — but still counts as processed, so this workflow's cursor advances past it.

**Consume by a per-workflow cursor, never by rewrite.** Bridge only appends to the one shared `inbox.jsonl`; this mod advances only `_bridge/.inbox-cursor.$SLUG` (the count of inbox lines this workflow's FO has processed). Each workflow's FO owns its own cursor, so several FOs draining the same inbox never clobber each other, and re-firing with no new lines is a no-op.

## Hook: startup

1. Write the heartbeat (see **Heartbeat** below) so Bridge shows this workflow attached as soon as the FO boots.
2. Drain any intent the captain queued while no FO was attached, so a freshly-booted FO picks up standing instructions before its first dispatch — run the **Drain** procedure below.

## Hook: idle

Refresh the heartbeat, then drain.

### Heartbeat

Bridge shows per-workflow FO liveness by reading `_bridge/fo.$SLUG.json`; it treats this workflow as live only when the `ts` is **fresh** (within 30 minutes) and **not in the future**. Stamp a present-time UTC timestamp on every tick so an attached FO keeps showing live:

```
mkdir -p _bridge
printf '{"session_id":"%s","ts":"%s","state":"idle"}\n' \
  "${CLAUDE_CODE_SESSION_ID:-}" "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > _bridge/fo.$SLUG.json
```

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
4. For each new record, in order, parse `kind` / `text` / `granted` / `target`. **Check the target first:** if `target` is present and is neither `"$SLUG"` nor `"all"`, this record is for another workflow's FO — skip it (it is not yours to act on); it still counts as processed (the cursor advances past it in step 5). Otherwise (target is `"$SLUG"`, `"all"`, or absent) act:
   - **`kind == "tell"`** — the captain sent you a message. Treat `text` as a directive or clarification for this tick: act on it as you would a captain instruction (commission or clear work, answer the implied question, adjust course), and acknowledge it to the captain.
   - **`kind == "conn"`** — a conn-handover change. `granted: true` → adopt the conn within the stated goal `text`: drive the entities the conn covers to done without stopping at their gates, per the conn rules in `first-officer-shared-core` (escalations remain non-delegable and still surface to the captain). `granted: false` → take the conn back: stop at every gate for the captain's call again.
5. Advance this workflow's cursor to the snapshot you read: `echo "$NEW" > _bridge/.inbox-cursor.$SLUG`.
6. Report to the captain: how many intents you drained (for this workflow) and what you did with each.

If a record is malformed (not valid JSON, or an unknown `kind`), skip it but still advance the cursor past it, and note the skip to the captain — never block the loop on a bad record.
