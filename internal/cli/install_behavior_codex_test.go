// ABOUTME: AC-5 behavioral codex install — a real isolated-CODEX_HOME run of
// ABOUTME: execHost.Install against a local-path marketplace, observing on-disk state.
package cli

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestCodexPluginInstallIsHostNative runs the REAL execHost{}.Install("codex", …)
// against an isolated CODEX_HOME + a local-path marketplace, then observes that
// (a) `codex plugin list --json` reports spacedock@spacedock with installed:true,
// and (b) the cache manifest the resolver looks for exists on disk. The source is
// the local marketplace path with an empty branch, so the install carries no
// --ref and stays hermetic/offline. This proves the real codex CLI accepts the
// 4-step Install sequence (the two cleanup steps run first against a fresh
// CODEX_HOME — plugin remove exits 0, marketplace remove exits 1 tolerated — then
// the two pin steps install) unattended, not just the stub. Skips when `codex` is
// not on PATH; kept hermetic by env isolation, not a mock.
func TestCodexPluginInstallIsHostNative(t *testing.T) {
	if _, err := exec.LookPath("codex"); err != nil {
		t.Skip("codex not on PATH; behavioral install test requires the host CLI")
	}

	tmp := t.TempDir()
	marketplace := buildLocalCodexMarketplace(t, tmp)
	codexHomeDir := filepath.Join(tmp, "codexhome")
	mustMkdir(t, codexHomeDir)
	t.Setenv("CODEX_HOME", codexHomeDir)

	// The real Install runs in-process through the production seam: the codex arm
	// shells the 4-step sequence against the isolated CODEX_HOME (CODEX_HOME is read
	// from the env by the codex CLI). Empty branch → no --ref → fully offline.
	out, err := execHost{}.Install("codex", marketplace, "")
	if err != nil {
		t.Fatalf("execHost.Install(codex) failed: %v\nout=%q", err, out)
	}

	// (a) The real codex CLI reports the plugin installed. The --json schema is
	// {"installed":[{"pluginId":…,"installed":true}]}; the installed flag is the
	// independent source of truth, not a substring of our source.
	listOut := runHost(t, mustLookPath(t, "codex"), os.Environ(), "plugin", "list", "--json")
	var listing struct {
		Installed []struct {
			PluginID  string `json:"pluginId"`
			Installed bool   `json:"installed"`
		} `json:"installed"`
	}
	if err := json.Unmarshal([]byte(listOut), &listing); err != nil {
		t.Fatalf("parse codex plugin list --json: %v\n%s", err, listOut)
	}
	found := false
	for _, e := range listing.Installed {
		if e.PluginID == "spacedock@spacedock" && e.Installed {
			found = true
		}
	}
	if !found {
		t.Fatalf("codex plugin list --json did not report spacedock@spacedock installed:true:\n%s", listOut)
	}

	// (b) The cache manifest exists at exactly the resolver's path
	// (<CODEX_HOME>/plugins/cache/spacedock/spacedock/<version>/.codex-plugin/plugin.json).
	// resolveCodexManifest must therefore resolve a non-empty manifest path.
	manifest, err := execHost{}.resolveCodexManifest()
	if err != nil {
		t.Fatalf("resolveCodexManifest after install: %v", err)
	}
	if manifest == "" {
		t.Fatalf("resolveCodexManifest returned empty after a real install")
	}
	if _, statErr := os.Stat(manifest); statErr != nil {
		t.Fatalf("resolved manifest %s does not exist: %v", manifest, statErr)
	}
}

// buildLocalCodexMarketplace writes a minimal valid local-path marketplace under
// root and returns the marketplace directory. Codex reads the marketplace
// manifest from .claude-plugin/marketplace.json (it reuses the claude manifest
// layout) and the plugin manifest from the plugin's .codex-plugin/plugin.json.
// The plugin manifest carries a requires-contract bracketing CONTRACT_VERSION.
func buildLocalCodexMarketplace(t *testing.T, root string) string {
	t.Helper()
	marketplace := filepath.Join(root, "marketplace")
	plugin := filepath.Join(marketplace, "spacedock")
	mustMkdir(t, filepath.Join(marketplace, ".claude-plugin"))
	mustMkdir(t, filepath.Join(plugin, ".codex-plugin"))
	mustMkdir(t, filepath.Join(plugin, "skills", "demo"))

	mustWrite(t, filepath.Join(marketplace, ".claude-plugin", "marketplace.json"), `{
  "name": "spacedock",
  "owner": { "name": "CL Kao" },
  "plugins": [
    { "name": "spacedock", "source": "./spacedock", "description": "test", "category": "workflow" }
  ]
}
`)
	mustWrite(t, filepath.Join(plugin, ".codex-plugin", "plugin.json"), `{ "name": "spacedock", "version": "0.0.0", "requires-contract": ">=1,<2", "skills": "./skills/" }
`)
	mustWrite(t, filepath.Join(plugin, "skills", "demo", "SKILL.md"), "---\nname: demo\ndescription: demo skill\n---\ndemo\n")
	return marketplace
}

// mustLookPath resolves bin on PATH or fails the test.
func mustLookPath(t *testing.T, bin string) string {
	t.Helper()
	p, err := exec.LookPath(bin)
	if err != nil {
		t.Fatalf("look path %s: %v", bin, err)
	}
	return p
}
