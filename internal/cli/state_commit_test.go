// ABOUTME: Real-git e2e for `state commit <slug>` — the rebase-HALT (AC-1), path-
// ABOUTME: scoped commit (AC-2), multi-writer happy path (AC-3), no-origin (AC-4).
package cli

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
	"github.com/spacedock-dev/spacedock/internal/testgit"
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

// writeFolderFile writes a file below a folder-form entity in the workflow's
// state checkout, creating parent directories as needed.
func writeFolderFile(t *testing.T, workflowDir, slug, rel, body string) {
	t.Helper()
	path := filepath.Join(workflowDir, ".spacedock-state", slug, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
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

	if code, out, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "-m", "A: scoped"); code != 0 {
		t.Fatalf("commit should succeed; got exit=%d stderr=%q", code, errOut)
	} else if want := "Committed and pushed first-task to spacedock-state/dev.\n"; out != want {
		t.Fatalf("direct-push prose changed: got=%q want=%q", out, want)
	}
	// The commit lists ONLY the entity path.
	names := strings.Fields(git(t, checkoutA, "show", "--name-only", "--pretty=format:", "HEAD"))
	if len(names) != 1 || names[0] != "first-task.md" {
		t.Fatalf("flat-form commit should contain exactly first-task.md; names=%q", names)
	}
	// The sibling stays untracked.
	if porcelain := git(t, checkoutA, "status", "--porcelain"); !strings.Contains(porcelain, "sibling-junk.md") {
		t.Fatalf("sibling should remain untracked after the scoped commit; porcelain:\n%s", porcelain)
	}
}

func TestStateCommitFlatIncludesExactCompanionDirectoryAndTrackedDeletions(t *testing.T) {
	_, workflowA, _, _ := twoHostStateWorkflow(t)
	checkout := filepath.Join(workflowA, ".spacedock-state")
	host := filepath.Dir(filepath.Dir(workflowA))
	const slug = "first-task"

	writeEntity(t, workflowA, slug, "---\nstatus: implementation\n---\n# Flat with room\n")
	companion := filepath.Join(checkout, slug, "review", "validation", "briefing-1")
	if err := os.MkdirAll(companion, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(companion, "request.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(checkout, "sibling-junk.md"), []byte("sibling\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code, _, errOut := runStateCommitCmd(t, host, workflowA, slug, "-m", "flat room unit"); code != 0 {
		t.Fatalf("flat companion commit exit=%d stderr=%q", code, errOut)
	}
	names := strings.Fields(git(t, checkout, "show", "--name-only", "--pretty=format:", "HEAD"))
	want := []string{"first-task.md", "first-task/review/validation/briefing-1/request.json"}
	if strings.Join(names, "\n") != strings.Join(want, "\n") {
		t.Fatalf("flat companion commit paths=%q want %q", names, want)
	}
	if porcelain := git(t, checkout, "status", "--porcelain"); !strings.Contains(porcelain, "sibling-junk.md") {
		t.Fatalf("sibling dirt was swept or lost:\n%s", porcelain)
	}

	if err := os.RemoveAll(filepath.Join(checkout, slug)); err != nil {
		t.Fatal(err)
	}
	writeEntity(t, workflowA, slug, "---\nstatus: validation\n---\n# Flat after provider\n")
	if code, _, errOut := runStateCommitCmd(t, host, workflowA, slug, "-m", "flat room deletion"); code != 0 {
		t.Fatalf("flat tracked deletion commit exit=%d stderr=%q", code, errOut)
	}
	names = strings.Fields(git(t, checkout, "show", "--name-only", "--pretty=format:", "HEAD"))
	if strings.Join(names, "\n") != strings.Join(want, "\n") {
		t.Fatalf("flat deletion commit paths=%q want %q", names, want)
	}
}

func TestStateCommitFlatCompanionRenameIncludesSourceAndDestination(t *testing.T) {
	_, workflow, _, _ := twoHostStateWorkflow(t)
	checkout := filepath.Join(workflow, ".spacedock-state")
	host := filepath.Dir(filepath.Dir(workflow))
	const slug = "first-task"

	oldPath := filepath.Join(checkout, slug, "review", "old.json")
	newPath := filepath.Join(checkout, slug, "review", "new.json")
	if err := os.MkdirAll(filepath.Dir(oldPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(oldPath, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code, _, errOut := runStateCommitCmd(t, host, workflow, slug, "-m", "seed flat companion"); code != 0 {
		t.Fatalf("seed companion exit=%d stderr=%q", code, errOut)
	}
	if err := os.Rename(oldPath, newPath); err != nil {
		t.Fatal(err)
	}
	if code, _, errOut := runStateCommitCmd(t, host, workflow, slug, "-m", "rename flat companion"); code != 0 {
		t.Fatalf("rename companion exit=%d stderr=%q", code, errOut)
	}
	nameStatus := strings.Fields(git(t, checkout, "show", "--name-status", "--no-renames", "--pretty=format:", "HEAD"))
	want := []string{"A", "first-task/review/new.json", "D", "first-task/review/old.json"}
	if strings.Join(nameStatus, "\n") != strings.Join(want, "\n") {
		t.Fatalf("rename commit paths=%q want %q", nameStatus, want)
	}
	if staged := strings.TrimSpace(git(t, checkout, "diff", "--cached", "--name-only")); staged != "" {
		t.Fatalf("rename left source staged outside commit: %q", staged)
	}
}

// TestStateCommitFolderIncludesWholeEntity pins AC-1 through AC-3: a folder-form
// entity is one commit unit. Its index, tracked report modification, tracked
// deletion, and untracked artifact land together while flat/folder siblings and
// unrelated top-level dirt remain untouched. Artifact-only dirt is not a false
// no-op, and a clean rerun is the existing no-op.
func TestStateCommitFolderIncludesWholeEntity(t *testing.T) {
	_, workflowA, _, _ := twoHostStateWorkflow(t)
	checkout := filepath.Join(workflowA, ".spacedock-state")
	host := filepath.Dir(filepath.Dir(workflowA))
	const slug = "folder-task"

	writeFolderFile(t, workflowA, slug, "index.md", "---\nstatus: ideation\n---\n# Folder\n")
	writeFolderFile(t, workflowA, slug, "reports/review.md", "baseline report\n")
	writeFolderFile(t, workflowA, slug, "artifacts/obsolete.md", "remove me\n")
	git(t, checkout, "add", "--", slug)
	git(t, checkout, "commit", "-q", "-m", "seed folder entity", "--", slug)
	git(t, checkout, "push", "-q", "origin", "HEAD")

	writeFolderFile(t, workflowA, slug, "index.md", "---\nstatus: implementation\n---\n# Folder\n")
	writeFolderFile(t, workflowA, slug, "reports/review.md", "updated tracked report\n")
	writeFolderFile(t, workflowA, slug, "artifacts/evidence.md", "new evidence\n")
	if err := os.Remove(filepath.Join(checkout, slug, "artifacts", "obsolete.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(checkout, "flat-sibling.md"), []byte("untracked flat sibling\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	writeFolderFile(t, workflowA, "folder-sibling", "index.md", "untracked folder sibling\n")
	if err := os.WriteFile(filepath.Join(checkout, "unrelated.txt"), []byte("untracked top-level path\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if code, _, errOut := runStateCommitCmd(t, host, workflowA, slug, "-m", "folder: scoped"); code != 0 {
		t.Fatalf("folder commit should succeed; exit=%d stderr=%q", code, errOut)
	}
	wantNames := []string{
		"folder-task/artifacts/evidence.md",
		"folder-task/artifacts/obsolete.md",
		"folder-task/index.md",
		"folder-task/reports/review.md",
	}
	gotNames := strings.Fields(git(t, checkout, "show", "--name-only", "--pretty=format:", "HEAD"))
	if strings.Join(gotNames, "\n") != strings.Join(wantNames, "\n") {
		t.Fatalf("folder commit paths mismatch\nwant: %q\n got: %q", wantNames, gotNames)
	}
	porcelain := git(t, checkout, "status", "--porcelain")
	for _, dirty := range []string{"flat-sibling.md", "folder-sibling/", "unrelated.txt"} {
		if !strings.Contains(porcelain, dirty) {
			t.Fatalf("sibling dirt %q should remain after scoped commit; porcelain:\n%s", dirty, porcelain)
		}
	}
	if strings.Contains(porcelain, slug+"/") {
		t.Fatalf("target folder should be clean after scoped commit; porcelain:\n%s", porcelain)
	}

	writeFolderFile(t, workflowA, slug, "artifacts/evidence.md", "artifact-only update\n")
	headBefore := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD"))
	if code, _, errOut := runStateCommitCmd(t, host, workflowA, slug, "-m", "folder: artifact only"); code != 0 {
		t.Fatalf("artifact-only commit should succeed; exit=%d stderr=%q", code, errOut)
	}
	if headAfter := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD")); headAfter == headBefore {
		t.Fatal("artifact-only folder dirt must advance HEAD, not return a false no-op")
	}
	if got := strings.TrimSpace(git(t, checkout, "show", "--name-only", "--pretty=format:", "HEAD")); got != slug+"/artifacts/evidence.md" {
		t.Fatalf("artifact-only commit should contain exactly the artifact; got %q", got)
	}
	code, stdout, errOut := runStateCommitCmd(t, host, workflowA, slug, "--json")
	if code != 0 || !strings.Contains(stdout, `"result": "no-op"`) {
		t.Fatalf("clean rerun should be no-op; exit=%d stdout=%q stderr=%q", code, stdout, errOut)
	}
}

// TestStateCommitFolderDeletion pins the Roborev deletion finding: the tracked
// index identifies a folder entity even after it disappears from the worktree,
// so deleting only index.md or the complete folder remains committable.
func TestStateCommitFolderDeletion(t *testing.T) {
	for _, tc := range []struct {
		name      string
		removeAll bool
		wantNames []string
	}{
		{name: "index", wantNames: []string{"deleted-folder/index.md"}},
		{name: "complete folder", removeAll: true, wantNames: []string{"deleted-folder/artifacts/evidence.md", "deleted-folder/index.md"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, workflow, _, _ := twoHostStateWorkflow(t)
			checkout := filepath.Join(workflow, ".spacedock-state")
			host := filepath.Dir(filepath.Dir(workflow))
			const slug = "deleted-folder"

			writeFolderFile(t, workflow, slug, "index.md", "---\nstatus: implementation\n---\n# Deleted Folder\n")
			writeFolderFile(t, workflow, slug, "artifacts/evidence.md", "tracked evidence\n")
			git(t, checkout, "add", "--", slug)
			git(t, checkout, "commit", "-q", "-m", "seed deletable folder", "--", slug)
			git(t, checkout, "push", "-q", "origin", "HEAD")

			if tc.removeAll {
				if err := os.RemoveAll(filepath.Join(checkout, slug)); err != nil {
					t.Fatal(err)
				}
			} else if err := os.Remove(filepath.Join(checkout, slug, "index.md")); err != nil {
				t.Fatal(err)
			}
			if code, _, errOut := runStateCommitCmd(t, host, workflow, slug, "-m", "delete folder paths"); code != 0 {
				t.Fatalf("tracked folder deletion should commit; exit=%d stderr=%q", code, errOut)
			}
			gotNames := strings.Fields(git(t, checkout, "show", "--name-only", "--pretty=format:", "HEAD"))
			if strings.Join(gotNames, "\n") != strings.Join(tc.wantNames, "\n") {
				t.Fatalf("folder deletion paths mismatch\nwant: %q\n got: %q", tc.wantNames, gotNames)
			}
			if porcelain := strings.TrimSpace(git(t, checkout, "status", "--porcelain", "--", slug)); porcelain != "" {
				t.Fatalf("folder deletion should leave target clean; porcelain=%q", porcelain)
			}
		})
	}
}

// TestStateCommitFlatDeletion retains the exact flat-file commit unit when the
// tracked file itself has been deleted from the worktree.
func TestStateCommitFlatDeletion(t *testing.T) {
	_, workflow, _, _ := twoHostStateWorkflow(t)
	checkout := filepath.Join(workflow, ".spacedock-state")
	host := filepath.Dir(filepath.Dir(workflow))
	if err := os.Remove(filepath.Join(checkout, "first-task.md")); err != nil {
		t.Fatal(err)
	}
	if code, _, errOut := runStateCommitCmd(t, host, workflow, "first-task", "-m", "delete flat entity"); code != 0 {
		t.Fatalf("tracked flat deletion should commit; exit=%d stderr=%q", code, errOut)
	}
	if got := strings.TrimSpace(git(t, checkout, "show", "--name-only", "--pretty=format:", "HEAD")); got != "first-task.md" {
		t.Fatalf("flat deletion should commit exactly first-task.md; got %q", got)
	}
}

// TestStateCommitTreatsSlugAsLiteralGitPathspec pins Roborev job 537: Git
// pathspec metacharacters are valid filename characters, not permission to sweep
// matching sibling entities or resolve a nonexistent wildcard alias.
func TestStateCommitTreatsSlugAsLiteralGitPathspec(t *testing.T) {
	t.Run("existing metacharacter slug", func(t *testing.T) {
		bare, workflow, _, stateBranch := twoHostStateWorkflow(t)
		checkout := filepath.Join(workflow, ".spacedock-state")
		host := filepath.Dir(filepath.Dir(workflow))
		const slug = ":(glob)scope*"
		const target = ":(glob)scope*.md"
		const trackedSibling = "scope-leak.md"
		const untrackedSibling = "scope-new.md"

		for path, body := range map[string]string{
			target:         "---\nstatus: ideation\n---\n# Literal Star\n",
			trackedSibling: "---\nstatus: ideation\n---\n# Tracked Sibling\n",
		} {
			if err := os.WriteFile(filepath.Join(checkout, path), []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		git(t, checkout, "--literal-pathspecs", "add", "--", target, trackedSibling)
		git(t, checkout, "--literal-pathspecs", "commit", "-q", "-m", "seed literal pathspec", "--", target, trackedSibling)
		git(t, checkout, "push", "-q", "origin", stateBranch)

		if err := os.WriteFile(filepath.Join(checkout, target), []byte("---\nstatus: implementation\n---\n# Literal Star\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(checkout, trackedSibling), []byte("dirty tracked sibling\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(checkout, untrackedSibling), []byte("dirty untracked sibling\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		if code, _, errOut := runStateCommitCmd(t, host, workflow, slug, "-m", "literal metacharacter slug"); code != 0 {
			t.Fatalf("metacharacter slug should commit literally; exit=%d stderr=%q", code, errOut)
		}
		if got := strings.TrimSpace(git(t, checkout, "show", "--name-only", "--pretty=format:", "HEAD")); got != target {
			t.Fatalf("metacharacter slug should commit exactly %q; got %q", target, got)
		}
		porcelain := git(t, checkout, "status", "--porcelain", "--untracked-files=all")
		for _, sibling := range []string{trackedSibling, untrackedSibling} {
			if !strings.Contains(porcelain, sibling) {
				t.Fatalf("matching sibling %q should remain dirty; porcelain:\n%s", sibling, porcelain)
			}
		}
		if got := showOriginFile(t, bare, stateBranch, trackedSibling); !strings.Contains(got, "Tracked Sibling") {
			t.Fatalf("origin tracked sibling should retain baseline; got %q", got)
		}
		if _, ok := gitOK(t, bare, "cat-file", "-e", stateBranch+":"+untrackedSibling); ok {
			t.Fatalf("untracked matching sibling %q must not reach origin", untrackedSibling)
		}
	})

	t.Run("nonexistent wildcard alias", func(t *testing.T) {
		bare, workflow, _, stateBranch := twoHostStateWorkflow(t)
		checkout := filepath.Join(workflow, ".spacedock-state")
		host := filepath.Dir(filepath.Dir(workflow))
		writeEntity(t, workflow, "wildcard-match", "---\nstatus: ideation\n---\n# Sibling\n")
		git(t, checkout, "add", "--", "wildcard-match.md")
		git(t, checkout, "commit", "-q", "-m", "seed wildcard sibling", "--", "wildcard-match.md")
		git(t, checkout, "push", "-q", "origin", stateBranch)
		writeEntity(t, workflow, "wildcard-match", "---\nstatus: implementation\n---\n# Dirty Sibling\n")
		headBefore := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD"))
		originBefore := strings.TrimSpace(git(t, bare, "rev-parse", stateBranch))

		code, _, errOut := runStateCommitCmd(t, host, workflow, ":(glob)wildcard*")
		if code == 0 || !strings.Contains(errOut, "no entity") {
			t.Fatalf("nonexistent wildcard alias should not resolve; exit=%d stderr=%q", code, errOut)
		}
		if got := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD")); got != headBefore {
			t.Fatalf("wildcard alias changed HEAD: before=%s after=%s", headBefore, got)
		}
		if got := strings.TrimSpace(git(t, bare, "rev-parse", stateBranch)); got != originBefore {
			t.Fatalf("wildcard alias changed origin: before=%s after=%s", originBefore, got)
		}
		if porcelain := git(t, checkout, "status", "--porcelain"); !strings.Contains(porcelain, "wildcard-match.md") {
			t.Fatalf("matching sibling should remain dirty; porcelain:\n%s", porcelain)
		}
	})
}

// TestStateCommitFolderMultiWriterHappyPath pins AC-5's disjoint-entity case:
// folder entities (including their artifacts) remain separate rebase units.
func TestStateCommitFolderMultiWriterHappyPath(t *testing.T) {
	bare, workflowA, workflowB, stateBranch := twoHostStateWorkflow(t)
	hostA := filepath.Dir(filepath.Dir(workflowA))
	hostB := filepath.Dir(filepath.Dir(workflowB))

	writeFolderFile(t, workflowA, "alpha-folder", "index.md", "---\nstatus: ideation\n---\n# Alpha\n")
	writeFolderFile(t, workflowA, "alpha-folder", "artifacts/a.md", "alpha evidence\n")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "alpha-folder", "-m", "A: add folder"); code != 0 {
		t.Fatalf("A folder commit should succeed; exit=%d stderr=%q", code, errOut)
	}

	writeFolderFile(t, workflowB, "beta-folder", "index.md", "---\nstatus: ideation\n---\n# Beta\n")
	writeFolderFile(t, workflowB, "beta-folder", "artifacts/b.md", "beta evidence\n")
	if code, _, errOut := runStateCommitCmd(t, hostB, workflowB, "beta-folder", "-m", "B: add folder"); code != 0 {
		t.Fatalf("B folder commit should rebase and push; exit=%d stderr=%q", code, errOut)
	}
	checkoutB := filepath.Join(workflowB, ".spacedock-state")
	if merges := strings.TrimSpace(git(t, checkoutB, "log", "--merges", "--oneline")); merges != "" {
		t.Fatalf("disjoint folder commits should retain linear history; merges:\n%s", merges)
	}
	for _, path := range []string{"alpha-folder/artifacts/a.md", "beta-folder/artifacts/b.md"} {
		if got := showOriginFile(t, bare, stateBranch, path); !strings.Contains(got, "evidence") {
			t.Fatalf("origin should contain %s; got %q", path, got)
		}
	}
}

// TestStateCommitFolderConflictHalts pins AC-5's same-folder case: concurrent
// edits to any shared nested path halt cleanly and name that path.
func TestStateCommitFolderConflictHalts(t *testing.T) {
	bare, workflowA, workflowB, stateBranch := twoHostStateWorkflow(t)
	checkoutA := filepath.Join(workflowA, ".spacedock-state")
	checkoutB := filepath.Join(workflowB, ".spacedock-state")
	hostA := filepath.Dir(filepath.Dir(workflowA))
	hostB := filepath.Dir(filepath.Dir(workflowB))
	const slug = "shared-folder"
	const report = "reports/review.md"

	writeFolderFile(t, workflowA, slug, "index.md", "---\nstatus: implementation\n---\n# Shared\n")
	writeFolderFile(t, workflowA, slug, report, "baseline\n")
	git(t, checkoutA, "add", "--", slug)
	git(t, checkoutA, "commit", "-q", "-m", "seed shared folder", "--", slug)
	git(t, checkoutA, "push", "-q", "origin", stateBranch)
	git(t, checkoutB, "pull", "-q", "--rebase", "origin", stateBranch)

	writeFolderFile(t, workflowA, slug, report, "host A\n")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, slug, "-m", "A: report"); code != 0 {
		t.Fatalf("A report commit should succeed; exit=%d stderr=%q", code, errOut)
	}
	writeFolderFile(t, workflowB, slug, report, "host B\n")
	code, _, errOut := runStateCommitCmd(t, hostB, workflowB, slug, "-m", "B: report")
	if code != 3 {
		t.Fatalf("same-folder nested conflict should HALT; exit=%d stderr=%q", code, errOut)
	}
	if !strings.Contains(errOut, slug+"/"+report) {
		t.Fatalf("HALT should name nested conflicting path; stderr:\n%s", errOut)
	}
	if porcelain := strings.TrimSpace(git(t, checkoutB, "status", "--porcelain")); porcelain != "" {
		t.Fatalf("HALT should abort rebase cleanly; porcelain=%q", porcelain)
	}
	localReport, err := os.ReadFile(filepath.Join(checkoutB, slug, report))
	if err != nil || string(localReport) != "host B\n" {
		t.Fatalf("HALT must preserve the local writer's nested edit; body=%q err=%v", localReport, err)
	}
	if got := showOriginFile(t, bare, stateBranch, slug+"/"+report); got != "host A\n" {
		t.Fatalf("origin must preserve peer's nested edit; got %q", got)
	}
}

// TestStateCommitRejectsNoncanonicalSlugWithoutSideEffects pins AC-6. The
// operand is one top-level entity slug, never an absolute/traversal/nested path.
func TestStateCommitRejectsNoncanonicalSlugWithoutSideEffects(t *testing.T) {
	bare, workflowA, _, stateBranch := twoHostStateWorkflow(t)
	checkout := filepath.Join(workflowA, ".spacedock-state")
	host := filepath.Dir(filepath.Dir(workflowA))
	writeFolderFile(t, workflowA, "folder-task", "index.md", "---\nstatus: ideation\n---\n# Folder\n")
	writeFolderFile(t, workflowA, "folder-task", "artifacts/evidence.md", "evidence\n")
	git(t, checkout, "add", "--", "folder-task")
	git(t, checkout, "commit", "-q", "-m", "seed nested target", "--", "folder-task")
	git(t, checkout, "push", "-q", "origin", stateBranch)

	cases := []string{
		"folder-task/artifacts/evidence",
		`folder-task\artifacts\evidence`,
		"roborev-workflow-setup-skill/artifacts/roborev-setup-skill/SKILL",
		".", "..", "../folder-task", "/folder-task", filepath.Join(checkout, "folder-task"),
	}
	evidencePath := filepath.Join(checkout, "folder-task", "artifacts", "evidence.md")
	for _, slug := range cases {
		t.Run(slug, func(t *testing.T) {
			headBefore := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD"))
			indexBefore := git(t, checkout, "diff", "--cached", "--binary")
			worktreeBefore := git(t, checkout, "status", "--porcelain=v1", "--untracked-files=all")
			originBefore := strings.TrimSpace(git(t, bare, "rev-parse", stateBranch))
			bytesBefore, err := os.ReadFile(evidencePath)
			if err != nil {
				t.Fatal(err)
			}

			code, _, errOut := runStateCommitCmd(t, host, workflowA, slug)
			if code == 0 || !strings.Contains(errOut, "invalid entity slug") {
				t.Fatalf("noncanonical slug should fail clearly; slug=%q exit=%d stderr=%q", slug, code, errOut)
			}
			if got := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD")); got != headBefore {
				t.Fatalf("invalid slug changed HEAD: before=%s after=%s", headBefore, got)
			}
			if got := git(t, checkout, "diff", "--cached", "--binary"); got != indexBefore {
				t.Fatalf("invalid slug changed index\nbefore=%q\nafter=%q", indexBefore, got)
			}
			if got := git(t, checkout, "status", "--porcelain=v1", "--untracked-files=all"); got != worktreeBefore {
				t.Fatalf("invalid slug changed worktree\nbefore=%q\nafter=%q", worktreeBefore, got)
			}
			if got, err := os.ReadFile(evidencePath); err != nil || string(got) != string(bytesBefore) {
				t.Fatalf("invalid slug changed worktree bytes: before=%q after=%q err=%v", bytesBefore, got, err)
			}
			if got := strings.TrimSpace(git(t, bare, "rev-parse", stateBranch)); got != originBefore {
				t.Fatalf("invalid slug changed origin: before=%s after=%s", originBefore, got)
			}
		})
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
	code, stdout, errOut := runStateCommitCmd(t, hostB, workflowB, "beta-task", "-m", "B: add beta")
	if code != 0 {
		t.Fatalf("B's commit should succeed via pull --rebase + re-push (exit 0); got exit=%d stderr=%q", code, errOut)
	}
	if want := "Committed beta-task, integrated peers' state, and pushed to spacedock-state/dev.\n"; stdout != want {
		t.Fatalf("rebase-push prose changed: got=%q want=%q", stdout, want)
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
	result := decodeOneJSON(t, stdout)
	if result["result"] != "local-only" || result["reason"] != "Committed first-task locally; no origin remote — state is local-only until an origin is configured." {
		t.Fatalf("no-origin result/prose contract changed: %#v", result)
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

	code, stdout, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task")
	if code != 0 {
		t.Fatalf("clean no-op commit should be exit 0; got exit=%d stderr=%q", code, errOut)
	}
	if want := "Nothing to commit for first-task — state checkout already up to date.\n"; stdout != want {
		t.Fatalf("clean no-op prose changed: got=%q want=%q", stdout, want)
	}
}

func TestStateCommitCleanActiveIntegratesPeerWithoutClaimingCommit(t *testing.T) {
	_, workflowA, workflowB, _ := twoHostStateWorkflow(t)
	hostA := filepath.Dir(filepath.Dir(workflowA))
	hostB := filepath.Dir(filepath.Dir(workflowB))

	writeEntity(t, workflowA, "peer-task", "---\nstatus: ideation\n---\n# Peer\n")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "peer-task", "-m", "peer state"); code != 0 {
		t.Fatalf("peer commit: exit=%d stderr=%q", code, errOut)
	}

	code, stdout, errOut := runStateCommitCmd(t, hostB, workflowB, "first-task")
	if code != 0 {
		t.Fatalf("clean active peer integration: exit=%d stderr=%q", code, errOut)
	}
	if want := "Nothing to commit for first-task — integrated peers' state; checkout is up to date.\n"; stdout != want {
		t.Fatalf("clean peer-only sync claimed a local commit: got=%q want=%q", stdout, want)
	}
	if _, err := os.Stat(filepath.Join(workflowB, ".spacedock-state", "peer-task.md")); err != nil {
		t.Fatalf("clean peer-only sync did not integrate peer state: %v", err)
	}
}

func TestStateCommitCleanActiveResumesEarlierFailedPush(t *testing.T) {
	bare, workflow, _, branch := twoHostStateWorkflow(t)
	host := filepath.Dir(filepath.Dir(workflow))
	checkout := filepath.Join(workflow, ".spacedock-state")
	hooks := t.TempDir()
	if err := os.WriteFile(filepath.Join(hooks, "pre-push"), []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, checkout, "config", "core.hooksPath", hooks)
	writeEntity(t, workflow, "first-task", "---\nstatus: implementation\n---\n# Locally committed\n")

	if code, _, _ := runStateCommitCmd(t, host, workflow, "first-task", "-m", "local before network failure"); code == 0 {
		t.Fatal("pre-push hook should fail the first publication")
	}
	localHead := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD"))
	if origin := strings.TrimSpace(git(t, bare, "rev-parse", branch)); origin == localHead {
		t.Fatal("failed push unexpectedly moved origin")
	}
	git(t, checkout, "config", "--unset", "core.hooksPath")

	if code, stdout, errOut := runStateCommitCmd(t, host, workflow, "first-task", "--json"); code != 0 {
		t.Fatalf("clean active retry must publish outstanding history: exit=%d stdout=%q stderr=%q", code, stdout, errOut)
	} else if result := decodeOneJSON(t, stdout); result["result"] != "pushed" || result["reason"] != "Published previously committed state for first-task to spacedock-state/dev." {
		t.Fatalf("clean active retry must describe publication without claiming a new commit: %#v", result)
	}
	if head := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD")); head != localHead {
		t.Fatalf("clean active retry created another commit: before=%s after=%s", localHead, head)
	}
	if origin := strings.TrimSpace(git(t, bare, "rev-parse", branch)); origin != localHead {
		t.Fatalf("clean active retry left origin behind: origin=%s local=%s", origin, localHead)
	}
}

func TestStateCommitCleanActiveSyncFailureDoesNotClaimNewCommit(t *testing.T) {
	_, workflowA, workflowB, _ := twoHostStateWorkflow(t)
	hostA := filepath.Dir(filepath.Dir(workflowA))
	hostB := filepath.Dir(filepath.Dir(workflowB))

	writeEntity(t, workflowA, "peer-task", "---\nstatus: ideation\n---\n# Peer\n")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "peer-task", "-m", "peer state"); code != 0 {
		t.Fatalf("peer commit: exit=%d stderr=%q", code, errOut)
	}
	writeEntity(t, workflowB, "peer-task", "---\nstatus: validation\n---\n# Untracked sibling collision\n")

	code, _, errOut := runStateCommitCmd(t, hostB, workflowB, "first-task")
	if code == 0 {
		t.Fatal("dirty sibling must block peer integration after non-fast-forward push")
	}
	if !strings.Contains(errOut, "no new commit was created in this invocation") ||
		strings.Contains(errOut, "new local commit remains recoverable") ||
		strings.Contains(errOut, "existing archive commit remains recoverable") {
		t.Fatalf("clean active sync failure misstated commit recovery:\n%s", errOut)
	}
}

func TestArchivedTargetCleanSeparatesGitFailureFromDirt(t *testing.T) {
	if clean, detail, err := archivedTargetClean(t.TempDir(), []string{"_archive/task.md"}); clean || detail != "" || err == nil {
		t.Fatalf("non-Git checkout must return a Git inspection error, got clean=%v detail=%q err=%v", clean, detail, err)
	}
}

func TestStateCommitRefusesActiveArchiveAndArchiveShapeCollisions(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, checkout string)
	}{
		{
			name: "active and archived",
			setup: func(t *testing.T, checkout string) {
				if err := os.MkdirAll(filepath.Join(checkout, "_archive"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(checkout, "_archive", "first-task.md"), []byte("---\nstatus: done\n---\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "archived flat and folder",
			setup: func(t *testing.T, checkout string) {
				if err := os.Remove(filepath.Join(checkout, "first-task.md")); err != nil {
					t.Fatal(err)
				}
				if err := os.MkdirAll(filepath.Join(checkout, "_archive", "first-task"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(checkout, "_archive", "first-task.md"), []byte("---\nstatus: done\n---\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(checkout, "_archive", "first-task", "index.md"), []byte("---\nstatus: done\n---\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bare, workflow, _, branch := twoHostStateWorkflow(t)
			host := filepath.Dir(filepath.Dir(workflow))
			checkout := filepath.Join(workflow, ".spacedock-state")
			tt.setup(t, checkout)
			headBefore := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD"))
			statusBefore := git(t, checkout, "status", "--porcelain=v1", "--untracked-files=all")
			originBefore := strings.TrimSpace(git(t, bare, "rev-parse", branch))

			code, _, errOut := runStateCommitCmd(t, host, workflow, "first-task")
			if code == 0 || !strings.Contains(errOut, "collision") {
				t.Fatalf("invalid identity shape must refuse clearly: exit=%d stderr=%q", code, errOut)
			}
			if got := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD")); got != headBefore {
				t.Fatalf("collision moved HEAD: before=%s after=%s", headBefore, got)
			}
			if got := git(t, checkout, "status", "--porcelain=v1", "--untracked-files=all"); got != statusBefore {
				t.Fatalf("collision changed index/worktree:\nbefore=%q\nafter=%q", statusBefore, got)
			}
			if got := strings.TrimSpace(git(t, bare, "rev-parse", branch)); got != originBefore {
				t.Fatalf("collision moved origin: before=%s after=%s", originBefore, got)
			}
		})
	}
}

func TestStateCommitRefusesDirtyArchivedEntityBeforePublication(t *testing.T) {
	for _, kind := range []string{"staged", "unstaged", "untracked"} {
		t.Run(kind, func(t *testing.T) {
			bare, workflow, _, branch := twoHostStateWorkflow(t)
			host := filepath.Dir(filepath.Dir(workflow))
			checkout := filepath.Join(workflow, ".spacedock-state")
			const slug = "folder-task"

			writeFolderFile(t, workflow, slug, "index.md", "---\nstatus: implementation\n---\n# Folder\n")
			if code, _, errOut := runStateCommitCmd(t, host, workflow, slug, "-m", "seed folder"); code != 0 {
				t.Fatalf("seed folder: exit=%d stderr=%q", code, errOut)
			}
			if err := os.MkdirAll(filepath.Join(checkout, "_archive"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Rename(filepath.Join(checkout, slug), filepath.Join(checkout, "_archive", slug)); err != nil {
				t.Fatal(err)
			}
			git(t, checkout, "add", "--", slug, "_archive/"+slug)
			git(t, checkout, "commit", "-q", "-m", "archive "+slug+" (merge guard)", "--", slug, "_archive/"+slug)

			index := filepath.Join(checkout, "_archive", slug, "index.md")
			switch kind {
			case "staged":
				if err := os.WriteFile(index, []byte("staged archive dirt\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				git(t, checkout, "add", "--", "_archive/"+slug)
			case "unstaged":
				if err := os.WriteFile(index, []byte("unstaged archive dirt\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			case "untracked":
				if err := os.WriteFile(filepath.Join(checkout, "_archive", slug, "artifact.txt"), []byte("untracked\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			headBefore := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD"))
			statusBefore := git(t, checkout, "status", "--porcelain=v1", "--untracked-files=all")
			originBefore := strings.TrimSpace(git(t, bare, "rev-parse", branch))
			code, _, errOut := runStateCommitCmd(t, host, workflow, slug)
			if code == 0 || !strings.Contains(errOut, "archived entity") || !strings.Contains(errOut, "dirty") {
				t.Fatalf("%s archived dirt must refuse: exit=%d stderr=%q", kind, code, errOut)
			}
			if got := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD")); got != headBefore {
				t.Fatalf("%s archived dirt moved HEAD: before=%s after=%s", kind, headBefore, got)
			}
			if got := git(t, checkout, "status", "--porcelain=v1", "--untracked-files=all"); got != statusBefore {
				t.Fatalf("%s refusal changed index/worktree:\nbefore=%q\nafter=%q", kind, statusBefore, got)
			}
			if got := strings.TrimSpace(git(t, bare, "rev-parse", branch)); got != originBefore {
				t.Fatalf("%s archived dirt moved origin: before=%s after=%s", kind, originBefore, got)
			}

			git(t, checkout, "reset", "--hard", "HEAD")
			if kind == "untracked" {
				if err := os.Remove(filepath.Join(checkout, "_archive", slug, "artifact.txt")); err != nil {
					t.Fatal(err)
				}
			}
			if code, stdout, errOut := runStateCommitCmd(t, host, workflow, slug, "--json"); code != 0 || decodeOneJSON(t, stdout)["result"] != "pushed" {
				t.Fatalf("clean folder archive should publish: exit=%d stdout=%q stderr=%q", code, stdout, errOut)
			}
			if code, stdout, errOut := runStateCommitCmd(t, host, workflow, slug, "--json"); code != 0 || decodeOneJSON(t, stdout)["result"] != "no-op" {
				t.Fatalf("published folder archive should no-op: exit=%d stdout=%q stderr=%q", code, stdout, errOut)
			}
		})
	}
}

// inlineGatedWorkflowReadme is an INLINE workflow (no `state:` key — entities live
// beside the README) whose validation stage is gated, the shape the live
// rejection-flow journey runs against.
const inlineGatedWorkflowReadme = `---
commissioned-by: spacedock@1
entity-type: task
id-style: slug
stages:
  states:
    - name: backlog
      initial: true
    - name: implementation
    - name: validation
      gate: true
      feedback-to: implementation
    - name: done
      terminal: true
---

# Inline Gated Workflow
`

const inlineGatedEntity = `---
id: rejection-task
status: validation
score: "0.90"
---
# Rejection Task

## Stage Report: validation

- DONE: Verify the corrected candidate carries the fix marker.
  Cycle 2 re-review confirms the marker is present.

### Summary

Validation is complete and ready for its decision gate.
`

// inlineGatedWorkflow builds a committed inline workflow inside its own Git work
// tree, carrying one folder-form entity parked at a gated validation stage with a
// mechanically complete report. Returns the workflow dir.
func inlineGatedWorkflow(t *testing.T) string {
	t.Helper()
	_, workflow := inlineGatedWorkflowAt(t, ".")
	return workflow
}

// inlineGatedWorkflowAt builds the inline fixture with the workflow dir at
// workflowRel below the repo root, and returns both. workflowRel "." puts the
// workflow AT the repo root; a nested path like "docs/dev" is the ordinary
// project shape, where Git reports staged names relative to the repo root rather
// than to the workflow. Returns (repoRoot, workflowDir).
func inlineGatedWorkflowAt(t *testing.T, workflowRel string) (string, string) {
	t.Helper()
	repo := filepath.Join(t.TempDir(), "repo")
	workflow := filepath.Join(repo, filepath.FromSlash(workflowRel))
	writeFileWithDirs(t, filepath.Join(workflow, "README.md"), inlineGatedWorkflowReadme)
	writeFileWithDirs(t, filepath.Join(workflow, "rejection-task", "index.md"), inlineGatedEntity)
	writeFileWithDirs(t, filepath.Join(workflow, "rejection-task", "inputs", "gate-validation", "gate-review.md"),
		"# Rejection Task — validation review\n\nSeed review.\n")
	testgit.InitRepo(t, repo, "-q")
	git(t, repo, "add", "-A")
	git(t, repo, "commit", "-q", "-m", "seed inline gated workflow")
	return repo, workflow
}

// runInlineCmd runs an arbitrary spacedock command with the workflow dir as cwd.
func runInlineCmd(t *testing.T, workflow string, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errBuf strings.Builder
	code = run(context.Background(), args, os.Environ(), workflow, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	return code, out.String(), errBuf.String()
}

// inlineReadyGates returns the `status --next --json` ready_gates rows.
func inlineReadyGates(t *testing.T, workflow string) []map[string]string {
	t.Helper()
	code, stdout, errOut := runInlineCmd(t, workflow, "status", "--workflow-dir", workflow, "--next", "--json")
	if code != 0 {
		t.Fatalf("status --next exit=%d stderr=%q", code, errOut)
	}
	var next struct {
		ReadyGates []map[string]string `json:"ready_gates"`
	}
	if err := json.Unmarshal([]byte(stdout), &next); err != nil {
		t.Fatalf("parse status --next: %v\n%s", err, stdout)
	}
	return next.ReadyGates
}

func runInlineGatePrepare(t *testing.T, workflow, review string) (code int, stdout, stderr string) {
	t.Helper()
	return runInlineCmd(t, workflow, "gate", "prepare", "rejection-task", "--workflow-dir", workflow,
		"--question", "Does the second validation confirm the fix marker?",
		"--artifact", review,
		"--summary", "The corrected candidate is ready for a decision.")
}

// TestStateCommitInlineUnblocksGatePreparationChain pins the whole causal chain the
// live rejection flow walks after the neutral round recorder writes: an inline
// workflow's entity and its fresh gate-review artifact are uncommitted, which makes
// `status --next` drop the candidate row (dirty bytes are never mechanical proof —
// TestGateReadinessRejectsDirtyAndMalformedColdReports pins that deliberately) and
// makes `gate prepare` refuse the artifact. `state commit` is the only durability
// verb the First Officer contract names, so if it no-ops here the run is stranded
// with no verb that can clear either refusal. The pre-condition assertions are the
// falsifying baseline: they are the exact refusals observed in CI, and they must
// hold BEFORE the verb runs and be gone after it.
func TestStateCommitInlineUnblocksGatePreparationChain(t *testing.T) {
	workflow := inlineGatedWorkflow(t)
	entity := filepath.Join(workflow, "rejection-task", "index.md")
	review := filepath.Join(workflow, "rejection-task", "inputs", "gate-validation", "gate-review.md")
	entityDirt := func() string {
		return strings.TrimSpace(git(t, workflow, "status", "--porcelain", "--untracked-files=all", "--", "rejection-task"))
	}

	// The durable state the round recorder plus the FO's fresh gate review leave:
	// entity edited, review artifact rewritten, nothing committed.
	writeFile(t, entity, inlineGatedEntity+"\n### Feedback Cycles\n\n- Cycle 1: REJECTED — validation reviewer; surface 1 marker vs estimate 1 (100%); AC unchanged\n")
	writeFile(t, review, "# Rejection Task — validation review (cycle 2)\n\nThe corrected candidate carries the fix marker.\n")

	if entityDirt() == "" {
		t.Fatal("fixture must start dirty; the whole chain is about clearing that dirt")
	}
	if ready := inlineReadyGates(t, workflow); len(ready) != 0 {
		t.Fatalf("baseline broken: a dirty entity must not be offered as a gate candidate; got %v", ready)
	}
	if code, _, errOut := runInlineGatePrepare(t, workflow, review); code == 0 ||
		!strings.Contains(errOut, "commit the exact source before preparation") {
		t.Fatalf("baseline broken: gate prepare must refuse an uncommitted artifact; exit=%d stderr=%q", code, errOut)
	}

	code, stdout, errOut := runStateCommitCmd(t, workflow, workflow, "rejection-task", "--json")
	if code != 0 || decodeOneJSON(t, stdout)["result"] != "committed" {
		t.Fatalf("inline state commit must commit; exit=%d stdout=%q stderr=%q", code, stdout, errOut)
	}

	if dirt := entityDirt(); dirt != "" {
		t.Fatalf("inline state commit left the entity dirty:\n%s", dirt)
	}
	ready := inlineReadyGates(t, workflow)
	if len(ready) != 1 || ready[0]["slug"] != "rejection-task" || ready[0]["readiness"] != "needs-preparation" {
		t.Fatalf("ready_gates=%v, want exactly one needs-preparation row for rejection-task", ready)
	}
	if code, stdout, errOut := runInlineGatePrepare(t, workflow, review); code != 0 || !strings.Contains(stdout, "state=open") {
		t.Fatalf("gate prepare must succeed on the committed artifact; exit=%d stdout=%q stderr=%q", code, stdout, errOut)
	}
}

// TestStateCommitInlineScopeAndPublishBoundary pins the two boundaries the inline
// commit must not cross: it stages only the named entity's commit unit, and it
// never pushes — an inline workflow repo is the caller's own code repo, so a push
// would move a branch the caller never asked this verb to touch. A workflow dir
// outside a Git work tree stays an exit-0 no-op that names the reason instead of
// reporting a commit that did not happen.
func TestStateCommitInlineScopeAndPublishBoundary(t *testing.T) {
	workflow := inlineGatedWorkflow(t)
	// Give the inline repo an origin, so "never pushes" is observable rather than
	// vacuous: the bare ref must not move across a successful commit.
	bare := filepath.Join(t.TempDir(), "origin.git")
	git(t, t.TempDir(), "init", "-q", "--bare", bare)
	git(t, workflow, "remote", "add", "origin", bare)
	git(t, workflow, "push", "-q", "origin", "HEAD:refs/heads/main")
	originBefore := strings.TrimSpace(git(t, bare, "rev-parse", "refs/heads/main"))

	writeFileWithDirs(t, filepath.Join(workflow, "sibling-task", "index.md"),
		"---\nid: sibling-task\nstatus: backlog\n---\n# Sibling\n")
	git(t, workflow, "add", "-A")
	git(t, workflow, "commit", "-q", "-m", "seed sibling")

	writeFile(t, filepath.Join(workflow, "rejection-task", "index.md"), inlineGatedEntity+"\ntarget dirt\n")
	writeFile(t, filepath.Join(workflow, "sibling-task", "index.md"),
		"---\nid: sibling-task\nstatus: backlog\n---\n# Sibling\n\nsibling dirt\n")

	if code, stdout, errOut := runStateCommitCmd(t, workflow, workflow, "rejection-task", "--json"); code != 0 ||
		decodeOneJSON(t, stdout)["result"] != "committed" {
		t.Fatalf("inline state commit exit=%d stdout=%q stderr=%q", code, stdout, errOut)
	}
	if dirt := strings.TrimSpace(git(t, workflow, "status", "--porcelain", "--", "sibling-task")); dirt == "" {
		t.Fatal("inline state commit swept up a dirty sibling entity; the commit must be path-scoped")
	}
	if got := strings.TrimSpace(git(t, bare, "rev-parse", "refs/heads/main")); got != originBefore {
		t.Fatalf("inline state commit pushed the caller's own repo: origin main before=%s after=%s", originBefore, got)
	}

	// A second run has nothing left to stage: a no-op, not a claimed commit.
	if code, stdout, errOut := runStateCommitCmd(t, workflow, workflow, "rejection-task", "--json"); code != 0 ||
		decodeOneJSON(t, stdout)["result"] != "no-op" {
		t.Fatalf("clean inline entity should no-op; exit=%d stdout=%q stderr=%q", code, stdout, errOut)
	}

	// Outside a Git work tree the verb has nothing to commit to, and says so.
	plain := filepath.Join(t.TempDir(), "plain")
	writeFileWithDirs(t, filepath.Join(plain, "README.md"), inlineGatedWorkflowReadme)
	writeFileWithDirs(t, filepath.Join(plain, "rejection-task", "index.md"), inlineGatedEntity)
	code, stdout, errOut := runStateCommitCmd(t, plain, plain, "rejection-task", "--json")
	result := decodeOneJSON(t, stdout)
	if code != 0 || result["result"] != "no-op" ||
		!strings.Contains(result["reason"].(string), "not inside a Git work tree") {
		t.Fatalf("non-repo workflow dir must be an honest no-op; exit=%d stdout=%q stderr=%q", code, stdout, errOut)
	}
}

// TestStateCommitInlineCommitsWhateverTheWorkflowDepth pins the inline commit
// against BOTH repo shapes, because the seam reads staged names from Git and
// compares them to entity paths, and those two are relative to different roots
// unless the workflow dir happens to be the repo root.
//
// The nested case is the ordinary project shape and it was silently broken: `git
// -C <workflow> diff --cached --name-only` answers with repo-root-relative names
// (docs/dev/rejection-task/index.md) while the entity pathspec is
// workflow-relative (rejection-task), so nothing matched, nothing was selected to
// commit, and the verb reported "already up to date" — after its own `git add`
// had staged the change. That is worse than the exit-0 no-op this entity set out
// to remove: it mutates the index and still claims there was nothing to do.
//
// Asserting a clean porcelain rather than just a fresh commit is what makes this
// falsifying: the pre-fix bug left the entity STAGED, which a commit-count check
// alone would miss. Dropping `--relative` from the diff reds the nested case and
// leaves the root case passing.
func TestStateCommitInlineCommitsWhateverTheWorkflowDepth(t *testing.T) {
	for name, workflowRel := range map[string]string{
		"workflow at repo root": ".",
		"workflow nested":       "docs/dev",
	} {
		t.Run(name, func(t *testing.T) {
			repo, workflow := inlineGatedWorkflowAt(t, workflowRel)
			entity := filepath.Join(workflow, "rejection-task", "index.md")
			writeFile(t, entity, inlineGatedEntity+"\nrecorded round, not yet durable\n")

			before := strings.TrimSpace(git(t, repo, "rev-parse", "HEAD"))
			code, stdout, errOut := runStateCommitCmd(t, workflow, workflow, "rejection-task", "--json")
			result := decodeOneJSON(t, stdout)
			if code != 0 || result["result"] != "committed" {
				t.Fatalf("inline state commit must report a commit, not a no-op; exit=%d stdout=%q stderr=%q", code, stdout, errOut)
			}

			// Empty porcelain, not just a new commit: the pre-fix bug left the
			// entity staged, which is invisible to a HEAD-moved check alone.
			if dirt := strings.TrimSpace(git(t, repo, "status", "--porcelain", "--untracked-files=all")); dirt != "" {
				t.Fatalf("entity left uncommitted (staged or modified) after state commit:\n%s", dirt)
			}
			if after := strings.TrimSpace(git(t, repo, "rev-parse", "HEAD")); after == before {
				t.Fatal("state commit reported a commit but HEAD did not move")
			}

			// The commit carries the entity, addressed from the repo root.
			entityRel, err := filepath.Rel(repo, entity)
			if err != nil {
				t.Fatal(err)
			}
			files := git(t, repo, "show", "--name-only", "--format=", "HEAD")
			if !strings.Contains(files, filepath.ToSlash(entityRel)) {
				t.Fatalf("commit does not contain %s; it contains:\n%s", entityRel, files)
			}
		})
	}
}
