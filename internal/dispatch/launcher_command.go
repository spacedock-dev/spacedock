// ABOUTME: Shared rendering for Spacedock helper invocations emitted into agent prompts.
// ABOUTME: Keeps launch-provided SPACEDOCK_BIN identity while preserving PATH fallback.
package dispatch

const launcherCommandExpr = `spacedock_launcher() { if [ -n "${SPACEDOCK_BIN:-}" ] && [ -x "${SPACEDOCK_BIN:-}" ]; then "$SPACEDOCK_BIN" "$@"; else spacedock "$@"; fi; }; spacedock_launcher`

// LauncherCommand returns the shell expression that resolves a Spacedock helper
// invocation: prefer an executable SPACEDOCK_BIN, otherwise fall back to the
// `spacedock` on PATH. It is the single source dispatch prompts and the
// safehouse env-scrub regression oracle both render, so neither drifts.
func LauncherCommand() string { return launcherCommandExpr }
