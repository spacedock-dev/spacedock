// ABOUTME: AC-1 live proof. A claude channel install against a local marketplace
// ABOUTME: for both channels must leave one enabled spacedock plugin.
package cli

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestClaudeInstallLeavesOnlySelectedChannelEnabled asserts AC-1 against the real
// claude CLI. After execHost.Install, `claude plugin list --json` reports exactly
// one enabled `spacedock@*` entry, and its id is the id of the selected channel.
// Claude reports this count, not Spacedock.
//
// The baseline moves the wrong way. If you delete the sibling uninstall step from
// installArgvSequence, the sibling-present subtest and the reverse subtest count 2.
// This is the bug that the captain reported, in a hermetic form.
//
// The three subtests give three directions. A co-installed sibling is the shape
// from the report. A fresh box shows that the tolerated non-zero exit of the absent
// sibling does not stop the run. Edge over stable shows that the map works in both
// directions.
//
// If `claude` is not on PATH, the test skips. The CI build job has no host CLI.
// Validation must run this test locally and must treat a SKIP as a failure.
func TestClaudeInstallLeavesOnlySelectedChannelEnabled(t *testing.T) {
	claudeBin, err := exec.LookPath("claude")
	if err != nil {
		t.Skip("claude not on PATH; channel-exclusivity test requires the host CLI")
	}

	cases := []struct {
		name      string
		devBranch string
		seedID    string // sibling installed before the run. Empty means a fresh box.
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
			// The isolation must apply to the process environment. execHost.Install
			// runs the host command with no explicit Env, so it inherits the
			// environment of the process (co_hosted_survival_test.go).
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

// enabledSpacedockIDs returns the id of every enabled `spacedock@*` entry that
// claude reports. This count is the quantity that AC-1 measures. It comes from the
// JSON of the host.
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

// buildNamedChannelMarketplace writes one local directory marketplace with the name
// of a channel (`spacedock` or `spacedock-edge`). The marketplace holds the single
// `spacedock` entry. Both channels use the same entry name. If both plugins stay
// installed, that entry name is ambiguous. The test calls this function once
// for each channel, so both ids resolve offline and both marketplaces can register
// at the same time. This function does not replace
// buildLocalMarketplaceWithDependent, which builds one marketplace with a co-hosted
// dependent for the cascade probe. The codex test shares that different fixture.
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
