#!/bin/sh
# ABOUTME: SessionStart(compact) reminder hook; SPACEDOCK_BIN gates it to launcher-marked sessions only — do not strip as dead weight.
[ -n "${SPACEDOCK_BIN:-}" ] || exit 0
src=$(cat | sed -n 's/.*"source"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -1)
[ "$src" = "compact" ] || exit 0
printf '%s\n' '{"hookSpecificOutput":{"hookEventName":"SessionStart","additionalContext":"COMPACTION BOUNDARY — your bindings are stale, not just your narrative.\n\nBefore any gate, merge, or state mutation, re-run the Startup procedure: resolve the binary, then `${SPACEDOCK_BIN:-spacedock} status --boot --identify --json`, and READ the whole boot record — not only the keys you came looking for. It carries the registered mods, ready-gate readiness, the binary version gate, PR state, and live team state. Any of those may contradict what your summary says.\n\nSpecifically distrust: which gates you believe are presented (check readiness, not memory), which workers you believe are alive, the SPACEDOCK_BIN path, and which contract version is installed."}}'
