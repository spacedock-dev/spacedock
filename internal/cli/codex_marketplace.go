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

// codexPluginDirIncludes lists the top-level checkout entries WriteCodexLocalMarketplace
// stages into the plugin dir — the plugin surface the host manifests, skills, and
// hook config reference. Everything else in a dev checkout (`.git`, `.worktrees`,
// `tmp/`, build artifacts, …) is excluded. Staging the whole checkout via a
// root-symlinked plugin dir was the measured root cause of the codex plugin-dir
// launch tax: codex `plugin add` copies the symlink target wholesale into its
// cache, so a 3.6 GB checkout was copied on every launch against a <1 MB plugin
// surface (skills/ alone is ~800 KB).
var codexPluginDirIncludes = []string{
	".codex-plugin",
	".claude-plugin",
	"skills",
	"agents",
	"hooks",
	"hooks.json",
}

// CodexMarketplaceInstall names the files WriteCodexLocalMarketplace laid down: the
// marketplace root to hand `codex plugin marketplace add`, the marketplace manifest
// it wrote, and the plugin symlink that points at the checkout.
type CodexMarketplaceInstall struct {
	MarketplaceRoot string
	ManifestPath    string
	PluginPath      string
}

// WriteCodexLocalMarketplace builds a local Codex marketplace under marketplaceRoot
// that exposes repoRoot as the `spacedock` plugin via a `plugins/spacedock`
// FILTERED COPY (codexPluginDirIncludes — not a symlink to the whole checkout) and
// a `.agents/plugins/marketplace.json` (`source: local`). marketplaceName is the
// marketplace's own `name` — the channel marketplace the binary resolves
// (`spacedock` stable / `spacedock-edge` edge), so an edge binary's install resolves
// through the channel id rather than silently landing on the wrong channel. The
// plugin ENTRY stays `spacedock` (equal to the plugin's manifest name) on every
// channel; only the marketplace name carries the channel. Each call re-stages the
// plugin dir from scratch (cheap: <1 MB), so a later launch always serves the
// checkout's current content — freshness must flow through this staged copy
// because codex's `plugin add` reads it, not the checkout in place.
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
	if err := os.RemoveAll(pluginPath); err != nil {
		return CodexMarketplaceInstall{}, fmt.Errorf("replace existing plugin path: %w", err)
	}
	if err := os.MkdirAll(pluginPath, 0o755); err != nil {
		return CodexMarketplaceInstall{}, fmt.Errorf("create plugin staging dir: %w", err)
	}
	if err := stageCodexPluginDir(pluginPath, repoRoot); err != nil {
		return CodexMarketplaceInstall{}, fmt.Errorf("stage checkout into marketplace: %w", err)
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

// stageCodexPluginDir copies codexPluginDirIncludes from src into the already-empty
// dest directory. Only the include list is copied — never the whole checkout — so a
// dev checkout's `.git`, `.worktrees`, and other multi-GB build/state directories
// never reach the codex plugin cache. A missing include entry (e.g. a plugin with no
// `agents/` dir) is skipped, not an error.
func stageCodexPluginDir(dest, src string) error {
	for _, name := range codexPluginDirIncludes {
		srcPath := filepath.Join(src, name)
		info, err := os.Lstat(srcPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("stat %s: %w", srcPath, err)
		}
		if err := copyPath(filepath.Join(dest, name), srcPath, info); err != nil {
			return fmt.Errorf("copy %s: %w", name, err)
		}
	}
	return nil
}

// copyPath copies srcPath (already Lstat'd as info) to destPath, recursing into
// directories and copying regular file content+mode. Symlinks are resolved and
// copied as the file/dir they point to — the codex plugin dir must not carry
// symlinks back out to the checkout, or codex's copy-on-add would still pull in
// whatever the link targets.
func copyPath(destPath, srcPath string, info os.FileInfo) error {
	if info.Mode()&os.ModeSymlink != 0 {
		resolved, err := filepath.EvalSymlinks(srcPath)
		if err != nil {
			return fmt.Errorf("resolve symlink: %w", err)
		}
		resolvedInfo, err := os.Stat(resolved)
		if err != nil {
			return fmt.Errorf("stat symlink target: %w", err)
		}
		return copyPath(destPath, resolved, resolvedInfo)
	}
	if info.IsDir() {
		if err := os.MkdirAll(destPath, info.Mode().Perm()); err != nil {
			return err
		}
		entries, err := os.ReadDir(srcPath)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			entryInfo, err := entry.Info()
			if err != nil {
				return err
			}
			if err := copyPath(filepath.Join(destPath, entry.Name()), filepath.Join(srcPath, entry.Name()), entryInfo); err != nil {
				return err
			}
		}
		return nil
	}
	return copyFile(destPath, srcPath, info.Mode().Perm())
}

// copyFile copies one regular file's content, setting perm on the destination.
func copyFile(destPath, srcPath string, perm os.FileMode) error {
	in, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(destPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
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
	// Preflight the checkout before creating the persistent marketplace directory,
	// re-staging its plugin dir, or asking Codex to remove/install any channel.
	if err := validateLocalSpacedockPlugin(checkout, "codex"); err != nil {
		return err
	}
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
		"Installed codex plugin from %s as %s.\n"+
			"Removed other Spacedock Codex channels so $spacedock:* resolves from this install.\n"+
			"version-masquerade advisory: the reported version reflects the checkout's "+
			"checked-in .codex-plugin/plugin.json, not necessarily its current HEAD.\n",
		checkout, channelPluginID(devBranch))
	return nil
}
