// ABOUTME: Detection-unit coverage for stateHasOrigin — the network-free
// ABOUTME: named-origin probe distinguishing a remote-backed state checkout from a local one.
package status

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// gitInit initializes a bare-minimum git repo at dir so a remote query resolves
// there. No remote is added, so `git remote get-url origin` exits non-zero.
func gitInitNoRemote(t *testing.T, dir string) {
	t.Helper()
	testgit.InitRepo(t, dir, "-q")
}

// TestStateHasOriginNoRemote (detection unit, seeds AC-1/AC-2 false path): a
// freshly-init'd checkout with no remote returns false — the spike's exit-2
// `No such remote 'origin'` case.
func TestStateHasOriginNoRemote(t *testing.T) {
	dir := t.TempDir()
	gitInitNoRemote(t, dir)
	if stateHasOrigin(dir) {
		t.Fatalf("stateHasOrigin(no-remote checkout) = true, want false")
	}
}

// TestStateHasOriginWithOrigin (detection unit, seeds AC-1/AC-2 true path): a
// checkout with a named `origin` remote returns true — the spike's exit-0 case.
func TestStateHasOriginWithOrigin(t *testing.T) {
	dir := t.TempDir()
	gitInitNoRemote(t, dir)
	upstream := filepath.Join(t.TempDir(), "upstream.git")
	if out, err := exec.Command("git", "init", "-q", "--bare", upstream).CombinedOutput(); err != nil {
		t.Fatalf("git init --bare: %v\n%s", err, out)
	}
	if out, err := exec.Command("git", "-C", dir, "remote", "add", "origin", upstream).CombinedOutput(); err != nil {
		t.Fatalf("git remote add origin: %v\n%s", err, out)
	}
	if !stateHasOrigin(dir) {
		t.Fatalf("stateHasOrigin(origin-backed checkout) = false, want true")
	}
}

// TestStateHasOriginNonRepo (detection unit, defensive): a directory that is not
// a git repo at all returns false — the probe must not panic or report true when
// the git command fails for a non-`origin` reason.
func TestStateHasOriginNonRepo(t *testing.T) {
	if stateHasOrigin(t.TempDir()) {
		t.Fatalf("stateHasOrigin(non-repo dir) = true, want false")
	}
}
