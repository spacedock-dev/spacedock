// ABOUTME: spacedock update — bring Spacedock, its extensions and their plugins to the
// ABOUTME: newest release on this binary's channel, then print what changed.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const installScriptURL = "https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh"

type updateDeps struct {
	env         []string
	dir         string
	host        hostOps
	selfPath    string
	lookPath    func(string) (string, error)
	fetchScript func(dst string) error
	runSelf     func(args ...string) (string, error)
	piOps       piRuntimeOps
}

func newUpdateCommand(ctx context.Context, env []string, dir string, stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "update",
		Short:              "Update Spacedock, its extensions and their plugins to the newest release",
		GroupID:            "setup",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				return cmd.Help()
			}
			if len(args) > 0 {
				fmt.Fprintf(stderr, "spacedock update: unknown argument %q\n", args[0])
				return exitCodeError{2}
			}
			applyDevBranchOverride(env)
			applyMarketplaceSourceOverride(env)
			selfPath, _ := resolvedLauncherBin()
			deps := updateDeps{env: env, dir: dir, host: execHost{}, selfPath: selfPath, lookPath: exec.LookPath, piOps: execPiRuntimeOps{},
				fetchScript: func(dst string) error {
					return exec.CommandContext(ctx, "curl", "-fsSL", "-o", dst, installScriptURL).Run()
				}}
			if selfPath != "" {
				deps.runSelf = func(args ...string) (string, error) {
					cmd := exec.CommandContext(ctx, selfPath, args...)
					cmd.Env = env
					out, err := cmd.CombinedOutput()
					return string(out), err
				}
			}
			if code := runUpdate(deps, stdout, stderr); code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
	setSetupHelp(cmd, stdout, `
Update every program and plugin on the channel of this binary, for each of Claude
Code and Codex that is on PATH and has the Spacedock plugin installed, and the Pi
packages when Pi holds the Spacedock package. Homebrew
owns its own programs: update prints the brew command and updates only the plugins.
SPACEDOCK_INSTALL_SCRIPT=PATH runs that script in place of the downloaded one.
`)
	return cmd
}

func runUpdate(d updateDeps, stdout, stderr io.Writer) int {
	inputs := versionInputs{env: d.env, dir: d.dir, host: d.host, lookPath: d.lookPath, selfPath: d.selfPath, selfVersion: displayVersion()}
	if d.selfPath != "" {
		inputs.selfVersion = programReportedVersion(d.selfPath, d.env, "spacedock ")
	}
	before := collectVersions(inputs)
	brew := d.selfPath != "" && isHomebrewPath(d.selfPath)

	owner := "installed with install.sh"
	if brew {
		owner = "installed with Homebrew: " + d.selfPath
	}
	fmt.Fprintf(stdout, "Updating Spacedock on the %s channel (%s)\n", before.channel, owner)

	hosts := planHostUpdates(before)
	if !brew {
		if err := runInstallScript(d, stdout, stderr); err != nil {
			fmt.Fprintf(stderr, "spacedock update: %v\n", err)
			return 1
		}
	}
	for _, h := range hosts {
		if !h.refresh {
			continue
		}
		// The new program owns the install sequence, so a companion it adds reaches
		// the host in the same run; Homebrew has not replaced the program yet.
		var out string
		if brew || d.runSelf == nil {
			out, h.err = d.host.Install(h.name, channelMarketplaceSource(devBranch), devBranch)
		} else {
			out, h.err = d.runSelf("install", "--host", h.name)
		}
		printInstallFailures(stdout, out)
		if h.err != nil && strings.TrimSpace(out) != "" {
			h.err = fmt.Errorf("%w\n%s", h.err, strings.TrimSpace(out))
		}
	}

	pi := planPiUpdate(d)
	if pi.refresh {
		var out string
		if brew || d.runSelf == nil {
			var buf bytes.Buffer
			if code := runInitWithPi(context.Background(), []string{"--host", "pi"}, d.host, d.piOps, d.env, &buf, &buf); code != 0 {
				pi.err = fmt.Errorf("spacedock install --host pi exited %d", code)
			}
			out = buf.String()
		} else {
			out, pi.err = d.runSelf("install", "--host", "pi")
		}
		printInstallFailures(stdout, out)
		if pi.err != nil && strings.TrimSpace(out) != "" {
			pi.err = fmt.Errorf("%w\n%s", pi.err, strings.TrimSpace(out))
		}
		pi.after = piPackageRows(pi.agentDir, pi.home)
	}

	if d.selfPath != "" && !brew {
		inputs.selfVersion = programReportedVersion(d.selfPath, d.env, "spacedock ")
	}
	after := collectVersions(inputs)
	for _, h := range hosts {
		h.after = after.hostPlugins[h.name]
	}
	printUpdateChanges(stdout, before, after, brewFormula(d.selfPath, brew), hosts, pi)
	code := 0
	if pi.err != nil {
		fmt.Fprintf(stderr, "spacedock update: the Pi package update failed: %v\n", pi.err)
		code = 1
	}
	for _, h := range hosts {
		switch {
		case !h.refresh:
		case h.err != nil:
			fmt.Fprintf(stderr, "spacedock update: the %s plugin update failed: %v\n", h.title, h.err)
			code = 1
		case len(h.after) == 0:
			fmt.Fprintf(stderr, "spacedock update: the Spacedock plugin is not installed on %s. Run: spacedock install --host %s\n", h.title, h.name)
			code = 1
		}
	}
	return code
}

// hostUpdate is one host's share of an update. A host that is absent or holds no
// Spacedock plugin is skipped, never installed into: update refreshes, install adds.
type hostUpdate struct {
	name, title   string
	refresh       bool
	skipped       string
	before, after []pieceRow
	err           error
}

func planHostUpdates(before versionReport) []*hostUpdate {
	var updates []*hostUpdate
	for _, status := range before.hosts {
		updates = append(updates, &hostUpdate{name: status.name, title: status.title, skipped: status.skipped,
			refresh: status.skipped == "", before: before.hostPlugins[status.name]})
	}
	return updates
}

// piUpdate is Pi's share of an update; Pi has no plugin inventory, so settings.json is read.
type piUpdate struct {
	agentDir, home string
	refresh        bool
	skipped        string
	before, after  []pieceRow
	err            error
}

func planPiUpdate(d updateDeps) *piUpdate {
	cfg := piRuntimeConfigFromEnv(d.env, d.dir, "")
	u := &piUpdate{agentDir: cfg.agentDir, home: cfg.home}
	if _, err := d.lookPath("pi"); err != nil {
		u.skipped = "pi is not on PATH, so its packages were skipped."
		return u
	}
	status := d.piOps.SpacedockPackageStatus(cfg.agentDir, cfg.home)
	switch {
	case !status.registered:
		u.skipped = "pi has no Spacedock package, so its packages were skipped."
	case !strings.HasPrefix(status.source, "git:"):
		u.skipped = fmt.Sprintf("the Spacedock package in pi is %s, not a published source, so its packages were skipped.", status.source)
	default:
		u.refresh = true
		u.before = piPackageRows(cfg.agentDir, cfg.home)
	}
	return u
}

func piPackageRows(agentDir, home string) []pieceRow {
	names := []string{"spacedock"}
	for _, ext := range extensions {
		if ext.piTitle != "" {
			names = append(names, ext.name)
		}
	}
	var rows []pieceRow
	for _, name := range names {
		_, ref, ok := piRegisteredPackage(agentDir, home, name)
		version := notInstalled
		if ok {
			version = refOrUnpinned(ref)
		}
		rows = append(rows, pieceRow{label: name + " package", version: version})
	}
	return rows
}

func pieceVersion(rows []pieceRow, label string) string {
	for _, r := range rows {
		if r.label == label {
			return r.version
		}
	}
	return ""
}

// runInstallScript installs into the directory of this binary, so the new program
// replaces it.
func runInstallScript(d updateDeps, stdout, stderr io.Writer) error {
	script := envMap(d.env)["SPACEDOCK_INSTALL_SCRIPT"]
	if script == "" {
		tmp, err := os.MkdirTemp("", "spacedock-update-")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmp)
		script = filepath.Join(tmp, "install.sh")
		if err := d.fetchScript(script); err != nil {
			return fmt.Errorf("could not download the install script, nothing was changed: %w", err)
		}
	}
	env := append(append([]string{}, d.env...), "SPACEDOCK_CHANNEL="+selectedChannelWord(), "SPACEDOCK_PLUGINS=none")
	if d.selfPath != "" {
		env = append(env, "SPACEDOCK_INSTALL_DIR="+filepath.Dir(d.selfPath))
	}
	cmd := exec.Command("sh", script)
	cmd.Env = env
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, stderr
	if err := cmd.Run(); err != nil {
		io.Copy(stdout, &out)
		return fmt.Errorf("the install script failed: %w", err)
	}
	printInstallFailures(stdout, out.String())
	return nil
}

func printInstallFailures(w io.Writer, out string) {
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "spacedock install:") {
			fmt.Fprintln(w, line)
		}
	}
}

func brewFormula(selfPath string, brew bool) string {
	switch {
	case !brew:
		return ""
	case isEdgeCaskPath(selfPath):
		return "spacedock@next"
	}
	return "spacedock"
}

var restartAdvice = map[string]string{
	"claude": "Open Claude Code sessions keep the old plugins until you restart them or run /reload-plugins.",
	"codex":  "Open Codex sessions keep the old plugins until you start a new Codex session.",
	"pi":     "Open Pi sessions keep the old skills until you run /reload.",
}

func printUpdateChanges(w io.Writer, before, after versionReport, formula string, hosts []*hostUpdate, pi *piUpdate) {
	printRow := func(label, b, a string) {
		if b == a {
			fmt.Fprintf(w, "  %-18s %-15s unchanged\n", label, a)
			return
		}
		fmt.Fprintf(w, "  %-18s %s -> %s\n", label, b, a)
	}
	row := func(label string) { printRow(label, before.version(label), after.version(label)) }
	fmt.Fprintln(w, "\nPrograms")
	if formula != "" {
		fmt.Fprintf(w, "  Homebrew owns spacedock and its extensions. Run: brew upgrade %s\n", formula)
	} else {
		row("spacedock")
		for _, ext := range extensions {
			row(ext.name + " program")
		}
	}
	for _, h := range hosts {
		fmt.Fprintf(w, "%s plugins\n", h.title)
		if !h.refresh {
			fmt.Fprintln(w, "  "+h.skipped)
			continue
		}
		for _, r := range h.after {
			printRow(r.label, pieceVersion(h.before, r.label), r.version)
		}
	}
	fmt.Fprintln(w, "Pi packages")
	if !pi.refresh {
		fmt.Fprintln(w, "  "+pi.skipped)
	}
	for _, r := range pi.after {
		printRow(r.label, pieceVersion(pi.before, r.label), r.version)
	}
	if len(after.problems) > 0 {
		fmt.Fprintln(w)
		for _, p := range after.problems {
			fmt.Fprintln(w, p)
		}
	}
	fmt.Fprintln(w)
	refreshed := false
	for _, h := range hosts {
		if h.refresh {
			refreshed = true
			fmt.Fprintln(w, restartAdvice[h.name])
		}
	}
	if pi.refresh {
		refreshed = true
		fmt.Fprintln(w, restartAdvice["pi"])
	}
	if !refreshed {
		fmt.Fprintln(w, "No plugins were updated. Run: spacedock install --host claude   or: spacedock install --host codex")
	}
	for _, ext := range extensions {
		if before.version(ext.name+" program") != after.version(ext.name+" program") {
			for _, line := range ext.afterUpdate {
				fmt.Fprintln(w, line)
			}
		}
	}
	if formula != "" {
		fmt.Fprintln(w, "Run `spacedock doctor` after `brew upgrade` to check that every piece matches.")
	}
}
