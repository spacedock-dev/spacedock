package ensigncycle

import (
	"os"
	"path/filepath"
	"strings"
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

func TestSeedCodexLiveConfigWritesExactMultiAgentFragment(t *testing.T) {
	codexHome := t.TempDir()

	if err := seedCodexLiveConfig(codexHome); err != nil {
		t.Fatalf("seedCodexLiveConfig errored: %v", err)
	}

	want := "[features.multi_agent_v2]\n" +
		"max_concurrent_threads_per_session = 16\n" +
		"tool_namespace = \"agents\"\n" +
		"hide_spawn_agent_metadata = false\n"
	got, err := os.ReadFile(filepath.Join(codexHome, "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("config.toml = %q, want exact multi_agent_v2 fragment %q", got, want)
	}

	entries, err := os.ReadDir(codexHome)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "config.toml" {
		t.Fatalf("isolated home entries = %#v, want only config.toml", entries)
	}
}

func TestSeedCodexLiveHomeCopiesOnlyAuthAndMinimalConfig(t *testing.T) {
	realHome := t.TempDir()
	operatorCodexDir := filepath.Join(realHome, ".codex")
	if err := os.MkdirAll(filepath.Join(operatorCodexDir, "plugins", "state"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(operatorCodexDir, "credentials.json"), []byte("operator-credential"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(operatorCodexDir, "config.toml"), []byte("operator-only = true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	auth := []byte(`{"auth_mode":"chatgpt","tokens":{"refresh_token":"test"}}`)
	if err := os.WriteFile(filepath.Join(operatorCodexDir, "auth.json"), auth, 0o600); err != nil {
		t.Fatal(err)
	}

	codexHome := t.TempDir()
	if err := seedCodexLiveConfig(codexHome); err != nil {
		t.Fatalf("seedCodexLiveConfig errored: %v", err)
	}
	if err := seedCodexLocalAuth(codexHome, realHome); err != nil {
		t.Fatalf("seedCodexLocalAuth errored: %v", err)
	}

	entries, err := os.ReadDir(codexHome)
	if err != nil {
		t.Fatal(err)
	}
	wantEntries := []string{"auth.json", "config.toml"}
	if len(entries) != len(wantEntries) {
		t.Fatalf("isolated home entries = %#v, want exactly %v", entries, wantEntries)
	}
	for i, entry := range entries {
		if entry.Name() != wantEntries[i] {
			t.Fatalf("isolated home entry %d = %q, want %q", i, entry.Name(), wantEntries[i])
		}
	}

	gotAuth, err := os.ReadFile(filepath.Join(codexHome, "auth.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(gotAuth) != string(auth) {
		t.Fatalf("auth.json = %q, want copied auth", gotAuth)
	}
	gotConfig, err := os.ReadFile(filepath.Join(codexHome, "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	wantConfig := "[features.multi_agent_v2]\n" +
		"max_concurrent_threads_per_session = 16\n" +
		"tool_namespace = \"agents\"\n" +
		"hide_spawn_agent_metadata = false\n"
	if string(gotConfig) != wantConfig {
		t.Fatalf("config.toml = %q, want minimal live fragment %q", gotConfig, wantConfig)
	}
	for _, name := range []string{"credentials.json", "plugins"} {
		if _, err := os.Stat(filepath.Join(codexHome, name)); !os.IsNotExist(err) {
			t.Fatalf("isolated home must not copy %s, stat err=%v", name, err)
		}
	}
}

func TestCodexLiveEnvOmitsEmptyAPIKey(t *testing.T) {
	env := codexLiveEnv("/tmp/codex-home", "/tmp/home", "", "")
	if _, ok := envValue(env, "OPENAI_API_KEY"); ok {
		t.Fatal("OPENAI_API_KEY must be omitted for local subscription auth")
	}
}

func TestCodexLiveEnvDropsForeignRuntimeMarkers(t *testing.T) {
	for key, value := range map[string]string{"CODEX_THREAD_ID": "codex", "CLAUDECODE": "claude", "PI_CODING_AGENT": "pi",
		"PI_CODING_AGENT_DIR": "/parent/pi", "CODEX_HOME": "/parent/codex", "HOME": "/parent/home", "OPENAI_API_KEY": "parent-key", "PATH": "/parent/bin"} {
		t.Setenv(key, value)
	}

	env := codexLiveEnv("/target/codex", "/target/home", "/spacedock/bin", "target-key")

	want := map[string]string{"CLAUDECODE": "", "PI_CODING_AGENT": "", "PI_CODING_AGENT_DIR": "",
		"CODEX_THREAD_ID": "codex", "CODEX_HOME": "/target/codex", "HOME": "/target/home",
		"OPENAI_API_KEY": "target-key", "PATH": "/spacedock/bin" + string(os.PathListSeparator) + "/parent/bin"}
	for key, value := range want {
		assertEnvValue(t, env, key, value)
	}
}

func TestCodexLiveHomeParentUsesUserCacheOutsideSystemTemp(t *testing.T) {
	cacheDir := filepath.Join(string(os.PathSeparator), "home", "runner", ".cache")
	parent, err := codexLiveIsolatedHomeParent(cacheDir)
	if err != nil {
		t.Fatalf("codexLiveIsolatedHomeParent errored: %v", err)
	}
	want := filepath.Join(cacheDir, "spacedock-live-codex")
	if parent != want {
		t.Fatalf("parent = %q, want %q", parent, want)
	}
	if tmp := filepath.Clean(os.TempDir()); parent == tmp || strings.HasPrefix(parent, tmp+string(os.PathSeparator)) {
		t.Fatalf("isolated CODEX_HOME parent %q is under system temp %q; Codex refuses helper aliases there", parent, tmp)
	}
}

func TestCodexLiveHomeParentCandidatesKeepFallbackOutsidePluginCheckoutWhenArtifactsAreTemp(t *testing.T) {
	cacheDir := filepath.Join(string(os.PathSeparator), "blocked", "cache")
	repoRoot := filepath.Join(string(os.PathSeparator), "home", "runner", "work", "spacedock", "spacedock")
	artifactRoot := filepath.Join(os.TempDir(), "spacedock-live-artifacts")

	got := codexLiveIsolatedHomeParentCandidates(cacheDir, repoRoot, artifactRoot)
	wantCache := filepath.Join(cacheDir, "spacedock-live-codex")
	wantRepo := filepath.Join(filepath.Dir(repoRoot), ".spacedock-live-codex", filepath.Base(repoRoot))
	if len(got) != 2 || got[0] != wantCache || got[1] != wantRepo {
		t.Fatalf("candidates = %#v, want cache then repo fallback %#v", got, []string{wantCache, wantRepo})
	}
	for _, parent := range got {
		if strings.HasPrefix(parent, artifactRoot) {
			t.Fatalf("temp artifact root %q must not be used as isolated CODEX_HOME parent; candidates=%#v", artifactRoot, got)
		}
		if parent == repoRoot || strings.HasPrefix(parent, repoRoot+string(os.PathSeparator)) {
			t.Fatalf("isolated CODEX_HOME parent %q is inside plugin checkout %q; codex plugin add copies the checkout into CODEX_HOME and can recurse until path limits fail", parent, repoRoot)
		}
	}
}

func TestCodexLiveHomeParentCandidatesRejectArtifactRootInsidePluginCheckout(t *testing.T) {
	cacheDir := filepath.Join(string(os.PathSeparator), "home", "runner", ".cache")
	repoRoot := filepath.Join(string(os.PathSeparator), "home", "runner", "work", "spacedock", "spacedock")
	artifactRoot := filepath.Join(repoRoot, "live-artifacts", "codex", "codex-shared-scenarios")

	got := codexLiveIsolatedHomeParentCandidates(cacheDir, repoRoot, artifactRoot)
	wantCache := filepath.Join(cacheDir, "spacedock-live-codex")
	wantRepo := filepath.Join(filepath.Dir(repoRoot), ".spacedock-live-codex", filepath.Base(repoRoot))
	if len(got) != 2 || got[0] != wantCache || got[1] != wantRepo {
		t.Fatalf("candidates = %#v, want cache then repo-adjacent fallback %#v", got, []string{wantCache, wantRepo})
	}
	for _, parent := range got {
		if parent == repoRoot || strings.HasPrefix(parent, repoRoot+string(os.PathSeparator)) {
			t.Fatalf("isolated CODEX_HOME parent %q is inside plugin checkout %q; artifact-backed homes there make codex plugin add copy its own cache", parent, repoRoot)
		}
	}
}

func codexLiveEnv(codexHome, home, pathPrefix, openAIAPIKey string) []string {
	env := cleanEnviron("CODEX_HOME", "HOME", "OPENAI_API_KEY", "PATH",
		"CLAUDECODE", "PI_CODING_AGENT", "PI_CODING_AGENT_DIR")
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
