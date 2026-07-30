//go:build live

// ABOUTME: Live proof that a headless First Officer without decision authority
// ABOUTME: drives to, binds, and presents a recorded gate, then stops open.
package ensigncycle

import "testing"

// TestLiveDefaultHeadlessStopsAtGate is AC-a's live proof: default headless `-p`
// with NO conn drives to a gate:true stage and EXITS reporting gate status, without
// greet-stopping, resolving the gate, or writing a verdict past it.
//
// This test is REGISTERED in .github/workflows/runtime-live-e2e.yml's -run list so
// it gates CI (the lean-boot lesson: a registered-but-never-run scenario does not
// gate). It is its own Test* func — a distinct fixture (entity at the initial stage)
// and a distinct prompt (no conn) from TestLiveEnsignCycle's conn-cue drive.
func TestLiveDefaultHeadlessStopsAtGate(t *testing.T) {
	// A wrong-root boot is the most specific diagnosis: a CI env leak lures the FO
	// off `root` into the real repo, where it finds nothing dispatchable and
	// greets-and-stops — which would otherwise look like a (wrong) AC-a pass
	// (nothing driven). Name it FIRST. Pass the symlink-RESOLVED fixture root: on
	// macOS t.TempDir() returns a `/var/folders/...` path while the FO's boot
	// command targets the EvalSymlinks-resolved `/private/var/folders/...` (the same
	// directory), so comparing the unresolved root would false-flag every local
	// macOS run as a wander. The CI Linux runner has no such symlink, so this is a
	// no-op there; resolving here keeps the detector accurate on BOTH.
	// The shared runner performs: registerClaudeLiveFailureDiagnostic(t, detectClaudeLiveFailureDiagnostic(stream, rootResolved))
	for _, scenario := range []struct {
		name, intent string
		run          func(*testing.T, liveDriver, sharedRuntimeScenario)
	}{
		{
			name:   "default-headless-recorded-gate-stop",
			intent: "drive, bind, commit, present, and stop open without decision authority",
			run:    runClaudeGateGuardrailScenario,
		},
		{
			name:   "default-headless-withdrawn-gate-recovery",
			intent: "prepare and commit a successor for a withdrawn attempt, present it, and stop open",
			run:    runClaudeWithdrawnGateRecoveryScenario,
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.run(t, newClaudeLiveRunner(t), sharedRuntimeScenario{
				name:          scenario.name,
				oldPythonTest: "tests/test_gate_guardrail.py",
				intent:        scenario.intent,
			})
		})
	}
}
