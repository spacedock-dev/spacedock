#!/usr/bin/env bash
# ABOUTME: Claude plugin hook wrapper for Bridge egress; observe-only and silent.
#
# Registered by hooks/hooks.json for Claude lifecycle events. The normalized
# Bridge schema and marker logic live in the spacedock binary, so host wrappers do
# not grow private JSON contracts.
set -u

if [ -n "${SPACEDOCK_BIN:-}" ] && [ -x "${SPACEDOCK_BIN:-}" ]; then
  bin="${SPACEDOCK_BIN}"
elif command -v spacedock >/dev/null 2>&1; then
  bin="spacedock"
else
  exit 0
fi

"$bin" bridge egress emit --host claude >/dev/null 2>&1 || :
exit 0
