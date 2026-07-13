// ABOUTME: Shared Git repository placement for split-root workflow state.
// ABOUTME: Keeps status, dispatch, and state verbs anchored to the main worktree.
package status

import (
	"fmt"
	"path/filepath"
	"strings"
)

// RepositoryPlacement describes where repository-wide files belong when dir is
// inside either the main worktree or a linked worktree. Prefix is dir's
// repository-relative path and CommonGitDir is the shared Git administration
// directory. InGit is false only when no enclosing .git entry exists.
type RepositoryPlacement struct {
	MainWorktreeRoot string
	Prefix           string
	CommonGitDir     string
	InGit            bool
	Linked           bool
}

// ResolveRepositoryPlacement resolves Git's authoritative worktree metadata.
// A conclusively non-Git directory is not an error. Once an enclosing .git entry
// exists, however, every metadata failure is returned: silently treating a
// broken linked worktree as a standalone directory could redirect mutations into
// a worktree-local copy.
func ResolveRepositoryPlacement(dir string) (RepositoryPlacement, error) {
	gitRoot := FindGitRoot(dir)
	if !hasGitEntry(gitRoot) {
		return RepositoryPlacement{}, nil
	}

	gitDirOut, err := runGitCmd(dir, "rev-parse", "--path-format=absolute", "--git-dir")
	if err != nil {
		return RepositoryPlacement{}, fmt.Errorf("resolve git directory for %s: %w", dir, err)
	}
	commonDirOut, err := runGitCmd(dir, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return RepositoryPlacement{}, fmt.Errorf("resolve git common directory for %s: %w", dir, err)
	}
	prefixOut, err := runGitCmd(dir, "rev-parse", "--show-prefix")
	if err != nil {
		return RepositoryPlacement{}, fmt.Errorf("resolve repository prefix for %s: %w", dir, err)
	}
	topOut, err := runGitCmd(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return RepositoryPlacement{}, fmt.Errorf("resolve worktree root for %s: %w", dir, err)
	}

	gitDir := RealpathOf(strings.TrimSpace(gitDirOut))
	commonDir := RealpathOf(strings.TrimSpace(commonDirOut))
	placement := RepositoryPlacement{
		MainWorktreeRoot: strings.TrimSpace(topOut),
		Prefix:           strings.TrimSpace(prefixOut),
		CommonGitDir:     commonDir,
		InGit:            true,
		Linked:           gitDir != commonDir,
	}
	if !placement.Linked {
		return placement, nil
	}

	worktreesOut, err := runGitCmd(dir, "worktree", "list", "--porcelain")
	if err != nil {
		return RepositoryPlacement{}, fmt.Errorf("resolve main worktree for %s: %w", dir, err)
	}
	for _, line := range strings.Split(worktreesOut, "\n") {
		if strings.HasPrefix(line, "worktree ") {
			placement.MainWorktreeRoot = strings.TrimSpace(strings.TrimPrefix(line, "worktree "))
			if placement.MainWorktreeRoot == "" {
				break
			}
			return placement, nil
		}
	}
	return RepositoryPlacement{}, fmt.Errorf("resolve main worktree for %s: git worktree list returned no worktree entries", dir)
}

// ResolveSplitRootCheckout anchors a workflow's state checkout under the main
// worktree. Standalone non-Git fixtures retain the historical plain join.
func ResolveSplitRootCheckout(workflowDir, relPath string) (string, error) {
	placement, err := ResolveRepositoryPlacement(workflowDir)
	if err != nil {
		return "", err
	}
	if !placement.InGit || !placement.Linked {
		return filepath.Join(workflowDir, relPath), nil
	}
	return filepath.Join(placement.MainWorktreeRoot, placement.Prefix, relPath), nil
}
