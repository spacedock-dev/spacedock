// ABOUTME: Offline replay of two real captured zero-discover boot streams through
// ABOUTME: detectBroadSearchAtBoot — proves the false-red/genuine-red boundary on real data.
package ensigncycle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDetectBroadSearchAtBootRealZeroDiscoverStreams replays two real captured
// sonnet zero-discover boot streams through detectBroadSearchAtBoot, the #462
// wrong-root-detector pattern (TestDetectWrongRootBootRealPR446Streams) applied to
// this detector.
//
//   - PR #467 (run 28587056405, artifact runtime-live-e2e-claude-live-sonnet id
//     8037986947): a CORRECT report-and-stop boot that ran `ls -la {root}` before
//     reporting no workflow and stopping (stream `result` event `subtype=success`).
//     The captain's 2026-07-02 decision made flat `ls` legal after zero-discover, so
//     this stream must return nil — the false red this task removes.
//   - PR #398 (run 27835552853 attempt 1, artifact runtime-live-e2e-claude-live-sonnet
//     id 7755173878): a GENUINE deviation — `find {root} -type f | head -30` — that
//     stays banned. This is the non-regression guard proving the alignment did not
//     loosen genuine detection (tq0 owns hardening that class further; this task must
//     not touch it).
//
// Both streams were reconstructed from CI artifacts by extracting the
// `zero_discover_live_test.go:82:`-prefixed t.Log lines from the run's
// live-e2e-detail.jsonl (ideation spike, 2026-07-02); their sha256 is recorded in
// the entity's Spike record for provenance.
func TestDetectBroadSearchAtBootRealZeroDiscoverStreams(t *testing.T) {
	cases := []struct {
		name        string
		fixture     string
		fixtureRoot string
		wantRed     bool
		wantNames   string
	}{
		{
			name:        "pr467_flat_ls_correct_boot_passes",
			fixture:     "claude_live_zero_discover_pr467_sonnet.stream.jsonl",
			fixtureRoot: "/tmp/TestLiveZeroDiscoverReportsAndStops2027276970/002",
			wantRed:     false,
		},
		{
			name:        "pr398_find_genuine_deviation_reds",
			fixture:     "claude_live_zero_discover_pr398_sonnet.stream.jsonl",
			fixtureRoot: "/tmp/TestLiveZeroDiscoverReportsAndStops2170189205/002",
			wantRed:     true,
			wantNames:   "find",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			stream, err := os.ReadFile(filepath.Join("testdata", c.fixture))
			if err != nil {
				t.Fatalf("open real-stream fixture %q: %v", c.fixture, err)
			}
			red := detectBroadSearchAtBoot(string(stream), c.fixtureRoot)
			if c.wantRed {
				if red == nil {
					t.Fatalf("detector passed the real %s stream — want a red naming %q", c.name, c.wantNames)
				}
				if c.wantNames != "" && !strings.Contains(red.Error(), c.wantNames) {
					t.Errorf("error must name the offending search %q: %v", c.wantNames, red)
				}
			} else if red != nil {
				t.Errorf("detector red the real %s stream (correct report-and-stop boot) — want nil: %v", c.name, red)
			}
		})
	}
}
