// ABOUTME: AC-1 structural guard — the legacy TeamCreate path is fully retired: no
// ABOUTME: skill directory, no load token, no select:TeamCreate probe, and the
// ABOUTME: normal-path dispatch contract inlines none of the legacy machinery.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// legacyMachineryStrings are the legacy TeamCreate/TeamDelete machinery markers the
// retirement removes entirely — the recovery ladder, the bounded terminal-teardown
// apparatus, the registry-desync issue, and the NAME_PATTERN naming constraint. They
// lived only in the now-deleted using-legacy-claude-team skill; a current-host
// session reads the normal-path dispatch contract (claude-fo-dispatch.md), so a
// re-inlined marker would be real context weight every session pays.
var legacyMachineryStrings = []string{
	"TERMINAL_TEARDOWN_BOUNDED",
	"36806",
	"priority-ordered ladder",
	"fresh-suffixed",
	"TeamCreate recovery procedure",
	"NAME_PATTERN",
}

// deletedFromMergedFloor are strings the rewrite REMOVES outright (not relocated):
// inter-agent communication works without the `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS`
// opt-in flag (OQ-2), so any captain-facing "set the flag to enable team mode" hint
// is stale. The whole env-var token must be absent from the normal-path dispatch
// contract — re-introducing it resurrects the obsolete advice. The CI lane in
// runtime-live-e2e.yml legitimately sets the flag (test config, not captain-facing
// skill prose), so it is out of this check's scope (skills/ only).
var deletedFromMergedFloor = []string{
	"CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS",
}

// normalPathDispatchContractPath is claude-fo-dispatch.md — the Claude FO dispatch
// contract a current-host session reads at first dispatch. It carries the
// inter-agent-communication shape and the idle/failure-retry discipline inline; it
// must inline none of the legacy TeamCreate machinery.
func normalPathDispatchContractPath(t *testing.T) string {
	return filepath.Join(skillsRoot(t), "first-officer", "references", "claude-fo-dispatch.md")
}

// legacyClaudeTeamSkillPath is the retired legacy skill path, asserted absent on
// disk by TestLegacyTeamCreatePathFullyRetired.
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

// TestNormalPathContractInlinesNoLegacyMachinery is AC-1's layering invariant
// (structural-ABSENCE): the normal-path dispatch contract claude-fo-dispatch.md
// inlines NONE of the legacy TeamCreate machinery markers, and does not carry the
// stale enable-the-flag hint. With the legacy skill retired, every current-host
// session reads only this contract, so a re-inlined marker is real context weight
// every such session pays. This asserts ABSENCE from the normal-path contract; the
// legacy skill's full retirement is bound by TestLegacyTeamCreatePathFullyRetired.
func TestNormalPathContractInlinesNoLegacyMachinery(t *testing.T) {
	contract := readFileString(t, normalPathDispatchContractPath(t))

	for _, marker := range legacyMachineryStrings {
		if strings.Contains(contract, marker) {
			t.Errorf("claude-fo-dispatch.md inlines legacy machinery marker %q — the legacy TeamCreate path is retired; it must not reappear", marker)
		}
	}
	for _, gone := range deletedFromMergedFloor {
		if strings.Contains(contract, gone) {
			t.Errorf("claude-fo-dispatch.md still carries %q — the rewrite removes this stale hint outright", gone)
		}
	}
}

// TestNormalPathContractNamesNoLegacySkill is the inverted load-point half: the
// legacy probe branch is deleted, so the normal-path dispatch contract no longer
// names the retired legacy skill token. A re-added token would dangle — no skill
// file resolves it — which boot_resident_closure_test.go's closure check also catches.
func TestNormalPathContractNamesNoLegacySkill(t *testing.T) {
	contract := readFileString(t, normalPathDispatchContractPath(t))
	if strings.Contains(contract, "spacedock:using-legacy-claude-team") {
		t.Errorf("claude-fo-dispatch.md still names spacedock:using-legacy-claude-team — the legacy probe branch is retired and must not name the deleted skill")
	}
}

// TestLegacyTeamCreatePathFullyRetired is AC-1's mechanical guard: the legacy
// TeamCreate path is fully retired from the shipped FO contract. It binds two
// independent sources — the filesystem (the skill directory is gone) and the
// contract text (no select:TeamCreate probe survives on the FO dispatch surface) —
// so it stays a structural guard, not a self-referential prose-grep. Re-adding the
// skill dir or the probe reds this test.
func TestLegacyTeamCreatePathFullyRetired(t *testing.T) {
	skillPath := legacyClaudeTeamSkillPath(t)
	if _, err := os.Stat(skillPath); !os.IsNotExist(err) {
		t.Errorf("using-legacy-claude-team/SKILL.md still resolves on disk (stat err=%v) — AC-1 retires the legacy skill entirely", err)
	}
	if _, err := os.Stat(filepath.Dir(skillPath)); !os.IsNotExist(err) {
		t.Errorf("the using-legacy-claude-team skill directory still exists (stat err=%v) — retire the directory, not just the file", err)
	}

	contract := readFileString(t, normalPathDispatchContractPath(t))
	if strings.Contains(contract, "select:TeamCreate") {
		t.Errorf("claude-fo-dispatch.md still carries a select:TeamCreate probe on the FO dispatch surface — the legacy TeamCreate probe is retired")
	}
}
