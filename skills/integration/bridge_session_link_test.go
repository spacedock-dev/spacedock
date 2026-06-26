// ABOUTME: Bridge session-marker smoke — the event hook deterministically records
// ABOUTME: a dispatched ensign's session→entity(+workflow) link from its entity-file
// ABOUTME: Read, so Bridge's "running" badge no longer depends on the ensign running
// ABOUTME: a first-action shell (which it skipped ~3/4 of the time). DRC running-badge.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// runEventHook feeds one hook payload (the JSON Claude Code pipes to the hook) to
// scripts/spacedock-bridge-events.sh and returns once it exits. cwd anchors the
// _bridge/ dir the hook writes to.
func runEventHook(t *testing.T, payload string) {
	t.Helper()
	hook := filepath.Join("..", "..", "scripts", "spacedock-bridge-events.sh")
	if _, err := os.Stat(hook); err != nil {
		t.Fatalf("hook script not found at %s: %v", hook, err)
	}
	cmd := exec.Command("bash", hook)
	cmd.Stdin = strings.NewReader(payload)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("event hook failed: %v\n%s", err, out)
	}
}

func readPayload(cwd, sid, agentType, filePath string) string {
	return `{"cwd":"` + cwd + `","session_id":"` + sid + `","agent_type":"` + agentType +
		`","hook_event_name":"PostToolUse","tool_name":"Read","tool_input":{"file_path":"` + filePath + `"}}`
}

// TestEventHookWritesSessionMarker locks the deterministic producer: an ensign's
// Read of its entity file records {session_id, entity, workflow} — carrying the
// workflow so Bridge's join is collision-free — and it is first-write-wins (a later
// sibling Read does not overwrite), while a non-ensign (FO) Read writes nothing.
func TestEventHookWritesSessionMarker(t *testing.T) {
	root := t.TempDir()
	ent := filepath.Join(root, "docs", "spacedock", "linear-drc-review", "drc-3467.md")

	// 1. Ensign reads its own entity file → marker written with entity + workflow.
	runEventHook(t, readPayload(root, "ses-1", "spacedock:ensign", ent))
	marker := filepath.Join(root, "_bridge", "sessions", "ses-1.json")
	data, err := os.ReadFile(marker)
	if err != nil {
		t.Fatalf("marker not written: %v", err)
	}
	got := string(data)
	for _, want := range []string{`"session_id":"ses-1"`, `"entity":"drc-3467"`, `"workflow":"linear-drc-review"`} {
		if !strings.Contains(got, want) {
			t.Errorf("marker missing %s\ngot: %s", want, got)
		}
	}

	// 2. First-write-wins: a later sibling Read in the same session must NOT overwrite.
	sibling := filepath.Join(root, "docs", "spacedock", "linear-drc-review", "drc-9999.md")
	runEventHook(t, readPayload(root, "ses-1", "spacedock:ensign", sibling))
	data2, _ := os.ReadFile(marker)
	if !strings.Contains(string(data2), `"entity":"drc-3467"`) {
		t.Errorf("sibling Read overwrote the marker (lost first-write-wins):\n%s", data2)
	}

	// 3. A non-ensign (FO) Read of an entity file writes no marker.
	runEventHook(t, readPayload(root, "fo-sess", "spacedock:first-officer", ent))
	if _, err := os.Stat(filepath.Join(root, "_bridge", "sessions", "fo-sess.json")); err == nil {
		t.Errorf("FO Read should not produce a session marker")
	}

	// 3b. RELATIVE entity path — the FO passes a repo-relative {entity_file_path}, so
	// the ensign's scoped Read carries "docs/spacedock/<wf>/<slug>.md" (no leading
	// slash). The hook must still record it (regression: the absolute-only pattern
	// missed every live ensign).
	runEventHook(t, readPayload(root, "ses-rel", "spacedock:ensign",
		"docs/spacedock/linear-drc-review/drc-7000.md"))
	relData, err := os.ReadFile(filepath.Join(root, "_bridge", "sessions", "ses-rel.json"))
	if err != nil {
		t.Fatalf("relative-path Read produced no marker (the live-ensign regression): %v", err)
	}
	for _, want := range []string{`"entity":"drc-7000"`, `"workflow":"linear-drc-review"`} {
		if !strings.Contains(string(relData), want) {
			t.Errorf("relative-path marker missing %s\ngot: %s", want, relData)
		}
	}

	// 4. _archive entity Reads are skipped (workflow dir starting with _ is rejected).
	arch := filepath.Join(root, "docs", "spacedock", "linear-drc-review", "_archive", "drc-1.md")
	runEventHook(t, readPayload(root, "ses-arch", "spacedock:ensign", arch))
	if _, err := os.Stat(filepath.Join(root, "_bridge", "sessions", "ses-arch.json")); err == nil {
		t.Errorf("_archive Read should not produce a session marker")
	}
}
