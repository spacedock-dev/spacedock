// ABOUTME: cobra command tree, grouped help, and exit-code behavior for spacedock.
// ABOUTME: status/dispatch forward argv verbatim; claude/codex use the Option-2 grammar.
package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
	"github.com/spacedock-dev/spacedock/internal/dispatch"
	"github.com/spacedock-dev/spacedock/internal/gates"
	"github.com/spacedock-dev/spacedock/internal/runtimehost"
	"github.com/spacedock-dev/spacedock/internal/safehouse"
	"github.com/spacedock-dev/spacedock/internal/status"
)

// Version is the single source of truth for the binary version. It is stamped by
// the release pipeline via -ldflags "-X
// github.com/spacedock-dev/spacedock/internal/cli.Version=$(git describe --tags --always)".
// It is a var (not a const) because the linker can only write package-level vars;
// a const is silently ignored by -X. The default is the `dev` sentinel so an
// unstamped `go build`/`go install` binary reads honestly as a dev build rather
// than impersonating a stale release; the git-describe tag overwrites it on a
// stamped release build.
var Version = "dev"

// envGetenv adapts the injected `KEY=VALUE` environment slice — the seam the CLI
// already threads everywhere instead of reading the process environment — into
// the getenv function the detection helpers take, so a test can pin the whole
// environment without touching the running machine's.
func envGetenv(env []string) func(string) string {
	vars := envMap(env)
	return func(key string) string { return vars[key] }
}

// tagline is the one-line product description rendered as the first help line.
const tagline = "spacedock — agentic workflow launcher"

// Run is the process entry point. status is routed to the native Go runner,
// which composes the definition root (README) and the entity root (the README's
// state: dir) itself; all other commands are handled directly. The vendored
// Python runner stays selectable through the injectable run() core.
func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	// The native binary is the Claude runtime's companion: the Claude first officer
	// invokes `spacedock status --boot` and `spacedock dispatch build` directly, so
	// the workflow surface is wired with the Claude team-state probe. claudeteam owns
	// the ~/.claude read; status/dispatch take it as an opaque value. A non-Claude
	// runtime entry point wires nil (host-neutral present:false / no bare-mode advisory).
	return run(context.Background(), args, os.Environ(), cwd(), os.Stdin, stdout, stderr,
		&status.NativeRunner{TeamStateProbe: claudeteam.Probe}, claudeteam.Probe)
}

// run is the injectable core. It depends only on the status.Runner interface,
// never on the vendored script or any exec detail, so the fake-runner tests can
// drive the status path with pinned env/cwd. cobra is wired INSIDE run so the
// package's public surface (Run) and the exit-code contract are unchanged: the
// command tree captures env/dir/stdin/stdout/stderr/runner in its RunE closures.
func run(ctx context.Context, args []string, env []string, dir string, stdin io.Reader, stdout io.Writer, stderr io.Writer, runner status.Runner, dispatchProbe claudeteam.TeamStateProbe) int {
	root := newRootCommand(ctx, args, env, dir, stdin, stdout, stderr, runner, dispatchProbe)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		return exitCodeFor(err)
	}
	return 0
}

// exitCodeError carries an explicit process exit code out of a RunE so the
// command tree can preserve the hand-rolled router's exit-code contract (status
// exit-1 surfacing, the front-door fail-fast exit 1) through cobra's single error
// return. cobra's own command-resolution errors (unknown command, unknown flag)
// carry no code and map to exit 2.
type exitCodeError struct{ code int }

func (e exitCodeError) Error() string { return fmt.Sprintf("exit %d", e.code) }

// exitCodeFor maps an Execute error to a process exit code. An explicit
// exitCodeError carries its own code (RunE already wrote any diagnostic); every
// other error is a cobra command/flag-resolution failure, which exits 2 to match
// the hand-rolled router's unknown-command contract (TestUnknownCommand).
func exitCodeFor(err error) int {
	var ec exitCodeError
	if errors.As(err, &ec) {
		return ec.code
	}
	return 2
}

// newRootCommand assembles the cobra tree. The root owns the grouped jargon-free
// help (AC-1) and the explicit `--version` handler with the `(contract N)` token
// (AC-5). SilenceErrors/SilenceUsage hand all output and exit-code control to this
// package: cobra never prints its own error or usage, so the unknown-command path
// emits the pinned message and exits 2 (the root RunE below), and the help is
// rendered solely by printHelp.
func newRootCommand(ctx context.Context, rawArgs []string, env []string, dir string, stdin io.Reader, stdout, stderr io.Writer, runner status.Runner, dispatchProbe claudeteam.TeamStateProbe) *cobra.Command {
	versionFlag := false

	root := &cobra.Command{
		Use:           "spacedock",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ArbitraryArgs,
		// An unknown subcommand carrying a trailing flag (e.g. the removed
		// `spacedock init --host claude`) must reach RunE so args[0] drives the
		// unknown-command diagnostic. Without this, the root flagset errors on the
		// unrecognized flag during parse and, under SilenceErrors/SilenceUsage,
		// exits 2 with no output — a silent exit. Whitelisting unknown flags lets
		// parsing fall through to RunE so the command token is reported.
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		CompletionOptions:  cobra.CompletionOptions{DisableDefaultCmd: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			if versionFlag {
				// The Channel line must report what the binary will ACTUALLY do, so
				// an explicit SPACEDOCK_DEV_BRANCH override (including empty) applies
				// before printVersion reads the package devBranch var — the same
				// helper every install path already calls.
				applyDevBranchOverride(env)
				printVersion(stdout, envGetenv(env), exec.LookPath)
				return nil
			}
			// No subcommand and no recognized flag: an unknown command token
			// (e.g. `spacedock bogus`) exits 2 with the pinned message.
			if len(args) > 0 {
				return unknownCommand(args[0], stderr)
			}
			// A leading unknown flag (e.g. `spacedock --bogus`, or the space-form
			// `spacedock --foo install` where `--foo` consumed `install` as its
			// value) reaches RunE with empty args because the UnknownFlags whitelist
			// strips it during parse. Left to the bare-help path it would exit 0 — a
			// silent usage error. Detect it from the raw argv and route to the
			// usage/exit-2 path; a bare `spacedock` (no leading flag) still helps + 0.
			if leadingUnknownFlag(rawArgs) {
				return unknownFlag(rawArgs[0], stderr)
			}
			printHelp(stdout)
			return nil
		},
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.Flags().BoolVar(&versionFlag, "version", false, "Print the spacedock version and contract level")
	root.SetHelpFunc(func(*cobra.Command, []string) { printHelp(stdout) })

	root.AddGroup(
		&cobra.Group{ID: "launch", Title: "Launch"},
		&cobra.Group{ID: "setup", Title: "Setup"},
		&cobra.Group{ID: "workflow", Title: "Workflow"},
	)

	root.AddCommand(
		newClaudeCommand(ctx, env, dir, stdout, stderr),
		newCodexCommand(ctx, env, dir, stdout, stderr),
		newPiCommand(ctx, env, dir, stdout, stderr),
		newInstallCommand(ctx, env, stdout, stderr),
		newDoctorCommand(ctx, env, stdout, stderr),
		newStatusCommand(ctx, env, dir, stdin, stdout, stderr, runner),
		newNewCommand(ctx, env, dir, stdin, stdout, stderr, runner),
		newStateCommand(ctx, env, dir, stdout, stderr),
		newMergeCommand(ctx, env, dir, stdout, stderr),
		newCompletionCommand(stdout, stderr),
		newDispatchCommand(dispatchProbe, env, stdin, stdout, stderr),
		newGateCommand(dir, stdout, stderr),
	)
	return root
}

// newGateCommand exposes only semantic decision sources; lifecycle mechanics,
// CAS values, and durable ids belong to the recorder.
func newGateCommand(dir string, stdout, stderr io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:                "gate prepare|withdraw|record|consume <entity>",
		Short:              "Prepare, withdraw, record, or consume durable gate resolutions",
		GroupID:            "workflow",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				fmt.Fprintln(stdout, "Usage: spacedock gate prepare <entity> --question TEXT --artifact REVIEW.md --summary TEXT [--reference FILE ...] [--workflow-dir DIR]\n       spacedock gate withdraw <entity> --reason TEXT [--workflow-dir DIR]\n       spacedock gate record <entity> --decision approve|revise|hold --actor ID [--reason TEXT] [--conn-quote TEXT --conn-source TEXT] [--consume] [--workflow-dir DIR]\n       spacedock gate record <entity> --round STAGE/CYCLE --briefing PATH/briefing.json --log PATH/briefing.review.jsonl [--workflow-dir DIR]\n       spacedock gate consume <entity> [--workflow-dir DIR]\n\nOn an approval whose target stage is terminal, consume spends nothing: it leaves the\napplication pending and reports route=approved-awaiting-merge. The terminal merge\nceremony (`spacedock merge guard <slug> --verdict passed|rejected`) is the sole terminal\nconsumer; `merge guard --rework` sends a failed delivery back through the declared\nfeedback-to (pending -> superseded, delivery state cleared).\n\n`gate record --consume` is the captain-approve fast path: close, sync, consume, sync\nin one call. `--consume` requires --decision approve and is rejected as a usage error\nwith --decision revise|hold. A delegated `--actor agent:first-officer` decision requires\n`--conn-quote` (the grant verbatim) and `--conn-source` (where it was given); those flags\nare refused with `--actor person:captain` or with `--round`. In a split-root workflow, a\nsuccessful close or consume ends with a machine-parseable `sync=.../phase=...` line;\nbranch on that final line plus the exit code, never on which prose lines printed.")
				return nil
			}
			if len(args) < 2 || (args[0] != "prepare" && args[0] != "withdraw" && args[0] != "record" && args[0] != "consume") {
				fmt.Fprintln(stderr, "spacedock gate: unknown subcommand (want: prepare|withdraw|record|consume)")
				return exitCodeError{2}
			}
			workflowDir := ""
			input := gates.RecordInput{}
			prepareInput := gates.PrepareInput{}
			questionCount, artifactCount, summaryCount, reasonCount := 0, 0, 0, 0
			consumeFlag := false
			for i := 2; i < len(args); i++ {
				if args[i] == "--consume" {
					consumeFlag = true
					continue
				}
				if i+1 >= len(args) {
					fmt.Fprintf(stderr, "Error: %s requires an argument\n", args[i])
					return exitCodeError{2}
				}
				switch args[i] {
				case "--workflow-dir":
					workflowDir = args[i+1]
				case "--briefing":
					input.BriefingPath = args[i+1]
				case "--actor":
					input.Actor = args[i+1]
				case "--decision":
					input.Decision = args[i+1]
				case "--reason":
					input.Reason = args[i+1]
					reasonCount++
				case "--conn-quote":
					input.ConnQuote = args[i+1]
				case "--conn-source":
					input.ConnSource = args[i+1]
				case "--round":
					input.Round = args[i+1]
				case "--log":
					input.LogPath = args[i+1]
				case "--question":
					prepareInput.Question = args[i+1]
					questionCount++
				case "--artifact":
					prepareInput.Artifact = args[i+1]
					artifactCount++
				case "--summary":
					prepareInput.Summary = args[i+1]
					summaryCount++
				case "--reference":
					prepareInput.References = append(prepareInput.References, args[i+1])
				default:
					fmt.Fprintf(stderr, "Error: unknown gate flag: %s\n", args[i])
					return exitCodeError{2}
				}
				i++
			}
			if consumeFlag && args[0] != "record" {
				fmt.Fprintf(stderr, "Error: --consume is only valid with gate record\n")
				return exitCodeError{2}
			}
			// --round is an advisory correction-round publication, not a close
			// or a consume attempt — mechanism 1 deliberately excludes it from
			// the implicit sync (Alternatives rejected 6: zero AC-1 benefit, a
			// new failure surface for a verb outside the measured ceremony
			// window), so it still needs an explicit `state commit` afterward,
			// same as before mechanism 1. --consume has nothing to sequence
			// here, so reject it as a usage error rather than silently
			// ignoring it.
			if consumeFlag && input.Round != "" {
				fmt.Fprintln(stderr, "Error: --consume is not valid with --round (no close/consume attempt to sequence; --round still needs an explicit state commit)")
				return exitCodeError{2}
			}
			if args[0] == "prepare" {
				if summaryCount != 1 {
					fmt.Fprintln(stderr, "gate prepare accepts --summary exactly once")
					return exitCodeError{1}
				}
				if !utf8.ValidString(prepareInput.Summary) {
					fmt.Fprintln(stderr, "--summary must be valid UTF-8")
					return exitCodeError{1}
				}
				if questionCount != 1 || artifactCount != 1 {
					fmt.Fprintln(stderr, "Error: gate prepare requires --question and --artifact exactly once")
					return exitCodeError{2}
				}
				if input != (gates.RecordInput{}) {
					fmt.Fprintln(stderr, "Error: gate prepare accepts only --question, --artifact, --summary, --reference, and --workflow-dir")
					return exitCodeError{2}
				}
			} else if questionCount != 0 || artifactCount != 0 || summaryCount != 0 || len(prepareInput.References) != 0 {
				fmt.Fprintf(stderr, "Error: gate %s does not accept prepare flags\n", args[0])
				return exitCodeError{2}
			}
			if args[0] == "withdraw" {
				allowed := input
				allowed.Reason = ""
				if reasonCount != 1 || allowed != (gates.RecordInput{}) {
					fmt.Fprintln(stderr, "Error: gate withdraw accepts exactly one --reason and optional --workflow-dir")
					return exitCodeError{2}
				}
			}
			definitionDir := workflowDir
			if definitionDir == "" {
				var code int
				definitionDir, code = status.ResolveWorkflowDir(dir, stderr)
				if code != 0 {
					return exitCodeError{code}
				}
			} else if !filepath.IsAbs(definitionDir) {
				definitionDir = filepath.Join(dir, definitionDir)
			}
			path, err := status.ResolveActivePath(definitionDir, dir, args[1], stderr)
			if err != nil {
				fmt.Fprintln(stderr, "Error:", err)
				return exitCodeError{1}
			}
			if args[0] == "prepare" {
				prepareInput.WorkflowDir = definitionDir
				prepareInput.LaunchDir = dir
				result, err := gates.Prepare(path, prepareInput)
				if err != nil {
					fmt.Fprintln(stderr, "Error:", err)
					return exitCodeError{1}
				}
				fmt.Fprintf(stdout, "room=%s\nbriefing=%s\ndigest=%s\nstate=%s\n", result.Room, result.Briefing, result.Digest, result.State)
				return nil
			}
			if args[0] == "withdraw" {
				s, err := gates.Withdraw(path, gates.WithdrawInput{Reason: input.Reason, WorkflowDir: definitionDir})
				if err != nil {
					fmt.Fprintln(stderr, "Error:", err)
					return exitCodeError{1}
				}
				fmt.Fprintf(stdout, "withdrawn gate=%s attempt=%s state=%s briefing=%s\n", s.Gate, s.Attempt, s.State, s.Briefing)
				return nil
			}
			if args[0] == "consume" {
				if input != (gates.RecordInput{}) {
					fmt.Fprintf(stderr, "Error: gate %s accepts only --workflow-dir\n", args[0])
					return exitCodeError{2}
				}
				// runGateConsumeAndSync renders the exact gate=.../consumed=.../route=
				// line this branch always has, then — mechanism 1 — a machine-parseable
				// sync=.../phase=consume line when (and only when) the call wrote (a
				// real advance or a stale-pending supersede); a terminal route or an
				// ineligible/blocked refusal writes nothing and stays byte-clean.
				if code := runGateConsumeAndSync(path, definitionDir, stdout, stderr); code != 0 {
					return exitCodeError{code}
				}
				return nil
			}
			// --conn-quote/--conn-source cite the grant behind a chat close; --round
			// and --briefing are advisory correction-round publication, not a chat
			// close, so a citation there is unrepresentable. Refuse before either
			// branch mutates anything.
			if (input.ConnQuote != "" || input.ConnSource != "") && (input.Round != "" || input.BriefingPath != "") {
				fmt.Fprintln(stderr, "Error: --conn-quote and --conn-source are not valid with --round or --briefing")
				return exitCodeError{2}
			}
			if input.Round != "" {
				input.WorkflowDir = definitionDir
				if err := gates.RecordSemantic(path, input); err != nil {
					fmt.Fprintln(stderr, "Error:", err)
					return exitCodeError{1}
				}
				summary, err := gates.ValidateRoundFile(path, input.Round)
				return printRound(summary, err, stdout, stderr)
			}
			if input.BriefingPath != "" {
				fmt.Fprintln(stderr, "Error: gate record --briefing requires --round")
				return exitCodeError{2}
			}
			if input.Decision == "" {
				fmt.Fprintln(stderr, "Error: gate record requires --decision")
				return exitCodeError{2}
			}
			if input.Actor == "" {
				fmt.Fprintln(stderr, "Error: --decision requires --actor ID")
				return exitCodeError{2}
			}
			// The two record shapes are disjoint by grammar: a delegated FO chat
			// decision always cites the grant it acted under; a captain decision
			// cites no grant. Refuse the incoherent shapes before any mutation —
			// the citation attributes, it never authorizes, so a warned-but-written
			// record would reproduce the exact ambiguity this grammar removes.
			if input.Actor == "agent:first-officer" && (strings.TrimSpace(input.ConnQuote) == "" || strings.TrimSpace(input.ConnSource) == "") {
				fmt.Fprintln(stderr, "Error: delegated First Officer decision requires --conn-quote and --conn-source")
				return exitCodeError{2}
			}
			if input.Actor == "person:captain" && (input.ConnQuote != "" || input.ConnSource != "") {
				fmt.Fprintln(stderr, "Error: --conn-quote and --conn-source are refused on a person:captain decision")
				return exitCodeError{2}
			}
			// --consume never softens a non-approve chat decision. This is a usage
			// error (exit 2) before any write.
			if consumeFlag && input.Decision != "approve" {
				fmt.Fprintln(stderr, "Error: --consume with --decision revise or hold is a usage error; the flag never softens a non-approve decision")
				return exitCodeError{2}
			}
			input.WorkflowDir = definitionDir
			s, err := gates.RecordSemanticSummary(path, input)
			if err != nil {
				fmt.Fprintln(stderr, "Error:", err)
				return exitCodeError{1}
			}
			fmt.Fprintf(stdout, "recorded gate=%s attempt=%s state=%s briefing=%s resolution=%s decision=%s\n", s.Gate, s.Attempt, s.State, s.Briefing, s.Resolution, s.Decision)

			// Mechanism 1: the close always wrote a Resolution, so it always syncs
			// (split-root only; inline prints nothing, exit 0 unchanged).
			entityStage := status.ParseFrontmatter(path)["status"]
			recordMsg := fmt.Sprintf("gate: record %s %s %s", status.EntitySlug(path), entityStage, s.Decision)
			if code := runGateSync(stdout, stderr, definitionDir, path, gateSyncPhaseRecord, recordMsg); code != 0 {
				return exitCodeError{code}
			}
			if !consumeFlag {
				return nil
			}
			// Sequence the existing consume handler in the same invocation.
			if code := runGateConsumeAndSync(path, definitionDir, stdout, stderr); code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
}

func printRound(s gates.RoundSummary, err error, stdout, stderr io.Writer) error {
	if err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return exitCodeError{1}
	}
	fmt.Fprintf(stdout, "round=%s stage=%s cycle=%d briefing=%s entries=%d\n", s.ID, s.Stage, s.Cycle, s.Briefing, len(s.Entries))
	for _, entry := range s.Entries {
		fmt.Fprintf(stdout, "entry=%s type=%s advisory=%t decision=%s\n", entry.ID, entry.Type, entry.Advisory, entry.Decision)
	}
	return nil
}

// newClaudeCommand wires `spacedock claude`. Flag parsing is disabled at the cobra
// layer so the post-subcommand argv reaches runClaude verbatim — runClaude owns the
// Option-2 grammar via parseFrontDoorArgs (ArgsLenAtDash). The flags are declared
// only so `--help` renders them (AC-4); `-h`/`--help` is intercepted here because
// DisableFlagParsing routes it to RunE rather than cobra's own help.
func newClaudeCommand(ctx context.Context, env []string, dir string, stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "claude [task] [-- claude-flags]",
		Short:              "Start Claude Code as your Spacedock first officer",
		GroupID:            "launch",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				return cmd.Help()
			}
			applyDevBranchOverride(env)
			applyMarketplaceSourceOverride(env)
			if code := runClaude(ctx, args, dir, execHost{}, exec.LookPath, stdout, stderr); code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
	setFrontDoorHelp(cmd, "claude", stdout)
	return cmd
}

// newCodexCommand mirrors newClaudeCommand for `spacedock codex`.
func newCodexCommand(ctx context.Context, env []string, dir string, stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "codex [task] [-- codex-flags]",
		Short:              "Start Codex as your Spacedock first officer",
		GroupID:            "launch",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				return cmd.Help()
			}
			applyDevBranchOverride(env)
			applyMarketplaceSourceOverride(env)
			if code := runCodex(ctx, args, dir, execHost{}, exec.LookPath, stdout, stderr); code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
	setFrontDoorHelp(cmd, "codex", stdout)
	return cmd
}

// newPiCommand wires `spacedock pi` to Pi's native skill/extension resource
// loading instead of Claude/Codex plugin or team-tool semantics.
func newPiCommand(ctx context.Context, env []string, dir string, stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "pi [task] [-- pi-flags]",
		Short:              "Start Pi as your Spacedock first officer",
		GroupID:            "launch",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				return cmd.Help()
			}
			if code := runPi(ctx, args, dir, env, execPiRuntimeOps{}, stdout, stderr); code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
	setPiHelp(cmd, stdout)
	return cmd
}

// newInstallCommand wires `spacedock install` (the renamed `init`). Behavior is
// unchanged from init: install the per-host plugin then run doctor (claude), or
// emit the documented codex add prose. DisableFlagParsing keeps the post-subcommand
// argv verbatim for the existing hand-parsed runInit (so `--host`/`--check` parse
// exactly as before); `-h`/`--help` is intercepted here.
func newInstallCommand(ctx context.Context, env []string, stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "install [--host claude|codex|pi] [--check]",
		Short:              "Install the Spacedock plugin for a host, then check it",
		GroupID:            "setup",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				return cmd.Help()
			}
			applyDevBranchOverride(env)
			applyMarketplaceSourceOverride(env)
			if code := runInitWithPi(ctx, args, execHost{}, execPiRuntimeOps{}, env, stdout, stderr); code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
	cmd.Flags().String("host", "claude", "Host to install the plugin for (claude, codex, or pi)")
	cmd.Flags().Bool("check", false, "Run the compatibility report without installing")
	setSetupHelp(cmd, stdout, `
Examples:
  spacedock install
  spacedock install --host codex
  spacedock install --host pi
  spacedock install --check
`)
	return cmd
}

// newDoctorCommand wires `spacedock doctor` with its existing hand-parsed
// `--host`/`--plugin-manifest` handling preserved verbatim.
func newDoctorCommand(ctx context.Context, env []string, stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "doctor [--host claude|codex|pi]",
		Short:              "Check the installed plugin and this binary are compatible",
		GroupID:            "setup",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				return cmd.Help()
			}
			applyDevBranchOverride(env)
			if code := runDoctorWithPi(ctx, args, execHost{}, execPiRuntimeOps{}, env, stdout, stderr); code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
	cmd.Flags().String("host", "claude", "Host to check (claude, codex, or pi)")
	cmd.Flags().String("plugin-manifest", "", "Read this manifest directly instead of resolving the installed plugin")
	setSetupHelp(cmd, stdout, `
Examples:
  spacedock doctor
  spacedock doctor --host codex
  spacedock doctor --host pi --plugin-dir ./checkout
`)
	return cmd
}

// newStatusCommand reparents `spacedock status` under cobra with flag parsing
// disabled, so its post-subcommand argv forwards VERBATIM to runStatus exactly as
// the hand-rolled router did — cobra never consumes, reorders, or validates a
// status flag (AC-5).
func newStatusCommand(ctx context.Context, env []string, dir string, stdin io.Reader, stdout, stderr io.Writer, runner status.Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "status [args]",
		Short:              "Show or update workflow state",
		GroupID:            "workflow",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				return cmd.Help()
			}
			if code := runStatus(ctx, args, env, dir, stdin, stdout, stderr, runner); code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
	setStatusHelp(cmd, stdout)
	return cmd
}

// newNewCommand wires `spacedock new [--folder] SLUG` as a pure alias for
// `status --new`: the post-subcommand argv (the optional --folder plus the slug)
// is prefixed with --new and forwarded verbatim to runStatus, so the body is read
// from stdin and the existing runNew atomic-create path is reused unchanged. With
// the discovery walk-up, `new` run inside a workflow needs no --workflow-dir.
// DisableFlagParsing keeps --folder reaching the runner intact (AC-3).
func newNewCommand(ctx context.Context, env []string, dir string, stdin io.Reader, stdout, stderr io.Writer, runner status.Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "new [--folder] SLUG",
		Short:              "Create an entity from a stdin body (auto-discovers the workflow)",
		GroupID:            "workflow",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				return cmd.Help()
			}
			aliased := append([]string{"--new"}, args...)
			if code := runStatus(ctx, aliased, env, dir, stdin, stdout, stderr, runner); code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
	setNewHelp(cmd, stdout)
	return cmd
}

// newStateCommand wires `spacedock state init|new|ready|sweep|commit` for split-root
// state-checkout management. Flag parsing is disabled so the post-subcommand argv
// (the optional --workflow-dir, a commit <slug>) reaches the handler verbatim.
// `init` resumes a cloned workflow's state checkout; `new` births one; `ready`
// integrates peers' state on boot; `sweep` lists merged-but-not-terminalized
// entities; `commit <slug>` path-scoped-commits and syncs one entity. An unknown or
// missing subcommand is a usage error (exit 2).
func newStateCommand(ctx context.Context, env []string, dir string, stdout, stderr io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:                "state init|new|ready|sweep|commit [--workflow-dir DIR]",
		Short:              "Manage a split-root workflow's state checkout (init/new births, ready/sweep/commit sync)",
		GroupID:            "workflow",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				return cmd.Help()
			}
			if len(args) == 0 {
				fmt.Fprintln(stderr, "spacedock state: unknown subcommand (want: init|new|ready|sweep|commit)")
				return exitCodeError{2}
			}
			var code int
			switch args[0] {
			case "init":
				code = runStateInit(ctx, args[1:], env, dir, stdout, stderr)
			case "new":
				code = runStateNew(ctx, args[1:], env, dir, stdout, stderr)
			case "ready":
				code = runStateReady(ctx, args[1:], env, dir, stdout, stderr)
			case "sweep":
				code = runStateSweep(ctx, args[1:], env, dir, stdout, stderr)
			case "commit":
				code = runStateCommit(ctx, args[1:], env, dir, stdout, stderr)
			default:
				fmt.Fprintln(stderr, "spacedock state: unknown subcommand (want: init|new|ready|sweep|commit)")
				return exitCodeError{2}
			}
			if code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
}

// newMergeCommand wires `spacedock merge guard <slug>` for the terminal
// merge-finalize ceremony. Flag parsing is disabled so the post-subcommand argv
// (the slug plus --workflow-dir / --verdict / --json / --quiet) reaches the
// handler verbatim. `guard` drives the atomic mod-block set->invoke->clear->
// terminalize sequence; an unknown or missing subcommand is a usage error (exit 2).
func newMergeCommand(ctx context.Context, env []string, dir string, stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "merge guard <slug> (--verdict passed|rejected | --rework)",
		Short:              "Run the terminal merge ceremony for an entity (finalize or rework)",
		GroupID:            "workflow",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				return cmd.Help()
			}
			if len(args) == 0 {
				fmt.Fprintln(stderr, "spacedock merge: unknown subcommand (want: guard)")
				return exitCodeError{2}
			}
			switch args[0] {
			case "guard":
				if code := status.MergeGuard(args[1:], dir, stdout, stderr); code != 0 {
					return exitCodeError{code}
				}
				return nil
			default:
				fmt.Fprintln(stderr, "spacedock merge: unknown subcommand (want: guard)")
				return exitCodeError{2}
			}
		},
	}
	setMergeHelp(cmd, stdout)
	return cmd
}

// newCompletionCommand wires `spacedock completion bash|zsh`, emitting a static
// completion script to stdout (exit 0). An unknown or missing shell prints the
// named usage error and returns 2 — the CLI-layer usage-error code, matching the
// unknown-command path. The static script (no dynamic slug completion: YAGNI)
// replaces cobra's default completion command, which is disabled at the root via
// CompletionOptions.DisableDefaultCmd (AC-3).
func newCompletionCommand(stdout, stderr io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:                "completion bash|zsh",
		Short:              "Print a bash or zsh completion script",
		GroupID:            "workflow",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				return cmd.Help()
			}
			if code := runCompletion(args, stdout, stderr); code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
}

// newDispatchCommand reparents `spacedock dispatch` under cobra with flag parsing
// disabled, forwarding its post-subcommand argv verbatim to dispatch.Run (AC-5).
func newDispatchCommand(probe claudeteam.TeamStateProbe, env []string, stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:                "dispatch build | show-stage-def",
		Short:              "Build worker dispatch artifacts",
		GroupID:            "workflow",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowLauncher, _ := resolvedDispatchLauncher(env)
			if code := dispatch.RunWithLauncher(probe, workflowLauncher, args, stdin, stdout, stderr); code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
}

// wantsHelp reports whether the operator asked for command help. Commands with
// DisableFlagParsing receive `-h`/`--help` as ordinary args, so each RunE checks
// for it before doing work. Only a leading help token counts: a `--help` after
// `--` is host passthrough, not a request for spacedock's help.
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

// unknownCommand writes the pinned unknown-command diagnostic plus the grouped
// help to stderr and returns the exit-2 carrier.
func unknownCommand(name string, stderr io.Writer) error {
	fmt.Fprintf(stderr, "unknown command: %s\n", name)
	printHelp(stderr)
	return exitCodeError{2}
}

// unknownFlag writes the pinned unknown-flag diagnostic plus the grouped help to
// stderr and returns the exit-2 carrier. It mirrors unknownCommand for a leading
// flag that resolves to no subcommand, so a stray `spacedock --bogus` is a loud
// usage error rather than a silent exit 0.
func unknownFlag(name string, stderr io.Writer) error {
	fmt.Fprintf(stderr, "unknown flag: %s\n", name)
	printHelp(stderr)
	return exitCodeError{2}
}

// leadingUnknownFlag reports whether the raw argv begins with a `-`-prefixed token
// that is not a recognized root flag. `--version` and `-h`/`--help` are handled
// before this is consulted (version returns early; help is intercepted by the
// help func), so any leading flag reaching the bare-help path is unknown — the
// UnknownFlags whitelist stripped it during parse, leaving empty positional args.
// A bare invocation (no args, or a leading non-flag token) is not a leading flag.
func leadingUnknownFlag(rawArgs []string) bool {
	return len(rawArgs) > 0 && strings.HasPrefix(rawArgs[0], "-")
}

// runStatus forwards the post-"status" argv verbatim to the runner and returns
// its exit code unmodified. The CLI adds nothing to and removes nothing from the
// runner's contract: it does not parse, reformat, interpret, or strip flags. If
// the runner itself cannot run (interpreter missing), surface a diagnostic and
// fail loudly with exit 1 — matching the script's own error exit code.
func runStatus(ctx context.Context, args []string, env []string, dir string, stdin io.Reader, stdout io.Writer, stderr io.Writer, runner status.Runner) int {
	code, err := runner.Run(ctx, status.Request{
		Args:   args,
		Dir:    dir,
		Env:    env,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return code
}

// applyDevBranchOverride lets SPACEDOCK_DEV_BRANCH override the pinned devBranch
// default (and the linker stamp). An UNSET env var leaves the default in place —
// the released binary keeps targeting `@next` — while an explicit value (including
// empty, to force the no-ref release path) wins.
func applyDevBranchOverride(env []string) {
	for _, kv := range env {
		if strings.HasPrefix(kv, "SPACEDOCK_DEV_BRANCH=") {
			devBranch = strings.TrimPrefix(kv, "SPACEDOCK_DEV_BRANCH=")
			return
		}
	}
}

// applyMarketplaceSourceOverride lets SPACEDOCK_MARKETPLACE_SOURCE override the
// default marketplace add source (`spacedock-dev/marketplace`) on the install
// paths. An UNSET env var leaves the default in place — the released binary keeps
// installing from the production marketplace — while an explicit value wins,
// pointing the install at a local/alternate marketplace for dogfooding.
func applyMarketplaceSourceOverride(env []string) {
	for _, kv := range env {
		if strings.HasPrefix(kv, "SPACEDOCK_MARKETPLACE_SOURCE=") {
			marketplaceSource = strings.TrimPrefix(kv, "SPACEDOCK_MARKETPLACE_SOURCE=")
			return
		}
	}
}

// cwd returns the working directory, falling back to "" so a getwd failure does
// not abort the command — the runner derives a scan root from --workflow-dir.
func cwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// frozenContractToken is D4's second cross-era sentinel: an integer-era FO
// reading `--version` parses `contract <N>` and checks it against its own
// contract range (`>=2,<3`). This binary carries no contract-integer mechanism at
// all — no constant, no compare math — but it still emits this LITERAL string so
// that old prose sees 3, which is at/above its upper bound, and aborts with
// "update the plugin": the correct remedy, since the skills genuinely are too old
// for a binary this new. Note the gate is prose executed by a model, so an ABSENT
// token is not a reliable abort — a model may reason "no token, nothing to check,
// proceed" where a deterministic parser would error. Emitting the literal keeps
// the abort correct in the old prose's own terms.
//
// The value must stay 3 — bumping it would false-green old skills against every
// future binary. It is pinned by internal/cli/version_session_test.go, which
// asserts the literal below the Sandbox line — internal/contractlint carries no
// reference to it at all.
//
// RETIREMENT CONDITION: removable once no plugin or binary predating #468 can
// still be running. Both reader populations are ours — old spacedock skills and
// old spacedock binaries, things we ship and tag — so this is a query against the
// Homebrew formula and the marketplace, not a guess.
const frozenContractToken = "contract 3"

// printVersion reports the binary version and, when this process is running
// inside an agent session, that session. The organising rule is: inside a
// session, report the session; outside one, report the version. Anything about
// what is INSTALLED belongs to `doctor`, which is already the version gate's own
// named remedy surface — so the per-host plugin block this used to print is gone,
// and with it the three host CLIs it shelled to produce it.
//
// Line 1 is `spacedock <version>` and nothing else. It is load-bearing: the
// FO/ensign version gate parses that token and aborts on any other shape (cobra's
// auto version-flag is deliberately NOT used). Line 2 is always
// `OS: <goos>/<goarch>` — in BOTH output shapes — so user issue reports carry
// the platform and later gate-logic versions can read the OS from `--version`
// once a compatible binary exists. Line 3 is always `Channel: <stable|edge>
// (<plugin-id>)` — in BOTH output shapes — reporting the EFFECTIVE devBranch: the
// `--version` path applies applyDevBranchOverride(env) before this call, so a
// SPACEDOCK_DEV_BRANCH override renders what the binary will really do. The
// frozen contract token moved BELOW line 1 — the integer-era prose says "run
// --version and parse contract <N>" and never pins it to line 1 — and prints
// inside a session only, since every integer-era reader is itself a session.
//
// Ambiguous markers are REPORTED, never guessed at, and never fail: refusing here
// would break the version gate and therefore every boot, including the nested-
// runtime marker leak that occurs in practice.
func printVersion(w io.Writer, getenv func(string) string, lookPath func(string) (string, error)) {
	fmt.Fprintf(w, "spacedock %s\n", displayVersion())
	fmt.Fprintf(w, "OS: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	// The channel word mirrors channelMarketplace's own two-way branch (main →
	// stable, anything else → edge); the parenthetical is the plugin id the
	// frontdoor re-ensures on every launch, so the incident's symptom (the
	// spacedock@spacedock reinstall loop) and its cause read off one line.
	channelWord := "edge"
	if devBranch == "main" {
		channelWord = "stable"
	}
	fmt.Fprintf(w, "Channel: %s (%s)\n", channelWord, channelPluginID(devBranch))

	host, markers, ambiguous := runtimehost.Detect(getenv)
	if !ambiguous && host == "" {
		// Outside every runtime — a human at a terminal. Three lines: the version
		// line, the OS line, and the Channel line, nothing else.
		return
	}

	if ambiguous {
		fmt.Fprintf(w, "Runtime: ambiguous (%s)\n", strings.Join(markers, ", "))
	} else {
		fmt.Fprintf(w, "Runtime: %s (%s)\n", host, strings.Join(markers, ", "))
	}

	insideName, inside := safehouse.Inside(getenv)
	available, _ := safehouse.Available(lookPath)
	fmt.Fprintf(w, "Sandbox: %s\n", safehouse.SessionState(insideName, inside, available))
	fmt.Fprintln(w, frozenContractToken)
}

// runCompletion emits a static shell-completion script for bash or zsh to
// stdout (exit 0). An unknown or missing shell prints the named usage error to
// stderr and returns 2 — the CLI-layer usage-error code, consistent with the
// unknown-command path, since completion is handled in-package and never reaches
// the native runner.
func runCompletion(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "Error: completion requires a shell: bash or zsh")
		return 2
	}
	switch args[0] {
	case "bash":
		fmt.Fprint(stdout, bashCompletion)
		return 0
	case "zsh":
		fmt.Fprint(stdout, zshCompletion)
		return 0
	default:
		fmt.Fprintln(stderr, "Error: completion requires a shell: bash or zsh")
		return 2
	}
}

// bashCompletion completes the top-level verbs and the common status flags. It
// is intentionally static (no dynamic slug completion): YAGNI.
const bashCompletion = `# spacedock bash completion
_spacedock() {
  local cur prev verbs status_flags
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
  verbs="claude codex pi install doctor status new state merge completion dispatch --version --help"
  status_flags="--workflow-dir --next --next-id --boot --identify --validate --archived --json --quiet --new --folder --set --where --archive --resolve --short-id --discover --root --page --limit"
  if [ "$COMP_CWORD" -eq 1 ]; then
    COMPREPLY=( $(compgen -W "$verbs" -- "$cur") )
    return 0
  fi
  case "${COMP_WORDS[1]}" in
    status) COMPREPLY=( $(compgen -W "$status_flags" -- "$cur") ) ;;
    state) COMPREPLY=( $(compgen -W "init new ready sweep commit --workflow-dir" -- "$cur") ) ;;
    merge) COMPREPLY=( $(compgen -W "guard" -- "$cur") ) ;;
    completion) COMPREPLY=( $(compgen -W "bash zsh" -- "$cur") ) ;;
  esac
}
complete -F _spacedock spacedock
`

// zshCompletion completes the top-level verbs and the common status flags.
const zshCompletion = `#compdef spacedock
# spacedock zsh completion
_spacedock() {
  local -a verbs status_flags
  verbs=(claude codex pi install doctor status new state merge completion dispatch --version --help)
  status_flags=(--workflow-dir --next --next-id --boot --identify --validate --archived --json --quiet --new --folder --set --where --archive --resolve --short-id --discover --root --page --limit)
  if (( CURRENT == 2 )); then
    compadd -- $verbs
    return
  fi
  case "${words[2]}" in
    status) compadd -- $status_flags ;;
    state) compadd -- init new ready sweep commit --workflow-dir ;;
    merge) compadd -- guard ;;
    completion) compadd -- bash zsh ;;
  esac
}
_spacedock "$@"
`
