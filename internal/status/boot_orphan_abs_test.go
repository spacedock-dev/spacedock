// ABOUTME: ORPHANS absolute-worktree parity — an absolute `worktree:` value is
// ABOUTME: used as-is for the DIR_EXISTS probe, matching os.path.join semantics.
package status

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// bootOrphanReadme defines a single worktree-bearing stage so --boot renders an
// ORPHANS table for an entity that carries a worktree field.
const bootOrphanReadme = `---
commissioned-by: spacedock@1
id-style: slug
stages:
  states:
    - name: build
      initial: true
      worktree: true
---

# Orphan Boot Workflow
`

// orphanDirExists returns the DIR_EXISTS cell of the single ORPHANS data row in
// a --boot rendering, or "" if no ORPHANS table is present.
func orphanDirExists(boot string) string {
	dir, _ := orphanExistence(boot)
	return dir
}

// orphanExistence returns the (DIR_EXISTS, BRANCH_EXISTS) cells of the single
// ORPHANS data row in a --boot rendering, or ("","") if no ORPHANS table is
// present. The ORPHANS columns are ID SLUG WORKTREE DIR_EXISTS BRANCH_EXISTS;
// the fixtures use space-free worktree values so strings.Fields cleanly splits
// them.
func orphanExistence(boot string) (dir, branch string) {
	lines := strings.Split(boot, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "ORPHANS") && i+2 < len(lines) {
			fields := strings.Fields(lines[i+2])
			if len(fields) >= 5 {
				return fields[3], fields[4]
			}
		}
	}
	return "", ""
}

// TestBootAbsoluteWorktreeDirExists is the carried boot-fix parity test: an
// absolute `worktree:` value resolves to that absolute dir (os.path.join drops
// git_root), so DIR_EXISTS is yes when it exists — both directly and vs the
// oracle. The pre-fix native code joined the absolute path onto git_root,
// yielding a non-existent path and DIR_EXISTS=no.
func TestBootAbsoluteWorktreeDirExists(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("git", "-C", root, "init").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v (%s)", err, out)
	}
	// An absolute worktree path that exists but is OUTSIDE the git root, so a
	// filepath.Join(git_root, wt) would point at a non-existent nested path.
	absWorktree := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), bootOrphanReadme)
	writeFile(t, filepath.Join(root, "feature.md"),
		"---\nstatus: build\nworktree: "+absWorktree+"\n---\n")

	env := pinnedEnv(t)
	args := []string{"--workflow-dir", root, "--boot"}

	nativeOut, nativeErr, nativeCode := runNative(t, root, env, args...)
	if nativeCode != 0 {
		t.Fatalf("native --boot exit=%d stderr=%q", nativeCode, nativeErr)
	}
	if got := orphanDirExists(nativeOut); got != "yes" {
		t.Fatalf("DIR_EXISTS for absolute existing worktree = %q, want \"yes\"\n%s", got, nativeOut)
	}

	// Oracle parity. The oracle is resolved in-tree, so this comparison always
	// runs (and hard-fails on a real divergence) on top of the direct DIR_EXISTS
	// assertion above.
	oracleOut, _, oracleCode := runOracle(t, root, env, args...)
	if oracleCode != 0 {
		t.Fatalf("oracle --boot exit=%d", oracleCode)
	}
	if got, want := orphanDirExists(nativeOut), orphanDirExists(oracleOut); got != want {
		t.Fatalf("ORPHANS DIR_EXISTS native=%q oracle=%q\n--- native ---\n%s\n--- oracle ---\n%s",
			got, want, nativeOut, oracleOut)
	}
}

// gitC runs a git subcommand in dir, failing the test on a non-zero exit so a
// fixture-setup failure surfaces with the git output instead of a downstream
// mystery.
func gitC(t *testing.T, dir string, args ...string) {
	t.Helper()
	full := append([]string{"-C", dir}, args...)
	if out, err := exec.Command("git", full...).CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v (%s)", strings.Join(args, " "), err, out)
	}
}

// splitRootWorktreeReadme declares a worktree-bearing stage on a split-root
// workflow (state: .spacedock-state), so --boot renders an ORPHANS row for an
// entity carrying a relative worktree value, resolved against the code repo root.
const splitRootWorktreeReadme = `---
commissioned-by: spacedock@1
id-style: slug
state: .spacedock-state
stages:
  states:
    - name: build
      initial: true
      worktree: true
---

# Split-Root Orphan Boot Workflow
`

// buildSplitRootWorktreeFixture materializes the spike topology that reproduces
// the misreport: a code repo with its own .git, a docs/dev definition dir
// declaring state: .spacedock-state, the state checkout materialized as a REAL
// git worktree (so it carries a .git pointer file — the trap FindGitRoot fell
// into), a relative-worktree entity under the state checkout, and a present code
// worktree at <coderoot>/.worktrees/feature-wt. Returns (defDir, wtDir).
//
// The .git pointer on the state checkout is load-bearing: FindGitRoot(entityDir)
// stops there instead of walking up to the code repo root, which is exactly the
// wrong-root the fix corrects by resolving against definitionDir instead.
func buildSplitRootWorktreeFixture(t *testing.T) (defDir, wtDir string) {
	t.Helper()
	coderoot := t.TempDir()
	gitC(t, coderoot, "init")
	gitC(t, coderoot, "config", "user.email", "test@example.com")
	gitC(t, coderoot, "config", "user.name", "test")
	// A worktree needs a committed HEAD to branch from.
	writeFile(t, filepath.Join(coderoot, "seed.txt"), "seed\n")
	gitC(t, coderoot, "add", "seed.txt")
	gitC(t, coderoot, "commit", "-m", "seed")

	defDir = filepath.Join(coderoot, "docs", "dev")
	writeFile(t, filepath.Join(defDir, "README.md"), splitRootWorktreeReadme)

	// The state checkout is a real linked worktree, so it carries a .git pointer
	// file rather than a .git directory — matching the live split-root topology.
	stateDir := filepath.Join(defDir, ".spacedock-state")
	gitC(t, coderoot, "worktree", "add", "--detach", stateDir)
	writeFile(t, filepath.Join(stateDir, "feature", "index.md"),
		"---\nstatus: build\nworktree: .worktrees/feature-wt\n---\n")

	// The present code worktree the entity points at, under the code repo root.
	// It is checked out on its own branch (not --detach), so `git worktree list
	// --porcelain` emits a `branch` line for it — which is what the
	// BRANCH_EXISTS probe keys off.
	wtDir = filepath.Join(coderoot, ".worktrees", "feature-wt")
	gitC(t, coderoot, "worktree", "add", "-b", "feature-wt", wtDir)
	return defDir, wtDir
}

// TestBootSplitRootWorktreeExistence is T-1/AC-1: under a split-root workflow
// whose state checkout is a git worktree (own .git pointer), an entity with a
// relative worktree pointing at a present code worktree must report
// DIR_EXISTS=yes / BRANCH_EXISTS=yes, resolved against the code repo root
// (FindGitRoot(definitionDir)) — NOT the state checkout (FindGitRoot(entityDir),
// which yields .spacedock-state/.worktrees/feature-wt and reports no/no).
//
// This is the missing case that let the misreport ship: the existing
// absolute-worktree test cannot exercise the wrong-root because os.path.join
// discards the root for an absolute value. Against today's entityDir resolution
// this asserts no/no (red); after the entityDir->definitionDir switch it asserts
// yes/yes (green). Native-only: split-root is an intentional native divergence
// with no oracle, consistent with the other split-root tests.
func TestBootSplitRootWorktreeExistence(t *testing.T) {
	defDir, wtDir := buildSplitRootWorktreeFixture(t)
	env := pinnedEnv(t)
	args := []string{"--workflow-dir", defDir, "--boot"}

	out, stderr, code := runNative(t, defDir, env, args...)
	if code != 0 {
		t.Fatalf("native --boot exit=%d stderr=%q", code, stderr)
	}
	dir, branch := orphanExistence(out)
	if dir != "yes" || branch != "yes" {
		t.Fatalf("present split-root worktree DIR_EXISTS=%q BRANCH_EXISTS=%q, want yes/yes\n%s", dir, branch, out)
	}

	// Removed-worktree variant: tear down the code worktree and prune, then the
	// same entity must report DIR_EXISTS=no (and BRANCH_EXISTS=no, since the
	// pruned worktree drops out of `git worktree list`).
	if err := os.RemoveAll(wtDir); err != nil {
		t.Fatal(err)
	}
	gitC(t, defDir, "worktree", "prune")

	out2, stderr2, code2 := runNative(t, defDir, env, args...)
	if code2 != 0 {
		t.Fatalf("native --boot (removed) exit=%d stderr=%q", code2, stderr2)
	}
	if dir2, branch2 := orphanExistence(out2); dir2 != "no" || branch2 != "no" {
		t.Fatalf("removed split-root worktree DIR_EXISTS=%q BRANCH_EXISTS=%q, want no/no\n%s", dir2, branch2, out2)
	}
}

// TestBootSingleRootWorktreeExistence is T-2/AC-2: a single-root workflow (no
// state: field, so definitionDir == entityDir) with a relative worktree pointing
// at a present code worktree reports DIR_EXISTS=yes, and a non-existent one
// reports no. This asserts the byte-identical-when-roots-coincide property — it
// passes both before and after the entityDir->definitionDir switch (a no-op when
// the two roots are the same dir), confirming no regression.
func TestBootSingleRootWorktreeExistence(t *testing.T) {
	root := t.TempDir()
	gitC(t, root, "init")
	gitC(t, root, "config", "user.email", "test@example.com")
	gitC(t, root, "config", "user.name", "test")
	writeFile(t, filepath.Join(root, "seed.txt"), "seed\n")
	gitC(t, root, "add", "seed.txt")
	gitC(t, root, "commit", "-m", "seed")
	writeFile(t, filepath.Join(root, "README.md"), bootOrphanReadme)

	env := pinnedEnv(t)
	args := []string{"--workflow-dir", root, "--boot"}

	// Present relative worktree -> yes.
	gitC(t, root, "worktree", "add", "--detach", filepath.Join(root, ".worktrees", "present-wt"))
	writeFile(t, filepath.Join(root, "present.md"),
		"---\nstatus: build\nworktree: .worktrees/present-wt\n---\n")

	out, stderr, code := runNative(t, root, env, args...)
	if code != 0 {
		t.Fatalf("native --boot exit=%d stderr=%q", code, stderr)
	}
	if got := orphanDirExists(out); got != "yes" {
		t.Fatalf("present single-root worktree DIR_EXISTS=%q, want yes\n%s", got, out)
	}

	// Non-existent relative worktree -> no.
	writeFile(t, filepath.Join(root, "present.md"),
		"---\nstatus: build\nworktree: .worktrees/missing-wt\n---\n")
	out2, stderr2, code2 := runNative(t, root, env, args...)
	if code2 != 0 {
		t.Fatalf("native --boot (missing) exit=%d stderr=%q", code2, stderr2)
	}
	if got := orphanDirExists(out2); got != "no" {
		t.Fatalf("non-existent single-root worktree DIR_EXISTS=%q, want no\n%s", got, out2)
	}
}
