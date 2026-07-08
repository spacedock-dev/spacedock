package ensigncycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/cli"
)

func TestWriteCodexLocalMarketplacePointsAtCurrentCheckout(t *testing.T) {
	repo := t.TempDir()
	marketplace := t.TempDir()

	install, err := cli.WriteCodexLocalMarketplace(marketplace, repo, "spacedock")
	if err != nil {
		t.Fatalf("WriteCodexLocalMarketplace errored: %v", err)
	}

	if install.MarketplaceRoot != marketplace {
		t.Fatalf("MarketplaceRoot = %q, want %q", install.MarketplaceRoot, marketplace)
	}
	if install.PluginPath != filepath.Join(marketplace, "plugins", "spacedock") {
		t.Fatalf("PluginPath = %q, want plugins/spacedock under marketplace", install.PluginPath)
	}
	target, err := os.Readlink(install.PluginPath)
	if err != nil {
		t.Fatalf("plugins/spacedock must be a symlink to the current checkout: %v", err)
	}
	if target != repo {
		t.Fatalf("plugins/spacedock symlink target = %q, want %q", target, repo)
	}

	data, err := os.ReadFile(install.ManifestPath)
	if err != nil {
		t.Fatalf("read marketplace manifest: %v", err)
	}
	if strings.Contains(string(data), "github.com") || strings.Contains(string(data), "next") {
		t.Fatalf("local marketplace manifest must not name remote next:\n%s", data)
	}

	var manifest struct {
		Name    string `json:"name"`
		Plugins []struct {
			Name   string `json:"name"`
			Source struct {
				Source string `json:"source"`
				Path   string `json:"path"`
			} `json:"source"`
		} `json:"plugins"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parse marketplace manifest: %v\n%s", err, data)
	}
	if manifest.Name != "spacedock" {
		t.Fatalf("manifest name = %q, want spacedock", manifest.Name)
	}
	if len(manifest.Plugins) != 1 {
		t.Fatalf("plugins count = %d, want 1", len(manifest.Plugins))
	}
	plugin := manifest.Plugins[0]
	if plugin.Name != "spacedock" {
		t.Fatalf("plugin name = %q, want spacedock", plugin.Name)
	}
	if plugin.Source.Source != "local" {
		t.Fatalf("plugin source = %q, want local", plugin.Source.Source)
	}
	if plugin.Source.Path != "./plugins/spacedock" {
		t.Fatalf("plugin path = %q, want ./plugins/spacedock", plugin.Source.Path)
	}
}
