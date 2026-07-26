// ABOUTME: AC-5 dev-build tests — the embed fallback unit, and a behavior
// ABOUTME: fixture that `go build`s the real binary and reads its --version.
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

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
// two branches: a marked release Version passes through unchanged; the
// unstamped `dev` sentinel resolves to the embedded checkout version plus +dev.
func TestDisplayVersionFallsBackToEmbedOnlyWhenUnstamped(t *testing.T) {
	t.Run("stamped version passes through", func(t *testing.T) {
		withVersion(t, "0.19.4")
		if got := displayVersion(); got != "0.19.4" {
			t.Fatalf("displayVersion() = %q, want the stamped Version unchanged", got)
		}
	})

	t.Run("unstamped dev resolves to checkout version + dev", func(t *testing.T) {
		withBuildIdentity(t, "dev", "false")
		want := checkoutManifestVersion(t) + "+dev"
		if got := displayVersion(); got != want {
			t.Fatalf("displayVersion() = %q, want %q", got, want)
		}
	})
}

// TestUnmarkedProvenanceDoesNotChangeCompatibilityIdentity proves AC-3 across
// tag, describe, revision, and dirty candidates.
func TestUnmarkedProvenanceDoesNotChangeCompatibilityIdentity(t *testing.T) {
	wantVersion := checkoutManifestVersion(t) + "+dev"
	manifest := filepath.Join(repoRootForDevBuild(t), ".claude-plugin", "plugin.json")
	for _, version := range []string{
		"0.27.0",
		"v0.27.0-pre0-17-gabcdef0",
		"abcdef012345",
		"0.27.0-5-gabcdef0-dirty",
	} {
		t.Run(version, func(t *testing.T) {
			withBuildIdentity(t, version, "false")
			if got := displayVersion(); got != wantVersion {
				t.Errorf("displayVersion() = %q, want invariant %q", got, wantVersion)
			}
			var doctor bytes.Buffer
			if code := contract.RunDoctor(manifest, "claude", displayVersion(), false, &doctor, &doctor); code != 0 {
				t.Errorf("doctor exit = %d, want 0", code)
			}
			wantDoctor := fmt.Sprintf("OK: spacedock binary %s and plugin %s are compatible.\n",
				wantVersion, checkoutManifestVersion(t))
			if got := doctor.String(); got != wantDoctor {
				t.Errorf("doctor output = %q, want %q", got, wantDoctor)
			}
		})
	}
}

func withBuildIdentity(t *testing.T, version, marker string) {
	t.Helper()
	origVersion, origMarker := Version, releaseBuild
	Version, releaseBuild = version, marker
	t.Cleanup(func() { Version, releaseBuild = origVersion, origMarker })
}

// TestSourceBuildCompatibilityIdentity is the AC-1/AC-2 real-binary proof:
// copied provenance is inert; an exact marked release remains load-bearing.
func TestSourceBuildCompatibilityIdentity(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go not on PATH; dev-build behavior fixture requires the toolchain")
	}

	repo := repoRootForDevBuild(t)
	manifestVersion := checkoutManifestVersion(t)
	major, minor, ok := contract.ParseMajorMinor(manifestVersion)
	if !ok {
		t.Fatalf("checkout manifest version %q has no parseable major.minor", manifestVersion)
	}
	sourceVersion := manifestVersion + "+dev"
	misleadingVersion := fmt.Sprintf("v%d.%d.0-pre0-17-gabcdef0", major, minor+1)
	releaseVersion := fmt.Sprintf("%d.%d.0-pre0", major, minor+1)

	for _, tt := range []struct {
		name           string
		ldflags        string
		wantVersion    string
		wantDoctorExit int
	}{
		{name: "plain source", wantVersion: sourceVersion},
		{
			name: "misleading future-minor version without release marker",
			ldflags: "-X github.com/spacedock-dev/spacedock/internal/cli.Version=" +
				misleadingVersion,
			wantVersion: sourceVersion,
		},
		{
			name: "exact marked future-minor release",
			ldflags: "-X github.com/spacedock-dev/spacedock/internal/cli.Version=" +
				releaseVersion + " -X github.com/spacedock-dev/spacedock/internal/cli.releaseBuild=true",
			wantVersion:    releaseVersion,
			wantDoctorExit: 1,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bin := filepath.Join(t.TempDir(), "spacedock")
			args := []string{"build", "-o", bin, "-ldflags", tt.ldflags, "./cmd/spacedock"}
			build := exec.Command(goBin, args...)
			build.Dir = repo
			if out, err := build.CombinedOutput(); err != nil {
				t.Fatalf("go %s: %v\n%s", strings.Join(args, " "), err, out)
			}

			out, err := exec.Command(bin, "--version").Output()
			if err != nil {
				t.Fatalf("%s --version: %v", bin, err)
			}
			firstLine := strings.SplitN(string(out), "\n", 2)[0]
			wantLine := "spacedock " + tt.wantVersion + " (contract 3)"
			if firstLine != wantLine {
				t.Errorf("--version line 1 = %q, want %q", firstLine, wantLine)
			}

			doctor := exec.Command(bin, "doctor", "--plugin-manifest",
				filepath.Join(repo, ".claude-plugin", "plugin.json"))
			doctorOut, doctorErr := doctor.CombinedOutput()
			if doctorErr != nil && doctor.ProcessState == nil {
				t.Fatalf("%s doctor: %v\n%s", bin, doctorErr, doctorOut)
			}
			if got := doctor.ProcessState.ExitCode(); got != tt.wantDoctorExit {
				t.Errorf("doctor exit = %d, want %d\n%s", got, tt.wantDoctorExit, doctorOut)
			}
			wantDoctor := fmt.Sprintf("OK: spacedock binary %s and plugin %s are compatible.",
				tt.wantVersion, manifestVersion)
			if tt.wantDoctorExit != 0 {
				wantDoctor = fmt.Sprintf("Spacedock version mismatch: binary %s, plugin %s.",
					tt.wantVersion, manifestVersion)
			}
			if got := string(doctorOut); !strings.HasPrefix(got, wantDoctor) {
				t.Errorf("doctor output = %q, want prefix %q", got, wantDoctor)
			}
		})
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
