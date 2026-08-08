// ABOUTME: mergeRootsGuard coverage (issue #485) — a --workflow-dir that cannot
// ABOUTME: possibly hold entities refuses with a named, actionable diagnostic.
package status

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// mergeGuardWorktreeReadme declares a split-root workflow whose terminal
// transition finalizes directly (no merge hook, no worktree-bearing stage) so
// the AC-1 real-worktree repro's corrected re-run finalizes in one MergeGuard
// call, matching the no-hook Phase C default path other merge_guard_test.go
// fixtures already exercise.
const mergeGuardWorktreeReadme = `---
commissioned-by: spacedock@1
id-style: sequential
state: .spacedock-state
stages:
  states:
    - name: implementation
      initial: true
    - name: done
      terminal: true
---

# Merge Guard Foreign-Cwd Fixture Workflow
`

// buildMergeGuardForeignCwdFixture materializes the issue-#485 repro topology: a
// code repo whose docs/dev declares a split-root state checkout, with
// docs/dev/README.md tracked in the code repo (so it appears in a linked
// worktree) and docs/dev/.spacedock-state its OWN separate git checkout (so it
// does NOT appear in a linked worktree — matching the live topology where the
// state checkout is a separate checkout that exists only under the main working
// copy), plus a real linked agent worktree of the code repo. Returns the
// main-checkout root, the definition dir, and the worktree dir.
func buildMergeGuardForeignCwdFixture(t *testing.T) (coderoot, defDir, wtDir string) {
	t.Helper()
	coderoot = t.TempDir()
	testgit.InitRepo(t, coderoot)

	defDir = filepath.Join(coderoot, "docs", "dev")
	writeFile(t, filepath.Join(defDir, "README.md"), mergeGuardWorktreeReadme)
	gitC(t, coderoot, "add", "docs/dev/README.md")
	gitC(t, coderoot, "commit", "-q", "-m", "seed")

	state := filepath.Join(defDir, ".spacedock-state")
	writeFile(t, filepath.Join(state, "010-repro.md"),
		"---\nid: \"010\"\nstatus: implementation\n---\n# Repro entity\n")
	testgit.InitRepo(t, state)
	gitC(t, state, "add", "010-repro.md")
	gitC(t, state, "commit", "-q", "-m", "seed")
	gitC(t, state, "branch", "-M", "spacedock-state/dev")

	wtDir = filepath.Join(coderoot, ".worktrees", "agent-x")
	gitC(t, coderoot, "worktree", "add", "--detach", wtDir)
	return coderoot, defDir, wtDir
}

// TestMergeGuardForeignCwdRefusalNamesWorkingFix is AC-1 (measures the
// end-value): from a linked-worktree cwd, the issue-#485 repro invocation
// (`merge guard <slug> --verdict passed --workflow-dir docs/dev`, with the
// split-root state checkout absent in the worktree) produces a refusal whose
// suggested corrected invocation actually works.
func TestMergeGuardForeignCwdRefusalNamesWorkingFix(t *testing.T) {
	coderoot, defDir, wtDir := buildMergeGuardForeignCwdFixture(t)

	var out, errBuf bytes.Buffer
	code := MergeGuard([]string{"010-repro", "--verdict", "passed", "--workflow-dir", "docs/dev"}, wtDir, &out, &errBuf)
	if code != 1 {
		t.Fatalf("foreign-cwd repro: exit = %d, want 1 (out=%q stderr=%q)", code, out.String(), errBuf.String())
	}
	stderr := errBuf.String()
	resolvedDefDir := filepath.Join(wtDir, "docs", "dev")
	missingState := filepath.Join(resolvedDefDir, ".spacedock-state")
	if !strings.Contains(stderr, resolvedDefDir) {
		t.Fatalf("refusal must name the resolved absolute path %q, got %q", resolvedDefDir, stderr)
	}
	if !strings.Contains(stderr, missingState) {
		t.Fatalf("refusal must name the missing state-checkout path %q, got %q", missingState, stderr)
	}
	// git rev-parse --git-common-dir resolves symlinks (e.g. macOS's /var ->
	// /private/var), so compare via realpathOf rather than raw string equality —
	// the same normalization TestDiscoverWorkflowDirWalksUp uses.
	hintFlag := "--workflow-dir " + realpathOf(defDir)
	if !strings.Contains(stderr, hintFlag) {
		t.Fatalf("refusal must name the corrected invocation %q, got %q", hintFlag, stderr)
	}
	if strings.Contains(stderr, "entity not found") {
		t.Fatalf("refusal must not fall through to the misleading entity-not-found message, got %q", stderr)
	}

	// Re-run with the suggested absolute --workflow-dir: the corrected invocation
	// must actually succeed and land the mutation in the main checkout's state dir.
	var out2, errBuf2 bytes.Buffer
	code2 := MergeGuard([]string{"010-repro", "--verdict", "passed", "--workflow-dir", defDir}, wtDir, &out2, &errBuf2)
	if code2 != 0 {
		t.Fatalf("corrected invocation: exit = %d, want 0 (stderr=%q)", code2, errBuf2.String())
	}
	archived := filepath.Join(coderoot, "docs", "dev", ".spacedock-state", "_archive", "010-repro.md")
	if !fileExists(archived) {
		t.Fatalf("corrected invocation must finalize+archive the entity at %q", archived)
	}
	if got := frontmatterField(t, archived, "status"); got != "done" {
		t.Fatalf("archived entity status = %q, want done", got)
	}
	if got := frontmatterField(t, archived, "verdict"); got != "PASSED" {
		t.Fatalf("archived entity verdict = %q, want PASSED (the schema-cased stored value)", got)
	}
}

// TestMergeGuardForeignCwdSuppressesHintWhenCandidateInvalid closes the
// negative-hint-suppression gap found by a detached adversarial audit (cycle 1):
// removing didYouMeanHint's validateRootsOrExit check made the hint fire for a
// syntactically-joined main-root candidate that does not actually resolve to a
// workflow, suggesting a corrected invocation that itself would not work. From a
// linked-worktree cwd, a relative --workflow-dir whose main-root-joined
// candidate does not exist in either checkout must refuse WITHOUT a "did you
// mean" line — the refusal stands on its own per the Proposed approach.
func TestMergeGuardForeignCwdSuppressesHintWhenCandidateInvalid(t *testing.T) {
	_, _, wtDir := buildMergeGuardForeignCwdFixture(t)

	var out, errBuf bytes.Buffer
	code := MergeGuard([]string{"010-repro", "--verdict", "passed", "--workflow-dir", "docs/bogus"}, wtDir, &out, &errBuf)
	if code != 1 {
		t.Fatalf("exit = %d, want 1 (out=%q stderr=%q)", code, out.String(), errBuf.String())
	}
	stderr := errBuf.String()
	resolved := filepath.Join(wtDir, "docs", "bogus")
	if !strings.Contains(stderr, resolved) {
		t.Fatalf("refusal must still name the nonexistent resolved dir %q, got %q", resolved, stderr)
	}
	if strings.Contains(stderr, "did you mean") {
		t.Fatalf("refusal must NOT suggest a did-you-mean hint when the main-root-joined candidate does not validate, got %q", stderr)
	}
}

// TestMergeGuardRefusesNonexistentWorkflowDir is AC-2: a --workflow-dir
// resolving to a nonexistent directory is refused with a diagnostic naming the
// as-passed spelling and the resolved absolute path — never falling through to
// "entity not found".
func TestMergeGuardRefusesNonexistentWorkflowDir(t *testing.T) {
	dir := t.TempDir()
	var out, errBuf bytes.Buffer
	code := MergeGuard([]string{"my-task", "--verdict", "passed", "--workflow-dir", "does/not/exist"}, dir, &out, &errBuf)
	if code != 1 {
		t.Fatalf("exit = %d, want 1 (stderr=%q)", code, errBuf.String())
	}
	stderr := errBuf.String()
	if !strings.Contains(stderr, "does/not/exist") {
		t.Fatalf("refusal must name the as-passed spelling, got %q", stderr)
	}
	resolved := filepath.Join(dir, "does/not/exist")
	if !strings.Contains(stderr, resolved) {
		t.Fatalf("refusal must name the resolved absolute path %q, got %q", resolved, stderr)
	}
	if strings.Contains(stderr, "entity not found") {
		t.Fatalf("must not fall through to entity not found, got %q", stderr)
	}
}

// TestMergeGuardRefusesMissingStateCheckout is AC-3: a resolved split-root
// workflow whose declared state checkout is missing is refused before entity
// resolution, naming the missing state-checkout path.
func TestMergeGuardRefusesMissingStateCheckout(t *testing.T) {
	def, state := buildSplitRoot(t, splitRootReadme, map[string]string{
		"add-login.md": "---\nstatus: ideation\n---\n",
	})
	if err := os.RemoveAll(state); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	code := MergeGuard([]string{"add-login", "--verdict", "passed", "--workflow-dir", def}, def, &out, &errBuf)
	if code != 1 {
		t.Fatalf("exit = %d, want 1 (stderr=%q)", code, errBuf.String())
	}
	stderr := errBuf.String()
	if !strings.Contains(stderr, state) {
		t.Fatalf("refusal must name the missing state-checkout path %q, got %q", state, stderr)
	}
	if strings.Contains(stderr, "entity not found") {
		t.Fatalf("must not fall through to entity not found, got %q", stderr)
	}
}

// TestMergeGuardWrongSlugKeepsEntityNotFound is AC-4: a wrong slug against a
// valid, resolvable workflow still reports the plain entity-not-found error —
// the roots guard does not shadow it.
func TestMergeGuardWrongSlugKeepsEntityNotFound(t *testing.T) {
	root, _, errOut, code := driveMergeGuard(t, "merge-local-workflow", "does-not-exist", "--verdict", "passed")
	if code != 1 {
		t.Fatalf("exit = %d, want 1 (root=%q stderr=%q)", code, root, errOut)
	}
	if !strings.Contains(errOut, "Error: entity not found: does-not-exist") {
		t.Fatalf("stderr = %q, want entity not found", errOut)
	}
}
