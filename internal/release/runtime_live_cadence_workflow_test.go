package release

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type liveCadenceLeg struct {
	Cadence     string `yaml:"cadence"`
	Model       string `yaml:"model"`
	Effort      string `yaml:"effort"`
	Environment string `yaml:"environment"`
}

func TestRuntimeLiveClaudeShimSetsMaximumEffort(t *testing.T) {
	workflow := readWorkflow(t, "runtime-live-e2e.yml")
	const start = `cat > "$shim_dir/claude" <<'SH'`
	startAt := strings.Index(workflow, start)
	bodyStart := startAt + len(start) + 1
	endAt := strings.Index(workflow[bodyStart:], "\n          SH\n")
	lines := strings.Split(workflow[bodyStart:bodyStart+endAt], "\n")
	for i := range lines {
		lines[i] = strings.TrimPrefix(lines[i], "          ")
	}

	root := t.TempDir()
	realClaude := filepath.Join(root, "real-claude")
	shimPath := filepath.Join(root, "claude")
	files := map[string]string{realClaude: "#!/usr/bin/env bash\nprintf '%s\\n' \"$*\"\n", shimPath: strings.Join(lines, "\n") + "\n"}
	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	var got []byte
	for _, args := range [][]string{{"--version"}, {"--model", "claude-sonnet-5", "--help"}} {
		cmd := exec.Command(shimPath, args...)
		cmd.Env = append(os.Environ(), "SPACEDOCK_CLAUDE_REAL_BIN="+realClaude)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Claude shim %v failed: %v\n%s", args, err, out)
		} else {
			got = append(got, out...)
		}
	}
	if want := "--effort max --version\n--effort max --model claude-sonnet-5 --help\n"; string(got) != want {
		t.Fatalf("Claude shim argv = %q, want %q", got, want)
	}
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
		Environment struct {
			Name string `yaml:"name"`
		} `yaml:"environment"`
		Strategy struct {
			Matrix struct {
				Cadence []string         `yaml:"cadence"`
				Model   []string         `yaml:"model"`
				Include []liveCadenceLeg `yaml:"include"`
				Exclude []liveCadenceLeg `yaml:"exclude"`
			} `yaml:"matrix"`
		} `yaml:"strategy"`
	} `yaml:"jobs"`
}

func TestRuntimeLiveWorkflowExpandsApprovedCadences(t *testing.T) {
	want := map[string][]liveCadenceLeg{
		"pull-request":     {{Cadence: "pull-request", Model: "claude-sonnet-5", Effort: "max", Environment: "CI-E2E"}},
		"sonnet":           {{Cadence: "sonnet", Model: "claude-sonnet-5", Effort: "max", Environment: "CI-E2E"}},
		"opus-pre-release": {{Cadence: "opus-pre-release", Model: "claude-opus-4-8", Effort: "max", Environment: "CI-E2E-OPUS"}},
	}
	workflow := readWorkflow(t, "runtime-live-e2e.yml")
	for cadence, expected := range want {
		got, approvals, err := expandClaudeCadence(workflow, cadence)
		if err != nil || fmt.Sprint(got) != fmt.Sprint(expected) || approvals != 2 {
			t.Fatalf("%s: legs=%#v approvals=%d err=%v, want %#v and two approvals", cadence, got, approvals, err, expected)
		}
	}
	mutations := [][2]string{
		{"cadence: pull-request\n            model: claude-opus-4-8", "cadence: never\n            model: claude-opus-4-8"},
		{"model: [claude-sonnet-5, claude-opus-4-8]", "model: [claude-sonnet-5, claude-sonnet-5, claude-opus-4-8]"},
		{"claude-sonnet-5", "sonnet"}, {"effort: max", "effort: high"},
	}
	for _, mutation := range mutations {
		legs, approvals, err := expandClaudeCadence(strings.ReplaceAll(workflow, mutation[0], mutation[1]), "pull-request")
		want := []liveCadenceLeg{{Cadence: "pull-request", Model: "claude-sonnet-5", Effort: "max", Environment: "CI-E2E"}}
		if err == nil && fmt.Sprint(legs) == fmt.Sprint(want) && approvals == 2 {
			t.Fatalf("mutation %q escaped the pull-request cadence guard", mutation[0])
		}
	}
}

func expandClaudeCadence(workflow, cadence string) ([]liveCadenceLeg, int, error) {
	var parsed liveCadenceWorkflow
	if err := yaml.Unmarshal([]byte(workflow), &parsed); err != nil {
		return nil, 0, err
	}
	input, ok := parsed.On.WorkflowDispatch.Inputs["live_cadence"]
	if !ok || input.Default != "sonnet" || fmt.Sprint(input.Options) != fmt.Sprint([]string{"sonnet", "opus-pre-release"}) {
		return nil, 0, fmt.Errorf("workflow_dispatch live_cadence choice is not the approved two-value surface")
	}
	job, ok := parsed.Jobs["claude-live"]
	if !ok || fmt.Sprint(job.Strategy.Matrix.Cadence) != fmt.Sprint([]string{"${{ github.event_name == 'pull_request' && 'pull-request' || inputs.live_cadence }}"}) {
		return nil, 0, fmt.Errorf("Claude matrix does not normalize the event to one cadence")
	}
	var legs []liveCadenceLeg
	for _, model := range job.Strategy.Matrix.Model {
		leg := liveCadenceLeg{Cadence: cadence, Model: model}
		for _, properties := range job.Strategy.Matrix.Include {
			if properties.Model == model {
				leg.Effort = properties.Effort
				leg.Environment = properties.Environment
			}
		}
		excluded := false
		for _, exclusion := range job.Strategy.Matrix.Exclude {
			if exclusion.Cadence == cadence && exclusion.Model == model {
				excluded = true
			}
		}
		if !excluded {
			legs = append(legs, leg)
		}
	}
	approvals := len(legs)
	if codex, ok := parsed.Jobs["codex-live"]; ok && codex.Environment.Name == "CI-E2E-CODEX" {
		approvals++
	}
	if _, ok := parsed.Jobs["pi-live"]; ok {
		approvals++
	}
	return legs, approvals, nil
}
