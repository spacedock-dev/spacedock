// ABOUTME: AC-2/AC-3 behavior fixtures for the FO version gate's install
// ABOUTME: journey — drive a shell mirror of the gate flow
// ABOUTME: (testdata/version_gate_flow.sh) with captive curl/install.sh/
// ABOUTME: spacedock commands and assert on exit codes and observable side
// ABOUTME: effects (sentinel, install-run count, SPACEDOCK_BIN repoint),
// ABOUTME: never on byte-for-byte prose (the FO's gate text is LLM prose).
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const curlInstallToken = "curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh"
const withdrawCapability = "spacedock gate withdraw <entity> --reason TEXT"

// gateFixtureDir builds the captive-command fixture tree:
//
//	<dir>/bin/curl            — fake curl: cats the captive install.sh to stdout
//	<dir>/home/.local/bin     — the captive install.sh's install target (off PATH)
//	<dir>/gate-cwd            — the gate's working directory
//
// and returns the dir. installScript is the captive install.sh body; it must
// honor $COUNT_FILE (append one line per run) as its observable side effect.
func gateFixtureDir(t *testing.T, installScript string) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "bin")
	curl := "#!/bin/sh\n# captive curl: ignore flags/URL, emit the captive installer\nexec cat \"$CAPTIVE_INSTALL\"\n"
	writeExe(t, filepath.Join(bin, "curl"), curl)
	captive := filepath.Join(dir, "install.sh")
	writeExe(t, captive, installScript)
	if err := os.MkdirAll(filepath.Join(dir, "gate-cwd"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CAPTIVE_INSTALL", captive)
	return dir
}

func writeExe(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}

// captiveInstall returns an install.sh body that drops a captive `spacedock`
// printing the given --version lines into $HOME/.local/bin, counts its runs in
// $COUNT_FILE, and honors the real install.sh's stderr contract
// (`install.sh: installed spacedock <version> to <dir>/spacedock`).
func captiveInstall(version string, extraVersionLines ...string) string {
	captiveBinary := "#!/bin/sh\n" +
		"[ -n \"${INVOCATION_LOG:-}\" ] && printf '%s|%s\\n' \"$0\" \"$*\" >> \"$INVOCATION_LOG\"\n" +
		"case \"$*\" in\n" +
		"--version) echo \"spacedock " + version + "\"\n"
	for _, l := range extraVersionLines {
		captiveBinary += "echo \"" + l + "\"\n"
	}
	captiveBinary += "  ;;\n" +
		"'gate --help') echo '" + withdrawCapability + "' ;;\n" +
		"'status --boot --identify --json') echo '{\"command\":\"status\",\"boot\":true}' ;;\n" +
		"*) exit 9 ;;\n" +
		"esac\n"
	return "#!/bin/sh\n" +
		"mkdir -p \"$HOME/.local/bin\"\n" +
		"cat > \"$HOME/.local/bin/spacedock\" <<'EOS'\n" + captiveBinary + "EOS\n" +
		"chmod +x \"$HOME/.local/bin/spacedock\"\n" +
		"echo run >> \"$COUNT_FILE\"\n" +
		"echo \"install.sh: installed spacedock " + version + " to $HOME/.local/bin/spacedock\" >&2\n"
}

// runGateFlow executes the gate-flow mirror with the given fixture state and
// returns combined output and exit code. PATH is minimal (captive bin dir plus
// /bin + /usr/bin) so no real `spacedock` can leak in; TMPDIR and HOME are
// fixture-scoped so the sentinel and install target stay in the temp dir.
func runGateFlow(t *testing.T, dir string, extraEnv []string) (string, int) {
	t.Helper()
	script, err := filepath.Abs(filepath.Join("testdata", "version_gate_flow.sh"))
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("sh", script)
	cmd.Dir = filepath.Join(dir, "gate-cwd")
	env := []string{
		"PATH=" + filepath.Join(dir, "bin") + ":/bin:/usr/bin",
		"TMPDIR=" + dir,
		"HOME=" + filepath.Join(dir, "home"),
		"COUNT_FILE=" + filepath.Join(dir, "install-count"),
		"CAPTIVE_INSTALL=" + os.Getenv("CAPTIVE_INSTALL"),
		"INVOCATION_LOG=" + filepath.Join(dir, "invocations"),
		"FIXTURE_OS=Linux",
		"REQUIRED_MINOR=0.27",
	}
	cmd.Env = append(env, extraEnv...)
	out, err := cmd.CombinedOutput()
	code := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		code = exitErr.ExitCode()
	} else if err != nil {
		t.Fatalf("run gate flow: %v", err)
	}
	return string(out), code
}

func installRunCount(t *testing.T, dir string) int {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, "install-count"))
	if os.IsNotExist(err) {
		return 0
	}
	if err != nil {
		t.Fatal(err)
	}
	return len(strings.Fields(string(b)))
}

const fixtureSessionID = "fixture-session-0001"

func sentinelPath(dir string) string {
	return filepath.Join(dir, "spacedock-install-attempted-"+fixtureSessionID)
}

func writeGateLauncher(t *testing.T, dir, version, help string) string {
	t.Helper()
	path := filepath.Join(dir, "selected", "spacedock")
	body := `#!/bin/sh
printf '%s|%s\n' "$0" "$*" >> "$INVOCATION_LOG"
case "$*" in
--version)
	printf '%s\n' 'spacedock ` + version + `'
	;;
'gate --help')
	printf '%s\n' '` + help + `'
	;;
'status --boot --identify --json')
	printf '%s\n' '{"command":"status","boot":true}'
	;;
*)
	printf '%s\n' "unexpected invocation: $*" >&2
	exit 9
	;;
esac
`
	writeExe(t, path, body)
	return path
}

func invocationLog(t *testing.T, dir string) []string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, "invocations"))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimSpace(string(b)), "\n")
}

func TestGateFlowRejectsStaleSameMinorBeforeBoot(t *testing.T) {
	for _, tc := range []struct {
		name   string
		os     string
		remedy string
	}{
		{name: "Darwin", os: "Darwin", remedy: "brew upgrade spacedock"},
		{name: "Linux", os: "Linux", remedy: curlInstallToken},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := gateFixtureDir(t, captiveInstall("0.27.0"))
			launcher := writeGateLauncher(t, dir, "0.27.0-pre2+dev", "stale gate help")
			out, code := runGateFlow(t, dir, []string{"SPACEDOCK_BIN=" + launcher, "FIXTURE_OS=" + tc.os})
			if code != 3 {
				t.Fatalf("exit = %d, want 3 for stale same-minor launcher:\n%s", code, out)
			}
			wantCalls := []string{launcher + "|--version", launcher + "|gate --help"}
			if got := invocationLog(t, dir); strings.Join(got, "\n") != strings.Join(wantCalls, "\n") {
				t.Fatalf("invocations = %#v, want %#v", got, wantCalls)
			}
			for _, want := range []string{launcher, "spacedock 0.27.0-pre2+dev", withdrawCapability, tc.remedy, "relaunch"} {
				if !strings.Contains(out, want) {
					t.Fatalf("stale failure missing %q:\n%s", want, out)
				}
			}
			for _, forbidden := range []string{"go build", "source", "checkout", "repository", "plugin refresh", "SPACEDOCK_BIN repointed"} {
				if strings.Contains(out, forbidden) {
					t.Fatalf("stale failure contains forbidden development/repoint advice %q:\n%s", forbidden, out)
				}
			}
		})
	}
}

func TestGateFlowCompatibleSameMinorProbesThenBootsOnce(t *testing.T) {
	dir := gateFixtureDir(t, captiveInstall("0.27.0"))
	launcher := writeGateLauncher(t, dir, "0.27.0-pre2+dev", "Usage:\n       "+withdrawCapability+" [--workflow-dir DIR]")
	out, code := runGateFlow(t, dir, []string{"SPACEDOCK_BIN=" + launcher})
	if code != 0 {
		t.Fatalf("exit = %d, want 0 for compatible same-minor launcher:\n%s", code, out)
	}
	wantCalls := []string{
		launcher + "|--version",
		launcher + "|gate --help",
		launcher + "|status --boot --identify --json",
	}
	if got := invocationLog(t, dir); strings.Join(got, "\n") != strings.Join(wantCalls, "\n") {
		t.Fatalf("invocations = %#v, want %#v", got, wantCalls)
	}
}

// TestGateFlowInstallAndResumeConvergence is AC-2 fixture 1: a captive install
// places a compatible binary at a known path; assert the install RAN (not just
// printed), the sentinel was created, SPACEDOCK_BIN was repointed, --version was
// re-checked, and the gate passed (exit 0).
func TestGateFlowInstallAndResumeConvergence(t *testing.T) {
	dir := gateFixtureDir(t, captiveInstall("0.27.0"))
	out, code := runGateFlow(t, dir, []string{"CLAUDE_CODE_SESSION_ID=" + fixtureSessionID})
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (convergence):\n%s", code, out)
	}
	if installRunCount(t, dir) != 1 {
		t.Fatalf("install ran %d times, want exactly 1", installRunCount(t, dir))
	}
	if _, err := os.Stat(sentinelPath(dir)); err != nil {
		t.Fatalf("sentinel not created at %s: %v", sentinelPath(dir), err)
	}
	installed := filepath.Join(dir, "home", ".local", "bin", "spacedock")
	if !strings.Contains(out, "SPACEDOCK_BIN repointed to "+installed) {
		t.Fatalf("SPACEDOCK_BIN repoint to the installed path not observed:\n%s", out)
	}
	if !strings.Contains(out, "gate passed after install: spacedock 0.27.0") {
		t.Fatalf("--version re-check convergence not observed:\n%s", out)
	}
}

// TestGateFlowFailureFallbackOneAttempt is AC-2 fixture 2 (B1/NF2): a captive
// install "succeeds" but yields an incompatible binary; assert exactly one
// install run, the sentinel created BEFORE the run, the --version re-check
// observed failing, hint-and-abort carrying the exact command + sentinel path +
// rm recovery, gate NOT passed (exit 3) — and a second gate entry observes the
// sentinel via test -f and runs ZERO further installs.
func TestGateFlowFailureFallbackOneAttempt(t *testing.T) {
	dir := gateFixtureDir(t, captiveInstall("0.20.0"))
	out, code := runGateFlow(t, dir, []string{"CLAUDE_CODE_SESSION_ID=" + fixtureSessionID})
	if code != 3 {
		t.Fatalf("exit = %d, want 3 (gate must not pass on a failed re-check):\n%s", code, out)
	}
	if installRunCount(t, dir) != 1 {
		t.Fatalf("install ran %d times, want exactly 1", installRunCount(t, dir))
	}
	if _, err := os.Stat(sentinelPath(dir)); err != nil {
		t.Fatalf("sentinel not created (create-before-run) at %s: %v", sentinelPath(dir), err)
	}
	if !strings.Contains(out, "--version re-check failed after install (spacedock 0.20.0)") {
		t.Fatalf("failed --version re-check not observed:\n%s", out)
	}
	if !strings.Contains(out, curlInstallToken) {
		t.Fatalf("fallback message must carry the exact OS-aware command:\n%s", out)
	}
	if !strings.Contains(out, "rm "+sentinelPath(dir)) {
		t.Fatalf("fallback message must print the sentinel path and its rm recovery:\n%s", out)
	}

	// Second simulated gate entry: the sentinel suppresses any further install.
	out2, code2 := runGateFlow(t, dir, []string{"CLAUDE_CODE_SESSION_ID=" + fixtureSessionID})
	if code2 != 3 {
		t.Fatalf("second entry exit = %d, want 3:\n%s", code2, out2)
	}
	if installRunCount(t, dir) != 1 {
		t.Fatalf("second entry re-ran the install (count=%d, want 1) — sentinel guardrail absent or ignored", installRunCount(t, dir))
	}
	if !strings.Contains(out2, "sentinel-blocked") || !strings.Contains(out2, "rm "+sentinelPath(dir)) {
		t.Fatalf("sentinel-blocked re-entry must name the sentinel and its rm recovery:\n%s", out2)
	}
}

// TestGateFlowUnsupportedOS covers D-5: a non-Darwin/Linux OS hint-and-aborts
// with the source-build note, exit 3, and never runs the install.
func TestGateFlowUnsupportedOS(t *testing.T) {
	dir := gateFixtureDir(t, captiveInstall("0.27.0"))
	out, code := runGateFlow(t, dir, []string{"FIXTURE_OS=FreeBSD"})
	if code != 3 {
		t.Fatalf("exit = %d, want 3:\n%s", code, out)
	}
	if !strings.Contains(out, "go build -o spacedock ./cmd/spacedock") || !strings.Contains(out, "unsupported OS FreeBSD") {
		t.Fatalf("unsupported-OS hint must name the source build and the OS:\n%s", out)
	}
	if installRunCount(t, dir) != 0 {
		t.Fatalf("install must not run on an unsupported OS")
	}
}

// TestGateFlowSandboxEnvMarker is AC-3: in the REAL install-offer state — a
// captive ABSENT binary producing NO --version output — with the sandbox env
// marker set, detection fires from the env marker alone, the message carries the
// exact command token plus the "outside the sandbox" instruction, and the
// install is not executed. A second arm (binary-present wrong-version class)
// asserts the ^Sandbox: corroboration agrees with — and a third that it never
// overrides — the env verdict.
func TestGateFlowSandboxEnvMarker(t *testing.T) {
	dir := gateFixtureDir(t, captiveInstall("0.27.0"))
	out, code := runGateFlow(t, dir, []string{
		"CLAUDE_CODE_SESSION_ID=" + fixtureSessionID,
		"APP_SANDBOX_CONTAINER_ID=agent-safehouse",
	})
	if code != 3 {
		t.Fatalf("exit = %d, want 3 (gate must not pass inside the sandbox):\n%s", code, out)
	}
	if installRunCount(t, dir) != 0 {
		t.Fatalf("install executed inside the sandbox (count=%d) — the silent no-op AC-3 exists to prevent", installRunCount(t, dir))
	}
	if !strings.Contains(out, curlInstallToken) || !strings.Contains(out, "outside the sandbox") {
		t.Fatalf("sandbox message must carry the exact command + the outside-the-sandbox instruction:\n%s", out)
	}

	// Binary-present wrong-version class, corroboration agrees.
	dir2 := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir2, "gate-cwd"), 0o755); err != nil {
		t.Fatal(err)
	}
	captiveOld := "#!/bin/sh\necho \"spacedock 0.20.0\"\necho \"OS: linux/amd64\"\necho \"Sandbox: inside (agent-safehouse)\"\n"
	writeExe(t, filepath.Join(dir2, "bin", "spacedock"), captiveOld)
	out2, code2 := runGateFlow(t, dir2, []string{"APP_SANDBOX_CONTAINER_ID=agent-safehouse"})
	if code2 != 3 {
		t.Fatalf("wrong-version class exit = %d, want 3:\n%s", code2, out2)
	}
	if !strings.Contains(out2, "corroboration: ^Sandbox: line agrees with the env verdict") {
		t.Fatalf("corroboration must agree with the env verdict:\n%s", out2)
	}
	if !strings.Contains(out2, "doctor") {
		t.Fatalf("wrong-version class must keep the doctor remedy pointer:\n%s", out2)
	}

	// Disagreement: env marker set, ^Sandbox: absent-or-contradicting — env wins.
	dir3 := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir3, "gate-cwd"), 0o755); err != nil {
		t.Fatal(err)
	}
	captiveOutside := "#!/bin/sh\necho \"spacedock 0.20.0\"\necho \"OS: linux/amd64\"\necho \"Sandbox: not sandboxed (safehouse available)\"\n"
	writeExe(t, filepath.Join(dir3, "bin", "spacedock"), captiveOutside)
	out3, _ := runGateFlow(t, dir3, []string{"APP_SANDBOX_CONTAINER_ID=agent-safehouse"})
	if !strings.Contains(out3, "DISAGREEMENT") || !strings.Contains(out3, "env check wins") {
		t.Fatalf("env verdict must win over a contradicting ^Sandbox: line, named in the message:\n%s", out3)
	}
}
