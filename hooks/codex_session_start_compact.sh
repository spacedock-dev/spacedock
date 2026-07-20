#!/bin/sh
# ABOUTME: Supplies model-visible compact recovery only to launcher-marked Codex sessions.
[ -n "${SPACEDOCK_BIN:-}" ] || exit 0
printf '%s\n' '{"hookSpecificOutput":{"hookEventName":"SessionStart","additionalContext":"Spacedock: reread the authoritative `spacedock:first-officer` contract and reconcile durable workflow state with live worker state before the next workflow effect."}}'
