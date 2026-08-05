// Package testgit scaffolds throwaway git repositories for tests.
package testgit

import (
	"os/exec"
	"testing"
)

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
