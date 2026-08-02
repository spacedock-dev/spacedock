// ABOUTME: Mechanism 1 (implicit split-root state sync at the end of `gate
// ABOUTME: record`/`gate consume`) and mechanism 2 (`gate record --consume`).
package cli

import (
	"fmt"
	"io"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"github.com/spacedock-dev/spacedock/internal/statesync"
	"github.com/spacedock-dev/spacedock/internal/status"
)

// gateSyncPhase names which half of the gate ceremony a sync line reports — the
// machine-parseable discriminator callers branch on, never on which prose lines
// printed.
type gateSyncPhase string

const (
	gateSyncPhaseRecord  gateSyncPhase = "record"
	gateSyncPhaseConsume gateSyncPhase = "consume"
)

// runGateSync is the gate-verb half of the mechanism-1 implicit-sync seam: for a
// split-root workflow it stages+commits the entity's path-scoped commit unit
// (skip if clean) and publishes via the shared push/rebase/HALT sequence
// (syncActiveEntity), then renders the final `sync=.../phase=...` line on
// stdout — the discriminator every caller branches on, plus the exit code
// (0 pushed/local-only/no-op, 1 failed, 3 halted). An inline workflow performs
// no sync and prints no line (byte-identical to pre-mechanism-1 output). Callers
// must only invoke this when the verb actually wrote (a close, an advance, a
// supersede) — a refusal path never reaches here, keeping refusal diagnostics
// byte-clean and side-effect-free.
func runGateSync(stdout, stderr io.Writer, definitionDir, path string, phase gateSyncPhase, msg string) int {
	checkout, branch, mode, code := resolveStateCheckout("gate "+string(phase), definitionDir, stderr)
	if code != 0 {
		return code
	}
	if mode == status.StateInline {
		return 0
	}
	slug := status.EntitySlug(path)
	_, outcome, err := syncActiveEntity(checkout, branch, slug, msg)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock gate %s: state sync failed: %v\n", phase, err)
		return 1
	}
	switch outcome.Result {
	case statesync.ResultHalted:
		writeSyncHaltStderr(stderr, "gate "+string(phase), branch, outcome)
		fmt.Fprintf(stdout, "sync=halted phase=%s\n", phase)
		return 3
	case statesync.ResultFailed:
		fmt.Fprintf(stderr, "spacedock gate %s: state publication failed; the write is already durable locally:\n%s\n", phase, outcome.Detail)
		fmt.Fprintf(stdout, "sync=failed phase=%s\n", phase)
		return 1
	default:
		fmt.Fprintf(stdout, "sync=%s phase=%s\n", outcome.Result, phase)
		return 0
	}
}

// runGateConsumeAndSync runs the standalone `gate consume` operation, rendering
// the exact pre-mechanism-1 result line (and route, where applicable) byte-
// identically, then — only when the call actually wrote (a real advance or a
// stale-pending supersede) — the mechanism-1 sync with phase=consume. It is
// called both for the standalone `gate consume` subcommand and for the second
// half of `gate record --consume`, so the two entry points share one write path
// and one sync path.
func runGateConsumeAndSync(path, definitionDir string, stdout, stderr io.Writer) int {
	result, err := gates.ConsumeAt(path, definitionDir)
	if err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return 1
	}
	fmt.Fprintf(stdout, "gate=%s attempt=%s application=%s/%s condition=%s eligible=%t consumed=%t target-stage=%s", result.Gate, result.Attempt, result.Action, result.ApplicationState, result.Condition, result.Eligible, result.Consumed, result.TargetStage)
	routed := !result.Consumed && result.Eligible && gates.ApprovedAwaitingMergeRoute(path, definitionDir)
	if routed {
		fmt.Fprintf(stdout, " route=%s", gates.RouteApprovedAwaitingMerge)
	}
	fmt.Fprintln(stdout)

	// A terminal-target route spends nothing (comment in gates.ConsumeAt): no
	// write, so no sync — matching mechanism 1's "only when the verb wrote" rule.
	// result.Wrote (not ApplicationState == "superseded"/"consumed") is the
	// correct signal: EvaluateEligibility copies the attempt's CURRENT
	// application state into ApplicationState on every read, including a pure
	// refusal against an already-superseded or already-consumed application —
	// checking ApplicationState alone would wrongly sync a repeat refusal.
	if result.Wrote {
		msg := fmt.Sprintf("gate: consume %s -> %s", status.EntitySlug(path), result.TargetStage)
		if code := runGateSync(stdout, stderr, definitionDir, path, gateSyncPhaseConsume, msg); code != 0 {
			return code
		}
	}
	if !result.Consumed && !routed {
		return 1
	}
	return 0
}
