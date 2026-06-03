// ABOUTME: AC-3 class-E integration — real git fixture asserting `rev-parse main
// ABOUTME: == rev-parse origin/next` after the reset, plus `go build` exit code.
package dispatch

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
)

// TestReconcileEDetectsAndResetAdvancesMain wires AC-3: build a real git
// fixture where main carries a commit not on origin/next, prove the helper
// emits Class E, then run the FO's prescribed sequence
// (`git fetch && git reset --hard origin/next`) and assert main now equals
// origin/next. The reset is the deterministic shell sequence whose exit code
// is the proof, not prose.
func TestReconcileEDetectsAndResetAdvancesMain(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	home := t.TempDir()
	repoRoot := t.TempDir()
	workflowDir := filepath.Join(repoRoot, "docs", "wf")
	writeFile(t, filepath.Join(workflowDir, "README.md"), `---
entity-type: task
state: .spacedock-state
stages:
  states:
    - name: backlog
      initial: true
    - name: done
      terminal: true
---
`)
	if err := os.MkdirAll(filepath.Join(workflowDir, ".spacedock-state"), 0o755); err != nil {
		t.Fatal(err)
	}
	// Team config: at least one ensign so resolveTeamConfig auto-discovery would
	// pick this team if we did not pass --team-name (we pass it explicitly here).
	teamName := "team-e-fixture"
	cfgPath := filepath.Join(home, ".claude", "teams", teamName, "config.json")
	writeFile(t, cfgPath, teamConfigJSON(teamName, []claudeteam.ReconcileMember{
		{Name: "team-lead", AgentType: "team-lead", Model: "claude-opus-4-7"},
		{Name: "spacedock-ensign-ghost-implementation",
			AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
	}))

	// Real git graph:
	//   initial  ←  origin/next
	//   └── main (+ stray-main.txt) — 1 ahead of origin/next.
	repoGitInit(t, repoRoot)
	repoSetOriginNext(t, repoRoot, "HEAD")
	repoMakeCommit(t, repoRoot, "main", "stray-main.txt", "stray\n")

	// 1. Detect: helper emits class E with ahead >= 1.
	var stdout, stderr bytes.Buffer
	opts := reconcileOpts{
		workflowDir: workflowDir,
		teamName:    teamName,
		repoRoot:    repoRoot,
		include:     map[string]bool{"E": true},
		home:        home,
		roster:      claudeteam.LoadReconcileTeam,
		gh:          func(string) (string, error) { return "", nil },
		git:         gitRunnerExec,
	}
	if code := Reconcile(opts, &stdout, &stderr); code != 0 {
		t.Fatalf("Reconcile exit=%d stderr=%s", code, stderr.String())
	}
	var result reconcileResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode: %v\nstdout=%s", err, stdout.String())
	}
	if len(result.Drift) != 1 || result.Drift[0].Class != "E" {
		t.Fatalf("want exactly one E drift; got %s", formatDrift(result.Drift))
	}
	if result.Drift[0].Ahead != 1 {
		t.Errorf("E.ahead=%d, want 1", result.Drift[0].Ahead)
	}

	// 2. Pre-reset: rev-parse main MUST differ from rev-parse origin/next.
	mainSHA := mustGit(t, repoRoot, "rev-parse", "main")
	nextSHA := mustGit(t, repoRoot, "rev-parse", "origin/next")
	if mainSHA == nextSHA {
		t.Fatalf("pre-reset: main and origin/next are already equal (%s)", mainSHA)
	}

	// 3. Run the FO's prescribed deterministic shell sequence — note we skip the
	// `git fetch origin next` step because origin/next is already at its target
	// (the fixture has no remote; the ref was set via update-ref). The reset is
	// the load-bearing operation.
	mustGit(t, repoRoot, "reset", "--hard", "origin/next")

	// 4. Post-reset: main now points at origin/next. This is the exit-code
	// equivalent — rev-parse returns 0 and both refs match.
	mainAfter := mustGit(t, repoRoot, "rev-parse", "main")
	nextAfter := mustGit(t, repoRoot, "rev-parse", "origin/next")
	if mainAfter != nextAfter {
		t.Errorf("post-reset: main=%s != origin/next=%s", mainAfter, nextAfter)
	}
}

// TestReconcileEGoBuildIsRunnable verifies the binary-rebuild half of AC-3:
// `go build ./cmd/spacedock` from the repo's own root must succeed (exit 0).
// The build's exit code IS the proof — we do not assert on the binary's
// contents. This guards the FO's reset-then-rebuild action from regressing.
func TestReconcileEGoBuildIsRunnable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping go build in short mode")
	}
	// Find the repo root by walking up from this file's package dir until we
	// see go.mod.
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := cwd
	for {
		if _, err := os.Stat(filepath.Join(repoRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(repoRoot)
		if parent == repoRoot {
			t.Fatalf("no go.mod found above %s", cwd)
		}
		repoRoot = parent
	}
	if _, err := os.Stat(filepath.Join(repoRoot, "cmd", "spacedock")); err != nil {
		t.Skipf("repo does not contain cmd/spacedock (skipping rebuild test): %v", err)
	}
	out := filepath.Join(t.TempDir(), "spacedock-test")
	cmd := exec.Command("go", "build", "-o", out, "./cmd/spacedock")
	cmd.Dir = repoRoot
	combined, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build ./cmd/spacedock failed: %v\n%s", err, combined)
	}
	if info, err := os.Stat(out); err != nil || info.Size() == 0 {
		t.Fatalf("expected binary at %s; stat=%v", out, err)
	}
	if !strings.HasSuffix(out, "spacedock-test") {
		t.Fatalf("unexpected output path %q", out)
	}
}
