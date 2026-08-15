// ABOUTME: AC-7 backstop — exercises codex's entry-name vs plugin.json name-match
// ABOUTME: constraint on the real CLI: old entry-name shape fails, new shape passes.
package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCodexEntryNameMustMatchPluginName is the AC-7 mechanism backstop for the
// integration gap that shipped the edge channel without a real `codex plugin add`.
// Codex enforces that a marketplace entry's name equals the plugin's own
// plugin.json `name` (`spacedock`). This test runs the REAL codex CLI against an
// isolated CODEX_HOME and two local marketplaces:
//
//   - OLD shape: an entry named `spacedock-edge` in a marketplace named `spacedock`
//     — `codex plugin add spacedock-edge@spacedock` MUST FAIL with the name-mismatch
//     error (the failing-today baseline the channel fix flips).
//   - NEW shape: the entry stays `spacedock`, the marketplace is named
//     `spacedock-edge` — `codex plugin add spacedock@spacedock-edge` MUST SUCCEED,
//     and the install lands at the channel cache path the resolver expects
//     (cache/spacedock-edge/spacedock/<ver>/.codex-plugin/plugin.json).
//
// Skips when `codex` is not on PATH; hermetic via CODEX_HOME isolation +
// local-path marketplaces (offline). This is the unit-level proof that serves AC-1;
// it does not substitute for AC-1's live-lane `codex plugin add`.
func TestCodexEntryNameMustMatchPluginName(t *testing.T) {
	codexBin, err := exec.LookPath("codex")
	if err != nil {
		t.Skip("codex not on PATH; name-match backstop requires the host CLI")
	}

	tmp := t.TempDir()
	codexHomeDir := filepath.Join(tmp, "codexhome")
	mustMkdir(t, codexHomeDir)
	t.Setenv("CODEX_HOME", codexHomeDir)

	// OLD shape: marketplace `spacedock`, entry `spacedock-edge` (!= plugin.json
	// name `spacedock`) — codex must reject the add.
	oldMkt := buildCodexNamedMarketplace(t, filepath.Join(tmp, "oldshape"), "spacedock", "spacedock-edge")
	runHost(t, codexBin, os.Environ(), "plugin", "marketplace", "add", oldMkt)
	out, err := exec.Command(codexBin, "plugin", "add", "spacedock-edge@spacedock").CombinedOutput()
	if err == nil {
		t.Fatalf("BASELINE REGRESSION: codex accepted the old entry-name shape spacedock-edge@spacedock; the name-match constraint this AC depends on no longer holds\n%s", out)
	}
	if !strings.Contains(string(out), "does not match") {
		t.Fatalf("old-shape add failed but not with the expected name-mismatch error:\n%s", out)
	}

	// NEW shape: marketplace `spacedock-edge`, entry `spacedock` (== plugin.json
	// name) — codex must accept the add.
	edgeMkt := buildCodexNamedMarketplace(t, filepath.Join(tmp, "edge"), "spacedock-edge", "spacedock")
	runHost(t, codexBin, os.Environ(), "plugin", "marketplace", "add", edgeMkt)
	if out, err := exec.Command(codexBin, "plugin", "add", "spacedock@spacedock-edge").CombinedOutput(); err != nil {
		t.Fatalf("NEW shape add of spacedock@spacedock-edge failed: %v\n%s", err, out)
	}

	// The install must land where the channel-aware resolver looks for the edge
	// channel: cache/spacedock-edge/spacedock/<ver>/.codex-plugin/plugin.json.
	saved := devBranch
	devBranch = "next"
	defer func() { devBranch = saved }()
	manifest, err := codexCacheManifest()
	if err != nil {
		t.Fatalf("codexCacheManifest (edge) after real install: %v", err)
	}
	if manifest == "" {
		t.Fatalf("codexCacheManifest (edge) returned empty after a real spacedock@spacedock-edge install")
	}
	if !strings.Contains(manifest, filepath.Join("cache", "spacedock-edge", "spacedock")) {
		t.Fatalf("edge install resolved to %q, want a path under cache/spacedock-edge/spacedock", manifest)
	}
	if _, statErr := os.Stat(manifest); statErr != nil {
		t.Fatalf("resolved edge manifest %s does not exist: %v", manifest, statErr)
	}
}

// buildCodexNamedMarketplace writes a minimal local marketplace under root with the
// given marketplace name and single-entry name, sourcing the plugin from a colocated
// checkout whose plugin.json name is always `spacedock`. Returns the marketplace dir.
// Codex reads the marketplace manifest from .claude-plugin/marketplace.json and the
// plugin manifest from the plugin's .codex-plugin/plugin.json.
func buildCodexNamedMarketplace(t *testing.T, root, marketplaceName, entryName string) string {
	t.Helper()
	marketplace := filepath.Join(root, "marketplace")
	plugin := filepath.Join(marketplace, "spacedock")
	mustMkdir(t, filepath.Join(marketplace, ".claude-plugin"))
	mustMkdir(t, filepath.Join(plugin, ".codex-plugin"))
	mustMkdir(t, filepath.Join(plugin, "skills", "demo"))

	mustWrite(t, filepath.Join(marketplace, ".claude-plugin", "marketplace.json"), `{
  "name": "`+marketplaceName+`",
  "owner": { "name": "CL Kao" },
  "plugins": [
    { "name": "`+entryName+`", "source": "./spacedock", "description": "test", "category": "workflow" }
  ]
}
`)
	mustWrite(t, filepath.Join(plugin, ".codex-plugin", "plugin.json"),
		`{ "name": "spacedock", "version": "0.0.0", "skills": "./skills/" }
`)
	mustWrite(t, filepath.Join(plugin, "skills", "demo", "SKILL.md"), "---\nname: demo\ndescription: demo skill\n---\ndemo\n")
	return marketplace
}
