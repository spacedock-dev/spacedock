// ABOUTME: AC-5 structural guard — the legacy TeamCreate machinery is layered
// ABOUTME: into a conditionally-loaded reference, with an externally-anchored
// ABOUTME: removal trigger; SKILL.md inlines none of it.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// legacyMachineryStrings are the legacy TeamCreate/TeamDelete machinery markers
// that AC-5 requires to live ONLY in references/legacy-teamcreate.md, never
// inline in the merged SKILL.md — the recovery ladder, the bounded
// terminal-teardown apparatus, the registry-desync issue, and the NAME_PATTERN
// naming constraint. Each is machinery a host that still has TeamCreate needs,
// so it MOVES to the reference rather than vanishing; a merged session that
// takes the no-match probe branch never loads it.
var legacyMachineryStrings = []string{
	"TERMINAL_TEARDOWN_BOUNDED",
	"36806",
	"priority-ordered ladder",
	"fresh-suffixed",
	"TeamCreate recovery procedure",
	"NAME_PATTERN",
}

// deletedFromMergedFloor are strings the merged rewrite REMOVES outright (not
// relocated): the stale `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` bare-mode UX
// hint, which OQ-2 proved points at a flag that no longer gates the merged
// concurrency channel. They must be absent from BOTH the SKILL.md and the legacy
// reference — re-introducing the hint anywhere resurrects the obsolete advice.
var deletedFromMergedFloor = []string{
	"CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 and restart",
}

func usingClaudeTeamSkillPath(t *testing.T) string {
	return filepath.Join(skillsRoot(t), "using-claude-team", "SKILL.md")
}

func legacyTeamCreateReferencePath(t *testing.T) string {
	return filepath.Join(skillsRoot(t), "using-claude-team", "references", "legacy-teamcreate.md")
}

func readFileString(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

// TestSkillMdInlinesNoLegacyMachinery is AC-5 half (a), the layering invariant:
// the merged using-claude-team/SKILL.md inlines NONE of the legacy TeamCreate
// machinery — every marker lives only in the conditionally-loaded reference. A
// merged session that takes the no-match probe branch never loads the reference,
// so a re-inlined marker is real context weight every merged session pays.
func TestSkillMdInlinesNoLegacyMachinery(t *testing.T) {
	skill := readFileString(t, usingClaudeTeamSkillPath(t))
	reference := readFileString(t, legacyTeamCreateReferencePath(t))

	for _, marker := range legacyMachineryStrings {
		if strings.Contains(skill, marker) {
			t.Errorf("SKILL.md inlines legacy machinery marker %q — it belongs only in references/legacy-teamcreate.md", marker)
		}
		if !strings.Contains(reference, marker) {
			t.Errorf("references/legacy-teamcreate.md is missing legacy machinery marker %q — the machinery must move there, not vanish", marker)
		}
	}
	for _, gone := range deletedFromMergedFloor {
		if strings.Contains(skill, gone) {
			t.Errorf("SKILL.md still carries %q — the merged rewrite removes this stale hint outright", gone)
		}
		if strings.Contains(reference, gone) {
			t.Errorf("references/legacy-teamcreate.md carries %q — this stale hint is removed, not relocated", gone)
		}
	}
}

// TestSkillMdNamesLegacyReference is the load-point half of the layering: the
// merged SKILL.md retains exactly the one probe branch that names the reference,
// so the conditional load actually resolves. (The reference-closure check in
// structural_checks_test.go proves the named file exists on disk.)
func TestSkillMdNamesLegacyReference(t *testing.T) {
	skill := readFileString(t, usingClaudeTeamSkillPath(t))
	if !strings.Contains(skill, "references/legacy-teamcreate.md") {
		t.Errorf("SKILL.md must name references/legacy-teamcreate.md so the legacy probe branch can load it")
	}
}

// TestLegacyRemovalTriggerIsExternallyAnchored is AC-5 half (b): the removal
// trigger in references/legacy-teamcreate.md references an external, checkable
// condition — the real SPACEDOCK_PINNED_CLAUDE_VERSION pin in the live-e2e
// workflow — and that pin currently still pins a team-tools-capable version
// (≤ 2.1.177, the last release exposing native team tools). When the pin moves
// past 2.1.177 this assertion goes red, which is the signal that the legacy
// branch has lost its only live consumer and the file may be deleted.
func TestLegacyRemovalTriggerIsExternallyAnchored(t *testing.T) {
	reference := readFileString(t, legacyTeamCreateReferencePath(t))
	if !strings.Contains(reference, "SPACEDOCK_PINNED_CLAUDE_VERSION") {
		t.Errorf("removal trigger must name the SPACEDOCK_PINNED_CLAUDE_VERSION pin (the external, checkable condition), not just prose")
	}
	if !strings.Contains(reference, "runtime-live-e2e.yml") {
		t.Errorf("removal trigger must name the workflow file that carries the pin")
	}

	workflow := readFileString(t, filepath.Join(repoRoot(t), ".github", "workflows", "runtime-live-e2e.yml"))
	if !strings.Contains(workflow, `SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"`) {
		t.Errorf("the live-e2e workflow no longer pins SPACEDOCK_PINNED_CLAUDE_VERSION to 2.1.177 — the legacy branch may have lost its live consumer; re-evaluate the legacy-teamcreate.md removal trigger")
	}
}
