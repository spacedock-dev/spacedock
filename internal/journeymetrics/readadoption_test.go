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

func TestFoldEnsignReadAdoption(t *testing.T) {
	// FO front-door stream with a known baseline: a plain `spacedock status` (no
	// --read) and one scoped Read (offset/limit). Baseline = (status_read=0,
	// scoped_read=1) — `status --read` adoption is principally an ENSIGN behavior,
	// so the FO front-door alone registers ~0.
	foStream := `{"type":"system","subtype":"init","model":"claude-sonnet-4-6"}
{"type":"assistant","message":{"id":"msg_fo1","model":"claude-sonnet-4-6","usage":{"input_tokens":100},"content":[{"type":"tool_use","id":"toolu_fo_status","name":"Bash","input":{"command":"spacedock status --json"}}]}}
{"type":"assistant","message":{"id":"msg_fo2","model":"claude-sonnet-4-6","usage":{"input_tokens":120},"content":[{"type":"tool_use","id":"toolu_fo_read","name":"Read","input":{"file_path":"/path/entity.md","offset":10,"limit":20}}]}}`
	parsed, err := ParseClaudeJSONL([]byte(foStream))
	if err != nil {
		t.Fatalf("ParseClaudeJSONL (FO front-door): %v", err)
	}
	baseline := parsed.Observation
	if baseline.StatusReadCalls != 0 || baseline.ScopedReadCalls != 1 {
		t.Fatalf("FO baseline counts = (%d,%d), want (0,1)", baseline.StatusReadCalls, baseline.ScopedReadCalls)
	}

	// A dispatched-ensign sub-agent transcript carrying one recognized-form
	// `spacedock_launcher status --read` Bash call + one scoped Read. The on-disk
	// agent-*.jsonl shape carries extra top-level fields (parentUuid, agentId, …)
	// the map[string]json.RawMessage decode ignores.
	ensignTranscript := `{"parentUuid":"p","agentId":"a1","type":"assistant","message":{"id":"msg_en1","model":"claude-sonnet-4-6","usage":{"input_tokens":80},"content":[{"type":"tool_use","id":"toolu_en_status","name":"Bash","input":{"command":"spacedock_launcher status --read /path/entity.md --json"}}]}}
{"parentUuid":"p","agentId":"a1","type":"assistant","message":{"id":"msg_en2","model":"claude-sonnet-4-6","usage":{"input_tokens":90},"content":[{"type":"tool_use","id":"toolu_en_read","name":"Read","input":{"file_path":"/path/entity.md","offset":40,"limit":12}}]}}`

	folded, err := FoldEnsignReadAdoption(baseline, [][]byte{[]byte(ensignTranscript)})
	if err != nil {
		t.Fatalf("FoldEnsignReadAdoption: %v", err)
	}
	// Both the ensign's status --read and its scoped Read fold onto the FO counts.
	if folded.StatusReadCalls != 1 || folded.ScopedReadCalls != 2 {
		t.Errorf("folded counts = (%d,%d), want (1,2)", folded.StatusReadCalls, folded.ScopedReadCalls)
	}
	// Non-tautological by perturbation: the fold moved BOTH counters strictly above
	// the FO-only baseline. With the fold removed (a no-op FoldEnsignReadAdoption)
	// the counts stay at the baseline (0,1) and this assertion fails — the baseline
	// can move the wrong way, so a passing assertion proves the ensign contribution
	// was actually folded in, not that the numbers happen to match.
	if folded.StatusReadCalls <= baseline.StatusReadCalls || folded.ScopedReadCalls <= baseline.ScopedReadCalls {
		t.Errorf("fold did not raise counts above FO-only baseline (%d,%d) -> (%d,%d)",
			baseline.StatusReadCalls, baseline.ScopedReadCalls, folded.StatusReadCalls, folded.ScopedReadCalls)
	}

	// Zero case: an ensign transcript carrying NEITHER a status --read (a plain
	// `spacedock status`) NOR a scoped Read (a whole-file Read) folds to exactly the
	// FO-only baseline.
	noAdoption := `{"type":"assistant","message":{"id":"msg_en3","model":"claude-sonnet-4-6","usage":{"input_tokens":70},"content":[{"type":"tool_use","id":"toolu_en_plain","name":"Bash","input":{"command":"spacedock status --json"}}]}}
{"type":"assistant","message":{"id":"msg_en4","model":"claude-sonnet-4-6","usage":{"input_tokens":75},"content":[{"type":"tool_use","id":"toolu_en_full","name":"Read","input":{"file_path":"/path/entity.md"}}]}}`
	zero, err := FoldEnsignReadAdoption(baseline, [][]byte{[]byte(noAdoption)})
	if err != nil {
		t.Fatalf("FoldEnsignReadAdoption (zero case): %v", err)
	}
	if zero.StatusReadCalls != baseline.StatusReadCalls || zero.ScopedReadCalls != baseline.ScopedReadCalls {
		t.Errorf("zero-case counts = (%d,%d), want FO-only baseline (%d,%d)",
			zero.StatusReadCalls, zero.ScopedReadCalls, baseline.StatusReadCalls, baseline.ScopedReadCalls)
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
		// The dominant real launcher forms the dispatch fetch commands emit: the
		// spacedock_launcher shell function and the ./spacedock dev build.
		{"spacedock_launcher function", "spacedock_launcher status --read /path/entity.md", true},
		{"dev build ./spacedock", "./spacedock status --read /path/entity.md", true},
		// The contract-canonical QUOTED launcher (shared-core invariant). The exact-
		// token map dropped it; normalization strips the surrounding quotes so it
		// counts — the M10 regression guard.
		{"quoted SPACEDOCK_BIN launcher", `"${SPACEDOCK_BIN:-spacedock}" status --read /path/entity.md`, true},
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
		// A launcher token alone, with no following status --read, does not match.
		{"launcher token alone", "spacedock_launcher", false},
		// $SPACEDOCK is a legacy non-canonical variable, NOT the SPACEDOCK_BIN family
		// the contract pins, so it stays unrecognized.
		{"legacy SPACEDOCK var", "$SPACEDOCK status --read /path/entity.md", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := commandInvokesStatusRead(tc.cmd); got != tc.want {
				t.Errorf("commandInvokesStatusRead(%q) = %v, want %v", tc.cmd, got, tc.want)
			}
		})
	}
}
