package contractlint

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDispatchAckMachineryIsAbsent(t *testing.T) {
	repo := repoRoot(t)

	data, err := os.ReadFile(filepath.Join(repo, "hooks.json"))
	if err != nil {
		t.Fatal(err)
	}
	var document struct {
		Hooks map[string]json.RawMessage `json:"hooks"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}
	if _, ok := document.Hooks["SessionStart"]; !ok {
		t.Error("hooks.json lost the unrelated SessionStart hook")
	}
	for _, event := range []string{"PreToolUse", "SubagentStart"} {
		if _, ok := document.Hooks[event]; ok {
			t.Errorf("hooks.json contains removed %s dispatch hook", event)
		}
	}

	// .claude-plugin/plugin.json wants hooks again as of
	// force-boot-at-compaction-boundary: a SessionStart(compact) reminder
	// hook, unrelated to the removed PreToolUse/SubagentStart dispatch-ack
	// machinery this test otherwise guards against. The event-name check
	// above and the forbidden-marker scan below still catch a regression of
	// the actual removed machinery.
	for _, manifest := range []struct {
		path      string
		wantHooks bool
	}{
		{path: filepath.Join(".claude-plugin", "plugin.json"), wantHooks: true},
		{path: filepath.Join(".codex-plugin", "plugin.json"), wantHooks: true},
	} {
		data, err := os.ReadFile(filepath.Join(repo, manifest.path))
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(data, &fields); err != nil {
			t.Fatal(err)
		}
		_, hasHooks := fields["hooks"]
		if hasHooks != manifest.wantHooks {
			t.Errorf("%s hooks activation = %t, want %t", manifest.path, hasHooks, manifest.wantHooks)
		}
	}

	if _, err := os.Stat(filepath.Join(repo, "internal", "dispatchack")); !os.IsNotExist(err) {
		t.Errorf("internal/dispatchack still exists: %v", err)
	}
	for _, name := range []string{
		filepath.Join("internal", "dispatch", "build.go"),
		filepath.Join("internal", "dispatch", "dispatch.go"),
		filepath.Join("internal", "status", "handlers.go"),
	} {
		data, err := os.ReadFile(filepath.Join(repo, name))
		if err != nil {
			t.Fatal(err)
		}
		for _, forbidden := range []string{"dispatchack", "dispatch_ack_", "dispatch-ack", "ack-hook", "sda-"} {
			if strings.Contains(string(data), forbidden) {
				t.Errorf("%s contains removed dispatch acknowledgment marker %q", name, forbidden)
			}
		}
	}
}
