// ABOUTME: Keeps the retired legacy TeamCreate skill path absent from the shipped topology.
package contractlint

import (
	"os"
	"path/filepath"
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
