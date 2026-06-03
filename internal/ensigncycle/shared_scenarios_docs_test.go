package ensigncycle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSharedScenarioDocsContract is the AC-6 guard: docs/dev/README.md documents
// the shared-scenario contract — how to add a scenario, what belongs in the
// host-neutral definition, what belongs in each runner, and the local Claude/Codex
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
		// How to add a shared scenario.
		"To add a shared runtime scenario",
		// The parity guard the contract leans on.
		"TestSharedScenarioRunnerCoverage",
		// Local live commands for BOTH hosts.
		"go test -tags live -run TestLiveClaudeSharedScenarios ./internal/ensigncycle -v",
		"go test -tags live -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v",
	}
	for _, clause := range mustContain {
		if !strings.Contains(doc, clause) {
			t.Errorf("docs/dev/README.md is missing the required shared-scenario contract clause: %q", clause)
		}
	}
}
