package ensigncycle

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

func TestDecideCodexLiveAuthForCI(t *testing.T) {
	for _, tc := range []struct {
		name, oauth, key string
		want             codexAuthMode
	}{
		{"oauth_preferred", `{"auth_mode":"chatgpt"}`, "sk-key", codexAuthOAuth},
		{"api_key_fallback", "", "sk-key", codexAuthAPIKey},
		{"missing_required", "", "", codexAuthFatal},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := decideCodexLiveAuthForCI(tc.oauth, tc.key, false, map[bool]string{true: "1", false: ""}[tc.want == codexAuthFatal])
			if got.mode != tc.want {
				t.Fatalf("mode = %d, want %d (%s)", got.mode, tc.want, got.message)
			}
		})
	}
}

func TestSeedCodexOAuthAuth(t *testing.T) {
	home := t.TempDir()
	payload := `{"auth_mode":"chatgpt","tokens":{"refresh_token":"sentinel"}}`
	if err := seedCodexOAuthAuth(home, payload); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(home, "auth.json"))
	if err != nil || string(got) != payload {
		t.Fatalf("auth.json = %q, err=%v", got, err)
	}
	if mode := fileMode(t, filepath.Join(home, "auth.json")); mode != 0o600 {
		t.Fatalf("auth mode = %o, want 600", mode)
	}
	if err := seedCodexOAuthAuth(home, "not-json"); err == nil {
		t.Fatal("malformed OAuth payload must fail")
	}
}

func fileMode(t *testing.T, path string) os.FileMode {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return info.Mode().Perm()
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

func TestCodexLiveWorkflowPinsOnlyExecToLuna(t *testing.T) {
	shim := codexLiveWorkflowExecShim(t)
	root := t.TempDir()
	realCodex := filepath.Join(root, "real-codex")
	logPath := filepath.Join(root, "codex-argv.log")
	if err := os.WriteFile(realCodex, []byte("#!/usr/bin/env bash\nset -euo pipefail\nprintf '%s\\n' \"$*\" >> \"$SPACEDOCK_CODEX_LOG\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	shimPath := filepath.Join(root, "codex")
	if err := os.WriteFile(shimPath, []byte(shim), 0o755); err != nil {
		t.Fatal(err)
	}

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(shimPath, args...)
		cmd.Env = append(os.Environ(),
			"SPACEDOCK_CODEX_REAL_BIN="+realCodex,
			"SPACEDOCK_CODEX_LOG="+logPath,
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Codex shim %v failed: %v\n%s", args, err, out)
		}
	}
	run("login", "--with-api-key")
	run("login", "status", "exec")
	run("plugin", "list", "exec")
	run("plugin", "add", "exec")
	run("exec", "--json", "prompt")
	run("--ask-for-approval", "on-request", "exec", "--json", "prompt")
	run("--dangerously-bypass-approvals-and-sandbox", "exec", "--json", "prompt")

	gotBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Split(strings.TrimSpace(string(gotBytes)), "\n")
	want := []string{
		"login --with-api-key",
		"login status exec",
		"plugin list exec",
		"plugin add exec",
		"exec --model gpt-5.6-luna -c model_reasoning_effort=\"max\" --json prompt",
		"--ask-for-approval on-request exec --model gpt-5.6-luna -c model_reasoning_effort=\"max\" --json prompt",
		"--dangerously-bypass-approvals-and-sandbox exec --model gpt-5.6-luna -c model_reasoning_effort=\"max\" --json prompt",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("Codex shim argv = %q, want %q", got, want)
	}
	for _, line := range got[4:] {
		if strings.Count(line, "--model gpt-5.6-luna") != 1 {
			t.Fatalf("pinned Codex exec argv = %q, want exactly one Luna model flag", line)
		}
		if strings.Count(line, `-c model_reasoning_effort="max"`) != 1 {
			t.Fatalf("pinned Codex exec argv = %q, want exactly one maximum-effort setting", line)
		}
	}
}

func TestRuntimeLiveClaudeShimSetsMaximumEffort(t *testing.T) {
	workflow, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "runtime-live-e2e.yml"))
	if err != nil {
		t.Fatal(err)
	}
	const start = `cat > "$shim_dir/claude" <<'SH'`
	bodyStart := strings.Index(string(workflow), start) + len(start) + 1
	endAt := strings.Index(string(workflow[bodyStart:]), "\n          SH\n")
	lines := strings.Split(string(workflow[bodyStart:bodyStart+endAt]), "\n")
	for i := range lines {
		lines[i] = strings.TrimPrefix(lines[i], "          ")
	}
	root := t.TempDir()
	realClaude := filepath.Join(root, "real-claude")
	shimPath := filepath.Join(root, "claude")
	files := map[string]string{realClaude: "#!/usr/bin/env bash\nprintf '%s\\n' \"$*\"\n", shimPath: strings.Join(lines, "\n") + "\n"}
	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	var got []byte
	for _, args := range [][]string{{"--version"}, {"--model", "claude-sonnet-5", "--help"}} {
		cmd := exec.Command(shimPath, args...)
		cmd.Env = append(os.Environ(), "SPACEDOCK_CLAUDE_REAL_BIN="+realClaude)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Claude shim %v failed: %v\n%s", args, err, out)
		} else {
			got = append(got, out...)
		}
	}
	if want := "--effort max --version\n--effort max --model claude-sonnet-5 --help\n"; string(got) != want {
		t.Fatalf("Claude shim argv = %q, want %q", got, want)
	}
}

func codexLiveWorkflowExecShim(t *testing.T) string {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve Codex live env test source")
	}
	workflowPath := filepath.Join(filepath.Dir(source), "..", "..", ".github", "workflows", "runtime-live-e2e.yml")
	workflow, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("read Codex live workflow: %v", err)
	}
	const start = `cat > "$shim_dir/codex" <<'SH'`
	startAt := strings.Index(string(workflow), start)
	if startAt < 0 {
		t.Fatalf("Codex live workflow has no exec shim heredoc")
	}
	bodyStart := startAt + len(start)
	if bodyStart < len(workflow) && workflow[bodyStart] == '\r' {
		bodyStart++
	}
	if bodyStart >= len(workflow) || workflow[bodyStart] != '\n' {
		t.Fatalf("Codex live workflow shim heredoc is not newline-delimited")
	}
	bodyStart++
	endAt := strings.Index(string(workflow[bodyStart:]), "\n          SH\n")
	if endAt < 0 {
		t.Fatalf("Codex live workflow exec shim heredoc has no terminator")
	}
	body := string(workflow[bodyStart : bodyStart+endAt])
	const yamlIndent = "          "
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, yamlIndent) {
			t.Fatalf("Codex live workflow shim line %d lacks YAML indentation: %q", i+1, line)
		}
		lines[i] = strings.TrimPrefix(line, yamlIndent)
	}
	return strings.Join(lines, "\n") + "\n"
}

func TestCodexLiveEnvOmitsEmptyAPIKey(t *testing.T) {
	env := codexLiveEnv("/tmp/codex-home", "/tmp/home", "", "", codexAuthOAuth)
	if _, ok := envValue(env, "OPENAI_API_KEY"); ok {
		t.Fatal("OPENAI_API_KEY must be omitted for local subscription auth")
	}
}

func TestCodexLiveEnvDropsForeignRuntimeMarkers(t *testing.T) {
	for key, value := range map[string]string{"CODEX_THREAD_ID": "codex", "CLAUDECODE": "claude", "PI_CODING_AGENT": "pi",
		"PI_CODING_AGENT_DIR": "/parent/pi", "CODEX_HOME": "/parent/codex", "HOME": "/parent/home", "OPENAI_API_KEY": "parent-key", "PATH": "/parent/bin"} {
		t.Setenv(key, value)
	}

	env := codexLiveEnv("/target/codex", "/target/home", "/spacedock/bin", "target-key", codexAuthAPIKey)

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

func codexLiveEnv(codexHome, home, pathPrefix, openAIAPIKey string, mode codexAuthMode) []string {
	env := cleanEnviron("CODEX_HOME", "HOME", "OPENAI_API_KEY", "CODEX_AUTH_JSON", "PATH",
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
	if mode == codexAuthAPIKey && openAIAPIKey != "" {
		env = append(env, "OPENAI_API_KEY="+openAIAPIKey)
	}
	return env
}
