package release

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type liveCadenceRow struct {
	Cadence     string `yaml:"cadence"`
	Model       string `yaml:"model"`
	Effort      string `yaml:"effort"`
	Environment string `yaml:"environment"`
}

type liveCadenceWorkflow struct {
	On struct {
		WorkflowDispatch struct {
			Inputs map[string]struct {
				Default string   `yaml:"default"`
				Options []string `yaml:"options"`
			} `yaml:"inputs"`
		} `yaml:"workflow_dispatch"`
	} `yaml:"on"`
	Jobs map[string]struct {
		If          string            `yaml:"if"`
		Env         map[string]string `yaml:"env"`
		Environment struct {
			Name string `yaml:"name"`
		} `yaml:"environment"`
		Strategy struct {
			Matrix map[string]yaml.Node `yaml:"matrix"`
		} `yaml:"strategy"`
	} `yaml:"jobs"`
}

const (
	pullRequestCadence = `${{ github.event_name == 'pull_request' && 'pull-request' || inputs.live_cadence }}`
	claudeCadenceModel = `${{ (github.event_name == 'pull_request' || inputs.live_cadence == 'sonnet') && 'claude-sonnet-5' || 'claude-opus-4-8' }}`
	claudeCadenceEnv   = `${{ (github.event_name == 'pull_request' || inputs.live_cadence == 'sonnet') && 'CI-E2E' || 'CI-E2E-OPUS' }}`
	claudeCadenceIf    = `${{ github.event_name == 'pull_request' || inputs.live_cadence == 'sonnet' || inputs.live_cadence == 'opus-pre-release' }}`
	codexCadenceIf     = `${{ github.event_name == 'pull_request' || inputs.live_cadence == 'sonnet' }}`
	piCadenceIf        = `${{ github.event_name == 'workflow_dispatch' && inputs.live_cadence == 'pi' }}`
)

func TestRuntimeLiveWorkflowHasOneExplicitClaudeCadence(t *testing.T) {
	workflow := readWorkflow(t, "runtime-live-e2e.yml")
	if err := assertOneClaudeCadence(workflow); err != nil {
		t.Fatal(err)
	}
	mutations := [][2]string{
		{"include:\n", "os: [ubuntu-latest, macos-latest]\n        include:\n"},
		{"include:\n", "model: [claude-opus-4-8]\n        include:\n"},
		{"include:\n", "exclude:\n          - model: claude-sonnet-5\n        include:\n"},
		{"environment: " + claudeCadenceEnv, "environment: " + claudeCadenceEnv + "\n          - model: claude-opus-4-8\n            effort: max\n            environment: CI-E2E-OPUS"},
		{"claude-sonnet-5", "sonnet"},
		{"effort: max", "effort: high"},
		{claudeCadenceEnv, "CI-E2E"},
		{"'pull-request'", "'sonnet'"},
	}
	for _, mutation := range mutations {
		if err := assertOneClaudeCadence(strings.Replace(workflow, mutation[0], mutation[1], 1)); err == nil {
			t.Fatalf("mutation %q escaped the one-row cadence guard", mutation[0])
		}
	}
}

func assertOneClaudeCadence(workflow string) error {
	var parsed liveCadenceWorkflow
	if err := yaml.Unmarshal([]byte(workflow), &parsed); err != nil {
		return err
	}
	input, ok := parsed.On.WorkflowDispatch.Inputs["live_cadence"]
	if !ok || input.Default != "sonnet" || fmt.Sprint(input.Options) != fmt.Sprint([]string{"sonnet", "opus-pre-release", "pi"}) {
		return fmt.Errorf("workflow_dispatch live_cadence choice is not the approved three-value surface")
	}
	job, ok := parsed.Jobs["claude-live"]
	if !ok {
		return fmt.Errorf("workflow lacks claude-live")
	}
	matrix := job.Strategy.Matrix
	include, ok := matrix["include"]
	if len(matrix) != 1 || !ok {
		return fmt.Errorf("Claude matrix keys = %v, want only include", sortedKeys(matrix))
	}
	var rows []liveCadenceRow
	if err := include.Decode(&rows); err != nil {
		return fmt.Errorf("decode Claude matrix include: %w", err)
	}
	want := liveCadenceRow{Cadence: pullRequestCadence, Model: claudeCadenceModel, Effort: "max", Environment: claudeCadenceEnv}
	if len(rows) != 1 || rows[0] != want {
		return fmt.Errorf("Claude matrix include = %#v, want one explicit row %#v", rows, want)
	}
	if job.If != claudeCadenceIf {
		return fmt.Errorf("Claude job if = %q, want exclusive cadence condition %q", job.If, claudeCadenceIf)
	}
	codex, ok := parsed.Jobs["codex-live"]
	if !ok || codex.Environment.Name != "CI-E2E-CODEX" || codex.If != codexCadenceIf {
		return fmt.Errorf("Codex lane environment/if = %q/%q, want CI-E2E-CODEX/%q", codex.Environment.Name, codex.If, codexCadenceIf)
	}
	pi, ok := parsed.Jobs["pi-live"]
	if !ok || pi.Environment.Name != "CI-E2E-PI" || pi.If != piCadenceIf {
		return fmt.Errorf("Pi lane environment/if = %q/%q, want CI-E2E-PI/%q", pi.Environment.Name, pi.If, piCadenceIf)
	}
	if pi.Env["OPENAI_API_KEY"] != `${{ secrets.OPENAI_API_KEY }}` || pi.Env["SPACEDOCK_PI_LIVE_CHILD_MODEL"] != "openai/gpt-5.6-luna:max" {
		return fmt.Errorf("Pi lane key/model = %q/%q, want stored OpenAI key and Luna/max", pi.Env["OPENAI_API_KEY"], pi.Env["SPACEDOCK_PI_LIVE_CHILD_MODEL"])
	}
	return nil
}

func sortedKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

type liveClaim struct{ step, selector, claim string }

var liveClaims = []liveClaim{
	{"Run live Claude E2E", "TestLiveCommon", "claude-common-journeys"},
	{"Run live Claude substrate proofs", "TestLiveMergedTeamModeDispatch", "claude-merged-agent-dispatch"},
	{"Run live Claude substrate proofs", "TestLiveBareReachable", "claude-bare-dispatch"},
	{"Run live Claude substrate proofs", "TestLiveBreakGlassShimRecovery", "claude-break-glass-recovery"},
	{"Verify Codex resolver against installed plugin", "TestCodexResolveManifestAgainstInstalledHost", "codex-current-checkout-manifest-resolution"},
	{"Run live Codex shared scenarios", "TestLiveCommon", "codex-common-journeys"},
	{"Run live Pi common journeys", "TestLiveCommon", "pi-common-journeys"},
	{"Run live Pi front-door smoke", "TestLivePiFrontDoorSmoke", "pi-front-door-substrate"},
}

var offlineControls = []string{
	"TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle", "TestCleanupKeepMovingRootRetainsOnlyFailures",
	"TestCodexLiveRunnerExecArgvEnablesMultiAgentV2",
	"TestCodexLiveRunnerUsesSpacedockFrontDoorBeforeHostArgs", "TestPiIntercomPackageRootDefaultsBesideSubagents",
	"TestPiLiveEnvDropsForeignRuntimeMarkers", "TestPiLiveEnvScrubsAmbientPiSubagentMarkers",
	"TestPiLiveSmokePromptRequiresExactStageReportHeading", "TestShallowBootFixtureContainsOnlyHeldGate",
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
		"Pi cadence loses exclusivity": func(s string) string {
			return strings.Replace(s, piCadenceIf, `${{ github.event_name == 'workflow_dispatch' }}`, 1)
		},
		"Opus cadence starts Codex": func(s string) string {
			return strings.Replace(s, codexCadenceIf, `${{ github.event_name == 'pull_request' || inputs.live_cadence != 'pi' }}`, 1)
		},
		"Pi cadence starts Claude": func(s string) string {
			return strings.Replace(s, claudeCadenceIf, `${{ github.event_name == 'pull_request' || inputs.live_cadence != 'none' }}`, 1)
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
			if assertOneClaudeCadence(adversarial) == nil && assertNamedLiveEvidence(adversarial) == nil && assertOfflineControls(adversarial) == nil {
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
	for _, dead := range []string{"TestLivePty", "TestLivePiRecordedGateLifecycle", "TestLivePiSubagentEnsignSmoke", "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS", "pty-team-mode", "Install tmux", "inputs.effort"} {
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
		return fmt.Errorf("the 9 deterministic controls are not an exact untagged command")
	}
	return nil
}

var runSelector = regexp.MustCompile(`(?:^|\s)-run\s+['"]?\^?([A-Za-z0-9_|]+)`)

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
