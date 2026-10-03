// ABOUTME: RUN-verified codex resolver test — execHost.ResolveManifest("codex")
// ABOUTME: against a real codex CLI in an isolated CODEX_HOME seeded by a local
// ABOUTME: marketplace install, never the operator's own plugin cache.
package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCodexResolveManifestAgainstInstalledHost drives the production codex
// resolver against the real `codex` CLI, but in a test-owned isolated CODEX_HOME:
// it installs `spacedock@spacedock` from a local-path marketplace and then
// requires ResolveManifest to return that isolated install's
// .codex-plugin/plugin.json. It never reads the operator's ~/.codex cache, so an
// ambient `spacedock@spacedock-local` install (which the resolver checks FIRST)
// can no longer make the resolver disagree with the test's own install check.
// codex 0.136.0 rejects `--json` (exit 2), so the resolver must use the supported
// `codex plugin list` text output. Skips when codex is absent.
func TestCodexResolveManifestAgainstInstalledHost(t *testing.T) {
	codexBin, err := exec.LookPath("codex")
	if err != nil {
		t.Skip("codex not on PATH; codex resolver test requires the host CLI")
	}
	// Verify the stable-channel resolver against a stable (`spacedock@spacedock`)
	// install: the resolver reads the package devBranch to pick the channel id, so
	// pin it to main to match the id this test installs.
	saved := devBranch
	devBranch = "main"
	defer func() { devBranch = saved }()

	tmp := t.TempDir()
	codexHomeDir := filepath.Join(tmp, "codexhome")
	mustMkdir(t, codexHomeDir)
	// Isolate the host home: both the resolver (codexHome) and the host CLI read
	// CODEX_HOME, so overriding the operator's ambient value is what makes this
	// deterministic on an operator machine that already has a codex plugin cache.
	t.Setenv("CODEX_HOME", codexHomeDir)

	// No install yet: the resolver must degrade to "" with no error.
	if path, err := (execHost{}).ResolveManifest("codex"); err != nil || path != "" {
		t.Fatalf("ResolveManifest on an empty isolated CODEX_HOME = (%q, %v), want empty no-error", path, err)
	}

	marketplace := buildLocalCodexMarketplace(t, filepath.Join(tmp, "marketplace-root"))
	if out, err := (execHost{}).Install("codex", marketplace, "main"); err != nil {
		t.Fatalf("codex install from the local marketplace failed: %v\n%s", err, out)
	}

	// Confirm the isolated host home actually registered the channel id before
	// requiring the resolver to find it.
	listOut := runHost(t, codexBin, os.Environ(), "plugin", "list")
	if !codexEntryInstalled(listOut, "spacedock@spacedock") {
		t.Fatalf("isolated CODEX_HOME did not register spacedock@spacedock:\n%s", listOut)
	}

	path, err := (execHost{}).ResolveManifest("codex")
	if err != nil {
		t.Fatalf("ResolveManifest(codex) errored: %v", err)
	}
	if path == "" {
		t.Fatal("spacedock@spacedock installed in the isolated CODEX_HOME, but resolver returned empty path")
	}
	if filepath.Base(path) != "plugin.json" || !strings.Contains(path, ".codex-plugin") {
		t.Fatalf("resolved codex manifest path is not a .codex-plugin/plugin.json: %q", path)
	}
	if !strings.HasPrefix(path, codexHomeDir) {
		t.Fatalf("resolved manifest %q is not under the isolated CODEX_HOME %q", path, codexHomeDir)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("resolved codex manifest does not exist: %q (%v)", path, err)
	}
}
