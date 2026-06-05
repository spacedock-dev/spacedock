package ensigncycle

import (
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

// Offline unit tests for the per-stage stall-watchdog. They drive synthetic
// streams (no model) so the liveness guard's behavior is pinned cheaply: a
// normal-cadence stream passes through; a stream that goes silent past the budget
// is killed + fails fast. They use a SHORT test budget so the suite stays fast —
// streamWithStallWatchdog takes the budget as an argument precisely so the offline
// test can shrink it while production uses stageStallTimeout (120s, the
// captain-approved exception pinned by TestStageStallTimeoutIsCaptainApprovedException).

// TestStallWatchdogPassesNormalCadenceStream: lines arriving within the budget are
// copied through and the watchdog never fires (clean EOF, nil error).
func TestStallWatchdogPassesNormalCadenceStream(t *testing.T) {
	const budget = 100 * time.Millisecond
	pr, pw := io.Pipe()
	killed := false

	go func() {
		// Five lines, each well within the budget — a live, progressing stream.
		for i := 0; i < 5; i++ {
			io.WriteString(pw, "event line\n")
			time.Sleep(budget / 4)
		}
		pw.Close()
	}()

	out, err := streamWithStallWatchdog(pr, budget, func() { killed = true })
	if err != nil {
		t.Fatalf("normal-cadence stream must not stall: %v", err)
	}
	if killed {
		t.Fatal("normal-cadence stream must not trigger the kill callback")
	}
	if got := strings.Count(out, "event line"); got != 5 {
		t.Fatalf("watchdog must copy all 5 lines through, got %d", got)
	}
}

// TestStallWatchdogKillsStalledStream: after one line the stream goes silent past
// the budget; the watchdog must invoke the kill callback and return a stallError.
func TestStallWatchdogKillsStalledStream(t *testing.T) {
	const budget = 80 * time.Millisecond
	pr, pw := io.Pipe()
	var once sync.Once
	killedCh := make(chan struct{})
	onStall := func() {
		once.Do(func() {
			close(killedCh)
			// A real runner kills the process here, which closes the pipe; mirror
			// that so the reader goroutine drains.
			pw.Close()
		})
	}

	go func() {
		io.WriteString(pw, "first event\n")
		// Then go silent far longer than the budget — a hung stage.
		time.Sleep(budget * 10)
		io.WriteString(pw, "too late\n")
	}()

	out, err := streamWithStallWatchdog(pr, budget, onStall)
	if err == nil {
		t.Fatal("a stalled stream must return a stall error")
	}
	if _, ok := err.(*stallError); !ok {
		t.Fatalf("expected *stallError, got %T: %v", err, err)
	}
	select {
	case <-killedCh:
	default:
		t.Fatal("a stalled stream must invoke the kill callback")
	}
	if !strings.Contains(out, "first event") {
		t.Fatal("watchdog must return the partial output captured before the stall")
	}
	if strings.Contains(err.Error(), "stalled") == false {
		t.Fatalf("stall error must name the stall: %q", err.Error())
	}
}
