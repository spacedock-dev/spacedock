package release

import (
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

func fixtureScenario(id string, turns int) journeymetrics.ScenarioLedgerEntry {
	return journeymetrics.ScenarioLedgerEntry{
		ScenarioID: id,
		Source:     "live-harness",
		Observations: []journeymetrics.Record{
			{ScenarioID: id, Runtime: "claude", Model: "claude-sonnet-4-6", Turns: turns, Outcome: journeymetrics.Outcome{Status: "passed"}},
		},
	}
}

// TestDiffAddedScenariosAcceptsExactlyOneAddedScenario is AC-4's primary safety
// proof: a rebuilt ledger that equals the original PLUS exactly one new
// shallow-boot-window entry, with every pre-existing scenario byte-identical,
// must pass.
func TestDiffAddedScenariosAcceptsExactlyOneAddedScenario(t *testing.T) {
	original := journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		fixtureScenario("gate-guardrail", 5),
		fixtureScenario("shallow-boot", 3),
	}}
	rebuilt := journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		fixtureScenario("gate-guardrail", 5),
		fixtureScenario("shallow-boot", 3),
		fixtureScenario("shallow-boot-window", 3),
	}}
	if err := DiffAddedScenarios(original, rebuilt, "shallow-boot-window"); err != nil {
		t.Fatalf("expected the exactly-one-addition case to pass, got: %v", err)
	}
}

// TestDiffAddedScenariosRejectsMutatedExistingScenario is AC-4's proof that a
// rebuild quirk mutating a pre-existing historical scenario (e.g. normalizeRecord
// restamping schema_version or recomputing token totals differently) is caught,
// not silently republished.
func TestDiffAddedScenariosRejectsMutatedExistingScenario(t *testing.T) {
	original := journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		fixtureScenario("gate-guardrail", 5),
	}}
	rebuilt := journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		fixtureScenario("gate-guardrail", 6), // mutated: Turns 5 -> 6
		fixtureScenario("shallow-boot-window", 3),
	}}
	err := DiffAddedScenarios(original, rebuilt, "shallow-boot-window")
	if err == nil {
		t.Fatal("expected the mutated pre-existing scenario to be rejected")
	}
	if !strings.Contains(err.Error(), "gate-guardrail") {
		t.Fatalf("error does not name the mutated scenario: %v", err)
	}
}

// TestDiffAddedScenariosRejectsMissingExistingScenario proves a rebuild that
// silently DROPPED a pre-existing scenario is caught.
func TestDiffAddedScenariosRejectsMissingExistingScenario(t *testing.T) {
	original := journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		fixtureScenario("gate-guardrail", 5),
		fixtureScenario("shallow-boot", 3),
	}}
	rebuilt := journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		fixtureScenario("shallow-boot", 3),
		fixtureScenario("shallow-boot-window", 3),
	}}
	err := DiffAddedScenarios(original, rebuilt, "shallow-boot-window")
	if err == nil {
		t.Fatal("expected the dropped pre-existing scenario to be rejected")
	}
	if !strings.Contains(err.Error(), "gate-guardrail") {
		t.Fatalf("error does not name the missing scenario: %v", err)
	}
}

// TestDiffAddedScenariosRejectsUnexpectedExtraScenario proves a rebuild that
// added something OTHER than (or in addition to) the expected scenario is
// caught, rather than treated as a benign superset.
func TestDiffAddedScenariosRejectsUnexpectedExtraScenario(t *testing.T) {
	original := journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		fixtureScenario("gate-guardrail", 5),
	}}
	rebuilt := journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		fixtureScenario("gate-guardrail", 5),
		fixtureScenario("shallow-boot-window", 3),
		fixtureScenario("mystery-scenario", 1),
	}}
	if err := DiffAddedScenarios(original, rebuilt, "shallow-boot-window"); err == nil {
		t.Fatal("expected an unexpected extra scenario beyond the named addition to be rejected")
	}
}

// TestDiffAddedScenariosRejectsMissingExpectedAddition proves a rebuild that did
// NOT actually add the expected scenario (e.g. the extraction silently no-opped)
// is caught rather than passing vacuously.
func TestDiffAddedScenariosRejectsMissingExpectedAddition(t *testing.T) {
	original := journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		fixtureScenario("gate-guardrail", 5),
	}}
	rebuilt := journeymetrics.Ledger{Scenarios: []journeymetrics.ScenarioLedgerEntry{
		fixtureScenario("gate-guardrail", 5),
	}}
	if err := DiffAddedScenarios(original, rebuilt, "shallow-boot-window"); err == nil {
		t.Fatal("expected a rebuild that did not add the expected scenario to be rejected")
	}
}
