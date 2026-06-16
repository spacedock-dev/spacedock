package journeymetrics

import "testing"

func TestParseClaudeSurfacesStatusReadAndScopedReadCounts(t *testing.T) {
	// WITH: one `status --read` Bash tool_use and one scoped Read (offset/limit).
	got, err := ParseClaudeJSONL(readTestdata(t, "claude_read_adoption.stream.jsonl"))
	if err != nil {
		t.Fatalf("ParseClaudeJSONL: %v", err)
	}
	if got.Observation.StatusReadCalls != 1 {
		t.Errorf("status_read_calls = %d, want 1", got.Observation.StatusReadCalls)
	}
	if got.Observation.ScopedReadCalls != 1 {
		t.Errorf("scoped_read_calls = %d, want 1", got.Observation.ScopedReadCalls)
	}

	// The counts flow through BuildRecord onto the emitted Record, exactly as
	// ToolCallsByName does — so the release ledger surfaces them.
	record := BuildRecord(JourneySpec{ScenarioID: "x", Source: "s"}, BehaviorResult{Passed: true}, got.Observation)
	if record.StatusReadCalls != 1 || record.ScopedReadCalls != 1 {
		t.Errorf("record counts = (%d,%d), want (1,1)", record.StatusReadCalls, record.ScopedReadCalls)
	}

	// WITHOUT: a plain `spacedock status` Bash and a full-file Read (no offset/limit).
	none, err := ParseClaudeJSONL(readTestdata(t, "claude_no_read_adoption.stream.jsonl"))
	if err != nil {
		t.Fatalf("ParseClaudeJSONL (no adoption): %v", err)
	}
	if none.Observation.StatusReadCalls != 0 {
		t.Errorf("status_read_calls without adoption = %d, want 0", none.Observation.StatusReadCalls)
	}
	if none.Observation.ScopedReadCalls != 0 {
		t.Errorf("scoped_read_calls without adoption = %d, want 0", none.Observation.ScopedReadCalls)
	}
}

func TestStatusReadCallsDedupAcrossDeltas(t *testing.T) {
	// Real runner streams are MULTI-DELTA: the same tool_use block (same tool id)
	// is carried across two assistant rows. The status-read count dedups on the
	// same toolID key the existing tool-call dedup uses, so a repeated delta does
	// not double-count. Shape mirrors TestParseClaudeTurnsMergesToolUseAcrossDeltas.
	stream := `{"type":"assistant","message":{"id":"msg_x","model":"claude-sonnet-4-6","usage":{"input_tokens":50},"content":[{"type":"tool_use","id":"toolu_sr","name":"Bash","input":{"command":"spacedock status --read /path/entity.md --json"}}]}}
{"type":"assistant","message":{"id":"msg_x","model":"claude-sonnet-4-6","usage":{"input_tokens":50},"content":[{"type":"tool_use","id":"toolu_sr","name":"Bash","input":{"command":"spacedock status --read /path/entity.md --json"}}]}}`

	got, err := ParseClaudeJSONL([]byte(stream))
	if err != nil {
		t.Fatalf("ParseClaudeJSONL: %v", err)
	}
	if got.Observation.StatusReadCalls != 1 {
		t.Errorf("status_read_calls = %d, want 1 (the repeated delta must not double-count)", got.Observation.StatusReadCalls)
	}
}

func TestCodexCharacterizationSurfacesStatusReadCalls(t *testing.T) {
	// Codex carries the shell command in tool_call.started arguments.cmd. The
	// Bash-arg detector reuses the shared commandInvokesStatusRead helper.
	withRead := `{"type":"session.started","session_id":"c","model":"gpt-5-codex"}
{"type":"tool_call.started","call_id":"call-1","name":"exec_command","arguments":{"cmd":"spacedock status --read /path/entity.md --json"}}
{"type":"tool_call.completed","call_id":"call-1","name":"exec_command","exit_code":0}`
	got, err := CharacterizeCodexExecJSONL([]byte(withRead))
	if err != nil {
		t.Fatalf("CharacterizeCodexExecJSONL: %v", err)
	}
	if got.StatusReadCalls != 1 {
		t.Errorf("codex status_read_calls = %d, want 1", got.StatusReadCalls)
	}

	plain := `{"type":"session.started","session_id":"c","model":"gpt-5-codex"}
{"type":"tool_call.started","call_id":"call-1","name":"exec_command","arguments":{"cmd":"spacedock status --json"}}`
	none, err := CharacterizeCodexExecJSONL([]byte(plain))
	if err != nil {
		t.Fatalf("CharacterizeCodexExecJSONL (plain): %v", err)
	}
	if none.StatusReadCalls != 0 {
		t.Errorf("codex status_read_calls for plain status = %d, want 0", none.StatusReadCalls)
	}
}

func TestCommandInvokesStatusRead(t *testing.T) {
	cases := []struct {
		name string
		cmd  string
		want bool
	}{
		// Launcher-agnostic positives: spacedock, sd, the ${SPACEDOCK_BIN:-spacedock}
		// fetch-command form, and a bare status invocation.
		{"spacedock launcher", "spacedock status --read /path/entity.md", true},
		{"sd launcher", "sd status --read /path/entity.md", true},
		{"SPACEDOCK_BIN launcher", "${SPACEDOCK_BIN:-spacedock} status --read /path/entity.md", true},
		{"bare status", "status --read /path/entity.md", true},
		// Flag order independence: --read after another flag still counts.
		{"flag order independent", "spacedock status --json --read /path/entity.md", true},
		// --json default appears too; the --read token anywhere after status counts.
		{"read then json", "spacedock status --read /path/entity.md --json", true},

		// Negatives: status without --read.
		{"plain status", "spacedock status", false},
		{"status boot", "spacedock status --boot", false},
		{"status discover json", "spacedock status --discover --json", false},
		// A non-status subcommand, even one that happens to carry --read elsewhere.
		{"non-status command", "spacedock dispatch show-stage-def --stage implementation", false},
		// status as a quoted argument to a different command is NOT a status invocation.
		{"echo quoted status read", "echo 'status --read'", false},
		// --read on a non-status command is not a status --read.
		{"readonly other flag", "spacedock dispatch build --readme", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := commandInvokesStatusRead(tc.cmd); got != tc.want {
				t.Errorf("commandInvokesStatusRead(%q) = %v, want %v", tc.cmd, got, tc.want)
			}
		})
	}
}
