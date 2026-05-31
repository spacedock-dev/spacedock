// ABOUTME: AC-3 extraction proof — the tag cut by AnnotatedTagArgs round-trips
// ABOUTME: through `git tag -l --format='%(contents:body)'`, the form CI feeds goreleaser.
package release

import (
	"os/exec"
	"strings"
	"testing"
)

// TestAnnotatedTagBodyRoundTrips locks AC-3's seam in a throwaway `git init`
// repo, exercising the REAL tag-cutting path: a tag cut via AnnotatedTagArgs
// (the single source of truth the CLI's cutAnnotatedTag also uses) must have its
// notes land in the tag BODY, so `git tag -l --format='%(contents:body)'`
// returns exactly those notes and is NON-EMPTY. This is the exact extraction
// release.yml feeds goreleaser via --release-notes; a single `-m` would fold the
// notes into the subject and leave the body empty (the seam bug this guards).
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
	// Cut the tag through the SAME arg builder the CLI uses — no hand-built
	// "subject\n\nbody" string the production code never emits.
	git(AnnotatedTagArgs("v9.9.9", "9.9.9", body)...)

	got := strings.TrimRight(git("tag", "-l", "--format=%(contents:body)", "v9.9.9"), "\n")
	if got != body {
		t.Errorf("contents:body round-trip mismatch:\n got: %q\nwant: %q", got, body)
	}
	if strings.TrimSpace(got) == "" {
		t.Errorf("contents:body is empty — notes folded into the subject (the seam bug)")
	}
	if strings.Contains(got, "Release 9.9.9") {
		t.Errorf("contents:body leaked the tag subject:\n%s", got)
	}
}
