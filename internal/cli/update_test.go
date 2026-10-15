package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type switchHost struct {
	*fakeHost
	flag  string
	after []pluginInventoryEntry
}

func (h *switchHost) PluginInventory(host string) ([]pluginInventoryEntry, error) {
	if _, err := os.Stat(h.flag); err == nil {
		return h.after, nil
	}
	return h.fakeHost.PluginInventory(host)
}

func claudeOnPath(string) (string, error) { return "/bin/claude", nil }

func updateDepsFor(s *stagedInstall, host hostOps, script string) updateDeps {
	env := append([]string{}, s.env...)
	if script != "" {
		env = append(env, "SPACEDOCK_INSTALL_SCRIPT="+script)
	}
	return updateDeps{env: env, dir: s.project, host: host, selfPath: filepath.Join(s.first, "spacedock"), lookPath: claudeOnPath,
		fetchScript: func(string) error { return errors.New("no network in tests") }}
}

func TestUpdateShowsEachPieceBeforeAndAfterAndIsIdempotent(t *testing.T) {
	s := stageInstall(t)
	flag := filepath.Join(t.TempDir(), "plugins-updated")
	seen := filepath.Join(t.TempDir(), "seen")
	beta100 := filepath.Join(s.dataHome, "beta", "1.0.0", "beta", "beta")
	writeProgram(t, beta100, "beta has no version flag")
	if err := os.Remove(filepath.Join(s.first, "spacedock-beta")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(beta100, filepath.Join(s.first, "spacedock-beta")); err != nil {
		t.Fatal(err)
	}
	s.host.inventory[2].Version = "1.0.0"
	host := &switchHost{fakeHost: s.host, flag: flag, after: []pluginInventoryEntry{
		s.host.inventory[0], s.host.inventory[1], {ID: "beta@spacedock-edge", Version: "1.0.1", Installed: true, Enabled: true}}}
	beta101 := filepath.Join(s.dataHome, "beta", "1.0.1", "beta", "beta")
	script := filepath.Join(t.TempDir(), "install.sh")
	mustWrite(t, script, "#!/bin/sh\n"+
		"echo \"$SPACEDOCK_CHANNEL $SPACEDOCK_INSTALL_DIR\" > '"+seen+"'\n"+
		"mkdir -p '"+filepath.Dir(beta101)+"' && echo x > '"+beta101+"'\n"+
		"ln -sf '"+beta101+"' '"+filepath.Join(s.first, "spacedock-beta")+"'\n"+
		"touch '"+flag+"'\n")

	var stdout, stderr bytes.Buffer
	if code := runUpdate(updateDepsFor(s, host, script), &stdout, &stderr); code != 0 {
		t.Fatalf("update exited %d: %s", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{
		"beta program       1.0.0 -> 1.0.1",
		"beta plugin        1.0.0 -> 1.0.1",
		"alpha plugin       0.16.0-beta.2   unchanged",
		"spacedock          0.28.0-pre6     unchanged",
		"spacedock beta --restart",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("update output lacks %q:\n%s", want, out)
		}
	}
	if got, _ := os.ReadFile(seen); strings.TrimSpace(string(got)) != "edge "+s.first {
		t.Errorf("install script saw %q, want the edge channel and the directory of the running binary %q", got, s.first)
	}

	stdout.Reset()
	if code := runUpdate(updateDepsFor(s, host, script), &stdout, &stderr); code != 0 {
		t.Fatalf("second update exited %d: %s", code, stderr.String())
	}
	if strings.Contains(stdout.String(), "->") {
		t.Errorf("second update still shows a change:\n%s", stdout.String())
	}
}

func TestUpdateOnHomebrewChangesNoProgramAndNamesTheBrewCommand(t *testing.T) {
	for path, formula := range map[string]string{
		"/opt/homebrew/Caskroom/spacedock/0.28.0/spacedock":      "brew upgrade spacedock\n",
		"/opt/homebrew/Caskroom/spacedock@next/0.28.0/spacedock": "brew upgrade spacedock@next\n",
	} {
		s := stageInstall(t)
		ran := filepath.Join(t.TempDir(), "ran")
		script := filepath.Join(t.TempDir(), "install.sh")
		mustWrite(t, script, "#!/bin/sh\ntouch '"+ran+"'\n")
		d := updateDepsFor(s, s.host, script)
		d.selfPath = path

		var stdout, stderr bytes.Buffer
		if code := runUpdate(d, &stdout, &stderr); code != 0 {
			t.Fatalf("update exited %d: %s", code, stderr.String())
		}
		if _, err := os.Stat(ran); err == nil {
			t.Errorf("%s: update ran the install script on a Homebrew install", path)
		}
		if !strings.Contains(stdout.String(), formula) {
			t.Errorf("%s: output lacks %q:\n%s", path, formula, stdout.String())
		}
		if len(s.host.installCmds) == 0 {
			t.Errorf("%s: update did not refresh the plugins", path)
		}
	}
}

func TestUpdateChangesNothingWhenTheInstallScriptCannotRun(t *testing.T) {
	s := stageInstall(t)

	t.Run("download fails", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		if code := runUpdate(updateDepsFor(s, s.host, ""), &stdout, &stderr); code != 1 {
			t.Fatalf("exit = %d, want 1", code)
		}
		if !strings.Contains(stderr.String(), "nothing was changed") || len(s.host.installCmds) != 0 {
			t.Errorf("stderr=%q installCmds=%v, want a refusal that changed no plugin", stderr.String(), s.host.installCmds)
		}
	})
	t.Run("script fails", func(t *testing.T) {
		script := filepath.Join(t.TempDir(), "install.sh")
		mustWrite(t, script, "#!/bin/sh\necho partial\nexit 3\n")
		var stdout, stderr bytes.Buffer
		if code := runUpdate(updateDepsFor(s, s.host, script), &stdout, &stderr); code != 1 {
			t.Fatalf("exit = %d, want 1", code)
		}
		if strings.Contains(stdout.String(), "Claude Code plugins") {
			t.Errorf("a failed update printed a result table:\n%s", stdout.String())
		}
	})
}

func TestUpdateFailsWhenThePluginStepLeavesNoSpacedockPlugin(t *testing.T) {
	s := stageInstall(t)
	s.host.inventory = s.host.inventory[1:]
	script := filepath.Join(t.TempDir(), "install.sh")
	mustWrite(t, script, "#!/bin/sh\nexit 0\n")

	var stdout, stderr bytes.Buffer
	if code := runUpdate(updateDepsFor(s, s.host, script), &stdout, &stderr); code != 1 {
		t.Fatalf("exit = %d, want 1 when no Spacedock plugin is installed afterwards", code)
	}
	if !strings.Contains(stderr.String(), "the Spacedock plugin is not installed") {
		t.Errorf("stderr = %q, want it to name the missing plugin", stderr.String())
	}
}

func TestUpdateOnStableListsNoExtension(t *testing.T) {
	s := stageInstall(t)
	devBranch = "main"
	s.host.inventory[0].ID = "spacedock@spacedock"
	script := filepath.Join(t.TempDir(), "install.sh")
	mustWrite(t, script, "#!/bin/sh\n:\n")

	var stdout, stderr bytes.Buffer
	if code := runUpdate(updateDepsFor(s, s.host, script), &stdout, &stderr); code != 0 {
		t.Fatalf("update exited %d: %s", code, stderr.String())
	}
	for _, label := range []string{"alpha program", "beta plugin"} {
		if strings.Contains(stdout.String(), label) {
			t.Errorf("stable update prints %q:\n%s", label, stdout.String())
		}
	}
}
