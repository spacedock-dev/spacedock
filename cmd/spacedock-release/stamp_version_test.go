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

// TestStampVersionCommandFailsOnCaseVariantVersionKey is the exact regression
// the validation reproduced: a named manifest whose top-level version key is
// spelled `Version` — `{"name":"fixture","Version":"1.0.0"}` — is not
// rewritten by the exact-lowercase-key stamp, so the command must read its own
// result back and fail instead of reporting a rewritten path over unchanged
// bytes.
//
// This test fails if the command trusts the bytes it just wrote (or a
// case-insensitive pre-check) rather than re-reading the target: StampVersion
// returns this input unchanged, so without the read-back the command exits 0 and
// reports a stamped path for a version it never installed.
func TestStampVersionCommandFailsOnCaseVariantVersionKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "uppercase.json")
	src := `{"name":"fixture","Version":"1.0.0"}` + "\n"
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := captureStdout(t, func() int { return stampVersion([]string{"7.8.9", path}) })
	if code == 0 {
		t.Fatalf("stamp-version exit = 0 for a manifest whose top-level version key is not lowercase `version`; want non-zero")
	}
	if strings.TrimSpace(out) != "" {
		t.Fatalf("stamp-version reported paths despite failing to stamp: %q", out)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != src {
		t.Fatalf("stamp-version rewrote a manifest it could not stamp:\nwant %q\ngot  %q", src, after)
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

// TestStampVersionReportsStampedPaths locks AC-1 and AC-2: one `stamp-version`
// call with NO target arguments rewrites every host descriptor — including the
// new Pi `.pi/plugin.json` — plus the FO prose, AND reports exactly those paths
// on stdout, one per line, so the release workflow captures the list from this
// run. The expected set is written out here independently of
// release.StampTargets: if the default list omitted the Pi descriptor,
// `.pi/plugin.json` would be neither reported nor rewritten and the report would
// not match.
func TestStampVersionReportsStampedPaths(t *testing.T) {
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

	out, code := captureStdout(t, func() int { return stampVersion([]string{"4.5.6"}) })
	if code != 0 {
		t.Fatalf("stamp-version with no target arguments exit = %d, want 0", code)
	}
	reported := strings.Fields(out)
	if len(reported) != len(targets) {
		t.Fatalf("stamp-version reported %d paths %v, want the %d default targets %v", len(reported), reported, len(targets), targets)
	}
	seen := map[string]bool{}
	for _, p := range reported {
		seen[p] = true
	}
	for _, rel := range targets {
		if !seen[rel] {
			t.Fatalf("stamp-version did not report default target %s: %v", rel, reported)
		}
		data, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "4.5.6") && !strings.Contains(string(data), "minor 4.5 ") {
			t.Fatalf("stamp-version reported %s but did not rewrite it: %s", rel, data)
		}
	}
}

// TestStampVersionCommandFailsOnVersionlessNamedManifest locks AC-3: a named
// JSON manifest with no usable top-level version exits non-zero, leaves the
// bytes untouched, and reports no path — the loud failure that replaces writing
// unchanged bytes and reporting success.
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
	if strings.TrimSpace(out) != "" {
		t.Fatalf("stamp-version reported paths despite failing to stamp: %q", out)
	}
	after, err := os.ReadFile(marketplace)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != src {
		t.Fatalf("stamp-version rewrote a manifest with no top-level version:\nwant %q\ngot  %q", src, after)
	}
}
