// ABOUTME: AC-4 channel-aware resolvers — an edge binary recognizes/refreshes an
// ABOUTME: installed edge plugin (spacedock@spacedock-edge), not the hardcoded stable id.
package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestCodexCacheManifestResolvesEdgeChannel locks AC-4 for the codex cache tail:
// an edge binary (devBranch=next) resolves the EDGE plugin's cached manifest under
// the channel marketplace's cache dir (cache/spacedock-edge/spacedock/<ver>/…), not
// the hardcoded stable cache/spacedock/spacedock path. A stable-only resolver would
// return "" here and the edge front door would never launch a present edge plugin.
func TestCodexCacheManifestResolvesEdgeChannel(t *testing.T) {
	saved := devBranch
	devBranch = "next"
	defer func() { devBranch = saved }()

	home := t.TempDir()
	t.Setenv("CODEX_HOME", home)
	// The edge plugin's cache layout: the marketplace name is the channel
	// (spacedock-edge), the entry is spacedock.
	manifest := filepath.Join(home, "plugins", "cache", "spacedock-edge", "spacedock", "0.13.0", ".codex-plugin", "plugin.json")
	if err := os.MkdirAll(filepath.Dir(manifest), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifest, []byte(`{"name":"spacedock"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := codexCacheManifest()
	if err != nil {
		t.Fatalf("codexCacheManifest errored: %v", err)
	}
	if got != manifest {
		t.Fatalf("codexCacheManifest (edge) = %q, want %q (edge resolves under the spacedock-edge marketplace cache)", got, manifest)
	}
}

// TestCodexCacheManifestStableUnchanged guards that the stable channel still
// resolves under cache/spacedock/spacedock — the channel-aware change must not
// regress the stable path.
func TestCodexCacheManifestStableUnchanged(t *testing.T) {
	saved := devBranch
	devBranch = "main"
	defer func() { devBranch = saved }()

	home := t.TempDir()
	t.Setenv("CODEX_HOME", home)
	manifest := filepath.Join(home, "plugins", "cache", "spacedock", "spacedock", "0.13.0", ".codex-plugin", "plugin.json")
	if err := os.MkdirAll(filepath.Dir(manifest), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifest, []byte(`{"name":"spacedock"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := codexCacheManifest()
	if err != nil {
		t.Fatalf("codexCacheManifest errored: %v", err)
	}
	if got != manifest {
		t.Fatalf("codexCacheManifest (stable) = %q, want %q", got, manifest)
	}
}

// TestResolveClaudeManifestMatchesEdgeChannel locks AC-4 for the claude resolver:
// with an edge binary (devBranch=next), resolveClaudeManifest must match the EDGE
// plugin id (spacedock@spacedock-edge) in the host listing, not only the stable
// spacedock@spacedock. A PATH-stub `claude` emits the listing JSON so the test runs
// with no real claude CLI.
func TestResolveClaudeManifestMatchesEdgeChannel(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	saved := devBranch
	devBranch = "next"
	defer func() { devBranch = saved }()

	installPath := t.TempDir()
	listing := `[{"id":"spacedock@spacedock-edge","installPath":"` + installPath + `","enabled":true}]`
	dir := writePluginListStub(t, "claude", listing)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	got, err := execHost{}.resolveClaudeManifest("claude")
	if err != nil {
		t.Fatalf("resolveClaudeManifest errored: %v", err)
	}
	want := filepath.Join(installPath, manifestSubpath("claude"))
	if got != want {
		t.Fatalf("resolveClaudeManifest (edge) = %q, want %q (edge binary must recognize the spacedock@spacedock-edge install)", got, want)
	}
}

// writePluginListStub writes a host CLI stub that prints the given JSON listing on
// `plugin list --json` (and empty otherwise), for resolver tests that need the
// host's `plugin list` output without a real CLI.
func writePluginListStub(t *testing.T, binName, listingJSON string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, binName)
	body := `#!/bin/sh
case "$*" in
"plugin list --json")
  cat <<'JSON'
` + listingJSON + `
JSON
  ;;
esac
`
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatalf("write plugin-list stub: %v", err)
	}
	return dir
}
