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
	const prepareToken = "exit=0\tgate prepare recorded-gate-task "
	prepare := strings.Index(log, prepareToken)
	commit := strings.LastIndex(log, "exit=0\tstate commit recorded-gate-task")
	head := strings.LastIndex(log, "state-head\t")
	const boundary = "gate hold crossed its committed no-authority boundary: "
	switch {
	case prepare < 0:
		return errGraded(boundary + "no successful gate prepare recorded")
	case commit < prepare:
		return errGraded(boundary + "state commit missing or before the successful gate prepare")
	case head < commit:
		return errGraded(boundary + "state-head missing or before the state commit")
	case strings.Count(log, prepareToken) != 1:
		return errGraded(boundary + "more than one successful gate prepare recorded")
	case strings.Contains(log[prepare:], " --decision "):
		return errGraded(boundary + "a decision was recorded after prepare")
	case strings.Contains(log[prepare:], "gate consume recorded-gate-task"):
		return errGraded(boundary + "the gate was consumed after prepare")
	case strings.Contains(log[prepare:], "dispatch build "):
		return errGraded(boundary + "a successor was dispatched after prepare")
	}
	return nil
}

func TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle(t *testing.T) {
	const prepared = "exit=1\tgate prepare recorded-gate-task validation\n" +
		"exit=0\tgate prepare recorded-gate-task validation\n" +
		"exit=0\tstate commit recorded-gate-task\n" +
		"state-head\tabc123\n"
	if err := assertRecordedGateHoldLog(prepared); err != nil {
		t.Fatalf("prepare-first hold log rejected: %v", err)
	}

	for name, tc := range map[string]struct {
		mutation string
		want     string
	}{
		"retired bind":      {strings.Replace(prepared, "exit=0\tgate prepare recorded-gate-task validation", "exit=0\tgate record recorded-gate-task --briefing briefing.md", 1), "no successful gate prepare recorded"},
		"missing commit":    {strings.Replace(prepared, "exit=0\tstate commit recorded-gate-task\n", "", 1), "state commit missing or before the successful gate prepare"},
		"decision":          {prepared + "exit=0\tgate record recorded-gate-task --decision approve\n", "a decision was recorded after prepare"},
		"consume":           {prepared + "exit=0\tgate consume recorded-gate-task\n", "the gate was consumed after prepare"},
		"successor build":   {prepared + "exit=0\tdispatch build successor\n", "a successor was dispatched after prepare"},
		"duplicate prepare": {prepared + "exit=0\tgate prepare recorded-gate-task validation\n", "more than one successful gate prepare recorded"},
	} {
		t.Run(name, func(t *testing.T) {
			err := assertRecordedGateHoldLog(tc.mutation)
			if err == nil {
				t.Fatal("mutated hold log unexpectedly accepted")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error=%q want substring %q", err.Error(), tc.want)
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
	runner = runner.withStubPATH(shimDir).(claudeLiveRunner)
	adapter := claudeRunnerAdapter{t: t, runner: runner}
	var fixture recordedGateFixture

	sc := livescenario.Scenario{
		Name:    "gate-held-via-primitive",
		Runbook: gatePrompt(dir),
		Setup: func(dir string) (string, error) {
			fixture = writePreparedRecordedGateFixtureAt(t, dir)
			return fixture.entity, nil
		},
		Assert: func(before, after livescenario.EntityState, observed string) error {
			expected, err := recordedGateHeldExpectation(fixture)
			if err != nil {
				return errGraded(err.Error())
			}
			if err := assertGateHeld(before.Body, after.Body, expected); err != nil {
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
