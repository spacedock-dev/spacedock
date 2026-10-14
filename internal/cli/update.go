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
			deps := updateDeps{env: env, dir: dir, host: execHost{}, selfPath: selfPath, lookPath: exec.LookPath,
				fetchScript: func(dst string) error {
					return exec.CommandContext(ctx, "curl", "-fsSL", "-o", dst, installScriptURL).Run()
				}}
			if code := runUpdate(deps, stdout, stderr); code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
	setSetupHelp(cmd, stdout, `
Update every program and plugin on the channel of this binary. Homebrew owns its
own programs: update prints the brew command and updates only the plugins.
SPACEDOCK_INSTALL_SCRIPT=PATH runs that script in place of the downloaded one.
`)
	return cmd
}

func runUpdate(d updateDeps, stdout, stderr io.Writer) int {
	inputs := versionInputs{env: d.env, dir: d.dir, host: d.host, selfPath: d.selfPath, selfVersion: displayVersion()}
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

	_, claudeErr := d.lookPath("claude")
	switch {
	case brew:
		if claudeErr == nil {
			out, err := d.host.Install("claude", channelMarketplaceSource(devBranch), devBranch)
			printInstallFailures(stdout, out)
			if err != nil {
				fmt.Fprintf(stderr, "spacedock update: plugin update failed: %v\n", err)
				return 1
			}
		}
	default:
		if err := runInstallScript(d, stdout, stderr); err != nil {
			fmt.Fprintf(stderr, "spacedock update: %v\n", err)
			return 1
		}
	}

	if d.selfPath != "" && !brew {
		inputs.selfVersion = programReportedVersion(d.selfPath, d.env, "spacedock ")
	}
	after := collectVersions(inputs)
	printUpdateChanges(stdout, before, after, brewFormula(d.selfPath, brew), claudeErr != nil)
	if claudeErr == nil && after.version("spacedock plugin") == notInstalled {
		fmt.Fprintln(stderr, "spacedock update: the Spacedock plugin is not installed. Run: spacedock install --host claude")
		return 1
	}
	return 0
}

// runInstallScript installs into the directory of this binary, so the new program
// replaces it. The script's stdout is shown only when it fails.
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
	env := append(append([]string{}, d.env...), "SPACEDOCK_CHANNEL="+selectedChannelWord())
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

func printUpdateChanges(w io.Writer, before, after versionReport, formula string, noClaude bool) {
	row := func(label string) {
		b, a := before.version(label), after.version(label)
		if b == a {
			fmt.Fprintf(w, "  %-18s %-15s unchanged\n", label, a)
			return
		}
		fmt.Fprintf(w, "  %-18s %s -> %s\n", label, b, a)
	}
	fmt.Fprintln(w, "\nPrograms")
	if formula != "" {
		fmt.Fprintf(w, "  Homebrew owns spacedock and its extensions. Run: brew upgrade %s\n", formula)
	} else {
		row("spacedock")
		for _, ext := range extensions {
			row(ext.name + " program")
		}
	}
	fmt.Fprintf(w, "Claude Code plugins (%s)\n", channelMarketplace(devBranch))
	if noClaude {
		fmt.Fprintln(w, "  claude is not on PATH, so the plugins were skipped. Run: spacedock install --host claude")
	} else {
		row("spacedock plugin")
		for _, ext := range extensions {
			row(ext.name + " plugin")
		}
	}
	if len(after.problems) > 0 {
		fmt.Fprintln(w)
		for _, p := range after.problems {
			fmt.Fprintln(w, p)
		}
	}
	fmt.Fprintln(w, "\nOpen Claude Code sessions keep the old plugins until you restart them or run /reload-plugins.")
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
