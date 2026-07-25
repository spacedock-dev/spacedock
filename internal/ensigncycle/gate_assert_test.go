package ensigncycle

import (
	"regexp"
	"strings"
	"testing"
)

var recordedGateStateHead = regexp.MustCompile(`(?:^|\n)state-head\t[0-9a-f]{40}(?:\s|$)`)

type recordedGateHost struct {
	extract                                                  func(string) string
	name, commit, review, decision, narration, failed, child string
}

func requireRecordedGate(t *testing.T, ok bool, format string, args ...any) {
	t.Helper()
	if !ok {
		t.Fatalf(format, args...)
	}
}

func TestAssertGateHeld(t *testing.T) {
	entity := "---\n" +
		"id: gate-check\n" +
		"title: Gate Check\n" +
		"status: review\n" +
		"completed:\n" +
		"verdict:\n" +
		"---\n" +
		"# Gate Check\n\n" +
		"## Stage Report: draft\n\n" +
		"- DONE: Draft exists\n" +
		"  fixture evidence\n" +
		"\n### Summary\n\nReady for review.\n"
	final := "Gate review: Gate Check - review\nRecommend approve.\nDecision: approve to enter done."

	before := recordedGateEntity()
	after := before + "\ngates:\n  records:\n    - id: gate:docs-dev:3k:validation\n      attempts:\n        - id: gate-attempt:3k-validation-1\n          state: open\n          briefing:\n            id: " + recordedGateBriefingID + "\n            digest: " + recordedGateDigest + "\n"
	requireRecordedGate(t, assertGateHeld(before, after, recordedGateReview()) == nil, "held gate failed")
	decision := "Decision ask: approve, revise with a concrete finding, or hold for a named prerequisite?"
	for name, line := range map[string]string{
		"baseline":        decision,
		"retained-claude": "Decision: approve to consume this authorization and advance recorded-gate-task from validation into the handoff stage for dispatch.",
		"semantic-label":  "Choose approve to enter handoff, revise with findings, or hold for a prerequisite.",
	} {
		review := strings.Replace(recordedGateReview(), decision, line, 1)
		requireRecordedGate(t, assertConciseRecordedGateReview(review) == nil, "%s decision prompt failed", name)
	}
	for name, line := range map[string]string{
		"missing":     "",
		"disposition": "Decision: continue into handoff after review.",
		"vague":       "Decision: approve with confidence.",
		"narration":   "The package can advance after a decision.",
	} {
		review := strings.Replace(recordedGateReview(), decision, line, 1)
		requireRecordedGate(t, assertConciseRecordedGateReview(review) != nil, "%s decision control qualified", name)
	}
	for name, tc := range map[string]struct{ after, review string }{"unbound": {before, recordedGateReview()}, "advanced": {strings.Replace(after, "status: validation", "status: handoff", 1), recordedGateReview()}, "resolution": {after + "\ntype: Resolution\n", recordedGateReview()}, "verdict": {strings.Replace(after, "verdict:\n", "verdict: passed\n", 1), recordedGateReview()}, "review": {after, "Gate review: legacy\nDecision: approve?"}, "legacy": {entity, final}} {
		requireRecordedGate(t, assertGateHeld(before, tc.after, tc.review) != nil, "%s control qualified", name)
	}
}
