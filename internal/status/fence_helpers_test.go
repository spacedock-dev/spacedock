// ABOUTME: AC-1 fence-extraction helper corruption tests — frontmatterSlice,
// ABOUTME: contentHasOpeningFence, splitLines, normalizeNewlines guarded directly.
package status

import (
	"testing"
)

// TestFrontmatterSliceHelperCorruption locks frontmatterSlice's contract: it
// returns exactly the bytes BETWEEN the first two `---` fences (the YAML body
// the reader feeds to yaml.v3), with CRLF/CR universal-newline normalized and
// a leading UTF-8 BOM stripped on the first line. Bugs in this helper would
// silently change what the reader sees BUT would NOT be caught by the
// migration-check walk (both sides of that walk go through this same helper
// — the audit's L1-b finding). These cases pin the helper directly.
func TestFrontmatterSliceHelperCorruption(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string // expected slice contents (LF-normalized, trailing \n appended)
	}{
		{
			name: "basic body bytes between fences",
			in:   "---\nid: a\ntitle: b\n---\nbody\n",
			want: "id: a\ntitle: b\n",
		},
		{
			name: "CRLF universal-newline normalized",
			in:   "---\r\nid: a\r\ntitle: b\r\n---\r\nbody\r\n",
			want: "id: a\ntitle: b\n",
		},
		{
			name: "leading BOM stripped on first line",
			in:   utf8BOM + "---\nid: a\n---\nbody\n",
			want: "id: a\n",
		},
		{
			name: "lone CR (old Mac) universal-newline normalized",
			in:   "---\rid: a\r---\rbody\r",
			want: "id: a\n",
		},
		{
			name: "missing close fence yields body to EOF",
			in:   "---\nid: a\nstill in body\n",
			want: "id: a\nstill in body\n",
		},
		{
			name: "no opening fence yields nil/empty",
			in:   "# heading\nid: a\n",
			want: "",
		},
		{
			name: "empty in-fence body",
			in:   "---\n---\nbody\n",
			want: "",
		},
		{
			name: "second --- after body is NOT a re-open",
			in:   "---\nid: a\n---\nbody\n---\ntail\n",
			want: "id: a\n",
		},
		{
			// Mirrors the sibling TestContentHasOpeningFenceHelperCorruption's
			// "leading blank lines skipped" case (the opening fence is on a
			// later line, not line 0). A bug-injection that reduced
			// frontmatterSlice to `if lines[0] != "---" { return nil }` would
			// pass the rest of this suite but lose the YAML body here — the
			// fence finder must scan past leading truly-empty lines, matching
			// contentHasOpeningFence's behavior on the same input shape.
			name: "leading blank lines skipped",
			in:   "\n\n---\nid: a\n---\nbody\n",
			want: "id: a\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := string(frontmatterSlice([]byte(tc.in)))
			if got != tc.want {
				t.Fatalf("frontmatterSlice(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// TestContentHasOpeningFenceHelperCorruption locks contentHasOpeningFence's
// contract directly. Bugs in this helper would cause the reader to either
// (a) silently return an empty fields map for a real entity or (b) mistake
// a non-entity for an entity. Both would survive the migration-check walk
// (which short-circuits on no-fence files), so it is pinned here.
func TestContentHasOpeningFenceHelperCorruption(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"plain opening fence", "---\nid: 1\n---\n", true},
		{"leading blank lines skipped", "\n\n---\nid: 1\n---\n", true},
		{"BOM-prefixed first line", utf8BOM + "---\nid: 1\n---\n", true},
		{"CRLF opening fence", "---\r\nid: 1\r\n---\r\n", true},
		{"whitespace-only first content line disqualifies", "   \n---\n", false},
		{"text before fence disqualifies", "hello\n---\n", false},
		{"missing close fence still recognized as opening", "---\nid: 1\n", true},
		{"empty input", "", false},
		{"single dash line is not a fence", "--\nid: 1\n", false},
		{"four-dash line is not a fence", "----\nid: 1\n", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := contentHasOpeningFence([]byte(tc.in))
			if got != tc.want {
				t.Fatalf("contentHasOpeningFence(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// TestSplitLinesNormalizeHelperCorruption locks splitLines + normalizeNewlines
// directly. Both are used by the fence finder AND by the mutator's in-fence
// splice; a bug here would mis-locate the close fence (returning the wrong
// body bytes) or fail to canonicalize CRLF, neither of which the
// migration-check walk could catch (it consumes the OUTPUT of these
// helpers, not their internals).
func TestSplitLinesNormalizeHelperCorruption(t *testing.T) {
	t.Run("LF-only file with terminal newline yields no empty trailing element", func(t *testing.T) {
		got := splitLines("a\nb\n")
		want := []string{"a", "b"}
		if !fenceHelpersEqualStrings(got, want) {
			t.Fatalf("splitLines(\"a\\nb\\n\") = %v, want %v", got, want)
		}
	})
	t.Run("LF-only file without terminal newline keeps last line", func(t *testing.T) {
		got := splitLines("a\nb")
		want := []string{"a", "b"}
		if !fenceHelpersEqualStrings(got, want) {
			t.Fatalf("splitLines(\"a\\nb\") = %v, want %v", got, want)
		}
	})
	t.Run("CRLF file is normalized line-by-line", func(t *testing.T) {
		got := splitLines("a\r\nb\r\n")
		want := []string{"a", "b"}
		if !fenceHelpersEqualStrings(got, want) {
			t.Fatalf("splitLines CRLF = %v, want %v", got, want)
		}
	})
	t.Run("lone CR (old Mac) is normalized", func(t *testing.T) {
		got := splitLines("a\rb\r")
		want := []string{"a", "b"}
		if !fenceHelpersEqualStrings(got, want) {
			t.Fatalf("splitLines lone CR = %v, want %v", got, want)
		}
	})
	t.Run("empty input yields nil", func(t *testing.T) {
		got := splitLines("")
		if got != nil {
			t.Fatalf("splitLines(\"\") = %v, want nil", got)
		}
	})
	t.Run("normalizeNewlines is identity on pure LF", func(t *testing.T) {
		in := "a\nb\nc\n"
		if got := normalizeNewlines(in); got != in {
			t.Fatalf("normalizeNewlines pure-LF mutated: %q", got)
		}
	})
	t.Run("normalizeNewlines collapses CRLF before lone CR", func(t *testing.T) {
		// `\r\n` is one newline (two bytes), `\r` alone is one newline. A
		// CRLF-pair must NOT yield two LFs.
		got := normalizeNewlines("a\r\nb\rc\r\nd")
		want := "a\nb\nc\nd"
		if got != want {
			t.Fatalf("normalizeNewlines mixed = %q, want %q", got, want)
		}
	})
}

// fenceHelpersEqualStrings is a slice equality helper kept local to this test file so the
// fence-helper tests stay self-contained.
func fenceHelpersEqualStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
