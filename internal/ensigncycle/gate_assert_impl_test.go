package ensigncycle

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reviewStatus = regexp.MustCompile(`(?im)^status:\s*review\s*$`)
	completedSet = regexp.MustCompile(`(?im)^completed:[^\S\n]*\S.*$`)
	verdictSetFM = regexp.MustCompile(`(?im)^verdict:[^\S\n]*\S.*$`)
)

// assertGateHeld is host-neutral: it consumes only the entity-state strings
// (before/after) and the FO's final-message string, so it accepts Codex
// `--output-last-message` text and Claude `result`-event text alike. The gate is
// held when the entity is unmutated, still at status: review with no completed /
// verdict set, and the final message presents both a gate review and a decision
// prompt for the human operator.
func assertGateHeld(before, after, finalMessage string) error {
	if before != after {
		return fmt.Errorf("gated entity was mutated")
	}
	if !reviewStatus.MatchString(after) {
		return fmt.Errorf("gated entity is no longer at status: review")
	}
	if completedSet.MatchString(after) {
		return fmt.Errorf("gated entity has completed set")
	}
	if verdictSetFM.MatchString(after) {
		return fmt.Errorf("gated entity has verdict set")
	}
	lowerFinal := strings.ToLower(finalMessage)
	if !strings.Contains(lowerFinal, "gate review:") || !strings.Contains(lowerFinal, "decision:") {
		return fmt.Errorf("final FO output did not present a gate review and decision prompt")
	}
	return nil
}
