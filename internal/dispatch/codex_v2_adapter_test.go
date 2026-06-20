package dispatch

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCodexMultiAgentV2SpawnInputMapsBuildOutput(t *testing.T) {
	helperPrompt := "Read /tmp/spacedock-dispatch/spacedock-ensign-status-validate-determinism-implementation.md and treat its content as your assignment."
	raw := []byte(`{
		"schema_version": 2,
		"subagent_type": "general",
		"description": "implementation worker",
		"fetch_commands": [],
		"dispatch_file_path": "/tmp/spacedock-dispatch/spacedock-ensign-status-validate-determinism-implementation.md",
		"prompt": ` + mustJSON(t, helperPrompt) + `,
		"model": "haiku",
		"name": "spacedock-ensign-status-validate-determinism-implementation"
	}`)

	spawn, err := CodexMultiAgentV2SpawnInput(raw)
	if err != nil {
		t.Fatalf("CodexMultiAgentV2SpawnInput: %v", err)
	}

	if spawn.TaskName != "spacedock_ensign_status_validate_determinism_implementation" {
		t.Fatalf("task_name = %q", spawn.TaskName)
	}
	if spawn.Message != helperPrompt {
		t.Fatalf("message was not byte-identical to helper prompt:\nwant %q\ngot  %q", helperPrompt, spawn.Message)
	}
	if spawn.ForkTurns != "" {
		t.Fatalf("fork_turns should be omitted by default, got %q", spawn.ForkTurns)
	}
	if spawn.Identity.Name != "spacedock-ensign-status-validate-determinism-implementation" ||
		spawn.Identity.Slug != "status-validate-determinism" ||
		spawn.Identity.Stage != "implementation" {
		t.Fatalf("identity did not preserve helper mapping: %+v", spawn.Identity)
	}

	encoded, err := json.Marshal(spawn.ToolArgs())
	if err != nil {
		t.Fatalf("marshal tool args: %v", err)
	}
	body := string(encoded)
	for _, want := range []string{`"task_name"`, `"message"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("tool args missing %s: %s", want, body)
		}
	}
	for _, banned := range []string{`"description"`, `"subagent_type"`, `"model"`, `"name"`} {
		if strings.Contains(body, banned) {
			t.Fatalf("tool args should omit unsupported helper field %q: %s", banned, body)
		}
	}
}

func TestCodexMultiAgentV2SpawnInputRejectsSanitizedCollision(t *testing.T) {
	first := []byte(`{"prompt":"one","name":"spacedock-ensign-a-b-implementation"}`)
	second := []byte(`{"prompt":"two","name":"spacedock-ensign-a--b-implementation"}`)
	seen := map[string]string{}
	if _, err := CodexMultiAgentV2SpawnInputWithSeen(first, seen); err != nil {
		t.Fatalf("first mapping errored: %v", err)
	}
	if _, err := CodexMultiAgentV2SpawnInputWithSeen(second, seen); err == nil {
		t.Fatal("expected sanitized task_name collision to fail")
	}
}

func mustJSON(t *testing.T, v string) string {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal test JSON: %v", err)
	}
	return string(raw)
}
