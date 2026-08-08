package release

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"
)

type liveClaim struct{ step, selector, claim string }

var liveClaims = []liveClaim{
	{"Run live Claude E2E", "TestLiveEnsignCycle", "full-ensign-cycle"},
	{"Run live Claude E2E", "TestLiveDefaultHeadlessStopsAtGate", "default-headless-gate-stop-and-withdrawn-recovery"},
	{"Run live Claude E2E", "TestLiveZeroDiscoverReportsAndStops", "zero-discovery"},
	{"Run live Claude shared scenarios", "TestLiveClaudeSharedScenarios", "claude-common-journeys"},
	{"Run live Claude substrate proofs", "TestLiveMergedTeamModeDispatch", "claude-merged-agent-dispatch"},
	{"Run live Claude substrate proofs", "TestLiveBareReachable", "claude-bare-dispatch"},
	{"Run live Claude substrate proofs", "TestLiveBreakGlassShimRecovery", "claude-break-glass-recovery"},
	{"Verify Codex resolver against installed plugin", "TestCodexResolveManifestAgainstInstalledHost", "codex-current-checkout-manifest-resolution"},
	{"Run live Codex shared scenarios", "TestLiveCodexSharedScenarios", "codex-common-journeys"},
}

var offlineControls = []string{
	"TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle", "TestClaudeRejectionFlowTODOModelScope",
	"TestClaudeSonnetGateGuardrailTODOModelScope", "TestClaudeTODOModelScope",
	"TestCleanupKeepMovingRootRetainsOnlyFailures", "TestCodexLiveRunnerExecArgvEnablesMultiAgentV2",
	"TestCodexLiveRunnerUsesSpacedockFrontDoorBeforeHostArgs", "TestPiIntercomPackageRootDefaultsBesideSubagents",
	"TestPiLiveEnvDropsForeignRuntimeMarkers", "TestPiLiveEnvScrubsAmbientPiSubagentMarkers",
	"TestPiLiveSmokePromptRequiresExactStageReportHeading", "TestPiSharedScenarioCoverage",
	"TestSharedScenarioRunnerCoverage", "TestShallowBootFixtureContainsOnlyHeldGate",
}

func TestRuntimeLiveWorkflowEverySelectedMinuteBuysNamedEvidence(t *testing.T) {
	if err := assertNamedLiveEvidence(readWorkflow(t, "runtime-live-e2e.yml")); err != nil {
		t.Fatal(err)
	}
}

func TestRuntimeLiveWorkflowRunsDeterministicControlsOffline(t *testing.T) {
	if err := assertOfflineControls(readWorkflow(t, "runtime-live-e2e.yml")); err != nil {
		t.Fatal(err)
	}
}

func TestRuntimeLiveWorkflowNamedEvidenceMutationControls(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	mutations := map[string]func(string) string{
		"removed claim selector": func(s string) string {
			return strings.Replace(s, "TestLiveBreakGlassShimRecovery", "TestLiveMergedTeamModeDispatch", 1)
		},
		"lost offline Pi coverage": func(s string) string {
			return strings.Replace(s, "|TestPiSharedScenarioCoverage", "", 1)
		},
		"legacy PTY flag": func(s string) string {
			return strings.Replace(s, "DISABLE_AUTOUPDATER: \"1\"", "DISABLE_AUTOUPDATER: \"1\"\n      CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS: \"1\"", 1)
		},
		"offline live tag": func(s string) string {
			return strings.Replace(s, "go test ./internal/ensigncycle -count=1 -run", "go test -tags live ./internal/ensigncycle -count=1 -run", 1)
		},
		"lost Codex PR consumer": removeLastCodexArtifactName,
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			adversarial := mutate(live)
			if assertNamedLiveEvidence(adversarial) == nil && assertOfflineControls(adversarial) == nil {
				t.Fatal("mutation escaped the workflow guards")
			}
		})
	}
}

func assertNamedLiveEvidence(workflow string) error {
	expected := map[string][]string{}
	claimOwners := map[string]string{}
	for _, item := range liveClaims {
		if claimOwners[item.claim] != "" {
			return fmt.Errorf("claim %q has two owners", item.claim)
		}
		claimOwners[item.claim] = item.selector
		expected[item.step] = append(expected[item.step], item.selector)
	}
	expected["Run deterministic live-harness controls offline"] = offlineControls
	for _, step := range parseWorkflowSteps(workflow) {
		got := selectedTests(step.run)
		want, owned := expected[step.name]
		if len(got) > 0 && !owned {
			return fmt.Errorf("step %q selects unowned tests %v", step.name, got)
		}
		if owned && strings.Join(got, "|") != strings.Join(sorted(want), "|") {
			return fmt.Errorf("step %q selects %v, want %v", step.name, got, sorted(want))
		}
		delete(expected, step.name)
	}
	if len(expected) != 0 {
		return fmt.Errorf("workflow lacks owned selector steps %v", expected)
	}
	for _, dead := range []string{"TestLivePty", "TestLivePiRecordedGateLifecycle", "TestLivePiSubagentEnsignSmoke", "TestLivePiFrontDoorSmoke", "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS", "pty-team-mode", "Install tmux", "journey-metrics/pi", "inputs.effort", "CI-E2E-PI"} {
		if activeYAMLText(workflow, dead) {
			return fmt.Errorf("workflow retains dead surface %q", dead)
		}
	}
	for _, job := range parseWorkflowJobs(workflow) {
		if job.name == "journey-delta-comment" && strings.Join(sorted(job.needs), ",") == "claude-live,codex-live" && strings.Count(workflow, "runtime-live-e2e-codex-live") >= 2 {
			return nil
		}
	}
	return fmt.Errorf("journey-delta-comment does not consume Claude and Codex metrics")
}

func assertOfflineControls(workflow string) error {
	step, ok := stepNamed(parseWorkflowSteps(workflow), "Run deterministic live-harness controls offline")
	if !ok || strings.Contains(step.run, "-tags") || strings.Join(selectedTests(step.run), "|") != strings.Join(sorted(offlineControls), "|") {
		return fmt.Errorf("the 14 deterministic controls are not an exact untagged command")
	}
	return nil
}

var runSelector = regexp.MustCompile(`(?:^|\s)-run\s+['"]?([A-Za-z0-9_|]+)`)

func selectedTests(script string) []string {
	seen := map[string]bool{}
	for _, command := range executableShellCommands(script) {
		if match := runSelector.FindStringSubmatch(command); len(match) == 2 {
			for _, name := range strings.Split(match[1], "|") {
				seen[name] = true
			}
		}
	}
	var names []string
	for name := range seen {
		names = append(names, name)
	}
	return sorted(names)
}

func stepNamed(steps []workflowStep, name string) (workflowStep, bool) {
	for _, step := range steps {
		if step.name == name {
			return step, true
		}
	}
	return workflowStep{}, false
}

func sorted(values []string) []string {
	copyOfValues := append([]string(nil), values...)
	sort.Strings(copyOfValues)
	return copyOfValues
}

func activeYAMLText(workflow, text string) bool {
	for _, line := range strings.Split(workflow, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && strings.Contains(line, text) {
			return true
		}
	}
	return false
}

func removeLastCodexArtifactName(workflow string) string {
	const name = "runtime-live-e2e-codex-live"
	index := strings.LastIndex(workflow, name)
	if index < 0 {
		return workflow
	}
	return workflow[:index] + "runtime-live-e2e-claude-live-sonnet" + workflow[index+len(name):]
}
