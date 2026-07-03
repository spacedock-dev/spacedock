// ABOUTME: Direct unit test for ManifestVersion — the version source --version's
// ABOUTME: per-runtime block reads; asserts the real fixture version, not just non-empty.
package contract

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TestManifestVersion drives ManifestVersion against a real fixture manifest and
// asserts the EXACT version string the file declares (0.19.8), so a hardcoded
// return would not pass — the test reads the manifest's own version. The
// missing-file case is errNoManifest (a distinct no-plugin state); an unparseable
// manifest is a parse error so the caller renders the bare version rather than a
// fabricated one.
func TestManifestVersion(t *testing.T) {
	got, err := ManifestVersion(filepath.Join("testdata", "compatible.json"))
	if err != nil {
		t.Fatalf("ManifestVersion(compatible.json) err = %v, want nil", err)
	}
	if got != "0.19.8" {
		t.Fatalf("ManifestVersion(compatible.json) = %q, want %q", got, "0.19.8")
	}
}

// TestManifestVersionMissingFile asserts a manifest path that does not exist yields
// errNoManifest, the no-plugin signal the --version probe collapses to an empty
// version (rendering `spacedock not installed`).
func TestManifestVersionMissingFile(t *testing.T) {
	_, err := ManifestVersion(filepath.Join("testdata", "does-not-exist.json"))
	if !errors.Is(err, errNoManifest) {
		t.Fatalf("ManifestVersion(missing) err = %v, want errNoManifest", err)
	}
}

// TestManifestVersionUnparseable writes a manifest file with invalid JSON and
// asserts ManifestVersion returns a (non-errNoManifest) parse error, so the
// --version probe renders the bare/absent version rather than a fabricated one.
func TestManifestVersionUnparseable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "plugin.json")
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	_, err := ManifestVersion(path)
	if err == nil {
		t.Fatalf("ManifestVersion(unparseable) err = nil, want a parse error")
	}
	if errors.Is(err, errNoManifest) {
		t.Fatalf("ManifestVersion(unparseable) err = errNoManifest, want a parse error (the file exists, it just does not parse)")
	}
}
