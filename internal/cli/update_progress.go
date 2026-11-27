package cli

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

type verboseInventory struct {
	hostOps
	progress *updateProgress
}

func (ops verboseInventory) PluginInventory(host string) ([]pluginInventoryEntry, error) {
	done := ops.progress.step("Listing " + hostTitle(host) + " plugins")
	defer done()
	return ops.hostOps.PluginInventory(host)
}

type updateLineWriter struct {
	out     io.Writer
	clock   elapsedTrace
	midLine bool
	peer    *updateLineWriter
}

func (w *updateLineWriter) Write(data []byte) (int, error) {
	if len(data) > 0 {
		// Either stream may be redirected to the same destination.
		w.peer.boundary()
	}
	written := 0
	for len(data) > 0 {
		if !w.midLine {
			if _, err := fmt.Fprintf(w.out, "%s ", w.clock.prefix(w.clock.now())); err != nil {
				return written, err
			}
			w.midLine = true
		}
		size := len(data)
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			size = i + 1
		}
		n, err := w.out.Write(data[:size])
		written += n
		if n > 0 && data[n-1] == '\n' {
			w.midLine = false
		}
		if err != nil {
			return written, err
		}
		if n != size {
			return written, io.ErrShortWrite
		}
		data = data[n:]
	}
	return written, nil
}

func (w *updateLineWriter) boundary() {
	if w.midLine {
		fmt.Fprintln(w.out)
		w.midLine = false
	}
}

type updateProgress struct {
	plain io.Writer
	trace *updateLineWriter
	last  []time.Time
}

func (p *updateProgress) at(now time.Time, line string) {
	p.trace.peer.boundary()
	p.trace.boundary()
	p.trace.clock.at(now, line)
}

func (p *updateProgress) step(label string) func() {
	if p.trace == nil {
		fmt.Fprintln(p.plain, label+"...")
		return func() {}
	}

	// Each open step tracks its own last child boundary.
	now := p.trace.clock.now()
	previous := p.last[len(p.last)-1]
	if previous.IsZero() {
		now, previous = p.trace.clock.zero, p.trace.clock.zero
	}
	scope := "gap"
	if len(p.last) > 1 {
		scope = "parent work"
	}
	p.at(now, fmt.Sprintf("%s... (%s %s)", label, scope, durationSeconds(now.Sub(previous))))
	p.last = append(p.last, now)
	return func() {
		end := p.trace.clock.now()
		p.last = p.last[:len(p.last)-1]
		p.last[len(p.last)-1] = end
		p.at(end, fmt.Sprintf("Elapsed: %s: %s", label, durationSeconds(end.Sub(now))))
	}
}

func (p *updateProgress) finish() {
	now := p.trace.clock.now()
	p.at(now, fmt.Sprintf("Total: %s (gap %s)", durationSeconds(now.Sub(p.trace.clock.zero)), durationSeconds(now.Sub(p.last[0]))))
}
