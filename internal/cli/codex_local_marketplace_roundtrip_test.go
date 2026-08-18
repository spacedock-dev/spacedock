// ABOUTME: Round-3 live proof — the codex mode-switch round trip (channel
// ABOUTME: install -> --plugin-dir dev install -> channel install) stays exclusive.
package cli

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCodexModeSwitchRoundTripPreservesExclusivity is the user journey the
// spacedock-local marketplace-name fix (captain-ruled option (c)) exists to
// protect: a captain runs a normal channel install, switches to `--plugin-dir`
// for local dev, then switches back. At every step, real codex must report
// EXACTLY ONE installed `spacedock@*` provider (codex's skill namespace is
// global, so two resolves $spacedock:* ambiguously) and
// execHost.ResolveManifest("codex") must resolve it — proving both halves of
// the round-3 fix hold across the switch, not just in one direction:
// codexPluginDirInstallArgvSequence removing the channel ids before adding the
// local one (channel -> local), and codexInstallArgvSequence removing the local
// id before adding the channel one (local -> channel, the "teach the gate the
// second id" + round-trip-exclusivity asks). Skips when codex is not on PATH;
// hermetic via CODEX_HOME isolation and local-path marketplace sources (no
// network).
func TestCodexModeSwitchRoundTripPreservesExclusivity(t *testing.T) {
	codexBin, err := exec.LookPath("codex")
	if err != nil {
		t.Skip("codex not on PATH; round-trip smoke requires the host CLI")
	}
	saved := devBranch
	devBranch = "main"
	defer func() { devBranch = saved }()

	tmp := t.TempDir()
	codexHomeDir := filepath.Join(tmp, "codexhome")
	mustMkdir(t, codexHomeDir)
	t.Setenv("CODEX_HOME", codexHomeDir)

	channelMarketplace := buildLocalCodexMarketplace(t, filepath.Join(tmp, "channel"))
	localCheckout := buildCodexPluginCheckout(t, filepath.Join(tmp, "local-checkout"), "9.9.9")

	// Step 1: a normal channel install (stable, devBranch=main).
	if out, err := (execHost{}).Install("codex", channelMarketplace, "main"); err != nil {
		t.Fatalf("channel install failed: %v\nout=%q", err, out)
	}
	assertExactlyOneCodexProvider(t, codexBin, "spacedock@spacedock")
	if manifest, err := (execHost{}).ResolveManifest("codex"); err != nil || manifest == "" {
		t.Fatalf("ResolveManifest after channel install: manifest=%q err=%v", manifest, err)
	}

	// Step 2: switch to --plugin-dir for local dev.
	if err := installCodexLocalPluginDir(execHost{}, localCheckout, io.Discard); err != nil {
		t.Fatalf("plugin-dir install failed: %v", err)
	}
	assertExactlyOneCodexProvider(t, codexBin, codexLocalPluginID)
	manifest, err := (execHost{}).ResolveManifest("codex")
	if err != nil {
		t.Fatalf("ResolveManifest after plugin-dir install: %v", err)
	}
	if manifest == "" {
		t.Fatal("ResolveManifest returned empty after a --plugin-dir install — the gate does not see the local id (round-3 point 2 regressed)")
	}
	if _, statErr := os.Stat(manifest); statErr != nil {
		t.Fatalf("resolved local manifest %s does not exist: %v", manifest, statErr)
	}

	// Step 3: switch back to the normal channel install.
	if out, err := (execHost{}).Install("codex", channelMarketplace, "main"); err != nil {
		t.Fatalf("channel re-install failed: %v\nout=%q", err, out)
	}
	assertExactlyOneCodexProvider(t, codexBin, "spacedock@spacedock")
	if manifest, err := (execHost{}).ResolveManifest("codex"); err != nil || manifest == "" {
		t.Fatalf("ResolveManifest after channel re-install: manifest=%q err=%v", manifest, err)
	}
}

// assertExactlyOneCodexProvider asserts `codex plugin list --json` reports
// EXACTLY one installed `spacedock@*` provider, and it is wantID — the
// plugin-level exclusivity property codex's global skill namespace requires.
func assertExactlyOneCodexProvider(t *testing.T, codexBin, wantID string) {
	t.Helper()
	listOut := runHost(t, codexBin, os.Environ(), "plugin", "list", "--json")
	var listing struct {
		Installed []struct {
			PluginID  string `json:"pluginId"`
			Installed bool   `json:"installed"`
		} `json:"installed"`
	}
	if err := json.Unmarshal([]byte(listOut), &listing); err != nil {
		t.Fatalf("parse codex plugin list --json: %v\n%s", err, listOut)
	}
	var installed []string
	for _, e := range listing.Installed {
		if e.Installed && strings.HasPrefix(e.PluginID, "spacedock@") {
			installed = append(installed, e.PluginID)
		}
	}
	if len(installed) != 1 || installed[0] != wantID {
		t.Fatalf("installed spacedock@* providers = %v, want exactly [%q]", installed, wantID)
	}
}
