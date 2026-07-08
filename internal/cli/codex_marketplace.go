// ABOUTME: Local Codex marketplace builder + the shared `--plugin-dir` install
// ABOUTME: helper that both `spacedock codex` and `spacedock install --host codex` call.
package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CodexMarketplaceInstall names the files WriteCodexLocalMarketplace laid down: the
// marketplace root to hand `codex plugin marketplace add`, the marketplace manifest
// it wrote, and the plugin symlink that points at the checkout.
type CodexMarketplaceInstall struct {
	MarketplaceRoot string
	ManifestPath    string
	PluginPath      string
}

// WriteCodexLocalMarketplace builds a local Codex marketplace under marketplaceRoot
// that exposes repoRoot as the `spacedock` plugin via a `plugins/spacedock` symlink
// and a `.agents/plugins/marketplace.json` (`source: local`). marketplaceName is the
// marketplace's own `name` — the channel marketplace the binary resolves
// (`spacedock` stable / `spacedock-edge` edge), so an edge binary's install resolves
// through the channel id rather than silently landing on the wrong channel. The
// plugin ENTRY stays `spacedock` (equal to the plugin's manifest name) on every
// channel; only the marketplace name carries the channel.
func WriteCodexLocalMarketplace(marketplaceRoot, repoRoot, marketplaceName string) (CodexMarketplaceInstall, error) {
	marketplaceRoot, err := filepath.Abs(marketplaceRoot)
	if err != nil {
		return CodexMarketplaceInstall{}, fmt.Errorf("resolve marketplace root: %w", err)
	}
	repoRoot, err = filepath.Abs(repoRoot)
	if err != nil {
		return CodexMarketplaceInstall{}, fmt.Errorf("resolve repo root: %w", err)
	}

	manifestDir := filepath.Join(marketplaceRoot, ".agents", "plugins")
	if err := os.MkdirAll(manifestDir, 0o755); err != nil {
		return CodexMarketplaceInstall{}, fmt.Errorf("create marketplace manifest dir: %w", err)
	}
	pluginsDir := filepath.Join(marketplaceRoot, "plugins")
	if err := os.MkdirAll(pluginsDir, 0o755); err != nil {
		return CodexMarketplaceInstall{}, fmt.Errorf("create marketplace plugins dir: %w", err)
	}

	pluginPath := filepath.Join(pluginsDir, "spacedock")
	if err := os.Remove(pluginPath); err != nil && !os.IsNotExist(err) {
		return CodexMarketplaceInstall{}, fmt.Errorf("replace existing plugin path: %w", err)
	}
	if err := os.Symlink(repoRoot, pluginPath); err != nil {
		return CodexMarketplaceInstall{}, fmt.Errorf("symlink current checkout into marketplace: %w", err)
	}

	manifest := struct {
		Name    string `json:"name"`
		Plugins []struct {
			Name   string `json:"name"`
			Source struct {
				Source string `json:"source"`
				Path   string `json:"path"`
			} `json:"source"`
			Policy struct {
				Installation   string `json:"installation"`
				Authentication string `json:"authentication"`
			} `json:"policy"`
			Category string `json:"category"`
		} `json:"plugins"`
	}{
		Name: marketplaceName,
		Plugins: []struct {
			Name   string `json:"name"`
			Source struct {
				Source string `json:"source"`
				Path   string `json:"path"`
			} `json:"source"`
			Policy struct {
				Installation   string `json:"installation"`
				Authentication string `json:"authentication"`
			} `json:"policy"`
			Category string `json:"category"`
		}{
			{
				Name:     "spacedock",
				Category: "workflow",
			},
		},
	}
	manifest.Plugins[0].Source.Source = "local"
	manifest.Plugins[0].Source.Path = "./plugins/spacedock"
	manifest.Plugins[0].Policy.Installation = "AVAILABLE"
	manifest.Plugins[0].Policy.Authentication = "ON_INSTALL"

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return CodexMarketplaceInstall{}, fmt.Errorf("encode marketplace manifest: %w", err)
	}
	data = append(data, '\n')

	manifestPath := filepath.Join(manifestDir, "marketplace.json")
	if err := os.WriteFile(manifestPath, data, 0o644); err != nil {
		return CodexMarketplaceInstall{}, fmt.Errorf("write marketplace manifest: %w", err)
	}
	return CodexMarketplaceInstall{
		MarketplaceRoot: marketplaceRoot,
		ManifestPath:    manifestPath,
		PluginPath:      pluginPath,
	}, nil
}

// installCodexLocalPluginDir builds a local marketplace from checkout (via
// WriteCodexLocalMarketplace, named for the binary's own channel) and installs it
// through the same Install() sequence a normal `spacedock codex` install uses, so
// the resulting plugin id and cache layout are indistinguishable from a marketplace
// install on the same channel. The marketplace root is persistent under the Codex
// home (one stable dir per channel, re-pointed on each call): codex records the
// marketplace source by path and re-loads it on every later `codex plugin` command,
// so a throwaway dir removed after install would hard-fail every subsequent codex
// invocation. It prints the version-masquerade advisory on every call: a
// `--plugin-dir` install reports the checkout's checked-in .codex-plugin/plugin.json
// version, not necessarily its current HEAD (the full stamping fix is deferred).
func installCodexLocalPluginDir(ops hostOps, checkout string, stderr io.Writer) error {
	marketplaceRoot := filepath.Join(codexHome(), "spacedock-plugin-dir", channelMarketplace(devBranch))
	if err := os.MkdirAll(marketplaceRoot, 0o755); err != nil {
		return fmt.Errorf("create local marketplace dir: %w", err)
	}

	install, err := WriteCodexLocalMarketplace(marketplaceRoot, checkout, channelMarketplace(devBranch))
	if err != nil {
		return fmt.Errorf("build local marketplace: %w", err)
	}
	if _, err := ops.Install("codex", install.MarketplaceRoot, devBranch); err != nil {
		return fmt.Errorf("install from local marketplace: %w", err)
	}
	fmt.Fprintf(stderr,
		"Installed codex plugin from %s.\n"+
			"version-masquerade advisory: the reported version reflects the checkout's "+
			"checked-in .codex-plugin/plugin.json, not necessarily its current HEAD.\n",
		checkout)
	return nil
}
