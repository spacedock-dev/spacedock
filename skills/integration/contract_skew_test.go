// ABOUTME: AC-1 skew test — the live repo manifest range rejects a contract-1
// ABOUTME: binary (the published v0.22.0) and admits the current CONTRACT_VERSION.
package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// liveRequiresContract reads the requires-contract range declared by the named
// repo plugin manifest (".claude-plugin" or ".codex-plugin").
func liveRequiresContract(t *testing.T, pluginDir string) string {
	t.Helper()
	path := filepath.Join(repoRoot(t), pluginDir, "plugin.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var m struct {
		RequiresContract string `json:"requires-contract"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return m.RequiresContract
}

// TestContract1BinaryRejectedByLiveRange locks AC-1: a contract-1 binary (the
// published v0.22.0, which reports `contract 1`) is REJECTED by the live skill
// range with the too-old-binary verdict, while the current binary's
// CONTRACT_VERSION is admitted as compatible. Both plugin manifests are exercised
// so the Claude and Codex ranges agree. The test reads the live range and the
// live CONTRACT_VERSION — leaving either side at contract 1 / `>=1,<2` makes the
// contract-1 case Compatible and fails here.
func TestContract1BinaryRejectedByLiveRange(t *testing.T) {
	const contract1 = 1 // the published v0.22.0 binary's reported contract
	for _, pluginDir := range []string{".claude-plugin", ".codex-plugin"} {
		raw := liveRequiresContract(t, pluginDir)

		stale := contract.Compare(contract1, raw, "claude", "0.22.0", "0.23.0")
		if stale.Verdict != contract.TooOldBinary {
			t.Fatalf("%s range %q: contract-1 binary verdict = %v, want too-old-binary (a contract-1 v0.22.0 binary must be rejected)",
				pluginDir, raw, stale.Verdict)
		}

		current := contract.Compare(contract.CONTRACT_VERSION, raw, "claude", "0.23.0", "0.23.0")
		if current.Verdict != contract.Compatible {
			t.Fatalf("%s range %q: current CONTRACT_VERSION=%d verdict = %v, want compatible",
				pluginDir, raw, contract.CONTRACT_VERSION, current.Verdict)
		}
	}
}
