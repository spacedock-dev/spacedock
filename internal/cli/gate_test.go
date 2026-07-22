package cli

import (
	"bytes"
	"context"
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
	writeFile(t, briefing, `{"type":"Briefing","id":"briefing:provider"}`)
	op := filepath.Join(root, "open.yml")
	writeFile(t, op, "operation: open\ngate-id: gate:design\nstage: ideation\nattempt-id: attempt:design-1\nbriefing: {id: 'briefing:design-1', room-ref: './review/ideation'}\n")

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--operation", op, "--briefing", briefing}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
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

func TestGateRecordOperationClosesMinimalBriefingAttempt(t *testing.T) {
	root, entity, briefing := crossGateFixture(t)
	const digest = "sha256:6b2c4f1388a58f42f7c8610f847ed9e7cce92758c00b201d4eb9f4f89dbedd8b"

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

	closeOperation := filepath.Join(root, "close.yml")
	writeFile(t, closeOperation, "operation: close\nexpected: {gate: 'gate:docs-dev:3k:ideation', attempt: 'gate-attempt:3k-ideation-10', briefing: 'briefing:docs-dev:3k:ideation:attempt-10:revision-18', digest: '"+digest+"'}\ngate-id: gate:docs-dev:3k:ideation\nattempt-id: gate-attempt:3k-ideation-10\nresult:\n  briefing-digest: "+digest+"\n  authorized-by: person:captain\n  entries:\n    - type: Resolution\n      id: resolution:docs-dev:3k:ideation:attempt-10\n      briefing: briefing:docs-dev:3k:ideation:attempt-10:revision-18\n      by: person:captain\n      at: 2026-07-22T00:00:00Z\n      decision: approve\n      reason: lgtm. add the fixture\n")
	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"gate", "record", "durable-gate-approval-pending-blockers", "--workflow-dir", root, "--operation", closeOperation}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("operation close exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if !strings.Contains(out.String(), "attempt=gate-attempt:3k-ideation-10 state=closed") || !strings.Contains(out.String(), "decision=approve") {
		t.Fatalf("operation close stdout=%q, want closed approved attempt", out.String())
	}
}

func TestLegacyGateOperationCharacterizesGlobalCurrentDefect(t *testing.T) {
	root, _, briefing := crossGateFixture(t)
	op := filepath.Join(root, "supersede.yml")
	writeFile(t, op, "operation: supersede\nexpected: {gate: 'gate:docs-dev:3k:validation', attempt: 'gate-attempt:3k-validation-1', briefing: 'briefing:docs-dev:3k:validation:attempt-1:revision-1', digest: 'sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac'}\ngate-id: gate:docs-dev:3k:ideation\nattempt-id: gate-attempt:3k-ideation-10\nbriefing: {id: 'briefing:docs-dev:3k:ideation:attempt-10:revision-18', room-ref: './review/ideation/briefing-18'}\n")

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "durable-gate-approval-pending-blockers", "--workflow-dir", root, "--operation", op, "--briefing", briefing}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 1 || !strings.Contains(errOut.String(), "pointer conflict") {
		t.Fatalf("legacy supersede exit=%d stdout=%q stderr=%q, want global-current pointer conflict", code, out.String(), errOut.String())
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
