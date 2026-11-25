// ABOUTME: AC-5 dev-build tests — the embed fallback unit, and a behavior
// ABOUTME: fixture that `go build`s the real binary and reads its --version.
package cli

import (
	"debug/buildinfo"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"
	"testing"

	spacedock "github.com/spacedock-dev/spacedock"
	"github.com/spacedock-dev/spacedock/internal/contract"
)

// TestEmbeddedManifestVersionReadsCheckoutVersion locks the embed fallback unit:
// embeddedManifestVersion() returns the SAME version this checkout's own
// .claude-plugin/plugin.json declares (read independently here, not assumed),
// so a future manifest bump is automatically reflected with no code change.
func TestEmbeddedManifestVersionReadsCheckoutVersion(t *testing.T) {
	want := checkoutManifestVersion(t)
	got, ok := embeddedManifestVersion()
	if !ok {
		t.Fatalf("embeddedManifestVersion() ok = false, want true")
	}
	if got != want {
		t.Fatalf("embeddedManifestVersion() = %q, want %q (this checkout's .claude-plugin/plugin.json version)", got, want)
	}
}

// TestDisplayVersionFallsBackToEmbedOnlyWhenUnstamped locks displayVersion()'s
// two branches: a stamped Version passes through unchanged (the embed is never
// consulted); the unstamped `dev` sentinel resolves to the embedded checkout
// version with a `+dev` build tag appended (D3).
func TestDisplayVersionFallsBackToEmbedOnlyWhenUnstamped(t *testing.T) {
	t.Run("stamped version passes through", func(t *testing.T) {
		withVersion(t, "0.19.4")
		if got := displayVersion(); got != "0.19.4" {
			t.Fatalf("displayVersion() = %q, want the stamped Version unchanged", got)
		}
	})

	t.Run("unstamped dev resolves to checkout version + dev", func(t *testing.T) {
		withVersion(t, "dev")
		withBuildInfo(t, nil, false)
		want := checkoutManifestVersion(t) + "+dev"
		if got := displayVersion(); got != want {
			t.Fatalf("displayVersion() = %q, want %q", got, want)
		}
	})
}

// TestUnstampedSourceBuildReportsCheckoutVersionPlusDev is the AC-5 behavior
// fixture: it `go build`s the REAL cmd/spacedock binary (no ldflags — the
// unstamped shape) and runs its `--version`, observing that line 1 reports
// the checkout version with dev metadata and available VCS revision. This is the actual command a source-build operator runs, not the
// package-internal displayVersion() call. Skipped when `go` is unavailable.
func TestUnstampedSourceBuildReportsCheckoutVersionPlusDev(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go not on PATH; dev-build behavior fixture requires the toolchain")
	}

	repo := repoRootForDevBuild(t)
	bin := filepath.Join(t.TempDir(), "spacedock-dev-build")
	build := exec.Command(goBin, "build", "-o", bin, "./cmd/spacedock")
	build.Dir = repo
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build ./cmd/spacedock: %v\n%s", err, out)
	}

	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		t.Fatalf("%s --version: %v", bin, err)
	}

	wantVersion := checkoutManifestVersion(t) + "+dev"
	info, err := buildinfo.ReadFile(bin)
	if err != nil {
		t.Fatal(err)
	}
	settings := map[string]string{}
	for _, setting := range info.Settings {
		settings[setting.Key] = setting.Value
	}
	if revision := settings["vcs.revision"]; regexp.MustCompile(`^[0-9a-fA-F]{7,}$`).MatchString(revision) {
		wantVersion += "." + strings.ToLower(revision[:7])
		if settings["vcs.modified"] == "true" {
			wantVersion += ".dirty"
		}
	}
	firstLine := strings.SplitN(string(out), "\n", 2)[0]
	wantLine := "spacedock " + wantVersion
	if firstLine != wantLine {
		t.Fatalf("unstamped build --version line 1 = %q, want %q", firstLine, wantLine)
	}
}

// checkoutManifestVersion reads the `version` field of this checkout's own
// .claude-plugin/plugin.json directly from disk — independent of the compiled
// go:embed constant, so a test comparing the two is not a tautology.
func checkoutManifestVersion(t *testing.T) string {
	t.Helper()
	path := filepath.Join(repoRootForDevBuild(t), ".claude-plugin", "plugin.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var m struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return m.Version
}

// repoRootForDevBuild resolves the module root (the parent of internal/cli/..)
// so the behavior fixture can `go build ./cmd/spacedock` from the real checkout.
func repoRootForDevBuild(t *testing.T) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func withBuildInfo(t *testing.T, info *debug.BuildInfo, ok bool) {
	t.Helper()
	old := readBuildInfo
	t.Cleanup(func() { readBuildInfo = old })
	readBuildInfo = func() (*debug.BuildInfo, bool) { return info, ok }
}

func TestDevVersionVCSMetadata(t *testing.T) {
	withVersion(t, "dev")
	manifest := checkoutManifestVersion(t)
	for _, tc := range []struct{ revision, modified, suffix string }{
		{"5210d36d0", "false", ".5210d36"},
		{"abcdef012", "false", ".abcdef0"},
		{"ABCDEF012", "true", ".abcdef0.dirty"},
		{"abcdef0", "TRUE", ".abcdef0"},
		{"", "true", ""}, {"123456", "true", ""},
		{"abcdef0z", "true", ""}, {"bad sha", "false", ""},
	} {
		t.Run(tc.revision+tc.modified, func(t *testing.T) {
			withBuildInfo(t, &debug.BuildInfo{Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: tc.revision}, {Key: "vcs.modified", Value: tc.modified},
			}}, true)
			got := displayVersion()
			if got != manifest+"+dev"+tc.suffix {
				t.Fatalf("version %q", got)
			}
			for _, version := range []string{manifest, got} {
				major, minor, ok := contract.ParseMajorMinor(version)
				if !ok || major != 0 || minor != 29 {
					t.Fatalf("compatibility %q: %d.%d %v", version, major, minor, ok)
				}
			}
		})
	}
	t.Run("missing build information", func(t *testing.T) {
		withBuildInfo(t, nil, false)
		if got := displayVersion(); got != manifest+"+dev" {
			t.Fatal(got)
		}
	})
	t.Run("malformed manifest", func(t *testing.T) {
		old := spacedock.PluginManifest
		t.Cleanup(func() { spacedock.PluginManifest = old })
		spacedock.PluginManifest = []byte("invalid")
		if got := displayVersion(); got != "dev" {
			t.Fatal(got)
		}
	})
	t.Run("stamped release bypasses build info", func(t *testing.T) {
		withVersion(t, "0.29.0")
		withBuildInfo(t, nil, false)
		readBuildInfo = func() (*debug.BuildInfo, bool) { t.Fatal("release read VCS"); return nil, false }
		if got := displayVersion(); got != "0.29.0" {
			t.Fatal(got)
		}
	})
}
