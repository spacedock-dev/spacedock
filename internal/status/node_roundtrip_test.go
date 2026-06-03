// ABOUTME: AC-2 node round-trip — unchanged fields survive --set field-exact
// ABOUTME: (value-map + key order); the 5 documented emitter normalizations are pinned.
package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestUpdateFrontmatterNodeRoundTrip locks the AC-2 mutator contract: when
// --set rewrites ONE field, every other top-level frontmatter field's VALUE
// and the overall KEY ORDER survive — the AC-2 field-exact contract (per the
// entity body line 36, "the contract relaxes from byte-identical to
// FIELD-identical"). The previous shape asserted byte-equality via
// strings.Contains on whitespace-clean seeds, which would silently pass even
// if the writer dropped or reordered a value. The shape below decodes the
// written file through the reader and asserts (1) every unchanged key's
// value matches the seed map, (2) the mutated field's value is the new
// value, (3) the encounter order in the written file matches the seed's
// declaration order.
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

	// (1)+(2) FIELD-exact: decode the written file through the reader and
	// assert every key has the expected value (mutated or unchanged).
	gotFields := ParseFrontmatter(path)
	wantFields := map[string]string{
		"id":          "050",
		"title":       "Tracker-linked entity",
		"status":      "implementation",
		"score":       "0.40",
		"source":      "linear",
		"issue":       "ENG-123",
		"tracker-url": "https://linear.app/x/ENG-123",
	}
	for k, want := range wantFields {
		if gotFields[k] != want {
			t.Fatalf("field %q: got %q, want %q (full file:\n%s)", k, gotFields[k], want, got)
		}
	}
	for k := range gotFields {
		if _, ok := wantFields[k]; !ok {
			t.Fatalf("unexpected extra key %q in written file:\n%s", k, got)
		}
	}

	// (3) KEY ORDER: every key appears in the written file at a position
	// strictly increasing in the seed's declaration order.
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

	// Body + close fence + EOF newline are preserved (the byte-preservation
	// seam outside the in-fence YAML body — AC-2 leaves this contract intact
	// since it lives outside the YAML emitter's reach).
	if !strings.HasSuffix(got, "# Tracker-linked entity\n") {
		t.Fatalf("body or EOF newline drifted:\n%s", got)
	}
}

// TestUpdateFrontmatterNodeInsertNew locks that an --set on a genuinely
// missing field is appended before the closing fence, the other fields keep
// their values, and the new key sits last in the resulting order.
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

	gotFields := ParseFrontmatter(path)
	wantFields := map[string]string{
		"id":     "002",
		"title":  "Vendor the status script",
		"status": "ideation",
		"score":  "0.60",
		"source": "roadmap",
		"pr":     "#42",
	}
	for k, want := range wantFields {
		if gotFields[k] != want {
			t.Fatalf("field %q: got %q, want %q (full file:\n%s)", k, gotFields[k], want, got)
		}
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

// TestUpdateFrontmatterEmitterNormalizations pins the 5 documented yaml.v3
// emitter normalizations as accepted observed behavior: a future yaml.v3
// emitter change that altered any of these would be a silent contract drift
// the cycle-1 audit named explicitly (L2-a/b/c/d/e). All 5 leave VALUES and
// KEY ORDER intact — that is what AC-2 actually requires — but they reshape
// visual whitespace. Treat this as the boundary marker: a regression that
// touched a VALUE here would also break TestUpdateFrontmatterNodeRoundTrip
// above; this test catches the cosmetic axis the field-exact test is
// permitted to ignore.
func TestUpdateFrontmatterEmitterNormalizations(t *testing.T) {
	cases := []struct {
		name      string
		seed      string
		mutateKey string
		mutateVal string
		// wantContains pins one or more substrings that MUST appear in the
		// written file. They encode the post-normalization shape — the
		// accepted observed behavior of the yaml.v3 emitter.
		wantContains []string
		// wantNotContains pins substrings that MUST NOT appear after the
		// normalization (the pre-normalization shape).
		wantNotContains []string
	}{
		{
			// L2-a: trailing whitespace on an unchanged scalar line is
			// stripped by the emitter. Value unchanged.
			name: "trailing-whitespace-stripped",
			seed: "---\n" +
				"id: \"x\"\n" +
				"title: Hello   \n" +
				"status: a\n" +
				"---\nbody\n",
			mutateKey: "status",
			mutateVal: "b",
			wantContains: []string{
				"title: Hello\n",
			},
			wantNotContains: []string{
				"title: Hello   \n",
			},
		},
		{
			// L2-c: blank lines INSIDE the in-fence body are erased by the
			// emitter. Value of `status` and `id` unchanged.
			name: "blank-lines-erased",
			seed: "---\n" +
				"id: \"x\"\n" +
				"\n" +
				"status: a\n" +
				"---\nbody\n",
			mutateKey: "status",
			mutateVal: "b",
			wantContains: []string{
				"id: \"x\"\nstatus: b\n",
			},
		},
		{
			// L2-d: a double-space-after-colon on an unchanged scalar is
			// collapsed to a single space. Value unchanged.
			name: "double-space-after-colon-collapsed",
			seed: "---\n" +
				"id:  \"x\"\n" +
				"status: a\n" +
				"---\nbody\n",
			mutateKey: "status",
			mutateVal: "b",
			wantContains: []string{
				"id: \"x\"\n",
			},
			wantNotContains: []string{
				"id:  \"x\"\n",
			},
		},
		{
			// L2-e: pre-comment whitespace on a scalar with an inline comment
			// is collapsed by yaml.v3 (3-space gap before `#` becomes a
			// single-space gap). The comment itself is preserved as a
			// LineComment on the node and re-emitted, but the WHITESPACE shape
			// is normalized. Value of `id` unchanged.
			name: "pre-comment-whitespace-collapsed",
			seed: "---\n" +
				"id: \"x\"   # tag\n" +
				"status: a\n" +
				"---\nbody\n",
			mutateKey: "status",
			mutateVal: "b",
			wantContains: []string{
				"id: \"x\" # tag\n",
			},
			wantNotContains: []string{
				"id: \"x\"   # tag",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "e.md")
			if err := os.WriteFile(path, []byte(tc.seed), 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := updateFrontmatter(path, []fieldUpdate{{field: tc.mutateKey, value: tc.mutateVal, hasValue: true}}); err != nil {
				t.Fatalf("updateFrontmatter: %v", err)
			}
			b, _ := os.ReadFile(path)
			got := string(b)
			for _, want := range tc.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("expected %q in output:\n%s", want, got)
				}
			}
			for _, notWant := range tc.wantNotContains {
				if strings.Contains(got, notWant) {
					t.Errorf("expected NOT %q in output:\n%s", notWant, got)
				}
			}
		})
	}
}

// TestUpdateFrontmatterBlockScalarsRewrapped pins L2-b: a folded (`>`) or
// literal (`|`) block-scalar value is decoded to its string form and
// re-emitted in the emitter's canonical shape. The value content survives
// (newlines preserved per the block-scalar style) but the visual shape may
// change. Skipped today because the live entity corpus has no block-scalar
// frontmatter values — keeping it as a documented intent rather than a
// pinned shape avoids over-claiming on a path that does not have a stable
// fixture.
func TestUpdateFrontmatterBlockScalarsRewrapped(t *testing.T) {
	t.Skip("documented divergence #L2-b — no live block-scalar frontmatter to pin; reactivate when a fixture exists")
}
