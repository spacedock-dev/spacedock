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
