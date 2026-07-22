package gates

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestEligibilityFailClosedTable(t *testing.T) {
	base := eligibleDocument()
	tests := []struct {
		name      string
		mutate    func(*Document)
		status    string
		current   bool
		eligible  bool
		condition string
	}{
		{name: "exact current pending approval", status: "ideation", current: true, eligible: true, condition: "approved-pending"},
		{name: "stale reviewed input", status: "ideation", current: false, condition: "stale"},
		{name: "superseded", status: "ideation", current: true, mutate: setApplicationState("superseded"), condition: "superseded"},
		{name: "consumed", status: "ideation", current: true, mutate: setApplicationState("consumed"), condition: "consumed"},
		{name: "active hold", status: "ideation", current: true, mutate: func(d *Document) {
			d.Records[0].Attempts[0].Application.ExecutionHold = &ExecutionHold{State: "active"}
		}, condition: "held"},
		{name: "unknown hold", status: "ideation", current: true, mutate: func(d *Document) {
			d.Records[0].Attempts[0].Application.ExecutionHold = &ExecutionHold{State: "mystery"}
		}, condition: "ineligible"},
		{name: "blocker present", status: "ideation", current: true, mutate: func(d *Document) {
			blockers := []Blocker{{ID: "blocker:x", State: "unsatisfied"}}
			d.Records[0].Attempts[0].Application.Blockers = &blockers
		}, condition: "blocked"},
		{name: "wrong stage", status: "validation", current: true, condition: "ineligible"},
		{name: "wrong decision", status: "ideation", current: true, mutate: func(d *Document) {
			d.Records[0].Attempts[0].Resolution.Decision = "revise"
		}, condition: "ineligible"},
		{name: "missing application", status: "ideation", current: true, mutate: func(d *Document) {
			d.Records[0].Attempts[0].Application = nil
		}, condition: "ineligible"},
		{name: "missing blockers field", status: "ideation", current: true, mutate: func(d *Document) {
			d.Records[0].Attempts[0].Application.Blockers = nil
		}, condition: "ineligible"},
		{name: "missing target", status: "ideation", current: true, mutate: func(d *Document) {
			d.Records[0].Attempts[0].Application.TargetStage = ""
		}, condition: "ineligible"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc := cloneDocument(t, base)
			if tc.mutate != nil {
				tc.mutate(doc)
			}
			got := EvaluateEligibility(doc, tc.status, tc.current)
			if got.Eligible != tc.eligible || got.Condition != tc.condition {
				t.Fatalf("eligibility = %#v, want eligible=%t condition=%s", got, tc.eligible, tc.condition)
			}
		})
	}
}

func TestRecordClosureShapesApplication(t *testing.T) {
	for _, tc := range []struct {
		decision, action, target, state string
	}{
		{decision: "approve", action: "advance", target: "implementation", state: "pending"},
		{decision: "revise", action: "feedback", target: "ideation", state: "pending"},
		{decision: "hold", action: "none", state: "not-applicable"},
	} {
		t.Run(tc.decision, func(t *testing.T) {
			root, entity := applicationWorkflow(t)
			reason := ""
			if tc.decision != "approve" {
				reason = "captain rationale"
			}
			if err := RecordSemantic(entity, RecordInput{Decision: tc.decision, Actor: "person:captain", Reason: reason, WorkflowDir: root}); err != nil {
				t.Fatal(err)
			}
			doc, _, err := Read(entity)
			if err != nil {
				t.Fatal(err)
			}
			app := doc.Records[0].Attempts[0].Application
			if app == nil || app.Action != tc.action || app.TargetStage != tc.target || app.State != tc.state {
				t.Fatalf("application = %#v, want %s/%s/%s", app, tc.action, tc.target, tc.state)
			}
			if tc.decision == "approve" && (app.Blockers == nil || len(*app.Blockers) != 0) {
				t.Fatalf("approval blockers = %#v, want explicit empty list", app.Blockers)
			}
		})
	}
}

func TestSecondApplicationOnClosedAttemptIsRefusedUnchanged(t *testing.T) {
	root, entity := applicationWorkflow(t)
	input := RecordInput{Decision: "approve", Actor: "person:captain", WorkflowDir: root}
	if err := RecordSemantic(entity, input); err != nil {
		t.Fatal(err)
	}
	before := readFile(t, entity)
	if err := RecordSemantic(entity, input); err == nil || !strings.Contains(err.Error(), "frozen closed") {
		t.Fatalf("second application error = %v, want frozen refusal", err)
	}
	if after := readFile(t, entity); after != before {
		t.Fatal("refused second application changed the gate record")
	}
}

func TestConsumeAdvancesAndSpendsAuthorizationOnce(t *testing.T) {
	root, entity := applicationWorkflow(t)
	if err := RecordSemantic(entity, RecordInput{Decision: "approve", Actor: "person:captain", WorkflowDir: root}); err != nil {
		t.Fatal(err)
	}
	effects := 0
	for pass := 0; pass < 3; pass++ {
		result, err := Consume(entity)
		if err != nil {
			t.Fatal(err)
		}
		if result.Consumed {
			effects++ // the existing caller advances/dispatches only on authorization.
		}
	}
	if effects != 1 {
		t.Fatalf("advance+dispatch effects = %d, want exactly 1", effects)
	}
	body := readFile(t, entity)
	if !strings.Contains(body, "status: implementation") || !strings.Contains(body, "state: consumed") {
		t.Fatalf("atomic co-write missing status or consumed application:\n%s", body)
	}
}

func TestConsumeStaleSupersedesWithoutEffectAndSuccessorLeavesOnePending(t *testing.T) {
	root, entity := applicationWorkflow(t)
	if err := RecordSemantic(entity, RecordInput{Decision: "approve", Actor: "person:captain", WorkflowDir: root}); err != nil {
		t.Fatal(err)
	}
	room := filepath.Join(filepath.Dir(entity), "review", "ideation", "briefing-1", "briefing.json")
	if err := os.WriteFile(room, []byte(completeBriefing("briefing:task:ideation:1", "drifted")), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := Consume(entity)
	if err != nil {
		t.Fatal(err)
	}
	if result.Consumed || result.Condition != "stale" {
		t.Fatalf("stale consume = %#v, want zero effect", result)
	}
	if result.ApplicationState != "superseded" {
		t.Fatalf("stale consume reported application state %q, want superseded", result.ApplicationState)
	}
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if got := doc.Records[0].Attempts[0].Application.State; got != "superseded" {
		t.Fatalf("stale application state = %q, want superseded", got)
	}

	nextBriefing := filepath.Join(filepath.Dir(entity), "review", "ideation", "briefing-2", "briefing.json")
	if err := os.MkdirAll(filepath.Dir(nextBriefing), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(nextBriefing, []byte(completeBriefing("briefing:task:ideation:2", "replacement")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := RecordBriefing(entity, nextBriefing); err != nil {
		t.Fatal(err)
	}
	if err := RecordSemantic(entity, RecordInput{Decision: "approve", Actor: "person:captain", WorkflowDir: root}); err != nil {
		t.Fatal(err)
	}
	doc, _, err = Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	pending := 0
	for _, attempt := range doc.Records[0].Attempts {
		if attempt.Application != nil && attempt.Application.State == "pending" {
			pending++
		}
	}
	if pending != 1 {
		t.Fatalf("pending applications across attempts = %d, want 1", pending)
	}
}

func TestConsumeRefusesTargetRemovedFromCurrentWorkflow(t *testing.T) {
	root, entity := applicationWorkflow(t)
	if err := RecordSemantic(entity, RecordInput{Decision: "approve", Actor: "person:captain", WorkflowDir: root}); err != nil {
		t.Fatal(err)
	}
	changed := "---\nstages:\n  states:\n    - name: ideation\n      initial: true\n    - name: validation\n---\n# Workflow\n"
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(changed), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := ConsumeAt(entity, root)
	if err != nil {
		t.Fatal(err)
	}
	if result.Consumed || result.Eligible || result.Condition != "ineligible" {
		t.Fatalf("removed-target consume = %#v, want fail-closed refusal", result)
	}
	body := readFile(t, entity)
	if !strings.Contains(body, "status: ideation") || !strings.Contains(body, "state: pending") {
		t.Fatalf("removed target changed entity:\n%s", body)
	}
}

func TestResolutionSummaryDoesNotHashBriefing(t *testing.T) {
	root, entity := applicationWorkflow(t)
	if err := RecordSemantic(entity, RecordInput{Decision: "approve", Actor: "person:captain", WorkflowDir: root}); err != nil {
		t.Fatal(err)
	}
	briefing := filepath.Join(root, "review", "ideation", "briefing-1", "briefing.json")
	if err := os.Remove(briefing); err != nil {
		t.Fatal(err)
	}
	summary, err := SummaryFile(entity)
	if err != nil || summary.Decision != "approve" {
		t.Fatalf("resolution-only summary = %#v, %v", summary, err)
	}
	eligibility, err := EligibilityFileAt(entity, root)
	if err != nil || eligibility.Condition != "ineligible" {
		t.Fatalf("explicit eligibility = %#v, %v, want fail-closed unknown", eligibility, err)
	}
	result, err := ConsumeAt(entity, root)
	if err != nil || result.Consumed || result.Condition != "ineligible" {
		t.Fatalf("missing-input consume = %#v, %v, want refusal", result, err)
	}
	doc, _, err := Read(entity)
	if err != nil || doc.Records[0].Attempts[0].Application.State != "pending" {
		t.Fatalf("missing input spent approval: state=%#v err=%v", doc, err)
	}
}

func TestConsumeCrashWindowsNeverReconsumeAuthorization(t *testing.T) {
	for _, crash := range []string{"after-consume-before-dispatch", "after-dispatch-before-observation"} {
		t.Run(crash, func(t *testing.T) {
			root, entity := applicationWorkflow(t)
			if err := RecordSemantic(entity, RecordInput{Decision: "approve", Actor: "person:captain", WorkflowDir: root}); err != nil {
				t.Fatal(err)
			}
			first, err := Consume(entity)
			if err != nil || !first.Consumed {
				t.Fatalf("initial consume = %#v, %v", first, err)
			}
			dispatches := 0
			if crash == "after-dispatch-before-observation" {
				dispatches++
			}
			// Recovery belongs to the ordinary at-least-once dispatch path. It
			// re-drives dispatch from advanced status without consuming again.
			dispatches++
			again, err := Consume(entity)
			if err != nil {
				t.Fatal(err)
			}
			if again.Consumed {
				t.Fatal("spent authorization was consumed twice")
			}
			wantDispatches := 1
			if crash == "after-dispatch-before-observation" {
				wantDispatches = 2
			}
			if dispatches != wantDispatches {
				t.Fatalf("dispatch retries = %d, want %d", dispatches, wantDispatches)
			}
		})
	}
}

func TestEightCanonicalApplicationShapesReplayByteIdentical(t *testing.T) {
	cases := []struct {
		name, decision, application string
	}{
		{"approval pending", "approve", "action: advance\n            target-stage: implementation\n            state: pending\n            blockers: []"},
		{"approval consumed", "approve", "action: advance\n            target-stage: implementation\n            state: consumed\n            blockers: []"},
		{"approval superseded", "approve", "action: advance\n            target-stage: implementation\n            state: superseded\n            blockers: []"},
		{"approval held", "approve", "action: advance\n            target-stage: implementation\n            state: pending\n            blockers: []\n            execution-hold: {id: hold:1, state: active, by: person:captain}"},
		{"portable hold", "hold", "action: none\n            state: not-applicable"},
		{"feedback pending", "revise", "action: feedback\n            target-stage: ideation\n            state: pending"},
		{"feedback consumed", "revise", "action: feedback\n            target-stage: implementation\n            state: consumed\n            feedback: {cycle: 1, finding-ref: resolution:1, finding-digest: 'sha256:" + strings.Repeat("2", 64) + "'}"},
		{"historical consumed without blockers", "approve", "action: advance\n            target-stage: implementation\n            state: consumed"},
	}
	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reason := ""
			if tc.decision != "approve" {
				reason = "\n            reason: recorded rationale"
			}
			source := "---\nstatus: ideation\ngates:\n  version: 1\n  current:\n    gate: gate:replay\n  records:\n    - id: gate:replay\n      stage: ideation\n      attempts:\n        - id: attempt:replay-" + string(rune('a'+i)) + "\n          briefing:\n            id: briefing:replay\n            digest: sha256:" + strings.Repeat("1", 64) + "\n            digest-domain: canonical-bytes\n            room-ref: ./review\n          resolution:\n            type: Resolution\n            id: resolution:replay\n            briefing: briefing:replay\n            by: person:captain\n            at: 2026-07-22T00:00:00Z\n            decision: " + tc.decision + reason + "\n          application:\n            " + tc.application + "\n---\n# Replay\n"
			doc, _, err := readData([]byte(source))
			if err != nil {
				t.Fatal(err)
			}
			block, err := yaml.Marshal(struct {
				Gates *Document `yaml:"gates"`
			}{Gates: doc})
			if err != nil {
				t.Fatal(err)
			}
			entity := filepath.Join(t.TempDir(), "entity.md")
			canonical := append([]byte("---\nstatus: ideation\n"), block...)
			canonical = append(canonical, []byte("title: Replay\n---\n# Replay\n")...)
			if err := os.WriteFile(entity, canonical, 0o644); err != nil {
				t.Fatal(err)
			}
			replayed, expected, err := Read(entity)
			if err != nil {
				t.Fatal(err)
			}
			if err := writeDocument(entity, expected, replayed); err != nil {
				t.Fatal(err)
			}
			after, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(after, canonical) {
				t.Fatal("canonical application replay changed bytes")
			}
		})
	}
}

func eligibleDocument() *Document {
	blockers := []Blocker{}
	return &Document{Version: 1, Current: Selection{Gate: "gate:task:ideation"}, Records: []GateRecord{{
		ID: "gate:task:ideation", Stage: "ideation", Attempts: []Attempt{{
			ID:          "attempt:1",
			Briefing:    Briefing{ID: "briefing:1", Digest: "sha256:" + strings.Repeat("1", 64), DigestDomain: "canonical-bytes", RoomRef: "./review"},
			Resolution:  &Resolution{Type: "Resolution", ID: "resolution:1", Briefing: "briefing:1", By: "person:captain", At: "now", Decision: "approve"},
			Application: &Application{Action: "advance", TargetStage: "implementation", State: "pending", Blockers: &blockers},
		}},
	}}}
}

func setApplicationState(state string) func(*Document) {
	return func(d *Document) { d.Records[0].Attempts[0].Application.State = state }
}

func applicationWorkflow(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	readme := "---\nstages:\n  states:\n    - name: ideation\n      initial: true\n      feedback-to: ideation\n    - name: implementation\n---\n# Workflow\n"
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	entity := filepath.Join(root, "task.md")
	briefing := filepath.Join(root, "review", "ideation", "briefing-1", "briefing.json")
	if err := os.MkdirAll(filepath.Dir(briefing), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(briefing, []byte(completeBriefing("briefing:task:ideation:1", "review")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(entity, []byte("---\nstatus: ideation\ntitle: Preserve formatting\n---\n# Task\nBody.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := RecordBriefing(entity, briefing); err != nil {
		t.Fatal(err)
	}
	return root, entity
}
