package dispatch

import (
	"reflect"
	"testing"
)

func TestCodexMultiAgentV2SpawnInputAlwaysIsolatesFreshDispatch(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{name: "absent", raw: `{"prompt":"pointer","name":"spacedock-ensign-task-validation"}`},
		{name: "all", raw: `{"prompt":"pointer","name":"spacedock-ensign-task-validation","fork_turns":"all"}`},
		{name: "numeric", raw: `{"prompt":"pointer","name":"spacedock-ensign-task-validation","fork_turns":3}`},
		{name: "future overrides", raw: `{"prompt":"pointer","name":"spacedock-ensign-task-validation","model":"gpt-5.5","reasoning_effort":"high"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			spawn, err := CodexMultiAgentV2SpawnInput([]byte(tc.raw))
			if err != nil {
				t.Fatalf("CodexMultiAgentV2SpawnInput: %v", err)
			}
			wantArgs := map[string]string{
				"task_name":  "spacedock_ensign_task_validation",
				"message":    "pointer",
				"fork_turns": "none",
			}
			if got := spawn.ToolArgs(); !reflect.DeepEqual(got, wantArgs) {
				t.Fatalf("tool args = %#v, want %#v", got, wantArgs)
			}
		})
	}
	if _, mutable := reflect.TypeOf(CodexMultiAgentV2Spawn{}).FieldByName("ForkTurns"); mutable {
		t.Fatal("CodexMultiAgentV2Spawn must not expose a mutable ForkTurns field")
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
