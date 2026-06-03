package ensigncycle

import (
	"fmt"
	"regexp"
	"strings"
)

const rejectionFixMarker = "shared-rejection-fix: applied"

var implementationStatus = regexp.MustCompile(`(?im)^status:\s*implementation\s*$`)

// assertRejectionFlow is host-neutral: it consumes the post-run entity-state
// string and an observed-output string, so it accepts either host's transcript.
// The flow held when the rejection follow-up applied the exact fix marker, left a
// second implementation stage report, routed the entity back to status:
// implementation, and the observed output surfaced both the rejection and the
// implementation follow-up.
func assertRejectionFlow(entity, observed string) error {
	if !strings.Contains(entity, rejectionFixMarker) {
		return fmt.Errorf("rejection follow-up did not apply the exact fix marker")
	}
	if strings.Count(entity, "## Stage Report: implementation") < 2 {
		return fmt.Errorf("rejection follow-up did not leave a new implementation stage report")
	}
	if !implementationStatus.MatchString(entity) {
		return fmt.Errorf("rejection follow-up did not route the entity back to status: implementation")
	}
	lowerObserved := strings.ToLower(observed)
	if !strings.Contains(lowerObserved, "reject") {
		return fmt.Errorf("FO output/log did not surface the rejection")
	}
	if !strings.Contains(lowerObserved, "implementation") {
		return fmt.Errorf("FO output/log did not surface the implementation follow-up")
	}
	return nil
}

// assertMergeHookGuardHeld is host-neutral: before/after entity state plus an
// observed-output string. The guard held when the entity was unmutated, still at
// status: implementation, and the observed output named both the merge hook and
// the terminal guard refusal — proving the FO could not bypass a registered merge
// hook by terminalizing without pr, mod-block, or force.
func assertMergeHookGuardHeld(before, after, observed string) error {
	if before != after {
		return fmt.Errorf("merge-hook guardrail scenario mutated the entity")
	}
	if !implementationStatus.MatchString(after) {
		return fmt.Errorf("merge-hook guardrail entity is no longer at status: implementation")
	}
	lowerObserved := strings.ToLower(observed)
	if !strings.Contains(lowerObserved, "merge hook") && !strings.Contains(lowerObserved, "merge-hook") {
		return fmt.Errorf("FO output/log did not mention the merge hook guard")
	}
	if !strings.Contains(lowerObserved, "cannot advance to terminal") {
		return fmt.Errorf("FO output/log did not include the terminal guard failure")
	}
	return nil
}
