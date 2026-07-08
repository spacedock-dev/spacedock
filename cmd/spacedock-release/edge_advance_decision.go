// ABOUTME: `edge-advance-decision` prints advance/skip gating the whole edge-advance
// ABOUTME: job; `edge-pre0-version` prints the auto-cut edge prerelease version.
package main

import (
	"fmt"
	"os"

	"github.com/spacedock-dev/spacedock/internal/release"
)

// edgeAdvanceDecision decides whether a tag advances the edge line or skips the
// whole edge-advance job, given `next`'s current manifest. It prints `advance`
// or `skip` to stdout (exit 0 for BOTH — a skip is a normal outcome, not an
// error) so release.yml's decision step can write advance=true|false to
// $GITHUB_OUTPUT and gate every downstream step on it. A missing/unparseable
// manifest or a malformed tag/version is a real error (non-zero), so a miswiring
// fails loud rather than silently skipping and leaving the edge channel broken.
func edgeAdvanceDecision(args []string) int {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "spacedock-release edge-advance-decision: need <tag> <next-plugin.json>")
		return 2
	}
	tag, manifestPath := args[0], args[1]
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", manifestPath, err)
		return 1
	}
	nextVersion, err := release.ManifestVersion(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", manifestPath, err)
		return 1
	}
	if nextVersion == "" {
		fmt.Fprintf(os.Stderr, "edge-advance-decision: %s carries no top-level version\n", manifestPath)
		return 1
	}
	advance, target, err := release.EdgeAdvanceDecision(tag, nextVersion)
	if err != nil {
		fmt.Fprintf(os.Stderr, "edge-advance-decision: %v\n", err)
		return 1
	}
	fmt.Fprintf(os.Stderr, "tag %s target edge version %s vs next %s\n", tag, target, nextVersion)
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
