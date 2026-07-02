// ABOUTME: Unit tests for the version-compatibility axis: major.minor parsing
// ABOUTME: and the minor-exact, both-directions compatibility compare.
package contract

import (
	"strings"
	"testing"
)

// TestParseMajorMinor covers the major.minor extraction: both published
// prerelease suffix styles, a build tag, bare stable versions, and the
// integer-era `dev` sentinel (no parseable major.minor).
func TestParseMajorMinor(t *testing.T) {
	cases := []struct {
		raw       string
		wantMajor int
		wantMinor int
		wantOk    bool
	}{
		{"0.22.0", 0, 22, true},
		{"0.23.0", 0, 23, true},
		{"0.24.0-pre1", 0, 24, true},
		{"0.23.0-pre.4", 0, 23, true},
		{"0.24.0-pre1+dev", 0, 24, true},
		{"1.2", 1, 2, true},
		{"1.2.3.4", 1, 2, true}, // extra components beyond major.minor.patch ignored
		{"dev", 0, 0, false},
		{"", 0, 0, false},
		{"0", 0, 0, false},      // no minor component
		{"a.2.0", 0, 0, false},  // non-integer major
		{"0.b.0", 0, 0, false},  // non-integer minor
		{"-1.2.0", 0, 0, false}, // negative major
	}
	for _, c := range cases {
		major, minor, ok := ParseMajorMinor(c.raw)
		if ok != c.wantOk {
			t.Errorf("ParseMajorMinor(%q) ok = %v, want %v", c.raw, ok, c.wantOk)
			continue
		}
		if !ok {
			continue
		}
		if major != c.wantMajor || minor != c.wantMinor {
			t.Errorf("ParseMajorMinor(%q) = (%d,%d), want (%d,%d)", c.raw, major, minor, c.wantMajor, c.wantMinor)
		}
	}
}

// TestCompare drives the minor-exact, both-directions comparison, asserting the
// verdict class and that each non-compatible verdict's message carries its
// pinned remedy substring.
func TestCompare(t *testing.T) {
	cases := []struct {
		name          string
		pluginVersion string
		binaryVersion string
		wantVerdict   Verdict
		wantPinned    string // substring the message must contain (empty = no check)
	}{
		{"compatible-exact", "0.20.0", "0.20.0", Compatible, ""},
		{"compatible-patch-skew-binary-ahead", "0.20.0", "0.20.3", Compatible, ""},
		{"compatible-patch-skew-plugin-ahead", "0.20.3", "0.20.0", Compatible, ""},
		{"compatible-prerelease-skew", "0.20.0-pre1", "0.20.0", Compatible, ""},
		{"too-old-binary", "0.21.0", "0.20.0", TooOldBinary, "Upgrade the binary to continue."},
		{"too-old-binary-major", "1.0.0", "0.20.0", TooOldBinary, "Upgrade the binary to continue."},
		{"too-old-plugin", "0.19.0", "0.20.0", TooOldPlugin, "Update the plugin to continue."},
		{"too-old-plugin-major", "0.20.0", "1.0.0", TooOldPlugin, "Update the plugin to continue."},
		{"malformed-empty-plugin", "", "0.20.0", MalformedVersion, "malformed plugin version"},
		{"malformed-garbage-plugin", "not-a-version", "0.20.0", MalformedVersion, "malformed plugin version"},
		{"dev-binary-too-old", "0.20.0", "dev", TooOldBinary, "Upgrade the binary to continue."},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := Compare("claude", c.pluginVersion, c.binaryVersion)
			if res.Verdict != c.wantVerdict {
				t.Fatalf("Compare(plugin=%q, binary=%q) verdict = %v, want %v", c.pluginVersion, c.binaryVersion, res.Verdict, c.wantVerdict)
			}
			if c.wantPinned != "" && !strings.Contains(res.Message, c.wantPinned) {
				t.Fatalf("Compare(plugin=%q, binary=%q) message = %q, want substring %q", c.pluginVersion, c.binaryVersion, res.Message, c.wantPinned)
			}
			// The compatible verdict carries no remedy; non-compatible ones do.
			if c.wantVerdict == Compatible && res.Message == "" {
				t.Fatalf("compatible verdict should still carry an OK message")
			}
		})
	}
}

// TestCompatibleUpgradeHint locks: when the binary and plugin display versions
// are both valid semver and the binary is strictly newer while the verdict is
// Compatible (same minor, binary ahead on patch), the Compatible message carries
// an opt-in upgrade hint naming a newer plugin AND the host-specific `spacedock
// install --host {host}` command — but the verdict stays Compatible and the OK
// line is preserved. The hint NEVER fires on equal versions or a non-semver
// (`dev`) binary version. The trigger (the semver skew) comes from the version
// inputs, not the message.
func TestCompatibleUpgradeHint(t *testing.T) {
	t.Run("behind-plugin-hints", func(t *testing.T) {
		for _, host := range []string{"claude", "codex"} {
			res := Compare(host, "0.20.0", "0.20.5")
			if res.Verdict != Compatible {
				t.Fatalf("host %s: verdict = %v, want Compatible (the hint must not change the verdict)", host, res.Verdict)
			}
			if !strings.Contains(res.Message, "OK: spacedock binary 0.20.5 and plugin 0.20.0 are compatible.") {
				t.Fatalf("host %s: OK line not preserved alongside the hint: %q", host, res.Message)
			}
			if !strings.Contains(res.Message, "newer plugin") {
				t.Fatalf("host %s: hint missing newer-plugin notice: %q", host, res.Message)
			}
			if !strings.Contains(res.Message, "spacedock install --host "+host) {
				t.Fatalf("host %s: hint missing host install command: %q", host, res.Message)
			}
		}
	})

	// Negative: equal versions carry no hint — there is nothing to upgrade to.
	t.Run("equal-version-no-hint", func(t *testing.T) {
		res := Compare("claude", "0.20.0", "0.20.0")
		if res.Verdict != Compatible {
			t.Fatalf("verdict = %v, want Compatible", res.Verdict)
		}
		if strings.Contains(res.Message, "newer plugin") || strings.Contains(res.Message, "spacedock install") {
			t.Fatalf("equal versions must not emit an upgrade hint: %q", res.Message)
		}
	})

	// Negative: an unstamped `dev` binary version is not valid semver — the hint
	// must not fire (no false "you must upgrade" against a dev build). Note: a raw
	// `dev` binaryVersion is TooOldBinary under Compare (D3), so this exercises
	// semverCompare's own defensive gate directly via a case that still lands
	// Compatible: a `dev`-suffixed but major.minor-parseable binary.
	t.Run("dev-suffixed-binary-no-hint", func(t *testing.T) {
		res := Compare("claude", "0.20.0", "0.20.0-pre1+dev")
		if res.Verdict != Compatible {
			t.Fatalf("verdict = %v, want Compatible", res.Verdict)
		}
		if strings.Contains(res.Message, "newer plugin") || strings.Contains(res.Message, "spacedock install") {
			t.Fatalf("a prerelease/build-suffixed binary version must not emit an upgrade hint: %q", res.Message)
		}
	})

	// Negative: a binary OLDER than the plugin (but still same-minor compatible)
	// carries no hint — the hint is for a behind plugin, not a behind binary.
	t.Run("older-binary-no-hint", func(t *testing.T) {
		res := Compare("claude", "0.20.5", "0.20.0")
		if res.Verdict != Compatible {
			t.Fatalf("verdict = %v, want Compatible", res.Verdict)
		}
		if strings.Contains(res.Message, "newer plugin") {
			t.Fatalf("a newer plugin than the binary must not emit the behind-plugin hint: %q", res.Message)
		}
	})

	// Positive across the single-vs-double-digit boundary: binary 0.20.10 is
	// numerically NEWER than plugin 0.20.9 and MUST hint — but lexically "0.20.10"
	// sorts BEFORE "0.20.9" ("1" < "9"), so a lexical-compare regression of
	// semverCompare would wrongly suppress the hint here. Pins the integer compare.
	t.Run("behind-plugin-double-digit-patch-hints", func(t *testing.T) {
		res := Compare("claude", "0.20.9", "0.20.10")
		if res.Verdict != Compatible {
			t.Fatalf("verdict = %v, want Compatible", res.Verdict)
		}
		if !strings.Contains(res.Message, "newer plugin") {
			t.Fatalf("binary 0.20.10 is numerically newer than plugin 0.20.9 — the hint MUST fire: %q", res.Message)
		}
	})

	// Negative mirror across the boundary: binary 0.20.9 is numerically OLDER than
	// plugin 0.20.10 and MUST NOT hint — but lexically "0.20.9" sorts AFTER
	// "0.20.10", so a lexical-compare regression would wrongly FIRE the hint on
	// this older binary. The two boundary cases together RED any lexical compare.
	t.Run("older-binary-double-digit-patch-no-hint", func(t *testing.T) {
		res := Compare("claude", "0.20.10", "0.20.9")
		if res.Verdict != Compatible {
			t.Fatalf("verdict = %v, want Compatible", res.Verdict)
		}
		if strings.Contains(res.Message, "newer plugin") {
			t.Fatalf("binary 0.20.9 is numerically older than plugin 0.20.10 — the hint MUST NOT fire: %q", res.Message)
		}
	})

	// Negative: a pre-release binary version (e.g. 0.20.5-rc1) is not clean
	// dotted-int semver — the conservative gate emits no hint. Pins that
	// parseDottedInts rejects a `-rc1` suffix rather than stripping it.
	t.Run("prerelease-binary-no-hint", func(t *testing.T) {
		res := Compare("claude", "0.20.0", "0.20.5-rc1")
		if res.Verdict != Compatible {
			t.Fatalf("verdict = %v, want Compatible", res.Verdict)
		}
		if strings.Contains(res.Message, "newer plugin") || strings.Contains(res.Message, "spacedock install") {
			t.Fatalf("a pre-release binary version must not emit an upgrade hint: %q", res.Message)
		}
	})
}

// TestCompareMessageShape locks the shared mismatch-message shape: the leading
// "Spacedock version mismatch" line names both display versions, and the message
// ends with the "Run spacedock doctor" pointer — for every mismatch class except
// malformed-version (which is a packaging bug, not a too-old install) and
// no-plugin-found.
func TestCompareMessageShape(t *testing.T) {
	for _, c := range []struct {
		plugin string
		binary string
	}{
		{"0.20.0", "0.19.4"}, // too-old-binary
		{"0.18.0", "0.19.4"}, // too-old-plugin
	} {
		res := Compare("claude", c.plugin, c.binary)
		header := "Spacedock version mismatch: binary " + c.binary + ", plugin " + c.plugin
		if !strings.Contains(res.Message, header) {
			t.Errorf("Compare(plugin=%q, binary=%q) message missing header %q: %q", c.plugin, c.binary, header, res.Message)
		}
		if !strings.Contains(res.Message, "Run spacedock doctor") {
			t.Errorf("Compare(plugin=%q, binary=%q) message missing doctor pointer: %q", c.plugin, c.binary, res.Message)
		}
	}
}

// TestCompareHostSubstitution verifies the host parameter is woven into the
// too-old-plugin remedy (the only place an install/update host appears). The
// remedy must name the live `spacedock install` command, not the removed `init`
// (which now exits 2) — the remedy a user hits at the gate must run.
func TestCompareHostSubstitution(t *testing.T) {
	for _, host := range []string{"claude", "codex"} {
		res := Compare(host, "0.18.0", "0.19.4")
		want := "spacedock install --host " + host
		if !strings.Contains(res.Message, want) {
			t.Errorf("too-old-plugin remedy for host %q missing %q: %q", host, want, res.Message)
		}
		if strings.Contains(res.Message, "spacedock init") {
			t.Errorf("too-old-plugin remedy for host %q names the removed init command: %q", host, res.Message)
		}
	}
}
