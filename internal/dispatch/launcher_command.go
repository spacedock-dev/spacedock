// ABOUTME: Shared rendering for Spacedock helper invocations emitted into agent prompts.
// ABOUTME: Keeps launch-provided SPACEDOCK_BIN identity while preserving PATH fallback.
package dispatch

const launcherCommandExpr = `spacedock_launcher() { if [ -n "${SPACEDOCK_BIN:-}" ] && [ -x "${SPACEDOCK_BIN:-}" ]; then "$SPACEDOCK_BIN" "$@"; else spacedock "$@"; fi; }; spacedock_launcher`

func launcherCommand() string { return launcherCommandExpr }
