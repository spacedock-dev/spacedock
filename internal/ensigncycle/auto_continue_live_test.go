//go:build live

// ABOUTME: AC-5 live regression — a real FO, given an implementation-ready dev
// ABOUTME: entity, must advance to validation and dispatch a fresh validator, not stop.
package ensigncycle

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/livescenario"
)

// TestLiveAutoContinueAfterImplementation is AC-5's live half: a real FO is
// pointed at a dev-shaped workflow whose one entity is parked at an
// implementation-ready state (status: implementation with a filed implementation
// stage report) under the NEUTRAL `Use $spacedock:first-officer` runbook — no
// "drive to done" coaching. The scenario grades on the DURABLE outcome via the
// shared, state-oriented assertAutoContinue: the FO must advance the entity past
// implementation AND leave a `## Stage Report: validation` behind — the durable
// footprint of a fresh validator it dispatched. A run that stops after the
// implementation report, leaving the entity at status: implementation, REDS the
// grade. Running this (`go test -tags live -run TestLiveAutoContinueAfterImplementation`)
// against a real credential produces the session artifact AC-5's `Verified by:
// live …` citation requires.
//
// Non-interactive `claude -p` puts the FO in single-entity mode, where it drives
// the parked entity all the way to terminal `done` and ARCHIVES it to
// `_archive/auto-continue-task.md`. That over-runs the minimum (it more than
// proves the FO did not stop after implementation), so the grade reads the entity
// from wherever it lands — the original path or the archive — via the captured
// workflow dir. The primitive's tolerant post-run read hands the Assert an empty
// after-body when the original path was archived; the Assert then resolves the
// real end-state. It reuses claudeRunnerAdapter + errGraded from
// livescenario_adapter_live_test.go and the assertAutoContinue grade shared with
// the offline negative case.
func TestLiveAutoContinueAfterImplementation(t *testing.T) {
	runner := newClaudeLiveRunner(t)
	// Implementation completion → validator dispatch → (single-entity) gate
	// auto-resolve → merge/terminalize runs TWO full agent runs serially (the FO
	// and the fresh validator), so the budget is generous.
	adapter := claudeRunnerAdapter{t: t, runner: runner, timeout: 15 * time.Minute}

	var workflowDir string
	sc := livescenario.Scenario{
		Name:    "auto-continue-after-implementation",
		Runbook: autoContinuePrompt(),
		Setup: func(dir string) (string, error) {
			// Capture the staged workflow dir so the Assert can find the entity even
			// after the FO archives it. Stage WITHOUT git-init; the adapter git-inits
			// once the primitive has captured the pre-run state (matching the live order).
			workflowDir = dir
			return writeAutoContinueWorkflowNoGit(dir)
		},
		Assert: func(before, after livescenario.EntityState, observed string) error {
			afterBody := resolveAutoContinueEndState(workflowDir, after.Body)
			if err := assertAutoContinue(before.Body, afterBody, observed); err != nil {
				return errGraded(err.Error())
			}
			return nil
		},
	}

	if err := livescenario.Run(context.Background(), t.TempDir(), sc, adapter); err != nil {
		t.Fatalf("live auto-continue scenario graded FAIL: %v", err)
	}
}

// resolveAutoContinueEndState returns the entity's durable end-state body. The FO
// may archive a terminalized entity, moving it out of its original path; in that
// case the primitive's after-body is empty (the original path is gone) and the
// real end-state lives at `_archive/auto-continue-task.md`. This reads that
// archived copy when present, otherwise falls back to the primitive's after-body
// (the entity stayed put — e.g. held at the validation gate). It NEVER fabricates
// state: a genuinely absent entity yields an empty body and the state-oriented
// assertAutoContinue reds.
func resolveAutoContinueEndState(workflowDir, afterBody string) string {
	if afterBody != "" {
		return afterBody
	}
	archived := filepath.Join(workflowDir, "_archive", "auto-continue-task.md")
	if data, err := os.ReadFile(archived); err == nil {
		return string(data)
	}
	return afterBody
}
