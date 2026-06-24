// ABOUTME: Fleet-mode smoke — the binary precondition the FO's `## Fleet Mode`
// ABOUTME: adopt-all relies on: `status --discover` lists every commissioned member
// ABOUTME: workflow and each member's `status` runs independently under one root.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// fleetMemberReadme is a single-root commissioned workflow README: it carries the
// `commissioned-by` marker `--discover` keys on AND a minimal valid stage set so
// `status` lists its entities. Single-root keeps the smoke test free of a state
// checkout — fleet mode's per-member split-root handling is exercised by the live
// cycle, not this binary-surface precondition check.
const fleetMemberReadme = `---
commissioned-by: spacedock@1.0
entity-type: task
entity-label: task
entity-label-plural: tasks
id-style: slug
stages:
  defaults:
    worktree: false
    concurrency: 1
  states:
    - name: backlog
      initial: true
    - name: done
      terminal: true
---

# Member Workflow

### backlog

Start.

- **Outputs:** seed.

### done

Terminal.
`

// stageFleet builds a fleet root holding TWO commissioned member workflows under
// one git repo, each with one flat backlog entity. Returns the root and the two
// member definition dirs. This is the multi-workflow topology the FO adopts as a
// member set when launched with a fleet directive (shared core `## Fleet Mode`).
func stageFleet(t *testing.T) (root, alphaDir, betaDir string) {
	t.Helper()
	root = t.TempDir()
	mk := func(slug, title string) string {
		dir := filepath.Join(root, slug)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(fleetMemberReadme), 0o644); err != nil {
			t.Fatal(err)
		}
		entity := "---\nid: \"\"\ntitle: " + title + "\nstatus: backlog\nscore: \"0.50\"\nsource: smoke\n---\n# " + title + "\n\nSeed entity.\n"
		if err := os.WriteFile(filepath.Join(dir, slug+"-1.md"), []byte(entity), 0o644); err != nil {
			t.Fatal(err)
		}
		return dir
	}
	alphaDir = mk("wf-alpha", "Alpha one")
	betaDir = mk("wf-beta", "Beta one")
	gitInitFixture(t, root)
	return root, alphaDir, betaDir
}

// TestFleetDiscoverListsAllMembers locks the fleet-mode adopt-all precondition:
// `status --discover` over a multi-workflow root returns BOTH commissioned member
// dirs (what the FO adopts as the member set instead of presenting the list), and
// each member's `status` lists its own entity independently under the one root —
// the per-member operation the round-robin event loop drives.
func TestFleetDiscoverListsAllMembers(t *testing.T) {
	root, alphaDir, betaDir := stageFleet(t)

	// --discover must enumerate BOTH members (adopt-all consumes this set).
	cmd := exec.Command(spacedockBinary(t), "status", "--discover", "--root", root)
	cmd.Env = append(os.Environ(), "HOME="+t.TempDir())
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("status --discover failed: %v\n%s", err, out)
	}
	discover := string(out)
	for _, member := range []string{"wf-alpha", "wf-beta"} {
		if !strings.Contains(discover, member) {
			t.Fatalf("--discover did not enumerate member %q (fleet adopt-all would miss it):\n%s", member, discover)
		}
	}

	// Each member's status runs independently, scoped by its own --workflow-dir —
	// the per-member iteration the fleet round-robin loop performs.
	alpha, code := runStatus(t, alphaDir)
	if code != 0 || !strings.Contains(alpha, "Alpha one") {
		t.Fatalf("member wf-alpha status (exit %d) missing its entity:\n%s", code, alpha)
	}
	beta, code := runStatus(t, betaDir)
	if code != 0 || !strings.Contains(beta, "Beta one") {
		t.Fatalf("member wf-beta status (exit %d) missing its entity:\n%s", code, beta)
	}

	// Per-member independence: alpha's view must not leak beta's entity.
	if strings.Contains(alpha, "Beta one") || strings.Contains(beta, "Alpha one") {
		t.Fatalf("member views are not independent — cross-member entity leak:\nalpha:\n%s\nbeta:\n%s", alpha, beta)
	}
}
