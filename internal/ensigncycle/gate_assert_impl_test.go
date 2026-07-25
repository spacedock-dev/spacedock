package ensigncycle

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	validatingStatus = regexp.MustCompile(`(?im)^status:\s*validation\s*$`)
	reviewStatus     = regexp.MustCompile(`(?im)^status:\s*review\s*$`)
	completedSet     = regexp.MustCompile(`(?im)^completed:[^\S\n]*\S.*$`)
	verdictSetFM     = regexp.MustCompile(`(?im)^verdict:[^\S\n]*\S.*$`)
)

func assertGateHeld(before, after, review string) error {
	if before == after || !validatingStatus.MatchString(after) || completedSet.MatchString(after) || verdictSetFM.MatchString(after) {
		return fmt.Errorf("gated entity is not held at its open validation boundary")
	}
	for _, exact := range []string{"state: open", "id: gate-attempt:recorded-gate-task-validation-1", "id: " + recordedGateBriefingID, "digest: " + recordedGateDigest} {
		if strings.Count(after, exact) != 1 {
			return fmt.Errorf("open bound entity count for %q is not 1", exact)
		}
	}
	for _, forbidden := range []string{"type: Resolution", "resolution:", "state: closed", "state: consumed", "target-stage:", "status: handoff", "status: done", recordedGateDispatchMarker} {
		if strings.Contains(after, forbidden) {
			return fmt.Errorf("open bound entity contains %q", forbidden)
		}
	}
	if err := assertConciseRecordedGateReview(review); err != nil {
		return fmt.Errorf("semantic gate review: %w", err)
	}
	return nil
}

type recordedGateCodexEvent struct {
	Type string `json:"type"`
	Item struct {
		Type             string `json:"type"`
		Text             string `json:"text"`
		Command          string `json:"command"`
		Status           string `json:"status"`
		AggregatedOutput string `json:"aggregated_output"`
		ExitCode         *int   `json:"exit_code"`
	} `json:"item"`
}
