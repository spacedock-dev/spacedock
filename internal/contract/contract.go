// ABOUTME: The version-compatibility axis — major.minor extraction and the
// ABOUTME: minor-exact binary<->plugin compatibility compare.
package contract

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// Verdict is the compatibility class produced by comparing a binary's display
// version against a plugin's declared display version.
type Verdict int

const (
	// Compatible means the binary and plugin share the same (major, minor).
	// Patch and prerelease skew within that minor are interchangeable.
	Compatible Verdict = iota
	// TooOldBinary means the binary's (major, minor) is behind the plugin's, or
	// the binary's version carries no parseable major.minor at all (the shape of
	// an integer-era `dev` source build). Remedy: rebuild/upgrade the binary.
	TooOldBinary
	// TooOldPlugin means the binary's (major, minor) is ahead of the plugin's.
	// Remedy: update/reinstall the plugin.
	TooOldPlugin
	// MalformedVersion means the manifest's version field is missing or does not
	// parse as a dotted major.minor(.patch) semver. A packaging bug, not a
	// too-old install; no upgrade remedy.
	MalformedVersion
	// NoPluginFound means no installed plugin manifest could be resolved for the
	// host. Distinct, non-fatal-by-default; reported rather than asserting
	// compatibility.
	NoPluginFound
)

// String renders the verdict's stable kebab-case token (the oracle string AC-1
// asserts on for compatible/no-plugin-found and embeds in the remedy lines).
func (v Verdict) String() string {
	switch v {
	case Compatible:
		return "compatible"
	case TooOldBinary:
		return "too-old-binary"
	case TooOldPlugin:
		return "too-old-plugin"
	case MalformedVersion:
		return "malformed-version"
	case NoPluginFound:
		return "no-plugin-found"
	default:
		return "unknown"
	}
}

// Result carries a comparison's verdict and the operator-facing message. For
// Compatible the message is a one-line "OK" report; for every mismatch it is the
// shared-shape actionable message with the per-class remedy. Hint is the opt-in
// upgrade hint for a compatible-but-behind plugin — empty unless the binary is a
// strictly-newer semver than the plugin. Doctor folds it into Message; the front
// door surfaces Hint alone (it stays silent on the bare OK line).
type Result struct {
	Verdict Verdict
	Message string
	Hint    string
}

// ParseMajorMinor extracts the (major, minor) pair from a dotted version string,
// cutting any prerelease/build suffix at the first "-" or "+" before parsing — so
// both published suffix styles (`0.24.0-pre1`, `0.23.0-pre.4`) and a build tag
// (`0.24.0+dev`) parse identically. It reports ok=false when fewer than two
// dotted integer components remain after the cut: the shape of the integer-era
// `dev` sentinel and any other non-semver token.
func ParseMajorMinor(v string) (major, minor int, ok bool) {
	v = strings.TrimSpace(v)
	if i := strings.IndexAny(v, "-+"); i >= 0 {
		v = v[:i]
	}
	parts := strings.SplitN(v, ".", 3)
	if len(parts) < 2 {
		return 0, 0, false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil || major < 0 {
		return 0, 0, false
	}
	minor, err = strconv.Atoi(parts[1])
	if err != nil || minor < 0 {
		return 0, 0, false
	}
	return major, minor, true
}

// Compare classifies a binary at binaryVersion against a plugin at
// pluginVersion, for the named host, by comparing their (major, minor) pairs
// (D1: minor-exact, both directions). Patch and prerelease skew within the same
// minor are interchangeable. NoPluginFound is produced by the caller (when the
// manifest is absent), not here — Compare always has a raw plugin version to
// evaluate.
func Compare(host, pluginVersion, binaryVersion string) Result {
	return compareNamed(host, "", pluginVersion, binaryVersion, false)
}

// compareNamed is Compare with an optional manifest path woven into the
// malformed-version message so a packaging bug names the offending file, and an
// edgeCask flag that names the `spacedock@next` formula in the too-old-binary
// remedy when the running binary is the edge cask install.
func compareNamed(host, manifestPath, pluginVersion, binaryVersion string, edgeCask bool) Result {
	bMajor, bMinor, bOk := ParseMajorMinor(binaryVersion)
	if !bOk {
		// A binary version with no parseable major.minor can only be an
		// integer-era source build (`dev`, pre-D3 embed) — treated as too-old,
		// with the existing remedy's "build from source" arm doubling as the
		// rebuild hint.
		return Result{Verdict: TooOldBinary, Message: mismatchMessage(binaryVersion, pluginVersion, "Upgrade the binary to continue.", tooOldBinaryRemedy(edgeCask, runtime.GOOS))}
	}
	pMajor, pMinor, pOk := ParseMajorMinor(pluginVersion)
	if !pOk {
		loc := manifestPath
		if loc == "" {
			loc = "the plugin manifest"
		}
		return Result{
			Verdict: MalformedVersion,
			Message: fmt.Sprintf(
				"malformed plugin version %q in %s: expected a dotted major.minor(.patch) semver. "+
					"This is a packaging bug — the plugin manifest is wrong, not your install.",
				strings.TrimSpace(pluginVersion), loc),
		}
	}
	switch {
	case bMajor < pMajor || (bMajor == pMajor && bMinor < pMinor):
		return Result{Verdict: TooOldBinary, Message: mismatchMessage(binaryVersion, pluginVersion, "Upgrade the binary to continue.", tooOldBinaryRemedy(edgeCask, runtime.GOOS))}
	case bMajor > pMajor || (bMajor == pMajor && bMinor > pMinor):
		return Result{Verdict: TooOldPlugin, Message: mismatchMessage(binaryVersion, pluginVersion, "Update the plugin to continue.", tooOldPluginRemedy(host))}
	default:
		msg := fmt.Sprintf("OK: spacedock binary %s and plugin %s are compatible.", binaryVersion, pluginVersion)
		hint := upgradeHint(host, pluginVersion, binaryVersion)
		if hint != "" {
			msg += "\n" + hint
		}
		return Result{Verdict: Compatible, Message: msg, Hint: hint}
	}
}

// upgradeHint returns the opt-in upgrade hint appended to a Compatible message
// when the binary's display semver is strictly newer than the plugin's — a
// plugin that still works (the contract is compatible) but is behind. It returns
// "" (no hint) unless BOTH versions are valid dotted-int semver and the binary
// is strictly greater: an unstamped `dev` binary, a non-semver, or an
// equal/older binary emits nothing, so the gate never fires a false "you must
// upgrade". The hint names the host-specific refresh command.
func upgradeHint(host, pluginVersion, binaryVersion string) string {
	if semverCompare(binaryVersion, pluginVersion) <= 0 {
		return ""
	}
	return fmt.Sprintf(
		"A newer plugin is available — run spacedock install --host %s to refresh.",
		host)
}

// semverCompare orders two dotted-int versions (e.g. `0.20.0`), returning -1, 0,
// or 1. It returns 0 (treat as not-greater, so no hint) when EITHER side is not
// a valid dotted-int version — the defensive gate that keeps a non-semver (`dev`)
// or empty value from triggering the upgrade hint. Unlike Compare's major.minor
// gate, a non-integer component here is a parse failure, not a truncation point:
// the hint must not fire on anything but a clean full-version semver skew.
func semverCompare(a, b string) int {
	an, aok := parseDottedInts(a)
	bn, bok := parseDottedInts(b)
	if !aok || !bok {
		return 0
	}
	for i := 0; i < len(an) || i < len(bn); i++ {
		var av, bv int
		if i < len(an) {
			av = an[i]
		}
		if i < len(bn) {
			bv = bn[i]
		}
		if av != bv {
			if av < bv {
				return -1
			}
			return 1
		}
	}
	return 0
}

// parseDottedInts splits a dotted version into its integer components, reporting
// ok=false when the string is empty or any component is not a non-negative
// integer. A pre-release/build suffix (e.g. `-rc1`) makes its component
// non-integer and so fails the parse — the conservative read for the upgrade
// hint, which only fires on a clean numeric skew.
func parseDottedInts(v string) ([]int, bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil, false
	}
	parts := strings.Split(v, ".")
	out := make([]int, len(parts))
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return nil, false
		}
		out[i] = n
	}
	return out, true
}

// mismatchMessage assembles the shared-shape mismatch message: a header naming
// the binary and plugin display versions, a one-line direction sentence, the
// per-class remedy block, and the doctor pointer.
func mismatchMessage(binaryVersion, pluginVersion, direction, remedy string) string {
	return fmt.Sprintf(
		"Spacedock version mismatch: binary %s, plugin %s. %s\n"+
			"%s\n"+
			"Run spacedock doctor for details.",
		binaryVersion, pluginVersion, direction, remedy)
}

// tooOldBinaryRemedy is the pinned too-old-binary remedy block. goos is
// threaded explicitly (the caller supplies runtime.GOOS) so tests pin the OS
// deterministically. On Linux the remedy leads with the documented curl|sh
// install path — a Linux host has no brew — and the edgeCask flag is ignored
// there: `spacedock@next` is a brew-only token, meaningless without Homebrew,
// so the curl|sh lead supersedes the edge-cask branch. On macOS it leads with
// the Homebrew upgrade, keeps the source-build fallback, and names the
// binary-vs-plugin distinction (refreshing the plugin instead is a different
// command); edgeCask names the edge `spacedock@next` cask as the upgrade target
// when the running binary was installed from it: that cask is a separate
// Homebrew token, so a plain `brew upgrade spacedock` is a no-op for an @next
// install. edgeCask=false leaves the block unchanged.
func tooOldBinaryRemedy(edgeCask bool, goos string) string {
	if goos == "linux" {
		return "  Install the latest binary: curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh\n" +
			"  Or build from source: go build -o spacedock ./cmd/spacedock\n" +
			"  Or refresh the plugin instead: spacedock install"
	}
	formula := "spacedock"
	if edgeCask {
		formula = "spacedock@next"
	}
	return "  Upgrade via Homebrew: brew upgrade " + formula + "\n" +
		"  Or build from source: go build -o spacedock ./cmd/spacedock\n" +
		"  Or refresh the plugin instead: spacedock install"
}

// tooOldPluginRemedy is the pinned too-old-plugin remedy line, parameterized by
// the detected host for the install hint.
func tooOldPluginRemedy(host string) string {
	return fmt.Sprintf("  Update the plugin: spacedock install --host %s", host)
}

// noPluginMessage is the pinned no-plugin-found report for a host. Not a
// mismatch — there is no plugin to compare against — so it stands alone, exit
// non-fatal by the caller's policy.
func noPluginMessage(host string) string {
	return fmt.Sprintf(
		"no installed Spacedock plugin found for host %s. Install it: spacedock install --host %s.",
		host, host)
}
