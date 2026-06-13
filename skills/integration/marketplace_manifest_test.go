// ABOUTME: Manifest tests — main carries a transitional bridge marketplace.json
// ABOUTME: (Model B removal deferred), and .codex-plugin/plugin.json brackets binary.
package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// TestMainCarriesMarketplaceBridgeManifest guards the transitional bridge: the
// released v0.20.0 binary resolves its install from main's
// .claude-plugin/marketplace.json (`claude plugin marketplace add
// spacedock-dev/spacedock@main`). Model B's removal of this manifest — it moves to the
// standalone marketplace repo, retiring the in-branch source.ref re-settle — is
// deferred until the v0.20.1 cutover ships the binary that points at that repo and
// existing v0.20.0 installs have migrated; until then the bridge MUST stay on main, or
// v0.20.0 installs break. The plugin's own .claude-plugin/plugin.json stays alongside it.
func TestMainCarriesMarketplaceBridgeManifest(t *testing.T) {
	marketplace := filepath.Join(repoRoot(t), ".claude-plugin", "marketplace.json")
	if _, err := os.Stat(marketplace); err != nil {
		t.Fatalf("bridge marketplace manifest %s missing: %v — the released v0.20.0 binary resolves its install from main's marketplace.json; removing it before the v0.20.1 cutover breaks v0.20.0 installs", marketplace, err)
	}

	plugin := filepath.Join(repoRoot(t), ".claude-plugin", "plugin.json")
	if _, err := os.Stat(plugin); err != nil {
		t.Fatalf("plugin manifest %s missing: %v", plugin, err)
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
