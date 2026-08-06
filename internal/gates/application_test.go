package gates

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// v1PilotManifest is the checked-in list of the active and archived pilot
// entities whose durable format is part of this unreleased-v1 cut.
//
//go:embed testdata/v1_pilot_manifest.txt
var v1PilotManifest string

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
		{name: "wrong stage", status: "validation", current: true, condition: "ineligible"},
		{name: "hold without application", status: "ideation", current: true, mutate: func(d *Document) {
			d.Records[0].Attempts[0].Resolution.Decision = "hold"
			d.Records[0].Attempts[0].Resolution.Reason = "wait"
			d.Records[0].Attempts[0].Application = nil
		}, condition: "not-applicable"},
		{name: "revise without application", status: "ideation", current: true, mutate: func(d *Document) {
			d.Records[0].Attempts[0].Resolution.Decision = "revise"
			d.Records[0].Attempts[0].Resolution.Reason = "changes requested"
			d.Records[0].Attempts[0].Application = nil
		}, condition: "feedback-pending"},
		{name: "wrong decision", status: "ideation", current: true, mutate: func(d *Document) {
			d.Records[0].Attempts[0].Resolution.Decision = "revise"
		}, condition: "ineligible"},
		{name: "missing application", status: "ideation", current: true, mutate: func(d *Document) {
			d.Records[0].Attempts[0].Application = nil
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
		decision, target, state string
		wantApplication         bool
	}{
		{decision: "approve", target: "implementation", state: "pending", wantApplication: true},
		{decision: "revise", wantApplication: false},
		{decision: "hold", wantApplication: false},
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
			doc, gatesNode, err := Read(entity)
			if err != nil {
				t.Fatal(err)
			}
			app := doc.Records[0].Attempts[0].Application
			if tc.wantApplication {
				if app == nil || app.TargetStage != tc.target || app.State != tc.state {
					t.Fatalf("application = %#v, want target=%s state=%s", app, tc.target, tc.state)
				}
				if got := firstApplicationNode(gatesNode); !sameStrings(yamlMappingKeys(got), []string{"state", "target-stage"}) {
					t.Fatalf("approval YAML application keys = %v, want [state target-stage]", yamlMappingKeys(got))
				}
			} else if app != nil {
				t.Fatalf("%s unexpectedly carries application %#v", tc.decision, app)
			} else if got := firstApplicationNode(gatesNode); got != nil {
				t.Fatalf("%s unexpectedly emitted an application YAML node with keys %v", tc.decision, yamlMappingKeys(got))
			}
		})
	}
}

func TestV1PilotManifestReadsAndValidates(t *testing.T) {
	paths := manifestPaths(v1PilotManifest)
	if len(paths) != 31 {
		t.Fatalf("pilot manifest has %d paths, want 31", len(paths))
	}
	seen := make(map[string]bool, len(paths))
	archives := 0
	for _, rel := range paths {
		if rel == "" || filepath.IsAbs(rel) || filepath.Clean(rel) != rel || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || seen[rel] {
			t.Fatalf("manifest contains invalid or duplicate path %q", rel)
		}
		seen[rel] = true
		if strings.HasPrefix(filepath.ToSlash(rel), "_archive/") {
			archives++
		}
	}
	if archives != 24 {
		t.Fatalf("pilot manifest has %d archived paths, want 24", archives)
	}
	stateRoot := v1PilotStateRoot()
	if stateRoot == "" {
		t.Skip("shared docs/dev/.spacedock-state checkout is not reachable from the code worktree")
	}
	for _, rel := range paths {
		rel, path := rel, filepath.Join(stateRoot, filepath.FromSlash(rel))
		t.Run(rel, func(t *testing.T) {
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("manifest path is not present: %v", err)
			}
			doc, gatesNode, err := Read(path)
			if err != nil {
				t.Fatalf("strict gates.Read: %v", err)
			}
			if err := Validate(doc); err != nil {
				t.Fatalf("gates.Validate: %v", err)
			}
			assertPilotApplicationNodes(t, gatesNode)
		})
	}
}

func manifestPaths(raw string) []string {
	var paths []string
	for _, line := range strings.Split(raw, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			paths = append(paths, line)
		}
	}
	return paths
}

func v1PilotStateRoot() string {
	if root := strings.TrimSpace(os.Getenv("SPACEDOCK_STATE_ROOT")); root != "" {
		return root
	}
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	for dir := cwd; ; dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, "docs", "dev", ".spacedock-state")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
	}
}

func assertPilotApplicationNodes(t *testing.T, gatesNode *yaml.Node) {
	t.Helper()
	records := yamlMappingValue(gatesNode, "records")
	if records == nil || records.Kind != yaml.SequenceNode {
		t.Fatalf("gates.records is not a sequence")
	}
	for _, record := range records.Content {
		attempts := yamlMappingValue(record, "attempts")
		if attempts == nil || attempts.Kind != yaml.SequenceNode {
			t.Fatalf("gate record attempts is not a sequence")
		}
		for _, attempt := range attempts.Content {
			resolution := yamlMappingValue(attempt, "resolution")
			application := yamlMappingValue(attempt, "application")
			decision := ""
			if resolution != nil {
				if n := yamlMappingValue(resolution, "decision"); n != nil {
					decision = n.Value
				}
			}
			if application == nil {
				continue
			}
			if decision != "approve" {
				t.Fatalf("non-approval decision %q carries application keys %v", decision, yamlMappingKeys(application))
			}
			if !sameStrings(yamlMappingKeys(application), []string{"state", "target-stage"}) {
				t.Fatalf("approval application keys = %v, want [state target-stage]", yamlMappingKeys(application))
			}
		}
	}
}

func firstApplicationNode(gatesNode *yaml.Node) *yaml.Node {
	records := yamlMappingValue(gatesNode, "records")
	if records == nil || len(records.Content) == 0 {
		return nil
	}
	attempts := yamlMappingValue(records.Content[0], "attempts")
	if attempts == nil || len(attempts.Content) == 0 {
		return nil
	}
	return yamlMappingValue(attempts.Content[0], "application")
}

func yamlMappingValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil {
		return nil
	}
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return nil
		}
		node = node.Content[0]
	}
	if node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

func yamlMappingKeys(node *yaml.Node) []string {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	keys := make([]string, 0, len(node.Content)/2)
	for i := 0; i+1 < len(node.Content); i += 2 {
		keys = append(keys, node.Content[i].Value)
	}
	sort.Strings(keys)
	return keys
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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
	validationRecord := "    - id: gate:task:validation\n      stage: validation\n      attempts:\n        - id: gate-attempt:task-validation-1\n          briefing: {id: 'briefing:task:validation:attempt-1:revision-1', digest: 'sha256:" + strings.Repeat("2", 64) + "', room-ref: ./review/validation/briefing-1}\n          resolution: {type: Resolution, id: resolution:validation:1, briefing: 'briefing:task:validation:attempt-1:revision-1', by: person:captain, at: now, decision: revise, reason: rework}\n"
	ideationPrior := "        - id: gate-attempt:task-ideation-1\n          briefing: {id: 'briefing:task:ideation:attempt-1:revision-1', digest: 'sha256:" + strings.Repeat("3", 64) + "', room-ref: ./review/ideation/briefing-1}\n          resolution: {type: Resolution, id: resolution:ideation:1, briefing: 'briefing:task:ideation:attempt-1:revision-1', by: person:captain, at: now, decision: hold, reason: wait}\n"
	body := readFile(t, entity)
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

// TestConsumeRepeatAfterStaleSupersedeReportsNoWrite pins the fix for the
// wrote-detection bug: EvaluateEligibility copies the attempt's CURRENT
// application state into ApplicationState on every read — including a pure
// refusal against an application that is ALREADY superseded from a prior
// call — so ApplicationState == "superseded" alone cannot distinguish a fresh
// write from a repeat no-op read. Only ConsumeResult.Wrote may be used for
// that; a caller checking ApplicationState instead would sync/commit on a
// repeat call that wrote nothing.
func TestConsumeRepeatAfterStaleSupersedeReportsNoWrite(t *testing.T) {
	root, entity := applicationWorkflow(t)
	if err := RecordSemantic(entity, RecordInput{Decision: "approve", Actor: "person:captain", WorkflowDir: root}); err != nil {
		t.Fatal(err)
	}
	room := filepath.Join(filepath.Dir(entity), "review", "ideation", "briefing-1", "briefing.json")
	if err := os.WriteFile(room, []byte(completeBriefing("briefing:task:ideation:1", "drifted")), 0o644); err != nil {
		t.Fatal(err)
	}
	first, err := Consume(entity)
	if err != nil {
		t.Fatal(err)
	}
	if !first.Wrote || first.ApplicationState != "superseded" {
		t.Fatalf("first (real) supersede = %#v, want Wrote=true ApplicationState=superseded", first)
	}
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}

	second, err := Consume(entity)
	if err != nil {
		t.Fatal(err)
	}
	if second.Wrote {
		t.Fatalf("repeat consume after supersede reported Wrote=true, want false (no write): %#v", second)
	}
	if second.ApplicationState != "superseded" {
		t.Fatalf("repeat consume ApplicationState = %q, want superseded (read of the already-superseded state)", second.ApplicationState)
	}
	after, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("repeat consume after supersede changed entity bytes despite Wrote=false")
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

func TestCanonicalApplicationShapesReplayByteIdentical(t *testing.T) {
	cases := []struct {
		name, decision, application, encoded string
	}{
		{"approval pending", "approve", "target-stage: implementation\n            state: pending", "              application:\n                target-stage: implementation\n                state: pending\n"},
		{"approval consumed", "approve", "target-stage: implementation\n            state: consumed", "              application:\n                target-stage: implementation\n                state: consumed\n"},
		{"approval superseded", "approve", "target-stage: implementation\n            state: superseded", "              application:\n                target-stage: implementation\n                state: superseded\n"},
	}
	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reason := ""
			if tc.decision != "approve" {
				reason = "\n            reason: recorded rationale"
			}
			source := "---\nstatus: ideation\ngates:\n  version: 1\n  records:\n    - id: gate:replay\n      stage: ideation\n      attempts:\n        - id: attempt:replay-" + string(rune('a'+i)) + "\n          briefing:\n            id: briefing:replay\n            digest: sha256:" + strings.Repeat("1", 64) + "\n            room-ref: ./review\n          resolution:\n            type: Resolution\n            id: resolution:replay\n            briefing: briefing:replay\n            by: person:captain\n            at: 2026-07-22T00:00:00Z\n            decision: " + tc.decision + reason + "\n          application:\n            " + tc.application + "\n---\n# Replay\n"
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

func TestApplicationExtensionShapesWarnWithoutMutation(t *testing.T) {
	for _, tc := range []struct {
		name, decision, application string
	}{
		{"action", "approve", "action: advance\n            target-stage: implementation\n            state: pending"},
		{"blockers", "approve", "target-stage: implementation\n            state: pending\n            blockers: []"},
		{"execution-hold", "approve", "target-stage: implementation\n            state: pending\n            execution-hold: {state: active}"},
		{"feedback", "approve", "target-stage: implementation\n            state: pending\n            feedback: {cycle: 1}"},
		{"not-applicable", "hold", "state: not-applicable"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source := "---\nstatus: ideation\ngates:\n  version: 1\n  records:\n    - id: gate:legacy\n      stage: ideation\n      attempts:\n        - id: attempt:legacy\n          briefing:\n            id: briefing:legacy:ideation:attempt-1:revision-1\n            digest: sha256:" + strings.Repeat("1", 64) + "\n            room-ref: ./review\n          resolution:\n            type: Resolution\n            id: resolution:legacy\n            briefing: briefing:legacy:ideation:attempt-1:revision-1\n            by: person:captain\n            at: 2026-07-22T00:00:00Z\n            decision: " + tc.decision + "\n            reason: legacy\n          application:\n            " + tc.application + "\n---\n# Legacy\n"
			before := []byte(source)
			doc, _, warnings, err := readDataDiagnostics(before)
			if tc.name == "not-applicable" {
				if err == nil {
					t.Fatal("invalid application state was accepted")
				}
			} else {
				if err != nil {
					t.Fatalf("application extension was rejected: %v", err)
				}
				if doc == nil || len(warnings) != 1 || warnings[0].Field == "" {
					t.Fatalf("warnings = %#v, doc=%#v; want one extension warning", warnings, doc)
				}
			}
			if !bytes.Equal(before, []byte(source)) {
				t.Fatal("read mutated fixture bytes")
			}
		})
	}
}

func TestReadDiagnosticsFiltersOnlyExactApplicationMappings(t *testing.T) {
	source := "---\nstatus: ideation\ngates:\n" +
		"  version: 1\n" +
		"  records:\n" +
		"    - id: gate:one\n" +
		"      stage: ideation\n" +
		"      attempts:\n" +
		"        - id: attempt:one\n" +
		"          briefing: {id: briefing:one, digest: sha256:" + strings.Repeat("1", 64) + ", room-ref: ./review}\n" +
		"          resolution: {type: Resolution, id: resolution:one, briefing: briefing:one, by: person:captain, at: 2026-07-22T00:00:00Z, decision: approve}\n" +
		"          application:\n" +
		"            nested: [one, two]\n" +
		"            target-stage: implementation\n" +
		"            state: pending\n" +
		"            feedback: {owner: old-producer}\n" +
		"            blockers: []\n" +
		"            action: advance\n" +
		"            action: duplicate-legacy-value\n" +
		"            execution-hold: true\n" +
		"    - id: gate:two\n" +
		"      stage: validation\n" +
		"      attempts:\n" +
		"        - id: attempt:two\n" +
		"          briefing: {id: briefing:two, digest: sha256:" + strings.Repeat("2", 64) + ", room-ref: ./review}\n" +
		"          resolution: {type: Resolution, id: resolution:two, briefing: briefing:two, by: person:captain, at: 2026-07-22T00:00:00Z, decision: approve}\n" +
		"          application: {target-stage: done, state: consumed, binding: {gate: old-producer}}\n" +
		"---\n# Task\n"
	doc, original, warnings, err := readDataDiagnostics([]byte(source))
	if err != nil {
		t.Fatalf("diagnostic read = %v", err)
	}
	path := filepath.Join(t.TempDir(), "task.md")
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := ReadDiagnostics(path); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("diagnostic read changed source bytes")
	}
	if len(warnings) != 6 {
		t.Fatalf("warnings = %#v, want six unknown application keys", warnings)
	}
	for i := 1; i < len(warnings); i++ {
		if warnings[i-1].Path > warnings[i].Path || warnings[i-1].Path == warnings[i].Path && warnings[i-1].Field > warnings[i].Field {
			t.Fatalf("warnings are not sorted: %#v", warnings)
		}
	}
	gotFields := make([]string, 0, len(warnings))
	for _, warning := range warnings {
		gotFields = append(gotFields, warning.Field)
	}
	if strings.Join(gotFields, ",") != "action,blockers,execution-hold,feedback,nested,binding" {
		t.Fatalf("warning fields = %v", gotFields)
	}
	if got := doc.Records[0].Attempts[0].Application; got == nil || got.TargetStage != "implementation" || got.State != "pending" {
		t.Fatalf("first canonical application = %#v", got)
	}
	if got := doc.Records[1].Attempts[0].Application; got == nil || got.TargetStage != "done" || got.State != "consumed" {
		t.Fatalf("second canonical application = %#v", got)
	}
	if original == nil || yamlMappingValue(original, "records") == nil {
		t.Fatal("diagnostic read did not return the original gates node")
	}
	originalBytes, err := yaml.Marshal(original)
	if err != nil || !bytes.Contains(originalBytes, []byte("action")) || !bytes.Contains(originalBytes, []byte("nested")) {
		t.Fatalf("original gates node was filtered: %s", originalBytes)
	}
}

func TestReadDiagnosticsRejectsNonMappingApplicationShapes(t *testing.T) {
	for _, tc := range []struct {
		name, value string
	}{
		{"null", "null"},
		{"sequence", "[legacy]"},
		{"scalar", "legacy"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source := "gates:\n  version: 1\n  records:\n    - id: gate:shape\n      stage: ideation\n      attempts:\n        - id: attempt:shape\n          briefing: {id: briefing:shape, digest: sha256:" + strings.Repeat("1", 64) + ", room-ref: ./review}\n          resolution: {type: Resolution, id: resolution:shape, briefing: briefing:shape, by: person:captain, at: 2026-07-22T00:00:00Z, decision: approve}\n          application: " + tc.value + "\n"
			if _, _, _, err := readDataDiagnostics([]byte(source)); err == nil {
				t.Fatal("non-mapping application was accepted")
			}
		})
	}
}

func TestReadDiagnosticsKeepsCanonicalAndBindingFailuresStrict(t *testing.T) {
	base := func(application, briefing string) string {
		return "gates:\n  version: 1\n  records:\n    - id: gate:strict\n      stage: ideation\n      attempts:\n        - id: attempt:strict\n          briefing: {id: briefing:strict, digest: " + briefing + ", room-ref: ./review}\n          resolution: {type: Resolution, id: resolution:strict, briefing: briefing:strict, by: person:captain, at: 2026-07-22T00:00:00Z, decision: approve}\n          application:\n" + application + "\n"
	}
	for _, tc := range []struct {
		name, source string
	}{
		{"missing target", base("            state: pending", "sha256:"+strings.Repeat("1", 64))},
		{"invalid state", base("            target-stage: implementation\n            state: waiting", "sha256:"+strings.Repeat("1", 64))},
		{"duplicate canonical", base("            target-stage: implementation\n            target-stage: other\n            state: pending", "sha256:"+strings.Repeat("1", 64))},
		{"bad binding", base("            target-stage: implementation\n            state: pending", "not-a-digest")},
		{"unknown outside application", base("            target-stage: implementation\n            state: pending\n          binding: wrong", "sha256:"+strings.Repeat("1", 64))},
		{"malformed YAML", "gates: [\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, _, err := readDataDiagnostics([]byte(tc.source)); err == nil {
				t.Fatal("strict defect was accepted")
			}
		})
	}
}

func eligibleDocument() *Document {
	return &Document{Version: 1, Records: []GateRecord{{
		ID: "gate:task:ideation", Stage: "ideation", Attempts: []Attempt{{
			ID:          "attempt:1",
			Briefing:    Briefing{ID: "briefing:1", Digest: "sha256:" + strings.Repeat("1", 64), RoomRef: "./review"},
			Resolution:  &Resolution{Type: "Resolution", ID: "resolution:1", Briefing: "briefing:1", By: "person:captain", At: "now", Decision: "approve"},
			Application: &Application{TargetStage: "implementation", State: "pending"},
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
		"  records:\n" +
		"    - id: gate:task:ideation\n" +
		"      stage: ideation\n" +
		"      attempts:\n" +
		"        - id: gate-attempt:task-ideation-1\n" +
		"          briefing: {id: 'briefing:task:ideation:attempt-1:revision-1', digest: '" + digest + "', room-ref: ./review/ideation/briefing-1/briefing.json}\n" +
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
	records := "    - id: gate:task:" + status + "\n      stage: " + status + "\n      attempts:\n        - id: gate-attempt:task-" + status + "-2\n          briefing: {id: '" + briefingID + "', digest: '" + digest + "', room-ref: ./review/" + status + "/briefing-2}\n"
	body := "---\nstatus: " + status + "\ngates:\n  version: 1\n  records:\n" + records + "---\n# Task\n"
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
		if result.Consumed || !result.Eligible || !ApprovedAwaitingMergeRoute(entity, root) {
			t.Fatalf("terminal consume = %#v routed=%t, want unconsumed approved-awaiting-merge routing", result, ApprovedAwaitingMergeRoute(entity, root))
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
	if !result.Consumed || ApprovedAwaitingMergeRoute(entity, root) {
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
	if err := validateApplicationMutation(oldNode, doc, attempt.ID, "pending", "consumed"); err != nil {
		t.Fatalf("envelope spend must reuse the pending->consumed guard: %v", err)
	}
	// A second spend is impossible: the once-consumed application is not
	// pending, so eligibility refuses it before any writer runs.
	consumedDoc := *doc
	consumed := EvaluateEligibility(&consumedDoc, "ideation", true)
	if consumed.Eligible || consumed.Condition != "consumed" {
		t.Fatalf("re-spend eligibility = %#v, want fail-closed consumed", consumed)
	}
	doc2, oldNode2, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	attempt2 := doc2.Records[0].Attempts[len(doc2.Records[0].Attempts)-1]
	attempt2.Application.State = "superseded"
	if err := validateApplicationMutation(oldNode2, doc2, attempt2.ID, "pending", "superseded"); err != nil {
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
