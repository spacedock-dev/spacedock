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

func recordedGateHeldEntity() string {
	gates := "gates:\n" +
		"  version: 1\n" +
		"  current:\n" +
		"    gate: gate:docs-dev:3k:validation\n" +
		"  records:\n" +
		"    - id: gate:docs-dev:3k:validation\n" +
		"      stage: validation\n" +
		"      attempts:\n" +
		"        - id: gate-attempt:3k-validation-1\n" +
		"          briefing:\n" +
		"            id: " + recordedGateBriefingID + "\n" +
		"            digest: " + recordedGateDigest + "\n" +
		"            digest-domain: raw-file-pin\n" +
		"            room-ref: rooms/validation/attempt-1/revision-1\n"
	return strings.Replace(recordedGateEntity(), "---\n# Recorded Gate Task", gates+"---\n# Recorded Gate Task", 1)
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
	after := recordedGateHeldEntity()
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
	resolved := strings.Replace(after, "          briefing:\n", "          resolution:\n            type: Resolution\n            id: resolution:docs-dev:3k:validation:1\n            briefing: "+recordedGateBriefingID+"\n            by: captain\n            at: 2026-07-28T00:00:00Z\n            decision: approve\n          briefing:\n", 1)
	applied := strings.Replace(resolved, "          briefing:\n", "          application:\n            action: advance\n            target-stage: handoff\n            state: pending\n          briefing:\n", 1)
	for name, tc := range map[string]struct{ after, review string }{
		"unbound": {before, recordedGateReview()}, "advanced": {strings.Replace(after, "status: validation", "status: handoff", 1), recordedGateReview()},
		"malformed-state": {strings.Replace(after, "          briefing:\n", "          state: open\n          briefing:\n", 1), recordedGateReview()}, "wrong-gate": {strings.ReplaceAll(after, "gate:docs-dev:3k:validation", "gate:docs-dev:wrong"), recordedGateReview()},
		"wrong-briefing": {strings.Replace(after, recordedGateBriefingID, "briefing:docs-dev:wrong", 1), recordedGateReview()}, "resolved": {resolved, recordedGateReview()},
		"applied": {applied, recordedGateReview()}, "verdict": {strings.Replace(after, "verdict:\n", "verdict: passed\n", 1), recordedGateReview()},
		"review": {after, "Gate review: legacy\nDecision: approve?"}, "legacy": {entity, final},
	} {
		requireRecordedGate(t, assertGateHeld(before, tc.after, tc.review) != nil, "%s control qualified", name)
	}
}

func TestAssertConciseRecordedGateReviewArchivedOpus(t *testing.T) {
	// PR #572 run 30325567515: the Opus root gate review is reproduced
	// verbatim so the offline matcher covers the exact live decision grammar.
	const review = "Bound Briefing digest read from entity state: `sha256:0a54f1baec...`. Now I'll present the gate review. This is the human decision boundary — I present and stop here without deciding or dispatching the successor.\n\n" +
		"---\n\n" +
		"**Gate review: Recorded Gate Task — validation**\n" +
		"Chosen direction: Validation replayed the retained evidence against the real command fixture and found the binary-owned gate-resolution recorder sound; it recommends proceeding.\n" +
		"Recommend **approve**.\n" +
		"Reviewed snapshot: `briefing:docs-dev:3k:validation:attempt-1:revision-1` — digest `sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac`\n\n" +
		"Checklist (from `## Stage Report: validation` in `.spacedock-state/recorded-gate-task/index.md` lines 15–22):\n" +
		"- DONE: Replayed retained evidence — real command fixture green\n\n" +
		"Reviewer findings:\n" +
		"- Material: The auto AC-scan flags **AC-1** (\"Successor dispatch requires consumed approval\") as citation-count 0 because the report never writes the literal token \"AC-1\". I read the evidence differently: the DONE item is a falsifiable replay — the real gate-resolution command fixture ran green — which is exactly the guardrail AC-1 names. Evidence present; label absent. Your call on whether that satisfies AC-1.\n\n" +
		"Assessment: 1 done, 0 skipped, 0 failed.\n\n" +
		"Decision: **approve** consumes the authorization and advances to `handoff` (successor dispatch is blocked until that consume, per AC-1); **reject** bounces to `implementation` with the concrete ask to tag the replay evidence explicitly to AC-1; **hold** keeps it at validation for a named prerequisite.\n\n" +
		"---\n\n" +
		"The gate is bound, committed, and open. I have not consumed approval, so the successor (`handoff`) stays undispatched — the recorded decision is yours. The validation worker remains alive pending your call."
	const decision = "Decision: **approve** consumes the authorization and advances to `handoff` (successor dispatch is blocked until that consume, per AC-1); **reject** bounces to `implementation` with the concrete ask to tag the replay evidence explicitly to AC-1; **hold** keeps it at validation for a named prerequisite."

	if err := assertConciseRecordedGateReview(review); err != nil {
		t.Fatalf("archived Opus gate review failed: %v", err)
	}

	nonActionable := strings.Replace(review, decision,
		"Decision: **approve** is recommended; **reject** is available; **hold** is available.", 1)
	if err := assertConciseRecordedGateReview(nonActionable); err == nil {
		t.Fatal("non-actionable gate review qualified")
	}
}
