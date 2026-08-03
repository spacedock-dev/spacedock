// Package testgit scaffolds throwaway git repositories for tests.
package testgit

import (
	"os/exec"
	"testing"
)

// InitRepo initializes a git repo at dir and persists user.name/user.email in
// the repo's own config, so a later plain `git commit` -- including one run by
// the code under test -- resolves an identity without the host's ambient git
// config. Extra args are appended to `git init` (e.g. "-b", "main").
func InitRepo(t testing.TB, dir string, initArgs ...string) {
	t.Helper()
	runGit(t, dir, append([]string{"init"}, initArgs...)...)
	runGit(t, dir, "config", "user.name", "Spacedock Test")
	runGit(t, dir, "config", "user.email", "spacedock@example.invalid")
}

func runGit(t testing.TB, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
