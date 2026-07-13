// ABOUTME: Adversarial Git porcelain coverage for main-worktree placement.
package status

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func placementGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git -C %q %v: %v\n%s", dir, args, err, out)
	}
	return string(out)
}

func TestPrimaryWorktreePorcelainZPreservesNewlinePath(t *testing.T) {
	want := "/tmp/main\nroot\\literal"
	raw := []byte("worktree " + want + "\x00HEAD deadbeef\x00branch refs/heads/main\x00\x00" +
		"worktree /tmp/linked\x00HEAD cafe\x00detached\x00\x00")
	got, err := primaryWorktreeFromPorcelainZ(raw)
	if err != nil {
		t.Fatalf("primaryWorktreeFromPorcelainZ: %v", err)
	}
	if got != want {
		t.Fatalf("primary path bytes changed: got %q want %q", got, want)
	}
}

func TestPrimaryWorktreePorcelainZRejectsBarePrimary(t *testing.T) {
	raw := []byte("worktree /tmp/repo.git\x00bare\x00\x00" +
		"worktree /tmp/linked\x00HEAD cafe\x00branch refs/heads/main\x00\x00")
	if _, err := primaryWorktreeFromPorcelainZ(raw); err == nil || !strings.Contains(err.Error(), "bare") {
		t.Fatalf("bare primary should fail closed, got %v", err)
	}
}

func TestParseWorktreePorcelainZRejectsIncompleteRecord(t *testing.T) {
	if _, err := ParseWorktreePorcelainZ([]byte("worktree /tmp/main\x00HEAD deadbeef\x00")); err == nil || !strings.Contains(err.Error(), "incomplete") {
		t.Fatalf("incomplete record should fail closed, got %v", err)
	}
}

func TestResolveRepositoryPlacementFromNewlineWorktree(t *testing.T) {
	base := t.TempDir()
	mainRoot := filepath.Join(base, "main\nroot")
	if err := os.MkdirAll(filepath.Join(mainRoot, "docs", "dev"), 0o755); err != nil {
		t.Fatal(err)
	}
	placementGit(t, mainRoot, "init", "-q")
	placementGit(t, mainRoot, "config", "user.email", "t@t")
	placementGit(t, mainRoot, "config", "user.name", "t")
	if err := os.WriteFile(filepath.Join(mainRoot, "docs", "dev", "README.md"), []byte("---\nstate: .state\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	placementGit(t, mainRoot, "add", "docs/dev/README.md")
	placementGit(t, mainRoot, "commit", "-q", "-m", "seed")
	linked := filepath.Join(base, "linked")
	placementGit(t, mainRoot, "worktree", "add", "-q", "-b", "linked", linked)

	workflow := filepath.Join(linked, "docs", "dev")
	got, err := ResolveSplitRootCheckout(workflow, ".state")
	if err != nil {
		t.Fatalf("ResolveSplitRootCheckout: %v", err)
	}
	want := filepath.Join(RealpathOf(mainRoot), "docs", "dev", ".state")
	if RealpathOf(got) != want {
		t.Fatalf("checkout = %q, want exact newline-root placement %q", got, want)
	}
}

func TestResolveRepositoryPlacementRejectsBarePrimary(t *testing.T) {
	base := t.TempDir()
	seed := filepath.Join(base, "seed")
	if err := os.MkdirAll(seed, 0o755); err != nil {
		t.Fatal(err)
	}
	placementGit(t, seed, "init", "-q")
	placementGit(t, seed, "config", "user.email", "t@t")
	placementGit(t, seed, "config", "user.name", "t")
	if err := os.WriteFile(filepath.Join(seed, "file"), []byte("seed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	placementGit(t, seed, "add", "file")
	placementGit(t, seed, "commit", "-q", "-m", "seed")
	bare := filepath.Join(base, "repo.git")
	if out, err := exec.Command("git", "clone", "-q", "--bare", seed, bare).CombinedOutput(); err != nil {
		t.Fatalf("clone bare: %v\n%s", err, out)
	}
	linked := filepath.Join(base, "linked")
	branch := strings.TrimSpace(placementGit(t, seed, "branch", "--show-current"))
	if out, err := exec.Command("git", "--git-dir", bare, "worktree", "add", "-q", linked, branch).CombinedOutput(); err != nil {
		t.Fatalf("bare worktree add: %v\n%s", err, out)
	}

	if _, err := ResolveRepositoryPlacement(linked); err == nil || !strings.Contains(err.Error(), "bare") {
		t.Fatalf("bare-primary repository should fail closed, got %v", err)
	}
}
