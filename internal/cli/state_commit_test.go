// ABOUTME: Real-git e2e for `state commit <slug>` — the rebase-HALT (AC-1), path-
// ABOUTME: scoped commit (AC-2), multi-writer happy path (AC-3), no-origin (AC-4).
package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// twoHostStateWorkflow builds a bare origin carrying a commissioned split-root
// workflow (code branch + orphan state branch + one seeded entity), then resumes
// two host clones (A and B) each with a linked-worktree state checkout on the
// shared spacedock-state/dev branch — the realistic two-writer split-root setup the
// sync verbs operate against. Returns each host's workflow dir and the state
// branch.
func twoHostStateWorkflow(t *testing.T) (bare, workflowA, workflowB, stateBranch string) {
	t.Helper()
	bare = filepath.Join(t.TempDir(), "origin.git")
	git(t, t.TempDir(), "init", "-q", "--bare", bare)

	hostA := filepath.Join(t.TempDir(), "hostA")
	git(t, t.TempDir(), "clone", "-q", bare, hostA)
	git(t, hostA, "config", "user.email", "a@t")
	git(t, hostA, "config", "user.name", "hostA")
	workflowA, stateBranch, _ = commissionSplitWorkflow(t, hostA)
	git(t, hostA, "push", "-q", "origin", "HEAD")

	// hostB clones the code branch, then `state init` fetches the orphan state branch
	// and adds B's own linked worktree at the gitignored state path.
	hostB := filepath.Join(t.TempDir(), "hostB")
	git(t, t.TempDir(), "clone", "-q", bare, hostB)
	git(t, hostB, "config", "user.email", "b@t")
	git(t, hostB, "config", "user.name", "hostB")
	workflowB = filepath.Join(hostB, "docs", "dev")
	stateInit(t, hostB, workflowB)
	return bare, workflowA, workflowB, stateBranch
}

// stateInit runs `spacedock state init` for workflowDir from the host clone, failing
// the test on a non-zero exit.
func stateInit(t *testing.T, hostDir, workflowDir string) {
	t.Helper()
	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "init", "--workflow-dir", workflowDir},
		os.Environ(), hostDir, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("state init exit=%d stderr=%q", code, errBuf.String())
	}
}

// writeEntity overwrites slug's entity file in the workflow's state checkout with
// body (no git ops — the verb does the staging/commit).
func writeEntity(t *testing.T, workflowDir, slug, body string) {
	t.Helper()
	checkout := filepath.Join(workflowDir, ".spacedock-state")
	if err := os.WriteFile(filepath.Join(checkout, slug+".md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// runStateCommitCmd runs `spacedock state commit slug` for workflowDir and returns
// the exit code plus captured stdout/stderr.
func runStateCommitCmd(t *testing.T, hostDir, workflowDir, slug string, extra ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errBuf strings.Builder
	args := append([]string{"state", "commit", slug, "--workflow-dir", workflowDir}, extra...)
	code = run(context.Background(), args, os.Environ(), hostDir, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	return code, out.String(), errBuf.String()
}

// TestStateCommitHaltsOnSameEntityConflict pins AC-1, the load-bearing halt: two
// writers edit the SAME entity's frontmatter; A pushes first; the verb runs as B
// and must exit 3, name the conflicting entity path on stderr, leave the checkout
// clean (rebase aborted), never force-push (a plain re-push stays rejected), and
// leave A's edit — not B's — on origin. Seeded by the ideation spike harness.
func TestStateCommitHaltsOnSameEntityConflict(t *testing.T) {
	bare, workflowA, workflowB, stateBranch := twoHostStateWorkflow(t)
	checkoutB := filepath.Join(workflowB, ".spacedock-state")

	// A edits first-task's frontmatter and pushes via the verb.
	writeEntity(t, workflowA, "first-task", "---\nstatus: implementation\n---\n# First Task (A)\n")
	hostA := filepath.Dir(filepath.Dir(workflowA))
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "-m", "A: -> implementation"); code != 0 {
		t.Fatalf("A's commit should succeed (exit 0); got exit=%d stderr=%q", code, errOut)
	}

	// B edits the SAME entity concurrently; the verb must HALT (exit 3).
	writeEntity(t, workflowB, "first-task", "---\nstatus: review\n---\n# First Task (B)\n")
	hostB := filepath.Dir(filepath.Dir(workflowB))
	code, _, errOut := runStateCommitCmd(t, hostB, workflowB, "first-task", "-m", "B: -> review")
	if code != 3 {
		t.Fatalf("same-entity conflict must HALT with exit 3; got exit=%d stderr=%q", code, errOut)
	}
	// stderr names the conflicting entity path.
	if !strings.Contains(errOut, "first-task.md") {
		t.Fatalf("HALT stderr should name the conflicting entity path; got:\n%s", errOut)
	}
	// The checkout is clean: rebase aborted, no rebase-in-progress dir.
	if porcelain := git(t, checkoutB, "status", "--porcelain"); strings.TrimSpace(porcelain) != "" {
		t.Fatalf("HALT must leave a clean checkout; porcelain:\n%s", porcelain)
	}
	gitDir := strings.TrimSpace(git(t, checkoutB, "rev-parse", "--git-path", "rebase-merge"))
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(checkoutB, gitDir)
	}
	if _, err := os.Stat(gitDir); err == nil {
		t.Fatalf("rebase still in progress after HALT (rebase-merge present at %s)", gitDir)
	}
	// B did NOT force-push — a plain push stays rejected, so A's edit survives.
	if _, ok := gitOK(t, checkoutB, "push", "origin", stateBranch); ok {
		t.Fatalf("a plain push after HALT must stay rejected; B must not have force-pushed")
	}
	// A's edit is what's on origin — no silent clobber of the peer's frontmatter.
	originFirst := showOriginFile(t, bare, stateBranch, "first-task.md")
	if !strings.Contains(originFirst, "status: implementation") {
		t.Fatalf("origin first-task should carry A's edit (status: implementation); got:\n%s", originFirst)
	}
}

// TestStateCommitHaltStderrCarriesPeerCommit pins AC-2 (D1): the exit-3 HALT
// stderr carries the peer commit that survived the aborted rebase — a populated,
// computed diagnostic rather than only an exit code. The peer commit is A's
// pushed HEAD sha (the pull's fetch phase updates origin/{branch} before the
// rebase conflicts; abort does not touch it).
func TestStateCommitHaltStderrCarriesPeerCommit(t *testing.T) {
	_, workflowA, workflowB, _ := twoHostStateWorkflow(t)
	checkoutA := filepath.Join(workflowA, ".spacedock-state")
	hostA := filepath.Dir(filepath.Dir(workflowA))

	writeEntity(t, workflowA, "first-task", "---\nstatus: implementation\n---\n# First Task (A)\n")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "-m", "A: -> implementation"); code != 0 {
		t.Fatalf("A's commit should succeed (exit 0); got exit=%d stderr=%q", code, errOut)
	}
	peerSHA := strings.TrimSpace(git(t, checkoutA, "rev-parse", "--short", "HEAD"))

	writeEntity(t, workflowB, "first-task", "---\nstatus: review\n---\n# First Task (B)\n")
	hostB := filepath.Dir(filepath.Dir(workflowB))
	code, _, errOut := runStateCommitCmd(t, hostB, workflowB, "first-task", "-m", "B: -> review")
	if code != 3 {
		t.Fatalf("same-entity conflict must HALT with exit 3; got exit=%d stderr=%q", code, errOut)
	}
	if !strings.Contains(errOut, "Peer commit: "+peerSHA) {
		t.Fatalf("HALT stderr should name the peer commit %q, got:\n%s", peerSHA, errOut)
	}
}

// TestStateCommitHaltJSONCarriesPeerCommit pins AC-2's --json requirement: the
// halt envelope carries peer_commit alongside the existing conflicting_paths.
func TestStateCommitHaltJSONCarriesPeerCommit(t *testing.T) {
	_, workflowA, workflowB, _ := twoHostStateWorkflow(t)
	checkoutA := filepath.Join(workflowA, ".spacedock-state")
	hostA := filepath.Dir(filepath.Dir(workflowA))

	writeEntity(t, workflowA, "first-task", "---\nstatus: implementation\n---\n# First Task (A)\n")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "-m", "A: -> implementation"); code != 0 {
		t.Fatalf("A's commit should succeed (exit 0); got exit=%d stderr=%q", code, errOut)
	}
	peerSHA := strings.TrimSpace(git(t, checkoutA, "rev-parse", "--short", "HEAD"))

	writeEntity(t, workflowB, "first-task", "---\nstatus: review\n---\n# First Task (B)\n")
	hostB := filepath.Dir(filepath.Dir(workflowB))
	code, stdout, errOut := runStateCommitCmd(t, hostB, workflowB, "first-task", "-m", "B: -> review", "--json")
	if code != 3 {
		t.Fatalf("same-entity conflict must HALT with exit 3; got exit=%d stderr=%q", code, errOut)
	}
	if !strings.Contains(stdout, `"peer_commit": "`+peerSHA+`"`) {
		t.Fatalf("--json halt envelope should carry peer_commit=%q, got:\n%s", peerSHA, stdout)
	}
}

// TestStateCommitIsPathScoped pins AC-2: a sibling dirty/untracked file in the
// state checkout is NOT swept into the commit (the verb stages exactly the entity,
// never `add -A`). This is the w4 2/3 `cd && git add -A` drift the verb deletes.
func TestStateCommitIsPathScoped(t *testing.T) {
	_, workflowA, _, _ := twoHostStateWorkflow(t)
	checkoutA := filepath.Join(workflowA, ".spacedock-state")
	hostA := filepath.Dir(filepath.Dir(workflowA))

	// Edit the entity AND drop a sibling untracked file in the same checkout.
	writeEntity(t, workflowA, "first-task", "---\nstatus: implementation\n---\n# First Task (A)\n")
	if err := os.WriteFile(filepath.Join(checkoutA, "sibling-junk.md"), []byte("untracked sibling\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "-m", "A: scoped"); code != 0 {
		t.Fatalf("commit should succeed; got exit=%d stderr=%q", code, errOut)
	}
	// The commit lists ONLY the entity path.
	names := git(t, checkoutA, "show", "--name-only", "--pretty=format:", "HEAD")
	if !strings.Contains(names, "first-task.md") {
		t.Fatalf("commit should include first-task.md; name-only:\n%s", names)
	}
	if strings.Contains(names, "sibling-junk.md") {
		t.Fatalf("path-scoped commit must NOT sweep the sibling; name-only:\n%s", names)
	}
	// The sibling stays untracked.
	if porcelain := git(t, checkoutA, "status", "--porcelain"); !strings.Contains(porcelain, "sibling-junk.md") {
		t.Fatalf("sibling should remain untracked after the scoped commit; porcelain:\n%s", porcelain)
	}
}

// TestStateCommitMultiWriterHappyPath pins AC-3: two writers commit DIFFERENT
// entities; A pushes first; the verb runs as B, its push is rejected non-ff, it
// pull --rebases the disjoint peer commit and re-pushes (exit 0) with both entities
// present and history linear (no merge commit).
func TestStateCommitMultiWriterHappyPath(t *testing.T) {
	bare, workflowA, workflowB, stateBranch := twoHostStateWorkflow(t)
	checkoutB := filepath.Join(workflowB, ".spacedock-state")
	hostA := filepath.Dir(filepath.Dir(workflowA))
	hostB := filepath.Dir(filepath.Dir(workflowB))

	// A commits a distinct NEW entity and pushes first.
	writeEntity(t, workflowA, "alpha-task", "---\nstatus: ideation\n---\n# Alpha (A)\n")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "alpha-task", "-m", "A: add alpha"); code != 0 {
		t.Fatalf("A's commit should succeed; exit=%d stderr=%q", code, errOut)
	}

	// B commits a DIFFERENT new entity — push rejected non-ff → pull --rebase → re-push.
	writeEntity(t, workflowB, "beta-task", "---\nstatus: ideation\n---\n# Beta (B)\n")
	code, _, errOut := runStateCommitCmd(t, hostB, workflowB, "beta-task", "-m", "B: add beta")
	if code != 0 {
		t.Fatalf("B's commit should succeed via pull --rebase + re-push (exit 0); got exit=%d stderr=%q", code, errOut)
	}
	// Both entities present in B's tree after the rebase.
	for _, slug := range []string{"alpha-task", "beta-task"} {
		if _, err := os.Stat(filepath.Join(checkoutB, slug+".md")); err != nil {
			t.Fatalf("B tree missing %s after commit: %v", slug, err)
		}
	}
	// Linear history: no merge commit on the state branch.
	if merges := git(t, checkoutB, "log", "--merges", "--oneline"); strings.TrimSpace(merges) != "" {
		t.Fatalf("history not linear; merge commits:\n%s", merges)
	}
	// B's entity reached origin.
	originBeta := showOriginFile(t, bare, stateBranch, "beta-task.md")
	if !strings.Contains(originBeta, "Beta (B)") {
		t.Fatalf("origin should carry B's beta-task; got:\n%s", originBeta)
	}
}

// TestStateCommitNoOriginLocalOnly pins AC-4: in a no-origin state checkout the verb
// commits path-scoped locally and reports local-only success (exit 0, --json result
// "local-only") without attempting push/pull.
func TestStateCommitNoOriginLocalOnly(t *testing.T) {
	_, workflowA, _, _ := twoHostStateWorkflow(t)
	checkoutA := filepath.Join(workflowA, ".spacedock-state")
	hostA := filepath.Dir(filepath.Dir(workflowA))

	// Drop the origin remote from the state checkout (linked worktree shares the
	// host's config, so remove origin on the host repo's config).
	git(t, checkoutA, "remote", "remove", "origin")
	headBefore := strings.TrimSpace(git(t, checkoutA, "rev-parse", "HEAD"))

	writeEntity(t, workflowA, "first-task", "---\nstatus: implementation\n---\n# First Task (A)\n")
	code, stdout, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "-m", "A: local", "--json")
	if code != 0 {
		t.Fatalf("no-origin commit should succeed local-only (exit 0); got exit=%d stderr=%q", code, errOut)
	}
	if !strings.Contains(stdout, `"result": "local-only"`) {
		t.Fatalf("no-origin commit should report result local-only; json:\n%s", stdout)
	}
	// The commit landed locally (HEAD advanced).
	headAfter := strings.TrimSpace(git(t, checkoutA, "rev-parse", "HEAD"))
	if headAfter == headBefore {
		t.Fatalf("no-origin commit should advance local HEAD; before=%s after=%s", headBefore, headAfter)
	}
}

// TestStateCommitNoOpWhenClean pins the clean no-op: committing a slug with no
// pending change is exit 0 with result "no-op", not an error.
func TestStateCommitNoOpWhenClean(t *testing.T) {
	_, workflowA, _, _ := twoHostStateWorkflow(t)
	hostA := filepath.Dir(filepath.Dir(workflowA))

	code, stdout, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "--json")
	if code != 0 {
		t.Fatalf("clean no-op commit should be exit 0; got exit=%d stderr=%q", code, errOut)
	}
	if !strings.Contains(stdout, `"result": "no-op"`) {
		t.Fatalf("clean commit should report result no-op; json:\n%s", stdout)
	}
}
