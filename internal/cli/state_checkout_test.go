// ABOUTME: Real-git e2e for main-worktree-anchored split-root checkout resolution
// ABOUTME: (issue #484) — worktree-cwd anchoring, no-origin/unreachable-origin
// ABOUTME: resume fallback, and targeted stale-worktree-registration repair.
package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// addAgentWorktree adds a linked git worktree at hostRoot/.worktrees/<name>, on a
// new branch off HEAD — the shape a Spacedock agent dispatch creates
// (`.worktrees/<worker>-<entity>/`). Returns the worktree's absolute root, whose
// docs/dev carries the same commissioned README as the main worktree (same repo
// content, a different checkout location) — the setup issue #484 observed the
// cwd-relative bug from.
func addAgentWorktree(t *testing.T, hostRoot, name string) string {
	t.Helper()
	wtRoot := filepath.Join(hostRoot, ".worktrees", name)
	git(t, hostRoot, "worktree", "add", "-q", "-b", "agent/"+name, wtRoot, "HEAD")
	return wtRoot
}

// noOriginSplitWorkflow births a standalone (no `origin` remote at all) split-root
// workflow at root/docs/dev via the real `state new` command, matching the
// no-origin carve-out `state commit`/`state new` already document. Returns the
// workflow dir and its state checkout path (present after birth).
func noOriginSplitWorkflow(t *testing.T) (root, workflowDir, statePath string) {
	t.Helper()
	root = t.TempDir()
	git(t, root, "init", "-q")
	git(t, root, "config", "user.email", "t@t")
	git(t, root, "config", "user.name", "t")
	workflowDir = filepath.Join(root, "docs", "dev")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workflowDir, "README.md"), []byte(splitWorkflowReadme), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", "docs/dev/README.md")
	git(t, root, "commit", "-q", "-m", "add split-root README")

	code, _, stderr := execStateNew(t, root, workflowDir)
	if code != 0 {
		t.Fatalf("noOriginSplitWorkflow: state new failed; exit=%d stderr=%q", code, stderr)
	}
	statePath = filepath.Join(workflowDir, ".spacedock-state")
	return root, workflowDir, statePath
}

// TestStateReadyFromWorktreeCwdPresentCheckoutResolvesMain pins test-plan item 1 /
// AC-2: `state ready`, invoked with cwd inside an agent worktree, resolves the
// checkout ALREADY PRESENT at the main worktree — no nested copy appears under
// the worktree. Cwd variant: bare discovery (no --workflow-dir); the process cwd
// (the worktree root, mirroring where an ensign's dispatch actually starts) walks
// up/scans down to the worktree's own copy of docs/dev.
func TestStateReadyFromWorktreeCwdPresentCheckoutResolvesMain(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	bare := filepath.Join(t.TempDir(), "origin.git")
	git(t, t.TempDir(), "init", "-q", "--bare", bare)

	hostClone := filepath.Join(t.TempDir(), "host")
	git(t, t.TempDir(), "clone", "-q", bare, hostClone)
	git(t, hostClone, "config", "user.email", "t@t")
	git(t, hostClone, "config", "user.name", "t")
	workflowDir, _, _ := commissionSplitWorkflow(t, hostClone)
	git(t, hostClone, "push", "-q", "origin", "HEAD")
	mainCheckout := filepath.Join(workflowDir, ".spacedock-state")

	wtRoot := addAgentWorktree(t, hostClone, "w1")
	nestedCheckout := filepath.Join(wtRoot, "docs", "dev", ".spacedock-state")

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "ready"},
		os.Environ(), wtRoot, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("state ready from worktree cwd (checkout present) should exit 0; got exit=%d stderr=%q", code, errBuf.String())
	}
	if _, err := os.Stat(mainCheckout); err != nil {
		t.Fatalf("main checkout should still be present at %s: %v", mainCheckout, err)
	}
	if _, err := os.Stat(nestedCheckout); !os.IsNotExist(err) {
		t.Fatalf("state ready must NOT create a nested checkout under the worktree; found at %s", nestedCheckout)
	}
}

// TestStateReadyFromWorktreeCwdAbsentCheckoutResumesMain pins test-plan item 2 /
// AC-2: `state ready` with the checkout ABSENT (fresh clone), invoked with cwd
// inside an agent worktree, resumes the checkout at the MAIN-anchored path, never
// under the worktree. Cwd variant: bare discovery from the worktree root, same as
// the present-checkout test above.
func TestStateReadyFromWorktreeCwdAbsentCheckoutResumesMain(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	bare := filepath.Join(t.TempDir(), "origin.git")
	git(t, t.TempDir(), "init", "-q", "--bare", bare)

	hostA := filepath.Join(t.TempDir(), "hostA")
	git(t, t.TempDir(), "clone", "-q", bare, hostA)
	git(t, hostA, "config", "user.email", "a@t")
	git(t, hostA, "config", "user.name", "a")
	commissionSplitWorkflow(t, hostA)
	git(t, hostA, "push", "-q", "origin", "HEAD")

	fresh := filepath.Join(t.TempDir(), "fresh")
	git(t, t.TempDir(), "clone", "-q", bare, fresh)
	git(t, fresh, "config", "user.email", "f@t")
	git(t, fresh, "config", "user.name", "f")
	mainCheckout := filepath.Join(fresh, "docs", "dev", ".spacedock-state")

	wtRoot := addAgentWorktree(t, fresh, "w1")
	nestedCheckout := filepath.Join(wtRoot, "docs", "dev", ".spacedock-state")

	if _, err := os.Stat(mainCheckout); !os.IsNotExist(err) {
		t.Fatalf("precondition: fresh clone should NOT yet have the state checkout (err=%v)", err)
	}

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "ready"},
		os.Environ(), wtRoot, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("state ready from worktree cwd (checkout absent) should resume (exit 0); got exit=%d stderr=%q", code, errBuf.String())
	}
	if _, err := os.Stat(mainCheckout); err != nil {
		t.Fatalf("state ready should have resumed the checkout at the MAIN-anchored path %s: %v", mainCheckout, err)
	}
	if _, err := os.Stat(nestedCheckout); !os.IsNotExist(err) {
		t.Fatalf("state ready must NOT resume a nested checkout under the worktree; found at %s", nestedCheckout)
	}
}

// TestStateReadyIssue484Repro pins test-plan item 3, AC-1 (the value-measure
// baseline) and AC-5: the full issue #484 repro — a no-origin repo whose split-root
// checkout directory was deleted while its worktree registration survived —
// converges from a worktree cwd. Cwd variant: bare discovery from the worktree
// root. Where 0.24.0-pre2 exits non-zero and proposes a worktree-nested path, this
// must exit 0 and restore the checkout at the main-root state path with the
// pre-existing entity still readable, and no `.spacedock-state` anywhere under
// `.worktrees/`.
func TestStateReadyIssue484Repro(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root, _, statePath := noOriginSplitWorkflow(t)

	// Seed an entity in the checkout so AC-1's "entity file readable" has
	// something concrete to assert on, then delete the checkout dir while
	// leaving the worktree registration behind (the reported repro: `git
	// worktree list --porcelain` still names the now-missing directory).
	entity := filepath.Join(statePath, "first-task.md")
	if err := os.WriteFile(entity, []byte("---\nstatus: ideation\n---\n# First Task\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, statePath, "add", "-A")
	git(t, statePath, "commit", "-q", "-m", "seed first-task")
	if err := os.RemoveAll(statePath); err != nil {
		t.Fatal(err)
	}
	if out := git(t, root, "worktree", "list", "--porcelain"); !strings.Contains(out, statePath) {
		t.Fatalf("precondition: worktree registration for %s should survive the directory deletion; list=%q", statePath, out)
	}

	wtRoot := addAgentWorktree(t, root, "w1")

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "ready"},
		os.Environ(), wtRoot, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("issue #484 repro should converge (exit 0); got exit=%d stdout=%q stderr=%q", code, out.String(), errBuf.String())
	}
	if _, err := os.Stat(filepath.Join(statePath, "first-task.md")); err != nil {
		t.Fatalf("restored checkout should carry the pre-existing entity: %v", err)
	}

	worktreesDir := filepath.Join(root, ".worktrees")
	walkErr := filepath.Walk(worktreesDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() && info.Name() == ".spacedock-state" {
			t.Fatalf("no .spacedock-state may appear anywhere under .worktrees/; found %s", path)
		}
		return nil
	})
	if walkErr != nil {
		t.Fatalf("walking %s failed: %v", worktreesDir, walkErr)
	}
}

// TestStateReadyNoOriginResumeFromRoot pins test-plan item 4, AC-3: a no-origin
// repo resumes an absent checkout from the local state branch when invoked from
// the repo root (no worktree involved) — exit 0, entities present, local-only
// wording in the output.
func TestStateReadyNoOriginResumeFromRoot(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root, workflowDir, statePath := noOriginSplitWorkflow(t)
	if err := os.RemoveAll(statePath); err != nil {
		t.Fatal(err)
	}

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "ready", "--workflow-dir", workflowDir},
		os.Environ(), root, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("no-origin resume from root should exit 0; got exit=%d stderr=%q", code, errBuf.String())
	}
	if _, err := os.Stat(statePath); err != nil {
		t.Fatalf("no-origin resume should have restored the checkout: %v", err)
	}
	if !strings.Contains(out.String(), "local-only") {
		t.Fatalf("no-origin resume output should say local-only; got stdout=%q", out.String())
	}
}

// TestStateReadyUnreachableOriginFallsBackToLocalBranch pins test-plan item 5
// (with-local-branch half), AC-4: origin is CONFIGURED but unreachable (a bad
// local path, never a real remote); resume falls back to the local state branch,
// exits 0, and warns that the fetch failed.
func TestStateReadyUnreachableOriginFallsBackToLocalBranch(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	bare := filepath.Join(t.TempDir(), "origin.git")
	git(t, t.TempDir(), "init", "-q", "--bare", bare)

	hostClone := filepath.Join(t.TempDir(), "host")
	git(t, t.TempDir(), "clone", "-q", bare, hostClone)
	git(t, hostClone, "config", "user.email", "t@t")
	git(t, hostClone, "config", "user.name", "t")
	workflowDir, _, _ := commissionSplitWorkflow(t, hostClone)
	git(t, hostClone, "push", "-q", "origin", "HEAD")
	statePath := filepath.Join(workflowDir, ".spacedock-state")

	if err := os.RemoveAll(statePath); err != nil {
		t.Fatal(err)
	}
	git(t, hostClone, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "does-not-exist.git"))

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "ready", "--workflow-dir", workflowDir},
		os.Environ(), hostClone, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("unreachable-origin resume with a local branch should exit 0; got exit=%d stdout=%q stderr=%q", code, out.String(), errBuf.String())
	}
	if _, err := os.Stat(statePath); err != nil {
		t.Fatalf("resume should have restored the checkout from the local branch: %v", err)
	}
	if !strings.Contains(out.String(), "Warning") || !strings.Contains(out.String(), "fetch") {
		t.Fatalf("resume should warn that the fetch failed; got stdout=%q", out.String())
	}
}

// TestStateReadyUnreachableOriginNoLocalBranchHintsMainAnchoredPath pins test-plan
// item 5 (without-local-branch half), AC-4: origin is configured but unreachable
// AND no local branch exists — resume cannot fall back to anything, exits
// non-zero, and the hint names the main-anchored checkout path via the generic
// manual-fallback wording (NOT the "state new" never-birthed wording — an
// unreachable origin is indeterminate, not proof the branch was never birthed).
func TestStateReadyUnreachableOriginNoLocalBranchHintsMainAnchoredPath(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := t.TempDir()
	git(t, root, "init", "-q")
	git(t, root, "config", "user.email", "t@t")
	git(t, root, "config", "user.name", "t")
	workflowDir := filepath.Join(root, "docs", "dev")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workflowDir, "README.md"), []byte(splitWorkflowReadme), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", "docs/dev/README.md")
	git(t, root, "commit", "-q", "-m", "add split-root README")
	git(t, root, "remote", "add", "origin", filepath.Join(t.TempDir(), "does-not-exist.git"))
	statePath := filepath.Join(workflowDir, ".spacedock-state")

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "ready", "--workflow-dir", workflowDir},
		os.Environ(), root, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code == 0 {
		t.Fatalf("unreachable origin with no local branch must fail, not silently succeed; stdout=%q", out.String())
	}
	if !strings.Contains(errBuf.String(), "Manual fallback") {
		t.Fatalf("hint should carry the manual-fallback wording; stderr=%q", errBuf.String())
	}
	if !strings.Contains(errBuf.String(), statePath) {
		t.Fatalf("hint should name the main-anchored checkout path %s; stderr=%q", statePath, errBuf.String())
	}
	if strings.Contains(errBuf.String(), "spacedock state new") {
		t.Fatalf("indeterminate (unreachable) origin must NOT claim the branch was never birthed; stderr=%q", errBuf.String())
	}
}

// TestStateReadyNeverBirthedFromWorktreeCwdHintsStateNewAtMainPath pins test-plan
// item 6 (AC-4) AND the ideation-gate-flagged AC-2 requirement: "one test forces
// the manual-fallback hint from a worktree cwd and asserts the hinted path." Cwd
// variant: EXPLICIT --workflow-dir naming a path physically inside the linked
// worktree (distinct from the bare-discovery-from-worktree-cwd variant exercised
// by the Present/Absent tests above) — the shape of a hook or script passing
// --workflow-dir explicitly from inside an agent worktree, issue #484's second
// named trigger ("hooks/scripts may invoke it from anywhere"). No origin and no
// local branch exist anywhere, so the branch was genuinely never birthed; the
// hint must name `spacedock state new` AND the MAIN-anchored path, never the
// worktree-nested one.
func TestStateReadyNeverBirthedFromWorktreeCwdHintsStateNewAtMainPath(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := t.TempDir()
	git(t, root, "init", "-q")
	git(t, root, "config", "user.email", "t@t")
	git(t, root, "config", "user.name", "t")
	workflowDir := filepath.Join(root, "docs", "dev")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workflowDir, "README.md"), []byte(splitWorkflowReadme), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", "docs/dev/README.md")
	git(t, root, "commit", "-q", "-m", "add split-root README")
	mainStatePath := filepath.Join(workflowDir, ".spacedock-state")

	wtRoot := addAgentWorktree(t, root, "w1")
	wtWorkflowDir := filepath.Join(wtRoot, "docs", "dev")
	nestedStatePath := filepath.Join(wtWorkflowDir, ".spacedock-state")

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "ready", "--workflow-dir", wtWorkflowDir},
		os.Environ(), wtRoot, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code == 0 {
		t.Fatalf("never-birthed branch must fail, not silently succeed; stdout=%q", out.String())
	}
	if !strings.Contains(errBuf.String(), "spacedock state new") {
		t.Fatalf("never-birthed hint should name `spacedock state new`; stderr=%q", errBuf.String())
	}
	if !strings.Contains(errBuf.String(), mainStatePath) {
		t.Fatalf("never-birthed hint should name the main-anchored path %s; stderr=%q", mainStatePath, errBuf.String())
	}
	if strings.Contains(errBuf.String(), nestedStatePath) {
		t.Fatalf("never-birthed hint must NOT name the worktree-nested path %s; stderr=%q", nestedStatePath, errBuf.String())
	}
}

// TestStateReadyPresentCheckoutRegressionUntouched pins test-plan item 7, AC-5's
// regression half: a PRESENT checkout's worktree registration and HEAD are
// untouched by `state ready` — the stale-registration repair must never fire when
// the directory already exists (dirExists guards it before repair runs).
func TestStateReadyPresentCheckoutRegressionUntouched(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	bare := filepath.Join(t.TempDir(), "origin.git")
	git(t, t.TempDir(), "init", "-q", "--bare", bare)

	hostClone := filepath.Join(t.TempDir(), "host")
	git(t, t.TempDir(), "clone", "-q", bare, hostClone)
	git(t, hostClone, "config", "user.email", "t@t")
	git(t, hostClone, "config", "user.name", "t")
	workflowDir, _, _ := commissionSplitWorkflow(t, hostClone)
	git(t, hostClone, "push", "-q", "origin", "HEAD")
	statePath := filepath.Join(workflowDir, ".spacedock-state")

	worktreesBefore := git(t, hostClone, "worktree", "list", "--porcelain")
	headBefore := strings.TrimSpace(git(t, statePath, "rev-parse", "HEAD"))

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "ready", "--workflow-dir", workflowDir},
		os.Environ(), hostClone, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("state ready on a present, up-to-date checkout should exit 0; got exit=%d stderr=%q", code, errBuf.String())
	}
	worktreesAfter := git(t, hostClone, "worktree", "list", "--porcelain")
	headAfter := strings.TrimSpace(git(t, statePath, "rev-parse", "HEAD"))
	if worktreesBefore != worktreesAfter {
		t.Fatalf("present-checkout state ready must not touch worktree registrations; before=%q after=%q", worktreesBefore, worktreesAfter)
	}
	if headBefore != headAfter {
		t.Fatalf("present-checkout state ready must not move HEAD; before=%s after=%s", headBefore, headAfter)
	}
}

// TestStateCommitFromWorktreeCwdLandsInMainCheckout pins test-plan item 8, AC-2:
// `state commit` (a second verb through the same shared resolver, not just
// `ready`) invoked with cwd inside an agent worktree lands its commit in the MAIN
// checkout, never a nested one. Cwd variant: bare discovery from the worktree
// root, matching TestStateReadyFromWorktreeCwd{Present,Absent}CheckoutResolvesMain.
func TestStateCommitFromWorktreeCwdLandsInMainCheckout(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	bare := filepath.Join(t.TempDir(), "origin.git")
	git(t, t.TempDir(), "init", "-q", "--bare", bare)

	hostClone := filepath.Join(t.TempDir(), "host")
	git(t, t.TempDir(), "clone", "-q", bare, hostClone)
	git(t, hostClone, "config", "user.email", "t@t")
	git(t, hostClone, "config", "user.name", "t")
	workflowDir, _, _ := commissionSplitWorkflow(t, hostClone)
	git(t, hostClone, "push", "-q", "origin", "HEAD")
	mainCheckout := filepath.Join(workflowDir, ".spacedock-state")

	writeEntity(t, workflowDir, "first-task", "---\nstatus: implementation\n---\n# First Task\n")

	wtRoot := addAgentWorktree(t, hostClone, "w1")
	nestedCheckout := filepath.Join(wtRoot, "docs", "dev", ".spacedock-state")

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "commit", "first-task", "-m", "from worktree cwd"},
		os.Environ(), wtRoot, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("state commit from worktree cwd should exit 0; got exit=%d stderr=%q", code, errBuf.String())
	}
	names := git(t, mainCheckout, "log", "--name-only", "--pretty=format:", "-1")
	if !strings.Contains(names, "first-task.md") {
		t.Fatalf("commit from worktree cwd should land in the main checkout's log; log:\n%s", names)
	}
	if _, err := os.Stat(nestedCheckout); !os.IsNotExist(err) {
		t.Fatalf("state commit must NOT create a nested checkout under the worktree; found at %s", nestedCheckout)
	}
}
