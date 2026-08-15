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

func TestGateConsumeRejectsBadRetainedRequestBinding(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("---\nid-style: slug\nstages:\n  states:\n    - name: ideation\n      initial: true\n      gate: true\n    - name: implementation\n---\n# Workflow\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	entity := "---\nid: task\nstatus: ideation\ngates:\n  version: 1\n  records:\n    - id: gate:task:ideation\n      stage: ideation\n      attempts:\n        - id: attempt:task:ideation\n          briefing: {id: briefing:task:ideation:attempt-1:revision-1, digest: sha256:" + strings.Repeat("1", 64) + ", request-digest: sha256:" + strings.Repeat("2", 64) + ", room-ref: ./review/ideation/briefing-1}\n          resolution: {type: Resolution, id: resolution:task:ideation, briefing: briefing:task:ideation:attempt-1:revision-1, by: person:captain, at: 2026-07-22T00:00:00Z, decision: approve}\n          application: {target-stage: implementation, state: pending}\n---\n# Task\n"
	if err := os.WriteFile(filepath.Join(root, "task.md"), []byte(entity), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "consume", "task", "--workflow-dir", root}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 1 || !strings.Contains(errOut.String(), "retained request.json") {
		t.Fatalf("bad retained binding exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
}
