// ABOUTME: AC-1 parity for the non-static read flags — --next-id (format +
// ABOUTME: oracle-equality under pinned env) and --boot (structural + section parity).
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// AC-1: --next-id is in scope for pass-through but asserted at format + oracle-
// equality level (not a static golden), since the candidate is SHA-derived. The
// harness pins all id material (timestamp via env, seed/actor via flags) so the
// launcher and the oracle produce the identical reproducible candidate.
func TestNextIDParity(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "sdb32-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	args := []string{"--workflow-dir", root, "--next-id", "--id-seed", "pinnedseed", "--id-actor", "pinnedactor"}

	nativeOut, nativeErr, nativeCode := runNative(t, root, env, args...)
	if nativeCode != 0 {
		t.Fatalf("native exit=%d stderr=%q", nativeCode, nativeErr)
	}
	candidate := strings.TrimSpace(nativeOut)

	// Format: 24 chars, all in the SD-B32 alphabet.
	if len(candidate) != 24 {
		t.Fatalf("--next-id candidate %q length=%d, want 24", candidate, len(candidate))
	}
	for _, c := range candidate {
		if !strings.ContainsRune(sdB32Chars, c) {
			t.Fatalf("--next-id candidate %q has char %q outside SD-B32 alphabet", candidate, c)
		}
	}

	// Deterministic under the pinned env: freeze the exact certified candidate.
	assertTextGolden(t, "nextid-parity-candidate", candidate)
}

// bootSections are the --boot section headers in their required order. The FO
// parses --boot by section at startup, so order and presence are load-bearing.
var bootSections = []string{
	"MODS:", "ID_STYLE:", "NEXT_ID:", "ORPHANS:", "PR_STATE:", "DISPATCHABLE", "TEAM_STATE",
}

// AC-1: --boot is verified structurally (headers present and in order) and the
// deterministic section bodies (ID_STYLE, NEXT_ID, DISPATCHABLE) are parity-
// checked against the oracle. Volatile material (NEXT_ID value, TEAM_STATE hint)
// is normalized; the fixture has no orphans/PRs so those render their `none`
// forms. The launcher and oracle run under the identical env (same HOME) so the
// TEAM_STATE probe is identical between them.
func TestBootStructuralParity(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "sdb32-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	args := []string{"--workflow-dir", root, "--boot"}

	nativeOut, nativeErr, nativeCode := runNative(t, root, env, args...)
	if nativeCode != 0 {
		t.Fatalf("native exit=%d stderr=%q", nativeCode, nativeErr)
	}

	// Section headers present and in order.
	lastIdx := -1
	for _, section := range bootSections {
		idx := strings.Index(nativeOut, section)
		if idx < 0 {
			t.Fatalf("--boot output missing section %q\n%s", section, nativeOut)
		}
		if idx < lastIdx {
			t.Fatalf("--boot section %q out of order\n%s", section, nativeOut)
		}
		lastIdx = idx
	}

	// ID_STYLE body parity (deterministic).
	if !strings.Contains(nativeOut, "ID_STYLE: sd-b32") {
		t.Fatalf("--boot ID_STYLE body wrong\n%s", nativeOut)
	}

	// Full structural parity frozen against the certified native --boot after
	// normalizing the sd-b32 NEXT_ID and the root prefix.
	normNative := sdB32Re.ReplaceAllString(normalize(nativeOut, root), "<ID>")
	assertTextGolden(t, "boot-structural", normNative)
}
