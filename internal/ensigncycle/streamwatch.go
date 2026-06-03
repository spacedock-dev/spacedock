// ABOUTME: Streaming per-step watcher over the FO's stream-json: bounds each live
// ABOUTME: stage by a no-progress quiet budget, replacing the monolithic timeout (Go port of FOStreamWatcher).
package ensigncycle

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// streamWatcher consumes the FO's stream-json JSONL as it arrives and bounds
// each stage-progress step by its own no-progress (quiet) budget. It is the Go
// port of the upstream Python FOStreamWatcher (scripts/test_lib.py), specialised
// to the v1 live cycle. expect / expectDispatchClose / expectExit each poll for
// new lines, tee them to the test log immediately so a hang leaves a partial
// transcript, and trip a labelled stepTimeout / stepFailure that names the
// stalled step.
//
// The one intentional delta from upstream expect(): the per-step budget is a
// NO-PROGRESS bound, not an overall cap. The deadline resets to now+budget on
// every drained line, so a live stage that legitimately runs for minutes never
// trips as long as the stream keeps moving; only genuine silence past the budget
// trips. This reconciles the captain's "no individual timeout > 60s" directive
// with stages that legitimately take minutes.
//
// It compiles under DEFAULT build tags (no //go:build live) so the offline unit
// test exercises it against synthetic logs with no model spend; the //go:build
// live live_test.go wires it to the real subprocess pipe.

const (
	// quietBudgetDefault is the per-step no-progress bound for expect and
	// expectDispatchClose — the deadline resets on any drained line, so this
	// caps stream SILENCE, not total stage wallclock.
	quietBudgetDefault = 60 * time.Second
	// exitBudgetDefault is how long expectExit waits for the FO to exit after
	// the last watched step matched, before killing it.
	exitBudgetDefault = 60 * time.Second
	// holdConfirmDefault is how long expectTerminalTeardownGrade watches AFTER the
	// terminal-status marker to confirm the FO HOLDS (no further teardown
	// tool_use). The marker + a clean hold over this window is the bounded-teardown
	// PASS; the launcher's kill() then reaps the subprocess (no self-exit needed).
	holdConfirmDefault = 30 * time.Second
	// pollIntervalDefault matches upstream POLL_INTERVAL_S = 0.2.
	pollIntervalDefault = 200 * time.Millisecond
	// transcriptTailLines bounds the tail carried in failure messages.
	transcriptTailLines = 20
)

// terminalTeardownMarker is the contract-mandated terminal-status sentinel the FO
// emits on cap-exhaustion of the bounded best-effort teardown (shared-core step
// 10 / the Claude `## Terminal Team Teardown` section). It is the load-bearing
// discriminator the offline grade and the live grade key on: NEITHER bug shape
// emits it — the pre-yy give-up ends the turn silently after one failed
// TeamDelete (sonnet_teamdelete_hang), and the post-yy retry-loop retries
// TeamDelete past the cap and never reaches a marker
// (sonnet_teamdelete_retryloop). Only the fix emits the marker and then HOLDs.
// Kept verbatim in sync with the contract files and the integration lint's
// terminalTeardownMarker constant.
const terminalTeardownMarker = "TERMINAL_TEARDOWN_BOUNDED: best-effort teardown exhausted; member(s) stuck in registry; holding for launcher."

// lineSource yields newly-completed stream-json lines since the previous drain.
// Partial trailing lines are held internally until their newline arrives — the
// port of upstream _drain_entries's buffer. drain never blocks.
type lineSource interface {
	drain() []string
}

// procPoller reports whether the FO subprocess has exited (and its code), and
// kills it on timeout. The port of the Python proc.poll() / proc.kill() surface;
// the offline unit test injects a fake (the port of _FakeProc), the live test
// wraps the real exec.Cmd.
type procPoller interface {
	poll() (code int, exited bool)
	kill()
}

// stepTimeout: no progress within the quiet budget (or no exit within the exit
// budget). Carries the step label and a transcript tail.
type stepTimeout struct {
	label string
	msg   string
}

func (e *stepTimeout) Error() string { return e.msg }

// stepFailure: the FO subprocess exited BEFORE the expected step matched. Carries
// the step label, the non-zero exit code, and a transcript tail.
type stepFailure struct {
	label    string
	exitCode int
	msg      string
}

func (e *stepFailure) Error() string { return e.msg }

// openDispatch tracks an in-flight ensign dispatch from its Agent tool_use to
// its matching close anchor.
type openDispatch struct {
	dispatchID string
	ensignName string
}

type streamWatcher struct {
	src  lineSource
	proc procPoller
	tee  func(string)

	quietBudget  time.Duration
	exitBudget   time.Duration
	pollInterval time.Duration

	transcript     []string
	openDispatches map[string]*openDispatch
	closedCount    int
}

// newStreamWatcher builds a watcher over the given line source and proc poller.
// tee receives every drained line immediately (the live test passes t.Log so a
// hang leaves a partial, step-naming transcript). Budgets default to the ≤60s
// constants; the offline unit test shrinks them for speed.
func newStreamWatcher(src lineSource, proc procPoller, tee func(string)) *streamWatcher {
	return &streamWatcher{
		src:            src,
		proc:           proc,
		tee:            tee,
		quietBudget:    quietBudgetDefault,
		exitBudget:     exitBudgetDefault,
		pollInterval:   pollIntervalDefault,
		openDispatches: map[string]*openDispatch{},
	}
}

// drainEntries reads the newly-available lines, tees them, records them in the
// bounded transcript, updates dispatch open/close tracking, and returns the
// parsed entries. Returns the count of lines drained so callers can reset the
// no-progress deadline on any activity.
func (w *streamWatcher) drainEntries() (entries []streamEntry, drained int) {
	lines := w.src.drain()
	for _, line := range lines {
		drained++
		w.tee(line)
		w.transcript = append(w.transcript, line)
		var e streamEntry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	w.updateDispatches(entries)
	return entries, drained
}

// transcriptTail returns the last transcriptTailLines streamed lines joined by
// newlines — the self-diagnosing tail carried in every failure message.
func (w *streamWatcher) transcriptTail() string {
	start := len(w.transcript) - transcriptTailLines
	if start < 0 {
		start = 0
	}
	return strings.Join(w.transcript[start:], "\n")
}

// expect returns the first drained entry where predicate is true, bounding the
// wait by a NO-PROGRESS quiet budget: the deadline resets to now+budget on every
// drained line. Trips stepTimeout on silence past the budget; stepFailure if the
// proc exits before a match (after a final drain).
func (w *streamWatcher) expect(predicate func(streamEntry) bool, budget time.Duration, label string) (streamEntry, error) {
	deadline := time.Now().Add(budget)
	for {
		entries, drained := w.drainEntries()
		if drained > 0 {
			deadline = time.Now().Add(budget)
		}
		for _, e := range entries {
			if safePredicate(predicate, e) {
				return e, nil
			}
		}

		if _, exited := w.proc.poll(); exited {
			// Final drain so a line that landed alongside exit still wins.
			entries, _ := w.drainEntries()
			for _, e := range entries {
				if safePredicate(predicate, e) {
					return e, nil
				}
			}
			code, _ := w.proc.poll()
			return streamEntry{}, &stepFailure{
				label:    label,
				exitCode: code,
				msg: fmt.Sprintf("FO subprocess exited (code=%d) before step %q matched.\nTranscript tail:\n%s",
					code, label, w.transcriptTail()),
			}
		}

		if time.Now().After(deadline) {
			return streamEntry{}, &stepTimeout{
				label: label,
				msg: fmt.Sprintf("step %q made no progress within %s (no-progress quiet budget).\nTranscript tail:\n%s",
					label, budget, w.transcriptTail()),
			}
		}
		time.Sleep(w.pollInterval)
	}
}

// expectDispatchClose blocks until the next ensign dispatch closes (any of the
// three mode anchors registered it), bounding the wait by the same no-progress
// quiet budget as expect. The flat-entity live cycle dispatches exactly one
// ensign, so the live test calls this once.
func (w *streamWatcher) expectDispatchClose(budget time.Duration, label string) error {
	baseline := w.closedCount
	deadline := time.Now().Add(budget)
	for {
		_, drained := w.drainEntries()
		if drained > 0 {
			deadline = time.Now().Add(budget)
		}
		if w.closedCount > baseline {
			return nil
		}

		if _, exited := w.proc.poll(); exited {
			w.drainEntries()
			if w.closedCount > baseline {
				return nil
			}
			code, _ := w.proc.poll()
			return &stepFailure{
				label:    label,
				exitCode: code,
				msg: fmt.Sprintf("FO subprocess exited (code=%d) before step %q closed.\nTranscript tail:\n%s",
					code, label, w.transcriptTail()),
			}
		}

		if time.Now().After(deadline) {
			return &stepTimeout{
				label: label,
				msg: fmt.Sprintf("step %q did not close within %s (no-progress quiet budget). Open dispatches: %v.\nTranscript tail:\n%s",
					label, budget, w.openDispatchNames(), w.transcriptTail()),
			}
		}
		time.Sleep(w.pollInterval)
	}
}

// expectExit waits up to budget for the FO subprocess to exit AFTER the last
// watched step matched, draining any final lines. On timeout it kills the
// subprocess and trips stepTimeout("expect_exit") so a hung claude does not
// outlive the test.
func (w *streamWatcher) expectExit(budget time.Duration) (int, error) {
	deadline := time.Now().Add(budget)
	for {
		w.drainEntries()
		if code, exited := w.proc.poll(); exited {
			w.drainEntries()
			return code, nil
		}
		if time.Now().After(deadline) {
			w.proc.kill()
			return 0, &stepTimeout{
				label: "expect_exit",
				msg: fmt.Sprintf("FO subprocess did not exit within %s.\nTranscript tail:\n%s",
					budget, w.transcriptTail()),
			}
		}
		time.Sleep(w.pollInterval)
	}
}

// expectTerminalTeardownGrade grades the FO's bounded best-effort terminal
// teardown WITHOUT requiring a clean self-exit (which is impossible: the harness
// will not let claude -p exit while the team's members[] is populated). It waits,
// bounded by the no-progress quiet budget, for the contract-mandated
// terminal-status MARKER to appear in the stream — the unique tail of the fix
// that NEITHER bug shape emits (the pre-yy give-up ends the turn with no marker;
// the post-yy retry-loop retries TeamDelete past the cap with no marker). Once
// the marker is observed, it confirms the HOLD over a bounded window: any
// teardown tool_use (TeamDelete or a shutdown_request SendMessage) emitted AFTER
// the marker FAILS the grade (the FO kept retrying instead of holding). On a clean
// hold it returns nil; the caller's deferred poller.kill() then reaps the
// subprocess and the cycle PASSES.
//
// This grade is the live realization of gradeTerminalTeardown (AC-2): it keys on
// the marker + hold, NOT on the shutdown_request/TeamDelete beats both bug shapes
// also emit — so it greens ONLY on the fix, not on the very run yy's fix failed.
// holdConfirm is the post-marker window over which a continued-retry FAILS; the
// quiet budget bounds how long the watcher waits for the marker itself. Both stay
// ≤60s so the AC-1 timeout guard is unaffected.
func (w *streamWatcher) expectTerminalTeardownGrade(quietBudget, holdConfirm time.Duration) error {
	// Phase 1: wait for the marker (an assistant text block, graded over raw
	// lines), bounded by the no-progress quiet budget.
	deadline := time.Now().Add(quietBudget)
	markerLine := -1
	for markerLine < 0 {
		_, drained := w.drainEntries()
		if drained > 0 {
			deadline = time.Now().Add(quietBudget)
		}
		for i := len(w.transcript) - 1; i >= 0 && i >= len(w.transcript)-drained; i-- {
			if lineHasMarker(w.transcript[i]) {
				markerLine = i
				break
			}
		}
		if markerLine >= 0 {
			break
		}
		if _, exited := w.proc.poll(); exited {
			// A final drain so a marker that landed alongside exit still wins.
			w.drainEntries()
			for i, line := range w.transcript {
				if lineHasMarker(line) {
					markerLine = i
					break
				}
			}
			if markerLine >= 0 {
				break
			}
			code, _ := w.proc.poll()
			return &stepFailure{
				label:    "terminal_teardown_grade",
				exitCode: code,
				msg: fmt.Sprintf("FO subprocess exited (code=%d) before emitting the terminal-status marker — a silent give-up or an unbounded retry loop.\nTranscript tail:\n%s",
					code, w.transcriptTail()),
			}
		}
		if time.Now().After(deadline) {
			return &stepTimeout{
				label: "terminal_teardown_grade",
				msg: fmt.Sprintf("the FO did not emit the terminal-status marker within %s (no-progress quiet budget) — it never reached the bounded-teardown terminus.\nTranscript tail:\n%s",
					quietBudget, w.transcriptTail()),
			}
		}
		time.Sleep(w.pollInterval)
	}

	// Phase 2: confirm the HOLD. Over holdConfirm, any teardown tool_use after the
	// marker means the FO did NOT hold — fail. The deadline resets on a hold-clean
	// line so a chatty hold still confirms; a teardown tool_use trips immediately.
	holdDeadline := time.Now().Add(holdConfirm)
	for {
		w.drainEntries()
		if ok, reason := gradeTerminalTeardown(w.transcript); !ok {
			return &stepFailure{
				label: "terminal_teardown_grade",
				msg: fmt.Sprintf("terminal teardown grade FAILED after the marker: %s.\nTranscript tail:\n%s",
					reason, w.transcriptTail()),
			}
		}
		if _, exited := w.proc.poll(); exited {
			// The launcher killed it (or it exited): the marker held to the end. PASS.
			return nil
		}
		if time.Now().After(holdDeadline) {
			// The hold held for the whole confirmation window with no further
			// teardown tool_use. PASS — the launcher's kill() reaps the subprocess.
			return nil
		}
		time.Sleep(w.pollInterval)
	}
}

// openDispatchNames lists the descriptions of dispatches still open, for the
// timeout message.
func (w *streamWatcher) openDispatchNames() []string {
	names := make([]string, 0, len(w.openDispatches))
	for _, d := range w.openDispatches {
		names = append(names, d.ensignName)
	}
	return names
}

// updateDispatches registers open dispatches from Agent(subagent_type=
// "spacedock:ensign") tool_use entries and closes them on any of the three mode
// anchors — the port of upstream _update_dispatch_budgets's open/close logic
// (the soft/hard budget + cooperative-shutdown machinery is out of scope here;
// the per-step quiet budget is the only timeout discipline).
func (w *streamWatcher) updateDispatches(entries []streamEntry) {
	for _, e := range entries {
		// Open: Agent(subagent_type="spacedock:ensign") assistant tool_use.
		if b := e.toolUseBlock(); b != nil && b.Name == "Agent" && b.Input.SubagentType == "spacedock:ensign" {
			id := b.ID
			if id == "" {
				id = fmt.Sprintf("anon-%d", len(w.openDispatches))
			}
			name := b.Input.Description
			if name == "" {
				name = "unnamed-ensign"
			}
			w.openDispatches[id] = &openDispatch{dispatchID: id, ensignName: name}
			continue
		}

		// Teams-mode-with-TTY close: system task_notification status=completed.
		if e.Type == "system" && e.Subtype == "task_notification" && e.Status == "completed" {
			if _, ok := w.openDispatches[e.ToolUseID]; ok {
				w.closeDispatch(e.ToolUseID)
			}
			continue
		}

		// User tool_result close anchors: bare-mode Done payload (keyed by the
		// Agent tool_use_id) or the headless inbox-poll Done: Bash result.
		if e.Type == "user" {
			for _, rb := range e.resultBlocks() {
				text := rb.flatText()
				// Bare-mode: the Agent tool_use_id's tool_result carrying a real
				// Done payload (NOT just a "Spawned successfully" spawn-ack).
				if _, ok := w.openDispatches[rb.ToolUseID]; ok && looksLikeBareDone(text) {
					w.closeDispatch(rb.ToolUseID)
					continue
				}
				// Headless inbox-poll: a Bash result whose body carries
				// `from: spacedock-ensign-…-{stage}` + `text: Done:`. Closes the
				// open dispatch whose description shares the stage substring.
				if sender := parseInboxDoneSender(text); sender != "" {
					if id := w.findOpenDispatchForSender(sender); id != "" {
						w.closeDispatch(id)
					}
				}
			}
		}
	}
}

func (w *streamWatcher) closeDispatch(id string) {
	if _, ok := w.openDispatches[id]; ok {
		delete(w.openDispatches, id)
		w.closedCount++
	}
}

// findOpenDispatchForSender returns the open dispatch whose ensignName shares a
// stage substring with the inbox sender (e.g. sender
// "spacedock-ensign-make-it-work-done" matches a dispatch described
// "…: done"). The port of upstream _find_open_dispatch_for_sender.
func (w *streamWatcher) findOpenDispatchForSender(sender string) string {
	s := strings.ToLower(sender)
	stages := []string{"implementation", "validation", "analysis", "design", "backlog", "ideation", "done", "work"}
	for id, d := range w.openDispatches {
		name := strings.ToLower(d.ensignName)
		for _, stage := range stages {
			if strings.Contains(s, stage) && strings.Contains(name, stage) {
				return id
			}
		}
	}
	return ""
}

// safePredicate runs the caller's predicate, treating a panic as no-match (the
// port of upstream's try/except around predicate(entry)).
func safePredicate(predicate func(streamEntry) bool, e streamEntry) (matched bool) {
	defer func() { _ = recover() }()
	return predicate(e)
}

// looksLikeBareDone reports whether an Agent tool_result text is a bare-mode
// completion payload rather than a teams-mode spawn ack. The port of upstream
// _looks_like_bare_done: a teams-mode Agent() tool_result is just a
// "Spawned successfully" ack; a bare-mode one is the worker's Done payload.
func looksLikeBareDone(text string) bool {
	return !strings.Contains(text, "Spawned successfully")
}

// parseInboxDoneSender returns the `from:` sender when text is an inbox-poll
// Done: entry (carries both `from:` and a `text: Done:` line). The port of
// upstream _parse_inbox_done_sender — the combined signature so plain agent
// names or unrelated Bash output don't accidentally close a dispatch.
func parseInboxDoneSender(text string) string {
	if !strings.Contains(text, "Done:") || !strings.Contains(text, "from:") {
		return ""
	}
	var sender string
	sawDone := false
	for _, line := range strings.Split(text, "\n") {
		s := strings.TrimSpace(line)
		if strings.HasPrefix(s, "from:") {
			sender = strings.TrimSpace(strings.TrimPrefix(s, "from:"))
		} else if strings.HasPrefix(s, "text:") && strings.Contains(s, "Done:") {
			sawDone = true
			break
		}
	}
	if sender != "" && sawDone {
		return sender
	}
	return ""
}

// --- parsed stream-json entry shapes -------------------------------------

// streamEntry is the parsed shape of one stream-json JSONL line — only the
// fields the watcher reads. claude's stream-json is standard Claude Code output
// (spacedock claude syscall.Exec-forwards --output-format stream-json verbatim),
// so these shapes match what the upstream Python watcher parses.
type streamEntry struct {
	Type      string         `json:"type"`
	Subtype   string         `json:"subtype"`
	Status    string         `json:"status"`
	ToolUseID string         `json:"tool_use_id"`
	Message   *streamMessage `json:"message"`
}

type streamMessage struct {
	Model   string        `json:"model"`
	Content []streamBlock `json:"content"`
}

type streamBlock struct {
	Type       string          `json:"type"`
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Input      streamToolInput `json:"input"`
	ToolUseID  string          `json:"tool_use_id"`
	RawContent json.RawMessage `json:"content"`
}

type streamToolInput struct {
	SubagentType string          `json:"subagent_type"`
	Description  string          `json:"description"`
	Message      json.RawMessage `json:"message"`
}

// toolUseBlock returns the first tool_use block of an assistant entry, or nil —
// the port of upstream _tool_use_block.
func (e streamEntry) toolUseBlock() *streamBlock {
	if e.Type != "assistant" || e.Message == nil {
		return nil
	}
	for i := range e.Message.Content {
		if e.Message.Content[i].Type == "tool_use" {
			return &e.Message.Content[i]
		}
	}
	return nil
}

// resultBlocks returns the tool_result blocks of a user entry.
func (e streamEntry) resultBlocks() []*streamBlock {
	if e.Type != "user" || e.Message == nil {
		return nil
	}
	var out []*streamBlock
	for i := range e.Message.Content {
		if e.Message.Content[i].Type == "tool_result" {
			out = append(out, &e.Message.Content[i])
		}
	}
	return out
}

// flatText flattens a tool_result content payload (a string OR a list of text
// blocks) into one string — the port of upstream _tool_result_text.
func (b *streamBlock) flatText() string {
	if len(b.RawContent) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(b.RawContent, &s) == nil {
		return s
	}
	var blocks []struct {
		Text string `json:"text"`
	}
	if json.Unmarshal(b.RawContent, &blocks) == nil {
		parts := make([]string, 0, len(blocks))
		for _, blk := range blocks {
			if blk.Text != "" {
				parts = append(parts, blk.Text)
			}
		}
		return strings.Join(parts, "\n")
	}
	return ""
}

// --- terminal-teardown grade ---------------------------------------------

// isTeardownToolUse reports whether the entry is one of the terminal-teardown
// tool calls: a `TeamDelete`, or a `SendMessage` carrying a `shutdown_request`.
// These are the calls the bounded teardown STOPS issuing once it emits the
// terminal-status marker — so a teardown tool_use AFTER the marker means the FO
// did NOT hold (it kept retrying), which fails the grade.
func isTeardownToolUse(e streamEntry) bool {
	b := e.toolUseBlock()
	if b == nil {
		return false
	}
	if b.Name == "TeamDelete" {
		return true
	}
	if b.Name == "SendMessage" && strings.Contains(string(b.Input.Message), "shutdown_request") {
		return true
	}
	return false
}

// lineHasMarker reports whether a raw stream-json line carries the verbatim
// terminal-status marker (emitted as an assistant text block). The grade greps
// the raw line — the same fixed substring the contract mandates — rather than
// parsing text-block content, so a marker the FO emits in any text position is
// found.
func lineHasMarker(line string) bool {
	return strings.Contains(line, terminalTeardownMarker)
}

// gradeTerminalTeardown decides whether a stream carries the bounded-best-effort
// terminus: the contract-mandated terminal-status marker FOLLOWED BY A HOLD (no
// further teardown tool_use). It returns ok=true ONLY when the marker appears AND
// no `TeamDelete`/`shutdown_request` tool_use appears on any LATER line. This is
// the load-bearing discriminator:
//
//   - the fix shape (bounded attempts → marker → hold) PASSES;
//   - the pre-yy give-up (sonnet_teamdelete_hang: one TeamDelete, ends the turn,
//     NO marker) FAILS — the marker is absent;
//   - the post-yy retry-loop (sonnet_teamdelete_retryloop: 6 TeamDelete past the
//     cap, NO marker, never holds) FAILS — the marker is absent.
//
// Grading on the marker+hold (NOT on the shutdown_request/TeamDelete beats both
// bug shapes ALSO emit) is what makes the grade green ONLY on the fix: a
// beats-present grade would green the very runs yy's fix failed. The reason
// string localizes which condition failed for a CI reader.
func gradeTerminalTeardown(lines []string) (ok bool, reason string) {
	markerLine := -1
	for i, line := range lines {
		if lineHasMarker(line) {
			markerLine = i
			break
		}
	}
	if markerLine < 0 {
		return false, "terminal-status marker absent — the FO never emitted the bounded-teardown terminus (a silent give-up or an unbounded retry loop)"
	}
	// After the marker, the FO must HOLD: no further teardown tool_use.
	for i := markerLine + 1; i < len(lines); i++ {
		var e streamEntry
		if json.Unmarshal([]byte(lines[i]), &e) != nil {
			continue
		}
		if isTeardownToolUse(e) {
			return false, fmt.Sprintf("teardown tool_use at line %d AFTER the marker (line %d) — the FO did not hold; it kept retrying", i, markerLine)
		}
	}
	return true, ""
}

// --- live pipe adapters (used only by the //go:build live live test) -----

// pipeLineSource is the live lineSource: a background goroutine reads the
// subprocess's StdoutPipe via a bufio.Scanner and pushes complete lines onto a
// channel; drain returns the lines buffered since the last poll. The Go port of
// upstream _drain_entries's read-and-hold-partial-lines logic — bufio.Scanner
// holds partial trailing lines until their newline arrives, so a half-written
// JSONL line is never parsed early.
type pipeLineSource struct {
	lines chan string
}

// newPipeLineSource starts reading r in the background and returns a lineSource
// over it. Closing r (subprocess exit) ends the goroutine.
func newPipeLineSource(r io.Reader) *pipeLineSource {
	p := &pipeLineSource{lines: make(chan string, 1024)}
	go func() {
		defer close(p.lines)
		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
		for scanner.Scan() {
			p.lines <- scanner.Text()
		}
	}()
	return p
}

func (p *pipeLineSource) drain() []string {
	var out []string
	for {
		select {
		case line, ok := <-p.lines:
			if !ok {
				return out
			}
			out = append(out, line)
		default:
			return out
		}
	}
}
