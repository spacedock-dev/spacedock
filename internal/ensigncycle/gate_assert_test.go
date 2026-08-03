package ensigncycle

import (
	"strings"
	"testing"
)

func requireRecordedGate(t *testing.T, ok bool, format string, args ...any) {
	t.Helper()
	if !ok {
		t.Fatalf(format, args...)
	}
}

func staticGateHeldExpectation() gateHeldExpectation {
	return gateHeldExpectation{
		gateID: "gate:docs-dev:3k:validation", attemptID: "gate-attempt:3k-validation-1",
		briefingID: recordedGateBriefingID, digest: recordedGateDigest,
	}
}

func recordedGateHeldEntity() string {
	gates := "gates:\n" +
		"  version: 1\n" +
		"  records:\n" +
		"    - id: gate:docs-dev:3k:validation\n" +
		"      stage: validation\n" +
		"      attempts:\n" +
		"        - id: gate-attempt:3k-validation-1\n" +
		"          briefing:\n" +
		"            id: " + recordedGateBriefingID + "\n" +
		"            digest: " + recordedGateDigest + "\n" +
		"            room-ref: rooms/validation/attempt-1/revision-1\n"
	return strings.Replace(recordedGateEntity(), "---\n# Recorded Gate Task", gates+"---\n# Recorded Gate Task", 1)
}

func TestAssertGateHeld(t *testing.T) {
	before := recordedGateEntity()
	after := recordedGateHeldEntity()
	expected := staticGateHeldExpectation()
	requireRecordedGate(t, assertGateHeld(before, after, expected) == nil, "held gate failed")
	resolved := strings.Replace(after, "          briefing:\n", "          resolution:\n            type: Resolution\n            id: resolution:docs-dev:3k:validation:1\n            briefing: "+recordedGateBriefingID+"\n            by: captain\n            at: 2026-07-28T00:00:00Z\n            decision: approve\n          briefing:\n", 1)
	applied := strings.Replace(resolved, "          briefing:\n", "          application:\n            action: advance\n            target-stage: handoff\n            state: pending\n          briefing:\n", 1)
	applied = strings.Replace(applied, "            action: advance\n", "", 1)
	for name, mutant := range map[string]string{
		"unbound":         before,
		"advanced":        strings.Replace(after, "status: validation", "status: handoff", 1),
		"malformed-state": strings.Replace(after, "          briefing:\n", "          state: open\n          briefing:\n", 1),
		"wrong-gate":      strings.ReplaceAll(after, "gate:docs-dev:3k:validation", "gate:docs-dev:wrong"),
		"wrong-briefing":  strings.Replace(after, recordedGateBriefingID, "briefing:docs-dev:wrong", 1),
		"resolved":        resolved,
		"applied":         applied,
		"verdict":         strings.Replace(after, "verdict:\n", "verdict: passed\n", 1),
	} {
		requireRecordedGate(t, assertGateHeld(before, mutant, expected) != nil, "%s control qualified", name)
	}
}

func TestAssertGateHeldAcceptsPreparedFixtureBinding(t *testing.T) {
	fixture := writePreparedRecordedGateFixture(t)
	before := readFile(t, fixture.entity)
	args := []string{
		"gate", "prepare", "recorded-gate-task",
		"--question", "Should the recorded validation gate advance?",
		"--artifact", fixture.gateReview,
		"--summary", "Exact recorded gate validation summary.",
	}
	for _, reference := range fixture.references {
		args = append(args, "--reference", reference)
	}
	args = append(args, "--workflow-dir", fixture.root)
	prepared := runRecordedGateCommand(buildRecordedGateBinary(t), fixture.root, "prepare", args...)
	if prepared.exit != 0 {
		t.Fatalf("prepare exit=%d\nstdout=%s\nstderr=%s", prepared.exit, prepared.stdout, prepared.stderr)
	}

	expected, err := recordedGateHeldExpectation(fixture)
	if err != nil {
		t.Fatal(err)
	}
	after := readFile(t, fixture.entity)
	if err := assertGateHeld(before, after, expected); err != nil {
		t.Fatalf("prepared fixture binding rejected: %v", err)
	}
	for name, mutant := range map[string]string{
		"attempt": strings.Replace(after, expected.attemptID, "gate-attempt:wrong", 1),
		"digest":  strings.Replace(after, expected.digest, "sha256:"+strings.Repeat("f", 64), 1),
	} {
		t.Run(name, func(t *testing.T) {
			if err := assertGateHeld(before, mutant, expected); err == nil {
				t.Fatal("prepared binding mutant graded PASS")
			}
		})
	}
}
