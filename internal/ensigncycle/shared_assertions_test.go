package ensigncycle

import "testing"

func TestAssertRejectionFlow(t *testing.T) {
	// The full two-cycle end-state: fix marker applied, two implementation reports
	// (original + cycle-2 rework), and two recorded `### Feedback Cycles` entries.
	entity := "---\nstatus: validation\n---\n" +
		rejectionFixMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: Initial implementation\n\n" +
		"## Stage Report: implementation\n\n- DONE: Applied rejection fix\n\n" +
		"### Feedback Cycles\n\n- Cycle 1: REJECTED\n- Cycle 2: PASSED\n"
	observed := "validation was REJECTED; routed follow-up to implementation"

	if err := assertRejectionFlow(entity, observed); err != nil {
		t.Fatalf("expected rejection flow to pass: %v", err)
	}
	if err := assertRejectionFlow("## Stage Report: implementation\n", observed); err == nil {
		t.Fatal("expected missing fix marker to fail")
	}
	// A single-cycle end-state: fix applied, two implementation reports, but only one
	// recorded cycle — the FO never drove the second validation round.
	singleCycle := "---\nstatus: implementation\n---\n" + rejectionFixMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: Initial\n\n" +
		"## Stage Report: implementation\n\n- DONE: Fix\n\n" +
		"### Feedback Cycles\n\n- Cycle 1: REJECTED\n"
	if err := assertRejectionFlow(singleCycle, observed); err == nil {
		t.Fatal("expected a single-cycle end-state (one recorded cycle) to fail")
	}
	// Two cycles recorded but only one implementation report — the rework never left
	// a second report.
	oneReport := "---\nstatus: validation\n---\n" + rejectionFixMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: Only one report\n\n" +
		"### Feedback Cycles\n\n- Cycle 1: REJECTED\n- Cycle 2: PASSED\n"
	if err := assertRejectionFlow(oneReport, observed); err == nil {
		t.Fatal("expected a single implementation report to fail")
	}
	if err := assertRejectionFlow(entity, "all quiet"); err == nil {
		t.Fatal("expected missing rejection output to fail")
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
