package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

func TestGateRecordBriefingReentersStageGateWhenAnotherGateIsSelected(t *testing.T) {
	root, entity, briefing := crossGateFixture(t)
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	beforeIdeationClosure := fixtureSection(t, string(before), "            - id: gate-attempt:3k-ideation-9\n", "        - id: gate:docs-dev:3k:validation\n")
	beforeValidationClosure := fixtureSection(t, string(before), "            - id: gate-attempt:3k-validation-1\n", "sprint: durable-decisions\n")
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
	if got := fixtureSection(t, afterText, "            - id: gate-attempt:3k-ideation-9\n", "        - id: gate:docs-dev:3k:validation\n"); !strings.HasPrefix(got, beforeIdeationClosure) {
		t.Fatal("closed ideation attempt 9 changed while appending attempt 10")
	}
	if got := fixtureSection(t, afterText, "            - id: gate-attempt:3k-validation-1\n", "sprint: durable-decisions\n"); got != beforeValidationClosure {
		t.Fatal("closed validation attempt changed during ideation re-entry")
	}
	if got := outsideFixtureGates(t, afterText); got != beforeOutsideGates {
		t.Fatal("record --briefing changed bytes outside gates")
	}
	ideation := fixtureSection(t, afterText, "        - id: gate:docs-dev:3k:ideation\n", "        - id: gate:docs-dev:3k:validation\n")
	if strings.Count(ideation, "            - id: gate-attempt:3k-ideation-") != 10 {
		t.Fatalf("ideation attempts did not grow to 10:\n%s", ideation)
	}
	for _, want := range []string{
		"gate: gate:docs-dev:3k:ideation",
		"attempt: gate-attempt:3k-ideation-10",
		"current-attempt: gate-attempt:3k-ideation-10",
		"id: briefing:docs-dev:3k:ideation:attempt-10:revision-18",
		"digest: sha256:6b2c4f1388a58f42f7c8610f847ed9e7cce92758c00b201d4eb9f4f89dbedd8b",
	} {
		if !strings.Contains(afterText, want) {
			t.Fatalf("recorded entity missing %q:\n%s", want, afterText)
		}
	}
	newAttempt := fixtureSection(t, afterText, "            - id: gate-attempt:3k-ideation-10\n", "        - id: gate:docs-dev:3k:validation\n")
	for _, forbidden := range []string{"sequence:", "previous-attempt:", "state:"} {
		if strings.Contains(newAttempt, forbidden) {
			t.Fatalf("minimal attempt minted %q:\n%s", forbidden, newAttempt)
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
	created, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	minimal := fixtureSection(t, string(created), "            - id: gate-attempt:3k-ideation-10\n", "        - id: gate:docs-dev:3k:validation\n")
	if strings.Contains(minimal, "state:") || strings.Contains(minimal, "resolution:") {
		t.Fatalf("semantic record did not create a minimal open attempt:\n%s", minimal)
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

func TestGateRecordConsumesExactResultOnlyWithCompleteAssociation(t *testing.T) {
	root, entity := semanticDecisionFixture(t)
	result := filepath.Join(root, "result.json")
	association := filepath.Join(root, "association.json")
	copyGateTestdata(t, result, "exact-review-v1-result.json")
	copyGateTestdata(t, association, "exact-review-v1-association.json")

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--result", result, "--actor", "person:reviewer", "--adoption-note", "Adopted by the captain"}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 2 || !strings.Contains(errOut.String(), "--association") {
		t.Fatalf("result without association exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	associationBytes, err := os.ReadFile(association)
	if err != nil {
		t.Fatal(err)
	}
	var primaryOnly map[string]any
	if err := json.Unmarshal(associationBytes, &primaryOnly); err != nil {
		t.Fatal(err)
	}
	presentation := primaryOnly["presentation"].([]any)
	primaryOnly["presentation"] = presentation[:1]
	primaryOnlyPath := filepath.Join(root, "primary-only-association.json")
	primaryOnlyBytes, err := json.Marshal(primaryOnly)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(primaryOnlyPath, primaryOnlyBytes, 0o644); err != nil {
		t.Fatal(err)
	}
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
	resolution := fixtureSection(t, string(body), "          resolution:\n", "title: Task\n")
	for _, wrapperField := range []string{"status:", "binding:", "actor:", "approver:", "resolutionId:"} {
		if strings.Contains(resolution, wrapperField) {
			t.Fatalf("provider wrapper field %q crossed the portable boundary:\n%s", wrapperField, resolution)
		}
	}
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
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: validation\n      initial: true\n---\n# Workflow\n")
	entity = filepath.Join(root, "task.md")
	writeFile(t, entity, "---\nstatus: validation\ngates:\n  version: 1\n  current:\n    gate: gate:docs-dev:3k:validation\n  records:\n    - id: gate:docs-dev:3k:validation\n      stage: validation\n      attempts:\n        - id: gate-attempt:3k-validation-1\n          briefing:\n            id: briefing:docs-dev:3k:validation:attempt-1:revision-1\n            digest: sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac\n            digest-domain: canonical-bytes\n            room-ref: ./review/validation/briefing-1\ntitle: Task\n---\n# Task\n")
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
	source, err := os.ReadFile(filepath.Join("..", "gates", "testdata", "cross-logical-gate-reentry.md"))
	if err != nil {
		t.Fatal(err)
	}
	entity = filepath.Join(root, "durable-gate-approval-pending-blockers.md")
	if err := os.WriteFile(entity, source, 0o644); err != nil {
		t.Fatal(err)
	}
	briefing = filepath.Join(root, "review", "ideation", "briefing-18", "revision-18.json")
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
