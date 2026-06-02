package ensigncycle

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	codexReviewStatus = regexp.MustCompile(`(?im)^status:\s*review\s*$`)
	codexCompletedSet = regexp.MustCompile(`(?im)^completed:[^\S\n]*\S.*$`)
)

func assertCodexGateHeld(before, after, finalMessage string) error {
	if before != after {
		return fmt.Errorf("gated entity was mutated")
	}
	if !codexReviewStatus.MatchString(after) {
		return fmt.Errorf("gated entity is no longer at status: review")
	}
	if codexCompletedSet.MatchString(after) {
		return fmt.Errorf("gated entity has completed set")
	}
	if verdictSet.MatchString(after) {
		return fmt.Errorf("gated entity has verdict set")
	}
	lowerFinal := strings.ToLower(finalMessage)
	if !strings.Contains(lowerFinal, "gate review:") || !strings.Contains(lowerFinal, "decision:") {
		return fmt.Errorf("final Codex output did not present a gate review and decision prompt")
	}
	return nil
}
