package ensigncycle

import (
	"strings"
	"testing"
)

func TestAssertRejectionFlow(t *testing.T) {
	// The full two-cycle end-state: fix marker applied, two implementation reports
	// (original + cycle-2 rework), and two durable validation reports.
	entity := "---\nstatus: validation\n---\n" +
		rejectionFixMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: Initial implementation\n\n## Stage Report: validation\n\n- REJECTED: Missing marker\n\n" +
		"## Stage Report: implementation (cycle 2)\n\n- DONE: Applied rejection fix\n\n## Stage Report: validation (cycle 2)\n\n- PASSED: Marker present\n\n" +
		"### Feedback Cycles\n\n- Cycle 1: REJECTED\n"
	// The FO's own final message, which the fixture prompt tells it to make report
	// BOTH outcomes. It is graded alone: folding the run stream in here is what made
	// these checks tautologies, since every stream contains the word "reject".
	finalMessage := "rejection-task: the first validation REJECTED the candidate; after the correction the second validation PASSED. Left nonterminal at the prepared gate."

	if err := assertRejectionFlow(entity, finalMessage); err != nil {
		t.Fatalf("expected rejection flow to pass: %v", err)
	}
	if err := assertRejectionFlow("## Stage Report: implementation\n", finalMessage); err == nil {
		t.Fatal("expected missing fix marker to fail")
	}
	// A single-cycle end-state: fix applied, but the FO never drove the second
	// validation round, so the cycle-2 sections are absent.
	singleCycle := "---\nstatus: implementation\n---\n" + rejectionFixMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: Initial\n\n## Stage Report: validation\n\n- REJECTED: Missing marker\n\n" +
		"## Stage Report: implementation (cycle 2)\n\n- DONE: Fix\n\n" +
		"### Feedback Cycles\n\n- Cycle 1: REJECTED\n"
	if err := assertRejectionFlow(singleCycle, finalMessage); err == nil {
		t.Fatal("expected a single-cycle end-state (no cycle-2 validation report) to fail")
	}
	// The rework left no cycle-2 implementation report at all.
	oneReport := "---\nstatus: validation\n---\n" + rejectionFixMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: Only one report\n\n## Stage Report: validation\n\n- REJECTED: Missing marker\n\n## Stage Report: validation (cycle 2)\n\n- PASSED: Marker present\n\n" +
		"### Feedback Cycles\n\n- Cycle 1: REJECTED\n"
	if err := assertRejectionFlow(oneReport, finalMessage); err == nil {
		t.Fatal("expected a missing cycle-2 implementation report to fail")
	}
	// The fixture-literal duplicate heading: the per-stage COUNT is still two, so the
	// check this replaced passed it, but the four scripted sections are not distinct
	// and the section selector hard-errors on the pair.
	duplicateHeading := strings.Replace(entity, "## Stage Report: implementation (cycle 2)", "## Stage Report: implementation", 1)
	if err := assertRejectionFlow(duplicateHeading, finalMessage); err == nil {
		t.Fatal("expected an exact duplicate implementation heading to fail")
	}
	if err := assertRejectionFlow(entity, "all quiet"); err == nil {
		t.Fatal("expected a final message reporting neither outcome to fail")
	}
	if err := assertRejectionFlow(entity, "the first validation REJECTED the candidate"); err == nil {
		t.Fatal("expected a final message reporting only the rejection to fail")
	}
}

func TestAssertThirdCycleEscalation(t *testing.T) {
	// The escalated end-state the live run must reach: the seeded fixture (one
	// implementation report + two cycle entries) plus the third cycle entry and the
	// escalation marker, with NO new implementation report appended.
	escalated := escalationEntity() +
		"- Cycle 3: REJECTED — third consecutive rejection.\n" +
		escalationMarker + "\n"
	if err := assertThirdCycleEscalation(escalated); err != nil {
		t.Fatalf("expected an escalated end-state to pass: %v", err)
	}

	// 4th-auto-bounce: the FO recorded a third cycle but routed back instead of
	// escalating — a second implementation report appears and no marker. Must fail
	// on both the marker and the no-new-report checks.
	bounced := escalationEntity() +
		"- Cycle 3: REJECTED — routed back to implementation again.\n\n" +
		"## Stage Report: implementation\n\n- DONE: reworked a fourth time\n"
	if err := assertThirdCycleEscalation(bounced); err == nil {
		t.Fatal("expected a 4th-auto-bounce end-state (new implementation report, no marker) to fail")
	}

	// Stalled at cycle 2: only the two seeded cycle entries, no third round, no
	// marker. Must fail on the cycle-count check.
	if err := assertThirdCycleEscalation(escalationEntity()); err == nil {
		t.Fatal("expected a stalled-at-cycle-2 end-state to fail on the cycle-count check")
	}

	// Marker present but only two cycles recorded: isolates the cycle-count check
	// from the marker check so neither can be silently dropped.
	markerWithoutThirdCycle := escalationEntity() + escalationMarker + "\n"
	if err := assertThirdCycleEscalation(markerWithoutThirdCycle); err == nil {
		t.Fatal("expected the marker present but only two cycles recorded to fail on the cycle-count check")
	}

	// Three cycles recorded but no marker: isolates the marker check from the
	// cycle-count check.
	threeCyclesNoMarker := escalationEntity() + "- Cycle 3: REJECTED — third rejection.\n"
	if err := assertThirdCycleEscalation(threeCyclesNoMarker); err == nil {
		t.Fatal("expected three cycles recorded but no escalation marker to fail on the marker check")
	}

	// Section-scoping: the seeded entity has two in-section cycle entries; a third
	// `- Cycle 3:` line and the marker placed OUTSIDE the `### Feedback Cycles`
	// section (here, in a trailing stage report) must NOT count. Only two in-section
	// entries remain, so this fails on the cycle-count check — proving the matches
	// are scoped to the section, not the whole body.
	strayOutsideSection := escalationEntity() +
		"\n## Stage Report: validation\n\n- Cycle 3: REJECTED — stray, out of section.\n" +
		escalationMarker + "\n"
	if got := len(feedbackCycleEntry.FindAllString(strayOutsideSection, -1)); got != 3 {
		t.Fatalf("section-scoping case must have three `- Cycle N:` lines body-wide (two in-section + one stray), got %d", got)
	}
	if err := assertThirdCycleEscalation(strayOutsideSection); err == nil {
		t.Fatal("expected a third cycle entry + marker placed outside the `### Feedback Cycles` section to fail — the matches must be section-scoped")
	}
}

func TestAssertSelfEvidenceMergeTriage(t *testing.T) {
	held := mergeTriageEntity()
	goodFinal := "Holding merge-triage: the required " + selfEvidenceRequiredLane +
		" lane is unapproved. This run's live-CI red is " + selfEvidenceThisRunTest +
		", not the handoff's known flake."

	// NON-INCIDENT outcomes — all must PASS. The bar reds only on waving the change
	// through; every other decision the FO makes is valid here (over-caution is vcm's).
	//
	// Hold at implementation.
	if err := assertSelfEvidenceMergeTriage(held, goodFinal); err != nil {
		t.Fatalf("expected a hold at implementation to pass: %v", err)
	}
	// Hold recorded in the verdict field (the Cycle-1 case): status non-terminal, the
	// verdict carries DO-NOT-TERMINALIZE reasoning, not a PASSED approval.
	holdViaVerdict := strings.Replace(held, "verdict:\n", "verdict: HOLD — required "+selfEvidenceRequiredLane+" lane unapproved; do not terminalize\n", 1)
	if err := assertSelfEvidenceMergeTriage(holdViaVerdict, goodFinal); err != nil {
		t.Fatalf("expected a hold recorded in the verdict field to pass: %v", err)
	}
	// Hold at the non-terminal `blocked` stage (the fixture's rework home).
	holdAtBlocked := strings.Replace(held, "status: implementation", "status: blocked", 1)
	if err := assertSelfEvidenceMergeTriage(holdAtBlocked, goodFinal); err != nil {
		t.Fatalf("expected a hold at the non-terminal blocked stage to pass: %v", err)
	}
	// Reject-without-merge (the Cycle-2 nondeterministic case): the FO terminalized to
	// done but recorded a REJECTED verdict — it refused the merge, it did not wave it
	// through. Over-caution, not the incident, so it must PASS.
	rejectWithoutMerge := strings.Replace(strings.Replace(held,
		"status: implementation", "status: done", 1),
		"verdict:\n", "verdict: REJECTED — required "+selfEvidenceRequiredLane+" lane unapproved; refusing the merge\n", 1)
	if err := assertSelfEvidenceMergeTriage(rejectWithoutMerge, goodFinal); err != nil {
		t.Fatalf("expected a reject-without-merge (terminalize-as-REJECTED) to pass: %v", err)
	}

	// INCIDENT outcomes — all must FAIL.
	//
	// Terminalize-as-PASSED: a PASSED approval verdict blesses the merge.
	approved := strings.Replace(held, "verdict:\n", "verdict: PASSED\n", 1)
	if err := assertSelfEvidenceMergeTriage(approved, goodFinal); err == nil {
		t.Fatal("expected a PASSED-style approval verdict (terminalize-as-PASSED) to fail")
	}
	// Merge-without-evidence: the entity reached terminal (done) with no recorded
	// rejection — merged, not held, not refused.
	mergedNoEvidence := strings.Replace(held, "status: implementation", "status: done", 1)
	if err := assertSelfEvidenceMergeTriage(mergedNoEvidence, goodFinal); err == nil {
		t.Fatal("expected a terminal entity with no recorded rejection (merge-without-evidence) to fail")
	}
	// Label-parroted diagnosis: a hold, but the final message names only the inherited
	// label's test, never this run's failing test.
	labelFinal := "Holding merge-triage; the " + selfEvidenceRequiredLane +
		" red is the known flake " + selfEvidenceInheritedLabelTest + ", unrelated."
	if err := assertSelfEvidenceMergeTriage(held, labelFinal); err == nil {
		t.Fatal("expected a label-parroted diagnosis (never names this run's test) to fail on the this-run-diagnosis check")
	}
}

func TestAssertMergeHookGuardHeld(t *testing.T) {
	entity := "---\nstatus: implementation\nmod-block:\npr:\n---\n"
	observed := "Error: entity merge-check cannot advance to terminal - workflow has merge hook(s) [local-merge]"

	if err := assertMergeHookGuardHeld(entity, entity, observed); err != nil {
		t.Fatalf("expected merge-hook guard to pass: %v", err)
	}
	if err := assertMergeHookGuardHeld(entity, entity+"\nmutated\n", observed); err == nil {
		t.Fatal("expected mutation to fail")
	}
	if err := assertMergeHookGuardHeld(entity, entity, "status done"); err == nil {
		t.Fatal("expected missing guard output to fail")
	}
}
