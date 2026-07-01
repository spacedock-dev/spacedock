#!/usr/bin/env bash
# ABOUTME: Codex plugin hook wrapper for Bridge egress; observe-only and silent.
#
# Codex loads hooks from hooks/codex-hooks.json, not hooks/hooks.json, because the
# shared Claude hook file uses async:true and Codex skips async command hooks.
# This wrapper intentionally delegates to the public Spacedock egress command so
# Codex packaging does not grow its own private event schema.
set -u

if [ -n "${SPACEDOCK_BIN:-}" ] && [ -x "${SPACEDOCK_BIN:-}" ]; then
  bin="${SPACEDOCK_BIN}"
elif command -v spacedock >/dev/null 2>&1; then
  bin="spacedock"
else
  exit 0
fi

"$bin" bridge egress emit --host codex >/dev/null 2>&1 || :
exit 0
