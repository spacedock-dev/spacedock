// ABOUTME: The contract-version axis — CONTRACT_VERSION, the half-open
// ABOUTME: requires-contract range grammar, and the five-verdict compatibility compare.
package contract

import (
	"fmt"
	"strconv"
	"strings"
)

// CONTRACT_VERSION is the monotonic integer naming the binary's observable
// contract the vendored skill surface depends on (the set of `spacedock`
// subcommands, flags, and output sections). It is distinct from the plugin's
// display semver and from the build version. Bump it only when a change to the
// binary alters the observable surface the FO/ensign contracts call — never as a
// side effect of a routine release bump (see the entity's OPEN-2 bump discipline).
const CONTRACT_VERSION = 1

// Verdict is the compatibility class produced by comparing a binary's contract
// version against a plugin's declared requires-contract range.
type Verdict int

const (
	// Compatible means lo <= C < hi: the binary's contract sits inside the
	// plugin's declared half-open range.
	Compatible Verdict = iota
	// TooOldBinary means C < lo: the installed binary predates the contract the
	// plugin needs. Remedy: rebuild/upgrade the binary.
	TooOldBinary
	// TooOldPlugin means C >= hi: the installed plugin predates the binary's
	// contract. Remedy: update/reinstall the plugin.
	TooOldPlugin
	// MalformedRange means the manifest's requires-contract does not parse as
	// ">=N,<M". A packaging bug, not a too-old install; no upgrade remedy.
	MalformedRange
	// NoPluginFound means no installed plugin manifest could be resolved for the
	// host. Distinct, non-fatal-by-default; reported rather than asserting
	// compatibility.
	NoPluginFound
	// PluginPredatesContract means the manifest has no requires-contract field at
	// all: the installed plugin predates the contract mechanism. Kin to
	// too-old-plugin, but with no range to name; remedy reinstalls via
	// `spacedock install` and omits the `plugin update` fallback (which no-ops on a
	// stale install).
	PluginPredatesContract
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
	case MalformedRange:
		return "malformed-range"
	case NoPluginFound:
		return "no-plugin-found"
	case PluginPredatesContract:
		return "plugin-predates-contract"
	default:
		return "unknown"
	}
}

// Result carries a comparison's verdict and the operator-facing message. For
// Compatible the message is a one-line "OK" report; for every mismatch it is the
// shared-shape actionable message with the per-class remedy.
type Result struct {
	Verdict Verdict
	Message string
}

// ParseRange parses a requires-contract value of the form ">=N,<M" into its
// half-open integer bounds [lo, hi). Surrounding whitespace is tolerated. Any
// other shape — missing a bound, the wrong operator, non-integer bounds, an
// empty or inverted interval, or extra clauses — is a parse error.
func ParseRange(raw string) (lo int, hi int, err error) {
	parts := strings.Split(raw, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected exactly two clauses %q: %q", ">=N,<M", raw)
	}
	loStr := strings.TrimSpace(parts[0])
	hiStr := strings.TrimSpace(parts[1])
	if !strings.HasPrefix(loStr, ">=") {
		return 0, 0, fmt.Errorf("lower bound must use >=: %q", loStr)
	}
	if !strings.HasPrefix(hiStr, "<") || strings.HasPrefix(hiStr, "<=") {
		return 0, 0, fmt.Errorf("upper bound must use <: %q", hiStr)
	}
	lo, err = strconv.Atoi(strings.TrimSpace(loStr[2:]))
	if err != nil {
		return 0, 0, fmt.Errorf("lower bound is not an integer: %q", loStr)
	}
	hi, err = strconv.Atoi(strings.TrimSpace(hiStr[1:]))
	if err != nil {
		return 0, 0, fmt.Errorf("upper bound is not an integer: %q", hiStr)
	}
	if lo >= hi {
		return 0, 0, fmt.Errorf("empty interval: lower bound %d not below upper bound %d", lo, hi)
	}
	return lo, hi, nil
}

// Compare classifies a binary at contract version c against a plugin's raw
// requires-contract range, for the named host and (pre-release) dev branch. It
// returns the verdict and the operator-facing message. pluginVersion and
// binaryVersion are the display semvers woven into the user-facing mismatch/OK
// lines. NoPluginFound is produced by the caller (when the manifest is absent),
// not here — Compare always has a raw range string to evaluate.
func Compare(c int, raw, host, branch, pluginVersion, binaryVersion string) Result {
	manifestNote := ""
	return compareWithManifest(c, raw, host, branch, manifestNote, pluginVersion, binaryVersion)
}

// compareWithManifest is Compare with an optional manifest path woven into the
// malformed-range message so a packaging bug names the offending file.
func compareWithManifest(c int, raw, host, branch, manifestPath, pluginVersion, binaryVersion string) Result {
	if strings.TrimSpace(raw) == "" {
		return Result{
			Verdict: PluginPredatesContract,
			Message: pluginPredatesContractRemedy(host, branch),
		}
	}
	lo, hi, err := ParseRange(raw)
	if err != nil {
		loc := manifestPath
		if loc == "" {
			loc = "the plugin manifest"
		}
		return Result{
			Verdict: MalformedRange,
			Message: fmt.Sprintf(
				"malformed contract range %q in %s: expected \">=N,<M\". "+
					"This is a packaging bug — the plugin manifest is wrong, not your install.",
				strings.TrimSpace(raw), loc),
		}
	}
	switch {
	case c < lo:
		return Result{Verdict: TooOldBinary, Message: mismatchMessage(binaryVersion, pluginVersion, "Upgrade the binary to continue.", tooOldBinaryRemedy())}
	case c >= hi:
		return Result{Verdict: TooOldPlugin, Message: mismatchMessage(binaryVersion, pluginVersion, "Update the plugin to continue.", tooOldPluginRemedy(host))}
	default:
		msg := fmt.Sprintf("OK: spacedock binary %s and plugin %s are compatible.", binaryVersion, pluginVersion)
		if hint := upgradeHint(host, pluginVersion, binaryVersion); hint != "" {
			msg += "\n" + hint
		}
		return Result{Verdict: Compatible, Message: msg}
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
// or empty value from triggering the upgrade hint. Unlike the cli resolver's
// lexical fallback, a non-integer component here is a parse failure, not a
// lexical tiebreak: the hint must not fire on anything but a clean semver skew.
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

// tooOldBinaryRemedy is the pinned too-old-binary remedy block: it leads with the
// Homebrew upgrade, keeps the source-build fallback, and names the binary-vs-plugin
// distinction (refreshing the plugin instead is a different command).
func tooOldBinaryRemedy() string {
	return "  Upgrade via Homebrew: brew upgrade spacedock\n" +
		"  Or build from source: go build -o spacedock ./cmd/spacedock\n" +
		"  Or refresh the plugin instead: spacedock install"
}

// tooOldPluginRemedy is the pinned too-old-plugin remedy line, parameterized by
// the detected host for the install hint.
func tooOldPluginRemedy(host string) string {
	return fmt.Sprintf("  Update the plugin: spacedock install --host %s", host)
}

// pluginPredatesContractRemedy is the pinned remedy for an installed plugin that
// predates the contract mechanism (no requires-contract field). It names the
// `spacedock install` one-liner — never raw `<host> plugin` commands — and omits the
// `plugin update` fallback, which no-ops on a stale already-installed plugin. The
// host is parameterized; the optional pre-release branch suffixes the reinstall
// source so a dev install reflects the branch (the default release path omits it).
func pluginPredatesContractRemedy(host, branch string) string {
	source := "spacedock-dev/spacedock"
	if branch != "" {
		source += "@" + branch
	}
	return fmt.Sprintf(
		"plugin-predates-contract: your installed Spacedock plugin is out of date "+
			"(predates this binary's contract). Upgrade it: spacedock install --host %s "+
			"(reinstalls from %s).",
		host, source)
}

// noPluginMessage is the pinned no-plugin-found report for a host. Not a
// mismatch — there is no range to compare — so it stands alone, exit non-fatal
// by the caller's policy.
func noPluginMessage(host string) string {
	return fmt.Sprintf(
		"no installed Spacedock plugin found for host %s. Install it: spacedock install --host %s.",
		host, host)
}
