// ABOUTME: AC-1/AC-3 — Class-D/E resolve the integration trunk from the README's
// ABOUTME: top-level trunk: key (sentinel ftrunk), never a hardcoded next/main.
package dispatch

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
)

// TestReconcileTrunkFromConfig (AC-1 + AC-3) builds a fixture whose README
// declares top-level `trunk: ftrunk` — a sentinel that is neither `next` nor
// `main` — and a git graph with an `origin/ftrunk` trunk ref. It asserts:
//   - Class-D fires when a worktree branch is behind origin/ftrunk;
//   - Class-E fires when local main carries commits not on origin/ftrunk;
//   - the emitted Class-E driftItem carries trunk == "ftrunk" (AC-3).
//
// The expected ref name comes from the fixture's declared trunk, NOT from any
// literal in reconcile.go: a regression that re-hardcodes `next` (or `main`)
// would query origin/next, find no such ref, and fail to detect against
// origin/ftrunk — reds the test against an independent oracle.
func TestReconcileTrunkFromConfig(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	const trunk = "ftrunk"
	f := newTrunkFixture(t, trunk)
	result := f.run()

	byClass := groupDriftByClass(result.Drift)

	// stale-branch: the stale worktree is behind origin/ftrunk.
	ds := byClass[classStaleBranch]
	if len(ds) != 1 {
		t.Fatalf("expected 1 stale-branch entry against origin/%s; got %d: %s",
			trunk, len(ds), formatDrift(byClass[classStaleBranch]))
	}
	if ds[0].Slug != "yaml-parser-migration" {
		t.Errorf("stale-branch.slug=%q, want yaml-parser-migration", ds[0].Slug)
	}
	if ds[0].Behind <= 0 {
		t.Errorf("stale-branch.behind=%d, want > 0 (behind origin/%s)", ds[0].Behind, trunk)
	}

	// local-main-drift: local main carries a commit not on origin/ftrunk, and the
	// drift item carries the resolved trunk (AC-3) so the FO remedy is data-driven.
	es := byClass[classLocalMainDrift]
	if len(es) != 1 {
		t.Fatalf("expected 1 local-main-drift entry against origin/%s; got %d: %s",
			trunk, len(es), formatDrift(byClass[classLocalMainDrift]))
	}
	if es[0].Ahead <= 0 {
		t.Errorf("local-main-drift.ahead=%d, want > 0", es[0].Ahead)
	}
	if es[0].Trunk != trunk {
		t.Errorf("local-main-drift.trunk=%q, want %q (the resolved trunk from the fixture README — AC-3)",
			es[0].Trunk, trunk)
	}
}

// trunkFixture is a minimal reconcile fixture parameterized by the integration
// trunk name. It builds origin/{trunk} as the trunk ref, one worktree behind it
// (class D), and a local main one commit ahead of it (class E). Only the
// git/filesystem classes D/E are exercised — no team roster is needed.
type trunkFixture struct {
	t           *testing.T
	home        string
	repoRoot    string
	workflowDir string
	stateRoot   string
	teamName    string
}

func newTrunkFixture(t *testing.T, trunk string) *trunkFixture {
	t.Helper()
	home := t.TempDir()
	repoRoot := t.TempDir()
	teamName := "team-trunk-fixture"

	// Workflow README declares the sentinel trunk as a top-level key.
	workflowDir := filepath.Join(repoRoot, "docs", "wf")
	writeFile(t, filepath.Join(workflowDir, "README.md"), readmeWithTrunk(trunk))
	stateRoot := filepath.Join(workflowDir, ".spacedock-state")
	if err := os.MkdirAll(stateRoot, 0o755); err != nil {
		t.Fatal(err)
	}

	// One entity whose worktree branch will lag origin/{trunk} → class D.
	wtStale := filepath.Join(repoRoot, ".worktrees", "spacedock-ensign-yaml-parser-migration")
	writeFile(t, filepath.Join(stateRoot, "yaml-parser-migration", "index.md"),
		reconcileEntityFM(map[string]string{
			"id":       "id-ypm",
			"title":    "yaml parser migration",
			"slug":     "yaml-parser-migration",
			"status":   "implementation",
			"worktree": filepath.Join(".worktrees", "spacedock-ensign-yaml-parser-migration"),
		}))

	// Team config so the sweep resolves a team identity (rosterTrusted) and runs
	// D/E. The single ensign matches the entity slug.
	cfgPath := filepath.Join(home, ".claude", "teams", teamName, "config.json")
	writeFile(t, cfgPath, teamConfigJSON(teamName, []claudeteam.ReconcileMember{
		{Name: "team-lead", AgentType: "team-lead", Model: "claude-opus-4-7"},
		{Name: "spacedock-ensign-yaml-parser-migration-implementation",
			AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
	}))

	// Git graph mirroring newReconcileFixture but with origin/{trunk} as the
	// trunk ref instead of origin/next:
	//   initial (seed) ── origin/{trunk} @ initial; wtStale @ initial
	//   main (+ second-commit) — 1 ahead of origin/{trunk} (class E)
	//   origin/{trunk} bumped (+{trunk}-bump) — wtStale now 1 behind (class D)
	repoGitInit(t, repoRoot)
	repoSetOriginTrunk(t, repoRoot, trunk, "HEAD")
	makeFixtureWorktree(t, repoRoot, wtStale, "ypm-branch", "refs/remotes/origin/"+trunk)
	repoMakeCommit(t, repoRoot, "main", "second-commit.txt", "two\n")
	repoBumpOriginTrunk(t, repoRoot, trunk)

	return &trunkFixture{
		t:           t,
		home:        home,
		repoRoot:    repoRoot,
		workflowDir: workflowDir,
		stateRoot:   stateRoot,
		teamName:    teamName,
	}
}

func (f *trunkFixture) run() reconcileResult {
	f.t.Helper()
	var stdout, stderr bytes.Buffer
	opts := reconcileOpts{
		workflowDir: f.workflowDir,
		teamName:    f.teamName,
		repoRoot:    f.repoRoot,
		include:     map[string]bool{classLingering: true, classSuperseded: true, classUnadvancedPR: true, classStaleBranch: true, classLocalMainDrift: true},
		home:        f.home,
		roster:      claudeteam.LoadReconcileTeam,
		gh:          func(string) (string, error) { return "", nil },
		git:         gitRunnerExec,
	}
	code := Reconcile(opts, &stdout, &stderr)
	if code != 0 {
		f.t.Fatalf("Reconcile exit=%d stderr=%s", code, stderr.String())
	}
	var result reconcileResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		f.t.Fatalf("decode result: %v\nstdout=%s", err, stdout.String())
	}
	return result
}

// repoSetOriginTrunk points refs/remotes/origin/{trunk} at commitish — the
// trunk-name-parameterized analogue of repoSetOriginNext.
func repoSetOriginTrunk(t *testing.T, dir, trunk, commitish string) {
	t.Helper()
	mustGit(t, dir, "update-ref", "refs/remotes/origin/"+trunk, commitish)
}

// repoBumpOriginTrunk advances origin/{trunk} by one commit (on a detached
// HEAD off the current trunk ref), so a worktree created before the bump lags
// by 1. Returns to main afterward — the trunk-name-parameterized analogue of
// repoBumpOriginNext.
func repoBumpOriginTrunk(t *testing.T, dir, trunk string) {
	t.Helper()
	mustGit(t, dir, "checkout", "-q", "refs/remotes/origin/"+trunk)
	writeFile(t, filepath.Join(dir, trunk+"-bump.txt"), "bump\n")
	mustGit(t, dir, "add", "-A")
	mustGit(t, dir, "commit", "-q", "-m", "trunk bump")
	mustGit(t, dir, "update-ref", "refs/remotes/origin/"+trunk, "HEAD")
	mustGit(t, dir, "checkout", "-q", "main")
}
