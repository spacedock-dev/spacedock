// ABOUTME: `edge-advance-decision` prints advance/skip gating the auto-pre0 cut;
// ABOUTME: `edge-pre0-version` prints its version; `highest-known-edge-version`
// ABOUTME: folds a git-tag scan into the "known version" edge-advance-decision compares against.
package main

import (
	"fmt"
	"os"

	"github.com/spacedock-dev/spacedock/internal/release"
)

// edgeAdvanceDecision decides whether a tag advances the edge line (fires the
// auto-pre0 cut) or skips it, given a manifest carrying the highest known
// release version — release.yml synthesizes that manifest from a
// highest-known-edge-version scan of git tag history (see that command's own
// doc; pre-round-3 it read `next`'s live manifest, which is no longer a
// meaningful source). It prints `advance` or `skip` to stdout (exit 0 for BOTH
// — a skip is a normal outcome, not an error) so release.yml's decision step
// can write advance=true|false to $GITHUB_OUTPUT. A missing/unparseable
// manifest or a malformed tag/version is a real error (non-zero), so a miswiring
// fails loud rather than silently skipping and leaving the edge line behind.
func edgeAdvanceDecision(args []string) int {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "spacedock-release edge-advance-decision: need <tag> <known-version-plugin.json>")
		return 2
	}
	tag, manifestPath := args[0], args[1]
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", manifestPath, err)
		return 1
	}
	knownVersion, err := release.ManifestVersion(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", manifestPath, err)
		return 1
	}
	if knownVersion == "" {
		fmt.Fprintf(os.Stderr, "edge-advance-decision: %s carries no top-level version\n", manifestPath)
		return 1
	}
	advance, target, err := release.EdgeAdvanceDecision(tag, knownVersion)
	if err != nil {
		fmt.Fprintf(os.Stderr, "edge-advance-decision: %v\n", err)
		return 1
	}
	fmt.Fprintf(os.Stderr, "tag %s target edge version %s vs highest known %s\n", tag, target, knownVersion)
	if advance {
		fmt.Println("advance")
	} else {
		fmt.Println("skip")
	}
	return 0
}

// edgePre0Version prints the auto-cut edge prerelease version (X.(Y+1).0-pre0)
// for a bare stable version, so release.yml's always-cut-pre0 step forms the
// `vX.(Y+1).0-pre0` annotated tag. It errors on a hyphenated/malformed input,
// since it runs only on the stable path that already guarantees a bare X.Y.Z.
func edgePre0Version(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "spacedock-release edge-pre0-version: need exactly one <stable-version> (e.g. 0.25.0)")
		return 2
	}
	version, err := release.Pre0EdgeVersion(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "edge-pre0-version: %v\n", err)
		return 1
	}
	fmt.Println(version)
	return 0
}

// highestKnownEdgeVersion prints the greatest version among its <tag> arguments
// (release.yml feeds it a `git tag --list 'v*'` scan, minus the tag being
// decided — see release.HighestKnownEdgeVersion's doc for why prereleases are
// INCLUDED and malformed entries are tolerated), ALSO comparing against the
// dev-preversion of the highest bare stable tag in that same argument list —
// the "next"-manifest notch release.HighestBareStableVersion's doc explains
// (restores what the retired stamp always sat one notch above, so an old-line
// patch cut while the newer line has only a pre0 tag correctly skips instead
// of out-ranking it and re-cutting a colliding pre0). Prints NOTHING (empty
// stdout) with exit 0 when no argument parses as a release tag — the
// empty-list and all-malformed cases are normal, not errors, so a `set -e`
// caller does not abort the job; release.yml's decision step checks for an
// empty result and fails CLOSED (skips the auto-pre0 cut) rather than
// treating "nothing to compare against" as "anything advances".
func highestKnownEdgeVersion(args []string) int {
	version, ok := release.HighestKnownEdgeVersion(args)
	if bare, bareOK := release.HighestBareStableVersion(args); bareOK {
		if notch, err := release.DevPreVersion(bare); err == nil {
			if !ok {
				version, ok = notch, true
			} else if cmp, cerr := release.ComparePreVersion(notch, version); cerr == nil && cmp > 0 {
				version = notch
			}
		}
	}
	if ok {
		fmt.Println(version)
	}
	return 0
}
