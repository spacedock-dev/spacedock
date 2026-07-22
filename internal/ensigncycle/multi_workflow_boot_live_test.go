//go:build live

package ensigncycle

import "testing"

func runMultiWorkflowBootJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario, build func(*testing.T, string) multiWorkflowBootFixture, assert func(multiWorkflowBootObservation) error) {
	t.Helper()
	projectRoot := t.TempDir()
	fixture := build(t, projectRoot)
	result := driver.run(t, scenario, projectRoot, multiWorkflowBootPrompt(projectRoot))
	observation := gatherMultiWorkflowBootObservation(t, fixture, result.commands, result.finalMessage)
	finishLiveScenario(t, driver, scenario, result,
		durableSemantic("multi-workflow-boot-violation", assert(observation)))
}
