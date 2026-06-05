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
// presence check does NOT prove the FO performs the sync. The git MECHANICS are
// already oracle-covered by real two-writer e2e — internal/cli state_sync_test.go
// (TestTwoWriterSyncHappyPath: push → non-ff rejection → pull --rebase → re-push;
// TestTwoWriterSameEntityConflictHalts: CONFLICT → rebase --abort, no force-push)
// and internal/dispatch build_statecommit_test.go. The remaining behavioral
// proof — that the FO actually ISSUES this sync at the contract points — rides
// task ev3e's halt drive (fo-halt-sync-journey-live-drives), where ev3e's ideation
// folded the sync/journey residual into the halt scenario. This lint guards the
// prose tokens.
func TestFOSyncProse(t *testing.T) {
	markNonAC(t, "behavioral-issuance rides task ev3e's halt drive (fo-halt-sync-journey-live-drives). Sync MECHANICS already oracle-covered: internal/cli state_sync_test.go (TestTwoWriterSyncHappyPath + TestTwoWriterSameEntityConflictHalts) + internal/dispatch build_statecommit_test.go (TestStateCommitGuidanceResolvesPaths)")
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
// half: the git MECHANICS are oracle-covered by the real two-writer e2e in
// internal/cli + build_statecommit_test.go; the remaining behavioral proof that
// the ensign ISSUES this sync after its state commits rides task ev3e's halt drive
// (fo-halt-sync-journey-live-drives). This lint guards the tokens.
func TestEnsignSyncProse(t *testing.T) {
	markNonAC(t, "behavioral-issuance rides task ev3e's halt drive (fo-halt-sync-journey-live-drives). Sync MECHANICS already oracle-covered: internal/cli state_sync_test.go (TestTwoWriterSyncHappyPath + TestTwoWriterSameEntityConflictHalts) + internal/dispatch build_statecommit_test.go (TestStateCommitGuidanceResolvesPaths)")
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
// work. The orphan-birth/resume MECHANICS are oracle-covered by command-level
// tests — internal/cli state_new_test.go (TestStateNewBirthsSplitRoot) +
// state_init_test.go (TestCommissionOrphanBranchScaffolding: orphan branch with a
// cleared tree as a linked worktree; TestStateInitInlineNoOp: the $inline branch).
// The remaining behavioral proof that the commission FLOW drives these journeys
// rides task ev3e's halt drive (fo-halt-sync-journey-live-drives). This lint guards
// the prose tokens.
func TestCommissionJourneyProse(t *testing.T) {
	markNonAC(t, "behavioral-issuance rides task ev3e's halt drive (fo-halt-sync-journey-live-drives). Journey MECHANICS already oracle-covered: internal/cli state_init_test.go (TestStateInitResumesFreshClone + TestCommissionOrphanBranchScaffolding + TestStateInitInlineNoOp) + state_new_test.go (TestStateNewBirthsSplitRoot)")
	text := readSkill(t, commissionSkillPath)
	assertAll(t, "commission SKILL.md (journeys)", text, []string{
		"checkout --orphan",
		"clear",
		"linked worktree",
		"spacedock state init",
		"$inline",
	})
}
