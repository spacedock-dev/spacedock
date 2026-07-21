// ABOUTME: Behavior fixtures for the version-bearing doctor messages — mismatch
// ABOUTME: and OK paths show plugin + binary display versions, never contract jargon.
package contract

import (
	"bytes"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// binaryVersionForTest is the binary display version threaded through RunDoctor
// in these fixtures. It lives outside the message under test, so a check that the
// message contains it tracks behavior (the version reached the message) rather
// than the message's spelling.
const binaryVersionForTest = "0.19.4"

// contractTokenPattern matches the internal "contract <N>" jargon the rewrite
// drops from the user-facing mismatch/OK lines. A paraphrase that re-introduces
// it FAILS these assertions.
var contractTokenPattern = regexp.MustCompile(`contract \d`)

// hasHalfOpenRange reports whether s carries a ">=N,<M" half-open range token
// (the other piece of contract jargon dropped from the user-facing lines).
func hasHalfOpenRange(s string) bool {
	return strings.Contains(s, ">=") && strings.Contains(s, ",<")
}

// TestMismatchShowsVersionsNotContract is AC-1: a too-old-binary and a
// too-old-plugin manifest exit 1 and emit a message that contains BOTH the
// plugin's display version (from the fixture) AND the binary display version
// (passed in), and contains NEITHER a "contract <N>" token NOR a ">=N,<M" range
// in the user-facing line. The oracle (the two versions) lives outside the
// message file, so this is not a tautology.
func TestMismatchShowsVersionsNotContract(t *testing.T) {
	cases := []struct {
		name          string
		manifest      string
		pluginVersion string // the fixture's display version
	}{
		{"too-old-binary", "too-old-binary.json", "0.20.0"},
		{"too-old-plugin", "too-old-plugin.json", "0.10.0"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			manifestPath := filepath.Join("testdata", c.manifest)
			code := RunDoctor(manifestPath, "claude", binaryVersionForTest, InstallSource{}, &stdout, &stderr)
			if code != 1 {
				t.Fatalf("exit = %d, want 1 (stderr=%q)", code, stderr.String())
			}
			out := stderr.String()
			if !strings.Contains(out, c.pluginVersion) {
				t.Fatalf("mismatch message must name the plugin version %q: %q", c.pluginVersion, out)
			}
			if !strings.Contains(out, binaryVersionForTest) {
				t.Fatalf("mismatch message must name the binary version %q: %q", binaryVersionForTest, out)
			}
			if contractTokenPattern.MatchString(out) {
				t.Fatalf("mismatch message must NOT carry a contract-N token: %q", out)
			}
			if hasHalfOpenRange(out) {
				t.Fatalf("mismatch message must NOT carry a >=N,<M range: %q", out)
			}
		})
	}
}

// TestCompatibleShowsVersions is the AC-2 OK half: the compatible fixture exits 0
// and stdout names both the plugin version and the binary version and carries no
// contract-N token.
func TestCompatibleShowsVersions(t *testing.T) {
	var stdout, stderr bytes.Buffer
	manifestPath := filepath.Join("testdata", "compatible.json")
	code := RunDoctor(manifestPath, "claude", binaryVersionForTest, InstallSource{}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "0.19.8") {
		t.Fatalf("OK message must name the plugin version %q: %q", "0.19.8", out)
	}
	if !strings.Contains(out, binaryVersionForTest) {
		t.Fatalf("OK message must name the binary version %q: %q", binaryVersionForTest, out)
	}
	if contractTokenPattern.MatchString(out) {
		t.Fatalf("OK message must NOT carry a contract-N token: %q", out)
	}
}

// TestTooOldBinaryRemedyLeadsWithBrew is the AC-2 remedy half: the too-old-binary
// remedy leads with `brew upgrade spacedock` and keeps the binary-vs-plugin
// distinction (`spacedock install` refreshes the plugin instead).
func TestTooOldBinaryRemedyLeadsWithBrew(t *testing.T) {
	var stdout, stderr bytes.Buffer
	manifestPath := filepath.Join("testdata", "too-old-binary.json")
	code := RunDoctor(manifestPath, "claude", binaryVersionForTest, InstallSource{}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}
	out := stderr.String()
	if !strings.Contains(out, "brew upgrade spacedock") {
		t.Fatalf("too-old-binary remedy must lead with `brew upgrade spacedock`: %q", out)
	}
	if !strings.Contains(out, "spacedock install") {
		t.Fatalf("too-old-binary remedy must name the plugin-refresh distinction `spacedock install`: %q", out)
	}
	if strings.Contains(out, "@next") {
		t.Fatalf("too-old-binary remedy must not pin a channel branch (the removed @branch shorthand): %q", out)
	}
}

// TestTooOldBinaryRemedyPerSource drives the remedy generator over the
// install-source matrix, asserting each source names the instruction that
// actually applies (AC-1) and that the edge/non-brew/host-only arms carry NO bare
// stable `brew upgrade spacedock` command (AC-2, word-boundary matched so
// `brew upgrade spacedock@next` does not count). The SourceUnknown arm is pinned
// byte-for-byte (AC-4). The oracle — the source Kind — lives outside the message,
// so a generator that ignored its input fails.
func TestTooOldBinaryRemedyPerSource(t *testing.T) {
	// bareStable matches the stable command only when `spacedock` is the whole
	// token — `brew upgrade spacedock@next` must NOT satisfy it.
	bareStable := regexp.MustCompile(`brew upgrade spacedock(\s|$)`)

	t.Run("brew-edge names the @next formula, no bare stable", func(t *testing.T) {
		got := tooOldBinaryRemedy(InstallSource{Kind: BrewEdge})
		if !strings.Contains(got, "brew upgrade spacedock@next") {
			t.Fatalf("edge remedy must name `brew upgrade spacedock@next`: %q", got)
		}
		if bareStable.MatchString(got) {
			t.Fatalf("edge remedy must NOT carry the bare stable `brew upgrade spacedock`: %q", got)
		}
		if !strings.Contains(got, "spacedock install") {
			t.Fatalf("edge remedy must keep the plugin-refresh line: %q", got)
		}
	})

	t.Run("non-brew rebuilds from source, no brew line", func(t *testing.T) {
		got := tooOldBinaryRemedy(InstallSource{Kind: NonBrew})
		if !strings.Contains(got, "go build -o spacedock ./cmd/spacedock") {
			t.Fatalf("non-brew remedy must name a source rebuild: %q", got)
		}
		if strings.Contains(got, "brew upgrade") {
			t.Fatalf("non-brew remedy must carry NO brew line: %q", got)
		}
	})

	t.Run("host-only tells the user to upgrade on the host", func(t *testing.T) {
		got := tooOldBinaryRemedy(InstallSource{Kind: BrewEdge, HostOnly: true})
		if !strings.Contains(got, "on your host") {
			t.Fatalf("host-only remedy must tell the user to upgrade on the host: %q", got)
		}
		if !strings.Contains(got, "brew upgrade spacedock@next") {
			t.Fatalf("host-only edge remedy must still name the @next formula (run on host): %q", got)
		}
		if bareStable.MatchString(got) {
			t.Fatalf("host-only edge remedy must NOT carry the bare stable command: %q", got)
		}
	})

	t.Run("source-unknown reproduces the generic block byte-for-byte", func(t *testing.T) {
		want := "  Upgrade via Homebrew: brew upgrade spacedock\n" +
			"  Or build from source: go build -o spacedock ./cmd/spacedock\n" +
			"  Or refresh the plugin instead: spacedock install"
		if got := tooOldBinaryRemedy(InstallSource{}); got != want {
			t.Fatalf("SourceUnknown remedy changed:\n got=%q\nwant=%q", got, want)
		}
	})
}
