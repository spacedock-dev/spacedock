// ABOUTME: The `.md` half of stamp-version — rewrites the FO shared-core's single
// ABOUTME: pinned "required binary minor" literal to the release's major.minor.
package release

import (
	"fmt"
	"regexp"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// proseMinorPattern matches the ONE pinned literal D5 requires the FO shared-core
// prose to carry: "These skills require binary minor <maj.min>". Every other
// clause in the gate step refers to "the required minor" instead of repeating
// the number, so this is the single point a hand edit or a stamp touches.
var proseMinorPattern = regexp.MustCompile(`(These skills require binary minor )(\d+\.\d+)`)

// StampProseVersion rewrites the single line in doc matching proseMinorPattern
// to the major.minor derived from releaseVersion, leaving the rest of the file
// untouched. It errors when releaseVersion has no parseable major.minor, or when
// the pattern matches zero or more than one time — the atomic-by-construction
// guarantee D5 requires: one literal, one whole-file rewrite, so no partial
// rewrite can yield self-contradictory instructions.
func StampProseVersion(doc []byte, releaseVersion string) ([]byte, error) {
	major, minor, ok := contract.ParseMajorMinor(releaseVersion)
	if !ok {
		return nil, fmt.Errorf("release version %q has no parseable major.minor", releaseVersion)
	}
	matches := proseMinorPattern.FindAllIndex(doc, -1)
	if len(matches) != 1 {
		return nil, fmt.Errorf(
			"expected exactly one %q line, found %d",
			"These skills require binary minor <maj.min>", len(matches))
	}
	replacement := []byte(fmt.Sprintf("${1}%d.%d", major, minor))
	return proseMinorPattern.ReplaceAll(doc, replacement), nil
}

// ProseMinor reads the stamped major.minor literal ("X.Y") back out of doc — the
// same pinned pattern StampProseVersion writes. Used by the tag gate and the
// internal/contractlint sync test to read the prose side of the prose/manifest
// binding. Errors under the same zero-or-multiple-match condition
// StampProseVersion enforces on write.
func ProseMinor(doc []byte) (string, error) {
	matches := proseMinorPattern.FindAllSubmatch(doc, -1)
	if len(matches) != 1 {
		return "", fmt.Errorf(
			"expected exactly one %q line, found %d",
			"These skills require binary minor <maj.min>", len(matches))
	}
	return string(matches[0][2]), nil
}
