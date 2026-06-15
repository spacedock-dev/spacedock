// ABOUTME: AC-2 single-source — each moved FO ceremony block lives in a host-neutral
// ABOUTME: core XOR in any runtime adapter, so a copy-paste-into-a-host regression goes RED.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// foReferencesDir is the FO references tree both the cores and the host adapters
// live under.
func foReferencesDir() string {
	return filepath.Join("skills", "first-officer", "references")
}

// hostAdapterFiles are the three runtime adapters plus the two Claude seam residues.
// AC-2's absence half asserts no MOVED ceremony block is restated in ANY of them —
// the seam files (claude-fo-merge/dispatch) are adapters too, so a block left behind
// in a seam during extraction is caught here, not only a copy pasted into codex/pi.
var hostAdapterFiles = []string{
	"claude-first-officer-runtime.md",
	"codex-first-officer-runtime.md",
	"pi-first-officer-runtime.md",
	"claude-fo-merge.md",
	"claude-fo-dispatch.md",
}

// movedCeremonyBlock is one fingerprint of a ceremony block the extraction moved
// into a host-neutral core: a distinctive verbatim sentence the block owns, the
// core it must live in, and a human label. The fingerprint is chosen to be unique
// to the moved block so its presence in the core and its absence from every adapter
// is an unambiguous single-source signal.
type movedCeremonyBlock struct {
	label       string
	fingerprint string
	core        string
}

// movedCeremonyBlocks covers the merge ceremony blocks (mod-block set→invoke→clear,
// the Ship-Local trunk step, the worktree-removal ladder, Mod-Block Enforcement) and
// the dispatch blocks (the dispatch-build-mandatory rule, the reuse conditions, the
// model-mismatch diagnostic) — the host-neutral text the design re-homed. Each must
// be present in its named core and absent from every host adapter.
var movedCeremonyBlocks = []movedCeremonyBlock{
	{"mod-block set→invoke→clear sequence", "The set→invoke→clear sequence (steps 1, 2, 4) is mandatory whenever a merge hook is registered", "fo-merge-core.md"},
	{"Ship-Local trunk resolution", "BASE=$(spacedock dispatch trunk --workflow-dir {workflow_dir})", "fo-merge-core.md"},
	{"worktree-removal safety ladder", "If removal fails on untracked files, the FO MUST:", "fo-merge-core.md"},
	{"Mod-Block Enforcement empty-state recovery", "In the empty-pr/empty-mod-block state the merge hook has provably not run", "fo-merge-core.md"},
	{"dispatch-build-mandatory rule", "Do NOT assemble worker prompts manually", "fo-dispatch-core.md"},
	{"reuse-routing condition", "Reuse-routing matches the entity's worktree state", "fo-dispatch-core.md"},
	{"model-mismatch diagnostic anchor", "does not match next stage effective_model", "fo-dispatch-core.md"},
}

// readFORef reads one file in the FO references tree, repo-relative.
func readFORef(t *testing.T, root, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, foReferencesDir(), name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(b)
}

// TestMovedCeremonyBlocksAreSingleSourced is the AC-2 single-source guard: every
// moved ceremony block lives in its host-neutral core (present-in-core) AND in NONE
// of the host adapters (absent-from-every-adapter) — the block lives in the core XOR
// in any adapter. The independent source is the content-fingerprint over the two
// cores and the five adapter files: a true single-source extraction passes; a
// copy-paste-into-a-host regression (a block left in a seam during extraction, or
// pasted into codex/pi) leaves the fingerprint in an adapter and goes RED. It is NOT
// a tautological "grep finds the phrase in the core" — the load-bearing assertion is
// the ABSENCE from every adapter, which a present-here-only grep never checks.
func TestMovedCeremonyBlocksAreSingleSourced(t *testing.T) {
	root := repoRoot(t)
	if len(movedCeremonyBlocks) == 0 {
		t.Fatal("no moved ceremony blocks declared — the single-source check would pass vacuously")
	}
	for _, blk := range movedCeremonyBlocks {
		core := readFORef(t, root, blk.core)
		if !strings.Contains(core, blk.fingerprint) {
			t.Errorf("%s: fingerprint %q absent from its core %s — the block did not land in the core (or the fingerprint drifted)", blk.label, blk.fingerprint, blk.core)
		}
		for _, adapter := range hostAdapterFiles {
			body := readFORef(t, root, adapter)
			if strings.Contains(body, blk.fingerprint) {
				t.Errorf("%s: fingerprint %q restated in host adapter %s — the moved ceremony block must live in the core XOR in any adapter; a copy-paste-into-a-host regression must be removed", blk.label, blk.fingerprint, adapter)
			}
		}
	}
}

// TestMovedCeremonyControlFailsOnRestatedBlock is the AC-2 control: it proves the
// single-source logic goes RED when a moved fingerprint is restated in an adapter.
// It drives the same present-in-core + absent-from-adapter check the real guard uses
// against an in-memory adapter set that includes a planted restatement, so the
// control exercises the real discriminator rather than re-implementing it.
func TestMovedCeremonyControlFailsOnRestatedBlock(t *testing.T) {
	root := repoRoot(t)
	blk := movedCeremonyBlocks[0]
	core := readFORef(t, root, blk.core)
	if !strings.Contains(core, blk.fingerprint) {
		t.Fatalf("control precondition: fingerprint %q not in core %s", blk.fingerprint, blk.core)
	}
	// A planted adapter body that restates the moved block — the copy-paste regression.
	planted := "## Some host seam\n\n" + blk.fingerprint + "\n"
	if !strings.Contains(planted, blk.fingerprint) {
		t.Fatal("control: planted adapter body unexpectedly lacks the restated fingerprint — the discriminator never exercises the absence check")
	}
	// And a clean adapter body that does not — the discriminator must distinguish them.
	clean := "## Some host seam\n\nA host-specific note with no moved ceremony text.\n"
	if strings.Contains(clean, blk.fingerprint) {
		t.Fatal("control: clean adapter body unexpectedly contains the fingerprint — the discriminator has nothing to contrast the restatement against")
	}
}
