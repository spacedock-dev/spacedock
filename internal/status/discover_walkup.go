// ABOUTME: Walk-up workflow discovery and state-checkout detection — find the
// ABOUTME: enclosing commissioned workflow, or name a misdirected state checkout.
package status

import (
	"io"
	"path/filepath"
	"strings"
)

// DiscoverWorkflowDir walks up from startDir to the nearest ancestor whose
// README.md frontmatter declares a `commissioned-by: spacedock@` field, the
// same predicate discoverWorkflows uses to recognize a workflow. The first
// match on the way up wins — when workflows are nested, this resolves to the
// innermost enclosing workflow. Returns (dir, true) on a match, ("", false) at
// the filesystem root with no match.
func DiscoverWorkflowDir(startDir string) (string, bool) {
	d, err := filepath.Abs(startDir)
	if err != nil {
		d = startDir
	}
	for {
		readme := filepath.Join(d, "README.md")
		if isRegularFile(readme) {
			fields := ParseFrontmatter(readme)
			if strings.HasPrefix(fields["commissioned-by"], "spacedock@") {
				return d, true
			}
		}
		parent := filepath.Dir(d)
		if parent == d {
			return "", false
		}
		d = parent
	}
}

// ResolveWorkflowDir resolves the workflow directory when no --workflow-dir (or
// PIPELINE_DIR/--root) override applies: walk up from dir to the enclosing
// commissioned workflow via DiscoverWorkflowDir, falling back to a downward scan
// from the git toplevel when the walk-up misses. Exactly one downward match
// resolves; zero or multiple matches write a named diagnostic to stderr and
// return a non-zero code. Shared by status dispatch, merge guard, and the state
// verbs so every no-flag invocation agrees on discovery.
func ResolveWorkflowDir(dir string, stderr io.Writer) (string, int) {
	if discovered, ok := DiscoverWorkflowDir(dir); ok {
		return discovered, 0
	}
	return discoverWorkflowDownward(dir, stderr)
}

// stateCheckoutParent reports whether pointedDir is the state checkout of an
// enclosing workflow. It walks up from pointedDir; at each ancestor A it reads
// A/README.md's `state:` field, and if that field's cleaned, non-escaping value
// resolves to the same realpath as pointedDir, A is the definition dir. This
// reuses the same `state:` validation rules as resolveRoots (reject absolute /
// `..`-escaping). Returns (defDir, true) on a match, ("", false) otherwise.
func stateCheckoutParent(pointedDir string) (string, bool) {
	target := realpathOf(pointedDir)
	d, err := filepath.Abs(pointedDir)
	if err != nil {
		d = pointedDir
	}
	for {
		readme := filepath.Join(d, "README.md")
		if isRegularFile(readme) {
			if mode, relPath, err := ClassifyState(ParseFrontmatter(readme)["state"]); err == nil && mode == StateSplitRoot {
				if realpathOf(filepath.Join(d, relPath)) == target {
					return d, true
				}
			}
		}
		parent := filepath.Dir(d)
		if parent == d {
			return "", false
		}
		d = parent
	}
}
