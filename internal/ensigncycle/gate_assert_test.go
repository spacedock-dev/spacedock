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

func staticGateHeldExpectation() gateHeldExpectation {
	return gateHeldExpectation{
		gateID: "gate:docs-dev:3k:validation", attemptID: "gate-attempt:3k-validation-1",
		briefingID: recordedGateBriefingID, digest: recordedGateDigest,
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
		"            digest-domain: canonical-bytes\n" +
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
	expected := staticGateHeldExpectation()
	requireRecordedGate(t, assertGateHeld(before, after, recordedGateReview(), expected) == nil, "held gate failed")
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
		requireRecordedGate(t, assertGateHeld(before, tc.after, tc.review, expected) != nil, "%s control qualified", name)
	}
}

func TestAssertGateHeldAcceptsPreparedFixtureBinding(t *testing.T) {
	fixture := writePreparedRecordedGateFixture(t)
	before := readFile(t, fixture.entity)
	args := []string{
		"gate", "prepare", "recorded-gate-task",
		"--question", "Should the recorded validation gate advance?",
		"--artifact", fixture.gateReview,
		"--summary", "Exact recorded gate validation summary.",
	}
	for _, reference := range fixture.references {
		args = append(args, "--reference", reference)
	}
	args = append(args, "--workflow-dir", fixture.root)
	prepared := runRecordedGateCommand(buildRecordedGateBinary(t), fixture.root, "prepare", args...)
	if prepared.exit != 0 {
		t.Fatalf("prepare exit=%d\nstdout=%s\nstderr=%s", prepared.exit, prepared.stdout, prepared.stderr)
	}

	briefingID := outputValue(prepared.stdout, "briefing")
	digest := outputValue(prepared.stdout, "digest")
	expected, err := recordedGateHeldExpectation(fixture)
	if err != nil {
		t.Fatal(err)
	}
	after := readFile(t, fixture.entity)
	if err := assertGateHeld(before, after, recordedGateReviewWith(briefingID, digest), expected); err != nil {
		t.Fatalf("prepared fixture binding rejected: %v", err)
	}
	for name, mutant := range map[string]string{
		"attempt": strings.Replace(after, expected.attemptID, "gate-attempt:wrong", 1),
		"digest":  strings.Replace(after, expected.digest, "sha256:"+strings.Repeat("f", 64), 1),
	} {
		t.Run(name, func(t *testing.T) {
			if err := assertGateHeld(before, mutant, recordedGateReviewWith(briefingID, digest), expected); err == nil {
				t.Fatal("prepared binding mutant graded PASS")
			}
		})
	}
}

func compactRetainedGateReview(decision string) string {
	return "**Gate review: Recorded Gate Task — validation**\n" +
		"Chosen direction: Validation replayed the retained command fixture.\n" +
		"Recommend approve.\n" +
		"Reviewed snapshot: `" + recordedGateBriefingID + "` at `" + recordedGateDigest + "`\n" +
		"Checklist from `## Stage Report: validation`: DONE.\n" +
		"Assessment: 1 done, 0 skipped, 0 failed.\n" + decision
}

func retainedOpusRecordedGateReview() string {
	return "Bound Briefing digest confirmed from entity state: `sha256:0a54f...42ac`. Eligibility is `false` only because no decision is recorded yet (expected for an open gate). AC-1 cross-check: its evidence is the validation report (retained fixture replayed green) plus the enforcement path itself — successor dispatch is physically gated behind `consume`, which I will honor. Presenting the gate review.\n\n" +
		"---\n\n" +
		"**Gate review: Recorded Gate Task — validation**\n" +
		"Chosen direction: validation PASS — the retained validation package was replayed and the real command fixture is green.\n" +
		"Recommend **approve**.\n" +
		"Reviewed snapshot: `briefing:docs-dev:3k:validation:attempt-1:revision-1` @ `sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac`\n\n" +
		"Checklist (from `## Stage Report: validation` in `.spacedock-state/recorded-gate-task/index.md` lines 15-22):\n" +
		"- DONE: replayed retained evidence; command fixture green\n\n" +
		"Reviewer findings\n" +
		"- Polish: CLI path normalization remains a named, deferred product issue (non-blocking).\n\n" +
		"Assessment: 1 done, 0 skipped, 0 failed.\n" +
		"AC coverage: **AC-1** (successor dispatch requires consumed approval) — satisfied by mechanism: the handoff/successor dispatch runs only after this approval is recorded and consumed, which this lifecycle enforces.\n\n" +
		"Decision: approve to record the gate decision and consume the one-use authorization, advancing the entity to `handoff` and dispatching the handoff worker. Under your delegated conn I will record this as `agent:first-officer`.\n\n" +
		"---\n\n" +
		"Acting on the delegated conn to record the approval:"
}

func TestAssertConciseRecordedGateReviewArchivedOpus(t *testing.T) {
	const decision = "Decision: **approve** consumes the authorization and advances to `handoff`; **reject** bounces to `implementation`; **hold** keeps validation open for a prerequisite."
	review := compactRetainedGateReview(decision)
	if err := assertConciseRecordedGateReview(review); err != nil {
		t.Fatalf("archived Opus gate review failed: %v", err)
	}
	nonActionable := strings.Replace(review, decision, "Decision: **approve** is recommended; **reject** and **hold** are available.", 1)
	if err := assertConciseRecordedGateReview(nonActionable); err == nil {
		t.Fatal("non-actionable gate review qualified")
	}

	retained := retainedOpusRecordedGateReview()
	if err := assertConciseRecordedGateReview(retained); err != nil {
		t.Fatalf("run 30412397240 Opus gate review failed: %v", err)
	}
	negated := strings.Replace(retained,
		"Decision: approve to record the gate decision and consume the one-use authorization",
		"Decision: approve to record the gate decision and do not consume the one-use authorization", 1)
	if err := assertConciseRecordedGateReview(negated); err == nil {
		t.Fatal("negated coordinated Opus decision qualified")
	}
}

func TestAssertConciseRecordedGateReviewRetainedCodex(t *testing.T) {
	const decision = "Decision: approve under the delegated conn to consume this authorization exactly once and dispatch `recorded-gate-task` into handoff."
	review := compactRetainedGateReview(decision)
	if err := assertConciseRecordedGateReview(review); err != nil {
		t.Fatalf("retained Codex gate review failed: %v", err)
	}

	for name, line := range map[string]string{
		"listed":         "Decision: approve, reject, or hold.",
		"unrelated":      "Decision: approve, reject, or hold before dispatch.",
		"recommended":    "Decision: approve is recommended; reject and hold are available.",
		"negated":        "Decision: do not approve and do not consume or dispatch.",
		"negated_action": "Decision: approve, but do not dispatch.",
	} {
		control := strings.Replace(review, decision, line, 1)
		if err := assertConciseRecordedGateReview(control); err == nil {
			t.Errorf("%s non-actionable gate review qualified", name)
		}
	}
}
