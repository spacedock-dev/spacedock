// ABOUTME: Host-specific validation and adjacent discovery for local Spacedock plugins.
// ABOUTME: A launcher selects its own checkout only when the host manifest and skills directory agree.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type localPluginManifest struct {
	Name   string `json:"name"`
	Skills string `json:"skills"`
}

func localPluginManifestPath(root, host string) (string, error) {
	switch host {
	case "claude":
		return filepath.Join(root, ".claude-plugin", "plugin.json"), nil
	case "codex":
		return filepath.Join(root, ".codex-plugin", "plugin.json"), nil
	default:
		return "", fmt.Errorf("unsupported plugin host %q", host)
	}
}

// validateLocalSpacedockPlugin qualifies root as a plugin for the named host.
// The host manifest must parse, identify Spacedock, configure a skills path, and
// point that path at an existing directory. Callers use the same validator for
// automatic discovery and explicit Codex installs so validation always precedes
// persistent marketplace mutation.
func validateLocalSpacedockPlugin(root, host string) error {
	manifestPath, err := localPluginManifestPath(root, host)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("invalid local Spacedock plugin %s: read manifest: %w", manifestPath, err)
	}
	var manifest localPluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return fmt.Errorf("invalid local Spacedock plugin %s: parse manifest: %w", manifestPath, err)
	}
	if manifest.Name != "spacedock" {
		return fmt.Errorf("invalid local Spacedock plugin %s: name is %q, want %q", manifestPath, manifest.Name, "spacedock")
	}
	if manifest.Skills == "" {
		return fmt.Errorf("invalid local Spacedock plugin %s: skills path is empty", manifestPath)
	}
	skillsPath := filepath.FromSlash(manifest.Skills)
	if !filepath.IsAbs(skillsPath) {
		skillsPath = filepath.Join(root, skillsPath)
	}
	info, err := os.Stat(skillsPath)
	if err != nil {
		return fmt.Errorf("invalid local Spacedock plugin %s: skills directory %s: %w", manifestPath, skillsPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("invalid local Spacedock plugin %s: skills path %s is not a directory", manifestPath, skillsPath)
	}
	return nil
}

// adjacentSpacedockPluginRoot returns the directory containing the resolved
// launcher only when it is a valid host-specific Spacedock plugin checkout.
// Invalid or release-style layouts are simply not adjacent providers; callers
// retain the installed resolver/gate behavior.
func adjacentSpacedockPluginRoot(host string) (string, bool) {
	bin, ok := resolvedLauncherBin()
	if !ok {
		return "", false
	}
	root := filepath.Dir(bin)
	if err := validateLocalSpacedockPlugin(root, host); err != nil {
		return "", false
	}
	return root, true
}
