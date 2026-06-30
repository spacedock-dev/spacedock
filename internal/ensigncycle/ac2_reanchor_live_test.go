//go:build live

// ABOUTME: Live proof for AC-2 — the importable re-anchor scenario runs against
// ABOUTME: the REAL Claude adapter and is graded on the observed gate REJECT.
package ensigncycle

import (
	"context"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/livescenario"
)

// TestLiveReanchorGateRejectsMeansOnlyRegressed is AC-2's behavioral proof: the
// importable AuthorACReanchorScenario runs against a real Claude agent through the
// existing launch adapter (claudeRunnerAdapter) and is graded on durable outcomes
// — the gated entity is left unmutated at its ideation gate and the observed gate
// review recommends REJECT with the end re-anchor / end-value-regression reasoning.
// Run it against a real credential:
// `go test -tags live -run TestLiveReanchorGateRejectsMeansOnlyRegressed ./internal/ensigncycle`.
func TestLiveReanchorGateRejectsMeansOnlyRegressed(t *testing.T) {
	runner := newClaudeLiveRunner(t)
	adapter := claudeRunnerAdapter{t: t, runner: runner}

	sc := livescenario.AuthorACReanchorScenario()

	if err := livescenario.Run(context.Background(), t.TempDir(), sc, adapter); err != nil {
		t.Fatalf("live re-anchor scenario graded FAIL: %v", err)
	}
}
