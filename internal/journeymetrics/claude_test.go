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

func readTestdata(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}
