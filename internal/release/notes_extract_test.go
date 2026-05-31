// ABOUTME: AC-3 extraction proof — an annotated tag's body round-trips through
// ABOUTME: `git tag -l --format='%(contents:body)'`, the form CI feeds goreleaser.
package release

import (
	"os/exec"
	"strings"
	"testing"
)

// TestAnnotatedTagBodyRoundTrips locks the AC-3 extraction half in a throwaway
// `git init` repo: an annotated tag whose message body carries the release
// notes is extracted byte-for-byte by `git tag -l --format='%(contents:body)'`
// — and that form strips any subject line, which is why release.yml uses it.
// This is the seam CI relies on to feed goreleaser --release-notes.
func TestAnnotatedTagBodyRoundTrips(t *testing.T) {
	dir := t.TempDir()
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		return string(out)
	}

	git("init", "-q")
	git("config", "user.email", "test@example.com")
	git("config", "user.name", "test")
	git("commit", "-q", "--allow-empty", "-m", "seed")

	body := "A clean release.\n- feat: the user value\n- fix: the other thing"
	// A subject line precedes a blank line and then the body, so we prove the
	// `:body` form drops the subject.
	git("tag", "-a", "v9.9.9", "-m", "Release 9.9.9\n\n"+body)

	got := strings.TrimRight(git("tag", "-l", "--format=%(contents:body)", "v9.9.9"), "\n")
	if got != body {
		t.Errorf("contents:body round-trip mismatch:\n got: %q\nwant: %q", got, body)
	}
	if strings.Contains(got, "Release 9.9.9") {
		t.Errorf("contents:body leaked the tag subject:\n%s", got)
	}
}
