// ABOUTME: Real-git proof that merge-guard archive finalization is durable on split-root state.
// ABOUTME: Covers two-host publication, archived resume, peer rebase, preflight HALT, and local-only.
package cli

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func mergeSentinelEntity(title string) string {
	return "---\nstatus: implementation\npr: pr-merge:99\n---\n# " + title + "\n"
}

func decodeOneJSON(t *testing.T, raw string) map[string]any {
	t.Helper()
	dec := json.NewDecoder(strings.NewReader(raw))
	var value map[string]any
	if err := dec.Decode(&value); err != nil {
		t.Fatalf("decode first JSON value: %v\n%s", err, raw)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		t.Fatalf("merge guard must emit exactly one JSON value followed by EOF; got %v (extra=%v)\n%s", err, extra, raw)
	}
	return value
}

func TestMergeGuardPublishesArchiveVisibleToFreshHost(t *testing.T) {
	bare, workflowA, workflowB, branch := twoHostStateWorkflow(t)
	hostA := filepath.Dir(filepath.Dir(workflowA))
	hostB := filepath.Dir(filepath.Dir(workflowB))
	checkoutA := filepath.Join(workflowA, ".spacedock-state")
	checkoutB := filepath.Join(workflowB, ".spacedock-state")

	writeEntity(t, workflowA, "first-task", mergeSentinelEntity("First Task"))
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "-m", "record merge sentinel"); code != 0 {
		t.Fatalf("push merge sentinel: exit=%d stderr=%q", code, errOut)
	}
	if code, _, errOut := runStateReadyCmd(t, hostB, workflowB); code != 0 {
		t.Fatalf("prepare fresh host: exit=%d stderr=%q", code, errOut)
	}
	writeEntity(t, workflowB, "peer-task", "---\nstatus: ideation\n---\n# Peer\n")
	if code, _, errOut := runStateCommitCmd(t, hostB, workflowB, "peer-task", "-m", "peer update"); code != 0 {
		t.Fatalf("push peer update: exit=%d stderr=%q", code, errOut)
	}

	stdout, stderr, code := runMergeCLI(t, hostA, "merge", "guard", "first-task", "--verdict", "passed", "--workflow-dir", workflowA, "--json")
	if code != 0 {
		t.Fatalf("merge guard: exit=%d stderr=%q", code, stderr)
	}
	result := decodeOneJSON(t, stdout)
	if result["result"] != "pushed" {
		t.Fatalf("origin-backed finalization must report pushed, got %#v", result)
	}
	localArchive := strings.TrimSpace(git(t, checkoutA, "rev-parse", "HEAD"))
	if origin := strings.TrimSpace(git(t, bare, "rev-parse", branch)); origin != localArchive {
		t.Fatalf("reported pushed before origin reached archive commit: local=%s origin=%s", localArchive, origin)
	}
	if merges := strings.TrimSpace(git(t, checkoutA, "log", "--merges", "--oneline")); merges != "" {
		t.Fatalf("peer integration must retain linear history:\n%s", merges)
	}

	if code, _, errOut := runStateReadyCmd(t, hostB, workflowB); code != 0 {
		t.Fatalf("fresh host state ready: exit=%d stderr=%q", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(checkoutB, "first-task.md")); !os.IsNotExist(err) {
		t.Fatalf("fresh host retained active merge sentinel: err=%v", err)
	}
	archived := filepath.Join(checkoutB, "_archive", "first-task.md")
	if body, err := os.ReadFile(archived); err != nil || !strings.Contains(string(body), "status: done") {
		t.Fatalf("fresh host missing terminal archive: body=%q err=%v", body, err)
	}
	for _, path := range []string{"_archive/first-task.md", "peer-task.md"} {
		if got := showOriginFile(t, bare, branch, path); got == "" {
			t.Fatalf("origin missing %s", path)
		}
	}
}

func TestArchivedStateCommitResumesInterruptedPublication(t *testing.T) {
	bare, workflowA, _, branch := twoHostStateWorkflow(t)
	hostA := filepath.Dir(filepath.Dir(workflowA))
	checkout := filepath.Join(workflowA, ".spacedock-state")

	body := strings.Replace(mergeSentinelEntity("First Task"), "pr: ", "worktree: .worktrees/first-task\npr: ", 1)
	writeEntity(t, workflowA, "first-task", body)
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "-m", "record merge sentinel"); code != 0 {
		t.Fatalf("push merge sentinel: exit=%d stderr=%q", code, errOut)
	}

	hooks := t.TempDir()
	prePush := filepath.Join(hooks, "pre-push")
	if err := os.WriteFile(prePush, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, checkout, "config", "core.hooksPath", hooks)
	stdout, _, code := runMergeCLI(t, hostA, "merge", "guard", "first-task", "--verdict", "passed", "--workflow-dir", workflowA)
	if code == 0 {
		t.Fatal("injected pre-push failure must leave a recoverable local archive commit")
	}
	if !strings.Contains(stdout, "State durability: unpublished") ||
		!strings.Contains(stdout, "Next: push; remove the worktree (`git worktree remove .worktrees/first-task`") {
		t.Fatalf("publication failure lost finalized durability/cleanup guidance:\n%s", stdout)
	}
	archiveHead := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD"))
	if _, err := os.Stat(filepath.Join(checkout, "_archive", "first-task.md")); err != nil {
		t.Fatalf("publication failure must retain local archive: %v", err)
	}
	git(t, checkout, "config", "--unset", "core.hooksPath")

	code, stdout, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "--json")
	if code != 0 {
		t.Fatalf("archived resume: exit=%d stderr=%q", code, errOut)
	}
	if got := decodeOneJSON(t, stdout)["result"]; got != "pushed" {
		t.Fatalf("archived resume must publish existing commit, got %v", got)
	}
	if head := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD")); head != archiveHead {
		t.Fatalf("archived resume created/replayed a commit: before=%s after=%s", archiveHead, head)
	}
	if origin := strings.TrimSpace(git(t, bare, "rev-parse", branch)); origin != archiveHead {
		t.Fatalf("archived resume did not publish archive commit: origin=%s archive=%s", origin, archiveHead)
	}
	if count := strings.TrimSpace(git(t, checkout, "rev-list", "--count", "--grep=archive first-task", "HEAD")); count != "1" {
		t.Fatalf("archive commit count=%s, want 1", count)
	}
	if code, stdout, errOut = runStateCommitCmd(t, hostA, workflowA, "first-task", "--json"); code != 0 || decodeOneJSON(t, stdout)["result"] != "no-op" {
		t.Fatalf("second archived resume must no-op: exit=%d stdout=%q stderr=%q", code, stdout, errOut)
	}
	if porcelain := strings.TrimSpace(git(t, checkout, "status", "--porcelain")); porcelain != "" {
		t.Fatalf("archived resume left checkout dirty: %s", porcelain)
	}
}

func TestMergeGuardAndStateCommitPreflightInterruptedArchiveRebase(t *testing.T) {
	bare, workflowA, workflowB, branch := twoHostStateWorkflow(t)
	hostA := filepath.Dir(filepath.Dir(workflowA))
	hostB := filepath.Dir(filepath.Dir(workflowB))
	checkoutB := filepath.Join(workflowB, ".spacedock-state")

	if err := os.MkdirAll(filepath.Join(checkoutB, "_archive"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(checkoutB, "first-task.md"), filepath.Join(checkoutB, "_archive", "first-task.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(checkoutB, "_archive", "first-task.md"),
		[]byte("---\nstatus: done\nverdict: passed\n---\n# Local archive\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(checkoutB, "guard-task.md"), []byte(mergeSentinelEntity("Guard Task")), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, checkoutB, "add", "--", "first-task.md", "_archive/first-task.md", "guard-task.md")
	git(t, checkoutB, "commit", "-q", "-m", "archive first-task (merge guard)", "--", "first-task.md", "_archive/first-task.md", "guard-task.md")
	localArchive := strings.TrimSpace(git(t, checkoutB, "rev-parse", "HEAD"))

	writeEntity(t, workflowA, "first-task", "---\nstatus: validation\n---\n# Peer edit\n")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "-m", "peer edits active task"); code != 0 {
		t.Fatalf("push peer conflict: exit=%d stderr=%q", code, errOut)
	}
	peer := strings.TrimSpace(git(t, bare, "rev-parse", "--short", branch))
	if _, ok := gitOK(t, checkoutB, "pull", "--rebase", "origin", branch); ok {
		t.Fatal("precondition: archive rename versus peer active edit must conflict")
	}

	stdout, errOut, code := runMergeCLI(t, hostB, "merge", "guard", "guard-task", "--verdict", "passed", "--workflow-dir", workflowB, "--json")
	if code != 3 {
		t.Fatalf("merge guard must HALT before mutation: exit=%d stdout=%q stderr=%q", code, stdout, errOut)
	}
	result := decodeOneJSON(t, stdout)
	reason, reasonOK := result["reason"].(string)
	if result["result"] != "halted" || !reasonOK || reason == "" {
		t.Fatalf("merge guard HALT JSON must carry outcome and reason, got %#v", result)
	}
	if _, err := os.Stat(filepath.Join(checkoutB, "guard-task.md")); err != nil {
		t.Fatalf("preflight HALT terminalized or moved guard-task: %v", err)
	}
	if _, err := os.Stat(filepath.Join(checkoutB, "_archive", "guard-task.md")); !os.IsNotExist(err) {
		t.Fatalf("preflight HALT archived guard-task: err=%v", err)
	}
	if head := strings.TrimSpace(git(t, checkoutB, "rev-parse", "HEAD")); head != localArchive {
		t.Fatalf("merge HALT did not restore local archive HEAD: got=%s want=%s", head, localArchive)
	}

	if _, ok := gitOK(t, checkoutB, "pull", "--rebase", "origin", branch); ok {
		t.Fatal("precondition: recreated archive conflict must fail")
	}
	code, stdout, errOut = runStateCommitCmd(t, hostB, workflowB, "first-task", "--json")
	if code != 3 {
		t.Fatalf("pre-existing rebase must HALT before archived resolution: exit=%d stdout=%q stderr=%q", code, stdout, errOut)
	}
	if !strings.Contains(errOut, "first-task.md") || !strings.Contains(errOut, "Peer commit: "+peer) {
		t.Fatalf("HALT must preserve path and peer evidence:\n%s", errOut)
	}
	if head := strings.TrimSpace(git(t, checkoutB, "rev-parse", "HEAD")); head != localArchive {
		t.Fatalf("rebase abort did not restore recoverable archive HEAD: got=%s want=%s", head, localArchive)
	}
	if rebasePath := strings.TrimSpace(git(t, checkoutB, "rev-parse", "--git-path", "rebase-merge")); func() bool {
		if !filepath.IsAbs(rebasePath) {
			rebasePath = filepath.Join(checkoutB, rebasePath)
		}
		_, err := os.Stat(rebasePath)
		return err == nil
	}() {
		t.Fatalf("HALT left rebase in progress at %s", rebasePath)
	}
}

func TestMergeGuardSiblingDirtNonFastForwardRemainsRecoverable(t *testing.T) {
	bare, workflowA, workflowB, branch := twoHostStateWorkflow(t)
	hostA := filepath.Dir(filepath.Dir(workflowA))
	hostB := filepath.Dir(filepath.Dir(workflowB))
	checkoutA := filepath.Join(workflowA, ".spacedock-state")

	writeEntity(t, workflowA, "sibling-task", "---\nstatus: ideation\n---\n# Sibling\n")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "sibling-task", "-m", "seed sibling"); code != 0 {
		t.Fatalf("seed sibling: exit=%d stderr=%q", code, errOut)
	}
	writeEntity(t, workflowA, "first-task", mergeSentinelEntity("First Task"))
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task", "-m", "record merge sentinel"); code != 0 {
		t.Fatalf("push merge sentinel: exit=%d stderr=%q", code, errOut)
	}
	if code, _, errOut := runStateReadyCmd(t, hostB, workflowB); code != 0 {
		t.Fatalf("prepare peer: exit=%d stderr=%q", code, errOut)
	}
	writeEntity(t, workflowB, "peer-task", "---\nstatus: ideation\n---\n# Peer\n")
	if code, _, errOut := runStateCommitCmd(t, hostB, workflowB, "peer-task", "-m", "peer update"); code != 0 {
		t.Fatalf("push peer update: exit=%d stderr=%q", code, errOut)
	}
	writeEntity(t, workflowA, "sibling-task", "---\nstatus: ideation\n---\n# Dirty sibling\n")

	stdout, stderr, code := runMergeCLI(t, hostA, "merge", "guard", "first-task", "--verdict", "passed", "--workflow-dir", workflowA, "--json")
	if code != 1 {
		t.Fatalf("dirty-sibling non-ff must retain archive with exit 1: exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if got := decodeOneJSON(t, stdout)["result"]; got != "unpublished" {
		t.Fatalf("publication failure must report unpublished, got %v", got)
	}
	archiveHead := strings.TrimSpace(git(t, checkoutA, "rev-parse", "HEAD"))
	if origin := strings.TrimSpace(git(t, bare, "rev-parse", branch)); origin == archiveHead {
		t.Fatalf("precondition: dirty sibling unexpectedly allowed archive publication")
	}
	if status := git(t, checkoutA, "status", "--porcelain"); !strings.Contains(status, "sibling-task.md") {
		t.Fatalf("publisher autostashed or lost sibling dirt:\n%s", status)
	}

	git(t, checkoutA, "restore", "--", "sibling-task.md")
	if code, _, errOut := runStateCommitCmd(t, hostA, workflowA, "first-task"); code != 0 {
		t.Fatalf("archived resume after settling sibling dirt: exit=%d stderr=%q", code, errOut)
	}
	if origin := strings.TrimSpace(git(t, bare, "rev-parse", branch)); origin != strings.TrimSpace(git(t, checkoutA, "rev-parse", "HEAD")) {
		t.Fatalf("settled archived resume did not publish: origin=%s local=%s", origin, archiveHead)
	}
	if got := showOriginFile(t, bare, branch, "peer-task.md"); !strings.Contains(got, "# Peer") {
		t.Fatalf("resume lost peer state: %q", got)
	}
}

func TestMergeGuardNoOriginReportsLocalOnly(t *testing.T) {
	_, workflowA, _, _ := twoHostStateWorkflow(t)
	hostA := filepath.Dir(filepath.Dir(workflowA))
	checkout := filepath.Join(workflowA, ".spacedock-state")

	writeEntity(t, workflowA, "first-task", mergeSentinelEntity("First Task"))
	git(t, checkout, "add", "--", "first-task.md")
	git(t, checkout, "commit", "-q", "-m", "record merge sentinel", "--", "first-task.md")
	git(t, checkout, "remote", "remove", "origin")

	stdout, stderr, code := runMergeCLI(t, hostA, "merge", "guard", "first-task", "--verdict", "passed", "--workflow-dir", workflowA, "--json")
	if code != 0 {
		t.Fatalf("no-origin merge guard: exit=%d stderr=%q", code, stderr)
	}
	if got := decodeOneJSON(t, stdout)["result"]; got != "local-only" {
		t.Fatalf("no-origin finalization must report local-only, got %v", got)
	}
	if _, err := os.Stat(filepath.Join(checkout, "_archive", "first-task.md")); err != nil {
		t.Fatalf("no-origin finalization must retain local archive commit: %v", err)
	}
}
