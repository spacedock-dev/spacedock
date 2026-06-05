package ensigncycle

import (
	"strings"
	"testing"
)

// AC-5 offline negative case: assertAutoContinue is behavior/state oriented, not
// a transcript-shape tautology. These cases build the SPECIFIC broken end-state
// the regression guards against — an FO that filed the implementation report and
// STOPPED, leaving the entity at status: implementation with no validation report
// — and prove the assertion reds even when the transcript narrates intent to
// continue. They are offline (default tag): the assertion is a pure function over
// entity-state + observed strings, so they spend no model.

func TestAutoContinueNegativeStoppedAfterImplementation(t *testing.T) {
	before := autoContinueEntity()

	// Baseline: an FO that truly advanced — status moved to validation and a
	// validation stage report was appended — GRADES PASS. Without this the negative
	// cases could pass against an assertion that always errors.
	advanced := strings.Replace(before, "status: implementation", "status: validation", 1) +
		"\n## Stage Report: validation\n\n- DONE: Verified the deliverable\n  PASSED.\n"
	if advanced == before {
		t.Fatal("fixture must contain `status: implementation` to advance")
	}
	if err := assertAutoContinue(before, advanced, "Advanced to validation and dispatched a fresh validator; gate presented."); err != nil {
		t.Fatalf("a truly-advanced end-state must grade PASS, got: %v", err)
	}

	// Broken: the FO stopped after filing the implementation report. The durable
	// state is byte-identical to the staged fixture (still status: implementation,
	// no validation report). Even WITH a transcript that narrates advancing to
	// validation, the state-oriented grade must catch the stop.
	stoppedObserved := "Implementation complete. Advancing to validation and dispatching a fresh validator."
	if err := assertAutoContinue(before, before, stoppedObserved); err == nil {
		t.Fatal("an FO that left the entity at status: implementation with no validation report must RED the grade even with a transcript narrating advancement")
	}

	// Broken: status advanced to validation in the frontmatter but NO validation
	// stage report was produced — a partial move (the FO bumped status but never
	// dispatched/ran the validator). This must still fail on the missing validation
	// report, not pass on the status bump alone.
	statusBumpedNoReport := strings.Replace(before, "status: implementation", "status: validation", 1)
	if statusBumpedNoReport == before {
		t.Fatal("fixture must contain `status: implementation` to bump")
	}
	if err := assertAutoContinue(before, statusBumpedNoReport, "advancing to validation"); err == nil {
		t.Fatal("a status bump with no validation stage report must RED the grade (the validator never ran)")
	}

	// Broken: a validation report was appended but the status was left at
	// implementation — the inverse partial move. Must fail on the status check.
	reportNoStatus := before + "\n## Stage Report: validation\n\n- DONE: Verified\n  PASSED.\n"
	if err := assertAutoContinue(before, reportNoStatus, "validated the deliverable"); err == nil {
		t.Fatal("a validation report with status left at implementation must RED the grade")
	}
}
