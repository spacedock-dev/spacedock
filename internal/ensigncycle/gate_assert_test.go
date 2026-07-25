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
	requireRecordedGate(t, assertConciseRecordedGateReview(strings.Replace(recordedGateReview(), "Decision ask: approve, revise with a concrete finding, or hold for a named prerequisite?", "Choose approve, request revisions, or hold.", 1)) == nil && assertConciseRecordedGateReview(strings.Replace(recordedGateReview(), "Decision ask: approve, revise with a concrete finding, or hold for a named prerequisite?", "", 1)) != nil && assertConciseRecordedGateReview(strings.Replace(recordedGateReview(), "Decision ask: approve, revise with a concrete finding, or hold for a named prerequisite?", "No decision is requested.", 1)) != nil && assertConciseRecordedGateReview(strings.Replace(recordedGateReview(), "Decision ask: approve, revise with a concrete finding, or hold for a named prerequisite?", "Can I provide more information?", 1)) != nil, "decision-ask semantic controls failed")
	for name, tc := range map[string]struct{ after, review string }{"unbound": {before, recordedGateReview()}, "advanced": {strings.Replace(after, "status: validation", "status: handoff", 1), recordedGateReview()}, "resolution": {after + "\ntype: Resolution\n", recordedGateReview()}, "verdict": {strings.Replace(after, "verdict:\n", "verdict: passed\n", 1), recordedGateReview()}, "review": {after, "Gate review: legacy\nDecision: approve?"}, "legacy": {entity, final}} {
		requireRecordedGate(t, assertGateHeld(before, tc.after, tc.review) != nil, "%s control qualified", name)
	}
}
