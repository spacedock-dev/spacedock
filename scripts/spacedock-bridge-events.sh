#!/usr/bin/env bash
# ABOUTME: Spacedock plugin hook — emit a normalized FO event line to _bridge/events.jsonl
# ABOUTME: for the external Bridge command-center UI to tail. Observe-only; never blocks the session.
#
# Registered for SessionStart/UserPromptSubmit/PostToolUse/Notification/Stop/SubagentStop
# in hooks/hooks.json (all async). It writes a stable, Spacedock-owned event contract so
# Bridge does not have to couple to Claude Code's internal transcript JSONL format.
#
# Events land in «session-cwd»/_bridge/events.jsonl — the same _bridge/ dir the bridge-inbox
# mod drains and that Bridge resolves from the repo root. agent_id/agent_type are empty for
# the main FO session and set for ensign subagents, so Bridge can tell FO vs ensign activity.
#
# Honesty + safety: this only observes. It must never alter the session, so it always exits 0
# and degrades to a silent no-op when jq is unavailable, the payload lacks a cwd, or the write
# fails — a telemetry side-channel must not be able to break the FO.
set -u

# No jq → no-op. (Spacedock already assumes jq for its gh-driven hooks.)
command -v jq >/dev/null 2>&1 || exit 0

payload="$(cat 2>/dev/null)" || exit 0
[ -n "$payload" ] || exit 0

cwd="$(printf '%s' "$payload" | jq -r '.cwd // empty' 2>/dev/null)" || exit 0
[ -n "$cwd" ] || exit 0

dir="$cwd/_bridge"
mkdir -p "$dir" 2>/dev/null || exit 0

ts="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

# Normalize to a stable line. Only generic, non-sensitive fields are emitted (no tool inputs
# or outputs, no prompt text) — liveness, not content.
printf '%s' "$payload" | jq -c \
  --arg ts "$ts" \
  '{
     ts: $ts,
     event: (.hook_event_name // "unknown"),
     session_id: (.session_id // ""),
     agent_id: (.agent_id // ""),
     agent_type: (.agent_type // ""),
     detail: {
       tool: (.tool_name // ""),
       source: (.source // "")
     }
   }' >> "$dir/events.jsonl" 2>/dev/null || exit 0

exit 0
