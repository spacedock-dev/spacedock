// ABOUTME: Migration-check walk-step — the non-entity-tree prune composition
// ABOUTME: shared by the live consistency check and its hermetic prune test.
package status

import (
	"os"
	"path/filepath"
)

// isMigrationCheckPrunedDir reports whether the migration-check walk skips a
// directory whole, given its path. Two non-entity trees are pruned, both holding
// frontmatter the migration check was never meant to govern:
//
//   - .spacedock-state — the gitignored split-root state checkout. On a dev
//     machine it materializes under docs/dev/ and holds ~100 machine-local entity
//     files (active entities plus _debriefs session records); on CI the tree is
//     absent and the prune is moot.
//   - docs/roadmap — the strategy layer (sprint definitions, staff reviews,
//     session debriefs). It owns outcome/scope/sequencing, NOT entity state
//     (that lives in docs/dev/.spacedock-state), so it carries no entity
//     frontmatter; matched as a `roadmap` dir directly under a `docs` dir so the
//     prune is location-independent (fires under the repo root AND under a temp
//     fixture root).
//
// The debriefs in both trees carry their own frontmatter shape (session-date,
// sequence, first-commit), whose bare-YAML date scalars decode as time.Time
// directly but as strings through the reader — expected for non-entity
// frontmatter and outside this check's scope. Pruning these trees wholesale
// matches the established sibling prunes in handlers.go (discoverIgnoreDirs) and
// boundary_guard_test.go.
func isMigrationCheckPrunedDir(path string) bool {
	name := filepath.Base(path)
	if name == ".spacedock-state" {
		return true
	}
	return name == "roadmap" && filepath.Base(filepath.Dir(path)) == "docs"
}

// migrationCheckWalkDir is the directory walk-step for the migration-check
// filepath.Walk. For a directory entry it returns filepath.SkipDir to prune a
// non-entity subtree wholesale (else nil to descend), and dir=true so the
// callback returns that action immediately. For a non-directory entry it returns
// dir=false and the callback proceeds to file handling. Sharing this step keeps
// the prune composition defined once — the live consistency check and the
// hermetic prune test drive the same walk-step, so a divergence in production's
// composition cannot pass the hermetic test unnoticed.
func migrationCheckWalkDir(path string, info os.FileInfo) (action error, dir bool) {
	if !info.IsDir() {
		return nil, false
	}
	if isMigrationCheckPrunedDir(path) {
		return filepath.SkipDir, true
	}
	return nil, true
}
