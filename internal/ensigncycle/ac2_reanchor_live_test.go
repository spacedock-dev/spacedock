//go:build live

// ABOUTME: Live proof for AC re-anchor behavior through the real Claude adapter.
// ABOUTME: The grade requires the stored revise/feedback/rework gate branch.
package ensigncycle

import (
	"context"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/livescenario"
)

// TestLiveReanchorGateRejectsMeansOnlyRegressed is the behavioral proof: the
// importable scenario runs against a real Claude agent and must store decision
// revise, action feedback, and target/status rework. Final-message narration is
// diagnostic only and cannot satisfy the durable oracle.
// Run it against a real credential:
// `go test -tags live -run TestLiveReanchorGateRejectsMeansOnlyRegressed ./internal/ensigncycle`.
func TestLiveReanchorGateRejectsMeansOnlyRegressed(t *testing.T) {
	runner := newClaudeLiveRunner(t)
	adapter := claudeRunnerAdapter{t: t, runner: runner}

	sc := livescenario.AuthorACReanchorScenario()

	if err := livescenario.Run(context.Background(), t.TempDir(), sc, adapter); err != nil {
		t.Fatalf("live re-anchor durable branch graded FAIL: %v", err)
	}
}
