package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

func TestJourneyCostsCommandWritesVersionedLedger(t *testing.T) {
	metricsDir := t.TempDir()
	writeMetricRecord(t, metricsDir, journeymetrics.Record{
		SchemaVersion:   journeymetrics.RecordSchemaVersion,
		ScenarioID:      "gate-guardrail",
		Source:          "live-harness",
		Mode:            "sonnet",
		Runtime:         "claude",
		Executor:        "llm",
		Host:            "claude",
		Model:           "claude-sonnet-4-6",
		MetricsState:    journeymetrics.StateMeasured,
		Outcome:         journeymetrics.Outcome{Status: "passed"},
		DurationMS:      12,
		Turns:           1,
		ToolCalls:       1,
		ToolCallsByName: map[string]int{"Bash": 1},
		Tokens:          journeymetrics.TokenTotals{Input: 1, Output: 2, Total: 3},
	})
	writeMetricRecord(t, metricsDir, journeymetrics.Record{
		SchemaVersion:   journeymetrics.RecordSchemaVersion,
		ScenarioID:      "gate-guardrail",
		Source:          "live-harness",
		Mode:            "gpt-5-codex",
		Runtime:         "codex",
		Executor:        "llm",
		Host:            "codex",
		Model:           "gpt-5-codex",
		MetricsState:    journeymetrics.StateCharacterized,
		Outcome:         journeymetrics.Outcome{Status: "passed"},
		DurationMS:      34,
		ToolCalls:       1,
		ToolCallsByName: map[string]int{"exec_command": 1},
	})

	out := filepath.Join(t.TempDir(), "journey-costs-v1.2.3.json")
	if code := journeyCosts([]string{"1.2.3", "--metrics-dir", metricsDir, "--out", out}); code != 0 {
		t.Fatalf("journeyCosts exit = %d, want 0", code)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(data), `"artifact": "journey-costs-v1.2.3.json"`) {
		t.Fatalf("ledger did not carry the versioned artifact name:\n%s", data)
	}
	if !strings.Contains(string(data), `"scenario_id": "gate-guardrail"`) {
		t.Fatalf("ledger did not include fixture scenario:\n%s", data)
	}
	if !strings.Contains(string(data), `"scenario_count": 1`) || !strings.Contains(string(data), `"observation_count": 2`) {
		t.Fatalf("ledger did not group Claude and Codex as one scenario with two observations:\n%s", data)
	}
	if strings.Contains(string(data), "claude-gate-guardrail") {
		t.Fatalf("ledger still contains a host-prefixed logical scenario id:\n%s", data)
	}
}

func TestJourneyCostsCommandRejectsEmptyAcceptedInput(t *testing.T) {
	out := filepath.Join(t.TempDir(), "journey-costs-v1.2.3.json")
	if code := journeyCosts([]string{"1.2.3", "--metrics-dir", t.TempDir(), "--out", out}); code == 0 {
		t.Fatalf("empty metrics dir unexpectedly succeeded")
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Fatalf("empty input wrote output file; stat err=%v", err)
	}
}

func TestJourneyCostsCommandRejectsMismatchedOutputFilename(t *testing.T) {
	metricsDir := t.TempDir()
	writeMetricRecord(t, metricsDir, journeymetrics.Record{
		SchemaVersion: journeymetrics.RecordSchemaVersion,
		ScenarioID:    "unit",
		Source:        "unit-test",
		Mode:          "fake",
		Runtime:       "go-test",
		Executor:      "codified",
		Host:          "go-test",
		Model:         "fake",
		MetricsState:  journeymetrics.StateMeasured,
		Outcome:       journeymetrics.Outcome{Status: "passed"},
	})

	out := filepath.Join(t.TempDir(), "journey-costs-v9.9.9.json")
	if code := journeyCosts([]string{"1.2.3", "--metrics-dir", metricsDir, "--out", out}); code == 0 {
		t.Fatalf("mismatched output filename unexpectedly succeeded")
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Fatalf("mismatched filename wrote output file; stat err=%v", err)
	}
}

func writeMetricRecord(t *testing.T, dir string, record journeymetrics.Record) {
	t.Helper()
	data, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	name := strings.Join([]string{record.ScenarioID, record.Runtime, record.Model}, "--") + ".json"
	if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
		t.Fatal(err)
	}
}
