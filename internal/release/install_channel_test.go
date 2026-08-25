// ABOUTME: Pins install.sh's SPACEDOCK_CHANNEL branches — the asset each channel
// ABOUTME: selects, the resolved-tag disclosure, and the typo abort.
package release

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func installScript() string { return filepath.Join("..", "..", "install.sh") }

// TestInstallChannelEdgeInstallsTheEdgeAsset locks AC-1's mechanism (and AC-4's
// converse) offline: from ONE dist holding BOTH channels' archives, the edge
// channel installs the edge binary. The two payloads print different markers, so
// this asserts which artifact actually landed and ran — not which name the
// script printed. The default channel's half of the choice is the pre-existing
// happy-path test, which now runs against the same both-assets dist.
func TestInstallChannelEdgeInstallsTheEdgeAsset(t *testing.T) {
	fx := buildInstallFixture(t)
	installDir, code := runInstall(t, installScript(), fx.dir, "edge")
	if code != 0 {
		t.Fatalf("edge install exited %d, want 0", code)
	}
	out, err := exec.Command(filepath.Join(installDir, "spacedock")).CombinedOutput()
	if err != nil {
		t.Fatalf("installed spacedock did not run: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), fx.edgeMarker) {
		t.Errorf("installed binary printed %q, want the edge payload %q", out, fx.edgeMarker)
	}
	if strings.Contains(string(out), fx.marker) {
		t.Errorf("installed binary printed %q — the STABLE payload, from an edge run", out)
	}
}

// TestInstallChannelEdgeTamperInstallsNothing locks AC-5: the edge asset flows
// through the same fail-closed checksum gate. The `_edge` tarball is swapped for
// a valid-but-different archive after checksums.txt is written, so only the hash
// mismatch can reject it — and the abort names that asset, which also proves the
// run had SELECTED the edge archive rather than failing to find one.
func TestInstallChannelEdgeTamperInstallsNothing(t *testing.T) {
	fx := buildInstallFixture(t)
	writeTarGz(t, fx.edgePath, "spacedock", []byte("#!/bin/sh\necho "+fx.edgeMarker+"-tampered\n"))
	if sha256OfFile(t, fx.edgePath) == fx.edgeChecksum {
		t.Fatal("tampered _edge tarball hash still equals the recorded checksum; the swap changed nothing")
	}

	installDir, code, _, stderr := runInstallCapture(t, installScript(), fx.dir, "edge")
	if code == 0 {
		t.Fatal("install.sh accepted a tampered _edge tarball (exit 0); the edge channel is not behind the fail-closed gate")
	}
	if want := "checksum mismatch for " + fx.edgeAsset; !strings.Contains(stderr, want) {
		t.Errorf("edge run aborted with %q, want %q — the rejection must come from the checksum gate", stderr, want)
	}
	if _, err := os.Stat(filepath.Join(installDir, "spacedock")); err == nil {
		t.Error("install.sh installed a binary despite the _edge checksum mismatch")
	}
}

// TestInstallChannelRejectsUnknownValue locks AC-6. The dist holds both assets,
// so a defaulting `case *)` would install stable and exit 0 — this abort is the
// only thing between a typo and silently installing the other channel.
func TestInstallChannelRejectsUnknownValue(t *testing.T) {
	fx := buildInstallFixture(t)
	installDir, code, _, stderr := runInstallCapture(t, installScript(), fx.dir, "egde")
	if code == 0 {
		t.Fatal("install.sh accepted SPACEDOCK_CHANNEL=egde (exit 0); a typo silently picks a channel")
	}
	for _, want := range []string{"egde", "stable", "edge"} {
		if !strings.Contains(stderr, want) {
			t.Errorf("abort message %q does not name %q", stderr, want)
		}
	}
	if _, err := os.Stat(installDir); err == nil {
		t.Error("install.sh created the install dir on an unknown channel")
	}
}

// TestInstallChannelPrintTarget locks AC-4: on the default and explicit-stable
// runs install.sh prints the EXACT five-line target the pre-channel script
// printed for the same dist, and edge differs only in the asset it names. The
// whole-stdout match also catches the disclosure line landing on stdout, which
// would break the `key=value` parser install_url_test.go depends on.
func TestInstallChannelPrintTarget(t *testing.T) {
	goos, arch := goosArch(t)
	fx := buildInstallFixture(t)
	for _, tc := range []struct{ name, channel, asset string }{
		{"unset channel names the unsuffixed asset", "", fx.asset},
		{"explicit stable names the unsuffixed asset", "stable", fx.asset},
		{"edge names the _edge asset", "edge", fx.edgeAsset},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, code := runInstallSh(t, installScript(), tc.channel,
				"SPACEDOCK_PRINT_TARGET=1", "SPACEDOCK_INSTALL_FROM="+fx.dir)
			if code != 0 {
				t.Fatalf("print-target exited %d: %s", code, stderr)
			}
			want := fmt.Sprintf("os=%s\narch=%s\nasset=%s\ntarball=%s\nchecksums=%s\n",
				goos, arch, tc.asset, filepath.Join(fx.dir, tc.asset), filepath.Join(fx.dir, "checksums.txt"))
			if stdout != want {
				t.Errorf("print-target stdout:\n%q\nwant:\n%q", stdout, want)
			}
			// The local-dist path resolves no version, so it discloses none.
			if strings.Contains(stderr, "resolved") {
				t.Errorf("local-dist run disclosed a resolved tag: %q", stderr)
			}
		})
	}
}

// TestInstallChannelURLBaseAssetName pins asset_name's channel branch OFFLINE.
// The local-dist source globs a directory and never calls asset_name, so without
// this the naming the live release path depends on would be held up only by the
// network test below, which skips on an unreachable API. The URL-base source
// builds the name from the pinned version, and print-target returns before any
// download — so the unreachable host is never contacted.
func TestInstallChannelURLBaseAssetName(t *testing.T) {
	goos, arch := goosArch(t)
	const base, ver = "https://example.invalid/dl", "1.2.3"
	for _, tc := range []struct{ name, channel, asset string }{
		{"unset channel builds the unsuffixed name", "", "spacedock_" + ver + "_" + goos + "_" + arch + ".tar.gz"},
		{"edge builds the _edge name", "edge", "spacedock_" + ver + "_" + goos + "_" + arch + "_edge.tar.gz"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, code := runInstallSh(t, installScript(), tc.channel,
				"SPACEDOCK_PRINT_TARGET=1", "SPACEDOCK_INSTALL_FROM="+base, "SPACEDOCK_INSTALL_VERSION="+ver)
			if code != 0 {
				t.Fatalf("print-target exited %d: %s", code, stderr)
			}
			want := fmt.Sprintf("os=%s\narch=%s\nasset=%s\ntarball=%s/%s\nchecksums=%s/checksums.txt\n",
				goos, arch, tc.asset, base, tc.asset, base)
			if stdout != want {
				t.Errorf("print-target stdout:\n%q\nwant:\n%q", stdout, want)
			}
		})
	}
}

// TestInstallChannelLiveEdgeResolution locks AC-3 and AC-2's production half
// against the live GitHub API: edge resolves the newest release INCLUDING
// prereleases, both channels disclose the tag they resolved before any download,
// and the constructed URL is one the release actually publishes. The
// byte-for-byte match against browser_download_url is what catches an unsuffixed
// asset name, which 404s on a prerelease; a name-shape check would not. Skips —
// never reds — when the API is unreachable or rate-limited, the pattern
// install_url_test.go already uses (a 403 produces the same resolution failure).
func TestInstallChannelLiveEdgeResolution(t *testing.T) {
	rel, ok := fetchNewestRelease(t)
	if !ok {
		t.Skip("github releases list API unreachable; skipping live edge resolution check")
	}
	goos, arch := goosArch(t)

	stdout, stderr, code := runInstallSh(t, installScript(), "edge", "SPACEDOCK_PRINT_TARGET=1")
	if code != 0 {
		t.Skipf("install.sh could not reach the GitHub releases API:\n%s", stderr)
	}
	if want := "install.sh: resolved edge channel to " + rel.TagName + "\n"; !strings.Contains(stderr, want) {
		t.Errorf("edge stderr = %q, want it to carry %q before any download", stderr, want)
	}
	if _, stableErr, stableCode := runInstallSh(t, installScript(), "", "SPACEDOCK_PRINT_TARGET=1"); stableCode == 0 &&
		!strings.Contains(stableErr, "install.sh: resolved stable channel to v") {
		t.Errorf("default-channel stderr = %q, want it to disclose the resolved stable tag", stableErr)
	}

	asset := "spacedock_" + strings.TrimPrefix(rel.TagName, "v") + "_" + goos + "_" + arch + "_edge.tar.gz"
	base := "https://github.com/spacedock-dev/spacedock/releases/download/" + rel.TagName
	want := fmt.Sprintf("os=%s\narch=%s\nasset=%s\ntarball=%s/%s\nchecksums=%s/checksums.txt\n",
		goos, arch, asset, base, asset, base)
	if stdout != want {
		t.Errorf("edge print-target stdout:\n%q\nwant:\n%q", stdout, want)
	}
	for _, a := range rel.Assets {
		if a.Name == asset {
			if a.BrowserDownloadURL != base+"/"+asset {
				t.Errorf("constructed URL %q != published browser_download_url %q", base+"/"+asset, a.BrowserDownloadURL)
			}
			return
		}
	}
	t.Errorf("newest release %s publishes no %s, so the edge channel would 404 on this host", rel.TagName, asset)
}

// fetchNewestRelease pulls the newest release INCLUDING prereleases — the
// endpoint the edge channel resolves, and the independent oracle for what it
// should have resolved. ok=false on any network, status, or decode failure so
// the caller skips rather than reds on flakiness or a rate limit.
func fetchNewestRelease(t *testing.T) (ghRelease, bool) {
	t.Helper()
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/spacedock-dev/spacedock/releases?per_page=1")
	if err != nil {
		return ghRelease{}, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ghRelease{}, false
	}
	var rels []ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rels); err != nil || len(rels) == 0 || rels[0].TagName == "" {
		return ghRelease{}, false
	}
	return rels[0], true
}

// runInstallCapture installs from a local dist into a fresh dir on the given
// channel, returning that dir plus the run's exit code and output streams.
func runInstallCapture(t *testing.T, script, distDir, channel string) (installDir string, code int, stdout, stderr string) {
	t.Helper()
	installDir = filepath.Join(t.TempDir(), "bin")
	stdout, stderr, code = runInstallSh(t, script, channel,
		"SPACEDOCK_INSTALL_FROM="+distDir, "SPACEDOCK_INSTALL_DIR="+installDir)
	return installDir, code, stdout, stderr
}

// runInstallSh runs install.sh with an explicit env — the process env minus every
// install override (SPACEDOCK_CHANNEL included, which scrubInstallEnv predates),
// plus the requested channel and extras. An empty channel leaves the variable
// UNSET, which is the default-channel case a leaked env var would otherwise turn
// into an edge run. stdout and stderr come back separately: the resolved-tag
// disclosure is stderr-only by design, and stdout must stay machine-parseable.
func runInstallSh(t *testing.T, script, channel string, extra ...string) (stdout, stderr string, code int) {
	t.Helper()
	env := append([]string{}, extra...)
	for _, kv := range scrubInstallEnv() {
		if !strings.HasPrefix(kv, "SPACEDOCK_CHANNEL=") {
			env = append(env, kv)
		}
	}
	if channel != "" {
		env = append(env, "SPACEDOCK_CHANNEL="+channel)
	}
	cmd := exec.Command("sh", script)
	cmd.Env = env
	var out, errOut bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errOut
	err := cmd.Run()
	t.Logf("install.sh (%s, channel=%q)\nstdout:\n%sstderr:\n%s", filepath.Base(script), channel, out.String(), errOut.String())
	var exitErr *exec.ExitError
	if err != nil && !errors.As(err, &exitErr) {
		t.Fatalf("install.sh failed to launch: %v", err)
	}
	return out.String(), errOut.String(), cmd.ProcessState.ExitCode()
}
