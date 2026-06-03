// ABOUTME: AC-1 machinery proof (offline) — drives the live watcher's
// ABOUTME: expectTerminalTeardownGrade over synthetic marker-emission / bug-shape streams, no model spend.
package ensigncycle

import (
	"errors"
	"testing"
	"time"
)

// markerLine is a synthetic assistant text-block carrying the verbatim
// terminal-status marker, as the FO emits it on cap-exhaustion.
const markerLine = `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"TERMINAL_TEARDOWN_BOUNDED: best-effort teardown exhausted; member(s) stuck in registry; holding for launcher."}]}}`

// markerThinkingLine is the marker emitted inside an assistant THINKING block —
// the shape the real sonnet FO produced on PR #285 (it narrated reaching the cap
// and emitting the marker in its thinking). The grade must accept this too.
const markerThinkingLine = `{"type":"assistant","message":{"role":"assistant","content":[{"type":"thinking","thinking":"I've hit the attempt cap. Emit TERMINAL_TEARDOWN_BOUNDED: best-effort teardown exhausted; member(s) stuck in registry; holding for launcher."}]}}`

// teamDeleteLine is a synthetic TeamDelete tool_use — a terminal-teardown attempt.
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

// TestTerminalTeardownGradePassesOnMarkerEmission proves the AC-1 machinery
// greens on the fix's unique signal: a BOUNDED count of TeamDelete attempts →
// the FO EMITS the terminal-status marker. No clean self-exit and no clean hold
// are required — the launcher (here, proc.setExited modelling the deferred
// poller.kill()) reaps the still-running subprocess and the grade returns nil as
// soon as the marker is authored.
func TestTerminalTeardownGradePassesOnMarkerEmission(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newGradeWatcher(src, proc)

	go func() {
		// Bounded attempts (≤cap): two TeamDelete failures.
		src.push(teamDeleteLine, activeMemberFailLine)
		time.Sleep(2 * w.pollInterval)
		src.push(teamDeleteLine, activeMemberFailLine)
		time.Sleep(2 * w.pollInterval)
		// Then the FO emits the marker — the grade passes here.
		src.push(markerLine)
		time.Sleep(2 * w.pollInterval)
		proc.setExited(-1)
	}()

	if err := w.expectTerminalTeardownGrade(w.quietBudget); err != nil {
		t.Fatalf("the bounded marker-emission tail must PASS the grade, got: %v", err)
	}
}

// TestTerminalTeardownGradePassesOnMarkerInThinking proves the grade accepts the
// marker emitted in an assistant THINKING block — the real PR #285 sonnet shape.
func TestTerminalTeardownGradePassesOnMarkerInThinking(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newGradeWatcher(src, proc)

	go func() {
		src.push(teamDeleteLine, activeMemberFailLine)
		time.Sleep(w.pollInterval)
		src.push(markerThinkingLine)
		time.Sleep(2 * w.pollInterval)
		proc.setExited(-1)
	}()

	if err := w.expectTerminalTeardownGrade(w.quietBudget); err != nil {
		t.Fatalf("a marker emitted in a thinking block must PASS the grade, got: %v", err)
	}
}

// TestTerminalTeardownGradePassesWhenTeardownContinuesAfterMarker is the PR #285
// reality check: the real sonnet FO emits the marker but RESUMES teardown on the
// next harness re-invoke (it does not cleanly hold). The grade keys on marker
// EMISSION, so continued teardown after the marker must STILL PASS — requiring a
// clean post-marker hold over-fit the offline recordings and failed the real
// producer (the whole reason this cycle exists).
func TestTerminalTeardownGradePassesWhenTeardownContinuesAfterMarker(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newGradeWatcher(src, proc)

	go func() {
		src.push(teamDeleteLine, activeMemberFailLine)
		time.Sleep(w.pollInterval)
		src.push(markerLine)
		time.Sleep(2 * w.pollInterval)
		// The FO resumes teardown after the marker (re-invoke reality) — still PASS.
		src.push(teamDeleteLine, activeMemberFailLine)
		time.Sleep(2 * w.pollInterval)
		proc.setExited(-1)
	}()

	if err := w.expectTerminalTeardownGrade(w.quietBudget); err != nil {
		t.Fatalf("continued teardown AFTER the marker must STILL PASS (the real FO does not hold), got: %v", err)
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
		// launcher kills the never-terminusing subprocess.
		for i := 0; i < 4; i++ {
			src.push(teamDeleteLine, activeMemberFailLine)
			time.Sleep(w.pollInterval)
		}
		proc.setExited(-1)
	}()

	err := w.expectTerminalTeardownGrade(w.quietBudget)
	if err == nil {
		t.Fatal("a stream with NO marker must RED the grade — the retry-loop bug must not pass")
	}
	var sf *stepFailure
	var st *stepTimeout
	if !errors.As(err, &sf) && !errors.As(err, &st) {
		t.Fatalf("want *stepFailure or *stepTimeout on a no-marker stream, got %T: %v", err, err)
	}
}

// TestTerminalTeardownGradeIgnoresMarkerInContractRead proves the grade keys on
// the FO AUTHORING the marker, not on the marker appearing in a contract-Read
// tool_result. The FO Reads the marker-bearing contract files at startup; that
// is not reaching the terminus. A stream whose ONLY marker is in a user
// tool_result (and never authored by the assistant) must RED.
func TestTerminalTeardownGradeIgnoresMarkerInContractRead(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newGradeWatcher(src, proc)

	// A user tool_result carrying the marker verbatim — the FO read the contract,
	// it did NOT emit the marker.
	contractReadLine := `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"r","content":[{"type":"text","text":"10. ...emit the marker TERMINAL_TEARDOWN_BOUNDED: best-effort teardown exhausted; member(s) stuck in registry; holding for launcher. ..."}]}]}}`
	go func() {
		src.push(contractReadLine)
		time.Sleep(w.pollInterval)
		src.push(teamDeleteLine, activeMemberFailLine)
		time.Sleep(w.pollInterval)
		proc.setExited(-1)
	}()

	err := w.expectTerminalTeardownGrade(w.quietBudget)
	if err == nil {
		t.Fatal("a marker that only appears in a contract-Read tool_result must RED — the FO did not author it")
	}
}
