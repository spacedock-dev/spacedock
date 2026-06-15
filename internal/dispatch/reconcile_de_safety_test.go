// ABOUTME: AC-1/AC-2 — Class-E direction table (never reset) and Class-D
// ABOUTME: ownership gating (rebase only an owned worktree) over seeded git fixtures.
package dispatch

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
)

// runClassE drives Reconcile over a seeded repo with only Class E included and
// returns the single E drift item (or fails). The repo is seeded by the caller
// so the git state is the independent oracle.
func runClassE(t *testing.T, home, repoRoot, workflowDir, teamName string) driftItem {
	t.Helper()
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
	return result.Drift[0]
}

// newClassERepo builds the shared scaffolding (workflow README with trunk: next,
// a state checkout, a team config with one ensign) and a git repo seeded so that
// main is AHEAD of origin/next by 1 (spike-1's ahead-only graph). The caller may
// further advance origin/next (behind) or main (diverged) before running.
func newClassERepo(t *testing.T) (home, repoRoot, workflowDir, teamName string) {
	t.Helper()
	home = t.TempDir()
	repoRoot = t.TempDir()
	teamName = "team-de-fixture"
	workflowDir = filepath.Join(repoRoot, "docs", "wf")
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
	cfgPath := filepath.Join(home, ".claude", "teams", teamName, "config.json")
	writeFile(t, cfgPath, teamConfigJSON(teamName, []claudeteam.ReconcileMember{
		{Name: "team-lead", AgentType: "team-lead", Model: "claude-opus-4-7"},
		{Name: "spacedock-ensign-ghost-implementation",
			AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
	}))
	if err := os.MkdirAll(filepath.Join(workflowDir, ".spacedock-state"), 0o755); err != nil {
		t.Fatal(err)
	}

	// ahead-only graph: initial ← origin/next; main +1 commit ahead.
	repoGitInit(t, repoRoot)
	repoSetOriginNext(t, repoRoot, "HEAD")
	repoMakeCommit(t, repoRoot, "main", "stray-main.txt", "stray\n")
	return home, repoRoot, workflowDir, teamName
}

// TestReconcileEAheadOnlyReportsNoReset (AC-1) — main carries unpushed commits
// (ahead>0, behind 0). The drift item must be report-only: Ahead==1, Behind==0,
// and reason contains no `reset` substring (case-insensitive).
func TestReconcileEAheadOnlyReportsNoReset(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	home, repoRoot, workflowDir, teamName := newClassERepo(t)
	d := runClassE(t, home, repoRoot, workflowDir, teamName)
	if d.Ahead != 1 {
		t.Errorf("ahead-only E.ahead=%d, want 1", d.Ahead)
	}
	if d.Behind != 0 {
		t.Errorf("ahead-only E.behind=%d, want 0", d.Behind)
	}
	if strings.Contains(strings.ToLower(d.Reason), "reset") {
		t.Errorf("ahead-only E.reason=%q must NOT contain 'reset'", d.Reason)
	}
}

// TestReconcileEBehindOnlyPrescribesFFMerge (AC-1) — origin/next is ahead of
// main (behind>0, ahead 0). The remedy is an ff-merge; reason carries no reset.
func TestReconcileEBehindOnlyPrescribesFFMerge(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	home, repoRoot, workflowDir, teamName := newClassERepo(t)
	// Reset main back to origin/next, then advance origin/next so main is behind.
	mustGit(t, repoRoot, "reset", "--hard", "origin/next")
	repoBumpOriginNext(t, repoRoot)
	d := runClassE(t, home, repoRoot, workflowDir, teamName)
	if d.Behind != 1 {
		t.Errorf("behind-only E.behind=%d, want 1", d.Behind)
	}
	if d.Ahead != 0 {
		t.Errorf("behind-only E.ahead=%d, want 0", d.Ahead)
	}
	if !strings.Contains(strings.ToLower(d.Reason), "ff-merge") {
		t.Errorf("behind-only E.reason=%q, want 'ff-merge'", d.Reason)
	}
	if strings.Contains(strings.ToLower(d.Reason), "reset") {
		t.Errorf("behind-only E.reason=%q must NOT contain 'reset'", d.Reason)
	}
}

// TestReconcileEDivergedReportsNoReset (AC-1) — main is both ahead and behind
// origin/next. Report-only: Ahead>0, Behind>0, no reset.
func TestReconcileEDivergedReportsNoReset(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	home, repoRoot, workflowDir, teamName := newClassERepo(t)
	// main already +1 ahead; advance origin/next so main is also behind → diverged.
	repoBumpOriginNext(t, repoRoot)
	d := runClassE(t, home, repoRoot, workflowDir, teamName)
	if d.Ahead <= 0 {
		t.Errorf("diverged E.ahead=%d, want > 0", d.Ahead)
	}
	if d.Behind <= 0 {
		t.Errorf("diverged E.behind=%d, want > 0", d.Behind)
	}
	if strings.Contains(strings.ToLower(d.Reason), "reset") {
		t.Errorf("diverged E.reason=%q must NOT contain 'reset'", d.Reason)
	}
}

// deOwnershipFixture seeds a repo with two behind worktrees: one whose slug an
// in-roster ensign decomposes to (owned), one whose slug no member resolves to
// (un-owned). teamName selects trusted vs untrusted roster.
type deOwnershipFixture struct {
	home        string
	repoRoot    string
	workflowDir string
	teamName    string
}

func newDEOwnershipFixture(t *testing.T) *deOwnershipFixture {
	t.Helper()
	home := t.TempDir()
	repoRoot := t.TempDir()
	teamName := "team-de-own-fixture"

	workflowDir := filepath.Join(repoRoot, "docs", "wf")
	writeFile(t, filepath.Join(workflowDir, "README.md"), `---
entity-type: task
state: .spacedock-state
trunk: next
stages:
  states:
    - name: backlog
      initial: true
    - name: implementation
      worktree: true
    - name: done
      terminal: true
---
`)
	stateRoot := filepath.Join(workflowDir, ".spacedock-state")
	if err := os.MkdirAll(stateRoot, 0o755); err != nil {
		t.Fatal(err)
	}

	// Two entities, each with a worktree that will lag origin/next.
	wtOwned := filepath.Join(repoRoot, ".worktrees", "spacedock-ensign-owned-slug")
	writeFile(t, filepath.Join(stateRoot, "owned-slug", "index.md"),
		reconcileEntityFM(map[string]string{
			"id":       "id-owned",
			"title":    "owned slug",
			"slug":     "owned-slug",
			"status":   "implementation",
			"worktree": filepath.Join(".worktrees", "spacedock-ensign-owned-slug"),
		}))
	wtUnowned := filepath.Join(repoRoot, ".worktrees", "spacedock-ensign-unowned-slug")
	writeFile(t, filepath.Join(stateRoot, "unowned-slug", "index.md"),
		reconcileEntityFM(map[string]string{
			"id":       "id-unowned",
			"title":    "unowned slug",
			"slug":     "unowned-slug",
			"status":   "implementation",
			"worktree": filepath.Join(".worktrees", "spacedock-ensign-unowned-slug"),
		}))

	// Roster: an ensign for owned-slug only; nothing resolves to unowned-slug.
	cfgPath := filepath.Join(home, ".claude", "teams", teamName, "config.json")
	writeFile(t, cfgPath, teamConfigJSON(teamName, []claudeteam.ReconcileMember{
		{Name: "team-lead", AgentType: "team-lead", Model: "claude-opus-4-7"},
		{Name: "spacedock-ensign-owned-slug-implementation",
			AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
	}))

	// Git graph: both worktrees created off origin/next at initial, then
	// origin/next bumped so both lag by 1.
	repoGitInit(t, repoRoot)
	repoSetOriginNext(t, repoRoot, "HEAD")
	makeFixtureWorktree(t, repoRoot, wtOwned, "owned-branch", "origin/next")
	makeFixtureWorktree(t, repoRoot, wtUnowned, "unowned-branch", "origin/next")
	repoBumpOriginNext(t, repoRoot)

	return &deOwnershipFixture{
		home:        home,
		repoRoot:    repoRoot,
		workflowDir: workflowDir,
		teamName:    teamName,
	}
}

func (f *deOwnershipFixture) runD(t *testing.T, teamName string) map[string]driftItem {
	t.Helper()
	var stdout, stderr bytes.Buffer
	opts := reconcileOpts{
		workflowDir: f.workflowDir,
		teamName:    teamName,
		repoRoot:    f.repoRoot,
		include:     map[string]bool{classStaleBranch: true},
		home:        f.home,
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
	bySlug := map[string]driftItem{}
	for _, d := range result.Drift {
		if d.Class == classStaleBranch {
			bySlug[d.Slug] = d
		}
	}
	return bySlug
}

// TestReconcileDOwnedRebasesUnownedReports (AC-2) — with the trusted roster, the
// owned-slug worktree gets Owned:true + a rebase remedy; the unowned-slug
// worktree gets Owned:false + a report-only reason with no pull/rebase verb.
func TestReconcileDOwnedRebasesUnownedReports(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	f := newDEOwnershipFixture(t)
	bySlug := f.runD(t, f.teamName)

	owned, ok := bySlug["owned-slug"]
	if !ok {
		t.Fatalf("missing D entry for owned-slug; got %v", bySlug)
	}
	if !owned.Owned {
		t.Errorf("owned-slug D.owned=%v, want true", owned.Owned)
	}
	if !strings.Contains(strings.ToLower(owned.Reason), "rebase") {
		t.Errorf("owned-slug D.reason=%q, want a rebase remedy", owned.Reason)
	}

	unowned, ok := bySlug["unowned-slug"]
	if !ok {
		t.Fatalf("missing D entry for unowned-slug; got %v", bySlug)
	}
	if unowned.Owned {
		t.Errorf("unowned-slug D.owned=%v, want false", unowned.Owned)
	}
	r := strings.ToLower(unowned.Reason)
	if strings.Contains(r, "pull") || strings.Contains(r, "rebase") {
		t.Errorf("unowned-slug D.reason=%q must NOT contain a pull/rebase verb", unowned.Reason)
	}
}

// TestReconcileDUntrustedRosterAllReports (AC-2) — with no team identity
// (untrusted roster), BOTH worktrees are Owned:false / report-only, because
// without a trusted roster we cannot prove we own anything.
func TestReconcileDUntrustedRosterAllReports(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	f := newDEOwnershipFixture(t)
	// teamName omitted and no session match → untrusted roster (git-only).
	bySlug := f.runD(t, "")

	for _, slug := range []string{"owned-slug", "unowned-slug"} {
		d, ok := bySlug[slug]
		if !ok {
			t.Fatalf("missing D entry for %s under untrusted roster; got %v", slug, bySlug)
		}
		if d.Owned {
			t.Errorf("%s D.owned=%v under untrusted roster, want false", slug, d.Owned)
		}
		r := strings.ToLower(d.Reason)
		if strings.Contains(r, "pull") || strings.Contains(r, "rebase") {
			t.Errorf("%s D.reason=%q under untrusted roster must NOT contain pull/rebase", slug, d.Reason)
		}
	}
}
