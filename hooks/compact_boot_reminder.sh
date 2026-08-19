#!/bin/sh
# ABOUTME: Claude Code SessionStart hook. Fires on every session start; acts
# ABOUTME: only on source=compact — the event after a COMPLETED compaction.
#
# PreCompact also fires on a compaction the host then refuses ("Not enough
# messages to compact"), so it is the wrong event for this: SessionStart with
# source=compact fires exactly once, after compaction actually finished. See
# docs/dev/.spacedock-state/force-boot-at-compaction-boundary.md for the
# incident this answers and the live capture that proves the event shape.
#
# A compacted First Officer keeps its narrative and loses its bindings: which
# binary, which contract version, which mods are registered, which workers
# are alive, and what durable state actually says. Stdout here is injected as
# context for the resumed turn — this hook cannot force the next tool call,
# only remind before it.
src=$(cat | sed -n 's/.*"source"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -1)
[ "$src" = "compact" ] || exit 0
cat <<'MSG'
COMPACTION BOUNDARY — your bindings are stale, not just your narrative.

Before any gate, merge, or state mutation, re-run the Startup procedure:
resolve the binary, then `${SPACEDOCK_BIN:-spacedock} status --boot --identify
--json`, and READ the whole boot record — not only the keys you came looking
for. It carries the registered mods, ready-gate readiness, the binary version
gate, PR state, and live team state. Any of those may contradict what your
summary says.

Specifically distrust: which gates you believe are presented (check
readiness, not memory), which workers you believe are alive, the
SPACEDOCK_BIN path, and which contract version is installed.
MSG
