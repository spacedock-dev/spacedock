// ABOUTME: AC-1 native read parity — NativeRunner stdout/stderr/exit equals the
// ABOUTME: oracle for every read subcommand, after the shared normalization.
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// nativeReadCases are the read subcommands compared native-vs-oracle on the
// seq-workflow fixture (sequential id-style, flat+folder, empty-score, etc.).
var nativeReadCases = []struct {
	name  string
	extra []string
}{
	{"default", nil},
	{"archived", []string{"--archived"}},
	{"next", []string{"--next"}},
	{"validate", []string{"--validate"}},
	{"resolve", []string{"--resolve", "003-wire-cli"}},
	{"short-id", []string{"--short-id", "003-wire-cli"}},
	{"where-status", []string{"--where", "status=ideation"}},
	// `worktree` is a non-default frontmatter key, so it appends as a single
	// extra in both runners. A default-named --fields (e.g. `source`) is NOT a
	// parity case: native de-dupes the duplicate column (captain-approved bug
	// fix) while the oracle still renders it twice — that deliberate divergence
	// is locked by TestFieldsDedupeNoDuplicateDefaultColumns instead.
	{"fields", []string{"--fields", "worktree"}},
	{"all-fields", []string{"--all-fields"}},
}

func TestNativeReadMatchesOracle(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)

	for _, tc := range nativeReadCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"--workflow-dir", root}, tc.extra...)

			nativeOut, nativeErr, nativeCode := runNative(t, root, env, args...)
			assertEnvelopeGolden(t, "native-read-"+tc.name, goldenEnvelope{
				stdout: normalize(nativeOut, root),
				stderr: normalize(nativeErr, root),
				exit:   nativeCode,
			})
		})
	}
}

// TestNativeNextIDFormatAndDeterminism exercises the sd-b32 minting path end to
// end through the runner: the native --next-id candidate is a valid 24-char
// sd-b32 token (FORMAT) and is reproducible across runs with identical inputs
// (DETERMINISM). The candidate's concrete bytes depend on realpathOf(workflowDir)
// so they vary per checkout path; the path-independent SHA-256 + 5-bit-extraction
// derivation is pinned by TestSDB32CandidateDerivationVector.
func TestNativeNextIDFormatAndDeterminism(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "sdb32-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	args := []string{"--workflow-dir", root, "--next-id", "--id-seed", "pinnedseed", "--id-actor", "pinnedactor"}

	out1, errOut, code := runNative(t, root, env, args...)
	if code != 0 {
		t.Fatalf("native exit=%d stderr=%q", code, errOut)
	}
	candidate := strings.TrimSpace(out1)
	if !sdB32IsValidCandidate(candidate) {
		t.Fatalf("candidate %q is not a 24-char sd-b32 token", candidate)
	}

	out2, _, code2 := runNative(t, root, env, args...)
	if code2 != 0 {
		t.Fatalf("second --next-id run exit=%d", code2)
	}
	if got := strings.TrimSpace(out2); got != candidate {
		t.Fatalf("--next-id not deterministic: run1=%q run2=%q", candidate, got)
	}
}

// TestNativeBootMatchesOracle freezes --boot structural + section output for the
// sd-b32 fixture against the certified native golden, masking the minted NEXT_ID
// line and the root prefix.
func TestNativeBootMatchesOracle(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "sdb32-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	args := []string{"--workflow-dir", root, "--boot"}

	nativeOut, nativeErr, nativeCode := runNative(t, root, env, args...)
	if nativeCode != 0 {
		t.Fatalf("native --boot exit=%d stderr=%q", nativeCode, nativeErr)
	}
	normNative := maskBootNextID(stripStateBackend(normalize(nativeOut, root)))
	assertTextGolden(t, "native-boot-sdb32", normNative)
}
