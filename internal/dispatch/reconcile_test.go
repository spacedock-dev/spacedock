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
	"sort"
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

	// 2. Entities — covering each AC-1 case plus L1-M1/M2 audit-fix rows.
	// alive: a clean (no-drift) entity with a clean worktree (no D, no C).
	wtClean := filepath.Join(repoRoot, ".worktrees", "spacedock-ensign-alive")
	writeFile(t, filepath.Join(stateRoot, "alive", "index.md"), reconcileEntityFM(map[string]string{
		"id":       "id-alive",
		"title":    "alive entity",
		"slug":     "alive",
		"status":   "implementation",
		"worktree": filepath.Join(".worktrees", "spacedock-ensign-alive"),
	}))
	// release-notes-local-summary: archived → its ensign is class A
	// (archived-branch path of classA).
	writeFile(t, filepath.Join(stateRoot, "_archive", "release-notes-local-summary", "index.md"),
		reconcileEntityFM(map[string]string{
			"id":     "id-rnls",
			"title":  "release notes",
			"slug":   "release-notes-local-summary",
			"status": "done",
		}))
	// terminal-not-yet-archived: active (NOT under _archive/) with status=done →
	// class A via the active-status=done branch (L1-M1 fix). Without this row
	// deleting reconcile.go's "active.status==done" emit block (the second of
	// classA's two disjoint paths) goes undetected.
	writeFile(t, filepath.Join(stateRoot, "terminal-not-yet-archived", "index.md"),
		reconcileEntityFM(map[string]string{
			"id":     "id-tnya",
			"title":  "terminal not yet archived",
			"slug":   "terminal-not-yet-archived",
			"status": "done",
		}))
	// cohort: same slug+stage as two ensigns — class B (loser flagged).
	writeFile(t, filepath.Join(stateRoot, "cohort", "index.md"), reconcileEntityFM(map[string]string{
		"id":     "id-cohort",
		"title":  "cohort entity",
		"slug":   "cohort",
		"status": "ideation",
	}))
	// pr-merged: pr set, status NOT done, gh MERGED → class C.
	writeFile(t, filepath.Join(stateRoot, "pr-merged", "index.md"), reconcileEntityFM(map[string]string{
		"id":     "id-prm",
		"title":  "pr merged",
		"slug":   "pr-merged",
		"status": "implementation",
		"pr":     "42",
	}))
	// pr-open: pr set, gh state=OPEN → no C entry (L1-M2 negative case for the
	// MERGED conjunct). Without this row, deleting the `state == "MERGED"`
	// conjunct in reconcile.go's classC goes undetected.
	writeFile(t, filepath.Join(stateRoot, "pr-open", "index.md"), reconcileEntityFM(map[string]string{
		"id":     "id-pro",
		"title":  "pr open",
		"slug":   "pr-open",
		"status": "implementation",
		"pr":     "43",
	}))
	// pr-merged-done: pr set, status=done, gh state=MERGED → no C entry
	// (L1-M2 negative case for the status!=done conjunct). Without this row,
	// deleting the `rec.status != "done"` conjunct in classC goes undetected.
	writeFile(t, filepath.Join(stateRoot, "pr-merged-done", "index.md"), reconcileEntityFM(map[string]string{
		"id":     "id-prmd",
		"title":  "pr merged done",
		"slug":   "pr-merged-done",
		"status": "done",
		"pr":     "44",
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

	// 3. Team config — team-lead (exempt), the alive ensign, A losers
	// (archived + terminal-not-yet-archived), B cohort (winner + loser-cycle1),
	// C ensigns (pr-merged + pr-open + pr-merged-done), and D ensign.
	// The pr-open and pr-merged-done ensigns are present so the helper sees
	// them as members; the assertions confirm they DO NOT emit C drift.
	cfgPath := filepath.Join(home, ".claude", "teams", teamName, "config.json")
	writeFile(t, cfgPath, teamConfigJSON(teamName, []claudeteam.ReconcileMember{
		{Name: "team-lead", AgentType: "team-lead", Model: "claude-opus-4-7"},
		{Name: "spacedock-ensign-alive-implementation",
			AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
		{Name: "spacedock-ensign-release-notes-local-summary-implementation",
			AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
		{Name: "spacedock-ensign-terminal-not-yet-archived-implementation",
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
		ghResponses: map[string]string{
			"42": "MERGED", // pr-merged: status!=done → emits C
			"43": "OPEN",   // pr-open: MERGED conjunct fails → no C
			"44": "MERGED", // pr-merged-done: status==done → no C
		},
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

	byClassAll := groupDriftByClass(result.Drift)
	for _, class := range []string{"A", "B", "C", "D", "E"} {
		if len(byClassAll[class]) == 0 {
			t.Errorf("expected at least one class %s drift entry; got: %s",
				class, formatDrift(result.Drift))
		}
	}

	// AC-1 wants one entry per class; L1-M1 splits A into two disjoint cases
	// (archived vs active-status=done), so A has TWO entries by design here.
	asByName := map[string]driftItem{}
	for _, a := range byClassAll["A"] {
		asByName[a.Name] = a
	}
	if len(asByName) != 2 {
		t.Errorf("expected 2 A entries (archived + active-status=done); got %d: %s",
			len(asByName), formatDrift(byClassAll["A"]))
	}
	// Archived-branch A entry — Reason must mention "archived".
	aArchived, ok := asByName["spacedock-ensign-release-notes-local-summary-implementation"]
	if !ok {
		t.Errorf("missing A entry for archived release-notes-local-summary; got names: %v",
			keysOfDriftMap(asByName))
	} else {
		if aArchived.Slug != "release-notes-local-summary" {
			t.Errorf("archived A.slug=%q", aArchived.Slug)
		}
		if !strings.Contains(aArchived.Reason, "archived") {
			t.Errorf("archived A.reason=%q, want substring 'archived'", aArchived.Reason)
		}
	}
	// Active-status=done branch A entry (L1-M1) — Reason must mention
	// "status=done" (distinct token from "archived" so an injected regression
	// that disables the active-status=done emit block is caught here).
	aActiveDone, ok := asByName["spacedock-ensign-terminal-not-yet-archived-implementation"]
	if !ok {
		t.Errorf("missing A entry for active terminal-not-yet-archived; got names: %v",
			keysOfDriftMap(asByName))
	} else {
		if aActiveDone.Slug != "terminal-not-yet-archived" {
			t.Errorf("active-done A.slug=%q", aActiveDone.Slug)
		}
		if !strings.Contains(aActiveDone.Reason, "status=done") {
			t.Errorf("active-done A.reason=%q, want substring 'status=done'", aActiveDone.Reason)
		}
		if strings.Contains(aActiveDone.Reason, "archived") {
			t.Errorf("active-done A.reason=%q, must NOT mention 'archived' (the two branches must remain distinguishable in output)",
				aActiveDone.Reason)
		}
	}

	bs := byClassAll["B"]
	if len(bs) != 1 {
		t.Errorf("expected 1 B entry; got %d: %s", len(bs), formatDrift(bs))
	}
	if len(bs) > 0 {
		b := bs[0]
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
	}

	// L1-M2: exactly ONE C entry — pr-merged only. pr-open (gh=OPEN) must NOT
	// emit C (MERGED conjunct check); pr-merged-done (status=done) must NOT
	// emit C (status!=done conjunct check). The exact-Reason check tightens
	// the assertion beyond the prior "contains 'merged'" substring.
	cs := byClassAll["C"]
	if len(cs) != 1 {
		t.Errorf("expected exactly 1 C entry (pr-merged only); got %d: %s",
			len(cs), formatDrift(cs))
	}
	if len(cs) > 0 {
		c := cs[0]
		if c.Slug != "pr-merged" || c.PR != "42" {
			t.Errorf("C {slug=%q pr=%q}, want {pr-merged 42}", c.Slug, c.PR)
		}
		const wantCReason = "PR merged but status=implementation"
		if c.Reason != wantCReason {
			t.Errorf("C.reason=%q, want exact %q", c.Reason, wantCReason)
		}
	}
	for _, c := range cs {
		if c.Slug == "pr-open" {
			t.Errorf("pr-open (gh=OPEN) must NOT emit C drift; got %+v", c)
		}
		if c.Slug == "pr-merged-done" {
			t.Errorf("pr-merged-done (status=done) must NOT emit C drift; got %+v", c)
		}
	}

	ds := byClassAll["D"]
	if len(ds) != 1 {
		t.Errorf("expected 1 D entry; got %d", len(ds))
	}
	if len(ds) > 0 {
		d := ds[0]
		if d.Slug != "yaml-parser-migration" {
			t.Errorf("D.slug=%q", d.Slug)
		}
		if d.Behind <= 0 {
			t.Errorf("D.behind=%d, want > 0", d.Behind)
		}
		if !strings.Contains(d.Worktree, "yaml-parser-migration") {
			t.Errorf("D.worktree=%q", d.Worktree)
		}
	}

	es := byClassAll["E"]
	if len(es) != 1 {
		t.Errorf("expected 1 E entry; got %d", len(es))
	}
	if len(es) > 0 {
		e := es[0]
		if e.Ahead <= 0 {
			t.Errorf("E.ahead=%d, want > 0", e.Ahead)
		}
		if !strings.Contains(e.Reason, "reset main") {
			t.Errorf("E.reason=%q, want reset main", e.Reason)
		}
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

// teamConfigJSON renders a config.json the helper can parse, with the default
// fixture lead session id.
func teamConfigJSON(name string, members []claudeteam.ReconcileMember) string {
	return teamConfigJSONWithSession(name, "session-fixture", members)
}

// teamConfigJSONWithSession renders a config.json carrying the given
// leadSessionId so session-scoping tests can seed a match or a decoy.
func teamConfigJSONWithSession(name, leadSessionID string, members []claudeteam.ReconcileMember) string {
	type m struct {
		Name      string `json:"name"`
		AgentType string `json:"agentType"`
		Model     string `json:"model"`
	}
	out := struct {
		Name          string `json:"name"`
		LeadSessionID string `json:"leadSessionId"`
		Members       []m    `json:"members"`
	}{Name: name, LeadSessionID: leadSessionID}
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

// groupDriftByClass groups drift items by their Class field, allowing
// multi-entry classes (used by the L1-M1 fixture which now has two A entries).
func groupDriftByClass(drift []driftItem) map[string][]driftItem {
	out := map[string][]driftItem{}
	for _, d := range drift {
		out[d.Class] = append(out[d.Class], d)
	}
	return out
}

// keysOfDriftMap returns the keys of a name->driftItem map sorted for
// deterministic error output.
func keysOfDriftMap(m map[string]driftItem) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func formatDrift(drift []driftItem) string {
	b, _ := json.MarshalIndent(drift, "", "  ")
	return string(b)
}
