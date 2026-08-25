// ABOUTME: Stable-regression guard predicate — a bare stable tag must not be
// ABOUTME: older than the release version that the `stable` channel ref carries.
package release

import (
	"fmt"
	"strings"
)

// StableRegressionDecision is the result of the comparison between a release tag
// and the release that the stable channel points at now. Pass tells the caller
// if the cut can continue. Reason gives the text for the step log.
type StableRegressionDecision struct {
	Pass   bool
	Reason string
}

// EvaluateStableRegressionGate compares a bare stable tag against the version in
// the plugin manifest that `refs/heads/stable` points at. The gate blocks a
// STRICTLY older tag only. An equal version must pass, because a re-run of a
// release that already reached `stable` is a supported recovery. This is the one
// place where the boundary differs from EdgeAdvanceDecision's strict `>`.
//
// The function gives an error for a tag that is not a bare X.Y.Z, and for a
// stable version that it cannot parse. A gate must fail loud on input that it
// cannot read. A silent pass here lets an old tag reach goreleaser, which moves
// /releases/latest and the stable Homebrew cask DOWN. No job re-run can undo a
// published release or a pushed cask commit.
func EvaluateStableRegressionGate(tagVersion, stableVersion string) (StableRegressionDecision, error) {
	tag := strings.TrimPrefix(strings.TrimSpace(tagVersion), "v")
	if _, pre, err := parsePreVersion(tag); err != nil {
		return StableRegressionDecision{}, fmt.Errorf("release tag: %w", err)
	} else if pre != "" {
		return StableRegressionDecision{}, fmt.Errorf(
			"release tag %q is a prerelease; the gate step runs on a bare vX.Y.Z tag only", tagVersion)
	}
	stable := strings.TrimPrefix(strings.TrimSpace(stableVersion), "v")
	cmp, err := ComparePreVersion(tag, stable)
	if err != nil {
		return StableRegressionDecision{}, fmt.Errorf("stable channel version: %w", err)
	}
	if cmp < 0 {
		return StableRegressionDecision{Reason: fmt.Sprintf(
			"tag %s is older than %s, the release that the stable channel points at; "+
				"a cut of this tag moves /releases/latest and the stable Homebrew cask DOWN",
			tagVersion, stable)}, nil
	}
	return StableRegressionDecision{Pass: true, Reason: fmt.Sprintf(
		"tag %s is not older than %s, the release that the stable channel points at",
		tagVersion, stable)}, nil
}
