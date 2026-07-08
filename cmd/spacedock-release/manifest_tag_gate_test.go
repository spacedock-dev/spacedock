package main

import (
	"os"
	"path/filepath"
	"testing"
)

// writeManifest writes a minimal plugin.json carrying the given version to a temp
// file and returns its path, so the manifest-tag-gate subcommand is exercised
// against a real file read rather than a stubbed version string.
func writeManifest(t *testing.T, version string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "plugin.json")
	body := `{"name": "spacedock", "version": "` + version + `"}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// TestManifestTagGateCommandPassesOnMatch — the subcommand exits 0 when the tag
// semver equals the manifest version (the stamp-then-tag ordering).
func TestManifestTagGateCommandPassesOnMatch(t *testing.T) {
	manifest := writeManifest(t, "0.22.0")
	if code := runManifestTagGate([]string{"v0.22.0", manifest}); code != 0 {
		t.Fatalf("manifest-tag-gate exit = %d, want 0 on a matching tag/manifest", code)
	}
}

// TestManifestTagGateCommandBlocksOnMismatch — the subcommand exits non-zero when
// the tagged commit's manifest still reads a prior release (the v0.20.0 inversion).
func TestManifestTagGateCommandBlocksOnMismatch(t *testing.T) {
	manifest := writeManifest(t, "0.19.9")
	if code := runManifestTagGate([]string{"v0.20.0", manifest}); code == 0 {
		t.Fatalf("manifest-tag-gate exit = 0 on a tag/manifest mismatch; want non-zero (cut blocked)")
	}
}

// TestManifestTagGateCommandChecksEveryManifest — both plugin manifests are
// checked, so a mismatch in either (e.g. the codex manifest lagging) blocks.
func TestManifestTagGateCommandChecksEveryManifest(t *testing.T) {
	good := writeManifest(t, "0.22.0")
	lagging := writeManifest(t, "0.21.0")
	if code := runManifestTagGate([]string{"v0.22.0", good, lagging}); code == 0 {
		t.Fatalf("manifest-tag-gate exit = 0 with a lagging second manifest; want non-zero")
	}
}

// TestManifestTagGateCommandRejectsMissingArgs — with no manifest argument the
// subcommand exits with a usage error and does not pass.
func TestManifestTagGateCommandRejectsMissingArgs(t *testing.T) {
	if code := runManifestTagGate([]string{"v0.22.0"}); code == 0 {
		t.Fatalf("manifest-tag-gate exit = 0 with no manifest argument; want non-zero")
	}
}

// TestManifestTagGateCommandBlocksUnreadableManifest — a manifest path that does
// not exist blocks the cut rather than silently passing.
func TestManifestTagGateCommandBlocksUnreadableManifest(t *testing.T) {
	if code := runManifestTagGate([]string{"v0.22.0", "/no/such/plugin.json"}); code == 0 {
		t.Fatalf("manifest-tag-gate exit = 0 on an unreadable manifest; want non-zero")
	}
}

// writeProse writes an FO shared-core prose fixture stamped at the given minor
// to a temp .md file and returns its path.
func writeProse(t *testing.T, minor string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "first-officer-shared-core.md")
	body := "These skills require binary minor " + minor + " (same major.minor; patch and prerelease skew are fine).\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// TestManifestTagGateCommandChecksProseMinor — D5: the subcommand also takes the
// FO shared-core `.md` prose file, gating on its stamped minor (not the full
// manifest version) against the tag's major.minor.
func TestManifestTagGateCommandChecksProseMinor(t *testing.T) {
	manifest := writeManifest(t, "0.24.0")
	prose := writeProse(t, "0.24")
	if code := runManifestTagGate([]string{"v0.24.0", manifest, prose}); code != 0 {
		t.Fatalf("manifest-tag-gate exit = %d, want 0 when both manifest and prose agree with the tag", code)
	}
}

// TestManifestTagGateCommandBlocksOnProseMinorMismatch — a stable tag whose
// major.minor disagrees with the prose-stamped minor (a forgotten prose stamp)
// blocks the cut even when the JSON manifest agrees.
func TestManifestTagGateCommandBlocksOnProseMinorMismatch(t *testing.T) {
	manifest := writeManifest(t, "0.24.0")
	staleProse := writeProse(t, "0.23") // forgotten stamp
	if code := runManifestTagGate([]string{"v0.24.0", manifest, staleProse}); code == 0 {
		t.Fatalf("manifest-tag-gate exit = 0 with a prose minor lagging the tag; want non-zero")
	}
}
