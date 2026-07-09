// ABOUTME: Offline proof of the boot-preamble-retry policy — retries a classified
// ABOUTME: fumble up to the budget, never for a boot-subject scenario or a clean boot.
package ensigncycle

import "testing"

// TestShouldRetryBootPreambleRetriesUpToBudget proves the retry policy's core
// value: a non-boot-subject scenario whose stream classifies as a boot-preamble
// fumble earns a retry on every attempt strictly under bootPreambleMaxAttempts,
// and stops earning one once the budget is exhausted — the scenario fails rather
// than retrying forever.
func TestShouldRetryBootPreambleRetriesUpToBudget(t *testing.T) {
	for attempt := 1; attempt < bootPreambleMaxAttempts; attempt++ {
		if !shouldRetryBootPreamble("filing", attempt, true) {
			t.Errorf("attempt %d/%d: want a retry (still under budget), got none", attempt, bootPreambleMaxAttempts)
		}
	}
	if shouldRetryBootPreamble("filing", bootPreambleMaxAttempts, true) {
		t.Errorf("attempt %d (the budget itself): want no further retry, got one", bootPreambleMaxAttempts)
	}
}

// TestShouldRetryBootPreambleNeverRetriesACleanBoot proves the retry never fires
// when the attempt did NOT classify as a preamble fumble — a clean boot, or a
// genuine downstream scenario-assertion failure (which classifyBootPreambleFailure
// never flags), must not trigger a wasted extra launch.
func TestShouldRetryBootPreambleNeverRetriesACleanBoot(t *testing.T) {
	if shouldRetryBootPreamble("filing", 1, false) {
		t.Error("a clean (non-preamble) attempt must never be retried")
	}
}

// TestShouldRetryBootPreambleOptsOutBootSubjectScenarios proves the opt-out: a
// boot-subject scenario (shallow-boot) never retries, even on attempt 1 with a
// classified fumble — a boot-preamble retry would mask exactly the lean-boot
// discipline regression the scenario exists to catch (the captain-confirmed
// 2026-07-08 shallow-boot finding).
func TestShouldRetryBootPreambleOptsOutBootSubjectScenarios(t *testing.T) {
	if shouldRetryBootPreamble("shallow-boot", 1, true) {
		t.Error("shallow-boot (a boot-subject scenario) must never retry a classified preamble fumble")
	}
}

// TestIsBootSubjectScenario locks the membership the retry policy keys on: only
// shallow-boot is boot-subject among the scenarios claudeLiveRunner.run launches
// (TestLiveDefaultHeadlessStopsAtGate and TestLiveZeroDiscoverReportsAndStops are
// the OTHER two ideation-stage opt-outs, but they never call claudeLiveRunner.run
// at all, so they never reach this policy).
func TestIsBootSubjectScenario(t *testing.T) {
	if !isBootSubjectScenario("shallow-boot") {
		t.Error("shallow-boot must be boot-subject")
	}
	for _, name := range []string{"gate-guardrail", "rejection-flow", "feedback-3-cycle-escalation",
		"merge-hook-guardrail", "filing", "self-evidence-merge-triage",
		"smallest-sufficient-mechanism", "keep-moving-posture"} {
		if isBootSubjectScenario(name) {
			t.Errorf("%q must not be boot-subject — it is one of the 8 non-boot-subject scenarios claudeScenarioRunners() maps", name)
		}
	}
}
