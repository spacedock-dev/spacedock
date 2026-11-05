//go:build live

package ensigncycle

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/livescenario"
)

// stubLiveDriver is the minimal liveDriver for grading-only live tests: the
// durable runner reaches only run (via the scenario adapter) and emitMetrics.
type stubLiveDriver struct{}

func (stubLiveDriver) run(*testing.T, sharedRuntimeScenario, string, string) liveResult {
	return liveResult{}
}
func (stubLiveDriver) emitMetrics(*testing.T, sharedRuntimeScenario, liveResult) {}
func (stubLiveDriver) gradeShallowBootObservation(*testing.T, liveResult)        {}
func (stubLiveDriver) prepareRecordedGate(*testing.T) (liveDriver, func(liveResult)) {
	return stubLiveDriver{}, func(liveResult) {}
}
func (stubLiveDriver) smallestMechanismTrace(liveResult, []string, []string) mechanismTrace {
	return mechanismTrace{}
}
func (stubLiveDriver) lifecycleStream(*testing.T, liveResult) string { return "" }
func (stubLiveDriver) model() string                                 { return "stub" }
func (stubLiveDriver) home() string                                  { return "" }
func (stubLiveDriver) withStubPATH(string) liveDriver                { return stubLiveDriver{} }

// TestDurableJourneyFailureConsultsGapBinding is the falsifier for the durable
// branch: a durable-path failure for a target with an xfail binding must grade
// XFAIL (owned), not FAIL. Before the fix runACValueReanchorJourney called
// t.Fatalf unconditionally, so the bound run below fataled the lane and this
// test could not pass.
func TestDurableJourneyFailureConsultsGapBinding(t *testing.T) {
	failure := durableSemantic("ac-value-reanchor-violation", errors.New("gate did not apply the feedback route"))

	// The gap-aware verdict: bound owns the finding; unbound reds the lane.
	bound := liveScenarioGrade(sharedRuntimeScenario{name: "ac-value-reanchor", gap: liveXFail("pi", "psvqjf0w8xh2txp9604gsvmz")}, failure)
	if bound.status != "xfail" {
		t.Fatalf("bound durable failure graded %q, want xfail", bound.status)
	}
	unbound := liveScenarioGrade(sharedRuntimeScenario{name: "ac-value-reanchor"}, failure)
	if unbound.status != "fail" {
		t.Fatalf("unbound durable failure graded %q, want fail", unbound.status)
	}

	// End to end through the durable runner: a failing durable assertion on a
	// bound journey must be owned, not fatal the lane.
	runACValueReanchorJourney(t, stubLiveDriver{},
		sharedRuntimeScenario{name: "ac-value-reanchor", gap: liveXFail("pi", "psvqjf0w8xh2txp9604gsvmz")},
		func() livescenario.Scenario {
			return livescenario.Scenario{Name: "ac-value-reanchor-test", Setup: func(dir string) (string, error) {
				path := filepath.Join(dir, "entity.md")
				return path, os.WriteFile(path, []byte("status: validation\n"), 0o644)
			}}
		},
		func(livescenario.Scenario, livescenario.EntityState, livescenario.EntityState, string) error {
			return errors.New("gate did not apply the feedback route")
		})
}

// Withdrawn recovery grades through its gap: Haiku's blank-owner binding owns a red (XFAIL); unbound Sonnet still FAILs.
func TestWithdrawnRecoveryGradesThroughGap(t *testing.T) {
	crossed := withdrawnRecoveryBoundary("exit=0\tgate prepare recorded-gate-task --x\nexit=0\tstate commit recorded-gate-task\nexit=0\tgate record recorded-gate-task\n")
	for target, want := range map[string]string{"claude-haiku": "xfail", "claude-sonnet": "fail"} {
		scenario := sharedRuntimeScenario{name: "withdrawn-gate-recovery", gap: selectLiveGap(withdrawnGateRecoveryGaps, target)}
		if got := liveScenarioGrade(scenario, crossed); got.status != want || len(got.codes) != 1 || got.codes[0] != "withdrawn-recovery-boundary" {
			t.Errorf("%s boundary crossing graded %+v, want %s withdrawn-recovery-boundary", target, got, want)
		}
	}
	if err := withdrawnRecoveryBoundary("exit=0\tgate prepare recorded-gate-task --x\nexit=0\tstate commit recorded-gate-task\n"); err != nil {
		t.Errorf("clean prepare/commit log graded %v", err)
	}
	// End to end: the stub run leaves recovery red; under Haiku's binding the runner must not fatal.
	runClaudeWithdrawnGateRecoveryScenario(t, stubLiveDriver{},
		sharedRuntimeScenario{name: "withdrawn-gate-recovery", gap: selectLiveGap(withdrawnGateRecoveryGaps, "claude-haiku")},
		writeWithdrawnGateFixture, assertWithdrawnGateRecovery)
}
