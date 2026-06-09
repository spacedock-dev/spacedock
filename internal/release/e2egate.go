// ABOUTME: Release-time e2e gate predicate — decides whether a Runtime Live E2E
// ABOUTME: run proves the live matrix passed for the exact commit being released.
package release

import (
	"encoding/json"
	"fmt"
)

// e2eRun is one entry of `gh run list --json databaseId,headSha,conclusion,status`.
// A run "proves" the live matrix passed for a commit only when its overall
// conclusion is success AND its headSha is the release commit. The spike
// established that a parked run (live lanes still awaiting per-environment
// approval) has an empty conclusion / waiting status and so never qualifies.
type e2eRun struct {
	DatabaseID int64  `json:"databaseId"`
	HeadSha    string `json:"headSha"`
	Conclusion string `json:"conclusion"`
	Status     string `json:"status"`
}

// Decision is the gate outcome over a `gh run list` result for a release commit.
// Pass gates whether goreleaser may run. Waived records that the pass came from
// the captain-waiver escape hatch rather than a matched run, so the step log /
// $GITHUB_STEP_SUMMARY can mark the cut as auditable-but-unverified. MatchedRunID
// is the qualifying run when the pass is genuine (0 when waived or blocked).
type Decision struct {
	Pass         bool
	Waived       bool
	MatchedRunID int64
	Reason       string
}

// EvaluateE2EGate decides whether the release commit may be cut. A non-empty
// waiverReason short-circuits to a waived pass BEFORE the run list is consulted,
// so an emergency cut survives an unusable `gh run list` result. Otherwise the
// gate passes only when runListJSON (the output of `gh run list --json
// databaseId,headSha,conclusion,status`) contains a run whose conclusion is
// "success" AND whose headSha equals releaseCommit — the exact green-for-commit
// run the spike proved distinguishes an approved full matrix from a parked,
// offline-only run. A malformed run list is an error (block), never a silent
// pass.
func EvaluateE2EGate(runListJSON []byte, releaseCommit, waiverReason string) (Decision, error) {
	if waiverReason != "" {
		return Decision{
			Pass:   true,
			Waived: true,
			Reason: fmt.Sprintf("e2e gate WAIVED for %s: %s", releaseCommit, waiverReason),
		}, nil
	}

	var runs []e2eRun
	if err := json.Unmarshal(runListJSON, &runs); err != nil {
		return Decision{}, fmt.Errorf("parse gh run list output: %w", err)
	}

	for _, run := range runs {
		if run.Conclusion == "success" && run.HeadSha == releaseCommit {
			return Decision{
				Pass:         true,
				MatchedRunID: run.DatabaseID,
				Reason:       fmt.Sprintf("green Runtime Live E2E run %d matches release commit %s", run.DatabaseID, releaseCommit),
			}, nil
		}
	}

	return Decision{
		Reason: fmt.Sprintf("no conclusion:success Runtime Live E2E run found for release commit %s", releaseCommit),
	}, nil
}
