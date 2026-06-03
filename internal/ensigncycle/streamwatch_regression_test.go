// ABOUTME: Offline regression replay of the archived sonnet stream that hung the
// ABOUTME: live cycle at expectExit (TeamDelete-fails-then-no-retry); pins the watcher's localization.
package ensigncycle

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestSonnetTeamDeleteHangReplay drips the captured sonnet stream-json from the
// failing live-e2e CI run (PR #277, run 26865865357, artifact
// runtime-live-e2e-claude-live-sonnet) through the SAME watcher the live test
// drives, and asserts the watcher localizes the FO terminal-teardown hang at
// expectExit. The recording is the real bug: the FO emits a premature
// shutdown_request right after dispatch, the dispatch-close fires on the keyed
// task_notification, then terminal TeamDelete fails with "active member(s)" and
// the FO ends its turn WITHOUT retrying TeamDelete to success, so the claude -p
// subprocess never exits and expectExit kills it.
//
// This is AC-2: a real captured stream + the real watcher, run offline with no
// model spend. It pins the watcher's localization (TeamCreate matches,
// dispatch-close fires, expectExit trips a stepTimeout that kills the
// never-exiting proc) and fails red if that localization ever regresses.
//
// HONEST LIMIT (recorded in the entity body): this fixture does NOT flip
// red->green on the FO contract fix. The fix changes a FUTURE live FO stream,
// not this recording — a recorded stream cannot demonstrate a behavior change in
// the producer. The fix's red->green proof is the live run (AC-1). What this
// guards is that the watcher accurately points at the hang, which is what
// localized the defect in the first place.
//
// The spike note in the entity body is load-bearing for the replay shape: an
// all-at-once push misrepresents the bug (expectDispatchClose would capture its
// baseline AFTER the close was already counted, so the close looks like it never
// fired). The replay MUST drip lines one-per-poll, which is exactly what
// drippingLineSource does — drain() yields at most one line per call.
func TestSonnetTeamDeleteHangReplay(t *testing.T) {
	lines := loadStreamFixture(t, "sonnet_teamdelete_hang.stream.jsonl")

	src := &drippingLineSource{lines: lines}
	proc := &fakeProc{} // never exits — the recorded FO subprocess hung at teardown.

	// Capture every drained line into a transcript so step 3 can assert the
	// watcher OBSERVED the captured TeamDelete-failure, not merely that a
	// never-exiting proc trips expectExit. Without this the exit assertion is
	// tautological: a fakeProc that never exits trips stepTimeout for ANY stream,
	// even one truncated to drop the failure evidence. The transcript anchors the
	// assertion to the recording's actual content (cycle-1 audit fix).
	rec := &transcriptRecorder{}

	// Shrink the budgets so the offline replay finishes in well under a second;
	// the production budgets are 60s. The poll interval matches the unit test's
	// fast cadence. One line drips per poll, so the whole 341-line stream drains
	// over ~341 polls before the steps reach their assertions.
	w := newStreamWatcher(src, proc, rec.tee)
	w.quietBudget = 2 * time.Second
	w.exitBudget = 150 * time.Millisecond
	w.pollInterval = time.Millisecond

	// 1. TeamCreate — the first progress beat. The captured stream opens with the
	//    FO's TeamCreate assistant tool_use (after the contract-gate Bash probes).
	//    isTeamCreate itself lives in the //go:build live live_test.go, so this
	//    default-tagged replay uses the equivalent isToolUse predicate.
	if _, err := w.expect(isToolUse("TeamCreate"), w.quietBudget, "TeamCreate"); err != nil {
		t.Fatalf("replay: TeamCreate must match in the captured stream: %v", err)
	}

	// 2. The single ensign dispatch closes. In this recording the close fires on
	//    the system task_notification status=completed keyed by the Agent
	//    tool_use_id — it fires because the PREMATURE shutdown_request made the
	//    ensign terminate early. The close IS observed; the watcher is correct here.
	if err := w.expectDispatchClose(w.quietBudget, "dispatch close"); err != nil {
		t.Fatalf("replay: the ensign dispatch close must fire (the close anchor is fine): %v", err)
	}
	if w.closedCount != 1 {
		t.Fatalf("replay: exactly one dispatch should have closed, got closedCount=%d", w.closedCount)
	}

	// 3. The recording must carry the diagnosed terminal-teardown failure, and the
	//    watcher must hang at exit BECAUSE of it. Two coupled assertions:
	//
	//    (a) the watcher OBSERVED the captured TeamDelete-failure: the post-close
	//        transcript contains the terminal TeamDelete tool_use AND its
	//        "active member(s)" failure result AND the turn ending with
	//        terminal_reason=completed (the FO going idle WITHOUT retrying). This
	//        ties the test to the recording's content — truncating the stream to
	//        drop L271-341 (the close, the TeamDelete, its failure, the no-retry
	//        turn-end) removes these substrings and turns the test RED.
	//
	//    (b) given that observed failure, expectExit hangs and kills the
	//        never-exiting proc. The hang is a CONSEQUENCE of the recorded
	//        no-retry-after-failure, not an artifact of the fake proc alone.
	if _, err := w.expectExit(w.exitBudget); err != nil {
		var st *stepTimeout
		if !errors.As(err, &st) {
			t.Fatalf("replay: expectExit must trip a *stepTimeout on the never-exiting proc, got %T: %v", err, err)
		}
		if st.label != "expect_exit" {
			t.Errorf("replay: stepTimeout.label = %q, want expect_exit", st.label)
		}
		// The localized failure message must name the exit step so a human reading
		// a CI failure sees "did not exit", not a mislocalized dispatch-close stall.
		if !strings.Contains(st.Error(), "did not exit") {
			t.Errorf("replay: the trip message must localize the hang at exit: %q", st.Error())
		}
	} else {
		t.Fatal("replay: expectExit must NOT return cleanly — the recorded subprocess hung at teardown")
	}
	if !proc.wasKilled() {
		t.Error("replay: expectExit must kill the never-exiting subprocess on timeout")
	}

	// (a) The recording's diagnosed terminal-teardown beats, observed AFTER the
	//     dispatch close. assertObservedAfterClose fails if any beat is missing or
	//     out of order — so a truncated/altered fixture that drops the failure
	//     evidence can no longer pass on the never-exiting proc alone.
	transcript := rec.lines()
	closeIdx := indexOfContainsAll(transcript, `"subtype":"task_notification"`, `"status":"completed"`)
	if closeIdx < 0 {
		t.Fatalf("replay: the dispatch-close task_notification must appear in the transcript")
	}
	postClose := transcript[closeIdx:]
	assertObservedInOrder(t, postClose,
		`"name":"TeamDelete"`,                         // the terminal teardown call
		`Cannot cleanup team with 1 active member(s)`, // its failure — the diagnosed race
		`"terminal_reason":"completed"`,               // the FO turn ends WITHOUT retrying
	)

	// And the recording must NOT contain a clean-exit beat after the failure —
	// no result entry reporting a zero exit that would mean the FO recovered. A
	// recovering (opus-like) stream would carry a later successful TeamDelete; its
	// absence here is the no-retry hang. (The failure result is the only
	// TeamDelete result in the post-close window.)
	if got := countOccurrences(postClose, `"name":"TeamDelete"`); got != 1 {
		t.Errorf("replay: the recording must show exactly one terminal TeamDelete (the failed, un-retried one), got %d — a retried/recovered stream is the opus path, not this sonnet hang", got)
	}
}

// loadStreamFixture reads the captured stream-json fixture under testdata/ and
// returns its lines. testdata/ is the conventional Go fixture directory the go
// tool ignores; the fixture is the verbatim stream the live test's t.Log tee
// recorded, one JSONL line per stream entry.
func loadStreamFixture(t *testing.T, name string) []string {
	t.Helper()
	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("open stream fixture %q: %v", name, err)
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	// The captured init entries carry the full tools/skills manifest (~99KB on
	// the largest line), so raise the scanner buffer well past bufio's 64KB cap.
	sc.Buffer(make([]byte, 0, 256*1024), 4*1024*1024)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan stream fixture %q: %v", name, err)
	}
	if len(lines) == 0 {
		t.Fatalf("stream fixture %q is empty", name)
	}
	return lines
}

// drippingLineSource replays a captured stream one line per drain() call — the
// faithful one-per-poll cadence the spike found necessary. An all-at-once push
// (the fakeLineSource shape) would let expectDispatchClose capture its baseline
// AFTER the close was already counted, so the close would look like it never
// fired; dripping one line per poll reproduces the live drainer's "new lines
// since the last poll" cadence, where the close lands on its own poll.
type drippingLineSource struct {
	mu    sync.Mutex
	lines []string
	next  int
}

// drain returns at most one line per call — the next un-replayed line, or
// nothing once the stream is exhausted (modelling the live drainer going quiet
// after the last line, which is exactly the recorded hang: the subprocess never
// exits and emits no further lines).
func (s *drippingLineSource) drain() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.next >= len(s.lines) {
		return nil
	}
	line := s.lines[s.next]
	s.next++
	return []string{line}
}

// transcriptRecorder captures every line the watcher tees, in drain order, so a
// test can assert which captured beats the watcher actually observed. The
// watcher tees from its poll loop (a goroutine in the live test); guard with a
// mutex even though this replay drives it single-goroutine.
type transcriptRecorder struct {
	mu  sync.Mutex
	rec []string
}

func (r *transcriptRecorder) tee(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rec = append(r.rec, line)
}

func (r *transcriptRecorder) lines() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, len(r.rec))
	copy(out, r.rec)
	return out
}

// indexOfContains returns the index of the first line containing sub, or -1.
func indexOfContains(lines []string, sub string) int {
	for i, l := range lines {
		if strings.Contains(l, sub) {
			return i
		}
	}
	return -1
}

// indexOfContainsAll returns the index of the first line containing every sub,
// or -1. Used where a single JSON line carries the anchor fields non-adjacently
// (the captured stream interleaves other keys between subtype and status).
func indexOfContainsAll(lines []string, subs ...string) int {
	for i, l := range lines {
		all := true
		for _, sub := range subs {
			if !strings.Contains(l, sub) {
				all = false
				break
			}
		}
		if all {
			return i
		}
	}
	return -1
}

// countOccurrences returns how many lines contain sub.
func countOccurrences(lines []string, sub string) int {
	n := 0
	for _, l := range lines {
		if strings.Contains(l, sub) {
			n++
		}
	}
	return n
}

// assertObservedInOrder fails unless each substring appears in lines, each at a
// position at or after the previous one — the diagnosed beats in their recorded
// order. A truncated or reordered fixture that drops the failure evidence fails
// here, so the exit-hang assertion is no longer tautological.
func assertObservedInOrder(t *testing.T, lines []string, subs ...string) {
	t.Helper()
	from := 0
	for _, sub := range subs {
		idx := indexOfContains(lines[from:], sub)
		if idx < 0 {
			t.Errorf("replay: the recording must contain the diagnosed beat %q in order after the prior beat", sub)
			return
		}
		from += idx
	}
}
