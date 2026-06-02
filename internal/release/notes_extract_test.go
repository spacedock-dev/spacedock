// ABOUTME: AC-3 extraction proof — the tag cut by AnnotatedTagArgs round-trips
// ABOUTME: through `git tag -l --format='%(contents:body)'`, the form CI feeds goreleaser.
package release

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

// guardConditionRe pulls the empty-body guard condition out of release.yml's
// "Extract release notes from the tag body" step, e.g. the `[ -z "..." ]` test
// inside `if <cond>; then`. The test runs the REAL guard string from the
// workflow, so a regression (e.g. back to the dead `[ ! -s release-notes.txt ]`
// form) breaks this test instead of silently shipping a blank-body Release.
var guardConditionRe = regexp.MustCompile(`(?m)^\s*if (\[.*\]); then\s*$`)

// TestReleaseYAMLGuardRejectsEmptyBody locks the release.yml empty-body guard:
// `git tag -l --format='%(contents:body)'` always appends a trailing newline, so
// the file is never zero-byte and a `[ ! -s ]` guard is dead code. For each tag
// shape the guard's comment claims to catch — lightweight, subject-only, and an
// empty second `-m` — the extracted body is whitespace-only and the guard MUST
// fail (exit non-zero); a real body MUST pass. The guard string is read from the
// workflow itself, not duplicated, so the test exercises what CI actually runs.
func TestReleaseYAMLGuardRejectsEmptyBody(t *testing.T) {
	guard := extractWorkflowGuard(t)

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

	// Lightweight, subject-only annotated, and empty-second-`-m` annotated tags
	// each yield a newline-only %(contents:body); a real body must survive.
	git("tag", "light")
	git("tag", "-a", "subjonly", "-m", "Release 1.0.0")
	git("tag", "-a", "emptybody", "-m", "Release 1.0.0", "-m", "")
	git(AnnotatedTagArgs("realbody", "1.0.0", "- a real user-facing change")...)

	cases := []struct {
		tag        string
		wantReject bool
	}{
		{"light", true},
		{"subjonly", true},
		{"emptybody", true},
		{"realbody", false},
	}
	for _, c := range cases {
		notesPath := filepath.Join(dir, "release-notes.txt")
		// Mirror CI exactly: extract %(contents:body) to release-notes.txt, then
		// evaluate the guard condition against that file.
		body := git("tag", "-l", "--format=%(contents:body)", c.tag)
		if err := os.WriteFile(notesPath, []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", notesPath, err)
		}
		cmd := exec.Command("sh", "-c", "if "+guard+"; then exit 1; fi")
		cmd.Dir = dir
		err := cmd.Run()
		fired := err != nil // guard's `if <cond>; then exit 1` ran exit 1
		if fired != c.wantReject {
			t.Errorf("tag %q: guard fired=%v, want reject=%v (body=%q)", c.tag, fired, c.wantReject, body)
		}
	}
}

// extractWorkflowGuard reads the `if [ ... ]; then` guard condition from the
// Extract-release-notes step of .github/workflows/release.yml, so the test runs
// the same expression CI does rather than a hand-copied duplicate.
func extractWorkflowGuard(t *testing.T) string {
	t.Helper()
	// The release package sits at internal/release; the workflow is at the repo
	// root under .github/workflows.
	path := filepath.Join("..", "..", ".github", "workflows", "release.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	m := guardConditionRe.FindSubmatch(data)
	if m == nil {
		t.Fatalf("no `if [ ... ]; then` guard found in %s", path)
	}
	return string(m[1])
}
