// ABOUTME: spacedock update — bring Spacedock, its extensions and their plugins to the
// ABOUTME: newest release on this binary's channel, then print what changed.
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/spacedock-dev/spacedock/internal/release"
)

const installScriptURL = "https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh"

type updateDeps struct {
	verbose     bool
	env         []string
	dir         string
	host        hostOps
	selfPath    string
	lookPath    func(string) (string, error)
	fetchScript func(dst string) error
	runSelf     func(args ...string) (string, error)
	piOps       piRuntimeOps
	ctx         context.Context
}

func newUpdateCommand(ctx context.Context, env []string, dir string, stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "update [--verbose]",
		Short:              "Update Spacedock, its extensions and their plugins to the newest release",
		GroupID:            "setup",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				return cmd.Help()
			}
			verbose := len(args) > 0 && args[0] == "--verbose"
			if verbose {
				args = args[1:]
			}
			if len(args) > 0 {
				fmt.Fprintf(stderr, "spacedock update: unknown argument %q\n", args[0])
				return exitCodeError{2}
			}
			applyDevBranchOverride(env)
			applyMarketplaceSourceOverride(env)
			selfPath, _ := resolvedLauncherBin()
			deps := updateDeps{verbose: verbose, env: env, dir: dir, host: execHost{}, selfPath: selfPath, lookPath: exec.LookPath, piOps: execPiRuntimeOps{}, ctx: ctx,
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
	cmd.Flags().Bool("verbose", false, "Prefix each line with the time since start and show how long each step took")
	setSetupHelp(cmd, stdout, `
Update every program and plugin on the channel of this binary, for each of Claude
Code and Codex that is on PATH and has the Spacedock plugin installed, and the Pi
packages when Pi holds the Spacedock package. A host holding other Spacedock pieces
without that plugin or package is named and fails the update. On Homebrew, update
can repair installed extension programs and their companions. Upgrade the Spacedock
cask separately.
SPACEDOCK_INSTALL_SCRIPT=PATH runs that script in place of the downloaded one.
`)
	return cmd
}

func runUpdate(d updateDeps, stdout, stderr io.Writer) int {
	progress := &updateProgress{plain: stderr}
	if d.verbose {
		zero := time.Now()
		verboseOut := &updateLineWriter{out: stdout, zero: zero}
		progress.trace = &updateLineWriter{out: stderr, zero: zero}
		verboseOut.peer, progress.trace.peer = progress.trace, verboseOut
		progress.last = []time.Time{{}}
		stdout, stderr = verboseOut, progress.trace
		defer progress.finish()
	}
	endBefore := progress.step("Checking installed versions before the update")
	inputs := versionInputs{env: d.env, dir: d.dir, host: d.host, lookPath: d.lookPath, selfPath: d.selfPath, selfVersion: displayVersion()}
	if d.verbose {
		inputs.host = verboseInventory{d.host, progress}
	}
	if d.piOps != nil {
		inputs.piStatus = d.piOps.SpacedockPackageStatus
	}
	if d.selfPath != "" {
		inputs.selfVersion = programReportedVersion(d.selfPath, d.env, "spacedock ")
	}
	before := collectVersions(inputs)
	endBefore()
	brew := d.selfPath != "" && isHomebrewPath(d.selfPath)

	owner := "installed with install.sh"
	if brew {
		owner = "installed with Homebrew: " + d.selfPath
	}
	fmt.Fprintf(stdout, "Updating Spacedock on the %s channel (%s)\n", before.channel, owner)

	hosts := planHostUpdates(before)
	ctx := d.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	repair := brewRepair{ctx, d.env, d.dir}
	var pi *piUpdate
	handled := map[string]bool{}
	var brewPlans []brewExtensionPlan
	var brewErrors []error
	code := 0
	if brew {
		endPi := progress.step("Checking which Pi packages are installed")
		pi = planPiUpdate(inputs)
		endPi()
		endPlan := progress.step("Checking release sources for Homebrew extension repairs")
		source, plans := repair.planExtensions(before, pi)
		endPlan()
		brewPlans = plans
		for _, plan := range brewPlans {
			handled[plan.ext.name] = true
			brewErrors = append(brewErrors, repair.applyExtension(plan, source, d.piOps, progress))
		}
	}
	addRefused := map[string]bool{}
	if !brew {
		if err := runInstallScript(d, stdout, stderr, progress); err != nil {
			fmt.Fprintf(stderr, "spacedock update: %v\n", err)
			return 1
		}
	}
	for _, h := range hosts {
		if !h.refresh {
			continue
		}
		// The new install.sh program owns its install sequence, including
		// companions introduced in that release.
		var out string
		endRefresh := progress.step("Refreshing " + h.title + " plugins")
		if brew {
			h.err = repair.refreshCorePlugin(h.name)
		} else if d.runSelf == nil {
			out, h.err = d.host.Install(h.name, channelMarketplaceSource(devBranch), devBranch)
		} else {
			out, h.err = d.runSelf("install", "--host", h.name)
		}
		endRefresh()
		printInstallFailures(stdout, out)
		addRefused[h.title] = strings.Contains(out, "spacedock install: "+h.name+" plugin marketplace add ")
		if h.err != nil && strings.TrimSpace(out) != "" {
			h.err = fmt.Errorf("%w\n%s", h.err, strings.TrimSpace(out))
		}
	}

	if pi == nil {
		endPiCheck := progress.step("Checking which Pi packages are installed")
		pi = planPiUpdate(inputs)
		endPiCheck()
	}
	if pi.refresh {
		endRefresh := progress.step("Refreshing Pi packages")
		var out string
		if brew {
			d.piOps.PiRemove(piSpacedockPackageSource)
			out, pi.err = d.piOps.PiInstall(piSpacedockPackagePin(devBranch))
		} else if d.runSelf == nil {
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
		endRefresh()
	}

	endFinal := progress.step("Checking installed versions after the update")
	if d.selfPath != "" && !brew {
		inputs.selfVersion = programReportedVersion(d.selfPath, d.env, "spacedock ")
	}
	after := collectVersions(inputs)
	_, piMissing := addPiPackages(&after, inputs)
	for i, plan := range brewPlans {
		if !repair.reportExtension(plan, brewErrors[i], stdout, stderr) {
			code = 1
		}
	}
	endFinal()
	for _, h := range hosts {
		h.after = after.hostPlugins[h.name]
	}
	printUpdateChanges(stdout, before, after, before.brewCask, hosts, pi)
	var unresolved []mixedExtension
	for _, m := range after.mixed {
		if !handled[m.name] {
			unresolved = append(unresolved, m)
		}
	}
	for _, line := range mixedAfterUpdate(unresolved, brew, addRefused) {
		fmt.Fprintln(stderr, line)
		code = 1
	}
	if pi.err != nil {
		fmt.Fprintf(stderr, "spacedock update: the Pi package update failed: %v\n", pi.err)
		code = 1
	} else if piMissing {
		fmt.Fprintln(stderr, "spacedock update: the Spacedock package is not installed on Pi. Run: spacedock install --host pi")
		code = 1
	}
	for _, h := range hosts {
		switch v := pieceVersion(h.after, "spacedock plugin"); {
		case !h.checked:
		case h.err != nil:
			fmt.Fprintf(stderr, "spacedock update: the %s plugin update failed: %v\n", h.title, h.err)
			code = 1
		case v == notInstalled || v == "":
			fmt.Fprintf(stderr, "spacedock update: the Spacedock plugin is not installed on %s. Run: spacedock install --host %s\n", h.title, h.name)
			code = 1
		}
	}
	return code
}

// hostUpdate is one host's share of an update. Only a host holding the Spacedock plugin
// is refreshed: update refreshes, install adds. A host holding any other piece from the
// marketplace is still checked.
type hostUpdate struct {
	name, title   string
	refresh       bool
	checked       bool
	skipped       string
	before, after []pieceRow
	err           error
}

func planHostUpdates(before versionReport) []*hostUpdate {
	var updates []*hostUpdate
	for _, status := range before.hosts {
		u := &hostUpdate{name: status.name, title: status.title, skipped: status.skipped,
			checked: status.skipped == "", before: before.hostPlugins[status.name]}
		u.refresh = u.checked && status.lookup(channelPluginID(devBranch)) != notInstalled
		if u.checked && !u.refresh {
			u.skipped = status.name + " has no Spacedock plugin, so its plugins were not updated."
		}
		updates = append(updates, u)
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

func planPiUpdate(in versionInputs) *piUpdate {
	cfg := piRuntimeConfigFromEnv(in.env, in.dir, "")
	u := &piUpdate{agentDir: cfg.agentDir, home: cfg.home}
	rows, skipped := piInventory(in)
	switch {
	case skipped != "":
		u.skipped = skipped
	case rows[0].version == notInstalled:
		u.skipped = "pi has no Spacedock package, so its packages were not updated."
	default:
		u.refresh, u.before = true, rows
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
func runInstallScript(d updateDeps, stdout, stderr io.Writer, progress *updateProgress) error {
	script := envMap(d.env)["SPACEDOCK_INSTALL_SCRIPT"]
	if script == "" {
		tmp, err := os.MkdirTemp("", "spacedock-update-")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmp)
		script = filepath.Join(tmp, "install.sh")
		endFetch := progress.step("Downloading the install script")
		err = d.fetchScript(script)
		endFetch()
		if err != nil {
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
	if d.verbose {
		cmd.WaitDelay = time.Second
	}
	endScript := progress.step("Running the install script")
	err := cmd.Run()
	endScript()
	// ErrWaitDelay means the script succeeded while a descendant held a pipe open.
	if err != nil && !errors.Is(err, exec.ErrWaitDelay) {
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

// mixedPiPackages compares each Pi extension package with its program. Only a
// version-tag ref names a release; an unpinned package or a branch ref (a listing
// may give `stable`) names none, so that package is not compared.
func mixedPiPackages(rows []pieceRow, after versionReport) []mixedExtension {
	var mixed []mixedExtension
	for _, ext := range extensions {
		if ext.piTitle == "" {
			continue
		}
		program := after.version(ext.name + " program")
		pkg := strings.TrimPrefix(pieceVersion(rows, ext.name+" package"), "v")
		_, notVersion := release.ComparePreVersion(pkg, pkg)
		switch {
		case program == notInstalled && after.programsOptedOut:
		case program == "unknown" || (pkg != notInstalled && notVersion != nil):
		case pkg == notInstalled && program == notInstalled:
		case pkg != program:
			mixed = append(mixed, mixedExtension{ext.name, "Pi", program, pkg})
		}
	}
	return mixed
}

func homebrewCask(selfPath string, brew bool) string {
	if !brew {
		return ""
	}
	if isEdgeCaskPath(selfPath) {
		return "spacedock@next"
	}
	return "spacedock"
}

func mixedAfterUpdate(mixed []mixedExtension, brew bool, addRefused map[string]bool) []string {
	var lines []string
	for _, m := range mixed {
		kind := "plugin"
		if m.host == "Pi" {
			kind = "package"
		}
		behind, step, next := "the program", "the install script's download of the "+m.name+" program", "Run `spacedock update` again."
		switch cmp, err := release.ComparePreVersion(m.program, m.plugin); {
		case m.program == notInstalled && brew:
			step = "Homebrew owns the program; see the MISSING line above"
		case m.program == notInstalled:
		case m.plugin == notInstalled && addRefused[m.host]:
			behind, step = "the "+m.host+" "+kind, "the "+m.host+" marketplace add, refused with the message above"
			next = "Change the entry Claude Code names above, then run `spacedock update`."
		case m.plugin == notInstalled || (err == nil && cmp > 0):
			behind, step = "the "+m.host+" "+kind, "the "+m.host+" "+kind+" install"
		case err != nil:
			behind, step = "the program or the "+kind, "a step of this update"
		case brew:
			step = "the Homebrew program is behind the plugin"
			next = "Check the installed formula and the plugin's release."
		}
		lines = append(lines, fmt.Sprintf("spacedock update: %s is still MIXED on %s: program %s, %s %s. Left behind: %s (failed: %s). %s",
			m.name, m.host, m.program, kind, m.plugin, behind, step, next))
	}
	return lines
}

var restartAdvice = map[string]string{
	"claude": "Open Claude Code sessions keep the old plugins until you restart them or run /reload-plugins.",
	"codex":  "Open Codex sessions keep the old plugins until you start a new Codex session.",
	"pi":     "Open Pi sessions keep the old skills until you run /reload.",
}

func printUpdateChanges(w io.Writer, before, after versionReport, cask string, hosts []*hostUpdate, pi *piUpdate) {
	printRow := func(label, b, a string) {
		if b == a {
			fmt.Fprintf(w, "  %-18s %-15s unchanged\n", label, a)
			return
		}
		fmt.Fprintf(w, "  %-18s %s -> %s\n", label, b, a)
	}
	row := func(label string) { printRow(label, before.version(label), after.version(label)) }
	fmt.Fprintln(w, "\nPrograms")
	if cask != "" {
		fmt.Fprintf(w, "  To upgrade Spacedock itself, run: brew upgrade --cask %s\n", cask)
	}
	row("spacedock")
	for _, ext := range extensions {
		row(ext.name + " program")
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
	for _, notice := range after.notices {
		fmt.Fprintln(w, notice)
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
}
