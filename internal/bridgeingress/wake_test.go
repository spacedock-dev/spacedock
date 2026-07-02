package bridgeingress

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWakeResumesFreshHeartbeatSessionAndAdvancesCursor(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root,
		`{"id":"i1","ts":"2026-07-02T11:59:00Z","kind":"tell","target":"all","target_set":["a","b"]}`,
	)
	writeHeartbeat(t, root, "a", "session-a", now.Add(-time.Minute))
	writeHeartbeat(t, root, "b", "session-a", now.Add(-time.Minute))

	var gotSession, gotPrompt string
	res := Wake(context.Background(), Options{
		Host: "codex",
		Root: root,
		Now:  func() time.Time { return now },
		Resume: func(_ context.Context, sessionID, prompt string) error {
			gotSession = sessionID
			gotPrompt = prompt
			return nil
		},
	})

	if res.Status != "woke" || res.Sessions != 1 {
		t.Fatalf("result = %+v, want woke one session", res)
	}
	if gotSession != "session-a" {
		t.Fatalf("session = %q, want session-a", gotSession)
	}
	for _, want := range []string{"Pending physical inbox lines for this session: 1", "Addressed workflow slugs: a,b", "_bridge/inbox.jsonl"} {
		if !strings.Contains(gotPrompt, want) {
			t.Fatalf("prompt missing %q:\n%s", want, gotPrompt)
		}
	}
	if cursor := readFile(t, filepath.Join(root, "_bridge", ".wake-cursor.codex")); cursor != "1\n" {
		t.Fatalf("cursor = %q, want 1", cursor)
	}
	var ev wakeEvent
	readLastJSON(t, filepath.Join(root, "_bridge", "wake-events.jsonl"), &ev)
	if ev.Status != "woke" || ev.SessionID != "session-a" || len(ev.Targets) != 2 {
		t.Fatalf("event = %+v, want woke session with targets", ev)
	}
}

func TestWakeNoSessionDoesNotAdvanceCursor(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root,
		`{"id":"i1","ts":"2026-07-02T11:59:00Z","kind":"tell","target":"all","target_set":["a"]}`,
	)

	res := Wake(context.Background(), Options{
		Host: "codex",
		Root: root,
		Now:  func() time.Time { return now },
		Resume: func(context.Context, string, string) error {
			t.Fatal("resume should not run without a fresh heartbeat")
			return nil
		},
	})

	if res.Status != "skipped-no-session" {
		t.Fatalf("result = %+v, want skipped-no-session", res)
	}
	if _, err := os.Stat(filepath.Join(root, "_bridge", ".wake-cursor.codex")); !os.IsNotExist(err) {
		t.Fatalf("cursor should not exist after no-session wake: %v", err)
	}
}

func TestWakeSkipsAlreadyWokenLines(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root,
		`{"id":"i1","ts":"2026-07-02T11:58:00Z","kind":"tell","target":"a"}`,
		`{"id":"i2","ts":"2026-07-02T11:59:00Z","kind":"tell","target":"a"}`,
	)
	if err := writeCursor(root, 1); err != nil {
		t.Fatal(err)
	}
	writeHeartbeat(t, root, "a", "session-a", now.Add(-time.Minute))

	var prompt string
	res := Wake(context.Background(), Options{
		Host: "codex",
		Root: root,
		Now:  func() time.Time { return now },
		Resume: func(_ context.Context, _ string, p string) error {
			prompt = p
			return nil
		},
	})

	if res.Status != "woke" {
		t.Fatalf("result = %+v, want woke", res)
	}
	if !strings.Contains(prompt, "Pending physical inbox lines for this session: 2") || strings.Contains(prompt, "1,2") {
		t.Fatalf("prompt did not scope to pending line 2:\n%s", prompt)
	}
	if cursor := readFile(t, filepath.Join(root, "_bridge", ".wake-cursor.codex")); cursor != "2\n" {
		t.Fatalf("cursor = %q, want 2", cursor)
	}
}

func writeInbox(t *testing.T, root string, lines ...string) {
	t.Helper()
	dir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(dir, "inbox.jsonl"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeHeartbeat(t *testing.T, root, slug, sessionID string, ts time.Time) {
	t.Helper()
	dir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := `{"session_id":"` + sessionID + `","ts":"` + ts.Format(time.RFC3339) + `","state":"idle"}`
	if err := os.WriteFile(filepath.Join(dir, "fo."+slug+".json"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func readLastJSON(t *testing.T, path string, out any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), out); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
}
