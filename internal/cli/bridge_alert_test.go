package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBridgeAlertPermissionHiddenCLIWritesAlert(t *testing.T) {
	root := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := run(context.Background(),
		[]string{
			"bridge", "alert", "permission",
			"--repo-root", root,
			"--id", "perm-1",
			"--host", "codex",
			"--workflow", "pr-review-queue",
			"--entity", "ship-1",
			"--session-id", "s1",
			"--reason", "sandbox blocked state gitdir",
			"--command", "git status",
			"--prefix-rule", "git,-C",
		},
		nil, filepath.Join(root, "elsewhere"), strings.NewReader(""), &stdout, &stderr, &fakeRunner{}, nil)

	if code != 0 {
		t.Fatalf("exit = %d, stderr=%q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	var result struct {
		ID     string `json:"id"`
		Queued bool   `json:"queued"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &result); err != nil {
		t.Fatalf("stdout JSON: %v\n%s", err, stdout.String())
	}
	if result.ID != "perm-1" || !result.Queued {
		t.Fatalf("result = %+v", result)
	}
	data, err := os.ReadFile(filepath.Join(root, "_bridge", "fo-alerts.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"reason":"sandbox blocked state gitdir"`) ||
		!strings.Contains(string(data), `"prefix_rule":["git","-C"]`) {
		t.Fatalf("alert file missing fields:\n%s", data)
	}
}
