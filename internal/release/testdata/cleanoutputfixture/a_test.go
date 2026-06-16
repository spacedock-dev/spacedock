// ABOUTME: Fixture test package proving the CI clean-output + jsonl-archive
// ABOUTME: shape; failures are gated by FIXTURE_FAIL so one binary runs green or red.
package cleanoutputfixture

import (
	"os"
	"testing"
)

// failOn reports whether the planted failures should fire. The AC test runs the
// same compiled binary twice — once with FIXTURE_FAIL unset (green) and once set
// (red) — so a single fixture exercises both the clean-green and failures-only-red
// surfaces without a second fixture package.
func failOn() bool { return os.Getenv("FIXTURE_FAIL") == "1" }

func TestAlphaPasses(t *testing.T) {
	if 1+1 != 2 {
		t.Fatal("arithmetic is broken")
	}
}

func TestGammaFails(t *testing.T) {
	if !failOn() {
		return
	}
	got := 7
	if got != 42 {
		t.Errorf("compute() = %d, want 42", got)
	}
}

func TestZetaSubtests(t *testing.T) {
	t.Run("case_good", func(t *testing.T) {})
	t.Run("case_bad", func(t *testing.T) {
		if !failOn() {
			return
		}
		t.Errorf("subtest assertion: got %q want %q", "x", "y")
	})
}

// streamFirehoseMarker is a distinctive token planted in every t.Log line so the
// AC test can assert the firehose is ABSENT from the clean stdout (and PRESENT in
// the jsonl archive).
const streamFirehoseMarker = "STREAM_FIREHOSE_LINE"

// TestPassingWithStreamFirehose mimics the live runner's stream watcher, which is
// `func(line string) { t.Log(line) }` — it t.Logs every host stream-json line,
// even on a PASSING test. Under -test.v those t.Log lines are arbitrary text with
// no go-test structural prefix, so a `grep -vE` over the scaffolding canNOT strip
// them; only a tool that routes per-test output off stdout (gotestsum) suppresses
// this firehose on a green run. This test ALWAYS passes — the firehose is the
// point, not a failure — so AC-1's green-run clean-stdout assertion exercises the
// real 143KB-class firehose, not just the go-test scaffolding.
func TestPassingWithStreamFirehose(t *testing.T) {
	for i := 0; i < 40; i++ {
		t.Logf(`%s {"type":"assistant","message":{"content":[{"type":"text","text":"lorem ipsum dolor sit amet line %d"}]}}`, streamFirehoseMarker, i)
	}
}
