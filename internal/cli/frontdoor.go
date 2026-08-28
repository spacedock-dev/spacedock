// ABOUTME: Host front doors (spacedock claude/codex) + init/doctor — the three
// ABOUTME: version-gate points wired through an injectable host-ops seam.
package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/pflag"

	"github.com/spacedock-dev/spacedock/internal/contract"
	"github.com/spacedock-dev/spacedock/internal/safehouse"
	"github.com/spacedock-dev/spacedock/internal/status"
)

// bootstrapPrompt is the fixed launch-and-go message appended as the last inner
// argv token so a fresh `spacedock claude` session starts the first officer
// rather than opening an idle agent. It is omitted when `--resume` is forwarded
// (a resume already carries its own session intent).
const bootstrapPrompt = "You totally got this. Take your time. I love you. And tell all subagents and team members you love them too."

// hostOps is the injectable seam the front-door, init, and doctor paths depend
// on. Production backs it with real `claude`/`codex` plugin commands and exec;
// tests back it with a fake that records interactions. Mirrors the status.Runner
// decoupling: the command logic never touches exec or the host CLI directly.
type hostOps interface {
	// ResolveManifest returns the installed plugin manifest path for host, or ""
	// when no plugin is installed (a distinct, non-error state). A non-nil error
	// means the host CLI itself failed.
	ResolveManifest(host string) (string, error)
	// PluginInventory supplies the launch gate and doctor's sibling-channel view.
	PluginInventory(host string) ([]pluginInventoryEntry, error)
	// Launch spawns argv with env as a resident child and waits (production) or
	// records it (test), returning the host's propagated exit code. The error is
	// reserved for a launch failure (host binary not found, fork failure), not a
	// non-zero host exit.
	Launch(argv []string, env []string) (int, error)
	// Install issues the host plugin commands to install/update the plugin from
	// source, returning combined output. devBranch selects the marketplace channel
	// entry the install targets (see channelEntry).
	Install(host, source, devBranch string) (string, error)
	// InstallCodexLocalPluginDir issues the codex `--plugin-dir` dev-install
	// sequence (codexPluginDirInstallArgvSequence) against source — the
	// dedicated `spacedock-local` marketplace, distinct from either real
	// channel's marketplace name — returning combined output.
	InstallCodexLocalPluginDir(source string) (string, error)
}

// pluginInventoryEntry normalizes the Claude and Codex plugin-list schemas.
type pluginInventoryEntry struct {
	ID        string
	Version   string
	Installed bool
	Enabled   bool
}

// devBranch is the binary's channel stamp: it selects which marketplace entry the
// install targets — `main` installs the stable `spacedock` entry, any other value
// (default `next`) installs the `spacedock-edge` entry tracking next HEAD (see
// channelEntry). It is a var (not a const) so the linker can stamp it per channel,
// mirroring Version, and so `SPACEDOCK_DEV_BRANCH` can override it; tests
// save/restore it.
var devBranch = "next"

var executablePath = os.Executable

const spacedockBinEnv = "SPACEDOCK_BIN"

// agentTeamsEnv is the env flag associated with claude's worker↔FO inter-agent communication
// (SendMessage/TeamCreate). The authoritative enabler is ~/.claude/settings.json
// (env.CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 + teammateMode:auto), which re-applies
// the flag to every child regardless of shell env; the launcher export below does
// NOT independently enable the feature.
const agentTeamsEnv = "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS"

func launchEnv(parent []string) []string {
	env := withoutEnv(parent, spacedockBinEnv)
	if bin, ok := resolvedLauncherBin(); ok {
		env = append(env, spacedockBinEnv+"="+bin)
	}
	return env
}

// withAgentTeams sets agentTeamsEnv=1 in the launched child env unless the parent
// already set it — an explicit operator value (even =0) is preserved. This is a
// best-effort export for users without the authoritative settings.json enabler
// (see agentTeamsEnv) and a no-op when settings already enable it; it does NOT
// independently enable inter-agent communication. Claude-only: codex/pi do not call this.
func withAgentTeams(env []string) []string {
	if hasEnv(env, agentTeamsEnv) {
		return env
	}
	return append(env, agentTeamsEnv+"=1")
}

func hasEnv(env []string, key string) bool {
	prefix := key + "="
	for _, entry := range env {
		if strings.HasPrefix(entry, prefix) {
			return true
		}
	}
	return false
}

// launcherBinEnvPassFlags returns the `--env-pass SPACEDOCK_BIN` safehouse flags
// that tell safehouse to forward SPACEDOCK_BIN from its (the launching process's)
// environment into the otherwise-sanitized sandbox, so the launcher binary the
// helper calls resolve survives the boundary. launchEnv already sets SPACEDOCK_BIN
// on the safehouse process; this flag carries it through. It is gated on the same
// resolvedLauncherBin() source as launchEnv and mirrors its omit-on-failure: when
// no binary resolves, no flag — never a stale pass-through. Returned as safehouse
// wrap flags (before `--`) so the inner program safehouse sees stays the real host
// (claude/codex), keeping its program-keyed profile auto-detection intact.
func launcherBinEnvPassFlags() []string {
	if _, ok := resolvedLauncherBin(); ok {
		return []string{"--env-pass", spacedockBinEnv}
	}
	return nil
}

func withoutEnv(env []string, key string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env))
	for _, entry := range env {
		if strings.HasPrefix(entry, prefix) {
			continue
		}
		out = append(out, entry)
	}
	return out
}

func resolvedLauncherBin() (string, bool) {
	p, err := executablePath()
	if err != nil || p == "" {
		return "", false
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", false
	}
	if resolved, err := filepath.EvalSymlinks(abs); err == nil && executableFile(resolved) {
		return resolved, true
	}
	if executableFile(abs) {
		return abs, true
	}
	return "", false
}

// resolvedDispatchLauncher selects the explicit launcher identity propagated to
// the agent session before falling back to this process's executable. Selection
// happens once while dispatch build runs; generated commands receive the
// resulting absolute path and never depend on the environment again.
func resolvedDispatchLauncher(env []string) (string, bool) {
	if candidate := envGetenv(env)(spacedockBinEnv); candidate != "" {
		if !filepath.IsAbs(candidate) {
			if abs, err := filepath.Abs(candidate); err == nil {
				candidate = abs
			}
		}
		if resolved, err := filepath.EvalSymlinks(candidate); err == nil && executableFile(resolved) {
			return resolved, true
		}
		if filepath.IsAbs(candidate) && executableFile(candidate) {
			return candidate, true
		}
	}
	return resolvedLauncherBin()
}

func executableFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return true
	}
	return info.Mode().Perm()&0o111 != 0
}

// noPluginRemedy is the host-correct manual remedy printed when no plugin is
// installed and the operator opted out of auto-install with --no-install. It
// names the host's own install and bootstrap commands — a codex run never says
// `claude`. The caller owns this message (not gateHost) because the response to
// NoPluginFound is non-uniform: auto-install by default, refuse-and-instruct
// under --no-install.
func noPluginRemedy(host string) string {
	return fmt.Sprintf(
		"Spacedock: no installed %s plugin found. "+
			"Run `spacedock install --host %s` (or `spacedock %s --skip-compat-check` to bootstrap).",
		host, host, host)
}

// launchBanner writes a short pre-launch orientation banner to w before the host
// is handed control: the spacedock version, the workflow detected from dir, the
// sandbox posture, and a one-line orientation pointer. Callers suppress the
// banner on a resume (the operator is continuing a session, not starting one).
//
// The Sandbox: line answers the LAUNCH question — will the launch this banner
// precedes be wrapped? — which is why the banner alone keeps safehouse.LaunchState
// while `--version` and `status --boot` take SessionState: those report a session
// that is already running, and a .safehouse profile says nothing about one.
// `selected` is whether this launch would be wrapped (a .safehouse profile or a
// --safehouse* flag), `available` is whether the safehouse binary resolves via
// lookPath, and `getenv` resolves whether this process is ALREADY inside a
// sandbox — the arm that stops the banner reporting `unavailable` from a launch
// made inside one. Both seams are injected so tests pin them.
func launchBanner(host, dir string, selected bool, getenv func(string) string, lookPath func(string) (string, error), w io.Writer) {
	label, value := detectedWorkflow(dir)
	available, _ := safehouse.Available(lookPath)
	insideName, inside := safehouse.Inside(getenv)
	fmt.Fprintf(w, "spacedock %s · launching %s as your first officer\n", displayVersion(), host)
	fmt.Fprintf(w, "%s: %s\n", label, value)
	fmt.Fprintf(w, "Sandbox: %s\n", safehouse.LaunchState(insideName, inside, selected, available))
	fmt.Fprintf(w, "%s is your first officer — ask it for the queue and next steps.\n", host)
}

// detectedWorkflow names the workflow the launch belongs to, found from ANY
// launch dir, returning the banner's line-2 label ("Workflow"/"Workflows") and its
// value so the label agrees with the content. When dir is inside a git repo, the
// whole repo is scanned downward (status.DiscoverWorkflows, which prunes the
// linked/agent-worktree + VCS noise), so launching from the repo root resolves the
// repo's real workflow rather than missing it. A single workflow is named by its
// path relative to the repo root (e.g. `docs/dev`); two or more list the detected
// paths space-separated in discovery order (`Workflows: docs/dev docs/user-testing`),
// the first officer disambiguating in-session. With no enclosing git repo, a
// bounded walk-up names a workflow at or above dir. The value is never the
// cwd-relative `.`/`..`; "none detected" is shown when no workflow is found.
func detectedWorkflow(dir string) (label, value string) {
	const noneDetected = "none detected"
	repoRoot := status.FindGitRoot(dir)
	gitEntry := filepath.Join(repoRoot, ".git")
	if dirExists(gitEntry) || fileExists(gitEntry) { // .git is a dir (repo) or a file (worktree gitlink)
		workflows := status.DiscoverWorkflows(repoRoot)
		switch len(workflows) {
		case 0:
			return "Workflow", noneDetected
		case 1:
			return "Workflow", workflowLabel(repoRoot, workflows[0])
		default:
			labels := make([]string, len(workflows))
			for i, wf := range workflows {
				labels[i] = workflowLabel(repoRoot, wf)
			}
			return "Workflows", strings.Join(labels, " ")
		}
	}
	// No enclosing git repo: a bounded walk-up names a workflow at or above dir.
	if workflowDir, ok := status.DiscoverWorkflowDir(dir); ok {
		return "Workflow", workflowLabel(repoRoot, workflowDir)
	}
	return "Workflow", noneDetected
}

// workflowLabel renders a workflow dir as a recognizable path relative to base
// (e.g. `docs/dev` relative to the repo root), falling back to the workflow dir's
// own base name. It never returns `.` or a `..`-escaping path. Both paths are
// resolved through symlinks first so a base that is a symlinked parent (e.g.
// macOS's /var → /private/var temp dirs) still yields the clean relative path
// rather than a spurious `..`-escape.
func workflowLabel(base, workflowDir string) string {
	if rel, err := filepath.Rel(realpath(base), realpath(workflowDir)); err == nil && rel != "." && !strings.HasPrefix(rel, "..") {
		return rel
	}
	return filepath.Base(workflowDir)
}

// realpath resolves path through symlinks, returning the original on error.
func realpath(path string) string {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		return resolved
	}
	return path
}

// gateHost resolves the installed manifest for host and compares its version
// against the binary's, returning the full Result so the caller can distinguish
// a missing plugin (NoPluginFound — recoverable by installing) from a version
// mismatch (too-old-binary fails fast; too-old-plugin is ALSO recoverable by
// installing, D6). Both no-plugin states (an empty manifestPath AND a
// resolved-but-missing manifest) collapse to NoPluginFound; a host-CLI resolve
// error maps to MalformedVersion so it stays a hard fail (a broken host CLI is
// not a missing plugin). The gate inspects the VERDICT, not a doctor exit code:
// RunDoctor maps no-plugin-found to exit 0 (a non-fatal report), so a non-empty
// installPath to a missing manifest would otherwise slip through as
// "compatible". gateHost prints the actionable remedy for too-old-binary and
// malformed-version (the two verdicts with exactly one caller response: fail
// fast) but stays SILENT for NoPluginFound and TooOldPlugin, whose message the
// caller owns — both auto-install by default and only refuse under --no-install,
// so the right wording (and whether to print at all before a silent heal, D6)
// depends on the caller's choice. On Compatible it stays silent on the bare OK
// line but surfaces the opt-in upgrade hint when the plugin is compatible yet
// behind a strictly-newer binary; the hint never blocks — the caller still
// proceeds to launch.
func gateHost(ops hostOps, host string, stderr io.Writer) contract.Result {
	manifestPath, err := ops.ResolveManifest(host)
	if err != nil {
		msg := fmt.Sprintf(
			"Spacedock: could not resolve the installed %s plugin (%v). "+
				"Run `spacedock doctor` or `spacedock install --host %s`.", host, err, host)
		fmt.Fprintln(stderr, msg)
		return contract.Result{Verdict: contract.MalformedVersion, Message: msg}
	}
	if manifestPath == "" {
		return contract.Result{Verdict: contract.NoPluginFound}
	}
	res := contract.ManifestVerdict(manifestPath, host, displayVersion(), runningEdgeCask())
	switch res.Verdict {
	case contract.NoPluginFound, contract.TooOldPlugin:
		return res
	case contract.Compatible:
		// Compatible but behind: surface the opt-in upgrade hint (the front door
		// stays silent on the bare OK line). The hint never blocks — the caller
		// still proceeds to launch on Compatible.
		if res.Hint != "" {
			fmt.Fprintln(stderr, res.Hint)
		}
		return res
	default:
		fmt.Fprintln(stderr, res.Message)
		return res
	}
}

// resolveHealableGate runs the version gate for host and, for its healable
// states (NoPluginFound, TooOldPlugin, or an enabled sibling), auto-installs and
// re-gates ONCE (the default) or refuses with the host-correct remedy
// (--no-install) — D6's shared "a single command yields a working session"
// contract for both front doors. A too-old-plugin heal announces "Refreshing"
// rather than "Installing" (gateHost prints no scary mismatch remedy before the
// silent heal). A second miss after the re-gate refuses rather than launching
// blind — one retry, no loop. It returns true when the caller should proceed to
// launch; false means the caller must return 1 (the refusal/failure message is
// already on stderr).
func resolveHealableGate(ops hostOps, host string, noInstall bool, stderr io.Writer) bool {
	res := gateHost(ops, host, stderr)
	announce := ""
	switch res.Verdict {
	case contract.Compatible:
		inventory, err := ops.PluginInventory(host)
		if err != nil {
			fmt.Fprintf(stderr, "Spacedock: could not verify the %s plugin enablement state: %v.\n", host, err)
			fmt.Fprintf(stderr, "Run `spacedock install --host %s` before launching.\n", host)
			return false
		}
		if _, conflict := enabledSiblingPlugin(inventory); !conflict {
			return true
		}
		announce = "Refreshing the " + host + " plugin to remove its enabled sibling…"
	case contract.NoPluginFound, contract.TooOldPlugin:
		announce = "Installing the " + host + " plugin…"
		if res.Verdict == contract.TooOldPlugin {
			announce = "Refreshing the " + host + " plugin…"
		}
	default:
		// too-old-binary / malformed-version: gateHost already printed the
		// remedy. Fail fast — auto-installing would not fix an incompatibility.
		return false
	}
	if noInstall {
		if res.Verdict == contract.Compatible {
			printSiblingRemedy(host, stderr)
		} else {
			printHealableRemedy(host, res, stderr)
		}
		return false
	}
	fmt.Fprintln(stderr, announce)
	if _, err := ops.Install(host, channelMarketplaceSource(devBranch), devBranch); err != nil {
		fmt.Fprintf(stderr, "spacedock %s: auto-install failed: %v\n", host, err)
		return false
	}
	regate := gateHost(ops, host, stderr)
	if regate.Verdict != contract.Compatible {
		printHealableRemedy(host, regate, stderr)
		return false
	}
	inventory, err := ops.PluginInventory(host)
	if _, conflict := enabledSiblingPlugin(inventory); err != nil || conflict {
		fmt.Fprintf(stderr, "Spacedock: the %s plugin repair did not leave one enabled channel.\n", host)
		printSiblingRemedy(host, stderr)
		return false
	}
	return true
}

func printSiblingRemedy(host string, stderr io.Writer) {
	fmt.Fprintf(stderr, "Run `spacedock install --host %s` to keep only the %s channel.\n", host, selectedChannelWord())
}

// printHealableRemedy prints the caller-owned remedy for a NoPluginFound or
// TooOldPlugin verdict: the host-correct install/bootstrap message for
// NoPluginFound (gateHost carries no message for it), or the gate's own
// mismatch message for TooOldPlugin.
func printHealableRemedy(host string, res contract.Result, stderr io.Writer) {
	if res.Verdict == contract.NoPluginFound {
		fmt.Fprintln(stderr, noPluginRemedy(host))
		return
	}
	fmt.Fprintln(stderr, res.Message)
}

// runClaude is the `spacedock claude` front door: version-gate, then launch the
// first officer. The gate fails fast on a too-old-binary mismatch, but the two
// healable verdicts (NoPluginFound, TooOldPlugin) auto-install the plugin and
// proceed to launch so the single command the user typed yields a working session
// (D6) — `--no-install` opts out,
// preserving the refuse-and-instruct behavior. The launch is interposed through
// `safehouse --trust-workdir-config [extra] -- claude --dangerously-skip-permissions …`
// when ANY of {a `.safehouse` profile in dir, the bare `--safehouse` flag, a
// `--safehouse-*` knob} is given; otherwise it is plain `claude --agent
// spacedock:first-officer …` (no skip-permissions in an unsandboxed launch). The
// `--safehouse-*` knobs translate into the safehouse `extra` slot. The bootstrap
// prompt is appended last (base, or base + " " + task when a task is fenced after
// `--`) unless a resume is forwarded. The gate is bypassed by an explicit
// `--skip-compat-check`, by a declared pre-fence `--plugin-dir`, or by a valid
// host-specific plugin checkout adjacent to the resolved launcher. Post-fence
// plugin directories remain native Claude additions and do not affect provider
// selection or gating. `lookPath` resolves the safehouse binary (default
// exec.LookPath; injected so tests pin not-found).
func runClaude(ctx context.Context, args []string, dir string, ops hostOps, lookPath func(string) (string, error), stdout, stderr io.Writer) int {
	fd, err := parseFrontDoorArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock claude: %v\n", err)
		return 1
	}
	extra, err := safehouse.TranslateFlags(fd.safehouseFlags)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock claude: %v\n", err)
		return 1
	}
	// A declared pre-fence plugin directory is the compatibility-preserving local
	// Spacedock override. Without one, a valid checkout beside the resolved
	// launcher supplies the same session-local provider automatically. Native
	// post-fence Claude plugin directories are additions only: they neither select
	// Spacedock nor suppress the installed compatibility gate.
	localSpacedock := fd.pluginDirPairs > 0
	if !localSpacedock {
		if adjacent, ok := adjacentSpacedockPluginRoot("claude"); ok {
			fd.passthrough = append([]string{"--plugin-dir", adjacent}, fd.passthrough...)
			localSpacedock = true
		}
	}
	if !fd.skipCheck && !localSpacedock {
		if !resolveHealableGate(ops, "claude", fd.noInstall, stderr) {
			return 1
		}
	}
	warnStrayPromptAfterDash(fd, "spacedock claude", stderr)

	wrap := safehouse.Present(dir) || fd.forceSafehouse || len(fd.safehouseFlags) > 0
	resume := containsResume(fd.passthrough)
	if !resume {
		launchBanner("claude", dir, wrap, os.Getenv, lookPath, stderr)
	}
	inner := []string{"claude"}
	if wrap {
		inner = append(inner, "--dangerously-skip-permissions")
	}
	inner = append(inner, "--agent", "spacedock:first-officer")
	// An unsandboxed launch has no safehouse isolation, so per-action permission
	// prompting is friction without a matching safety gain: start the first officer
	// in auto permission-mode. Suppressed when the operator already chose a mode and
	// on a resume (which rides its own session intent, like the bootstrap prompt).
	// The sandboxed arm's --dangerously-skip-permissions above already covers its
	// posture, so this is the !wrap counterpart, not a replacement.
	if !wrap && !resume && !passthroughHasFlag(fd.passthrough, "--permission-mode") {
		inner = append(inner, "--permission-mode", "auto")
	}
	inner = append(inner, fd.passthrough...)
	if !resume {
		inner = append(inner, launchPrompt(bootstrapPrompt, fd))
	}

	argv := inner
	if wrap {
		if ok, hint := safehouse.Available(lookPath); !ok {
			fmt.Fprintln(stderr, hint)
			return 1
		}
		// Forward SPACEDOCK_BIN through safehouse's env sanitization with
		// `--env-pass`: launchEnv sets it on the safehouse process, this flag carries
		// it into the sandbox. It rides the safehouse flags (before `--`), so the
		// inner program stays `claude` and safehouse's program-keyed profile
		// auto-detection still fires. Omitted when the bin cannot be resolved.
		argv = safehouse.Wrap(inner, append(launcherBinEnvPassFlags(), extra...))
	}

	code, err := ops.Launch(argv, withAgentTeams(launchEnv(os.Environ())))
	if err != nil {
		fmt.Fprintf(stderr, "spacedock claude: launch failed: %v\n", err)
		return 1
	}
	return code
}

// warnStrayPromptAfterDash emits an advisory stderr warning when a bare positional
// appears after `--` with no task before it — almost always an operator prompt that
// silently degrades to host passthrough so the spacedock launch prompt is never
// prepended to it. The warning names the stray positional and the corrected form
// (put the prompt BEFORE `--`). It does NOT alter the assembled host argv; the
// launch is byte-identical with or without this call. `name` is the front-door
// verb so the message names a runnable fix.
func warnStrayPromptAfterDash(fd frontDoorArgs, name string, stderr io.Writer) {
	pos, ok := strayPromptAfterDash(fd)
	if !ok {
		return
	}
	fmt.Fprintf(stderr,
		"%s: warning: positional %q after `--` is forwarded to the host as-is; "+
			"the spacedock launch prompt was NOT prepended to it. "+
			"To make it the launch prompt, put it BEFORE `--`: `%s %q -- …`\n",
		name, pos, name, pos)
}

// launchPrompt returns the inner-argv launch prompt: `base + " " + task` when the
// operator fenced a task after `--`, otherwise the bare base prompt. Claude and
// Codex suppress it on their respective resume forms.
func launchPrompt(base string, fd frontDoorArgs) string {
	if fd.hasTask {
		return base + " " + fd.task
	}
	return base
}

// hasPluginDir reports whether the Claude passthrough carries a `--plugin-dir`
// flag (either `--plugin-dir P` or `--plugin-dir=P`). Its presence relaxes
// Claude's version gate because the local checkout supersedes the installed plugin.
func hasPluginDir(passthrough []string) bool {
	for _, a := range passthrough {
		if a == "--plugin-dir" || strings.HasPrefix(a, "--plugin-dir=") {
			return true
		}
	}
	return false
}

// takeCodexPluginDir consumes only the declared pre-fence --plugin-dir pairs
// parseFrontDoorArgs injected at the front of passthrough. Codex itself rejects
// that Spacedock-owned flag, so those owned pairs become a local-marketplace
// install. Tokens written after the fence are forwarded host argv, even if they
// spell `--plugin-dir`, and are never consumed here.
func takeCodexPluginDir(passthrough []string, ownedPairs int) (dir string, rest []string, found bool) {
	if ownedPairs == 0 {
		return "", passthrough, false
	}
	rest = passthrough
	for range ownedPairs {
		if len(rest) < 2 || rest[0] != "--plugin-dir" {
			return "", passthrough, false
		}
		dir = rest[1] // the last declared pre-fence checkout wins
		rest = rest[2:]
	}
	return dir, rest, true
}

// passthroughHasFlag reports whether the operator already supplied any of the
// named host flags in the passthrough, in either `--flag value` or `--flag=value`
// form. The unsandboxed launchers consult it before injecting their default
// permission/approval flag so an operator-supplied one is never duplicated
// (operator wins). Mirrors hasPluginDir, generalized over a flag set.
func passthroughHasFlag(passthrough []string, names ...string) bool {
	for _, a := range passthrough {
		for _, name := range names {
			if a == name || strings.HasPrefix(a, name+"=") {
				return true
			}
		}
	}
	return false
}

// containsResume reports whether the operator forwarded any of claude's
// session-resume forms (which carry their own session intent, so the bootstrap
// prompt is suppressed): `--resume`, `--resume=<id>`, `-r`, `--continue`, `-c`.
func containsResume(args []string) bool {
	for _, a := range args {
		switch a {
		case "--resume", "-r", "--continue", "-c":
			return true
		}
		if strings.HasPrefix(a, "--resume=") {
			return true
		}
	}
	return false
}

// codexBootstrapPrompt is the fixed launch-and-go message appended as the last
// inner argv token so a fresh `spacedock codex` session starts the first
// officer. Codex has no `--agent` analog (spike-confirmed: no agent/skill-select
// flag on the top-level, `exec`, or `plugin` surfaces), so the only FO-selection
// injection point is the positional prompt — this prompt names the
// `spacedock:first-officer` skill explicitly.
const codexBootstrapPrompt = "You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Assume $spacedock:first-officer for the entire session."

var codexCollaborationLayer = []string{
	"-c", "agents.enabled=true",
	"-c", "features.multi_agent=true",
	"-c", `features.multi_agent_v2={max_concurrent_threads_per_session=16,tool_namespace="agents",hide_spawn_agent_metadata=false}`,
}

func hasCodexCollaborationOverride(args []string) bool {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-c", "--config":
			if i+1 < len(args) && reservedCodexConfig(args[i+1]) {
				return true
			}
			i++
		case "--enable", "--disable":
			if i+1 < len(args) && reservedCodexFeature(args[i+1]) {
				return true
			}
			i++
		default:
			for _, prefix := range []string{"-c", "--config="} {
				if strings.HasPrefix(arg, prefix) && reservedCodexConfig(strings.TrimPrefix(arg, prefix)) {
					return true
				}
			}
			for _, prefix := range []string{"--enable=", "--disable="} {
				if strings.HasPrefix(arg, prefix) && reservedCodexFeature(strings.TrimPrefix(arg, prefix)) {
					return true
				}
			}
		}
	}
	return false
}

func reservedCodexConfig(assignment string) bool {
	assignment = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(assignment), "="))
	key, _, ok := strings.Cut(assignment, "=")
	if !ok {
		return false
	}
	parts := strings.Split(key, ".")
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) >= 2 && part[0] == part[len(part)-1] && (part[0] == '\'' || part[0] == '"') {
			if part[0] == '\'' {
				part = part[1 : len(part)-1]
			} else if unquoted, err := strconv.Unquote(part); err == nil {
				part = unquoted
			}
		}
		parts[i] = part
	}
	key = strings.Join(parts, ".")
	return key == "agents.enabled" || key == "features.multi_agent" ||
		key == "features.multi_agent_v2" || strings.HasPrefix(key, "features.multi_agent_v2.")
}

func reservedCodexFeature(feature string) bool {
	feature = strings.TrimSpace(feature)
	return feature == "multi_agent" || feature == "multi_agent_v2"
}

// runCodex is the `spacedock codex` front door: version-gate, then launch the
// first officer. The gate fails fast on a too-old-binary mismatch, but the two
// healable verdicts (NoPluginFound, TooOldPlugin) auto-install the plugin and
// proceed to launch so the single command the user typed yields a working session
// (D6) — `--no-install` opts out,
// preserving the refuse-and-instruct behavior. The launch is interposed through
// `safehouse --trust-workdir-config [extra] -- codex --dangerously-bypass-approvals-and-sandbox …`
// when ANY of {a `.safehouse` profile in dir, the bare `--safehouse` flag, a
// `--safehouse-*` knob} is given — safehouse is the sandbox, so codex's own
// sandbox is bypassed. Otherwise the launch is plain `codex …` keeping codex's own
// sandbox (the bypass flag is omitted: it is safe only when safehouse provides the
// sandbox). After Spacedock consumes its own pre-fence flags, an exact `resume`
// token in Codex's forwarded argv suppresses the banner, default approval mode,
// and bootstrap prompt; every other forwarded argv retains the normal fresh-launch
// posture. A declared pre-fence `--plugin-dir`, or a valid host-specific plugin
// checkout adjacent to the resolved launcher when no explicit override exists,
// installs the local checkout before the normal gate. Forwarded post-fence
// `--plugin-dir` is rejected because Codex has no such native session flag.
// `--skip-compat-check` alone bypasses the gate. `lookPath` resolves the safehouse
// binary (default exec.LookPath; injected so tests pin not-found).
func runCodex(ctx context.Context, args []string, dir string, ops hostOps, lookPath func(string) (string, error), stdout, stderr io.Writer) int {
	fd, err := parseFrontDoorArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock codex: %v\n", err)
		return 1
	}
	extra, err := safehouse.TranslateFlags(fd.safehouseFlags)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock codex: %v\n", err)
		return 1
	}
	// Codex has no native session-local plugin flag. Consume only Spacedock's
	// declared pre-fence override; reject a post-fence occurrence truthfully before
	// any persistent work. When there is no explicit override, use a qualified
	// checkout beside the resolved launcher. Explicit selection always wins.
	pluginDir, rest, explicitPluginDir := takeCodexPluginDir(fd.passthrough, fd.pluginDirPairs)
	fd.passthrough = rest
	if hasCodexCollaborationOverride(fd.passthrough) {
		fmt.Fprintln(stderr, "spacedock codex: collaboration settings are managed by Spacedock; remove the forwarded override")
		return 1
	}
	if hasPluginDir(fd.passthrough) {
		fmt.Fprintln(stderr, "spacedock codex: Codex does not accept forwarded --plugin-dir; place the Spacedock checkout flag before `--`, or install additional Codex plugins with `codex plugin add`")
		return 1
	}
	if explicitPluginDir && pluginDir == "" {
		fmt.Fprintln(stderr, "spacedock codex: --plugin-dir requires a checkout path")
		return 1
	}
	if !explicitPluginDir {
		pluginDir, _ = adjacentSpacedockPluginRoot("codex")
	}
	if pluginDir != "" {
		if err := installCodexLocalPluginDir(ops, pluginDir, stderr); err != nil {
			fmt.Fprintf(stderr, "spacedock codex: %v\n", err)
			return 1
		}
	}
	// Once Spacedock-owned plugin-dir handling has been consumed, suppress the
	// bootstrap posture only for Codex's exact resume command token. This is
	// deliberately not an option parser or command grammar.
	resume := slices.Contains(fd.passthrough, "resume")
	// The gate fails fast on too-old-binary, but the two healable verdicts
	// (NoPluginFound, TooOldPlugin) auto-install the codex plugin and proceed to
	// launch so the single command the user typed yields a working session —
	// `--no-install` opts out, preserving the refuse-and-instruct behavior. This
	// mirrors runClaude.
	if !fd.skipCheck {
		if !resolveHealableGate(ops, "codex", fd.noInstall, stderr) {
			return 1
		}
	}
	wrap := safehouse.Present(dir) || fd.forceSafehouse || len(fd.safehouseFlags) > 0
	if !resume {
		launchBanner("codex", dir, wrap, os.Getenv, lookPath, stderr)
	}
	inner := append([]string{"codex"}, codexCollaborationLayer...)
	if wrap {
		inner = append(inner, "--dangerously-bypass-approvals-and-sandbox")
	}
	// An unsandboxed launch has no safehouse isolation; codex has no single
	// auto-mode flag, so its nearest analog to claude's auto permission-mode is
	// `--ask-for-approval on-request` (the model decides when to escalate).
	// The sandboxed arm's bypass flag above already covers its posture, so this is
	// the !wrap counterpart. An exact resume keeps its established prompt-free
	// posture, and an operator-provided approval mode always wins.
	if !wrap && !resume && !passthroughHasFlag(fd.passthrough, "--ask-for-approval", "-a") {
		inner = append(inner, "--ask-for-approval", "on-request")
	}
	inner = append(inner, fd.passthrough...)
	if !resume {
		inner = append(inner, launchPrompt(codexBootstrapPrompt, fd))
	}

	argv := inner
	if wrap {
		if ok, hint := safehouse.Available(lookPath); !ok {
			fmt.Fprintln(stderr, hint)
			return 1
		}
		// Forward SPACEDOCK_BIN through safehouse's env sanitization with
		// `--env-pass`: launchEnv sets it on the safehouse process, this flag carries
		// it into the sandbox. It rides the safehouse flags (before `--`), so the
		// inner program stays `codex` and safehouse's program-keyed profile
		// auto-detection still fires. Omitted when the bin cannot be resolved.
		argv = safehouse.Wrap(inner, append(launcherBinEnvPassFlags(), extra...))
	}

	code, err := ops.Launch(argv, launchEnv(os.Environ()))
	if err != nil {
		fmt.Fprintf(stderr, "spacedock codex: launch failed: %v\n", err)
		return 1
	}
	return code
}

// claudeValueTakingHostFlags identifies Claude's space-form values for its
// advisory-only stray-prompt scan. Codex has no corresponding grammar table; its
// bootstrap posture is decided solely by exact `resume` token membership.
var claudeValueTakingHostFlags = map[string]bool{
	"-p": true, "--print": true,
	"--model":                true,
	"--mcp-config":           true,
	"--permission-mode":      true,
	"--add-dir":              true,
	"--append-system-prompt": true,
	"--settings":             true,
	"--session-id":           true,
	"--output-format":        true,
}

// skipInjectedPrefix returns the passthrough slice past the spacedock-injected
// leading `--plugin-dir <dir>` pairs. parseFrontDoorArgs re-prepends each
// before-`--` `--plugin-dir` as a `--plugin-dir <dir>` pair at the FRONT of
// fd.passthrough; that prefix is spacedock-owned, not operator after-`--` tokens,
// so Claude's advisory scan runs against the real after-`--` tokens. `--plugin-dir`
// is the only flag parseFrontDoorArgs re-prepends (the safehouse knobs live in
// fd.safehouseFlags, the booleans are consumed), so it is the complete injected
// prefix set. Skipping a leading `--plugin-dir <dir>` pair is correct regardless of
// origin: the dir is the flag's value, never a stray prompt.
func skipInjectedPrefix(passthrough []string) []string {
	for len(passthrough) >= 2 && passthrough[0] == "--plugin-dir" {
		passthrough = passthrough[2:]
	}
	return passthrough
}

// strayPromptAfterDash is a pure read over the parsed front-door args that returns
// the first stray positional in the after-`--` passthrough — an operator prompt
// that was placed after `--` and so silently degrades to host passthrough instead
// of being prepended to the spacedock launch prompt. It returns ("", false) when
// nothing is stray.
//
// It fires only when the operator gave no task before `--` (hasTask == false): a
// task before `--` means the operator already placed their prompt, so a positional
// after `--` is a deliberate host positional. The scan first skips the
// spacedock-injected leading `--plugin-dir <dir>` prefix, then runs every check
// against the real Claude after-`--` tokens. A token is a candidate when it is
// non-flag (does not start with `-`, and is not the bare `--` separator). A
// candidate is reported as stray only when we can be confident it is NOT a host
// flag's value:
//   - preceding a recognized value-taking host flag → it is that flag's value, skip;
//   - preceding an UNRECOGNIZED `-`-prefixed flag → it MIGHT be that flag's value, so
//     we suppress rather than give the actively-wrong "put X before --" advice;
//   - otherwise (first token, or preceded by a positional or recognized boolean
//     flag) → confidently stray.
//
// The check never alters the assembled argv — runClaude only writes the warning to
// stderr.
func strayPromptAfterDash(fd frontDoorArgs) (positional string, ok bool) {
	if fd.hasTask {
		return "", false
	}
	tokens := skipInjectedPrefix(fd.passthrough)
	for i, tok := range tokens {
		if tok == "--" || strings.HasPrefix(tok, "-") {
			continue
		}
		if i > 0 {
			prev := tokens[i-1]
			if claudeValueTakingHostFlags[prev] {
				continue // the value of a recognized value-taking host flag (space form)
			}
			// An equals-form flag (`--flag=value`) carries its own value, so it never
			// consumes the next token. Only a space-form `-`-prefixed flag is ambiguous:
			// it could be a value-taking flag we don't recognize, so suppress rather
			// than risk prescribing the wrong correction for what may be its value.
			if strings.HasPrefix(prev, "-") && prev != "--" && !strings.Contains(prev, "=") {
				continue
			}
		}
		return tok, true
	}
	return "", false
}

// frontDoorArgs is the parsed front-door grammar. The launchers read it to
// assemble the inner host argv and decide the safehouse wrap.
type frontDoorArgs struct {
	// passthrough is the host-only argv (claude/codex flags), in operator order.
	passthrough []string
	// pluginDirPairs counts the pre-fence --plugin-dir pairs re-injected at the
	// front of passthrough. It lets Codex consume only Spacedock-owned flags while
	// preserving post-fence argv exactly.
	pluginDirPairs int
	// task is the launch-prompt override (the bare text after the `--` fence);
	// hasTask distinguishes an explicit empty task from "no fence given".
	task    string
	hasTask bool
	// forceSafehouse is set by the bare `--safehouse` front-door flag.
	forceSafehouse bool
	// safehouseFlags are the de-prefixed `--safehouse-<key>=…` knob tokens, fed to
	// safehouse.TranslateFlags. Their presence also implies sandbox-on.
	safehouseFlags []string
	// skipCheck is set by `--skip-compat-check` (bypasses the version gate).
	skipCheck bool
	// noInstall is set by `--no-install` (opt out of the no-plugin auto-install,
	// preserving the refuse-and-instruct behavior).
	noInstall bool
}

// frontDoorFlags binds the spacedock-owned front-door flags onto a pflag.FlagSet
// so cobra owns their vocabulary natively: the three value-taking safehouse knobs
// are StringArray (accept both `--flag value` and `--flag=value`, accumulate on
// repeat), and the bare `--safehouse`/`--skip-compat-check` are Bool. The
// returned pointers are read back by parseFrontDoorArgs after Parse. The same
// binding feeds the per-command cobra help (AC-4), so the help and the parser
// never drift.
type frontDoorFlags struct {
	safehouse *bool
	skipCheck *bool
	noInstall *bool
	enable    *[]string
	addDirs   *[]string
	addDirsRO *[]string
	pluginDir *[]string
}

func bindFrontDoorFlags(fs *pflag.FlagSet) frontDoorFlags {
	return frontDoorFlags{
		safehouse: fs.Bool("safehouse", false,
			"Force the safehouse sandbox wrap even without a .safehouse profile in the directory"),
		skipCheck: fs.Bool("skip-compat-check", false,
			"Bypass the version gate and launch without resolving the installed plugin (bootstrap)"),
		noInstall: fs.Bool("no-install", false,
			"Do not auto-install the plugin when none is found; refuse and print install instructions instead"),
		enable: fs.StringArray("safehouse-enable", nil,
			"Enable a safehouse capability (KEY[,KEY]); repeatable; e.g. --safehouse-enable ssh,docker"),
		addDirs: fs.StringArray("safehouse-add-dirs", nil,
			"Grant safehouse read-write access to a directory; repeatable"),
		addDirsRO: fs.StringArray("safehouse-add-dirs-ro", nil,
			"Grant safehouse read-only access to a directory; repeatable"),
		pluginDir: fs.StringArray("plugin-dir", nil,
			"Select a local Spacedock checkout before -- (relaxes the version gate); repeatable"),
	}
}

// parseFrontDoorArgs parses the Option-2 front-door grammar in one pass via a
// pflag.FlagSet. cobra owns the spacedock flags wherever they appear before `--`;
// the non-flag positionals before `--` join (single space) into the launch task;
// everything after `--` forwards verbatim to the host as passthrough. This is the
// grammar inversion: host flags now ride AFTER `--`, the task before — the prompt
// is always spacedock-constructed and never adjacent to a value-taking host flag,
// so no dangling host flag can swallow it. The collected safehouse-knob values are
// re-prefixed to the `key=value` token form safehouse.TranslateFlags owns, so the
// safehouse vocabulary (the comma-split on enable, the unknown-key error) is
// unchanged — the space-form bug dies because cobra consumes the value as the
// flag's argument instead of leaking it to passthrough.
func parseFrontDoorArgs(args []string) (fd frontDoorArgs, err error) {
	fs := pflag.NewFlagSet("spacedock-front-door", pflag.ContinueOnError)
	fs.SetOutput(io.Discard)
	flags := bindFrontDoorFlags(fs)
	if err := fs.Parse(args); err != nil {
		return frontDoorArgs{}, err
	}

	fd.forceSafehouse = *flags.safehouse
	fd.skipCheck = *flags.skipCheck
	fd.noInstall = *flags.noInstall
	for _, v := range *flags.enable {
		fd.safehouseFlags = append(fd.safehouseFlags, "enable="+v)
	}
	for _, v := range *flags.addDirs {
		fd.safehouseFlags = append(fd.safehouseFlags, "add-dirs="+v)
	}
	for _, v := range *flags.addDirsRO {
		fd.safehouseFlags = append(fd.safehouseFlags, "add-dirs-ro="+v)
	}

	// ArgsLenAtDash is the count of positionals seen before `--` (or -1 when no
	// `--` was given). Without a `--`, every positional is the task and nothing
	// forwards. With a `--`, the pre-dash positionals join into the task and the
	// post-dash positionals forward verbatim as host passthrough.
	positionals := fs.Args()
	dash := fs.ArgsLenAtDash()
	var taskTokens []string
	if dash < 0 {
		taskTokens = positionals
	} else {
		taskTokens = positionals[:dash]
		fd.passthrough = positionals[dash:]
	}
	if len(taskTokens) > 0 {
		fd.task = strings.Join(taskTokens, " ")
		fd.hasTask = true
	}

	// --plugin-dir is the one host flag spacedock parses before `--`: pflag knows
	// its arity (one value, repeatable), so the dirs are captured correctly in
	// space/equals/repeated forms. Re-inject each as a `--plugin-dir <dir>` pair at
	// the FRONT of passthrough, ahead of any after-`--` tokens, and record that owned
	// prefix separately. Claude forwards it; Codex consumes only this declared prefix
	// into its local marketplace while preserving all post-fence argv verbatim.
	if dirs := *flags.pluginDir; len(dirs) > 0 {
		fd.pluginDirPairs = len(dirs)
		front := make([]string, 0, len(dirs)*2+len(fd.passthrough))
		for _, d := range dirs {
			front = append(front, "--plugin-dir", d)
		}
		fd.passthrough = append(front, fd.passthrough...)
	}
	return fd, nil
}
