package release

import (
	"strings"
	"testing"
)

// The manifest/tag agreement guard compares two INDEPENDENT values: the git tag
// semver and the version stamped into the tagged commit's plugin.json. They come
// from different sources (the annotated tag vs the manifest blob) and the real
// release history proves they can disagree — so this is a divergeable check, not
// a tautology. The fixtures below encode that history: the v0.20.0 cut tagged a
// commit whose plugin.json still read 0.19.9 (the pre-stamp/post-stamp inversion
// the reconciled releasing.md removes), whereas the v0.22.0 cut tagged a commit
// whose plugin.json already read 0.22.0.
const (
	divergedTag      = "v0.20.0"
	divergedManifest = "0.19.9"
	agreedTag        = "v0.22.0"
	agreedManifest   = "0.22.0"
)

// TestManifestTagGateBlocksDivergedHistory locks the block case against the real
// v0.20.0 history: a tag whose semver does not equal the tagged commit's manifest
// version must NOT pass — that is exactly the ungated pre-stamp inversion the
// reconciled procedure forbids. The two values are independent (tag object vs
// manifest blob), so a manifest left at the prior release reds the gate.
func TestManifestTagGateBlocksDivergedHistory(t *testing.T) {
	dec := EvaluateManifestTagGate(divergedTag, divergedManifest)
	if dec.Pass {
		t.Fatalf("gate passed a tag/manifest mismatch (%s vs %s); the v0.20.0 inversion must block", divergedTag, divergedManifest)
	}
	if !strings.Contains(dec.Reason, divergedTag) || !strings.Contains(dec.Reason, divergedManifest) {
		t.Errorf("block reason does not cite both diverging values; got: %s", dec.Reason)
	}
}

// TestManifestTagGatePassesWhenStamped locks the pass case against the real
// v0.22.0 history: when the tagged commit's manifest already carries the tag's
// semver (the stamp-then-tag ordering the reconciled procedure documents), the
// gate passes.
func TestManifestTagGatePassesWhenStamped(t *testing.T) {
	dec := EvaluateManifestTagGate(agreedTag, agreedManifest)
	if !dec.Pass {
		t.Fatalf("gate blocked a matching tag/manifest pair (%s vs %s): %s", agreedTag, agreedManifest, dec.Reason)
	}
}

// TestManifestTagGateNormalizesTagPrefix proves the comparison is on semver, not
// raw strings: the `v` prefix on the git tag is stripped before the compare, so a
// `vX.Y.Z` tag matches an `X.Y.Z` manifest. Without normalization every release
// would falsely diverge.
func TestManifestTagGateNormalizesTagPrefix(t *testing.T) {
	if dec := EvaluateManifestTagGate("v1.2.3", "1.2.3"); !dec.Pass {
		t.Fatalf("gate treated the tag `v` prefix as a mismatch: %s", dec.Reason)
	}
	if dec := EvaluateManifestTagGate("1.2.3", "1.2.3"); !dec.Pass {
		t.Fatalf("gate blocked a bare (already-normalized) tag/manifest match: %s", dec.Reason)
	}
}

// TestManifestTagGateBlocksOffByOne proves the off-by-one inversion (tag one ahead
// of the manifest, the exact shape of the v0.20.1..v0.20.3 post-stamp lag) blocks.
func TestManifestTagGateBlocksOffByOne(t *testing.T) {
	if dec := EvaluateManifestTagGate("v0.20.1", "0.20.0"); dec.Pass {
		t.Fatalf("gate passed a tag one minor-patch ahead of the manifest (the post-stamp lag inversion)")
	}
}

// TestManifestTagGateBlocksEmptyManifest proves a manifest with no resolvable
// version (a stamp that never ran) blocks rather than silently passing.
func TestManifestTagGateBlocksEmptyManifest(t *testing.T) {
	if dec := EvaluateManifestTagGate("v0.22.0", ""); dec.Pass {
		t.Fatalf("gate passed an empty manifest version")
	}
}
