package release

import (
	"testing"

	"gopkg.in/yaml.v3"
)

// TestRuntimeLiveWorkflowJourneyDeltaJobHasPRCommentPermission is AC-3's
// workflow-shape guard: the job that posts the per-PR journey-cost delta
// comment must declare its own `permissions: pull-requests: write` (the
// workflow-level default is `contents: read` only, which cannot post a PR
// comment), must run only on `pull_request` (a workflow_dispatch run has no PR
// to comment on), and must actually invoke the journey-delta CLI.
func TestRuntimeLiveWorkflowJourneyDeltaJobHasPRCommentPermission(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")

	var doc struct {
		Jobs map[string]struct {
			If          string            `yaml:"if"`
			Permissions map[string]string `yaml:"permissions"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal([]byte(live), &doc); err != nil {
		t.Fatalf("parse runtime-live-e2e.yml: %v", err)
	}

	job, ok := doc.Jobs["journey-delta-comment"]
	if !ok {
		t.Fatal("runtime-live-e2e.yml has no journey-delta-comment job")
	}
	if job.Permissions["pull-requests"] != "write" {
		t.Fatalf("journey-delta-comment job permissions = %+v, want pull-requests: write", job.Permissions)
	}
	if job.If != "github.event_name == 'pull_request'" {
		t.Fatalf("journey-delta-comment job if = %q, want it scoped to pull_request events", job.If)
	}
	if !workflowHasExecutableCommandContaining(live, "go run ./cmd/spacedock-release journey-delta") {
		t.Fatal("runtime-live-e2e.yml journey-delta-comment job does not invoke the journey-delta CLI")
	}
}
