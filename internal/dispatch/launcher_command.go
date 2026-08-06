// ABOUTME: Shared rendering for the host-launch safehouse regression fixture.
// ABOUTME: Dispatch artifacts do not use this ambient launcher expression.
package dispatch

const launcherCommandExpr = `spacedock_launcher() { if [ -n "${SPACEDOCK_BIN:-}" ] && [ -x "${SPACEDOCK_BIN:-}" ]; then "$SPACEDOCK_BIN" "$@"; else spacedock "$@"; fi; }; spacedock_launcher`

// LauncherCommand returns the host-launch environment fallback expression used
// by the safehouse propagation fixture. Worker dispatches pin an absolute path.
func LauncherCommand() string { return launcherCommandExpr }
