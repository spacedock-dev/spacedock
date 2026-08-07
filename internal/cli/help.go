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
  merge       guard <slug>           Run the terminal merge-finalize ceremony for an entity
  completion  bash|zsh               Print a bash or zsh completion script
  dispatch    build | show-stage-def Build worker dispatch artifacts
  gate        record | validate | consume
                                      Record, inspect, or consume durable gate resolutions

Run "spacedock <command> --help" for details.  ·  --version prints the version.
`

func printHelp(w io.Writer) {
	fmt.Fprint(w, topLevelHelp)
}

// setFrontDoorHelp installs a per-command help renderer for claude/codex (AC-4):
// the sandbox knobs, --skip-compat-check, --plugin-dir, the `--` host-flag
// forwarding explanation, and an Examples block. The flags are declared on the
// command (declareFrontDoorHelpFlags) only so FlagUsages renders them — the real
// parsing is parseFrontDoorArgs's. A per-command HelpFunc is set so the root's
// grouped HelpFunc is not inherited (cobra walks to the parent only when a child
// has none).
func setFrontDoorHelp(cmd *cobra.Command, host string, w io.Writer) {
	declareFrontDoorHelpFlags(cmd)
	cmd.SetHelpFunc(func(c *cobra.Command, _ []string) {
		pluginBehavior := `A valid Spacedock checkout beside the resolved launcher is selected automatically.
Use --plugin-dir before -- to override that selection with another local checkout.
For Claude, --plugin-dir after -- is a native additional session plugin; it does
not replace Spacedock or bypass the installed-plugin compatibility gate.`
		forwardingExamples := `Claude model/session flags and additional --plugin-dir entries.`
		lastExample := `spacedock claude --safehouse-add-dirs ~/scratch -- --plugin-dir ./additional-plugin`
		if host == "codex" {
			pluginBehavior = `A valid Spacedock checkout beside the resolved launcher is selected automatically.
Use --plugin-dir before -- to override that selection with another local checkout.
Codex installs the selected checkout through its persistent local marketplace.
Codex has no forwarded --plugin-dir flag; install additional plugins with
"codex plugin add" instead.`
			forwardingExamples = `Codex model/session flags. A forwarded --plugin-dir is rejected.`
			lastExample = `spacedock codex -- --model gpt-5.6-sol`
		}
		fmt.Fprint(w, tagline+`

Usage:
  spacedock `+host+` [task] [spacedock-flags] [-- `+host+`-flags]

Start `+hostTitle(host)+` as your Spacedock first officer. The optional task is the
launch prompt; supported host arguments after -- forward verbatim to `+host+`.

`+pluginBehavior+`

Flags:
`)
		fmt.Fprint(w, c.Flags().FlagUsages())
		fmt.Fprint(w, `
Forwarding:
  Tokens before -- are spacedock's (the task + the flags above). Tokens after --
  forward verbatim where supported to `+host+`: `+forwardingExamples+`

Examples:
  spacedock `+host+`
  spacedock `+host+` "review the open PRs"
  spacedock `+host+` --plugin-dir ./checkout
  `+lastExample+`
`)
	})
}

// setPiHelp installs the Pi-specific launch help. Pi loads explicit skills and
// extensions instead of a Claude/Codex plugin manifest. The safehouse subset
// mirrors the claude/codex help (same flag names) so operator muscle memory
// transfers across hosts; --skip-compat-check / --no-install are NOT declared
// because pi has no version gate. The flags are declared on the command only so
// FlagUsages renders them (the pi command is DisableFlagParsing, exactly like
// claude/codex); parsePiFrontDoorArgs owns the real parsing.
func setPiHelp(cmd *cobra.Command, w io.Writer) {
	cmd.Flags().StringArray("plugin-dir", nil, "Load a local Spacedock skill checkout")
	cmd.Flags().Bool("safehouse", false, "Force the safehouse sandbox wrap even without a .safehouse profile in the directory")
	cmd.Flags().StringArray("safehouse-enable", nil, "Enable a safehouse capability (KEY[,KEY]); repeatable; e.g. --safehouse-enable ssh,docker")
	cmd.Flags().StringArray("safehouse-add-dirs", nil, "Grant safehouse read-write access to a directory; repeatable")
	cmd.Flags().StringArray("safehouse-add-dirs-ro", nil, "Grant safehouse read-only access to a directory; repeatable")
	cmd.SetHelpFunc(func(c *cobra.Command, _ []string) {
		fmt.Fprint(w, tagline+`

Usage:
  spacedock pi [task] [--plugin-dir <checkout>] [--safehouse-* knobs] [-- pi-flags]

Start Pi as your Spacedock first officer by loading the Pi-native pi-subagents
extension/skill. The Spacedock first-officer/ensign skills are discovered from
the installed Spacedock package (run: spacedock install --host pi); the optional
--plugin-dir loads a local checkout as a dev override. The optional task is
appended to the launch prompt; everything after -- forwards verbatim to pi.

When any --safehouse-* knob (or a .safehouse profile, or the bare --safehouse) is
present, the launch is wrapped as `+"`"+`safehouse --trust-workdir-config [extra] -- pi …`+"`"+`
so a sandboxed session can grant additional directory access. No inner
permission flag is added — pi's posture is its --tools/--exclude-tools passthrough
and the safehouse isolation itself.

Flags:
`)
		fmt.Fprint(w, c.Flags().FlagUsages())
		fmt.Fprint(w, `
Forwarding:
  Tokens before -- are spacedock's (the task + --plugin-dir + the --safehouse-*
  knobs). Tokens after -- forward verbatim to pi, e.g. --model, --print, or
  --session-dir.

Examples:
  spacedock pi --plugin-dir ./checkout
  spacedock pi "drive the workflow" --plugin-dir ./checkout -- --model google/gemini
  spacedock pi --safehouse-add-dirs ~/scratch -- --model google/gemini
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

// setNewHelp installs a per-command help renderer for `new`: its own synopsis and
// flag surface (--workflow-dir, --folder, --id-seed, --id-actor) instead of the
// root's grouped menu (cobra walks to the parent HelpFunc only when a child has
// none).
func setNewHelp(cmd *cobra.Command, w io.Writer) {
	cmd.SetHelpFunc(func(c *cobra.Command, _ []string) {
		fmt.Fprint(w, c.Short+`

Usage:
  spacedock new [--folder] SLUG < body

Create an entity from a stdin body. Run from the project root: new auto-discovers
the lone commissioned workflow. If the repo holds more than one, new reports the
candidates and you pass --workflow-dir to pick one.

Flags:
  --workflow-dir DIR   Target this workflow explicitly (skips auto-discovery)
  --folder             Write the entity in folder form (SLUG/index.md)
  --id-seed SEED       Seed material for the minted id (id-style: sd-b32 only)
  --id-actor ACTOR     Actor material for the minted id (id-style: sd-b32 only)

Examples:
  spacedock new my-task < body.md
  spacedock new --folder my-task < body.md
  spacedock new my-task --workflow-dir docs/dev < body.md
`)
	})
}

// setStatusHelp installs a per-command help renderer for `status`: the query
// flag surface instead of the root's grouped menu (cobra walks to the parent
// HelpFunc only when a child has none). It names --where as THE entity query,
// the one-clause-per-flag AND rule (a repeated-clause string is an error, not an
// AND — the #314 fix), the canonical known-field list, and --archived's
// active-plus-archived semantics (it composes with --where; it does not swap
// scope) — the discoverability surface `status --help` did not previously
// provide (before this, RunE had no wantsHelp guard, so --help fell through to
// the entity listing).
func setStatusHelp(cmd *cobra.Command, w io.Writer) {
	cmd.SetHelpFunc(func(c *cobra.Command, _ []string) {
		fmt.Fprint(w, c.Short+`

Usage:
  spacedock status --workflow-dir DIR [query flags]

Query flags:
  --where FIELD=VALUE  Filter entities — THE entity query; repeat the flag to AND clauses
  --archived           Include archived entities (active PLUS archived, not archived-only)
  --fields a,b,c       Project to these fields    --all-fields  Show every stored field
  --next               Dispatchable entities      --boot        Startup roll-up
  --resolve REF        Resolve slug/id/prefix     --next-id     Preview the next id
  --validate           Check workflow state       --json        Machine-readable output

Operators (one clause per flag):
  field=value   equals             field!=value  not equals
  field=        field is empty     field!=       field is non-empty

Repeat --where to AND clauses. Two clauses in one string is an error, not an AND.

Known fields are this workflow's entity frontmatter keys plus the canonical set
(id, slug, status, title, score, source, worktree, pr, started, completed,
verdict, mod-block, archived, issue). An unknown field is an error that lists them.

Examples:
  spacedock status --workflow-dir docs/dev --where status=ideation
  spacedock status --workflow-dir docs/dev --where sprint=X --where 'sprint-readiness!=defer'
  spacedock status --workflow-dir docs/dev --where sprint=X --archived
  spacedock status --workflow-dir docs/dev --where sprint=X --archived --fields slug,status,verdict,archived
`)
	})
}

// setMergeHelp installs a per-command help renderer for `merge`: the `guard`
// subcommand synopsis, its flag surface, and the three-phase ceremony it drives,
// instead of the root's grouped menu (cobra walks to the parent HelpFunc only when
// a child has none).
func setMergeHelp(cmd *cobra.Command, w io.Writer) {
	cmd.SetHelpFunc(func(c *cobra.Command, _ []string) {
		fmt.Fprint(w, c.Short+`

Usage:
  spacedock merge guard <slug> --verdict passed|rejected [--workflow-dir DIR]
  spacedock merge guard <slug> --rework [--workflow-dir DIR]

Drive the terminal merge ceremony for an entity as one ordered envelope: arm the
mod-block before the hook, detect hook completion by the state delta, then
terminalize, archive with a path-scoped commit, and publish split-root state.
The verb owns the sequence so the steps cannot be combined, skipped, or reordered;
it does NOT invoke the merge hook or make the merge verdict (you pass that in with
--verdict).

merge guard is the sole terminal consumer of a gate approval whose target is the
terminal stage: "gate consume" leaves such an approval pending (route
approved-awaiting-merge), and guard's successful delivery spends it — the
mod-block is cleared in its own step first, then application, terminal status,
verdict, and completed move in ONE locked write; the pr merge sentinel is
retained through archive as durable delivery proof. A non-forced terminal
"status --set" while that approval is pending is refused in favor of this
ceremony.

--rework is the delivery-requires-rework outcome: when delivery fails beyond
retry (e.g. the PR is closed unmerged), it writes the approval pending->superseded,
routes the entity back through the record stage's declared feedback-to, and clears
delivery state (mod-block/pr). It refuses without a pending terminal approval, or
when no defined non-terminal feedback-to is declared. Retryable delivery trouble
(CI red, push failure) takes neither flag: guard without delivery proof writes
nothing and reports armed/blocked, and the approval stays pending.

Re-run guard after invoking the hook: it resumes from the entity's current state
(armed -> blocked on an open PR, or armed -> finalized once the merge has landed).
After interrupted publication, spacedock state commit <slug> resumes the archived slug without creating another archive commit.
Exit 3 means Git rebase conflict: the rebase was aborted; stop and surface its path/peer evidence instead of rerunning guard.

Flags:
  --verdict passed|rejected   The merge decision (required unless --rework; a
                              verdict-less finalize is refused)
  --rework                    Delivery requires rework: supersede the pending terminal
                              approval and route back through the declared feedback-to
  --workflow-dir DIR          Target this workflow explicitly (skips auto-discovery).
                              A relative DIR resolves against the current directory;
                              from anywhere else (e.g. an agent worktree) pass an
                              absolute path.
  --json                      Emit the phase signal as JSON
  --quiet                     Emit a terse machine-readable phase signal

Examples:
  spacedock merge guard my-task --verdict passed
  spacedock merge guard my-task --verdict rejected --workflow-dir docs/dev
  spacedock merge guard my-task --rework
`)
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
