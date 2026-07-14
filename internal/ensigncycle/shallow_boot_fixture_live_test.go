//go:build live

package ensigncycle

import (
	"os"
	"path/filepath"
	"testing"
)

// writeShallowBootWorkflow seeds the shallow-boot fixture under root and returns
// the entity/archive paths plus the stub-gh dir. It registers the canonical
// pr-merge mod verbatim (so the boot JSON reports the startup/idle/merge hooks and
// an accidental engage can take the real advancement path) and seeds the
// gate-check + merged-pr entities. Live-tagged because it copies the canonical mod via repoRoot,
// which is a live-only helper.
func writeShallowBootWorkflow(t *testing.T, root string) shallowBootFixture {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), shallowBootReadme())
	modsDir := filepath.Join(root, "_mods")
	if err := os.MkdirAll(modsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	prMergeSrc := filepath.Join(repoRoot(t), "mods", "pr-merge.md")
	prMergeBody, err := os.ReadFile(prMergeSrc)
	if err != nil {
		t.Fatalf("read canonical pr-merge mod %s: %v", prMergeSrc, err)
	}
	writeFile(t, filepath.Join(modsDir, "pr-merge.md"), string(prMergeBody))

	gatePath := filepath.Join(root, "gate-check.md")
	writeFile(t, gatePath, shallowBootGateEntity())
	mergedPath := filepath.Join(root, "merged-pr.md")
	writeFile(t, mergedPath, shallowBootMergedEntity())
	gitInit(t, root)

	return shallowBootFixture{
		gateEntityPath:   gatePath,
		mergedEntityPath: mergedPath,
		mergedArchive:    filepath.Join(root, "_archive", "merged-pr.md"),
		stubGhDir:        writeStubMergedGh(t),
	}
}

func writeMergeModRecoveryWorkflow(t *testing.T, root string) mergeModRecoveryFixture {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), shallowBootReadme())
	modsDir := filepath.Join(root, "_mods")
	if err := os.MkdirAll(modsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	prMergeSrc := filepath.Join(repoRoot(t), "mods", "pr-merge.md")
	prMergeBody, err := os.ReadFile(prMergeSrc)
	if err != nil {
		t.Fatalf("read canonical pr-merge mod %s: %v", prMergeSrc, err)
	}
	writeFile(t, filepath.Join(modsDir, "pr-merge.md"), string(prMergeBody))
	entityPath := filepath.Join(root, "merge-recovery.md")
	writeFile(t, entityPath, mergeModRecoveryEntity())
	gitInit(t, root)
	return mergeModRecoveryFixture{
		entityPath:  entityPath,
		archivePath: filepath.Join(root, "_archive", "merge-recovery.md"),
	}
}

// writeStubMergedGh writes a `gh` shim that reports MERGED for `gh pr view`, so the
// engage's live PR-state probe and the pr-merge startup hook see a merged PR
// deterministically (offline, no real PR). Read-only boot must not invoke it.
// Returns the dir to prepend to PATH.
func writeStubMergedGh(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	// gh pr view {n} --json state --jq .state -> MERGED; any other gh subcommand
	// (e.g. repo view) prints an empty line so it does not hard-error the FO.
	script := "#!/bin/sh\n" +
		"case \"$1 $2\" in\n" +
		"  \"pr view\") echo MERGED ;;\n" +
		"  *) echo \"\" ;;\n" +
		"esac\n"
	path := filepath.Join(dir, "gh")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}
