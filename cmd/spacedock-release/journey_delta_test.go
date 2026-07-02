package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

// TestJourneyDeltaCommandReusesStickyCommentOnRepeatedRuns is AC-3's proof that
// the SAME comment id/marker is reused across repeated runs (an update, not a
// second comment): it drives journeyDelta twice — once with an initial PR-run
// fixture, once with an updated one carrying different turns/tokens — and
// asserts BOTH stubbed `gh pr comment` invocations carry the identical sticky
// argv shape (--edit-last --create-if-none against the same PR number), which is
// the mechanism that makes gh edit its own last comment rather than posting a
// new one.
func TestJourneyDeltaCommandReusesStickyCommentOnRepeatedRuns(t *testing.T) {
	ledgerPath := writePreviousLedger(t, journeymetrics.Ledger{
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
	})

	var calls [][]string
	stub := func(args []string) error {
		calls = append(calls, args)
		return nil
	}

	metricsDirA := t.TempDir()
	writeCurrentRecord(t, metricsDirA, journeymetrics.Record{
		ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6",
		MetricsState: journeymetrics.StateMeasured, Outcome: journeymetrics.Outcome{Status: "passed"},
		Turns: 20, Tokens: journeymetrics.TokenTotals{Total: 1600},
	})
	if code := journeyDelta([]string{ledgerPath, "--metrics-dir", metricsDirA, "--pr", "42"}, stub); code != 0 {
		t.Fatalf("journeyDelta (first run) exit = %d, want 0", code)
	}

	metricsDirB := t.TempDir()
	writeCurrentRecord(t, metricsDirB, journeymetrics.Record{
		ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6",
		MetricsState: journeymetrics.StateMeasured, Outcome: journeymetrics.Outcome{Status: "passed"},
		Turns: 25, Tokens: journeymetrics.TokenTotals{Total: 1700},
	})
	if code := journeyDelta([]string{ledgerPath, "--metrics-dir", metricsDirB, "--pr", "42"}, stub); code != 0 {
		t.Fatalf("journeyDelta (second run) exit = %d, want 0", code)
	}

	if len(calls) != 2 {
		t.Fatalf("gh pr comment invocations = %d, want 2", len(calls))
	}
	for i, args := range calls {
		joined := strings.Join(args, " ")
		for _, want := range []string{"pr comment 42", "--edit-last", "--create-if-none"} {
			if !strings.Contains(joined, want) {
				t.Fatalf("call %d argv %v missing %q — repeated runs must use the identical sticky-update shape", i, args, want)
			}
		}
	}
}

// TestJourneyDeltaCommandRejectsMissingArgs proves the CLI's argument-validation
// exit code, mirroring the sibling journey-costs command's shape.
func TestJourneyDeltaCommandRejectsMissingArgs(t *testing.T) {
	if code := journeyDelta([]string{"ledger.json", "--metrics-dir", t.TempDir()}, func([]string) error { return nil }); code != 2 {
		t.Fatalf("journeyDelta with a missing --pr exit = %d, want 2", code)
	}
}

func writePreviousLedger(t *testing.T, ledger journeymetrics.Ledger) string {
	t.Helper()
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

func writeCurrentRecord(t *testing.T, dir string, record journeymetrics.Record) {
	t.Helper()
	writeMetricRecord(t, dir, record)
}
