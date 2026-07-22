// ABOUTME: AC-1 offline proof — every shared-scenario fixture README carries the
// ABOUTME: commissioned-by marker discovery requires; the zero-discovery fixture stays markerless.
package ensigncycle

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// fixtureDiscoveryReadmes names every live-runner fixture README whose scenario
// expects a discoverable workflow — the Spike promoted to a repo test. Before the
// commissioned-by marker landed on all of them, the FO's nondeterministic choice
// between `--boot --identify` (gated on this marker) and `--workflow-dir .`
// (marker-blind) could false-negative a scenario whose fixture lacked it; only
// filingReadme and smallestMechanismReadme carried it from the start. The
// zero-discovery fixture (TestLiveZeroDiscoverReportsAndStops) is the sole,
// deliberate exception — its SUBJECT is non-discovery — and is proved the other way
// in TestZeroDiscoverFixtureStaysUndiscoverable below, not listed here.
func fixtureDiscoveryReadmes() map[string]func() string {
	return map[string]func() string{
		"gate-guardrail":                gateReadme,
		"rejection-flow":                rejectionReadme,
		"feedback-3-cycle-escalation":   escalationReadme,
		"merge-hook-guardrail":          mergeHookGuardReadme,
		"filing":                        filingReadme,
		"shallow-boot":                  shallowBootReadme,
		"multi-workflow-boot":           func() string { return multiWorkflowBootReadme("alpha") },
		"self-evidence-merge-triage":    mergeTriageReadme,
		"smallest-sufficient-mechanism": smallestMechanismReadme,
		"keep-moving-posture":           keepMovingReadme,
		"gate-stop":                     gateStopReadme,
	}
}

// TestSharedScenarioFixturesAreDiscoverable is AC-1's CLI/deterministic proof: from
// every fixture root, discovery (the predicate `--boot --identify` resolves
// through) finds the workflow — so the FO's nondeterministic choice between
// `--boot --identify` and `--workflow-dir .` can no longer produce a
// scenario-failing false-negative. status.DiscoverWorkflowDir is the same
// predicate both discoverWorkflows (handlers.go) and `--boot --identify` use, so
// this is a same-mechanism CLI-level assertion, not a proxy.
func TestSharedScenarioFixturesAreDiscoverable(t *testing.T) {
	for name, readme := range fixtureDiscoveryReadmes() {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme()), 0o644); err != nil {
				t.Fatal(err)
			}
			if _, found := status.DiscoverWorkflowDir(root); !found {
				t.Fatalf("fixture %q's README is not discoverable (missing `commissioned-by: spacedock@` marker) — the FO's boot-path choice between --boot --identify and --workflow-dir . can false-negative this scenario", name)
			}
		})
	}
}

// TestZeroDiscoverFixtureStaysUndiscoverable guards Fix B's over-application: the
// zero-discovery scenario's SUBJECT is non-discovery (TestLiveZeroDiscoverReportsAndStops
// proves an FO report-and-stops on a genuinely empty root), so its root must stay
// undiscoverable — gaining a marker here would silently disable the very scenario
// it is meant to exercise.
func TestZeroDiscoverFixtureStaysUndiscoverable(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".gitkeep"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, found := status.DiscoverWorkflowDir(root); found {
		t.Fatal("zero-discovery fixture root is discoverable — it must stay markerless so TestLiveZeroDiscoverReportsAndStops still exercises a real zero-discover boot")
	}
}
