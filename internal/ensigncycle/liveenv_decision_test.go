// ABOUTME: Offline unit test for the live cycle's auth + HOME-isolation decision
// ABOUTME: tree, mirroring the Python test_isolated_claude_env_* cases (no model).
package ensigncycle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// These mirror ~/git/spacedock/tests/test_test_lib_helpers.py's three
// test_isolated_claude_env_* cases against the Go port. They run under the
// DEFAULT build tags (no //go:build live) so `go test ./...` covers them; they
// spend NO model — only the credential-selection logic is exercised, against a
// fake HOME so the real ~/.claude is never read.

// hasEnvKV reports whether the KEY=VALUE pair (any value) for key is present in
// env, and returns its value.
func envValue(env []string, key string) (string, bool) {
	prefix := key + "="
	for _, kv := range env {
		if strings.HasPrefix(kv, prefix) {
			return strings.TrimPrefix(kv, prefix), true
		}
	}
	return "", false
}

// TestDecideClaudeEnv covers the pure decision tree directly: it returns the
// chosen auth mode (and OAuth token) from explicit inputs, with no skip/abort
// and no temp dir, so all three branches — including the neither-credential
// case the live wrapper turns into t.Skip — are assertable in one offline test.
func TestDecideClaudeEnv(t *testing.T) {
	// (a) operator-local: a non-empty benchmark-token under the fake HOME wins,
	// even when ANTHROPIC_API_KEY is also set — the token is authoritative and
	// the key is dropped (mode carries no key).
	t.Run("token_file_present_picks_oauth_and_drops_key", func(t *testing.T) {
		fakeHome := t.TempDir()
		claudeDir := filepath.Join(fakeHome, ".claude")
		if err := os.MkdirAll(claudeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		writeFile(t, filepath.Join(claudeDir, "benchmark-token"), "sk-oauth-test-token\n")

		d := decideClaudeEnv(fakeHome, "sk-api-should-be-dropped")

		if d.mode != authOAuthToken {
			t.Fatalf("mode = %d, want authOAuthToken", d.mode)
		}
		if d.oauthToken != "sk-oauth-test-token" {
			t.Errorf("oauthToken = %q, want trimmed sk-oauth-test-token", d.oauthToken)
		}
	})

	// (b) CI: no token file but ANTHROPIC_API_KEY present -> passthrough mode.
	t.Run("api_key_only_passes_through", func(t *testing.T) {
		fakeHome := t.TempDir() // no .claude/benchmark-token under it

		d := decideClaudeEnv(fakeHome, "sk-ci-api-key")

		if d.mode != authAPIKey {
			t.Fatalf("mode = %d, want authAPIKey", d.mode)
		}
		if d.oauthToken != "" {
			t.Errorf("oauthToken = %q, want empty for the API-key path", d.oauthToken)
		}
	})

	// (c) neither credential -> authNone (the live wrapper turns this into
	// t.Skip, NOT t.Fatal).
	t.Run("neither_credential_is_none", func(t *testing.T) {
		fakeHome := t.TempDir()

		d := decideClaudeEnv(fakeHome, "")

		if d.mode != authNone {
			t.Fatalf("mode = %d, want authNone", d.mode)
		}
	})

	// An empty benchmark-token file is treated as absent (the original strips
	// then checks for non-empty), so a bare API key still wins path (b).
	t.Run("empty_token_file_falls_through_to_api_key", func(t *testing.T) {
		fakeHome := t.TempDir()
		claudeDir := filepath.Join(fakeHome, ".claude")
		if err := os.MkdirAll(claudeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		writeFile(t, filepath.Join(claudeDir, "benchmark-token"), "   \n")

		d := decideClaudeEnv(fakeHome, "sk-ci-api-key")

		if d.mode != authAPIKey {
			t.Fatalf("mode = %d, want authAPIKey for an empty token file", d.mode)
		}
	})
}

// TestResolveClaudeConfigDir covers the CLAUDE_CONFIG_DIR resolution the child
// env carries: the parent's value passes through when set (the CI path — the
// archivable $RUNNER_TEMP dir the upload step grabs), else a deterministic path
// under the isolated HOME (the local path — a fresh, inspectable config dir).
func TestResolveClaudeConfigDir(t *testing.T) {
	// CI path: the job env sets CLAUDE_CONFIG_DIR to the archivable per-job dir;
	// the child must keep it verbatim so claude writes projects/<slug>/*.jsonl
	// there (not under the reaped isolated HOME) for the upload step to capture.
	t.Run("parent_value_passes_through", func(t *testing.T) {
		isolatedHome := t.TempDir()
		parent := "/runner/_temp/spacedock-claude-config/sonnet"

		got := resolveClaudeConfigDir(parent, isolatedHome)

		if got != parent {
			t.Errorf("resolveClaudeConfigDir = %q, want the parent value %q", got, parent)
		}
	})

	// Local path: no CLAUDE_CONFIG_DIR in the parent env (an operator run with no
	// CI env) -> a deterministic dir under the isolated HOME so local runs still
	// land the jsonl somewhere inspectable.
	t.Run("unset_parent_defaults_under_isolated_home", func(t *testing.T) {
		isolatedHome := t.TempDir()

		got := resolveClaudeConfigDir("", isolatedHome)

		want := filepath.Join(isolatedHome, ".claude")
		if got != want {
			t.Errorf("resolveClaudeConfigDir = %q, want the isolated-home default %q", got, want)
		}
	})
}

// TestIsolatedClaudeEnvHonorsParentConfigDir asserts the concrete child env keeps
// the parent's CLAUDE_CONFIG_DIR (the CI archivable path) rather than dropping or
// overriding it — the contract the upload step depends on.
func TestIsolatedClaudeEnvHonorsParentConfigDir(t *testing.T) {
	fakeHome := t.TempDir() // no benchmark-token -> API-key path
	t.Setenv("ANTHROPIC_API_KEY", "sk-ci-api-key")
	t.Setenv("CLAUDECODE", "1")
	parentConfig := "/runner/_temp/spacedock-claude-config/sonnet"
	t.Setenv("CLAUDE_CONFIG_DIR", parentConfig)

	env := isolatedClaudeEnv(t, fakeHome)

	cfg, ok := envValue(env, "CLAUDE_CONFIG_DIR")
	if !ok || cfg != parentConfig {
		t.Errorf("CLAUDE_CONFIG_DIR = %q (present=%v), want the parent CI path %q", cfg, ok, parentConfig)
	}
}

// TestIsolatedClaudeEnvDefaultsConfigDirUnderHome asserts that when the parent
// has no CLAUDE_CONFIG_DIR (local operator run), the child env gets a
// deterministic config dir under the fresh isolated HOME.
func TestIsolatedClaudeEnvDefaultsConfigDirUnderHome(t *testing.T) {
	fakeHome := t.TempDir() // no benchmark-token -> API-key path
	t.Setenv("ANTHROPIC_API_KEY", "sk-ci-api-key")
	t.Setenv("CLAUDECODE", "1")
	os.Unsetenv("CLAUDE_CONFIG_DIR")

	env := isolatedClaudeEnv(t, fakeHome)

	home, ok := envValue(env, "HOME")
	if !ok {
		t.Fatal("HOME must be set to the fresh isolated dir")
	}
	cfg, ok := envValue(env, "CLAUDE_CONFIG_DIR")
	if !ok {
		t.Fatal("CLAUDE_CONFIG_DIR must default to a dir under the isolated HOME")
	}
	want := filepath.Join(home, ".claude")
	if cfg != want {
		t.Errorf("CLAUDE_CONFIG_DIR = %q, want the isolated-home default %q", cfg, want)
	}
}

// TestIsolatedClaudeEnvOAuthPath asserts the concrete child env the OAuth path
// produces: CLAUDE_CODE_OAUTH_TOKEN set, ANTHROPIC_API_KEY ABSENT (dropped),
// CLAUDECODE dropped, and HOME pointing at a fresh dir that is NOT the real one.
// Mirrors test_isolated_claude_env_injects_oauth_token_when_token_file_present.
func TestIsolatedClaudeEnvOAuthPath(t *testing.T) {
	fakeHome := t.TempDir()
	claudeDir := filepath.Join(fakeHome, ".claude")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(claudeDir, "benchmark-token"), "sk-oauth-test-token\n")

	// Set both credential keys + CLAUDECODE in the parent so we can prove the
	// helper drops the API key + CLAUDECODE and overrides HOME.
	t.Setenv("ANTHROPIC_API_KEY", "sk-api-should-be-dropped")
	t.Setenv("CLAUDECODE", "1")

	env := isolatedClaudeEnv(t, fakeHome)

	if tok, ok := envValue(env, "CLAUDE_CODE_OAUTH_TOKEN"); !ok || tok != "sk-oauth-test-token" {
		t.Errorf("CLAUDE_CODE_OAUTH_TOKEN = %q (present=%v), want sk-oauth-test-token", tok, ok)
	}
	if _, ok := envValue(env, "ANTHROPIC_API_KEY"); ok {
		t.Error("ANTHROPIC_API_KEY must be dropped on the OAuth path")
	}
	if _, ok := envValue(env, "CLAUDECODE"); ok {
		t.Error("CLAUDECODE must be dropped so the child takes the front-door path")
	}
	home, ok := envValue(env, "HOME")
	if !ok {
		t.Fatal("HOME must be set to the fresh isolated dir")
	}
	if home == fakeHome {
		t.Errorf("HOME = %q must be a fresh dir, not the real home %q", home, fakeHome)
	}
	if info, err := os.Stat(home); err != nil || !info.IsDir() {
		t.Errorf("HOME %q must be an existing dir: err=%v", home, err)
	}
}

// TestIsolatedClaudeEnvAPIKeyPath asserts the CI path: ANTHROPIC_API_KEY passed
// through, no OAuth token, CLAUDECODE dropped, HOME a fresh dir. Mirrors
// test_isolated_claude_env_preserves_api_key_when_no_token_file (fakeHome has no
// benchmark-token).
func TestIsolatedClaudeEnvAPIKeyPath(t *testing.T) {
	fakeHome := t.TempDir() // no .claude/benchmark-token
	t.Setenv("ANTHROPIC_API_KEY", "sk-ci-api-key")
	t.Setenv("CLAUDECODE", "1")

	env := isolatedClaudeEnv(t, fakeHome)

	if key, ok := envValue(env, "ANTHROPIC_API_KEY"); !ok || key != "sk-ci-api-key" {
		t.Errorf("ANTHROPIC_API_KEY = %q (present=%v), want sk-ci-api-key", key, ok)
	}
	if _, ok := envValue(env, "CLAUDE_CODE_OAUTH_TOKEN"); ok {
		t.Error("CLAUDE_CODE_OAUTH_TOKEN must be absent on the API-key path")
	}
	if _, ok := envValue(env, "CLAUDECODE"); ok {
		t.Error("CLAUDECODE must be dropped so the child takes the front-door path")
	}
	home, ok := envValue(env, "HOME")
	if !ok {
		t.Fatal("HOME must be set to the fresh isolated dir")
	}
	if home == fakeHome {
		t.Errorf("HOME = %q must be a fresh dir, not the real home %q", home, fakeHome)
	}
}
