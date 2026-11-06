// ABOUTME: dispatch build --model/--reasoning-effort — explicit Codex worker settings pass
// ABOUTME: through opaquely, and any other host or mode refuses them before mutation.
package dispatch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func codexOverrideFixture(t *testing.T) []string {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeNonWorktreeStages())
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", ""))
	gitInit(t, root)
	return []string{"build", "--workflow-dir", root, "--entity-path", entityPath, "--stage", "implementation", "--checklist-file", "-"}
}

func TestBuildCodexModelAndEffortPassthrough(t *testing.T) {
	base := codexOverrideFixture(t)
	for _, tc := range []struct {
		name       string
		extra      []string
		wantModel  string
		wantEffort string
	}{
		{"space form", []string{"--model", "gpt-5.5", "--reasoning-effort", "high"}, "gpt-5.5", "high"},
		{"equals form", []string{"--model=gpt-5.5-mini", "--reasoning-effort=low"}, "gpt-5.5-mini", "low"},
		{"model only", []string{"--model", "gpt-5.5"}, "gpt-5.5", ""},
		{"effort only", []string{"--reasoning-effort", "medium"}, "", "medium"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := append(append([]string{}, base...), "--host", "codex")
			res := runNative("- a", append(args, tc.extra...)...)
			if res.exit != 0 {
				t.Fatalf("exit=%d stderr=%q", res.exit, res.stderr)
			}
			var raw map[string]json.RawMessage
			if err := json.Unmarshal([]byte(res.stdout), &raw); err != nil {
				t.Fatalf("stdout is not JSON: %v", err)
			}
			if tc.wantModel == "" {
				if string(raw["model"]) != "null" {
					t.Errorf("model = %s, want null", raw["model"])
				}
			} else if string(raw["model"]) != `"`+tc.wantModel+`"` {
				t.Errorf("model = %s, want %q", raw["model"], tc.wantModel)
			}
			effort, has := raw["reasoning_effort"]
			if tc.wantEffort == "" && has {
				t.Errorf("reasoning_effort must be absent, got %s", effort)
			}
			if tc.wantEffort != "" && string(effort) != `"`+tc.wantEffort+`"` {
				t.Errorf("reasoning_effort = %s, want %q", effort, tc.wantEffort)
			}
			if tc.wantModel != "" && !strings.Contains(res.stderr, "[build] explicit model="+tc.wantModel) {
				t.Errorf("stderr must note the model: %q", res.stderr)
			}
			if tc.wantEffort != "" && !strings.Contains(res.stderr, tc.wantEffort) {
				t.Errorf("stderr must note the effort: %q", res.stderr)
			}
			if _, err := CodexMultiAgentV2SpawnInput([]byte(res.stdout)); err != nil {
				t.Errorf("adapter rejects emitted envelope: %v", err)
			}
		})
	}
}

func TestBuildCodexWithoutOverridesKeepsEnvelope(t *testing.T) {
	base := codexOverrideFixture(t)
	res := runNative("- a", append(base, "--host", "codex")...)
	if res.exit != 0 {
		t.Fatalf("exit=%d stderr=%q", res.exit, res.stderr)
	}
	if !strings.Contains(res.stdout, `"model": null`) || strings.Contains(res.stdout, "reasoning_effort") {
		t.Errorf("unflagged codex envelope changed:\n%s", res.stdout)
	}
}

func TestBuildModelAndEffortRefusals(t *testing.T) {
	base := codexOverrideFixture(t)
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"model on claude", []string{"--host", "claude", "--model", "gpt-5.5"}, "--model"},
		{"effort on claude", []string{"--host", "claude", "--reasoning-effort", "high"}, "--reasoning-effort"},
		{"model on pi", []string{"--host", "pi", "--model", "gpt-5.5"}, "--model"},
		{"empty model", []string{"--host", "codex", "--model", ""}, "--model"},
		{"empty model equals", []string{"--host", "codex", "--model="}, "--model"},
		{"empty effort equals", []string{"--host", "codex", "--reasoning-effort="}, "--reasoning-effort"},
		{"model with advance", []string{"--host", "codex", "--model", "gpt-5.5", "--advance"}, "--advance"},
		{"effort with advance", []string{"--host", "codex", "--reasoning-effort", "high", "--advance"}, "--advance"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res := runNative("- a", append(append([]string{}, base...), tc.args...)...)
			if res.exit == 0 {
				t.Fatalf("accepted; stdout=%q", res.stdout)
			}
			if res.stdout != "" {
				t.Errorf("refusal emitted stdout: %q", res.stdout)
			}
			if !strings.Contains(res.stderr, tc.want) {
				t.Errorf("stderr %q does not name %q", res.stderr, tc.want)
			}
		})
	}
}

func TestBuildModelRefusedBeforeStampMutation(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	_, workflowDir, statePath, entityPath := stampFixture(t, "implementation", "implementation", false)
	before, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatal(err)
	}
	headBefore := strings.TrimSpace(gitOutput(t, statePath, "rev-parse", "HEAD"))

	args := append(buildStampArgs(workflowDir, entityPath, "implementation", "-"), "--host", "claude", "--model", "gpt-5.5")
	res := runNative("- a\n", args...)
	if res.exit == 0 {
		t.Fatalf("accepted; stdout=%q", res.stdout)
	}
	after, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Error("refused build mutated entity frontmatter")
	}
	if headAfter := strings.TrimSpace(gitOutput(t, statePath, "rev-parse", "HEAD")); headAfter != headBefore {
		t.Error("refused build advanced the state checkout HEAD")
	}
}
