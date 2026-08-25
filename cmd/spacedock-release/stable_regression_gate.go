// ABOUTME: `spacedock-release stable-regression-gate <tag> <stable-manifest>` —
// ABOUTME: blocks a cut whose tag is older than the release `stable` points at.
package main

import (
	"fmt"
	"os"

	"github.com/spacedock-dev/spacedock/internal/release"
)

// runStableRegressionGate reads the plugin manifest that the `stable` channel ref
// points at, and compares its version against the release tag. It records the
// result to $GITHUB_STEP_SUMMARY. The exit code is 0 when the tag is the same
// version or newer. It is 1 when the tag is older, or the manifest is
// unreadable. It is 2 on a usage error.
//
// release.yml runs this in the e2e-gate job, which goreleaser needs. A non-zero
// exit therefore stops the cut before goreleaser publishes the release and bumps
// the stable Homebrew cask.
func runStableRegressionGate(args []string) int {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "spacedock-release stable-regression-gate: need <tag> <stable-plugin.json>")
		return 2
	}
	tag, path := args[0], args[1]
	data, err := os.ReadFile(path)
	if err != nil {
		return blockStableRegression(fmt.Sprintf(
			"stable-regression gate BLOCKED for %s: cannot read %s (%v)", tag, path, err))
	}
	stableVersion, err := release.ManifestVersion(data)
	if err != nil {
		return blockStableRegression(fmt.Sprintf(
			"stable-regression gate BLOCKED for %s: cannot parse %s (%v)", tag, path, err))
	}
	dec, err := release.EvaluateStableRegressionGate(tag, stableVersion)
	if err != nil {
		return blockStableRegression(fmt.Sprintf("stable-regression gate BLOCKED for %s: %v", tag, err))
	}
	if !dec.Pass {
		return blockStableRegression("stable-regression gate BLOCKED: " + dec.Reason)
	}
	fmt.Println(dec.Reason)
	recordStableRegressionSummary("stable-regression gate PASSED: " + dec.Reason)
	return 0
}

// blockStableRegression writes reason to stderr and to the step summary, then
// gives the exit code that stops the cut.
func blockStableRegression(reason string) int {
	fmt.Fprintln(os.Stderr, reason)
	recordStableRegressionSummary(reason)
	return 1
}

// recordStableRegressionSummary appends the gate result to $GITHUB_STEP_SUMMARY
// when that variable is set, so each cut keeps an auditable record. Outside CI
// the variable is empty and this function does nothing.
func recordStableRegressionSummary(reason string) {
	path := os.Getenv("GITHUB_STEP_SUMMARY")
	if path == "" {
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "### Release stable-regression gate\n\n%s\n", reason)
}
