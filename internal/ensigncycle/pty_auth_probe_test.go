//go:build live

// ABOUTME: A cheap short-budget probe that dumps the painted interactive-claude pane (escape codes
// ABOUTME: preserved) so a human can read WHICH pre-transcript screen blocks the Linux-CI pty harness.
package ensigncycle

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// ptyAuthProbeBudget is the SHORT wall-clock budget the probe waits for a session
// jsonl before it dumps the live pane. It is deliberately a fraction of the harness'
// 4-minute ptyBootBudget: the blocking screen paints in the first seconds, so a ~25s
// wait is more than enough to confirm the jsonl never appears while keeping a failing
// probe run to seconds, not minutes (AC-5 of the parent diagnosis).
const ptyAuthProbeBudget = 25 * time.Second

// TestPtyAuthProbe settles the auth-vs-onboarding fork the 4-minute harness timeout
// cannot: it launches interactive `spacedock claude` in tmux with EXACTLY the CI
// child env shape (env-only credential via isolatedClaudeEnv, fresh per-run
// CLAUDE_CONFIG_DIR, seedInteractiveClaudeConfig applied, and NO stored-login seed —
// the same shape Linux CI carries because seedStoredLoginCredential is a macOS-only
// no-op there), waits a SHORT budget for the FO session jsonl, and on the expected
// timeout IMMEDIATELY dumps the painted pane WITH escape codes preserved plus the
// pane's pid/command/title and the tmux session liveness.
//
// A human reads probe-pane.ansi to discriminate the blocking screen:
//   - a "Claude API" / "Invalid API key" / login banner => AUTH rejects the env
//     credential (cause a);
//   - a theme / trust-folder / login-method picker => an ONBOARDING gate the
//     .claude.json seed misses (cause c);
//   - a genuinely blank pane while pane-info.txt shows node/claude still running =>
//     render-blank or stdin-wait (investigate further).
//
// It is gated OFF in the normal suite by SPACEDOCK_PTY_AUTH_PROBE=1 (it spends a
// model and needs tmux), and inherits the AC-6 skip-not-fatal-without-auth gate from
// newPtyLiveDriver -> isolatedClaudeEnv.
func TestPtyAuthProbe(t *testing.T) {
	if os.Getenv("SPACEDOCK_PTY_AUTH_PROBE") != "1" {
		t.Skip("set SPACEDOCK_PTY_AUTH_PROBE=1 to run the interactive-pty auth probe (spends a model + needs tmux)")
	}

	driver := newPtyLiveDriver(t)

	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeRealisticLifecycle())
	writeFile(t, filepath.Join(root, "make-it-work.md"), entityFixture())
	gitInit(t, root)

	artifactDir := filepath.Join(driver.artifactRoot, "auth-probe")
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Per-run CLAUDE_CONFIG_DIR nested under the driver's base, mirroring
	// launchAndSend's isolation. A COPY of the env, never a mutation of the shared
	// driver env.
	env := driver.env
	configDir, _ := envValue(driver.env, "CLAUDE_CONFIG_DIR")
	if configDir != "" {
		configDir = filepath.Join(configDir, "auth-probe")
		env = withClaudeConfigDir(driver.env, configDir)
	}
	effectiveConfigDir := configDirOrDefault(configDir, driver.homeDir)

	// The same symlink-resolved cwd the seed + the projects-dir name key on.
	resolvedCwd := root
	if r, err := filepath.EvalSymlinks(root); err == nil {
		resolvedCwd = r
	}

	// Seed ONLY the onboarding-skip .claude.json (exactly what the harness does), and
	// DELIBERATELY skip seedStoredLoginCredential so the child carries only the env
	// credential — the precise Linux-CI shape this probe exists to reproduce and read.
	if err := seedInteractiveClaudeConfig(effectiveConfigDir, resolvedCwd); err != nil {
		t.Fatalf("seed interactive claude config: %v", err)
	}

	session := fmt.Sprintf("sdpty-auth-probe-%d", time.Now().UnixNano())
	proc := newTmuxLiveProc(session)
	defer proc.kill()

	// The launch line is the harness' exact interactive front door (the `-p` front
	// door minus `-p`): spacedock-owned flags before `--`, host flags after. The
	// nested-session markers are unset on the launch command itself.
	launch := shellJoin(append([]string{"env"}, unsetNestedSessionArgs(driver.binary, "claude",
		"--plugin-dir", driver.pluginDir,
		"--skip-contract-check",
		"--",
		"--permission-mode", "bypassPermissions",
		"--model", driver.modelName,
	)...))
	// The child env rides on per-session `-e KEY=VAL` flags (a pre-existing tmux
	// server would otherwise drop the command process env), exactly as launchAndSend.
	args := []string{"new-session", "-d", "-s", session, "-x", "220", "-y", "50", "-c", root}
	for _, kv := range env {
		args = append(args, "-e", kv)
	}
	args = append(args, launch)
	if out, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		t.Fatalf("tmux new-session for the auth probe failed: %v\n%s", err, out)
	}

	projectsDir := filepath.Join(effectiveConfigDir, "projects", encodeProjectDir(resolvedCwd))
	t.Logf("[pty auth-probe] polling for session jsonl under: %s (budget %s)", projectsDir, ptyAuthProbeBudget)

	sessionFile, waitErr := waitForSessionFile(projectsDir, proc, ptyAuthProbeBudget)
	if waitErr == nil {
		// REFUTATION outcome: the interactive child DID write a transcript under the
		// env-only credential. Capture the same diagnostics for the record (the human
		// reads whether it is a real authenticated boot), but this is the path that
		// would overturn the auth/onboarding fork.
		t.Logf("[pty auth-probe] UNEXPECTED: a session jsonl appeared under the env-only credential: %s", sessionFile)
	} else {
		t.Logf("[pty auth-probe] no session jsonl within budget (%v) — dumping the live pane", waitErr)
	}

	// Whether the jsonl appeared or not, dump the live pane WITH escape codes, the
	// pane info, and the session liveness — this IS the diagnostic payload.
	dumpProbeArtifacts(t, session, projectsDir, proc, artifactDir, sessionFile, waitErr)
}

// dumpProbeArtifacts writes the discriminating evidence to artifactDir:
//   - probe-pane.ansi: `tmux capture-pane -p -e` (escape codes preserved) of the
//     active pane — the painted blocking screen a human reads;
//   - probe-pane-fo.txt: the title-resolved FO pane (plain), for when the FO pane
//     titled itself (it usually has NOT, which is itself signal);
//   - probe-pane-info.txt: `tmux list-panes` pid/command/title — distinguishes a
//     live node/claude waiting on stdin from an exited child;
//   - probe-state.txt: a one-line synthesis (session alive?, jsonl seen?, the wait
//     error) so the artifact is self-describing without the CI log.
//
// It also logs the captures so the evidence survives even if the artifact upload is
// skipped. It NEVER fatals: the probe's job is to emit evidence, not to pass/fail.
func dumpProbeArtifacts(t *testing.T, session, projectsDir string, proc *tmuxLiveProc, artifactDir, sessionFile string, waitErr error) {
	t.Helper()

	// Escape-code-preserving capture of the ACTIVE pane (`-e`): the harness'
	// diagnostic captureFOPane uses `-p` plain AND title-resolves the FO pane, which
	// returns "" when no FO-titled pane exists — masking a painted auth/login banner.
	// The probe captures the active pane raw so a pre-FO blocking screen is visible.
	paneANSI, paneErr := exec.Command("tmux", "capture-pane", "-p", "-e", "-t", session).CombinedOutput()
	if paneErr != nil {
		paneANSI = []byte(fmt.Sprintf("capture-pane -e failed: %v\n%s", paneErr, paneANSI))
	}
	writeProbeArtifact(t, artifactDir, "probe-pane.ansi", paneANSI)

	// The plain, title-resolved FO pane (the harness' own diagnostic view) for
	// comparison — empty here is the harness' observed 0-byte signature.
	foPane, foErr := exec.Command("tmux", "capture-pane", "-p", "-t", session).CombinedOutput()
	if foErr != nil {
		foPane = []byte(fmt.Sprintf("capture-pane failed: %v\n%s", foErr, foPane))
	}
	writeProbeArtifact(t, artifactDir, "probe-pane-fo.txt", foPane)

	// pane pid/command/title: a live `node`/`claude` waiting on stdin vs an exited
	// child is the auth-banner (process blocked) vs immediate-crash discriminator.
	paneInfo, infoErr := exec.Command("tmux", "list-panes", "-t", session,
		"-F", "#{pane_pid} #{pane_current_command} #{pane_title}").CombinedOutput()
	if infoErr != nil {
		paneInfo = []byte(fmt.Sprintf("list-panes failed: %v\n%s", infoErr, paneInfo))
	}
	writeProbeArtifact(t, artifactDir, "probe-pane-info.txt", paneInfo)

	_, sessionExited := proc.poll()
	var state strings.Builder
	fmt.Fprintf(&state, "tmux session: %s\n", session)
	fmt.Fprintf(&state, "session alive at dump: %t\n", !sessionExited)
	fmt.Fprintf(&state, "session jsonl seen: %t (%q)\n", waitErr == nil, sessionFile)
	fmt.Fprintf(&state, "wait error: %v\n", waitErr)
	fmt.Fprintf(&state, "polled projects dir: %s\n", projectsDir)
	writeProbeArtifact(t, artifactDir, "probe-state.txt", []byte(state.String()))

	t.Logf("[pty auth-probe] state:\n%s", state.String())
	t.Logf("[pty auth-probe] pane info: %s", strings.TrimSpace(string(paneInfo)))
	t.Logf("[pty auth-probe] painted pane (escape codes preserved):\n%s", string(paneANSI))
	t.Logf("[pty auth-probe] artifacts: %s", artifactDir)
}

// writeProbeArtifact writes one probe artifact, logging (never fataling) on a write
// error so a single failed write does not lose the rest of the evidence.
func writeProbeArtifact(t *testing.T, dir, name string, data []byte) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
		t.Logf("[pty auth-probe] write %s: %v", name, err)
	}
}
