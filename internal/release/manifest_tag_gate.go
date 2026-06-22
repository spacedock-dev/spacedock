// ABOUTME: Release-time guard predicate — the tagged commit's plugin.json version
// ABOUTME: must equal the tag's semver, so a tag/manifest mismatch blocks the cut.
package release

import (
	"fmt"
	"strings"
)

// ManifestTagDecision is the outcome of comparing a release tag's semver against
// the version stamped into the tagged commit's plugin manifest. Pass gates whether
// the cut may proceed; Reason explains the match or the mismatch for the step log.
type ManifestTagDecision struct {
	Pass   bool
	Reason string
}

// EvaluateManifestTagGate decides whether a release may be cut by comparing two
// INDEPENDENT values: the git tag's semver (tagVersion, with or without a leading
// `v`) and the version field read from the tagged commit's plugin.json
// (manifestVersion). They originate from different artifacts — the annotated tag
// and the manifest blob — and the real release history proves they can diverge
// (the v0.20.0 cut tagged a commit whose manifest still read 0.19.9). The gate
// passes only when the normalized tag semver equals the manifest version, so a
// pre-stamp/post-stamp inversion that leaves the manifest at the prior release is
// caught before goreleaser fires. An empty manifest version (a stamp that never
// ran) never passes.
func EvaluateManifestTagGate(tagVersion, manifestVersion string) ManifestTagDecision {
	tag := strings.TrimPrefix(strings.TrimSpace(tagVersion), "v")
	manifest := strings.TrimSpace(manifestVersion)
	if manifest == "" {
		return ManifestTagDecision{
			Reason: fmt.Sprintf("tagged commit's plugin.json has no version; tag %s expects %s", tagVersion, tag),
		}
	}
	if tag != manifest {
		return ManifestTagDecision{
			Reason: fmt.Sprintf("tag %s does not match tagged commit's plugin.json version %s; stamp the manifest to %s and tag THAT commit", tagVersion, manifest, tag),
		}
	}
	return ManifestTagDecision{
		Pass:   true,
		Reason: fmt.Sprintf("tag %s matches tagged commit's plugin.json version %s", tagVersion, manifest),
	}
}
