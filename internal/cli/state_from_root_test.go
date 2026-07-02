// ABOUTME: Real-git e2e — bare `state ready/sweep/commit/init/new` (no --workflow-dir)
// ABOUTME: auto-discover the lone nested workflow from the repo toplevel (AC-1, AC-2).
package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// TestStateReadyFromRootResolves pins AC-1: bare `state ready`, run with the repo
// toplevel as cwd and no --workflow-dir, resolves the single nested workflow via
// the downward-scan fallback and performs the same boot-pull integration as the
// explicit --workflow-dir path.
func TestStateReadyFromRootResolves(t *testing.T) {
	_, workflowA, workflowB, _ := twoHostStateWorkflow(t)
	checkoutB := filepath.Join(workflowB, ".spacedock-state")
	hostA := filepath.Dir(filepath.Dir(workflowA))
	hostB := filepath.Dir(filepath.Dir(workflowB))

	writeEntity(t, workflowA, "alpha-task", "---\nstatus: ideation\n---\n# Alpha (A)\n")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "alpha-task", "-m", "A: add alpha"); code != 0 {
		t.Fatalf("A's commit should succeed; exit=%d stderr=%q", code, errOut)
	}

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "ready"},
		os.Environ(), hostB, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("bare state ready from root should exit 0; got exit=%d stderr=%q", code, errBuf.String())
	}
	if _, err := os.Stat(filepath.Join(checkoutB, "alpha-task.md")); err != nil {
		t.Fatalf("bare state ready from root should have integrated A's alpha-task: %v", err)
	}
}

// TestStateSweepFromRootResolves pins AC-1: bare `state sweep` from the repo
// toplevel resolves the nested workflow and stays read-only (HEAD unchanged),
// mirroring TestStateSweepIsReadOnly with no --workflow-dir.
func TestStateSweepFromRootResolves(t *testing.T) {
	_, workflowA, _, _ := twoHostStateWorkflow(t)
	checkoutA := filepath.Join(workflowA, ".spacedock-state")
	hostA := filepath.Dir(filepath.Dir(workflowA))

	headBefore := strings.TrimSpace(git(t, checkoutA, "rev-parse", "HEAD"))

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "sweep", "--json"},
		os.Environ(), hostA, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("bare state sweep from root should exit 0; got exit=%d stderr=%q", code, errBuf.String())
	}
	if !strings.Contains(out.String(), `"command": "state sweep"`) {
		t.Fatalf("sweep --json should carry the command envelope; json:\n%s", out.String())
	}
	headAfter := strings.TrimSpace(git(t, checkoutA, "rev-parse", "HEAD"))
	if headAfter != headBefore {
		t.Fatalf("bare state sweep from root must be read-only; HEAD changed %s -> %s", headBefore, headAfter)
	}
}

// TestStateCommitFromRootResolves pins AC-1: bare `state commit <slug>` from the
// repo toplevel resolves the nested workflow and lands the path-scoped commit in
// the state checkout's git log, same as the explicit --workflow-dir path.
func TestStateCommitFromRootResolves(t *testing.T) {
	bare, workflowA, _, stateBranch := twoHostStateWorkflow(t)
	checkoutA := filepath.Join(workflowA, ".spacedock-state")
	hostA := filepath.Dir(filepath.Dir(workflowA))

	writeEntity(t, workflowA, "first-task", "---\nstatus: implementation\n---\n# First Task (A)\n")

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "commit", "first-task", "-m", "A: from root"},
		os.Environ(), hostA, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("bare state commit from root should exit 0; got exit=%d stderr=%q", code, errBuf.String())
	}
	names := git(t, checkoutA, "show", "--name-only", "--pretty=format:", "HEAD")
	if !strings.Contains(names, "first-task.md") {
		t.Fatalf("bare commit from root should include first-task.md; name-only:\n%s", names)
	}
	originFirst := showOriginFile(t, bare, stateBranch, "first-task.md")
	if !strings.Contains(originFirst, "status: implementation") {
		t.Fatalf("origin first-task should carry the from-root commit; got:\n%s", originFirst)
	}
}

// TestStateInitFromRootResolves pins the AC-1 extension to `state init`: a fresh
// clone with an absent state checkout resumes it via bare `state init` run from
// the repo toplevel, no --workflow-dir.
func TestStateInitFromRootResolves(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	bare := filepath.Join(t.TempDir(), "origin.git")
	git(t, t.TempDir(), "init", "-q", "--bare", bare)

	hostClone := filepath.Join(t.TempDir(), "host")
	git(t, t.TempDir(), "clone", "-q", bare, hostClone)
	git(t, hostClone, "config", "user.email", "t@t")
	git(t, hostClone, "config", "user.name", "t")
	commissionSplitWorkflow(t, hostClone)
	git(t, hostClone, "push", "-q", "origin", "HEAD")

	fresh := filepath.Join(t.TempDir(), "fresh")
	git(t, t.TempDir(), "clone", "-q", bare, fresh)
	git(t, fresh, "config", "user.email", "t@t")
	git(t, fresh, "config", "user.name", "t")
	freshState := filepath.Join(fresh, "docs", "dev", ".spacedock-state")

	if _, err := os.Stat(freshState); !os.IsNotExist(err) {
		t.Fatalf("precondition: fresh clone should NOT have the state path yet (err=%v)", err)
	}

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "init"},
		os.Environ(), fresh, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("bare state init from root should exit 0; got exit=%d stderr=%q", code, errBuf.String())
	}
	if _, err := os.Stat(freshState); err != nil {
		t.Fatalf("bare state init from root should have created the state worktree: %v", err)
	}
}

// TestStateNewFromRootResolves pins the AC-1 extension to `state new`: an
// already-present split-root README (not yet birthed) is birthed by bare
// `state new` run from the repo toplevel, no --workflow-dir.
func TestStateNewFromRootResolves(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	hostClone, workflowDir := writeSplitReadmeRepo(t)
	statePath := filepath.Join(workflowDir, ".spacedock-state")

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "new"},
		os.Environ(), hostClone, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("bare state new from root should exit 0; got exit=%d stdout=%q stderr=%q", code, out.String(), errBuf.String())
	}
	if _, err := os.Stat(statePath); err != nil {
		t.Fatalf("bare state new from root should have created the state worktree: %v", err)
	}
}

// TestStateVerbsFromRootRefuseAmbiguousWorkflows pins AC-2: with two commissioned
// workflows nested under a plain toplevel, the bare state verbs refuse (non-zero
// exit) and stderr names both candidate dirs and instructs --workflow-dir — the
// same refusal spacedock new emits (TestNewFromRootMultiWorkflowRefuses).
func TestStateVerbsFromRootRefuseAmbiguousWorkflows(t *testing.T) {
	root := t.TempDir()
	for _, rel := range []string{"docs/dev", "docs/ops"} {
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(full, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(full, "README.md"), []byte(splitWorkflowReadme), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git(t, root, "init", "-q")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "two commissioned workflows")

	dev := filepath.Join(root, "docs", "dev")
	ops := filepath.Join(root, "docs", "ops")
	if real, err := filepath.EvalSymlinks(dev); err == nil {
		dev = real
	}
	if real, err := filepath.EvalSymlinks(ops); err == nil {
		ops = real
	}

	for _, args := range [][]string{
		{"state", "sweep"},
		{"state", "ready"},
		{"state", "commit", "any-slug", "-m", "msg"},
	} {
		var out, errBuf strings.Builder
		code := run(context.Background(), args, os.Environ(), root, nil, &out, &errBuf, &status.NativeRunner{}, nil)
		if code == 0 {
			t.Fatalf("bare %v with two workflows must refuse (non-zero)", args)
		}
		errOut := errBuf.String()
		if !strings.Contains(errOut, dev) || !strings.Contains(errOut, ops) {
			t.Fatalf("bare %v ambiguity error must name both candidates, got: %q", args, errOut)
		}
		if !strings.Contains(errOut, "--workflow-dir") {
			t.Fatalf("bare %v ambiguity error must instruct --workflow-dir, got: %q", args, errOut)
		}
	}
}
