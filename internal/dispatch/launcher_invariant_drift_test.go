// ABOUTME: Behavioral AC-3 — with a SPACEDOCK_BIN binary AND a DIFFERENT spacedock
// ABOUTME: on PATH both present, every helper invocation resolves the SPACEDOCK_BIN one.
package dispatch

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The launcher drift bug: a session launched with SPACEDOCK_BIN set runs its
// version gate on that binary, then later resolves a helper call (status, state
// ready, …) against a DIFFERENT `spacedock` on $PATH — a stale brew install whose
// command surface differs, yielding false subcommand-missing results. This test
// reproduces the exact condition the bug needs (both binaries present, different
// identities) and asserts the resolved launcher is the SPACEDOCK_BIN one for every
// helper subcommand, not just for a status probe.

// resolveHelperBinary runs `launcherExpr <helper-args>` with the given SPACEDOCK_BIN
// and a PATH whose `spacedock` is a distinct binary, returning what the resolved
// binary printed (each stub echoes its own identity tag plus its first arg).
func resolveHelperBinary(t *testing.T, launcherExpr, spacedockBin, pathDir string, helperArgs ...string) string {
	t.Helper()
	cmd := exec.Command("sh", "-c", launcherExpr+" "+strings.Join(helperArgs, " "))
	cmd.Env = append(environWithoutSpacedockBin(), "SPACEDOCK_BIN="+spacedockBin)
	cmd.Env = append(cmd.Env, "PATH="+pathDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("resolve helper %q failed: %v\n%s", helperArgs, err, out)
	}
	return strings.TrimSpace(string(out))
}

// twoBinaryFixture writes an executable SPACEDOCK_BIN binary tagged `env-bin` and a
// DIFFERENT `spacedock` on a fresh PATH dir tagged `path-bin`, so a resolved call's
// output reveals which binary ran. Returns (spacedockBin, pathDir).
func twoBinaryFixture(t *testing.T) (string, string) {
	t.Helper()
	envBin := writeExecutable(t, t.TempDir(), "spacedock-launch", "#!/bin/sh\necho env-bin:$1\n")
	pathDir := t.TempDir()
	writeExecutable(t, pathDir, "spacedock", "#!/bin/sh\necho path-bin:$1\n")
	return envBin, pathDir
}

// driftingLauncher is the BUG shape used as the RED control: it ignores SPACEDOCK_BIN
// and always runs bare `spacedock` from PATH — exactly the drift the invariant bans.
const driftingLauncher = `spacedock`

// helperSubcommands are the post-gate helper calls the regression touched (the false
// subcommand-missing surfaced on `state ready`/`state sweep`).
var helperSubcommands = [][]string{
	{"status", "--boot"},
	{"state", "ready"},
	{"state", "sweep"},
	{"dispatch", "build"},
	{"merge", "guard"},
}

// TestLauncherInvariantNoPathDriftAcrossHelpers (AC-3, behavioral) proves the real
// LauncherCommand() resolves the SPACEDOCK_BIN binary for EVERY helper subcommand
// when a different `spacedock` is also on PATH. The driftingLauncher control proves
// the assertion actually catches the bug (it RUNS the wrong binary), so a green real
// result is not vacuous.
func TestLauncherInvariantNoPathDriftAcrossHelpers(t *testing.T) {
	envBin, pathDir := twoBinaryFixture(t)

	for _, helper := range helperSubcommands {
		name := strings.Join(helper, "_")

		// RED control: the drift shape resolves the PATH binary — the bug. The
		// assertion MUST observe `path-bin` here, proving it can tell the binaries
		// apart and would fail a real drift.
		t.Run("drift_control_"+name, func(t *testing.T) {
			got := resolveHelperBinary(t, driftingLauncher, envBin, pathDir, helper...)
			if !strings.HasPrefix(got, "path-bin:") {
				t.Fatalf("drift control for %v resolved %q, expected the PATH binary (`path-bin:`) — the assertion cannot distinguish the binaries", helper, got)
			}
		})

		// GREEN: the real launcher resolves the SPACEDOCK_BIN binary despite the
		// different PATH `spacedock` — no mid-session drift.
		t.Run(name, func(t *testing.T) {
			got := resolveHelperBinary(t, LauncherCommand(), envBin, pathDir, helper...)
			if !strings.HasPrefix(got, "env-bin:") {
				t.Fatalf("helper %v resolved %q — the launcher drifted to the PATH `spacedock` instead of the pinned SPACEDOCK_BIN binary (`env-bin:`)", helper, got)
			}
		})
	}
}

// TestLauncherInvariantResolvesAbsoluteSpacedockBinPath asserts the resolved binary
// is the EXACT SPACEDOCK_BIN path, not merely "some env-bin" — the strongest form of
// AC-3's "keeps resolving the SPACEDOCK_BIN path" claim.
func TestLauncherInvariantResolvesAbsoluteSpacedockBinPath(t *testing.T) {
	pathDir := t.TempDir()
	writeExecutable(t, pathDir, "spacedock", "#!/bin/sh\necho path-bin\n")
	// The env binary prints its own resolved path ($0) so we can assert identity.
	envBin := writeExecutable(t, t.TempDir(), "spacedock-launch", "#!/bin/sh\necho \"$0\"\n")

	got := resolveHelperBinary(t, LauncherCommand(), envBin, pathDir, "state", "ready")
	if got != envBin {
		t.Fatalf("launcher resolved %q, want the SPACEDOCK_BIN path %q", got, filepath.Clean(envBin))
	}
}
