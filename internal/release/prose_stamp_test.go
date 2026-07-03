// ABOUTME: The `.md` prose-stamp handler — rewrites the single pinned
// ABOUTME: "required binary minor" literal, erroring on zero or multi-match.
package release

import (
	"strings"
	"testing"
)

// TestStampProseVersionRewritesSingleLiteral locks D5: the pinned literal
// ("These skills require binary minor X.Y") is rewritten to the release
// version's major.minor, and every other line in the document is left
// byte-for-byte untouched.
func TestStampProseVersionRewritesSingleLiteral(t *testing.T) {
	src := "# First Officer Shared Core\n\n" +
		"1. **Binary version gate.** ... These skills require binary minor 0.23 (same major.minor; patch and prerelease skew are fine). Abort by class:\n" +
		"   - the required minor is stamped above.\n"

	out, err := StampProseVersion([]byte(src), "0.24.0")
	if err != nil {
		t.Fatalf("StampProseVersion: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "These skills require binary minor 0.24 ") {
		t.Fatalf("stamped doc missing the rewritten literal: %q", got)
	}
	if strings.Contains(got, "minor 0.23") {
		t.Fatalf("stamped doc still carries the old literal: %q", got)
	}
	// Every other line survives untouched.
	if !strings.Contains(got, "the required minor is stamped above.") {
		t.Fatalf("unrelated line was lost or altered: %q", got)
	}
}

// TestStampProseVersionCutsPrereleaseSuffix locks that a prerelease release
// version (the shape release.yml's next-branch stamp actually uses) stamps
// only its major.minor, dropping the -preN suffix.
func TestStampProseVersionCutsPrereleaseSuffix(t *testing.T) {
	src := "These skills require binary minor 0.23 (blah).\n"
	out, err := StampProseVersion([]byte(src), "0.25.0-pre1")
	if err != nil {
		t.Fatalf("StampProseVersion: %v", err)
	}
	if !strings.Contains(string(out), "These skills require binary minor 0.25 ") {
		t.Fatalf("stamped doc = %q, want minor 0.25", out)
	}
}

// TestStampProseVersionErrorsOnZeroMatches locks the atomic-by-construction
// guarantee: a document with no pinned literal errors rather than silently
// no-oping (which would leave the prose stale with no signal).
func TestStampProseVersionErrorsOnZeroMatches(t *testing.T) {
	_, err := StampProseVersion([]byte("no pinned literal here\n"), "0.24.0")
	if err == nil {
		t.Fatalf("StampProseVersion with no literal = nil error, want an error")
	}
}

// TestStampProseVersionErrorsOnMultipleMatches locks the other half of the
// atomic guarantee: TWO pinned literals (a paraphrase gone wrong, or an
// accidental duplicate) errors rather than picking one silently, since a
// partial rewrite would leave self-contradictory instructions.
func TestStampProseVersionErrorsOnMultipleMatches(t *testing.T) {
	src := "These skills require binary minor 0.23 (first).\n" +
		"These skills require binary minor 0.23 (duplicate).\n"
	_, err := StampProseVersion([]byte(src), "0.24.0")
	if err == nil {
		t.Fatalf("StampProseVersion with two literals = nil error, want an error")
	}
}

// TestStampProseVersionErrorsOnUnparseableRelease locks that a release version
// with no parseable major.minor (a malformed CI input) errors rather than
// writing a garbage literal.
func TestStampProseVersionErrorsOnUnparseableRelease(t *testing.T) {
	src := "These skills require binary minor 0.23 (blah).\n"
	_, err := StampProseVersion([]byte(src), "not-a-version")
	if err == nil {
		t.Fatalf("StampProseVersion with an unparseable release version = nil error, want an error")
	}
}

// TestProseMinorRoundTripsStampProseVersion locks that ProseMinor reads back
// exactly what StampProseVersion wrote.
func TestProseMinorRoundTripsStampProseVersion(t *testing.T) {
	src := "These skills require binary minor 0.23 (blah).\n"
	out, err := StampProseVersion([]byte(src), "0.24.3")
	if err != nil {
		t.Fatalf("StampProseVersion: %v", err)
	}
	got, err := ProseMinor(out)
	if err != nil {
		t.Fatalf("ProseMinor: %v", err)
	}
	if got != "0.24" {
		t.Fatalf("ProseMinor = %q, want %q", got, "0.24")
	}
}

// TestProseMinorErrorsOnZeroOrMultipleMatches mirrors the stamp-side guard on
// the read side — the sync test and the tag gate must not silently accept a
// drifted or duplicated literal.
func TestProseMinorErrorsOnZeroOrMultipleMatches(t *testing.T) {
	if _, err := ProseMinor([]byte("no literal here\n")); err == nil {
		t.Fatalf("ProseMinor with no literal = nil error, want an error")
	}
	dup := "These skills require binary minor 0.23 (a).\nThese skills require binary minor 0.23 (b).\n"
	if _, err := ProseMinor([]byte(dup)); err == nil {
		t.Fatalf("ProseMinor with two literals = nil error, want an error")
	}
}
