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
