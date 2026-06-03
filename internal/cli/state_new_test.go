// ABOUTME: Real-git e2e for `spacedock state new` — birth the orphan state branch
// ABOUTME: + linked worktree around an already-present split-root README, no mocks.
package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// writeSplitReadmeRepo sets up the PRE-birth state `state new` onboards: a cloned
// repo whose code branch carries a split-root README at docs/dev, committed and
// pushed, but with NO .gitignore state entry and NO orphan branch yet. Returns the
// host clone root and its workflow dir. This is the input commission has NOT yet
// run the birth half on.
func writeSplitReadmeRepo(t *testing.T) (hostClone, workflowDir string) {
	t.Helper()
	bare := filepath.Join(t.TempDir(), "origin.git")
	git(t, t.TempDir(), "init", "-q", "--bare", bare)

	hostClone = filepath.Join(t.TempDir(), "host")
	git(t, t.TempDir(), "clone", "-q", bare, hostClone)
	git(t, hostClone, "config", "user.email", "t@t")
	git(t, hostClone, "config", "user.name", "t")

	workflowDir = filepath.Join(hostClone, "docs", "dev")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workflowDir, "README.md"), []byte(splitWorkflowReadme), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, hostClone, "add", "docs/dev/README.md")
	git(t, hostClone, "commit", "-q", "-m", "add split-root README")
	git(t, hostClone, "push", "-q", "origin", "HEAD")
	return hostClone, workflowDir
}

// execStateNew runs `spacedock state new --workflow-dir workflowDir` from dir,
// returning exit code, stdout, stderr.
func execStateNew(t *testing.T, dir, workflowDir string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errBuf strings.Builder
	code = run(context.Background(), []string{"state", "new", "--workflow-dir", workflowDir},
		os.Environ(), dir, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	return code, out.String(), errBuf.String()
}

// TestStateNewBirthsSplitRoot pins AC-1: one `state new` against a code branch
// carrying a split-root README births the orphan branch, appends the .gitignore
// entry, and checks out the linked state worktree. After it returns the orphan
// branch is in local refs, worktree list has the state path, .gitignore carries
// the path line, and boot shows split-root + present: true.
func TestStateNewBirthsSplitRoot(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	hostClone, workflowDir := writeSplitReadmeRepo(t)
	statePath := filepath.Join(workflowDir, ".spacedock-state")
	stateBranch := "spacedock-state/dev"

	code, stdout, stderr := execStateNew(t, hostClone, workflowDir)
	if code != 0 {
		t.Fatalf("state new exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}

	// Orphan branch present in local refs.
	if out, ok := gitOK(t, hostClone, "rev-parse", "--verify", "refs/heads/"+stateBranch); !ok {
		t.Fatalf("orphan branch %q not in local refs; rev-parse=%q", stateBranch, out)
	}
	// State path is a linked worktree.
	if out := git(t, hostClone, "worktree", "list"); !strings.Contains(out, statePath) {
		t.Fatalf("state path is not a linked worktree; worktree list=%q", out)
	}
	// .gitignore carries the state path as an exact line — not a substring of a
	// mangled path (a symlink-resolution bug once produced `../../…/docs/dev/.spacedock-state/`,
	// which a substring check would have passed).
	gi, err := os.ReadFile(filepath.Join(hostClone, ".gitignore"))
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}
	if !gitignoreHasLine(string(gi), "docs/dev/.spacedock-state/") {
		t.Fatalf(".gitignore missing the exact state path line; got %q", string(gi))
	}
	// The entry actually functions as an ignore rule (zero-churn R1: the state
	// checkout never shows in the code branch porcelain). git check-ignore exits 0
	// only when the path is ignored.
	if _, ok := gitOK(t, hostClone, "check-ignore", "docs/dev/.spacedock-state"); !ok {
		t.Fatalf("git check-ignore says the state path is NOT ignored — the .gitignore entry is wrong")
	}
	// Boot renders against the state checkout: split-root + present: true.
	bootOut := runStatusBoot(t, workflowDir)
	if !strings.Contains(bootOut, "STATE_BACKEND: split-root") {
		t.Fatalf("boot should show split-root; got\n%s", bootOut)
	}
	if !strings.Contains(bootOut, "present: true") {
		t.Fatalf("boot should show present: true after birth; got\n%s", bootOut)
	}
}

// TestStateNewRoundTripsWithStateInit pins AC-2: after `state new` on repo A and a
// code+orphan push, a fresh clone of A runs `state init` and lands in the same
// working state — worktree list has the state path, boot present: true, status
// renders. state new and state init are inverses reading the same README.
func TestStateNewRoundTripsWithStateInit(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	hostClone, workflowDir := writeSplitReadmeRepo(t)

	code, stdout, stderr := execStateNew(t, hostClone, workflowDir)
	if code != 0 {
		t.Fatalf("state new exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	// Commit the .gitignore edit (the operator's step) and push the code branch so
	// a fresh clone sees the split-root README and ignore entry. The orphan branch
	// was best-effort pushed by state new.
	git(t, hostClone, "add", ".gitignore")
	git(t, hostClone, "commit", "-q", "-m", "gitignore state path")
	git(t, hostClone, "push", "-q", "origin", "HEAD")

	// Fresh clone resumes via state init.
	origin := git(t, hostClone, "remote", "get-url", "origin")
	bare := strings.TrimSpace(origin)
	fresh := filepath.Join(t.TempDir(), "fresh")
	git(t, t.TempDir(), "clone", "-q", bare, fresh)
	git(t, fresh, "config", "user.email", "t@t")
	git(t, fresh, "config", "user.name", "t")
	freshWorkflow := filepath.Join(fresh, "docs", "dev")
	freshState := filepath.Join(freshWorkflow, ".spacedock-state")

	if _, err := os.Stat(freshState); !os.IsNotExist(err) {
		t.Fatalf("fresh clone should NOT have the state path yet (err=%v)", err)
	}

	var out, errBuf strings.Builder
	initCode := run(context.Background(), []string{"state", "init", "--workflow-dir", freshWorkflow},
		os.Environ(), fresh, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if initCode != 0 {
		t.Fatalf("state init exit=%d stdout=%q stderr=%q", initCode, out.String(), errBuf.String())
	}
	if _, err := os.Stat(freshState); err != nil {
		t.Fatalf("state init did not create the state worktree: %v", err)
	}
	bootOut := runStatusBoot(t, freshWorkflow)
	if !strings.Contains(bootOut, "present: true") {
		t.Fatalf("post-init boot should show present: true; got\n%s", bootOut)
	}
}

// TestStateNewRefusesAlreadyBirthed pins AC-3: state new refuses (a) an occupied
// state path, (b) a local orphan branch, and (c) an orphan on origin — each
// non-zero with a tailored message, case (c) naming `state init`. No raw git
// `already exists` / `fatal:` leaks to stderr.
func TestStateNewRefusesAlreadyBirthed(t *testing.T) {
	t.Run("occupied-state-path", func(t *testing.T) {
		t.Setenv("HOME", t.TempDir())
		hostClone, workflowDir := writeSplitReadmeRepo(t)
		// First birth succeeds.
		if code, _, stderr := execStateNew(t, hostClone, workflowDir); code != 0 {
			t.Fatalf("first state new should succeed; exit=%d stderr=%q", code, stderr)
		}
		// Second birth: the path now exists → refusal.
		code, _, stderr := execStateNew(t, hostClone, workflowDir)
		if code == 0 {
			t.Fatalf("second state new on an occupied path must refuse (non-zero)")
		}
		assertNoGitFatal(t, stderr)
		if !strings.Contains(stderr, "already") {
			t.Fatalf("occupied-path refusal should mention the path already exists; stderr=%q", stderr)
		}
	})

	t.Run("local-orphan-branch", func(t *testing.T) {
		t.Setenv("HOME", t.TempDir())
		hostClone, workflowDir := writeSplitReadmeRepo(t)
		// Birth the orphan branch in a temp worktree but DO NOT check out the state
		// path, so only the local-orphan precheck (not the path guard) trips.
		tmpWT := filepath.Join(t.TempDir(), "orphan")
		git(t, hostClone, "worktree", "add", "--detach", tmpWT)
		git(t, tmpWT, "checkout", "--orphan", "spacedock-state/dev")
		git(t, tmpWT, "rm", "-rf", "--cached", ".")
		git(t, tmpWT, "commit", "-q", "--allow-empty", "-m", "seed")
		git(t, hostClone, "worktree", "remove", "--force", tmpWT)

		code, _, stderr := execStateNew(t, hostClone, workflowDir)
		if code == 0 {
			t.Fatalf("state new with a local orphan branch must refuse (non-zero)")
		}
		assertNoGitFatal(t, stderr)
		if !strings.Contains(stderr, "state init") {
			t.Fatalf("local-orphan refusal should redirect to `state init`; stderr=%q", stderr)
		}
	})

	t.Run("remote-orphan-branch", func(t *testing.T) {
		t.Setenv("HOME", t.TempDir())
		hostClone, workflowDir := writeSplitReadmeRepo(t)
		// Push an orphan branch to origin from a temp worktree, then drop the local
		// ref so only the remote ls-remote precheck trips.
		tmpWT := filepath.Join(t.TempDir(), "orphan")
		git(t, hostClone, "worktree", "add", "--detach", tmpWT)
		git(t, tmpWT, "checkout", "--orphan", "spacedock-state/dev")
		git(t, tmpWT, "rm", "-rf", "--cached", ".")
		git(t, tmpWT, "commit", "-q", "--allow-empty", "-m", "seed")
		git(t, tmpWT, "push", "-q", "origin", "spacedock-state/dev")
		git(t, hostClone, "worktree", "remove", "--force", tmpWT)
		git(t, hostClone, "branch", "-D", "spacedock-state/dev")

		code, _, stderr := execStateNew(t, hostClone, workflowDir)
		if code == 0 {
			t.Fatalf("state new with a remote orphan branch must refuse (non-zero)")
		}
		assertNoGitFatal(t, stderr)
		if !strings.Contains(stderr, "state init") {
			t.Fatalf("remote-orphan refusal should name `state init` as the right command; stderr=%q", stderr)
		}
	})
}

// TestStateNewNoRemoteBestEffort pins AC-4: in a repo with no origin, state new
// births the orphan branch + worktree locally, exits 0, and warns the operator to
// push the orphan branch to share state.
func TestStateNewNoRemoteBestEffort(t *testing.T) {
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

	code, stdout, stderr := execStateNew(t, root, workflowDir)
	if code != 0 {
		t.Fatalf("no-remote state new must exit 0 (best-effort push); exit=%d stderr=%q", code, stderr)
	}
	assertNoGitFatal(t, stderr)
	statePath := filepath.Join(workflowDir, ".spacedock-state")
	if _, err := os.Stat(statePath); err != nil {
		t.Fatalf("no-remote state new did not create the state worktree: %v", err)
	}
	combined := stdout + stderr
	if !strings.Contains(combined, "push") {
		t.Fatalf("no-remote state new should warn the operator to push the orphan branch; got stdout=%q stderr=%q", stdout, stderr)
	}
}

// TestStateNewInlineErrors pins AC-5: state new against an inline / non-split
// README exits non-zero with a message explaining state new only births split-root
// workflows. (state init prints a benign one-liner; new is a mutating verb, so a
// mismatched backend is an operator error.)
func TestStateNewInlineErrors(t *testing.T) {
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
	git(t, root, "init", "-q")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "init")

	code, _, stderr := execStateNew(t, root, wf)
	if code == 0 {
		t.Fatalf("inline state new must exit non-zero (mutating verb on a non-split backend)")
	}
	if !strings.Contains(stderr, "split-root") {
		t.Fatalf("inline state new should explain it only births split-root workflows; stderr=%q", stderr)
	}
}

// gitignoreHasLine reports whether content carries want as a complete line, so a
// mangled path that merely contains want as a substring does not pass.
func gitignoreHasLine(content, want string) bool {
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) == want {
			return true
		}
	}
	return false
}

// assertNoGitFatal fails if stderr leaked a raw git porcelain fatal — state new
// must translate git failures into operator-facing messages, never pass through
// `fatal:` or `already exists`.
func assertNoGitFatal(t *testing.T, stderr string) {
	t.Helper()
	for _, leak := range []string{"fatal:", "already exists"} {
		if strings.Contains(stderr, leak) {
			t.Fatalf("stderr leaked raw git output %q: %s", leak, stderr)
		}
	}
}

// TestLsRemoteHasBranch is the pure table test for the remote-orphan precheck's
// ls-remote parser: it must match the full refs/heads/<branch> ref-path (branch
// names carry slashes), treat empty output as no-match, and not be fooled by a
// similarly-prefixed branch.
func TestLsRemoteHasBranch(t *testing.T) {
	cases := []struct {
		name   string
		output string
		branch string
		want   bool
	}{
		{"empty output", "", "spacedock-state/dev", false},
		{"exact match", "abc123\trefs/heads/spacedock-state/dev\n", "spacedock-state/dev", true},
		{"no trailing newline", "abc123\trefs/heads/spacedock-state/dev", "spacedock-state/dev", true},
		{"different branch", "abc123\trefs/heads/main\n", "spacedock-state/dev", false},
		{"prefix is not a match", "abc123\trefs/heads/spacedock-state/dev-2\n", "spacedock-state/dev", false},
		{"match among several heads", "a\trefs/heads/main\nb\trefs/heads/spacedock-state/dev\n", "spacedock-state/dev", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := lsRemoteHasBranch(tc.output, tc.branch); got != tc.want {
				t.Fatalf("lsRemoteHasBranch(%q, %q) = %v, want %v", tc.output, tc.branch, got, tc.want)
			}
		})
	}
}
