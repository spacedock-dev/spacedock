package release

import (
	"strings"
	"testing"
)

// runListJSON renders the shape `gh run list --json databaseId,headSha,conclusion,status`
// returns: a JSON array of run objects. Tests feed constructed fixtures so the
// expected pass/block decision comes from the fixture, never from the workflow.
const greenForCommitJSON = `[
  {"databaseId": 27050060639, "headSha": "abcdef1234567890abcdef1234567890abcdef12", "conclusion": "success", "status": "completed"}
]`

const parkedRunJSON = `[
  {"databaseId": 27118281803, "headSha": "abcdef1234567890abcdef1234567890abcdef12", "conclusion": "", "status": "waiting"}
]`

const greenWrongCommitJSON = `[
  {"databaseId": 27050060639, "headSha": "0000000000000000000000000000000000000000", "conclusion": "success", "status": "completed"}
]`

const emptyRunListJSON = `[]`

const releaseCommit = "abcdef1234567890abcdef1234567890abcdef12"

// TestE2EGatePassesForGreenRunOnReleaseCommit — AC-2 pass case: a run with
// conclusion==success AND headSha==release commit passes the gate.
func TestE2EGatePassesForGreenRunOnReleaseCommit(t *testing.T) {
	dec, err := EvaluateE2EGate([]byte(greenForCommitJSON), releaseCommit, "")
	if err != nil {
		t.Fatalf("EvaluateE2EGate errored on a well-formed run list: %v", err)
	}
	if !dec.Pass {
		t.Fatalf("gate blocked a green run on the release commit: %s", dec.Reason)
	}
	if dec.Waived {
		t.Fatalf("gate reported a waiver for a genuinely green run: %s", dec.Reason)
	}
	if !strings.Contains(dec.Reason, releaseCommit) {
		t.Errorf("pass reason does not cite the matched commit %q: %s", releaseCommit, dec.Reason)
	}
}

// TestE2EGateBlocksParkedRun — AC-2 block case: a waiting/empty-conclusion run
// (offline job green, live lanes still awaiting environment approval) is NOT a
// pass. This is the exact false-pass the spike disproved.
func TestE2EGateBlocksParkedRun(t *testing.T) {
	dec, err := EvaluateE2EGate([]byte(parkedRunJSON), releaseCommit, "")
	if err != nil {
		t.Fatalf("EvaluateE2EGate errored on a well-formed run list: %v", err)
	}
	if dec.Pass {
		t.Fatalf("gate passed a parked run (conclusion empty / status waiting) for the release commit")
	}
}

// TestE2EGateBlocksGreenRunOnWrongCommit — AC-2 block case: a fully green run
// that ran for a DIFFERENT commit does not satisfy the gate; the green live run
// must be bound to the exact tagged commit.
func TestE2EGateBlocksGreenRunOnWrongCommit(t *testing.T) {
	dec, err := EvaluateE2EGate([]byte(greenWrongCommitJSON), releaseCommit, "")
	if err != nil {
		t.Fatalf("EvaluateE2EGate errored on a well-formed run list: %v", err)
	}
	if dec.Pass {
		t.Fatalf("gate passed a green run whose headSha did not match the release commit")
	}
}

// TestE2EGateBlocksEmptyRunList — AC-2 block case: no Runtime Live E2E run at
// all for the line being released ⇒ block the cut.
func TestE2EGateBlocksEmptyRunList(t *testing.T) {
	dec, err := EvaluateE2EGate([]byte(emptyRunListJSON), releaseCommit, "")
	if err != nil {
		t.Fatalf("EvaluateE2EGate errored on an empty run list: %v", err)
	}
	if dec.Pass {
		t.Fatalf("gate passed with no Runtime Live E2E run for the release commit")
	}
}

// TestE2EGateWaiverPassesWhenSet — AC-3: a non-empty waiver reason passes the
// gate even with no matching run, and records the reason for the audit log.
func TestE2EGateWaiverPassesWhenSet(t *testing.T) {
	const reason = "emergency hotfix cut; live matrix unavailable (captain: clkao)"
	dec, err := EvaluateE2EGate([]byte(emptyRunListJSON), releaseCommit, reason)
	if err != nil {
		t.Fatalf("EvaluateE2EGate errored evaluating a waived gate: %v", err)
	}
	if !dec.Pass {
		t.Fatalf("waived gate did not pass: %s", dec.Reason)
	}
	if !dec.Waived {
		t.Fatalf("waived gate did not flag itself as waived: %s", dec.Reason)
	}
	if !strings.Contains(dec.Reason, reason) {
		t.Errorf("waived gate did not record the waiver reason; got: %s", dec.Reason)
	}
}

// TestE2EGateWaiverEnforcesWhenUnset — AC-3: an empty waiver reason does NOT
// bypass the gate; a parked-only run still blocks.
func TestE2EGateWaiverEnforcesWhenUnset(t *testing.T) {
	dec, err := EvaluateE2EGate([]byte(parkedRunJSON), releaseCommit, "")
	if err != nil {
		t.Fatalf("EvaluateE2EGate errored: %v", err)
	}
	if dec.Pass {
		t.Fatalf("an unset waiver still bypassed the gate on a parked run")
	}
	if dec.Waived {
		t.Fatalf("an unset waiver reported itself as a waiver")
	}
}

// TestE2EGateRejectsMalformedRunList — a run list that is not valid JSON must be
// an error (block), never a silent pass.
func TestE2EGateRejectsMalformedRunList(t *testing.T) {
	if _, err := EvaluateE2EGate([]byte(`not json`), releaseCommit, ""); err == nil {
		t.Fatalf("EvaluateE2EGate accepted malformed run-list JSON without error")
	}
}

// TestE2EGateWaiverPassesEvenWithMalformedRunList — AC-3: an explicit captain
// waiver short-circuits BEFORE the run list is consulted, so an emergency cut is
// possible even when `gh run list` output is unusable. The waiver is the
// auditable escape hatch; it must not be defeated by a query failure.
func TestE2EGateWaiverPassesEvenWithMalformedRunList(t *testing.T) {
	const reason = "registry outage; gh unavailable"
	dec, err := EvaluateE2EGate([]byte(`not json`), releaseCommit, reason)
	if err != nil {
		t.Fatalf("waiver did not short-circuit a malformed run list: %v", err)
	}
	if !dec.Pass || !dec.Waived {
		t.Fatalf("waiver did not pass over a malformed run list: pass=%v waived=%v reason=%s", dec.Pass, dec.Waived, dec.Reason)
	}
}

// TestE2EGatePicksMatchingRunAmongMany — when the run list carries several runs
// (the query is not limited to one), the gate passes if ANY run is green for the
// release commit, even if earlier entries are parked or on other commits.
func TestE2EGatePicksMatchingRunAmongMany(t *testing.T) {
	manyJSON := `[
  {"databaseId": 3, "headSha": "0000000000000000000000000000000000000000", "conclusion": "success", "status": "completed"},
  {"databaseId": 2, "headSha": "abcdef1234567890abcdef1234567890abcdef12", "conclusion": "", "status": "waiting"},
  {"databaseId": 1, "headSha": "abcdef1234567890abcdef1234567890abcdef12", "conclusion": "success", "status": "completed"}
]`
	dec, err := EvaluateE2EGate([]byte(manyJSON), releaseCommit, "")
	if err != nil {
		t.Fatalf("EvaluateE2EGate errored: %v", err)
	}
	if !dec.Pass {
		t.Fatalf("gate did not find the green run for the release commit among many: %s", dec.Reason)
	}
	if dec.MatchedRunID != 1 {
		t.Errorf("gate matched run %d; want the green-for-commit run 1", dec.MatchedRunID)
	}
}
