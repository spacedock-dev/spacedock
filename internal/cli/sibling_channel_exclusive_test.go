// ABOUTME: AC-1 live proof — a claude channel install against a two-channel local
// ABOUTME: marketplace must leave exactly one enabled spacedock plugin: the selected one.
package cli

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestClaudeInstallLeavesOnlySelectedChannelEnabled locks AC-1 against the real
// claude CLI: after execHost.Install, the COUNT of enabled `spacedock@*` entries in
// `claude plugin list --json` — claude's own reporting, not a Spacedock value — is
// exactly 1, and it is the selected channel's id. The baseline moves the wrong way:
// delete the sibling-uninstall step from installArgvSequence and the sibling-present
// and reverse subtests count 2, which is the captain's reported bug reproduced
// hermetically. Three directions: a co-installed sibling (the reported shape), a
// fresh box (proving the absent sibling's tolerated non-zero exit does not abort the
// run), and edge-over-stable (proving the mapping works both ways rather than
// hard-coding one channel). Skips when `claude` is not on PATH — CI's build job has
// no host CLI, so validation MUST run this locally and treat a SKIP as a failure.
func TestClaudeInstallLeavesOnlySelectedChannelEnabled(t *testing.T) {
	claudeBin, err := exec.LookPath("claude")
	if err != nil {
		t.Skip("claude not on PATH; channel-exclusivity test requires the host CLI")
	}

	cases := []struct {
		name      string
		devBranch string
		seedID    string // sibling installed before the run; empty means fresh box
		wantID    string
	}{
		{name: "sibling-present", devBranch: "main", seedID: "spacedock@spacedock-edge", wantID: "spacedock@spacedock"},
		{name: "fresh-box", devBranch: "main", wantID: "spacedock@spacedock"},
		{name: "reverse", devBranch: "next", seedID: "spacedock@spacedock", wantID: "spacedock@spacedock-edge"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()
			stable := buildNamedChannelMarketplace(t, tmp, "spacedock")
			edge := buildNamedChannelMarketplace(t, tmp, "spacedock-edge")
			configDir := filepath.Join(tmp, "config")
			cacheDir := filepath.Join(tmp, "cache")
			mustMkdir(t, configDir)
			mustMkdir(t, cacheDir)
			// Isolation must land on the process environment: execHost.Install shells
			// out with no explicit Env and inherits it (co_hosted_survival_test.go).
			t.Setenv("CLAUDE_CONFIG_DIR", configDir)
			t.Setenv("CLAUDE_CODE_PLUGIN_CACHE_DIR", cacheDir)
			env := os.Environ()

			source := stable
			if tc.devBranch != "main" {
				source = edge
			}
			if tc.seedID != "" {
				seedSource := edge
				if tc.seedID == "spacedock@spacedock" {
					seedSource = stable
				}
				runHost(t, claudeBin, env, "plugin", "marketplace", "add", seedSource)
				runHost(t, claudeBin, env, "plugin", "install", tc.seedID)
				if got := enabledSpacedockIDs(t, claudeBin, env); len(got) != 1 || got[0] != tc.seedID {
					t.Fatalf("seed left enabled %v, want exactly [%s] — the fixture, not the install, is wrong", got, tc.seedID)
				}
			}

			out, err := execHost{}.Install("claude", source, tc.devBranch)
			if err != nil {
				t.Fatalf("channel Install failed: %v\nout=%q", err, out)
			}

			got := enabledSpacedockIDs(t, claudeBin, env)
			if len(got) != 1 || got[0] != tc.wantID {
				t.Fatalf("enabled spacedock plugins after install = %v, want exactly [%s]\ninstall output:\n%s", got, tc.wantID, out)
			}
		})
	}
}

// enabledSpacedockIDs returns the ids of every enabled `spacedock@*` entry claude
// reports — the AC-1 quantity, read from the host's own JSON.
func enabledSpacedockIDs(t *testing.T, claudeBin string, env []string) []string {
	t.Helper()
	listOut := runHost(t, claudeBin, env, "plugin", "list", "--json")
	var entries []struct {
		ID      string `json:"id"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.Unmarshal([]byte(listOut), &entries); err != nil {
		t.Fatalf("parse plugin list --json: %v\n%s", err, listOut)
	}
	var ids []string
	for _, e := range entries {
		if e.Enabled && strings.HasPrefix(e.ID, "spacedock@") {
			ids = append(ids, e.ID)
		}
	}
	return ids
}

// buildNamedChannelMarketplace writes one local directory marketplace named for a
// channel (`spacedock` or `spacedock-edge`), hosting the single `spacedock` entry —
// the entry name is the same on both channels, which is why leaving both installed
// is ambiguous in the first place. Called once per channel so both ids resolve
// offline and both marketplaces can be registered side by side.
// buildLocalMarketplaceWithDependent is left alone rather than generalized: it
// builds one marketplace with a co-hosted dependent for the cascade probe, a
// different fixture shape that the codex test shares.
func buildNamedChannelMarketplace(t *testing.T, root, name string) string {
	t.Helper()
	marketplace := filepath.Join(root, name)
	plugin := filepath.Join(marketplace, "spacedock")
	mustMkdir(t, filepath.Join(marketplace, ".claude-plugin"))
	mustMkdir(t, filepath.Join(plugin, ".claude-plugin"))
	mustMkdir(t, filepath.Join(plugin, "skills", "demo"))

	mustWrite(t, filepath.Join(marketplace, ".claude-plugin", "marketplace.json"), `{
  "name": "`+name+`",
  "owner": { "name": "CL Kao" },
  "plugins": [
    { "name": "spacedock", "source": "./spacedock", "description": "test", "category": "workflow" }
  ]
}
`)
	mustWrite(t, filepath.Join(plugin, ".claude-plugin", "plugin.json"),
		`{ "name": "spacedock", "version": "`+displayVersion()+`", "skills": "./skills/" }`+"\n")
	mustWrite(t, filepath.Join(plugin, "skills", "demo", "SKILL.md"), "---\nname: demo\ndescription: demo skill\n---\ndemo\n")
	return marketplace
}
