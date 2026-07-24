package ensigncycle

import (
	"strings"
	"testing"
)

func TestAssertGateHeld(t *testing.T) {
	before := recordedGateEntity()
	after := before + "\ngates:\n  records:\n    - id: gate:docs-dev:3k:validation\n      attempts:\n        - id: gate-attempt:3k-validation-1\n          state: open\n          briefing:\n            id: " + recordedGateBriefingID + "\n            digest: " + recordedGateDigest + "\n"
	if err := assertGateHeld(before, after, recordedGateReview()); err != nil {
		t.Fatalf("gate-held assertion errored on held gate: %v", err)
	}

	t.Run("rejects_unbound_entity", func(t *testing.T) {
		if err := assertGateHeld(before, before, recordedGateReview()); err == nil {
			t.Fatal("expected byte-identical unbound state to fail")
		}
	})

	t.Run("rejects_advanced_status", func(t *testing.T) {
		advanced := strings.Replace(after, "status: validation", "status: handoff", 1)
		if err := assertGateHeld(before, advanced, recordedGateReview()); err == nil {
			t.Fatal("expected status: done to fail the gate-held assertion")
		}
	})

	t.Run("rejects_resolution", func(t *testing.T) {
		closed := after + "\ntype: Resolution\n"
		if err := assertGateHeld(before, closed, recordedGateReview()); err == nil {
			t.Fatal("expected a Resolution to fail the open-bound assertion")
		}
	})

	t.Run("rejects_set_verdict", func(t *testing.T) {
		selfApproved := strings.Replace(after, "verdict:\n", "verdict: passed\n", 1)
		if err := assertGateHeld(before, selfApproved, recordedGateReview()); err == nil {
			t.Fatal("expected set verdict to fail the gate-held assertion")
		}
	})

	t.Run("rejects_noncanonical_output", func(t *testing.T) {
		if err := assertGateHeld(before, after, "Gate review: legacy\nDecision: approve?"); err == nil {
			t.Fatal("expected legacy marker output to fail")
		}
	})
}
