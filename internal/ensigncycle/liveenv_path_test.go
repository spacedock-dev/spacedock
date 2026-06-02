// ABOUTME: Offline unit test for the child-env PATH-prepend helper that puts the
// ABOUTME: built spacedock binary's dir first so the FO resolves `spacedock` by name.
package ensigncycle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
