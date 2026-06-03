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
		JourneyID:       "claude-gate-guardrail",
		Source:          "live-harness",
		Host:            "claude",
		Model:           "sonnet",
		MetricsState:    journeymetrics.StateMeasured,
		Outcome:         journeymetrics.Outcome{Status: "passed"},
		DurationMS:      12,
		Turns:           1,
		ToolCalls:       1,
		ToolCallsByName: map[string]int{"Bash": 1},
		Tokens:          journeymetrics.TokenTotals{Input: 1, Output: 2, Total: 3},
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
	if !strings.Contains(string(data), `"journey_id": "claude-gate-guardrail"`) {
		t.Fatalf("ledger did not include fixture journey:\n%s", data)
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
		JourneyID:     "unit",
		Source:        "unit-test",
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
	if err := os.WriteFile(filepath.Join(dir, record.JourneyID+".json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
}
