package status

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"gopkg.in/yaml.v3"
)

func TestStatusTextAndJSONProjectApprovedPendingApplication(t *testing.T) {
	root := t.TempDir()
	readme := "---\nid-style: slug\nstages:\n  states:\n    - name: ideation\n      initial: true\n    - name: implementation\n---\n# Workflow\n"
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	briefing := []byte(`{"type":"Briefing","version":"1","id":"briefing:task:1","question":"approve?","artifacts":[{"id":"artifact:1","uri":"artifact.md","rev":"sha256:` + strings.Repeat("a", 64) + `"}]}`)
	digest, err := gates.CanonicalDigest(briefing)
	if err != nil {
		t.Fatal(err)
	}
	room := filepath.Join(root, "review", "ideation", "briefing-1")
	if err := os.MkdirAll(room, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(room, "briefing.json"), briefing, 0o644); err != nil {
		t.Fatal(err)
	}
	body := "---\nid: task\nstatus: ideation\ntitle: Task\ngates:\n  version: 1\n  current: {gate: 'gate:task:ideation'}\n  records:\n    - id: gate:task:ideation\n      stage: ideation\n      attempts:\n        - id: attempt:task-1\n          briefing: {id: 'briefing:task:1', digest: '" + digest + "', digest-domain: canonical-bytes, room-ref: ./review/ideation/briefing-1}\n          resolution: {type: Resolution, id: 'resolution:task-1', briefing: 'briefing:task:1', by: 'person:captain', at: '2026-07-22T00:00:00Z', decision: approve}\n          application: {action: advance, target-stage: implementation, state: pending, blockers: []}\n---\n# Task\n"
	if err := os.WriteFile(filepath.Join(root, "task.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	args := []string{"--workflow-dir", root, "--fields", "gate-decision,gate-application,gate-condition,gate-eligible"}
	text, stderr, code := runNative(t, root, nil, args...)
	if code != 0 || !strings.Contains(text, "approve") || !strings.Contains(text, "advance/pending") || !strings.Contains(text, "approved-pending") {
		t.Fatalf("text status exit=%d stderr=%q output=%q", code, stderr, text)
	}
	jsonOut, stderr, code := runNative(t, root, nil, append(args, "--json")...)
	if code != 0 || !strings.Contains(jsonOut, `"gate-application":"advance/pending"`) || !strings.Contains(jsonOut, `"gate-condition":"approved-pending"`) || !strings.Contains(jsonOut, `"gate-eligible":"true"`) {
		t.Fatalf("json status exit=%d stderr=%q output=%q", code, stderr, jsonOut)
	}
	changed := strings.Replace(readme, "name: implementation", "name: validation", 1)
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(changed), 0o644); err != nil {
		t.Fatal(err)
	}
	jsonOut, stderr, code = runNative(t, root, nil, append(args, "--json")...)
	if code != 0 || !strings.Contains(jsonOut, `"gate-condition":"ineligible"`) || !strings.Contains(jsonOut, `"gate-eligible":"false"`) {
		t.Fatalf("changed-taxonomy status exit=%d stderr=%q output=%q", code, stderr, jsonOut)
	}
}

func TestStatusTextAndJSONProjectAllRecordedResolutionStates(t *testing.T) {
	root := t.TempDir()
	readme := "---\nid-style: slug\nstages:\n  states:\n    - name: ideation\n      initial: true\n---\n# Workflow\n"
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, decision := range []string{"approve", "hold", "revise"} {
		reason := ""
		if decision != "approve" {
			reason = ", reason: 'recorded reason'"
		}
		body := "---\nstatus: ideation\ntitle: " + decision + "\ngates:\n  version: 1\n  current: {gate: 'gate:" + decision + "'}\n  records:\n    - id: gate:" + decision + "\n      stage: ideation\n      attempts:\n        - id: attempt:" + decision + "-1\n          briefing: {id: 'briefing:" + decision + "-1', digest: 'sha256:" + strings.Repeat("2", 64) + "', digest-domain: canonical-bytes, room-ref: ./review}\n          resolution: {type: Resolution, id: 'resolution:" + decision + "-1', briefing: 'briefing:" + decision + "-1', by: 'person:captain', at: '2026-07-22T00:00:00Z', decision: " + decision + reason + "}\n---\n"
		if err := os.WriteFile(filepath.Join(root, decision+".md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	args := []string{"--workflow-dir", root, "--fields", "gate-state,gate-decision,gate-resolution"}
	text, stderr, code := runNative(t, root, nil, args...)
	if code != 0 {
		t.Fatalf("text status exit=%d stderr=%q", code, stderr)
	}
	jsonOut, stderr, code := runNative(t, root, nil, append(args, "--json")...)
	if code != 0 {
		t.Fatalf("json status exit=%d stderr=%q", code, stderr)
	}
	for _, decision := range []string{"approve", "hold", "revise"} {
		if !strings.Contains(text, decision) || !strings.Contains(jsonOut, `"gate-decision":"`+decision+`"`) {
			t.Errorf("decision %s missing from status outputs:\ntext=%s\njson=%s", decision, text, jsonOut)
		}
	}
}

func TestUnrelatedSetPreservesGatesAndStatusProjectsResolution(t *testing.T) {
	path := filepath.Join(t.TempDir(), "task.md")
	body := "---\nid: task\nstatus: ideation\nscore: '0.5'\ngates:\n  version: 1\n  current: {gate: 'gate:design'}\n  records:\n    - id: gate:design\n      stage: ideation\n      attempts:\n        - id: attempt:design-1\n          briefing: {id: 'briefing:design-1', digest: 'sha256:" + strings.Repeat("1", 64) + "', digest-domain: canonical-bytes, room-ref: ./review}\n          resolution: {type: Resolution, id: 'resolution:design-1', briefing: 'briefing:design-1', by: 'person:captain', at: '2026-07-22T00:00:00Z', decision: hold, reason: 'wait'}\n          application: {action: none, state: not-applicable}\n---\n# Task\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	before := nestedGateValue(t, path)
	if _, err := updateFrontmatter(path, []fieldUpdate{{field: "score", value: "0.6", hasValue: true}}); err != nil {
		t.Fatal(err)
	}
	after := nestedGateValue(t, path)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("unrelated --set changed gates:\nbefore=%#v\nafter=%#v", before, after)
	}
	e := newEntity(ParseFrontmatter(path), "task", path, "active")
	if e.fields["gate-state"] != "closed" || e.fields["gate-decision"] != "hold" || e.fields["gate-resolution"] != "resolution:design-1" {
		t.Fatalf("derived gate status = %#v", e.fields)
	}
}

func nestedGateValue(t *testing.T, path string) any {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var fm map[string]any
	if err := yaml.Unmarshal(frontmatterSlice(b), &fm); err != nil {
		t.Fatal(err)
	}
	return fm["gates"]
}
