package testgit

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// scrubbedIdentityEnv disables git's ambient/global identity fallback so a
// subprocess commit can only succeed if the repo itself carries a persisted
// user.name/user.email. HOME-scrubbing alone is not sufficient: git falls back
// to auto-detecting user.email from username@hostname, which succeeds on a
// developer machine and only fails on a runner with a different hostname
// shape. user.useConfigOnly=true removes that fallback deterministically.
func scrubbedIdentityEnv() []string {
	return append(os.Environ(),
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_SYSTEM=/dev/null",
		"GIT_CONFIG_COUNT=1",
		"GIT_CONFIG_KEY_0=user.useConfigOnly",
		"GIT_CONFIG_VALUE_0=true",
	)
}

func commitWithScrubbedIdentity(dir string) error {
	cmd := exec.Command("git", "commit", "--allow-empty", "-m", "scrubbed-identity probe")
	cmd.Dir = dir
	cmd.Env = scrubbedIdentityEnv()
	return cmd.Run()
}

// TestInitRepoSurvivesScrubbedIdentity is AC-3's positive case: a repo
// scaffolded by InitRepo persists identity in its own config, so a plain
// commit succeeds even with no ambient identity available to the subprocess.
func TestInitRepoSurvivesScrubbedIdentity(t *testing.T) {
	dir := t.TempDir()
	InitRepo(t, dir)

	if err := commitWithScrubbedIdentity(dir); err != nil {
		t.Fatalf("commit in InitRepo'd repo under scrubbed identity: %v", err)
	}
}

// TestPlainGitInitFailsUnderScrubbedIdentity is AC-3's paired negative case.
// Without it, the positive case above would pass on any developer machine via
// hostname auto-detection -- the exact way the original bug escaped local
// testing. A repo made by a bare `git init` (no InitRepo call, no persisted
// identity) must fail the same commit.
func TestPlainGitInitFailsUnderScrubbedIdentity(t *testing.T) {
	dir := t.TempDir()
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}

	if err := commitWithScrubbedIdentity(dir); err == nil {
		t.Fatal("commit in a plain git-init repo (no InitRepo) unexpectedly succeeded under scrubbed identity -- falsifiability check failed")
	}
}

// TestInitRepoPersistsUserEmailInRepoConfig asserts the identity resolves
// from the repo's own local config after a single InitRepo call, not from
// any inherited/global config.
func TestInitRepoPersistsUserEmailInRepoConfig(t *testing.T) {
	dir := t.TempDir()
	InitRepo(t, dir)

	cmd := exec.Command("git", "config", "--local", "user.email")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git config --local user.email: %v", err)
	}
	if got, want := strings.TrimSpace(string(out)), "spacedock@example.invalid"; got != want {
		t.Fatalf("git config --local user.email = %q, want %q", got, want)
	}
}
