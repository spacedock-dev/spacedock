// ABOUTME: Unix resident-launcher signal model + exit-code mapping — absorb the
// ABOUTME: terminal signals, forward the externally-targeted ones, 128+signum.
//go:build unix

package cli

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// forwardHostSignals arms the launcher's signal disposition for the
// shared-foreground-group model and returns a forward function and a stop
// function the caller defers. Arming (signal.Notify) runs here, before the
// caller's cmd.Start(), so a signal delivered during the fork/exec window queues
// on the buffered channel rather than hitting the launcher's default disposition.
// The caller invokes forward AFTER cmd.Start() returns: the go-statement that
// spawns the pump then carries a happens-before edge from Start's write of
// cmd.Process to the goroutine's read of it, so the read is synchronized (and any
// queued early signal drains from the buffered channel once the pump runs).
//
// The host shares the launcher's process group (no Setpgid), which the shell
// already made the terminal's foreground group, so the kernel delivers
// terminal-generated signals to the host directly. The launcher therefore:
//   - ABSORBS SIGINT/SIGQUIT (via a Go handler, not SIG_IGN — SIG_IGN is
//     inherited across exec and would silently un-interruptible a default-
//     disposition host). The host already got its own copy from the foreground
//     group, so re-forwarding would double-deliver; the launcher just stays
//     resident. (SIGINT-absorb is spike-validated; SIGQUIT-absorb is by analogy
//     — same terminal-generated, shared-group reasoning.)
//   - FORWARDS SIGTERM/SIGHUP, which reach only the launcher when sent to its
//     pid (a `kill` or a zellij pane close), so relaying them tears the pair
//     down together. (SIGTERM-forward is spike-validated; SIGHUP-forward is by
//     analogy — same pid-targeted reasoning. A controlling-terminal hangup that
//     hits the whole foreground group makes the host receive SIGHUP twice, which
//     is a harmless idempotent teardown.)
//
// SIGTSTP/SIGCONT/SIGWINCH are left at their default disposition: Ctrl-Z then
// suspends the whole foreground group and the shell's fg/bg resumes it (correct
// shared-group job control), and SIGWINCH's default is ignore so the launcher
// survives while the host gets its own copy from the foreground group.
func forwardHostSignals(cmd *exec.Cmd) (forward func(), stop func()) {
	ch := make(chan os.Signal, 16)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)
	forward = func() {
		go func() {
			for sig := range ch {
				switch sig {
				case syscall.SIGTERM, syscall.SIGHUP:
					if cmd.Process != nil {
						_ = cmd.Process.Signal(sig)
					}
				default:
					// SIGINT/SIGQUIT: the kernel already delivered this to the host
					// via the shared foreground group. Absorb it here so the launcher
					// stays resident; do NOT re-forward.
				}
			}
		}()
	}
	stop = func() {
		signal.Stop(ch)
		close(ch)
	}
	return forward, stop
}

// hostExitCode maps a finished host process to the exit code the launcher
// propagates. A nil waitErr is a clean exit 0. A non-zero exit yields the host's
// status; a signal death refines to the shell's 128+signum convention via the
// WaitStatus. Any other Wait anomaly collapses to 1.
func hostExitCode(ps *os.ProcessState, waitErr error) int {
	if waitErr == nil {
		return 0
	}
	if ee, ok := waitErr.(*exec.ExitError); ok {
		if ws, ok := ee.Sys().(syscall.WaitStatus); ok {
			if ws.Signaled() {
				return 128 + int(ws.Signal())
			}
			return ws.ExitStatus()
		}
		if ps != nil {
			return ps.ExitCode()
		}
	}
	return 1
}
