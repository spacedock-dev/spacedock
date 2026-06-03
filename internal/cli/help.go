// ABOUTME: Grouped jargon-free top-level help and the per-command help renderers.
// ABOUTME: Renders the Launch/Setup/Workflow groups, terse one-liners, and footer.
package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// topLevelHelp is the grouped, jargon-free, META-free top-level help (AC-1). The
// command list is hand-rendered (not cobra's default `Available Commands` block)
// so the three groups carry terse aligned one-liners and the footer holds
// `--version` and the per-command help pointer instead of dedicated command rows.
// No internal jargon (`front door`, `contract-gated`, `META`) appears.
const topLevelHelp = tagline + `

Launch
  claude  [task] [-- claude-flags]   Start Claude Code as your Spacedock first officer
  codex   [task] [-- codex-flags]    Start Codex as your Spacedock first officer
  pi      [task] [-- pi-flags]       Start Pi as your Spacedock first officer
Setup
  install  [--host claude|codex|pi]  Install the Spacedock plugin for a host, then check it
  doctor   [--host claude|codex|pi]  Check the installed plugin and this binary are compatible
Workflow
  status      [args]                 Show or update workflow state
  new         [--folder] SLUG        Create an entity from a stdin body (auto-discovers the workflow)
  state       init                   Initialize a cloned split-root workflow's state checkout
  completion  bash|zsh               Print a bash or zsh completion script
  dispatch    build | show-stage-def Build worker dispatch artifacts

Run "spacedock <command> --help" for details.  ·  --version prints the version.
`

func printHelp(w io.Writer) {
	fmt.Fprint(w, topLevelHelp)
}

// setFrontDoorHelp installs a per-command help renderer for claude/codex (AC-4):
// the sandbox knobs, --skip-contract-check, --plugin-dir, the `--` host-flag
// forwarding explanation, and an Examples block. The flags are declared on the
// command (declareFrontDoorHelpFlags) only so FlagUsages renders them — the real
// parsing is parseFrontDoorArgs's. A per-command HelpFunc is set so the root's
// grouped HelpFunc is not inherited (cobra walks to the parent only when a child
// has none).
func setFrontDoorHelp(cmd *cobra.Command, host string, w io.Writer) {
	declareFrontDoorHelpFlags(cmd)
	cmd.SetHelpFunc(func(c *cobra.Command, _ []string) {
		fmt.Fprint(w, tagline+`

Usage:
  spacedock `+host+` [task] [spacedock-flags] [-- `+host+`-flags]

Start `+hostTitle(host)+` as your Spacedock first officer. The optional task is the
launch prompt; everything after -- forwards verbatim to `+host+`.

A --plugin-dir launch loads a local plugin checkout and relaxes the contract gate,
so it does not require a prior "spacedock install". --plugin-dir is accepted both
before -- (as a spacedock-parsed flag, repeatable) and after -- (forwarded verbatim).

Flags:
`)
		fmt.Fprint(w, c.Flags().FlagUsages())
		fmt.Fprint(w, `
Forwarding:
  Tokens before -- are spacedock's (the task + the flags above). Tokens after --
  forward verbatim to `+host+`, e.g. `+host+` model/session flags and --plugin-dir.

Examples:
  spacedock `+host+`
  spacedock `+host+` "review the open PRs"
  spacedock `+host+` --plugin-dir ./checkout
  spacedock `+host+` --safehouse-add-dirs ~/scratch -- --plugin-dir ./checkout
`)
	})
}

// setPiHelp installs the Pi-specific launch help. Pi loads explicit skills and
// extensions instead of a Claude/Codex plugin manifest.
func setPiHelp(cmd *cobra.Command, w io.Writer) {
	cmd.Flags().String("plugin-dir", "", "Load a local Spacedock skill checkout")
	cmd.SetHelpFunc(func(c *cobra.Command, _ []string) {
		fmt.Fprint(w, tagline+`

Usage:
  spacedock pi [task] [--plugin-dir <checkout>] [-- pi-flags]

Start Pi as your Spacedock first officer by loading the Pi-native pi-subagents
extension/skill and the Spacedock first-officer/ensign skills. The optional task
is appended to the launch prompt; everything after -- forwards verbatim to pi.

Flags:
`)
		fmt.Fprint(w, c.Flags().FlagUsages())
		fmt.Fprint(w, `
Forwarding:
  Tokens before -- are spacedock's (the task + --plugin-dir). Tokens after --
  forward verbatim to pi, e.g. --model, --print, or --session-dir.

Examples:
  spacedock pi --plugin-dir ./checkout
  spacedock pi "drive the workflow" --plugin-dir ./checkout -- --model google/gemini
`)
	})
}

// setSetupHelp installs a per-command help renderer for install/doctor: the
// command's own flags and an Examples block. A per-command HelpFunc is set so the
// root's grouped HelpFunc is not inherited.
func setSetupHelp(cmd *cobra.Command, w io.Writer, examples string) {
	cmd.SetHelpFunc(func(c *cobra.Command, _ []string) {
		fmt.Fprintf(w, "%s\n\nUsage:\n  spacedock %s\n\nFlags:\n", c.Short, c.Use)
		fmt.Fprint(w, c.Flags().FlagUsages())
		fmt.Fprint(w, examples)
	})
}

// declareFrontDoorHelpFlags registers the spacedock-owned front-door flags onto a
// command's flag set purely so `--help` renders them (AC-4). The flags are never
// parsed by cobra here (the command has DisableFlagParsing); parseFrontDoorArgs
// owns the real parsing. bindFrontDoorFlags is the single source for the flag set
// so the help and the parser never drift.
func declareFrontDoorHelpFlags(cmd *cobra.Command) {
	bindFrontDoorFlags(cmd.Flags())
}

// hostTitle returns the display name for the host in help prose.
func hostTitle(host string) string {
	if host == "codex" {
		return "Codex"
	}
	return "Claude Code"
}
