// ABOUTME: End-to-end env-scrub smoke: through a safehouse that STRIPS SPACEDOCK_BIN,
// ABOUTME: the argv re-assert lets the real spacedock_launcher resolve the launched binary.
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

// scrubbingSafehouse writes a fake `safehouse` that UNSETS SPACEDOCK_BIN before
// exec'ing the inner argv past `--` — the env-scrubbing behavior the real
// safehouse exhibits across its sandbox boundary. The env-only propagation
// (launchEnv setting SPACEDOCK_BIN on the outer process) cannot survive this, so
// the value must ride the inner argv to reach the host.
func scrubbingSafehouse(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "safehouse")
	body := `#!/bin/sh
unset SPACEDOCK_BIN
while [ "$#" -gt 0 ] && [ "$1" != "--" ]; do shift; done
if [ "$1" = "--" ]; then shift; fi
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

// TestSafehouseScrubPreservedThroughArgvReassert is the AC-5 regression oracle.
// Through a safehouse that STRIPS SPACEDOCK_BIN, the inner argv re-assert
// (`env SPACEDOCK_BIN=<bin> …`, the prefix runClaude/runCodex prepend on the wrap
// path) lets the real spacedock_launcher resolve the LAUNCHED binary. The control
// case — the SAME wrap WITHOUT the argv prefix (env-only propagation) — falls back
// to the PATH `spacedock`, proving the test fails for the right reason absent the
// fix rather than passing trivially.
func TestSafehouseScrubPreservedThroughArgvReassert(t *testing.T) {
	dir := t.TempDir()
	safehousePath := scrubbingSafehouse(t, dir)
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
	// coreutils the probe/wrapper need still resolve (the re-assert itself uses the
	// absolute /usr/bin/env, so it needs no PATH entry).
	base := withoutEnv(withoutEnv(os.Environ(), "PATH"), spacedockBinEnv)
	scrubEnv := append(base,
		"PATH="+pathDir+":/usr/bin:/bin",
		spacedockBinEnv+"="+launchedBin)

	t.Run("argv re-assert survives the scrub", func(t *testing.T) {
		// The inner argv as runClaude composes it on the wrap path: the production
		// launcherBinArgvPrefix (driven by executablePath → launchedBin) ahead of the
		// inner host (here the probe). Exercising the real prefix builder, not a copy.
		withExecutablePath(t, launchedBin, nil)
		inner := append(launcherBinArgvPrefix(), probe)
		argv := safehouse.Wrap(inner, nil)
		out := runWrapped(t, safehousePath, argv[1:], scrubEnv)
		if out != "LAUNCHED" {
			t.Fatalf("probe resolved %q, want LAUNCHED (the launched binary survived the scrub via argv)", out)
		}
	})

	t.Run("without the re-assert it falls back to PATH (regression control)", func(t *testing.T) {
		// The SAME wrap WITHOUT the argv prefix — env-only propagation. The scrub
		// drops SPACEDOCK_BIN, so the probe falls back to the PATH spacedock. This is
		// the bug the argv re-assert fixes; it proves the oracle is not trivial.
		argv := safehouse.Wrap([]string{probe}, nil)
		out := runWrapped(t, safehousePath, argv[1:], scrubEnv)
		if out != "PATH-FALLBACK" {
			t.Fatalf("env-only probe resolved %q, want PATH-FALLBACK (the scrub must drop SPACEDOCK_BIN)", out)
		}
	})
}

// runWrapped execs the fake safehouse with the wrapped argv (everything after the
// `safehouse` token) under the scrubbing env, returning the probe's single
// trimmed output line.
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
