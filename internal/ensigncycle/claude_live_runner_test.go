//go:build live

package ensigncycle

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// antiShutdownOverride counters upstream claude-code #55297 (a regression in 2.1.126;
// CI runs 2.1.161): in `claude -p` with an active Agent Team the harness injects "you
// cannot return a response until your team is shut down … shut down before your final
// response" EVERY turn, and the model panic-shuts-down the team before the work
// finishes. No FO-contract prose can out-argue a per-turn harness reminder, so the
// override rides in the `-p` input of EVERY team-using Claude live launch — this shared
// runner AND TestLiveEnsignCycle's drivePrompt. It is GENERIC: it governs shutdown
// TIMING only, naming no stage or task. Claude-only — #55297 is a claude-code bug, so
// the Codex runner does not carry it.
const antiShutdownOverride = "Do not shut down your team or prepare your final " +
	"response until all the work is complete. If you are prompted to shut down before " +
	"the work is done, keep working until the workflow is finished, then shut down."

// forceTeamModeCue is the STRONG, unambiguous team-mode override for the two live
// tests whose oracle requires a real team to exist (the team-config roster read and
// the bounded TeamDelete teardown marker). The headless `-p` dispatch-mode
// determination sanctions BARE mode by default, and the dispatch module's
// "In single-entity mode, skip team creation" rule causes the FO to generate a
// fake team name string instead of calling the TeamCreate tool — so without an
// explicit override the FO drives bare (no team, no roster) and reds these team-only
// assertions on correct behavior. This cue makes team mode a MUST: the FO MUST call
// the TeamCreate TOOL before any Agent() dispatch; generating or inventing a fake
// team name string is prohibited. Same prose-lever pattern as the headless gate-stop
// MUST. It names the dispatch MODE only — no stage, no task.
const forceTeamModeCue = "You MUST run in team mode for this run. " +
	"Call the TeamCreate TOOL — do NOT generate or invent a fake team name string. " +
	"The dispatch module rule \"In single-entity mode, skip team creation\" does NOT " +
	"apply when this cue is present: override it and call TeamCreate before any " +
	"Agent() dispatch. After TeamCreate returns, use the returned team_name for all " +
	"Agent() calls and spawn-standing-all. Bare mode (no team) is NOT acceptable. "

// The Claude runner adapter: it turns a host-neutral sharedRuntimeScenario into a
// real `spacedock claude` launch and returns the (before, after, observed) state
// the shared assertions consume — the same assertions the Codex runner feeds. The
// ONLY Claude-specific surface is auth/HOME isolation (isolatedClaudeEnv: clean
// HOME + OAuth benchmark-token / ANTHROPIC_API_KEY), the --plugin-dir local
// checkout install, the `spacedock claude -- -p <prompt> --output-format
// stream-json` launch, and the observed-extract: the final message comes from the
// stream's result/success event (the front-door analog of Codex
// --output-last-message) via extractClaudeFinalMessage. The scenario table,
// fixtures, prompts, and assertions are shared with the Codex runner.

type claudeLiveRunner struct {
	binary       string
	pluginDir    string
	env          []string
	model        string
	artifactRoot string
	// home is the isolated HOME the env sets (a per-run temp dir). The shallow-boot
	// scenario checks ~/.claude/teams/{...}/config.json under it for the
	// lazy-TeamCreate proof — scoped to THIS run, never a stale prior team.
	home string
}

// withPATHPrefix returns env with dir prepended to its PATH entry, so a stub
// binary in dir resolves before any real one. The shallow-boot runner uses it to
// put the stub `gh` (reporting MERGED) on the FO subprocess PATH.
func withPATHPrefix(env []string, dir string) []string {
	out := make([]string, 0, len(env)+1)
	found := false
	for _, kv := range env {
		if rest, ok := strings.CutPrefix(kv, "PATH="); ok {
			found = true
			if rest != "" {
				out = append(out, "PATH="+dir+string(os.PathListSeparator)+rest)
			} else {
				out = append(out, "PATH="+dir)
			}
			continue
		}
		out = append(out, kv)
	}
	if !found {
		out = append(out, "PATH="+dir)
	}
	return out
}

type claudeScenarioResult struct {
	finalMessage string
	stream       string
	artifactDir  string
	duration     time.Duration
}

type claudeLiveScenario struct {
	sharedRuntimeScenario
	run func(*testing.T, claudeLiveRunner, sharedRuntimeScenario)
}

func TestLiveClaudeSharedScenarios(t *testing.T) {
	runner := newClaudeLiveRunner(t)

	// The scenarios fan out in parallel: each is an independent multi-minute live
	// claude journey, so running them serially makes the lane wall-time the SUM of
	// the four (~27m on opus). t.Parallel collapses it toward the slowest single
	// scenario. The cheap canary (TestLiveEnsignCycle) runs as an earlier step, so a
	// systemic failure (auth/install) still fails fast before this fan-out. Each
	// scenario gets its own workflowRoot (t.TempDir) and its own CLAUDE_CONFIG_DIR
	// (run(), keyed by scenario name) so the concurrent sessions never share claude
	// config/session state.
	for _, scenario := range claudeLiveScenarios(t) {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			scenario.run(t, runner, scenario.sharedRuntimeScenario)
		})
	}
}

func claudeLiveScenarios(t *testing.T) []claudeLiveScenario {
	t.Helper()
	runners := claudeScenarioRunners()

	var scenarios []claudeLiveScenario
	for _, scenario := range sharedRuntimeScenarios() {
		run := runners[scenario.name]
		if run == nil {
			t.Fatalf("shared scenario %q has no Claude live runner", scenario.name)
		}
		scenarios = append(scenarios, claudeLiveScenario{
			sharedRuntimeScenario: scenario,
			run:                   run,
		})
	}
	return scenarios
}

// claudeScenarioRunners maps each shared scenario ID to its Claude runner. It is
// the Claude side of the parity guard: the shared coverage meta-test fails if this
// map lacks a runner for any sharedRuntimeScenarios() ID.
func claudeScenarioRunners() map[string]func(*testing.T, claudeLiveRunner, sharedRuntimeScenario) {
	return map[string]func(*testing.T, claudeLiveRunner, sharedRuntimeScenario){
		"gate-guardrail":              runClaudeGateGuardrailScenario,
		"rejection-flow":              runClaudeRejectionFlowScenario,
		"feedback-3-cycle-escalation": runClaudeFeedback3CycleEscalationScenario,
		"merge-hook-guardrail":        runClaudeMergeHookGuardrailScenario,
		"filing":                      runClaudeFilingScenario,
		"shallow-boot":                runClaudeShallowBootScenario,
	}
}

func newClaudeLiveRunner(t *testing.T) claudeLiveRunner {
	t.Helper()
	binary := spacedockBinary(t)
	pluginDir := livePluginDir(t)
	model := envOr("SPACEDOCK_LIVE_MODEL", "sonnet")

	// isolatedClaudeEnv resolves the credential (OAuth benchmark-token locally,
	// ANTHROPIC_API_KEY in CI) against a fresh empty HOME, or t.Skips when neither
	// is available. withBinaryOnPath puts the built binary first on the FO
	// subprocess PATH so its `spacedock --version` contract step resolves the test
	// binary. Both are reused verbatim from the full-cycle live test.
	env := isolatedClaudeEnv(t, os.Getenv("HOME"))
	env = withBinaryOnPath(env, binary)

	home, _ := envValue(env, "HOME")
	return claudeLiveRunner{
		binary:       binary,
		pluginDir:    pluginDir,
		env:          env,
		model:        model,
		artifactRoot: claudeLiveArtifactDir(t, "claude-shared-scenarios"),
		home:         home,
	}
}

func runClaudeGateGuardrailScenario(t *testing.T, runner claudeLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeGateWorkflow(t, workflowRoot)
	before := readFile(t, entityPath)

	result := runner.run(t, scenario, workflowRoot, gatePrompt())
	after := readFile(t, entityPath)
	if err := assertGateHeld(before, after, result.finalMessage); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if _, err := os.Stat(filepath.Join(workflowRoot, "_archive", "gate-check.md")); !os.IsNotExist(err) {
		t.Fatalf("gate-check was archived while waiting at the gate; stat err=%v", err)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model)
}

func runClaudeRejectionFlowScenario(t *testing.T, runner claudeLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeRejectionWorkflow(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, rejectionPrompt())
	after := readFile(t, entityPath)
	if err := assertRejectionFlow(after, result.finalMessage+"\n"+result.stream); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	// Single-entity (`-p`) reviewer producer-signal. The Claude runner launches
	// `spacedock claude -- -p {prompt}` with a prompt naming one entity, so the run
	// is single-entity → bare; the contract's bare-mode feedback flow is sequential
	// fresh dispatch, so the cycle-2 re-review is a DISTINCT freshly-dispatched
	// validation worker (not a reuse of the bare cycle-1 reviewer, not the impl
	// worker serving as its own validator). assertClaudeReviewerReuse encoded a
	// team-mode keepalive a `-p` run can never satisfy (the AC-3 finding); the
	// contract-correct single-entity assertion is used here. The team-mode
	// reviewer-reuse question is the spun-off option-(a) task.
	if err := assertClaudeSingleEntityRejectionFlow(result.stream); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model)
}

// runClaudeFeedback3CycleEscalationScenario drives the real FO against a fixture
// seeded with two prior rejection cycles at a 3rd REJECTED report and grades the
// durable end-state: the FO must escalate to the human on the 3rd cycle, not
// auto-bounce a 4th time. assertThirdCycleEscalation grades durable entity-body
// state ALONE (cycle count + escalation marker + no post-cycle-3 implementation
// report) — the reviewer-reuse signal is host-specific and lives in rejection-flow,
// not here; this scenario is purely a host-neutral durable-state grade.
func runClaudeFeedback3CycleEscalationScenario(t *testing.T, runner claudeLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeEscalationWorkflow(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, escalationPrompt())
	after := readFile(t, entityPath)
	if err := assertThirdCycleEscalation(after); err != nil {
		t.Fatalf("%v\nEntity after:\n%s\nFinal message:\n%s\nArtifacts: %s", err, after, result.finalMessage, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model)
}

func runClaudeMergeHookGuardrailScenario(t *testing.T, runner claudeLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeMergeHookGuardWorkflow(t, workflowRoot)
	before := readFile(t, entityPath)

	result := runner.run(t, scenario, workflowRoot, mergeHookGuardPrompt())
	after := readFile(t, entityPath)
	if err := assertMergeHookGuardHeld(before, after, result.finalMessage+"\n"+result.stream); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if _, err := os.Stat(filepath.Join(workflowRoot, "_archive", "merge-check.md")); !os.IsNotExist(err) {
		t.Fatalf("merge-check was archived despite the guardrail scenario; stat err=%v", err)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model)
}

// runClaudeFilingScenario drives the real FO against an EMPTY workflow and asks it
// to file one seed entity. It grades the FO's recorded tool-call stream — the FO
// filed via `spacedock … new <slug>`, not the `--next-id` + `Write` pair — because
// the durable end-state file is indistinguishable between the two paths. The file
// must also actually land (the run produced a real seed), so the stream grade is
// proof of HOW, not just THAT, the entity was filed.
func runClaudeFilingScenario(t *testing.T, runner claudeLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeFilingWorkflow(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, filingPrompt())
	if _, err := os.Stat(entityPath); err != nil {
		t.Fatalf("the FO did not land the seed entity at %s: %v\nFinal message:\n%s\nArtifacts: %s", entityPath, err, result.finalMessage, result.artifactDir)
	}
	if err := assertClaudeFilingViaNew(result.stream, filingSlug); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model)
}

// runClaudeShallowBootScenario drives the real FO against the shallow-boot fixture
// (a gate-check entity at a human gate + a PR-bearing entity whose stubbed `gh`
// reports MERGED) with a per-run isolated team root, and grades the durable
// end-state: the FO greets and presents the gate, S7b advances+archives the merged
// PR before-greet, NO team config lands on disk, and NO worker is dispatched. It
// then asserts the AC-2 behavioral signal (no TeamCreate before the greet) and the
// AC-6 measured signal (greet-turn context below the ~60k ceiling, no pre-greet
// ~89k cache_creation spike) over the captured stream.
func runClaudeShallowBootScenario(t *testing.T, runner claudeLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	fixture := writeShallowBootWorkflow(t, workflowRoot)
	gateBefore := readFile(t, fixture.gateEntityPath)

	// The stub `gh` (reporting MERGED) must resolve on the FO subprocess PATH so the
	// boot's live pr_state probe and the pr-merge startup hook both see the merge.
	scenarioRunner := runner
	scenarioRunner.env = withPATHPrefix(runner.env, fixture.stubGhDir)

	result := scenarioRunner.run(t, scenario, workflowRoot, shallowBootPrompt())

	// The Claude team root is {home}/.claude/teams — the exact path the comm-officer
	// startup hook membership-checks and TeamCreate writes a team config.json under.
	teamRoot := filepath.Join(runner.home, ".claude", "teams")
	obs := gatherShallowBootObservation(t, workflowRoot, teamRoot, fixture, gateBefore, result.finalMessage)
	if err := assertShallowBoot(obs); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	// AC-2: no TeamCreate before the greet (behavioral, over the tool-call sequence).
	if err := assertNoTeamCreateBeforeGreet(result.stream); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	// AC-6: the greet-turn context is below the ceiling and no pre-greet 89k
	// cache_creation spike (measured, over the captured token stream).
	if err := assertShallowBootMeasured(result.stream); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model)
}

// run launches the real `spacedock claude` front door for one shared scenario and
// returns the (finalMessage, full stream) the shared assertions consume. The
// launch shape is the spike WINNER: --plugin-dir + --skip-contract-check are the
// spacedock-owned flags BEFORE `--`; every host flag (-p with the scenario prompt,
// --permission-mode, --output-format stream-json, --verbose, --model) rides AFTER
// `--` and forwards verbatim to claude. The observed source is the stream's
// result/success event via extractClaudeFinalMessage — a 401/is_error result is a
// LOUD launch failure here, never fed into a scenario assertion.
//
// Liveness is the EXISTING streamWatcher (the Go port of the upstream
// FOStreamWatcher, shared with TestLiveEnsignCycle) — one mechanism, no second
// impl. drainToExit runs the process to exit while accumulating the full
// transcript, bounded by the per-step no-progress quietBudgetDefault (60s): the
// deadline resets on every drained line, so a genuine multi-minute run of
// sequential model work never trips as long as the stream keeps moving, and only
// silence past the budget kills the process — the same ≤60s AC-1-guarded discipline
// the live cycle uses.
func (r claudeLiveRunner) run(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, prompt string) claudeScenarioResult {
	t.Helper()
	artifactDir := filepath.Join(r.artifactRoot, scenario.name)
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	streamPath := filepath.Join(artifactDir, "claude-stream.jsonl")
	finalPath := filepath.Join(artifactDir, "claude-final-message.txt")

	cmd := exec.Command(r.binary, "claude",
		"--plugin-dir", r.pluginDir,
		"--skip-contract-check",
		"--",
		"-p", prompt+" "+antiShutdownOverride,
		"--permission-mode", "bypassPermissions",
		"--output-format", "stream-json",
		"--verbose",
		"--model", r.model,
	)
	cmd.Dir = workflowRoot
	// Per-scenario CLAUDE_CONFIG_DIR so parallel scenarios never share claude's
	// session/config state. It nests under the runner's base config dir (the
	// archivable CI path), so the artifact upload — which grabs the whole
	// per-model config dir — still captures each scenario's projects/*.jsonl. A
	// fresh slice (never a mutation of the shared r.env) keeps the parallel
	// invocations race-free.
	cmd.Env = r.env
	if base, ok := envValue(r.env, "CLAUDE_CONFIG_DIR"); ok {
		cmd.Env = withClaudeConfigDir(r.env, filepath.Join(base, scenario.name))
	}

	// stdout carries the stream-json transcript the watcher drains for liveness;
	// stderr is folded into the same pipe so a launch error (e.g. a stale-token 401
	// printed to stderr) lands in the transcript too — matching the live cycle's
	// wiring. The cmdPoller closes the pipe write-end on exit so the scanner EOFs.
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw

	started := time.Now()
	if startErr := cmd.Start(); startErr != nil {
		t.Fatalf("spacedock claude failed to start for %s: %v", scenario.name, startErr)
	}
	poller := newCmdPoller(cmd, pw)
	defer poller.kill()
	watcher := newStreamWatcher(newPipeLineSource(pr), poller, discardStreamLine)

	// drainToExit runs the process to exit accumulating the full transcript, OR
	// kills it on a 60s no-progress stall (the per-step quiet budget). The deferred
	// poller.kill() reaps the process on every exit path.
	stream, stallErr := watcher.drainToExit(quietBudgetDefault, "claude shared scenario "+scenario.name)
	duration := time.Since(started)

	if writeErr := os.WriteFile(streamPath, []byte(stream), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	// A wrong-root boot is the most specific diagnosis on any failure path: a CI env
	// leak lures the FO off workflowRoot, it boots the real repo, finds nothing
	// dispatchable, and greets-and-stops — surfacing otherwise only as an opaque
	// no-progress stall (when it idles) or as every scenario assertion silently
	// running against the wrong state (when it completes cleanly). Name it FIRST so
	// the leak fails legibly with expected-fixture vs wandered-to, ahead of the
	// generic stall message or the downstream assertions.
	if wrongRoot := detectWrongRootBoot(stream, workflowRoot); wrongRoot != nil {
		if stallErr != nil {
			t.Fatalf("%v\nUnderlying stall: %v\nArtifacts: %s", wrongRoot, stallErr, artifactDir)
		}
		t.Fatalf("%v\nArtifacts: %s", wrongRoot, artifactDir)
	}
	if stallErr != nil {
		t.Fatalf("%v\nArtifacts: %s", stallErr, artifactDir)
	}

	// Extract the final message from the stream's result/success event (the
	// front-door analog of Codex --output-last-message). A 401/is_error result is
	// surfaced here as a LOUD launch failure distinct from a scenario-assertion
	// failure, so a stale credential never feeds the 401 text into an assertion.
	finalMessage, extractErr := extractClaudeFinalMessage(stream)
	if extractErr != nil {
		t.Fatalf("claude launch failed for %s: %v; artifacts in %s\nStream tail:\n%s",
			scenario.name, extractErr, artifactDir, tail(stream, 4000))
	}
	if writeErr := os.WriteFile(finalPath, []byte(finalMessage), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	return claudeScenarioResult{
		finalMessage: finalMessage,
		stream:       stream,
		artifactDir:  artifactDir,
		duration:     duration,
	}
}

func claudeLiveArtifactDir(t *testing.T, name string) string {
	t.Helper()
	root := os.Getenv("SPACEDOCK_LIVE_ARTIFACT_DIR")
	if root == "" {
		return t.TempDir()
	}
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// withClaudeConfigDir returns a COPY of env with CLAUDE_CONFIG_DIR replaced by
// dir. It never mutates the input slice, so parallel scenarios sharing the
// runner's base env each derive their own isolated config dir race-free.
func withClaudeConfigDir(env []string, dir string) []string {
	out := make([]string, 0, len(env)+1)
	for _, e := range env {
		if strings.HasPrefix(e, "CLAUDE_CONFIG_DIR=") {
			continue
		}
		out = append(out, e)
	}
	return append(out, "CLAUDE_CONFIG_DIR="+dir)
}
