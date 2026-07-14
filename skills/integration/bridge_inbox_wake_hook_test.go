// ABOUTME: Claude Bridge-ingress wake wiring — the Stop hook that drains queued
// ABOUTME: captain intent must be SYNCHRONOUS (an async hook cannot return a decision).
package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestClaudeStopHookRegistersSynchronousInboxCheck locks the Claude durable-wake
// wiring: hooks/hooks.json must register scripts/spacedock-bridge-inbox-check.sh on
// Stop WITHOUT async, or the block decision that keeps a parked FO draining is
// dropped (async Stop hooks are fire-and-forget). The pre-existing async egress
// Stop hook may coexist; this asserts the check hook specifically is synchronous.
func TestClaudeStopHookRegistersSynchronousInboxCheck(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "hooks", "hooks.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read hooks: %v", err)
	}

	var cfg struct {
		Hooks map[string][]struct {
			Hooks []map[string]any `json:"hooks"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse hooks: %v", err)
	}

	const script = "spacedock-bridge-inbox-check.sh"
	var found, sync bool
	for _, group := range cfg.Hooks["Stop"] {
		for _, h := range group.Hooks {
			cmd, _ := h["command"].(string)
			if !strings.Contains(cmd, script) {
				continue
			}
			found = true
			if async, ok := h["async"].(bool); !ok || !async {
				sync = true
			}
		}
	}
	if !found {
		t.Fatalf("Stop hook does not register %s — a parked Claude FO is never nudged to drain:\n%s", script, data)
	}
	if !sync {
		t.Fatalf("%s Stop hook must be synchronous (no async:true); an async hook cannot return the block decision that forces a drain", script)
	}

	scriptPath := filepath.Join(root, "scripts", script)
	info, err := os.Stat(scriptPath)
	if err != nil {
		t.Fatalf("wake hook script missing: %v", err)
	}
	if info.Mode()&0o111 == 0 {
		t.Fatalf("%s is not executable (mode %v)", script, info.Mode())
	}
}
