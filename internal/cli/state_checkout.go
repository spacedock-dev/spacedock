// ABOUTME: Split-root checkout path resolution — anchors the checkout at the
// ABOUTME: repository's main worktree regardless of the invoking cwd.
package cli

import "github.com/spacedock-dev/spacedock/internal/status"

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
