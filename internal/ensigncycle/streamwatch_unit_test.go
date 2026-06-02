// ABOUTME: Offline unit test for the live cycle's streamWatcher: synthetic JSONL
// ABOUTME: + a fake exit-poller exercise the no-progress quiet trip, reset-on-activity, and early-exit failure (no live model).
package ensigncycle

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// These run under the DEFAULT build tags (no //go:build live) so `go test ./...`
// covers them; they spend NO model. They mirror the upstream Python
// test_fo_stream_watcher.py split: a synthetic JSONL line source + a fake
// exit-poller (the Go port of _FakeProc) drive the watcher's expect loop. The
// load-bearing semantic these prove is the one place the Go port intentionally
// differs from upstream expect(): the per-step budget is a NO-PROGRESS quiet
// bound — the deadline resets to now+budget on every drained line — not a
// total-stage cap.

// fakeLineSource hands out synthetic stream-json lines on demand, the offline
// stand-in for the StdoutPipe line drainer. push() appends lines a later drain()
// will return; each drain() returns only the lines pushed since the previous
// drain (the "new complete lines since last poll" contract the real drainer also
// honors). holdPartial models a half-written trailing line the drainer holds
// until its newline arrives.
type fakeLineSource struct {
	mu      sync.Mutex
	pending []string
}

func (s *fakeLineSource) push(lines ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pending = append(s.pending, lines...)
}

func (s *fakeLineSource) drain() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := s.pending
	s.pending = nil
	return out
}

// fakeProc is the Go port of the Python _FakeProc: poll() reports the exit code
// once setExited has flipped it, and reports not-exited before then. Guarded by
// a mutex because the watcher polls it from a goroutine while the test flips it.
type fakeProc struct {
	mu     sync.Mutex
	code   int
	exited bool
	killed bool
}

func (p *fakeProc) poll() (int, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.code, p.exited
}

func (p *fakeProc) setExited(code int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.code = code
	p.exited = true
}

func (p *fakeProc) kill() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.killed = true
}

func (p *fakeProc) wasKilled() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.killed
}

// A tiny predicate over the parsed entry: matches an assistant tool_use for the
// named tool. The real watcher uses the same parsed-entry shape against the live
// stream-json.
func isToolUse(name string) func(streamEntry) bool {
	return func(e streamEntry) bool {
		b := e.toolUseBlock()
		return b != nil && b.Name == name
	}
}

// newTestWatcher builds a watcher with a SHRUNK quiet/exit budget and a fast
// poll interval so the offline tests run in well under a second; the production
// budgets are 60s.
func newTestWatcher(src *fakeLineSource, proc *fakeProc) *streamWatcher {
	w := newStreamWatcher(src, proc, func(string) {})
	w.quietBudget = 150 * time.Millisecond
	w.exitBudget = 150 * time.Millisecond
	w.pollInterval = 5 * time.Millisecond
	return w
}

// TestExpectReturnsMatchedEntry: expect returns the matched entry on success.
func TestExpectReturnsMatchedEntry(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newTestWatcher(src, proc)

	src.push(toolUseLine("Bash", `"command":"echo one"`))
	src.push(toolUseLine("TeamCreate", `"team_name":"t"`))

	entry, err := w.expect(isToolUse("TeamCreate"), w.quietBudget, "TeamCreate")
	if err != nil {
		t.Fatalf("expect returned error: %v", err)
	}
	if b := entry.toolUseBlock(); b == nil || b.Name != "TeamCreate" {
		t.Fatalf("expect returned the wrong entry: %+v", entry)
	}
}

// TestExpectQuietTimeoutCarriesLabel: a stream that goes silent before the
// predicate matches trips StepTimeout carrying the step label. This is the
// localized no-progress trip — not a total-stage cap.
func TestExpectQuietTimeoutCarriesLabel(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newTestWatcher(src, proc)

	src.push(toolUseLine("TeamCreate", `"team_name":"t"`))

	start := time.Now()
	_, err := w.expect(isToolUse("Agent"), w.quietBudget, "dispatch close")
	elapsed := time.Since(start)

	var st *stepTimeout
	if !errors.As(err, &st) {
		t.Fatalf("want *stepTimeout, got %T: %v", err, err)
	}
	if st.label != "dispatch close" {
		t.Errorf("stepTimeout.label = %q, want %q", st.label, "dispatch close")
	}
	if !strings.Contains(st.Error(), "dispatch close") {
		t.Errorf("error message must name the step: %q", st.Error())
	}
	// Tripped on the quiet budget, NOT a total-stage cap: it returns within
	// ~1x the (shrunk) quiet budget once the stream goes silent.
	if elapsed > 3*w.quietBudget {
		t.Errorf("quiet trip took %s, far beyond ~1x the %s budget — not localized", elapsed, w.quietBudget)
	}
}

// TestExpectResetsDeadlineOnActivity: a noisy-but-unmatched stream that keeps
// emitting lines PAST the quiet budget does NOT trip — the deadline resets on
// every drained line — and expect returns once the predicate finally matches.
// This is the single intentional delta from upstream's overall budget.
func TestExpectResetsDeadlineOnActivity(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newTestWatcher(src, proc)

	// Drive unrelated lines for well past the quiet budget, one per poll, so the
	// total elapsed exceeds quietBudget several times over. A total-stage cap
	// would trip; the no-progress bound must not, because each drained line
	// resets the deadline.
	matched := make(chan error, 1)
	go func() {
		_, err := w.expect(isToolUse("Agent"), w.quietBudget, "agent dispatch")
		matched <- err
	}()

	noiseEnd := time.Now().Add(4 * w.quietBudget)
	for time.Now().Before(noiseEnd) {
		src.push(toolUseLine("Bash", `"command":"heartbeat"`))
		time.Sleep(w.pollInterval)
	}
	// Now emit the matching line; expect must return cleanly despite total
	// elapsed already exceeding the quiet budget many times.
	src.push(toolUseLine("Agent", `"subagent_type":"spacedock:ensign"`))

	select {
	case err := <-matched:
		if err != nil {
			t.Fatalf("noisy-but-progressing stream must NOT trip the quiet budget; got %v", err)
		}
	case <-time.After(2 * w.quietBudget):
		t.Fatalf("expect did not return after the matching line arrived")
	}
}

// TestExpectStepFailureOnEarlyExit: the subprocess exits before the predicate
// matches → StepFailure carrying the non-zero exit code + the label.
func TestExpectStepFailureOnEarlyExit(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newTestWatcher(src, proc)

	src.push(toolUseLine("Bash", `"command":"noise"`))

	go func() {
		time.Sleep(2 * w.pollInterval)
		proc.setExited(1)
	}()

	_, err := w.expect(isToolUse("Agent"), 5*time.Second, "agent dispatch")
	var sf *stepFailure
	if !errors.As(err, &sf) {
		t.Fatalf("want *stepFailure, got %T: %v", err, err)
	}
	if sf.exitCode != 1 {
		t.Errorf("stepFailure.exitCode = %d, want 1", sf.exitCode)
	}
	if sf.label != "agent dispatch" {
		t.Errorf("stepFailure.label = %q, want %q", sf.label, "agent dispatch")
	}
	if !strings.Contains(sf.Error(), "code=1") {
		t.Errorf("error must carry the exit code: %q", sf.Error())
	}
}

// TestExpectFinalDrainBeforePoll: when the matching line is already drained
// and the proc has exited, expect must drain-then-check so the match wins over
// the exit (the ordering the Python port guarantees).
func TestExpectFinalDrainBeforePoll(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newTestWatcher(src, proc)

	src.push(toolUseLine("Agent", `"subagent_type":"spacedock:ensign"`))
	proc.setExited(0)

	entry, err := w.expect(isToolUse("Agent"), w.quietBudget, "agent dispatch")
	if err != nil {
		t.Fatalf("expect must match the final-flushed line before reporting exit; got %v", err)
	}
	if b := entry.toolUseBlock(); b == nil || b.Name != "Agent" {
		t.Fatalf("wrong entry returned: %+v", entry)
	}
}

// TestExpectFailureCarriesTranscriptTail: a hang leaves a step-naming tail that
// includes the lines streamed so far (AC-2's diagnosability proof, offline).
func TestExpectFailureCarriesTranscriptTail(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newTestWatcher(src, proc)

	// Advance through TeamCreate, then go silent before the dispatch-close.
	src.push(toolUseLine("TeamCreate", `"team_name":"t"`))
	if _, err := w.expect(isToolUse("TeamCreate"), w.quietBudget, "TeamCreate"); err != nil {
		t.Fatalf("TeamCreate should have matched: %v", err)
	}

	_, err := w.expect(isToolUse("Agent"), w.quietBudget, "dispatch close")
	var st *stepTimeout
	if !errors.As(err, &st) {
		t.Fatalf("want *stepTimeout, got %T: %v", err, err)
	}
	msg := st.Error()
	if !strings.Contains(msg, "dispatch close") {
		t.Errorf("tail-bearing error must name the stalled step: %q", msg)
	}
	if !strings.Contains(msg, "TeamCreate") {
		t.Errorf("transcript tail must include the last reached step (TeamCreate): %q", msg)
	}
}

// TestExpectDispatchCloseOnThreeAnchors: a dispatch opens on an
// Agent(subagent_type="spacedock:ensign") tool_use and closes on each of the
// three mode anchors the port carries — bare-mode Agent tool_result, teams-mode
// task_notification status=completed, and the headless inbox-poll Done: Bash
// tool_result. Each anchor, on its own, must close the open dispatch.
func TestExpectDispatchCloseOnThreeAnchors(t *testing.T) {
	cases := []struct {
		name      string
		openTUID  string
		closeLine string
	}{
		{
			name:      "bare_mode_agent_tool_result",
			openTUID:  "toolu_bare",
			closeLine: bareDoneToolResultLine("toolu_bare"),
		},
		{
			name:      "teams_task_notification_completed",
			openTUID:  "toolu_teams",
			closeLine: taskNotificationCompletedLine("toolu_teams"),
		},
		{
			name:      "headless_inbox_poll_done",
			openTUID:  "toolu_headless",
			closeLine: inboxPollDoneLine("spacedock-ensign-make-it-work-done"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := &fakeLineSource{}
			proc := &fakeProc{}
			w := newTestWatcher(src, proc)

			src.push(agentDispatchLine(tc.openTUID, "Create a greeting file: done"))
			src.push(tc.closeLine)

			if err := w.expectDispatchClose(w.quietBudget, "dispatch close"); err != nil {
				t.Fatalf("dispatch close on the %s anchor did not register: %v", tc.name, err)
			}
		})
	}
}

// TestExpectDispatchCloseSpawnAckDoesNotClose: a teams-mode Agent tool_result
// that is just the "Spawned successfully" ack must NOT close the dispatch — only
// a real Done payload (or the other anchors) does. Guards the bare-vs-teams
// discriminator the port carries.
func TestExpectDispatchCloseSpawnAckDoesNotClose(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newTestWatcher(src, proc)

	src.push(agentDispatchLine("toolu_x", "Create a greeting file: done"))
	src.push(spawnAckToolResultLine("toolu_x"))

	err := w.expectDispatchClose(w.quietBudget, "dispatch close")
	var st *stepTimeout
	if !errors.As(err, &st) {
		t.Fatalf("a spawn-ack must NOT close the dispatch; want *stepTimeout, got %T: %v", err, err)
	}
}

// TestExpectExitWaitsThenKills: expectExit returns the exit code once the proc
// exits; on a hung proc it trips StepTimeout and kills it.
func TestExpectExitWaitsThenKills(t *testing.T) {
	t.Run("exits_cleanly", func(t *testing.T) {
		src := &fakeLineSource{}
		proc := &fakeProc{}
		w := newTestWatcher(src, proc)
		go func() {
			time.Sleep(2 * w.pollInterval)
			proc.setExited(0)
		}()
		code, err := w.expectExit(w.exitBudget)
		if err != nil {
			t.Fatalf("expectExit on a clean exit returned error: %v", err)
		}
		if code != 0 {
			t.Errorf("exit code = %d, want 0", code)
		}
	})

	t.Run("hung_proc_trips_and_is_killed", func(t *testing.T) {
		src := &fakeLineSource{}
		proc := &fakeProc{} // never exits
		w := newTestWatcher(src, proc)

		_, err := w.expectExit(w.exitBudget)
		var st *stepTimeout
		if !errors.As(err, &st) {
			t.Fatalf("want *stepTimeout on a hung exit, got %T: %v", err, err)
		}
		if st.label != "expect_exit" {
			t.Errorf("stepTimeout.label = %q, want expect_exit", st.label)
		}
		if !proc.wasKilled() {
			t.Error("expectExit must kill a hung subprocess on timeout")
		}
	})
}

// --- synthetic stream-json line builders ---------------------------------
//
// These emit the standard Claude Code stream-json shapes the watcher parses.
// `spacedock claude` syscall.Exec-replaces itself with `claude`, forwarding
// --output-format stream-json verbatim, so the live pipe carries exactly these
// shapes (Spike in the entity body). Each builder returns one JSONL line.

// toolUseLine builds an assistant tool_use entry. extraInputJSON is raw JSON
// key/value fragments (e.g. `"command":"echo one"`) spliced into the input
// object — enough for the predicate to discriminate on the tool name + a field.
func toolUseLine(toolName, extraInputJSON string) string {
	return `{"type":"assistant","message":{"model":"claude-sonnet-4-6","content":[` +
		`{"type":"tool_use","id":"toolu_anon","name":` + strconv.Quote(toolName) +
		`,"input":{` + extraInputJSON + `}}]}}`
}

// agentDispatchLine opens a dispatch: an assistant Agent tool_use with
// subagent_type=spacedock:ensign and a description, keyed by tool_use_id.
func agentDispatchLine(toolUseID, description string) string {
	return `{"type":"assistant","message":{"model":"claude-sonnet-4-6","content":[` +
		`{"type":"tool_use","id":` + strconv.Quote(toolUseID) + `,"name":"Agent","input":{` +
		`"subagent_type":"spacedock:ensign","description":` + strconv.Quote(description) + `}}]}}`
}

// bareDoneToolResultLine is the bare-mode close anchor: the user tool_result for
// the Agent tool_use_id carrying the teammate's Done payload directly (a list of
// text blocks, NOT a "Spawned successfully" spawn-ack).
func bareDoneToolResultLine(toolUseID string) string {
	return `{"type":"user","message":{"content":[` +
		`{"type":"tool_result","tool_use_id":` + strconv.Quote(toolUseID) +
		`,"content":[{"type":"text","text":"Done: completed done. Report written."}]}]}}`
}

// spawnAckToolResultLine is the teams-mode Agent tool_result that is JUST the
// spawn ack — it must NOT close the dispatch.
func spawnAckToolResultLine(toolUseID string) string {
	return `{"type":"user","message":{"content":[` +
		`{"type":"tool_result","tool_use_id":` + strconv.Quote(toolUseID) +
		`,"content":[{"type":"text","text":"Spawned successfully. agent_id: a123-xyz"}]}]}}`
}

// taskNotificationCompletedLine is the teams-mode-with-TTY close anchor: a
// system task_notification with status=completed keyed by the Agent tool_use_id.
func taskNotificationCompletedLine(toolUseID string) string {
	return `{"type":"system","subtype":"task_notification","status":"completed","tool_use_id":` +
		strconv.Quote(toolUseID) + `}`
}

// inboxPollDoneLine is the headless `claude -p` close anchor: a Bash tool_result
// whose body carries `from: spacedock-ensign-…-{stage}` and `text: Done:`,
// surfaced by the FO's inbox-poll script. Closes any open dispatch whose
// description shares the stage substring.
func inboxPollDoneLine(sender string) string {
	body := "team: t\nfrom: " + sender + "\ntimestamp: now\nsummary: s\ntext: Done: completed done."
	return `{"type":"user","message":{"content":[` +
		`{"type":"tool_result","tool_use_id":"toolu_bash","content":[{"type":"text","text":` +
		strconv.Quote(body) + `}]}]}}`
}
