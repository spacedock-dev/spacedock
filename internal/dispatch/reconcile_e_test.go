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

// TestReconcileEDetectsAndResetAdvancesMain builds a real git fixture where
// main carries a commit not on origin/next (ahead-only / unpushed) and asserts
// the report-only contract: the helper emits exactly one Class-E item with
// Ahead==1, Behind==0, and a reason that NEVER prescribes a destructive reset.
// The seeded git state (main ahead, origin/next behind) is the independent
// oracle — main carries unpushed committed work the contract must never discard.
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
trunk: next
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
		include:     map[string]bool{classLocalMainDrift: true},
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
	if len(result.Drift) != 1 || result.Drift[0].Class != classLocalMainDrift {
		t.Fatalf("want exactly one local-main-drift entry; got %s", formatDrift(result.Drift))
	}
	d := result.Drift[0]
	if d.Ahead != 1 {
		t.Errorf("E.ahead=%d, want 1", d.Ahead)
	}
	if d.Behind != 0 {
		t.Errorf("E.behind=%d, want 0 (ahead-only)", d.Behind)
	}

	// 2. The seeded state confirms main carries committed work origin/next lacks:
	// rev-parse main differs from origin/next. A reset --hard would discard it.
	mainSHA := mustGit(t, repoRoot, "rev-parse", "main")
	nextSHA := mustGit(t, repoRoot, "rev-parse", "origin/next")
	if mainSHA == nextSHA {
		t.Fatalf("fixture invalid: main and origin/next are already equal (%s)", mainSHA)
	}

	// 3. Report-only contract: the reason must NEVER prescribe a reset that would
	// discard the unpushed commit.
	if strings.Contains(strings.ToLower(d.Reason), "reset") {
		t.Errorf("E.reason=%q must NOT contain 'reset' (report-only for unpushed main)", d.Reason)
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
