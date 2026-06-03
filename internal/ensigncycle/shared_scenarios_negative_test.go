package ensigncycle

import (
	"strings"
	"testing"
)

// AC-5: each shared scenario's assertion is behavior/state oriented, not a
// transcript-shape tautology. For every shared scenario these cases build the
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

func TestRejectionFlowNegativeMissingRoute(t *testing.T) {
	rejectedObserved := "validation was REJECTED; routing the finding back to implementation"

	// Missing rejection route: the FO never applied the fix nor routed the entity
	// back. The fixture entity as-staged is still at status: validation with only
	// the original implementation report — exactly the not-routed-back state. (The
	// marker substring appears once, quoted inside the REJECTED finding; the
	// assertion is not satisfied by that quote because it also requires a SECOND
	// implementation report and status: implementation, which the un-routed state
	// lacks — so a tautological marker-only assertion would falsely pass here while
	// the behavior-oriented one correctly fails.)
	notRouted := rejectionEntity()
	if !strings.Contains(notRouted, "status: validation") {
		t.Fatal("rejection fixture must start at status: validation")
	}
	if strings.Count(notRouted, "## Stage Report: implementation") != 1 {
		t.Fatal("un-routed rejection fixture must carry exactly one implementation report")
	}
	if err := assertRejectionFlow(notRouted, rejectedObserved); err == nil {
		t.Fatal("expected an un-routed rejection (still at validation, only one implementation report) to fail assertRejectionFlow")
	}

	// A partial route — fix marker applied and a second implementation report, but
	// the FO left status at validation (forgot to route the frontmatter back) —
	// must still fail on the status check, not pass on transcript shape alone.
	partialBody := rejectionEntity() +
		"\n" + rejectionFixMarker + "\n\n## Stage Report: implementation\n\n- DONE: applied fix\n"
	if strings.Count(partialBody, "## Stage Report: implementation") < 2 {
		t.Fatal("partial-route body must carry a second implementation report")
	}
	if err := assertRejectionFlow(partialBody, rejectedObserved); err == nil {
		t.Fatal("expected a fix applied but status left at validation to fail assertRejectionFlow")
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
}
