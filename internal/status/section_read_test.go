// ABOUTME: --read AC1-AC5 — frontmatter parity, heading-map byte location,
// ABOUTME: fenced-code skip, entity-ref resolution, and flag incompatibility.
package status

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// fixturePath is the committed known-structure markdown the oracle-based --read
// tests slice against (FM + multi-level + fenced-code + stage-report sections).
func fixturePath(t *testing.T) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("testdata", "section-reader", "fixture.md"))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// fixtureLines returns the fixture's lines via splitLines (the same view the
// reader indexes), so an `offset` slices these directly: lines[offset-1:...].
func fixtureLines(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile(fixturePath(t))
	if err != nil {
		t.Fatal(err)
	}
	return splitLines(string(data))
}

// TestReadFrontmatterParity (AC1) asserts the emitted frontmatter equals
// ParseFrontmatter over the same file — every top-level scalar key/value.
func TestReadFrontmatterParity(t *testing.T) {
	path := fixturePath(t)
	sr, ok := readSections(path)
	if !ok {
		t.Fatalf("readSections(%q) failed", path)
	}
	want := ParseFrontmatter(path)
	if len(sr.frontmatter) != len(want) {
		t.Fatalf("frontmatter key count = %d, want %d", len(sr.frontmatter), len(want))
	}
	for k, v := range want {
		if sr.frontmatter[k] != v {
			t.Errorf("frontmatter[%q] = %q, want %q", k, sr.frontmatter[k], v)
		}
	}
}

// TestReadHeadingMapLocatesBytes (AC2) is the in-process equivalent of the
// spike's Read(offset,limit) exercise: for every emitted heading, slicing the
// fixture by offset/lines must equal the section's known text — heading line
// through the line before the next heading of level <= its own (final to EOF).
// The oracle is the fixture's structure (cat -n line numbers), computed here as
// expected (offset, lines) pairs, never a prose match.
func TestReadHeadingMapLocatesBytes(t *testing.T) {
	path := fixturePath(t)
	lines := fixtureLines(t)
	sr, ok := readSections(path)
	if !ok {
		t.Fatalf("readSections(%q) failed", path)
	}

	// Expected structure, by the fixture's cat -n line numbers (the external
	// oracle): heading text -> {level, offset, lines}.
	want := []heading{
		{text: "Problem", level: 2, offset: 10, lines: 4},
		{text: "Proposed approach", level: 2, offset: 14, lines: 13},
		{text: "Sub-detail", level: 3, offset: 18, lines: 9},
		{text: "Stage Report: ideation", level: 2, offset: 27, lines: 3},
	}
	if len(sr.headings) != len(want) {
		t.Fatalf("heading count = %d, want %d (%+v)", len(sr.headings), len(want), sr.headings)
	}
	for i, w := range want {
		got := sr.headings[i]
		if got.text != w.text || got.level != w.level || got.offset != w.offset || got.lines != w.lines {
			t.Fatalf("heading[%d] = %+v, want %+v", i, got, w)
		}
		// Slice the fixture by the emitted offset/lines (1-based, Read semantics)
		// and assert it equals the section's bytes recomputed independently from
		// the known structure (offset through next-heading-of-level<=own minus 1).
		gotSlice := strings.Join(lines[got.offset-1:got.offset-1+got.lines], "\n")
		end := w.offset - 1 + w.lines
		wantSlice := strings.Join(lines[w.offset-1:end], "\n")
		if gotSlice != wantSlice {
			t.Fatalf("heading[%d] %q slice mismatch\n--- got ---\n%s\n--- want ---\n%s", i, w.text, gotSlice, wantSlice)
		}
	}
}

// TestReadFencedHeadingsSkipped (AC3) asserts a `#`-prefixed line inside a code
// fence does NOT appear in the heading map.
func TestReadFencedHeadingsSkipped(t *testing.T) {
	sr, ok := readSections(fixturePath(t))
	if !ok {
		t.Fatal("readSections failed")
	}
	for _, h := range sr.headings {
		if strings.Contains(h.text, "not a heading") {
			t.Fatalf("fenced heading leaked into map: %+v", h)
		}
	}
	// The fixture's fence holds a `#` and a `##` line; the only real `##`/`###`
	// headings are the four asserted in AC2. A fence leak would inflate the count.
	if len(sr.headings) != 4 {
		t.Fatalf("heading count = %d, want 4 (a fenced `#` line leaked)", len(sr.headings))
	}
}

func TestFindSectionSpansUsesFenceSafeHeadingOwnership(t *testing.T) {
	data := []byte("## Pärent\r\nα\r\n### Child\rβ\r\n## Sibling\r\nbody\vcontinued\n```md\n## Fenced\n```\n")
	spans, err := FindSectionSpans(data, []string{"Sibling", "Pärent", "Child"})
	if err != nil {
		t.Fatal(err)
	}
	wantSlices := []string{
		"## Sibling\r\nbody\vcontinued\n```md\n## Fenced\n```\n",
		"## Pärent\r\nα\r\n### Child\rβ\r\n",
		"### Child\rβ\r\n",
	}
	for i, span := range spans {
		if got := string(data[span.Start:span.End]); got != wantSlices[i] {
			t.Fatalf("span[%d] %s [%d,%d) slice = %q, want %q",
				i, span.Heading, span.Start, span.End, got, wantSlices[i])
		}
	}
	if _, err := FindSectionSpans(data, []string{"Fenced"}); err == nil ||
		!strings.Contains(err.Error(), "matches 0 headings") {
		t.Fatalf("fenced selector error = %v, want missing", err)
	}

	ambiguous := []byte("## Same\none\n## Same\ntwo\n")
	if _, err := FindSectionSpans(ambiguous, []string{"Same"}); err == nil ||
		!strings.Contains(err.Error(), "matches 2 headings") {
		t.Fatalf("ambiguous selector error = %v", err)
	}
}

// TestReadEntityRefResolvesLikeResolve (AC4) drives the runner: a valid slug
// yields the map for that entity's file, and an unknown ref exits 1 with the
// resolver's error shape.
func TestReadEntityRefResolvesLikeResolve(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)

	t.Run("valid slug yields its entity's map", func(t *testing.T) {
		out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", "003-wire-cli", "--json")
		if code != 0 {
			t.Fatalf("exit=%d stderr=%q", code, stderr)
		}
		var env map[string]json.RawMessage
		if err := json.Unmarshal([]byte(out), &env); err != nil {
			t.Fatalf("output is not JSON: %v\n%s", err, out)
		}
		var cmd string
		json.Unmarshal(env["command"], &cmd)
		if cmd != "read" {
			t.Fatalf("command = %q, want read", cmd)
		}
		// The resolved path must be the entity's own file (folder-form
		// 003-wire-cli/index.md), proving --read resolved the ref to its file.
		var path string
		json.Unmarshal(env["path"], &path)
		if !strings.Contains(path, "003-wire-cli") {
			t.Fatalf("resolved path = %q, want it under 003-wire-cli", path)
		}
		// frontmatter must carry the resolved entity's own id/slug, not another's.
		var fm map[string]string
		json.Unmarshal(env["frontmatter"], &fm)
		if fm["slug"] != "" && fm["slug"] != "003-wire-cli" {
			t.Fatalf("frontmatter slug = %q, want 003-wire-cli or absent", fm["slug"])
		}
	})

	t.Run("unknown ref exits 1 with resolver error", func(t *testing.T) {
		out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", "no-such-entity")
		if code != 1 {
			t.Fatalf("exit=%d (want 1) stdout=%q", code, out)
		}
		if !strings.Contains(stderr, "unknown reference") {
			t.Fatalf("stderr = %q, want resolver 'unknown reference' error", stderr)
		}
	})
}

// TestReadPathFormResolvesPlainFile asserts the path form: --read of a plain
// filesystem path (not a tracked entity) reads it directly. Drives the binary
// against the committed fixture and checks the heading map matches AC2.
func TestReadPathFormResolvesPlainFile(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixturePath(t), "--json")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	var doc struct {
		Command    string `json:"command"`
		TotalLines string `json:"total_lines"`
		Headings   []struct {
			Text   string `json:"text"`
			Level  string `json:"level"`
			Offset string `json:"offset"`
			Lines  string `json:"lines"`
		} `json:"headings"`
	}
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	if doc.Command != "read" {
		t.Fatalf("command = %q, want read", doc.Command)
	}
	if len(doc.Headings) != 4 {
		t.Fatalf("heading count = %d, want 4", len(doc.Headings))
	}
	// Every value is a string (the all-strings contract): level/offset/lines
	// must parse as ints.
	for _, h := range doc.Headings {
		for _, v := range []string{h.Level, h.Offset, h.Lines} {
			if _, err := strconv.Atoi(v); err != nil {
				t.Fatalf("non-int string field %q in %+v", v, h)
			}
		}
	}
	// total_lines is the exact append offset; the fixture is 29 lines.
	if doc.TotalLines != "29" {
		t.Fatalf("total_lines = %q, want 29", doc.TotalLines)
	}
}

// TestReadFlagIncompatibility (AC5) asserts --read combined with each
// conflicting action flag exits 1 with the established "cannot be combined
// with" message naming --read (the --read branch is checked first, so it owns
// the message for every pairing).
func TestReadFlagIncompatibility(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)

	cases := []struct {
		name  string
		extra []string
	}{
		{"set", []string{"--set", "003-wire-cli", "status=done"}},
		{"next", []string{"--next"}},
		{"boot", []string{"--boot"}},
		{"validate", []string{"--validate"}},
		{"resolve", []string{"--resolve", "003-wire-cli"}},
		{"next-id", []string{"--next-id"}},
		{"archive", []string{"--archive", "003-wire-cli"}},
		{"where", []string{"--where", "status=ideation"}},
		{"archived", []string{"--archived"}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"--workflow-dir", root, "--read", "003-wire-cli"}, tc.extra...)
			out, stderr, code := runNative(t, root, env, args...)
			if code != 1 {
				t.Fatalf("exit=%d (want 1) stdout=%q", code, out)
			}
			if !strings.Contains(stderr, "cannot be combined with") {
				t.Fatalf("stderr = %q, want 'cannot be combined with' message", stderr)
			}
		})
	}
}
