package release

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseLiveModelsExcludesCadencePolicy(t *testing.T) {
	got := parseLiveModels("claude.future=fixture-model\nclaude.allowed.sonnet=future\ncodex.exec=fixture-codex\npi.oauth=fixture-pi:max\n")
	want := map[string]string{
		"claude.future": "fixture-model",
		"codex.exec":    "fixture-codex",
		"pi.oauth":      "fixture-pi:max",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseLiveModels() = %v, want %v", got, want)
	}
}

// TestClaudeLiveKeepsCadenceBoundaryAndSelectedModel pins the claude-live legs
// per event: a pull request runs the selected (Sonnet) model blocking plus the
// registered Haiku model continue-on-error; a dispatch runs only the selected
// model, blocking. The matrix model axis is expanded from its event branches
// and each leg is classified by the job's continue-on-error rule.
func TestClaudeLiveKeepsCadenceBoundaryAndSelectedModel(t *testing.T) {
	var doc struct {
		Jobs map[string]struct {
			Name            string
			Needs           any
			ContinueOnError string `yaml:"continue-on-error"`
			Outputs         map[string]string
			Strategy        struct {
				Matrix struct {
					Model   string
					Include []map[string]string
				}
			}
		}
	}
	if err := yaml.Unmarshal([]byte(readWorkflow(t, "runtime-live-e2e.yml")), &doc); err != nil {
		t.Fatalf("parse runtime-live-e2e.yml: %v", err)
	}
	job := doc.Jobs["claude-live"]
	if job.Needs != "offline" {
		t.Errorf("claude-live needs = %v; want offline so cadence refusal precedes approval", job.Needs)
	}
	if want := "claude-live (${{ matrix.cadence }}, ${{ matrix.model }}, ${{ matrix.effort }}, ${{ matrix.environment }})"; job.Name != want {
		t.Errorf("claude-live name = %q; want %q so the Sonnet check name is unchanged", job.Name, want)
	}
	// An include entry that names a model creates a leg on every event.
	if inc := job.Strategy.Matrix.Include; len(inc) != 1 || inc[0]["model"] != "" || inc[0]["effort"] != "max" {
		t.Errorf("claude-live include = %v; want one model-free entry carrying effort max", inc)
	}
	const rule = "${{ matrix.model != needs.offline.outputs.claude_selected }}"
	if job.ContinueOnError != rule {
		t.Errorf("claude-live continue-on-error = %q; want %q", job.ContinueOnError, rule)
	}
	branches := regexp.MustCompile(`^\$\{\{ fromJSON\(github\.event_name == 'pull_request' && format\('([^']*)', ([^)]*)\) \|\| format\('([^']*)', ([^)]*)\)\) \}\}$`).FindStringSubmatch(job.Strategy.Matrix.Model)
	if branches == nil {
		t.Fatalf("claude-live matrix model = %q; want a pull_request/dispatch fromJSON(format(...)) choice", job.Strategy.Matrix.Model)
	}
	// legs expands one branch into resolver output names, through the offline
	// job's outputs, marking each leg blocking or not by the continue-on-error rule.
	legs := func(format, args string) []string {
		for i, arg := range strings.Split(args, ", ") {
			output := doc.Jobs["offline"].Outputs[strings.TrimPrefix(arg, "needs.offline.outputs.")]
			resolved := strings.TrimSuffix(strings.TrimPrefix(output, "${{ steps.live_models.outputs."), " }}")
			format = strings.ReplaceAll(format, fmt.Sprintf("{%d}", i), resolved)
		}
		var models []string
		if err := json.Unmarshal([]byte(format), &models); err != nil {
			t.Fatalf("matrix model branch %q: %v", format, err)
		}
		for i, model := range models {
			if model == "claude_selected" {
				models[i] += " blocking"
			} else {
				models[i] += " continue-on-error"
			}
		}
		return models
	}
	if got, want := legs(branches[1], branches[2]), []string{"claude_selected blocking", "claude_haiku continue-on-error"}; !reflect.DeepEqual(got, want) {
		t.Errorf("pull_request claude-live legs = %v; want %v", got, want)
	}
	if got, want := legs(branches[3], branches[4]), []string{"claude_selected blocking"}; !reflect.DeepEqual(got, want) {
		t.Errorf("workflow_dispatch claude-live legs = %v; want %v", got, want)
	}
	if _, rest, ok := strings.Cut(liveModelsText, "\nclaude.allowed.sonnet="); !ok || strings.Contains(","+strings.SplitN(rest, "\n", 2)[0]+",", ",opus,") {
		t.Errorf("routine sonnet cadence policy is missing or permits opus: %q", strings.SplitN(rest, "\n", 2)[0])
	}
}

// TestLiveModelResolverEmitsPullRequestLegs runs the offline job's resolver
// step: a pull request resolves Sonnet as the selected model and emits the
// Haiku output the second leg reads, and it refuses a registry without Haiku.
func TestLiveModelResolverEmitsPullRequestLegs(t *testing.T) {
	script := extractStepRun(t, readWorkflow(t, "runtime-live-e2e.yml"), "Resolve live lane models")
	run := func(registry string) (string, error) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "internal", "release"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "internal", "release", "live_models.txt"), []byte(registry), 0o644); err != nil {
			t.Fatal(err)
		}
		outputs := filepath.Join(dir, "github-output")
		cmd := exec.Command("bash", "-c", script)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GITHUB_OUTPUT="+outputs, "LIVE_EVENT=pull_request", "LIVE_CADENCE=", "CLAUDE_MODEL_KEY=")
		if out, err := cmd.CombinedOutput(); err != nil {
			return string(out), err
		}
		got, err := os.ReadFile(outputs)
		return string(got), err
	}
	got, err := run(liveModelsText)
	if err != nil {
		t.Fatalf("resolver failed on the shipped registry: %v\n%s", err, got)
	}
	for _, want := range []string{"claude_selected=" + LiveModels["claude.sonnet"], "claude_haiku=" + LiveModels["claude.haiku"]} {
		if !strings.Contains("\n"+got, "\n"+want+"\n") {
			t.Errorf("pull_request resolver outputs missing %q:\n%s", want, got)
		}
	}
	withoutHaiku := "claude.sonnet=fixture-sonnet\nclaude.opus=fixture-opus\nclaude.allowed.sonnet=sonnet\nclaude.allowed.opus-pre-release=sonnet,opus\n"
	if got, err := run(withoutHaiku); err == nil || !strings.Contains(got, "claude.haiku") {
		t.Errorf("resolver accepted a pull request registry without claude.haiku: err=%v\n%s", err, got)
	}
}
