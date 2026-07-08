// ABOUTME: `spacedock-release stamp-version` CLI test — one invocation stamps a
// ABOUTME: JSON manifest AND the FO prose file, dispatching by file extension.
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestStampVersionCommandStampsManifestAndProseInOneInvocation is D5's atomic
// multi-file round-trip: one `stamp-version` call over a `.json` manifest AND a
// `.md` prose fixture rewrites the manifest's `version` field AND the prose's
// pinned minor literal — the shape release.yml's stamp steps actually invoke
// (the plugin manifests plus the FO shared-core file in one command).
func TestStampVersionCommandStampsManifestAndProseInOneInvocation(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "plugin.json")
	prosePath := filepath.Join(dir, "first-officer-shared-core.md")

	if err := os.WriteFile(manifestPath, []byte(`{"name": "spacedock", "version": "0.23.0", "skills": "./skills/"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(prosePath, []byte("These skills require binary minor 0.23 (blah).\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if code := stampVersion([]string{"0.24.0", manifestPath, prosePath}); code != 0 {
		t.Fatalf("stampVersion exit = %d, want 0", code)
	}

	manifestOut, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifestOut), `"version": "0.24.0"`) {
		t.Fatalf("manifest not stamped: %s", manifestOut)
	}

	proseOut, err := os.ReadFile(prosePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(proseOut), "These skills require binary minor 0.24 ") {
		t.Fatalf("prose not stamped: %s", proseOut)
	}
	if strings.Contains(string(proseOut), "minor 0.23") {
		t.Fatalf("prose still carries the old minor: %s", proseOut)
	}
}

// TestStampVersionCommandErrorsOnUnstampableProse locks that a `.md` argument
// with no pinned literal (or a duplicated one) errors the whole invocation
// rather than silently leaving the prose untouched.
func TestStampVersionCommandErrorsOnUnstampableProse(t *testing.T) {
	dir := t.TempDir()
	prosePath := filepath.Join(dir, "no-literal.md")
	if err := os.WriteFile(prosePath, []byte("no pinned literal here\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := stampVersion([]string{"0.24.0", prosePath}); code == 0 {
		t.Fatalf("stampVersion exit = 0 on an unstampable prose file; want non-zero")
	}
}
