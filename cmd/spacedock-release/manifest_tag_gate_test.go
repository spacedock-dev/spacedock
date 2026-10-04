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

// TestManifestTagGateCommandRejectsMissingTag — with no tag argument the
// subcommand exits with a usage error and does not pass.
func TestManifestTagGateCommandRejectsMissingTag(t *testing.T) {
	if code := runManifestTagGate(nil); code == 0 {
		t.Fatalf("manifest-tag-gate exit = 0 with no tag argument; want non-zero")
	}
}

// TestManifestTagGateDefaultsToStampTargets locks AC-2: with no explicit files
// the gate reads the ONE authoritative release.StampTargets list — the same list
// the no-argument stamp rewrites and reports. A fully-stamped target set must
// pass (a gate default wider than the stamp list would fail to read a path and
// block), and diverging ANY target alone must block (a gate default omitting it
// would still pass while that file disagrees with the tag). The expected set is
// written out here independently of release.StampTargets and of the stamp
// report.
func TestManifestTagGateDefaultsToStampTargets(t *testing.T) {
	targets := []string{
		".claude-plugin/plugin.json",
		".codex-plugin/plugin.json",
		".pi/plugin.json",
		"skills/first-officer/references/first-officer-shared-core.md",
	}
	dir := t.TempDir()
	for _, rel := range targets {
		writeStampTarget(t, filepath.Join(dir, rel), "1.2.3")
	}
	chdirTemp(t, dir)

	if code := runManifestTagGate([]string{"v1.2.3"}); code != 0 {
		t.Fatalf("manifest-tag-gate with no explicit files exit = %d, want 0 for a fully stamped target set", code)
	}
	for _, rel := range targets {
		writeStampTarget(t, filepath.Join(dir, rel), "2.0.0") // diverge from the tag
		if code := runManifestTagGate([]string{"v1.2.3"}); code == 0 {
			t.Fatalf("manifest-tag-gate passed while %s diverged from the tag; the default list omits it", rel)
		}
		writeStampTarget(t, filepath.Join(dir, rel), "1.2.3") // restore for the next target
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
