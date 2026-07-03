// ABOUTME: Release-time guard predicate — the tagged commit's plugin.json version
// ABOUTME: must equal the tag's semver, so a tag/manifest mismatch blocks the cut.
package release

import (
	"fmt"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/contract"
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

// EvaluateProseMinorTagGate is the D5 prose half of the tag gate: it compares
// the git tag's major.minor against the FO shared-core's stamped minor literal
// (proseMinor, read via ProseMinor) — MINOR-only, since the prose carries no
// patch. Same prerelease carve-out as the JSON gate (release.yml skips this
// check entirely on a `-pre` tag; this function stays strict for defense).
func EvaluateProseMinorTagGate(tagVersion, proseMinor string) ManifestTagDecision {
	tag := strings.TrimPrefix(strings.TrimSpace(tagVersion), "v")
	tagMajor, tagMinor, ok := contract.ParseMajorMinor(tag)
	if !ok {
		return ManifestTagDecision{
			Reason: fmt.Sprintf("tag %s does not parse as major.minor", tagVersion),
		}
	}
	tagMinorStr := fmt.Sprintf("%d.%d", tagMajor, tagMinor)
	prose := strings.TrimSpace(proseMinor)
	if prose == "" {
		return ManifestTagDecision{
			Reason: fmt.Sprintf("prose file has no stamped minor; tag %s expects %s", tagVersion, tagMinorStr),
		}
	}
	if tagMinorStr != prose {
		return ManifestTagDecision{
			Reason: fmt.Sprintf("tag %s (minor %s) does not match prose-stamped minor %s; stamp the prose to %s and tag THAT commit", tagVersion, tagMinorStr, prose, tagMinorStr),
		}
	}
	return ManifestTagDecision{
		Pass:   true,
		Reason: fmt.Sprintf("tag %s matches prose-stamped minor %s", tagVersion, prose),
	}
}
