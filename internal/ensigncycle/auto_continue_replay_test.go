// ABOUTME: Offline replay of a real captured auto-continue claude stream through the
// ABOUTME: dispatch-evidence lifecycle check — the non-regression control for AC-2.
package ensigncycle

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

// autoContinueReplayFixture is a REAL claude-sonnet-5 auto-continue stream, captured
// from the single-root leg of GitHub Actions run 31915540750 on main (artifact
// runtime-live-e2e-claude-live-claude-sonnet-5, id 9255098892) — a green run whose FO
// advanced to validation, dispatched a fresh validator, and presented the gate.
//
// Provenance is recorded in the entity's Spike section at ideation and re-verified
// byte-for-byte against the artifact at capture time.
const (
	autoContinueReplayFixture = "claude_live_auto_continue_run31915540750_sonnet.stream.jsonl"
	autoContinueReplayDigest  = "465e849afda41a39c8443bae98119e8ab806a0af0a116f54f5321761d7eab287"
	autoContinueReplayBytes   = 430087
)

// TestAutoContinueReplayRealClaudeStream is AC-2: the tightened grading produces no false
// red on real observed conforming behavior. It follows zero_discover_replay_test.go —
// replay a captured stream from a green CI run through the check that grades it.
//
// This is the first regression guard over the claude branch of assertWorkerLifecycle for
// this journey: the check was previously reached only through an optional interface that
// only codex implemented, so the claude dialect's parse had never run in CI at all. A
// synthetic stream was rejected deliberately — it would prove the parser matches what the
// author imagined the runtime emits, not what it emits.
//
// What would make this fail: an ordering rule the real stream violates, a spawn/completion
// shape the parser stops recognizing, or a claude dialect change. Each is a true signal
// that the live grader has gone blind — the failure this entity exists to prevent.
func TestAutoContinueReplayRealClaudeStream(t *testing.T) {
	stream := readFile(t, filepath.Join("testdata", autoContinueReplayFixture))
	report := autoContinueGatedEndState(false)

	if err := assertWorkerLifecycle(stream, report, "validation", "gate prepare"); err != nil {
		t.Fatalf("the real captured claude stream from green run 31915540750 graded RED: %v", err)
	}
}

// TestAutoContinueReplayRealStreamThroughDispatchEvidence replays the same real stream
// through the whole dispatch-evidence check — lifecycle parse, durable-commit lookup, and
// gate-open state — against state staged the way a conforming run leaves it: the
// end-to-end half of AC-2, not just the stream parse in isolation. Both fixture layouts
// are replayed because each host runs both and the report lookup differs.
func TestAutoContinueReplayRealStreamThroughDispatchEvidence(t *testing.T) {
	stream := readFile(t, filepath.Join("testdata", autoContinueReplayFixture))

	for _, variant := range []struct {
		name      string
		splitRoot bool
	}{
		{name: "single-root", splitRoot: false},
		{name: "split-root", splitRoot: true},
	} {
		t.Run(variant.name, func(t *testing.T) {
			stateRoot, entityPath := stageAutoContinueEndState(t, variant.splitRoot, false)
			if err := assertAutoContinueDispatchEvidence(t, stream, stateRoot, entityPath); err != nil {
				t.Fatalf("real captured leg graded RED on the %s layout: %v", variant.name, err)
			}
		})
	}
}

// TestAutoContinueReplayFixtureIsTheCapturedArtifact pins the fixture's identity to the
// digest recorded in the entity at ideation. The replay control is only evidence if the
// bytes are the ones CI actually produced; a regenerated or hand-edited stream would
// silently turn AC-2 back into the synthetic check it was chosen over. A digest mismatch
// means the fixture was rewritten, not that the runtime changed.
func TestAutoContinueReplayFixtureIsTheCapturedArtifact(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", autoContinueReplayFixture))
	if err != nil {
		t.Fatal(err)
	}
	if got := hex.EncodeToString(sha256Sum(raw)); got != autoContinueReplayDigest {
		t.Fatalf("replay fixture sha256 = %s, want the captured %s — the artifact was replaced", got, autoContinueReplayDigest)
	}
	if len(raw) != autoContinueReplayBytes {
		t.Fatalf("replay fixture is %d bytes, want the captured %d", len(raw), autoContinueReplayBytes)
	}
}

func sha256Sum(b []byte) []byte {
	sum := sha256.Sum256(b)
	return sum[:]
}
