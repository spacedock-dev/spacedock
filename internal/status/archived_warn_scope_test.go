// ABOUTME: --validate warn tiers skip archived scope (publish-only, unfixable),
// ABOUTME: while archived structural gate errors stay fatal.
package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// archivedEntityFixture plants a single entity under `_archive/` so discovery
// assigns it archived scope, with the caller's gate application block and
// frontmatter tail spliced in.
func archivedEntityFixture(t *testing.T, frontmatterTail, application string) string {
	t.Helper()
	root := t.TempDir()
	readme := "---\nid-style: slug\nstages:\n  states:\n    - name: ideation\n      initial: true\n      gate: true\n    - name: implementation\n---\n# Workflow\n"
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "_archive"), 0o755); err != nil {
		t.Fatal(err)
	}
	entity := "---\nid: task\nstatus: ideation\ntitle: Task\n" + frontmatterTail + "gates:\n  version: 1\n  records:\n    - id: gate:task:ideation\n      stage: ideation\n      attempts:\n        - id: attempt:task:ideation\n          briefing: {id: briefing:task:ideation:attempt-1:revision-1, digest: sha256:" + strings.Repeat("1", 64) + ", room-ref: ./review}\n          resolution: {type: Resolution, id: resolution:task:ideation, briefing: briefing:task:ideation:attempt-1:revision-1, by: person:captain, at: 2026-07-22T00:00:00Z, decision: approve}\n          application:\n" + application + "\n---\n# Task\n"
	if err := os.WriteFile(filepath.Join(root, "_archive", "task.md"), []byte(entity), 0o644); err != nil {
		t.Fatal(err)
	}
	gitInit(t, root)
	return root
}

// Archived scope is publish-only, so neither warn channel can ever be cleared by
// a tool-mediated write there. Asserting on the `scope=archived` evidence token
// rather than each message text covers both channels at once — and any warn
// channel added later. The fixture carries one defect per channel, so this fails
// if either guard is dropped or its predicate inverted.
func TestValidateSkipsArchivedWarnChannels(t *testing.T) {
	root := archivedEntityFixture(t, "verdict: MAYBE\n", "            target-stage: implementation\n            state: consumed\n            action: advance\n            blockers: []")
	out, stderr, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--validate")
	if code != 0 || strings.TrimSpace(out) != "VALID" {
		t.Fatalf("--validate exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	if strings.Contains(stderr, "scope=archived") {
		t.Fatalf("--validate reported archived-scope findings: %q", stderr)
	}
}

// The guards sit at the two warn emission points, not on the entity slice at the
// call site: `gateValidationDiagnostics` returns errors and warnings from one
// loop, so a slice-level filter would take this error down with the warnings and
// turn exit 1 into exit 0.
func TestValidateKeepsArchivedGateErrorsFatal(t *testing.T) {
	root := archivedEntityFixture(t, "", "            target-stage: implementation\n            state: consumed\n          unrelated: true")
	_, stderr, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--validate")
	if code != 1 || !strings.Contains(stderr, "Error: invalid gates:") {
		t.Fatalf("archived structural gate error exit=%d stderr=%q", code, stderr)
	}
}
