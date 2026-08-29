// ABOUTME: Derives the Spacedock Pi package install source from the running
// launcher's release identity — release-shaped binaries pin, dev builds float.
package cli

import (
	"regexp"
	"runtime/debug"
	"strings"
)

// piReleaseRef returns the git ref a release-shaped binary pins the Spacedock
// Pi package source to, resolved in order (ideation, pin-pi-package-to-binary-
// release):
//
//  1. the linker-stamped internal/cli.Version — release artifacts set it
//     (`-X ...Version={{ .Version }}`, stable and edge channels alike);
//  2. when Version is the "dev" sentinel, the Go build-info main module
//     version when it is a semver tag — this covers `go install …@vX.Y.Z`
//     proxy builds, which carry no ldflags but ARE release-shaped (the module
//     proxy embeds the tagged manifest, which is why displayVersion reports
//     them as X.Y.Z+dev);
//  3. otherwise "" — the dev sentinel: a plain `go build` / `go install
//     ./cmd/spacedock` checkout build has no release identity, keeps the
//     unpinned floating source, and performs no repair (no pin target).
//
// A release-shaped binary (stamped or proxy-tagged) never floats; the returned
// ref is the git tag the package source pins to.
func piReleaseRef() string {
	return piReleaseRefFrom(Version, piBuildInfoMainVersion())
}

// piBuildInfoMainVersion is the seam for the Go build-info main module version;
// a var so tests can stub the proxy-build identity.
var piBuildInfoMainVersion = func() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return strings.TrimSpace(bi.Main.Version)
}

// piReleaseRefFrom is the pure derivation over (linker stamp, build-info main
// version). A value counts as a release ref when it parses as a semver version
// (optionally "v"-prefixed, optional -preN prerelease) and is not a Go module
// pseudo-version (which go install emits for branch/HEAD resolutions like
// @next — those are dev-shaped and must float).
func piReleaseRefFrom(linkerStamp, buildInfoVersion string) string {
	if ref := piSemverTag(linkerStamp); ref != "" {
		return ref
	}
	// A proxy build self-reports Version == "dev" but embeds the tagged
	// manifest; the build info carries the resolved module version. A checkout
	// build reports "" or "(devel)" here, which the semver gate rejects.
	if bi := buildInfoVersion; bi != "" && bi != "(devel)" && bi != "devel" {
		return piSemverTag(bi)
	}
	return ""
}

func piSemverTag(v string) string {
	v = strings.TrimSpace(v)
	if piPseudoVersion.MatchString(v) {
		return "" // Go module pseudo-version (branch/HEAD resolution): dev-shaped
	}
	v = strings.TrimPrefix(v, "v")
	if v == "" {
		return ""
	}
	base := v
	if i := strings.IndexAny(v, "-+"); i >= 0 {
		base = v[:i]
	}
	parts := strings.Split(base, ".")
	if len(parts) != 3 {
		return ""
	}
	for _, p := range parts {
		if p == "" || len(p) > 9 {
			return ""
		}
		for _, r := range p {
			if r < '0' || r > '9' {
				return ""
			}
		}
	}
	// Go module pseudo-versions resolve no tag; the dev sentinel applies.
	// (Matched above before normalization.)
	return "v" + v
}

// piPseudoVersion matches Go module pseudo-versions — the resolution go
// install emits for branch/HEAD specs like @next (vX.Y.Z-0.<ts>-<hash>, and
// the base-tag form vX.Y.Z-pre.0.<ts>-<hash>). Both end in a 14-digit
// timestamp followed by the commit hash; they resolve no tag, so they carry
// no release identity: the dev sentinel applies.
var piPseudoVersion = regexp.MustCompile(`\d{14}-[0-9a-f]{7,12}(?:\+[^ ]+)?$`)

// piPinnedSource returns the install source for the Spacedock Pi package given
// the binary's release ref: a ref-pinned source for a release-shaped binary
// (never floats across tags), or the bare floating source for the dev sentinel.
func piPinnedSource(releaseRef string) string {
	if releaseRef == "" {
		return piSpacedockPackageSource
	}
	return piSpacedockPackageSource + "@" + releaseRef
}

// piGitSourceRef parses a Pi settings `packages` entry as a git source. ok is
// true only for `git:` sources; ref is the entry's @ref ("" when the git
// source is unpinned). Non-git entries (file:, npm:, local paths) return
// ok=false — they are user-managed and the repair never rewrites them.
func piGitSourceRef(source string) (ref string, ok bool) {
	const prefix = "git:"
	if !strings.HasPrefix(source, "git:") {
		return "", false
	}
	spec := strings.TrimPrefix(source, "git:")
	if i := strings.LastIndex(spec, "@"); i > 0 {
		return spec[i+1:], true
	}
	return "", true
}

// piPackageNeedsRepair reports whether the front door's one repair attempt
// must run for the discovered package status under a binary whose release ref
// is releaseRef. The repair triggers only for a release-shaped binary whose
// entry is missing, unpinned, or pinned to another ref — the ref delta IS the
// version mismatch. A non-git entry is user-managed (never rewritten), and a
// dev-sentinel binary never repairs (no pin target; a floating reinstall would
// clobber a pinned entry with an unpinned one).
func piPackageNeedsRepair(st piPackageStatus, releaseRef string) (bool, string) {
	if releaseRef == "" {
		return false, "" // dev sentinel: no pin target, never repairs
	}
	if !st.registered {
		return true, "missing"
	}
	entryRef, isGit := piGitSourceRef(st.source)
	if !isGit {
		return false, "" // user-managed entry: no clobber in either direction
	}
	if entryRef == "" {
		return true, "unpinned"
	}
	if entryRef != releaseRef {
		return true, "wrong release line"
	}
	return false, ""
}
