// ABOUTME: Native dispatch harness — drives the native dispatch surface in-process
// ABOUTME: with split stdout/stderr, plus the fixture git-init and write helpers.
package dispatch

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
)

// runResult is one run's three channels.
type runResult struct {
	stdout string
	stderr string
	exit   int
}

// runNative drives the native dispatch surface in-process with the given args
// and stdin, capturing stdout and stderr into separate buffers. HOME is pinned
// via the process env (set by the caller through t.Setenv) so the bare-mode
// team-evidence probe is hermetic.
func runNative(stdin string, args ...string) runResult {
	var stdout, stderr bytes.Buffer
	exit := Run(claudeteam.Probe, args, strings.NewReader(stdin), &stdout, &stderr)
	return runResult{stdout.String(), stderr.String(), exit}
}

// readDispatchBody reads the dispatch body file the run wrote (the path is
// deterministic from the derived name). Returns "" when the run wrote none.
func readDispatchBody(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read dispatch body %s: %v", path, err)
	}
	return string(b)
}

// gitInit initializes a git repo at dir so find_git_root resolves there.
func gitInit(t *testing.T, dir string) {
	t.Helper()
	for _, args := range [][]string{
		{"init", "-q"},
		{"-c", "user.email=t@t", "-c", "user.name=t", "add", "-A"},
		{"-c", "user.email=t@t", "-c", "user.name=t", "commit", "-q", "-m", "init"},
	} {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
}

// writeFile writes content to path, creating parent dirs.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
