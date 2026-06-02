// ABOUTME: AC-1 fixture for `dispatch reconcile` — five disjoint drift classes
// ABOUTME: detected from a tmpdir team config + workflow + worktree, plus a flip.
package dispatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
)

// reconcileFixture is one self-contained fixture that holds a tmp HOME (so the
// team config glob is hermetic), a tmp git repo with a state checkout, a couple
// of worktrees, a stubbed gh, and a tiny shim around the helper's gitRunner so
// classD/classE can be exercised against a real git tree.
type reconcileFixture struct {
	t           *testing.T
	home        string
	repoRoot    string
	workflowDir string
	stateRoot   string
	teamName    string
	cfgPath     string
	wtClean     string // worktree for a clean (D-class-skipped) entity
	wtStale     string // worktree for a class-D entity (HEAD behind origin/next)
	ghResponses map[string]string
}

func newReconcileFixture(t *testing.T) *reconcileFixture {
	t.Helper()
	home := t.TempDir()
	repoRoot := t.TempDir()
	teamName := "team-fixture-001"

	// 1. Workflow under the repo (split-root: README declares state: .spacedock-state).
	workflowDir := filepath.Join(repoRoot, "docs", "wf")
	writeFile(t, filepath.Join(workflowDir, "README.md"), `---
entity-type: task
state: .spacedock-state
stages:
  states:
    - name: backlog
      initial: true
    - name: ideation
    - name: implementation
      worktree: true
    - name: validation
      worktree: true
      fresh: true
    - name: done
      terminal: true
---
`)
	stateRoot := filepath.Join(workflowDir, ".spacedock-state")
	if err := os.MkdirAll(stateRoot, 0o755); err != nil {
		t.Fatal(err)
	}

	// 2. Entities — five active slugs covering each AC-1 case.
	// alive: a clean (no-drift) entity with a clean worktree (no D, no C).
	wtClean := filepath.Join(repoRoot, ".worktrees", "spacedock-ensign-alive")
	writeFile(t, filepath.Join(stateRoot, "alive", "index.md"), reconcileEntityFM(map[string]string{
		"id":       "id-alive",
		"title":    "alive entity",
		"slug":     "alive",
		"status":   "implementation",
		"worktree": filepath.Join(".worktrees", "spacedock-ensign-alive"),
	}))
	// release-notes-local-summary: archived → its ensign is class A.
	writeFile(t, filepath.Join(stateRoot, "_archive", "release-notes-local-summary", "index.md"),
		reconcileEntityFM(map[string]string{
			"id":     "id-rnls",
			"title":  "release notes",
			"slug":   "release-notes-local-summary",
			"status": "done",
		}))
	// cohort: same slug+stage as two ensigns — class B (loser flagged).
	writeFile(t, filepath.Join(stateRoot, "cohort", "index.md"), reconcileEntityFM(map[string]string{
		"id":     "id-cohort",
		"title":  "cohort entity",
		"slug":   "cohort",
		"status": "ideation",
	}))
	// pr-merged: pr set, status NOT done — class C when gh returns MERGED.
	writeFile(t, filepath.Join(stateRoot, "pr-merged", "index.md"), reconcileEntityFM(map[string]string{
		"id":     "id-prm",
		"title":  "pr merged",
		"slug":   "pr-merged",
		"status": "implementation",
		"pr":     "42",
	}))
	// yaml-parser-migration: worktree branch is BEHIND origin/next — class D.
	wtStale := filepath.Join(repoRoot, ".worktrees", "spacedock-ensign-yaml-parser-migration")
	writeFile(t, filepath.Join(stateRoot, "yaml-parser-migration", "index.md"),
		reconcileEntityFM(map[string]string{
			"id":       "id-ypm",
			"title":    "yaml parser migration",
			"slug":     "yaml-parser-migration",
			"status":   "implementation",
			"worktree": filepath.Join(".worktrees", "spacedock-ensign-yaml-parser-migration"),
		}))

	// 3. Team config — team-lead (exempt), one alive ensign, A loser, B cohort
	// (winner + loser-cycle1), C and D ensigns.
	cfgPath := filepath.Join(home, ".claude", "teams", teamName, "config.json")
	writeFile(t, cfgPath, teamConfigJSON(teamName, []claudeteam.ReconcileMember{
		{Name: "team-lead", AgentType: "team-lead", Model: "claude-opus-4-7"},
		{Name: "spacedock-ensign-alive-implementation",
			AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
		{Name: "spacedock-ensign-release-notes-local-summary-implementation",
			AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
		{Name: "spacedock-ensign-cohort-ideation",
			AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
		{Name: "spacedock-ensign-cohort-ideation-cycle2",
			AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
		{Name: "spacedock-ensign-pr-merged-implementation",
			AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
		{Name: "spacedock-ensign-yaml-parser-migration-implementation",
			AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
	}))

	// 4. Real git graph with main, origin/next, and two worktree branches at
	// distinct refs. The graph we build:
	//
	//   initial (seed.txt)
	//   ├── main (+ second-commit.txt) — 1 ahead of origin/next (class E)
	//   └── refs/remotes/origin/next
	//       └── (initial) — then bumped to:
	//           └── next-bump.txt — wtClean is HERE, wtStale is at initial
	//
	// So wtClean is up to date with origin/next; wtStale is 1 behind (class D).
	repoGitInit(t, repoRoot)
	repoSetOriginNext(t, repoRoot, "HEAD") // origin/next = initial commit
	// wtStale: created BEFORE bumping origin/next, so it lags by 1.
	makeFixtureWorktree(t, repoRoot, wtStale, "ypm-branch", "origin/next")
	// main: diverge with a non-origin/next commit so class E fires.
	repoMakeCommit(t, repoRoot, "main", "second-commit.txt", "two\n")
	// Bump origin/next to a new commit (not a parent of main).
	repoBumpOriginNext(t, repoRoot)
	// wtClean: created AFTER bumping origin/next, so HEAD..origin/next == 0.
	makeFixtureWorktree(t, repoRoot, wtClean, "alive-branch", "origin/next")

	return &reconcileFixture{
		t:           t,
		home:        home,
		repoRoot:    repoRoot,
		workflowDir: workflowDir,
		stateRoot:   stateRoot,
		teamName:    teamName,
		cfgPath:     cfgPath,
		wtClean:     wtClean,
		wtStale:     wtStale,
		ghResponses: map[string]string{"42": "MERGED"},
	}
}

// run drives the helper through Reconcile() with the fixture's wiring and
// returns the parsed JSON envelope.
func (f *reconcileFixture) run() reconcileResult {
	f.t.Helper()
	var stdout, stderr bytes.Buffer
	opts := reconcileOpts{
		workflowDir: f.workflowDir,
		teamName:    f.teamName,
		repoRoot:    f.repoRoot,
		include:     map[string]bool{"A": true, "B": true, "C": true, "D": true, "E": true},
		home:        f.home,
		roster:      claudeteam.LoadReconcileTeam,
		gh: func(pr string) (string, error) {
			if state, ok := f.ghResponses[pr]; ok {
				return state, nil
			}
			return "", fmt.Errorf("no fixture gh response for pr=%s", pr)
		},
		git: gitRunnerExec,
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

// TestReconcileFiveClasses runs the full sweep on the synthetic fixture and
// asserts exactly one entry of each of A, B, C, D, E with the documented
// slug/name/reason fields.
func TestReconcileFiveClasses(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	f := newReconcileFixture(t)
	result := f.run()

	if result.Command != "reconcile" {
		t.Errorf("command=%q, want reconcile", result.Command)
	}
	if result.TeamName != f.teamName {
		t.Errorf("team_name=%q, want %q", result.TeamName, f.teamName)
	}

	byClass := indexDriftByClass(t, result.Drift)
	if len(byClass) != 5 {
		t.Fatalf("expected 5 classes, got %d (%v): %s",
			len(byClass), keys(byClass), formatDrift(result.Drift))
	}
	a := requireOne(t, byClass, "A")
	if a.Name != "spacedock-ensign-release-notes-local-summary-implementation" {
		t.Errorf("A.name=%q, want release-notes-local-summary-implementation", a.Name)
	}
	if a.Slug != "release-notes-local-summary" {
		t.Errorf("A.slug=%q", a.Slug)
	}
	if !strings.Contains(a.Reason, "archived") {
		t.Errorf("A.reason=%q, want archived", a.Reason)
	}

	b := requireOne(t, byClass, "B")
	if b.Slug != "cohort" || b.Stage != "ideation" {
		t.Errorf("B {slug=%q stage=%q}, want {cohort ideation}", b.Slug, b.Stage)
	}
	// The loser is the unsuffixed (=cycle 1) name; the winner is cycle2.
	if b.Name != "spacedock-ensign-cohort-ideation" {
		t.Errorf("B.name=%q, want spacedock-ensign-cohort-ideation (loser)", b.Name)
	}
	if !strings.Contains(b.Reason, "spacedock-ensign-cohort-ideation-cycle2") {
		t.Errorf("B.reason=%q, want winner cycle2 named", b.Reason)
	}

	c := requireOne(t, byClass, "C")
	if c.Slug != "pr-merged" || c.PR != "42" {
		t.Errorf("C {slug=%q pr=%q}, want {pr-merged 42}", c.Slug, c.PR)
	}
	if !strings.Contains(c.Reason, "merged") {
		t.Errorf("C.reason=%q, want merged", c.Reason)
	}

	d := requireOne(t, byClass, "D")
	if d.Slug != "yaml-parser-migration" {
		t.Errorf("D.slug=%q", d.Slug)
	}
	if d.Behind <= 0 {
		t.Errorf("D.behind=%d, want > 0", d.Behind)
	}
	if !strings.Contains(d.Worktree, "yaml-parser-migration") {
		t.Errorf("D.worktree=%q", d.Worktree)
	}

	e := requireOne(t, byClass, "E")
	if e.Ahead <= 0 {
		t.Errorf("E.ahead=%d, want > 0", e.Ahead)
	}
	if !strings.Contains(e.Reason, "reset main") {
		t.Errorf("E.reason=%q, want reset main", e.Reason)
	}
}

// TestReconcileFlipReclassifies mutates the fixture (archive `alive` →
// alive-implementation becomes class A; the prior Class A entity remains A)
// then re-runs to confirm the classification is data-driven, not hard-coded.
func TestReconcileFlipReclassifies(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	f := newReconcileFixture(t)
	// Before flip: alive is not in any drift class.
	before := f.run()
	for _, d := range before.Drift {
		if d.Name == "spacedock-ensign-alive-implementation" {
			t.Fatalf("alive ensign should be clean before flip; got %+v", d)
		}
	}

	// Flip: archive `alive` by moving its dir into _archive.
	src := filepath.Join(f.stateRoot, "alive")
	dst := filepath.Join(f.stateRoot, "_archive", "alive")
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(src, dst); err != nil {
		t.Fatal(err)
	}
	after := f.run()
	foundA := 0
	for _, d := range after.Drift {
		if d.Class == "A" && d.Name == "spacedock-ensign-alive-implementation" {
			foundA++
		}
	}
	if foundA != 1 {
		t.Errorf("after flip: expected 1 A entry for alive, got %d. drift=%s",
			foundA, formatDrift(after.Drift))
	}
	// The original archived entity's A entry still fires too.
	origA := 0
	for _, d := range after.Drift {
		if d.Class == "A" && d.Slug == "release-notes-local-summary" {
			origA++
		}
	}
	if origA != 1 {
		t.Errorf("original A entry missing after flip; drift=%s", formatDrift(after.Drift))
	}
}

// TestReconcileIncludeScope confirms --include scopes the sweep — passing only
// "A,B" emits only class-A/B entries.
func TestReconcileIncludeScope(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	f := newReconcileFixture(t)
	var stdout, stderr bytes.Buffer
	opts := reconcileOpts{
		workflowDir: f.workflowDir,
		teamName:    f.teamName,
		repoRoot:    f.repoRoot,
		include:     map[string]bool{"A": true, "B": true},
		home:        f.home,
		roster:      claudeteam.LoadReconcileTeam,
		gh: func(pr string) (string, error) {
			return f.ghResponses[pr], nil
		},
		git: gitRunnerExec,
	}
	if code := Reconcile(opts, &stdout, &stderr); code != 0 {
		t.Fatalf("Reconcile exit=%d stderr=%s", code, stderr.String())
	}
	var result reconcileResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	for _, d := range result.Drift {
		if d.Class != "A" && d.Class != "B" {
			t.Errorf("--include=A,B emitted %s drift: %+v", d.Class, d)
		}
	}
}

// TestReconcileMissingWorkflowDir surfaces a setup failure (exit 1).
func TestReconcileMissingWorkflowDir(t *testing.T) {
	home := t.TempDir()
	var stdout, stderr bytes.Buffer
	opts := reconcileOpts{
		workflowDir: filepath.Join(home, "no-such-dir"),
		teamName:    "any",
		home:        home,
		roster:      claudeteam.LoadReconcileTeam,
		gh:          func(string) (string, error) { return "", nil },
		git:         gitRunnerExec,
	}
	code := Reconcile(opts, &stdout, &stderr)
	if code != 1 {
		t.Errorf("missing workflow dir exit=%d, want 1; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "workflow directory not found") {
		t.Errorf("missing workflow stderr=%q", stderr.String())
	}
}

// TestReconcileUsageError covers exit 2 paths (missing --workflow-dir, bad --include).
func TestReconcileUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runReconcile([]string{}, &stdout, &stderr)
	if code != 2 {
		t.Errorf("missing --workflow-dir exit=%d, want 2", code)
	}
	stdout.Reset()
	stderr.Reset()
	code = runReconcile([]string{"--workflow-dir", "/tmp", "--include", "X"}, &stdout, &stderr)
	if code != 2 {
		t.Errorf("bad --include exit=%d, want 2", code)
	}
}

// Helpers ----------------------------------------------------------------------

// reconcileEntityFM renders a minimal frontmatter file body for the fixture.
// status/worktree/pr fields render unquoted (the parser tolerates either; the
// live FO emits unquoted), title in plain text.
func reconcileEntityFM(fields map[string]string) string {
	var b strings.Builder
	b.WriteString("---\n")
	for _, k := range []string{"id", "title", "slug", "status", "worktree", "pr"} {
		v, ok := fields[k]
		if !ok {
			continue
		}
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(v)
		b.WriteString("\n")
	}
	b.WriteString("---\n\nEntity body.\n")
	return b.String()
}

// teamConfigJSON renders a config.json the helper can parse.
func teamConfigJSON(name string, members []claudeteam.ReconcileMember) string {
	type m struct {
		Name      string `json:"name"`
		AgentType string `json:"agentType"`
		Model     string `json:"model"`
	}
	out := struct {
		Name          string `json:"name"`
		LeadSessionID string `json:"leadSessionId"`
		Members       []m    `json:"members"`
	}{Name: name, LeadSessionID: "session-fixture"}
	for _, x := range members {
		out.Members = append(out.Members, m{x.Name, x.AgentType, x.Model})
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	return string(b) + "\n"
}

// repoGitInit initializes the fixture repo with one commit on main and sets
// origin (a local file:// remote at a sibling dir) and a "next" branch.
func repoGitInit(t *testing.T, dir string) {
	t.Helper()
	mustGit(t, dir, "init", "-q", "-b", "main")
	mustGit(t, dir, "config", "user.email", "t@t")
	mustGit(t, dir, "config", "user.name", "t")
	writeFile(t, filepath.Join(dir, "seed.txt"), "seed\n")
	mustGit(t, dir, "add", "-A")
	mustGit(t, dir, "commit", "-q", "-m", "seed")
}

// repoSetOriginNext faked-remotes the origin/next ref to commit-ish by writing
// it through git's plumbing — no remote configured, just the ref the sweep reads.
func repoSetOriginNext(t *testing.T, dir, commitish string) {
	t.Helper()
	mustGit(t, dir, "update-ref", "refs/remotes/origin/next", commitish)
}

// repoMakeCommit adds a file and commits it on `branch` (creating or switching
// to it), without affecting origin/next. The commit lands on the named branch.
func repoMakeCommit(t *testing.T, dir, branch, file, content string) {
	t.Helper()
	mustGit(t, dir, "checkout", "-q", branch)
	writeFile(t, filepath.Join(dir, file), content)
	mustGit(t, dir, "add", "-A")
	mustGit(t, dir, "commit", "-q", "-m", "commit "+file)
}

// repoBumpOriginNext creates a new commit on a throwaway branch and points
// origin/next at it. The class-D stale-worktree was created BEFORE this, so its
// HEAD lags origin/next by 1.
func repoBumpOriginNext(t *testing.T, dir string) {
	t.Helper()
	// Park on a detached HEAD to avoid touching main; create a temp branch.
	mustGit(t, dir, "checkout", "-q", "refs/remotes/origin/next")
	writeFile(t, filepath.Join(dir, "next-bump.txt"), "bump\n")
	mustGit(t, dir, "add", "-A")
	mustGit(t, dir, "commit", "-q", "-m", "next bump")
	mustGit(t, dir, "update-ref", "refs/remotes/origin/next", "HEAD")
	// Return to main so subsequent helpers run cleanly.
	mustGit(t, dir, "checkout", "-q", "main")
}

// makeFixtureWorktree creates a worktree at path on a new branch pointing at
// commit-ish. The branch is created fresh to avoid colliding with main.
func makeFixtureWorktree(t *testing.T, repoRoot, wtPath, branch, commitish string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(wtPath), 0o755); err != nil {
		t.Fatal(err)
	}
	mustGit(t, repoRoot, "worktree", "add", "-q", "-b", branch, wtPath, commitish)
}

// mustGit runs a git command in dir, failing the test on a non-zero exit.
func mustGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-C", dir}, args...)
	out, err := exec.Command("git", full...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %s in %s: %v\n%s", strings.Join(args, " "), dir, err, out)
	}
	return strings.TrimSpace(string(out))
}

// hasGit reports whether the test environment has `git` on PATH. Reconcile
// tests skip when git is missing.
func hasGit(t *testing.T) bool {
	t.Helper()
	_, err := exec.LookPath("git")
	return err == nil
}

// indexDriftByClass groups drift items by their Class field and FAILS the test
// if any class appears more than once (every fixture entity is unique by design).
func indexDriftByClass(t *testing.T, drift []driftItem) map[string]driftItem {
	t.Helper()
	out := map[string]driftItem{}
	for _, d := range drift {
		if _, dup := out[d.Class]; dup {
			t.Errorf("duplicate drift for class %s: %s", d.Class, formatDrift(drift))
		}
		out[d.Class] = d
	}
	return out
}

func requireOne(t *testing.T, m map[string]driftItem, class string) driftItem {
	t.Helper()
	d, ok := m[class]
	if !ok {
		t.Fatalf("missing class %s drift; got: %v", class, keys(m))
	}
	return d
}

func keys(m map[string]driftItem) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func formatDrift(drift []driftItem) string {
	b, _ := json.MarshalIndent(drift, "", "  ")
	return string(b)
}
