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

// EdgeAdvanceDecision decides whether the edge-advance job should ADVANCE `next`
// (reconcile + stamp + calendar bump + auto-cut pre0) or SKIP the whole job for
// the given tag, given `next`'s current manifest version. It computes the tag's
// TARGET EDGE VERSION — the version `next` would carry after this tag's
// reconcile: dev-preversion(tag) for a bare stable vX.Y.Z (`next` is stamped
// PAST the release), or the tag's own version for a `-pre` tag (`next` inherits
// it through the `-X theirs` reconcile) — and returns advance iff that target is
// STRICTLY GREATER than nextVersion. A patch cut from an older line yields a
// target at or below `next`'s current line, so it skips, leaving `next`'s tip
// untouched: no `-X theirs` clobber, no manifest/gate-line rewind, no colliding
// pre0 tag, no calendar re-pull. The boundary is strict `>`, not `>=`: a patch
// whose target exactly equals `next` (dev-preversion(vX.Y.1) == next) must skip.
func EdgeAdvanceDecision(tag, nextVersion string) (advance bool, targetEdge string, err error) {
	version := strings.TrimPrefix(tag, "v")
	if strings.Contains(version, "-") {
		// A prerelease tag: `next` would inherit the tag's own version.
		targetEdge = version
	} else {
		// A bare stable tag: `next` is stamped PAST it to the dev pre-version.
		targetEdge, err = DevPreVersion(version)
		if err != nil {
			return false, "", err
		}
	}
	cmp, err := ComparePreVersion(targetEdge, nextVersion)
	if err != nil {
		return false, "", err
	}
	return cmp > 0, targetEdge, nil
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
