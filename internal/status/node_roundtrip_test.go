// ABOUTME: AC-2 node round-trip — unchanged scalar nodes survive --set
// ABOUTME: byte-identically and key order is preserved through the mutator.
package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestUpdateFrontmatterNodeRoundTrip locks the AC-2 mutator contract: when
// --set rewrites ONE field, every other top-level frontmatter field's bytes
// are preserved (field-exact, not necessarily byte-exact across the whole
// file — the rewritten field's line is what changes) AND the order of keys
// is stable. This guards the yaml.Node parse->mutate-target->marshal seam:
// unchanged scalar nodes must re-marshal identically.
func TestUpdateFrontmatterNodeRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "e.md")
	const seed = "---\n" +
		"id: \"050\"\n" +
		"title: Tracker-linked entity\n" +
		"status: ideation\n" +
		"score: \"0.40\"\n" +
		"source: linear\n" +
		"issue: ENG-123\n" +
		"tracker-url: https://linear.app/x/ENG-123\n" +
		"---\n" +
		"# Tracker-linked entity\n"
	if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := updateFrontmatter(path, []fieldUpdate{{field: "status", value: "implementation", hasValue: true}}); err != nil {
		t.Fatalf("updateFrontmatter: %v", err)
	}
	b, _ := os.ReadFile(path)
	got := string(b)

	// Unchanged scalar lines survive byte-identically — quoted strings keep
	// their quote shape (id, score), unknown fields keep their value verbatim.
	for _, line := range []string{
		"id: \"050\"",
		"title: Tracker-linked entity",
		"score: \"0.40\"",
		"source: linear",
		"issue: ENG-123",
		"tracker-url: https://linear.app/x/ENG-123",
	} {
		if !strings.Contains(got, line) {
			t.Fatalf("expected unchanged line %q in:\n%s", line, got)
		}
	}
	// The rewritten field landed.
	if !strings.Contains(got, "status: implementation") {
		t.Fatalf("expected mutated status line in:\n%s", got)
	}
	// Body + close fence + EOF newline are preserved.
	if !strings.HasSuffix(got, "# Tracker-linked entity\n") {
		t.Fatalf("body or EOF newline drifted:\n%s", got)
	}

	// Key order is stable: id < title < status < score < source < issue < tracker-url.
	wantOrder := []string{"id:", "title:", "status:", "score:", "source:", "issue:", "tracker-url:"}
	lastIdx := -1
	for _, key := range wantOrder {
		idx := strings.Index(got, key)
		if idx < 0 {
			t.Fatalf("key %q missing from output:\n%s", key, got)
		}
		if idx < lastIdx {
			t.Fatalf("key %q reordered (idx %d < last %d):\n%s", key, idx, lastIdx, got)
		}
		lastIdx = idx
	}
}

// TestUpdateFrontmatterNodeInsertNew locks that an --set on a genuinely
// missing field is appended before the closing fence, the other keys keep
// their bytes, and the new key sits last in the resulting order.
func TestUpdateFrontmatterNodeInsertNew(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "e.md")
	const seed = "---\n" +
		"id: \"002\"\n" +
		"title: Vendor the status script\n" +
		"status: ideation\n" +
		"score: \"0.60\"\n" +
		"source: roadmap\n" +
		"---\nbody\n"
	if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := updateFrontmatter(path, []fieldUpdate{{field: "pr", value: "#42", hasValue: true}}); err != nil {
		t.Fatalf("updateFrontmatter: %v", err)
	}
	b, _ := os.ReadFile(path)
	got := string(b)

	for _, line := range []string{
		"id: \"002\"",
		"title: Vendor the status script",
		"status: ideation",
		"score: \"0.60\"",
		"source: roadmap",
	} {
		if !strings.Contains(got, line) {
			t.Fatalf("expected unchanged line %q in:\n%s", line, got)
		}
	}
	if !strings.Contains(got, "pr: \"#42\"") && !strings.Contains(got, "pr: '#42'") && !strings.Contains(got, "pr: #42") {
		t.Fatalf("expected appended pr line in:\n%s", got)
	}
	// pr lands after source (the prior last key), before the close fence.
	srcIdx := strings.Index(got, "source: roadmap")
	prIdx := strings.Index(got, "pr:")
	closeIdx := strings.Index(got[srcIdx:], "\n---")
	if prIdx < srcIdx {
		t.Fatalf("pr appeared before source:\n%s", got)
	}
	if closeIdx < 0 || srcIdx+closeIdx < prIdx {
		t.Fatalf("pr appeared after the close fence:\n%s", got)
	}
}
