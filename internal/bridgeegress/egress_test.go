package bridgeegress

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestEmitWritesEventSchemaAndClaudeMarker(t *testing.T) {
	root := t.TempDir()
	entityPath := filepath.Join(root, "docs", "spacedock", "linear-drc-review", "drc-3467", "index.md")
	Emit([]byte(`{
		"cwd":`+quote(root)+`,
		"hook_event_name":"PostToolUse",
		"session_id":"ses-1",
		"agent_type":"spacedock:ensign",
		"tool_name":"Read",
		"source":"hook",
		"tool_input":{"file_path":`+quote(entityPath)+`}
	}`), fixedOptions("claude"))

	var event Event
	readLastEvent(t, root, &event)
	if event.Timestamp != "2026-07-01T01:02:03Z" {
		t.Fatalf("timestamp = %q", event.Timestamp)
	}
	if event.TS != event.Timestamp {
		t.Fatalf("ts = %q, want timestamp %q", event.TS, event.Timestamp)
	}
	if event.Host != "claude" || event.Event != "PostToolUse" || event.SessionID != "ses-1" {
		t.Fatalf("event identity mismatch: %+v", event)
	}
	if event.ActorID != "ses-1" {
		t.Fatalf("actor_id = %q, want Claude session id", event.ActorID)
	}
	if event.AgentType != "spacedock:ensign" || event.Detail.Tool != "Read" || event.Detail.Source != "hook" {
		t.Fatalf("event detail mismatch: %+v", event)
	}

	var marker Marker
	readMarker(t, root, "ses-1", &marker)
	if marker.Host != "claude" || marker.SessionID != "ses-1" || marker.ActorID != "ses-1" {
		t.Fatalf("marker identity mismatch: %+v", marker)
	}
	if marker.Workflow != "linear-drc-review" || marker.Entity != "drc-3467" {
		t.Fatalf("marker target mismatch: %+v", marker)
	}
}

func TestEmitMalformedOrIncompletePayloadNoops(t *testing.T) {
	root := t.TempDir()
	for _, input := range [][]byte{
		[]byte(`{`),
		[]byte(`{"cwd":` + quote(root) + `}`),
		[]byte(`{"event":"SessionStart"}`),
	} {
		Emit(input, fixedOptions("claude"))
	}
	if _, err := os.Stat(filepath.Join(root, "_bridge", "events.jsonl")); !os.IsNotExist(err) {
		t.Fatalf("events.jsonl exists after malformed/incomplete payload: %v", err)
	}
}

func TestEmitInvalidActorIDStillWritesEventButNoMarker(t *testing.T) {
	root := t.TempDir()
	Emit([]byte(`{
		"cwd":`+quote(root)+`,
		"event":"PostToolUse",
		"session_id":"bad/id",
		"agent_type":"spacedock:ensign",
		"tool_name":"Read",
		"tool_input":{"file_path":"docs/spacedock/wf/task.md"}
	}`), fixedOptions("claude"))

	var event Event
	readLastEvent(t, root, &event)
	if event.ActorID != "" {
		t.Fatalf("unsafe actor_id = %q, want empty", event.ActorID)
	}
	if _, err := os.Stat(filepath.Join(root, "_bridge", "sessions")); !os.IsNotExist(err) {
		t.Fatalf("sessions dir exists for unsafe actor id: %v", err)
	}
}

func TestEmitDoesNotCreateMarkerFromIncidentalFilePath(t *testing.T) {
	root := t.TempDir()
	Emit([]byte(`{
		"cwd":`+quote(root)+`,
		"hook_event_name":"PostToolUse",
		"session_id":"ses-edit",
		"agent_type":"spacedock:ensign",
		"tool_name":"Edit",
		"tool_input":{"file_path":"docs/spacedock/wf/task.md"}
	}`), fixedOptions("claude"))

	var event Event
	readLastEvent(t, root, &event)
	if event.Event != "PostToolUse" || event.Detail.Tool != "Edit" {
		t.Fatalf("event mismatch: %+v", event)
	}
	if _, err := os.Stat(filepath.Join(root, "_bridge", "sessions", "ses-edit.json")); !os.IsNotExist(err) {
		t.Fatalf("Edit file_path should not create a marker: %v", err)
	}
}

func TestEmitMarkerFirstWriteWins(t *testing.T) {
	root := t.TempDir()
	first := `{"cwd":` + quote(root) + `,"hook_event_name":"PostToolUse","session_id":"ses-1","agent_type":"spacedock:ensign","tool_name":"Read","tool_input":{"file_path":"docs/spacedock/wf/first.md"}}`
	second := `{"cwd":` + quote(root) + `,"hook_event_name":"PostToolUse","session_id":"ses-1","agent_type":"spacedock:ensign","tool_name":"Read","tool_input":{"file_path":"docs/spacedock/wf/second.md"}}`

	Emit([]byte(first), fixedOptions("claude"))
	Emit([]byte(second), fixedOptions("claude"))

	var marker Marker
	readMarker(t, root, "ses-1", &marker)
	if marker.Entity != "first" {
		t.Fatalf("marker overwritten: %+v", marker)
	}
}

func TestEmitSkipsArchiveMarkers(t *testing.T) {
	root := t.TempDir()
	Emit([]byte(`{
		"cwd":`+quote(root)+`,
		"event":"PostToolUse",
		"session_id":"ses-arch",
		"agent_type":"spacedock:ensign",
		"tool_input":{"file_path":"docs/spacedock/wf/_archive/old.md"}
	}`), fixedOptions("claude"))
	if _, err := os.Stat(filepath.Join(root, "_bridge", "sessions", "ses-arch.json")); !os.IsNotExist(err) {
		t.Fatalf("archive marker exists: %v", err)
	}
}

func TestDeriveEntitySupportsWorkflowLocalSplitRootAndFolderForm(t *testing.T) {
	root := t.TempDir()
	workflowDir := filepath.Join(root, "docs", "dev")
	if err := os.MkdirAll(filepath.Join(workflowDir, ".spacedock-state", "wire-egress"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workflowDir, "README.md"), []byte("---\nstate: .spacedock-state\n---\n# Dev\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(workflowDir, ".spacedock-state", "wire-egress", "index.md")

	workflow, entity, ok := DeriveEntity(root, entityPath)
	if !ok {
		t.Fatalf("DeriveEntity did not recognize split-root folder path")
	}
	if workflow != "dev" || entity != "wire-egress" {
		t.Fatalf("DeriveEntity = (%q,%q), want (dev,wire-egress)", workflow, entity)
	}
}

func TestDeriveEntitySupportsDotStateWithoutReadme(t *testing.T) {
	workflow, entity, ok := DeriveEntity("/repo", "/repo/docs/dev/.spacedock-state/flat-task.md")
	if !ok || workflow != "dev" || entity != "flat-task" {
		t.Fatalf("DeriveEntity fallback = (%q,%q,%v), want (dev,flat-task,true)", workflow, entity, ok)
	}
}

func TestDeriveEntityRejectsAbsolutePathOutsideCWD(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "docs", "spacedock", "wf", "task.md")
	workflow, entity, ok := DeriveEntity(root, outside)
	if ok {
		t.Fatalf("DeriveEntity accepted outside path as (%q,%q)", workflow, entity)
	}
}

func TestEmitExplicitEntityPathCombinesSessionAndAgent(t *testing.T) {
	root := t.TempDir()
	Emit([]byte(`{
		"cwd":`+quote(root)+`,
		"event":"SubagentStart",
		"session_id":"parent",
		"agent_id":"agent-7",
		"agent_type":"spacedock:ensign",
		"entity_path":"docs/spacedock/dev/task.md"
	}`), fixedOptions("codex"))

	var marker Marker
	readMarker(t, root, "parent.agent-7", &marker)
	if marker.ActorID != "parent.agent-7" || marker.SessionID != "parent" || marker.AgentID != "agent-7" {
		t.Fatalf("explicit entity_path marker identity mismatch: %+v", marker)
	}
}

func TestEmitNormalizesPiLifecycleEvents(t *testing.T) {
	root := t.TempDir()
	cases := []struct {
		raw  string
		want string
	}{
		{raw: "session_start", want: "SessionStart"},
		{raw: "session_shutdown", want: "Stop"},
		{raw: "agent_start", want: "SubagentStart"},
		{raw: "agent_end", want: "SubagentStop"},
		{raw: "turn_start", want: "UserPromptSubmit"},
		{raw: "turn_end", want: "Stop"},
		{raw: "tool_execution_start", want: "PostToolUse"},
		{raw: "tool_execution_end", want: "PostToolUse"},
		{raw: "tool_call", want: "PostToolUse"},
		{raw: "tool_result", want: "PostToolUse"},
		{raw: "future_pi_event", want: "future_pi_event"},
	}

	for _, tc := range cases {
		Emit([]byte(`{"cwd":`+quote(root)+`,"event":`+quote(tc.raw)+`,"session_id":"pi-ses"}`), fixedOptions("pi"))

		var event Event
		readLastEvent(t, root, &event)
		if event.Host != "pi" || event.Event != tc.want {
			t.Fatalf("Pi event %q normalized to host=%q event=%q, want pi/%q", tc.raw, event.Host, event.Event, tc.want)
		}
	}
}

func TestEmitAppendsAndTruncatesEvents(t *testing.T) {
	root := t.TempDir()
	opts := fixedOptions("claude")
	opts.MaxLines = 3
	opts.KeepLines = 2
	for _, event := range []string{"one", "two", "three", "four"} {
		Emit([]byte(`{"cwd":`+quote(root)+`,"event":`+quote(event)+`,"session_id":"ses"}`), opts)
	}

	data, err := os.ReadFile(filepath.Join(root, "_bridge", "events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Fatalf("kept %d lines, want 2:\n%s", len(lines), data)
	}
	var first, second Event
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(lines[1]), &second); err != nil {
		t.Fatal(err)
	}
	if first.Event != "three" || second.Event != "four" {
		t.Fatalf("kept events = %q,%q; want three,four", first.Event, second.Event)
	}
}

func fixedOptions(host string) Options {
	return Options{
		Host: host,
		Now: func() time.Time {
			return time.Date(2026, 7, 1, 1, 2, 3, 0, time.UTC)
		},
	}
}

func readLastEvent(t *testing.T, root string, out *Event) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "_bridge", "events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), out); err != nil {
		t.Fatalf("unmarshal event: %v\n%s", err, lines[len(lines)-1])
	}
}

func readMarker(t *testing.T, root, actorID string, out *Marker) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "_bridge", "sessions", actorID+".json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, out); err != nil {
		t.Fatalf("unmarshal marker: %v\n%s", err, data)
	}
}

func quote(s string) string {
	data, _ := json.Marshal(s)
	return string(data)
}
