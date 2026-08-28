// ABOUTME: spacedock install/doctor command paths — install the per-host plugin via
// ABOUTME: the host plugin mechanism (claude) or emit the documented add prose (codex).
package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// defaultMarketplaceSource is the production marketplace add source: the standalone
// marketplace repo (NOT the plugin repo). The channel selects a branch of it via
// channelMarketplaceSource — stable installs from the bare repo (root marketplace.json
// named `spacedock`), edge from `@edge` (root marketplace.json named `spacedock-edge`)
// — each a marketplace.json with the single `spacedock` entry. The version pin lives in
// the manifest, not an @ref on the install command.
const defaultMarketplaceSource = "spacedock-dev/marketplace"

// marketplaceSource is the marketplace add source. It is a var (not a const) so
// `SPACEDOCK_MARKETPLACE_SOURCE` can override it — pointing the install at a
// local/alternate marketplace to dogfood a channel fix before it reaches the
// production marketplace; tests save/restore it. When overridden it replaces the
// resolved source verbatim (channelMarketplaceSource appends no channel ref to an
// operator's chosen source); unset, it is the production default the channel branches.
var marketplaceSource = defaultMarketplaceSource

// runInit installs/updates the per-host plugin (claude or codex) then runs doctor.
// `--check` runs the report without installing. Both hosts install programmatically
// through the host plugin mechanism (claude `install` / codex `add`) — no skill-file
// copies — which is what makes Skill()/--agent spacedock:first-officer resolve.
func runInit(ctx context.Context, args []string, ops hostOps, stdout, stderr io.Writer) int {
	host, check, code := parseInitArgs(args, stderr)
	if code != 0 {
		return code
	}

	switch host {
	case "claude":
		if !check {
			out, err := ops.Install(host, channelMarketplaceSource(devBranch), devBranch)
			if err != nil {
				fmt.Fprintf(stderr, "spacedock install: host install failed: %v\n", err)
				return 1
			}
			if out != "" {
				fmt.Fprintln(stdout, out)
			}
		}
		return runDoctor(ctx, []string{"--host", "claude"}, ops, stdout, stderr)
	case "codex":
		if check {
			return runDoctor(ctx, []string{"--host", "codex"}, ops, stdout, stderr)
		}
		// `install --host codex` drives the install seam (marketplace add + plugin
		// add, re-pinning the source) and runs doctor — the same programmatic path
		// as the claude arm, whether or not a plugin is already present. A fresh box
		// installs; a present plugin is refreshed. The install verb is codex's `add`
		// (supplied by codexInstallArgvSequence), and the channel marketplace the
		// binary's devBranch selects is the install target.
		out, err := ops.Install("codex", channelMarketplaceSource(devBranch), devBranch)
		if err != nil {
			fmt.Fprintf(stderr, "spacedock install: host install failed: %v\n", err)
			return 1
		}
		if out != "" {
			fmt.Fprintln(stdout, out)
		}
		return runDoctor(ctx, []string{"--host", "codex"}, ops, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "spacedock install: unknown host %q (want claude or codex)\n", host)
		return 2
	}
}

// runDoctor is the `spacedock doctor` command path. With `--plugin-manifest PATH`
// it reads that manifest directly (used by fixtures and operators); otherwise it
// resolves the installed manifest via the host-ops seam. A resolved compatible
// manifest exits 0; a mismatch exits 1; no installed plugin is a non-fatal
// report (exit 0).
func runDoctor(ctx context.Context, args []string, ops hostOps, stdout, stderr io.Writer) int {
	manifestPath, host, code := parseDoctorArgs(args, stderr)
	if code != 0 {
		return code
	}

	if manifestPath != "" {
		return contract.RunDoctor(manifestPath, host, displayVersion(), runningEdgeCask(), stdout, stderr)
	}

	resolved, err := ops.ResolveManifest(host)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock doctor: could not resolve the installed %s plugin: %v\n", host, err)
		return 1
	}
	// An empty resolved path is the no-plugin-found report; RunDoctor renders it
	// from a non-existent path as a non-fatal report. Inventory is deliberately
	// queried afterward, and its report never changes the compatibility exit code.
	compatCode := contract.RunDoctor(resolved, host, displayVersion(), runningEdgeCask(), stdout, stderr)
	inventory, err := ops.PluginInventory(host)
	if err != nil {
		fmt.Fprintf(stdout, "INCOMPLETE: doctor checked compatibility but did not read the %s plugin enablement state: %v\n", host, err)
		return compatCode
	}
	printSiblingConflict(host, inventory, stdout)
	return compatCode
}

func printSiblingConflict(host string, inventory []pluginInventoryEntry, stdout io.Writer) {
	selectedID := channelPluginID(devBranch)
	selected, selectedOK := installedPlugin(inventory, selectedID)
	sibling, siblingOK := enabledSiblingPlugin(inventory)
	if !selectedOK || !siblingOK {
		return
	}
	selectedState := "disabled"
	if selected.Enabled {
		selectedState = "enabled"
	}

	fmt.Fprintf(stdout, "CONFLICT: %s can load a different Spacedock plugin than doctor checked.\n", host)
	fmt.Fprintf(stdout, "  checked: %s %s (installed, %s)\n", selected.ID, selected.Version, selectedState)
	fmt.Fprintf(stdout, "  sibling: %s %s (installed, enabled)\n", sibling.ID, sibling.Version)
	fmt.Fprintf(stdout, "Run `spacedock install --host %s` to keep only the %s channel.\n", host, selectedChannelWord())
}

func selectedChannelWord() string {
	if devBranch == "main" {
		return "stable"
	}
	return "edge"
}

func enabledSiblingPlugin(inventory []pluginInventoryEntry) (pluginInventoryEntry, bool) {
	siblingBranch := "main"
	if devBranch == "main" {
		siblingBranch = "next"
	}
	sibling, ok := installedPlugin(inventory, channelPluginID(siblingBranch))
	return sibling, ok && sibling.Enabled
}

func installedPlugin(inventory []pluginInventoryEntry, id string) (pluginInventoryEntry, bool) {
	for _, entry := range inventory {
		if entry.ID == id && entry.Installed {
			return entry, true
		}
	}
	return pluginInventoryEntry{}, false
}

// parseInitArgs reads `--host claude|codex` (default claude) and `--check`. A
// missing --host value is a usage error (exit 2).
func parseInitArgs(args []string, stderr io.Writer) (host string, check bool, code int) {
	host = "claude"
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--host":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "spacedock install: --host requires a value (claude or codex)")
				return "", false, 2
			}
			host = args[i+1]
			i++
		case "--check":
			check = true
		default:
			fmt.Fprintf(stderr, "spacedock install: unknown argument %q\n", args[i])
			return "", false, 2
		}
	}
	return host, check, 0
}

// parseDoctorArgs reads `--plugin-manifest PATH` (optional explicit manifest)
// and `--host claude|codex` (default claude).
func parseDoctorArgs(args []string, stderr io.Writer) (manifestPath, host string, code int) {
	host = "claude"
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--plugin-manifest":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "spacedock doctor: --plugin-manifest requires a path")
				return "", "", 2
			}
			manifestPath = args[i+1]
			i++
		case "--host":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "spacedock doctor: --host requires a value (claude or codex)")
				return "", "", 2
			}
			host = args[i+1]
			i++
		default:
			fmt.Fprintf(stderr, "spacedock doctor: unknown argument %q\n", args[i])
			return "", "", 2
		}
	}
	return manifestPath, host, 0
}
