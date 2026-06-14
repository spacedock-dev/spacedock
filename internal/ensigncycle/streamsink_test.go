// ABOUTME: Pins the shared-scenario runners' stream sink to a silent discard —
// ABOUTME: the watcher still drains/records every line while nothing tees to stdout.
package ensigncycle

import (
	"strings"
	"sync"
	"testing"
	"time"
)

// streamLineSink is the destination a stream tee forwards each parsed line to. The
// shared-scenario runners' discardStreamLine deliberately does NOT forward here —
// the full stream is persisted to the per-scenario artifact (claude-stream.jsonl /
// the Codex equivalent) and every failure message already carries the bounded
// transcriptTail()/tail(stream,…), so teeing every parsed line is redundant and
// bloats CI output (a single failed step dumped ~143KB of jsonl, burying the
// actual failure line). Reverting discardStreamLine to forward each line here is
// exactly the re-flood regression TestStreamSinkDiscardsLines guards against:
// tests install a counter as the sink, and the discard run must leave it at zero.
var streamLineSink = func(string) {}

// discardStreamLine is the stream sink the shared-scenario live runners hand
// newStreamWatcher. It drops the line — it does NOT forward to streamLineSink — so
// the watcher still drains and records the full transcript while nothing is teed
// out. (See streamLineSink for why the artifact + failure-tail make the tee
// redundant.)
func discardStreamLine(string) {}

// emitCounter counts every line a sink forwards to it — a stand-in for stdout.
// Mutex-guarded because the watcher tees from its poll-loop goroutine.
type emitCounter struct {
	mu    sync.Mutex
	count int
}

func (s *emitCounter) tee(string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count++
}

func (s *emitCounter) seen() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.count
}

// drainedThenExitProc reports exited once its dripping source has replayed every
// line — modelling the live subprocess that closes its stream pipe and exits
// after the final stream entry, so drainToExit drains the whole fixture then
// returns cleanly.
type drainedThenExitProc struct {
	src *drippingLineSource
}

func (p *drainedThenExitProc) poll() (int, bool) {
	p.src.mu.Lock()
	defer p.src.mu.Unlock()
	return 0, p.src.next >= len(p.src.lines)
}

func (p *drainedThenExitProc) kill() {}

// drainFixture runs the captured 341-line stream through newStreamWatcher with
// the given tee, draining to clean exit, and returns the count of lines the
// watcher recorded into its transcript. The watcher appends every drained line to
// transcript INDEPENDENT of the tee, so this count proves the parse/record path
// regardless of what the tee does.
func drainFixture(t *testing.T, tee func(string)) (recorded int) {
	t.Helper()
	lines := loadStreamFixture(t, "sonnet_teamdelete_hang.stream.jsonl")
	if len(lines) != 341 {
		t.Fatalf("fixture line count = %d, want 341", len(lines))
	}
	src := &drippingLineSource{lines: lines}
	proc := &drainedThenExitProc{src: src}
	w := newStreamWatcher(src, proc, tee)
	w.quietBudget = 2 * time.Second
	w.exitBudget = 2 * time.Second
	w.pollInterval = time.Millisecond
	if _, err := w.drainToExit(w.quietBudget, "stream sink discard"); err != nil {
		t.Fatalf("drainToExit over the captured fixture must reach exit: %v", err)
	}
	return len(strings.Split(w.fullTranscript(), "\n"))
}

// TestStreamSinkDiscardsLines drives newStreamWatcher with the SAME callback the
// shared-scenario runners pass (discardStreamLine) over the real 341-line captured
// fixture and asserts that the sink wired BEHIND discardStreamLine (streamLineSink)
// receives ZERO lines, while the watcher's transcript still holds all 341 lines
// (drain/parse/record path intact).
//
// The recorder is installed AS streamLineSink — the destination a forwarding tee
// would route to. The production discardStreamLine does not forward there, so the
// recorder stays 0. Reverting discardStreamLine's body to forward each line to
// streamLineSink (the re-flood regression) makes the recorder non-zero and reds
// the first assertion — the control run below proves that path drives it to 341,
// so the zero is the discard, not a dead counter. (AC-1)
func TestStreamSinkDiscardsLines(t *testing.T) {
	// Install a recorder behind the tee. Restore the production no-op sink after, so
	// the swap does not leak into sibling tests.
	rec := &emitCounter{}
	saved := streamLineSink
	streamLineSink = rec.tee
	defer func() { streamLineSink = saved }()

	// Discard run: the watcher's tee is discardStreamLine, the exact symbol both
	// shared-scenario runners hand newStreamWatcher. It does not forward to
	// streamLineSink, so the recorder behind the sink must stay at zero.
	recorded := drainFixture(t, discardStreamLine)
	if got := rec.seen(); got != 0 {
		t.Errorf("discardStreamLine must forward 0 lines to streamLineSink, got %d — the per-line tee re-flood is back", got)
	}
	if recorded != 341 {
		t.Errorf("with discardStreamLine as the tee, the watcher transcript must still record all 341 fixture lines, got %d", recorded)
	}

	// Control: a forwarding tee that DOES route each line to streamLineSink — the
	// old per-line t.Log shape the runners moved off of. The same recorder reaches
	// every fixture line, proving the recorder fires when the tee forwards. This is
	// the wiring the discard assertion guards against regressing back to.
	rec2 := &emitCounter{}
	streamLineSink = rec2.tee
	drainFixture(t, func(line string) { streamLineSink(line) })
	if rec2.seen() != 341 {
		t.Errorf("control: a forwarding tee must route all 341 lines to streamLineSink, got %d", rec2.seen())
	}
}
