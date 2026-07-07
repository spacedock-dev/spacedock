// ABOUTME: Unit proof for the edge-advance line-ordering guard — ComparePreVersion
// ABOUTME: ordering, the advance/skip decision table, and the pre0/required-minor algebra.
package release

import (
	"strings"
	"testing"
)

// TestComparePreVersionOrders locks the prerelease-aware ordering the decision
// guard needs and contract.semverCompare cannot give (it is dotted-int only and
// fails to parse `-preN`). Core X.Y.Z dominates; a release outranks the same
// core's prerelease (semver §11); two `-preN` order by NUMERIC N, so pre0 < pre1
// and pre2 < pre10 (a lexical compare would wrongly put pre10 before pre2).
func TestComparePreVersionOrders(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"0.26.0-pre1", "0.26.0", -1}, // release > prerelease of same core
		{"0.26.0", "0.26.0-pre1", 1},
		{"0.26.0-pre1", "0.26.0-pre2", -1},  // pre1 < pre2
		{"0.26.0-pre0", "0.26.0-pre1", -1},  // the recursion-skip's pre0 < pre1
		{"0.26.0-pre2", "0.26.0-pre10", -1}, // numeric, not lexical
		{"0.27.0-pre1", "0.26.0", 1},        // core dominates the prerelease label
		{"0.26.0", "0.27.0-pre1", -1},
		{"1.0.0", "0.27.0-pre1", 1}, // major dominates
		{"0.26.0-pre1", "0.26.0-pre1", 0},
		{"0.26.0", "0.26.0", 0},
	}
	for _, c := range cases {
		got, err := ComparePreVersion(c.a, c.b)
		if err != nil {
			t.Errorf("ComparePreVersion(%q, %q) errored: %v", c.a, c.b, err)
			continue
		}
		if got != c.want {
			t.Errorf("ComparePreVersion(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
		// Antisymmetry: swapping the operands negates the sign.
		rev, err := ComparePreVersion(c.b, c.a)
		if err != nil {
			t.Errorf("ComparePreVersion(%q, %q) errored: %v", c.b, c.a, err)
			continue
		}
		if rev != -c.want {
			t.Errorf("ComparePreVersion(%q, %q) = %d, want %d (antisymmetry)", c.b, c.a, rev, -c.want)
		}
	}
}

// TestComparePreVersionRejectsMalformed proves a version that is not X.Y.Z or
// X.Y.Z-<prerelease> errors rather than silently comparing as 0 (which would
// make the decision guard skip on garbage instead of failing loudly).
func TestComparePreVersionRejectsMalformed(t *testing.T) {
	for _, bad := range []string{"", "0.26", "0.26.0.1", "v0.26.0", "garbage", "0.x.0"} {
		if _, err := ComparePreVersion(bad, "0.26.0"); err == nil {
			t.Errorf("ComparePreVersion(%q, …) accepted a malformed version", bad)
		}
		if _, err := ComparePreVersion("0.26.0", bad); err == nil {
			t.Errorf("ComparePreVersion(…, %q) accepted a malformed version", bad)
		}
	}
}

// TestEdgeAdvanceDecision locks the guard's forward-fires / backward-skips table,
// including the three patch-equality skip-cases: a patch cut whose target edge
// version (dev-preversion of the patch) EQUALS next's current manifest must SKIP,
// because the boundary is strict `>` (an `>=` boundary would fire on that
// equality and clobber next via the `-X theirs` reconcile).
func TestEdgeAdvanceDecision(t *testing.T) {
	cases := []struct {
		name        string
		tag         string
		next        string
		wantAdvance bool
		wantTarget  string
	}{
		{"forward stable", "v0.26.0", "0.26.0-pre1", true, "0.27.0-pre1"},
		{"forward prerelease", "v0.26.0-pre2", "0.26.0-pre1", true, "0.26.0-pre2"},
		{"major bump", "v1.0.0", "0.27.0-pre1", true, "1.1.0-pre1"},
		{"old-line patch, next ahead", "v0.25.1", "0.27.0-pre1", false, "0.26.0-pre1"},
		{"pre0 vs next pre1", "v0.26.0-pre0", "0.26.0-pre1", false, "0.26.0-pre0"},
		// The three patch-equality skip-cases: dev-preversion(tag) == next.
		{"equality skip: v0.25.1 vs 0.26.0-pre1", "v0.25.1", "0.26.0-pre1", false, "0.26.0-pre1"},
		{"equality skip: v0.26.1 vs 0.27.0-pre1", "v0.26.1", "0.27.0-pre1", false, "0.27.0-pre1"},
		{"equality skip: v1.0.9 vs 1.1.0-pre1", "v1.0.9", "1.1.0-pre1", false, "1.1.0-pre1"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			advance, target, err := EdgeAdvanceDecision(c.tag, c.next)
			if err != nil {
				t.Fatalf("EdgeAdvanceDecision(%q, %q) errored: %v", c.tag, c.next, err)
			}
			if advance != c.wantAdvance {
				t.Errorf("EdgeAdvanceDecision(%q, %q) advance = %v, want %v", c.tag, c.next, advance, c.wantAdvance)
			}
			if target != c.wantTarget {
				t.Errorf("EdgeAdvanceDecision(%q, %q) target = %q, want %q", c.tag, c.next, target, c.wantTarget)
			}
		})
	}
}

// TestAutoPre0MinorEqualsRequiredMinor is AC-1's version-algebra proof: for a
// latest-line stable cut vX.Y.0, the required binary minor stamped into next's
// prose (dev-preversion → X.(Y+1)) and the auto-cut pre0 tag's own minor
// (Pre0EdgeVersion → X.(Y+1)) are EQUAL BY CONSTRUCTION — they share the same
// X.(Y+1).0 core, differing only in the -pre1/-pre0 label. So a pre0 edge binary
// under pre1-stamped skills passes the minor-exact boot gate. No goreleaser tag
// resolution is modeled here — the real binary-version measurement is the CI
// dry-run spike; this pins only the algebra.
func TestAutoPre0MinorEqualsRequiredMinor(t *testing.T) {
	for _, stable := range []string{"0.25.0", "0.26.0", "1.0.0", "2.9.0", "0.99.0"} {
		dev, err := DevPreVersion(stable)
		if err != nil {
			t.Fatalf("DevPreVersion(%q): %v", stable, err)
		}
		pre0, err := Pre0EdgeVersion(stable)
		if err != nil {
			t.Fatalf("Pre0EdgeVersion(%q): %v", stable, err)
		}
		if !strings.HasSuffix(dev, "-pre1") {
			t.Errorf("DevPreVersion(%q) = %q, want a -pre1 suffix", stable, dev)
		}
		if !strings.HasSuffix(pre0, "-pre0") {
			t.Errorf("Pre0EdgeVersion(%q) = %q, want a -pre0 suffix", stable, pre0)
		}
		devCore := strings.TrimSuffix(dev, "-pre1")
		pre0Core := strings.TrimSuffix(pre0, "-pre0")
		if devCore != pre0Core {
			t.Errorf("stable %q: required-minor core %q != pre0 core %q — the boot gate would mismatch", stable, devCore, pre0Core)
		}
	}
}
