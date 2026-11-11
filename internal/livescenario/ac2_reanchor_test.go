// ABOUTME: AC re-anchor scenario tests pin the durable gate-decision oracle.
// ABOUTME: Correct, wrong, unchanged, and stable-identity branches run offline.
package livescenario

import (
	"fmt"
	"os"
	"strings"
	"testing"

	statuspkg "github.com/spacedock-dev/spacedock/internal/status"
)

func acReanchorDecisionBody(status, decision string) string {
	reason, application := ", reason: measurements regressed", ""
	if decision == "approve" {
		reason, application = "", "\n          application: {target-stage: accepted, state: pending}"
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
          briefing: {id: "briefing:validation:1", digest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", room-ref: review/validation/1}
          resolution: {type: Resolution, id: "resolution:validation:1", briefing: "briefing:validation:1", by: "person:captain", at: "2026-08-03T00:00:00Z", decision: %s%s}%s
---
# Contract Measurement
`, status, decision, reason, application)
}

func assertACReanchorBody(t *testing.T, body string) error {
	t.Helper()
	root := t.TempDir()
	sc := AuthorACReanchorScenario()
	entity, err := sc.Setup(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(entity, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return sc.Assert(EntityState{}, EntityState{Body: body}, "Decision: approve")
}

func TestACReanchorDurableDecisionBranches(t *testing.T) {
	correct := acReanchorDecisionBody("validation", "revise")
	if err := assertACReanchorBody(t, correct); err != nil {
		t.Fatalf("closed revise held at validation should pass independently of narration: %v", err)
	}

	wrong := map[string]string{
		"approve decision": acReanchorDecisionBody("rework", "approve"),
		"rework status":    acReanchorDecisionBody("rework", "revise"),
		"accepted target":  acReanchorDecisionBody("accepted", "revise"),
		"malformed gate":   strings.Replace(correct, "type: Resolution, ", "", 1),
		"later approve": strings.Replace(correct, "---\n# Contract", `        - id: gate-attempt:ac2-design-proof-validation-2
          briefing: {id: "briefing:validation:2", digest: "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", room-ref: review/validation/2}
          resolution: {type: Resolution, id: "resolution:validation:2", briefing: "briefing:validation:2", by: "person:captain", at: "2026-08-03T01:00:00Z", decision: approve}
          application: {target-stage: accepted, state: pending}
---
# Contract`, 1),
		"duplicate validation": strings.Replace(correct, "---\n# Contract", `    - id: gate:ac2-design-proof:validation-duplicate
      stage: validation
      attempts:
        - id: gate-attempt:ac2-design-proof-validation-2
          briefing: {id: "briefing:validation:2", digest: "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", room-ref: review/validation/2}
          resolution: {type: Resolution, id: "resolution:validation:2", briefing: "briefing:validation:2", by: "person:captain", at: "2026-08-03T01:00:00Z", decision: approve}
          application: {target-stage: accepted, state: pending}
---
# Contract`, 1),
	}
	for name, body := range wrong {
		t.Run(name, func(t *testing.T) {
			if err := assertACReanchorBody(t, body); err == nil {
				t.Fatal("wrong durable decision branch unexpectedly passed")
			}
		})
	}
}

func TestACReanchorRejectsUnchangedNarratedResult(t *testing.T) {
	if err := assertACReanchorBody(t, acReanchorEntity); err == nil {
		t.Fatal("unchanged durable state unexpectedly passed on narrated rejection")
	}
}

func TestACReanchorFixtureIdentity(t *testing.T) {
	if got, want := AuthorACReanchorScenario().Name, "ac-reanchor/means-pass-value-regressed"; got != want {
		t.Fatalf("scenario identity = %q, want %q", got, want)
	}
}

func TestACReanchorFixtureIsDiscoverable(t *testing.T) {
	root := t.TempDir()
	if _, err := writeACReanchorFixture(root); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"ac2-design-proof-review.md", "ac2-design-proof-reference.md"} {
		if _, err := os.Stat(root + "/" + name); err != nil {
			t.Fatalf("AC re-anchor fixture missing committed gate package file %q: %v", name, err)
		}
	}
	if got, found := statuspkg.DiscoverWorkflowDir(root); !found {
		t.Fatal("AC re-anchor fixture is not discoverable from its workflow root")
	} else if got != root {
		t.Fatalf("DiscoverWorkflowDir(%q) = %q, want %q", root, got, root)
	}
}
