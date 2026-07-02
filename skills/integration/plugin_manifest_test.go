// ABOUTME: Structural test over the vendored repo plugin manifest — the manifest
// ABOUTME: version parses as major.minor and the D4 tombstone stays frozen.
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
// `version` field (the compatibility declaration itself — D2, no separate
// requires-contract range is read by this binary) parses as a well-formed
// major.minor semver, and the D4 cross-era tombstone (requires-contract
// ">=3,<4") stays present and frozen so an integer-era binary still aborts with
// the correct "upgrade the binary" remedy.
func TestVendoredManifestVersionParses(t *testing.T) {
	manifestPath := filepath.Join(repoRoot(t), ".claude-plugin", "plugin.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read vendored manifest %s: %v", manifestPath, err)
	}
	var m struct {
		Name             string `json:"name"`
		Version          string `json:"version"`
		RequiresContract string `json:"requires-contract"`
		Skills           string `json:"skills"`
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
	const frozenTombstone = ">=3,<4"
	if m.RequiresContract != frozenTombstone {
		t.Fatalf("manifest requires-contract = %q, want the frozen D4 tombstone %q (never edited again)", m.RequiresContract, frozenTombstone)
	}
}
