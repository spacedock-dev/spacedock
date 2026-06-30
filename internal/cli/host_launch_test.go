// ABOUTME: AC-1/AC-2 resident-launcher behavior tests — a re-exec helper-process
// ABOUTME: harness proves host PPID parentage, the exec baseline, and exit codes.
//go:build unix

package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

// launchTestRoleEnv selects which helper-process role the re-exec'd test binary
// plays. The default (unset) runs the normal test suite; the launcher/execlauncher
// roles spawn the stub through the real launch path / the exec baseline, and the
// stub role reports its process identity + TTY state to a logfile.
const launchTestRoleEnv = "SPACEDOCK_LAUNCH_TEST_ROLE"

// TestMain dispatches the helper-process roles before the test suite runs. Each
// role reads its logfile path and mode from os.Args[1]/os.Args[2] and os.Exits;
// only the default arm runs the tests. This is the os/exec helper-process pattern
// — the test binary re-execs itself instead of compiling separate fixtures, so
// execHost.Launch runs in a genuine child process under the test's control.
func TestMain(m *testing.M) {
	switch os.Getenv(launchTestRoleEnv) {
	case "stub":
		stubHostRole()
	case "launcher":
		launcherRole()
	case "execlauncher":
		execLauncherRole()
	default:
		os.Exit(m.Run())
	}
}

// launcherRole spawns the stub through the production resident-parent path
// (execHost.Launch) and propagates its exit code — the SUT under AC-1/AC-2.
func launcherRole() {
	self, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "launcher role: executable:", err)
		os.Exit(1)
	}
	argv := []string{self, os.Args[1], os.Args[2]}
	env := append(withoutEnv(os.Environ(), launchTestRoleEnv), launchTestRoleEnv+"=stub")
	code, err := execHost{}.Launch(argv, env)
	if err != nil {
		fmt.Fprintln(os.Stderr, "launcher role: launch:", err)
		os.Exit(1)
	}
	os.Exit(code)
}

// execLauncherRole is the syscall.Exec baseline (today's pre-change model): it
// replaces its own image with the stub, so no resident parent remains. AC-1 reads
// the contrast — under exec the stub's pid IS the launcher's pid.
func execLauncherRole() {
	self, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "execlauncher role: executable:", err)
		os.Exit(1)
	}
	argv := []string{self, os.Args[1], os.Args[2]}
	env := append(withoutEnv(os.Environ(), launchTestRoleEnv), launchTestRoleEnv+"=stub")
	if err := syscall.Exec(self, argv, env); err != nil {
		fmt.Fprintln(os.Stderr, "execlauncher role: exec:", err)
		os.Exit(1)
	}
}

// winsize mirrors struct winsize (TIOCGWINSZ): rows, cols, xpix, ypix.
type winsize struct{ Row, Col, Xpix, Ypix uint16 }

// stubIsattyAndSize runs TIOCGWINSZ on fd. Success is a genuine isatty test (a
// pipe/file errors out); the size proves SIGWINCH-relevant terminal state is
// visible to the child.
func stubIsattyAndSize(fd uintptr) (bool, winsize) {
	var ws winsize
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws)))
	return errno == 0, ws
}

// stubHostRole is the stub "host": it logs its process identity + TTY state, then
// either exits immediately with a chosen code, dies by a named signal, or waits
// in a signal loop (recording every signal it receives). Args: <logPath> <mode>.
func stubHostRole() {
	logPath := os.Args[1]
	mode := os.Args[2]
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "stub role: open log:", err)
		os.Exit(99)
	}
	defer lf.Close()
	logf := func(format string, a ...any) {
		fmt.Fprintf(lf, "stub "+format+"\n", a...)
		lf.Sync()
	}

	atty, ws := stubIsattyAndSize(os.Stdin.Fd())
	started := fmt.Sprintf("STARTED pid=%d ppid=%d pgid=%d isatty=%v winsize=%dx%d",
		os.Getpid(), os.Getppid(), syscall.Getpgrp(), atty, ws.Row, ws.Col)

	if sigName, ok := strings.CutPrefix(mode, "killself:"); ok {
		// Genuine signal death (no handler installed): the launcher's Wait sees
		// WIFSIGNALED and must render 128+signum. Stronger than an explicit
		// os.Exit(128+n) — it exercises the launcher's signal-status branch.
		logf("%s", started)
		sig := signalByName(sigName)
		logf("KILLSELF %v", sig)
		_ = syscall.Kill(os.Getpid(), sig)
		time.Sleep(2 * time.Second)
		logf("EXIT 125 (signal did not land)")
		os.Exit(125)
	}

	if mode != "wait" {
		logf("%s", started)
		code, _ := strconv.Atoi(mode)
		logf("EXIT %d (immediate)", code)
		os.Exit(code)
	}

	// Register the signal handlers BEFORE announcing STARTED, so a terminal signal
	// the test injects the instant it observes STARTED cannot land before the stub
	// is listening — which would hit the default disposition (SIGINT killing the
	// stub pre-log, SIGWINCH silently ignored) and lose the evidence.
	ch := make(chan os.Signal, 8)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGWINCH, syscall.SIGHUP)
	logf("%s", started)
	deadline := time.After(10 * time.Second)
	for {
		select {
		case sig := <-ch:
			logf("SIGNAL %v", sig)
			switch sig {
			case syscall.SIGINT:
				logf("EXIT 130 (saw SIGINT)")
				os.Exit(130)
			case syscall.SIGTERM:
				logf("EXIT 143 (saw SIGTERM)")
				os.Exit(143)
			case syscall.SIGHUP:
				logf("EXIT 129 (saw SIGHUP)")
				os.Exit(129)
				// SIGWINCH: keep running, just record it.
			}
		case <-deadline:
			logf("EXIT 124 (timeout, no signal)")
			os.Exit(124)
		}
	}
}

func signalByName(name string) syscall.Signal {
	switch name {
	case "SIGINT":
		return syscall.SIGINT
	case "SIGTERM":
		return syscall.SIGTERM
	case "SIGHUP":
		return syscall.SIGHUP
	case "SIGQUIT":
		return syscall.SIGQUIT
	case "SIGKILL":
		return syscall.SIGKILL
	default:
		return syscall.SIGTERM
	}
}

// stubRecord is the parsed STARTED line the stub writes to its logfile.
type stubRecord struct {
	pid, ppid  int
	isatty     bool
	rows, cols int
}

// newStubLog returns a fresh temp logfile path the stub appends to.
func newStubLog(t *testing.T) string {
	t.Helper()
	return t.TempDir() + "/stub.log"
}

// roleCmd builds (but does not start) an *exec.Cmd that re-execs the test binary
// in the given helper role with the stub's log path + mode as its args.
func roleCmd(t *testing.T, role, logPath, mode string) *exec.Cmd {
	t.Helper()
	self, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable: %v", err)
	}
	cmd := exec.Command(self, logPath, mode)
	cmd.Env = append(withoutEnv(os.Environ(), launchTestRoleEnv), launchTestRoleEnv+"="+role)
	return cmd
}

// readLog returns the current logfile contents (empty if absent).
func readLog(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

// waitForLog polls the logfile until needle appears or the timeout elapses.
func waitForLog(t *testing.T, path, needle string) bool {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if strings.Contains(readLog(t, path), needle) {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return false
}

// readStubRecord parses the stub's STARTED line into a stubRecord.
func readStubRecord(t *testing.T, path string) stubRecord {
	t.Helper()
	var rec stubRecord
	body := readLog(t, path)
	sc := bufio.NewScanner(strings.NewReader(body))
	for sc.Scan() {
		line := sc.Text()
		if !strings.Contains(line, "STARTED") {
			continue
		}
		for _, tok := range strings.Fields(line) {
			k, v, ok := strings.Cut(tok, "=")
			if !ok {
				continue
			}
			switch k {
			case "pid":
				rec.pid, _ = strconv.Atoi(v)
			case "ppid":
				rec.ppid, _ = strconv.Atoi(v)
			case "isatty":
				rec.isatty = v == "true"
			case "winsize":
				if r, c, ok := strings.Cut(v, "x"); ok {
					rec.rows, _ = strconv.Atoi(r)
					rec.cols, _ = strconv.Atoi(c)
				}
			}
		}
		return rec
	}
	t.Fatalf("no STARTED line in stub log: %q", body)
	return rec
}

// exitCodeOf returns the finished command's exit code. The launcher always exits
// normally (it converts the host's status itself), so ProcessState.ExitCode()
// carries the propagated code.
func exitCodeOf(cmd *exec.Cmd) int {
	if cmd.ProcessState == nil {
		return -1
	}
	return cmd.ProcessState.ExitCode()
}

// TestLaunchResidentParent is AC-1: while the host runs, the launcher stays
// resident as the host's parent (host.ppid == launcher.pid, host.pid !=
// launcher.pid), where the syscall.Exec baseline has no resident parent at all
// (host.pid == launcher.pid). Both observed from the stub's reported identity.
func TestLaunchResidentParent(t *testing.T) {
	t.Run("resident-parent model: host is a child of the live launcher", func(t *testing.T) {
		log := newStubLog(t)
		cmd := roleCmd(t, "launcher", log, "wait")
		if err := cmd.Start(); err != nil {
			t.Fatalf("start launcher: %v", err)
		}
		t.Cleanup(func() { _ = cmd.Process.Kill() })
		if !waitForLog(t, log, "STARTED") {
			t.Fatalf("stub never started; log=%q", readLog(t, log))
		}
		rec := readStubRecord(t, log)
		launcherPid := cmd.Process.Pid
		if rec.ppid != launcherPid {
			t.Fatalf("host.ppid=%d, want launcher.pid=%d (host must be a child of the resident launcher)", rec.ppid, launcherPid)
		}
		if rec.pid == launcherPid {
			t.Fatalf("host.pid=%d == launcher.pid — host must be a distinct child, not the launcher's image", rec.pid)
		}
		// Tear down: SIGTERM the launcher, which forwards to the host.
		_ = cmd.Process.Signal(syscall.SIGTERM)
		_ = cmd.Wait()
	})

	t.Run("exec baseline: no resident parent — host IS the launcher image", func(t *testing.T) {
		log := newStubLog(t)
		cmd := roleCmd(t, "execlauncher", log, "wait")
		if err := cmd.Start(); err != nil {
			t.Fatalf("start execlauncher: %v", err)
		}
		t.Cleanup(func() { _ = cmd.Process.Kill() })
		if !waitForLog(t, log, "STARTED") {
			t.Fatalf("stub never started; log=%q", readLog(t, log))
		}
		rec := readStubRecord(t, log)
		launcherPid := cmd.Process.Pid
		if rec.pid != launcherPid {
			t.Fatalf("baseline host.pid=%d, want == launcher.pid=%d (exec replaces the image)", rec.pid, launcherPid)
		}
		if rec.ppid == launcherPid {
			t.Fatalf("baseline host.ppid=%d == launcher.pid — exec must leave NO resident parent (the metric moves the wrong way)", rec.ppid)
		}
		_ = cmd.Process.Signal(syscall.SIGTERM)
		_ = cmd.Wait()
	})
}

// TestLaunchExitCodePropagation is AC-2: the launcher exits with the host's exit
// code — clean 0, a nonzero N, and a genuine signal death rendered 128+signum.
// The signal-death cases kill the host with no handler (true WIFSIGNALED), so the
// launcher's signal-status branch is exercised, not an explicit host os.Exit.
func TestLaunchExitCodePropagation(t *testing.T) {
	cases := []struct {
		name string
		mode string
		want int
	}{
		{"clean exit 0", "0", 0},
		{"nonzero exit 7", "7", 7},
		{"nonzero exit 42", "42", 42},
		{"SIGTERM death renders 143", "killself:SIGTERM", 143},
		{"SIGINT death renders 130", "killself:SIGINT", 130},
		{"SIGKILL death renders 137", "killself:SIGKILL", 137},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			log := newStubLog(t)
			cmd := roleCmd(t, "launcher", log, tc.mode)
			if err := cmd.Start(); err != nil {
				t.Fatalf("start launcher: %v", err)
			}
			_ = cmd.Wait()
			if got := exitCodeOf(cmd); got != tc.want {
				t.Fatalf("launcher exit=%d, want %d (stub log=%q)", got, tc.want, readLog(t, log))
			}
		})
	}
}
