// ABOUTME: Migration-check walk-step — the .spacedock-state prune composition
// ABOUTME: shared by the live consistency check and its hermetic prune test.
package status

import (
	"os"
	"path/filepath"
)

// isMigrationCheckPrunedDir reports whether the migration-check walk skips a
// directory whole. The .spacedock-state tree is the gitignored split-root state
// checkout: on a dev machine it materializes under docs/dev/ and holds ~100
// machine-local entity files (active entities plus _debriefs session records)
// the migration check was never meant to govern; on CI the tree is absent and
// the prune is moot. The debriefs in particular carry their own frontmatter
// shape (session-date, sequence, first-commit), whose bare-YAML date scalars
// decode as time.Time directly but as strings through the reader — expected for
// non-entity frontmatter and outside this check's scope. Pruning the tree
// wholesale matches the established sibling prunes in handlers.go
// (discoverIgnoreDirs) and boundary_guard_test.go.
func isMigrationCheckPrunedDir(name string) bool {
	return name == ".spacedock-state"
}

// migrationCheckWalkDir is the directory walk-step for the migration-check
// filepath.Walk. For a directory entry it returns filepath.SkipDir to prune the
// .spacedock-state subtree wholesale (else nil to descend), and dir=true so the
// callback returns that action immediately. For a non-directory entry it returns
// dir=false and the callback proceeds to file handling. Sharing this step keeps
// the prune composition defined once — the live consistency check and the
// hermetic prune test drive the same walk-step, so a divergence in production's
// composition cannot pass the hermetic test unnoticed.
func migrationCheckWalkDir(info os.FileInfo) (action error, dir bool) {
	if !info.IsDir() {
		return nil, false
	}
	if isMigrationCheckPrunedDir(info.Name()) {
		return filepath.SkipDir, true
	}
	return nil, true
}
