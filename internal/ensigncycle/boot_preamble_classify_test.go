// ABOUTME: AC-2 offline proof — a contract-file find-hunt classifies through the
// ABOUTME: runner's failure classification as a distinct broad-search preamble diagnosis.
package ensigncycle

import (
	"strings"
	"testing"
)

// TestClassifyBootPreambleFailureBroadSearch is AC-2's offline classification
// proof: a captured boot stream carrying a contract-file `find /` hunt (the
// ledger's instance 7/8/11/16 shape — after resolving its workflow the FO greps
// the filesystem for a contract reference file instead of proceeding on the
// fixture/plugin content already read) classifies through
// classifyBootPreambleFailure as a distinct broad-search preamble diagnosis
// naming the sweep — never silence (which would leave the caller's opaque stall
// path or a downstream scenario assertion to misdiagnose it).
func TestClassifyBootPreambleFailureBroadSearch(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveClaudeSharedScenariosself-evidence-merge-triage1234567890/001"
	stream := strings.Join([]string{
		streamLine(`spacedock status --boot --workflow-dir ` + fixtureRoot),
		streamLine(`find / -maxdepth 6 -iname "fo-merge-core.md"`),
	}, "\n")

	err := classifyBootPreambleFailure(stream, fixtureRoot)
	if err == nil {
		t.Fatal("classifier passed a contract-file find-hunt boot — want a distinct broad-search preamble diagnosis")
	}
	if !strings.Contains(err.Error(), "FO broad-searched the filesystem at boot") {
		t.Errorf("diagnosis must name the broad-search preamble class, got: %v", err)
	}
}

// TestClassifyBootPreambleFailureWrongRootTakesPriority proves wrong-root wander
// is classified first, ahead of broad-search: a stream carrying both a wrong-root
// cd AND a subsequent find-hunt returns the wrong-root diagnosis — the earliest,
// most specific signature of the two.
func TestClassifyBootPreambleFailureWrongRootTakesPriority(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveClaudeSharedScenariosfiling1234567890/001"
	const realRepo = "/home/runner/work/spacedock/spacedock"
	stream := strings.Join([]string{
		streamLine(`cd ` + realRepo + ` && spacedock status --discover`),
		streamLine(`find / -maxdepth 6 -iname "fo-merge-core.md"`),
	}, "\n")

	err := classifyBootPreambleFailure(stream, fixtureRoot)
	if err == nil {
		t.Fatal("classifier passed a wrong-root + find-hunt boot — want the wrong-root diagnosis")
	}
	if !strings.Contains(err.Error(), "FO booted the wrong root") {
		t.Errorf("wrong-root wander must classify first, got: %v", err)
	}
}

// TestClassifyBootPreambleFailureCleanBootPasses proves a clean, on-fixture boot
// with no wander and no sweep classifies as nil, leaving the caller's normal
// stall/assertion path intact — the classifier must not false-red an ordinary run.
func TestClassifyBootPreambleFailureCleanBootPasses(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveClaudeSharedScenariosgate-guardrail1234567890/001"
	stream := strings.Join([]string{
		streamLine(`spacedock --version`),
		streamLine(`spacedock status --boot --workflow-dir ` + fixtureRoot),
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Read","input":{"file_path":"` + fixtureRoot + `/README.md"}}]}}`,
	}, "\n")

	if err := classifyBootPreambleFailure(stream, fixtureRoot); err != nil {
		t.Errorf("classifier red a clean, on-fixture boot: %v", err)
	}
}
