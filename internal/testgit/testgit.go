// Package testgit scaffolds throwaway git repositories for tests.
package testgit

import (
	"os/exec"
	"testing"
)

// InitRepo initializes a git repo at dir and persists identity plus automatic-
// maintenance disablement in the repo's own config. Throwaway test repositories
// must not leave detached maintenance writing into a TempDir during cleanup.
// maintenance.auto is the modern switch; gc.auto=0 covers older Git versions.
// Extra args are appended to `git init` (e.g. "-b", "main").
func InitRepo(t testing.TB, dir string, initArgs ...string) {
	t.Helper()
	runGit(t, dir, append([]string{"init"}, initArgs...)...)
	runGit(t, dir, "config", "--local", "user.name", "Spacedock Test")
	runGit(t, dir, "config", "--local", "user.email", "spacedock@example.invalid")
	runGit(t, dir, "config", "--local", "maintenance.auto", "false")
	runGit(t, dir, "config", "--local", "gc.auto", "0")
}

func runGit(t testing.TB, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
