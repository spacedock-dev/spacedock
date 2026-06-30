// ABOUTME: AC-3 resident-launcher PTY test — under a real controlling terminal,
// ABOUTME: Ctrl-C/resize reach the host while the launcher survives and propagates.
//go:build unix

package cli

import (
	"io"
	"os"
	"os/exec"
	"syscall"
	"testing"

	"github.com/creack/pty"
)

// startLauncherOnPTY starts the launcher role on a fresh controlling terminal
// sized 24x80 and drains the master so the slave never blocks. The launcher
// becomes the session leader (pty.Start sets Setsid/Setctty), so its process
// group is the terminal's foreground group; the stub it spawns (no Setpgid)
// shares that group. The caller waits on cmd and reads the stub's logfile.
func startLauncherOnPTY(t *testing.T, log, mode string) (*os.File, *exec.Cmd) {
	t.Helper()
	cmd := roleCmd(t, "launcher", log, mode)
	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 24, Cols: 80})
	if err != nil {
		t.Fatalf("pty.StartWithSize: %v", err)
	}
	t.Cleanup(func() { _ = ptmx.Close() })
	t.Cleanup(func() { _ = cmd.Process.Kill() })
	go func() { _, _ = io.Copy(io.Discard, ptmx) }()
	return ptmx, cmd
}

// TestLaunchPTYTerminalSignals is AC-3: under a real controlling terminal the
// host's stdin is a TTY at the real window size, a terminal Ctrl-C reaches the
// host while the launcher survives to propagate, a resize delivers SIGWINCH, and
// an external SIGTERM to the launcher pid is forwarded — all observed via the
// stub's logfile + the launcher's exit code.
func TestLaunchPTYTerminalSignals(t *testing.T) {
	t.Run("host stdin is a real TTY at the real window size, and Ctrl-C reaches it", func(t *testing.T) {
		log := newStubLog(t)
		ptmx, cmd := startLauncherOnPTY(t, log, "wait")
		if !waitForLog(t, log, "STARTED") {
			t.Fatalf("stub never started; log=%q", readLog(t, log))
		}
		rec := readStubRecord(t, log)
		if !rec.isatty {
			t.Fatalf("host stdin is not a TTY (isatty=false); STARTED=%q", readLog(t, log))
		}
		if rec.rows != 24 || rec.cols != 80 {
			t.Fatalf("host winsize=%dx%d, want 24x80 (host must see the real terminal size)", rec.rows, rec.cols)
		}
		// Ctrl-C through the line discipline → SIGINT to the foreground group.
		if _, err := ptmx.Write([]byte{0x03}); err != nil {
			t.Fatalf("write Ctrl-C: %v", err)
		}
		if !waitForLog(t, log, "SIGNAL SIGINT") {
			t.Fatalf("host never saw SIGINT after Ctrl-C; log=%q", readLog(t, log))
		}
		_ = cmd.Wait()
		if c := exitCodeOf(cmd); c != 130 {
			t.Fatalf("launcher exit=%d, want 130 (it must survive Ctrl-C and propagate the host's exit)", c)
		}
	})

	t.Run("terminal resize delivers SIGWINCH to the host", func(t *testing.T) {
		log := newStubLog(t)
		ptmx, cmd := startLauncherOnPTY(t, log, "wait")
		if !waitForLog(t, log, "STARTED") {
			t.Fatalf("stub never started; log=%q", readLog(t, log))
		}
		if err := pty.Setsize(ptmx, &pty.Winsize{Rows: 40, Cols: 120}); err != nil {
			t.Fatalf("pty.Setsize: %v", err)
		}
		// Match the resize-triggered SIGWINCH specifically (its new dimensions),
		// via a platform-stable marker — syscall.Signal.String() renders SIGWINCH
		// differently on Linux ("window changed") vs darwin ("window size changes").
		if !waitForLog(t, log, "SIGNAL SIGWINCH winsize=40x120") {
			t.Fatalf("host never saw the resize SIGWINCH (40x120); log=%q", readLog(t, log))
		}
		_ = cmd.Process.Signal(syscall.SIGTERM)
		_ = cmd.Wait()
	})

	t.Run("external SIGTERM to the launcher pid is forwarded to the host", func(t *testing.T) {
		log := newStubLog(t)
		_, cmd := startLauncherOnPTY(t, log, "wait")
		if !waitForLog(t, log, "STARTED") {
			t.Fatalf("stub never started; log=%q", readLog(t, log))
		}
		// Pid-targeted: only the launcher receives it; it must forward to the host.
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			t.Fatalf("signal launcher: %v", err)
		}
		if !waitForLog(t, log, "SIGNAL SIGTERM") {
			t.Fatalf("host never saw the forwarded SIGTERM; log=%q", readLog(t, log))
		}
		_ = cmd.Wait()
		if c := exitCodeOf(cmd); c != 143 {
			t.Fatalf("launcher exit=%d, want 143 (forwarded SIGTERM tears the pair down)", c)
		}
	})
}
