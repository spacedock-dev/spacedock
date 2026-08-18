// ABOUTME: End-to-end env-pass smoke: through a safehouse that scrubs SPACEDOCK_BIN
// ABOUTME: by default, --env-pass SPACEDOCK_BIN forwards it so the real launcher resolves the launched binary.
package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/dispatch"
	"github.com/spacedock-dev/spacedock/internal/safehouse"
)

// envPassSafehouse writes a fake `safehouse` that SCRUBS named test variables by
// default, then honors comma-separated `--env-pass NAMES` and Safehouse's native
// SAFEHOUSE_ENV_PASS configuration. Each named value is forwarded from the host
// env it inherited, on top of sanitized defaults. The fake lets the smoke prove
// that a named allowlist matters rather than passing through ambient process env.
func envPassSafehouse(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "safehouse")
	body := `#!/bin/sh
# Capture the inherited value before scrubbing (the host env safehouse sees).
saved_bin="${SPACEDOCK_BIN-}"
saved_zellij_pane_id="${ZELLIJ_PANE_ID-}"
saved_zellij_session_name="${ZELLIJ_SESSION_NAME-}"
saved_tmux="${TMUX-}"
saved_tmux_pane="${TMUX_PANE-}"
saved_extra_target="${EXTRA_TARGET-}"
pass="${SAFEHOUSE_ENV_PASS-}"
unset SPACEDOCK_BIN ZELLIJ_PANE_ID ZELLIJ_SESSION_NAME TMUX TMUX_PANE EXTRA_TARGET
append_pass() {
  if [ -n "$pass" ]; then pass="$pass,$1"; else pass="$1"; fi
}
while [ "$#" -gt 0 ] && [ "$1" != "--" ]; do
  case "$1" in
    --env-pass) shift; append_pass "$1" ;;
    --env-pass=*) append_pass "${1#--env-pass=}" ;;
  esac
  shift
done
if [ "$1" = "--" ]; then shift; fi
# Honor each named pass-through: forward values from the saved host env.
case ",$pass," in *,SPACEDOCK_BIN,*) export SPACEDOCK_BIN="$saved_bin" ;; esac
case ",$pass," in *,ZELLIJ_PANE_ID,*) export ZELLIJ_PANE_ID="$saved_zellij_pane_id" ;; esac
case ",$pass," in *,ZELLIJ_SESSION_NAME,*) export ZELLIJ_SESSION_NAME="$saved_zellij_session_name" ;; esac
case ",$pass," in *,TMUX,*) export TMUX="$saved_tmux" ;; esac
case ",$pass," in *,TMUX_PANE,*) export TMUX_PANE="$saved_tmux_pane" ;; esac
case ",$pass," in *,EXTRA_TARGET,*) export EXTRA_TARGET="$saved_extra_target" ;; esac
exec "$@"
`
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// launcherProbe writes a fake host whose only job is to resolve a spacedock
// helper through the REAL spacedock_launcher expression (dispatch.LauncherCommand)
// and print whatever that resolution runs. It is the in-session helper call:
// prefer an executable SPACEDOCK_BIN, else fall back to `spacedock` on PATH.
func launcherProbe(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "probe")
	body := "#!/bin/sh\n" + dispatch.LauncherCommand() + " status\n"
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// markerBinary writes an executable that prints a single marker line and ignores
// its args — stands in for a spacedock binary so the probe's resolution is
// observable (which binary ran), not just present.
func markerBinary(t *testing.T, path, marker string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("#!/bin/sh\nprintf '%s\\n' "+marker+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
}

// TestSafehouseEnvPassForwardsSpacedockBin is the AC-5 regression oracle. Through
// a safehouse that SCRUBS SPACEDOCK_BIN by default, the `--env-pass SPACEDOCK_BIN`
// wrap flag the front door adds forwards the launching binary so the real
// spacedock_launcher resolves it. The control case — the SAME wrap WITHOUT
// `--env-pass` (env-only propagation) — falls back to the PATH `spacedock`,
// proving the test fails for the right reason absent the flag rather than passing
// trivially. The inner program after `--` is the host probe directly (no wrapper),
// matching the production argv shape.
func TestSafehouseEnvPassForwardsSpacedockBin(t *testing.T) {
	dir := t.TempDir()
	safehousePath := envPassSafehouse(t, dir)
	probe := launcherProbe(t, dir)

	// The launched binary (what SPACEDOCK_BIN points at) and a DIFFERENT spacedock
	// on PATH, each printing a distinct marker so the resolution is observable.
	launchedBin := filepath.Join(dir, "spacedock-launched")
	markerBinary(t, launchedBin, "LAUNCHED")
	pathDir := filepath.Join(dir, "pathbin")
	if err := os.Mkdir(pathDir, 0o755); err != nil {
		t.Fatal(err)
	}
	markerBinary(t, filepath.Join(pathDir, "spacedock"), "PATH-FALLBACK")
	// Build a clean env (no inherited PATH or SPACEDOCK_BIN to collide): pathDir
	// holds the ONLY `spacedock` on PATH; the system dirs follow so `sh` and the
	// coreutils the probe/wrapper need still resolve. SPACEDOCK_BIN is set on the
	// safehouse process, exactly as launchEnv sets it on the launched process.
	base := withoutEnv(withoutEnv(os.Environ(), "PATH"), spacedockBinEnv)
	env := append(base,
		"PATH="+pathDir+":/usr/bin:/bin",
		spacedockBinEnv+"="+launchedBin)

	t.Run("--env-pass forwards SPACEDOCK_BIN through the scrub", func(t *testing.T) {
		// The wrap as runClaude composes it: launcherBinEnvPassFlags (driven by
		// executablePath → launchedBin) is handed to the Safehouse wrapper. The
		// inner program after `--` is the probe directly (no /usr/bin/env wrapper).
		withExecutablePath(t, launchedBin, nil)
		argv := safehouse.Wrap([]string{probe}, launcherBinEnvPassFlags())
		out := runWrapped(t, safehousePath, argv[1:], env)
		if out != "LAUNCHED" {
			t.Fatalf("probe resolved %q, want LAUNCHED (--env-pass must forward SPACEDOCK_BIN through the scrub)", out)
		}
	})

	t.Run("without --env-pass it falls back to PATH (regression control)", func(t *testing.T) {
		// The SAME wrap WITHOUT --env-pass — env-only propagation. The scrub drops
		// SPACEDOCK_BIN, so the probe falls back to the PATH spacedock. This is the bug
		// the --env-pass flag fixes; it proves the oracle is not trivial.
		argv := safehouse.Wrap([]string{probe}, nil)
		out := runWrapped(t, safehousePath, argv[1:], env)
		if out != "PATH-FALLBACK" {
			t.Fatalf("env-only probe resolved %q, want PATH-FALLBACK (the scrub must drop SPACEDOCK_BIN absent --env-pass)", out)
		}
	})
}

// TestSafehouseEnvPassForwardsTerminalTargetingMetadata proves the presence
// filter through a scrubbing wrapper for two host families: a tmux-pair
// parent's child sees the pair's values and no Zellij names, and the
// Zellij-pair case still forwards — each with the other family's names never
// presented as empty (they're never named in --env-pass to begin with).
func TestSafehouseEnvPassForwardsTerminalTargetingMetadata(t *testing.T) {
	dir := t.TempDir()
	safehousePath := envPassSafehouse(t, dir)
	probe := terminalTargetingProbe(t, dir)

	t.Run("tmux pair forwards exact inherited values; Zellij names stay unset", func(t *testing.T) {
		setTmuxTargetingEnv(t)
		parent := []string{
			"PATH=/usr/bin:/bin",
			"TMUX=/tmp/tmux-501/default,12345,0",
			"TMUX_PANE=%3",
		}
		argv := safehouse.Wrap([]string{probe}, []string{"--env-pass", spacedockBinEnv})
		out := runWrapped(t, safehousePath, argv[1:], parent)
		want := "TMUX=/tmp/tmux-501/default,12345,0\nTMUX_PANE=%3\nZELLIJ_PANE_ID=<unset>\nZELLIJ_SESSION_NAME=<unset>"
		if out != want {
			t.Fatalf("terminal metadata = %q, want %q", out, want)
		}
	})

	t.Run("Zellij pair still forwards; tmux names stay unset", func(t *testing.T) {
		setZellijTargetingEnv(t)
		parent := []string{
			"PATH=/usr/bin:/bin",
			"ZELLIJ_PANE_ID=51",
			"ZELLIJ_SESSION_NAME=excellent-pheasant",
		}
		argv := safehouse.Wrap([]string{probe}, []string{"--env-pass", spacedockBinEnv})
		out := runWrapped(t, safehousePath, argv[1:], parent)
		want := "TMUX=<unset>\nTMUX_PANE=<unset>\nZELLIJ_PANE_ID=51\nZELLIJ_SESSION_NAME=excellent-pheasant"
		if out != want {
			t.Fatalf("terminal metadata = %q, want %q", out, want)
		}
	})

	t.Run("native global allowlist composes with the built-in allowance", func(t *testing.T) {
		setZellijTargetingEnv(t)
		probe := extraTargetProbe(t, dir)
		env := []string{
			"PATH=/usr/bin:/bin",
			"ZELLIJ_PANE_ID=51",
			"ZELLIJ_SESSION_NAME=excellent-pheasant",
			"SAFEHOUSE_ENV_PASS=EXTRA_TARGET",
			"EXTRA_TARGET=operator-choice",
		}
		argv := safehouse.Wrap([]string{probe}, []string{"--env-pass", spacedockBinEnv})
		out := runWrapped(t, safehousePath, argv[1:], env)
		if out != "EXTRA_TARGET=operator-choice" {
			t.Fatalf("global Safehouse env pass-through = %q, want %q", out, "EXTRA_TARGET=operator-choice")
		}
	})
}

func setZellijTargetingEnv(t *testing.T) {
	t.Helper()
	clearTerminalTargetingEnv(t)
	t.Setenv("ZELLIJ_PANE_ID", "51")
	t.Setenv("ZELLIJ_SESSION_NAME", "excellent-pheasant")
}

func setTmuxTargetingEnv(t *testing.T) {
	t.Helper()
	clearTerminalTargetingEnv(t)
	t.Setenv("TMUX", "/tmp/tmux-501/default,12345,0")
	t.Setenv("TMUX_PANE", "%3")
}

func clearTerminalTargetingEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{"ZELLIJ_PANE_ID", "ZELLIJ_SESSION_NAME", "TMUX", "TMUX_PANE"} {
		value, present := os.LookupEnv(key)
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if present {
				_ = os.Setenv(key, value)
				return
			}
			_ = os.Unsetenv(key)
		})
	}
}

func terminalTargetingProbe(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "terminal-probe")
	body := `#!/bin/sh
printf 'TMUX=%s\n' "${TMUX-<unset>}"
printf 'TMUX_PANE=%s\n' "${TMUX_PANE-<unset>}"
printf 'ZELLIJ_PANE_ID=%s\n' "${ZELLIJ_PANE_ID-<unset>}"
printf 'ZELLIJ_SESSION_NAME=%s\n' "${ZELLIJ_SESSION_NAME-<unset>}"
`
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func extraTargetProbe(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "extra-target-probe")
	body := `#!/bin/sh
printf 'EXTRA_TARGET=%s\n' "${EXTRA_TARGET-<unset>}"
`
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// runWrapped execs the fake safehouse with the wrapped argv (everything after the
// `safehouse` token) under the env, returning the probe's single trimmed output
// line.
func runWrapped(t *testing.T, safehousePath string, args []string, env []string) string {
	t.Helper()
	cmd := exec.Command(safehousePath, args...)
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("wrapped exec failed: %v\n%s", err, out)
	}
	return strings.TrimSpace(string(out))
}
