#!/usr/bin/env bash
# ABOUTME: Claude Stop-hook wrapper — blocks the stop when Bridge intent is queued.
#
# Registered by hooks/hooks.json on Stop (synchronous, NOT async: an async hook
# cannot return a decision). It reads the Stop payload on stdin and delegates the
# decision to the spacedock binary, which resolves this session's workflow slug
# and emits {"decision":"block","reason":...} when captain intent is pending, or
# {} to let the session stop. This is the Claude durable-wake path: a parked FO
# drains queued intent in-session at the turn boundary, with no unsafe external
# session resume. Any failure degrades to no output (the session stops normally).
set -u

if [ -n "${SPACEDOCK_BIN:-}" ] && [ -x "${SPACEDOCK_BIN:-}" ]; then
  bin="${SPACEDOCK_BIN}"
elif command -v spacedock >/dev/null 2>&1; then
  bin="spacedock"
else
  exit 0
fi

"$bin" bridge inbox check --host claude 2>/dev/null || :
exit 0
