package bridgeingress

import (
	"strings"
	"testing"
	"time"
)

func TestCheckBlocksWhenIntentQueuedForSession(t *testing.T) {
	root := t.TempDir()
	now := time.Now().UTC()
	writeInbox(t, root, `{"id":"i1","ts":"2026-07-03T11:00:00Z","kind":"tell","text":"hi","target":"alpha"}`)
	writeHeartbeat(t, root, "alpha", "sess-a", now)

	d := Check(CheckOptions{Host: "claude", Root: root, SessionID: "sess-a"})
	if d.Decision != "block" {
		t.Fatalf("decision = %q, want block", d.Decision)
	}
	for _, want := range []string{"alpha", "spacedock bridge inbox drain", "1 queued"} {
		if !strings.Contains(d.Reason, want) {
			t.Fatalf("reason missing %q:\n%s", want, d.Reason)
		}
	}
}

func TestCheckEmptyWhenStopHookActive(t *testing.T) {
	root := t.TempDir()
	writeInbox(t, root, `{"id":"i1","ts":"2026-07-03T11:00:00Z","kind":"tell","target":"alpha"}`)
	writeHeartbeat(t, root, "alpha", "sess-a", time.Now().UTC())
	d := Check(CheckOptions{Host: "claude", Root: root, SessionID: "sess-a", StopHookActive: true})
	if d.Decision != "" {
		t.Fatalf("stop_hook_active must not block again, got %+v", d)
	}
}

func TestCheckEmptyWhenSessionUnknown(t *testing.T) {
	root := t.TempDir()
	writeInbox(t, root, `{"id":"i1","ts":"2026-07-03T11:00:00Z","kind":"tell","target":"alpha"}`)
	writeHeartbeat(t, root, "alpha", "sess-a", time.Now().UTC())
	// A different session id must not block for a sibling FO's intent.
	d := Check(CheckOptions{Host: "claude", Root: root, SessionID: "sess-OTHER"})
	if d.Decision != "" {
		t.Fatalf("unknown session must not block, got %+v", d)
	}
}

func TestCheckEmptyWhenNothingPending(t *testing.T) {
	root := t.TempDir()
	now := time.Now().UTC()
	writeInbox(t, root, `{"id":"i1","ts":"2026-07-03T11:00:00Z","kind":"tell","target":"alpha"}`)
	writeHeartbeat(t, root, "alpha", "sess-a", now)
	writeInboxCursor(t, root, "alpha", "1") // already drained
	d := Check(CheckOptions{Host: "claude", Root: root, SessionID: "sess-a"})
	if d.Decision != "" {
		t.Fatalf("no pending intent must not block, got %+v", d)
	}
}

func TestCheckFromReaderParsesStopPayload(t *testing.T) {
	root := t.TempDir()
	now := time.Now().UTC()
	writeInbox(t, root, `{"id":"i1","ts":"2026-07-03T11:00:00Z","kind":"tell","target":"alpha"}`)
	writeHeartbeat(t, root, "alpha", "sess-a", now)

	payload := `{"cwd":"` + root + `","session_id":"sess-a","stop_hook_active":false,"hook_event_name":"Stop"}`
	d := CheckFromReader(strings.NewReader(payload), CheckOptions{Host: "claude"})
	if d.Decision != "block" {
		t.Fatalf("reader path should block, got %+v", d)
	}

	// stop_hook_active in the payload suppresses a re-block.
	payload2 := `{"cwd":"` + root + `","session_id":"sess-a","stop_hook_active":true}`
	if d2 := CheckFromReader(strings.NewReader(payload2), CheckOptions{Host: "claude"}); d2.Decision != "" {
		t.Fatalf("payload stop_hook_active must suppress block, got %+v", d2)
	}
}
