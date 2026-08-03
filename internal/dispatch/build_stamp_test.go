// ABOUTME: `dispatch build --stamp` (mechanism 3) — stamps+commit+worktree-add,
// ABOUTME: the status-match refusal, idempotent re-run, and the --stamp: stderr contract.
package dispatch

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// stampFixture builds a single-repo split-root workflow: a main repo with a
// README declaring `state: .spacedock-state`, and the state checkout as a
// linked worktree on branch spacedock-state/dev (StateBranch's default for a
// workflow dir named "dev") seeded with one flat entity at status. No origin
// remote — statesync.Publish degrades to local-only, which is enough to prove
// mechanism 3 wires stamp+commit+worktree-add correctly without re-proving
// statesync's own push/rebase/HALT behavior (already covered in
// internal/cli/state_commit_test.go).
func stampFixture(t *testing.T, status_, worktreeStage string, worktree bool) (mainRepo, workflowDir, statePath, entityPath string) {
	t.Helper()
	root := t.TempDir()
	mainRepo = filepath.Join(root, "main")
	if err := os.MkdirAll(mainRepo, 0o755); err != nil {
		t.Fatal(err)
	}
	testgit.InitRepo(t, mainRepo, "-q")

	workflowDir = filepath.Join(mainRepo, "docs", "dev")
	writeFile(t, filepath.Join(workflowDir, "README.md"), readmeWorktree(true))
	runGitFatal(t, mainRepo, "add", "-A")
	runGitFatal(t, mainRepo, "commit", "-q", "-m", "init workflow")

	statePath = filepath.Join(workflowDir, "state-checkout")
	stateBranch := "spacedock-state/dev"
	runGitFatal(t, mainRepo, "worktree", "add", "-b", stateBranch, statePath)

	entityPath = filepath.Join(statePath, "thing.md")
	worktreeField := ""
	if worktree {
		worktreeField = ".worktrees/spacedock-ensign-thing"
	}
	writeFile(t, entityPath, entityFM("Thing", status_, worktreeField))
	runGitFatal(t, statePath, "add", "-A")
	runGitFatal(t, statePath, "commit", "-q", "-m", "seed thing")

	return mainRepo, workflowDir, statePath, entityPath
}

func runGitFatal(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
	}
}

// buildStampArgs returns the flag/file-mode argv for a --stamp build.
func buildStampArgs(workflowDir, entityPath, stage, checklistFile string) []string {
	return []string{"build", "--stamp", "--workflow-dir", workflowDir,
		"--entity-path", entityPath, "--stage", stage, "--checklist-file", checklistFile}
}

func writeChecklist(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "impl.checklist")
	writeFile(t, path, "- a\n- b\n")
	return path
}

// TestStampStagesCommitsAndCreatesWorktree pins mechanism 3's happy path: a
// worktree-declaring stage gets `started` and `worktree=` stamped through the
// native status --set machinery, the state checkout commits the stamp (no
// origin -> local-only, exit 0 continues to envelope assembly), and the CODE
// worktree is created at the main repo root on branch {worker_key}/{slug}.
func TestStampStagesCommitsAndCreatesWorktree(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	mainRepo, workflowDir, statePath, entityPath := stampFixture(t, "implementation", "implementation", false)
	checklistFile := writeChecklist(t)

	native := runNative("", buildStampArgs(workflowDir, entityPath, "implementation", checklistFile)...)
	if native.exit != 0 {
		t.Fatalf("--stamp build exit=%d stdout=%q stderr=%q", native.exit, native.stdout, native.stderr)
	}

	fields := status.ParseFrontmatter(entityPath)
	if fields["started"] == "" {
		t.Errorf("--stamp did not stamp started; frontmatter=%#v", fields)
	}
	wantWorktree := ".worktrees/spacedock-ensign-thing"
	if fields["worktree"] != wantWorktree {
		t.Errorf("--stamp worktree=%q, want %q", fields["worktree"], wantWorktree)
	}

	// The state checkout committed the stamp and is clean (local-only, no
	// origin).
	if porcelain := gitOutput(t, statePath, "status", "--porcelain"); strings.TrimSpace(porcelain) != "" {
		t.Errorf("state checkout dirty after --stamp: %q", porcelain)
	}
	log := gitOutput(t, statePath, "log", "--oneline", "-1")
	if !strings.Contains(log, "dispatch: thing entering implementation") {
		t.Errorf("state checkout commit message = %q, want the dispatch: <slug> entering <stage> shape", log)
	}

	// The CODE worktree was created at the main repo root on the expected branch.
	worktreePath := filepath.Join(mainRepo, wantWorktree)
	if info, err := os.Stat(worktreePath); err != nil || !info.IsDir() {
		t.Fatalf("--stamp did not create worktree at %s: %v", worktreePath, err)
	}
	branch := strings.TrimSpace(gitOutput(t, worktreePath, "rev-parse", "--abbrev-ref", "HEAD"))
	if branch != "spacedock-ensign/thing" {
		t.Errorf("worktree branch = %q, want spacedock-ensign/thing", branch)
	}

	// The build still assembled and emitted a normal spawn envelope.
	if !strings.Contains(native.stdout, `"schema_version"`) {
		t.Errorf("--stamp build emitted no envelope on success:\nstdout=%s\nstderr=%s", native.stdout, native.stderr)
	}
}

// TestStampNonWorktreeStageOnlyStampsStarted pins the non-worktree-stage case:
// only `started` is stamped (no worktree: field, no CODE worktree created).
func TestStampNonWorktreeStageOnlyStampsStarted(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	mainRepo, workflowDir, _, entityPath := stampFixture(t, "backlog", "", false)
	checklistFile := writeChecklist(t)

	native := runNative("", buildStampArgs(workflowDir, entityPath, "backlog", checklistFile)...)
	if native.exit != 0 {
		t.Fatalf("--stamp build exit=%d stdout=%q stderr=%q", native.exit, native.stdout, native.stderr)
	}

	fields := status.ParseFrontmatter(entityPath)
	if fields["started"] == "" {
		t.Errorf("--stamp did not stamp started on a non-worktree stage")
	}
	if fields["worktree"] != "" {
		t.Errorf("--stamp stamped worktree=%q on a non-worktree stage", fields["worktree"])
	}
	if _, err := os.Stat(filepath.Join(mainRepo, ".worktrees")); !os.IsNotExist(err) {
		t.Errorf(".worktrees was created for a non-worktree stage")
	}
}

// TestStampRefusesStatusStageMismatchWithoutMutation pins mechanism 3's status
// guard: --stamp never advances status, so a mismatch refuses outright (no
// envelope, no frontmatter write, no commit) rather than stamping a stale
// dispatch.
func TestStampRefusesStatusStageMismatchWithoutMutation(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	_, workflowDir, statePath, entityPath := stampFixture(t, "backlog", "", false)
	checklistFile := writeChecklist(t)
	before, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatal(err)
	}
	headBefore := strings.TrimSpace(gitOutput(t, statePath, "rev-parse", "HEAD"))

	native := runNative("", buildStampArgs(workflowDir, entityPath, "implementation", checklistFile)...)
	if native.exit == 0 {
		t.Fatalf("--stamp accepted a status/stage mismatch: stdout=%q", native.stdout)
	}
	if !strings.HasPrefix(native.stderr, "dispatch build --stamp:") {
		t.Errorf("stamp failure stderr missing the dispatch build --stamp: prefix: %q", native.stderr)
	}
	if native.stdout != "" {
		t.Errorf("refused --stamp emitted an envelope on stdout: %q", native.stdout)
	}
	after, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Error("refused --stamp mutated the entity frontmatter")
	}
	if headAfter := strings.TrimSpace(gitOutput(t, statePath, "rev-parse", "HEAD")); headAfter != headBefore {
		t.Error("refused --stamp advanced the state checkout HEAD")
	}
}

// TestStampIdempotentReRunSkipsAlreadyStamped pins the re-entrancy case: a
// second --stamp build against an already-stamped entity (started set, worktree
// present) performs no frontmatter mutation and no new state commit, and treats
// the existing worktree as a skip rather than an error (`git worktree add`
// otherwise fatals "already exists").
func TestStampIdempotentReRunSkipsAlreadyStamped(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	_, workflowDir, statePath, entityPath := stampFixture(t, "implementation", "implementation", false)
	checklistFile := writeChecklist(t)

	first := runNative("", buildStampArgs(workflowDir, entityPath, "implementation", checklistFile)...)
	if first.exit != 0 {
		t.Fatalf("first --stamp build exit=%d stderr=%q", first.exit, first.stderr)
	}
	headAfterFirst := strings.TrimSpace(gitOutput(t, statePath, "rev-parse", "HEAD"))
	fieldsAfterFirst := status.ParseFrontmatter(entityPath)

	second := runNative("", buildStampArgs(workflowDir, entityPath, "implementation", checklistFile)...)
	if second.exit != 0 {
		t.Fatalf("second --stamp build exit=%d stdout=%q stderr=%q", second.exit, second.stdout, second.stderr)
	}
	fieldsAfterSecond := status.ParseFrontmatter(entityPath)
	if fieldsAfterSecond["started"] != fieldsAfterFirst["started"] {
		t.Errorf("idempotent re-run changed started: before=%q after=%q", fieldsAfterFirst["started"], fieldsAfterSecond["started"])
	}
	if headAfterSecond := strings.TrimSpace(gitOutput(t, statePath, "rev-parse", "HEAD")); headAfterSecond != headAfterFirst {
		t.Errorf("idempotent re-run created a new state commit: before=%s after=%s", headAfterFirst, headAfterSecond)
	}
	if !strings.Contains(second.stdout, `"schema_version"`) {
		t.Errorf("idempotent re-run did not still assemble an envelope:\nstdout=%s\nstderr=%s", second.stdout, second.stderr)
	}
}

// TestStampAdvanceIncompatible pins the usage-error guard: --stamp and --advance
// both name it a reuse-vs-fresh-dispatch pair that never mixes (a reuse advance
// presupposes an already-stamped live worker, so the post-gate reuse path needs
// no stamps at all).
func TestStampAdvanceIncompatible(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	_, workflowDir, _, entityPath := stampFixture(t, "implementation", "implementation", false)
	checklistFile := writeChecklist(t)

	args := append(buildStampArgs(workflowDir, entityPath, "implementation", checklistFile), "--advance")
	native := runNative("", args...)
	if native.exit != 2 {
		t.Fatalf("--stamp --advance exit=%d, want 2; stdout=%q stderr=%q", native.exit, native.stdout, native.stderr)
	}
	if !strings.Contains(native.stderr, "incompatible") {
		t.Errorf("--stamp --advance stderr missing an incompatibility diagnostic: %q", native.stderr)
	}
}

// TestStampFailureStderrDiscriminatesFromAssemblyFailure pins the stderr
// discrimination contract fo-dispatch-core's block clause relies on: a
// --stamp-phase failure carries the "dispatch build --stamp:" prefix (remedy +
// rerun, never break-glass); an ordinary assembly failure (missing checklist,
// no --stamp involved) does not carry that prefix (break-glass eligible).
func TestStampFailureStderrDiscriminatesFromAssemblyFailure(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	t.Run("stamp failure carries the prefix", func(t *testing.T) {
		_, workflowDir, _, entityPath := stampFixture(t, "backlog", "", false)
		checklistFile := writeChecklist(t)
		native := runNative("", buildStampArgs(workflowDir, entityPath, "implementation", checklistFile)...)
		if native.exit == 0 {
			t.Fatal("expected a stamp status-mismatch failure")
		}
		if !strings.HasPrefix(native.stderr, "dispatch build --stamp:") {
			t.Errorf("stamp failure stderr = %q, want the dispatch build --stamp: prefix", native.stderr)
		}
	})

	t.Run("assembly failure without --stamp carries no prefix", func(t *testing.T) {
		root := t.TempDir()
		workflowDir := filepath.Join(root, "wf")
		writeFile(t, filepath.Join(workflowDir, "README.md"), readmeWorktree(false))
		entityPath := filepath.Join(workflowDir, "thing.md")
		writeFile(t, entityPath, entityFM("Thing", "backlog", ""))
		gitInit(t, root)

		emptyChecklist := filepath.Join(t.TempDir(), "empty.checklist")
		writeFile(t, emptyChecklist, "")
		native := runNative("", "build", "--workflow-dir", workflowDir,
			"--entity-path", entityPath, "--stage", "backlog", "--checklist-file", emptyChecklist)
		if native.exit == 0 {
			t.Fatal("expected an empty-checklist assembly failure")
		}
		if strings.HasPrefix(native.stderr, "dispatch build --stamp:") {
			t.Errorf("assembly failure stderr wrongly carries the --stamp: prefix: %q", native.stderr)
		}
	})
}

// gitOutput runs a git command in dir and returns combined output, failing the
// test on a non-zero exit.
func gitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
	}
	return string(out)
}

// stampFixtureWithOrigin mirrors stampFixture but gives the state checkout a
// real bare origin (pushed once), so a pre-push hook can force a genuine
// publish failure — needed to exercise the retry-after-sync-failure path,
// since statesync.Publish degrades to local-only (no push attempted at all)
// when there is no origin remote.
func stampFixtureWithOrigin(t *testing.T, status_ string, worktree bool) (mainRepo, workflowDir, statePath, entityPath, bareDir string) {
	t.Helper()
	bareDir = filepath.Join(t.TempDir(), "origin.git")
	runGitFatal(t, t.TempDir(), "init", "-q", "--bare", bareDir)

	mainRepo = filepath.Join(t.TempDir(), "host")
	runGitFatal(t, t.TempDir(), "clone", "-q", bareDir, mainRepo)
	runGitFatal(t, mainRepo, "config", "user.email", "t@t")
	runGitFatal(t, mainRepo, "config", "user.name", "t")

	workflowDir = filepath.Join(mainRepo, "docs", "dev")
	writeFile(t, filepath.Join(workflowDir, "README.md"), readmeWorktree(true))
	runGitFatal(t, mainRepo, "add", "-A")
	runGitFatal(t, mainRepo, "commit", "-q", "-m", "init workflow")
	runGitFatal(t, mainRepo, "push", "-q", "origin", "HEAD")

	statePath = filepath.Join(workflowDir, "state-checkout")
	stateBranch := "spacedock-state/dev"
	runGitFatal(t, mainRepo, "worktree", "add", "-b", stateBranch, statePath)

	entityPath = filepath.Join(statePath, "thing.md")
	worktreeField := ""
	if worktree {
		worktreeField = ".worktrees/spacedock-ensign-thing"
	}
	writeFile(t, entityPath, entityFM("Thing", status_, worktreeField))
	runGitFatal(t, statePath, "add", "-A")
	runGitFatal(t, statePath, "commit", "-q", "-m", "seed thing")
	runGitFatal(t, statePath, "push", "-q", "origin", stateBranch)

	return mainRepo, workflowDir, statePath, entityPath, bareDir
}

// blockPush installs a pre-push hook that always fails in checkout, forcing a
// publish attempt to fail without touching origin. Returns a restore func.
func blockPush(t *testing.T, checkout string) (restore func()) {
	t.Helper()
	hooks := t.TempDir()
	if err := os.WriteFile(filepath.Join(hooks, "pre-push"), []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	runGitFatal(t, checkout, "config", "core.hooksPath", hooks)
	return func() { runGitFatal(t, checkout, "config", "--unset", "core.hooksPath") }
}

// TestStampRetriesSyncOnRetryEvenWhenAlreadyStamped pins finding 1: a --stamp
// retry after an earlier sync=failed must not skip synchronization just
// because started/worktree are already stamped from the failed attempt — it
// must still attempt to publish (and succeed, once the block is lifted)
// before creating the worktree, rather than silently proceeding against
// unresolved divergent state.
func TestStampRetriesSyncOnRetryEvenWhenAlreadyStamped(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	mainRepo, workflowDir, statePath, entityPath, bareDir := stampFixtureWithOrigin(t, "implementation", false)
	checklistFile := writeChecklist(t)

	restore := blockPush(t, statePath)
	first := runNative("", buildStampArgs(workflowDir, entityPath, "implementation", checklistFile)...)
	if first.exit != 1 {
		t.Fatalf("first --stamp (blocked push) exit=%d stdout=%q stderr=%q, want 1", first.exit, first.stdout, first.stderr)
	}
	fieldsAfterFirst := status.ParseFrontmatter(entityPath)
	if fieldsAfterFirst["started"] == "" {
		t.Fatal("first --stamp did not commit the stamp locally before the sync failure")
	}
	worktreePath := filepath.Join(mainRepo, ".worktrees", "spacedock-ensign-thing")
	if _, err := os.Stat(worktreePath); err == nil {
		t.Fatal("first --stamp created the worktree despite the sync failure")
	}

	restore()
	second := runNative("", buildStampArgs(workflowDir, entityPath, "implementation", checklistFile)...)
	if second.exit != 0 {
		t.Fatalf("retried --stamp exit=%d stdout=%q stderr=%q, want 0 once the block is lifted", second.exit, second.stdout, second.stderr)
	}
	if porcelain := gitOutput(t, statePath, "status", "--porcelain"); strings.TrimSpace(porcelain) != "" {
		t.Errorf("state checkout still dirty after the retried --stamp: %q", porcelain)
	}
	// The retry's own commit-or-noop step no-ops (nothing new to stage), but
	// Publish must still have run and pushed the FIRST attempt's local commit —
	// the bug this pins is exactly that a retry skipped this step entirely.
	originHead := strings.TrimSpace(gitOutput(t, bareDir, "rev-parse", "spacedock-state/dev"))
	localHead := strings.TrimSpace(gitOutput(t, statePath, "rev-parse", "HEAD"))
	if originHead != localHead {
		t.Errorf("origin HEAD %s != local HEAD %s after the retried --stamp; the first attempt's commit was never published", originHead, localHead)
	}
	if _, err := os.Stat(worktreePath); err != nil {
		t.Errorf("retried --stamp did not create the worktree: %v", err)
	}
}

// TestStampRefusesMismatchedEntityPath pins finding 2: a caller-supplied
// entity_path that shares a slug with the real canonical entity but is not
// itself that file (a stray duplicate elsewhere) must be refused — not
// validated against the wrong file while status --set mutates the real one.
func TestStampRefusesMismatchedEntityPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	_, workflowDir, statePath, entityPath := stampFixture(t, "implementation", "implementation", false)
	checklistFile := writeChecklist(t)

	// A stray duplicate named identically (same slug "thing") but living
	// outside the discoverable state checkout.
	strayPath := filepath.Join(t.TempDir(), "thing.md")
	body, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, strayPath, string(body))
	before := string(body)

	native := runNative("", buildStampArgs(workflowDir, strayPath, "implementation", checklistFile)...)
	if native.exit == 0 {
		t.Fatalf("--stamp accepted a mismatched entity_path sharing a slug with the canonical entity: stdout=%q", native.stdout)
	}
	if !strings.HasPrefix(native.stderr, "dispatch build --stamp:") {
		t.Errorf("mismatched-path failure stderr missing the --stamp: prefix: %q", native.stderr)
	}
	after, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != before {
		t.Error("refused --stamp still mutated the real canonical entity")
	}
	if porcelain := gitOutput(t, statePath, "status", "--porcelain"); strings.TrimSpace(porcelain) != "" {
		t.Errorf("refused --stamp left the state checkout dirty: %q", porcelain)
	}
}

// TestStampCommitsInlineBeforeWorktreeCreation pins finding 3: for an inline
// (single-root) workflow, --stamp must commit the entity's dirty state in the
// main repo BEFORE creating the worktree — a worktree checks out committed
// HEAD only, so an uncommitted stamp (or an uncommitted earlier gate
// record/consume write) would otherwise hand a freshly-dispatched worker a
// stale, pre-decision copy of the entity.
func TestStampCommitsInlineBeforeWorktreeCreation(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	workflowDir := filepath.Join(root, "wf")
	writeFile(t, filepath.Join(workflowDir, "README.md"), readmeWorktree(false))
	writeFile(t, filepath.Join(root, ".gitignore"), ".worktrees/\n")
	entityPath := filepath.Join(workflowDir, "thing.md")
	// Simulate a just-recorded gate decision: status already at the target
	// stage, dirty/uncommitted, exactly as gate consume would leave it for an
	// inline workflow (mechanism 1 never syncs inline).
	writeFile(t, entityPath, entityFM("Thing", "implementation", ""))
	gitInit(t, root)
	// gitInit's own commit passes -c user.name/user.email inline, but the
	// inline --stamp commit under test (stampCommitInline) runs a plain `git
	// commit` that relies on the repo's configured identity — CI runners carry
	// no global git identity, so that commit fails there without this.
	runGitFatal(t, root, "config", "user.name", "Spacedock Test")
	runGitFatal(t, root, "config", "user.email", "spacedock@example.invalid")
	writeFile(t, entityPath, strings.Replace(entityFM("Thing", "implementation", ""), "Body.", "Body (post-decision, uncommitted).", 1))
	if headBefore, dirty := gitOutput(t, root, "rev-parse", "HEAD"), gitOutput(t, root, "status", "--porcelain"); strings.TrimSpace(dirty) == "" {
		t.Fatalf("fixture setup did not leave the entity dirty relative to HEAD=%s", strings.TrimSpace(headBefore))
	}
	checklistFile := writeChecklist(t)

	native := runNative("", buildStampArgs(workflowDir, entityPath, "implementation", checklistFile)...)
	if native.exit != 0 {
		t.Fatalf("--stamp build exit=%d stdout=%q stderr=%q", native.exit, native.stdout, native.stderr)
	}
	if porcelain := gitOutput(t, root, "status", "--porcelain"); strings.TrimSpace(porcelain) != "" {
		t.Errorf("main repo still dirty after inline --stamp: %q", porcelain)
	}

	worktreePath := filepath.Join(root, ".worktrees", "spacedock-ensign-thing")
	// The worktree checks out the whole git root's tree, so the entity lives
	// under the same workflowDir-relative path as in the main checkout (root/wf/
	// thing.md), not directly at the worktree's own root.
	worktreeEntity := filepath.Join(worktreePath, "wf", "thing.md")
	body, err := os.ReadFile(worktreeEntity)
	if err != nil {
		t.Fatalf("worktree missing its own entity copy: %v", err)
	}
	if !strings.Contains(string(body), "post-decision, uncommitted") {
		t.Errorf("worktree's entity copy is stale (built before the inline commit):\n%s", body)
	}
	fields := status.ParseFrontmatter(worktreeEntity)
	if fields["started"] == "" {
		t.Errorf("worktree's entity copy is missing the started stamp: %#v", fields)
	}
}

// TestStampRefusesWorktreePathOnWrongBranch pins finding 4: an existing
// filesystem object at the expected worktree path must not be treated as an
// idempotent skip unless it is a REGISTERED worktree on exactly
// {worker_key}/{slug} — a stray directory, or a worktree left on the wrong
// branch, is an error, not a silent skip.
func TestStampRefusesWorktreePathOnWrongBranch(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	mainRepo, workflowDir, _, entityPath := stampFixture(t, "implementation", "implementation", false)
	checklistFile := writeChecklist(t)

	worktreePath := filepath.Join(mainRepo, ".worktrees", "spacedock-ensign-thing")
	runGitFatal(t, mainRepo, "worktree", "add", "-b", "some-other-branch", worktreePath)

	native := runNative("", buildStampArgs(workflowDir, entityPath, "implementation", checklistFile)...)
	if native.exit == 0 {
		t.Fatalf("--stamp silently accepted a worktree path already occupied by the wrong branch: stdout=%q", native.stdout)
	}
	if !strings.HasPrefix(native.stderr, "dispatch build --stamp:") {
		t.Errorf("wrong-branch worktree failure stderr missing the --stamp: prefix: %q", native.stderr)
	}
	branch := strings.TrimSpace(gitOutput(t, worktreePath, "rev-parse", "--abbrev-ref", "HEAD"))
	if branch != "some-other-branch" {
		t.Errorf("refused --stamp mutated the existing worktree's branch: now %q", branch)
	}
}
