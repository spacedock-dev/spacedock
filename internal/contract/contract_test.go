// ABOUTME: Unit tests for the contract-version axis: CONTRACT_VERSION sanity,
// ABOUTME: half-open range parsing, and the five-verdict compatibility compare.
package contract

import (
	"strings"
	"testing"
)

// TestContractVersionIsPositiveInteger guards the axis invariant: CONTRACT_VERSION
// is a positive integer (the monotonic compatibility axis, not semver).
func TestContractVersionIsPositiveInteger(t *testing.T) {
	if CONTRACT_VERSION < 1 {
		t.Fatalf("CONTRACT_VERSION = %d, want >= 1", CONTRACT_VERSION)
	}
}

// TestParseRange covers the half-open range grammar ">=N,<M": accepted forms,
// the bracket bounds, and every malformed shape that must be rejected.
func TestParseRange(t *testing.T) {
	cases := []struct {
		raw     string
		wantLo  int
		wantHi  int
		wantErr bool
	}{
		{">=1,<2", 1, 2, false},
		{">=1,<3", 1, 3, false},
		{">=10,<11", 10, 11, false},
		{">= 1 , < 2", 1, 2, false}, // surrounding whitespace tolerated
		{"", 0, 0, true},
		{">=1", 0, 0, true},        // missing upper bound
		{"<2", 0, 0, true},         // missing lower bound
		{">1,<2", 0, 0, true},      // wrong lower operator
		{">=1,<=2", 0, 0, true},    // wrong upper operator
		{">=a,<2", 0, 0, true},     // non-integer lower
		{">=1,<b", 0, 0, true},     // non-integer upper
		{">=2,<2", 0, 0, true},     // empty interval (lo >= hi)
		{">=3,<2", 0, 0, true},     // inverted interval
		{">=1,<2,<3", 0, 0, true},  // extra clause
		{"1,2", 0, 0, true},        // no operators
		{">=1.0,<2.0", 0, 0, true}, // float bounds
	}
	for _, c := range cases {
		lo, hi, err := ParseRange(c.raw)
		if c.wantErr {
			if err == nil {
				t.Errorf("ParseRange(%q) = (%d,%d,nil), want error", c.raw, lo, hi)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseRange(%q) unexpected error: %v", c.raw, err)
			continue
		}
		if lo != c.wantLo || hi != c.wantHi {
			t.Errorf("ParseRange(%q) = (%d,%d), want (%d,%d)", c.raw, lo, hi, c.wantLo, c.wantHi)
		}
	}
}

// TestCompare drives the five-verdict comparison over a contract C and a raw
// range, asserting the verdict class and that each non-compatible verdict's
// message carries its pinned remedy substring.
func TestCompare(t *testing.T) {
	cases := []struct {
		name        string
		contract    int
		raw         string
		wantVerdict Verdict
		wantPinned  string // substring the message must contain (empty = no check)
	}{
		{"compatible-exact", 1, ">=1,<2", Compatible, ""},
		{"compatible-forward-tolerant", 2, ">=1,<3", Compatible, ""},
		{"compatible-lower-edge", 1, ">=1,<3", Compatible, ""},
		{"too-old-binary", 1, ">=2,<3", TooOldBinary, "Upgrade the binary to continue."},
		{"too-old-plugin-at-hi", 2, ">=1,<2", TooOldPlugin, "Update the plugin to continue."},
		{"too-old-plugin-above-hi", 5, ">=1,<2", TooOldPlugin, "Update the plugin to continue."},
		{"malformed", 1, ">=1", MalformedRange, "malformed contract range"},
		{"predates-contract-empty", 1, "", PluginPredatesContract, "spacedock install --host claude"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := Compare(c.contract, c.raw, "claude", "0.12.1", "0.19.4")
			if res.Verdict != c.wantVerdict {
				t.Fatalf("Compare(%d,%q) verdict = %v, want %v", c.contract, c.raw, res.Verdict, c.wantVerdict)
			}
			if c.wantPinned != "" && !strings.Contains(res.Message, c.wantPinned) {
				t.Fatalf("Compare(%d,%q) message = %q, want substring %q", c.contract, c.raw, res.Message, c.wantPinned)
			}
			// The compatible verdict carries no remedy; non-compatible ones do.
			if c.wantVerdict == Compatible && res.Message == "" {
				t.Fatalf("compatible verdict should still carry an OK message")
			}
		})
	}
}

// TestCompatibleUpgradeHint locks AC-3: when the binary and plugin display
// versions are both valid semver and the binary is strictly newer while the
// contract is compatible, the Compatible message carries an opt-in upgrade hint
// naming a newer plugin AND the host-specific `spacedock install --host {host}`
// command — but the verdict stays Compatible and the OK line is preserved. The
// hint NEVER fires on equal versions or a non-semver (`dev`) binary version. The
// trigger (the semver skew) comes from the version inputs, not the message.
func TestCompatibleUpgradeHint(t *testing.T) {
	t.Run("behind-plugin-hints", func(t *testing.T) {
		for _, host := range []string{"claude", "codex"} {
			res := Compare(CONTRACT_VERSION, ">=2,<3", host, "0.19.8", "0.20.0")
			if res.Verdict != Compatible {
				t.Fatalf("host %s: verdict = %v, want Compatible (the hint must not change the verdict)", host, res.Verdict)
			}
			if !strings.Contains(res.Message, "OK: spacedock binary 0.20.0 and plugin 0.19.8 are compatible.") {
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
		res := Compare(CONTRACT_VERSION, ">=2,<3", "claude", "0.20.0", "0.20.0")
		if res.Verdict != Compatible {
			t.Fatalf("verdict = %v, want Compatible", res.Verdict)
		}
		if strings.Contains(res.Message, "newer plugin") || strings.Contains(res.Message, "spacedock install") {
			t.Fatalf("equal versions must not emit an upgrade hint: %q", res.Message)
		}
	})

	// Negative: an unstamped `dev` binary version is not valid semver — the hint
	// must not fire (no false "you must upgrade" against a dev build).
	t.Run("dev-binary-no-hint", func(t *testing.T) {
		res := Compare(CONTRACT_VERSION, ">=2,<3", "claude", "0.19.8", "dev")
		if res.Verdict != Compatible {
			t.Fatalf("verdict = %v, want Compatible", res.Verdict)
		}
		if strings.Contains(res.Message, "newer plugin") || strings.Contains(res.Message, "spacedock install") {
			t.Fatalf("dev binary version must not emit an upgrade hint: %q", res.Message)
		}
	})

	// Negative: a binary OLDER than the plugin (but still contract-compatible)
	// carries no hint — the hint is for a behind plugin, not a behind binary.
	t.Run("older-binary-no-hint", func(t *testing.T) {
		res := Compare(CONTRACT_VERSION, ">=2,<3", "claude", "0.21.0", "0.20.0")
		if res.Verdict != Compatible {
			t.Fatalf("verdict = %v, want Compatible", res.Verdict)
		}
		if strings.Contains(res.Message, "newer plugin") {
			t.Fatalf("a newer plugin than the binary must not emit the behind-plugin hint: %q", res.Message)
		}
	})

	// Positive across the single-vs-double-digit boundary: binary 0.10.0 is
	// numerically NEWER than plugin 0.9.0 and MUST hint — but lexically "0.10.0"
	// sorts BEFORE "0.9.0" ("1" < "9"), so a lexical-compare regression of
	// semverCompare would wrongly suppress the hint here. Pins the integer compare.
	t.Run("behind-plugin-double-digit-minor-hints", func(t *testing.T) {
		res := Compare(CONTRACT_VERSION, ">=2,<3", "claude", "0.9.0", "0.10.0")
		if res.Verdict != Compatible {
			t.Fatalf("verdict = %v, want Compatible", res.Verdict)
		}
		if !strings.Contains(res.Message, "newer plugin") {
			t.Fatalf("binary 0.10.0 is numerically newer than plugin 0.9.0 — the hint MUST fire: %q", res.Message)
		}
	})

	// Negative mirror across the boundary: binary 0.9.0 is numerically OLDER than
	// plugin 0.10.0 and MUST NOT hint — but lexically "0.9.0" sorts AFTER "0.10.0",
	// so a lexical-compare regression would wrongly FIRE the hint on this older
	// binary. The two boundary cases together RED any lexical compare.
	t.Run("older-binary-double-digit-minor-no-hint", func(t *testing.T) {
		res := Compare(CONTRACT_VERSION, ">=2,<3", "claude", "0.10.0", "0.9.0")
		if res.Verdict != Compatible {
			t.Fatalf("verdict = %v, want Compatible", res.Verdict)
		}
		if strings.Contains(res.Message, "newer plugin") {
			t.Fatalf("binary 0.9.0 is numerically older than plugin 0.10.0 — the hint MUST NOT fire: %q", res.Message)
		}
	})

	// Negative: a pre-release binary version (e.g. 0.20.0-rc1) is not clean
	// dotted-int semver — the conservative gate emits no hint. Pins that
	// parseDottedInts rejects a `-rc1` suffix rather than stripping it.
	t.Run("prerelease-binary-no-hint", func(t *testing.T) {
		res := Compare(CONTRACT_VERSION, ">=2,<3", "claude", "0.19.8", "0.20.0-rc1")
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
// malformed-range (which is a packaging bug, not a too-old install) and
// no-plugin-found.
func TestCompareMessageShape(t *testing.T) {
	const pluginVersion, binaryVersion = "0.18.0", "0.19.4"
	for _, c := range []struct {
		contract int
		raw      string
	}{
		{1, ">=2,<3"}, // too-old-binary
		{2, ">=1,<2"}, // too-old-plugin
	} {
		res := Compare(c.contract, c.raw, "claude", pluginVersion, binaryVersion)
		header := "Spacedock version mismatch: binary " + binaryVersion + ", plugin " + pluginVersion
		if !strings.Contains(res.Message, header) {
			t.Errorf("Compare(%d,%q) message missing header %q: %q", c.contract, c.raw, header, res.Message)
		}
		if !strings.Contains(res.Message, "Run spacedock doctor") {
			t.Errorf("Compare(%d,%q) message missing doctor pointer: %q", c.contract, c.raw, res.Message)
		}
	}
}

// TestPluginPredatesContractRemedy locks the verdict for an absent/empty
// requires-contract: it names the `spacedock install --host <host>` one-liner and
// OMITS the `plugin update` fallback that reusing too-old-plugin would drag in
// (that fallback no-ops on a stale install). Post-decouple the remedy carries NO
// reinstall-source parenthetical: `spacedock install` auto-selects the channel from
// the binary's own devBranch stamp (entry name), so there is no plugin-repo name
// and no `@branch` shorthand to name. A whitespace-only value routes here too; a
// non-empty unparseable value still reads as a packaging bug.
func TestPluginPredatesContractRemedy(t *testing.T) {
	for _, raw := range []string{"", "   "} {
		res := Compare(1, raw, "claude", "0.12.1", "0.19.4")
		if res.Verdict != PluginPredatesContract {
			t.Fatalf("Compare(1,%q) verdict = %v, want plugin-predates-contract", raw, res.Verdict)
		}
		if !strings.Contains(res.Message, "spacedock install --host claude") {
			t.Errorf("predates-contract remedy missing install one-liner: %q", res.Message)
		}
		if strings.Contains(res.Message, "spacedock-dev/spacedock") {
			t.Errorf("predates-contract remedy names the plugin repo; post-decouple the reinstall source is not named in the remedy: %q", res.Message)
		}
		if strings.Contains(res.Message, "@next") {
			t.Errorf("predates-contract remedy carries the removed @branch shorthand: %q", res.Message)
		}
		if strings.Contains(res.Message, "plugin update") {
			t.Errorf("predates-contract remedy must omit the no-op `plugin update` fallback: %q", res.Message)
		}
	}

	// A non-empty unparseable value is still a packaging bug, not predates-contract.
	bug := Compare(1, ">=1", "claude", "0.12.1", "0.19.4")
	if bug.Verdict != MalformedRange {
		t.Fatalf("Compare(1,%q) verdict = %v, want malformed-range", ">=1", bug.Verdict)
	}
	if !strings.Contains(bug.Message, "This is a packaging bug") {
		t.Errorf("non-empty malformed should keep the packaging-bug message: %q", bug.Message)
	}
}

// TestCompareHostSubstitution verifies the host parameter is woven into the
// too-old-plugin remedy (the only place an install/update host appears). The
// remedy must name the live `spacedock install` command, not the removed `init`
// (which now exits 2) — the remedy a user hits at the gate must run.
func TestCompareHostSubstitution(t *testing.T) {
	for _, host := range []string{"claude", "codex"} {
		res := Compare(2, ">=1,<2", host, "0.18.0", "0.19.4")
		want := "spacedock install --host " + host
		if !strings.Contains(res.Message, want) {
			t.Errorf("too-old-plugin remedy for host %q missing %q: %q", host, want, res.Message)
		}
		if strings.Contains(res.Message, "spacedock init") {
			t.Errorf("too-old-plugin remedy for host %q names the removed init command: %q", host, res.Message)
		}
	}
}
