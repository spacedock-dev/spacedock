// ABOUTME: Pure boot-preamble retry policy — bounds how many fresh launches a
// ABOUTME: classified preamble fumble earns, opting boot-subject scenarios out.
package ensigncycle

// bootPreambleMaxAttempts bounds the boot-preamble retry: a launch whose stream
// classifies as a boot-preamble fumble (wrong-root wander or a broad-search
// contract-file hunt) — not the scenario's own assertion — is retried with a
// fresh launch, up to this many total attempts, before the scenario fails. AC-3's
// live measurement (PR #490, cycle 1) found the contract-file-hunt sub-class
// persisting at 33% (8/24) on the sonnet lane even with fix C's classification
// wired in: classification alone reds a fumble legibly but does not prevent it.
// This is the boot-preamble-retry lever the ideation-stage Proposed approach
// named and explicitly deferred, held for exactly this measured condition
// ("Held as the escalation lever if the live rate (AC-3) stays high after A+B+C
// land"). It retries ONLY on a launch classifyBootPreambleFailure actually flags —
// the same detectors AC-1/AC-2 already prove have teeth (mutation-tested in
// cycle-1 validation) — never on a generic stall or a downstream scenario
// assertion, so a genuine scenario regression is never silently re-run away. A
// boot fumble is not scenario evidence; a scenario assertion is.
const bootPreambleMaxAttempts = 3

// isBootSubjectScenario reports whether the scenario's own SUBJECT is unaided
// boot behavior, in which case a boot-preamble retry would mask exactly the
// regression the scenario exists to catch (the captain confirmed 2026-07-08 that
// shallow-boot's `find {fixture root}` self-orientation call, reachable only
// through this classifier, is a genuine lean-boot discipline lapse, not a
// detector bug). Of the ideation-stage opt-out set (TestLiveDefaultHeadlessStopsAtGate,
// shallow-boot, TestLiveZeroDiscoverReportsAndStops), only shallow-boot launches
// through claudeLiveRunner.run — the other two are standalone Test funcs with
// their own inline launch, never reaching this retry policy.
func isBootSubjectScenario(name string) bool {
	return name == "shallow-boot"
}

// shouldRetryBootPreamble is the retry policy claudeLiveRunner.run's loop follows
// after each launch attempt: given the scenario name, the attempt number just
// completed, and whether THAT attempt's stream classified as a boot-preamble
// fumble (classifyBootPreambleFailure != nil), it reports whether run should
// retry with a fresh launch. Extracted as a pure function so the POLICY — not
// just the underlying detectors — is unit-testable offline, with no subprocess
// and no model spend.
func shouldRetryBootPreamble(scenarioName string, attempt int, preambleClassified bool) bool {
	if !preambleClassified {
		return false
	}
	if isBootSubjectScenario(scenarioName) {
		return false
	}
	return attempt < bootPreambleMaxAttempts
}
