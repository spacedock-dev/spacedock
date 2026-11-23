// ABOUTME: Offline controls and pure helpers for the isolated-daemon Paseo live harness.
package ensigncycle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/release"
)

func defaultPaseoLiveProvider() string {
	return "codex/" + release.LiveModels["codex.exec"]
}

func paseoDaemonConfigEntries(listenAddr string) [][2]string {
	return [][2]string{
		{"daemon.listen", listenAddr},
		{"daemon.relay.enabled", "false"},
		{"features.webUi.enabled", "false"},
		{"features.dictation.enabled", "false"},
		{"features.voiceMode.enabled", "false"},
	}
}

func paseoLiveEnv(paseoHome, spacedockBin, codexHome, openAIKey string) []string {
	env := cleanEnviron("PORT", "CLAUDECODE", "CODEX_THREAD_ID", "PI_CODING_AGENT", "PI_CODING_AGENT_DIR", "CODEX_HOME", "CODEX_AUTH_JSON", "OPENAI_API_KEY", "GIT_SSH_COMMAND")
	env = append(dropEnvPrefix(env, "PASEO_"), "PASEO_HOME="+paseoHome, "GIT_SSH_COMMAND=ssh -o BatchMode=yes")
	if codexHome != "" {
		env = append(env, "CODEX_HOME="+codexHome)
	}
	if openAIKey != "" {
		env = append(env, "OPENAI_API_KEY="+openAIKey)
	}
	return withBinaryOnPath(env, spacedockBin)
}

func paseoRunArgs(workspaceID, provider, thinking, title string, labels, envVars []string, prompt string) []string {
	args := []string{"run", "--no-wait", "--workspace", workspaceID, "--provider", provider}
	if strings.TrimSpace(thinking) != "" {
		args = append(args, "--thinking", thinking)
	}
	args = append(args, "--title", title)
	for _, label := range labels {
		args = append(args, "--label", label)
	}
	for _, kv := range envVars {
		args = append(args, "--env", kv)
	}
	return append(args, "--json", prompt)
}

// runPaseoJSON invokes paseo with separate stdout and stderr buffers so benign
// exit-0 stderr lines (such as "Using workspace wks_...") never corrupt JSON decoding.
func runPaseoJSON(paseoBin, paseoHome string, env []string, dst any, args ...string) (string, error) {
	cmd := exec.Command(paseoBin, append([]string{"--home", paseoHome}, args...)...)
	cmd.Env = env
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return stderr.String(), fmt.Errorf("paseo %s failed: %w\nstdout:\n%s\nstderr:\n%s", strings.Join(args, " "), err, stdout.String(), stderr.String())
	}
	if dst != nil {
		if err := json.Unmarshal(stdout.Bytes(), dst); err != nil {
			return stderr.String(), fmt.Errorf("decode paseo %s stdout JSON: %w\nstdout:\n%s\nstderr:\n%s", strings.Join(args, " "), err, stdout.String(), stderr.String())
		}
	}
	return stderr.String(), nil
}

func decidePaseoLivePreconditions(paseoBinFound bool, provider string, codexAuth codexLiveAuthDecision, required string) (string, string) {
	if !paseoBinFound {
		msg := "paseo CLI not found on PATH (set SPACEDOCK_PASEO_BIN or install @getpaseo/cli)"
		if required != "" {
			return "fatal", msg + " while SPACEDOCK_PASEO_LIVE_REQUIRED is set"
		}
		return "skip", msg
	}
	if strings.HasPrefix(provider, "codex/") {
		if codexAuth.mode == codexAuthSkip {
			return "skip", codexAuth.message
		}
		if codexAuth.mode == codexAuthFatal {
			return "fatal", codexAuth.message
		}
	}
	return "run", ""
}

func TestPaseoLiveConfigEnvAndRunArgs(t *testing.T) {
	wantModel := release.LiveModels["codex.exec"]
	if wantModel == "" || defaultPaseoLiveProvider() != "codex/"+wantModel {
		t.Fatalf("defaultPaseoLiveProvider() = %q, want codex/%s", defaultPaseoLiveProvider(), wantModel)
	}
	wantCfg := map[string]string{
		"daemon.listen": "127.0.0.1:41234", "daemon.relay.enabled": "false",
		"features.webUi.enabled": "false", "features.dictation.enabled": "false", "features.voiceMode.enabled": "false",
	}
	entries := paseoDaemonConfigEntries("127.0.0.1:41234")
	if len(entries) != len(wantCfg) {
		t.Fatalf("config count = %d, want %d", len(entries), len(wantCfg))
	}
	for _, kv := range entries {
		if strings.HasPrefix(kv[0], "--") || wantCfg[kv[0]] != kv[1] {
			t.Errorf("config %s=%q, want %q", kv[0], kv[1], wantCfg[kv[0]])
		}
	}

	for k, v := range map[string]string{
		"PASEO_AGENT_ID": "a", "PASEO_HOME": "/old", "PASEO_HOST": "127.0.0.1:6767", "PORT": "6767",
		"CLAUDECODE": "1", "CODEX_THREAD_ID": "t1", "PI_CODING_AGENT": "1", "PI_CODING_AGENT_DIR": "/pi", "GITHUB_WORKSPACE": "/repo",
	} {
		t.Setenv(k, v)
	}
	env := paseoLiveEnv("/iso/paseo", "/bin/dir/spacedock", "/iso/codex", "")
	for _, scrubbed := range []string{"PASEO_AGENT_ID", "PASEO_HOST", "PORT", "CLAUDECODE", "CODEX_THREAD_ID", "PI_CODING_AGENT", "PI_CODING_AGENT_DIR", "GITHUB_WORKSPACE"} {
		if v, ok := envValue(env, scrubbed); ok {
			t.Errorf("paseoLiveEnv leaked %s=%q", scrubbed, v)
		}
	}
	assertEnvValue(t, env, "PASEO_HOME", "/iso/paseo")
	assertEnvValue(t, env, "CODEX_HOME", "/iso/codex")
	assertEnvValue(t, env, "GIT_SSH_COMMAND", "ssh -o BatchMode=yes")
	if pathVal, ok := envValue(env, "PATH"); !ok || !strings.HasPrefix(pathVal, "/bin/dir"+string(os.PathListSeparator)) {
		t.Errorf("PATH = %q, want /bin/dir prefix", pathVal)
	}

	withThink := strings.Join(paseoRunArgs("wks_1", "antigravity/gemini-4-neon", "high", "T", []string{"role=ensign"}, []string{"SPACEDOCK_BIN=/x/s"}, "p"), "\x00")
	if !strings.Contains(withThink, "--provider\x00antigravity/gemini-4-neon\x00--thinking\x00high") || !strings.Contains(withThink, "--json\x00p") {
		t.Fatalf("unexpected run args: %q", withThink)
	}
	for _, a := range paseoRunArgs("wks_1", "codex/m", "", "T", nil, nil, "p") {
		if a == "--thinking" {
			t.Fatal("empty thinking must omit --thinking")
		}
	}
}

func TestPaseoRunJSONAndPreconditions(t *testing.T) {
	if runtime.GOOS != "windows" {
		dir := t.TempDir()
		okBin, failBin := filepath.Join(dir, "paseo-ok"), filepath.Join(dir, "paseo-fail")
		_ = os.WriteFile(okBin, []byte("#!/bin/sh\necho 'Using workspace wks_abc' >&2\necho '{\"agentId\":\"agt-1\",\"status\":\"running\"}'\n"), 0o755)
		_ = os.WriteFile(failBin, []byte("#!/bin/sh\necho '{\"error\":{\"code\":\"AGENT_CREATE_FAILED\"}}' >&2\nexit 1\n"), 0o755)
		var res struct{ AgentID, Status string }
		stderr, err := runPaseoJSON(okBin, filepath.Join(dir, "home"), os.Environ(), &res, "run", "--json", "hello")
		if err != nil || res.AgentID != "agt-1" || !strings.Contains(stderr, "Using workspace wks_abc") {
			t.Fatalf("runPaseoJSON ok = (%+v, stderr=%q, err=%v)", res, stderr, err)
		}
		if _, err := runPaseoJSON(failBin, filepath.Join(dir, "home"), os.Environ(), &res, "run", "--json", "hello"); err == nil || !strings.Contains(err.Error(), "AGENT_CREATE_FAILED") {
			t.Fatalf("expected stderr error, got %v", err)
		}
	}
	skipAuth, fatalAuth, okAuth := decideCodexLiveAuthForCI("", "", false, ""), decideCodexLiveAuthForCI("", "", false, "1"), decideCodexLiveAuthForCI("", "sk-test", false, "1")
	for _, tc := range []struct {
		bin      bool
		provider string
		auth     codexLiveAuthDecision
		req      string
		want     string
	}{
		{false, "codex/m", okAuth, "", "skip"}, {false, "codex/m", okAuth, "1", "fatal"},
		{true, "codex/m", skipAuth, "", "skip"}, {true, "codex/m", fatalAuth, "1", "fatal"},
		{true, "antigravity/gemini-4-neon", skipAuth, "", "run"},
	} {
		if got, _ := decidePaseoLivePreconditions(tc.bin, tc.provider, tc.auth, tc.req); got != tc.want {
			t.Errorf("precondition(%v,%s,%s) = %q, want %q", tc.bin, tc.provider, tc.req, got, tc.want)
		}
	}
}
