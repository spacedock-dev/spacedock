package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// stubClaudeInventory puts a `claude` on PATH that answers `plugin list --json` with
// listing (exit 1 when listing is empty), exits 1 on every uninstall with claude's
// own stderr line, and logs each call. It returns the log path.
func stubClaudeInventory(t *testing.T, listing string) string {
	t.Helper()
	dir := t.TempDir()
	log := filepath.Join(dir, "calls.log")
	listFile := filepath.Join(dir, "list.json")
	if listing != "" {
		mustWrite(t, listFile, listing)
	}
	script := "#!/bin/sh\necho \"$*\" >> '" + log + "'\n" +
		"case \"$*\" in\n" +
		"\"plugin list --json\") [ -f '" + listFile + "' ] && { cat '" + listFile + "'; exit 0; }; exit 1 ;;\n" +
		"\"plugin uninstall \"*) echo \"✘ Failed to uninstall plugin \\\"$3\\\": Plugin \\\"$3\\\" not found in installed plugins\" >&2; exit 1 ;;\n" +
		"esac\nexit 0\n"
	mustWrite(t, filepath.Join(dir, "claude"), script)
	if err := os.Chmod(filepath.Join(dir, "claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return log
}

func installWithInventory(t *testing.T, listing string) (out string, calls string) {
	t.Helper()
	log := stubClaudeInventory(t, listing)
	out, err := execHost{}.Install("claude", "spacedock-dev/marketplace", "next")
	if err != nil {
		t.Fatalf("Install returned %v\nout=%q", err, out)
	}
	data, _ := os.ReadFile(log)
	return out, string(data)
}

func TestInstallOnACleanMachineRunsNoUninstallAndPrintsNoFailure(t *testing.T) {
	out, calls := installWithInventory(t, `[]`)
	if strings.Contains(out, "Failed to uninstall") || strings.Contains(calls, "plugin uninstall") {
		t.Errorf("a clean machine ran an uninstall.\nout=%q\ncalls=%q", out, calls)
	}
	if !strings.Contains(calls, "plugin install spacedock@spacedock-edge") {
		t.Errorf("the channel plugin was not installed: %q", calls)
	}
}

func TestInstallStillShowsAFailedRemovalOfAnInstalledPlugin(t *testing.T) {
	out, calls := installWithInventory(t, `[{"id":"subspace@spacedock-edge","scope":"user"}]`)
	if !strings.Contains(out, "spacedock install: claude plugin uninstall subspace@spacedock-edge failed") {
		t.Errorf("the failed removal of a listed plugin is not shown:\n%s", out)
	}
	if strings.Count(calls, "plugin uninstall") != 1 || !strings.Contains(calls, "plugin install spacedock@spacedock-edge") {
		t.Errorf("want exactly the listed uninstall, then the install.\ncalls=%q", calls)
	}
}

func TestInstallRemovesAPluginWhoseMarketplaceListingIsGone(t *testing.T) {
	_, calls := installWithInventory(t, `[{"id":"subspace@spacedock-edge","scope":"user","errors":["Plugin subspace not found in marketplace spacedock-edge"]}]`)
	if !strings.Contains(calls, "plugin uninstall subspace@spacedock-edge") {
		t.Errorf("a listed plugin with errors was skipped: %q", calls)
	}
}

func TestInstallLeavesAProjectScopeCopyAlone(t *testing.T) {
	_, calls := installWithInventory(t, `[{"id":"subspace@spacedock-edge","scope":"project"}]`)
	if strings.Contains(calls, "plugin uninstall") {
		t.Errorf("a project-scope copy was uninstalled at user scope: %q", calls)
	}
}

func TestInstallRunsEveryUninstallWhenTheInventoryIsUnreadable(t *testing.T) {
	for name, listing := range map[string]string{"list fails": "", "not json": "plugins: none"} {
		t.Run(name, func(t *testing.T) {
			_, calls := installWithInventory(t, listing)
			if got := strings.Count(calls, "plugin uninstall"); got != 3+len(extensions) {
				t.Errorf("uninstalls = %d, want all %d.\ncalls=%q", got, 3+len(extensions), calls)
			}
		})
	}
}

func TestStableClaudeInstallRefreshesOnlyWhenNeeded(t *testing.T) {
	for _, tc := range []struct {
		name, plugins, markets, fail                   string
		refreshMarket, refreshCompass, refreshSubspace bool
	}{
		{"empty", "[]", "[]", "", false, false, false},
		{"bundled plugins in a fresh profile", `[{"id":"cc-plugin-privacy@builtin","scope":"user"}]`, "[]", "", false, false, false},
		{"installed plugin named builtin", `[{"id":"cc-plugin-privacy@builtin","installPath":"/somewhere"}]`, "[]", "", true, true, true},
		{"unknown entry", `[{}]`, "[]", "", true, true, true},
		{"custom plugin cache", "[]", "[]", "", true, true, true},
		{"custom plugin seed", "[]", "[]", "", true, true, true},
		{"interrupted install left a cache", "[]", "[]", "", true, true, true},
		{"existing market without plugins", "[]", `[{"name":"spacedock"}]`, "", true, true, true},
		{"unreadable markets", "[]", "null", "", true, true, true},
		{"unreadable plugins", "null", "[]", "", true, true, true},
		{"installed plugin", `[{"id":"compass@spacedock","scope":"user"}]`, "[]", "", true, true, true},
		{"project plugin", `[{"id":"compass@spacedock","scope":"project"}]`, "[]", "", true, true, true},
		{"compass install fails", "[]", "[]", "plugin install compass@spacedock", false, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			config := t.TempDir()
			t.Setenv("CLAUDE_CONFIG_DIR", config)
			t.Setenv("CLAUDE_CODE_PLUGIN_CACHE_DIR", "")
			t.Setenv("CLAUDE_CODE_PLUGIN_SEED_DIR", "")
			switch tc.name {
			case "custom plugin cache":
				t.Setenv("CLAUDE_CODE_PLUGIN_CACHE_DIR", t.TempDir())
			case "custom plugin seed":
				t.Setenv("CLAUDE_CODE_PLUGIN_SEED_DIR", t.TempDir())
			case "interrupted install left a cache":
				if err := os.MkdirAll(filepath.Join(config, "plugins", "cache", "spacedock"), 0755); err != nil {
					t.Fatal(err)
				}
			}
			log := filepath.Join(dir, "calls")
			t.Setenv("CALLS", log)
			t.Setenv("PLUGINS", tc.plugins)
			t.Setenv("MARKETS", tc.markets)
			t.Setenv("FAIL", tc.fail)
			t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
			body := `#!/bin/sh
echo "$*" >> "$CALLS"
case "$*" in
"plugin list --json") echo "$PLUGINS" ;;
"plugin marketplace list --json") echo "$MARKETS" ;;
"$FAIL") exit 1 ;;
esac
`
			if err := os.WriteFile(filepath.Join(dir, "claude"), []byte(body), 0755); err != nil {
				t.Fatal(err)
			}
			if _, err := (execHost{}).Install("claude", defaultMarketplaceSource, "main"); err != nil {
				t.Fatal(err)
			}
			calls, err := os.ReadFile(log)
			if err != nil {
				t.Fatal(err)
			}
			for call, want := range map[string]bool{
				"plugin marketplace add " + defaultMarketplaceSource: true,
				"plugin install spacedock@spacedock":                 true,
				"plugin install compass@spacedock":                   true,
				"plugin install subspace@spacedock":                  true,
				"plugin marketplace update spacedock":                tc.refreshMarket,
				"plugin update compass@spacedock":                    tc.refreshCompass,
				"plugin update subspace@spacedock":                   tc.refreshSubspace,
			} {
				if strings.Contains(string(calls), call+"\n") != want {
					t.Errorf("call %q, want present=%v; calls:\n%s", call, want, calls)
				}
			}
		})
	}
}
