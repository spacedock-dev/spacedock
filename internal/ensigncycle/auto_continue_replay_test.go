// ABOUTME: Offline replay of a real captured auto-continue claude stream through the
// ABOUTME: dispatch-evidence lifecycle check — the non-regression control for AC-2.
package ensigncycle

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
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

// TestAutoContinueGateStateAcrossEntityCopyPlacements pins every placement of the
// gate record observed live, because the FO does not reliably keep the gate record
// and the validation report in the same copy of a worktree-backed entity:
//
//   - base-only: no worktree in play; both live in the base entity.
//   - worktree-both: claude run 2 of the AC-4 loop put the report AND the gate record
//     in the worktree copy, so reading gates from the base path red a conforming run.
//   - split: codex files the report in the worktree copy but the gate record in the
//     base copy, so reading gates from the worktree path red a conforming run.
//
// Reading either path alone therefore reds a real host's conforming run — the AC-2
// false-red class. Each placement is asserted in BOTH directions so tolerating the
// placement cannot blind the bypass check.
//
// The first attempt at this guard tested only worktree-both, the placement its author
// had hypothesised, and so confirmed the expectation instead of varying the dimension
// that actually differs. The authorized single-path remedy shipped on that proof and
// red codex live. The table below varies placement deliberately.
func TestAutoContinueGateStateAcrossEntityCopyPlacements(t *testing.T) {
	stream := readFile(t, filepath.Join("testdata", autoContinueReplayFixture))
	const worktree = ".worktrees/spacedock-ensign-auto-continue-task"

	stage := func(t *testing.T, placement string, resolved bool) (string, string) {
		t.Helper()
		root := t.TempDir()
		entityPath, err := writeAutoContinueWorkflowNoGit(root)
		if err != nil {
			t.Fatal(err)
		}
		gated := autoContinueGatedEndState(resolved)
		ungated := strings.Replace(autoContinueEntity(), "status: implementation", "status: validation", 1) +
			"\n## Stage Report: validation\n\n- DONE: Verify the implementation against AC-1\n  PASSED.\n"
		worktreeCopy := filepath.Join(root, worktree, filepath.Base(entityPath))

		switch placement {
		case "base-only":
			writeFile(t, entityPath, gated)
		case "worktree-both":
			writeFile(t, entityPath, strings.Replace(ungated, "worktree:\n", "worktree: "+worktree+"\n", 1))
			writeFile(t, worktreeCopy, gated)
		case "split-report-worktree-gate-base":
			writeFile(t, entityPath, strings.Replace(gated, "worktree:\n", "worktree: "+worktree+"\n", 1))
			writeFile(t, worktreeCopy, ungated)
		default:
			t.Fatalf("unknown placement %q", placement)
		}
		gitInit(t, root)
		return root, entityPath
	}

	for _, placement := range []string{"base-only", "worktree-both", "split-report-worktree-gate-base"} {
		t.Run(placement, func(t *testing.T) {
			t.Run("open_gate_grades_green", func(t *testing.T) {
				root, entityPath := stage(t, placement, false)
				if err := assertAutoContinueDispatchEvidence(t, stream, root, entityPath); err != nil {
					t.Fatalf("conforming run graded RED with the gate record placed %s: %v", placement, err)
				}
			})
			t.Run("resolved_gate_still_reds", func(t *testing.T) {
				root, entityPath := stage(t, placement, true)
				err := assertAutoContinueDispatchEvidence(t, stream, root, entityPath)
				if err == nil {
					t.Fatalf("a resolved gate placed %s graded GREEN — tolerating placement must not blind the bypass check", placement)
				}
				if code := gradedCode(err); code != autoContinueBypassCode {
					t.Fatalf("bypass placed %s graded under %q, want %q", placement, code, autoContinueBypassCode)
				}
			})
		})
	}
}
