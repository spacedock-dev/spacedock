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
	// devBranch=main installs the stable `spacedock` channel and resolves it under
	// cache/spacedock/spacedock — the resolver reads the package devBranch, so it
	// must match the channel this test installs.
	saved := devBranch
	devBranch = "main"
	defer func() { devBranch = saved }()

	tmp := t.TempDir()
	marketplace := buildLocalCodexMarketplace(t, tmp)
	codexHomeDir := filepath.Join(tmp, "codexhome")
	mustMkdir(t, codexHomeDir)
	t.Setenv("CODEX_HOME", codexHomeDir)

	// The real Install runs in-process through the production seam: the codex arm
	// shells the 4-step sequence against the isolated CODEX_HOME (CODEX_HOME is read
	// from the env by the codex CLI). devBranch=main selects the stable `spacedock`
	// entry the fixture marketplace defines; the local-path source keeps it offline.
	out, err := execHost{}.Install("codex", marketplace, "main")
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

// TestCodexInitRefreshAdvancesBehindPlugin is the AC-2 live proof for the wired
// codex arm: seed an older plugin (0.0.1) via the production Install seam, then
// run that same seam again — against the SAME marketplace path, now rewritten to
// 0.0.2 — the path runInit's codex arm now drives on a present plugin. The
// marketplace path must stay the same source string across both calls: real
// codex measurably refuses `plugin marketplace add` when a marketplace name is
// already registered from a DIFFERENT source (probe 3), so this mirrors how a
// real git-tracked marketplace source behaves (a stable source URL; new content
// arrives via a git pull at that same URL, never a new URL) rather than exercising
// the source-migration case this task explicitly does not cover. The resolved
// cache manifest must advance to 0.0.2, read from the on-disk manifest (not the
// install command's claim). This pins the spike's finding as a regression test
// on the production refresh path. Skips when `codex` is not on PATH; hermetic
// via CODEX_HOME isolation + a local-path marketplace (empty branch → no --ref →
// offline).
func TestCodexInitRefreshAdvancesBehindPlugin(t *testing.T) {
	if _, err := exec.LookPath("codex"); err != nil {
		t.Skip("codex not on PATH; refresh-advances smoke requires the host CLI")
	}
	// devBranch=main installs and resolves the stable `spacedock` channel; the
	// resolver reads the package devBranch, so it must match the seeded channel.
	saved := devBranch
	devBranch = "main"
	defer func() { devBranch = saved }()

	tmp := t.TempDir()
	marketplaceRoot := filepath.Join(tmp, "marketplace-root")
	codexHomeDir := filepath.Join(tmp, "codexhome")
	mustMkdir(t, codexHomeDir)
	t.Setenv("CODEX_HOME", codexHomeDir)

	// Seed the behind install (0.0.1) through the production seam. devBranch=main
	// selects the stable `spacedock` entry the fixture defines.
	marketplace := buildCodexMarketplaceAtVersion(t, marketplaceRoot, "0.0.1")
	if out, err := (execHost{}).Install("codex", marketplace, "main"); err != nil {
		t.Fatalf("seed Install(codex, 0.0.1) failed: %v\nout=%q", err, out)
	}
	if got := resolvedCodexManifestVersion(t); got != "0.0.1" {
		t.Fatalf("after seed, resolved manifest version = %q, want 0.0.1", got)
	}

	// Refresh-on-present (0.0.2) — the wired runInit codex arm calls exactly this
	// Install seam when a plugin is already resolved. Same marketplaceRoot, content
	// rewritten in place: the source string `marketplace` is byte-identical to the
	// seed call.
	marketplace = buildCodexMarketplaceAtVersion(t, marketplaceRoot, "0.0.2")
	if out, err := (execHost{}).Install("codex", marketplace, "main"); err != nil {
		t.Fatalf("refresh Install(codex, 0.0.2) failed: %v\nout=%q", err, out)
	}
	if got := resolvedCodexManifestVersion(t); got != "0.0.2" {
		t.Fatalf("after refresh, resolved manifest version = %q, want 0.0.2 (refresh did not advance the behind plugin)", got)
	}
}

// resolvedCodexManifestVersion resolves the cached spacedock@spacedock manifest
// via the production resolver and returns its on-disk version field — the
// independent source of truth that an install advanced the plugin, not the
// command's stdout.
func resolvedCodexManifestVersion(t *testing.T) string {
	t.Helper()
	manifest, err := execHost{}.resolveCodexManifest()
	if err != nil {
		t.Fatalf("resolveCodexManifest: %v", err)
	}
	if manifest == "" {
		t.Fatalf("resolveCodexManifest returned empty after an install")
	}
	data, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatalf("read resolved manifest %s: %v", manifest, err)
	}
	var m struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse resolved manifest %s: %v", manifest, err)
	}
	return m.Version
}

// buildLocalCodexMarketplace writes a minimal valid local-path marketplace under
// root and returns the marketplace directory. Codex reads the marketplace
// manifest from .claude-plugin/marketplace.json (it reuses the claude manifest
// layout) and the plugin manifest from the plugin's .codex-plugin/plugin.json.
// The plugin manifest carries the display version the doctor verdict reads.
func buildLocalCodexMarketplace(t *testing.T, root string) string {
	return buildCodexMarketplaceAtVersion(t, root, "0.0.0")
}

// buildCodexMarketplaceAtVersion is buildLocalCodexMarketplace parameterized by
// the plugin's display version, so a test can seed a behind plugin then refresh
// it from a newer marketplace and observe the resolved version advance.
func buildCodexMarketplaceAtVersion(t *testing.T, root, version string) string {
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
	mustWrite(t, filepath.Join(plugin, ".codex-plugin", "plugin.json"),
		`{ "name": "spacedock", "version": "`+version+`", "skills": "./skills/" }
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
