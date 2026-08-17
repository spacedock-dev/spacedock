// ABOUTME: AC-5 live proof — the probe-1 lane re-run against the NEW non-destructive
// ABOUTME: install sequences: a co-hosted dependent plugin must survive a channel install.
package cli

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestClaudeChannelInstallLeavesCoHostedPluginInstalled locks AC-5 on claude: a
// dependent plugin installed from the SAME shared marketplace as spacedock (the
// subspace/cargento shape — round 2's probe 1 measured the OLD `plugin
// marketplace remove` step cascade-uninstalling it) must still be
// installed/enabled after a full channel install through the production
// non-destructive sequence (installArgvSequence, no `marketplace remove` step).
// Baseline that moves the wrong way: the pre-round-2 sequence, which removed the
// shared marketplace record and took the dependent down with it. Skips when
// `claude` is not on PATH; kept hermetic via CLAUDE_CONFIG_DIR isolation.
func TestClaudeChannelInstallLeavesCoHostedPluginInstalled(t *testing.T) {
	claudeBin, err := exec.LookPath("claude")
	if err != nil {
		t.Skip("claude not on PATH; co-hosted survival test requires the host CLI")
	}

	tmp := t.TempDir()
	marketplace := buildLocalMarketplaceWithDependent(t, tmp, "claude")
	configDir := filepath.Join(tmp, "config")
	cacheDir := filepath.Join(tmp, "cache")
	mustMkdir(t, configDir)
	mustMkdir(t, cacheDir)
	// t.Setenv (not just an env slice passed to runHost) — execHost.Install shells
	// out via exec.Command with no explicit Env, so it inherits the process
	// environment. Isolation must land on the process itself or Install() runs
	// against the real ambient claude config instead of this fixture.
	t.Setenv("CLAUDE_CONFIG_DIR", configDir)
	t.Setenv("CLAUDE_CODE_PLUGIN_CACHE_DIR", cacheDir)
	env := os.Environ()

	// Seed: a dependent plugin installed from the shared marketplace — probe 1's
	// setup (a second plugin sourced from the same marketplace spacedock installs
	// from).
	runHost(t, claudeBin, env, "plugin", "marketplace", "add", marketplace)
	runHost(t, claudeBin, env, "plugin", "install", "dependent@spacedock")

	// The full channel install, through the production seam — the exact argv
	// this session's Install() now issues, not a re-derived shape.
	out, err := execHost{}.Install("claude", marketplace, "main")
	if err != nil {
		t.Fatalf("channel Install failed: %v\nout=%q", err, out)
	}

	listOut := runHost(t, claudeBin, env, "plugin", "list", "--json")
	var entries []struct {
		ID          string `json:"id"`
		InstallPath string `json:"installPath"`
		Enabled     bool   `json:"enabled"`
	}
	if err := json.Unmarshal([]byte(listOut), &entries); err != nil {
		t.Fatalf("parse plugin list --json: %v\n%s", err, listOut)
	}
	found := false
	for _, e := range entries {
		if e.ID == "dependent@spacedock" && e.InstallPath != "" && e.Enabled {
			found = true
		}
	}
	if !found {
		t.Fatalf("co-hosted dependent@spacedock is not installed/enabled after the channel install (the cascade-uninstall probe 1 measured has regressed):\n%s", listOut)
	}
}

// TestCodexChannelInstallLeavesCoHostedPluginInstalled locks AC-5 on codex,
// mirroring the claude case above: a dependent plugin installed from the SAME
// marketplace name as spacedock's channel (the captain's real
// `subspace@spacedock` shape) must still be installed after a full channel
// install through codexInstallArgvSequence — which removes PLUGIN content, never
// the shared marketplace record. Skips when `codex` is not on PATH; hermetic via
// CODEX_HOME isolation.
func TestCodexChannelInstallLeavesCoHostedPluginInstalled(t *testing.T) {
	codexBin, err := exec.LookPath("codex")
	if err != nil {
		t.Skip("codex not on PATH; co-hosted survival test requires the host CLI")
	}
	saved := devBranch
	devBranch = "main"
	defer func() { devBranch = saved }()

	tmp := t.TempDir()
	marketplace := buildLocalMarketplaceWithDependent(t, tmp, "codex")
	codexHomeDir := filepath.Join(tmp, "codexhome")
	mustMkdir(t, codexHomeDir)
	t.Setenv("CODEX_HOME", codexHomeDir)

	// Seed: the dependent, installed from the same marketplace name the channel
	// install below targets.
	runHost(t, codexBin, os.Environ(), "plugin", "marketplace", "add", marketplace)
	runHost(t, codexBin, os.Environ(), "plugin", "add", "dependent@spacedock")

	out, err := execHost{}.Install("codex", marketplace, "main")
	if err != nil {
		t.Fatalf("channel Install failed: %v\nout=%q", err, out)
	}

	listOut := runHost(t, codexBin, os.Environ(), "plugin", "list", "--json")
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
		if e.PluginID == "dependent@spacedock" && e.Installed {
			found = true
		}
	}
	if !found {
		t.Fatalf("co-hosted dependent@spacedock is not installed after the channel install (out=%q):\n%s", out, listOut)
	}
}

// buildLocalMarketplaceWithDependent writes a local-path marketplace named
// `spacedock` (the stable channel) hosting TWO plugin entries — `spacedock`
// itself and a `dependent` plugin sourced from the same marketplace, the shape
// this task's cascade-uninstall harm was measured against (subspace/cargento
// co-hosted on the real `spacedock` marketplace). host selects the manifest
// subdirectory (`.claude-plugin` / `.codex-plugin`); codex reads the marketplace
// manifest from `.claude-plugin/marketplace.json` regardless (it reuses the
// claude layout, per buildLocalCodexMarketplace).
func buildLocalMarketplaceWithDependent(t *testing.T, root, host string) string {
	t.Helper()
	marketplace := filepath.Join(root, "marketplace")
	plugin := filepath.Join(marketplace, "spacedock")
	dependent := filepath.Join(marketplace, "dependent")
	manifestDir := "." + host + "-plugin"

	mustMkdir(t, filepath.Join(marketplace, ".claude-plugin"))
	mustMkdir(t, filepath.Join(plugin, manifestDir))
	mustMkdir(t, filepath.Join(plugin, "skills", "demo"))
	mustMkdir(t, filepath.Join(dependent, manifestDir))
	mustMkdir(t, filepath.Join(dependent, "skills", "demo"))

	mustWrite(t, filepath.Join(marketplace, ".claude-plugin", "marketplace.json"), `{
  "name": "spacedock",
  "owner": { "name": "CL Kao" },
  "plugins": [
    { "name": "spacedock", "source": "./spacedock", "description": "test", "category": "workflow" },
    { "name": "dependent", "source": "./dependent", "description": "co-hosted dependent", "category": "workflow" }
  ]
}
`)
	mustWrite(t, filepath.Join(plugin, manifestDir, "plugin.json"),
		`{ "name": "spacedock", "version": "`+displayVersion()+`", "skills": "./skills/" }`+"\n")
	mustWrite(t, filepath.Join(plugin, "skills", "demo", "SKILL.md"), "---\nname: demo\ndescription: demo skill\n---\ndemo\n")
	mustWrite(t, filepath.Join(dependent, manifestDir, "plugin.json"),
		`{ "name": "dependent", "version": "1.0.0", "skills": "./skills/" }`+"\n")
	mustWrite(t, filepath.Join(dependent, "skills", "demo", "SKILL.md"), "---\nname: demo\ndescription: demo skill\n---\ndemo\n")
	return marketplace
}
