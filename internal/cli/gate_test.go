package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"github.com/spacedock-dev/spacedock/internal/status"
)

func TestGateRecordAndValidateCLILeaveStatusUntouched(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: ideation\n      initial: true\n---\n# Workflow\n")
	entity := filepath.Join(root, "task.md")
	writeFile(t, entity, "---\nstatus: ideation\ntitle: Task\n---\n# Task\n")
	briefing := filepath.Join(root, "briefing.json")
	writeFile(t, briefing, `{"type":"Briefing","version":"1","id":"briefing:provider","question":"Review task","artifacts":[{"id":"artifact:primary","uri":"artifact.md","rev":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}]}`)
	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--briefing", briefing}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("record exit=%d stderr=%q", code, errOut.String())
	}
	b, _ := os.ReadFile(entity)
	if !strings.Contains(string(b), "status: ideation") || strings.Contains(string(b), "application:") {
		t.Fatalf("recorder changed lifecycle/application state:\n%s", b)
	}
	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"gate", "validate", "task", "--workflow-dir", root}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 || !strings.Contains(out.String(), "state=open") || !strings.Contains(out.String(), "decision=") {
		t.Fatalf("validate exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
}

func TestGateRoundRecordAndValidateCLI(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: implementation\n      initial: true\n---\n# Workflow\n")
	entityDir := filepath.Join(root, "task")
	if err := os.MkdirAll(entityDir, 0o755); err != nil {
		t.Fatal(err)
	}
	entity := filepath.Join(entityDir, "index.md")
	writeFile(t, entity, "---\nid: task\nstatus: implementation\ntitle: Task\n---\n# Task\n")
	copyGateTestdata(t, filepath.Join(entityDir, "candidate.patch"), filepath.Join("advisory-round", "candidate.patch"))
	inputs := filepath.Join(entityDir, "inputs")
	if err := os.MkdirAll(inputs, 0o755); err != nil {
		t.Fatal(err)
	}
	briefing := filepath.Join(inputs, "briefing.json")
	log := filepath.Join(inputs, "briefing.review.jsonl")
	feedback := filepath.Join(inputs, "feedback-cycle.txt")
	copyGateTestdata(t, briefing, filepath.Join("advisory-round", "briefing.json"))
	copyGateTestdata(t, log, filepath.Join("advisory-round", "briefing.review.jsonl"))
	writeFile(t, feedback, "- Cycle 1: REJECTED — Roborev; surface 1/580 vs estimate 580 (100%); AC unchanged\n")
	var out, errOut bytes.Buffer
	invoke := func(args ...string) int {
		out.Reset()
		errOut.Reset()
		return run(context.Background(), args, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	}
	args := []string{"gate", "record", "task", "--workflow-dir", root, "--round", "implementation/1", "--briefing", briefing, "--log", log, "--feedback-cycle", feedback}
	if code := invoke(args...); code != 0 || !strings.Contains(out.String(), "round=round:task:implementation:1") || !strings.Contains(out.String(), "triage=all-declines") {
		t.Fatalf("round record exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if code := invoke("gate", "validate", "task", "--workflow-dir", root, "--round", "implementation/1"); code != 0 ||
		strings.Count(out.String(), "advisory=true") != 2 || !strings.Contains(out.String(), "annotation:job-592") {
		t.Fatalf("round validate exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
}

func TestGateImplicitWorkflowUsesOwningDefinitionFromNestedDirectory(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\ncommissioned-by: spacedock@0.1.0\nid-style: slug\nstages:\n  states:\n    - name: ideation\n      initial: true\n    - name: implementation\n---\n# Workflow\n")
	entity := filepath.Join(root, "task.md")
	writeFile(t, entity, "---\nstatus: ideation\ntitle: Task\n---\n# Task\n")
	briefing := filepath.Join(root, "review", "ideation", "briefing-1", "briefing.json")
	if err := os.MkdirAll(filepath.Dir(briefing), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, briefing, `{"type":"Briefing","version":"1","id":"briefing:task:1","question":"Review task","artifacts":[{"id":"artifact:primary","uri":"artifact.md","rev":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}]}`)
	nested := filepath.Join(root, "nested", "deeper")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "task", "--briefing", briefing}, nil, nested, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("implicit workflow record exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if _, _, err := gates.Read(entity); err != nil {
		t.Fatalf("owning workflow entity was not recorded: %v", err)
	}
}

func TestGatePresentationRemainsOutsideBinary(t *testing.T) {
	root := t.TempDir()
	before, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "review", "task"}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 2 || !strings.Contains(errOut.String(), "unknown subcommand (want: record|validate)") {
		t.Fatalf("gate review exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	after, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(before) != len(after) {
		t.Fatalf("rejected presentation verb changed working directory: before=%v after=%v", before, after)
	}
}

func TestGateRecordRejectsNonCanonicalBriefingBasenameBeforeMutation(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: ideation\n      initial: true\n---\n# Workflow\n")
	entity := filepath.Join(root, "task.md")
	writeFile(t, entity, "---\nstatus: ideation\ntitle: Task\n---\n# Task\n")
	room := filepath.Join(root, "review", "ideation", "briefing-1")
	if err := os.MkdirAll(room, 0o755); err != nil {
		t.Fatal(err)
	}
	noncanonical := filepath.Join(room, "revision-1.json")
	copyGateTestdata(t, noncanonical, "revision-18.json")
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--briefing", noncanonical}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 1 || !strings.Contains(errOut.String(), "named briefing.json") {
		t.Fatalf("noncanonical basename exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	after, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(after, before) {
		t.Fatal("noncanonical Briefing basename changed the entity")
	}
	if _, err := os.Stat(entity + ".gates.lock"); !os.IsNotExist(err) {
		t.Fatalf("noncanonical Briefing basename left lock residue: %v", err)
	}
}

func TestGateRecordBriefingReentersStageGateWhenAnotherGateIsSelected(t *testing.T) {
	root, entity, briefing := crossGateFixture(t)
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	beforeDoc, _, err := gates.Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	expectedIdeationClosure := beforeDoc.Records[0].Attempts[0]
	expectedIdeationClosure.Application.State = "superseded"
	beforeIdeationClosure, _ := json.Marshal(expectedIdeationClosure)
	beforeValidationClosure, _ := json.Marshal(beforeDoc.Records[1].Attempts[0])
	beforeOutsideGates := outsideFixtureGates(t, string(before))

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "durable-gate-approval-pending-blockers", "--workflow-dir", root, "--briefing", briefing}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("record --briefing exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	after, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	afterText := string(after)
	afterDoc, _, err := gates.Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	gotIdeationClosure, _ := json.Marshal(afterDoc.Records[0].Attempts[0])
	gotValidationClosure, _ := json.Marshal(afterDoc.Records[1].Attempts[0])
	if !bytes.Equal(gotIdeationClosure, beforeIdeationClosure) {
		t.Fatal("closed ideation attempt 9 changed beyond superseding its pending application")
	}
	if !bytes.Equal(gotValidationClosure, beforeValidationClosure) {
		t.Fatal("closed validation attempt changed during ideation re-entry")
	}
	if got := outsideFixtureGates(t, afterText); got != beforeOutsideGates {
		t.Fatal("record --briefing changed bytes outside gates")
	}
	if afterDoc.Current.Gate != "gate:docs-dev:3k:ideation" || len(afterDoc.Records[0].Attempts) != 2 {
		t.Fatalf("ideation successor not selected: %#v", afterDoc)
	}
	for _, want := range []string{
		"gate: gate:docs-dev:3k:ideation",
		"id: briefing:docs-dev:3k:ideation:attempt-10:revision-18",
		"digest: sha256:6b2c4f1388a58f42f7c8610f847ed9e7cce92758c00b201d4eb9f4f89dbedd8b",
	} {
		if !strings.Contains(afterText, want) {
			t.Fatalf("recorded entity missing %q:\n%s", want, afterText)
		}
	}
	for _, forbidden := range []string{"current-attempt:", "sequence:", "previous-attempt:"} {
		if strings.Contains(afterText, forbidden) {
			t.Fatalf("canonical writer emitted prototype field %q:\n%s", forbidden, afterText)
		}
	}
}

func TestGateRecordDecisionClosesMinimalBriefingAttempt(t *testing.T) {
	root, entity, briefing := crossGateFixture(t)

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "durable-gate-approval-pending-blockers", "--workflow-dir", root, "--briefing", briefing}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("record --briefing exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	created, _, err := gates.Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if got := created.Records[0].Attempts[1].Resolution; got != nil {
		t.Fatalf("semantic record did not create an open attempt: %#v", got)
	}

	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"gate", "record", "durable-gate-approval-pending-blockers", "--workflow-dir", root, "--decision", "approve", "--actor", "person:captain", "--reason", "lgtm. add the fixture"}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("decision close exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if !strings.Contains(out.String(), "attempt=gate-attempt:3k-ideation-10 state=closed") || !strings.Contains(out.String(), "decision=approve") {
		t.Fatalf("operation close stdout=%q, want closed approved attempt", out.String())
	}
}

func TestGateEligibilityAndConsumeAuthorizeOnce(t *testing.T) {
	root, entity, briefing := crossGateFixture(t)
	var out, errOut bytes.Buffer
	invoke := func(args ...string) int {
		out.Reset()
		errOut.Reset()
		return run(context.Background(), args, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	}
	if code := invoke("gate", "record", "durable-gate-approval-pending-blockers", "--workflow-dir", root, "--briefing", briefing); code != 0 {
		t.Fatalf("record briefing exit=%d stderr=%q", code, errOut.String())
	}
	if code := invoke("gate", "record", "durable-gate-approval-pending-blockers", "--workflow-dir", root, "--decision", "approve", "--actor", "person:captain"); code != 0 {
		t.Fatalf("record approval exit=%d stderr=%q", code, errOut.String())
	}
	if code := invoke("gate", "eligibility", "durable-gate-approval-pending-blockers", "--workflow-dir", root); code != 0 || !strings.Contains(out.String(), "application=advance/pending condition=approved-pending eligible=true target-stage=validation") {
		t.Fatalf("eligibility exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if code := invoke("gate", "consume", "durable-gate-approval-pending-blockers", "--workflow-dir", root); code != 0 || !strings.Contains(out.String(), "consumed=true") {
		t.Fatalf("consume exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	body, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "status: validation") || !strings.Contains(string(body), "state: consumed") {
		t.Fatalf("consume did not co-write status and application:\n%s", body)
	}
	if code := invoke("gate", "consume", "durable-gate-approval-pending-blockers", "--workflow-dir", root); code != 1 || !strings.Contains(out.String(), "condition=consumed") || !strings.Contains(out.String(), "consumed=false") {
		t.Fatalf("repeat consume exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
}

func TestGateRecordConsumesExactResultOnlyWithCompleteAssociation(t *testing.T) {
	root, entity, briefing := unboundSemanticDecisionFixture(t)
	result := filepath.Join(root, "result.json")
	association := filepath.Join(root, "association.json")
	copyGateTestdata(t, result, "exact-review-v1-result.json")
	copyGateTestdata(t, association, "exact-review-v1-association.json")

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--briefing", briefing}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("bind retained briefing.json exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	bound, _, err := gates.Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if got := bound.Records[0].Attempts[0].Briefing.RoomRef; got != "./review/validation/briefing-1" {
		t.Fatalf("accepted binding room-ref = %q", got)
	}
	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--result", result, "--actor", "person:reviewer", "--adoption-note", "Adopted by the captain"}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 2 || !strings.Contains(errOut.String(), "--association") {
		t.Fatalf("result without association exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	primaryOnlyPath := filepath.Join(root, "primary-only-association.json")
	copyGateTestdata(t, primaryOnlyPath, "exact-review-v1-association-truncated.json")
	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--result", result, "--association", primaryOnlyPath, "--actor", "person:reviewer", "--adoption-note", "Adopted by the captain"}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 1 || !strings.Contains(errOut.String(), "complete presentation mapping") {
		t.Fatalf("primary-only association exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	afterFailed, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(afterFailed, before) {
		t.Fatal("rejected primary-only association changed the entity")
	}

	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--result", result, "--association", association, "--actor", "person:reviewer", "--adoption-note", "Adopted by the captain"}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("record result exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	body, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"id: resolution:actor-1784675152206198000", "briefing: briefing:docs-dev:3k:validation:attempt-1:revision-1", "decision: revise", "adoption-note: Adopted by the captain"} {
		if !strings.Contains(string(body), want) {
			t.Fatalf("recorded Result missing %q:\n%s", want, body)
		}
	}
	doc, _, err := gates.Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	resolutionBytes, err := json.Marshal(doc.Records[0].Attempts[0].Resolution)
	if err != nil {
		t.Fatal(err)
	}
	var portable map[string]any
	if err := json.Unmarshal(resolutionBytes, &portable); err != nil {
		t.Fatal(err)
	}
	for _, wrapperField := range []string{"status", "binding", "actor", "approver", "resolutionId"} {
		if _, ok := portable[wrapperField]; ok {
			t.Fatalf("provider wrapper field %q crossed the portable boundary: %s", wrapperField, resolutionBytes)
		}
	}
}

func unboundSemanticDecisionFixture(t *testing.T) (root, entity, briefing string) {
	t.Helper()
	root = t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: validation\n      initial: true\n    - name: done\n      terminal: true\n---\n# Workflow\n")
	entity = filepath.Join(root, "task.md")
	writeFile(t, entity, "---\nstatus: validation\ntitle: Task\n---\n# Task\n")
	briefing = filepath.Join(root, "review", "validation", "briefing-1", "briefing.json")
	if err := os.MkdirAll(filepath.Dir(briefing), 0o755); err != nil {
		t.Fatal(err)
	}
	copyGateTestdata(t, briefing, "exact-validation-briefing.json")
	return root, entity, briefing
}

func TestGateRecordChatDecisionAndRejectsOperationInterface(t *testing.T) {
	root, entity := semanticDecisionFixture(t)
	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--decision", "approve", "--actor", "agent:first-officer", "--reason", "All retained ACs reproduced", "--directive", "Captain: approve after the reset"}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("record decision exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	body, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"by: agent:first-officer", "decision: approve", "reason: All retained ACs reproduced", "adoption-note: 'Captain: approve after the reset'"} {
		if !strings.Contains(string(body), want) {
			t.Fatalf("chat Resolution missing %q:\n%s", want, body)
		}
	}

	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--operation", "old.yml"}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 2 || !strings.Contains(errOut.String(), "unknown gate flag: --operation") {
		t.Fatalf("legacy operation exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
}

func semanticDecisionFixture(t *testing.T) (root, entity string) {
	t.Helper()
	root = t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: validation\n      initial: true\n    - name: done\n      terminal: true\n---\n# Workflow\n")
	entity = filepath.Join(root, "task.md")
	writeFile(t, entity, "---\nstatus: validation\ngates:\n  version: 1\n  current:\n    gate: gate:docs-dev:3k:validation\n  records:\n    - id: gate:docs-dev:3k:validation\n      stage: validation\n      attempts:\n        - id: gate-attempt:3k-validation-1\n          briefing:\n            id: briefing:docs-dev:3k:validation:attempt-1:revision-1\n            digest: sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac\n            digest-domain: canonical-bytes\n            room-ref: ./review/validation/briefing-1\ntitle: Task\n---\n# Task\n")
	briefing := filepath.Join(root, "review", "validation", "briefing-1", "briefing.json")
	if err := os.MkdirAll(filepath.Dir(briefing), 0o755); err != nil {
		t.Fatal(err)
	}
	copyGateTestdata(t, briefing, "exact-validation-briefing.json")
	return root, entity
}

func copyGateTestdata(t *testing.T, destination, name string) {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("..", "gates", "testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destination, body, 0o644); err != nil {
		t.Fatal(err)
	}
}

func crossGateFixture(t *testing.T) (root, entity, briefing string) {
	t.Helper()
	root = t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: ideation\n      initial: true\n    - name: validation\n---\n# Workflow\n")
	entity = filepath.Join(root, "durable-gate-approval-pending-blockers.md")
	digest := func(char string) string { return "sha256:" + strings.Repeat(char, 64) }
	writeFile(t, entity, "---\nstatus: ideation\ngates:\n  version: 1\n  current:\n    gate: gate:docs-dev:3k:validation\n  records:\n    - id: gate:docs-dev:3k:ideation\n      stage: ideation\n      attempts:\n        - id: gate-attempt:3k-ideation-9\n          briefing:\n            id: briefing:ideation:9\n            digest: "+digest("1")+"\n            digest-domain: raw-file-pin\n            room-ref: ./review/ideation/9\n          resolution:\n            type: Resolution\n            id: resolution:ideation:9\n            briefing: briefing:ideation:9\n            by: person:captain\n            at: 2026-07-22T00:00:00Z\n            decision: approve\n          application:\n            action: advance\n            target-stage: validation\n            state: pending\n            blockers: [{id: blocker:preserve-me, state: unsatisfied}]\n    - id: gate:docs-dev:3k:validation\n      stage: validation\n      attempts:\n        - id: gate-attempt:3k-validation-1\n          briefing:\n            id: briefing:validation:1\n            digest: "+digest("2")+"\n            digest-domain: raw-file-pin\n            room-ref: ./review/validation/1\n          resolution:\n            type: Resolution\n            id: resolution:validation:1\n            briefing: briefing:validation:1\n            by: person:captain\n            at: 2026-07-22T00:00:00Z\n            decision: revise\n            reason: Re-enter ideation.\nsprint: durable-decisions\ntitle: Task\n---\n# Task\n")
	briefing = filepath.Join(root, "review", "ideation", "briefing-18", "briefing.json")
	if err := os.MkdirAll(filepath.Dir(briefing), 0o755); err != nil {
		t.Fatal(err)
	}
	briefingBytes, err := os.ReadFile(filepath.Join("..", "gates", "testdata", "revision-18.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(briefing, briefingBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	return root, entity, briefing
}

func fixtureSection(t *testing.T, body, start, end string) string {
	t.Helper()
	startAt := strings.Index(body, start)
	if startAt < 0 {
		t.Fatalf("fixture section start %q not found", start)
	}
	endAt := strings.Index(body[startAt+len(start):], end)
	if endAt < 0 {
		t.Fatalf("fixture section end %q not found", end)
	}
	return body[startAt : startAt+len(start)+endAt]
}

func outsideFixtureGates(t *testing.T, body string) string {
	t.Helper()
	start := strings.Index(body, "gates:\n")
	end := strings.Index(body, "sprint: durable-decisions\n")
	if start < 0 || end < 0 || end <= start {
		t.Fatal("fixture gates boundaries not found")
	}
	return body[:start] + body[end:]
}
