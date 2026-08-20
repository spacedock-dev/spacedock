// ABOUTME: Keeps the retired legacy TeamCreate skill path absent from the shipped topology,
// ABOUTME: and no shipped skill text re-teaching a TeamCreate imperative call.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLegacyTeamCreateSkillPathIsRetired(t *testing.T) {
	skillDir := filepath.Join(skillsRoot(t), "using-legacy-claude-team")
	if _, err := os.Stat(filepath.Join(skillDir, "SKILL.md")); !os.IsNotExist(err) {
		t.Errorf("legacy TeamCreate SKILL.md still resolves on disk: %v", err)
	}
	if _, err := os.Stat(skillDir); !os.IsNotExist(err) {
		t.Errorf("legacy TeamCreate skill directory still exists: %v", err)
	}
}

// TestNoShippedSkillInstructsTeamCreate is AC-4's extended invariant: no shipped
// skill file re-teaches the retired legacy dispatch mode by instructing an actual
// `TeamCreate(...)` call. This is the exact regression a review pass missed once
// already — commission's Step 3 survived #549's sweep instructing `TeamCreate(...)`
// for over a cycle before this task caught it. The check is scoped to the literal
// call syntax (`TeamCreate(`), not the bare word: legitimate explanatory prose
// elsewhere (e.g. claude-fo-dispatch.md's reconcile section, "no TeamCreate name
// to pass") mentions the retired tool by name without instructing its use, and
// must not false-positive here.
func TestNoShippedSkillInstructsTeamCreate(t *testing.T) {
	root := skillsRoot(t)
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if strings.Contains(string(data), "TeamCreate(") {
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				rel = path
			}
			t.Errorf("skills/%s instructs a TeamCreate(...) call — legacy team mode is retired; no shipped skill should re-teach it", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
}
