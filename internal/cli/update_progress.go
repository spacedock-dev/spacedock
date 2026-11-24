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

func updateSeconds(d time.Duration) string {
	ms := d.Milliseconds()
	return fmt.Sprintf("%d.%03ds", ms/1000, ms%1000)
}

type updateLineWriter struct {
	out     io.Writer
	zero    time.Time
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
			if _, err := fmt.Fprintf(w.out, "[+%s] ", updateSeconds(time.Since(w.zero))); err != nil {
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
	fmt.Fprintf(p.trace.out, "[+%s] %s\n", updateSeconds(now.Sub(p.trace.zero)), line)
}

func (p *updateProgress) step(label string) func() {
	if p.trace == nil {
		fmt.Fprintln(p.plain, label+"...")
		return func() {}
	}

	// Each open step tracks its own last child boundary.
	now := time.Now()
	previous := p.last[len(p.last)-1]
	if previous.IsZero() {
		now, previous = p.trace.zero, p.trace.zero
	}
	scope := "gap"
	if len(p.last) > 1 {
		scope = "parent work"
	}
	p.at(now, fmt.Sprintf("%s... (%s %s)", label, scope, updateSeconds(now.Sub(previous))))
	p.last = append(p.last, now)
	return func() {
		end := time.Now()
		p.last = p.last[:len(p.last)-1]
		p.last[len(p.last)-1] = end
		p.at(end, fmt.Sprintf("Elapsed: %s: %s", label, updateSeconds(end.Sub(now))))
	}
}

func (p *updateProgress) finish() {
	now := time.Now()
	p.at(now, fmt.Sprintf("Total: %s (gap %s)", updateSeconds(now.Sub(p.trace.zero)), updateSeconds(now.Sub(p.last[0]))))
}
