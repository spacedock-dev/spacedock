// ABOUTME: AC-5 structural guard — the legacy TeamCreate machinery is layered
// ABOUTME: into a conditionally-loaded skill, with an externally-anchored removal
// ABOUTME: trigger; the normal-path dispatch contract inlines none of it.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// legacyMachineryStrings are the legacy TeamCreate/TeamDelete machinery markers
// that AC-5 requires to live ONLY in the conditionally-loaded
// using-legacy-claude-team skill, never inline in the normal-path dispatch
// contract (claude-fo-dispatch.md) every current-host session reads — the
// recovery ladder, the bounded terminal-teardown apparatus, the registry-desync
// issue, and the NAME_PATTERN naming constraint. A current-host session takes the
// no-match probe branch and never loads the legacy skill, so a re-inlined marker
// would be real context weight every such session pays.
var legacyMachineryStrings = []string{
	"TERMINAL_TEARDOWN_BOUNDED",
	"36806",
	"priority-ordered ladder",
	"fresh-suffixed",
	"TeamCreate recovery procedure",
	"NAME_PATTERN",
}

// deletedFromMergedFloor are strings the rewrite REMOVES outright (not relocated):
// the back-channel works without the `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` opt-in
// flag (OQ-2), so any captain-facing "set the flag to enable team mode" hint is
// stale. The whole env-var token must be absent from BOTH the normal-path dispatch
// contract and the legacy skill — re-introducing it anywhere resurrects the
// obsolete advice. The CI lane in runtime-live-e2e.yml legitimately sets the flag
// (test config, not captain-facing skill prose), so it is out of this check's
// scope (skills/ only).
var deletedFromMergedFloor = []string{
	"CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS",
}

// normalPathDispatchContractPath is claude-fo-dispatch.md — the Claude FO dispatch
// contract a current-host session reads at first dispatch. It carries the
// worker-back-channel shape and the idle/degraded discipline inline; it must
// inline none of the legacy TeamCreate machinery.
func normalPathDispatchContractPath(t *testing.T) string {
	return filepath.Join(skillsRoot(t), "first-officer", "references", "claude-fo-dispatch.md")
}

// legacyClaudeTeamSkillPath is the conditionally-loaded legacy skill that owns the
// full TeamCreate lifecycle, read only on a probe match.
func legacyClaudeTeamSkillPath(t *testing.T) string {
	return filepath.Join(skillsRoot(t), "using-legacy-claude-team", "SKILL.md")
}

func readFileString(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

// TestNormalPathContractInlinesNoLegacyMachinery is AC-5's layering invariant
// (structural-ABSENCE): the normal-path dispatch contract claude-fo-dispatch.md
// inlines NONE of the legacy TeamCreate machinery markers, and neither contract
// carries the stale enable-the-flag hint. A current-host session that takes the
// no-match probe branch never loads the legacy skill, so a re-inlined marker is
// real context weight every such session pays. This asserts ABSENCE from the
// normal-path contract only — it does not grep the legacy skill for its own prose
// (that would be a tautological prose-grep); the legacy machinery's continued
// existence is bound by TestLegacyConsumerRetiredButPathLives (the skill file and
// its dispatch-contract load line still resolve), not by a string match here.
func TestNormalPathContractInlinesNoLegacyMachinery(t *testing.T) {
	contract := readFileString(t, normalPathDispatchContractPath(t))
	legacySkill := readFileString(t, legacyClaudeTeamSkillPath(t))

	for _, marker := range legacyMachineryStrings {
		if strings.Contains(contract, marker) {
			t.Errorf("claude-fo-dispatch.md inlines legacy machinery marker %q — it belongs only in the conditionally-loaded using-legacy-claude-team skill", marker)
		}
	}
	for _, gone := range deletedFromMergedFloor {
		if strings.Contains(contract, gone) {
			t.Errorf("claude-fo-dispatch.md still carries %q — the rewrite removes this stale hint outright", gone)
		}
		if strings.Contains(legacySkill, gone) {
			t.Errorf("using-legacy-claude-team/SKILL.md carries %q — this stale hint is removed, not relocated", gone)
		}
	}
}

// TestNormalPathContractNamesLegacySkill is the load-point half of the layering:
// the normal-path dispatch contract retains exactly the probe branch that names
// the legacy skill, so the conditional load actually resolves. (The reference-
// closure check in boot_resident_closure_test.go proves the named skill exists on
// disk via the boot-resident runtime adapter that also names it.)
func TestNormalPathContractNamesLegacySkill(t *testing.T) {
	contract := readFileString(t, normalPathDispatchContractPath(t))
	if !strings.Contains(contract, "spacedock:using-legacy-claude-team") {
		t.Errorf("claude-fo-dispatch.md must name spacedock:using-legacy-claude-team so the legacy probe branch can load it")
	}
}

// TestLegacyConsumerRetiredButPathLives is AC-5's removal-trigger binding after
// the #395 pin retired. The legacy removal trigger has TWO conditions: (1) no live
// lane drives the legacy branch, AND (2) no runtime the FO targets still exposes
// TeamCreate. Unpinning the live lane (this task) satisfies only (1) — a user on
// installer-`stable` (2.1.170, still TeamCreate-capable) still hits the legacy
// path, so condition (2) is unmet and the path is NOT dead code. This test binds
// two independent sources so it stays a structural guard, not a self-referential
// prose-grep:
//   - the REAL .github/workflows/runtime-live-e2e.yml NO LONGER pins a
//     TeamCreate-capable version (SPACEDOCK_PINNED_CLAUDE_VERSION absent) — the
//     live consumer is RETIRED; AND
//   - the legacy skill and its dispatch-contract load line STILL EXIST — the
//     best-effort legacy path is RETAINED, not deleted.
//
// The externally-checkable proxy for legacy DELETION is no longer this CI pin (now
// retired); it is the future condition "even installer-`stable` Claude is ≥ the
// merged floor" — a separate trigger, out of this task's scope.
func TestLegacyConsumerRetiredButPathLives(t *testing.T) {
	workflow := readFileString(t, filepath.Join(repoRoot(t), ".github", "workflows", "runtime-live-e2e.yml"))
	if strings.Contains(workflow, "SPACEDOCK_PINNED_CLAUDE_VERSION") {
		t.Errorf("the live-e2e workflow still pins SPACEDOCK_PINNED_CLAUDE_VERSION — the #395 pin should be retired so the live lane floats to the merged floor (the live legacy consumer is retired)")
	}

	contract := readFileString(t, normalPathDispatchContractPath(t))
	if !strings.Contains(contract, "spacedock:using-legacy-claude-team") {
		t.Errorf("claude-fo-dispatch.md no longer names spacedock:using-legacy-claude-team — the best-effort legacy path must be RETAINED (a user on installer-stable still hits it); deletion is a separate trigger")
	}

	legacySkill := readFileString(t, legacyClaudeTeamSkillPath(t))
	if strings.TrimSpace(legacySkill) == "" {
		t.Errorf("using-legacy-claude-team/SKILL.md is empty or missing — the best-effort legacy path must be RETAINED until even installer-stable Claude clears the merged floor")
	}
}
