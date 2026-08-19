#!/bin/sh
# ABOUTME: SessionStart hook, shared by the Claude and Codex plugin manifests
# ABOUTME: via hooks.json. Fires on every session start; acts only on
# ABOUTME: source=compact — the event after a COMPLETED compaction, not
# ABOUTME: PreCompact (which also fires on a compaction the host then
# ABOUTME: refuses). Gated on SPACEDOCK_BIN: that variable means "launcher-
# ABOUTME: marked session" — a plugin install by itself should not start
# ABOUTME: handing out First Officer instructions after every compaction;
# ABOUTME: only a session actually running a spacedock-launched host does.
# ABOUTME: Do not strip the gate as dead weight.
#
# A compacted First Officer keeps its narrative and loses its bindings: which
# binary, which contract version, which mods are registered, which workers
# are alive, and what durable state actually says. See
# docs/dev/.spacedock-state/force-boot-at-compaction-boundary.md for the
# incident this answers and the live capture that proves the event shape.
# Stdout here is injected as context for the resumed turn — this hook cannot
# force the next tool call, only remind before it.
[ -n "${SPACEDOCK_BIN:-}" ] || exit 0
src=$(cat | sed -n 's/.*"source"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -1)
[ "$src" = "compact" ] || exit 0
printf '%s\n' '{"hookSpecificOutput":{"hookEventName":"SessionStart","additionalContext":"COMPACTION BOUNDARY — your bindings are stale, not just your narrative.\n\nBefore any gate, merge, or state mutation, re-run the Startup procedure: resolve the binary, then `${SPACEDOCK_BIN:-spacedock} status --boot --identify --json`, and READ the whole boot record — not only the keys you came looking for. It carries the registered mods, ready-gate readiness, the binary version gate, PR state, and live team state. Any of those may contradict what your summary says.\n\nSpecifically distrust: which gates you believe are presented (check readiness, not memory), which workers you believe are alive, the SPACEDOCK_BIN path, and which contract version is installed."}}'
