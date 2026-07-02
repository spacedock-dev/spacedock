package bridgealert

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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
	if got.RequestID != "perm-1" {
		t.Fatalf("request id = %q, want alias", got.RequestID)
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

func TestAppendPermissionMissingReasonIsNonBlocking(t *testing.T) {
	got, err := AppendPermission(PermissionOptions{Root: t.TempDir(), ID: "perm-1"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Queued || got.Error == "" {
		t.Fatalf("result = %+v, want non-queued error result", got)
	}
}

func TestAppendPermissionWriteFailureIsNonBlocking(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "_bridge"), []byte("not a dir"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := AppendPermission(PermissionOptions{Root: root, ID: "perm-1", Reason: "sandbox blocked"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Queued || got.Error == "" {
		t.Fatalf("result = %+v, want non-queued error result", got)
	}
}

func TestAppendPermissionNormalizesOneLineSummaries(t *testing.T) {
	root := t.TempDir()
	got, err := AppendPermission(PermissionOptions{
		Root:    root,
		ID:      "perm-1",
		Reason:  "sandbox\nblocked\tstate gitdir",
		Command: "git status\n" + strings.Repeat("x", 800),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !got.Queued {
		t.Fatalf("result = %+v, want queued", got)
	}
	data, err := os.ReadFile(filepath.Join(root, "_bridge", "fo-alerts.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var alert PermissionAlert
	if err := json.Unmarshal(bytes.TrimSpace(data), &alert); err != nil {
		t.Fatal(err)
	}
	if alert.Reason != "sandbox blocked state gitdir" {
		t.Fatalf("reason = %q", alert.Reason)
	}
	if strings.ContainsAny(alert.Command, "\n\t") || len(alert.Command) > maxCommandLen {
		t.Fatalf("command = %q len=%d", alert.Command, len(alert.Command))
	}
}
