package ensigncycle

import "testing"

func TestAssertRejectionFlow(t *testing.T) {
	entity := "---\nstatus: implementation\n---\n" +
		rejectionFixMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: Initial implementation\n\n" +
		"## Stage Report: implementation\n\n- DONE: Applied rejection fix\n"
	observed := "validation was REJECTED; routed follow-up to implementation"

	if err := assertRejectionFlow(entity, observed); err != nil {
		t.Fatalf("expected rejection flow to pass: %v", err)
	}
	if err := assertRejectionFlow("## Stage Report: implementation\n", observed); err == nil {
		t.Fatal("expected missing fix marker to fail")
	}
	if err := assertRejectionFlow("---\nstatus: validation\n---\n"+rejectionFixMarker+"\n\n## Stage Report: implementation\n", observed); err == nil {
		t.Fatal("expected missing implementation status to fail")
	}
	if err := assertRejectionFlow("---\nstatus: implementation\n---\n"+rejectionFixMarker+"\n\n## Stage Report: implementation\n", observed); err == nil {
		t.Fatal("expected missing follow-up implementation report to fail")
	}
	if err := assertRejectionFlow(entity, "all quiet"); err == nil {
		t.Fatal("expected missing rejection output to fail")
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
