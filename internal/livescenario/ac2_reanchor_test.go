// ABOUTME: AC re-anchor scenario tests pin the durable gate-decision oracle.
// ABOUTME: Correct, wrong, unchanged, and stable-identity branches run offline.
package livescenario

import (
	"fmt"
	"testing"
)

func acReanchorDecisionBody(status, decision, action, target string) string {
	application := ""
	if action != "" {
		application = fmt.Sprintf("          application:\n            action: %s\n            target-stage: %s\n            state: pending\n", action, target)
	}
	return fmt.Sprintf(`---
id: ac2-design-proof
status: %s
gates:
  version: 1
  records:
    - id: gate:ac2-design-proof:validation
      stage: validation
      attempts:
        - id: gate-attempt:ac2-design-proof-validation-1
          resolution:
            decision: %s
%s
---
# Contract Measurement
`, status, decision, application)
}

func TestACReanchorDurableDecisionBranches(t *testing.T) {
	sc := AuthorACReanchorScenario()
	correct := acReanchorDecisionBody("rework", "revise", "", "")
	if err := sc.Assert(EntityState{}, EntityState{Body: correct}, "Decision: approve"); err != nil {
		t.Fatalf("revise/feedback/rework durable branch should pass independently of narration: %v", err)
	}

	wrong := map[string]string{
		"approve decision": acReanchorDecisionBody("rework", "approve", "feedback", "rework"),
		"advance action":   acReanchorDecisionBody("rework", "revise", "advance", "rework"),
		"accepted target":  acReanchorDecisionBody("accepted", "revise", "feedback", "accepted"),
	}
	for name, body := range wrong {
		t.Run(name, func(t *testing.T) {
			if err := sc.Assert(EntityState{}, EntityState{Body: body}, "Decision: REJECT"); err == nil {
				t.Fatal("wrong durable decision branch unexpectedly passed")
			}
		})
	}
}

func TestACReanchorRejectsUnchangedNarratedResult(t *testing.T) {
	sc := AuthorACReanchorScenario()
	unchanged := EntityState{Body: acReanchorEntity}
	if err := sc.Assert(unchanged, unchanged, "Gate review: REJECT because the end value regressed."); err == nil {
		t.Fatal("unchanged durable state unexpectedly passed on narrated rejection")
	}
}

func TestACReanchorFixtureIdentity(t *testing.T) {
	if got, want := AuthorACReanchorScenario().Name, "ac-reanchor/means-pass-value-regressed"; got != want {
		t.Fatalf("scenario identity = %q, want %q", got, want)
	}
}
