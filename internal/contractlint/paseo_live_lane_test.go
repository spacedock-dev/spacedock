// ABOUTME: Structural checks for the optional Paseo live E2E workflow and single-source model IDs.
package contractlint

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/release"
	"gopkg.in/yaml.v3"
)

type paseoStep struct {
	Run string `yaml:"run"`
}

type paseoWorkflowDoc struct {
	On          map[string]any    `yaml:"on"`
	Permissions map[string]string `yaml:"permissions"`
	Jobs        map[string]struct {
		Needs           any               `yaml:"needs"`
		ContinueOnError bool              `yaml:"continue-on-error"`
		Env             map[string]string `yaml:"env"`
		Environment     map[string]string `yaml:"environment"`
		Outputs         map[string]string `yaml:"outputs"`
		Steps           []paseoStep       `yaml:"steps"`
	} `yaml:"jobs"`
}

func TestPaseoLiveWorkflowSecurityAndModelSSOT(t *testing.T) {
	repo := repoRoot(t)
	wfBytes := mustRead(t, filepath.Join(repo, ".github", "workflows", "paseo-live-e2e.yml"))
	wfText, docsText := string(wfBytes), string(mustRead(t, filepath.Join(repo, "docs", "runtime-live-ci.md")))

	var doc paseoWorkflowDoc
	if err := yaml.Unmarshal(wfBytes, &doc); err != nil {
		t.Fatalf("parse paseo-live-e2e.yml: %v", err)
	}
	if len(doc.Permissions) != 1 || doc.Permissions["contents"] != "read" {
		t.Errorf("permissions = %v, want only contents: read", doc.Permissions)
	}
	hasOn := func(k string) bool { _, ok := doc.On[k]; return ok }
	if !hasOn("pull_request") || hasOn("pull_request_target") || !hasOn("workflow_dispatch") {
		t.Errorf("triggers mismatch: pull_request=%v pull_request_target=%v workflow_dispatch=%v", hasOn("pull_request"), hasOn("pull_request_target"), hasOn("workflow_dispatch"))
	}

	offline, ok := doc.Jobs["offline"]
	if !ok || offline.Environment["name"] != "" || offline.Outputs["codex_exec"] == "" || !strings.Contains(joinPaseoRuns(offline.Steps), "internal/release/live_models.txt") {
		t.Fatalf("offline job invalid: present=%v env=%v outputs=%v", ok, offline.Environment, offline.Outputs)
	}

	live, ok := doc.Jobs["paseo-live"]
	if !ok || live.Needs != "offline" || !live.ContinueOnError || live.Environment["name"] != "CI-E2E-CODEX" {
		t.Fatalf("paseo-live config mismatch: present=%v needs=%v continue-on-error=%v env=%v", ok, live.Needs, live.ContinueOnError, live.Environment)
	}
	if live.Env["CODEX_AUTH_JSON"] != "${{ secrets.CODEX_AUTH_JSON }}" || live.Env["OPENAI_API_KEY"] != "${{ secrets.OPENAI_API_KEY }}" {
		t.Errorf("paseo-live.env secret bindings = %v", live.Env)
	}
	liveRun := joinPaseoRuns(live.Steps)
	for _, want := range []string{"CODEX_AUTH_JSON", "OPENAI_API_KEY", "SPACEDOCK_BIN=", "SPACEDOCK_REPO_ROOT="} {
		if !strings.Contains(liveRun, want) {
			t.Errorf("paseo-live steps missing %q", want)
		}
	}
	var foundSmokeStep bool
	for _, s := range live.Steps {
		if strings.Contains(s.Run, "-run '^TestLivePaseoDispatchSmoke$'") {
			foundSmokeStep = true
			if !strings.Contains(s.Run, "SPACEDOCK_PASEO_LIVE_REQUIRED=1") || !strings.Contains(s.Run, "needs.offline.outputs.codex_exec") {
				t.Errorf("paseo-live smoke step missing required env or model output:\n%s", s.Run)
			}
		}
	}
	if !foundSmokeStep {
		t.Error("paseo-live missing TestLivePaseoDispatchSmoke run step")
	}

	if codexModel := release.LiveModels["codex.exec"]; codexModel == "" || strings.Contains(wfText, codexModel) || strings.Contains(docsText, codexModel) {
		t.Errorf("invalid or hardcoded codex.exec model %q in workflow/docs", codexModel)
	}
}

func joinPaseoRuns(steps []paseoStep) string {
	var b strings.Builder
	for _, s := range steps {
		b.WriteString(s.Run)
		b.WriteByte('\n')
	}
	return b.String()
}
