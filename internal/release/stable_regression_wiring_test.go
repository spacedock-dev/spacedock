// ABOUTME: Structural guards for the release steps this task adds and moves.
// ABOUTME: Each step must sit in the job that gives it its blocking power.
package release

import (
	"strings"
	"testing"
)

// The names of the release.yml steps that this task adds, moves, or renames.
// The channel-agreement parsers and the edge-advance gate guard key on these, so
// a rename in the workflow reds a test instead of going unnoticed.
const (
	stableRegressionGateStepName = "Gate the cut on the tag not regressing the stable channel"
	stableAdvanceStepName        = "Advance the stable channel ref to the tagged commit"
	mainStampStepName            = "Stamp main to the next edge prerelease version"
)

// stepInJob reports whether the named job owns the named step.
func stepInJob(workflow, jobName, stepName string) bool {
	for _, job := range parseWorkflowJobs(workflow) {
		if job.name != jobName {
			continue
		}
		for _, step := range job.steps {
			if step.name == stepName {
				return true
			}
		}
	}
	return false
}

// TestReleaseStepsSitInTheirOwningJobs binds each step to the job that gives it
// its power. The regression gate must sit in `e2e-gate`, which goreleaser needs,
// because a gate inside goreleaser starts the job that it must stop. The `main`
// stamp must sit in `edge-advance`, which owns the latest-line decision and cuts
// the matching pre0 binary. The stable-ref advance must stay in `goreleaser`,
// which publishes the release the ref points at.
//
// Per the README's release-machinery proof posture, this structural check and the
// condition-equality check in edge_advance_wiring_test.go are the proof for the
// YAML wiring. The in-situ shell behavior is observed at the next real cut.
func TestReleaseStepsSitInTheirOwningJobs(t *testing.T) {
	workflow := readWorkflow(t, "release.yml")
	for _, want := range []struct{ job, step string }{
		{"e2e-gate", stableRegressionGateStepName},
		{"goreleaser", stableAdvanceStepName},
		{"edge-advance", mainStampStepName},
	} {
		if !stepInJob(workflow, want.job, want.step) {
			t.Errorf("release.yml job %q has no step named %q", want.job, want.step)
		}
	}

	// Adversarial twin: rename the gate step away and the guard must red. Without
	// this, a guard that silently found nothing can still pass.
	renamed := strings.Replace(workflow,
		"      - name: "+stableRegressionGateStepName+"\n",
		"      - name: Some unrelated step\n", 1)
	if renamed == workflow {
		t.Fatal("fixture workflow missing the gate step name to rename")
	}
	if stepInJob(renamed, "e2e-gate", stableRegressionGateStepName) {
		t.Fatal("presence guard still found the gate step after the rename")
	}
}
