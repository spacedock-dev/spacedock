//go:build live

// ABOUTME: Live proof that the promoted livescenario primitive runs against the
// ABOUTME: REAL Claude launch adapter — AC-2's runnable-by-a-real-agent half.
package ensigncycle

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/livescenario"
)

// claudeRunnerAdapter wraps the existing package-private claudeLiveRunner (the
// launch + observed-extract adapter 8y built) as a livescenario.Runner. This is
// the seam AC-2 names: the runner adapter STAYS in ensigncycle; what changes is
// that a scenario is now authored against the importable livescenario primitive
// rather than buried in this package's test files. The adapter copies the staged
// dir's contents into the runner's own workflow root and forwards the launch.
type claudeRunnerAdapter struct {
	t       *testing.T
	runner  claudeLiveRunner
	timeout time.Duration
}

func (a claudeRunnerAdapter) Launch(ctx context.Context, dir, entityPath, runbook string) (string, error) {
	// gitInit so the FO front door sees a workflow git root (the live adapter's
	// fixtures are always git-initialized).
	gitInit(a.t, dir)
	scenario := sharedRuntimeScenario{name: "livescenario-primitive", timeout: a.timeout}
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
	adapter := claudeRunnerAdapter{t: t, runner: runner, timeout: 2 * time.Minute}

	sc := livescenario.Scenario{
		Name:    "gate-held-via-primitive",
		Runbook: gatePrompt(),
		Setup: func(dir string) (string, error) {
			// Reuse 8y's host-neutral gate fixture writer — the SETUP half of the
			// triple — proving the primitive composes the existing setup substrate.
			return writeGateWorkflowNoGit(dir)
		},
		Assert: func(before, after livescenario.EntityState, observed string) error {
			if after.Body != before.Body {
				return errGraded("gated entity was mutated during the run")
			}
			if !strings.Contains(after.Body, "status: review") {
				return errGraded("gated entity left its review status")
			}
			if !strings.Contains(observed, "Gate review") || !strings.Contains(observed, "Decision") {
				return errGraded("observed output did not present the gate review + decision")
			}
			return nil
		},
	}

	if err := livescenario.Run(context.Background(), t.TempDir(), sc, adapter); err != nil {
		t.Fatalf("live primitive scenario graded FAIL: %v", err)
	}
}

func errGraded(msg string) error { return &gradedErr{msg} }

type gradedErr struct{ msg string }

func (e *gradedErr) Error() string { return e.msg }

// writeGateWorkflowNoGit stages the host-neutral gate fixture (README + entity)
// into dir WITHOUT git-initializing — the adapter's Launch git-inits once the
// primitive has captured the pre-run state, matching the live adapter's order.
func writeGateWorkflowNoGit(dir string) (string, error) {
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(gateReadme()), 0o644); err != nil {
		return "", err
	}
	entityPath := filepath.Join(dir, "gate-check.md")
	if err := os.WriteFile(entityPath, []byte(gateEntity()), 0o644); err != nil {
		return "", err
	}
	return entityPath, nil
}
