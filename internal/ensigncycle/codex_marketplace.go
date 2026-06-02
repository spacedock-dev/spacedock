package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type codexMarketplaceInstall struct {
	marketplaceRoot string
	manifestPath    string
	pluginPath      string
}

func writeCodexLocalMarketplace(marketplaceRoot, repoRoot string) (codexMarketplaceInstall, error) {
	marketplaceRoot, err := filepath.Abs(marketplaceRoot)
	if err != nil {
		return codexMarketplaceInstall{}, fmt.Errorf("resolve marketplace root: %w", err)
	}
	repoRoot, err = filepath.Abs(repoRoot)
	if err != nil {
		return codexMarketplaceInstall{}, fmt.Errorf("resolve repo root: %w", err)
	}

	manifestDir := filepath.Join(marketplaceRoot, ".agents", "plugins")
	if err := os.MkdirAll(manifestDir, 0o755); err != nil {
		return codexMarketplaceInstall{}, fmt.Errorf("create marketplace manifest dir: %w", err)
	}
	pluginsDir := filepath.Join(marketplaceRoot, "plugins")
	if err := os.MkdirAll(pluginsDir, 0o755); err != nil {
		return codexMarketplaceInstall{}, fmt.Errorf("create marketplace plugins dir: %w", err)
	}

	pluginPath := filepath.Join(pluginsDir, "spacedock")
	if err := os.Remove(pluginPath); err != nil && !os.IsNotExist(err) {
		return codexMarketplaceInstall{}, fmt.Errorf("replace existing plugin path: %w", err)
	}
	if err := os.Symlink(repoRoot, pluginPath); err != nil {
		return codexMarketplaceInstall{}, fmt.Errorf("symlink current checkout into marketplace: %w", err)
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
		Name: "spacedock",
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
		return codexMarketplaceInstall{}, fmt.Errorf("encode marketplace manifest: %w", err)
	}
	data = append(data, '\n')

	manifestPath := filepath.Join(manifestDir, "marketplace.json")
	if err := os.WriteFile(manifestPath, data, 0o644); err != nil {
		return codexMarketplaceInstall{}, fmt.Errorf("write marketplace manifest: %w", err)
	}
	return codexMarketplaceInstall{
		marketplaceRoot: marketplaceRoot,
		manifestPath:    manifestPath,
		pluginPath:      pluginPath,
	}, nil
}
