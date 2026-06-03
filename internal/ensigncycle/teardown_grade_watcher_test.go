// ABOUTME: AC-1 machinery proof (offline) — drives the live watcher's
// ABOUTME: expectTerminalTeardownGrade over synthetic marker+hold / bug-shape streams, no model spend.
package ensigncycle

import (
	"errors"
	"testing"
	"time"
)

// markerLine is a synthetic assistant text-block carrying the verbatim
// terminal-status marker, as the FO emits it on cap-exhaustion.
const markerLine = `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"TERMINAL_TEARDOWN_BOUNDED: best-effort teardown exhausted; member(s) stuck in registry; holding for launcher."}]}}`

// teamDeleteLine is a synthetic TeamDelete tool_use — a terminal-teardown
// attempt the FO must STOP issuing once it emits the marker.
const teamDeleteLine = `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","id":"tu","name":"TeamDelete","input":{}}]}}`

// activeMemberFailLine is the TeamDelete failure result — the dead-but-listed
// member the bounded loop attempts past.
const activeMemberFailLine = `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"tu","content":[{"type":"text","text":"{\"success\":false,\"message\":\"Cannot cleanup team with 1 active member(s): x\"}"}]}]}}`

// newGradeWatcher builds a watcher with shrunk budgets so the offline grade
// finishes in well under a second.
func newGradeWatcher(src *fakeLineSource, proc *fakeProc) *streamWatcher {
	w := newStreamWatcher(src, proc, func(string) {})
	w.quietBudget = 200 * time.Millisecond
	w.exitBudget = 200 * time.Millisecond
	w.pollInterval = 5 * time.Millisecond
	return w
}

// TestTerminalTeardownGradePassesOnMarkerThenHold proves the AC-1 machinery
// greens on the fix's unique tail: a BOUNDED count of TeamDelete attempts → the
// terminal-status marker → a HOLD. No clean self-exit is required — the launcher
// (here, proc.setExited modelling the deferred poller.kill()) reaps the
// subprocess and the grade returns nil.
func TestTerminalTeardownGradePassesOnMarkerThenHold(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newGradeWatcher(src, proc)

	go func() {
		// Bounded attempts (≤cap): two TeamDelete failures.
		src.push(teamDeleteLine, activeMemberFailLine)
		time.Sleep(2 * w.pollInterval)
		src.push(teamDeleteLine, activeMemberFailLine)
		time.Sleep(2 * w.pollInterval)
		// Then STOP and emit the marker, then HOLD (no further teardown).
		src.push(markerLine)
		time.Sleep(2 * w.pollInterval)
		// The launcher kills the held subprocess — the grade must PASS.
		proc.setExited(-1)
	}()

	if err := w.expectTerminalTeardownGrade(w.quietBudget, w.exitBudget); err != nil {
		t.Fatalf("the bounded marker+hold tail must PASS the grade, got: %v", err)
	}
}

// TestTerminalTeardownGradeFailsWhenMarkerNeverEmitted proves the grade REDs the
// bug shapes' shared property — no marker. Models the post-yy retry-loop: the FO
// keeps attempting TeamDelete and the launcher kills it WITHOUT a marker ever
// appearing. The grade must trip a stepFailure/stepTimeout, not pass.
func TestTerminalTeardownGradeFailsWhenMarkerNeverEmitted(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newGradeWatcher(src, proc)

	go func() {
		// An unbounded retry loop: TeamDelete failures, NO marker, then the
		// launcher kills the never-holding subprocess.
		for i := 0; i < 4; i++ {
			src.push(teamDeleteLine, activeMemberFailLine)
			time.Sleep(w.pollInterval)
		}
		proc.setExited(-1)
	}()

	err := w.expectTerminalTeardownGrade(w.quietBudget, w.exitBudget)
	if err == nil {
		t.Fatal("a stream with NO marker must RED the grade — the retry-loop bug must not pass")
	}
	var sf *stepFailure
	var st *stepTimeout
	if !errors.As(err, &sf) && !errors.As(err, &st) {
		t.Fatalf("want *stepFailure or *stepTimeout on a no-marker stream, got %T: %v", err, err)
	}
}

// TestTerminalTeardownGradeFailsWhenRetryContinuesAfterMarker proves the HOLD
// half live: even WITH the marker, a teardown tool_use emitted AFTER it (the FO
// kept retrying instead of holding) must RED. This is the discriminator a
// marker-presence-only grade would miss.
func TestTerminalTeardownGradeFailsWhenRetryContinuesAfterMarker(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newGradeWatcher(src, proc)

	go func() {
		src.push(teamDeleteLine, activeMemberFailLine)
		time.Sleep(w.pollInterval)
		src.push(markerLine)
		time.Sleep(2 * w.pollInterval)
		// The FO did NOT hold — it kept retrying after the marker.
		src.push(teamDeleteLine, activeMemberFailLine)
	}()

	err := w.expectTerminalTeardownGrade(w.quietBudget, w.exitBudget)
	if err == nil {
		t.Fatal("a teardown tool_use AFTER the marker must RED — the FO must hold, not keep retrying")
	}
	var sf *stepFailure
	if !errors.As(err, &sf) {
		t.Fatalf("want *stepFailure on a continued-retry-after-marker stream, got %T: %v", err, err)
	}
}
