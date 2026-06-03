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

	// Shrink the budgets so the offline replay finishes in well under a second;
	// the production budgets are 60s. The poll interval matches the unit test's
	// fast cadence. One line drips per poll, so the whole 341-line stream drains
	// over ~341 polls before the steps reach their assertions.
	w := newStreamWatcher(src, proc, func(string) {})
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

	// 3. expectExit hangs and kills the never-exiting proc. This is the localized
	//    failure: terminal TeamDelete failed with "active member(s)" and the FO
	//    never retried it to success, so the recorded subprocess never exits.
	//    expectExit must trip a stepTimeout("expect_exit") AND kill the proc.
	_, err := w.expectExit(w.exitBudget)
	var st *stepTimeout
	if !errors.As(err, &st) {
		t.Fatalf("replay: expectExit must trip a *stepTimeout on the never-exiting proc, got %T: %v", err, err)
	}
	if st.label != "expect_exit" {
		t.Errorf("replay: stepTimeout.label = %q, want expect_exit", st.label)
	}
	if !proc.wasKilled() {
		t.Error("replay: expectExit must kill the never-exiting subprocess on timeout")
	}

	// The localized failure message must name the exit step so a human reading a
	// CI failure sees "did not exit", not a mislocalized dispatch-close stall.
	if !strings.Contains(st.Error(), "did not exit") {
		t.Errorf("replay: the trip message must localize the hang at exit: %q", st.Error())
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
