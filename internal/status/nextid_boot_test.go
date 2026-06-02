// ABOUTME: AC-1 checks for the non-static read flags — --next-id (format +
// ABOUTME: determinism + a path-independent derivation vector) and --boot structure.
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// sdB32IsValidCandidate reports whether s is a syntactically valid sd-b32 id:
// exactly sdB32IDLength chars, all in the alphabet.
func sdB32IsValidCandidate(s string) bool {
	if len(s) != sdB32IDLength {
		return false
	}
	for _, c := range s {
		if !strings.ContainsRune(sdB32Chars, c) {
			return false
		}
	}
	return true
}

// TestNextIDFormatAndDeterminism asserts the two machine-independent properties
// the --next-id candidate must hold (the captain steer: logical soundness, not
// byte-parity with the retired Python oracle): (a) FORMAT — a 24-char valid
// sd-b32 token — and (b) DETERMINISM — the candidate is a pure function of the
// digest material, so two runs with identical inputs yield the identical id. The
// candidate hashes realpathOf(workflowDir) (identity.go:183), so its concrete
// bytes vary per checkout path; freezing a static golden was machine-dependent
// (red-by-construction on a fresh clone). The path-independent derivation is
// asserted by TestSDB32CandidateDerivationVector instead.
func TestNextIDFormatAndDeterminism(t *testing.T) {
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
		t.Fatalf("--next-id candidate %q is not a 24-char sd-b32 token", candidate)
	}

	// Determinism: a second run over identical inputs reproduces the candidate.
	out2, _, code2 := runNative(t, root, env, args...)
	if code2 != 0 {
		t.Fatalf("second --next-id run exit=%d", code2)
	}
	if got := strings.TrimSpace(out2); got != candidate {
		t.Fatalf("--next-id not deterministic: run1=%q run2=%q", candidate, got)
	}
}

// TestSDB32CandidateDerivationVector pins the SHA-256 + big-endian 5-bit-window
// minting against a FIXED, path-INDEPENDENT digest-material vector, so the same
// expected id reproduces on any machine. The workflow path is a literal that
// realpathOf returns unchanged (its `/spacedock/...` prefix does not exist, so
// realpathOf cleans it to itself on every host), and seed/actor/timestamp/nonce
// are fixed. The expected id was derived from an independent reference
// implementation of encode_sd_b32_digest with the correct 5-bit mask (31); the
// 31->30 mask flip at identity.go:151 changes the output, so this vector is the
// minting-regression detector that the old static golden was meant to be — minus
// the machine coupling.
func TestSDB32CandidateDerivationVector(t *testing.T) {
	const (
		fixedWorkflow  = "/spacedock/derivation-vector/wf"
		fixedSeed      = "fixed-seed"
		fixedActor     = "fixed-actor"
		fixedTimestamp = "2026-01-01T00:00:00.000000Z"
		fixedNonce     = 0
		// Derived independently (mask 31) over the fixed digest material; a 31->30
		// flip yields a different string, so this catches the masking regression.
		wantID = "n1y6x7mw00etcc0v1c9mzwwr"
	)
	got := sdB32Candidate(fixedWorkflow, fixedSeed, fixedActor, fixedTimestamp, fixedNonce, env{})
	if got != wantID {
		t.Fatalf("sd-b32 derivation vector mismatch: got %q want %q (5-bit mask regression?)", got, wantID)
	}
}

// bootSections are the --boot section headers in their required order. The FO
// parses --boot by section at startup, so order and presence are load-bearing.
var bootSections = []string{
	"MODS:", "ID_STYLE:", "NEXT_ID:", "ORPHANS:", "PR_STATE:", "DISPATCHABLE", "TEAM_STATE",
}

// AC-1: --boot is verified structurally (headers present and in order) and the
// deterministic section bodies (ID_STYLE, NEXT_ID, DISPATCHABLE) are frozen
// against the certified native golden. Volatile material (the minted NEXT_ID
// value) is masked; the fixture has no orphans/PRs so those render their `none`
// forms, and the pinned (empty) HOME makes the TEAM_STATE probe deterministic.
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

	// Structural --boot frozen against the certified native output: the minted
	// NEXT_ID line is masked (path-dependent), every other token — including the
	// fixed fixture id in the DISPATCHABLE table — freezes literally.
	normNative := maskBootNextID(normalize(nativeOut, root))
	assertTextGolden(t, "boot-structural", normNative)
}
