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
		Mode:            journeymetrics.ModeLLMLive,
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
		Mode:            journeymetrics.ModeLLMLive,
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
		Mode:          journeymetrics.ModeCodified,
		Runtime:       "go-test",
		Executor:      "codified",
		Host:          "go-test",
		Model:         "fake-model",
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

// TestJourneyCostsCommandAggregatesPerRunSubdirectories proves the AC-2 fix: once
// release.yml downloads each discovered run into its OWN subdirectory instead of
// flat-copying, journeymetrics.ReadRecordsDir's existing recursive walk aggregates
// both runs' observations for the SAME scenario/model without collision, and each
// observation carries the distinct, non-empty run_id/run_url of the source run that
// produced it — proving "traceable to more than one run" against actual
// provenance data, not just a count.
func TestJourneyCostsCommandAggregatesPerRunSubdirectories(t *testing.T) {
	metricsDir := t.TempDir()
	runADir := filepath.Join(metricsDir, "27931963802")
	runBDir := filepath.Join(metricsDir, "28432388663")

	base := journeymetrics.Record{
		SchemaVersion: journeymetrics.RecordSchemaVersion,
		ScenarioID:    "shallow-boot-window",
		Source:        "live-harness",
		Mode:          journeymetrics.ModeLLMLive,
		Runtime:       "claude",
		Executor:      "llm",
		Host:          "claude",
		Model:         "claude-sonnet-4-6",
		MetricsState:  journeymetrics.StateMeasured,
		Outcome:       journeymetrics.Outcome{Status: "passed"},
	}
	runA := base
	runA.Turns, runA.Tokens = 18, journeymetrics.TokenTotals{CacheCreation: 1562}
	runA.RunID = "27931963802"
	runA.RunURL = "https://github.com/spacedock-dev/spacedock/actions/runs/27931963802"
	runA.CapturedAt = "2026-06-20T00:00:00Z"
	runB := base
	runB.Turns, runB.Tokens = 20, journeymetrics.TokenTotals{CacheCreation: 1576}
	runB.RunID = "28432388663"
	runB.RunURL = "https://github.com/spacedock-dev/spacedock/actions/runs/28432388663"
	runB.CapturedAt = "2026-06-27T00:00:00Z"

	writeMetricRecord(t, runADir, runA)
	writeMetricRecord(t, runBDir, runB)

	out := filepath.Join(t.TempDir(), "journey-costs-v1.2.3.json")
	if code := journeyCosts([]string{"1.2.3", "--metrics-dir", metricsDir, "--out", out}); code != 0 {
		t.Fatalf("journeyCosts exit = %d, want 0", code)
	}
	ledger := readLedger(t, out)
	entry := findScenario(t, ledger, "shallow-boot-window")
	if len(entry.Observations) != 2 {
		t.Fatalf("per-run-subdirectory shape must yield 2 observations for shallow-boot-window, got %d: %+v", len(entry.Observations), entry.Observations)
	}
	seenRunIDs := map[string]bool{}
	for _, obs := range entry.Observations {
		if obs.RunID == "" || obs.RunURL == "" {
			t.Fatalf("observation missing run provenance: %+v", obs)
		}
		seenRunIDs[obs.RunID] = true
	}
	if !seenRunIDs["27931963802"] || !seenRunIDs["28432388663"] {
		t.Fatalf("observations did not carry the two distinct source run ids, got %+v", seenRunIDs)
	}
}

func readLedger(t *testing.T, path string) journeymetrics.Ledger {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var ledger journeymetrics.Ledger
	if err := json.Unmarshal(data, &ledger); err != nil {
		t.Fatalf("parse ledger %s: %v\n%s", path, err, data)
	}
	return ledger
}

func findScenario(t *testing.T, ledger journeymetrics.Ledger, scenarioID string) journeymetrics.ScenarioLedgerEntry {
	t.Helper()
	for _, entry := range ledger.Scenarios {
		if entry.ScenarioID == scenarioID {
			return entry
		}
	}
	t.Fatalf("ledger has no scenario %q; scenarios: %+v", scenarioID, ledger.Scenarios)
	return journeymetrics.ScenarioLedgerEntry{}
}

func writeMetricRecord(t *testing.T, dir string, record journeymetrics.Record) {
	t.Helper()
	data, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	name := strings.Join([]string{record.ScenarioID, record.Runtime, record.Model}, "--") + ".json"
	if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
		t.Fatal(err)
	}
}
