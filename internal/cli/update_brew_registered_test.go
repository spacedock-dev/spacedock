package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

func TestBrewUpdateRefreshesOnlySelectedPluginFromRegisteredGitMarketplace(t *testing.T) {
	for _, name := range []string{"published tag", "dirty Claude snapshot"} {
		t.Run(name, func(t *testing.T) {
			root := stageBrewRepair(t, "spacedock", "3.0.1", "3.0.1", "")
			if err := os.Remove(filepath.Join(root, "bin", "spacedock-compass")); err != nil {
				t.Fatal(err)
			}
			writeProgram(t, filepath.Join(root, "bin", "spacedock-review"), "1.1.0")
			publisher := filepath.Join(root, "market")
			listing := filepath.Join(".claude-plugin", "marketplace.json")
			writeListing := func(version string) {
				mustWrite(t, filepath.Join(publisher, listing),
					`{"name":"spacedock","plugins":[{"name":"unrelated","source":"./unrelated"},{"name":"subspace","version":"`+version+`","source":{"source":"url","url":"https://github.com/example-fixture/subspace.git","ref":"v`+version+`"}}]}`)
			}
			writeListing("1.0.0")
			testgit.InitRepo(t, publisher, "-q", "-b", "main")
			runGit := func(dir string, args ...string) {
				t.Helper()
				cmd := exec.Command("git", args...)
				cmd.Dir = dir
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("git %v: %s (%v)", args, out, err)
				}
			}
			runGit(publisher, "add", listing)
			runGit(publisher, "commit", "-qm", "Initial listing")
			t.Setenv("SPACEDOCK_MARKETPLACE_SOURCE", "example-fixture/market")
			t.Setenv("GIT_ALLOW_PROTOCOL", "file")
			t.Setenv("GIT_CONFIG_COUNT", "1")
			t.Setenv("GIT_CONFIG_KEY_0", "url.file://"+publisher+".insteadOf")
			t.Setenv("GIT_CONFIG_VALUE_0", "https://github.com/example-fixture/market.git")
			claude := filepath.Join(root, "claude-profile")
			codex := filepath.Join(root, "codex-profile")
			t.Setenv("CLAUDE_CONFIG_DIR", claude)
			t.Setenv("CODEX_HOME", codex)
			paths := map[string]string{
				"claude": filepath.Join(claude, "plugins", "marketplaces", "spacedock"),
				"codex":  filepath.Join(codex, ".tmp", "marketplaces", "spacedock"),
			}
			for _, host := range []string{"claude", "codex"} {
				dir := paths[host]
				runGit(root, "clone", "-q", "--depth=1", "https://github.com/example-fixture/market.git", dir)
				mustWrite(t, filepath.Join(root, host+"-registered"), dir)
				mustWrite(t, filepath.Join(root, host+"-subspace"), "1.0.0")
			}
			if name == "dirty Claude snapshot" {
				mustWrite(t, filepath.Join(paths["claude"], "local-notes"), "preserve")
			}
			writeListing("1.1.0")
			runGit(publisher, "add", listing)
			runGit(publisher, "commit", "-qm", "Publish the new tag")

			if name == "dirty Claude snapshot" {
				out := runBrewRepairFixture(t, root, 1)
				if !strings.Contains(out, "marketplace snapshot has local changes") {
					t.Fatalf("missing reason for refusing the dirty checkout:\n%s", out)
				}
				assertBrewRepairWrites(t, root, "codex subspace")
				return
			}
			runBrewRepairFixture(t, root, 0)
			assertBrewRepairWrites(t, root, "claude subspace\ncodex subspace")
			runBrewRepairFixture(t, root, 0)
			assertBrewRepairWrites(t, root, "claude subspace\ncodex subspace")
		})
	}
}
