// ABOUTME: AC-4(ii) round-trip parity — a #-bearing --set value is written
// ABOUTME: quoted and reads back whole, byte-identically on launcher and oracle.
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestCommentValueRoundTripParity locks option C end-to-end: a --set value
// containing a space-then-`#` is written quoted (so the reader's comment-strip
// does not truncate it) and reads back whole. The narration, on-disk
// frontmatter, and the round-trip --fields read are frozen against the certified
// native output, so the reader-strip AND the writer-quote are pinned together.
func TestCommentValueRoundTripParity(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "source=consolidates #223, #217")
	nOut, nErr, nCode := runNative(t, root, env, args...)
	if nCode != 0 {
		t.Fatalf("exit: native=%d (%q)", nCode, nErr)
	}
	// Narration uses the raw (unquoted) value.
	assertEnvelopeGolden(t, "comment-roundtrip-set", goldenEnvelope{
		stdout: normalize(nOut, root), stderr: normalize(nErr, root), exit: nCode,
	})

	// On-disk frontmatter carries the QUOTED value.
	fm := readFrontmatter(t, filepath.Join(root, "002-vendor-script.md"))
	assertTextGolden(t, "comment-roundtrip-frontmatter", normalize(fm, root))
	if !strings.Contains(fm, `source: "consolidates #223, #217"`) {
		t.Fatalf("expected quoted #-bearing value on disk, got:\n%s", fm)
	}

	// Round-trip read yields the whole value.
	got := ParseFrontmatter(filepath.Join(root, "002-vendor-script.md"))["source"]
	if got != "consolidates #223, #217" {
		t.Fatalf("round-trip read = %q, want whole value", got)
	}
	// A --fields source read renders the #-bearing value back whole.
	read, _, _ := runNative(t, root, env, "--workflow-dir", root, "--fields", "source")
	assertEnvelopeGolden(t, "comment-roundtrip-fields-read", goldenEnvelope{stdout: normalize(read, root)})
	if !strings.Contains(read, "#223") {
		t.Fatalf("did not read the #-bearing value back whole:\n%s", read)
	}
}
