# Bridge seam (Spacedock side)

The Bridge command center observes and steers a Spacedock First Officer through
append-only files under a fleet's `_bridge/` directory. **Bridge owns the
schema**; the authoritative, file-by-file contract is Bridge's
`docs/seam-contract.md` (in the `bridge` repo). This document is the
Spacedock-local overview: what Spacedock produces, and how.

## The one rule

Bridge consumes **files, not verbs.** Every FO→Bridge signal (heartbeat, drain
cursor, terminal ack, gate-review, status, permission alert) is a direct file
write the FO performs itself, per the `bridge-seam` mod's `## Agent Prompt` and
the per-host `## Bridge seam` runtime sections. The consolidation retired the
packaged `inbox drain|ack|commit`, `alert`, and `initiate` verbs — an agent with
only Bridge's `docs/seam-contract.md` can implement the producer side.

## What Spacedock still ships (the irreducible core)

Three things genuinely need the harness and remain as `spacedock bridge …`
entrypoints — all hook- or daemon-invoked, never called by the FO:

| Entrypoint | Invoked by | Purpose |
|---|---|---|
| `bridge egress emit --host <h>` | plugin lifecycle hooks | normalize a harness turn payload → one `_bridge/events.jsonl` line; derive the `_bridge/sessions/<id>.json` marker on an ensign's entity Read |
| `bridge inbox check --host claude` | synchronous Claude `Stop` hook | return `{"decision":"block"}` while this session has queued intent, so a Claude FO drains before stopping (its only durable wake — a parked Claude transcript cannot be resumed out-of-band) |
| `bridge ingress wake --host codex` | Bridge / daemon | resume a parked Codex session via `codex exec resume` and prompt it to drain |

Plus one contract-conformant FO mod (`mods/bridge-seam.md`) carrying the file
protocol, and the FO/ensign prose that binds each host's producer.

## Per-host capability matrix

| Capability | Claude | Codex | Pi |
|---|---|---|---|
| `events.jsonl` liveness/activity egress | PRESENT (`hooks/hooks.json`) | PACKAGED (`hooks/codex-hooks.json`) | PACKAGED (`.pi/extensions/spacedock.ts`) |
| deterministic `sessions/<id>.json` running marker | PRESENT | not yet claimed | not yet claimed |
| durable wake | in-session (sync Stop hook) | external (`ingress wake`) | none yet (queued count until self-drain) |
| heartbeat session-id source | `$CLAUDE_CODE_SESSION_ID` | `$CODEX_THREAD_ID` | runtime id (often empty today) |

**Load-bearing:** the `_bridge/fo.<slug>.json` heartbeat MUST carry the harness
session id. The Claude Stop-hook check resolves which slugs belong to a stopping
session by matching the Stop payload's `session_id` against the heartbeats; a
heartbeat that omits or mismatches that id makes the check resolve nothing, so a
Claude FO is never blocked and queued captain intent is delivered never.

## Hook registrations

- `.claude-plugin/plugin.json` → `hooks/hooks.json`: six async egress events
  (SessionStart, UserPromptSubmit, PostToolUse, Notification, Stop, SubagentStop)
  via `scripts/spacedock-bridge-events.sh`, plus a **synchronous** Stop entry
  (`scripts/spacedock-bridge-inbox-check.sh`) — the block-until-drained wake.
- `.codex-plugin/plugin.json` → `hooks/codex-hooks.json`: six non-async
  command hooks calling `bridge egress emit --host codex`. No Stop-block (Codex
  uses external wake).
- `.pi/extensions/spacedock.ts`: forwards Pi lifecycle payloads to
  `bridge egress emit --host pi`.

## Check the inbox line cap

`bridge inbox check` scans `_bridge/inbox.jsonl` with a 1 MiB per-line cap; an
oversized intent line makes the check honestly emit `{}` (no block) rather than
error. Pair with Bridge's own reader cap follow-up.

## Not shipped here (intentional)

`fo-feed.jsonl` is not produced — it is redundant with git narration plus the
marker-derived feed (see Bridge `docs/seam-contract.md` §5). Getting the
`bridge-seam` mod into non-dogfood workflows still needs commission/refit
scaffolding (a named follow-up); today it ships canonically at `mods/` and as a
dogfood copy under `docs/dev/_mods/`.
