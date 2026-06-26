// ABOUTME: Contract anchor — the FO records its launch branch and HALTS if the
// ABOUTME: shared working tree's branch changes underneath it, instead of
// ABOUTME: dispatching/committing into a tree a concurrent actor switched (RC2b).
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFOPinsLaunchBranch locks RC2b: at boot the FO records the working tree's
// branch and re-checks it before dispatch / state-changing git ops, halting on a
// change. Pinned on the branch NAME so a same-branch fast-forward is not a halt.
// Without this, a concurrent branch switch (the DRC-3653 incident) silently moves
// the FO onto the wrong branch and deletes its tracked READMEs.
func TestFOPinsLaunchBranch(t *testing.T) {
	path := filepath.Join(repoRoot(t), "skills", "first-officer", "references", "first-officer-shared-core.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}
	content := string(data)

	for _, r := range []string{
		"launch branch",                // the recorded baseline
		"before each dispatch",         // when it is re-checked
		"branch changed underneath me", // the halt surface
	} {
		if !strings.Contains(content, r) {
			t.Errorf("first-officer-shared-core.md no longer pins the launch branch (RC2b): missing %q.\n"+
				"Without it, a concurrent branch switch moves the FO onto the wrong branch and deletes its tracked workflow files.", r)
		}
	}

	// Must pin on the branch NAME, not the commit — a same-branch fast-forward is
	// normal and must not trigger a false halt.
	if !strings.Contains(content, "branch NAME") {
		t.Errorf("the branch-pin must be on the branch NAME (a same-branch fast-forward is not a halt) — missing that qualification.")
	}
}
