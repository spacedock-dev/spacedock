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

// envPassSafehouse writes a fake `safehouse` that SCRUBS SPACEDOCK_BIN by default
// (the env-sanitization the real safehouse does across its boundary) BUT honors
// `--env-pass NAMES` (comma-separated): each named var is forwarded from the host
// env it inherited, on top of the sanitized defaults. This models the real
// `--env-pass` flag the front door now passes. The inner program after `--` runs
// with SPACEDOCK_BIN present only when `--env-pass SPACEDOCK_BIN` was given.
func envPassSafehouse(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "safehouse")
	body := `#!/bin/sh
# Capture the inherited value before scrubbing (the host env safehouse sees).
saved_bin="${SPACEDOCK_BIN:-}"
unset SPACEDOCK_BIN
pass=""
while [ "$#" -gt 0 ] && [ "$1" != "--" ]; do
  if [ "$1" = "--env-pass" ]; then shift; pass="$1"; fi
  shift
done
if [ "$1" = "--" ]; then shift; fi
# Honor --env-pass SPACEDOCK_BIN: forward the named var from the saved host env.
case ",$pass," in *,SPACEDOCK_BIN,*) export SPACEDOCK_BIN="$saved_bin" ;; esac
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
		// The wrap as runClaude composes it: the production launcherBinEnvPassFlags
		// (driven by executablePath → launchedBin) in the safehouse extra slot. The
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
