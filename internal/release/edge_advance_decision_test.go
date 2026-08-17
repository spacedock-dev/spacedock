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

// TestHighestKnownEdgeVersionIncludesPrereleases is the round-3 replacement
// design's load-bearing case (team-lead's worked example): a candidate pool
// mixing a bare stable tag and a newer PRERELEASE tag must find the
// prerelease — filtering to stable-only would find the wrong, lower "known"
// version and let an old-line patch wrongly out-rank it. Also proves the "v"
// prefix is optional and stripped consistently.
func TestHighestKnownEdgeVersionIncludesPrereleases(t *testing.T) {
	got, ok := HighestKnownEdgeVersion([]string{"v0.26.0", "v0.27.0-pre7", "v0.25.1"})
	if !ok {
		t.Fatal("HighestKnownEdgeVersion reported ok=false on a valid candidate list")
	}
	if got != "0.27.0-pre7" {
		t.Fatalf("HighestKnownEdgeVersion = %q, want 0.27.0-pre7 (the prerelease must win over the bare stable v0.26.0)", got)
	}

	// pre9 vs pre10 must order numerically, not lexically — the exact distinction
	// ComparePreVersion exists for; a naive string-max would get this backwards.
	got, ok = HighestKnownEdgeVersion([]string{"v0.27.0-pre2", "v0.27.0-pre10", "v0.27.0-pre9"})
	if !ok || got != "0.27.0-pre10" {
		t.Fatalf("HighestKnownEdgeVersion(pre2,pre10,pre9) = (%q,%v), want (0.27.0-pre10,true)", got, ok)
	}
}

// TestHighestKnownEdgeVersionSkipsMalformed proves a candidate list mixing
// garbage/non-release entries (the shape a broader `git tag --list 'v*'` scan
// can produce — non-semver tags, empty strings) with valid ones still finds the
// correct highest, ignoring what doesn't parse rather than erroring the whole
// scan.
func TestHighestKnownEdgeVersionSkipsMalformed(t *testing.T) {
	got, ok := HighestKnownEdgeVersion([]string{"not-a-version", "", "v0.26.0", "vSomeOtherTag", "v0.25.1"})
	if !ok {
		t.Fatal("HighestKnownEdgeVersion reported ok=false with a valid candidate present")
	}
	if got != "0.26.0" {
		t.Fatalf("HighestKnownEdgeVersion = %q, want 0.26.0 (malformed entries must be skipped, not fatal)", got)
	}
}

// TestHighestKnownEdgeVersionFailsClosedWhenNothingParses is the round-3 new
// failure mode team-lead flagged: an empty candidate list, and a list of only
// malformed/non-release entries, must both report ok=false — the caller (the
// release.yml decision step) is responsible for treating that as "skip the
// auto-pre0 cut", never as "nothing to compare against, so anything advances".
// A missed pre0 is recoverable by hand; a wrongly-cut lower one is not.
func TestHighestKnownEdgeVersionFailsClosedWhenNothingParses(t *testing.T) {
	if _, ok := HighestKnownEdgeVersion(nil); ok {
		t.Fatal("HighestKnownEdgeVersion(nil) reported ok=true; want false (fail-closed)")
	}
	if _, ok := HighestKnownEdgeVersion([]string{}); ok {
		t.Fatal("HighestKnownEdgeVersion([]) reported ok=true; want false (fail-closed)")
	}
	if _, ok := HighestKnownEdgeVersion([]string{"not-a-version", "", "garbage-tag"}); ok {
		t.Fatal("HighestKnownEdgeVersion(all-malformed) reported ok=true; want false (fail-closed)")
	}
}

// TestOldLinePatchSkipsAgainstHighestKnownEdgeVersion is the ported successor
// to the retired TestEdgeAdvancePatchDoesNotRegressNext (deleted with
// edge_advance_noregress_test.go — the reconcile fixture it drove no longer
// exists). What that test proved and this one re-proves against the NEW
// mechanism: an old-line patch tag, cut after a newer prerelease line already
// exists, must SKIP — so the auto-pre0 step never fires and no colliding,
// lower pre0 tag is cut. This wires the two round-3 pieces end to end: a git-tag
// candidate pool (as release.yml's decision step would scan) feeds
// HighestKnownEdgeVersion, whose result feeds EdgeAdvanceDecision.
func TestOldLinePatchSkipsAgainstHighestKnownEdgeVersion(t *testing.T) {
	// Mirrors the deleted fixture's scenario: the 0.27 line is already ahead
	// (via a prerelease tag, not a next-branch manifest) when an old-line 0.25
	// patch is cut.
	candidates := []string{"v0.24.0", "v0.25.0", "v0.26.0", "v0.27.0-pre1", "v0.27.0-pre7"}
	patchTag := "v0.25.1"

	known, ok := HighestKnownEdgeVersion(candidates)
	if !ok || known != "0.27.0-pre7" {
		t.Fatalf("HighestKnownEdgeVersion(%v) = (%q,%v), want (0.27.0-pre7,true)", candidates, known, ok)
	}
	advance, target, err := EdgeAdvanceDecision(patchTag, known)
	if err != nil {
		t.Fatalf("EdgeAdvanceDecision(%q, %q) errored: %v", patchTag, known, err)
	}
	if advance {
		t.Fatalf("EdgeAdvanceDecision(%q, %q) advance = true, want false (old-line patch, target %q must not out-rank known %q — a colliding lower pre0 would be cut)", patchTag, known, target, known)
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
