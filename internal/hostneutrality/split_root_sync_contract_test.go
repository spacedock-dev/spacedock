// ABOUTME: B5/B6 contract skill-text tests — assert the FO halt-gate + push/pull
// ABOUTME: sync + rebase-conflict halt prose is present at the FO and ensign homes.
package hostneutrality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// foCorePath is the FO shared-core contract.
var foCorePath = filepath.Join("..", "..", "skills", "first-officer", "references",
	"first-officer-shared-core.md")

// ensignCorePath is the ensign shared-core contract.
var ensignCorePath = filepath.Join("..", "..", "skills", "ensign", "references",
	"ensign-shared-core.md")

// commissionSkillPath is the commission SKILL.md.
var commissionSkillPath = filepath.Join("..", "..", "skills", "commission", "SKILL.md")

// readSkill reads a skill file, failing the test on error.
func readSkill(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

// assertAll fails with a per-token diagnostic for each token missing from text.
func assertAll(t *testing.T, name, text string, tokens []string) {
	t.Helper()
	for _, tok := range tokens {
		if !strings.Contains(text, tok) {
			t.Errorf("%s missing required contract token: %q", name, tok)
		}
	}
}

// TestFOHaltGateProse is a non-AC text-consistency lint: it asserts the FO core
// carries the boot halt-gate prose keyed on the boot fields (split-root &&
// entity_dir_present false → HALT, point at `spacedock state init`). Per the proof
// policy this presence check does NOT prove the FO halts; an inverted clause keeps
// every token. The MECHANISM is proven by command-level tests — the binary EMITS
// the halt signal (internal/status TestBootJSONStateBackendEntityDirAbsent observes
// `entity_dir_present: false` + `state_backend: split-root`) and the `spacedock
// state init` recovery WORKS (internal/cli TestStateInitResumesFreshClone) — but
// the OWED behavioral proof that the FO actually HALTs on that signal (rather than
// running state init silently and proceeding) is a live drive, tracked as task
// ev3e (fo-halt-sync-journey-live-drives). This lint guards the prose tokens until
// ev3e lands the live oracle.
func TestFOHaltGateProse(t *testing.T) {
	markNonAC(t, "OWED live drive: task ev3e (fo-halt-sync-journey-live-drives). Mechanism today: internal/status TestBootJSONStateBackendEntityDirAbsent (binary emits the halt signal) + internal/cli TestStateInitResumesFreshClone (recovery works)")
	text := readSkill(t, foCorePath)
	assertAll(t, "FO core (B5 halt-gate)", text, []string{
		"state_backend",
		"entity_dir_present",
		"split-root",
		"spacedock state init",
		"HALT",
	})
}

// TestFOSyncProse is a non-AC text-consistency lint: it asserts the FO core
// carries the B6 sync prose (pull --rebase, push origin, the M-3 rebase-conflict
// halt: abort + no force-push + no auto-resolve). Per the proof policy this
// presence check does NOT prove the FO performs the sync. The git MECHANISM is
// proven by real two-writer e2e — internal/cli TestTwoWriterSyncHappyPath
// (push → non-ff rejection → pull --rebase → re-push on real clones) and
// TestTwoWriterSameEntityConflictHalts (CONFLICT → rebase --abort, no force-push)
// — but the OWED proof that the FO/ensign actually RUNS this sync at the contract
// points is a live drive, tracked as task ev3e (fo-halt-sync-journey-live-drives).
// This lint guards the prose tokens until ev3e lands the live oracle.
func TestFOSyncProse(t *testing.T) {
	markNonAC(t, "OWED live drive: task ev3e (fo-halt-sync-journey-live-drives). Mechanism today: internal/cli TestTwoWriterSyncHappyPath + TestTwoWriterSameEntityConflictHalts (real 2-writer git push/pull-rebase/conflict-halt)")
	text := readSkill(t, foCorePath)
	assertAll(t, "FO core (B6 sync)", text, []string{
		"pull --rebase",
		"push origin",
		"rebase --abort",
		"--force",
		"auto-resolve",
		"must NOT",
	})
}

// TestEnsignSyncProse is a non-AC text-consistency lint: it asserts the ensign
// core carries the B6 sync prose (push origin, pull --rebase, the M-3
// rebase-conflict halt) alongside the path-scoped rule. Same disposition as the FO
// half: the git MECHANISM is proven by the real two-writer e2e in internal/cli,
// but the OWED proof that the ensign RUNS this sync after its state commits is a
// live drive, tracked as task ev3e (fo-halt-sync-journey-live-drives). This lint
// guards the tokens until ev3e lands the live oracle.
func TestEnsignSyncProse(t *testing.T) {
	markNonAC(t, "OWED live drive: task ev3e (fo-halt-sync-journey-live-drives). Mechanism today: internal/cli TestTwoWriterSyncHappyPath + TestTwoWriterSameEntityConflictHalts (real 2-writer git push/pull-rebase/conflict-halt)")
	text := readSkill(t, ensignCorePath)
	assertAll(t, "ensign core (B6 sync)", text, []string{
		"push origin",
		"pull --rebase",
		"rebase --abort",
		"--force",
		"auto-resolve",
	})
}

// TestCommissionJourneyProse is a non-AC text-consistency lint: it asserts the
// commission SKILL.md carries the journey-1 orphan-branch mechanics (clear
// inherited tree, linked worktree, state-init pointer) and the journey-2 $inline
// prose. Per the proof policy this presence check does NOT prove the mechanics
// work. The orphan-birth/resume MECHANISM is proven by command-level tests —
// internal/cli TestStateNewBirthsSplitRoot + TestCommissionOrphanBranchScaffolding
// (orphan branch created with a cleared tree as a linked worktree) and
// TestStateInitInlineNoOp ($inline branch) — but the OWED proof that the
// commission FLOW drives these journeys is a live drive, tracked as task ev3e
// (fo-halt-sync-journey-live-drives). This lint guards the prose tokens until ev3e
// lands the live oracle.
func TestCommissionJourneyProse(t *testing.T) {
	markNonAC(t, "OWED live drive: task ev3e (fo-halt-sync-journey-live-drives). Mechanism today: internal/cli TestStateNewBirthsSplitRoot + TestCommissionOrphanBranchScaffolding + TestStateInitInlineNoOp (real orphan-birth/resume + inline e2e)")
	text := readSkill(t, commissionSkillPath)
	assertAll(t, "commission SKILL.md (journeys)", text, []string{
		"checkout --orphan",
		"clear",
		"linked worktree",
		"spacedock state init",
		"$inline",
	})
}
