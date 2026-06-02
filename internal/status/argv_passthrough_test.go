// ABOUTME: AC-1 argv-passthrough parity — the launcher forwards argv verbatim so
// ABOUTME: the vendored script applies its own semantics (mid-set truncation, unknowns).
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestMidSetUnknownTokenTruncates locks the mid-set truncation case: a token
// starting with -- terminates the --set field-list parse, so
// `--set <slug> --bogus status=done` drops status=done, exits 1 with the
// "requires at least one field=value" error, and leaves the entity unchanged.
func TestMidSetUnknownTokenTruncates(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "--bogus", "status=done")
	nOut, nErr, nCode := runNative(t, root, env, args...)

	if nCode != 1 {
		t.Fatalf("native exit=%d, want 1 (truncated --set has no field=value)", nCode)
	}
	if !strings.Contains(nErr, "requires at least one field=value") {
		t.Fatalf("native stderr=%q, want the truncation error", nErr)
	}
	if nOut != "" {
		t.Fatalf("stdout must be empty: native=%q", nOut)
	}
	assertEnvelopeGolden(t, "argv-midset-truncate", goldenEnvelope{
		stdout: normalize(nOut, root), stderr: normalize(nErr, root), exit: nCode,
	})
	// status=done was dropped: the entity remains in ideation.
	fm := readFrontmatter(t, filepath.Join(root, "002-vendor-script.md"))
	if !strings.Contains(fm, "status: ideation") {
		t.Fatalf("status=done should have been dropped; frontmatter:\n%s", fm)
	}
}

// TestUnknownTopLevelFlagFallsThrough locks the other passthrough case: an
// unrecognized top-level flag is not rejected; the reader ignores it and renders
// the default table at exit 0. Frozen against the certified-parity native output.
func TestUnknownTopLevelFlagFallsThrough(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	args := []string{"--workflow-dir", root, "--bogus-top-level"}

	nOut, nErr, nCode := runNative(t, root, env, args...)

	if nCode != 0 {
		t.Fatalf("native exit=%d stderr=%q, want 0 (unknown top-level flag falls through)", nCode, nErr)
	}
	// Default table renders (header present, entities present).
	if !strings.Contains(nOut, "001-design-seam") {
		t.Fatalf("expected default table with entities, got:\n%s", nOut)
	}
	assertEnvelopeGolden(t, "argv-unknown-toplevel", goldenEnvelope{
		stdout: normalize(nOut, root), stderr: normalize(nErr, root), exit: nCode,
	})
}
