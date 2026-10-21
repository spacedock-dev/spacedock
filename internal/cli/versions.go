// ABOUTME: The version of every installed piece (Spacedock, its plugin, each extension
// ABOUTME: program and plugin), read by `doctor` and `update`.
package cli

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/spacedock-dev/spacedock/internal/release"
)

// extension is one row of the extension table. program is the `spacedock-NAME` link on
// PATH. A tree extension answers no `--version`, so its version is the directory that
// its link resolves into. check adds the rows and problems that only this extension
// knows how to find, and afterUpdate is printed when its program changed in an update.
type extension struct {
	name        string
	program     string
	tree        bool
	check       func(in versionInputs) (rows []pieceRow, problems []string)
	afterUpdate []string
}

// extensions lists the extensions Spacedock installs. install.sh holds the matching
// EXTRAS rows, and the two must name the same programs.
var extensions = []extension{}

// channelExtensions is the table as this channel applies it: stable installs no
// extension unless SPACEDOCK_EXTRAS asks, so it expects and updates none.
func channelExtensions(devBranch string) []extension {
	if devBranch == "main" {
		return nil
	}
	return extensions
}

const notInstalled = "not installed"

type pieceRow struct{ label, version, where string }

type versionReport struct {
	channel                string
	rows                   []pieceRow
	hosts                  []hostStatus
	hostPlugins            map[string][]pieceRow
	problems               []string
	needsUpdate, extraCopy bool
	pluginsBehind          []string
}

type versionInputs struct {
	env         []string
	dir         string
	host        hostOps
	lookPath    func(string) (string, error)
	selfPath    string
	selfVersion string
}

// hostStatus is what one plugin host holds: its plugin versions by id, or why it has
// none to show. skipped is set for a host that is absent, holds no Spacedock plugin, or
// fails to list its plugins; only the last still shows rows, at "unknown", because the
// host is there and its state is not.
type hostStatus struct {
	name, title string
	skipped     string
	err         error
	lookup      func(id string) string
}

func (h hostStatus) listed() bool { return h.skipped == "" || h.err != nil }

func readHosts(in versionInputs) []hostStatus {
	var hosts []hostStatus
	for _, h := range []struct{ name, title string }{{"claude", "Claude Code"}, {"codex", "Codex"}} {
		status := hostStatus{name: h.name, title: h.title}
		if _, err := in.lookPath(h.name); err != nil {
			status.skipped = h.name + " is not on PATH, so its plugins were skipped."
			hosts = append(hosts, status)
			continue
		}
		status.lookup, status.err = pluginVersionLookup(in.host, h.name)
		switch {
		case status.err != nil:
			status.skipped = fmt.Sprintf("%s plugin list failed (%v), so its plugins were skipped.", h.name, status.err)
		case status.lookup(channelPluginID(devBranch)) == notInstalled:
			status.skipped = h.name + " has no Spacedock plugin, so its plugins were skipped."
		}
		hosts = append(hosts, status)
	}
	return hosts
}

func (r versionReport) version(label string) string {
	for _, row := range r.rows {
		if row.label == label {
			return row.version
		}
	}
	return ""
}

func collectVersions(in versionInputs) versionReport {
	report := versionReport{channel: selectedChannelWord(), hosts: readHosts(in), hostPlugins: map[string][]pieceRow{}}
	add := func(label, version, where string) {
		report.rows = append(report.rows, pieceRow{label, version, where})
	}
	addPlugin := func(h hostStatus, label, version, where string) {
		report.hostPlugins[h.name] = append(report.hostPlugins[h.name], pieceRow{label, version, where})
		add(label+" ("+h.title+")", version, where)
	}
	problem := func(format string, a ...any) {
		report.problems = append(report.problems, fmt.Sprintf(format, a...))
	}
	pathEnv := envMap(in.env)["PATH"]

	add("spacedock", in.selfVersion, in.selfPath)
	spacedockID := channelPluginID(devBranch)
	for _, h := range report.hosts {
		if h.listed() {
			addPlugin(h, "spacedock plugin", h.lookup(spacedockID), spacedockID)
		}
	}

	for _, ext := range channelExtensions(devBranch) {
		programPath, programVersion := extensionProgram(ext, in.env, pathEnv)
		add(ext.name+" program", programVersion, programPath)
		id := extensionPluginID(ext.name)
		if programVersion == notInstalled {
			problem("MISSING: the %s program is not on PATH.", ext.name)
			report.needsUpdate = true
		}
		for _, h := range report.hosts {
			if !h.listed() {
				continue
			}
			pluginVer := h.lookup(id)
			addPlugin(h, ext.name+" plugin", pluginVer, id)
			switch {
			case programVersion == notInstalled:
			case h.err == nil && pluginVer == notInstalled:
				problem("MISSING: the %s plugin is not installed on %s from %s.", ext.name, h.title, extensionMarketplace)
				report.needsUpdate = true
			case programVersion != "unknown" && pluginVer != "unknown" && programVersion != pluginVer:
				problem("MIXED: the %s program %s does not match its plugin %s on %s.", ext.name, programVersion, pluginVer, h.title)
				if cmp, err := release.ComparePreVersion(programVersion, pluginVer); err == nil && cmp > 0 {
					if !slices.Contains(report.pluginsBehind, ext.name) {
						report.pluginsBehind = append(report.pluginsBehind, ext.name)
					}
				} else {
					report.needsUpdate = true
				}
			}
		}
	}

	for _, ext := range channelExtensions(devBranch) {
		if ext.check == nil {
			continue
		}
		rows, problems := ext.check(in)
		report.rows = append(report.rows, rows...)
		report.problems = append(report.problems, problems...)
	}

	if copies := spacedockCopies(pathEnv); len(copies) > 1 {
		lines := []string{fmt.Sprintf("MIXED: %d copies of spacedock on PATH. The first one runs:", len(copies))}
		for i, c := range copies {
			suffix := ""
			switch {
			case i == 0:
				suffix = "  (runs)"
			case isHomebrewPath(c):
				suffix = "  Homebrew"
			}
			lines = append(lines, fmt.Sprintf("  %s    %s%s", c, programReportedVersion(c, in.env, "spacedock "), suffix))
		}
		report.problems = append(report.problems, strings.Join(lines, "\n"))
		report.extraCopy = true
	}
	return report
}

func pluginVersionLookup(h hostOps, host string) (lookup func(id string) string, err error) {
	inventory, err := h.PluginInventory(host)
	return func(id string) string {
		if err != nil {
			return "unknown"
		}
		if entry, ok := installedPlugin(inventory, id); ok {
			return entry.Version
		}
		return notInstalled
	}, err
}

func printVersionReport(w io.Writer, report versionReport) {
	fmt.Fprintf(w, "\nVersions (%s channel)\n", report.channel)
	width := 18
	for _, row := range report.rows {
		width = max(width, len(row.label))
	}
	for _, row := range report.rows {
		fmt.Fprintf(w, "  %-*s %-15s %s\n", width, row.label, row.version, row.where)
	}
	if len(report.problems) == 0 {
		fmt.Fprintln(w, "\nNo problems found.")
		return
	}
	fmt.Fprintln(w)
	for _, p := range report.problems {
		fmt.Fprintln(w, p)
	}
	fmt.Fprintf(w, "\nProblems found: %d.", len(report.problems))
	if report.needsUpdate {
		fmt.Fprint(w, " Run `spacedock update` to bring every program and plugin to the same release.")
	}
	if len(report.pluginsBehind) > 0 {
		fmt.Fprintf(w, " The %s plugin waits on its listing in the %s marketplace, which each extension's release keeps current. Run `spacedock update` after the listing shows the new version.",
			strings.Join(report.pluginsBehind, " and "), extensionMarketplace)
	}
	if report.extraCopy {
		fmt.Fprint(w, " Remove the copies of spacedock that you do not want to run.")
	}
	fmt.Fprintln(w)
}

func extensionProgram(ext extension, env []string, pathEnv string) (path, version string) {
	found := pathCopies(pathEnv, ext.program)
	if len(found) == 0 {
		return "", notInstalled
	}
	if ext.tree {
		return found[0], treeVersion(found[0], ext.name, env)
	}
	return found[0], programReportedVersion(found[0], env, "")
}

func treeVersion(link, name string, env []string) string {
	vars := envMap(env)
	root := vars["XDG_DATA_HOME"]
	if root == "" {
		root = filepath.Join(vars["HOME"], ".local", "share")
	}
	root = filepath.Join(root, name)
	if resolved, err := filepath.EvalSymlinks(root); err == nil {
		root = resolved
	}
	target, err := filepath.EvalSymlinks(link)
	if err != nil {
		return "unknown"
	}
	rel, err := filepath.Rel(root, target)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "unknown"
	}
	return strings.Split(rel, string(filepath.Separator))[0]
}

func programReportedVersion(path string, env []string, prefix string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, "--version")
	cmd.Env = env
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	line, _, _ := strings.Cut(strings.TrimSpace(string(out)), "\n")
	line = strings.TrimSpace(strings.TrimPrefix(line, prefix))
	if line == "" {
		return "unknown"
	}
	return line
}

func pathCopies(pathEnv, name string) []string {
	var found []string
	for _, dir := range filepath.SplitList(pathEnv) {
		if !filepath.IsAbs(dir) {
			continue
		}
		candidate := filepath.Join(dir, name)
		if executableFile(candidate) {
			found = append(found, candidate)
		}
	}
	return found
}

func spacedockCopies(pathEnv string) []string {
	var copies []string
	seen := map[string]bool{}
	for _, c := range pathCopies(pathEnv, "spacedock") {
		real := c
		if resolved, err := filepath.EvalSymlinks(c); err == nil {
			real = resolved
		}
		if !seen[real] {
			seen[real] = true
			copies = append(copies, c)
		}
	}
	return copies
}

func isHomebrewPath(path string) bool {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		path = resolved
	}
	for _, seg := range strings.Split(filepath.ToSlash(path), "/") {
		if seg == "Caskroom" || seg == "Cellar" {
			return true
		}
	}
	return false
}

func versionsWanted(args []string) bool {
	for i, a := range args {
		if a == "--plugin-manifest" || a == "--plugin-dir" || (a == "--host" && i+1 < len(args) && args[i+1] != "claude" && args[i+1] != "codex") {
			return false
		}
	}
	return true
}

func reportVersions(args, env []string, dir string, stdout io.Writer) {
	if !versionsWanted(args) {
		return
	}
	selfPath, _ := resolvedLauncherBin()
	printVersionReport(stdout, collectVersions(versionInputs{env: env, dir: dir, host: execHost{}, lookPath: exec.LookPath, selfPath: selfPath, selfVersion: displayVersion()}))
}
