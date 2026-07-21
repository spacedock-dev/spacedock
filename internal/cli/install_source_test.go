// ABOUTME: Tables for install-source detection (real resolved path → source kind)
// ABOUTME: and the end-to-end doctor remedy render per source, plus a live cask smoke.
package cli

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// brewFound stubs a reachable Homebrew (`brew` resolves); brewMissing stubs the
// sandbox/minimal case where `brew` is stripped from PATH.
func brewFound(string) (string, error)   { return "/opt/homebrew/bin/brew", nil }
func brewMissing(string) (string, error) { return "", errors.New("exec: \"brew\": not found in $PATH") }

// TestDetectInstallSource is AC-3: detectInstallSource classifies the resolved
// running-binary path. A `…/Caskroom/<token>/…` path is the brew kind named by
// the token; `brew` unreachable flips a brew kind to HostOnly; a resolved
// non-Caskroom path is non-brew; an empty path is the generic fallback. The
// oracle is the input path + brew stub, both external to the code under test.
func TestDetectInstallSource(t *testing.T) {
	cases := []struct {
		name     string
		execPath string
		brew     func(string) (string, error)
		want     contract.InstallSource
	}{
		{
			name:     "edge cask, brew reachable",
			execPath: "/opt/homebrew/Caskroom/spacedock@next/0.26.0-pre0/spacedock",
			brew:     brewFound,
			want:     contract.InstallSource{Kind: contract.BrewEdge},
		},
		{
			name:     "stable cask, brew reachable",
			execPath: "/usr/local/Caskroom/spacedock/0.25.0/spacedock",
			brew:     brewFound,
			want:     contract.InstallSource{Kind: contract.BrewStable},
		},
		{
			name:     "edge cask, brew stripped (sandbox)",
			execPath: "/opt/homebrew/Caskroom/spacedock@next/0.26.0-pre0/spacedock",
			brew:     brewMissing,
			want:     contract.InstallSource{Kind: contract.BrewEdge, HostOnly: true},
		},
		{
			name:     "source checkout build",
			execPath: "/Users/x/git/spacedock/spacedock",
			brew:     brewFound,
			want:     contract.InstallSource{Kind: contract.NonBrew},
		},
		{
			name:     "unknown cask token falls back",
			execPath: "/opt/homebrew/Caskroom/spacedock@beta/0.1.0/spacedock",
			brew:     brewFound,
			want:     contract.InstallSource{Kind: contract.SourceUnknown},
		},
		{
			name:     "unresolvable path",
			execPath: "",
			brew:     brewFound,
			want:     contract.InstallSource{Kind: contract.SourceUnknown},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := detectInstallSource(c.execPath, c.brew, devBranch)
			if got != c.want {
				t.Fatalf("detectInstallSource(%q) = %+v, want %+v", c.execPath, got, c.want)
			}
		})
	}
}

// TestTooOldBinaryRemedyRendersPerSource is AC-1/AC-2 end-to-end: for each install
// source, a too-old-binary manifest driven through contract.RunDoctor emits the
// remedy that matches the source. It proves src threads the whole doctor path
// (verdict → message → remedy), and that no edge/non-brew row carries the bare
// stable command (word-boundary matched). The manifest version (0.99.0 > the
// 0.1.0 binary) forces the too-old-binary verdict; the source is the moving
// oracle, external to the message.
func TestTooOldBinaryRemedyRendersPerSource(t *testing.T) {
	dir := t.TempDir()
	manifest := filepath.Join(dir, "plugin.json")
	mustWrite(t, manifest, `{"version":"0.99.0"}`+"\n")
	const oldBinary = "0.1.0"

	// bareStable matches the stable command only when `spacedock` is the whole
	// token — `brew upgrade spacedock@next` must NOT satisfy it (AC-2).
	bareStable := regexp.MustCompile(`brew upgrade spacedock(\s|$)`)

	cases := []struct {
		name         string
		src          contract.InstallSource
		wantContains []string
		wantAbsent   []string
		wantStable   bool // whether the bare stable command should be present
	}{
		{
			name:         "brew stable",
			src:          contract.InstallSource{Kind: contract.BrewStable},
			wantContains: []string{"brew upgrade spacedock", "spacedock install"},
			wantStable:   true,
		},
		{
			name:         "brew edge",
			src:          contract.InstallSource{Kind: contract.BrewEdge},
			wantContains: []string{"brew upgrade spacedock@next", "spacedock install"},
			wantStable:   false,
		},
		{
			name:         "non-brew",
			src:          contract.InstallSource{Kind: contract.NonBrew},
			wantContains: []string{"go build -o spacedock", "spacedock install"},
			wantAbsent:   []string{"brew upgrade"},
			wantStable:   false,
		},
		{
			name:         "sandbox host-only",
			src:          contract.InstallSource{Kind: contract.BrewEdge, HostOnly: true},
			wantContains: []string{"on your host", "brew upgrade spacedock@next"},
			wantStable:   false,
		},
		// The SourceUnknown arm through RunDoctor is covered by the contract
		// package's TestTooOldBinaryRemedyLeadsWithBrew — not repeated here.
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := contract.RunDoctor(manifest, "claude", oldBinary, c.src, &stdout, &stderr)
			if code != 1 {
				t.Fatalf("exit = %d, want 1 (stderr=%q)", code, stderr.String())
			}
			out := stderr.String()
			for _, want := range c.wantContains {
				if !strings.Contains(out, want) {
					t.Fatalf("remedy for %s must contain %q: %q", c.name, want, out)
				}
			}
			for _, absent := range c.wantAbsent {
				if strings.Contains(out, absent) {
					t.Fatalf("remedy for %s must NOT contain %q: %q", c.name, absent, out)
				}
			}
			if got := bareStable.MatchString(out); got != c.wantStable {
				t.Fatalf("remedy for %s: bare stable `brew upgrade spacedock` present=%v, want %v: %q", c.name, got, c.wantStable, out)
			}
		})
	}
}

// TestLiveInstalledBinaryRemedyMatchesCask is the machine-grounded value check:
// on a box where the PATH `spacedock` resolves under a Homebrew Caskroom segment,
// the in-process detector + remedy render name that cask's own formula — so a
// `@next` box gets `brew upgrade spacedock@next`, never the bare stable command.
// It reads the REAL resolved path (not a fixture) and skips when the PATH binary
// is not a Caskroom install, keeping the suite portable per "no hidden machine
// dependencies". It grounds on the running machine yet exercises THIS code's
// detection+render in-process (not the possibly-stale installed binary).
func TestLiveInstalledBinaryRemedyMatchesCask(t *testing.T) {
	bin, err := exec.LookPath("spacedock")
	if err != nil {
		t.Skip("spacedock not on PATH; live cask smoke requires the installed binary")
	}
	resolved, err := filepath.EvalSymlinks(bin)
	if err != nil {
		t.Skipf("cannot resolve %s through symlinks: %v", bin, err)
	}
	token, isCask := caskToken(resolved)
	if !isCask {
		t.Skipf("PATH spacedock (%s) is not a Caskroom install; nothing to assert", resolved)
	}
	if token != "spacedock" && token != "spacedock@next" {
		t.Skipf("PATH spacedock cask token %q is not a known channel; nothing to assert", token)
	}
	src := detectInstallSource(resolved, exec.LookPath, devBranch)
	if src.Kind != contract.BrewStable && src.Kind != contract.BrewEdge {
		t.Fatalf("Caskroom path %s classified as %+v, want a brew kind", resolved, src)
	}

	dir := t.TempDir()
	manifest := filepath.Join(dir, "plugin.json")
	mustWrite(t, manifest, `{"version":"999.0.0"}`+"\n")
	var stdout, stderr bytes.Buffer
	if code := contract.RunDoctor(manifest, "claude", "0.1.0", src, &stdout, &stderr); code != 1 {
		t.Fatalf("doctor exit = %d, want 1 (stderr=%q)", code, stderr.String())
	}
	wantFormula := "brew upgrade " + token
	if !strings.Contains(stderr.String(), wantFormula) {
		t.Fatalf("remedy for cask token %q must name %q: %q", token, wantFormula, stderr.String())
	}
}
