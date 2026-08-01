// ABOUTME: Pi-subagents stage dispatch wrapper invariants.
// ABOUTME: Proves the wrapper adds only transport fields around canonical assignments.
package piruntime

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSubagentStageDispatchAddsOnlyPiTransportFields(t *testing.T) {
	assignment := "Read /tmp/spacedock-dispatches/spacedock-ensign-fixture-implementation.md and treat its content as your assignment."

	wrapped := SubagentStageDispatch(assignment, "implementation", "fixture implementation")
	if wrapped.Context != "fresh" {
		t.Fatalf("context = %q, want fresh", wrapped.Context)
	}
	if wrapped.Task != assignment {
		t.Fatalf("wrapper rewrote assignment:\nwant: %s\n got: %s", assignment, wrapped.Task)
	}
	if wrapped.Phase != "implementation" || wrapped.Label != "fixture implementation" {
		t.Fatalf("phase/label not preserved: %#v", wrapped)
	}

	payload, err := json.Marshal(wrapped)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(payload), "acceptance") {
		t.Fatalf("stage wrapper must not contain same-agent acceptance contract: %s", payload)
	}
}

func TestSubagentDispatchSpawnFieldsRoundTrip(t *testing.T) {
	// Default ensign dispatch: the build artifact's agent/skill bind the spawn.
	wrapped := SubagentStageDispatch("assignment", "implementation", "label")
	wrapped.Agent = "worker"
	wrapped.Skill = "ensign"
	payload, err := json.Marshal(wrapped)
	if err != nil {
		t.Fatal(err)
	}
	var back SubagentDispatch
	if err := json.Unmarshal(payload, &back); err != nil {
		t.Fatal(err)
	}
	if back.Agent != "worker" || back.Skill != "ensign" {
		t.Fatalf("agent/skill did not round-trip: %#v from %s", back, payload)
	}

	// Stage agent override: skill is omitted from the JSON entirely.
	override := SubagentStageDispatch("assignment", "implementation", "label")
	override.Agent = "custom:agent"
	payload, err = json.Marshal(override)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(payload), `"agent":"custom:agent"`) {
		t.Fatalf("override agent missing from payload: %s", payload)
	}
	if strings.Contains(string(payload), `"skill"`) {
		t.Fatalf("skill must be omitted on an agent-override dispatch: %s", payload)
	}
}
