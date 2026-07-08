// ABOUTME: Split-root checkout path resolution — anchors the checkout at the
// ABOUTME: repository's main worktree regardless of the invoking cwd.
package cli

import (
	"fmt"
	"path/filepath"
	"strings"

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
func resolveSplitRootCheckout(workflowDir, relPath string) string {
	if !isLinkedWorktree(workflowDir) {
		return filepath.Join(workflowDir, relPath)
	}
	mainRoot, err := mainWorktreeRoot(workflowDir)
	if err != nil {
		return filepath.Join(workflowDir, relPath)
	}
	prefix, err := repoRelPrefix(workflowDir)
	if err != nil {
		return filepath.Join(workflowDir, relPath)
	}
	return filepath.Join(mainRoot, prefix, relPath)
}

// isLinkedWorktree reports whether dir sits inside a linked (non-main) git
// worktree. `--git-dir` and `--git-common-dir` agree only in the main worktree
// (both resolve to the repo's single `.git`); a linked worktree's git-dir is
// `.git/worktrees/<name>` while its git-common-dir stays the shared `.git`. Both
// are read with `--path-format=absolute` and realpath-normalized (macOS `/var` ->
// `/private/var`) so the comparison never depends on cwd-relative spelling. A
// non-git dir (the standalone test fixtures) reports false — the caller's
// fallback is the historical unanchored join.
func isLinkedWorktree(dir string) bool {
	gitDirOK, gitDirOut := runGit(dir, "rev-parse", "--path-format=absolute", "--git-dir")
	commonDirOK, commonDirOut := runGit(dir, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if !gitDirOK || !commonDirOK {
		return false
	}
	return status.RealpathOf(strings.TrimSpace(gitDirOut)) != status.RealpathOf(strings.TrimSpace(commonDirOut))
}

// mainWorktreeRoot returns the absolute path of the repository's main worktree,
// as seen from any linked worktree — `git worktree list --porcelain`'s first
// entry is always the main worktree, regardless of which worktree invokes it.
func mainWorktreeRoot(dir string) (string, error) {
	ok, out := runGit(dir, "worktree", "list", "--porcelain")
	if !ok {
		return "", fmt.Errorf("git worktree list --porcelain failed:\n%s", out)
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "worktree ") {
			return strings.TrimSpace(line[len("worktree "):]), nil
		}
	}
	return "", fmt.Errorf("git worktree list --porcelain returned no worktree entries")
}
