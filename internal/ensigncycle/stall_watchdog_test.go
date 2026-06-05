package ensigncycle

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

// The per-stage stall-watchdog for the live host launches. The captain ruling
// bans long basket/binary timeouts: a live scenario may take minutes of GENUINE
// sequential model work (a 4-stage rejection-flow shares one `spacedock claude -p`
// process because reviewer-reuse needs the reviewer alive across the route-back),
// but every UNIT of progress — each FO turn, each ensign stage — completes well
// under 60s, and the host stream emits frequent liveness events (Claude
// thinking_tokens, Codex item.* events). So the rule-compliant guard is not a
// whole-run deadline but a STALL detector: reset a timer on every received stream
// line, and kill the process only when NO line arrives for stallTimeout. That
// catches a true hang (the opus-panic signature was a stalled stage) precisely,
// while never false-killing slow-but-live thinking.

// stageStallTimeout is the per-stage liveness budget: if the host stream emits no
// line for this long the stage is treated as hung and the process is killed.
const stageStallTimeout = 60 * time.Second

// streamWithStallWatchdog copies r line-by-line into a returned buffer, resetting
// a stall timer on every line. If no line arrives for stallTimeout it invokes
// onStall (the caller's process-kill) and returns the bytes read so far plus a
// non-nil stall error. On clean EOF it returns the full output and a nil error.
// It is host-neutral: the Claude and Codex runners both feed it their process
// stdout pipe and a kill closure, so the 60s liveness guard is identical across
// hosts and exercised offline without a model.
func streamWithStallWatchdog(r io.Reader, stallTimeout time.Duration, onStall func()) (string, error) {
	lines := make(chan string)
	readErr := make(chan error, 1)

	go func() {
		scanner := bufio.NewScanner(r)
		// Stream-json lines can be large (a full transcript event); raise the
		// scanner's token cap so a long line is delivered, not silently dropped as
		// a fake stall.
		scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		readErr <- scanner.Err()
		close(lines)
	}()

	var buf []byte
	timer := time.NewTimer(stallTimeout)
	defer timer.Stop()
	for {
		select {
		case line, ok := <-lines:
			if !ok {
				// Reader drained. Surface any scan error; otherwise clean EOF.
				if err := <-readErr; err != nil {
					return string(buf), err
				}
				return string(buf), nil
			}
			buf = append(buf, line...)
			buf = append(buf, '\n')
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(stallTimeout)
		case <-timer.C:
			onStall()
			// Drain the reader goroutine so it does not leak once the process is
			// killed (the pipe closes, the scanner ends).
			go func() {
				for range lines {
				}
			}()
			return string(buf), &stallError{stallTimeout: stallTimeout, partial: string(buf)}
		}
	}
}

// stallError is returned when the host stream goes silent past the stall budget —
// a hung stage, distinct from a clean failure. The runner surfaces it as a
// fail-fast "stage stalled" so a true hang is killed in stallTimeout, not after a
// multi-minute basket.
type stallError struct {
	stallTimeout time.Duration
	partial      string
}

func (e *stallError) Error() string {
	return fmt.Sprintf("stage stalled >%s with no stream progress (hang); killed the host process", e.stallTimeout)
}
