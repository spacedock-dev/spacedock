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

	// Build a minimal checkout carrying both the plugin surface (staged) and
	// checkout-only cruft (excluded) — proves the copy is filtered, not a symlink
	// to the whole checkout.
	mustMkdirAll(t, filepath.Join(repo, ".codex-plugin"))
	mustWriteFile(t, filepath.Join(repo, ".codex-plugin", "plugin.json"), `{"name":"spacedock","version":"0.0.0","skills":"./skills/"}`)
	mustMkdirAll(t, filepath.Join(repo, "skills", "demo"))
	mustWriteFile(t, filepath.Join(repo, "skills", "demo", "SKILL.md"), "---\nname: demo\ndescription: demo\n---\ndemo\n")
	mustMkdirAll(t, filepath.Join(repo, ".git"))
	mustWriteFile(t, filepath.Join(repo, ".git", "HEAD"), "ref: refs/heads/main\n")
	mustMkdirAll(t, filepath.Join(repo, ".worktrees", "other"))
	mustWriteFile(t, filepath.Join(repo, ".worktrees", "other", "marker"), "not plugin surface\n")

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
	info, err := os.Lstat(install.PluginPath)
	if err != nil {
		t.Fatalf("stat plugins/spacedock: %v", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Fatalf("plugins/spacedock must be a real staged directory, not a symlink to the checkout")
	}
	if !info.IsDir() {
		t.Fatalf("plugins/spacedock must be a directory")
	}
	staged, err := os.ReadFile(filepath.Join(install.PluginPath, ".codex-plugin", "plugin.json"))
	if err != nil {
		t.Fatalf("staged plugin manifest missing: %v", err)
	}
	if !strings.Contains(string(staged), `"name":"spacedock"`) {
		t.Fatalf("staged plugin manifest content wrong:\n%s", staged)
	}
	if _, err := os.Stat(filepath.Join(install.PluginPath, "skills", "demo", "SKILL.md")); err != nil {
		t.Fatalf("staged plugin dir missing skills/demo/SKILL.md: %v", err)
	}
	if _, err := os.Stat(filepath.Join(install.PluginPath, ".git")); !os.IsNotExist(err) {
		t.Fatalf("staged plugin dir must exclude .git; stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(install.PluginPath, ".worktrees")); !os.IsNotExist(err) {
		t.Fatalf("staged plugin dir must exclude .worktrees; stat err = %v", err)
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

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
