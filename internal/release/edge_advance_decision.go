// ABOUTME: Edge-advance line-ordering guard — orders X.Y.Z / X.Y.Z-preN versions
// ABOUTME: and decides whether a tag advances `next` or skips the whole job.
package release

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// preVersionRe matches X.Y.Z with an optional semver prerelease suffix
// (`-<identifiers>`). The prerelease is captured whole and ordered per semver
// §11 by comparePrerelease; build metadata is not part of the release scheme
// this tool cuts, so `+meta` is not accepted.
var preVersionRe = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z.-]+))?$`)

// parsePreVersion splits a version into its X.Y.Z core and its prerelease string
// (empty for a plain stable). It errors on anything that is not X.Y.Z or
// X.Y.Z-<prerelease>, so the decision guard fails loudly on a malformed version
// rather than treating it as equal (which would silently skip).
func parsePreVersion(v string) (core [3]int, pre string, err error) {
	m := preVersionRe.FindStringSubmatch(v)
	if m == nil {
		return core, "", fmt.Errorf("version %q is not X.Y.Z or X.Y.Z-preN", v)
	}
	for i := 0; i < 3; i++ {
		core[i], _ = strconv.Atoi(m[i+1])
	}
	return core, m[4], nil
}

// ComparePreVersion orders two release versions of the shape X.Y.Z or
// X.Y.Z-<prerelease>, returning -1, 0, or 1 (a<b, a==b, a>b). The X.Y.Z core
// dominates; when cores are equal a version WITHOUT a prerelease outranks one
// WITH it (0.26.0 > 0.26.0-pre1, semver §11), and two prereleases compare by
// their dot-separated identifiers with numeric identifiers ordered numerically
// (so pre0 < pre1 and pre2 < pre10 — the exact distinction the recursion-skip
// needs and contract.semverCompare, dotted-int only, cannot make). Errors on a
// version neither side can parse.
func ComparePreVersion(a, b string) (int, error) {
	ca, pa, err := parsePreVersion(a)
	if err != nil {
		return 0, err
	}
	cb, pb, err := parsePreVersion(b)
	if err != nil {
		return 0, err
	}
	for i := 0; i < 3; i++ {
		if ca[i] != cb[i] {
			return sign(ca[i] - cb[i]), nil
		}
	}
	return comparePrerelease(pa, pb), nil
}

// comparePrerelease orders two prerelease strings. An absent prerelease outranks
// a present one (a stable release is greater than its own prereleases, semver
// §11), so the empty string ranks highest. Two present prereleases compare
// naturally: embedded digit runs order NUMERICALLY, so this scheme's `preN`
// labels sort pre0 < pre1 < pre2 < pre10 rather than lexically (`pre10 < pre2`),
// which is the exact distinction the recursion-skip needs and would silently
// break past pre9 under a plain ASCII compare.
func comparePrerelease(a, b string) int {
	if a == b {
		return 0
	}
	if a == "" {
		return 1 // stable release > any prerelease
	}
	if b == "" {
		return -1
	}
	return naturalCompare(a, b)
}

// naturalCompare orders two strings chunk by chunk, where a chunk is a maximal
// run of digits or of non-digits: two digit runs compare numerically, two
// non-digit runs by ASCII, and a shorter string ranks below a longer one when
// all shared chunks are equal. A digit run and a non-digit run at the same
// position fall back to a byte compare of the remaining strings (a shape this
// scheme's `preN` prereleases never produce).
func naturalCompare(a, b string) int {
	for a != "" && b != "" {
		aDigit, bDigit := isDigit(a[0]), isDigit(b[0])
		if aDigit != bDigit {
			return strings.Compare(a, b)
		}
		aChunk, aRest := splitChunk(a, aDigit)
		bChunk, bRest := splitChunk(b, bDigit)
		if aDigit {
			an, _ := strconv.Atoi(aChunk)
			bn, _ := strconv.Atoi(bChunk)
			if an != bn {
				return sign(an - bn)
			}
		} else if c := strings.Compare(aChunk, bChunk); c != 0 {
			return c
		}
		a, b = aRest, bRest
	}
	return sign(len(a) - len(b))
}

// splitChunk peels the leading maximal run of digits (or of non-digits, per
// digit) off s, returning the chunk and the remainder.
func splitChunk(s string, digit bool) (chunk, rest string) {
	i := 0
	for i < len(s) && isDigit(s[i]) == digit {
		i++
	}
	return s[:i], s[i:]
}

func isDigit(b byte) bool { return b >= '0' && b <= '9' }

func sign(n int) int {
	switch {
	case n < 0:
		return -1
	case n > 0:
		return 1
	default:
		return 0
	}
}

// EdgeAdvanceDecision decides whether a tag ADVANCES the edge line (currently:
// the auto-pre0 cut) or SKIPS, given knownVersion — the highest version already
// known to have been released (originally `next`'s current manifest version,
// pre-round-3; now the caller feeds HighestKnownEdgeVersion's result over git
// tag history — see docs/releasing.md "Advancing the Edge Line"). The function
// itself is agnostic to where knownVersion comes from; only the caller changed.
// It computes the tag's TARGET EDGE VERSION: dev-preversion(tag) for a bare
// stable vX.Y.Z (the edge line would be stamped PAST the release), or the tag's
// own version for a `-pre` tag — and returns advance iff that target is
// STRICTLY GREATER than knownVersion. A patch cut from an older line yields a
// target at or below the known line, so it skips. The boundary is strict `>`,
// not `>=`: a patch whose target exactly equals knownVersion must skip.
func EdgeAdvanceDecision(tag, knownVersion string) (advance bool, targetEdge string, err error) {
	version := strings.TrimPrefix(tag, "v")
	if strings.Contains(version, "-") {
		// A prerelease tag: the edge line would inherit the tag's own version.
		targetEdge = version
	} else {
		// A bare stable tag: the edge line is stamped PAST it to the dev pre-version.
		targetEdge, err = DevPreVersion(version)
		if err != nil {
			return false, "", err
		}
	}
	cmp, err := ComparePreVersion(targetEdge, knownVersion)
	if err != nil {
		return false, "", err
	}
	return cmp > 0, targetEdge, nil
}

// HighestKnownEdgeVersion returns the greatest version among candidates — the
// git-tag-sourced replacement for the retired `next`-manifest read
// EdgeAdvanceDecision's second argument used to come from (see
// docs/releasing.md "Advancing the Edge Line"). Each candidate is a release tag
// name (`v` prefix optional): "X.Y.Z" or "X.Y.Z-preN". PRERELEASE tags are
// INCLUDED in the comparison, not filtered to bare stable ones — the retired
// `next`-manifest read carried prerelease versions too (`next` tracked whatever
// `-preN` last reconciled), so including them here preserves the same
// protection: an old-line stable patch must still lose to an already-cut
// newer-line prerelease, not just to a newer bare stable. Worked example:
// cutting v0.25.3 while v0.27.0-pre7 already exists must find 0.27.0-pre7 as
// the known version (target 0.26.0-pre0 loses to it and skips) — filtering to
// stable-only would find only 0.26.0, a much closer and wrong call. A candidate
// that fails to parse (a malformed tag name, or any non-release ref shape the
// caller's own git-tag filter let through) is skipped, not an error — the
// caller controls the candidate list's shape, so a stray non-matching value
// here is defensive, not expected. ok is false when NO candidate parses (an
// empty list, or a list of only unparseable entries) — the caller MUST
// fail-closed (skip the auto-pre0 cut) on ok=false rather than treat "nothing
// to compare against" as "anything advances": a missed pre0 is recoverable by
// hand, a wrongly-cut lower one publishes a regression.
func HighestKnownEdgeVersion(candidates []string) (version string, ok bool) {
	for _, c := range candidates {
		v := strings.TrimPrefix(strings.TrimSpace(c), "v")
		if v == "" {
			continue
		}
		if _, _, err := parsePreVersion(v); err != nil {
			continue
		}
		if !ok {
			version, ok = v, true
			continue
		}
		if cmp, err := ComparePreVersion(v, version); err == nil && cmp > 0 {
			version = v
		}
	}
	return version, ok
}

// Pre0EdgeVersion computes the auto-cut edge prerelease version for a latest-line
// stable cut — X.(Y+1).0-pre0 — derived from the SAME dev-preversion the stable
// path stamps into `next` (X.(Y+1).0-pre1), swapping only the trailing label. So
// the pre0 tag's minor equals the required binary minor next's skills demand BY
// CONSTRUCTION, and the pre0 edge binary passes the minor-exact boot gate under
// the pre1-stamped skills (AC-1). release.yml's stable path prepends `v` to form
// the annotated auto-tag.
func Pre0EdgeVersion(stableVersion string) (string, error) {
	dev, err := DevPreVersion(stableVersion)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(dev, "-pre1") + "-pre0", nil
}
