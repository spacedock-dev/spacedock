//go:build live

// ABOUTME: Live proof that the promoted livescenario primitive runs against the
// ABOUTME: REAL Claude launch adapter — AC-2's runnable-by-a-real-agent half.
package ensigncycle

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/livescenario"
)

func assertRecordedGateHoldLog(log string) error {
	const successfulPrepare = "exit=0\tgate prepare recorded-gate-task "
	prepare := strings.Index(log, successfulPrepare)
	commit := recordedGateSuccessfulStateCommitAfter(log, "recorded-gate-task", prepare)
	head := strings.LastIndex(log, "state-head\t")
	if prepare < 0 || commit < prepare || head < commit || strings.Count(log, successfulPrepare) != 1 ||
		strings.Contains(log[prepare:], " --decision ") || strings.Contains(log[prepare:], "gate consume recorded-gate-task") || strings.Contains(log[prepare:], "dispatch build ") {
		return errGraded("gate hold crossed its committed no-authority boundary")
	}
	return nil
}

func TestAssertRecordedGateHoldLog(t *testing.T) {
	valid := strings.Join([]string{
		"begin\tgate prepare recorded-gate-task --question Approve? --artifact gate-review.md",
		"exit=0\tgate prepare recorded-gate-task --question Approve? --artifact gate-review.md",
		"begin\tstate commit recorded-gate-task",
		"exit=0\tstate commit recorded-gate-task",
		"state-head\t0123456789abcdef0123456789abcdef01234567",
	}, "\n")
	if err := assertRecordedGateHoldLog(valid); err != nil {
		t.Fatalf("successful prepare plus binding commit must hold at the no-authority boundary: %v", err)
	}
	optionBeforeEntity := strings.ReplaceAll(
		valid,
		"state commit recorded-gate-task",
		`state commit --workflow-dir "$WD" recorded-gate-task`,
	)
	if err := assertRecordedGateHoldLog(optionBeforeEntity); err != nil {
		t.Fatalf("supported option-before-entity binding commit must hold at the no-authority boundary: %v", err)
	}

	controls := map[string]string{
		"absent prepare":       strings.Replace(valid, "gate prepare", "gate inspect", 2),
		"failed prepare":       strings.Replace(valid, "exit=0\tgate prepare", "exit=1\tgate prepare", 1),
		"other entity":         strings.ReplaceAll(valid, "recorded-gate-task", "other-task"),
		"absent commit":        strings.Replace(valid, "exit=0\tstate commit recorded-gate-task", "exit=0\tstate inspect recorded-gate-task", 1),
		"failed commit":        strings.Replace(valid, "exit=0\tstate commit recorded-gate-task", "exit=1\tstate commit recorded-gate-task", 1),
		"absent state head":    strings.Replace(valid, "state-head\t", "missing-head\t", 1),
		"duplicate prepare":    valid + "\nexit=0\tgate prepare recorded-gate-task --question Again?",
		"legacy briefing bind": strings.Replace(valid, "gate prepare recorded-gate-task --question Approve? --artifact gate-review.md", "gate record recorded-gate-task --briefing gate-review.md", 2),
		"decision":             valid + "\nexit=0\tgate record recorded-gate-task --decision approve",
		"consume":              valid + "\nexit=0\tgate consume recorded-gate-task",
		"dispatch":             valid + "\nexit=0\tdispatch build --entity-path recorded-gate-task/index.md",
	}
	for name, log := range controls {
		t.Run(name, func(t *testing.T) {
			if err := assertRecordedGateHoldLog(log); err == nil {
				t.Fatal("control crossed the recorded gate hold oracle")
			}
		})
	}
}

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
	runner = withRecordedGateShellShim(t, runner, shimDir).(claudeLiveRunner)
	adapter := claudeRunnerAdapter{t: t, runner: runner}

	sc := livescenario.Scenario{
		Name:    "gate-held-via-primitive",
		Runbook: gatePrompt(dir),
		Setup: func(dir string) (string, error) {
			return writePreparedRecordedGateFixtureAt(t, dir).entity, nil
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
