package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func lookupRun(t *testing.T) (*extrasRun, string) {
	t.Helper()
	r := newExtrasRun(t, "1.0.1")
	osName, arch := goosArch(t)
	return r, r.stubCurl(osName, arch)
}

func TestInstallStableResolvesEachRepoFromTheRedirectWithoutTheAPI(t *testing.T) {
	r, log := lookupRun(t)
	r.addProgram("claude", "exit 0")

	stderr, code := r.runWith("SPACEDOCK_CHANNEL=stable")
	urls, _ := os.ReadFile(log)
	if code != 0 || !strings.Contains(stderr, "installed alpha v0.16.0") {
		t.Fatalf("exit=%d; alpha was not installed from the redirect:\n%s", code, stderr)
	}
	if strings.Contains(string(urls), "api.github.com") || !strings.Contains(string(urls), "https://github.com/spacedock-dev/alpha/releases/latest") {
		t.Errorf("stable must read the redirect and make no API request; urls:\n%s", urls)
	}
}

func TestInstallResolvesTheNewestTagByVersionWhenGitHubDoesNotAnswer(t *testing.T) {
	cases := []struct {
		name, channel, tags, want string
	}{
		{"a release outranks its own prereleases", "edge", "v0.16.0 v0.16.0-beta.5 v0.16.0-beta.1", "v0.16.0"},
		{"prerelease numbers compare as numbers", "edge", "v0.28.0-pre9 v0.28.0-pre10 v0.27.3", "v0.28.0-pre10"},
		{"minor versions compare as numbers and non-version tags are ignored", "edge", "v1.9.0 v1.10.0 v1.10.0-rc.1 archive/v0", "v1.10.0"},
		{"stable skips prereleases", "stable", "v0.28.0-pre5 v0.27.3 v0.9.9 v0.27.10", "v0.27.10"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := lookupRun(t)
			stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "SPACEDOCK_CHANNEL="+tc.channel,
				"STUB_API=fail", "STUB_REDIRECT=norelease", "STUB_TAGS="+tc.tags)
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
