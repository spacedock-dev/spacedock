package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

// TestJourneyDeltaCommandRejectsMissingArgs proves the CLI's argument-validation
// exit code, mirroring the sibling journey-costs command's shape.
func TestJourneyDeltaCommandRejectsMissingArgs(t *testing.T) {
	stubFind := func(string) (string, error) { return "", nil }
	stubPost := func([]string) error { return nil }
	if code := journeyDelta([]string{"ledger.json", "--metrics-dir", t.TempDir()}, stubFind, stubPost); code != 2 {
		t.Fatalf("journeyDelta with a missing --pr exit = %d, want 2", code)
	}
}

// TestJourneyDeltaCommandCreatesNewCommentWhenNoneExists and
// TestJourneyDeltaCommandUpdatesExistingCommentWhenFound are the find-by-marker
// replacement for --edit-last: journeyDelta's own branching logic (not a
// generic argv pass-through) must call the CREATE args when find reports no
// existing comment, and the UPDATE args (targeting the exact found id) when it
// does — proving the job posts once and then edits the SAME comment in place,
// located by content rather than by "the poster's last comment on the PR."
func TestJourneyDeltaCommandCreatesNewCommentWhenNoneExists(t *testing.T) {
	ledgerPath := writeDeltaFixtureLedger(t)
	metricsDir := t.TempDir()
	writeMetricRecord(t, metricsDir, journeymetrics.Record{
		ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6",
		MetricsState: journeymetrics.StateMeasured, Outcome: journeymetrics.Outcome{Status: "passed"},
		Turns: 20, Tokens: journeymetrics.TokenTotals{Total: 1600},
	})

	find := func(string) (string, error) { return "", nil }
	var calls [][]string
	post := func(args []string) error {
		calls = append(calls, args)
		return nil
	}

	if code := journeyDelta([]string{ledgerPath, "--metrics-dir", metricsDir, "--pr", "42"}, find, post); code != 0 {
		t.Fatalf("journeyDelta exit = %d, want 0", code)
	}
	if len(calls) != 1 {
		t.Fatalf("post invocations = %d, want 1", len(calls))
	}
	joined := strings.Join(calls[0], " ")
	if !strings.Contains(joined, "pr comment 42") || strings.Contains(joined, "PATCH") {
		t.Fatalf("with no existing comment, journeyDelta must CREATE (pr comment), got: %v", calls[0])
	}
}

func TestJourneyDeltaCommandUpdatesExistingCommentWhenFound(t *testing.T) {
	ledgerPath := writeDeltaFixtureLedger(t)
	metricsDir := t.TempDir()
	writeMetricRecord(t, metricsDir, journeymetrics.Record{
		ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6",
		MetricsState: journeymetrics.StateMeasured, Outcome: journeymetrics.Outcome{Status: "passed"},
		Turns: 25, Tokens: journeymetrics.TokenTotals{Total: 1700},
	})

	find := func(string) (string, error) { return "987654321", nil }
	var calls [][]string
	post := func(args []string) error {
		calls = append(calls, args)
		return nil
	}

	if code := journeyDelta([]string{ledgerPath, "--metrics-dir", metricsDir, "--pr", "42"}, find, post); code != 0 {
		t.Fatalf("journeyDelta exit = %d, want 0", code)
	}
	if len(calls) != 1 {
		t.Fatalf("post invocations = %d, want 1", len(calls))
	}
	joined := strings.Join(calls[0], " ")
	if !strings.Contains(joined, "issues/comments/987654321") || !strings.Contains(joined, "PATCH") {
		t.Fatalf("with an existing comment id, journeyDelta must UPDATE that exact comment via gh api PATCH, got: %v", calls[0])
	}
}

func writeDeltaFixtureLedger(t *testing.T) string {
	t.Helper()
	ledger := journeymetrics.Ledger{
		Scenarios: []journeymetrics.ScenarioLedgerEntry{
			{
				ScenarioID: "shallow-boot-window",
				Observations: []journeymetrics.Record{
					{
						ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6",
						Turns: 18, Tokens: journeymetrics.TokenTotals{Total: 1562}, CapturedAt: "2026-06-20T00:00:00Z",
						RunURL: "https://github.com/spacedock-dev/spacedock/actions/runs/27931963802",
					},
				},
			},
		},
	}
	data, err := json.Marshal(ledger)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "journey-costs-previous.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
