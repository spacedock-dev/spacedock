// ABOUTME: Real-git e2e for `state ready` — boot pull integration + the same-entity
// ABOUTME: boot-conflict HALT (AC-5), inline no-op and absent-checkout resume (AC-6).
package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// runStateReadyCmd runs `spacedock state ready` for workflowDir and returns the
// exit code plus captured stdout/stderr.
func runStateReadyCmd(t *testing.T, hostDir, workflowDir string, extra ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errBuf strings.Builder
	args := append([]string{"state", "ready", "--workflow-dir", workflowDir}, extra...)
	code = run(context.Background(), args, os.Environ(), hostDir, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	return code, out.String(), errBuf.String()
}

// TestStateReadyIntegratesPeerState pins AC-5 clean case: A commits+pushes a new
// entity via the verb; B runs `state ready`, whose single pull --rebase integrates
// A's commit (exit 0) and leaves A's entity present in B's checkout.
func TestStateReadyIntegratesPeerState(t *testing.T) {
	_, workflowA, workflowB, _ := twoHostStateWorkflow(t)
	checkoutB := filepath.Join(workflowB, ".spacedock-state")
	hostA := filepath.Dir(filepath.Dir(workflowA))
	hostB := filepath.Dir(filepath.Dir(workflowB))

	// A pushes a distinct new entity.
	writeEntity(t, workflowA, "alpha-task", "---\nstatus: ideation\n---\n# Alpha (A)\n")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "alpha-task", "-m", "A: add alpha"); code != 0 {
		t.Fatalf("A's commit should succeed; exit=%d stderr=%q", code, errOut)
	}

	// B's checkout does not yet have alpha-task; `state ready` pulls it in.
	if _, err := os.Stat(filepath.Join(checkoutB, "alpha-task.md")); err == nil {
		t.Fatalf("precondition: B should NOT yet have alpha-task before ready")
	}
	code, _, errOut := runStateReadyCmd(t, hostB, workflowB)
	if code != 0 {
		t.Fatalf("state ready clean integration should exit 0; got exit=%d stderr=%q", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(checkoutB, "alpha-task.md")); err != nil {
		t.Fatalf("state ready should pull A's alpha-task into B's checkout: %v", err)
	}
}

// TestStateReadyHaltsOnBootConflict pins AC-5 conflict case: B has an unpushed local
// commit on the SAME entity A pushed; `state ready`'s boot pull --rebase conflicts
// and must HALT identically to commit — exit 3, checkout left clean (rebase aborted).
func TestStateReadyHaltsOnBootConflict(t *testing.T) {
	bare, workflowA, workflowB, stateBranch := twoHostStateWorkflow(t)
	checkoutB := filepath.Join(workflowB, ".spacedock-state")
	hostA := filepath.Dir(filepath.Dir(workflowA))

	// A edits first-task and pushes.
	writeEntity(t, workflowA, "first-task", "---\nstatus: implementation\n---\n# First Task (A)\n")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "-m", "A: -> implementation"); code != 0 {
		t.Fatalf("A's commit should succeed; exit=%d stderr=%q", code, errOut)
	}

	// B makes a CONFLICTING local commit on the same entity (committed, NOT pushed),
	// then runs `state ready` — the boot pull --rebase conflicts.
	writeEntity(t, workflowB, "first-task", "---\nstatus: review\n---\n# First Task (B)\n")
	git(t, checkoutB, "add", "first-task.md")
	git(t, checkoutB, "commit", "-q", "-m", "B: -> review", "--", "first-task.md")
	hostB := filepath.Dir(filepath.Dir(workflowB))

	code, _, errOut := runStateReadyCmd(t, hostB, workflowB)
	if code != 3 {
		t.Fatalf("same-entity boot conflict must HALT with exit 3; got exit=%d stderr=%q", code, errOut)
	}
	if porcelain := git(t, checkoutB, "status", "--porcelain"); strings.TrimSpace(porcelain) != "" {
		t.Fatalf("ready HALT must leave a clean checkout; porcelain:\n%s", porcelain)
	}
	// A's edit survives on origin.
	originFirst := showOriginFile(t, bare, stateBranch, "first-task.md")
	if !strings.Contains(originFirst, "status: implementation") {
		t.Fatalf("origin first-task should carry A's edit after ready HALT; got:\n%s", originFirst)
	}
}

// TestStateReadyHaltStderrCarriesRemediationAndPeerCommit pins AC-2 (D1) for the
// boot-conflict HALT path: identical remediation to `state commit`'s exit-3 halt
// — the peer commit, the FO's next-action line, and the never-force line.
func TestStateReadyHaltStderrCarriesRemediationAndPeerCommit(t *testing.T) {
	_, workflowA, workflowB, _ := twoHostStateWorkflow(t)
	checkoutA := filepath.Join(workflowA, ".spacedock-state")
	checkoutB := filepath.Join(workflowB, ".spacedock-state")
	hostA := filepath.Dir(filepath.Dir(workflowA))
	hostB := filepath.Dir(filepath.Dir(workflowB))

	writeEntity(t, workflowA, "first-task", "---\nstatus: implementation\n---\n# First Task (A)\n")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "-m", "A: -> implementation"); code != 0 {
		t.Fatalf("A's commit should succeed; exit=%d stderr=%q", code, errOut)
	}
	peerSHA := strings.TrimSpace(git(t, checkoutA, "rev-parse", "--short", "HEAD"))

	writeEntity(t, workflowB, "first-task", "---\nstatus: review\n---\n# First Task (B)\n")
	git(t, checkoutB, "add", "first-task.md")
	git(t, checkoutB, "commit", "-q", "-m", "B: -> review", "--", "first-task.md")

	code, _, errOut := runStateReadyCmd(t, hostB, workflowB)
	if code != 3 {
		t.Fatalf("same-entity boot conflict must HALT with exit 3; got exit=%d stderr=%q", code, errOut)
	}
	if !strings.Contains(errOut, "Peer commit: "+peerSHA) {
		t.Fatalf("ready HALT stderr should name the peer commit %q, got:\n%s", peerSHA, errOut)
	}
	if !strings.Contains(errOut, "Next: HALT dispatch — do not dispatch against this state tree. Surface the conflicting path(s) and peer commit to the operator and stop.") {
		t.Fatalf("ready HALT stderr should name the FO's next action, got:\n%s", errOut)
	}
	if !strings.Contains(errOut, "Never `git push --force`/`--force-with-lease`; never re-run with `-X ours`/`-X theirs`; never discard either side.") {
		t.Fatalf("ready HALT stderr should carry the never-force/never-auto-resolve line, got:\n%s", errOut)
	}
}

// TestStateReadyInlineNoOp pins AC-6 inline case: an inline workflow is a clean
// no-op (exit 0, nothing to sync, no git network op).
func TestStateReadyInlineNoOp(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := t.TempDir()
	wf := filepath.Join(root, "wf")
	if err := os.MkdirAll(wf, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wf, "README.md"),
		[]byte("---\ncommissioned-by: spacedock@1\nid-style: slug\nstate: $inline\nstages:\n  states:\n    - name: ideation\n      initial: true\n    - name: done\n      terminal: true\n---\n\n# Inline WF\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	testgit.InitRepo(t, root, "-q")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "init")

	code, stdout, errOut := runStateReadyCmd(t, root, wf)
	if code != 0 {
		t.Fatalf("inline state ready should be a clean no-op (exit 0); got exit=%d stderr=%q", code, errOut)
	}
	if stdout == "" {
		t.Fatalf("inline state ready should print a one-liner; stdout empty")
	}
}

// TestStateReadyResumesAbsentCheckout pins AC-6 resume case: a split-root checkout
// that is absent on a fresh clone is resumed by `state ready` (present afterward),
// reusing the `state init` fetch + worktree-add path.
func TestStateReadyResumesAbsentCheckout(t *testing.T) {
	bare := filepath.Join(t.TempDir(), "origin.git")
	git(t, t.TempDir(), "init", "-q", "--bare", bare)

	hostA := filepath.Join(t.TempDir(), "hostA")
	git(t, t.TempDir(), "clone", "-q", bare, hostA)
	git(t, hostA, "config", "user.email", "a@t")
	git(t, hostA, "config", "user.name", "hostA")
	commissionSplitWorkflow(t, hostA)
	git(t, hostA, "push", "-q", "origin", "HEAD")

	// Fresh clone: code branch only, state checkout absent.
	fresh := filepath.Join(t.TempDir(), "fresh")
	git(t, t.TempDir(), "clone", "-q", bare, fresh)
	git(t, fresh, "config", "user.email", "f@t")
	git(t, fresh, "config", "user.name", "fresh")
	freshWorkflow := filepath.Join(fresh, "docs", "dev")
	freshState := filepath.Join(freshWorkflow, ".spacedock-state")

	if _, err := os.Stat(freshState); !os.IsNotExist(err) {
		t.Fatalf("precondition: fresh clone should NOT yet have the state checkout (err=%v)", err)
	}
	code, stdout, errOut := runStateReadyCmd(t, fresh, freshWorkflow)
	if code != 0 {
		t.Fatalf("state ready on an absent checkout should resume it (exit 0); got exit=%d stderr=%q", code, errOut)
	}
	if _, err := os.Stat(freshState); err != nil {
		t.Fatalf("state ready should have resumed the absent checkout: %v", err)
	}
	// D1(c): the resume path carries the re-boot-after-resume sequencing the
	// «state.ensure-ready» prose used to own — the FO must re-invoke the boot read
	// before the greet since the checkout was just linked.
	if !strings.Contains(stdout, "checkout resumed — re-run `spacedock status --boot` before the greet.") {
		t.Fatalf("resumed checkout should print the re-boot-before-greet line; stdout:\n%s", stdout)
	}
}
