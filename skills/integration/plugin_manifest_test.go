// ABOUTME: Structural test over the vendored repo plugin manifest — name, skills
// ABOUTME: path, and a version that parses as major.minor under minor coupling.
package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// repoRoot is the project root: the parent of the skills/ dir this test package
// lives inside. The vendored plugin manifest lives at repoRoot/.claude-plugin/.
func repoRoot(t *testing.T) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// TestVendoredManifestVersionParses locks the manifest<->binary drift check
// under minor-version coupling: the vendored .claude-plugin/plugin.json's
// `version` field (the compatibility declaration itself — D2) parses as a
// well-formed major.minor semver. That version's binding to the FO shared-core's
// stamped minor is pinned by the internal/contractlint sync test — the
// sanctioned home for a check that reads the prose file this manifest binds
// against.
func TestVendoredManifestVersionParses(t *testing.T) {
	manifestPath := filepath.Join(repoRoot(t), ".claude-plugin", "plugin.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read vendored manifest %s: %v", manifestPath, err)
	}
	var m struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Skills  string `json:"skills"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse vendored manifest %s: %v", manifestPath, err)
	}
	if m.Name != "spacedock" {
		t.Errorf("manifest name = %q, want spacedock", m.Name)
	}
	if m.Skills != "./skills/" {
		t.Errorf("manifest skills = %q, want ./skills/", m.Skills)
	}
	if _, _, ok := contract.ParseMajorMinor(m.Version); !ok {
		t.Fatalf("manifest version %q does not parse as major.minor", m.Version)
	}
}
