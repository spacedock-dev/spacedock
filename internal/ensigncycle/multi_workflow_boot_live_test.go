//go:build live

package ensigncycle

import "testing"

func runMultiWorkflowBootJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario, build func(*testing.T, string) multiWorkflowBootFixture, assert func(multiWorkflowBootObservation) error) {
	t.Helper()
	projectRoot := t.TempDir()
	fixture := build(t, projectRoot)
	ledger := newTestInvocationLedger(t, spacedockBinary(t))
	driver = driver.withInvocationLedger(t, ledger)
	result := driver.run(t, scenario, projectRoot, multiWorkflowBootPrompt(projectRoot))
	invocations := ledger.read(t)
	writeInvocationLedgerArtifact(t, result.artifactDir, invocations)
	observation := gatherMultiWorkflowBootObservation(t, fixture, invocations, result.finalMessage)
	finishLiveScenario(t, driver, scenario, result,
		durableSemantic("multi-workflow-boot-violation", assert(observation)))
}
