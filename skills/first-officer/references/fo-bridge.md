# FO Bridge — Permission-Block Alerts

Deferred reference for the one Bridge signal the boot-resident core does not carry: the FO→captain **permission-block alert**. Loaded when the FO hits a permission or sandbox block, not at boot. The heartbeat, the captain-intent drain, the terminal acks, and the gate/status cards all live in the `bridge-seam` mod (`_mods/bridge-seam.md`, its `## Agent Prompt`); this file covers only the alert append. The full seam is Bridge's `docs/seam-contract.md` (§2.8 is the alert shape); the Spacedock-local overview is `docs/dev/bridge-seam.md`.

## Permission Blocks and Bridge Alerts

When a host sandbox or permission boundary blocks a workflow action that would otherwise be valid to attempt, surface it as a high-priority Bridge interrupt **before parking the loop**. There is no packaged permission-alert verb — this is a direct file write, exactly like the rest of the seam: append one JSONL line to `_bridge/fo-alerts.jsonl` (relative to the fleet root where the FO launched), using the §2.8 shape verbatim:

```json
{"schema":1,"id":"al_3c","ts":"2026-07-14T18:08:00Z","kind":"permission-request","workflow":"linear-drc-ship","entity":"drc-3467","reason":"rm outside repo","command":"rm -rf /tmp/x","prefix_rule":["rm -rf /tmp/"],"status":"open"}
```

- `id` — a non-empty stable correlator you mint (e.g. `al_<short-hex>`). Bridge **drops** any line whose `id` is empty or whose `kind` is not `permission-request`. This `id` is the join key the captain's answer echoes back.
- `kind` — always `permission-request` (the only alert kind today).
- `reason` / `command` — one-line why-blocked and the command summary the captain reads.
- `prefix_rule` — optional proposed allow-rule prefix (a string array) for an `approve-rule` answer.
- `workflow` / `entity` / `host` / `session_id` — optional provenance/scoping; pass them so the card attributes the block.
- `status` — always write `open`; Bridge overlays the resolution itself from the matching inbox `permission-decision`.

## Answering the block

The captain's answer arrives as a `permission-decision` intent on `_bridge/inbox.jsonl` whose `request_id` matches your alert `id`. You drain and apply it through the ordinary **Drain procedure** in the `bridge-seam` mod — no verb here either. The `value` decides the retry:

- `deny` → do not retry; append a `permission-ack` with `status:"denied"`.
- `approve-once` → retry the exact blocked action once through the runtime's escalation path; append `permission-ack` `status:"accepted"` first (or `status:"blocked"` if the retry cannot start).
- `approve-rule` → retry with the alert's `prefix_rule` if present, else treat as `approve-once`; ack `accepted`.

A Bridge approval is FO intent, not a bypass of a host-native security prompt: if the host still presents its own approval dialog, honor that dialog normally.
