#!/usr/bin/env bash
# ABOUTME: Spacedock plugin hook — emit a normalized FO event line to _bridge/events.jsonl
# ABOUTME: for the external Bridge command-center UI to tail. Observe-only; never blocks the session.
#
# Registered for SessionStart/UserPromptSubmit/PostToolUse/Notification/Stop/SubagentStop
# in hooks/hooks.json (all async). It writes a stable, Spacedock-owned event contract so
# Bridge does not have to couple to Claude Code's internal transcript JSONL format.
#
# This script is the CLAUDE producer-binding for the egress surface; the harness-neutral
# schema for events.jsonl + the _bridge/sessions/<sid>.json marker (and where Codex/Pi are
# ABSENT/TODO) is docs/dev/bridge-egress-contract.md. Keep the emitted line shapes in sync
# with that contract.
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

# ── Session→entity marker (DETERMINISTIC running-badge source) ───────────────────────
# Bridge maps a live working session to its ship via «cwd»/_bridge/sessions/<session_id>.json.
# Derive it HERE — from the hook, which fires on every tool call — instead of relying on the
# ensign to run a first-action shell (it skips that ~3/4 of the time). On an ENSIGN's first
# Read of its entity file (.../docs/spacedock/<workflow>/<slug>.md, flat or <slug>/index.md),
# record {session_id, entity, workflow}: the path carries BOTH the workflow (so Bridge's join
# is collision-free across workflows that reuse a ticket id) and the full slug. First-write-
# wins per session, so the ensign's OWN entity — read before any duplicate-check sibling reads
# — is what gets recorded. Observe-only and best-effort: every step degrades to a no-op.
m_sid="$(printf '%s' "$payload" | jq -r '.session_id // empty' 2>/dev/null)"
m_type="$(printf '%s' "$payload" | jq -r '.agent_type // empty' 2>/dev/null)"
m_evt="$(printf '%s' "$payload" | jq -r '.hook_event_name // empty' 2>/dev/null)"
m_tool="$(printf '%s' "$payload" | jq -r '.tool_name // empty' 2>/dev/null)"
case "$m_sid" in *[!A-Za-z0-9._-]*) m_sid="" ;; esac   # unsafe id → skip (never escape _bridge/)
if [ -n "$m_sid" ] && [ "$m_type" = "spacedock:ensign" ] && [ "$m_evt" = "PostToolUse" ] && [ "$m_tool" = "Read" ]; then
  marker="$dir/sessions/$m_sid.json"
  if [ ! -f "$marker" ]; then   # first-write-wins → the ensign's own entity (read first)
    fp="$(printf '%s' "$payload" | jq -r '.tool_input.file_path // empty' 2>/dev/null)"
    case "$fp" in
      */docs/spacedock/*/*.md|docs/spacedock/*/*.md)   # absolute OR repo-relative (the FO passes a relative {entity_file_path})
        # Workflow = the path segment right AFTER docs/spacedock/ — robust to a
        # split-root entity nested under <wf>/.spacedock-state/<slug>.md (taking the
        # entity's parent dir would wrongly yield ".spacedock-state").
        rest="${fp##*/docs/spacedock/}"   # absolute / nested: strip through the last /docs/spacedock/
        rest="${rest#docs/spacedock/}"    # repo-relative: strip the leading docs/spacedock/
        wf="${rest%%/*}"                  # first segment after docs/spacedock/ = the workflow dir
        if [ "$(basename "$fp")" = "index.md" ]; then   # folder entity: .../<slug>/index.md
          slug="$(basename "$(dirname "$fp")")"
        else                                            # flat entity: .../<slug>.md
          slug="$(basename "$fp" .md)"
        fi
        case "$fp"   in */_archive/*) wf="" ;; esac                  # never mark an archived (terminal) entity
        case "$wf"   in ""|_*|*[!A-Za-z0-9._-]*) wf="" ;; esac        # skip _-dirs, unsafe
        case "$slug" in ""|README|*[!A-Za-z0-9._-]*) slug="" ;; esac
        if [ -n "$slug" ] && [ -n "$wf" ]; then
          mkdir -p "$dir/sessions" 2>/dev/null &&
          printf '{"session_id":"%s","entity":"%s","workflow":"%s"}\n' "$m_sid" "$slug" "$wf" > "$marker" 2>/dev/null || :
        fi
        ;;
    esac
  fi
fi

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
