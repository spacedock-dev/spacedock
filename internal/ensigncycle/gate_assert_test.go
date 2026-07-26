package ensigncycle

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
)

var recordedGateStateHead = regexp.MustCompile(`(?:^|\n)state-head\t[0-9a-f]{40}(?:\s|$)`)

type recordedGateHost struct {
	extract                                                  func(string) string
	name, commit, review, decision, narration, failed, child string
}

func requireRecordedGate(t *testing.T, ok bool, format string, args ...any) {
	t.Helper()
	if !ok {
		t.Fatalf(format, args...)
	}
}

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

	before := recordedGateEntity()
	after := recordedGateHeldEntity(before, recordedGateBriefingID, recordedGateDigest)
	requireRecordedGate(t, assertGateHeld(before, after, recordedGateReview()) == nil, "held gate failed")
	decision := "Decision ask: approve, revise with a concrete finding, or hold for a named prerequisite?"
	for name, line := range map[string]string{
		"baseline":              decision,
		"markdown-action-verbs": "Decision: **approve** consumes the validation authorization and advances toward handoff; **reject** bounces it back to implementation with concrete asks.",
		"retained-claude":       "Decision: approve to consume this authorization and advance recorded-gate-task from validation into the handoff stage for dispatch.",
		"semantic-label":        "Choose approve to enter handoff, revise with findings, or hold for a prerequisite.",
	} {
		review := strings.Replace(recordedGateReview(), decision, line, 1)
		requireRecordedGate(t, assertConciseRecordedGateReview(review) == nil, "%s decision prompt failed", name)
	}
	for name, line := range map[string]string{
		"missing":     "",
		"disposition": "Decision: continue into handoff after review.",
		"vague":       "Decision: approve with confidence.",
		"narration":   "The package can advance after a decision.",
	} {
		review := strings.Replace(recordedGateReview(), decision, line, 1)
		requireRecordedGate(t, assertConciseRecordedGateReview(review) != nil, "%s decision control qualified", name)
	}
	for name, tc := range map[string]struct{ after, review string }{"unbound": {before, recordedGateReview()}, "advanced": {strings.Replace(after, "status: validation", "status: handoff", 1), recordedGateReview()}, "resolution": {strings.Replace(after, "\n---\n", "\nresolution:\n  type: Resolution\n---\n", 1), recordedGateReview()}, "verdict": {strings.Replace(after, "verdict:\n", "verdict: passed\n", 1), recordedGateReview()}, "review": {after, "Gate review: legacy\nDecision: approve?"}, "legacy": {entity, final}} {
		requireRecordedGate(t, assertGateHeld(before, tc.after, tc.review) != nil, "%s control qualified", name)
	}
}

func TestAssertGateHeldUsesDynamicPreparedBinding(t *testing.T) {
	fixture := writePreparedRecordedGateFixture(t)
	writeFile(t, fixture.gateReview, recordedGateSourceReview()+"\nDynamic prepared-binding sentinel.\n")
	gateReviewRel, err := filepath.Rel(fixture.stateRoot, fixture.gateReview)
	if err != nil {
		t.Fatal(err)
	}
	gitCommitPathScoped(t, fixture.stateRoot, gateReviewRel, "vary prepared source identity")

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
	briefingID := outputValue(prepared.stdout, "briefing")
	digest := outputValue(prepared.stdout, "digest")
	if briefingID == "" || digest == "" || digest == recordedGateDigest {
		t.Fatalf("prepare did not produce a distinct dynamic binding: briefing=%q digest=%q", briefingID, digest)
	}
	room := outputValue(prepared.stdout, "room")
	briefingBytes, err := os.ReadFile(filepath.Join(room, "gate-briefing.json"))
	if err != nil {
		t.Fatal(err)
	}
	if recomputed, err := gates.CanonicalDigest(briefingBytes); err != nil || recomputed != digest {
		t.Fatalf("prepared digest=%q recomputed=%q err=%v", digest, recomputed, err)
	}

	after := readFile(t, fixture.entity)
	review := recordedGateReviewWith(briefingID, digest)
	if err := assertGateHeld(before, after, review); err != nil {
		t.Fatalf("dynamic prepared binding rejected: %v", err)
	}
	for name, mutant := range map[string]string{
		"wrong-briefing": strings.Replace(review, briefingID, "briefing:wrong", 1),
		"wrong-digest":   strings.Replace(review, digest, "sha256:"+strings.Repeat("f", 64), 1),
	} {
		t.Run(name, func(t *testing.T) {
			if err := assertGateHeld(before, after, mutant); err == nil {
				t.Fatal("identity mutant graded PASS")
			}
		})
	}
}

func TestAssertGateHeldUsesCurrentRetainedGateBinding(t *testing.T) {
	before := recordedGateEntity()
	after := recordedGateHeldEntity(before, recordedGateBriefingID, recordedGateDigest)
	retained := "    - id: gate:recorded-gate-task:implementation\n" +
		"      stage: implementation\n" +
		"      attempts:\n" +
		"        - id: gate-attempt:recorded-gate-task-implementation-1\n" +
		"          briefing:\n" +
		"            id: briefing:recorded-gate-task:implementation:attempt-1:revision-1\n" +
		"            digest: sha256:" + strings.Repeat("1", 64) + "\n"
	after = strings.Replace(after, "  records:\n", "  records:\n"+retained, 1)

	if err := assertGateHeld(before, after, recordedGateReview()); err != nil {
		t.Fatalf("current validation binding was shadowed by an older retained gate: %v", err)
	}
}
