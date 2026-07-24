//go:build live

// ABOUTME: Live proof that the promoted livescenario primitive runs against the
// ABOUTME: REAL Claude launch adapter — AC-2's runnable-by-a-real-agent half.
package ensigncycle

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/livescenario"
)

// claudeRunnerAdapter wraps the existing package-private claudeLiveRunner (the
// launch + observed-extract adapter 8y built) as a livescenario.Runner. This is
// the seam AC-2 names: the runner adapter STAYS in ensigncycle; what changes is
// that a scenario is now authored against the importable livescenario primitive
// rather than buried in this package's test files. The adapter copies the staged
// dir's contents into the runner's own workflow root and forwards the launch.
// Liveness is the claudeLiveRunner's own per-stage no-progress quiet budget (the
// shared streamWatcher) — no per-call basket timeout (those are banned).
type claudeRunnerAdapter struct {
	t      *testing.T
	runner claudeLiveRunner
}

func (a claudeRunnerAdapter) Launch(ctx context.Context, dir, entityPath, runbook string) (string, error) {
	scenario := sharedRuntimeScenario{name: "livescenario-primitive"}
	result := a.runner.run(a.t, scenario, dir, runbook+" "+antiShutdownOverride)
	return result.finalMessage + "\n" + result.stream, nil
}

// TestLivePrimitiveRunsAgainstClaudeAdapter is AC-2's live half: a scenario
// authored via the importable livescenario primitive (NOT buried in this
// package) runs against a real Claude agent through the existing launch adapter
// and is graded on durable outcomes — the gated entity stays parked at review
// and the observed output presents the gate review + decision. Running this test
// (`go test -tags live -run TestLivePrimitiveRunsAgainstClaudeAdapter`) against
// a real credential produces the session/ci-run artifact that satisfies AC-2's
// own `Verified by: live …` citation under AC-1's gate — p4 eating its own gate.
func TestLivePrimitiveRunsAgainstClaudeAdapter(t *testing.T) {
	runner := newClaudeLiveRunner(t)
	dir := t.TempDir()
	commandLog := filepath.Join(dir, "evidence", "command.log")
	shimDir := writeRecordedGateLoggingShim(t, buildRecordedGateBinary(t), commandLog)
	runner = runner.withStubPATH(shimDir).(claudeLiveRunner)
	adapter := claudeRunnerAdapter{t: t, runner: runner}

	sc := livescenario.Scenario{
		Name:    "gate-held-via-primitive",
		Runbook: gatePrompt(dir),
		Setup: func(dir string) (string, error) {
			return writeRecordedGateFixtureAt(t, dir).entity, nil
		},
		Assert: func(before, after livescenario.EntityState, observed string) error {
			if err := assertGateHeld(before.Body, after.Body, recordedGateReviewFromClaudeStream(observed)); err != nil {
				return errGraded(err.Error())
			}
			if err := assertRecordedGateHoldLog(readFile(t, commandLog)); err != nil {
				return errGraded(err.Error())
			}
			return nil
		},
	}

	if err := livescenario.Run(context.Background(), dir, sc, adapter); err != nil {
		t.Fatalf("live primitive scenario graded FAIL: %v", err)
	}
}

func errGraded(msg string) error { return &gradedErr{msg} }

type gradedErr struct{ msg string }

func (e *gradedErr) Error() string { return e.msg }
