package ensigncycle

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// seedScenarioBlock captures the IDs declared between the machine-readable
// `<!-- seed-scenarios -->` markers in scenario-testing-principles.md.
var seedScenarioBlock = regexp.MustCompile(`(?s)<!-- seed-scenarios -->(.*?)<!-- /seed-scenarios -->`)

// seedScenarioID matches one backtick-quoted scenario ID at the start of a seed
// bullet line: `- ` + "`id`".
var seedScenarioID = regexp.MustCompile("(?m)^- `([^`]+)`")

// TestSeedScenariosDocLock is the AC-5 doc-lock: it binds the machine-readable
// `<!-- seed-scenarios -->` block in docs/specs/scenario-testing-principles.md to
// the code-backed sharedRuntimeScenarios() table. The doc declares the seed IDs as
// the human-readable face of a code-backed truth, so the two sets must be equal;
// this test reds on drift in EITHER direction — a scenario added, dropped, or
// renamed on one side without the other. It is offline (default tag): it reads the
// doc and the table, spending no model.
func TestSeedScenariosDocLock(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	docPath := filepath.Join(wd, "..", "..", "docs", "specs", "scenario-testing-principles.md")
	b, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("read docs/specs/scenario-testing-principles.md: %v", err)
	}

	block := seedScenarioBlock.FindStringSubmatch(string(b))
	if block == nil {
		t.Fatal("scenario-testing-principles.md is missing the `<!-- seed-scenarios -->` machine-readable block")
	}
	var docIDs []string
	for _, m := range seedScenarioID.FindAllStringSubmatch(block[1], -1) {
		docIDs = append(docIDs, m[1])
	}
	if len(docIDs) == 0 {
		t.Fatal("the `<!-- seed-scenarios -->` block declared no scenario IDs")
	}

	var codeIDs []string
	for _, scenario := range sharedRuntimeScenarios() {
		codeIDs = append(codeIDs, scenario.name)
	}

	sort.Strings(docIDs)
	sort.Strings(codeIDs)
	if strings.Join(docIDs, ",") != strings.Join(codeIDs, ",") {
		t.Fatalf("seed-scenarios doc block drifted from sharedRuntimeScenarios():\n  doc:  %v\n  code: %v", docIDs, codeIDs)
	}
}

// TestSharedScenarioDocsContract is the AC-5 README-contract guard: docs/dev/README.md documents
// the shared-scenario contract — how to add a scenario, what belongs in the
// host-neutral definition, what belongs in each runner, and the local Claude/Codex/Pi
// live commands. The README IS the claim here (the contract is the evergreen doc),
// so a presence check over its real text is proof at the claim's own level: it
// fails if a future edit drops a required clause. The required-clause set is
// deliberately specific (the actual command strings and section anchors) so it
// cannot be satisfied by an unrelated mention.
func TestSharedScenarioDocsContract(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	readmePath := filepath.Join(wd, "..", "..", "docs", "dev", "README.md")
	b, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("read docs/dev/README.md: %v", err)
	}
	doc := string(b)

	// The contract clauses the README must carry. Each is a concrete fact a reader
	// reproduces, not a vague topic mention.
	mustContain := []string{
		// Section anchors.
		"## Runtime Live CI",
		"### Shared runtime scenarios",
		"### Local live execution",
		// What belongs in the host-neutral definition vs. each runner.
		"sharedRuntimeScenario",
		"runner adapter",
		"codexScenarioRunners()",
		"claudeScenarioRunners()",
		"pi_shared_coverage_test.go",
		// How to add a shared scenario.
		"To add a shared runtime scenario",
		// The parity guard the contract leans on.
		"TestSharedScenarioRunnerCoverage",
		// Local live commands for BOTH hosts. `-count=1` defeats a stale Go test
		// cache replaying a prior pass without launching the model (the false-green
		// that bit the live gate); `-timeout 40m` is a LOOSE BACKSTOP above the full
		// 4-scenario serial-suite wall-time (~27m opus) — the real guard is the 120s
		// per-stage stall-watchdog, and 40m keeps the suite off Go's too-short 10m
		// default binary timeout.
		"go test -tags live -count=1 -timeout 40m -run TestLiveClaudeSharedScenarios ./internal/ensigncycle -v",
		"go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v",
		"go test -tags live -count=1 -run TestLivePiFrontDoorSmoke ./internal/ensigncycle -v",
	}
	for _, clause := range mustContain {
		if !strings.Contains(doc, clause) {
			t.Errorf("docs/dev/README.md is missing the required shared-scenario contract clause: %q", clause)
		}
	}
}
