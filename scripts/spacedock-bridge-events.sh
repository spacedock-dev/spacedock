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
events="$dir/events.jsonl"
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
   }' >> "$events" 2>/dev/null || exit 0

# Best-effort size cap: a liveness side-channel must not grow without bound (PostToolUse
# fires on every tool call). When the log passes max_lines, keep only the most recent
# keep_lines. Lock-free and best-effort by design — this is liveness, not a ledger, so
# losing a few lines to a concurrent append during the rare trim is acceptable, and every
# step degrades to a no-op ($$ keeps the temp unique across concurrent async hooks; any
# failure leaves the existing log untouched). Never block or break the FO.
max_lines=2000
keep_lines=1000
n="$(wc -l < "$events" 2>/dev/null || echo 0)"
if [ "${n:-0}" -gt "$max_lines" ] 2>/dev/null; then
  tmp="$events.tmp.$$"
  if tail -n "$keep_lines" "$events" > "$tmp" 2>/dev/null; then
    mv -f "$tmp" "$events" 2>/dev/null || rm -f "$tmp" 2>/dev/null
  else
    rm -f "$tmp" 2>/dev/null
  fi
fi

exit 0
