// ABOUTME: Shared Git repository placement for split-root workflow state.
// ABOUTME: Keeps status, dispatch, and state verbs anchored to the main worktree.
package status

import (
	"bytes"
	"fmt"
	"os"
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
	gitDirOut, err := runGitCmd(dir, "rev-parse", "--path-format=absolute", "--git-dir")
	if err != nil {
		// A path reached through an external symlink has no lexical .git ancestor,
		// so Git is the primary classifier. Only downgrade a failed Git probe to
		// standalone when neither the lexical nor canonical path has repository
		// metadata; a broken linked checkout must still fail closed.
		lexicalRoot := FindGitRoot(dir)
		canonicalRoot := FindGitRoot(RealpathOf(dir))
		if !hasGitEntry(lexicalRoot) && !hasGitEntry(canonicalRoot) {
			return RepositoryPlacement{}, nil
		}
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

	gitDir := RealpathOf(TrimGitLineTerminator(gitDirOut))
	commonDir := RealpathOf(TrimGitLineTerminator(commonDirOut))
	placement := RepositoryPlacement{
		MainWorktreeRoot: TrimGitLineTerminator(topOut),
		Prefix:           TrimGitLineTerminator(prefixOut),
		CommonGitDir:     commonDir,
		InGit:            true,
		Linked:           gitDir != commonDir,
	}
	if !placement.Linked {
		return placement, nil
	}

	worktreesOut, err := runGitCmd(dir, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		return RepositoryPlacement{}, fmt.Errorf("resolve main worktree for %s: %w", dir, err)
	}
	mainRoot, err := primaryWorktreeFromPorcelainZ([]byte(worktreesOut))
	if err != nil {
		return RepositoryPlacement{}, fmt.Errorf("resolve main worktree for %s: %w", dir, err)
	}
	if !filepath.IsAbs(mainRoot) {
		return RepositoryPlacement{}, fmt.Errorf("resolve main worktree for %s: primary worktree path is not absolute", dir)
	}
	info, err := os.Stat(mainRoot)
	if err != nil || !info.IsDir() || !hasGitEntry(mainRoot) {
		return RepositoryPlacement{}, fmt.Errorf("resolve main worktree for %s: primary worktree %q is not usable", dir, mainRoot)
	}
	if RealpathOf(mainRoot) == commonDir {
		return RepositoryPlacement{}, fmt.Errorf("resolve main worktree for %s: primary worktree resolves to the common Git directory", dir)
	}
	placement.MainWorktreeRoot = mainRoot
	return placement, nil
}

// TrimGitLineTerminator removes Git's one output-record LF without trimming
// path bytes. A path or prefix may legitimately begin or end with spaces, tabs,
// carriage returns, or even a newline; only the command's final LF is framing.
func TrimGitLineTerminator(out string) string {
	return strings.TrimSuffix(out, "\n")
}

// primaryWorktreeFromPorcelainZ parses one complete primary record from
// `git worktree list --porcelain -z`. Under -z, fields are NUL-delimited and
// records end with a second NUL; path bytes are never C-quoted, so embedded
// newlines and backslashes remain unambiguous. The first record is authoritative:
// a bare primary is rejected rather than skipped in favor of a linked worktree.
func primaryWorktreeFromPorcelainZ(raw []byte) (string, error) {
	records, err := ParseWorktreePorcelainZ(raw)
	if err != nil {
		return "", err
	}
	if len(records) == 0 {
		return "", fmt.Errorf("git worktree list returned no primary record")
	}
	primary := records[0]
	if primary.Bare {
		return "", fmt.Errorf("primary worktree record is bare")
	}
	if primary.Prunable {
		return "", fmt.Errorf("primary worktree record is prunable")
	}
	return primary.Path, nil
}

// WorktreeRecord is one complete `git worktree list --porcelain -z` record.
type WorktreeRecord struct {
	Path     string
	Branch   string
	Bare     bool
	Prunable bool
}

// ParseWorktreePorcelainZ decodes complete NUL-delimited porcelain records.
// It preserves path bytes exactly and rejects truncated or structurally
// ambiguous records instead of guessing at repository placement.
func ParseWorktreePorcelainZ(raw []byte) ([]WorktreeRecord, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	separator := []byte{0, 0}
	if !bytes.HasSuffix(raw, separator) {
		return nil, fmt.Errorf("git worktree list returned an incomplete record")
	}
	parts := bytes.Split(raw, separator)
	records := make([]WorktreeRecord, 0, len(parts)-1)
	for i, part := range parts[:len(parts)-1] {
		if len(part) == 0 {
			return nil, fmt.Errorf("git worktree list returned an empty record at index %d", i)
		}
		var record WorktreeRecord
		pathSeen := false
		for _, field := range bytes.Split(part, []byte{0}) {
			if len(field) == 0 {
				return nil, fmt.Errorf("worktree record %d has an empty field", i)
			}
			switch {
			case bytes.Equal(field, []byte("bare")):
				record.Bare = true
			case bytes.HasPrefix(field, []byte("prunable")):
				record.Prunable = true
			case bytes.HasPrefix(field, []byte("branch ")):
				record.Branch = string(bytes.TrimPrefix(field, []byte("branch ")))
			case bytes.HasPrefix(field, []byte("worktree ")):
				if pathSeen {
					return nil, fmt.Errorf("worktree record %d has multiple paths", i)
				}
				pathSeen = true
				record.Path = string(bytes.TrimPrefix(field, []byte("worktree ")))
			}
		}
		if !pathSeen || record.Path == "" {
			return nil, fmt.Errorf("worktree record %d has no path", i)
		}
		records = append(records, record)
	}
	return records, nil
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
