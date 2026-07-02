package bridgealert

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAppendPermissionWritesAlert(t *testing.T) {
	root := t.TempDir()
	got, err := AppendPermission(PermissionOptions{
		Root:       root,
		Now:        func() time.Time { return time.Date(2026, 7, 2, 5, 0, 0, 0, time.UTC) },
		ID:         "perm-1",
		Workflow:   "pr-review-queue",
		Entity:     "datarecce-recce-pr-1",
		Host:       "codex",
		SessionID:  "s1",
		Reason:     "sandbox blocked state gitdir",
		Command:    "git status",
		PrefixRule: []string{"git", "-C"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "perm-1" || !got.Queued {
		t.Fatalf("result = %+v", got)
	}
	data, err := os.ReadFile(filepath.Join(root, "_bridge", "fo-alerts.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var alert PermissionAlert
	if err := json.Unmarshal(bytes.TrimSpace(data), &alert); err != nil {
		t.Fatalf("alert JSON: %v\n%s", err, data)
	}
	if alert.Kind != "permission-request" || alert.Severity != "blocked" || alert.Status != "open" {
		t.Fatalf("alert = %+v", alert)
	}
	if alert.TS != "2026-07-02T05:00:00Z" || alert.Workflow != "pr-review-queue" || alert.Command != "git status" {
		t.Fatalf("alert fields = %+v", alert)
	}
	if len(alert.PrefixRule) != 2 || alert.PrefixRule[0] != "git" {
		t.Fatalf("prefix rule = %+v", alert.PrefixRule)
	}
}

func TestAppendPermissionRequiresReason(t *testing.T) {
	if _, err := AppendPermission(PermissionOptions{Root: t.TempDir()}); err == nil {
		t.Fatal("expected missing reason to fail")
	}
}
