package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func gateApplicationWarningFixture(t *testing.T, application string) string {
	t.Helper()
	root := t.TempDir()
	readme := "---\nid-style: slug\nstages:\n  states:\n    - name: ideation\n      initial: true\n      gate: true\n    - name: implementation\n---\n# Workflow\n"
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	entity := "---\nid: task\nstatus: ideation\ntitle: Task\ngates:\n  version: 1\n  records:\n    - id: gate:task:ideation\n      stage: ideation\n      attempts:\n        - id: attempt:task:ideation\n          briefing: {id: briefing:task:ideation:attempt-1:revision-1, digest: sha256:" + strings.Repeat("1", 64) + ", room-ref: ./review}\n          resolution: {type: Resolution, id: resolution:task:ideation, briefing: briefing:task:ideation:attempt-1:revision-1, by: person:captain, at: 2026-07-22T00:00:00Z, decision: approve}\n          application:\n" + application + "\n---\n# Task\n"
	if err := os.WriteFile(filepath.Join(root, "task.md"), []byte(entity), 0o644); err != nil {
		t.Fatal(err)
	}
	gitInit(t, root)
	return root
}

func TestStatusValidateReportsGateApplicationExtensionsAsWarnings(t *testing.T) {
	root := gateApplicationWarningFixture(t, "            target-stage: implementation\n            state: pending\n            action: advance\n            blockers: []\n            execution-hold: {state: active}\n            feedback: {owner: old-producer}\n            nested: [legacy]")
	out, stderr, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--validate")
	if code != 0 || strings.TrimSpace(out) != "VALID" {
		t.Fatalf("--validate exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	for _, field := range []string{"action", "blockers", "execution-hold", "feedback", "nested"} {
		if !strings.Contains(stderr, "Warning: unknown gate application field '"+field+"'") {
			t.Fatalf("--validate missing warning for %s: %q", field, stderr)
		}
	}
	if strings.Count(stderr, "Warning: unknown gate application field ") != 5 {
		t.Fatalf("--validate warning count = %d, stderr=%q", strings.Count(stderr, "Warning: unknown gate application field "), stderr)
	}

	ordinary, ordinaryErr, ordinaryCode := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--fields", "id,gate-application,gate-condition")
	if ordinaryCode != 0 || !strings.Contains(ordinary, "advance/pending") || strings.Contains(ordinaryErr, "unknown gate application field") {
		t.Fatalf("ordinary status changed warning surface: exit=%d stdout=%q stderr=%q", ordinaryCode, ordinary, ordinaryErr)
	}
}

func TestStatusValidateKeepsUnrelatedGateUnknownFieldsFatal(t *testing.T) {
	root := gateApplicationWarningFixture(t, "            target-stage: implementation\n            state: pending\n            action: advance\n          unrelated: true")
	_, stderr, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--validate")
	if code != 1 || !strings.Contains(stderr, "Error: invalid gates:") {
		t.Fatalf("unrelated gate key exit=%d stderr=%q", code, stderr)
	}
}
