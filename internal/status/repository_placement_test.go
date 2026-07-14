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

func TestTrimGitLineTerminatorPreservesWhitespacePathBytes(t *testing.T) {
	for _, tc := range []struct {
		in, want string
	}{
		{" /tmp/leading-and-trailing \n", " /tmp/leading-and-trailing "},
		{"\t/tmp/tabbed\t\n", "\t/tmp/tabbed\t"},
		{"/tmp/path-ending-in-newline\n\n", "/tmp/path-ending-in-newline\n"},
		{"", ""},
	} {
		if got := TrimGitLineTerminator(tc.in); got != tc.want {
			t.Fatalf("TrimGitLineTerminator(%q)=%q want %q", tc.in, got, tc.want)
		}
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

func TestResolveRepositoryPlacementFromExternalSymlinkIntoLinkedWorktree(t *testing.T) {
	base := t.TempDir()
	mainRoot := filepath.Join(base, "main")
	workflowRel := filepath.Join("docs", "dev")
	if err := os.MkdirAll(filepath.Join(mainRoot, workflowRel), 0o755); err != nil {
		t.Fatal(err)
	}
	placementGit(t, mainRoot, "init", "-q")
	placementGit(t, mainRoot, "config", "user.email", "t@t")
	placementGit(t, mainRoot, "config", "user.name", "t")
	if err := os.WriteFile(filepath.Join(mainRoot, workflowRel, "README.md"), []byte("---\nstate: .state\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	placementGit(t, mainRoot, "add", "docs/dev/README.md")
	placementGit(t, mainRoot, "commit", "-q", "-m", "seed")
	linked := filepath.Join(base, "linked")
	placementGit(t, mainRoot, "worktree", "add", "-q", "-b", "linked-symlink", linked)
	external := filepath.Join(t.TempDir(), "external-workflow")
	if err := os.Symlink(filepath.Join(linked, workflowRel), external); err != nil {
		t.Fatal(err)
	}

	placement, err := ResolveRepositoryPlacement(external)
	if err != nil {
		t.Fatal(err)
	}
	if !placement.InGit || !placement.Linked {
		t.Fatalf("external symlink placement=%+v, want linked Git worktree", placement)
	}
	got, err := ResolveSplitRootCheckout(external, ".state")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(mainRoot, workflowRel, ".state")
	if RealpathOf(got) != RealpathOf(want) {
		t.Fatalf("external symlink checkout=%q want main-root checkout %q", got, want)
	}
}

func TestResolveSplitRootCheckoutRejectsCanonicalEscapes(t *testing.T) {
	base := t.TempDir()
	mainRoot := filepath.Join(base, "main")
	if err := os.MkdirAll(mainRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	placementGit(t, mainRoot, "init", "-q")
	placementGit(t, mainRoot, "config", "user.email", "t@t")
	placementGit(t, mainRoot, "config", "user.name", "t")

	outsidePrefix := filepath.Join(base, "outside-prefix")
	if err := os.MkdirAll(outsidePrefix, 0o755); err != nil {
		t.Fatal(err)
	}
	workflowRel := "workflow"
	if err := os.Symlink(outsidePrefix, filepath.Join(mainRoot, workflowRel)); err != nil {
		t.Fatal(err)
	}
	safeRel := "safe-workflow"
	if err := os.MkdirAll(filepath.Join(mainRoot, safeRel), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(mainRoot, safeRel, "README.md"), []byte("---\nstate: .state\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	placementGit(t, mainRoot, "add", workflowRel, filepath.Join(safeRel, "README.md"))
	placementGit(t, mainRoot, "commit", "-q", "-m", "seed main prefix symlink")

	linked := filepath.Join(base, "linked")
	placementGit(t, mainRoot, "worktree", "add", "-q", "-b", "linked-prefix-directory", linked)
	if err := os.Remove(filepath.Join(linked, workflowRel)); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(linked, workflowRel), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(linked, workflowRel, "README.md"), []byte("---\nstate: .state\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	placementGit(t, linked, "add", "-A")
	placementGit(t, linked, "commit", "-q", "-m", "replace prefix symlink with directory")

	t.Run("main prefix symlink", func(t *testing.T) {
		if got, err := ResolveSplitRootCheckout(filepath.Join(linked, workflowRel), ".state"); err == nil || !strings.Contains(err.Error(), "canonical main worktree") {
			t.Fatalf("symlink-escaped main prefix resolved to %q with err=%v", got, err)
		}
	})

	t.Run("valid linked prefix", func(t *testing.T) {
		got, err := ResolveSplitRootCheckout(filepath.Join(linked, safeRel), ".state")
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(mainRoot, safeRel, ".state")
		if RealpathOf(got) != RealpathOf(want) {
			t.Fatalf("valid linked checkout=%q want=%q", got, want)
		}
	})

	t.Run("state child symlink", func(t *testing.T) {
		outsideState := filepath.Join(base, "outside-state")
		if err := os.MkdirAll(outsideState, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outsideState, filepath.Join(mainRoot, safeRel, ".state")); err != nil {
			t.Fatal(err)
		}
		if got, err := ResolveSplitRootCheckout(filepath.Join(linked, safeRel), ".state"); err == nil || !strings.Contains(err.Error(), "canonical main worktree") {
			t.Fatalf("symlink-escaped state child resolved to %q with err=%v", got, err)
		}
	})

	t.Run("state relative traversal", func(t *testing.T) {
		escapingState := filepath.Join("..", "..", "outside-state")
		if got, err := ResolveSplitRootCheckout(filepath.Join(linked, safeRel), escapingState); err == nil || !strings.Contains(err.Error(), "canonical main worktree") {
			t.Fatalf("parent-escaped state path resolved to %q with err=%v", got, err)
		}
	})
}

func TestResolveSplitRootCheckoutPreservesWhitespacePathBytes(t *testing.T) {
	base := t.TempDir()
	mainRoot := filepath.Join(base, " main root ")
	workflowRel := filepath.Join(" docs ", " dev ")
	if err := os.MkdirAll(filepath.Join(mainRoot, workflowRel), 0o755); err != nil {
		t.Fatal(err)
	}
	placementGit(t, mainRoot, "init", "-q")
	placementGit(t, mainRoot, "config", "user.email", "t@t")
	placementGit(t, mainRoot, "config", "user.name", "t")
	if err := os.WriteFile(filepath.Join(mainRoot, workflowRel, "README.md"), []byte("---\nstate: .state\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	placementGit(t, mainRoot, "add", ".")
	placementGit(t, mainRoot, "commit", "-q", "-m", "seed")
	linked := filepath.Join(base, " linked root ")
	placementGit(t, mainRoot, "worktree", "add", "-q", "-b", "linked-whitespace", linked)

	workflow := filepath.Join(linked, workflowRel)
	placement, err := ResolveRepositoryPlacement(workflow)
	if err != nil {
		t.Fatal(err)
	}
	wantPrefix := filepath.ToSlash(workflowRel) + "/"
	if placement.Prefix != wantPrefix {
		t.Fatalf("prefix bytes=%q want %q", placement.Prefix, wantPrefix)
	}
	got, err := ResolveSplitRootCheckout(workflow, ".state")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(mainRoot, workflowRel, ".state")
	if RealpathOf(got) != RealpathOf(want) {
		t.Fatalf("whitespace checkout=%q want exact %q", got, want)
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
