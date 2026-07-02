package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

const shallowBootWindowFixtureStream = `{"type":"system","subtype":"init","model":"claude-test-model"}
{"type":"assistant","message":{"id":"msg_1","model":"claude-test-model","usage":{"input_tokens":10,"output_tokens":5,"cache_creation_input_tokens":100,"cache_read_input_tokens":50},"content":[{"type":"tool_use","id":"toolu_1","name":"Bash","input":{"command":"echo hi"}}]}}
{"type":"assistant","message":{"id":"msg_2","model":"claude-test-model","usage":{"input_tokens":20,"output_tokens":8,"cache_creation_input_tokens":200,"cache_read_input_tokens":90},"content":[{"type":"text","text":"Hello!"}]}}
`

// TestShallowBootWindowRecordCommandExtractsFromArchivedStream is AC-4's proof
// that the CLI applies AC-1's extraction logic (BuildShallowBootWindowRecord) to
// a standalone archived claude-stream.jsonl — the mechanism the release-ledger
// backfill uses instead of duplicating extraction logic. The fixture's greet
// turn is turns[1] (the only turn with no dispatch tool_use), so Turns == 2 and
// Tokens carries turns[1]'s full usage.
func TestShallowBootWindowRecordCommandExtractsFromArchivedStream(t *testing.T) {
	streamPath := filepath.Join(t.TempDir(), "claude-stream.jsonl")
	if err := os.WriteFile(streamPath, []byte(shallowBootWindowFixtureStream), 0o644); err != nil {
		t.Fatal(err)
	}
	outDir := t.TempDir()

	if code := shallowBootWindowRecord([]string{streamPath, "--model", "claude-test-model", "--out", outDir}); code != 0 {
		t.Fatalf("shallowBootWindowRecord exit = %d, want 0", code)
	}

	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected exactly one emitted record file, got %d: %+v", len(entries), entries)
	}
	data, err := os.ReadFile(filepath.Join(outDir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	var record journeymetrics.Record
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	if record.ScenarioID != "shallow-boot-window" {
		t.Fatalf("ScenarioID = %q, want shallow-boot-window", record.ScenarioID)
	}
	if record.Turns != 2 {
		t.Fatalf("Turns = %d, want 2", record.Turns)
	}
	want := journeymetrics.TokenTotals{Input: 20, Output: 8, CacheCreation: 200, CacheRead: 90, Total: 318}
	if record.Tokens != want {
		t.Fatalf("Tokens = %+v, want %+v", record.Tokens, want)
	}
}

// TestShallowBootWindowRecordCommandRejectsMissingArgs proves a miswired
// invocation (missing the required stream/--model/--out) is a usage error, not a
// silent no-op that would leave the AC-4 backfill believing extraction ran.
func TestShallowBootWindowRecordCommandRejectsMissingArgs(t *testing.T) {
	if code := shallowBootWindowRecord(nil); code == 0 {
		t.Fatalf("shallowBootWindowRecord exit = 0 with no arguments; want non-zero")
	}
}

// TestShallowBootWindowRecordCommandRejectsMissingStreamFile proves an
// unreadable stream path fails loud instead of emitting an empty/zero record.
func TestShallowBootWindowRecordCommandRejectsMissingStreamFile(t *testing.T) {
	code := shallowBootWindowRecord([]string{filepath.Join(t.TempDir(), "does-not-exist.jsonl"), "--model", "claude-test-model", "--out", t.TempDir()})
	if code == 0 {
		t.Fatalf("shallowBootWindowRecord exit = 0 with a missing stream file; want non-zero")
	}
}

// TestShallowBootWindowRecordCommandRejectsStreamWithNoAssistantTurns proves a
// stream that produced no assistant turns at all (nothing to extract a
// boot-window observation from) is rejected rather than silently skipped.
func TestShallowBootWindowRecordCommandRejectsStreamWithNoAssistantTurns(t *testing.T) {
	streamPath := filepath.Join(t.TempDir(), "claude-stream.jsonl")
	if err := os.WriteFile(streamPath, []byte(`{"type":"system","subtype":"init","model":"claude-test-model"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code := shallowBootWindowRecord([]string{streamPath, "--model", "claude-test-model", "--out", t.TempDir()})
	if code == 0 {
		t.Fatalf("shallowBootWindowRecord exit = 0 on a stream with no assistant turns; want non-zero")
	}
}
