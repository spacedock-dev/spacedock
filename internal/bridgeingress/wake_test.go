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

func TestWakeResumesFreshHeartbeatSessionWithoutAdvancingWakeCursor(t *testing.T) {
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
	if _, err := os.Stat(filepath.Join(root, "_bridge", ".wake-cursor.codex")); !os.IsNotExist(err) {
		t.Fatalf("wake cursor should not be written after resume launch: %v", err)
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

func TestWakeRetriesAlreadyStartedUndeliveredLines(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root,
		`{"id":"i1","ts":"2026-07-02T11:58:00Z","kind":"tell","target":"a"}`,
		`{"id":"i2","ts":"2026-07-02T11:59:00Z","kind":"tell","target":"a"}`,
	)
	writeWakeCursor(t, root, "1")
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
	if !strings.Contains(prompt, "Pending physical inbox lines for this session: 1,2") {
		t.Fatalf("prompt did not retry undelivered lines despite prior wake cursor:\n%s", prompt)
	}
}

func TestWakeSkipsDeliveredByInboxCursor(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root,
		`{"id":"i1","ts":"2026-07-02T11:58:00Z","kind":"tell","target":"a"}`,
	)
	writeInboxCursor(t, root, "a", "1")
	writeHeartbeat(t, root, "a", "session-a", now.Add(-time.Minute))

	res := Wake(context.Background(), Options{
		Host: "codex",
		Root: root,
		Now:  func() time.Time { return now },
		Resume: func(context.Context, string, string) error {
			t.Fatal("resume should not run for cursor-delivered line")
			return nil
		},
	})

	if res.Status != "noop" {
		t.Fatalf("result = %+v, want noop", res)
	}
}

func TestWakeSkipsDeliveredByReplyAck(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root,
		`{"id":"i1","ts":"2026-07-02T11:58:00Z","kind":"decision","target":"a"}`,
	)
	writeReplies(t, root,
		`{"schema":1,"ts":"2026-07-02T12:00:00Z","kind":"decision-ack","target":"a","in_reply_to_id":"i1","in_reply_to_line":1,"in_reply_to_ts":"2026-07-02T11:58:00Z","intent_kind":"decision","status":"applied"}`,
	)
	writeHeartbeat(t, root, "a", "session-a", now.Add(-time.Minute))

	res := Wake(context.Background(), Options{
		Host: "codex",
		Root: root,
		Now:  func() time.Time { return now },
		Resume: func(context.Context, string, string) error {
			t.Fatal("resume should not run for ack-delivered line")
			return nil
		},
	})

	if res.Status != "noop" {
		t.Fatalf("result = %+v, want noop", res)
	}
}

func TestWakeOnlyTargetsUndeliveredMembers(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root,
		`{"id":"i1","ts":"2026-07-02T11:58:00Z","kind":"tell","target":"all","target_set":["a","b"]}`,
	)
	writeInboxCursor(t, root, "a", "1")
	writeHeartbeat(t, root, "a", "session-a", now.Add(-time.Minute))
	writeHeartbeat(t, root, "b", "session-b", now.Add(-time.Minute))

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

	if res.Status != "woke" || res.Sessions != 1 || gotSession != "session-b" {
		t.Fatalf("result=%+v session=%q, want only pending target b", res, gotSession)
	}
	if strings.Contains(gotPrompt, "Addressed workflow slugs: a") || !strings.Contains(gotPrompt, "Addressed workflow slugs: b") {
		t.Fatalf("prompt should name only pending target b:\n%s", gotPrompt)
	}
}

func TestWakeStaleHeartbeatWithSessionIsResumable(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root,
		`{"id":"i1","ts":"2026-07-02T11:59:00Z","kind":"tell","target":"a"}`,
	)
	writeHeartbeat(t, root, "a", "session-a", now.Add(-2*time.Hour))

	var gotSession string
	res := Wake(context.Background(), Options{
		Host: "codex",
		Root: root,
		Now:  func() time.Time { return now },
		Resume: func(_ context.Context, sessionID, _ string) error {
			gotSession = sessionID
			return nil
		},
	})

	if res.Status != "woke" || gotSession != "session-a" {
		t.Fatalf("result=%+v session=%q, want stale heartbeat resumable", res, gotSession)
	}
}

func TestWakeSessionMarkerWithoutHeartbeatIsResumable(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root,
		`{"id":"i1","ts":"2026-07-02T11:59:00Z","kind":"tell","target":"a"}`,
	)
	writeSessionMarker(t, root, "session-a", "a")

	var gotSession string
	res := Wake(context.Background(), Options{
		Host: "codex",
		Root: root,
		Now:  func() time.Time { return now },
		Resume: func(_ context.Context, sessionID, _ string) error {
			gotSession = sessionID
			return nil
		},
	})

	if res.Status != "woke" || gotSession != "session-a" {
		t.Fatalf("result=%+v session=%q, want session marker resumable", res, gotSession)
	}
}

func TestWakeMultipleTargetsCoalesceOneSession(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root,
		`{"id":"i1","ts":"2026-07-02T11:59:00Z","kind":"tell","target":"all","target_set":["a","b"]}`,
	)
	writeHeartbeat(t, root, "a", "session-a", now.Add(-time.Minute))
	writeHeartbeat(t, root, "b", "session-a", now.Add(-time.Minute))

	var calls int
	var prompt string
	res := Wake(context.Background(), Options{
		Host: "codex",
		Root: root,
		Now:  func() time.Time { return now },
		Resume: func(_ context.Context, _ string, p string) error {
			calls++
			prompt = p
			return nil
		},
	})

	if res.Status != "woke" || res.Sessions != 1 || calls != 1 {
		t.Fatalf("result=%+v calls=%d, want one resumed session", res, calls)
	}
	if !strings.Contains(prompt, "Addressed workflow slugs: a,b") {
		t.Fatalf("prompt missing coalesced targets:\n%s", prompt)
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

func writeInboxCursor(t *testing.T, root, slug, content string) {
	t.Helper()
	dir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".inbox-cursor."+slug), []byte(content+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeWakeCursor(t *testing.T, root, content string) {
	t.Helper()
	dir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".wake-cursor.codex"), []byte(content+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeReplies(t *testing.T, root string, lines ...string) {
	t.Helper()
	dir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(dir, "fo-replies.jsonl"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeSessionMarker(t *testing.T, root, sessionID, workflow string) {
	t.Helper()
	dir := filepath.Join(root, "_bridge", "sessions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := `{"session_id":"` + sessionID + `","workflow":"` + workflow + `"}`
	if err := os.WriteFile(filepath.Join(dir, sessionID+".json"), []byte(body), 0o600); err != nil {
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
