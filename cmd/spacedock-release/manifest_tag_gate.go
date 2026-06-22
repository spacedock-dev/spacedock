// ABOUTME: `spacedock-release manifest-tag-gate <tag> <plugin.json>...` — blocks
// ABOUTME: the cut unless every tagged manifest's version equals the tag semver.
package main

import (
	"fmt"
	"os"

	"github.com/spacedock-dev/spacedock/internal/release"
)

// runManifestTagGate asserts the tag semver equals the version stamped into each
// plugin manifest the cutter is about to tag. It reads the manifest version from
// each file (an independent value from the tag), runs the pure decision predicate,
// records the outcome to $GITHUB_STEP_SUMMARY, and returns the process exit code:
// 0 when every manifest matches the tag, 1 when any diverges or cannot be read, 2
// on a usage error. This is the divergeable guard behind the reconciled
// releasing.md's stamp-then-tag ordering: a tag-vs-manifest mismatch (the pre-stamp
// inversion) is caught before goreleaser fires.
func runManifestTagGate(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "spacedock-release manifest-tag-gate: need <tag> <plugin.json> [<plugin.json> ...]")
		return 2
	}
	tag, manifests := args[0], args[1:]
	for _, path := range manifests {
		data, err := os.ReadFile(path)
		if err != nil {
			reason := fmt.Sprintf("manifest-tag gate BLOCKED for %s: cannot read %s (%v)", tag, path, err)
			fmt.Fprintln(os.Stderr, reason)
			recordManifestTagSummary(reason)
			return 1
		}
		version, err := release.ManifestVersion(data)
		if err != nil {
			reason := fmt.Sprintf("manifest-tag gate BLOCKED for %s: cannot parse %s (%v)", tag, path, err)
			fmt.Fprintln(os.Stderr, reason)
			recordManifestTagSummary(reason)
			return 1
		}
		dec := release.EvaluateManifestTagGate(tag, version)
		if !dec.Pass {
			reason := fmt.Sprintf("%s (%s)", dec.Reason, path)
			fmt.Fprintf(os.Stderr, "spacedock-release manifest-tag-gate: %s\n", reason)
			recordManifestTagSummary("manifest-tag gate BLOCKED: " + reason)
			return 1
		}
		fmt.Printf("%s (%s)\n", dec.Reason, path)
	}
	recordManifestTagSummary(fmt.Sprintf("manifest-tag gate PASSED for %s: all %d manifest(s) match the tag semver", tag, len(manifests)))
	return 0
}

// recordManifestTagSummary appends the gate decision to $GITHUB_STEP_SUMMARY when
// set, so every cut leaves an auditable record. Outside CI the env var is unset
// and this is a no-op.
func recordManifestTagSummary(reason string) {
	path := os.Getenv("GITHUB_STEP_SUMMARY")
	if path == "" {
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "### Release manifest-tag gate\n\n%s\n", reason)
}
