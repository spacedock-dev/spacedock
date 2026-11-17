// ABOUTME: Best-effort, daily Spacedock release hint for interactive launches:
// ABOUTME: the front-door hook, the off switches, quiet modes and terminal guard.
package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/creack/pty"

	"github.com/spacedock-dev/spacedock/internal/release"
	"github.com/spacedock-dev/spacedock/internal/updatecheck"
)

type suiteDeps struct {
	client   *http.Client
	now      func() time.Time
	terminal func(io.Writer) bool
}

// suiteCheckDeps is the test seam for the network client, clock and terminal check.
var suiteCheckDeps = defaultSuiteCheckDeps

func defaultSuiteCheckDeps() *suiteDeps {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableCompression = true // send no header but the User-Agent
	return &suiteDeps{client: &http.Client{Transport: transport}, now: time.Now, terminal: suiteCheckTerminal}
}

const suiteCheckDeadline = 2 * time.Second

// suiteVersionHint checks Spacedock's own release once per UTC day. Extensions
// are not checked: update installs the version their listing names, so a newer
// extension release could name an update that `spacedock update` cannot apply.
func suiteVersionHint(ctx context.Context, fd frontDoorArgs, env []string, stdout io.Writer) {
	getenv := envGetenv(env)
	if fd.noUpdateCheck || getenv("SPACEDOCK_NO_UPDATE_CHECK") == "1" || !suiteInteractive(fd.passthrough) {
		return
	}
	d := suiteCheckDeps()
	defer d.client.CloseIdleConnections()
	if !d.terminal(stdout) {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, suiteCheckDeadline)
	defer cancel()
	if !updatecheck.ClaimDay(getenv, d.now()) {
		return
	}
	version := displayVersion()
	latest := updatecheck.Latest(ctx, d.client, selectedChannelWord(), updatecheck.UserAgent(version, runtime.GOOS, runtime.GOARCH, getenv))
	installed := strings.TrimPrefix(strings.SplitN(version, "+", 2)[0], "v")
	if cmp, err := release.ComparePreVersion(latest, installed); latest != "" && err == nil && cmp > 0 {
		fmt.Fprintf(stdout, "Spacedock %s is available. Run `spacedock update`.\n", latest)
	}
}

// A window-size query succeeds only on a terminal; pty reports an error for
// pipes, files and on platforms without one.
func suiteCheckTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	_, err := pty.GetsizeFull(f)
	return err == nil
}

func suiteInteractive(args []string) bool {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-p") && !strings.HasPrefix(arg, "--") {
			return false
		}
		flag, _, _ := strings.Cut(arg, "=")
		switch flag {
		case "--json", "--jsonl", "--output-format", "--input-format", "--print", "--mode", "--version", "-v", "-V", "--help", "-h", "exec", "e", "review", "completion", "completions", "--non-interactive":
			return false
		}
	}
	return true
}
