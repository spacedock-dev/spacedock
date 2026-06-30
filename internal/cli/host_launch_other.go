// ABOUTME: Non-unix compile shim for the resident launcher — the portable
// ABOUTME: spawn-wait core builds, with no signal model and a plain exit code.
//go:build !unix

package cli

import (
	"os"
	"os/exec"
)

// forwardHostSignals is a no-op off unix: the resident-launcher signal model
// (terminal-generated delivery through a shared foreground process group) is a
// unix contract. `spacedock <host>` is not a supported launch path on Windows
// (syscall.Exec returned EWINDOWS even before this change); only compile-safety
// is in scope.
func forwardHostSignals(cmd *exec.Cmd) func() {
	return func() {}
}

// hostExitCode returns the portable exit code for a finished host process. A nil
// waitErr is a clean exit 0; otherwise ProcessState.ExitCode() carries the host's
// status (-1 when unavailable), with no unix signal-death refinement.
func hostExitCode(ps *os.ProcessState, waitErr error) int {
	if waitErr == nil {
		return 0
	}
	if ps != nil {
		return ps.ExitCode()
	}
	return 1
}
