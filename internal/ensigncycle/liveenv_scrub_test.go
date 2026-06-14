// ABOUTME: Offline test that the live-cycle child env scrubs the CI repo-naming
// ABOUTME: vars (GITHUB_*/RUNNER_*) that lure the FO off its launch cwd (no model).
package ensigncycle

import (
	"testing"
)

// ciRepoNamingVars are the CI-runner-injected GITHUB_*/RUNNER_* vars the live-test
// child env must NOT carry. On a GitHub Actions runner os.Environ() includes the
// whole family; GITHUB_WORKSPACE (= /home/runner/work/spacedock/spacedock, the REAL
// repo) is the proven lure that made the FO `cd` away from its launch cwd and boot
// the real docs/dev workflow instead of its tmpdir fixture (PR #365 opus
// TestLiveEnsignCycle). In real `spacedock claude` use the cwd IS the project and no
// such CI var exists, so scrubbing the family reproduces a production-clean child
// env. The list is the GITHUB_*/RUNNER_* names a runner sets that name or point at
// the real repo / the runner workspace; ANTHROPIC_API_KEY, CLAUDE_CONFIG_DIR, and
// PATH are NOT in it (they are the credential / archive / launcher the test needs).
var ciRepoNamingVars = []string{
	"GITHUB_WORKSPACE",
	"GITHUB_REPOSITORY",
	"GITHUB_REPOSITORY_OWNER",
	"GITHUB_ACTION_PATH",
	"GITHUB_SERVER_URL",
	"RUNNER_WORKSPACE",
	"RUNNER_TEMP",
}

// TestCleanEnvironScrubsCIRepoNamingVars asserts cleanEnviron drops every
// GITHUB_*/RUNNER_* repo-naming var from the parent env, so the live FO subprocess
// never sees GITHUB_WORKSPACE (or its family) and cannot be lured off its launch
// cwd. It seeds the parent env with the family, builds the child env via
// cleanEnviron, and asserts none survive.
func TestCleanEnvironScrubsCIRepoNamingVars(t *testing.T) {
	for _, key := range ciRepoNamingVars {
		t.Setenv(key, "/home/runner/work/spacedock/spacedock")
	}

	env := cleanEnviron("CLAUDECODE", "HOME", "CLAUDE_CONFIG_DIR")

	for _, key := range ciRepoNamingVars {
		if v, ok := envValue(env, key); ok {
			t.Errorf("cleanEnviron leaked CI repo-naming var %s=%q — it lures the FO off its launch cwd", key, v)
		}
	}
}

// TestIsolatedClaudeEnvScrubsCIRepoNamingVars asserts the same scrub through the
// concrete child env both Claude live lanes build (TestLiveEnsignCycle and
// TestLiveClaudeSharedScenarios). It covers the API-key (CI) auth path — the path
// the leaked CI env actually rides on — and confirms the credential the test needs
// survives the scrub.
func TestIsolatedClaudeEnvScrubsCIRepoNamingVars(t *testing.T) {
	fakeHome := t.TempDir() // no benchmark-token -> API-key (CI) path
	t.Setenv("ANTHROPIC_API_KEY", "sk-ci-api-key")
	t.Setenv("CLAUDECODE", "1")
	for _, key := range ciRepoNamingVars {
		t.Setenv(key, "/home/runner/work/spacedock/spacedock")
	}

	env := isolatedClaudeEnv(t, fakeHome)

	for _, key := range ciRepoNamingVars {
		if v, ok := envValue(env, key); ok {
			t.Errorf("isolatedClaudeEnv leaked CI repo-naming var %s=%q to the child FO", key, v)
		}
	}
	// The credential the CI auth path needs MUST survive the scrub.
	if key, ok := envValue(env, "ANTHROPIC_API_KEY"); !ok || key != "sk-ci-api-key" {
		t.Errorf("ANTHROPIC_API_KEY = %q (present=%v), want it to survive the scrub", key, ok)
	}
}
