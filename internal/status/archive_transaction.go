// ABOUTME: Explicit split-root retirement owns the archive move and its durability.
package status

import (
	"fmt"
	"io"
	"path/filepath"
)

func runArchiveTransaction(roots roots, slug string, force, quiet, asJSON bool, stdout, stderr io.Writer) int {
	mode, _, err := ClassifyState(ParseFrontmatter(filepath.Join(roots.definitionDir, "README.md"))["state"])
	if err != nil {
		return errExit(stderr, err.Error())
	}
	if mode == StateInline {
		return runArchive(roots.definitionDir, roots.entityDir, roots.entityDirSpelling, slug, force, quiet, asJSON, stdout, stderr)
	}
	// Mutation preflight handles an existing rebase before checking HEAD's branch.
	if rc := preflightArchiveState(roots, slug, "archive", asJSON, stdout, stderr); rc != 0 {
		return rc
	}
	snapshot, err := captureArchiveState(roots.entityDir, slug)
	if err != nil {
		return errExit(stderr, fmt.Sprintf("archive: failed to snapshot %s: %v", slug, err))
	}
	if rc := runArchive(roots.definitionDir, roots.entityDir, roots.entityDirSpelling, slug, force, true, false, io.Discard, stderr); rc != 0 {
		return rc
	}
	if rc := commitArchiveMove(roots.entityDir, slug, snapshot, "retirement", stderr); rc != 0 {
		if err := rollbackArchive(roots.entityDir, slug, snapshot); err != nil {
			fmt.Fprintf(stderr, "archive: CRITICAL: rollback failed for %s: %v\n", slug, err)
		}
		return rc
	}
	durability, _, rc := publishArchive(roots, slug, "retirement", stdout, stderr)
	if asJSON {
		emitJSON(stdout, newJSONObj().set("command", "archive").set("slug", slug).set("result", durability))
	} else if !quiet {
		fmt.Fprintf(stdout, "archived: %s (retired). State durability: %s.\n", slug, durability)
	}
	return rc
}
