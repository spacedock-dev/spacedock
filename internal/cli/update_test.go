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

func hostsOnPath(names ...string) func(string) (string, error) {
	return func(name string) (string, error) {
		for _, n := range names {
			if n == name {
				return "/bin/" + name, nil
			}
		}
		return "", errors.New("not on PATH")
	}
}

func claudeOnPath(name string) (string, error) { return hostsOnPath("claude")(name) }

type twoHosts struct {
	*fakeHost
	codexPlugins    []pluginInventoryEntry
	codexInstalls   int
	codexInstallErr error
}

func (h *twoHosts) PluginInventory(host string) ([]pluginInventoryEntry, error) {
	if host == "codex" {
		return h.codexPlugins, nil
	}
	return h.fakeHost.PluginInventory(host)
}

func (h *twoHosts) Install(host, source, branch string) (string, error) {
	if host == "codex" {
		h.codexInstalls++
		return "", h.codexInstallErr
	}
	return h.fakeHost.Install(host, source, branch)
}

func codexPlugins(betaVersion string) []pluginInventoryEntry {
	return []pluginInventoryEntry{
		{ID: "spacedock@spacedock-edge", Version: "0.28.0-pre6", Installed: true, Enabled: true},
		{ID: "alpha@spacedock", Version: "0.16.0-beta.2", Installed: true, Enabled: true},
		{ID: "beta@spacedock", Version: betaVersion, Installed: true, Enabled: true},
	}
}

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
		s.host.inventory[0], s.host.inventory[1], {ID: "beta@spacedock", Version: "1.0.1", Installed: true, Enabled: true}}}
	beta101 := filepath.Join(s.dataHome, "beta", "1.0.1", "beta", "beta")
	script := filepath.Join(t.TempDir(), "install.sh")
	mustWrite(t, script, "#!/bin/sh\n"+
		"echo \"$SPACEDOCK_CHANNEL $SPACEDOCK_INSTALL_DIR $SPACEDOCK_PLUGINS\" > '"+seen+"'\n"+
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
	if got, _ := os.ReadFile(seen); strings.TrimSpace(string(got)) != "edge "+s.first+" none" {
		t.Errorf("install script saw %q, want the edge channel, the directory of the running binary %q and no plugin step", got, s.first)
	}

	stdout.Reset()
	if code := runUpdate(updateDepsFor(s, host, script), &stdout, &stderr); code != 0 {
		t.Fatalf("second update exited %d: %s", code, stderr.String())
	}
	if strings.Contains(stdout.String(), "->") {
		t.Errorf("second update still shows a change:\n%s", stdout.String())
	}
	if strings.Contains(stdout.String(), "spacedock beta --restart") {
		t.Errorf("second update tells the user to restart beta, which did not change:\n%s", stdout.String())
	}
}

func TestUpdateShowsTheInstallScriptsFailedPluginLinesWhenTheScriptSucceeds(t *testing.T) {
	s := stageInstall(t)
	script := filepath.Join(t.TempDir(), "install.sh")
	const failed = "spacedock install: claude plugin update beta@spacedock failed: exit status 1"
	mustWrite(t, script, "#!/bin/sh\necho 'Updating plugin noise'\necho '"+failed+"'\nexit 0\n")

	var stdout, stderr bytes.Buffer
	if code := runUpdate(updateDepsFor(s, s.host, script), &stdout, &stderr); code != 0 {
		t.Fatalf("update exited %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), failed) {
		t.Errorf("update output lacks the install script's failure line:\n%s", stdout.String())
	}
	if strings.Contains(stdout.String(), "Updating plugin noise") {
		t.Errorf("update output repeats the install script's other output:\n%s", stdout.String())
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
	s.host.inventoryAfterInstall = s.host.inventory[1:]
	script := filepath.Join(t.TempDir(), "install.sh")
	mustWrite(t, script, "#!/bin/sh\nexit 0\n")

	var stdout, stderr bytes.Buffer
	if code := runUpdate(updateDepsFor(s, s.host, script), &stdout, &stderr); code != 1 {
		t.Fatalf("exit = %d, want 1 when no Spacedock plugin is installed afterwards", code)
	}
	if !strings.Contains(stderr.String(), "the Spacedock plugin is not installed on Claude Code") {
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

func TestUpdateRefreshesTheCodexPluginsWithTheNewProgramAndLabelsEachHost(t *testing.T) {
	s := stageInstall(t)
	host := &twoHosts{fakeHost: s.host, codexPlugins: codexPlugins("1.0.0")}
	script := filepath.Join(t.TempDir(), "install.sh")
	mustWrite(t, script, "#!/bin/sh\n:\n")
	d := updateDepsFor(s, host, script)
	d.lookPath = hostsOnPath("claude", "codex")
	var ran []string
	d.runSelf = func(args ...string) (string, error) {
		ran = append(ran, strings.Join(args, " "))
		if args[2] != "codex" {
			return "", nil
		}
		host.codexPlugins = codexPlugins("1.0.1")
		return "noise\nspacedock install: codex plugin add beta@spacedock failed: exit status 1\n", nil
	}

	var stdout, stderr bytes.Buffer
	if code := runUpdate(d, &stdout, &stderr); code != 0 {
		t.Fatalf("update exited %d: %s", code, stderr.String())
	}
	if strings.Join(ran, "; ") != "install --host claude; install --host codex" || host.codexInstalls != 0 {
		t.Errorf("runSelf calls = %v, in-process codex installs = %d, want the new program to install for each host", ran, host.codexInstalls)
	}
	out := stdout.String()
	claude, codex, _ := strings.Cut(out, "Codex plugins\n")
	for _, want := range []string{"Claude Code plugins\n", "alpha plugin       0.16.0-beta.2   unchanged"} {
		if !strings.Contains(claude, want) {
			t.Errorf("Claude section lacks %q:\n%s", want, out)
		}
	}
	if !strings.Contains(codex, "beta plugin        1.0.0 -> 1.0.1") {
		t.Errorf("Codex section lacks the beta change:\n%s", out)
	}
	for _, want := range []string{"spacedock plugin   ",
		"spacedock install: codex plugin add beta@spacedock failed", "/reload-plugins", "start a new Codex session"} {
		if !strings.Contains(out, want) {
			t.Errorf("update output lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(claude, "->") || strings.Contains(out, "noise") {
		t.Errorf("Claude rows show a change Codex made, or the new program's other output was printed:\n%s", out)
	}
}

func TestUpdateOnHomebrewRefreshesCodexInProcess(t *testing.T) {
	s := stageInstall(t)
	host := &twoHosts{fakeHost: s.host, codexPlugins: codexPlugins("1.0.0")}
	d := updateDepsFor(s, host, "")
	d.selfPath = "/opt/homebrew/Caskroom/spacedock/0.28.0/spacedock"
	d.lookPath = hostsOnPath("codex")
	d.runSelf = func(...string) (string, error) {
		t.Fatal("Homebrew has not replaced the program; update must not run it")
		return "", nil
	}

	var stdout, stderr bytes.Buffer
	if code := runUpdate(d, &stdout, &stderr); code != 0 {
		t.Fatalf("update exited %d: %s", code, stderr.String())
	}
	if host.codexInstalls != 1 || len(s.host.installCmds) != 0 {
		t.Errorf("codex installs = %d, claude install commands = %v, want one Codex refresh and none for the absent claude", host.codexInstalls, s.host.installCmds)
	}
	if !strings.Contains(stdout.String(), "claude is not on PATH, so its plugins were skipped.") {
		t.Errorf("output does not name the skipped claude:\n%s", stdout.String())
	}
	if strings.Contains(stdout.String(), "/reload-plugins") {
		t.Errorf("a Claude restart hint was printed though Claude was not refreshed:\n%s", stdout.String())
	}
}

func TestUpdateSkipsAHostThatIsAbsentOrHoldsNoSpacedockPlugin(t *testing.T) {
	cases := []struct {
		name    string
		onPath  []string
		plugins []pluginInventoryEntry
		want    string
	}{
		{"codex absent", []string{"claude"}, codexPlugins("1.0.0"), "codex is not on PATH, so its plugins were skipped."},
		{"codex without the plugin", []string{"claude", "codex"}, codexPlugins("1.0.0")[1:], "codex has no Spacedock plugin, so its plugins were skipped."},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := stageInstall(t)
			host := &twoHosts{fakeHost: s.host, codexPlugins: tc.plugins}
			script := filepath.Join(t.TempDir(), "install.sh")
			mustWrite(t, script, "#!/bin/sh\n:\n")
			d := updateDepsFor(s, host, script)
			d.lookPath = hostsOnPath(tc.onPath...)
			var ran []string
			d.runSelf = func(args ...string) (string, error) { ran = append(ran, strings.Join(args, " ")); return "", nil }

			var stdout, stderr bytes.Buffer
			if code := runUpdate(d, &stdout, &stderr); code != 0 {
				t.Fatalf("update exited %d: %s", code, stderr.String())
			}
			if !strings.Contains(stdout.String(), tc.want) || strings.Join(ran, "; ") != "install --host claude" {
				t.Errorf("want %q and only the claude install, got %v:\n%s", tc.want, ran, stdout.String())
			}
			if strings.Contains(stdout.String(), "start a new Codex session") {
				t.Errorf("a Codex restart hint was printed though Codex was skipped:\n%s", stdout.String())
			}
		})
	}
}

func TestUpdateWithNeitherHostStillUpdatesProgramsAndNamesTheFix(t *testing.T) {
	s := stageInstall(t)
	ran := filepath.Join(t.TempDir(), "ran")
	script := filepath.Join(t.TempDir(), "install.sh")
	mustWrite(t, script, "#!/bin/sh\ntouch '"+ran+"'\n")
	d := updateDepsFor(s, &twoHosts{fakeHost: s.host}, script)
	d.lookPath = hostsOnPath()

	var stdout, stderr bytes.Buffer
	if code := runUpdate(d, &stdout, &stderr); code != 0 {
		t.Fatalf("update exited %d: %s", code, stderr.String())
	}
	if _, err := os.Stat(ran); err != nil {
		t.Errorf("the install script did not run, so no program was updated")
	}
	for _, want := range []string{"claude is not on PATH", "codex is not on PATH", "spacedock install --host claude", "spacedock install --host codex"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("output lacks %q:\n%s", want, stdout.String())
		}
	}
}

func TestUpdateFailsAndNamesCodexWhenItsPluginUpdateFails(t *testing.T) {
	s := stageInstall(t)
	host := &twoHosts{fakeHost: s.host, codexPlugins: codexPlugins("1.0.0"), codexInstallErr: errors.New("codex plugin add spacedock@spacedock-edge: exit status 1")}
	d := updateDepsFor(s, host, "")
	d.selfPath = "/opt/homebrew/Caskroom/spacedock/0.28.0/spacedock"
	d.lookPath = hostsOnPath("codex")

	var stdout, stderr bytes.Buffer
	if code := runUpdate(d, &stdout, &stderr); code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "the Codex plugin update failed") {
		t.Errorf("stderr = %q, want it to name Codex", stderr.String())
	}
}

func TestUpdateLeavesAClaudeWithoutTheSpacedockPluginAlone(t *testing.T) {
	s := stageInstall(t)
	s.host.inventory = s.host.inventory[1:]
	host := &twoHosts{fakeHost: s.host, codexPlugins: codexPlugins("1.0.0")}
	seen := filepath.Join(t.TempDir(), "seen")
	script := filepath.Join(t.TempDir(), "install.sh")
	mustWrite(t, script, "#!/bin/sh\necho \"$SPACEDOCK_PLUGINS\" > '"+seen+"'\n")
	d := updateDepsFor(s, host, script)
	d.lookPath = hostsOnPath("claude", "codex")
	var ran []string
	d.runSelf = func(args ...string) (string, error) { ran = append(ran, strings.Join(args, " ")); return "", nil }

	var stdout, stderr bytes.Buffer
	if code := runUpdate(d, &stdout, &stderr); code != 0 {
		t.Fatalf("update exited %d: %s", code, stderr.String())
	}
	if strings.Join(ran, "; ") != "install --host codex" {
		t.Errorf("installs = %v, want Codex only", ran)
	}
	if got, _ := os.ReadFile(seen); strings.TrimSpace(string(got)) != "none" {
		t.Errorf("the install script was not told to skip its plugin step: %q", got)
	}
	if !strings.Contains(stdout.String(), "claude has no Spacedock plugin, so its plugins were skipped.") {
		t.Errorf("output does not name the skipped claude:\n%s", stdout.String())
	}
}
