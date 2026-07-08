// ABOUTME: AC-1's divergeable proof — restamping ONLY the release version (no
// ABOUTME: hand-edited integer) moves the SAME binary from compatible to too-old.
package release

import (
	"encoding/json"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// TestStampVersionFlipsCompatibilityFloor is AC-1's load-bearing fixture test:
// restamping ONLY the release version — StampVersion on a manifest copy,
// StampProseVersion on a prose copy, no other edit — moves the compatibility
// floor. The SAME fixed binary version flips from Compatible to TooOldBinary
// purely because the stamp advanced the manifest/prose minor past it; no
// hand-maintained integer exists anywhere in this path to forget.
func TestStampVersionFlipsCompatibilityFloor(t *testing.T) {
	const binaryVersion = "0.24.5" // fixed throughout; never re-stamped

	manifestFixture := []byte(`{"name":"spacedock","version":"0.0.0","skills":"./skills/"}`)
	proseFixture := []byte("Some preamble.\n" +
		"These skills require binary minor 0.0 (same major.minor; patch and prerelease skew are fine).\n" +
		"More text.\n")

	// Stamp both copies to the binary's own minor (0.24.0): compatible.
	beforeManifest, err := StampVersion(manifestFixture, "0.24.0")
	if err != nil {
		t.Fatalf("StampVersion (before): %v", err)
	}
	beforeProse, err := StampProseVersion(proseFixture, "0.24.0")
	if err != nil {
		t.Fatalf("StampProseVersion (before): %v", err)
	}

	beforePluginVersion := mustManifestVersion(t, beforeManifest)
	beforeProseMinor := mustProseMinorField(t, beforeProse)

	beforeManifestVerdict := contract.Compare("claude", beforePluginVersion, binaryVersion)
	if beforeManifestVerdict.Verdict != contract.Compatible {
		t.Fatalf("before restamp: manifest verdict = %v, want Compatible", beforeManifestVerdict.Verdict)
	}
	beforeProseVerdict := contract.Compare("claude", beforeProseMinor+".0", binaryVersion)
	if beforeProseVerdict.Verdict != contract.Compatible {
		t.Fatalf("before restamp: prose gate verdict = %v, want Compatible", beforeProseVerdict.Verdict)
	}

	// Restamp ONLY the release version, on the ORIGINAL fixture copies (no
	// other edit), to one minor ahead of the fixed binary.
	afterManifest, err := StampVersion(manifestFixture, "0.25.0")
	if err != nil {
		t.Fatalf("StampVersion (after): %v", err)
	}
	afterProse, err := StampProseVersion(proseFixture, "0.25.0")
	if err != nil {
		t.Fatalf("StampProseVersion (after): %v", err)
	}

	afterPluginVersion := mustManifestVersion(t, afterManifest)
	afterProseMinor := mustProseMinorField(t, afterProse)

	afterManifestVerdict := contract.Compare("claude", afterPluginVersion, binaryVersion)
	if afterManifestVerdict.Verdict != contract.TooOldBinary {
		t.Fatalf("after restamp: manifest verdict = %v, want TooOldBinary — the SAME binary %s must now be rejected", afterManifestVerdict.Verdict, binaryVersion)
	}
	afterProseVerdict := contract.Compare("claude", afterProseMinor+".0", binaryVersion)
	if afterProseVerdict.Verdict != contract.TooOldBinary {
		t.Fatalf("after restamp: prose gate verdict = %v, want TooOldBinary — the SAME binary %s must now be rejected", afterProseVerdict.Verdict, binaryVersion)
	}
}

func mustManifestVersion(t *testing.T, manifest []byte) string {
	t.Helper()
	var m struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(manifest, &m); err != nil {
		t.Fatalf("parse stamped manifest: %v", err)
	}
	return m.Version
}

func mustProseMinorField(t *testing.T, prose []byte) string {
	t.Helper()
	minor, err := ProseMinor(prose)
	if err != nil {
		t.Fatalf("read stamped prose minor: %v", err)
	}
	return minor
}
