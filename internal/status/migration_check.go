// ABOUTME: Migration-check walk-step — the non-entity-tree prune composition
// ABOUTME: shared by the live consistency check and its hermetic prune test.
package status

import (
	"os"
	"path/filepath"
)

// isMigrationCheckPrunedDir reports whether the migration-check walk skips a
// directory whole, given its path. One non-entity tree is pruned: the
// .spacedock-state split-root state checkout. It is gitignored on the code
// branch; on a dev machine it materializes under docs/dev/ and holds the
// machine-local entity files plus non-entity records — _debriefs session notes
// and the _roadmap strategy layer (sprint definitions, staff reviews, debriefs).
// On CI the tree is absent and the prune is moot.
//
// Those records carry their own frontmatter shape (session-date, sequence,
// first-commit), whose bare-YAML date scalars decode as time.Time directly but
// as strings through the reader — expected for non-entity frontmatter and
// outside this check's scope. The tree is pruned wholesale: the migration check
// governs fixture frontmatter on the code branch, not the state checkout.
func isMigrationCheckPrunedDir(path string) bool {
	return filepath.Base(path) == ".spacedock-state"
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
