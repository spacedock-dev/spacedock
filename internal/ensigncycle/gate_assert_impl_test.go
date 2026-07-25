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

func recordedGateHeldEntity(before, briefingID, digest string) string {
	gatesRecord := "gates:\n" +
		"  version: 1\n" +
		"  current:\n" +
		"    gate: gate:recorded-gate-task:validation\n" +
		"  records:\n" +
		"    - id: gate:recorded-gate-task:validation\n" +
		"      stage: validation\n" +
		"      attempts:\n" +
		"        - id: gate-attempt:recorded-gate-task-validation-1\n" +
		"          briefing:\n" +
		"            id: " + briefingID + "\n" +
		"            digest: " + digest + "\n"
	return strings.Replace(before, "\n---\n", "\n"+gatesRecord+"---\n", 1)
}

func assertGateHeld(before, after, review string) error {
	if before == after || !validatingStatus.MatchString(after) || completedSet.MatchString(after) || verdictSetFM.MatchString(after) {
		return fmt.Errorf("gated entity is not held at its open validation boundary")
	}
	authority := after
	if strings.HasPrefix(after, "---\n") {
		if end := strings.Index(after[len("---\n"):], "\n---\n"); end >= 0 {
			authority = after[:len("---\n")+end]
		}
	}
	briefingID := firstRecordedGateMatch(authority, `(?m)^\s+id: (briefing:[^\s]+)$`)
	digest := firstRecordedGateMatch(authority, `(?m)^\s+digest: (sha256:[0-9a-f]{64})$`)
	if briefingID == "" || digest == "" {
		return fmt.Errorf("open bound entity has no durable Briefing identity and digest")
	}
	for _, exact := range []string{
		"gate: gate:recorded-gate-task:validation",
		"id: gate-attempt:recorded-gate-task-validation-1",
		"id: " + briefingID,
		"digest: " + digest,
	} {
		if strings.Count(authority, exact) != 1 {
			return fmt.Errorf("open bound entity count for %q is not 1", exact)
		}
	}
	for _, forbidden := range []string{"type: Resolution", "resolution:", "application:", "state: closed", "state: consumed", "target-stage:", "status: handoff", "status: done", recordedGateDispatchMarker} {
		if strings.Contains(authority, forbidden) {
			return fmt.Errorf("open bound entity contains %q", forbidden)
		}
	}
	if err := assertConciseRecordedGateReview(review); err != nil {
		return fmt.Errorf("semantic gate review: %w", err)
	}
	if reviewTokenCount(review, briefingID) != 1 || reviewTokenCount(review, digest) != 1 {
		return fmt.Errorf("gate review does not name the exact durable Briefing identity and digest once")
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
