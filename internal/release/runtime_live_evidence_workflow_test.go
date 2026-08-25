package release

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// liveCadenceWorkflow reads only the job-level env of the live jobs — the
// surface the retired/unavailable-secret bans below inspect.
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
		Env map[string]string `yaml:"env"`
	} `yaml:"jobs"`
}

// assertLiveSecretsBansHold enforces three retired/unavailable-secret
// prohibitions in the live workflow's job-level env: codex-live must not
// carry CODEX_HOME (runner.temp is unavailable in job-level env — a
// GHA-semantics fact, so isolated CODEX_HOME must be set in a step instead),
// and pi-live must not carry the retired PI_OPENAI_CODEX_AUTH_JSON secret or
// the retired SPACEDOCK_PI_LIVE_CHILD_MODEL override.
func assertLiveSecretsBansHold(workflow string) error {
	var parsed liveCadenceWorkflow
	if err := yaml.Unmarshal([]byte(workflow), &parsed); err != nil {
		return err
	}
	codex, ok := parsed.Jobs["codex-live"]
	if !ok {
		return fmt.Errorf("workflow lacks codex-live job")
	}
	if _, present := codex.Env["CODEX_HOME"]; present {
		return fmt.Errorf("Codex HOME must be initialized in a step, not with the unavailable job-level runner.temp context")
	}
	pi, ok := parsed.Jobs["pi-live"]
	if !ok {
		return fmt.Errorf("workflow lacks pi-live job")
	}
	if _, present := pi.Env["PI_OPENAI_CODEX_AUTH_JSON"]; present {
		return fmt.Errorf("Pi lane must not use the dedicated PI_OPENAI_CODEX_AUTH_JSON secret")
	}
	if _, present := pi.Env["SPACEDOCK_PI_LIVE_CHILD_MODEL"]; present {
		return fmt.Errorf("Pi lane must derive its provider model from auth mode")
	}
	parallel, ok := parsed.On.WorkflowDispatch.Inputs["live_parallel"]
	if !ok || parallel.Default != "4" {
		return fmt.Errorf("workflow_dispatch live_parallel input = %#v, want a string input defaulting to \"4\"", parallel)
	}
	if pi.Env["SPACEDOCK_PI_LIVE_PARALLEL"] != `${{ inputs.live_parallel }}` {
		return fmt.Errorf("Pi lane SPACEDOCK_PI_LIVE_PARALLEL = %q, want the live_parallel input", pi.Env["SPACEDOCK_PI_LIVE_PARALLEL"])
	}
	return nil
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

func TestRuntimeLiveWorkflowSecretBansHoldOnRealFile(t *testing.T) {
	if err := assertLiveSecretsBansHold(readWorkflow(t, "runtime-live-e2e.yml")); err != nil {
		t.Fatal(err)
	}
}

func TestRuntimeLiveWorkflowNamedEvidenceMutationControls(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	mutations := map[string]func(string) string{
		"legacy PTY flag": func(s string) string {
			return strings.Replace(s, "DISABLE_AUTOUPDATER: \"1\"", "DISABLE_AUTOUPDATER: \"1\"\n      CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS: \"1\"", 1)
		},
		"Pi parallel override lost": func(s string) string {
			return strings.Replace(s, `-parallel "$PI_LIVE_PARALLEL" `, "", 1)
		},
		"offline live tag": func(s string) string {
			return strings.Replace(s, "go test ./internal/ensigncycle -count=1 -run", "go test -tags live ./internal/ensigncycle -count=1 -run", 1)
		},
		"job-level CODEX_HOME reintroduced": func(s string) string {
			return strings.Replace(s, "SPACEDOCK_CODEX_LIVE_REQUIRED: \"1\"", "SPACEDOCK_CODEX_LIVE_REQUIRED: \"1\"\n      CODEX_HOME: ${{ runner.temp }}/codex-home", 1)
		},
		"retired PI_OPENAI_CODEX_AUTH_JSON reintroduced": func(s string) string {
			return strings.Replace(s, "SPACEDOCK_PI_LIVE_REQUIRED: \"1\"", "SPACEDOCK_PI_LIVE_REQUIRED: \"1\"\n      PI_OPENAI_CODEX_AUTH_JSON: ${{ secrets.CODEX_AUTH_JSON }}", 1)
		},
		"retired SPACEDOCK_PI_LIVE_CHILD_MODEL reintroduced": func(s string) string {
			return strings.Replace(s, "SPACEDOCK_PI_LIVE_REQUIRED: \"1\"", "SPACEDOCK_PI_LIVE_REQUIRED: \"1\"\n      SPACEDOCK_PI_LIVE_CHILD_MODEL: claude-opus-4-8", 1)
		},
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			adversarial := mutate(live)
			if adversarial == live {
				t.Fatalf("mutation %q did not change the fixture workflow", name)
			}
			if assertNamedLiveEvidence(adversarial) == nil && assertOfflineControls(adversarial) == nil && assertLiveSecretsBansHold(adversarial) == nil {
				t.Fatal("mutation escaped the workflow guards")
			}
		})
	}
}

// assertNamedLiveEvidence bans a step from selecting live tests it does not
// own: every step that passes a `-run` selector must be a registered claim
// owner (from liveClaims, or the deterministic offline-controls step), and
// the workflow must not retain any of a fixed list of dead surfaces.
func assertNamedLiveEvidence(workflow string) error {
	owned := map[string]bool{"Run deterministic live-harness controls offline": true}
	claimOwners := map[string]string{}
	for _, item := range liveClaims {
		if claimOwners[item.claim] != "" {
			return fmt.Errorf("claim %q has two owners", item.claim)
		}
		claimOwners[item.claim] = item.selector
		owned[item.step] = true
	}
	for _, step := range parseWorkflowSteps(workflow) {
		got := selectedTests(step.run)
		if step.name == "Run live Pi common journeys" && !strings.Contains(step.run, `-parallel "$PI_LIVE_PARALLEL"`) {
			return fmt.Errorf("Pi common journeys step lost the explicit -parallel override")
		}
		if len(got) > 0 && !owned[step.name] {
			return fmt.Errorf("step %q selects unowned tests %v", step.name, got)
		}
	}
	for _, dead := range []string{"TestLivePty", "TestLivePiRecordedGateLifecycle", "TestLivePiSubagentEnsignSmoke", "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS", "pty-team-mode", "Install tmux", "inputs.effort"} {
		if activeYAMLText(workflow, dead) {
			return fmt.Errorf("workflow retains dead surface %q", dead)
		}
	}
	return nil
}

// assertOfflineControls bans the deterministic live-harness controls step
// from running under the `live` build tag — it must stay executable offline.
func assertOfflineControls(workflow string) error {
	step, ok := stepNamed(parseWorkflowSteps(workflow), "Run deterministic live-harness controls offline")
	if !ok {
		return fmt.Errorf("workflow lacks the deterministic live-harness controls step")
	}
	if strings.Contains(step.run, "-tags") {
		return fmt.Errorf("deterministic live-harness controls step must run untagged (offline), not gated behind -tags live")
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
