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
	"github.com/spacedock-dev/spacedock/internal/testgit"
)

const testWorkflowLauncher = "/opt/spacedock/bin/spacedock"

// TestMain clears CLAUDE_CODE_SESSION_ID for the whole package before any test
// runs. Merged-mode dispatch (build.go) keys the dispatch filename on this var
// when set, so a developer shell that happens to export it (e.g. a real Claude
// Code session) would otherwise leak its live session id into every golden-
// comparison fixture that exercises a non-bare host=claude dispatch — flaky
// everywhere except that one shell. Individual tests that need a deterministic
// non-empty session id (e.g. the disambiguator tests) still set their own via
// t.Setenv, which scopes and restores per-test over this package-wide clear.
func TestMain(m *testing.M) {
	os.Unsetenv("CLAUDE_CODE_SESSION_ID")
	os.Exit(m.Run())
}

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
	if len(args) > 0 && args[0] == "build" {
		return runNativeWithDefaultClaudeHost(stdin, args...)
	}
	return runNativePreservingHostEnv(stdin, args...)
}

// runNativePreservingHostEnv drives the native dispatch surface without
// normalizing host-marker environment. Host-resolution tests use this helper to
// assert CODEX_THREAD_ID / CLAUDECODE behavior directly.
func runNativePreservingHostEnv(stdin string, args ...string) runResult {
	return runNativeWithLauncher(stdin, testWorkflowLauncher, args...)
}

func runNativeWithLauncher(stdin, workflowLauncher string, args ...string) runResult {
	var stdout, stderr bytes.Buffer
	exit := RunWithLauncher(claudeteam.Probe, workflowLauncher, args, strings.NewReader(stdin), &stdout, &stderr)
	return runResult{stdout.String(), stderr.String(), exit}
}

// runNativeWithDefaultClaudeHost keeps legacy build fixtures deterministic when
// the developer's real shell happens to carry Codex or Pi runtime markers.
func runNativeWithDefaultClaudeHost(stdin string, args ...string) runResult {
	oldCodex, hadCodex := os.LookupEnv("CODEX_THREAD_ID")
	oldClaude, hadClaude := os.LookupEnv("CLAUDECODE")
	oldPiAgent, hadPiAgent := os.LookupEnv("PI_CODING_AGENT")
	oldPiAgentDir, hadPiAgentDir := os.LookupEnv("PI_CODING_AGENT_DIR")
	os.Unsetenv("CODEX_THREAD_ID")
	os.Setenv("CLAUDECODE", "1")
	os.Unsetenv("PI_CODING_AGENT")
	os.Unsetenv("PI_CODING_AGENT_DIR")
	defer restoreEnv("CODEX_THREAD_ID", oldCodex, hadCodex)
	defer restoreEnv("CLAUDECODE", oldClaude, hadClaude)
	defer restoreEnv("PI_CODING_AGENT", oldPiAgent, hadPiAgent)
	defer restoreEnv("PI_CODING_AGENT_DIR", oldPiAgentDir, hadPiAgentDir)
	return runNativePreservingHostEnv(stdin, args...)
}

func restoreEnv(key, value string, had bool) {
	if had {
		os.Setenv(key, value)
		return
	}
	os.Unsetenv(key)
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
	testgit.InitRepo(t, dir, "-q")
	for _, args := range [][]string{
		{"-c", "user.email=t@t", "-c", "user.name=t", "add", "-A"},
		{"-c", "user.email=t@t", "-c", "user.name=t", "commit", "-q", "-m", "init"},
	} {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
}

// gitInitBare initializes a git repo at dir with no seed commit — enough for a
// remote query (`remote get-url origin`) to resolve there. Unlike gitInit it does
// not add/commit, so it works on an empty directory.
func gitInitBare(t *testing.T, dir string) {
	t.Helper()
	testgit.InitRepo(t, dir, "-q")
}

// gitAddOrigin adds a named `origin` remote to the repo at dir, pointing at a
// throwaway bare upstream, so stateHasOrigin reports true. The probe is
// `remote get-url origin` (network-free), so the upstream is never contacted and
// needs no content — it exists only to give the remote a valid URL.
func gitAddOrigin(t *testing.T, dir string) {
	t.Helper()
	upstream := filepath.Join(t.TempDir(), "upstream.git")
	if out, err := exec.Command("git", "init", "-q", "--bare", upstream).CombinedOutput(); err != nil {
		t.Fatalf("git init --bare: %v\n%s", err, out)
	}
	if out, err := exec.Command("git", "-C", dir, "remote", "add", "origin", upstream).CombinedOutput(); err != nil {
		t.Fatalf("git remote add origin: %v\n%s", err, out)
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
