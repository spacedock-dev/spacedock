---
name: bridge-inbox
description: Drain captain intent queued by the Bridge command-center UI (_bridge/inbox.jsonl) and surface it to the FO each tick
version: 0.1.0
fo-realm: "FO realm — the FO maintains this file directly; it is FO process (the seam to the Bridge UI), not product built under the dev workflow."
---

# Bridge Inbox

[Bridge](https://github.com/spacedock-dev/bridge) is a read-only command-center UI over this fleet. It cannot push into a running FO session (a Claude Code session has no inbound API), so it writes captain intent to a durable, **append-only** inbox at `_bridge/inbox.jsonl` (relative to the FO's working directory — the repo root where `spacedock claude` was launched, the same root Bridge resolves). This hook drains that inbox on the FO's own loop ticks and acts on the intent.

**This is a pull, not a push.** Delivery latency is one FO loop cadence: the captain's intent is read whenever the FO next idles or boots, never instantly. Bridge's UI says as much ("queued — read on the FO's next tick"); do not promise synchronous delivery.

**Inbox record schema** (one JSON object per line, written by Bridge):

```
{"ts": "<rfc3339>", "kind": "tell" | "conn", "text": "<string>", "granted": <bool, conn only>}
```

**Consume by cursor, never by rewrite.** Bridge only appends; this hook only advances `_bridge/.inbox-cursor` (the count of inbox lines already processed). Neither side rewrites `inbox.jsonl`, so concurrent Bridge appends and FO reads never clobber, and re-firing the hook with no new lines is a no-op.

## Hook: startup

Drain any intent the captain queued while no FO session was attached, so a freshly-booted FO picks up standing instructions before its first dispatch. Run the same drain procedure as the idle hook below.

## Hook: idle

Drain newly-queued captain intent, if any:

1. If `_bridge/inbox.jsonl` does not exist, skip — no Bridge is attached.
2. Read the cursor (lines already processed): `CURSOR=$(cat _bridge/.inbox-cursor 2>/dev/null || echo 0)`.
3. Snapshot the current line count and read exactly the new records (bounding the read so a concurrent Bridge append can't make the cursor skip a line):
   ```
   NEW=$(wc -l < _bridge/inbox.jsonl | tr -d ' ')
   sed -n "$((CURSOR + 1)),${NEW}p" _bridge/inbox.jsonl
   ```
   If `NEW` is not greater than `CURSOR`, there is nothing new — skip (idempotent).
4. For each new record, in order, parse `kind` / `text` / `granted` and act:
   - **`kind == "tell"`** — the captain sent you a message. Treat `text` as a directive or clarification for this tick: act on it as you would a captain instruction (commission or clear work, answer the implied question, adjust course), and acknowledge it to the captain.
   - **`kind == "conn"`** — a conn-handover change. `granted: true` → adopt the conn within the stated goal `text`: drive the entities the conn covers to done without stopping at their gates, per the conn rules in `first-officer-shared-core` (escalations remain non-delegable and still surface to the captain). `granted: false` → take the conn back: stop at every gate for the captain's call again.
5. Advance the cursor to the snapshot you read: `echo "$NEW" > _bridge/.inbox-cursor`.
6. Report to the captain: how many intents you drained and what you did with each.

If a record is malformed (not valid JSON, or an unknown `kind`), skip it but still advance the cursor past it, and note the skip to the captain — never block the loop on a bad record.
