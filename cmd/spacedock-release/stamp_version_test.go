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

// chdirTemp runs the rest of the test from dir. The no-argument stamp and gate
// defaults are repo-root-relative, so a test exercising them must run where the
// target files live.
func chdirTemp(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(old) })
}

// writeStampTarget stages one stamp-target file at path carrying oldVersion,
// choosing JSON-descriptor or FO-prose content by extension.
func writeStampTarget(t *testing.T, path, oldVersion string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	var body string
	if strings.HasSuffix(path, ".md") {
		majorMinor := oldVersion
		if i := strings.LastIndex(oldVersion, "."); i > 0 {
			majorMinor = oldVersion[:i]
		}
		body = "These skills require binary minor " + majorMinor + " (same major.minor; patch and prerelease skew are fine).\n"
	} else {
		body = `{"name": "spacedock", "version": "` + oldVersion + `"}` + "\n"
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestStampVersionCommandDefaultsToStampTargets locks AC-1: one `stamp-version`
// call with NO target arguments stamps every host descriptor — including the new
// Pi `.pi/plugin.json` — plus the FO prose. The expected target set is written
// out here independently of release.StampTargets: if the default list omitted
// the Pi descriptor, `.pi/plugin.json` would stay at its old version and fail.
func TestStampVersionCommandDefaultsToStampTargets(t *testing.T) {
	dir := t.TempDir()
	targets := []string{
		".claude-plugin/plugin.json",
		".codex-plugin/plugin.json",
		".pi/plugin.json",
		"skills/first-officer/references/first-officer-shared-core.md",
	}
	for _, rel := range targets {
		writeStampTarget(t, filepath.Join(dir, rel), "0.0.0")
	}
	chdirTemp(t, dir)

	if code := stampVersion([]string{"1.2.3"}); code != 0 {
		t.Fatalf("stamp-version with no target arguments exit = %d, want 0", code)
	}
	for _, rel := range targets {
		data, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "1.2.3") && !strings.Contains(string(data), "minor 1.2 ") {
			t.Fatalf("%s was not stamped to 1.2.3 by the default target list: %s", rel, data)
		}
	}
}

// TestStampPathsMatchesDefaultStampList locks AC-2: the list `stamp-paths`
// prints is the list the no-argument stamp rewrites. Every printed path is
// staged old, then stamped; a stamp default list that diverged (or a printed
// path the default omitted) leaves that file unrewritten and fails.
func TestStampPathsMatchesDefaultStampList(t *testing.T) {
	out, code := captureStdout(t, func() int { return stampPaths(nil) })
	if code != 0 {
		t.Fatalf("stamp-paths exit = %d, want 0", code)
	}
	listed := strings.Fields(out)
	if len(listed) == 0 {
		t.Fatal("stamp-paths printed nothing")
	}
	foundPi := false
	for _, p := range listed {
		if p == ".pi/plugin.json" {
			foundPi = true
		}
	}
	if !foundPi {
		t.Fatalf("stamp-paths does not list .pi/plugin.json: %v", listed)
	}

	dir := t.TempDir()
	for _, rel := range listed {
		writeStampTarget(t, filepath.Join(dir, rel), "0.0.0")
	}
	chdirTemp(t, dir)
	if code := stampVersion([]string{"4.5.6"}); code != 0 {
		t.Fatalf("stamp-version over the stamp-paths list exit = %d, want 0", code)
	}
	for _, rel := range listed {
		data, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "4.5.6") && !strings.Contains(string(data), "minor 4.5 ") {
			t.Fatalf("stamp-paths listed %s but the default stamp did not rewrite it: %s", rel, data)
		}
	}
}

// TestStampVersionCommandFailsOnVersionlessNamedManifest locks AC-3: a named
// JSON manifest with no usable top-level version exits non-zero, leaves the
// bytes untouched, and prints NO success line — the loud failure that replaces
// writing unchanged bytes and reporting success.
func TestStampVersionCommandFailsOnVersionlessNamedManifest(t *testing.T) {
	dir := t.TempDir()
	marketplace := filepath.Join(dir, "marketplace.json")
	// The only version is nested on the plugin entry; there is no top level.
	src := `{"name": "spacedock", "plugins": [{"name": "spacedock", "version": "0.0.1"}]}` + "\n"
	if err := os.WriteFile(marketplace, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := captureStdout(t, func() int { return stampVersion([]string{"1.2.3", marketplace}) })
	if code == 0 {
		t.Fatalf("stamp-version exit = 0 on a named manifest with no top-level version; want non-zero")
	}
	if strings.Contains(out, "stamped") {
		t.Fatalf("stamp-version printed a success line for an unstamped manifest: %q", out)
	}
	after, err := os.ReadFile(marketplace)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != src {
		t.Fatalf("stamp-version rewrote a manifest with no top-level version:\nwant %q\ngot  %q", src, after)
	}
}
