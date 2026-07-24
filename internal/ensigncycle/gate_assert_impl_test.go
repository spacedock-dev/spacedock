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

// assertGateHeld grades the v1 no-authority boundary: one committed open
// binding and one semantic review, without close, consume, advance, or dispatch.
func assertGateHeld(before, after, review string) error {
	if before == after {
		return fmt.Errorf("gated entity stayed unbound")
	}
	if !validatingStatus.MatchString(after) {
		return fmt.Errorf("gated entity is no longer validating")
	}
	if completedSet.MatchString(after) {
		return fmt.Errorf("gated entity has completed set")
	}
	if verdictSetFM.MatchString(after) {
		return fmt.Errorf("gated entity has verdict set")
	}
	for _, exact := range []string{
		"state: open", "id: gate-attempt:3k-validation-1",
		"id: " + recordedGateBriefingID,
		"digest: " + recordedGateDigest,
	} {
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
