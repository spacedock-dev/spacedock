// ABOUTME: AC-2 manifest tests — the plugin branch carries no marketplace.json
// ABOUTME: (Model B), and .codex-plugin/plugin.json requires-contract brackets binary.
package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// TestPluginBranchCarriesNoMarketplaceManifest locks AC-2: under Model B the
// marketplace manifest moved OUT of the plugin branch into a separate marketplace
// repo (the two channels — stable pinned to a release tag, edge tracking next HEAD
// — are entries of one source in THAT repo's manifest). So the plugin branch
// carries NO .claude-plugin/marketplace.json; with it gone, there is no in-branch
// source.ref for a release to re-settle. The plugin's own
// .claude-plugin/plugin.json stays — only the marketplace manifest moved.
func TestPluginBranchCarriesNoMarketplaceManifest(t *testing.T) {
	marketplace := filepath.Join(repoRoot(t), ".claude-plugin", "marketplace.json")
	if _, err := os.Stat(marketplace); err == nil {
		t.Fatalf("%s is present; Model B moves the marketplace manifest to the separate marketplace repo, so the plugin branch must carry no marketplace.json", marketplace)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat %s: %v", marketplace, err)
	}

	plugin := filepath.Join(repoRoot(t), ".claude-plugin", "plugin.json")
	if _, err := os.Stat(plugin); err != nil {
		t.Fatalf("plugin manifest %s missing: %v (only marketplace.json should move out)", plugin, err)
	}
}

// TestCodexManifestBracketsContractVersion locks AC-2's Codex half: the new
// .codex-plugin/plugin.json carries a requires-contract that parses via the real
// ParseRange and brackets the binary's CONTRACT_VERSION (so `spacedock doctor
// --host codex` resolves a compatible manifest), names the plugin `spacedock`,
// and points skills at ./skills/.
func TestCodexManifestBracketsContractVersion(t *testing.T) {
	path := filepath.Join(repoRoot(t), ".codex-plugin", "plugin.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read codex manifest %s: %v", path, err)
	}
	var m struct {
		Name             string `json:"name"`
		RequiresContract string `json:"requires-contract"`
		Skills           string `json:"skills"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse codex manifest: %v", err)
	}
	if m.Name != "spacedock" {
		t.Errorf("codex manifest name = %q, want spacedock", m.Name)
	}
	if m.Skills != "./skills/" {
		t.Errorf("codex manifest skills = %q, want ./skills/", m.Skills)
	}
	if m.RequiresContract == "" {
		t.Fatalf("codex manifest has no requires-contract (AC-2 requires it for doctor --host codex)")
	}
	lo, hi, err := contract.ParseRange(m.RequiresContract)
	if err != nil {
		t.Fatalf("codex requires-contract %q does not parse: %v", m.RequiresContract, err)
	}
	if !(lo <= contract.CONTRACT_VERSION && contract.CONTRACT_VERSION < hi) {
		t.Fatalf("codex requires-contract %s does not bracket CONTRACT_VERSION=%d", m.RequiresContract, contract.CONTRACT_VERSION)
	}
}
