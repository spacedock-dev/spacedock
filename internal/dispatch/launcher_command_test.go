// ABOUTME: Locks the env-aware launcher command token generated into dispatch prompts.
// ABOUTME: Ensures ensign fetch commands prefer SPACEDOCK_BIN with spacedock PATH fallback.
package dispatch

import "testing"

func TestLauncherCommandUsesSpacedockBinFallbackExpression(t *testing.T) {
	if got, want := launcherCommand(), "${SPACEDOCK_BIN:-spacedock}"; got != want {
		t.Fatalf("launcherCommand() = %q, want %q", got, want)
	}
}
