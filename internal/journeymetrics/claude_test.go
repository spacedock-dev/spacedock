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

func readTestdata(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}
