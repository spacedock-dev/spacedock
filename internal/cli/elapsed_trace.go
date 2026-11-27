package cli

import (
	"fmt"
	"io"
	"time"
)

type elapsedTrace struct {
	out  io.Writer
	now  func() time.Time
	zero time.Time
}

func newElapsedTrace(out io.Writer, now func() time.Time) elapsedTrace {
	return elapsedTrace{out: out, now: now, zero: now()}
}

func (t elapsedTrace) prefix(now time.Time) string {
	return "[+" + durationSeconds(now.Sub(t.zero)) + "]"
}

func (t elapsedTrace) at(now time.Time, line string) {
	fmt.Fprintf(t.out, "%s %s\n", t.prefix(now), line)
}

func durationSeconds(d time.Duration) string {
	ms := d.Milliseconds()
	return fmt.Sprintf("%d.%03ds", ms/1000, ms%1000)
}
