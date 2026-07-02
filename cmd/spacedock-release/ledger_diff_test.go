package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

func writeLedgerFile(t *testing.T, path string, ledger journeymetrics.Ledger) {
	t.Helper()
	data, err := json.Marshal(ledger)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func ledgerDiffFixtureScenario(id string, turns int) journeymetrics.ScenarioLedgerEntry {
	return journeymetrics.ScenarioLedgerEntry{
		ScenarioID: id,
		Source:     "live-harness",
		Observations: []journeymetrics.Record{
			{ScenarioID: id, Runtime: "claude", Model: "claude-sonnet-4-6", Turns: turns, Outcome: journeymetrics.Outcome{Status: "passed"}},
		},
	}
}

// TestLedgerDiffCommandAcceptsExactlyOneAddedScenario is AC-4's CLI-level proof
// that the pre-upload safeguard (release.DiffAddedScenarios) is actually wired
// into the `ledger-diff` subcommand the backfill invokes before any `gh release
// upload --clobber`, not just exercised as an internal package function.
func TestLedgerDiffCommandAcceptsExactlyOneAddedScenario(t *testing.T) {
	dir := t.TempDir()
	originalPath := filepath.Join(dir, "original.json")
	rebuiltPath := filepath.Join(dir, "rebuilt.json")
	writeLedgerFile(t, originalPath, journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		ledgerDiffFixtureScenario("gate-guardrail", 5),
	}})
	writeLedgerFile(t, rebuiltPath, journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		ledgerDiffFixtureScenario("gate-guardrail", 5),
		ledgerDiffFixtureScenario("shallow-boot-window", 3),
	}})

	if code := ledgerDiff([]string{originalPath, rebuiltPath, "--added", "shallow-boot-window"}); code != 0 {
		t.Fatalf("ledgerDiff exit = %d, want 0 for a rebuild that only adds the named scenario", code)
	}
}

// TestLedgerDiffCommandRejectsMutatedExistingScenario proves the CLI surfaces a
// non-zero exit — the signal the backfill's upload step must gate on — when the
// rebuild alters a pre-existing historical scenario.
func TestLedgerDiffCommandRejectsMutatedExistingScenario(t *testing.T) {
	dir := t.TempDir()
	originalPath := filepath.Join(dir, "original.json")
	rebuiltPath := filepath.Join(dir, "rebuilt.json")
	writeLedgerFile(t, originalPath, journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		ledgerDiffFixtureScenario("gate-guardrail", 5),
	}})
	writeLedgerFile(t, rebuiltPath, journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		ledgerDiffFixtureScenario("gate-guardrail", 6), // mutated
		ledgerDiffFixtureScenario("shallow-boot-window", 3),
	}})

	if code := ledgerDiff([]string{originalPath, rebuiltPath, "--added", "shallow-boot-window"}); code == 0 {
		t.Fatal("ledgerDiff exit = 0 for a rebuild that mutated a pre-existing scenario; want non-zero")
	}
}

// TestLedgerDiffCommandRejectsMissingArgs proves a miswired invocation is a
// usage error rather than a silent pass that would let an upload proceed
// unguarded.
func TestLedgerDiffCommandRejectsMissingArgs(t *testing.T) {
	if code := ledgerDiff(nil); code == 0 {
		t.Fatalf("ledgerDiff exit = 0 with no arguments; want non-zero")
	}
	if code := ledgerDiff([]string{"only-one-arg.json"}); code == 0 {
		t.Fatalf("ledgerDiff exit = 0 with a single positional argument; want non-zero")
	}
}

// TestLedgerDiffCommandRejectsMissingFile proves an unreadable ledger path fails
// loud instead of comparing against an empty ledger.
func TestLedgerDiffCommandRejectsMissingFile(t *testing.T) {
	dir := t.TempDir()
	rebuiltPath := filepath.Join(dir, "rebuilt.json")
	writeLedgerFile(t, rebuiltPath, journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		ledgerDiffFixtureScenario("gate-guardrail", 5),
	}})

	code := ledgerDiff([]string{filepath.Join(dir, "does-not-exist.json"), rebuiltPath, "--added", "shallow-boot-window"})
	if code == 0 {
		t.Fatalf("ledgerDiff exit = 0 with a missing original ledger file; want non-zero")
	}
}
