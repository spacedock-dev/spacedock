//go:build live

// ABOUTME: Live BEHAVIORAL test that a zero-discover FO boot reports-and-stops —
// ABOUTME: no TeamCreate, no broad filesystem sweep (gated, -tags live only).
package ensigncycle

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestLiveZeroDiscoverReportsAndStops is the behavioral half of the lean-boot
// hardening AC-1's offline detector cannot stand in for: a REAL `spacedock claude`
// FO booted against a ZERO-discover fixture (a git-init'd tmpdir with NO
// `commissioned-by: spacedock@` README), driven through the front door exactly like
// TestLiveEnsignCycle. The captain observed (2026-06-14) an FO, after a zero
// `status --discover`, run a broad find/grep filesystem sweep to hunt a workflow
// instead of obeying the contract's terminal zero branch (Startup step 2: zero →
// report no workflow found and STOP). This test proves the real boot:
//
//	(a) reaches its greet/no-workflow report WITHOUT a TeamCreate (no workflow to
//	    drive → no team is engaged), and
//	(b) takes NO broad filesystem sweep — detectBroadSearchAtBoot(transcript, root)
//	    returns nil. On a sweep the detector reds and names the offending command.
//
// The proof is the captured transcript driven through the AC-1 detector — a real
// model boot, observed behavior, not prose. The `//go:build live` tag keeps it out
// of the offline suite; it runs only under `-tags live` with auth present.
func runRetiredClaudeZeroDiscoveryWrapper(t *testing.T) {
	binary := spacedockBinary(t)
	pluginDir := livePluginDir(t)
	model := envOr("SPACEDOCK_LIVE_MODEL", "sonnet")

	childEnv := isolatedClaudeEnv(t, os.Getenv("HOME"))
	childEnv = withBinaryOnPath(childEnv, binary)

	// The zero-discover fixture: a git-init'd root with NO commissioned README.
	// status --discover gates on `commissioned-by: spacedock@` frontmatter
	// (livefixture_discover_test.go), so this root yields zero workflows (empty
	// stdout, exit 0) — the exact condition whose terminal zero branch this guards.
	// The inert .gitkeep gives gitInit's `add -A` something to stage (an empty dir
	// stages nothing, so the init commit would fail "nothing to commit") while
	// keeping the root workflow-less — it carries no commissioned-by frontmatter.
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".gitkeep"), "")
	gitInit(t, root)

	task := "Use $spacedock:first-officer for this whole run."

	// The same front door as TestLiveEnsignCycle: --plugin-dir loads the local v1
	// plugin checkout and relaxes the contract gate; every host flag rides after `--`.
	// drivePrompt keeps the generic anti-early-shutdown override so the boot timing is
	// governed the same way; with no workflow the FO never creates a team.
	drivePrompt := "Drive the workflow. " + antiShutdownOverride
	cmd := exec.Command(binary, "claude",
		"--plugin-dir", pluginDir,
		"--skip-compat-check",
		"--",
		"-p", drivePrompt,
		"--permission-mode", "bypassPermissions",
		"--output-format", "stream-json",
		"--verbose",
		"--model", model,
		task,
	)
	cmd.Dir = root
	cmd.Env = childEnv

	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw
	if err := cmd.Start(); err != nil {
		t.Fatalf("spacedock claude failed to start: %v", err)
	}

	poller := newCmdPoller(cmd, pw)
	watcher := newStreamWatcher(newPipeLineSource(pr), poller, func(line string) { t.Log(line) })
	defer poller.kill()

	// A zero-discover boot has NO watched progress beat (no TeamCreate, no dispatch):
	// the FO greets, reports no workflow found, and ends its turn — the subprocess
	// exits on its own. drainToExit runs it to exit while accumulating the FULL
	// transcript, bounded by the no-progress quiet budget so a silent hang trips fast
	// and localized rather than running to the suite timeout.
	transcript, err := watcher.drainToExit(quietBudgetDefault, "zero-discover boot")
	if err != nil {
		t.Fatalf("zero-discover boot made no stream progress / did not exit: %v", err)
	}

	// (a) NO TeamCreate: with no workflow to drive, the FO must not engage a team.
	for _, line := range strings.Split(transcript, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}
		var e streamEntry
		if json.Unmarshal([]byte(line), &e) != nil {
			continue
		}
		if isTeamCreate(e) {
			t.Fatalf("zero-discover boot engaged a TeamCreate — the FO must report no workflow found and STOP, not create a team\nTranscript:\n%s", transcript)
		}
	}

	// (b) NO broad filesystem sweep: the detector (AC-1) returns nil on a clean
	// report-and-stop boot and reds a sweep, naming the offending command.
	if sweep := detectBroadSearchAtBoot(transcript, root); sweep != nil {
		t.Fatalf("zero-discover boot broad-searched the filesystem instead of report-and-stop: %v\nTranscript:\n%s", sweep, transcript)
	}
}
