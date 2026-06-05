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
// every token. The behavioral halves are proven by real command-level tests: the
// binary EMITS the halt signal (internal/status TestBootJSONStateBackendEntityDirAbsent
// observes `entity_dir_present: false` + `state_backend: split-root` for an
// uninitialized split-root) and the recovery the FO points at WORKS
// (internal/cli TestStateInitResumesFreshClone drives `spacedock state init` and
// observes the state checkout materialize). This lint guards the prose tokens.
func TestFOHaltGateProse(t *testing.T) {
	markNonAC(t, "internal/status TestBootJSONStateBackendEntityDirAbsent (binary emits the halt signal) + internal/cli TestStateInitResumesFreshClone (the state-init recovery works)")
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
// presence check does NOT prove the FO performs the sync. The behavior is proven
// by real two-writer git e2e tests: internal/cli TestTwoWriterSyncHappyPath
// (push → non-ff rejection → pull --rebase → re-push, observed on real clones) and
// TestTwoWriterSameEntityConflictHalts (CONFLICT detected, rebase --abort, no
// force-push). This lint guards the prose tokens.
func TestFOSyncProse(t *testing.T) {
	markNonAC(t, "internal/cli TestTwoWriterSyncHappyPath + TestTwoWriterSameEntityConflictHalts (real 2-writer git push/pull-rebase/conflict-halt e2e)")
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
// rebase-conflict halt) alongside the path-scoped rule. Same behavioral oracle as
// the FO half: the real two-writer git e2e in internal/cli proves the push/pull-
// rebase/conflict-halt behavior the prose describes; this lint guards the tokens.
func TestEnsignSyncProse(t *testing.T) {
	markNonAC(t, "internal/cli TestTwoWriterSyncHappyPath + TestTwoWriterSameEntityConflictHalts (real 2-writer git push/pull-rebase/conflict-halt e2e)")
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
// work. The behavior is proven by real command-level tests that drive the orphan
// birth/resume: internal/cli TestStateNewBirthsSplitRoot + TestCommissionOrphanBranchScaffolding
// (the orphan branch is created with a cleared tree as a linked worktree) and
// TestStateInitInlineNoOp (the $inline branch). This lint guards the prose tokens.
func TestCommissionJourneyProse(t *testing.T) {
	markNonAC(t, "internal/cli TestStateNewBirthsSplitRoot + TestCommissionOrphanBranchScaffolding + TestStateInitInlineNoOp (real orphan-birth/resume + inline command e2e)")
	text := readSkill(t, commissionSkillPath)
	assertAll(t, "commission SKILL.md (journeys)", text, []string{
		"checkout --orphan",
		"clear",
		"linked worktree",
		"spacedock state init",
		"$inline",
	})
}
