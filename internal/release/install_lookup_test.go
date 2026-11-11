package release

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Removing the route, redirect following, or edge selector breaks this proof.
func TestInstallLookupRoute(t *testing.T) {
	for _, tc := range []struct{ channel, request string }{
		{"stable", "https://spacedock.md/latest-release/spacedock|spacedock-install/unknown (linux; amd64)|1|%{url_effective}|GET\n"},
		{"edge", "https://spacedock.md/latest-release/spacedock?prerelease=1|spacedock-install/unknown (linux; amd64)|1||GET\n"},
	} {
		t.Run(tc.channel, func(t *testing.T) {
			r, log := lookupRun(t)
			r.addProgram("uname", "case \"$1\" in -s) echo Linux ;; -m) echo x86_64 ;; esac")
			stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "STUB_ROUTE=ok", "SPACEDOCK_CHANNEL="+tc.channel)
			if code != 0 || !strings.Contains(stderr, "resolved "+tc.channel+" channel to v9.9.9\n") {
				t.Fatalf("exit=%d; route tag was not resolved:\n%s", code, stderr)
			}
			requests, err := os.ReadFile(log + ".requests")
			if err != nil || string(requests) != tc.request {
				t.Errorf("requests=%q (%v), want %q", requests, err, tc.request)
			}
		})
	}
}

// Using GitHub first or dropping its fallback breaks the literal request order.
func TestInstallLookupRouteFailureUsesDirectGitHub(t *testing.T) {
	for _, tc := range []struct{ channel, urls string }{
		{"stable", "https://spacedock.md/latest-release/spacedock\nhttps://github.com/spacedock-dev/spacedock/releases/latest\n"},
		{"edge", "https://spacedock.md/latest-release/spacedock?prerelease=1\nhttps://api.github.com/repos/spacedock-dev/spacedock/releases?per_page=1\n"},
	} {
		t.Run(tc.channel, func(t *testing.T) {
			for _, route := range []string{"fail", "malformed", "traversal", "whitespace", "metacharacters"} {
				t.Run(route, func(t *testing.T) {
					r, log := lookupRun(t)
					stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "SPACEDOCK_CHANNEL="+tc.channel, "STUB_ROUTE="+route)
					if code != 0 || !strings.Contains(stderr, "resolved "+tc.channel+" channel to v9.9.9\n") {
						t.Fatalf("exit=%d; direct GitHub fallback failed:\n%s", code, stderr)
					}
					urls, err := os.ReadFile(log)
					if err != nil || string(urls) != tc.urls {
						t.Errorf("urls=%q (%v), want %q", urls, err, tc.urls)
					}
				})
			}
		})
	}
}

// A separate count request or routing only Spacedock breaks this proof.
func TestInstallLookupEachProductOnce(t *testing.T) {
	r, log := lookupRun(t)
	stderr, code := r.runWith("STUB_ROUTE=ok", "SPACEDOCK_PLUGINS=none")
	if code != 0 {
		t.Fatalf("exit=%d:\n%s", code, stderr)
	}
	urls, err := os.ReadFile(log)
	if err != nil {
		t.Fatal(err)
	}
	var lookups []string
	for _, url := range strings.Fields(string(urls)) {
		if !strings.Contains(url, "/releases/download/") {
			lookups = append(lookups, url)
		}
	}
	if got := strings.Join(lookups, "\n"); got != "https://spacedock.md/latest-release/spacedock\nhttps://spacedock.md/latest-release/alpha\nhttps://spacedock.md/latest-release/beta" {
		t.Errorf("release lookups=%q", got)
	}
	requests, _ := os.ReadFile(log + ".requests")
	for _, request := range strings.Split(strings.TrimSpace(string(requests)), "\n") {
		if strings.HasPrefix(request, "https://spacedock.md/") && !strings.Contains(request, "|spacedock-install/unknown (") {
			t.Errorf("lookup used the newly installed version instead of the original client version: %s", request)
		}
	}
}

// Changing version discovery, normalization, or marker rules breaks a literal UA.
func TestInstallLookupUserAgent(t *testing.T) {
	for _, tc := range []struct {
		name, os, arch, version, ci, internal, want string
		inPath, readsInput                          bool
	}{
		{name: "installed version", os: "Darwin", arch: "aarch64", version: "spacedock v0.27.3\nOS: darwin/arm64", want: "spacedock-install/0.27.3 (darwin; arm64)"},
		{name: "version cannot read installer input", os: "Linux", arch: "amd64", version: "spacedock 0.27.3", readsInput: true, want: "spacedock-install/0.27.3 (linux; amd64)"},
		{name: "version from PATH", os: "Linux", arch: "x86_64", version: "spacedock 0.28.0-pre10", inPath: true, want: "spacedock-install/0.28.0-pre10 (linux; amd64)"},
		{name: "case folded platform", os: "dArWiN", arch: "X64", want: "spacedock-install/unknown (darwin; amd64)"},
		{name: "CI one", os: "Linux", arch: "amd64", ci: "1", want: "spacedock-install/unknown (linux; amd64; ci)"},
		{name: "CI true trimmed", os: "Linux", arch: "amd64", ci: " \tTrUe\n", want: "spacedock-install/unknown (linux; amd64; ci)"},
		{name: "CI yes", os: "Linux", arch: "amd64", ci: "YES", want: "spacedock-install/unknown (linux; amd64; ci)"},
		{name: "CI on and internal", os: "Linux", arch: "amd64", ci: "on", internal: "1", want: "spacedock-install/unknown (linux; amd64; ci; internal)"},
		{name: "internal only", os: "Linux", arch: "amd64", ci: "false", internal: "1", want: "spacedock-install/unknown (linux; amd64; internal)"},
		{name: "markers require explicit values", os: "Linux", arch: "amd64", ci: "enabled", internal: "true", want: "spacedock-install/unknown (linux; amd64)"},
		{name: "unrecognized version", os: "Linux", arch: "arm64", version: "not a version", want: "spacedock-install/unknown (linux; arm64)"},
		{name: "development version is unknown", os: "Linux", arch: "arm64", version: "spacedock dev", ci: "true", want: "spacedock-install/unknown (linux; arm64; ci)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, log := lookupRun(t)
			r.addProgram("uname", "case \"$1\" in -s) echo '"+tc.os+"' ;; -m) echo '"+tc.arch+"' ;; esac")
			if tc.version != "" {
				r.stdin = "installer input\n"
				body := "printf '%s\\n' '" + tc.version + "'"
				if tc.readsInput {
					body = "if read -r input; then exit 1; fi\n" + body
				}
				if tc.inPath {
					r.addProgram("spacedock", body)
				} else if err := os.WriteFile(filepath.Join(r.bin, "spacedock"), []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "STUB_ROUTE=ok", "CI="+tc.ci, "SPACEDOCK_INTERNAL="+tc.internal)
			if code != 0 {
				t.Fatalf("exit=%d:\n%s", code, stderr)
			}
			requests, err := os.ReadFile(log + ".requests")
			if err != nil {
				t.Fatal(err)
			}
			fields := strings.Split(strings.TrimSpace(string(requests)), "|")
			if len(fields) != 5 || fields[1] != tc.want {
				t.Errorf("requests=%q, want User-Agent %q", requests, tc.want)
			}
		})
	}
}

func lookupRun(t *testing.T) (*extrasRun, string) {
	t.Helper()
	r := newExtrasRun(t, "1.0.1")
	osName, arch := goosArch(t)
	return r, r.stubCurl(osName, arch)
}

func TestInstallStableResolvesEachRepoFromTheRedirectWithoutTheAPI(t *testing.T) {
	r, log := lookupRun(t)
	r.addProgram("claude", "exit 0")

	stderr, code := r.runWith("SPACEDOCK_CHANNEL=stable", "STUB_ROUTE=ok")
	urls, _ := os.ReadFile(log)
	if code != 0 || !strings.Contains(stderr, "installed alpha v0.16.0") {
		t.Fatalf("exit=%d; alpha was not installed from the redirect:\n%s", code, stderr)
	}
	if strings.Contains(string(urls), "api.github.com") || strings.Contains(string(urls), "/releases/latest") || !strings.Contains(string(urls), "https://spacedock.md/latest-release/alpha\n") {
		t.Errorf("stable must use the route without direct GitHub fallback; urls:\n%s", urls)
	}
}

func TestInstallResolvesTheNewestTagByVersionWhenGitHubDoesNotAnswer(t *testing.T) {
	cases := []struct {
		name, channel, tags, want, api string
	}{
		{"a release outranks its own prereleases", "edge", "v0.16.0 v0.16.0-beta.5 v0.16.0-beta.1", "v0.16.0", "fail"},
		{"prerelease numbers compare as numbers", "edge", "v0.28.0-pre9 v0.28.0-pre10 v0.27.3", "v0.28.0-pre10", "fail"},
		{"minor versions compare as numbers and non-version tags are ignored", "edge", "v1.9.0 v1.10.0 v1.10.0-rc.1 archive/v0", "v1.10.0", "fail"},
		{"null API tag falls back to git", "edge", "v0.16.0", "v0.16.0", "null"},
		{"stable skips prereleases", "stable", "v0.28.0-pre5 v0.27.3 v0.9.9 v0.27.10", "v0.27.10", "fail"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := lookupRun(t)
			stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "SPACEDOCK_CHANNEL="+tc.channel,
				"STUB_ROUTE=malformed", "STUB_API="+tc.api, "STUB_REDIRECT=norelease", "STUB_TAGS="+tc.tags)
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

	stderr, code := r.runWith("SPACEDOCK_INSTALL_FROM="+r.dist, "SPACEDOCK_EXTRAS=1", "STUB_FAIL_REPO=spacedock-dev/alpha")
	if code != 1 {
		t.Errorf("exit = %d, want 1 when an extension release cannot be looked up", code)
	}
	for _, want := range []string{"rate limit", "spacedock-dev/alpha", "installed beta v1.0.1", "installed spacedock"} {
		if !strings.Contains(stderr, want) {
			t.Errorf("stderr lacks %q:\n%s", want, stderr)
		}
	}
	if !strings.Contains(r.calls(), "install --host claude") {
		t.Errorf("the plugin step did not run after the failed lookup: %q", r.calls())
	}
	if _, err := os.Stat(filepath.Join(r.bin, "alpha-tui")); err == nil {
		t.Errorf("alpha was installed although its release could not be looked up")
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

// Removing the platform defaults makes an isolated lookup fail under set -u.
func TestInstallLookupUserAgentWithoutPlatform(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("..", "..", "install.sh"))
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("sh", "-c", strings.TrimSuffix(string(src), "main \"$@\"\n")+"\nunset os arch\nlookup_user_agent\n")
	cmd.Env = []string{"PATH=/usr/bin:/bin", "HOME=" + t.TempDir(), "SPACEDOCK_INSTALL_DIR=" + t.TempDir()}
	out, err := cmd.CombinedOutput()
	if err != nil || string(out) != "spacedock-install/unknown (unknown; unknown)\n" {
		t.Fatalf("isolated lookup=%q (%v)", out, err)
	}
}
