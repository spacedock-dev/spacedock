package ensigncycle

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDecideCodexLiveAuth(t *testing.T) {
	t.Run("api_key_runs", func(t *testing.T) {
		d := decideCodexLiveAuth("sk-test", false, "")
		if d.mode != codexAuthAPIKey {
			t.Fatalf("mode = %d, want codexAuthAPIKey", d.mode)
		}
		if d.message != "" {
			t.Fatalf("message = %q, want empty", d.message)
		}
	})

	t.Run("local_subscription_auth_runs_without_api_key", func(t *testing.T) {
		d := decideCodexLiveAuth("", true, "")
		if d.mode != codexAuthLocal {
			t.Fatalf("mode = %d, want codexAuthLocal", d.mode)
		}
		if d.message != "" {
			t.Fatalf("message = %q, want empty", d.message)
		}
	})

	t.Run("ci_required_ignores_local_subscription_auth", func(t *testing.T) {
		d := decideCodexLiveAuth("", true, "1")
		if d.mode != codexAuthFatal {
			t.Fatalf("mode = %d, want codexAuthFatal", d.mode)
		}
		if d.message == "" {
			t.Fatal("fatal decision must carry a clear message")
		}
	})

	t.Run("missing_local_auth_skips", func(t *testing.T) {
		d := decideCodexLiveAuth("", false, "")
		if d.mode != codexAuthSkip {
			t.Fatalf("mode = %d, want codexAuthSkip", d.mode)
		}
		if d.message == "" {
			t.Fatal("skip decision must carry a message")
		}
	})

	t.Run("missing_required_auth_fails", func(t *testing.T) {
		d := decideCodexLiveAuth("", false, "1")
		if d.mode != codexAuthFatal {
			t.Fatalf("mode = %d, want codexAuthFatal", d.mode)
		}
		if d.message == "" {
			t.Fatal("fatal decision must carry a clear message")
		}
	})
}

func TestCodexLocalAuthAvailable(t *testing.T) {
	realHome := t.TempDir()
	if codexLocalAuthAvailable(realHome) {
		t.Fatal("auth must be unavailable before auth.json exists")
	}

	authPath := filepath.Join(realHome, ".codex", "auth.json")
	if err := os.MkdirAll(filepath.Dir(authPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(authPath, []byte("  \n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if codexLocalAuthAvailable(realHome) {
		t.Fatal("blank auth.json must not count as available")
	}

	if err := os.WriteFile(authPath, []byte(`{"auth_mode":"chatgpt"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if !codexLocalAuthAvailable(realHome) {
		t.Fatal("non-empty auth.json must count as available")
	}
}

func TestSeedCodexLocalAuthCopiesOnlyAuthIntoIsolatedHome(t *testing.T) {
	realHome := t.TempDir()
	authPath := filepath.Join(realHome, ".codex", "auth.json")
	if err := os.MkdirAll(filepath.Dir(authPath), 0o755); err != nil {
		t.Fatal(err)
	}
	want := []byte(`{"auth_mode":"chatgpt","tokens":{"refresh_token":"test"}}`)
	if err := os.WriteFile(authPath, want, 0o600); err != nil {
		t.Fatal(err)
	}

	codexHome := t.TempDir()
	if err := seedCodexLocalAuth(codexHome, realHome); err != nil {
		t.Fatalf("seedCodexLocalAuth errored: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(codexHome, "auth.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("copied auth.json mismatch: got %q want %q", got, want)
	}
	if _, err := os.Stat(filepath.Join(codexHome, "plugins")); !os.IsNotExist(err) {
		t.Fatalf("seed must not copy or link real plugin state, stat err=%v", err)
	}
}

func TestCodexLiveEnvOmitsEmptyAPIKey(t *testing.T) {
	env := codexLiveEnv("/tmp/codex-home", "/tmp/home", "", "")
	if _, ok := envValue(env, "OPENAI_API_KEY"); ok {
		t.Fatal("OPENAI_API_KEY must be omitted for local subscription auth")
	}
}

func codexLiveEnv(codexHome, home, pathPrefix, openAIAPIKey string) []string {
	env := cleanEnviron("CODEX_HOME", "HOME", "OPENAI_API_KEY", "CLAUDECODE")
	path := os.Getenv("PATH")
	if pathPrefix != "" {
		path = pathPrefix + string(os.PathListSeparator) + path
	}
	env = append(env,
		"CODEX_HOME="+codexHome,
		"HOME="+home,
		"PATH="+path,
	)
	if openAIAPIKey != "" {
		env = append(env, "OPENAI_API_KEY="+openAIAPIKey)
	}
	return env
}
