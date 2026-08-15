package status

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestStatusProjectsSharedGateReadinessReducer(t *testing.T) {
	root, _ := buildSplitRoot(t, identifyReadyGatesReadme, map[string]string{
		"sp.md": "---\nid: sp\nstatus: validation\nscore: 100\n---\n# incomplete\n",
		"mf.md": openGateEntity("mf", "validation", "90"),
		"r4.md": openGateEntity("r4", "validation", "80"),
		"wd.md": withdrawnGateEntity("wd", "validation", "75"),
		"2n.md": approvedGateEntity("2n", "validation", "done", "70"),
		"ax.md": approvedGateEntity("ax", "validation", "implementation", "65"),
		"qc.md": "---\nid: qc\nstatus: validation\nscore: 60\n---\n# incomplete\n",
	})
	args := []string{"--workflow-dir", root, "--fields", "id,gate-readiness", "--json"}
	out, errOut, code := runNative(t, root, pinnedEnv(t), args...)
	if code != 0 {
		t.Fatalf("status exit=%d stderr=%q", code, errOut)
	}
	var result struct {
		Entities []map[string]string `json:"entities"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"sp": "validating", "mf": "awaiting-captain", "r4": "awaiting-captain",
		"wd": "withdrawn-awaiting-prepare", "2n": "approved-awaiting-merge",
		"ax": "approved-awaiting-advance", "qc": "validating",
	}
	for _, entity := range result.Entities {
		id := entity["id"]
		if entity["gate-readiness"] != want[id] {
			t.Errorf("%s gate-readiness = %q, want %q", id, entity["gate-readiness"], want[id])
		}
	}

	text, errOut, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--fields", "gate-readiness")
	if code != 0 || !strings.Contains(text, "approved-awaiting-m") || !strings.Contains(text, "awaiting-captain") {
		t.Fatalf("human gate-readiness exit=%d stderr=%q output=%q", code, errOut, text)
	}
	all, errOut, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--all-fields", "--json")
	var allResult struct {
		Entities []map[string]string `json:"entities"`
	}
	if err := json.Unmarshal([]byte(all), &allResult); err != nil {
		t.Fatalf("decode --all-fields output: %v\n%s", err, all)
	}
	seen := map[string]bool{}
	for _, entity := range allResult.Entities {
		seen[entity["gate-readiness"]] = true
	}
	if code != 0 || !seen["validating"] || !seen["approved-awaiting-advance"] {
		t.Fatalf("--all-fields gate-readiness exit=%d stderr=%q output=%q", code, errOut, all)
	}
}

func TestStatusAllFieldsProjectsValidatingWithoutGateRecord(t *testing.T) {
	root, _ := buildSplitRoot(t, identifyReadyGatesReadme, map[string]string{
		"sp.md": "---\nstatus: validation\n---\n# still validating\n",
		"qc.md": "---\nstatus: validation\n---\n# still validating\n",
	})
	out, errOut, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--all-fields", "--json")
	if code != 0 {
		t.Fatalf("--all-fields exit=%d stderr=%q", code, errOut)
	}
	var result struct {
		Entities []map[string]string `json:"entities"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Entities) != 2 {
		t.Fatalf("entities = %#v, want two validating entities", result.Entities)
	}
	for _, entity := range result.Entities {
		if entity["gate-readiness"] != "validating" {
			t.Fatalf("%s gate-readiness = %q, want validating", entity["slug"], entity["gate-readiness"])
		}
	}
}

func TestUnrelatedSetPreservesGates(t *testing.T) {
	path := filepath.Join(t.TempDir(), "task.md")
	body := "---\nid: task\nstatus: ideation\nscore: '0.5'\ngates:\n  version: 1\n  records:\n    - id: gate:design\n      stage: ideation\n      attempts:\n        - id: attempt:design-1\n          briefing: {id: 'briefing:design-1', digest: 'sha256:" + strings.Repeat("1", 64) + "', room-ref: ./review}\n          resolution: {type: Resolution, id: 'resolution:design-1', briefing: 'briefing:design-1', by: 'person:captain', at: '2026-07-22T00:00:00Z', decision: hold, reason: 'wait'}\n          application: {action: none, state: not-applicable}\n---\n# Task\n"
	body = strings.Replace(body, "\n          application: {action: none, state: not-applicable}", "", 1)
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
