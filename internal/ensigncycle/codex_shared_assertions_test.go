package ensigncycle

import (
	"fmt"
	"regexp"
	"strings"
)

const codexRejectionFixMarker = "codex-rejection-fix: applied"

var codexImplementationStatus = regexp.MustCompile(`(?im)^status:\s*implementation\s*$`)

func assertCodexRejectionFlow(entity, observed string) error {
	if !strings.Contains(entity, codexRejectionFixMarker) {
		return fmt.Errorf("rejection follow-up did not apply the exact fix marker")
	}
	if strings.Count(entity, "## Stage Report: implementation") < 2 {
		return fmt.Errorf("rejection follow-up did not leave a new implementation stage report")
	}
	if !codexImplementationStatus.MatchString(entity) {
		return fmt.Errorf("rejection follow-up did not route the entity back to status: implementation")
	}
	lowerObserved := strings.ToLower(observed)
	if !strings.Contains(lowerObserved, "reject") {
		return fmt.Errorf("Codex output/log did not surface the rejection")
	}
	if !strings.Contains(lowerObserved, "implementation") {
		return fmt.Errorf("Codex output/log did not surface the implementation follow-up")
	}
	return nil
}

func assertCodexMergeHookGuardHeld(before, after, observed string) error {
	if before != after {
		return fmt.Errorf("merge-hook guardrail scenario mutated the entity")
	}
	if !codexImplementationStatus.MatchString(after) {
		return fmt.Errorf("merge-hook guardrail entity is no longer at status: implementation")
	}
	lowerObserved := strings.ToLower(observed)
	if !strings.Contains(lowerObserved, "merge hook") && !strings.Contains(lowerObserved, "merge-hook") {
		return fmt.Errorf("Codex output/log did not mention the merge hook guard")
	}
	if !strings.Contains(lowerObserved, "cannot advance to terminal") {
		return fmt.Errorf("Codex output/log did not include the terminal guard failure")
	}
	return nil
}
