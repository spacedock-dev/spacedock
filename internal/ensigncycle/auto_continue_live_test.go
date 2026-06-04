//go:build live

// ABOUTME: AC-5 live regression — a real FO, given an implementation-ready dev
// ABOUTME: entity, must advance to validation and dispatch a fresh validator, not stop.
package ensigncycle

import (
	"context"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/livescenario"
)

// TestLiveAutoContinueAfterImplementation is AC-5's live half: a real FO is
// pointed at a dev-shaped workflow whose one entity is parked at an
// implementation-ready state (status: implementation with a filed implementation
// stage report) under the NEUTRAL `Use $spacedock:first-officer` runbook — no
// "drive to done" coaching. The scenario grades on the DURABLE outcome: the FO
// must advance the entity past implementation (to validation) AND leave a
// `## Stage Report: validation` behind — the durable footprint of a fresh
// validator it dispatched. A run that stops after the implementation report,
// leaving the entity at status: implementation, REDS the grade. Running this
// (`go test -tags live -run TestLiveAutoContinueAfterImplementation`) against a
// real credential produces the session/ci-run artifact AC-5's `Verified by: live …`
// citation requires.
//
// It reuses claudeRunnerAdapter + errGraded from livescenario_adapter_live_test.go
// (same package, same build tag) and the assertAutoContinue grade shared with the
// offline negative case — so the live half and the no-model negative half grade on
// the identical durable-outcome assertion.
func TestLiveAutoContinueAfterImplementation(t *testing.T) {
	runner := newClaudeLiveRunner(t)
	// Implementation completion → validator dispatch + gate presentation is a
	// multi-step journey (verify report, advance, dispatch a fresh validator into a
	// new worktree, the validator runs, present the gate); size the timeout like the
	// rejection-flow scenario, which is the structural twin (mid-lifecycle dispatch).
	adapter := claudeRunnerAdapter{t: t, runner: runner, timeout: 4 * time.Minute}

	sc := livescenario.Scenario{
		Name:    "auto-continue-after-implementation",
		Runbook: autoContinuePrompt(),
		Setup: func(dir string) (string, error) {
			// Stage the dev-shaped fixture WITHOUT git-init; the adapter git-inits once
			// the primitive has captured the pre-run state (matching the live order).
			return writeAutoContinueWorkflowNoGit(dir)
		},
		Assert: func(before, after livescenario.EntityState, observed string) error {
			if err := assertAutoContinue(before.Body, after.Body, observed); err != nil {
				return errGraded(err.Error())
			}
			return nil
		},
	}

	if err := livescenario.Run(context.Background(), t.TempDir(), sc, adapter); err != nil {
		t.Fatalf("live auto-continue scenario graded FAIL: %v", err)
	}
}
