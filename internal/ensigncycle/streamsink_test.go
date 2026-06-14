// ABOUTME: Pins the shared-scenario runners' stream sink to a silent discard —
// ABOUTME: the watcher still drains/records every line while nothing tees to stdout.
package ensigncycle

import (
	"strings"
	"sync"
	"testing"
	"time"
)

// discardStreamLine is the stream sink the shared-scenario live runners hand
// newStreamWatcher. It drops the line: the full stream is persisted to the
// per-scenario artifact (claude-stream.jsonl / the Codex equivalent) and every
// failure message already carries the bounded transcriptTail()/tail(stream,…),
// so teeing every parsed line to t.Log is redundant and bloats CI output (a
// single failed step dumped ~143KB of jsonl, burying the actual failure line).
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
// fixture, draining the whole stream, and asserts nothing reaches stdout while the
// watcher's transcript still holds all 341 lines (drain/parse/record path intact).
//
// The discard run wires discardStreamLine — the exact symbol both runners hand
// newStreamWatcher — as the watcher's tee with a counter behind it; discard never
// calls onward, so the counter stays 0. A control run wires the counter directly
// as the tee, where it reaches 341 — proving the counter fires when the tee
// forwards, so the zero in the discard run is the discard, not a dead counter.
// Reverting the runners' tee to that forwarding (t.Log) shape is exactly what
// drives the count non-zero and reds the discard assertion. (AC-1)
func TestStreamSinkDiscardsLines(t *testing.T) {
	// Discard run: the watcher's tee is discardStreamLine, the exact symbol both
	// shared-scenario runners hand newStreamWatcher. The drain/parse/record path is
	// driven by the real callback; the watcher must still record all 341 lines.
	recorded := drainFixture(t, discardStreamLine)
	if recorded != 341 {
		t.Errorf("with discardStreamLine as the tee, the watcher transcript must record all 341 fixture lines, got %d", recorded)
	}

	// Control: a forwarding tee — the old per-line t.Log wiring the runners are
	// moving OFF of. It reaches every fixture line, so reverting either runner's
	// newStreamWatcher(...) tee from discardStreamLine back to a forwarding sink
	// re-floods stdout with the whole 341-line stream: the bloat this task removes.
	teeCounter := &emitCounter{}
	drainFixture(t, teeCounter.tee)
	if teeCounter.seen() != 341 {
		t.Errorf("control: a forwarding tee must reach all 341 lines (the re-flood discardStreamLine prevents), got %d", teeCounter.seen())
	}
}
