// ABOUTME: Contract anchor — an ensign commits only to its isolated target (a
// ABOUTME: worktree or a split-root state checkout), NEVER at the bare repo root.
// ABOUTME: A bare-root commit on a shared working tree lands the entity on whatever
// ABOUTME: branch a concurrent actor switched it to (the DRC-3653 collision). RC1.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEnsignNeverCommitsAtBareRepoRoot locks RC1: the ensign contract forbids
// `git add`/`git commit` at the bare repo root, and a single-root non-worktree
// stage (no worktree, no `state:` checkout) has no commit target — it writes the
// entity in place and signals. Without this, "MUST commit before signaling" drove
// a review ensign to commit at the shared repo root, which a concurrent branch
// switch had moved to an unrelated branch.
func TestEnsignNeverCommitsAtBareRepoRoot(t *testing.T) {
	path := filepath.Join(repoRoot(t), "skills", "ensign", "references", "ensign-shared-core.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read ensign shared core: %v", err)
	}
	content := string(data)

	for _, r := range []string{
		"bare repo root",                // the forbidden target, named
		"Single-Root, No Commit Target", // the carve-out subsection for non-worktree single-root stages
		"no ensign commit target",       // the explicit "nothing to commit" statement
	} {
		if !strings.Contains(content, r) {
			t.Errorf("ensign-shared-core.md no longer forbids the bare-root commit (RC1): missing %q.\n"+
				"A non-worktree single-root ensign must NOT git-commit at the repo root — it pollutes whatever branch the shared tree is on.", r)
		}
	}

	// The "MUST commit before signaling" rule must be qualified — an unqualified
	// "MUST commit" is what pushed the ensign to commit with no valid target.
	if strings.Contains(content, "MUST commit before signaling completion.**") {
		t.Errorf("the commit rule is unqualified again ('MUST commit before signaling completion.') — it must be scoped to a worktree/state-checkout target, never the bare repo root.")
	}
}
