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

func TestRecordRequiresCanonicalBriefingAtActionableCurrentStage(t *testing.T) {
	for _, tc := range []struct {
		name, status, briefingID, stageFlags, want string
	}{
		{"cross-stage", "implementation", "briefing:task:validation:attempt-1:revision-1", "      gate: true\n", "Briefing stage validation does not match current workflow stage implementation"},
		{"unqualified", "implementation", "briefing:legacy", "      gate: true\n", "Briefing id briefing:legacy is not a canonical stage-qualified v1 identity"},
		{"malformed", "implementation", "briefing:task:implementation:attempt-0:revision-1", "      gate: true\n", "Briefing id briefing:task:implementation:attempt-0:revision-1 is not a canonical stage-qualified v1 identity"},
		{"non-gated", "implementation", "briefing:task:implementation:attempt-1:revision-1", "", "current workflow stage implementation is not an actionable gate:true stage"},
		{"terminal", "done", "briefing:task:done:attempt-1:revision-1", "      gate: true\n      terminal: true\n", "current workflow stage done is not an actionable gate:true stage"},
	} {
		for _, source := range []string{"chat", "room"} {
			t.Run(tc.name+"/"+source, func(t *testing.T) {
				root, entity := recordStageFixture(t, tc.status, tc.briefingID, tc.stageFlags)
				before := readFile(t, entity)
				input := RecordInput{Decision: "hold", Actor: "person:captain", Reason: "wait", WorkflowDir: root}
				if source == "room" {
					input = RecordInput{RoomPath: filepath.Join(root, "missing-room"), WorkflowDir: root}
				}
				err := RecordSemantic(entity, input)
				if err == nil || err.Error() != tc.want {
					t.Fatalf("record error = %v, want %q", err, tc.want)
				}
				if after := readFile(t, entity); after != before {
					t.Fatal("refused record changed entity bytes")
				}
				if _, err := os.Stat(entity + ".gates.lock"); !os.IsNotExist(err) {
					t.Fatalf("refused record left lock residue: %v", err)
				}
			})
		}
	}
}

func TestRecordCanonicalSuccessorAndCrossGateReentry(t *testing.T) {
	root, entity := recordStageFixture(t, "ideation", "briefing:org:task:ideation:attempt-2:revision-3", "      gate: true\n")
	validationRecord := "    - id: gate:task:validation\n      stage: validation\n      attempts:\n        - id: gate-attempt:task-validation-1\n          briefing: {id: 'briefing:task:validation:attempt-1:revision-1', digest: 'sha256:" + strings.Repeat("2", 64) + "', digest-domain: canonical-bytes, room-ref: ./review/validation/briefing-1}\n          resolution: {type: Resolution, id: resolution:validation:1, briefing: 'briefing:task:validation:attempt-1:revision-1', by: person:captain, at: now, decision: revise, reason: rework}\n"
	ideationPrior := "        - id: gate-attempt:task-ideation-1\n          briefing: {id: 'briefing:task:ideation:attempt-1:revision-1', digest: 'sha256:" + strings.Repeat("3", 64) + "', digest-domain: canonical-bytes, room-ref: ./review/ideation/briefing-1}\n          resolution: {type: Resolution, id: resolution:ideation:1, briefing: 'briefing:task:ideation:attempt-1:revision-1', by: person:captain, at: now, decision: hold, reason: wait}\n"
	body := strings.Replace(readFile(t, entity), "current: {gate: 'gate:task:ideation'}", "current: {gate: 'gate:task:validation'}", 1)
	body = strings.Replace(body, "    - id: gate:task:ideation\n", validationRecord+"    - id: gate:task:ideation\n", 1)
	body = strings.Replace(body, "        - id: gate-attempt:task-ideation-2\n", ideationPrior+"        - id: gate-attempt:task-ideation-2\n", 1)
	if err := os.WriteFile(entity, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	validationBefore := marshalAttempt(t, doc.Records[0].Attempts[0])

	if err := RecordSemantic(entity, RecordInput{Decision: "hold", Actor: "person:captain", Reason: "wait", WorkflowDir: root}); err != nil {
		t.Fatal(err)
	}
	doc, _, err = Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Current.Gate != "gate:task:ideation" {
		t.Fatalf("current gate = %q, want re-entered ideation", doc.Current.Gate)
	}
	if got := marshalAttempt(t, doc.Records[0].Attempts[0]); got != validationBefore {
		t.Fatal("cross-gate re-entry modified the formerly selected validation record")
	}
	attempt := doc.Records[1].Attempts[1]
	if attempt.Resolution == nil || attempt.Resolution.Briefing != "briefing:org:task:ideation:attempt-2:revision-3" {
		t.Fatalf("canonical successor closure = %#v", attempt.Resolution)
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

func TestConsumeStaleSupersedesWithoutEffect(t *testing.T) {
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
			// Recovery belongs to the ordinary at-least-once dispatch path. It
			// may re-drive dispatch, but never consumes the authorization again,
			// regardless of whether dispatch had started before the crash.
			again, err := Consume(entity)
			if err != nil {
				t.Fatal(err)
			}
			if again.Consumed {
				t.Fatal("spent authorization was consumed twice")
			}
		})
	}
}

func TestEightCanonicalApplicationShapesReplayByteIdentical(t *testing.T) {
	cases := []struct {
		name, decision, application, encoded string
	}{
		{"approval pending", "approve", "action: advance\n            target-stage: implementation\n            state: pending\n            blockers: []", "              application:\n                action: advance\n                target-stage: implementation\n                state: pending\n                blockers: []\n"},
		{"approval consumed", "approve", "action: advance\n            target-stage: implementation\n            state: consumed\n            blockers: []", "              application:\n                action: advance\n                target-stage: implementation\n                state: consumed\n                blockers: []\n"},
		{"approval superseded", "approve", "action: advance\n            target-stage: implementation\n            state: superseded\n            blockers: []", "              application:\n                action: advance\n                target-stage: implementation\n                state: superseded\n                blockers: []\n"},
		{"approval held", "approve", "action: advance\n            target-stage: implementation\n            state: pending\n            blockers: []\n            execution-hold: {id: hold:1, state: active, by: person:captain}", "              application:\n                action: advance\n                target-stage: implementation\n                state: pending\n                blockers: []\n                execution-hold:\n                    id: hold:1\n                    state: active\n                    by: person:captain\n"},
		{"portable hold", "hold", "action: none\n            state: not-applicable", "              application:\n                action: none\n                state: not-applicable\n"},
		{"feedback pending", "revise", "action: feedback\n            target-stage: ideation\n            state: pending", "              application:\n                action: feedback\n                target-stage: ideation\n                state: pending\n"},
		{"feedback consumed", "revise", "action: feedback\n            target-stage: implementation\n            state: consumed\n            feedback: {cycle: 1, finding-ref: resolution:1, finding-digest: 'sha256:" + strings.Repeat("2", 64) + "'}", "              application:\n                action: feedback\n                target-stage: implementation\n                state: consumed\n                feedback:\n                    cycle: 1\n                    finding-ref: resolution:1\n                    finding-digest: sha256:" + strings.Repeat("2", 64) + "\n"},
		{"historical consumed without blockers", "approve", "action: advance\n            target-stage: implementation\n            state: consumed", "              application:\n                action: advance\n                target-stage: implementation\n                state: consumed\n"},
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
			if !bytes.Contains(block, []byte(tc.encoded)) {
				t.Fatalf("canonical application encoding changed:\n%s", block)
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
	readme := "---\nstages:\n  states:\n    - name: ideation\n      initial: true\n      gate: true\n      feedback-to: ideation\n    - name: implementation\n---\n# Workflow\n"
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	entity := filepath.Join(root, "task.md")
	briefing := filepath.Join(root, "review", "ideation", "briefing-1", "briefing.json")
	if err := os.MkdirAll(filepath.Dir(briefing), 0o755); err != nil {
		t.Fatal(err)
	}
	briefingBytes := []byte(completeBriefing("briefing:task:ideation:attempt-1:revision-1", "review"))
	if err := os.WriteFile(briefing, briefingBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	digest, err := CanonicalDigest(briefingBytes)
	if err != nil {
		t.Fatal(err)
	}
	entityBytes := "---\nstatus: ideation\ntitle: Preserve formatting\ngates:\n" +
		"  version: 1\n" +
		"  current: {gate: 'gate:task:ideation'}\n" +
		"  records:\n" +
		"    - id: gate:task:ideation\n" +
		"      stage: ideation\n" +
		"      attempts:\n" +
		"        - id: gate-attempt:task-ideation-1\n" +
		"          briefing: {id: 'briefing:task:ideation:attempt-1:revision-1', digest: '" + digest + "', digest-domain: canonical-bytes, room-ref: ./review/ideation/briefing-1/briefing.json}\n" +
		"---\n# Task\nBody.\n"
	if err := os.WriteFile(entity, []byte(entityBytes), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, entity
}

func recordStageFixture(t *testing.T, status, briefingID, stageFlags string) (string, string) {
	t.Helper()
	root := t.TempDir()
	readme := "---\nstages:\n  states:\n    - name: " + status + "\n" + stageFlags + "    - name: next\n---\n# Workflow\n"
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	entity := filepath.Join(root, "task.md")
	digest := "sha256:" + strings.Repeat("1", 64)
	records := "    - id: gate:task:" + status + "\n      stage: " + status + "\n      attempts:\n        - id: gate-attempt:task-" + status + "-2\n          briefing: {id: '" + briefingID + "', digest: '" + digest + "', digest-domain: canonical-bytes, room-ref: ./review/" + status + "/briefing-2}\n"
	body := "---\nstatus: " + status + "\ngates:\n  version: 1\n  current: {gate: 'gate:task:" + status + "'}\n  records:\n" + records + "---\n# Task\n"
	if err := os.WriteFile(entity, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, entity
}

// applicationTerminalWorkflow mirrors applicationWorkflow but targets the
// terminal stage: ideation (gate, feedback-to: ideation) -> done (terminal).
func applicationTerminalWorkflow(t *testing.T) (string, string) {
	t.Helper()
	root, entity := applicationWorkflow(t)
	readme := "---\nstages:\n  states:\n    - name: ideation\n      initial: true\n      gate: true\n      feedback-to: ideation\n    - name: done\n      terminal: true\n---\n# Workflow\n"
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, entity
}

// TestConsumeTerminalTargetRoutesWithoutSpending pins the sole-consumer rule:
// a terminal-target approval stays pending at consume, status is untouched, the
// approved-awaiting-merge route is returned, and a repeated consume re-returns
// the same route (routing is an at-least-once effect; the authority never
// moves). Anything else (a spent application, a done status, a missing route)
// is the spend-at-consume hole this design removes.
func TestConsumeTerminalTargetRoutesWithoutSpending(t *testing.T) {
	root, entity := applicationTerminalWorkflow(t)
	if err := RecordSemantic(entity, RecordInput{Decision: "approve", Actor: "person:captain", WorkflowDir: root}); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	for pass := 0; pass < 3; pass++ {
		result, err := ConsumeAt(entity, root)
		if err != nil {
			t.Fatal(err)
		}
		if result.Consumed || result.Route != RouteApprovedAwaitingMerge || !result.Eligible {
			t.Fatalf("terminal consume = %#v, want unconsumed route %q", result, RouteApprovedAwaitingMerge)
		}
	}
	after, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatalf("terminal consume wrote the entity:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if !strings.Contains(string(after), "status: ideation") || !strings.Contains(string(after), "state: pending") {
		t.Fatalf("terminal consume must leave status at the gated stage and the approval pending:\n%s", after)
	}
}

// TestConsumeNonTerminalStillSpendsOnce pins the unchanged non-terminal arm:
// existing approvals keep spending at consume (TestConsumeAdvancesAndSpendsAuthorizationOnce
// covers the advance); this one locks in that AdvanceTargetTerminal's terminal
// flag drives the routing split, not any stage-name heuristic.
func TestConsumeNonTerminalStillSpendsOnce(t *testing.T) {
	root, entity := applicationWorkflow(t)
	if err := RecordSemantic(entity, RecordInput{Decision: "approve", Actor: "person:captain", WorkflowDir: root}); err != nil {
		t.Fatal(err)
	}
	result, err := ConsumeAt(entity, root)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Consumed || result.Route != "" {
		t.Fatalf("non-terminal consume = %#v, want spent with no route", result)
	}
}

// TestTerminalSpendAndReworkGuardReuse pins the envelope's exactly-once reuse:
// the delivery envelope spends pending->consumed and the --rework route
// supersedes pending->superseded through the SAME guarded mutation; a second
// spend or a spend of an already-superseded application is refused unchanged.
func TestTerminalSpendAndReworkGuardReuse(t *testing.T) {
	root, entity := applicationWorkflow(t)
	if err := RecordSemantic(entity, RecordInput{Decision: "approve", Actor: "person:captain", WorkflowDir: root}); err != nil {
		t.Fatal(err)
	}
	doc, oldNode, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	attempt := doc.Records[0].Attempts[len(doc.Records[0].Attempts)-1]
	attempt.Application.State = "consumed"
	if err := ValidateApplicationMutation(oldNode, doc, attempt.ID, "pending", "consumed"); err != nil {
		t.Fatalf("envelope spend must reuse the pending->consumed guard: %v", err)
	}
	// A second spend is impossible: the once-consumed application is not
	// pending, so eligibility refuses it before any writer runs.
	consumedDoc := *doc
	consumed := EvaluateEligibility(&consumedDoc, "implementation", true)
	if consumed.Eligible || consumed.Condition != "consumed" {
		t.Fatalf("re-spend eligibility = %#v, want fail-closed consumed", consumed)
	}
	doc2, oldNode2, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	attempt2 := doc2.Records[0].Attempts[len(doc2.Records[0].Attempts)-1]
	attempt2.Application.State = "superseded"
	if err := ValidateApplicationMutation(oldNode2, doc2, attempt2.ID, "pending", "superseded"); err != nil {
		t.Fatalf("rework route must reuse the pending->superseded guard: %v", err)
	}
	// A superseded application stays non-eligible (fail-closed, as today):
	// superseded authority is never re-spent.
	supersededDoc := *doc2
	elig := EvaluateEligibility(&supersededDoc, "ideation", true)
	if elig.Eligible || elig.Condition != "superseded" {
		t.Fatalf("superseded application eligibility = %#v, want fail-closed superseded", elig)
	}
}
