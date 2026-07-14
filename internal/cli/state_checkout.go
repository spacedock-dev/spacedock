// ABOUTME: Split-root checkout path resolution — anchors the checkout at the
// ABOUTME: repository's main worktree regardless of the invoking cwd.
package cli

import (
	"fmt"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// resolveSplitRootCheckout maps (workflowDir, relPath) to the split-root state
// checkout's absolute path. The checkout is a process-wide singleton that only
// ever exists at the repository's MAIN worktree — its directory is gitignored, so
// a linked (agent) worktree's copy of the same relative path never exists. When
// workflowDir sits inside a linked worktree, the checkout is re-anchored at
// <main-worktree-root>/<repo-relative-prefix>/<relPath>. Outside a linked
// worktree — the main worktree itself, or workflowDir not inside a git repository
// at all (the standalone test fixtures) — this is a plain
// filepath.Join(workflowDir, relPath).
func resolveSplitRootCheckout(workflowDir, relPath string) (string, error) {
	return status.ResolveSplitRootCheckout(workflowDir, relPath)
}

// validateExistingStateCheckout proves that an existing state directory is the
// exact registered linked-worktree root for branch in workflowDir's repository.
// Without this guard, `git -C statePath` can climb into the main worktree when a
// plain directory occupies statePath and mutate the code branch.
func validateExistingStateCheckout(workflowDir, statePath, branch string) error {
	placement, err := status.ResolveRepositoryPlacement(workflowDir)
	if err != nil {
		return err
	}
	if !placement.InGit {
		// Preserve standalone fixtures/workflows where the state checkout is its
		// own repository. Exact-root proof still prevents Git from climbing into
		// an unrelated parent, but there is no shared worktree registry or state
		// branch contract to enforce outside a repository-backed workflow.
		ok, topOut := runGit(statePath, "rev-parse", "--path-format=absolute", "--show-toplevel")
		if !ok {
			return fmt.Errorf("standalone state checkout %s is not a Git repository root:\n%s", statePath, topOut)
		}
		if top := status.RealpathOf(status.TrimGitLineTerminator(topOut)); top != status.RealpathOf(statePath) {
			return fmt.Errorf("standalone state checkout %s resolves through Git root %s", statePath, top)
		}
		return nil
	}
	ok, out := runGit(workflowDir, "worktree", "list", "--porcelain", "-z")
	if !ok {
		return fmt.Errorf("listing worktree registrations failed:\n%s", out)
	}
	records, err := status.ParseWorktreePorcelainZ([]byte(out))
	if err != nil {
		return fmt.Errorf("parsing worktree registrations: %w", err)
	}
	target := status.RealpathOf(statePath)
	wantBranch := "refs/heads/" + branch
	matches := 0
	for _, record := range records {
		if status.RealpathOf(record.Path) != target {
			continue
		}
		matches++
		if record.Bare || record.Prunable {
			return fmt.Errorf("state checkout %s has an unusable worktree registration", statePath)
		}
		if record.Branch != wantBranch {
			return fmt.Errorf("state checkout %s is on %q, expected %q", statePath, record.Branch, wantBranch)
		}
	}
	if matches != 1 {
		return fmt.Errorf("state checkout %s has %d exact worktree registrations, expected 1", statePath, matches)
	}

	ok, topOut := runGit(statePath, "rev-parse", "--path-format=absolute", "--show-toplevel")
	if !ok {
		return fmt.Errorf("state checkout %s is not a usable worktree root:\n%s", statePath, topOut)
	}
	if top := status.RealpathOf(status.TrimGitLineTerminator(topOut)); top != target {
		return fmt.Errorf("state checkout %s resolves through worktree root %s", statePath, top)
	}
	ok, commonOut := runGit(statePath, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if !ok {
		return fmt.Errorf("state checkout %s has unreadable common Git metadata:\n%s", statePath, commonOut)
	}
	if common := status.RealpathOf(status.TrimGitLineTerminator(commonOut)); common != status.RealpathOf(placement.CommonGitDir) {
		return fmt.Errorf("state checkout %s belongs to a different Git repository", statePath)
	}
	ok, gitDirOut := runGit(statePath, "rev-parse", "--path-format=absolute", "--git-dir")
	if !ok {
		return fmt.Errorf("state checkout %s has unreadable Git metadata:\n%s", statePath, gitDirOut)
	}
	if status.RealpathOf(status.TrimGitLineTerminator(gitDirOut)) == status.RealpathOf(placement.CommonGitDir) {
		return fmt.Errorf("state checkout %s is not a linked worktree", statePath)
	}
	return nil
}
