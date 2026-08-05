package testgit

import (
	"os"
	"os/exec"
	"path/filepath"
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

// TestInitRepoDisablesAutomaticMaintenanceLocally guards throwaway repositories
// against Git commands spawning detached maintenance while a test is tearing its
// TempDir down. maintenance.auto is the modern causal switch; gc.auto=0 is the
// compatible fallback for older Git versions that predate `git maintenance`.
// Both keys are ordinary repo-local config, so old Git safely stores the newer
// key even when it does not act on it.
func TestInitRepoDisablesAutomaticMaintenanceLocally(t *testing.T) {
	dir := t.TempDir()
	InitRepo(t, dir)

	for key, want := range map[string]string{
		"maintenance.auto": "false",
		"gc.auto":          "0",
	} {
		cmd := exec.Command("git", "config", "--local", "--get", key)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git config --local --get %s: %v\n%s", key, err, out)
		}
		if got := strings.TrimSpace(string(out)); got != want {
			t.Errorf("git config --local %s = %q, want %q", key, got, want)
		}
	}

	plain := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = plain
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("plain git init: %v\n%s", err, out)
	}
	for _, key := range []string{"maintenance.auto", "gc.auto"} {
		cmd := exec.Command("git", "config", "--local", "--get", key)
		cmd.Dir = plain
		if out, err := cmd.CombinedOutput(); err == nil {
			t.Errorf("plain git init unexpectedly stored %s=%q; removal control is not discriminating", key, strings.TrimSpace(string(out)))
		}
	}
}

// TestInitRepoSuppressesTraceableAutomaticMaintenance proves the config is the
// causal seam when this Git version exposes automatic work in TRACE2. The
// positive control reenables maintenance but keeps it synchronous, so the test
// never creates the detached cleanup race it guards against. Older Git versions
// that store the keys but do not expose either child command skip this behavioral
// supplement; the version-neutral local-config test above remains mandatory.
func TestInitRepoSuppressesTraceableAutomaticMaintenance(t *testing.T) {
	enabled := t.TempDir()
	InitRepo(t, enabled)
	runGit(t, enabled, "config", "--local", "maintenance.auto", "true")
	runGit(t, enabled, "config", "--local", "maintenance.autoDetach", "false")
	runGit(t, enabled, "config", "--local", "maintenance.strategy", "gc")
	runGit(t, enabled, "config", "--local", "gc.auto", "1")
	enabledTrace := commitTrace2(t, enabled, "enabled")
	if !traceHasAutomaticMaintenance(enabledTrace) {
		t.Skip("installed Git does not expose automatic maintenance or gc as a TRACE2 child")
	}

	disabled := t.TempDir()
	InitRepo(t, disabled)
	disabledTrace := commitTrace2(t, disabled, "disabled")
	if traceHasAutomaticMaintenance(disabledTrace) {
		t.Fatalf("InitRepo spawned automatic maintenance despite repo-local disablement:\n%s", disabledTrace)
	}
}

func commitTrace2(t *testing.T, dir, label string) string {
	t.Helper()
	tracePath := filepath.Join(t.TempDir(), label+"-trace.json")
	cmd := exec.Command("git", "commit", "--allow-empty", "-q", "-m", label)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TRACE2_EVENT="+tracePath)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit with TRACE2: %v\n%s", err, out)
	}
	trace, err := os.ReadFile(tracePath)
	if err != nil {
		t.Fatalf("read TRACE2 event file: %v", err)
	}
	return string(trace)
}

func traceHasAutomaticMaintenance(trace string) bool {
	return strings.Contains(trace, `"argv":["git","maintenance","run","--auto`) ||
		strings.Contains(trace, `"argv":["git","gc","--auto`)
}
