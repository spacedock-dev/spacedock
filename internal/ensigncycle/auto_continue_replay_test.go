// ABOUTME: Offline replay of a real captured auto-continue claude stream through the
// ABOUTME: dispatch-evidence lifecycle check — the non-regression control for AC-2.
package ensigncycle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
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

// TestAutoContinueGateStateVocabulary asserts every value gates.attemptState can
// produce, because "not open" and "bypassed" are different claims and conflating
// them accuses correct behavior.
//
// The vocabulary is closed and enumerated from source: attemptState has exactly four
// returns — open, withdrawn, closed, invalid — and CurrentSummary yields "" when a
// record has no current attempt. Summary.State is assigned from attemptState alone
// (model.go), so no caller widens the range. Each of the five is asserted below.
//
// The case that forced this: claude run 2 of the AC-4 loop red under
// `human-gate-bypassed` on state `withdrawn`, against an FO that had prepared the
// gate against the stale base copy, caught the mistake before presenting, withdrawn
// that room and re-prepared correctly against the worktree. Withdrawal records no
// decision. Only a Resolution accuses.
func TestAutoContinueGateStateVocabulary(t *testing.T) {
	stream := readFile(t, filepath.Join("testdata", autoContinueReplayFixture))
	const worktree = ".worktrees/spacedock-ensign-auto-continue-task"

	body := func(frontmatter string) string {
		e := strings.Replace(autoContinueEntity(), "status: implementation", "status: validation", 1)
		e = strings.Replace(e, "worktree:\n---\n", "worktree: "+worktree+"\n"+frontmatter+"---\n", 1)
		return e + "\n## Stage Report: validation\n\n- DONE: Verify the implementation against AC-1\n  PASSED.\n"
	}
	noRecord := func() string {
		e := strings.Replace(autoContinueEntity(), "status: implementation", "status: validation", 1)
		e = strings.Replace(e, "worktree:\n", "worktree: "+worktree+"\n", 1)
		return e + "\n## Stage Report: validation\n\n- DONE: Verify the implementation against AC-1\n  PASSED.\n"
	}

	for _, c := range []struct {
		name         string
		base, wtCopy string
		wantRed      bool
		wantBypass   bool
	}{
		// open: the conforming end state.
		{name: "open", base: noRecord(), wtCopy: body(autoContinueGateFrontmatter(false))},
		// closed: the bypass this entity exists to catch, and it must not be
		// maskable by an open attempt in the other copy.
		{name: "closed_beside_open", base: body(autoContinueGateFrontmatter(true)), wtCopy: body(autoContinueGateFrontmatter(false)),
			wantRed: true, wantBypass: true},
		// withdrawn: sanctioned retraction. The live shape that forced this test.
		{name: "withdrawn_beside_open", base: body(autoContinueWithdrawnGateFrontmatter()), wtCopy: body(autoContinueGateFrontmatter(false))},
		// withdrawn everywhere: no gate was left for the captain, but nobody
		// resolved anything, so it must not be graded as a bypass.
		{name: "withdrawn_only", base: body(autoContinueWithdrawnGateFrontmatter()), wtCopy: body(autoContinueWithdrawnGateFrontmatter()),
			wantRed: true},
		// invalid: carries a Resolution, so an open copy must NOT rescue it, but the
		// record is too incoherent to accuse with.
		{name: "invalid_beside_open", base: body(autoContinueInvalidGateFrontmatter()), wtCopy: body(autoContinueGateFrontmatter(false)),
			wantRed: true},
		// no record in any copy: the gate state is unknowable.
		{name: "no_record_anywhere", base: noRecord(), wtCopy: noRecord(), wantRed: true},
	} {
		t.Run(c.name, func(t *testing.T) {
			root := t.TempDir()
			entityPath, err := writeAutoContinueWorkflowNoGit(root)
			if err != nil {
				t.Fatal(err)
			}
			writeFile(t, entityPath, c.base)
			writeFile(t, filepath.Join(root, worktree, filepath.Base(entityPath)), c.wtCopy)
			gitInit(t, root)

			err = assertAutoContinueDispatchEvidence(t, stream, root, entityPath)
			switch {
			case !c.wantRed && err != nil:
				t.Fatalf("conforming shape graded RED: %v", err)
			case c.wantRed && err == nil:
				t.Fatal("expected RED, graded GREEN")
			case c.wantRed && c.wantBypass && gradedCode(err) != autoContinueBypassCode:
				t.Fatalf("graded under %q, want %q", gradedCode(err), autoContinueBypassCode)
			case c.wantRed && !c.wantBypass && gradedCode(err) == autoContinueBypassCode:
				t.Fatalf("state %q must not be accused of a bypass — no Resolution was recorded against an unapproved gate: %v", c.name, err)
			}
		})
	}
}

// TestAutoContinueGateFixturesParseAsIntended pins that each frontmatter helper
// really produces the state it is named for. Without this the vocabulary table can
// pass vacuously: a fixture that fails to decode is skipped by
// autoContinueGateStates, so its case asserts nothing while appearing to pass. That
// happened, and it was caught only because a falsification failed to fail.
func TestAutoContinueGateFixturesParseAsIntended(t *testing.T) {
	const rejectedByDecoder = "<rejected by decoder>"

	for _, c := range []struct{ name, frontmatter, want string }{
		{"open", autoContinueGateFrontmatter(false), "open"},
		{"closed", autoContinueGateFrontmatter(true), "closed"},
		{"withdrawn", autoContinueWithdrawnGateFrontmatter(), "withdrawn"},
		// `invalid` is unreachable as a STATE: the decoder rejects a conflicting
		// withdrawal+resolution before CurrentSummary ever sees it. Pinned as a
		// rejection so a future loosening of the decoder shows up here.
		{"invalid", autoContinueInvalidGateFrontmatter(), rejectedByDecoder},
	} {
		t.Run(c.name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, "entity.md")
			writeFile(t, path, "---\nid: auto-continue-task\ntitle: T\nstatus: validation\ncompleted:\nverdict:\nworktree:\n"+c.frontmatter+"---\n# T\n")
			doc, _, err := gates.Read(path)
			if c.want == rejectedByDecoder {
				if err == nil {
					t.Fatalf("fixture %q now decodes as state %q — the decoder no longer rejects a conflicting withdrawal+resolution, so `invalid` can reach the gate check and needs a branch there", c.name, gates.CurrentSummary(doc, "validation").State)
				}
				return
			}
			if err != nil {
				t.Fatalf("fixture %q does not decode, so any test using it asserts nothing: %v", c.name, err)
			}
			if got := gates.CurrentSummary(doc, "validation").State; got != c.want {
				t.Fatalf("fixture %q parses as state %q, want %q", c.name, got, c.want)
			}
		})
	}
}

// autoContinueRevalidateFixture is a REAL claude-sonnet-5 stream from the split-root
// leg of run 4 of this entity's AC-4 loop — the run whose FO dispatched three agents
// and was counted as zero.
//
// It is the line-type-filtered form of the capture (rows carrying a `tool_use` block
// or a `task_notification`), the fallback the entity pre-approved for an oversized
// stream. Full capture: 975079 bytes / 718 lines, sha256
// 2895ae231957a06a772e667571a5f362779245e804b0ac3333bc8b0ffc2ebd83. Filtered: 149187
// bytes / 122 lines, sha256
// 870418e9dca8b163eda734e840079037b1999dd507a10debc0453e85f67d02f2. Capture-time
// identical-verdict proof: assertWorkerLifecycle returns the same verdict on both for
// stage `validation` AND stage `implementation`; only the line indices in the message
// shift, which filtering necessarily changes.
const autoContinueRevalidateFixture = "claude_live_auto_continue_run4_revalidate_cycle_sonnet.stream.jsonl"

// TestClaudeSpawnIsForStageAnchorsOnTheDispatchPointer pins the matcher against the
// three real Agent blocks that defeated description-only matching. Their descriptions
// were composed by a live FO; their dispatch pointers were generated by
// `spacedock dispatch build`.
//
// The reviser is the case that matters most: it must NOT count as a validation spawn,
// and its pointer is what says so. Without the pointer the matcher counts zero
// validators; with a pointer that ignored the stage suffix it would count three.
func TestClaudeSpawnIsForStageAnchorsOnTheDispatchPointer(t *testing.T) {
	const dispatch = "/tmp/spacedock-dispatch/spacedock-ensign-auto-continue-task-"
	for _, c := range []struct {
		name, description, prompt string
		wantValidation            bool
	}{
		{"validator", "Validate auto-continue-task implementation", dispatch + "validation.md", true},
		{"reviser", "Revise auto-continue-task implementation", dispatch + "implementation.md", false},
		{"revalidator", "Re-validate auto-continue-task cycle 2", dispatch + "validation.md", true},
		{"description_alone_still_counts", "Validation worker for auto-continue-task", "", true},
		{"neither_signal", "Tidy the workspace", dispatch + "backlog.md", false},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := claudeSpawnIsForStage(c.description, c.prompt, "validation"); got != c.wantValidation {
				t.Fatalf("claudeSpawnIsForStage(%q, ...) = %v, want %v", c.description, got, c.wantValidation)
			}
		})
	}
}

// TestAutoContinueRevalidateStreamCountsBothValidators replays the real run-4 stream
// and pins what the matcher now sees: two validation dispatches, not zero. The
// remaining red on that stream is the reject-then-revalidate cycle tripping
// `spawns != 1`, which belongs to the shared helper's ordering semantics and is not
// this entity's to change.
//
// No-regression evidence for the pointer anchor is the 14-leg sweep over this
// entity's AC-4 loop recorded in the stage report: every leg that passed under
// description-only matching counted exactly 1 under both signals, so the anchor adds
// only spawns the dispatch contract itself labels. The committed run-31915540750
// replay above is the standing guard for that direction.
func TestAutoContinueRevalidateStreamCountsBothValidators(t *testing.T) {
	stream := readFile(t, filepath.Join("testdata", autoContinueRevalidateFixture))

	got := countClaudeStageSpawns(t, stream, "validation")
	if got != 2 {
		t.Fatalf("validation spawns in the real re-validate stream = %d, want 2 (the validator and the re-validator; the reviser must not count)", got)
	}
	if descOnly := countClaudeStageSpawnsDescriptionOnly(t, stream, "validation"); descOnly != 0 {
		t.Fatalf("description-only matching = %d, want 0 — the fixture no longer reproduces the false negative it exists to pin", descOnly)
	}
}

// countClaudeStageSpawns and countClaudeStageSpawnsDescriptionOnly walk the stream's
// Agent blocks the way assertWorkerLifecycle does, so the test observes the matcher
// rather than restating it.
func countClaudeStageSpawns(t *testing.T, stream, stage string) int {
	t.Helper()
	return countAgentBlocks(t, stream, func(description, prompt string) bool {
		return claudeSpawnIsForStage(description, prompt, stage)
	})
}

func countClaudeStageSpawnsDescriptionOnly(t *testing.T, stream, stage string) int {
	t.Helper()
	return countAgentBlocks(t, stream, func(description, _ string) bool {
		return strings.Contains(strings.ToLower(description), stage)
	})
}

func countAgentBlocks(t *testing.T, stream string, match func(description, prompt string) bool) int {
	t.Helper()
	type block struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		Input struct{ Description, Prompt string }
	}
	count := 0
	for _, line := range strings.Split(stream, "\n") {
		var event struct {
			Message *struct{ Content []block }
		}
		if json.Unmarshal([]byte(line), &event) != nil || event.Message == nil {
			continue
		}
		for _, item := range event.Message.Content {
			if item.Type == "tool_use" && item.Name == "Agent" && match(item.Input.Description, item.Input.Prompt) {
				count++
			}
		}
	}
	return count
}
