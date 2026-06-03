package journeymetrics

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAggregateLedgerMatchesGolden(t *testing.T) {
	records := []Record{
		{
			SchemaVersion:   RecordSchemaVersion,
			ScenarioID:      "gate-guardrail",
			Source:          "live-harness",
			Mode:            ModeLLMLive,
			Runtime:         "codex",
			Executor:        "llm",
			Host:            "codex",
			Model:           "gpt-5-codex",
			MetricsState:    StateCharacterized,
			Outcome:         Outcome{Status: "passed"},
			DurationMS:      2345,
			ToolCalls:       1,
			ToolCallsByName: map[string]int{"exec_command": 1},
		},
		{
			SchemaVersion:   RecordSchemaVersion,
			ScenarioID:      "gate-guardrail",
			Source:          "live-harness",
			Mode:            ModeLLMLive,
			Runtime:         "claude",
			Executor:        "llm",
			Host:            "claude",
			Model:           "claude-sonnet-4-6",
			MetricsState:    StateMeasured,
			Outcome:         Outcome{Status: "passed"},
			DurationMS:      1234,
			Turns:           2,
			ToolCalls:       1,
			ToolCallsByName: map[string]int{"Bash": 1},
			Tokens: TokenTotals{
				Input:         100,
				Output:        50,
				CacheCreation: 20,
				CacheRead:     10,
				Total:         180,
			},
			TotalCostUSD: 0.42,
			ModelUsage: map[string]ModelUsage{
				"claude-sonnet-4-6": {
					Tokens:  TokenTotals{Input: 100, Output: 50, CacheCreation: 20, CacheRead: 10, Total: 180},
					CostUSD: 0.42,
				},
			},
		},
	}

	ledger, err := AggregateLedger("0.20.0", records, time.Date(2026, 6, 3, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("AggregateLedger: %v", err)
	}
	got, err := MarshalLedger(ledger)
	if err != nil {
		t.Fatalf("MarshalLedger: %v", err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "golden", "journey-costs.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("ledger mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}
