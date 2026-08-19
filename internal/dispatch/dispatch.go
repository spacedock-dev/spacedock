// ABOUTME: dispatch command router — build + show-stage-def, the native surface
// ABOUTME: of the self-hosted dispatch path; everything else is a usage error.
package dispatch

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
)

// Run routes launcher-independent dispatch subcommands. Artifact builds must use
// RunWithLauncher so they can bind the CLI-resolved executable and fail closed.
func Run(probe claudeteam.TeamStateProbe, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return RunWithLauncher(probe, "", args, stdin, stdout, stderr)
}

// RunWithLauncher routes a `spacedock dispatch <subcommand> [flags]` invocation. build and
// show-stage-def are the host-neutral surface (assembled here); context-budget,
// list-standing, show-standing, and spawn-standing-all are the Claude-coupled surface
// (their ~/.claude and standing-mod reads live in internal/claudeteam). An unknown
// subcommand fails with exit 2 and a usage diagnostic on stderr. probe is the
// host-supplied team-state probe gating the bare-mode advisory (nil on a non-Claude
// host → no advisory).
func RunWithLauncher(probe claudeteam.TeamStateProbe, workflowLauncher string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	switch args[0] {
	case "build":
		if wantsHelp(args[1:]) {
			printBuildUsage(stdout)
			return 0
		}
		opts, code := parseBuildOptions(args[1:], stderr)
		if code != 0 {
			return code
		}
		if !opts.PrintSchema && opts.ValidateOnly == "" && opts.WorkflowDir == "" {
			fmt.Fprintln(stderr, "error: dispatch build requires --workflow-dir")
			return 2
		}
		if opts.Advance && opts.BareMode {
			fmt.Fprintln(stderr, "error: dispatch build --advance is incompatible with --bare-mode (a reuse advance presupposes an addressable worker; bare mode has none)")
			return 2
		}
		if opts.Stamp && opts.Advance {
			fmt.Fprintln(stderr, "error: dispatch build --stamp is incompatible with --advance (a reuse advance presupposes an already-stamped live worker; the post-gate reuse path needs no stamps)")
			return 2
		}
		return runBuild(probe, workflowLauncher, opts, stdin, stdout, stderr)
	case "show-stage-def":
		if wantsHelp(args[1:]) {
			printShowStageDefUsage(stdout)
			return 0
		}
		workflowDir, stage, code := requireStageFlags(args[1:], stderr)
		if code != 0 {
			return code
		}
		return runShowStageDef(workflowDir, stage, stdout, stderr)
	case "trunk":
		workflowDir, code := requireSubcommandFlag(args[1:], "trunk", "--workflow-dir", stderr)
		if code != 0 {
			return code
		}
		return runDispatchTrunk(workflowDir, stdout, stderr)
	case "context-budget":
		name, code := requireSubcommandFlag(args[1:], "context-budget", "--name", stderr)
		if code != 0 {
			return code
		}
		return claudeteam.ContextBudget(os.Getenv("HOME"), name, stdout, stderr)
	case "list-standing":
		workflowDir, code := requireSubcommandFlag(args[1:], "list-standing", "--workflow-dir", stderr)
		if code != 0 {
			return code
		}
		return runListStanding(workflowDir, stdout, stderr)
	case "show-standing":
		workflowDir, code := requireSubcommandFlag(args[1:], "show-standing", "--workflow-dir", stderr)
		if code != 0 {
			return code
		}
		return runShowStanding(workflowDir, stdout, stderr)
	case "spawn-standing-all":
		flags := parseFlags(args[1:], map[string]bool{"--workflow-dir": true})
		wd, okWD := flags["--workflow-dir"]
		if !okWD {
			fmt.Fprintln(stderr, "error: dispatch spawn-standing-all requires --workflow-dir")
			return 2
		}
		return runSpawnStandingAll(wd, stdout, stderr)
	case "reconcile":
		return runReconcile(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "error: unknown dispatch subcommand: %s\n", args[0])
		printUsage(stderr)
		return 2
	}
}

type buildOptions struct {
	WorkflowDir         string
	Host                string
	EntityPath          string
	Stage               string
	ChecklistFile       string
	ScopeNotesFile      string
	FeedbackContextFile string
	BareMode            bool
	FeedbackReflow      bool
	Advance             bool
	Stamp               bool
	PrintSchema         bool
	ValidateOnly        string
	requestFlagProvided bool
}

func (o buildOptions) hasRequestFlags() bool {
	return o.requestFlagProvided
}

func parseBuildOptions(args []string, stderr io.Writer) (buildOptions, int) {
	var opts buildOptions
	valueFlags := map[string]*string{
		"--workflow-dir":          &opts.WorkflowDir,
		"--host":                  &opts.Host,
		"--entity-path":           &opts.EntityPath,
		"--stage":                 &opts.Stage,
		"--checklist-file":        &opts.ChecklistFile,
		"--scope-notes-file":      &opts.ScopeNotesFile,
		"--feedback-context-file": &opts.FeedbackContextFile,
		"--validate-only":         &opts.ValidateOnly,
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--team-name" || strings.HasPrefix(a, "--team-name=") {
			fmt.Fprintln(stderr, "error: unknown flag --team-name: legacy TeamCreate-registry dispatch mode is retired; omit --team-name, auto-team is the only claude shape")
			return opts, 2
		}
		if eq := indexByte(a, '='); eq > 0 {
			name := a[:eq]
			if dst, ok := valueFlags[name]; ok {
				*dst = a[eq+1:]
				if isBuildRequestFlag(name) {
					opts.requestFlagProvided = true
				}
			}
			continue
		}
		if dst, ok := valueFlags[a]; ok {
			if i+1 >= len(args) {
				fmt.Fprintf(stderr, "error: dispatch build requires value for %s\n", a)
				return opts, 2
			}
			*dst = args[i+1]
			if isBuildRequestFlag(a) {
				opts.requestFlagProvided = true
			}
			i++
			continue
		}
		switch a {
		case "--bare-mode":
			opts.BareMode = true
			opts.requestFlagProvided = true
		case "--feedback-reflow":
			opts.FeedbackReflow = true
			opts.requestFlagProvided = true
		case "--advance":
			opts.Advance = true
			opts.requestFlagProvided = true
		case "--stamp":
			opts.Stamp = true
			opts.requestFlagProvided = true
		case "--print-schema":
			opts.PrintSchema = true
		}
	}
	return opts, 0
}

func isBuildRequestFlag(name string) bool {
	switch name {
	case "--entity-path", "--stage", "--checklist-file", "--scope-notes-file",
		"--feedback-context-file":
		return true
	default:
		return false
	}
}

// wantsHelp reports whether args request subcommand help. Stop at -- so future
// passthrough forms can carry a literal --help without changing this contract.
func wantsHelp(args []string) bool {
	for _, a := range args {
		if a == "--" {
			return false
		}
		if a == "-h" || a == "--help" {
			return true
		}
	}
	return false
}

// requireFlag returns the value of a single required `name value` flag. A
// missing flag, missing value, or trailing junk is a usage error (exit 2).
func requireFlag(args []string, name string, stderr io.Writer) (string, int) {
	val, ok := parseFlags(args, map[string]bool{name: true})[name]
	if !ok {
		fmt.Fprintf(stderr, "error: dispatch build requires %s\n", name)
		return "", 2
	}
	return val, 0
}

// requireSubcommandFlag returns the value of a single required flag for a
// claude-coupled subcommand, with a usage error (exit 2) naming the subcommand
// and flag when missing. The diagnostic is the native CLI's own ergonomic — the
// command-logic loud-failures (not this argument-parse error) are what the parity
// harness byte-compares.
func requireSubcommandFlag(args []string, subcommand, name string, stderr io.Writer) (string, int) {
	val, ok := parseFlags(args, map[string]bool{name: true})[name]
	if !ok {
		fmt.Fprintf(stderr, "error: dispatch %s requires %s\n", subcommand, name)
		return "", 2
	}
	return val, 0
}

// requireStageFlags returns the --workflow-dir and --stage values show-stage-def
// requires, with a usage error (exit 2) when either is missing.
func requireStageFlags(args []string, stderr io.Writer) (string, string, int) {
	flags := parseFlags(args, map[string]bool{"--workflow-dir": true, "--stage": true})
	wd, okWD := flags["--workflow-dir"]
	stage, okStage := flags["--stage"]
	if !okWD || !okStage {
		fmt.Fprintln(stderr, "error: dispatch show-stage-def requires --workflow-dir and --stage")
		return "", "", 2
	}
	return wd, stage, 0
}

// parseFlags reads `--flag value` pairs from args for the flags in want. Unknown
// flags and bare arguments are ignored — the required-flag checks above surface
// the actionable error, matching the oracle's argparse required-flag behavior.
func parseFlags(args []string, want map[string]bool) map[string]string {
	out := map[string]string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if want[a] && i+1 < len(args) {
			out[a] = args[i+1]
			i++
			continue
		}
		if eq := indexByte(a, '='); eq > 0 && want[a[:eq]] {
			out[a[:eq]] = a[eq+1:]
		}
	}
	return out
}

// indexByte returns the index of the first b in s, or -1.
func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

func printUsage(w io.Writer) {
	fmt.Fprint(w, `spacedock dispatch assembles ensign dispatch artifacts.

Usage:
  spacedock dispatch build --workflow-dir DIR        (stdin JSON -> stdout JSON)
  spacedock dispatch build --workflow-dir DIR --entity-path FILE --stage STAGE --checklist-file FILE [--host claude|codex|pi]
  spacedock dispatch build --print-schema
  spacedock dispatch build --validate-only FILE
  spacedock dispatch show-stage-def --workflow-dir DIR --stage STAGE
  spacedock dispatch trunk --workflow-dir DIR
  spacedock dispatch reconcile --workflow-dir DIR [--team-name NAME] [--repo-root DIR] [--include lingering,superseded,un-advanced-pr,stale-branch,local-main-drift]
`)
}

func printBuildUsage(w io.Writer) {
	fmt.Fprint(w, `Usage:
  spacedock dispatch build --workflow-dir DIR < request.json                                                   (stdin mode)
  spacedock dispatch build --workflow-dir DIR --entity-path FILE --stage STAGE --checklist-file FILE [flags]   (flag/file mode)
  spacedock dispatch build --print-schema
  spacedock dispatch build --validate-only FILE

Build an ensign dispatch artifact and write the JSON envelope to stdout. The
request comes from a JSON object on stdin OR from flags/files. The two are
selected by the rule below and never merged.

Input mode selection:
  If ANY request flag is present the request is read from flags/files and stdin
  is IGNORED (flag/file mode); otherwise the request is read as a JSON object on
  stdin (stdin JSON mode). Request flags:
    --entity-path  --stage  --checklist-file  --scope-notes-file
    --feedback-context-file  --bare-mode  --feedback-reflow  --advance  --stamp
  Flag/file mode requires --entity-path, --stage, and --checklist-file; any
  request flag with one of the three missing fails:
    error: flag/file input requires --entity-path, --stage, and --checklist-file
  Because --advance is a request flag, piping JSON on stdin together with
  --advance is NOT accepted -- it selects flag/file mode and ignores the piped
  JSON. Pass a reuse-advance request in flag/file form (see the --advance
  example below).

Flags:
  --workflow-dir DIR            Workflow definition directory containing README.md (both modes).
  --host HOST                   Override the runtime host (claude|codex|pi). Defaults to the detected runtime (both modes).
  --entity-path FILE            Entity file for this dispatch (flag/file mode).
  --stage STAGE                 Stage name to dispatch (flag/file mode).
  --checklist-file FILE         File of checklist lines, one per line (flag/file mode).
  --scope-notes-file FILE       Optional scope-notes file (flag/file mode).
  --feedback-context-file FILE  Optional feedback-context file; required with --feedback-reflow (flag/file mode).
  --bare-mode                   Emit the bare sequential shape (no name, no run_in_background); unsupported on host=codex.
  --feedback-reflow             Route a rejection back to its feedback-to target stage; requires --feedback-context-file.
  --advance                     Emit a reuse-advance pointer message for a live worker instead of a spawn envelope. Incompatible with --bare-mode.
  --stamp                       Fold the ordinary post-gate dispatch steps (started/worktree frontmatter stamps, state commit+sync, worktree creation) into this build, before assembling the envelope. Refuses (no mutation) unless the entity's status already equals --stage. Incompatible with --advance.
  --print-schema                Print the stdin request JSON schema and exit.
  --validate-only FILE          Validate a request JSON file without writing a dispatch; exit 0 on success.

Stdin JSON request fields (stdin JSON mode):
  schema_version  Dispatch schema version. The current supported value is 2.
  entity_path     Path to the entity file for this dispatch.
  workflow_dir    Workflow directory for the dispatch request.
  stage           Stage name to dispatch.
  checklist       Array of checklist strings for the dispatched worker.
  (optional: scope_notes, feedback_context, bare_mode, is_feedback_reflow, advance, host)

Examples:
  stdin JSON mode:
  {"schema_version":2,"entity_path":"thing.md","workflow_dir":".","stage":"implementation","checklist":["DONE: run tests"]}

  flag/file mode:
  spacedock dispatch build --workflow-dir . --entity-path thing.md --stage implementation --checklist-file impl.checklist

  reuse-advance (flag/file mode):
  spacedock dispatch build --workflow-dir . --entity-path thing.md --stage validation --checklist-file validation.checklist --advance
`)
}

func printShowStageDefUsage(w io.Writer) {
	fmt.Fprint(w, `Usage:
  spacedock dispatch show-stage-def --workflow-dir DIR --stage STAGE

Print a stage's workflow README section followed by its declared context sections.

Flags:
  --workflow-dir DIR   Workflow definition directory containing README.md.
  --stage STAGE        Stage name whose definition should be printed.
`)
}
