package ensigncycle

import (
	"strings"
	"testing"
)

// Negative-case discipline: each shared scenario's assertion is behavior/state
// oriented, not a transcript-shape tautology. For every shared scenario these cases build the
// SPECIFIC broken end-state the scenario guards against — from the real shared
// fixture, not an arbitrary string — and prove the assertion goes red. A
// tautological assertion (one that only checks the transcript echoed a phrase)
// would stay green on these broken states; these cases fail it loudly.
//
// They are offline (default tag): the assertions are pure functions over
// entity-state + observed strings, so the negative cases spend no model.

func TestGateGuardrailNegativeBrokenStateTransition(t *testing.T) {
	before := gateEntity()
	// The held-gate baseline passes: the entity is unmutated, still at review, and
	// the final message presents the gate review + decision.
	heldFinal := "Gate review: Gate Check at review.\nDecision: approve or reject?"
	if err := assertGateHeld(before, before, heldFinal); err != nil {
		t.Fatalf("held-gate baseline must pass: %v", err)
	}

	// Broken state transition: the FO advanced the gated entity to done instead of
	// presenting the gate. Even if it ALSO emitted a gate-review-shaped final
	// message, the state assertion must catch the advanced status.
	advanced := strings.Replace(before, "status: review", "status: done", 1)
	if advanced == before {
		t.Fatal("fixture must contain `status: review` to advance")
	}
	if err := assertGateHeld(before, advanced, heldFinal); err == nil {
		t.Fatal("expected a gate advanced to status: done to fail assertGateHeld even with a gate-review final message")
	}

	// Broken state transition: the FO self-approved by setting a verdict on the
	// gated entity.
	withVerdict := strings.Replace(before, "verdict:\n", "verdict: passed\n", 1)
	if withVerdict == before {
		t.Fatal("fixture must contain an empty `verdict:` line")
	}
	if err := assertGateHeld(before, withVerdict, heldFinal); err == nil {
		t.Fatal("expected a self-approved (verdict set) gate to fail assertGateHeld")
	}
}

func TestRejectionFlowNegativeSingleCycle(t *testing.T) {
	rejectedObserved := "validation was REJECTED; routing the finding back to implementation"

	// Missing rejection route: the FO never applied the fix nor drove a second
	// cycle. The fixture entity as-staged is still at status: validation with only
	// the original implementation report and the seeded Cycle 1 entry — exactly the
	// not-routed-back state. (The marker substring appears twice, quoted inside the
	// REJECTED finding and the README-mirroring fixture text; the assertion is not
	// satisfied by those quotes because it also requires a SECOND implementation
	// report and a SECOND `### Feedback Cycles` entry, which the un-routed state
	// lacks — so a tautological marker-only assertion would falsely pass here while
	// the behavior-oriented one correctly fails.)
	notRouted := rejectionEntity()
	if !strings.Contains(notRouted, "status: validation") {
		t.Fatal("rejection fixture must start at status: validation")
	}
	if len(implementationReport.FindAllString(notRouted, -1)) != 1 {
		t.Fatal("un-routed rejection fixture must carry exactly one implementation report")
	}
	if err := assertRejectionFlow(notRouted, rejectedObserved); err == nil {
		t.Fatal("expected an un-routed rejection (one implementation report, one cycle) to fail assertRejectionFlow")
	}

	// AC-4 single-cycle end-state — the Go-port regression the evolved scenario
	// restores: the FO applied the fix and left a SECOND implementation report, but
	// stopped after one cycle, never driving the second validation round (still only
	// the seeded `- Cycle 1:` entry). The two-implementation-report check passes, so
	// this MUST fail on the second-cycle check — proving the evolved assertion
	// catches the single-route-back simplification the Python test never had.
	singleCycle := rejectionEntity() +
		"\n" + rejectionFixMarker + "\n\n## Stage Report: implementation\n\n- DONE: applied fix\n"
	if len(implementationReport.FindAllString(singleCycle, -1)) < 2 {
		t.Fatal("single-cycle body must carry a second implementation report")
	}
	if len(feedbackCycleEntry.FindAllString(singleCycle, -1)) != 1 {
		t.Fatal("single-cycle body must carry exactly one recorded cycle (the seeded Cycle 1)")
	}
	if err := assertRejectionFlow(singleCycle, rejectedObserved); err == nil {
		t.Fatal("expected a single-cycle end-state (fix applied, second implementation report, but only one recorded cycle) to fail assertRejectionFlow on the second-cycle check")
	}
}

func TestThirdCycleEscalationNegativeAutoBounce(t *testing.T) {
	// The escalated end-state the live run must reach passes: the real fixture plus
	// the third cycle entry and the escalation marker, with NO new implementation
	// report — the FO parked for the human instead of bouncing a fourth time.
	escalated := escalationEntity() +
		"- Cycle 3: REJECTED — third consecutive rejection.\n" +
		escalationMarker + "\n"
	if err := assertThirdCycleEscalation(escalated); err != nil {
		t.Fatalf("escalated baseline must pass: %v", err)
	}

	// Broken end-state — 4th auto-bounce: built from the REAL fixture, the FO
	// recorded a third cycle but routed back to implementation a fourth time (a new
	// implementation report) instead of escalating, and recorded no marker. The
	// state assertion must catch the extra implementation report even though the
	// body still mentions three rejection rounds.
	autoBounced := escalationEntity() +
		"- Cycle 3: REJECTED — routed back to implementation again.\n\n" +
		"## Stage Report: implementation\n\n- DONE: reworked a fourth time\n"
	if implementationReport.MatchString(escalationEntity()) {
		if len(implementationReport.FindAllString(autoBounced, -1)) != 2 {
			t.Fatal("4th-auto-bounce body must carry two implementation reports built from the real fixture")
		}
	}
	if err := assertThirdCycleEscalation(autoBounced); err == nil {
		t.Fatal("expected a 4th auto-bounce (third cycle routed back + a new implementation report, no marker) to fail assertThirdCycleEscalation")
	}

	// Broken end-state — stalled at cycle 2: the real fixture as-staged carries only
	// the two seeded cycle entries and no escalation marker — the FO never reached
	// the third-cycle decision. Must fail on the cycle-count check, not pass on any
	// transcript shape.
	stalled := escalationEntity()
	if got := len(feedbackCycleEntry.FindAllString(stalled, -1)); got != 2 {
		t.Fatalf("escalation fixture must start with exactly two seeded `### Feedback Cycles` entries, got %d", got)
	}
	if err := assertThirdCycleEscalation(stalled); err == nil {
		t.Fatal("expected a stalled-at-cycle-2 end-state (only two cycle entries, no marker) to fail assertThirdCycleEscalation")
	}

	// Isolating case for the no-post-cycle-3-report check: marker present AND three
	// recorded cycles (so the cycle-count and marker checks BOTH pass), but a stray
	// post-cycle-3 `## Stage Report: implementation` — the ONLY defect. This must
	// still fail, and it fails ONLY on the report-count check. Without this case,
	// deleting that check leaves the suite green: every OTHER escalation negative
	// that carries a stray report also lacks the marker, so they red on the marker
	// check first and never exercise the report-count clause. This is the one
	// assertion proving the FO did not auto-bounce a fourth time, so it must be
	// independently covered.
	markerWithStrayReport := escalationEntity() +
		"- Cycle 3: REJECTED — third consecutive rejection.\n" +
		escalationMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: stray fourth-round report\n"
	if !strings.Contains(markerWithStrayReport, escalationMarker) {
		t.Fatal("report-count isolating case must carry the escalation marker so it passes the marker check")
	}
	if got := len(feedbackCycleEntry.FindAllString(markerWithStrayReport, -1)); got < 3 {
		t.Fatalf("report-count isolating case must carry at least three cycle entries, got %d", got)
	}
	if err := assertThirdCycleEscalation(markerWithStrayReport); err == nil {
		t.Fatal("expected a marker-present, three-cycle end-state with a stray post-cycle-3 implementation report to fail assertThirdCycleEscalation on the report-count check")
	}

	// Isolating case for the park-not-advance (non-terminal) check: marker present,
	// three recorded cycles, exactly one implementation report (so the cycle-count,
	// marker, and report-count checks ALL pass), but the FO terminalized the entity
	// to status: done — escalate-then-terminalize, auto-resolving instead of parking
	// for the human (the escalation prompt says "do not advance to done"). This must
	// still fail, and it fails ONLY on the non-terminal check.
	markerButTerminalized := strings.Replace(escalationEntity(), "status: validation", "status: done", 1) +
		"- Cycle 3: REJECTED — third consecutive rejection.\n" +
		escalationMarker + "\n"
	if !strings.Contains(markerButTerminalized, "status: done") {
		t.Fatal("non-terminal isolating case must carry status: done")
	}
	if got := len(implementationReport.FindAllString(markerButTerminalized, -1)); got != 1 {
		t.Fatalf("non-terminal isolating case must carry exactly one implementation report so only the non-terminal check rejects it, got %d", got)
	}
	if err := assertThirdCycleEscalation(markerButTerminalized); err == nil {
		t.Fatal("expected a marker-present but terminalized (status: done) end-state to fail assertThirdCycleEscalation on the park-not-advance check")
	}
}

func TestMergeHookGuardrailNegativeBypass(t *testing.T) {
	before := mergeHookGuardEntity()
	guardObserved := "Error: entity merge-check cannot advance to terminal - workflow has merge hook(s) [local-merge]"
	// The held-guard baseline passes: entity unmutated, still implementation, and
	// the observed output named the merge hook + terminal guard refusal.
	if err := assertMergeHookGuardHeld(before, before, guardObserved); err != nil {
		t.Fatalf("held merge-hook guard baseline must pass: %v", err)
	}

	// Merge-hook bypass: the FO terminalized the entity to done despite the
	// registered hook. The state assertion must catch the advanced status even if
	// the observed transcript still mentions a merge hook.
	bypassed := strings.Replace(before, "status: implementation", "status: done", 1)
	if bypassed == before {
		t.Fatal("merge-hook fixture must contain `status: implementation`")
	}
	if err := assertMergeHookGuardHeld(before, bypassed, guardObserved); err == nil {
		t.Fatal("expected a terminalized (status: done) entity to fail assertMergeHookGuardHeld even with merge-hook mention in the transcript")
	}

	// Bypass with no guard error in the observed output: the FO advanced and the
	// run never surfaced the guard refusal. Must fail on the missing guard signal.
	if err := assertMergeHookGuardHeld(before, before, "terminalized merge-check to done"); err == nil {
		t.Fatal("expected a run with no terminal-guard refusal in observed to fail assertMergeHookGuardHeld")
	}

	// Isolating case for the `cannot advance to terminal` guard-error check: the
	// observed output mentions a merge hook (so the mention check passes) but never
	// reports the terminal-guard refusal — the FO touched the hook yet the guard
	// never FIRED. This must still fail, and it fails ONLY on the guard-error
	// check. Without this case, deleting that check leaves the suite green: every
	// other merge-hook negative observed string also lacks the merge-hook mention,
	// so they trip the earlier mention check and never exercise the guard-error
	// clause. This is the one assertion proving the merge guard actually fired, so
	// it must be independently covered.
	hookMentionedNoGuard := "Inspected startup: workflow registers a merge hook [local-merge]. Proceeding without terminalization."
	if err := assertMergeHookGuardHeld(before, before, hookMentionedNoGuard); err == nil {
		t.Fatal("expected observed that mentions a merge hook but omits the terminal-guard refusal to fail assertMergeHookGuardHeld on the guard-error check")
	}
}
