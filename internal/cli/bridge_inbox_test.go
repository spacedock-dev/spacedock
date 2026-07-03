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

// TestBridgeInboxDrainAckCommitCLI drives the full FO-facing verb cycle through
// the hidden CLI to lock the argv grammar and JSON result shapes.
func TestBridgeInboxDrainAckCommitCLI(t *testing.T) {
	root := t.TempDir()
	bridgeDir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(bridgeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bridgeDir, "inbox.jsonl"),
		[]byte(`{"id":"i1","ts":"2026-07-03T11:00:00Z","kind":"tell","text":"go","target":"alpha"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	drain := runInbox(t, root, "drain", "--host", "claude", "--repo-root", root, "--slug", "alpha", "--session-id", "sess-a")
	if drain["status"] != "ok" || drain["count"].(float64) != 1 {
		t.Fatalf("drain = %v, want ok/count 1", drain)
	}

	ack := runInbox(t, root, "ack", "--repo-root", root, "--slug", "alpha", "--line", "1", "--id", "i1",
		"--ts", "2026-07-03T11:00:00Z", "--kind", "tell", "--status", "answered")
	if ack["appended"] != true || ack["kind"] != "reply" {
		t.Fatalf("ack = %v, want appended reply", ack)
	}

	commit := runInbox(t, root, "commit", "--repo-root", root, "--slug", "alpha", "--cursor", "1")
	if commit["cursor"].(float64) != 1 {
		t.Fatalf("commit = %v, want cursor 1", commit)
	}

	redrain := runInbox(t, root, "drain", "--host", "claude", "--repo-root", root, "--slug", "alpha")
	if redrain["count"].(float64) != 0 {
		t.Fatalf("re-drain = %v, want count 0", redrain)
	}
}

func TestBridgeInboxCheckReadsStdin(t *testing.T) {
	root := t.TempDir()
	bridgeDir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(bridgeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(bridgeDir, "inbox.jsonl"),
		[]byte(`{"id":"i1","ts":"2026-07-03T11:00:00Z","kind":"tell","target":"alpha"}`+"\n"), 0o644)
	os.WriteFile(filepath.Join(bridgeDir, "fo.alpha.json"),
		[]byte(`{"session_id":"sess-a","host":"claude","ts":"2026-07-03T11:00:00Z","state":"idle"}`), 0o644)

	var stdout, stderr bytes.Buffer
	payload := `{"cwd":"` + root + `","session_id":"sess-a","stop_hook_active":false}`
	code := run(context.Background(),
		[]string{"bridge", "inbox", "check", "--host", "claude"},
		nil, filepath.Join(root, "elsewhere"), strings.NewReader(payload), &stdout, &stderr, &fakeRunner{}, nil)
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	var got map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &got); err != nil {
		t.Fatalf("stdout JSON: %v\n%s", err, stdout.String())
	}
	if got["decision"] != "block" {
		t.Fatalf("decision = %v, want block", got["decision"])
	}
}

func runInbox(t *testing.T, root string, args ...string) map[string]any {
	t.Helper()
	var stdout, stderr bytes.Buffer
	full := append([]string{"bridge", "inbox"}, args...)
	code := run(context.Background(), full, nil, filepath.Join(root, "elsewhere"),
		strings.NewReader(""), &stdout, &stderr, &fakeRunner{}, nil)
	if code != 0 {
		t.Fatalf("%v exit = %d stderr=%s", args, code, stderr.String())
	}
	var got map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &got); err != nil {
		t.Fatalf("%v stdout JSON: %v\n%s", args, err, stdout.String())
	}
	return got
}
