// ABOUTME: Offline unit test for the child-env PATH-prepend helper that puts the
// ABOUTME: built spacedock binary's dir first so the FO resolves `spacedock` by name.
package ensigncycle

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestClaudeBashEnvFileRestoresStubPATHInInteractiveBash locks the interactive
// Claude transport seam: Claude Code starts Bash tool invocations with its own
// shell environment, then sources CLAUDE_ENV_FILE. The file must restore the
// harness PATH so the fixture gh wins over an operator-installed gh. Each
// session gets a distinct file that is removed with that session's test cleanup.
func TestClaudeBashEnvFileRestoresStubPATHInInteractiveBash(t *testing.T) {
	stubDir := t.TempDir()
	realDir := t.TempDir()
	stubGh := filepath.Join(stubDir, "gh")
	realGh := filepath.Join(realDir, "gh")
	for path, body := range map[string]string{
		stubGh: "#!/bin/sh\necho stub\n",
		realGh: "#!/bin/sh\necho real\n",
	} {
		if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	desiredPath := strings.Join([]string{stubDir, realDir}, string(os.PathListSeparator))
	var firstEnvFile string
	t.Run("first session", func(t *testing.T) {
		env, envFile := withClaudeBashEnvFile(t, []string{
			"HOME=/isolated/home",
			"PATH=" + desiredPath,
			"CLAUDE_ENV_FILE=/stale/session-file",
		}, t.TempDir())
		firstEnvFile = envFile
		if got, ok := envValue(env, "CLAUDE_ENV_FILE"); !ok || got != envFile {
			t.Fatalf("CLAUDE_ENV_FILE = %q (present=%v), want %q", got, ok, envFile)
		}

		// Model the later Bash tool command: it begins without the harness PATH,
		// then sources the file named by CLAUDE_ENV_FILE.
		cmd := exec.Command("/bin/sh", "-c", `. "$CLAUDE_ENV_FILE"; command -v gh`)
		cmd.Env = []string{"PATH=" + realDir, "CLAUDE_ENV_FILE=" + envFile}
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("source CLAUDE_ENV_FILE: %v\n%s", err, out)
		}
		if got := strings.TrimSpace(string(out)); got != stubGh {
			t.Fatalf("gh resolved to %q after sourcing CLAUDE_ENV_FILE, want stub %q", got, stubGh)
		}
	})
	if _, err := os.Stat(firstEnvFile); !os.IsNotExist(err) {
		t.Fatalf("first session env file survived cleanup: stat err = %v", err)
	}

	var secondEnvFile string
	t.Run("second session", func(t *testing.T) {
		_, secondEnvFile = withClaudeBashEnvFile(t, []string{"PATH=" + desiredPath}, t.TempDir())
		if filepath.Dir(secondEnvFile) == filepath.Dir(firstEnvFile) {
			t.Fatalf("session env files share directory %q, want per-session isolation", filepath.Dir(secondEnvFile))
		}
	})
	if _, err := os.Stat(secondEnvFile); !os.IsNotExist(err) {
		t.Fatalf("second session env file survived cleanup: stat err = %v", err)
	}
}

// TestWithBinaryOnPath covers the seam that fixed the live-e2e binary-not-on-PATH
// defect (CI run 26839572693): the FO subprocess inherits PATH verbatim from the
// child env, and the built binary's directory was never on it, so the FO's first
// contract step `spacedock --version` hit `command not found`. withBinaryOnPath
// prepends filepath.Dir(binary) as the FIRST PATH element so an exec.LookPath-style
// resolution finds the just-built binary ahead of any other `spacedock`. No model.
func TestWithBinaryOnPath(t *testing.T) {
	// A binary under a temp dir whose directory is NOT already on the child PATH:
	// this reproduces the CI defect (the runner PATH had the claude dir but never
	// the built-spacedock dir).
	binDir := t.TempDir()
	binary := filepath.Join(binDir, "spacedock")

	// A child env with a PATH that does NOT contain binDir, plus a couple of
	// unrelated keys we expect the helper to preserve untouched.
	otherDirs := strings.Join([]string{"/usr/local/bin", "/usr/bin", "/bin"}, string(os.PathListSeparator))
	childEnv := []string{
		"HOME=/some/isolated/home",
		"PATH=" + otherDirs,
		"ANTHROPIC_API_KEY=sk-ci-api-key",
	}

	got := withBinaryOnPath(childEnv, binary)

	pathVal, ok := envValue(got, "PATH")
	if !ok {
		t.Fatal("PATH must be present in the augmented child env")
	}
	dirs := strings.Split(pathVal, string(os.PathListSeparator))
	if dirs[0] != binDir {
		t.Errorf("PATH[0] = %q, want the binary's dir %q first so `spacedock` resolves by name", dirs[0], binDir)
	}
	// The original PATH entries must be preserved after the prepended dir, so any
	// other tool (claude) the FO shells still resolves.
	if !strings.Contains(pathVal, otherDirs) {
		t.Errorf("PATH = %q must still contain the original entries %q", pathVal, otherDirs)
	}
	// Unrelated keys pass through untouched.
	if home, ok := envValue(got, "HOME"); !ok || home != "/some/isolated/home" {
		t.Errorf("HOME = %q (present=%v), want it preserved untouched", home, ok)
	}
	if key, ok := envValue(got, "ANTHROPIC_API_KEY"); !ok || key != "sk-ci-api-key" {
		t.Errorf("ANTHROPIC_API_KEY = %q (present=%v), want it preserved untouched", key, ok)
	}
}

// TestWithBinaryOnPathNoExistingPath asserts the helper still works when the child
// env carries no PATH entry at all: it synthesizes a PATH=<binary dir> entry so the
// FO can resolve the binary by name even from an env that dropped PATH.
func TestWithBinaryOnPathNoExistingPath(t *testing.T) {
	binDir := t.TempDir()
	binary := filepath.Join(binDir, "spacedock")

	childEnv := []string{"HOME=/some/isolated/home"}

	got := withBinaryOnPath(childEnv, binary)

	pathVal, ok := envValue(got, "PATH")
	if !ok {
		t.Fatal("PATH must be synthesized when the child env has none")
	}
	dirs := strings.Split(pathVal, string(os.PathListSeparator))
	if dirs[0] != binDir {
		t.Errorf("PATH[0] = %q, want the binary's dir %q", dirs[0], binDir)
	}
}
