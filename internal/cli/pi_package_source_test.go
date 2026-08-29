package cli

import "testing"

// TestPiReleaseRefDerivation pins the ordered ref derivation (linker stamp,
// then build-info main version, then the dev sentinel) for every identity
// class the derivation must distinguish. Falsifier: changing any derivation
// arm flips the row for that identity — notably the proxy-build row (the
// go-install @vX.Y.Z build that must pin) and the pseudo-version rows (go
// install @next resolutions, which must float).
func TestPiReleaseRefDerivation(t *testing.T) {
	const pseudo1 = "v0.28.0-pre0.0.20260828120000-abcdef123456" // @next over a pre tag
	const pseudo2 = "v0.27.1-0.20260828120000-abcdef123456"      // no base tag
	const dirtyPseudo = "v0.28.0-pre0.0.20260828165724-81e3386e8234+dirty"
	for _, tc := range []struct {
		name, stamp, buildInfo, want string
	}{
		{"linker stamp wins", "0.27.2", "v0.26.0", "v0.27.2"},
		{"v-prefixed stamp normalizes", "v0.27.2", "", "v0.27.2"},
		{"prerelease stamp", "0.28.0-pre1", "(devel)", "v0.28.0-pre1"},
		{"proxy build pins its tag", "dev", "v0.27.2", "v0.27.2"},
		{"checkout build is the dev sentinel", "dev", "(devel)", ""},
		{"empty build info is the dev sentinel", "dev", "", ""},
		{"next pseudo-version floats", "dev", pseudo1, ""},
		{"plain pseudo-version floats", "dev", pseudo2, ""},
		{"dirty checkout pseudo-version floats", "dev", dirtyPseudo, ""},
		{"non-semver build info floats", "dev", "not-a-version", ""},
		{"empty stamp uses build info", "", "v0.27.2", "v0.27.2"},
	} {
		if got := piReleaseRefFrom(tc.stamp, tc.buildInfo); got != tc.want {
			t.Fatalf("%s: piReleaseRefFrom(%q,%q)=%q want %q", tc.name, tc.stamp, tc.buildInfo, got, tc.want)
		}
	}
}

// TestPiPinnedSourceAndRefParsing pins the source-string construction and the
// git/non-git entry classifier the repair trigger and AC-5's no-clobber
// guarantee depend on. Falsifier: dropping the @ref suffix, floating a pinned
// ref, or classifying a git entry as non-git (or vice versa) flips a row.
func TestPiPinnedSourceAndRefParsing(t *testing.T) {
	const repo = "git:github.com/spacedock-dev/spacedock"
	for _, tc := range []struct{ ref, want string }{
		{"v0.27.2", repo + "@v0.27.2"},
		{"v0.28.0-pre0", repo + "@v0.28.0-pre0"},
		{"", repo}, // dev sentinel keeps the bare floating source
	} {
		if got := piPinnedSource(tc.ref); got != tc.want {
			t.Fatalf("piPinnedSource(%q)=%q want %q", tc.ref, got, tc.want)
		}
	}
	for _, tc := range []struct {
		source, ref string
		isGit       bool
	}{
		{"git:github.com/spacedock-dev/spacedock", "", true},
		{"git:github.com/spacedock-dev/spacedock@v0.27.2", "v0.27.2", true},
		{"file:/tmp/repo", "", false},
		{"npm:spacedock", "", false},
		{"/usr/local/spacedock", "", false},
		{"", "", false},
	} {
		if ref, ok := piGitSourceRef(tc.source); ref != tc.ref || ok != tc.isGit {
			t.Fatalf("piGitSourceRef(%q)=(%q,%v) want (%q,%v)", tc.source, ref, ok, tc.ref, tc.isGit)
		}
	}
}

// TestPiPackageNeedsRepairTable pins the trigger decision: a release-shaped
// binary repairs a missing, unpinned, or wrong-ref git entry; its own ref,
// any non-git entry, and the dev sentinel never trigger.
func TestPiPackageNeedsRepairTable(t *testing.T) {
	missing := piPackageStatus{}
	unpinned := piPackageStatus{registered: true, source: "git:github.com/spacedock-dev/spacedock"}
	wrongLine := piPackageStatus{registered: true, source: "git:github.com/spacedock-dev/spacedock@v0.28.0-pre1"}
	ownRef := piPackageStatus{registered: true, source: "git:github.com/spacedock-dev/spacedock@v0.27.2"}
	userEntry := piPackageStatus{registered: true, source: "file:/x"}
	for _, tc := range []struct {
		name       string
		st         piPackageStatus
		relRef     string
		want       bool
		wantReason string
	}{
		{"missing entry, release-shaped", missing, "v0.27.2", true, "missing"},
		{"unpinned git entry", unpinned, "v0.27.2", true, "unpinned"},
		{"wrong release line", wrongLine, "v0.27.2", true, "wrong release line"},
		{"pinned to own ref", ownRef, "v0.27.2", false, ""},
		{"non-git entry is user-managed", userEntry, "v0.27.2", false, ""},
		{"dev sentinel never repairs", missing, "", false, ""},
	} {
		got, reason := piPackageNeedsRepair(tc.st, tc.relRef)
		if got != tc.want || (got && reason != tc.wantReason) {
			t.Fatalf("%s: repair=(%v,%q) want (%v,%q)", tc.name, got, reason, tc.want, tc.wantReason)
		}
	}
}
