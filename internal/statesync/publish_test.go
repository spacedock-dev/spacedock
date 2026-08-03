// ABOUTME: Real-git tests for the shared state publisher's ref safety and no-op outcome.
// ABOUTME: Invalid branch input cannot mutate refs; already-published HEAD is durable success.
package statesync

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
	}
	return strings.TrimSpace(string(out))
}

func publishedRepo(t *testing.T) (repo, bare, branch string) {
	t.Helper()
	repo = t.TempDir()
	bare = filepath.Join(t.TempDir(), "origin.git")
	testgit.InitRepo(t, repo, "-q")
	if err := os.WriteFile(filepath.Join(repo, "task.md"), []byte("state\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, repo, "add", "task.md")
	git(t, repo, "commit", "-q", "-m", "seed")
	branch = "spacedock-state/dev"
	git(t, repo, "branch", "-M", branch)
	git(t, repo, "init", "-q", "--bare", bare)
	git(t, repo, "remote", "add", "origin", bare)
	git(t, repo, "push", "-q", "origin", "HEAD:refs/heads/"+branch)
	git(t, repo, "fetch", "-q", "origin", branch+":refs/remotes/origin/"+branch)
	return repo, bare, branch
}

func TestPublishRejectsInvalidBranchBeforeRefMutation(t *testing.T) {
	repo, bare, _ := publishedRepo(t)
	headBefore := git(t, repo, "rev-parse", "HEAD")
	refsBefore := git(t, bare, "show-ref")
	statusBefore := git(t, repo, "status", "--porcelain=v1", "--untracked-files=all")

	outcome := Publish(repo, "--mirror")
	if outcome.Result != ResultFailed || !strings.Contains(outcome.Detail, "invalid state branch") {
		t.Fatalf("option-like branch must fail validation, got %#v", outcome)
	}
	if got := git(t, repo, "rev-parse", "HEAD"); got != headBefore {
		t.Fatalf("invalid branch moved HEAD: before=%s after=%s", headBefore, got)
	}
	if got := git(t, bare, "show-ref"); got != refsBefore {
		t.Fatalf("invalid branch mutated remote refs:\nbefore=%s\nafter=%s", refsBefore, got)
	}
	if got := git(t, repo, "status", "--porcelain=v1", "--untracked-files=all"); got != statusBefore {
		t.Fatalf("invalid branch mutated index/worktree:\nbefore=%q\nafter=%q", statusBefore, got)
	}
}

func TestPublishReturnsNoOpWhenOriginAlreadyEqualsHead(t *testing.T) {
	repo, _, branch := publishedRepo(t)
	if outcome := Publish(repo, branch); outcome.Result != ResultNoOp {
		t.Fatalf("already-published HEAD must be a no-op, got %#v", outcome)
	}
}

func TestPublishVerifiesNoOpWithoutRemoteTrackingRef(t *testing.T) {
	repo, _, branch := publishedRepo(t)
	git(t, repo, "update-ref", "-d", "refs/remotes/origin/"+branch)

	if outcome := Publish(repo, branch); outcome.Result != ResultNoOp {
		t.Fatalf("remote-equal HEAD with no tracking ref must still be a verified no-op, got %#v", outcome)
	}
}

func TestPublishRejectsWrongRootAndWrongBranchBeforeMutation(t *testing.T) {
	for _, kind := range []string{"nested root", "wrong branch"} {
		t.Run(kind, func(t *testing.T) {
			repo, bare, branch := publishedRepo(t)
			checkout := repo
			if kind == "nested root" {
				checkout = filepath.Join(repo, "nested")
				if err := os.Mkdir(checkout, 0o755); err != nil {
					t.Fatal(err)
				}
			} else {
				git(t, repo, "checkout", "-q", "-b", "wrong-state-branch")
			}
			headBefore := git(t, repo, "rev-parse", "HEAD")
			refsBefore := git(t, bare, "show-ref")
			statusBefore := git(t, repo, "status", "--porcelain=v1", "--untracked-files=all")

			if outcome := Publish(checkout, branch); outcome.Result != ResultFailed {
				t.Fatalf("%s must fail preflight, got %#v", kind, outcome)
			} else if kind == "nested root" && (!strings.Contains(outcome.Detail, checkout) || !strings.Contains(outcome.Detail, repo)) {
				t.Fatalf("nested-root refusal must name configured and observed roots, got %#v", outcome)
			} else if kind == "wrong branch" && (!strings.Contains(outcome.Detail, "wrong-state-branch") || !strings.Contains(outcome.Detail, "git branch -M "+branch)) {
				t.Fatalf("wrong-branch refusal must name actual branch and remedy, got %#v", outcome)
			}
			if got := git(t, repo, "rev-parse", "HEAD"); got != headBefore {
				t.Fatalf("%s moved HEAD: before=%s after=%s", kind, headBefore, got)
			}
			if got := git(t, bare, "show-ref"); got != refsBefore {
				t.Fatalf("%s mutated remote refs:\nbefore=%s\nafter=%s", kind, refsBefore, got)
			}
			if got := git(t, repo, "status", "--porcelain=v1", "--untracked-files=all"); got != statusBefore {
				t.Fatalf("%s mutated index/worktree:\nbefore=%q\nafter=%q", kind, statusBefore, got)
			}
		})
	}
}

func TestPreflightLeavesWrongBranchRebaseUntouched(t *testing.T) {
	repo, _, branch := publishedRepo(t)
	base := git(t, repo, "rev-parse", "HEAD")
	git(t, repo, "checkout", "-q", "-b", "wrong-target")
	if err := os.WriteFile(filepath.Join(repo, "task.md"), []byte("peer\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, repo, "commit", "-qam", "peer")
	git(t, repo, "checkout", "-q", "-b", "wrong-rebase", base)
	if err := os.WriteFile(filepath.Join(repo, "task.md"), []byte("local\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, repo, "commit", "-qam", "local")
	cmd := exec.Command("git", "-C", repo, "rebase", "wrong-target")
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("precondition: wrong-branch rebase should conflict\n%s", out)
	}
	headBefore := git(t, repo, "rev-parse", "HEAD")
	statusBefore := git(t, repo, "status", "--porcelain=v1", "--untracked-files=all")

	outcome := Preflight(repo, branch)
	if outcome.Result != ResultFailed || !strings.Contains(outcome.Detail, "left untouched") {
		t.Fatalf("wrong-branch rebase must refuse without abort, got %#v", outcome)
	}
	if got := git(t, repo, "rev-parse", "HEAD"); got != headBefore {
		t.Fatalf("preflight moved wrong-branch rebase HEAD: before=%s after=%s", headBefore, got)
	}
	if got := git(t, repo, "status", "--porcelain=v1", "--untracked-files=all"); got != statusBefore {
		t.Fatalf("preflight changed wrong-branch conflict:\nbefore=%q\nafter=%q", statusBefore, got)
	}
	rebasePath := git(t, repo, "rev-parse", "--git-path", "rebase-merge")
	if !filepath.IsAbs(rebasePath) {
		rebasePath = filepath.Join(repo, rebasePath)
	}
	if _, err := os.Stat(rebasePath); err != nil {
		t.Fatalf("preflight aborted wrong-branch rebase: %v", err)
	}
}
