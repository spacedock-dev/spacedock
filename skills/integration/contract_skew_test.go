// ABOUTME: Minor-skew test — the live repo manifest's own version rejects an
// ABOUTME: old-minor binary and admits a binary sharing its own minor.
package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// liveManifestVersion reads the `version` field declared by the named repo
// plugin manifest (".claude-plugin" or ".codex-plugin").
func liveManifestVersion(t *testing.T, pluginDir string) string {
	t.Helper()
	path := filepath.Join(repoRoot(t), pluginDir, "plugin.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var m struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return m.Version
}

// TestOldMinorBinaryRejectedByLiveManifest locks D1: a binary one minor BEHIND
// the live repo's own manifest version is REJECTED with the too-old-binary
// verdict, while a binary sharing the manifest's own minor is admitted as
// compatible. Both plugin manifests are exercised so the Claude and Codex
// declarations agree. The test reads the LIVE manifest version — a future
// minor bump that forgets to keep the binary in step still classifies an old
// binary as too-old-binary here (the property this locks), so the test
// self-adjusts with every release rather than pinning a stale integer.
func TestOldMinorBinaryRejectedByLiveManifest(t *testing.T) {
	for _, pluginDir := range []string{".claude-plugin", ".codex-plugin"} {
		pluginVersion := liveManifestVersion(t, pluginDir)
		major, minor, ok := contract.ParseMajorMinor(pluginVersion)
		if !ok {
			t.Fatalf("%s version %q does not parse as major.minor", pluginDir, pluginVersion)
		}
		if minor == 0 {
			t.Fatalf("%s minor is 0 — no room for an old-minor binary fixture one minor behind", pluginDir)
		}
		oldBinary := fmt.Sprintf("%d.%d.0", major, minor-1)

		stale := contract.Compare("claude", pluginVersion, oldBinary)
		if stale.Verdict != contract.TooOldBinary {
			t.Fatalf("%s: old-minor binary %s verdict = %v, want too-old-binary (an old-minor binary must be rejected)",
				pluginDir, oldBinary, stale.Verdict)
		}

		current := contract.Compare("claude", pluginVersion, pluginVersion)
		if current.Verdict != contract.Compatible {
			t.Fatalf("%s: same-minor binary %s verdict = %v, want compatible",
				pluginDir, pluginVersion, current.Verdict)
		}
	}
}
