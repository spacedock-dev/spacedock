// ABOUTME: Shared rendering for Spacedock helper invocations emitted into agent prompts.
// ABOUTME: Keeps launch-provided SPACEDOCK_BIN identity while preserving PATH fallback.
package dispatch

const launcherCommandExpr = "${SPACEDOCK_BIN:-spacedock}"

func launcherCommand() string { return launcherCommandExpr }
