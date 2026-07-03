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

func TestDrainReturnsAddressedRecordsAndStampsHostHeartbeat(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 3, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root,
		`{"id":"i1","ts":"2026-07-03T11:00:00Z","kind":"tell","text":"do the thing","target_set":["alpha"]}`,
		`{"id":"i2","ts":"2026-07-03T11:01:00Z","kind":"tell","text":"not yours","target_set":["beta"]}`,
		`{"id":"i3","ts":"2026-07-03T11:02:00Z","kind":"decision","entity":"e1","field":"status","value":"approved","target":"alpha"}`,
	)

	res := Drain(DrainOptions{Host: "claude", Root: root, Slug: "alpha", SessionID: "sess-1", Now: func() time.Time { return now }})

	if res.Status != "ok" {
		t.Fatalf("status = %q, want ok", res.Status)
	}
	if !res.Heartbeat {
		t.Fatal("heartbeat not written")
	}
	if res.Cursor != 0 || res.HighWater != 3 {
		t.Fatalf("cursor=%d high_water=%d, want 0/3", res.Cursor, res.HighWater)
	}
	if res.Count != 2 || len(res.Records) != 2 {
		t.Fatalf("count=%d, want 2 addressed records; got %+v", res.Count, res.Records)
	}
	if res.Records[0].ID != "i1" || res.Records[0].Text != "do the thing" || res.Records[0].Line != 1 {
		t.Fatalf("record[0] = %+v, want i1 line 1", res.Records[0])
	}
	if res.Records[1].ID != "i3" || res.Records[1].Kind != "decision" || res.Records[1].Field != "status" || res.Records[1].Value != "approved" {
		t.Fatalf("record[1] = %+v, want decision i3 with field/value", res.Records[1])
	}

	// Heartbeat carries the host so Bridge can route wake per harness (WS-4).
	var hb heartbeatOut
	data, err := os.ReadFile(filepath.Join(root, "_bridge", "fo.alpha.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &hb); err != nil {
		t.Fatal(err)
	}
	if hb.Host != "claude" || hb.SessionID != "sess-1" || hb.State != "idle" {
		t.Fatalf("heartbeat = %+v, want claude/sess-1/idle", hb)
	}

	// Drain never advances the cursor.
	if _, err := os.Stat(filepath.Join(root, "_bridge", ".inbox-cursor.alpha")); !os.IsNotExist(err) {
		t.Fatalf("drain must not write a cursor: %v", err)
	}
}

func TestDrainPreservesRawTimestamp(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 3, 12, 0, 0, 0, time.UTC)
	// A non-UTC offset ts: drain must round-trip it verbatim, not reformat to UTC,
	// or wake's ts-based replyKey fallback (id-less records) would mismatch.
	writeInbox(t, root, `{"id":"i1","ts":"2026-07-03T20:00:00+08:00","kind":"tell","text":"hi","target":"alpha"}`)
	res := Drain(DrainOptions{Host: "claude", Root: root, Slug: "alpha", Now: func() time.Time { return now }})
	if res.Count != 1 {
		t.Fatalf("want 1 record, got %+v", res.Records)
	}
	if res.Records[0].TS != "2026-07-03T20:00:00+08:00" {
		t.Fatalf("ts = %q, want the verbatim on-disk offset ts", res.Records[0].TS)
	}
}

func TestDrainRespectsCursorAndAckIdempotency(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 3, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root,
		`{"id":"i1","ts":"2026-07-03T11:00:00Z","kind":"tell","text":"first","target":"alpha"}`,
		`{"id":"i2","ts":"2026-07-03T11:01:00Z","kind":"tell","text":"second","target":"alpha"}`,
	)
	// Cursor past line 1: only line 2 remains.
	writeInboxCursor(t, root, "alpha", "1")

	res := Drain(DrainOptions{Host: "claude", Root: root, Slug: "alpha", Now: func() time.Time { return now }})
	if res.Count != 1 || res.Records[0].ID != "i2" {
		t.Fatalf("with cursor=1 want only i2, got %+v", res.Records)
	}

	// Ack line 2 without committing the cursor: a re-drain must not resurface it.
	if ar := Ack(AckOptions{Host: "claude", Root: root, Slug: "alpha", Line: 2, ID: "i2", TS: "2026-07-03T11:01:00Z", IntentKind: "tell", Status: "answered"}); !ar.Appended {
		t.Fatalf("ack failed: %+v", ar)
	}
	res2 := Drain(DrainOptions{Host: "claude", Root: root, Slug: "alpha", Now: func() time.Time { return now }})
	if res2.Count != 0 {
		t.Fatalf("acked record must not resurface (idempotency), got %+v", res2.Records)
	}
}

func TestAckIsCompactSingleLineAndMapsKind(t *testing.T) {
	root := t.TempDir()
	cases := []struct{ intent, wantKind string }{
		{"tell", "reply"},
		{"conn", "conn-ack"},
		{"decision", "decision-ack"},
		{"permission-decision", "permission-ack"},
	}
	for _, c := range cases {
		res := Ack(AckOptions{Host: "claude", Root: root, Slug: "alpha", Line: 1, ID: "x", TS: "2026-07-03T11:00:00Z", IntentKind: c.intent, Status: "applied", Text: "multi\nline\ttext"})
		if !res.Appended || res.Kind != c.wantKind {
			t.Fatalf("intent %q -> %+v, want kind %q", c.intent, res, c.wantKind)
		}
	}

	raw, err := os.ReadFile(filepath.Join(root, "_bridge", "fo-replies.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(string(raw), "\n"), "\n")
	if len(lines) != len(cases) {
		t.Fatalf("want %d compact lines, got %d:\n%s", len(cases), len(lines), raw)
	}
	for _, line := range lines {
		var rec map[string]any
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			t.Fatalf("line is not valid JSON: %q (%v)", line, err)
		}
		if v, _ := rec["schema"].(float64); v != 1 {
			t.Fatalf("schema != 1 in %q", line)
		}
		if strings.Contains(line, "\n") || strings.Contains(line, "\t") {
			t.Fatalf("reply line not single-line/compact: %q", line)
		}
	}
	// The multi-line --text must be flattened to a single line.
	if strings.Contains(string(raw), "multi\nline") {
		t.Fatalf("text was not flattened:\n%s", raw)
	}
}

func TestAckOmitsGrantedWhenAbsent(t *testing.T) {
	root := t.TempDir()
	Ack(AckOptions{Host: "claude", Root: root, Slug: "alpha", Line: 1, ID: "x", IntentKind: "tell", Status: "answered"})
	raw, _ := os.ReadFile(filepath.Join(root, "_bridge", "fo-replies.jsonl"))
	if strings.Contains(string(raw), "granted") {
		t.Fatalf("absent --granted must be omitted, got %s", raw)
	}
	tr := true
	Ack(AckOptions{Host: "claude", Root: root, Slug: "alpha", Line: 2, ID: "y", IntentKind: "conn", Status: "accepted", Granted: &tr})
	raw2, _ := os.ReadFile(filepath.Join(root, "_bridge", "fo-replies.jsonl"))
	if !strings.Contains(string(raw2), `"granted":true`) {
		t.Fatalf("present --granted must be serialized, got %s", raw2)
	}
}

func TestAck_NoCorrelatorFailsLoudly(t *testing.T) {
	root := t.TempDir()
	// No id AND no ts: the reply has no strong correlator, so it must fail loudly at
	// write time rather than append an orphan the Bridge reader would silently drop.
	ar := Ack(AckOptions{Host: "claude", Root: root, Slug: "alpha", Line: 1, IntentKind: "tell", Status: "answered"})
	if ar.Appended {
		t.Fatalf("no-correlator ack must not append, got %+v", ar)
	}
	if ar.Error == "" {
		t.Fatalf("no-correlator ack must return an error, got %+v", ar)
	}
	// fo-replies.jsonl must not exist / must be untouched.
	if _, err := os.Stat(filepath.Join(root, "_bridge", "fo-replies.jsonl")); !os.IsNotExist(err) {
		t.Fatalf("no-correlator ack must not write fo-replies.jsonl: %v", err)
	}
}

func TestAck_IdOnlyReplyIsCorrelatable(t *testing.T) {
	root := t.TempDir()
	// id present, ts empty: id alone is a strong correlator, so the ack appends.
	ar := Ack(AckOptions{Host: "claude", Root: root, Slug: "alpha", Line: 2, ID: "i2", IntentKind: "tell", Status: "answered"})
	if !ar.Appended || ar.Error != "" {
		t.Fatalf("id-only ack must append, got %+v", ar)
	}
	raw, err := os.ReadFile(filepath.Join(root, "_bridge", "fo-replies.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var rec replyOut
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(raw))), &rec); err != nil {
		t.Fatalf("reply line not valid JSON: %q (%v)", raw, err)
	}
	if rec.InReplyToID != "i2" {
		t.Fatalf("in_reply_to_id = %q, want i2", rec.InReplyToID)
	}
	if rec.InReplyToTS != "" {
		t.Fatalf("empty --ts must echo as empty in_reply_to_ts, got %q", rec.InReplyToTS)
	}
}

func TestDrain_EmitsActingAck(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 3, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root, `{"id":"i1","ts":"2026-07-03T11:00:00Z","kind":"tell","text":"do it","target":"alpha"}`)

	res := Drain(DrainOptions{Host: "claude", Root: root, Slug: "alpha", SessionID: "sess-1", Now: func() time.Time { return now }})
	if res.Count != 1 {
		t.Fatalf("want 1 drained record, got %+v", res.Records)
	}
	dr := res.Records[0]

	raw, err := os.ReadFile(filepath.Join(root, "_bridge", "fo-replies.jsonl"))
	if err != nil {
		t.Fatalf("drain must auto-write an acting ack: %v", err)
	}
	lines := strings.Split(strings.TrimRight(string(raw), "\n"), "\n")
	if len(lines) != 1 {
		t.Fatalf("want exactly one acting ack, got %d:\n%s", len(lines), raw)
	}
	var rec replyOut
	if err := json.Unmarshal([]byte(lines[0]), &rec); err != nil {
		t.Fatalf("acting ack not valid JSON: %q (%v)", lines[0], err)
	}
	if rec.Status != "acting" {
		t.Fatalf("status = %q, want acting", rec.Status)
	}
	// The acting ack matches the drained record's id/line/ts/kind.
	if rec.InReplyToID != dr.ID || rec.InReplyToLine != dr.Line || rec.InReplyToTS != dr.TS {
		t.Fatalf("acting ack correlators = %q/%d/%q, want %q/%d/%q", rec.InReplyToID, rec.InReplyToLine, rec.InReplyToTS, dr.ID, dr.Line, dr.TS)
	}
	if rec.IntentKind != dr.Kind || rec.Kind != "reply" {
		t.Fatalf("acting ack kinds = intent %q / reply %q, want %q / reply", rec.IntentKind, rec.Kind, dr.Kind)
	}
}

func TestCommitAdvancesCursorMonotonically(t *testing.T) {
	root := t.TempDir()
	if r := Commit(CommitOptions{Root: root, Slug: "alpha", Cursor: 5}); r.Cursor != 5 {
		t.Fatalf("commit = %+v, want 5", r)
	}
	if got := inboxCursor(root, "alpha"); got != 5 {
		t.Fatalf("cursor file = %d, want 5", got)
	}
	// A stale/lower commit must not regress the cursor.
	if r := Commit(CommitOptions{Root: root, Slug: "alpha", Cursor: 2}); r.Cursor != 5 {
		t.Fatalf("regressing commit = %+v, want held at 5", r)
	}
	if got := inboxCursor(root, "alpha"); got != 5 {
		t.Fatalf("cursor regressed to %d, want 5", got)
	}
}

func TestDrainInvalidSlugRejected(t *testing.T) {
	root := t.TempDir()
	if r := Drain(DrainOptions{Host: "claude", Root: root, Slug: "../escape"}); r.Status != "failed" {
		t.Fatalf("unsafe slug must fail, got %+v", r)
	}
}

// TestPackagedDrainAckCommitSatisfiesWakeGuard is the load-bearing integration:
// after the packaged drain/ack/commit cycle, the Codex wake path (which owns the
// double-delivery guard) must treat the intent as delivered and NOT re-wake. This
// proves the packaged ack is byte-compatible with wake's replyKey/loadReplies.
func TestPackagedDrainAckCommitSatisfiesWakeGuard(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 7, 3, 12, 0, 0, 0, time.UTC)
	writeInbox(t, root, `{"id":"i1","ts":"2026-07-03T11:00:00Z","kind":"tell","text":"hello","target":"alpha"}`)
	writeHeartbeat(t, root, "alpha", "sess-a", now)

	// FO drains, acts, acks (without committing yet): the reply alone must satisfy
	// the wake guard even before the cursor advances.
	res := Drain(DrainOptions{Host: "codex", Root: root, Slug: "alpha", Now: func() time.Time { return now }})
	if res.Count != 1 {
		t.Fatalf("want 1 drained record, got %+v", res.Records)
	}
	rec := res.Records[0]
	if ar := Ack(AckOptions{Host: "codex", Root: root, Slug: "alpha", Line: rec.Line, ID: rec.ID, TS: rec.TS, IntentKind: rec.Kind, Status: "answered"}); !ar.Appended {
		t.Fatalf("ack failed: %+v", ar)
	}

	woke := false
	wr := Wake(context.Background(), Options{Host: "codex", Root: root, Now: func() time.Time { return now }, Resume: func(context.Context, string, string) error {
		woke = true
		return nil
	}})
	if woke {
		t.Fatal("wake re-resumed an already-acked intent — replyKey byte-compat broken")
	}
	if wr.Status != "noop" {
		t.Fatalf("wake status = %q, want noop (nothing pending)", wr.Status)
	}
}
