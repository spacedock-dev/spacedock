// ABOUTME: Offline guard that the live cycle's workflow README fixture is
// ABOUTME: discoverable, so the FO's `status --discover` finds it (no model).
package ensigncycle

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// TestLiveFixtureIsDiscoverable guards the SECOND pre-TeamCreate blocker the audit
// found on CI run 26839572693's opus leg: even with `spacedock` on PATH, the FO
// runs `status --discover` (FO contract step 3) and gets ZERO workflows when the
// staged README lacks `commissioned-by: spacedock@`, then reports "no workflow
// found" and exits before TeamCreate. Discovery gates on that frontmatter field —
// both discoverWorkflows (internal/status/handlers.go) and DiscoverWorkflowDir
// (internal/status/discover_walkup.go) require strings.HasPrefix(commissioned-by,
// "spacedock@"). This test writes the SAME README the live cycle stages
// (readmeNonWorktree) and asserts DiscoverWorkflowDir — which uses the same
// predicate as the FO's --discover — resolves the fixture dir. Red without the
// commissioned-by line in the fixture, green with it. No model.
func TestLiveFixtureIsDiscoverable(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readmeNonWorktree()), 0o644); err != nil {
		t.Fatal(err)
	}

	dir, found := status.DiscoverWorkflowDir(root)
	if !found {
		t.Fatal("live fixture README is not discoverable: DiscoverWorkflowDir found no workflow " +
			"(the fixture needs a `commissioned-by: spacedock@` frontmatter line, or the FO's " +
			"`status --discover` reports 0 workflows and exits before TeamCreate)")
	}
	wantReal, err := filepath.EvalSymlinks(root)
	if err != nil {
		wantReal = root
	}
	gotReal, err := filepath.EvalSymlinks(dir)
	if err != nil {
		gotReal = dir
	}
	if gotReal != wantReal {
		t.Errorf("DiscoverWorkflowDir = %q, want the fixture root %q", dir, root)
	}
}
