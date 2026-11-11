package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Each route answer below must be rejected so the direct GitHub lookup still
// resolves v9.9.9. Accepting the answer instead resolves the bad tag.
func TestInstallLookupRejectsAnUnusableRouteAnswer(t *testing.T) {
	for _, tc := range []struct{ name, channel, route string }{
		{"stable route fails", "stable", ""},
		{"stable tag path traversal", "stable", "v1.2.3/../../evil"},
		{"stable tag on two lines", "stable", "v1.2.3\nv8.8.8"},
		{"stable prerelease tag", "stable", "v8.8.8-pre1"},
		{"stable foreign host", "stable", "foreign"},
		{"edge route fails", "edge", ""},
		{"edge tag path traversal", "edge", "v1.2.3/../../evil"},
		{"edge HTML page", "edge", "html"},
		{"edge foreign host", "edge", "foreign"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := lookupRun(t)
			stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "SPACEDOCK_CHANNEL="+tc.channel, "STUB_ROUTE="+tc.route)
			if code != 0 || !strings.Contains(stderr, "resolved "+tc.channel+" channel to v9.9.9\n") {
				t.Errorf("exit=%d, want the direct GitHub tag v9.9.9:\n%s", code, stderr)
			}
		})
	}
}

// Changing version discovery, normalization, or marker rules breaks a literal UA.
func TestInstallLookupUserAgent(t *testing.T) {
	for _, tc := range []struct {
		name, os, arch, version, ci, internal, want string
		inPath, readsInput, fails, slow             bool
	}{
		{name: "installed version", os: "Darwin", arch: "aarch64", version: "spacedock v0.27.3\nOS: darwin/arm64", want: "spacedock-install/0.27.3 (darwin; arm64)"},
		{name: "version from PATH", os: "Linux", arch: "x86_64", version: "spacedock 0.28.0-pre10", inPath: true, want: "spacedock-install/0.28.0-pre10 (linux; amd64)"},
		{name: "version cannot read installer input", os: "Linux", arch: "amd64", version: "spacedock 0.27.3", readsInput: true, want: "spacedock-install/0.27.3 (linux; amd64)"},
		{name: "CI value is trimmed and case folded", os: "Linux", arch: "amd64", ci: " \tTrUe\n", want: "spacedock-install/unknown (linux; amd64; ci)"},
		{name: "internal marker", os: "Linux", arch: "amd64", ci: "false", internal: "1", want: "spacedock-install/unknown (linux; amd64; internal)"},
		{name: "markers require explicit values", os: "Linux", arch: "amd64", ci: "enabled; evil", internal: "true", want: "spacedock-install/unknown (linux; amd64)"},
		{name: "unrecognized version", os: "Linux", arch: "arm64", version: "not a version", want: "spacedock-install/unknown (linux; arm64)"},
		{name: "failed version command", os: "Linux", arch: "amd64", version: "spacedock 1.2.3", fails: true, want: "spacedock-install/unknown (linux; amd64)"},
		{name: "slow version command", os: "Linux", arch: "amd64", slow: true, want: "spacedock-install/unknown (linux; amd64)"},
		{name: "version cannot inject a curl option", os: "Linux", arch: "amd64", version: "spacedock 1.2.3' -A evil; $(exit 1)", want: "spacedock-install/unknown (linux; amd64)"},
		{name: "maximum version length", os: "Linux", arch: "amd64", version: "spacedock 1.2.3-" + strings.Repeat("a", 26), want: "spacedock-install/1.2.3-" + strings.Repeat("a", 26) + " (linux; amd64)"},
		{name: "version exceeds contract limit", os: "Linux", arch: "amd64", version: "spacedock 1.2.3-" + strings.Repeat("a", 27), ci: "1", want: "spacedock-install/unknown (linux; amd64; ci)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, curlLog := lookupRun(t)
			r.addProgram("uname", "case \"$1\" in -s) echo '"+tc.os+"' ;; -m) echo '"+tc.arch+"' ;; esac")
			if tc.version != "" || tc.slow {
				r.stdin = "installer input\n"
				body := "printf '%s\\n' '" + strings.ReplaceAll(tc.version, "'", "'\"'\"'") + "'"
				if tc.fails {
					body += "\nexit 1"
				}
				if tc.slow {
					body = "exec sleep 8"
				}
				if tc.readsInput {
					body = "if read -r input; then exit 1; fi\n" + body
				}
				if tc.inPath {
					r.addProgram("spacedock", body)
				} else if err := os.WriteFile(filepath.Join(r.bin, "spacedock"), []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			start := time.Now()
			stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "CI="+tc.ci, "SPACEDOCK_INTERNAL="+tc.internal)
			if code != 0 {
				t.Fatalf("exit=%d:\n%s", code, stderr)
			}
			if tc.slow && time.Since(start) > 4*time.Second {
				t.Error("version probe delayed the lookup for more than four seconds")
			}
			agents, _ := os.ReadFile(curlLog + ".ua")
			if string(agents) != tc.want+"\n" {
				t.Errorf("User-Agents sent to the route=%q, want %q", agents, tc.want+"\n")
			}
		})
	}
}

// A pinned inspection has no release lookup and must not execute an old binary.
func TestInstallLookupSkipsVersionProbeForPinnedSource(t *testing.T) {
	r := newExtrasRun(t, "1.0.1")
	r.addProgram("spacedock", "echo probed > '"+r.log+"'")
	stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "SPACEDOCK_INSTALL_FROM="+r.dist)
	if code != 0 {
		t.Fatalf("exit=%d: %s", code, stderr)
	}
	if _, err := os.Stat(r.log); !os.IsNotExist(err) {
		t.Errorf("pinned inspection executed spacedock: %v", err)
	}
}

func lookupRun(t *testing.T) (*extrasRun, string) {
	t.Helper()
	r := newExtrasRun(t, "1.0.1")
	osName, arch := goosArch(t)
	return r, r.stubCurl(osName, arch)
}

func TestInstallResolvesTheNewestTagByVersionWhenGitHubDoesNotAnswer(t *testing.T) {
	cases := []struct {
		name, channel, tags, want, api, redirect string
	}{
		{"a release outranks its own prereleases", "edge", "v0.16.0 v0.16.0-beta.5 v0.16.0-beta.1", "v0.16.0", "fail", "norelease"},
		{"prerelease numbers compare as numbers", "edge", "v0.28.0-pre9 v0.28.0-pre10 v0.27.3", "v0.28.0-pre10", "fail", "norelease"},
		{"minor versions compare as numbers and non-version tags are ignored", "edge", "v1.9.0 v1.10.0 v1.10.0-rc.1 archive/v0", "v1.10.0", "fail", "norelease"},
		{"null API tag falls back to git", "edge", "v0.16.0", "v0.16.0", "null", "norelease"},
		{"unsafe API tag falls back to git", "edge", "v0.16.0", "v0.16.0", "traversal", "norelease"},
		{"stable skips prereleases", "stable", "v0.28.0-pre5 v0.27.3 v0.9.9 v0.27.10", "v0.27.10", "fail", "norelease"},
		{"direct redirect rejects a foreign host", "stable", "v0.27.3", "v0.27.3", "fail", "foreign"},
		{"direct redirect validates the complete tag path", "stable", "v0.27.3", "v0.27.3", "fail", "nested"},
		{"direct redirect rejects an unsafe tag", "stable", "v0.27.3", "v0.27.3", "fail", "traversal"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := lookupRun(t)
			stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "SPACEDOCK_CHANNEL="+tc.channel,
				"STUB_API="+tc.api, "STUB_REDIRECT="+tc.redirect, "STUB_TAGS="+tc.tags)
			if code != 0 || !strings.Contains(stderr, "resolved "+tc.channel+" channel to "+tc.want+"\n") {
				t.Errorf("exit=%d, want %s:\n%s", code, tc.want, stderr)
			}
			if !strings.Contains(stderr, "reading the tags with git ls-remote") {
				t.Errorf("the fallback is not announced:\n%s", stderr)
			}
		})
	}
}

func TestInstallFailsNamingTheRateLimitWhenEveryRouteFails(t *testing.T) {
	for _, tc := range []struct{ name, channel, tags string }{
		{"stable without a tag", "stable", ""},
		{"edge without a tag", "edge", ""},
		{"stable with only a prerelease tag", "stable", "v0.28.0-pre5"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := lookupRun(t)
			stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "SPACEDOCK_CHANNEL="+tc.channel,
				"STUB_API=fail", "STUB_REDIRECT=fail", "STUB_TAGS="+tc.tags)
			if code != 1 || !strings.Contains(stderr, "rate limit") || !strings.Contains(stderr, "spacedock-dev/spacedock") {
				t.Errorf("exit=%d, want 1 and a message naming the rate limit and the repo:\n%s", code, stderr)
			}
			if strings.Contains(stderr, "resolved ") {
				t.Errorf("a run that resolved nothing claims a version:\n%s", stderr)
			}
		})
	}
}

func TestInstallFinishesTheOtherStepsThenFailsWhenAnExtensionLookupFails(t *testing.T) {
	r, _ := lookupRun(t)
	r.addProgram("claude", "exit 0")

	stderr, code := r.runWith("SPACEDOCK_INSTALL_FROM="+r.dist, "SPACEDOCK_EXTRAS=1", "STUB_LISTING=noalpha")
	if code != 1 {
		t.Errorf("exit = %d, want 1 when an extension release cannot be looked up", code)
	}
	for _, want := range []string{"names no release of alpha", "installed beta v1.0.1", "installed spacedock"} {
		if !strings.Contains(stderr, want) {
			t.Errorf("stderr lacks %q:\n%s", want, stderr)
		}
	}
	if !strings.Contains(r.calls(), "install --host claude") {
		t.Errorf("the plugin step did not run after the failed lookup: %q", r.calls())
	}
	if _, err := os.Stat(filepath.Join(r.bin, "alpha-tui")); err == nil {
		t.Errorf("alpha was installed although the listing names no release of it")
	}
}

func TestInstallExitsZeroWhenAnExtensionDownloadFailsAfterItsTagResolved(t *testing.T) {
	r, _ := lookupRun(t)
	r.addProgram("claude", "exit 0")
	osName, arch := goosArch(t)
	os.Remove(filepath.Join(r.extras, "alpha_0.16.0_"+osName+"_"+arch+".tar.gz"))

	stderr, code := r.runWith("SPACEDOCK_INSTALL_FROM="+r.dist, "SPACEDOCK_EXTRAS=1")
	if code != 0 || !strings.Contains(stderr, "alpha: download failed") {
		t.Errorf("exit=%d, want 0 with the download failure reported:\n%s", code, stderr)
	}
}
