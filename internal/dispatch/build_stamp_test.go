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
	runGitFatal(t, mainRepo, "init", "-q")
	runGitFatal(t, mainRepo, "config", "user.email", "t@t")
	runGitFatal(t, mainRepo, "config", "user.name", "t")

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
