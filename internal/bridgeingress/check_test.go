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
	// The block reason must speak the direct _bridge/ file protocol (seam-contract
	// §3): name the inbox, the per-slug cursor, and the ack file — and the count.
	for _, want := range []string{"alpha", "1 queued", "_bridge/inbox.jsonl", "_bridge/.inbox-cursor.<slug>", "_bridge/fo-replies.jsonl"} {
		if !strings.Contains(d.Reason, want) {
			t.Fatalf("reason missing %q:\n%s", want, d.Reason)
		}
	}
}

// TestCheckReasonNamesNoDroppedVerb pins that the Stop-hook block reason never
// instructs a dropped `spacedock bridge inbox drain|ack|commit` (or alert /
// initiate) verb — FO judgment is direct file writes, not a CLI verb (scope 2).
func TestCheckReasonNamesNoDroppedVerb(t *testing.T) {
	root := t.TempDir()
	writeInbox(t, root, `{"id":"i1","ts":"2026-07-03T11:00:00Z","kind":"tell","target":"alpha"}`)
	writeHeartbeat(t, root, "alpha", "sess-a", time.Now().UTC())

	d := Check(CheckOptions{Host: "claude", Root: root, SessionID: "sess-a"})
	if d.Decision != "block" {
		t.Fatalf("decision = %q, want block", d.Decision)
	}
	for _, banned := range []string{
		"spacedock bridge inbox drain",
		"spacedock bridge inbox ack",
		"spacedock bridge inbox commit",
		"bridge alert",
		"bridge initiate",
	} {
		if strings.Contains(d.Reason, banned) {
			t.Fatalf("reason names dropped verb %q:\n%s", banned, d.Reason)
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

// TestCheckEmptyWhenHeartbeatMissingSessionID documents the review-B2 failure
// mode the mod protocol prevents: a heartbeat written WITHOUT a session_id
// resolves no slug for the stopping session, so the check never blocks and the
// queued intent sits undelivered. The load-bearing fix lives in the producer
// (the mod's startup/idle hooks MUST write the harness session id); this test
// pins that the reader honestly resolves nothing rather than fabricating a block.
func TestCheckEmptyWhenHeartbeatMissingSessionID(t *testing.T) {
	root := t.TempDir()
	now := time.Now().UTC()
	writeInbox(t, root, `{"id":"i1","ts":"2026-07-03T11:00:00Z","kind":"tell","target":"alpha"}`)
	writeHeartbeat(t, root, "alpha", "", now) // heartbeat with an EMPTY session_id

	d := Check(CheckOptions{Host: "claude", Root: root, SessionID: "sess-a"})
	if d.Decision != "" {
		t.Fatalf("missing heartbeat session_id must resolve nothing (no block), got %+v", d)
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
