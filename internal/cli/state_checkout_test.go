// ABOUTME: Real-git e2e for main-worktree-anchored split-root checkout resolution
// ABOUTME: (issue #484) — worktree-cwd anchoring, no-origin/unreachable-origin
// ABOUTME: resume fallback, and targeted stale-worktree-registration repair.
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/dispatch"
	"github.com/spacedock-dev/spacedock/internal/status"
)

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

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

func TestStateReadyJSONAbsentCheckoutIsAtomic(t *testing.T) {
	cases := []struct {
		name  string
		setup func(t *testing.T) (root, workflowDir, statePath string)
	}{
		{
			name: "local-only with stale registration",
			setup: func(t *testing.T) (string, string, string) {
				root, workflowDir, statePath := noOriginSplitWorkflow(t)
				if err := os.RemoveAll(statePath); err != nil {
					t.Fatal(err)
				}
				return root, workflowDir, statePath
			},
		},
		{
			name: "origin-backed fresh clone",
			setup: func(t *testing.T) (string, string, string) {
				bare := filepath.Join(t.TempDir(), "origin.git")
				git(t, t.TempDir(), "init", "-q", "--bare", bare)
				hostA := filepath.Join(t.TempDir(), "host-a")
				git(t, t.TempDir(), "clone", "-q", bare, hostA)
				git(t, hostA, "config", "user.email", "a@t")
				git(t, hostA, "config", "user.name", "a")
				commissionSplitWorkflow(t, hostA)
				git(t, hostA, "push", "-q", "origin", "HEAD")
				fresh := filepath.Join(t.TempDir(), "fresh")
				git(t, t.TempDir(), "clone", "-q", bare, fresh)
				return fresh, filepath.Join(fresh, "docs", "dev"), filepath.Join(fresh, "docs", "dev", ".spacedock-state")
			},
		},
		{
			name: "unreachable origin local fallback warning suppressed",
			setup: func(t *testing.T) (string, string, string) {
				bare := filepath.Join(t.TempDir(), "origin.git")
				git(t, t.TempDir(), "init", "-q", "--bare", bare)
				host := filepath.Join(t.TempDir(), "host")
				git(t, t.TempDir(), "clone", "-q", bare, host)
				git(t, host, "config", "user.email", "t@t")
				git(t, host, "config", "user.name", "t")
				workflowDir, _, _ := commissionSplitWorkflow(t, host)
				statePath := filepath.Join(workflowDir, ".spacedock-state")
				if err := os.RemoveAll(statePath); err != nil {
					t.Fatal(err)
				}
				git(t, host, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "missing.git"))
				return host, workflowDir, statePath
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("HOME", t.TempDir())
			root, workflowDir, statePath := tc.setup(t)
			var out, errBuf strings.Builder
			code := run(context.Background(), []string{"state", "ready", "--workflow-dir", workflowDir, "--json"},
				os.Environ(), root, nil, &out, &errBuf, &status.NativeRunner{}, nil)
			if code != 0 {
				t.Fatalf("state ready --json exit=%d stderr=%q stdout=%q", code, errBuf.String(), out.String())
			}
			if errBuf.String() != "" {
				t.Fatalf("successful state ready --json wrote stderr=%q", errBuf.String())
			}
			if !json.Valid([]byte(out.String())) {
				t.Fatalf("stdout is not one atomic JSON document: %q", out.String())
			}
			var envelope map[string]any
			if err := json.Unmarshal([]byte(out.String()), &envelope); err != nil {
				t.Fatal(err)
			}
			if envelope["command"] != "state ready" || envelope["result"] != "ready" {
				t.Fatalf("unexpected envelope: %#v", envelope)
			}
			if _, err := os.Stat(statePath); err != nil {
				t.Fatalf("state checkout not restored: %v", err)
			}
		})
	}
}

func TestConcurrentStateReadySerializesRepairAndCreation(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root, workflowDir, statePath := noOriginSplitWorkflow(t)
	entity := filepath.Join(statePath, "race-task.md")
	if err := os.WriteFile(entity, []byte("---\nstatus: implementation\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, statePath, "add", "race-task.md")
	git(t, statePath, "commit", "-q", "-m", "seed race task")
	wantHead := strings.TrimSpace(git(t, statePath, "rev-parse", "HEAD"))
	if err := os.RemoveAll(statePath); err != nil {
		t.Fatal(err)
	}

	const callers = 12
	start := make(chan struct{})
	results := make(chan string, callers)
	var wg sync.WaitGroup
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			var out, errBuf strings.Builder
			code := run(context.Background(), []string{"state", "ready", "--workflow-dir", workflowDir},
				os.Environ(), root, nil, &out, &errBuf, &status.NativeRunner{}, nil)
			results <- strings.Join([]string{fmt.Sprint(code), out.String(), errBuf.String()}, "|")
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	for result := range results {
		if !strings.HasPrefix(result, "0|") {
			t.Fatalf("concurrent state ready failed: %q", result)
		}
	}
	if got := strings.TrimSpace(git(t, statePath, "rev-parse", "HEAD")); got != wantHead {
		t.Fatalf("state HEAD changed: got %s want %s", got, wantHead)
	}
	if _, err := os.Stat(entity); err != nil {
		t.Fatalf("entity missing after concurrent resume: %v", err)
	}
	registrations := 0
	for _, line := range strings.Split(git(t, root, "worktree", "list", "--porcelain"), "\n") {
		if strings.HasPrefix(line, "worktree ") && status.RealpathOf(strings.TrimSpace(strings.TrimPrefix(line, "worktree "))) == status.RealpathOf(statePath) {
			registrations++
		}
	}
	if registrations != 1 {
		t.Fatalf("state checkout registrations=%d, want exactly 1", registrations)
	}
}

func runConcurrentStateReadyAtAbsentBoundary(t *testing.T, root, workflowDir string, callers int) []string {
	t.Helper()
	arrived := make(chan struct{}, callers)
	release := make(chan struct{})
	stateReadyObservationHook = func() {
		arrived <- struct{}{}
		<-release
	}
	t.Cleanup(func() { stateReadyObservationHook = nil })

	results := make(chan string, callers)
	var wg sync.WaitGroup
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var out, errBuf strings.Builder
			code := run(context.Background(), []string{"state", "ready", "--workflow-dir", workflowDir},
				os.Environ(), root, nil, &out, &errBuf, &status.NativeRunner{}, nil)
			results <- strings.Join([]string{fmt.Sprint(code), out.String(), errBuf.String()}, "|")
		}()
	}
	for i := 0; i < callers; i++ {
		<-arrived
	}
	close(release)
	wg.Wait()
	close(results)
	var collected []string
	for result := range results {
		collected = append(collected, result)
	}
	return collected
}

func TestConcurrentTwoWorkflowResumeOutcomesArePathIsolated(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := t.TempDir()
	git(t, root, "init", "-q")
	git(t, root, "config", "user.email", "t@t")
	git(t, root, "config", "user.name", "t")

	workflowA := filepath.Join(root, "docs", "alpha")
	workflowB := filepath.Join(root, "docs", "beta")
	for _, workflowDir := range []string{workflowA, workflowB} {
		if err := os.MkdirAll(workflowDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(workflowDir, "README.md"), []byte(splitWorkflowReadme), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git(t, root, "add", "docs")
	git(t, root, "commit", "-q", "-m", "two split workflows")

	stateA := filepath.Join(workflowA, ".spacedock-state")
	stateB := filepath.Join(workflowB, ".spacedock-state")
	branchA := "spacedock-state/alpha"
	branchB := "spacedock-state/beta"
	git(t, root, "branch", branchA)
	git(t, root, "branch", branchB)
	git(t, root, "worktree", "add", "-q", stateA, branchA)
	git(t, root, "worktree", "add", "-q", stateB, branchB)
	if err := os.RemoveAll(stateA); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(stateB); err != nil {
		t.Fatal(err)
	}

	type call struct{ workflow string }
	calls := []call{{workflowA}, {workflowA}, {workflowB}, {workflowB}}
	arrived := make(chan struct{}, len(calls))
	release := make(chan struct{})
	stateReadyObservationHook = func() {
		arrived <- struct{}{}
		<-release
	}
	t.Cleanup(func() { stateReadyObservationHook = nil })
	results := make(chan string, len(calls))
	var wg sync.WaitGroup
	for _, c := range calls {
		c := c
		wg.Add(1)
		go func() {
			defer wg.Done()
			var out, errBuf strings.Builder
			code := run(context.Background(), []string{"state", "ready", "--workflow-dir", c.workflow},
				os.Environ(), root, nil, &out, &errBuf, &status.NativeRunner{}, nil)
			results <- fmt.Sprintf("%d|%s|%s", code, out.String(), errBuf.String())
		}()
	}
	for range calls {
		<-arrived
	}
	close(release)
	wg.Wait()
	close(results)
	stateReadyObservationHook = nil
	for result := range results {
		if !strings.HasPrefix(result, "0|") {
			t.Fatalf("two-workflow concurrent resume failed: %q", result)
		}
	}

	outcomeA, err := stateResumeOutcomePath(workflowA, stateA)
	if err != nil {
		t.Fatal(err)
	}
	outcomeB, err := stateResumeOutcomePath(workflowB, stateB)
	if err != nil {
		t.Fatal(err)
	}
	if outcomeA == outcomeB {
		t.Fatalf("distinct canonical checkouts share outcome path %s", outcomeA)
	}
	if result, err := readStateResumeOutcome(workflowA, stateA); err != nil || result != "ready" {
		t.Fatalf("workflow A outcome=%q err=%v", result, err)
	}
	if result, err := readStateResumeOutcome(workflowB, stateB); err != nil || result != "ready" {
		t.Fatalf("workflow B outcome=%q err=%v", result, err)
	}
	if err := writeStateResumeOutcome(workflowA, stateA, "failed"); err != nil {
		t.Fatal(err)
	}
	if result, err := readStateResumeOutcome(workflowB, stateB); err != nil || result != "ready" {
		t.Fatalf("workflow A overwrite leaked into B: outcome=%q err=%v", result, err)
	}
}

func originBackedCheckoutBehindRemote(t *testing.T) (root, workflowDir, statePath, wantHead string) {
	t.Helper()
	bare, workflowA, workflowB, stateBranch := twoHostStateWorkflow(t)
	hostB := filepath.Dir(filepath.Dir(workflowB))
	writeEntity(t, workflowB, "first-task", "---\nstatus: ideation\n---\n# Peer state\n")
	if code, _, errOut := runStateCommitCmd(t, hostB, workflowB, "first-task", "-m", "peer state"); code != 0 {
		t.Fatalf("peer state commit exit=%d stderr=%q", code, errOut)
	}
	wantHead = strings.TrimSpace(git(t, bare, "rev-parse", stateBranch))
	statePath = filepath.Join(workflowA, ".spacedock-state")
	if err := os.RemoveAll(statePath); err != nil {
		t.Fatal(err)
	}
	return filepath.Dir(filepath.Dir(workflowA)), workflowA, statePath, wantHead
}

func TestStateInitRestoredLocalBranchIntegratesFetchedRemoteState(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root, workflowDir, statePath, wantHead := originBackedCheckoutBehindRemote(t)
	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "init", "--workflow-dir", workflowDir},
		os.Environ(), root, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("state init exit=%d stdout=%q stderr=%q", code, out.String(), errBuf.String())
	}
	if got := strings.TrimSpace(git(t, statePath, "rev-parse", "HEAD")); got != wantHead {
		t.Fatalf("restored state HEAD=%s, want fetched peer HEAD=%s", got, wantHead)
	}
	body, err := os.ReadFile(filepath.Join(statePath, "first-task.md"))
	if err != nil || !strings.Contains(string(body), "# Peer state") {
		t.Fatalf("restored state omitted peer bytes: err=%v body=%q", err, body)
	}
}

func TestStateReadyOriginConvergesBeforeFinalPathPublication(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root, workflowDir, statePath, wantHead := originBackedCheckoutBehindRemote(t)
	branch, err := status.StateBranch(workflowDir)
	if err != nil {
		t.Fatal(err)
	}
	observed := make(chan string, 1)
	release := make(chan struct{})
	stateResumeBeforePublishHook = func(path string) {
		if status.RealpathOf(path) != status.RealpathOf(statePath) {
			return
		}
		if _, err := os.Lstat(statePath); !os.IsNotExist(err) {
			observed <- fmt.Sprintf("final path visible before publication: %v", err)
		} else if ok, head := runGit(root, "rev-parse", branch); !ok || strings.TrimSpace(head) != wantHead {
			observed <- fmt.Sprintf("branch not converged privately: ok=%v head=%q want=%s", ok, head, wantHead)
		} else {
			observed <- ""
		}
		<-release
	}
	t.Cleanup(func() { stateResumeBeforePublishHook = nil })

	type resumeResult struct {
		code int
		err  string
	}
	result := make(chan resumeResult, 1)
	go func() {
		var out, errBuf strings.Builder
		code := run(context.Background(), []string{"state", "ready", "--workflow-dir", workflowDir},
			os.Environ(), root, nil, &out, &errBuf, &status.NativeRunner{}, nil)
		result <- resumeResult{code: code, err: errBuf.String()}
	}()
	if problem := <-observed; problem != "" {
		close(release)
		t.Fatal(problem)
	}
	close(release)
	got := <-result
	if got.code != 0 {
		t.Fatalf("private origin convergence failed: %+v", got)
	}
	if head := strings.TrimSpace(git(t, statePath, "rev-parse", "HEAD")); head != wantHead {
		t.Fatalf("published HEAD=%s want privately converged %s", head, wantHead)
	}
}

func TestConcurrentStateReadyReachableOriginConvergesAtomically(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root, workflowDir, statePath, wantHead := originBackedCheckoutBehindRemote(t)
	for _, result := range runConcurrentStateReadyAtAbsentBoundary(t, root, workflowDir, 8) {
		if !strings.HasPrefix(result, "0|") {
			t.Fatalf("reachable-origin concurrent ready failed: %q", result)
		}
	}
	if got := strings.TrimSpace(git(t, statePath, "rev-parse", "HEAD")); got != wantHead {
		t.Fatalf("concurrent ready HEAD=%s, want integrated peer HEAD=%s", got, wantHead)
	}
}

func TestConcurrentStateReadyUnreachableOriginWaitersDoNotRepull(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	bare := filepath.Join(t.TempDir(), "origin.git")
	git(t, t.TempDir(), "init", "-q", "--bare", bare)
	root := filepath.Join(t.TempDir(), "host")
	git(t, t.TempDir(), "clone", "-q", bare, root)
	git(t, root, "config", "user.email", "t@t")
	git(t, root, "config", "user.name", "t")
	workflowDir, _, _ := commissionSplitWorkflow(t, root)
	statePath := filepath.Join(workflowDir, ".spacedock-state")
	wantHead := strings.TrimSpace(git(t, statePath, "rev-parse", "HEAD"))
	if err := os.RemoveAll(statePath); err != nil {
		t.Fatal(err)
	}
	git(t, root, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "missing.git"))

	for _, result := range runConcurrentStateReadyAtAbsentBoundary(t, root, workflowDir, 8) {
		if !strings.HasPrefix(result, "0|") {
			t.Fatalf("unreachable-origin concurrent ready failed: %q", result)
		}
	}
	if got := strings.TrimSpace(git(t, statePath, "rev-parse", "HEAD")); got != wantHead {
		t.Fatalf("fallback state HEAD=%s, want preserved local HEAD=%s", got, wantHead)
	}
}

func TestConcurrentStateReadyCreatorPullFailureNeverLooksReady(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	bare, workflowA, workflowB, stateBranch := twoHostStateWorkflow(t)
	stateA := filepath.Join(workflowA, ".spacedock-state")
	stateB := filepath.Join(workflowB, ".spacedock-state")

	if err := os.WriteFile(filepath.Join(stateA, "first-task.md"), []byte("---\nstatus: ideation\n---\n# Local divergence\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, stateA, "add", "first-task.md")
	git(t, stateA, "commit", "-q", "-m", "local divergence")
	wantLocalHead := strings.TrimSpace(git(t, stateA, "rev-parse", "HEAD"))

	if err := os.WriteFile(filepath.Join(stateB, "first-task.md"), []byte("---\nstatus: ideation\n---\n# Remote divergence\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, stateB, "add", "first-task.md")
	git(t, stateB, "commit", "-q", "-m", "remote divergence")
	git(t, stateB, "push", "origin", stateBranch)
	if got := strings.TrimSpace(git(t, bare, "rev-parse", stateBranch)); got == wantLocalHead {
		t.Fatal("precondition: remote and local state heads must diverge")
	}
	if err := os.RemoveAll(stateA); err != nil {
		t.Fatal(err)
	}

	rootA := filepath.Dir(filepath.Dir(workflowA))
	for _, result := range runConcurrentStateReadyAtAbsentBoundary(t, rootA, workflowA, 6) {
		if !strings.HasPrefix(result, "1|") || (!strings.Contains(result, "pull --rebase") && !strings.Contains(result, "did not converge")) {
			t.Fatalf("failed creator/waiter must report convergence failure, got %q", result)
		}
	}
	if _, err := os.Stat(stateA); !os.IsNotExist(err) {
		t.Fatalf("failed private convergence published a checkout: %v", err)
	}
	records, err := status.ParseWorktreePorcelainZ([]byte(git(t, rootA, "worktree", "list", "--porcelain", "-z")))
	if err != nil {
		t.Fatal(err)
	}
	registrations := 0
	prunable := false
	for _, record := range records {
		if status.RealpathOf(record.Path) == status.RealpathOf(stateA) {
			registrations++
			prunable = record.Prunable
		}
	}
	if registrations != 1 || !prunable {
		t.Fatalf("failed private convergence registration=(count=%d prunable=%v), want original stale registration", registrations, prunable)
	}
	if got := strings.TrimSpace(git(t, rootA, "rev-parse", stateBranch)); got != wantLocalHead {
		t.Fatalf("failed resume changed local branch: got %s want %s", got, wantLocalHead)
	}
}

func TestStateReadyPrivatePublishPreservesConcurrentStateWriter(t *testing.T) {
	testPrivateStateResumePreservesConcurrentStateWriter(t, "ready")
}

func TestStateInitPrivatePublishPreservesConcurrentStateWriter(t *testing.T) {
	testPrivateStateResumePreservesConcurrentStateWriter(t, "init")
}

func testPrivateStateResumePreservesConcurrentStateWriter(t *testing.T, verb string) {
	t.Setenv("HOME", t.TempDir())
	root, workflowDir, statePath := noOriginSplitWorkflow(t)
	entityBody := []byte("---\nstatus: ideation\n---\n# Concurrent writer\n")
	if err := os.WriteFile(filepath.Join(statePath, "first-task.md"), entityBody, 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, statePath, "add", "first-task.md")
	git(t, statePath, "commit", "-q", "-m", "seed writer entity")
	if err := os.RemoveAll(statePath); err != nil {
		t.Fatal(err)
	}

	reached := make(chan struct{})
	release := make(chan struct{})
	stateResumeBeforePublishHook = func(path string) {
		if status.RealpathOf(path) == status.RealpathOf(statePath) {
			close(reached)
			<-release
		}
	}
	t.Cleanup(func() { stateResumeBeforePublishHook = nil })

	type readyResult struct {
		code        int
		stdout, err string
	}
	result := make(chan readyResult, 1)
	go func() {
		var out, errBuf strings.Builder
		code := run(context.Background(), []string{"state", verb, "--workflow-dir", workflowDir},
			os.Environ(), root, nil, &out, &errBuf, &status.NativeRunner{}, nil)
		result <- readyResult{code: code, stdout: out.String(), err: errBuf.String()}
	}()
	<-reached

	if _, err := os.Lstat(statePath); !os.IsNotExist(err) {
		close(release)
		t.Fatalf("final state path was visible before private publication: %v", err)
	}
	if err := os.MkdirAll(statePath, 0o755); err != nil {
		close(release)
		t.Fatal(err)
	}
	entity := filepath.Join(statePath, "first-task.md")
	if err := os.WriteFile(entity, entityBody, 0o644); err != nil {
		close(release)
		t.Fatal(err)
	}
	var statusOut, statusErr strings.Builder
	code := run(context.Background(), []string{"status", "--workflow-dir", workflowDir, "--set", "first-task", "status=done"},
		os.Environ(), root, nil, &statusOut, &statusErr, &status.NativeRunner{}, nil)
	if code != 0 {
		close(release)
		t.Fatalf("concurrent status writer exit=%d stderr=%q", code, statusErr.String())
	}
	close(release)
	got := <-result
	if got.code != 1 || !strings.Contains(got.err, "appeared concurrently") {
		t.Fatalf("private publication should preserve concurrent writer: %+v", got)
	}
	if status.ParseFrontmatter(entity)["status"] != "done" {
		t.Fatalf("private publication deleted or reverted concurrent state writer: %q", string(mustReadFile(t, entity)))
	}
}

func TestStateReadyDoesNotDeleteConcurrentStateNew(t *testing.T) {
	testStateResumeDoesNotDeleteConcurrentStateNew(t, "ready")
}

func TestStateInitDoesNotDeleteConcurrentStateNew(t *testing.T) {
	testStateResumeDoesNotDeleteConcurrentStateNew(t, "init")
}

func testStateResumeDoesNotDeleteConcurrentStateNew(t *testing.T, verb string) {
	t.Setenv("HOME", t.TempDir())
	root, workflowDir := writeSplitReadmeRepo(t)
	statePath := filepath.Join(workflowDir, ".spacedock-state")
	reached := make(chan struct{})
	release := make(chan struct{})
	stateResumeBeforeRestoreHook = func(path string) {
		if status.RealpathOf(path) == status.RealpathOf(statePath) {
			close(reached)
			<-release
		}
	}
	t.Cleanup(func() { stateResumeBeforeRestoreHook = nil })

	type readyResult struct {
		code        int
		stdout, err string
	}
	result := make(chan readyResult, 1)
	go func() {
		var out, errBuf strings.Builder
		code := run(context.Background(), []string{"state", verb, "--workflow-dir", workflowDir},
			os.Environ(), root, nil, &out, &errBuf, &status.NativeRunner{}, nil)
		result <- readyResult{code: code, stdout: out.String(), err: errBuf.String()}
	}()
	<-reached
	newCode, newOut, newErr := execStateNew(t, root, workflowDir)
	if newCode != 0 {
		close(release)
		t.Fatalf("concurrent state new exit=%d stdout=%q stderr=%q", newCode, newOut, newErr)
	}
	branchHead := strings.TrimSpace(git(t, statePath, "rev-parse", "HEAD"))
	close(release)
	got := <-result
	if got.code != 1 || !strings.Contains(got.err, "appeared concurrently") {
		t.Fatalf("state ready should fail closed around concurrent birth: %+v", got)
	}
	if head := strings.TrimSpace(git(t, statePath, "rev-parse", "HEAD")); head != branchHead {
		t.Fatalf("resume mutated concurrent state new checkout: got HEAD %s want %s", head, branchHead)
	}
	if !gitignoreHasLine(string(mustReadFile(t, filepath.Join(root, ".gitignore"))), "docs/dev/.spacedock-state/") {
		t.Fatal("resume deleted concurrent state new's uncommitted code-branch edit")
	}
}

func TestStateReadyDoesNotDeleteConcurrentDirectWorktree(t *testing.T) {
	testStateResumeDoesNotDeleteConcurrentDirectWorktree(t, "ready")
}

func TestStateInitDoesNotDeleteConcurrentDirectWorktree(t *testing.T) {
	testStateResumeDoesNotDeleteConcurrentDirectWorktree(t, "init")
}

func testStateResumeDoesNotDeleteConcurrentDirectWorktree(t *testing.T, verb string) {
	t.Setenv("HOME", t.TempDir())
	root, workflowDir, statePath := noOriginSplitWorkflow(t)
	if err := os.RemoveAll(statePath); err != nil {
		t.Fatal(err)
	}
	reached := make(chan struct{})
	release := make(chan struct{})
	stateResumeBeforePublishHook = func(path string) {
		if status.RealpathOf(path) == status.RealpathOf(statePath) {
			close(reached)
			<-release
		}
	}
	t.Cleanup(func() { stateResumeBeforePublishHook = nil })

	type resumeResult struct {
		code int
		err  string
	}
	result := make(chan resumeResult, 1)
	go func() {
		var out, errBuf strings.Builder
		code := run(context.Background(), []string{"state", verb, "--workflow-dir", workflowDir},
			os.Environ(), root, nil, &out, &errBuf, &status.NativeRunner{}, nil)
		result <- resumeResult{code: code, err: errBuf.String()}
	}()
	<-reached
	git(t, root, "worktree", "add", "-q", "--force", "-b", "concurrent-state-writer", statePath, "HEAD")
	payload := filepath.Join(statePath, "uncommitted-writer-state")
	if err := os.WriteFile(payload, []byte("preserve me\n"), 0o644); err != nil {
		close(release)
		t.Fatal(err)
	}
	close(release)
	got := <-result
	if got.code != 1 || !strings.Contains(got.err, "appeared concurrently") {
		t.Fatalf("resume should fail closed around direct worktree publication: %+v", got)
	}
	if branch := strings.TrimSpace(git(t, statePath, "branch", "--show-current")); branch != "concurrent-state-writer" {
		t.Fatalf("resume replaced concurrent direct worktree branch %q", branch)
	}
	if body := string(mustReadFile(t, payload)); body != "preserve me\n" {
		t.Fatalf("resume changed concurrent uncommitted bytes: %q", body)
	}
}

func TestStateNewFromMainDoesNotPersistSharedExclude(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root, _, _ := noOriginSplitWorkflow(t)
	exclude, err := os.ReadFile(filepath.Join(root, ".git", "info", "exclude"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(exclude), "docs/dev/.spacedock-state/") {
		t.Fatalf("main-worktree birth left hidden shared exclude rule: %q", exclude)
	}
	tracked, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil || !strings.Contains(string(tracked), "docs/dev/.spacedock-state/") {
		t.Fatalf("main birth missing tracked ignore: err=%v body=%q", err, tracked)
	}
}

func TestLinkedWorktreeMetadataFailureFailsClosed(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root, _, statePath := noOriginSplitWorkflow(t)
	wantHead := strings.TrimSpace(git(t, statePath, "rev-parse", "HEAD"))
	wtRoot := addAgentWorktree(t, root, "broken-metadata")
	wtWorkflow := filepath.Join(wtRoot, "docs", "dev")
	if err := os.WriteFile(filepath.Join(wtRoot, ".git"), []byte("gitdir: /definitely/missing\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "ready", "--workflow-dir", wtWorkflow},
		os.Environ(), wtRoot, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 1 || !strings.Contains(errBuf.String(), "cannot resolve state checkout") {
		t.Fatalf("metadata failure must fail closed; exit=%d stdout=%q stderr=%q", code, out.String(), errBuf.String())
	}
	if _, err := os.Stat(filepath.Join(wtWorkflow, ".spacedock-state")); !os.IsNotExist(err) {
		t.Fatalf("metadata failure created nested state path: %v", err)
	}
	if got := strings.TrimSpace(git(t, statePath, "rev-parse", "HEAD")); got != wantHead {
		t.Fatalf("main state checkout mutated: got HEAD %s want %s", got, wantHead)
	}
}

func TestStateMutatorsRejectInvalidExistingCheckoutBeforeGitParentDiscovery(t *testing.T) {
	for _, occupant := range []string{"plain-directory", "wrong-branch-worktree"} {
		for _, verb := range []string{"ready", "init", "commit"} {
			t.Run(occupant+"/"+verb, func(t *testing.T) {
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
				git(t, root, "commit", "-q", "-m", "seed workflow")
				git(t, root, "branch", "spacedock-state/dev")
				statePath := filepath.Join(workflowDir, ".spacedock-state")
				if occupant == "wrong-branch-worktree" {
					git(t, root, "worktree", "add", "-q", "-b", "wrong-state", statePath, "HEAD")
				} else if err := os.MkdirAll(statePath, 0o755); err != nil {
					t.Fatal(err)
				}
				entity := filepath.Join(statePath, "first-task.md")
				if err := os.WriteFile(entity, []byte("---\nstatus: ideation\n---\n# Must not reach main\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				beforeHead := strings.TrimSpace(git(t, root, "rev-parse", "HEAD"))
				beforeStatus := git(t, root, "status", "--porcelain", "--untracked-files=all")

				args := []string{"state", verb, "--workflow-dir", workflowDir}
				if verb == "commit" {
					args = []string{"state", "commit", "first-task", "--workflow-dir", workflowDir, "-m", "must not commit main"}
				}
				var out, errBuf strings.Builder
				code := run(context.Background(), args, os.Environ(), root, nil, &out, &errBuf, &status.NativeRunner{}, nil)
				if code != 1 || !strings.Contains(errBuf.String(), "refusing invalid state checkout") {
					t.Fatalf("%s accepted %s: exit=%d stdout=%q stderr=%q", verb, occupant, code, out.String(), errBuf.String())
				}
				if afterHead := strings.TrimSpace(git(t, root, "rev-parse", "HEAD")); afterHead != beforeHead {
					t.Fatalf("%s moved main HEAD: got %s want %s", verb, afterHead, beforeHead)
				}
				if afterStatus := git(t, root, "status", "--porcelain", "--untracked-files=all"); afterStatus != beforeStatus {
					t.Fatalf("%s changed main worktree: before=%q after=%q", verb, beforeStatus, afterStatus)
				}
			})
		}
	}
}

func TestStatusAndSweepFromLinkedWorktreeUseMainState(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root, _, statePath := noOriginSplitWorkflow(t)
	entity := filepath.Join(statePath, "placement-task.md")
	if err := os.WriteFile(entity, []byte("---\nid: placement\ntitle: Placement\nstatus: implementation\npr: '#42'\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, statePath, "add", "placement-task.md")
	git(t, statePath, "commit", "-q", "-m", "seed placement task")
	wtRoot := addAgentWorktree(t, root, "placement")
	wtWorkflow := filepath.Join(wtRoot, "docs", "dev")

	var statusOut, statusErr strings.Builder
	code := run(context.Background(), []string{"status", "--workflow-dir", wtWorkflow, "--set", "placement-task", "status=ideation"},
		os.Environ(), wtRoot, nil, &statusOut, &statusErr, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("linked-worktree status mutation exit=%d stderr=%q", code, statusErr.String())
	}
	if got := status.ParseFrontmatter(entity)["status"]; got != "ideation" {
		t.Fatalf("main state entity status=%q, want ideation", got)
	}
	if _, err := os.Stat(filepath.Join(wtWorkflow, ".spacedock-state")); !os.IsNotExist(err) {
		t.Fatalf("status created or used nested checkout: %v", err)
	}
	var bootOut, bootErr strings.Builder
	code = run(context.Background(), []string{"status", "--workflow-dir", wtWorkflow, "--boot", "--json"},
		os.Environ(), wtRoot, nil, &bootOut, &bootErr, &status.NativeRunner{}, nil)
	if code != 0 || !json.Valid([]byte(bootOut.String())) {
		t.Fatalf("linked-worktree boot exit=%d stdout=%q stderr=%q", code, bootOut.String(), bootErr.String())
	}
	if !strings.Contains(bootOut.String(), statePath) || strings.Contains(bootOut.String(), filepath.Join(wtWorkflow, ".spacedock-state")) {
		t.Fatalf("boot did not report exact main-root entity directory: %s", bootOut.String())
	}

	var sweepOut, sweepErr strings.Builder
	code = dispatch.Sweep(wtWorkflow, func(string) (string, error) { return "MERGED", nil }, true, &sweepOut, &sweepErr)
	if code != 0 || !json.Valid([]byte(sweepOut.String())) {
		t.Fatalf("linked-worktree sweep exit=%d stdout=%q stderr=%q", code, sweepOut.String(), sweepErr.String())
	}
	if !strings.Contains(sweepOut.String(), `"slug": "placement-task"`) {
		t.Fatalf("sweep did not read main state checkout: %s", sweepOut.String())
	}
}

func TestStateNewFromLinkedWorktreeIgnoresPhysicalMainCheckout(t *testing.T) {
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
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".worktrees/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", "docs/dev/README.md", ".gitignore")
	git(t, root, "commit", "-q", "-m", "add workflow")
	wtRoot := addAgentWorktree(t, root, "birth")
	wtWorkflow := filepath.Join(wtRoot, "docs", "dev")

	code, _, stderr := execStateNew(t, wtRoot, wtWorkflow)
	if code != 0 {
		t.Fatalf("state new from linked worktree exit=%d stderr=%q", code, stderr)
	}
	if _, err := os.Stat(filepath.Join(workflowDir, ".spacedock-state")); err != nil {
		t.Fatalf("main-anchored state checkout absent: %v", err)
	}
	exclude, err := os.ReadFile(filepath.Join(root, ".git", "info", "exclude"))
	if err != nil || !strings.Contains(string(exclude), "/docs/dev/.spacedock-state/") {
		t.Fatalf("shared exclude does not protect physical checkout: err=%v body=%q", err, string(exclude))
	}
	trackedIgnore, err := os.ReadFile(filepath.Join(wtRoot, ".gitignore"))
	if err != nil || !strings.Contains(string(trackedIgnore), "docs/dev/.spacedock-state/") {
		t.Fatalf("tracked ignore missing on invoking branch: err=%v body=%q", err, string(trackedIgnore))
	}
	if got := strings.TrimSpace(git(t, root, "status", "--porcelain", "--untracked-files=all")); got != "" {
		t.Fatalf("main worktree sees physical state checkout as dirty: %q", got)
	}
	if got := strings.TrimSpace(git(t, wtRoot, "status", "--porcelain", "--untracked-files=all")); got != "M .gitignore" && got != " M .gitignore" {
		t.Fatalf("linked worktree should carry only the intentional .gitignore edit, got %q", got)
	}
	git(t, wtRoot, "add", ".gitignore")
	git(t, wtRoot, "commit", "-q", "-m", "ignore state checkout")
	if got := strings.TrimSpace(git(t, root, "status", "--porcelain", "--untracked-files=all")); got != "" {
		t.Fatalf("main worktree dirty after linked ignore commit: %q", got)
	}
	if got := strings.TrimSpace(git(t, wtRoot, "status", "--porcelain", "--untracked-files=all")); got != "" {
		t.Fatalf("linked worktree dirty after ignore commit: %q", got)
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
