// ABOUTME: Structural checks for the optional Paseo live E2E workflow and single-source model IDs.
package contractlint

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/release"
	"gopkg.in/yaml.v3"
)

type paseoWorkflowDoc struct {
	On          map[string]any    `yaml:"on"`
	Permissions map[string]string `yaml:"permissions"`
	Jobs        map[string]struct {
		Needs           any  `yaml:"needs"`
		ContinueOnError bool `yaml:"continue-on-error"`
		Environment     struct {
			Name string `yaml:"name"`
		} `yaml:"environment"`
		Outputs map[string]string `yaml:"outputs"`
		Steps   []struct {
			Run string `yaml:"run"`
		} `yaml:"steps"`
	} `yaml:"jobs"`
}

func TestPaseoLiveWorkflowSecurityAndModelSSOT(t *testing.T) {
	repo := repoRoot(t)
	wfBytes := mustRead(t, filepath.Join(repo, ".github", "workflows", "paseo-live-e2e.yml"))
	wfText := string(wfBytes)
	docsText := string(mustRead(t, filepath.Join(repo, "docs", "runtime-live-ci.md")))

	var doc paseoWorkflowDoc
	if err := yaml.Unmarshal(wfBytes, &doc); err != nil {
		t.Fatalf("parse paseo-live-e2e.yml: %v", err)
	}
	if doc.Permissions["contents"] != "read" {
		t.Errorf("permissions.contents = %q, want read", doc.Permissions["contents"])
	}
	if _, ok := doc.On["pull_request"]; !ok {
		t.Error("workflow must trigger on pull_request")
	}
	if _, ok := doc.On["pull_request_target"]; ok {
		t.Error("workflow must not trigger on pull_request_target")
	}
	if _, ok := doc.On["workflow_dispatch"]; !ok {
		t.Error("workflow must trigger on workflow_dispatch")
	}

	offline, ok := doc.Jobs["offline"]
	if !ok || offline.Environment.Name != "" || offline.Outputs["codex_exec"] == "" {
		t.Fatalf("offline job invalid: present=%v env=%q outputs=%v", ok, offline.Environment.Name, offline.Outputs)
	}
	if !strings.Contains(joinPaseoRuns(offline.Steps), "internal/release/live_models.txt") {
		t.Error("offline job must read internal/release/live_models.txt")
	}

	live, ok := doc.Jobs["paseo-live"]
	if !ok {
		t.Fatal("missing paseo-live job")
	}
	if live.Needs != "offline" || !live.ContinueOnError || live.Environment.Name != "CI-E2E-CODEX" {
		t.Errorf("paseo-live config mismatch: needs=%v continue-on-error=%v env=%q", live.Needs, live.ContinueOnError, live.Environment.Name)
	}
	liveRun := joinPaseoRuns(live.Steps)
	for _, want := range []string{
		"CODEX_AUTH_JSON", "OPENAI_API_KEY", "SPACEDOCK_BIN=", "SPACEDOCK_REPO_ROOT=",
		"needs.offline.outputs.codex_exec", "-run '^TestLivePaseoDispatchSmoke$'",
	} {
		if !strings.Contains(liveRun, want) {
			t.Errorf("paseo-live steps missing %q", want)
		}
	}

	codexModel := release.LiveModels["codex.exec"]
	if codexModel == "" {
		t.Fatal("release.LiveModels[codex.exec] is empty")
	}
	if strings.Contains(wfText, codexModel) || strings.Contains(docsText, codexModel) {
		t.Errorf("hardcoded model literal %q found in workflow or docs", codexModel)
	}
}

func joinPaseoRuns(steps []struct {
	Run string `yaml:"run"`
}) string {
	var b strings.Builder
	for _, s := range steps {
		b.WriteString(s.Run)
		b.WriteByte('\n')
	}
	return b.String()
}
