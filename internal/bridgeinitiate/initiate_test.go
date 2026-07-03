package bridgeinitiate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAppendInitiationWritesRecord(t *testing.T) {
	root := t.TempDir()
	got, err := AppendInitiation(InitiationOptions{
		Root:      root,
		Now:       func() time.Time { return time.Date(2026, 7, 3, 5, 0, 0, 0, time.UTC) },
		ID:        "gate-ship-a-ideation",
		Kind:      "gate-review",
		Workflow:  "pr-review-queue",
		Entity:    "ship-a",
		ShipID:    "pr-review-queue/ship-a",
		Host:      "claude",
		SessionID: "s1",
		Headline:  "ship-a ideation ready for a call",
		Body:      "chosen direction: reuse existing writer",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "gate-ship-a-ideation" || !got.Queued {
		t.Fatalf("result = %+v", got)
	}
	if got.RequestID != "gate-ship-a-ideation" {
		t.Fatalf("request id = %q, want default to id for gate-review", got.RequestID)
	}

	data, err := os.ReadFile(filepath.Join(root, "_bridge", "fo-initiate.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var rec InitiationRecord
	if err := json.Unmarshal(bytes.TrimSpace(data), &rec); err != nil {
		t.Fatalf("record JSON: %v\n%s", err, data)
	}
	if rec.Schema != 1 || rec.Status != "open" {
		t.Fatalf("record = %+v", rec)
	}
	if rec.Kind != "gate-review" || rec.RequestID != "gate-ship-a-ideation" {
		t.Fatalf("record = %+v", rec)
	}
	if rec.TS != "2026-07-03T05:00:00Z" || rec.Headline != "ship-a ideation ready for a call" {
		t.Fatalf("record fields = %+v", rec)
	}
	if rec.ShipID != "pr-review-queue/ship-a" || rec.Workflow != "pr-review-queue" {
		t.Fatalf("record routing = %+v", rec)
	}
}

func TestAppendInitiationDefaultsRequestIDOnlyForGate(t *testing.T) {
	root := t.TempDir()
	got, err := AppendInitiation(InitiationOptions{
		Root:     root,
		ID:       "status-1",
		Kind:     "status",
		Headline: "advancing ship-b",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !got.Queued {
		t.Fatalf("result = %+v", got)
	}
	if got.RequestID != "" {
		t.Fatalf("request id = %q, want empty for non-gate kind", got.RequestID)
	}
}

func TestAppendInitiationMissingIDIsLoudError(t *testing.T) {
	got, err := AppendInitiation(InitiationOptions{
		Root:     t.TempDir(),
		Kind:     "status",
		Headline: "no id here",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Queued || got.Error == "" {
		t.Fatalf("result = %+v, want non-queued error result", got)
	}
}

func TestAppendInitiationInvalidKindIsLoudError(t *testing.T) {
	got, err := AppendInitiation(InitiationOptions{
		Root:     t.TempDir(),
		ID:       "x1",
		Kind:     "chatter",
		Headline: "bad kind",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Queued || got.Error == "" {
		t.Fatalf("result = %+v, want non-queued error result", got)
	}
}

func TestAppendInitiationMissingHeadlineIsLoudError(t *testing.T) {
	got, err := AppendInitiation(InitiationOptions{
		Root:     t.TempDir(),
		ID:       "x1",
		Kind:     "status",
		Headline: "   ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Queued || got.Error == "" {
		t.Fatalf("result = %+v, want non-queued error result", got)
	}
}

func TestAppendInitiationNormalizesBounds(t *testing.T) {
	root := t.TempDir()
	got, err := AppendInitiation(InitiationOptions{
		Root:     root,
		ID:       "x1",
		Kind:     "reco",
		Headline: "line one\ntwo\tthree",
		Body:     "body\n" + strings.Repeat("z", maxBodyLen+500),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !got.Queued {
		t.Fatalf("result = %+v", got)
	}
	data, err := os.ReadFile(filepath.Join(root, "_bridge", "fo-initiate.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var rec InitiationRecord
	if err := json.Unmarshal(bytes.TrimSpace(data), &rec); err != nil {
		t.Fatal(err)
	}
	if rec.Headline != "line one two three" {
		t.Fatalf("headline = %q", rec.Headline)
	}
	if strings.ContainsAny(rec.Body, "\n\t") || len(rec.Body) > maxBodyLen {
		t.Fatalf("body len=%d %q", len(rec.Body), rec.Body)
	}
}

func TestAppendInitiation_AnchorsAtRepoRoot(t *testing.T) {
	root := t.TempDir()
	if _, err := AppendInitiation(InitiationOptions{
		Root:     root,
		ID:       "anchor-1",
		Kind:     "status",
		Headline: "anchoring check",
	}); err != nil {
		t.Fatal(err)
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}
	// The file must land at filepath.Abs(root)/_bridge/fo-initiate.jsonl — the
	// same anchor bridgealert.AppendPermission uses — NOT a canonicalBridgeRoot
	// git-root walk-up.
	want := filepath.Join(absRoot, "_bridge", "fo-initiate.jsonl")
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("expected file at %s: %v", want, err)
	}
}

func TestTruncateInitiate_KeepsOpenGate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fo-initiate.jsonl")

	var lines []string
	// An OLD still-open gate-review that must survive truncation.
	openGate := InitiationRecord{
		Schema: 1, ID: "gate-old", TS: "2026-07-03T00:00:00Z",
		Kind: "gate-review", Headline: "old open gate", RequestID: "gate-old", Status: "open",
	}
	lines = append(lines, mustJSON(t, openGate))
	// A resolved old gate that is NOT protected (should be dropped).
	resolvedGate := InitiationRecord{
		Schema: 1, ID: "gate-done", TS: "2026-07-03T00:00:01Z",
		Kind: "gate-review", Headline: "resolved gate", RequestID: "gate-done", Status: "resolved",
	}
	lines = append(lines, mustJSON(t, resolvedGate))
	// Enough status records to exceed maxLines and push the gate before the tail.
	for i := 0; i < defaultMaxLines+200; i++ {
		lines = append(lines, mustJSON(t, InitiationRecord{
			Schema: 1, ID: fmt.Sprintf("status-%d", i), TS: "2026-07-03T01:00:00Z",
			Kind: "status", Headline: fmt.Sprintf("noise %d", i), Status: "open",
		}))
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	truncateInitiate(path)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	kept := strings.Split(strings.TrimSpace(string(data)), "\n")
	// keepLines tail + exactly the one protected open gate.
	if len(kept) != defaultKeepLines+1 {
		t.Fatalf("kept %d lines, want %d", len(kept), defaultKeepLines+1)
	}

	var foundOpenGate, foundResolvedGate bool
	for _, line := range kept {
		rec, ok := parseRecord(line)
		if !ok {
			t.Fatalf("unparseable kept line: %s", line)
		}
		if rec.ID == "gate-old" {
			foundOpenGate = true
		}
		if rec.ID == "gate-done" {
			foundResolvedGate = true
		}
	}
	if !foundOpenGate {
		t.Fatal("still-open gate-review was evicted by truncation")
	}
	if foundResolvedGate {
		t.Fatal("resolved gate-review should not be retained past the cap")
	}
}

func TestTruncateInitiate_EvictsGateResolvedViaInbox(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fo-initiate.jsonl")

	var lines []string
	// Both gates are written status "open" — the production writer NEVER stamps
	// "resolved" on disk. Resolution is signalled only by an inbox decision intent.
	undecidedGate := InitiationRecord{
		Schema: 1, ID: "gate-open", TS: "2026-07-03T00:00:00Z",
		Kind: "gate-review", Headline: "undecided gate", RequestID: "gate-open", Status: "open",
	}
	decidedGate := InitiationRecord{
		Schema: 1, ID: "gate-decided", TS: "2026-07-03T00:00:01Z",
		Kind: "gate-review", Headline: "captain decided this", RequestID: "gate-decided", Status: "open",
	}
	lines = append(lines, mustJSON(t, undecidedGate), mustJSON(t, decidedGate))
	for i := 0; i < defaultMaxLines+200; i++ {
		lines = append(lines, mustJSON(t, InitiationRecord{
			Schema: 1, ID: fmt.Sprintf("status-%d", i), TS: "2026-07-03T01:00:00Z",
			Kind: "status", Headline: fmt.Sprintf("noise %d", i), Status: "open",
		}))
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// A captain decision intent for gate-decided lands in the sibling inbox.
	inbox := `{"schema":1,"kind":"decision","request_id":"gate-decided","verdict":"approve"}` + "\n"
	if err := os.WriteFile(filepath.Join(dir, "inbox.jsonl"), []byte(inbox), 0o644); err != nil {
		t.Fatal(err)
	}

	truncateInitiate(path)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var foundUndecided, foundDecided bool
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		rec, ok := parseRecord(line)
		if !ok {
			t.Fatalf("unparseable kept line: %s", line)
		}
		switch rec.ID {
		case "gate-open":
			foundUndecided = true
		case "gate-decided":
			foundDecided = true
		}
	}
	if !foundUndecided {
		t.Fatal("undecided open gate-review was evicted by truncation")
	}
	if foundDecided {
		t.Fatal("gate-review resolved by an inbox decision should be evicted past the cap")
	}
}

func mustJSON(t *testing.T, rec InitiationRecord) string {
	t.Helper()
	data, err := json.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
