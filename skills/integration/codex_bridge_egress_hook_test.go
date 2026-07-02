// ABOUTME: Codex Bridge egress packaging tests — Codex must use its own non-async
// ABOUTME: hooks and call the shared Spacedock egress command without plugin-root state.
package integration

import (
	"encoding/json"
	"os"
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

func TestCodexBridgeHooksAreNonAsyncAndCallEgressDirectly(t *testing.T) {
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
				for _, forbidden := range []string{"PLUGIN_ROOT", "scripts/codex-bridge-events.sh"} {
					if strings.Contains(cmd, forbidden) {
						t.Fatalf("codex hook %s command must not depend on plugin checkout paths (%q): %q", event, forbidden, cmd)
					}
				}
				for _, want := range []string{"SPACEDOCK_BIN", "bridge egress emit --host codex"} {
					if !strings.Contains(cmd, want) {
						t.Fatalf("codex hook %s command %q missing %q", event, cmd, want)
					}
				}
			}
		}
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
