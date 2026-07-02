package journeymetrics

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseClaudeTranscriptUsesTerminalUsageAndDedupesSplitAssistantRows(t *testing.T) {
	data := readTestdata(t, "claude_terminal_split.stream.jsonl")

	got, err := ParseClaudeJSONL(data)
	if err != nil {
		t.Fatalf("ParseClaudeJSONL: %v", err)
	}

	if got.Observation.Turns != 2 {
		t.Errorf("turns = %d, want 2", got.Observation.Turns)
	}
	if got.Observation.ToolCalls != 1 {
		t.Errorf("tool calls = %d, want 1", got.Observation.ToolCalls)
	}
	if got.Observation.ToolCallsByName["Bash"] != 1 {
		t.Errorf("Bash calls = %d, want 1", got.Observation.ToolCallsByName["Bash"])
	}
	if got.Observation.Tokens != (TokenTotals{Input: 500, Output: 70, CacheCreation: 30, CacheRead: 40, Total: 640}) {
		t.Errorf("terminal tokens = %+v, want terminal result usage", got.Observation.Tokens)
	}
	if got.Observation.TotalCostUSD != 0.123 {
		t.Errorf("total cost = %v, want 0.123", got.Observation.TotalCostUSD)
	}
	if got.AssistantUsage != (TokenTotals{Input: 300, Output: 30, CacheCreation: 5, CacheRead: 20, Total: 355}) {
		t.Errorf("assistant diagnostics = %+v, want deduped assistant totals", got.AssistantUsage)
	}
	model := got.Observation.ModelUsage["claude-sonnet-4-6"]
	if model.Tokens != (TokenTotals{Input: 500, Output: 70, CacheCreation: 30, CacheRead: 40, Total: 640}) {
		t.Errorf("model tokens = %+v, want terminal modelUsage totals", model.Tokens)
	}
	if model.CostUSD != 0.123 {
		t.Errorf("model cost = %v, want 0.123", model.CostUSD)
	}
}

func TestParseClaudeTranscriptFallsBackToAssistantUsageWhenTerminalResultMissing(t *testing.T) {
	data := readTestdata(t, "claude_no_terminal.stream.jsonl")

	got, err := ParseClaudeJSONL(data)
	if err != nil {
		t.Fatalf("ParseClaudeJSONL: %v", err)
	}

	if got.Observation.Turns != 2 {
		t.Errorf("turns = %d, want 2", got.Observation.Turns)
	}
	if got.Observation.ToolCalls != 2 {
		t.Errorf("tool calls = %d, want 2", got.Observation.ToolCalls)
	}
	if got.Observation.ToolCallsByName["Read"] != 1 || got.Observation.ToolCallsByName["Bash"] != 1 {
		t.Errorf("tool calls by name = %+v, want Read=1 Bash=1", got.Observation.ToolCallsByName)
	}
	if got.Observation.Tokens != (TokenTotals{Input: 18, Output: 8, Total: 26}) {
		t.Errorf("fallback tokens = %+v, want deduped assistant usage", got.Observation.Tokens)
	}
}

func TestParseClaudeJSONLSkipsNonJSONStderrLines(t *testing.T) {
	base := readTestdata(t, "claude_no_terminal.stream.jsonl")

	// The front-door launch banner is written to stderr, which the live claude
	// runner folds into the same pipe as the stdout stream-json (to capture
	// launch errors like a 401). The parser must skip these non-JSON lines, not
	// error on them — a banner prefix must yield the same parse as the bare stream.
	banner := "spacedock dev · launching claude as your first officer\n" +
		"Workflow: none detected\n" +
		"claude is your first officer — ask it for the queue and next steps.\n"
	withBanner := append([]byte(banner), base...)

	want, err := ParseClaudeJSONL(base)
	if err != nil {
		t.Fatalf("baseline ParseClaudeJSONL: %v", err)
	}
	got, err := ParseClaudeJSONL(withBanner)
	if err != nil {
		t.Fatalf("ParseClaudeJSONL with stderr banner prefix: %v", err)
	}
	if got.Observation.Turns != want.Observation.Turns {
		t.Errorf("turns with banner = %d, want %d (banner must not change the parse)", got.Observation.Turns, want.Observation.Turns)
	}
	if got.Observation.ToolCalls != want.Observation.ToolCalls {
		t.Errorf("tool calls with banner = %d, want %d", got.Observation.ToolCalls, want.Observation.ToolCalls)
	}
	if got.Observation.Tokens != want.Observation.Tokens {
		t.Errorf("tokens with banner = %+v, want %+v", got.Observation.Tokens, want.Observation.Tokens)
	}
}

func TestParseClaudeTurnsPreservesPerTurnContextAndDedupes(t *testing.T) {
	data := readTestdata(t, "claude_terminal_split.stream.jsonl")

	turns, err := ParseClaudeTurns(data)
	if err != nil {
		t.Fatalf("ParseClaudeTurns: %v", err)
	}
	// Two distinct assistant messages (msg_1 appears twice — deduped, msg_2 once);
	// the terminal result row is not a turn.
	if len(turns) != 2 {
		t.Fatalf("turns = %d, want 2 (deduped msg_1 + msg_2, no result row)", len(turns))
	}
	// msg_1 carried a Bash tool_use with its own per-turn usage (NOT the run sum).
	if turns[0].ID != "msg_1" {
		t.Errorf("turn[0].ID = %q, want msg_1", turns[0].ID)
	}
	if turns[0].Usage != (TokenTotals{Input: 100, Output: 10, CacheCreation: 5, CacheRead: 20, Total: 135}) {
		t.Errorf("turn[0].Usage = %+v, want the single turn's usage (not the run sum)", turns[0].Usage)
	}
	if turns[0].Context() != 125 { // input 100 + cache_read 20 + cache_creation 5
		t.Errorf("turn[0].Context() = %d, want 125 (input+cache_read+cache_creation)", turns[0].Context())
	}
	if len(turns[0].ToolNames) != 1 || turns[0].ToolNames[0] != "Bash" {
		t.Errorf("turn[0].ToolNames = %v, want [Bash]", turns[0].ToolNames)
	}
	if turns[1].ID != "msg_2" {
		t.Errorf("turn[1].ID = %q, want msg_2", turns[1].ID)
	}
	if len(turns[1].ToolNames) != 0 {
		t.Errorf("turn[1].ToolNames = %v, want none (text-only turn)", turns[1].ToolNames)
	}
}

func TestParseClaudeTurnsMergesToolUseAcrossDeltas(t *testing.T) {
	// Real runner streams are MULTI-DELTA: the same message id appears on several
	// `assistant` rows, the first delta carries a `thinking` block with no tool_use,
	// and the tool_use block lands on a LATER delta. The per-delta `usage` is
	// identical across deltas. ParseClaudeTurns must MERGE the later-delta tool_use
	// names into the turn — taking only the first delta (and skipping the rest) drops
	// the tool_use entirely, which would make assertNoTeamCreateBeforeGreet blind to a
	// TeamCreate (a hollow lazy-TeamCreate proof). Shape mirrors the committed real
	// captures (testdata/sonnet_teamdelete_hang.stream.jsonl).
	stream := `{"type":"assistant","message":{"id":"msg_x","model":"claude-sonnet-4-6","usage":{"input_tokens":50,"cache_read_input_tokens":1000,"cache_creation_input_tokens":89000},"content":[{"type":"thinking","thinking":"deciding to create a team"}]}}
{"type":"assistant","message":{"id":"msg_x","model":"claude-sonnet-4-6","usage":{"input_tokens":50,"cache_read_input_tokens":1000,"cache_creation_input_tokens":89000},"content":[{"type":"tool_use","id":"toolu_tc","name":"TeamCreate","input":{"team_name":"eager"}}]}}`

	turns, err := ParseClaudeTurns([]byte(stream))
	if err != nil {
		t.Fatalf("ParseClaudeTurns: %v", err)
	}
	if len(turns) != 1 {
		t.Fatalf("turns = %d, want 1 (the two deltas are one message)", len(turns))
	}
	// The tool_use lands on the SECOND delta — it must still be surfaced.
	if len(turns[0].ToolNames) != 1 || turns[0].ToolNames[0] != "TeamCreate" {
		t.Fatalf("turn[0].ToolNames = %v, want [TeamCreate] (merged from the later delta) — taking only the first delta loses it", turns[0].ToolNames)
	}
	// Usage is consistent across deltas; the merged turn keeps it.
	if turns[0].Usage.CacheCreation != 89000 {
		t.Errorf("turn[0].Usage.CacheCreation = %d, want 89000 (consistent across deltas)", turns[0].Usage.CacheCreation)
	}
}

func TestParseClaudeTurnsExtractsSkillNamesAcrossDeltas(t *testing.T) {
	// A Skill tool_use carries the invoked skill in input.skill (e.g. the FO's
	// `Skill(skill="spacedock:present-gate")`). Like ToolNames and ReadTargets, the
	// block can land on a LATER delta of a multi-delta message, so ParseClaudeTurns
	// must MERGE it — a first-delta-only parse would make the skill invocation
	// invisible to a caller asserting on the greet's Skill sequence.
	stream := `{"type":"assistant","message":{"id":"msg_s","model":"claude-opus-4-8","usage":{"input_tokens":8,"cache_read_input_tokens":5000},"content":[{"type":"thinking","thinking":"present the ready gate"}]}}
{"type":"assistant","message":{"id":"msg_s","model":"claude-opus-4-8","usage":{"input_tokens":8,"cache_read_input_tokens":5000},"content":[{"type":"tool_use","id":"toolu_skill","name":"Skill","input":{"skill":"spacedock:present-gate"}}]}}`

	turns, err := ParseClaudeTurns([]byte(stream))
	if err != nil {
		t.Fatalf("ParseClaudeTurns: %v", err)
	}
	if len(turns) != 1 {
		t.Fatalf("turns = %d, want 1 (the two deltas are one message)", len(turns))
	}
	// The Skill lands on the SECOND delta — its skill argument must be surfaced.
	if len(turns[0].SkillNames) != 1 || turns[0].SkillNames[0] != "spacedock:present-gate" {
		t.Fatalf("turn[0].SkillNames = %v, want [spacedock:present-gate] (merged from the later delta) — taking only the first delta loses it", turns[0].SkillNames)
	}
	// A Skill block is not a read; it must not pollute ReadTargets.
	if len(turns[0].ReadTargets) != 0 {
		t.Errorf("turn[0].ReadTargets = %v, want none (a Skill invocation is not a read)", turns[0].ReadTargets)
	}
}

func readTestdata(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}
