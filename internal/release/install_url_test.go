package release

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestInstallScriptResolvesLiveReleaseAsset locks AC-3: install.sh's production
// path (OS/arch detection + GitHub latest-release resolution + asset-URL
// construction) is exercised against the LIVE GitHub API, and the URL it would
// fetch is asserted against the release's actually-published asset for this host.
//
// This is the installer's ONLY network path; AC-2's local-dist override isolates
// the extract/verify/install logic from network flakiness, and this test isolates
// the network URL construction here. The test drives install.sh itself (via its
// SPACEDOCK_PRINT_TARGET inspection mode) rather than re-deriving the mapping in
// Go, so it proves the real installer's logic — no duplicate-and-drift.
//
// When this host's os/arch tarball is not yet published in the live release
// (the linux gap this task closes ships before linux assets exist), the test
// still asserts the construction is correct against goreleaser's naming for the
// resolved tag, and confirms the constructed URL points at that tag's asset
// namespace. When the asset IS published (darwin today), it asserts an exact
// match against the live browser_download_url. The test skips when the GitHub
// API is unreachable so transient network failure does not red the suite.
func TestInstallScriptResolvesLiveReleaseAsset(t *testing.T) {
	rel, ok := fetchLatestRelease(t)
	if !ok {
		t.Skip("github releases API unreachable; skipping live AC-3 URL check")
	}

	got := runInstallPrintTarget(t)

	wantOS, wantArch := goosArch(t)
	if got["os"] != wantOS {
		t.Errorf("install.sh detected os=%q, want %q for this host", got["os"], wantOS)
	}
	if got["arch"] != wantArch {
		t.Errorf("install.sh detected arch=%q, want %q for this host", got["arch"], wantArch)
	}

	ver := strings.TrimPrefix(rel.TagName, "v")
	wantAsset := "spacedock_" + ver + "_" + wantOS + "_" + wantArch + ".tar.gz"
	if got["asset"] != wantAsset {
		t.Errorf("constructed asset name = %q, want %q (goreleaser template for live tag %s)",
			got["asset"], wantAsset, rel.TagName)
	}
	wantURL := "https://github.com/spacedock-dev/spacedock/releases/download/" + rel.TagName + "/" + wantAsset
	if got["tarball"] != wantURL {
		t.Errorf("constructed tarball URL = %q, want %q", got["tarball"], wantURL)
	}

	// When the live release publishes this host's asset, the constructed URL must
	// match its browser_download_url byte-for-byte. (Pre-ship, a linux host's
	// asset is absent; the construction assertions above still bind it.)
	for _, a := range rel.Assets {
		if a.Name == wantAsset {
			if a.BrowserDownloadURL != got["tarball"] {
				t.Errorf("constructed URL %q != live published browser_download_url %q",
					got["tarball"], a.BrowserDownloadURL)
			}
			return
		}
	}
	t.Logf("note: %s is not yet published in live release %s (expected pre-ship for this host's os/arch); construction asserted against goreleaser naming",
		wantAsset, rel.TagName)
}

// goosArch maps the Go runtime's GOOS/GOARCH to the goreleaser goos/goarch
// tokens install.sh derives from uname — the independent oracle the test checks
// install.sh's own detection against.
func goosArch(t *testing.T) (string, string) {
	t.Helper()
	os := runtime.GOOS
	switch os {
	case "darwin", "linux":
	default:
		t.Skipf("unsupported host OS %q for install.sh", os)
	}
	var arch string
	switch runtime.GOARCH {
	case "amd64":
		arch = "amd64"
	case "arm64":
		arch = "arm64"
	default:
		t.Skipf("unsupported host arch %q for install.sh", runtime.GOARCH)
	}
	return os, arch
}

// runInstallPrintTarget runs install.sh in its inspection mode (no download, no
// install) and parses the `key=value` lines it prints into a map.
func runInstallPrintTarget(t *testing.T) map[string]string {
	t.Helper()
	script := filepath.Join("..", "..", "install.sh")
	cmd := exec.Command("sh", script)
	cmd.Env = append([]string{"SPACEDOCK_PRINT_TARGET=1"}, scrubInstallEnv()...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("install.sh SPACEDOCK_PRINT_TARGET=1 failed: %v\n%s", err, out)
	}
	fields := map[string]string{}
	for _, line := range strings.Split(string(out), "\n") {
		if k, v, ok := strings.Cut(line, "="); ok {
			fields[k] = v
		}
	}
	for _, k := range []string{"os", "arch", "asset", "tarball"} {
		if fields[k] == "" {
			t.Fatalf("install.sh print-target output missing %q:\n%s", k, out)
		}
	}
	return fields
}

// scrubInstallEnv returns the process env with the install overrides removed so
// the print-target run takes the PRODUCTION GitHub-resolution path, not a
// local-dist override that a parent test or shell might have set.
func scrubInstallEnv() []string {
	var env []string
	for _, kv := range os.Environ() {
		switch {
		case strings.HasPrefix(kv, "SPACEDOCK_INSTALL_FROM="),
			strings.HasPrefix(kv, "SPACEDOCK_INSTALL_VERSION="),
			strings.HasPrefix(kv, "SPACEDOCK_PRINT_TARGET="):
			continue
		}
		env = append(env, kv)
	}
	return env
}

type ghRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// fetchLatestRelease pulls the live latest-release JSON. ok=false on any network
// or decode failure so the caller skips rather than reds on transient flakiness.
func fetchLatestRelease(t *testing.T) (ghRelease, bool) {
	t.Helper()
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/spacedock-dev/spacedock/releases/latest")
	if err != nil {
		return ghRelease{}, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ghRelease{}, false
	}
	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return ghRelease{}, false
	}
	if rel.TagName == "" {
		return ghRelease{}, false
	}
	return rel, true
}
