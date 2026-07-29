package gitsource

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInspectAndResolveUseExactLocalGitObjectsAcrossMovedRoots(t *testing.T) {
	mainRoot := initRepository(t, "main", map[string]string{
		"docs/review file.md": "main review\n",
	})
	stateRoot := initRepository(t, "state", map[string]string{
		"task/index.md": "state entity\n",
	})
	mainCommit := gitOutput(t, mainRoot, "rev-parse", "HEAD")
	stateCommit := gitOutput(t, stateRoot, "rev-parse", "HEAD")

	roots := Roots{Main: mainRoot, State: stateRoot}
	mainSource, err := Inspect(roots, filepath.Join(mainRoot, "docs", "review file.md"))
	if err != nil {
		t.Fatal(err)
	}
	stateSource, err := Inspect(roots, filepath.Join(stateRoot, "task", "index.md"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "git-root://main/" + mainCommit + "/docs/review%20file.md"; mainSource.URI != want {
		t.Fatalf("main URI=%q want %q", mainSource.URI, want)
	}
	if want := "git-root://state/" + stateCommit + "/task/index.md"; stateSource.URI != want {
		t.Fatalf("state URI=%q want %q", stateSource.URI, want)
	}
	if mainSource.Rev != RawDigest([]byte("main review\n")) || stateSource.Rev != RawDigest([]byte("state entity\n")) {
		t.Fatalf("raw revisions not derived from selected bytes: main=%s state=%s", mainSource.Rev, stateSource.Rev)
	}

	moved := t.TempDir()
	movedMain := filepath.Join(moved, "unrelated-main")
	movedState := filepath.Join(moved, "unrelated-state")
	if err := os.Rename(mainRoot, movedMain); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(stateRoot, movedState); err != nil {
		t.Fatal(err)
	}
	movedRoots := Roots{Main: movedMain, State: movedState}
	for _, tc := range []struct {
		source Source
		want   string
	}{
		{mainSource, "main review\n"},
		{stateSource, "state entity\n"},
	} {
		got, err := Resolve(movedRoots, tc.source.URI, tc.source.Rev)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != tc.want {
			t.Fatalf("resolved bytes=%q want %q", got, tc.want)
		}
	}
}

func TestInspectAcceptsDetachedCommitAndRejectsMutableOrForeignSelections(t *testing.T) {
	mainRoot := initRepository(t, "main", map[string]string{"review.md": "review\n"})
	stateRoot := initRepository(t, "state", map[string]string{"task.md": "task\n"})
	roots := Roots{Main: mainRoot, State: stateRoot}
	linkedRoot := filepath.Join(t.TempDir(), "linked-main")
	gitRun(t, mainRoot, "worktree", "add", "-q", "--detach", linkedRoot, "HEAD")
	source, err := Inspect(roots, filepath.Join(linkedRoot, "review.md"))
	if err != nil {
		t.Fatalf("clean detached linked-worktree source rejected: %v", err)
	}
	if !strings.HasPrefix(source.URI, "git-root://main/") {
		t.Fatalf("linked worktree classified outside main: %q", source.URI)
	}

	if err := os.WriteFile(filepath.Join(linkedRoot, "review.md"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Inspect(roots, filepath.Join(linkedRoot, "review.md")); err == nil || !strings.Contains(err.Error(), "commit the exact source") {
		t.Fatalf("dirty source error=%v", err)
	}
	if err := os.WriteFile(filepath.Join(linkedRoot, "untracked.md"), []byte("new\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Inspect(roots, filepath.Join(linkedRoot, "untracked.md")); err == nil || !strings.Contains(err.Error(), "commit the exact source") {
		t.Fatalf("untracked source error=%v", err)
	}
	if err := os.Symlink(filepath.Join(stateRoot, "task.md"), filepath.Join(stateRoot, "link.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := Inspect(roots, filepath.Join(stateRoot, "link.md")); err == nil || !strings.Contains(err.Error(), "non-symlink regular file") {
		t.Fatalf("symlink source error=%v", err)
	}

	foreign := initRepository(t, "foreign", map[string]string{"other.md": "other\n"})
	if _, err := Inspect(roots, filepath.Join(foreign, "other.md")); err == nil || !strings.Contains(err.Error(), "workflow Git root") {
		t.Fatalf("foreign source error=%v", err)
	}
}

func TestInspectDistinguishesSplitStateWorktreeFromMainLinkedWorktree(t *testing.T) {
	mainRoot := initRepository(t, "main", map[string]string{"review.md": "review\n"})
	stateRoot := filepath.Join(t.TempDir(), "state-worktree")
	implementationRoot := filepath.Join(t.TempDir(), "implementation-worktree")
	gitRun(t, mainRoot, "worktree", "add", "-q", "--detach", stateRoot, "HEAD")
	gitRun(t, mainRoot, "worktree", "add", "-q", "--detach", implementationRoot, "HEAD")
	roots := Roots{Main: mainRoot, State: stateRoot}

	stateSource, err := Inspect(roots, filepath.Join(stateRoot, "review.md"))
	if err != nil {
		t.Fatal(err)
	}
	mainSource, err := Inspect(roots, filepath.Join(implementationRoot, "review.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(stateSource.URI, "git-root://state/") {
		t.Fatalf("split state worktree classified outside state: %q", stateSource.URI)
	}
	if !strings.HasPrefix(mainSource.URI, "git-root://main/") {
		t.Fatalf("linked implementation worktree classified outside main: %q", mainSource.URI)
	}
}

func TestResolveRejectsNoncanonicalOrUnverifiableGitRootCoordinates(t *testing.T) {
	mainRoot := initRepository(t, "main", map[string]string{"review.md": "review\n"})
	roots := Roots{Main: mainRoot}
	source, err := Inspect(roots, filepath.Join(mainRoot, "review.md"))
	if err != nil {
		t.Fatal(err)
	}
	commit := gitOutput(t, mainRoot, "rev-parse", "HEAD")
	bad := []string{
		"git-root://unknown/" + commit + "/review.md",
		"git-root://main/" + commit[:12] + "/review.md",
		"git-root://main/" + commit + "/review%2emd",
		"git-root://main/" + commit + "/../review.md",
		"git-root://main/" + commit + "/review.md?ref=main",
		"git-root://user@main/" + commit + "/review.md",
	}
	for _, uri := range bad {
		t.Run(uri, func(t *testing.T) {
			if _, err := Resolve(roots, uri, source.Rev); err == nil {
				t.Fatalf("Resolve(%q) succeeded", uri)
			}
		})
	}
	if _, err := Resolve(roots, source.URI, "sha256:"+strings.Repeat("0", 64)); err == nil || !strings.Contains(err.Error(), "raw SHA-256") {
		t.Fatalf("raw digest mismatch error=%v", err)
	}

	// A commit may be packed by a user's Git configuration; absence is covered by
	// the syntactically valid, nonexistent full object id below in either case.
	_ = os.Remove(filepath.Join(mainRoot, ".git", "objects", commit[:2], commit[2:]))
	missing := "git-root://main/" + strings.Repeat("f", len(commit)) + "/review.md"
	if _, err := Resolve(roots, missing, source.Rev); err == nil {
		t.Fatal("missing local object resolved")
	}
}

func TestInspectTreatsRepositoryPathAsLiteralGitPathspec(t *testing.T) {
	mainRoot := initRepository(t, "main", map[string]string{
		":(glob)review*.md": "literal review\n",
		"review-leak.md":    "must not match\n",
	})
	source, err := Inspect(Roots{Main: mainRoot}, filepath.Join(mainRoot, ":(glob)review*.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(source.URI, ":%28glob%29review%2A.md") {
		t.Fatalf("literal metacharacter path was not canonically escaped: %q", source.URI)
	}
	body, err := Resolve(Roots{Main: mainRoot}, source.URI, source.Rev)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "literal review\n" {
		t.Fatalf("resolved pathspec neighbor instead of literal file: %q", body)
	}
}

func TestSameLogicalRevisionIgnoresUnrelatedCommitButNotPathOrBytes(t *testing.T) {
	mainRoot := initRepository(t, "main", map[string]string{
		"review.md":     "review\n",
		"same-bytes.md": "review\n",
		"other.md":      "other\n",
	})
	roots := Roots{Main: mainRoot}
	before, err := Inspect(roots, filepath.Join(mainRoot, "review.md"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(mainRoot, "other.md"), []byte("changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, mainRoot, "add", "other.md")
	gitRun(t, mainRoot, "commit", "-q", "-m", "unrelated")
	after, err := Inspect(roots, filepath.Join(mainRoot, "review.md"))
	if err != nil {
		t.Fatal(err)
	}
	same, err := SameLogicalRevision(before, after)
	if err != nil {
		t.Fatal(err)
	}
	if !same {
		t.Fatalf("same logical revision rejected:\nbefore=%+v\nafter=%+v", before, after)
	}

	sameBytesAtOtherPath, err := Inspect(roots, filepath.Join(mainRoot, "same-bytes.md"))
	if err != nil {
		t.Fatal(err)
	}
	if same, err := SameLogicalRevision(before, sameBytesAtOtherPath); err != nil || same {
		t.Fatalf("different logical path with identical raw bytes accepted: same=%t err=%v", same, err)
	}

	other, err := Inspect(roots, filepath.Join(mainRoot, "other.md"))
	if err != nil {
		t.Fatal(err)
	}
	if same, err := SameLogicalRevision(before, other); err != nil || same {
		t.Fatalf("different logical path accepted: same=%t err=%v", same, err)
	}
	changed := after
	changed.Rev = RawDigest([]byte("different\n"))
	if same, err := SameLogicalRevision(before, changed); err != nil || same {
		t.Fatalf("different raw revision accepted: same=%t err=%v", same, err)
	}
}

func initRepository(t *testing.T, name string, files map[string]string) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), name)
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "init", "-q")
	gitRun(t, root, "config", "user.name", "Spacedock Test")
	gitRun(t, root, "config", "user.email", "spacedock@example.invalid")
	for name, body := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitRun(t, root, "add", ".")
	gitRun(t, root, "commit", "-q", "-m", "fixture")
	return root
}

func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func gitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(bytes.TrimSpace(out)))
}
