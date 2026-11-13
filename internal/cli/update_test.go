package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
		if h.codexInstallErr == nil {
			h.codexPlugins = codexPlugins("1.0.1")
		}
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
	for _, ch := range []struct{ name, branch, pluginID, word string }{
		{"edge", "next", "spacedock@spacedock-edge", "edge"},
		{"stable", "main", "spacedock@spacedock", "stable"},
	} {
		t.Run(ch.name, func(t *testing.T) {
			s := stageInstall(t)
			devBranch = ch.branch
			s.host.inventory[0].ID = ch.pluginID
			flag := filepath.Join(t.TempDir(), "plugins-updated")
			seen := filepath.Join(t.TempDir(), "seen")
			beta100 := filepath.Join(s.dataHome, "beta", "1.0.0", "beta", "beta")
			writeProgramExiting(t, beta100, 2, "unrecognized arguments: --version")
			if err := os.Remove(filepath.Join(s.first, "spacedock-beta")); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(beta100, filepath.Join(s.first, "spacedock-beta")); err != nil {
				t.Fatal(err)
			}
			s.host.inventory[2].Version = "1.0.0"
			host := &switchHost{fakeHost: s.host, flag: flag, after: []pluginInventoryEntry{
				{ID: ch.pluginID, Version: "0.28.0-pre7", Installed: true, Enabled: true},
				s.host.inventory[1], {ID: "beta@spacedock", Version: "1.0.1", Installed: true, Enabled: true}}}
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
				"spacedock plugin   0.28.0-pre6 -> 0.28.0-pre7",
				"spacedock          0.28.0-pre6     unchanged",
				"spacedock beta --restart",
			} {
				if !strings.Contains(out, want) {
					t.Errorf("update output lacks %q:\n%s", want, out)
				}
			}
			if got, _ := os.ReadFile(seen); strings.TrimSpace(string(got)) != ch.word+" "+s.first+" none" {
				t.Errorf("install script saw %q, want the %s channel, the directory of the running binary %q and no plugin step", got, ch.word, s.first)
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
		})
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
		"/opt/homebrew/Caskroom/spacedock/0.28.0/spacedock":      "brew upgrade spacedock spacedock-dev/tap/subspace-beta spacedock-dev/tap/compass\n",
		"/opt/homebrew/Caskroom/spacedock@next/0.28.0/spacedock": "brew upgrade spacedock@next spacedock-dev/tap/subspace-beta spacedock-dev/tap/compass\n",
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

func TestUpdateSkipsAHostThatIsAbsentOrHoldsNothingFromTheMarketplace(t *testing.T) {
	cases := []struct {
		name    string
		onPath  []string
		plugins []pluginInventoryEntry
		want    string
	}{
		{"codex absent", []string{"claude"}, codexPlugins("1.0.0"), "codex is not on PATH, so its plugins were skipped."},
		{"codex holding nothing from the marketplace", []string{"claude", "codex"}, []pluginInventoryEntry{}, "codex holds nothing from the Spacedock marketplace, so its plugins were skipped."},
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

func TestUpdateInstallsNothingIntoAClaudeWithoutTheSpacedockPluginAndFails(t *testing.T) {
	s := stageInstall(t)
	s.host.inventory = s.host.inventory[1:]
	host := &twoHosts{fakeHost: s.host, codexPlugins: codexPlugins("1.0.0")}
	seen := filepath.Join(t.TempDir(), "seen")
	script := filepath.Join(t.TempDir(), "install.sh")
	mustWrite(t, script, "#!/bin/sh\necho \"$SPACEDOCK_PLUGINS\" > '"+seen+"'\n")
	d := updateDepsFor(s, host, script)
	d.lookPath = hostsOnPath("claude", "codex")
	var ran []string
	d.runSelf = func(args ...string) (string, error) {
		ran = append(ran, strings.Join(args, " "))
		host.codexPlugins = codexPlugins("1.0.1")
		return "", nil
	}

	var stdout, stderr bytes.Buffer
	if code := runUpdate(d, &stdout, &stderr); code != 1 {
		t.Fatalf("update exited %d, want 1: %s", code, stderr.String())
	}
	if strings.Join(ran, "; ") != "install --host codex" {
		t.Errorf("installs = %v, want Codex only", ran)
	}
	if got, _ := os.ReadFile(seen); strings.TrimSpace(string(got)) != "none" {
		t.Errorf("the install script was not told to skip its plugin step: %q", got)
	}
	if want := "the Spacedock plugin is not installed on Claude Code. Run: spacedock install --host claude"; !strings.Contains(stderr.String(), want) {
		t.Errorf("stderr lacks %q:\n%s", want, stderr.String())
	}
}

func TestUpdateFailsAndRefreshesNoPluginWhenAnExtensionLookupFails(t *testing.T) {
	s := stageInstall(t)
	root := t.TempDir()
	dist, stubs := filepath.Join(root, "dist"), filepath.Join(root, "stubs")
	writeProgram(t, filepath.Join(root, "pkg", "spacedock"), "spacedock 9.9.9")
	writeProgramExiting(t, filepath.Join(stubs, "curl"), 22, "")
	writeProgram(t, filepath.Join(stubs, "git"), "")
	archive := "spacedock_9.9.9_" + runtime.GOOS + "_" + runtime.GOARCH + "_edge.tar.gz"
	mustMkdir(t, dist)
	if out, err := exec.Command("tar", "-czf", filepath.Join(dist, archive), "-C", filepath.Join(root, "pkg"), "spacedock").CombinedOutput(); err != nil {
		t.Fatalf("tar: %v\n%s", err, out)
	}
	data, err := os.ReadFile(filepath.Join(dist, archive))
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	mustWrite(t, filepath.Join(dist, "checksums.txt"), hex.EncodeToString(sum[:])+"  "+archive+"\n")
	script, err := filepath.Abs(filepath.Join("..", "..", "install.sh"))
	if err != nil {
		t.Fatal(err)
	}
	d := updateDepsFor(s, s.host, script)
	d.env = append(d.env, "PATH="+stubs+":/usr/bin:/bin", "SPACEDOCK_INSTALL_FROM="+dist, "SPACEDOCK_EXTRAS=1")

	var stdout, stderr bytes.Buffer
	if code := runUpdate(d, &stdout, &stderr); code != 1 {
		t.Fatalf("exit = %d, want 1; stderr:\n%s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "the spacedock marketplace listing names no release of subspace") {
		t.Errorf("stderr does not name the unreadable listing:\n%s", stderr.String())
	}
	if strings.Contains(stdout.String(), "Claude Code plugins") || len(s.host.installCmds) != 0 {
		t.Errorf("a failed lookup still refreshed plugins or printed a result table: installCmds=%v\n%s", s.host.installCmds, stdout.String())
	}
}

func TestUpdateFailsAndNamesTheExtensionWhenAProgramDownloadFailedButThePluginsMoved(t *testing.T) {
	s := stageInstall(t)
	s.host.inventory[2].Version = "1.0.2"
	host := &switchHost{fakeHost: s.host, flag: "", after: nil}
	script := filepath.Join(t.TempDir(), "install.sh")
	mustWrite(t, script, "#!/bin/sh\necho 'beta: download failed - skipped' >&2\n")

	var stdout, stderr bytes.Buffer
	if code := runUpdate(updateDepsFor(s, host, script), &stdout, &stderr); code != 1 {
		t.Fatalf("update exited %d, want 1 while beta is MIXED; stderr:\n%s", code, stderr.String())
	}
	want := "spacedock update: beta is still MIXED on Claude Code: program 1.0.1, plugin 1.0.2. Left behind: the program (failed: the install script's download of the beta program)."
	if !strings.Contains(stderr.String(), want) {
		t.Errorf("stderr lacks %q:\n%s", want, stderr.String())
	}
}

func TestUpdateFailsAndNamesAnExtensionWhoseProgramIsMissingBesideItsPlugin(t *testing.T) {
	s := stageInstall(t)
	if err := os.Remove(filepath.Join(s.first, "spacedock-alpha")); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(t.TempDir(), "install.sh")
	mustWrite(t, script, "#!/bin/sh\n:\n")

	var stdout, stderr bytes.Buffer
	if code := runUpdate(updateDepsFor(s, s.host, script), &stdout, &stderr); code != 1 {
		t.Fatalf("update exited %d, want 1; stderr:\n%s", code, stderr.String())
	}
	want := "spacedock update: alpha is still MIXED on Claude Code: program not installed, plugin 0.16.0-beta.2. Left behind: the program (failed: the install script's download of the alpha program)."
	if !strings.Contains(stderr.String(), want) {
		t.Errorf("stderr lacks %q:\n%s", want, stderr.String())
	}
}

type execInstallHost struct{ *fakeHost }

func (h execInstallHost) Install(host, source, branch string) (string, error) {
	return execHost{}.Install(host, source, branch)
}

func TestUpdateNamesARefusedMarketplaceAddWithClaudesOwnMessage(t *testing.T) {
	s := stageInstall(t)
	s.host.inventory = s.host.inventory[:1]
	refusal := `✘ Failed to add marketplace: Cannot add marketplace "spacedock": its source doesn't match its extraKnownMarketplaces entry`
	bin := t.TempDir()
	mustWrite(t, filepath.Join(bin, "refusal"), "Adding marketplace…"+refusal+"\n\n")
	mustWrite(t, filepath.Join(bin, "claude"), "#!/bin/sh\ncase \"$*\" in\n"+
		"'plugin list --json') echo '[]' ;;\n"+
		"'plugin marketplace add spacedock-dev/marketplace') cat \"${0%/*}/refusal\"; exit 1 ;;\n"+
		"esac\n")
	if err := os.Chmod(filepath.Join(bin, "claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+":/usr/bin:/bin")
	script := filepath.Join(t.TempDir(), "install.sh")
	mustWrite(t, script, "#!/bin/sh\n:\n")

	var stdout, stderr bytes.Buffer
	if code := runUpdate(updateDepsFor(s, execInstallHost{s.host}, script), &stdout, &stderr); code != 1 {
		t.Fatalf("update exited %d, want 1; stderr:\n%s", code, stderr.String())
	}
	if want := "spacedock install: claude plugin marketplace add spacedock-dev/marketplace failed: exit status 1: Adding marketplace…" + refusal + "\n"; !strings.Contains(stdout.String(), want) {
		t.Errorf("stdout lacks Claude Code's refusal %q:\n%s", want, stdout.String())
	}
	if want := "alpha is still MIXED on Claude Code: program 0.16.0-beta.2, plugin not installed. Left behind: the Claude Code plugin (failed: the Claude Code marketplace add"; !strings.Contains(stderr.String(), want) {
		t.Errorf("stderr lacks %q:\n%s", want, stderr.String())
	}
}

func TestDoctorAndUpdateCheckEveryHostThatHoldsAPieceFromTheMarketplace(t *testing.T) {
	for _, host := range []string{"claude", "codex", "pi"} {
		for _, spacedock := range []string{"present", "absent"} {
			t.Run(host+" with the Spacedock plugin "+spacedock, func(t *testing.T) {
				s := stageInstall(t)
				extensions[0].piTitle = "Alpha"
				bin, agentDir := t.TempDir(), t.TempDir()
				ids := []string{"alpha@spacedock"}
				if spacedock == "present" {
					ids = append(ids, channelPluginID(devBranch))
				}
				var entries []string
				for _, id := range ids {
					entries = append(entries, `{"id":"`+id+`","pluginId":"`+id+`","version":"0.16.0-beta.1","installPath":"/p","installed":true,"enabled":true}`)
				}
				list := map[string]string{"claude": "[]", "codex": `{"installed":[]}`, "pi": ""}
				switch host {
				case "claude":
					list["claude"] = "[" + strings.Join(entries, ",") + "]"
				case "codex":
					list["codex"] = `{"installed":[` + strings.Join(entries, ",") + `]}`
				case "pi":
					packages := `"` + piTestAlpha + `@v0.16.0-beta.1"`
					if spacedock == "present" {
						packages += `,"` + piSpacedockPackageSource + `@stable"`
					}
					writeFile(t, filepath.Join(agentDir, "settings.json"), `{"packages":[`+packages+`]}`)
				}
				mustWrite(t, filepath.Join(bin, host), "#!/bin/sh\n[ \"$*\" = 'plugin list --json' ] && echo '"+list[host]+"'\nexit 0\n")
				if err := os.Chmod(filepath.Join(bin, host), 0o755); err != nil {
					t.Fatal(err)
				}
				t.Setenv("PATH", bin+":/usr/bin:/bin")
				env := append(append([]string{}, s.env...), "PI_CODING_AGENT_DIR="+agentDir)

				var doctor bytes.Buffer
				reportVersions(nil, env, s.project, &doctor)
				if strings.Contains(doctor.String(), "No problems found.") {
					t.Errorf("doctor reads clean while the alpha %s piece differs from its program:\n%s", host, doctor.String())
				}
				script := filepath.Join(t.TempDir(), "install.sh")
				mustWrite(t, script, "#!/bin/sh\n:\n")
				d := updateDeps{env: append(env, "SPACEDOCK_INSTALL_SCRIPT="+script), dir: s.project, host: execHost{}, selfPath: filepath.Join(s.first, "spacedock"),
					lookPath: exec.LookPath, piOps: execPiRuntimeOps{}, fetchScript: func(string) error { return errors.New("no network in tests") },
					runSelf: func(...string) (string, error) { return "", nil }}
				var stdout, stderr bytes.Buffer
				if code := runUpdate(d, &stdout, &stderr); code == 0 {
					t.Errorf("update exited 0 while the alpha %s piece differs from its program:\n%s%s", host, stdout.String(), stderr.String())
				}
				remedy := "spacedock install --host " + host
				if spacedock == "absent" && (!strings.Contains(doctor.String(), remedy) || !strings.Contains(stderr.String(), remedy)) {
					t.Errorf("doctor or update does not tell the user to run %q:\n%s\n%s", remedy, doctor.String(), stderr.String())
				}
			})
		}
	}
}
