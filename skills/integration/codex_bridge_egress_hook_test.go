// ABOUTME: Codex Bridge egress packaging tests — Codex must use its own non-async
// ABOUTME: hooks and a silent wrapper around the shared Spacedock egress command.
package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCodexManifestPointsAtCodexBridgeHooks(t *testing.T) {
	manifestPath := filepath.Join(repoRoot(t), ".codex-plugin", "plugin.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read codex manifest: %v", err)
	}

	var manifest struct {
		Hooks string `json:"hooks"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parse codex manifest: %v", err)
	}
	if manifest.Hooks != "./hooks/codex-hooks.json" {
		t.Fatalf("codex manifest hooks = %q, want ./hooks/codex-hooks.json", manifest.Hooks)
	}
	if manifest.Hooks == "./hooks/hooks.json" {
		t.Fatalf("codex manifest must not reuse hooks/hooks.json; that file contains async Claude hooks")
	}

	if _, err := os.Stat(filepath.Join(repoRoot(t), strings.TrimPrefix(manifest.Hooks, "./"))); err != nil {
		t.Fatalf("codex manifest hooks target is not present: %v", err)
	}
}

func TestCodexBridgeHooksAreNonAsyncAndCallPluginRootWrapper(t *testing.T) {
	path := filepath.Join(repoRoot(t), "hooks", "codex-hooks.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read codex hooks: %v", err)
	}

	var cfg struct {
		Hooks map[string][]struct {
			Hooks []map[string]any `json:"hooks"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse codex hooks: %v", err)
	}
	if strings.Contains(string(data), `"async"`) {
		t.Fatalf("codex hooks must not contain any async field; Codex skips async command hooks:\n%s", data)
	}

	for _, event := range []string{"SessionStart", "UserPromptSubmit", "PostToolUse", "SubagentStart", "SubagentStop", "Stop"} {
		groups := cfg.Hooks[event]
		if len(groups) == 0 {
			t.Fatalf("codex hooks missing %s", event)
		}
		for _, group := range groups {
			if len(group.Hooks) == 0 {
				t.Fatalf("codex hook %s has no command handlers", event)
			}
			for _, handler := range group.Hooks {
				if _, ok := handler["async"]; ok {
					t.Fatalf("codex hook %s contains async field; Codex skips async command hooks", event)
				}
				if handler["type"] != "command" {
					t.Fatalf("codex hook %s handler type = %v, want command", event, handler["type"])
				}
				cmd, ok := handler["command"].(string)
				if !ok || cmd == "" {
					t.Fatalf("codex hook %s handler has no command: %#v", event, handler)
				}
				if strings.Contains(cmd, "CLAUDE_PLUGIN_ROOT") {
					t.Fatalf("codex hook %s command must not depend on Claude env: %q", event, cmd)
				}
				for _, want := range []string{"PLUGIN_ROOT", "scripts/codex-bridge-events.sh"} {
					if !strings.Contains(cmd, want) {
						t.Fatalf("codex hook %s command %q missing %q", event, cmd, want)
					}
				}
			}
		}
	}
}

func TestCodexBridgeWrapperCallsSharedEgressEmitterSilently(t *testing.T) {
	root := repoRoot(t)
	wrapper := filepath.Join(root, "scripts", "codex-bridge-events.sh")
	fi, err := os.Stat(wrapper)
	if err != nil {
		t.Fatalf("codex wrapper missing: %v", err)
	}
	if fi.Mode()&0o111 == 0 {
		t.Fatalf("codex wrapper must be executable: mode %v", fi.Mode())
	}

	binDir := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "argv.log")
	fake := filepath.Join(binDir, "spacedock")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nprintf '%s\\n' \"$*\" > \"$SPACEDOCK_FAKE_ARGV_LOG\"\nprintf 'stdout leak\\n'\nprintf 'stderr leak\\n' >&2\nexit 42\n"), 0o755); err != nil {
		t.Fatalf("write fake spacedock: %v", err)
	}

	cmd := exec.Command("bash", wrapper)
	cmd.Stdin = strings.NewReader(`{"hook_event_name":"SessionStart","session_id":"codex-parent-session","cwd":"/repo/spacedock"}`)
	cmd.Env = append(os.Environ(),
		"PATH="+binDir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"SPACEDOCK_BIN=",
		"SPACEDOCK_FAKE_ARGV_LOG="+logPath,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("codex wrapper must remain observe-only even when emitter fails: %v\n%s", err, out)
	}
	if string(out) != "" {
		t.Fatalf("codex wrapper must be silent; got %q", out)
	}

	argv, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("fake spacedock was not invoked: %v", err)
	}
	if got, want := strings.TrimSpace(string(argv)), "bridge egress emit --host codex"; got != want {
		t.Fatalf("spacedock argv = %q, want %q", got, want)
	}
}

func TestCodexBridgeEgressMinimalPayloadFixture(t *testing.T) {
	path := filepath.Join(repoRoot(t), "skills", "integration", "testdata", "codex", "bridge-egress-minimal-session-start.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read codex fixture: %v", err)
	}

	var payload struct {
		SessionID     string `json:"session_id"`
		CWD           string `json:"cwd"`
		HookEventName string `json:"hook_event_name"`
		Source        string `json:"source"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("parse codex fixture: %v", err)
	}
	if payload.SessionID == "" || payload.CWD == "" || payload.HookEventName != "SessionStart" || payload.Source != "startup" {
		t.Fatalf("minimal Codex SessionStart fixture lost required fields: %+v", payload)
	}
	if strings.Contains(string(data), `"tool_name"`) || strings.Contains(string(data), `"Read"`) {
		t.Fatalf("minimal Codex fixture must not imply Read/PostToolUse marker support before live proof:\n%s", data)
	}
}
