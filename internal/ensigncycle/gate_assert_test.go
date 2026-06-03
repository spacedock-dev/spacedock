package ensigncycle

import "testing"

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

	if err := assertGateHeld(entity, entity, final); err != nil {
		t.Fatalf("gate-held assertion errored on held gate: %v", err)
	}

	t.Run("rejects_mutated_entity", func(t *testing.T) {
		after := entity + "\nFO edited this file.\n"
		if err := assertGateHeld(entity, after, final); err == nil {
			t.Fatal("expected mutation to fail the gate-held assertion")
		}
	})

	t.Run("rejects_advanced_status", func(t *testing.T) {
		after := "---\n" +
			"id: gate-check\n" +
			"title: Gate Check\n" +
			"status: done\n" +
			"completed:\n" +
			"verdict:\n" +
			"---\n"
		if err := assertGateHeld(after, after, final); err == nil {
			t.Fatal("expected status: done to fail the gate-held assertion")
		}
	})

	t.Run("rejects_set_verdict", func(t *testing.T) {
		after := "---\n" +
			"id: gate-check\n" +
			"title: Gate Check\n" +
			"status: review\n" +
			"completed:\n" +
			"verdict: passed\n" +
			"---\n"
		if err := assertGateHeld(after, after, final); err == nil {
			t.Fatal("expected set verdict to fail the gate-held assertion")
		}
	})

	t.Run("rejects_missing_gate_output", func(t *testing.T) {
		if err := assertGateHeld(entity, entity, "No work available."); err == nil {
			t.Fatal("expected missing gate output to fail the gate-held assertion")
		}
	})
}
