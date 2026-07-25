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
	binary          string
	pluginDir       string
	extraPluginDirs []string
	env             []string
	modelName       string
	artifactRoot    string
	// homeDir is the isolated HOME the env sets (a per-run temp dir). The
	// shallow-boot scenario checks ~/.claude/teams/{...}/config.json under it for the
	// lazy-TeamCreate proof — scoped to THIS run, never a stale prior team.
	homeDir string
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

// liveDriver is the transport seam: it turns a prompt + workflow root into the
// observed (finalMessage, stream) the shared assertions consume. The headless
// `-p` runner (claudeLiveRunner) and the pty/tmux runner (ptyLiveDriver) are two
// implementations; the scenario orchestration and assertions do not know which
// transport ran. model and home expose the two concrete facts the per-scenario
// orchestration needs beyond the launch itself — the metrics tag (model) and the
// isolated team root the shallow-boot scenario probes (home/.claude/teams).
// withStubPATH returns a driver copy whose launched FO subprocess resolves a stub
// binary in dir first (the shallow-boot scenario's stub `gh` reporting MERGED).
type liveDriver interface {
	run(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, prompt string) liveResult
	model() string
	home() string
	withStubPATH(dir string) liveDriver
}

// liveResult is the host-neutral observed state the shared assertions consume.
// headless `-p`: finalMessage is the stream's result/success event, stream is the
// stream-json transcript. pty/tmux: finalMessage is the FO-pane final text, stream
// is the session jsonl written under CLAUDE_CONFIG_DIR (the SAME stream-json
// dialect the assertions grade). The transport is invisible to the assertions.
type liveResult struct {
	finalMessage string
	stream       string
	artifactDir  string
	duration     time.Duration
	// configDir and cwd locate the dispatched-ensign sub-agent transcripts on disk
	// (under {configDir}/projects/{encode(cwd)}/{FO-session-id}/subagents), so the
	// journey-metrics fold can observe the ensign's --read adoption. cwd is the
	// EvalSymlinks-resolved FO working dir — the form Claude Code encodes into the
	// projects path. Empty on transports that do not record them (the pty driver),
	// so the fold no-ops to FO-front-door counts.
	configDir string
	cwd       string
}

type claudeLiveScenario struct {
	sharedRuntimeScenario
	run func(*testing.T, liveDriver, sharedRuntimeScenario)
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

// claudeScenarioRunners maps each shared scenario ID to its runner. The runners
// take the liveDriver seam, not a concrete runner, so the SAME map drives the
// headless `-p` runner and the pty/tmux driver. It is the parity guard: the shared
// coverage meta-test fails if this map lacks a runner for any
// sharedRuntimeScenarios() ID.
func claudeScenarioRunners() map[string]func(*testing.T, liveDriver, sharedRuntimeScenario) {
	return map[string]func(*testing.T, liveDriver, sharedRuntimeScenario){
		"gate-guardrail":                runClaudeGateGuardrailScenario,
		"recorded-gate-lifecycle":       runClaudeRecordedGateLifecycleScenario,
		"rejection-flow":                runClaudeRejectionFlowScenario,
		"feedback-3-cycle-escalation":   runClaudeFeedback3CycleEscalationScenario,
		"merge-hook-guardrail":          runClaudeMergeHookGuardrailScenario,
		"filing":                        runClaudeFilingScenario,
		"shallow-boot":                  runClaudeShallowBootScenario,
		"self-evidence-merge-triage":    runClaudeSelfEvidenceMergeTriageScenario,
		"smallest-sufficient-mechanism": runClaudeSmallestSufficientMechanismScenario,
		"keep-moving-posture":           runClaudeKeepMovingScenario,
	}
}

func runClaudeRecordedGateLifecycleScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	if copied, ok := runner.(claudeLiveRunner); ok {
		copied.pluginDir = t.TempDir()
		if err := copyTree(runner.(claudeLiveRunner).pluginDir, copied.pluginDir); err != nil {
			t.Fatal(err)
		}
		runner = copied
	}
	fixture := writePreparedRecordedGateFixture(t)
	before := readFile(t, fixture.entity)
	commandLog := filepath.Join(fixture.root, "evidence", "command.log")
	shimDir := writeRecordedGateLoggingShim(t, buildRecordedGateBinary(t), commandLog)
	shellEnvDir := t.TempDir()
	bashEnv := filepath.Join(shellEnvDir, "recorded-gate-env.sh")
	writeFile(t, bashEnv, "export SPACEDOCK_BIN="+filepath.Join(shimDir, "spacedock")+"\n")
	writeFile(t, filepath.Join(shellEnvDir, ".zshenv"), readFile(t, bashEnv))
	runner = runner.withStubPATH(shimDir)
	switch copied := runner.(type) {
	case claudeLiveRunner:
		copied.env = withRecordedGateEnv(copied.env, "BASH_ENV", bashEnv)
		copied.env = withRecordedGateEnv(copied.env, "ZDOTDIR", shellEnvDir)
		runner = copied
	case ptyLiveDriver:
		copied.env = withRecordedGateEnv(copied.env, "BASH_ENV", bashEnv)
		copied.env = withRecordedGateEnv(copied.env, "ZDOTDIR", shellEnvDir)
		runner = copied
	}
	result := runner.run(t, scenario, fixture.root, recordedGatePrompt(fixture.root))
	writeFile(t, filepath.Join(result.artifactDir, "command.log"), readFile(t, commandLog))
	commandLog = filepath.Join(result.artifactDir, "command.log")
	if strings.Contains(result.stream, "multiple runtime markers are set (CODEX_THREAD_ID, CLAUDECODE)") {
		t.Fatalf("recorded gate lifecycle hit mixed-marker ambiguity\nArtifacts: %s", result.artifactDir)
	}
	for _, line := range strings.Split(readFile(t, commandLog), "\n") {
		if strings.Contains(line, "\tdispatch build ") && (strings.Contains(line, " --host ") || strings.Contains(line, " --host=")) {
			t.Fatalf("recorded gate lifecycle recovered with an explicit host: %s\nArtifacts: %s", line, result.artifactDir)
		}
	}
	if copied, ok := runner.(claudeLiveRunner); ok && (!strings.Contains(result.stream, copied.pluginDir) || !strings.Contains(result.stream, "# First Officer Gate Lifecycle")) {
		t.Fatalf("recorded gate lifecycle did not load the copied skill body\nArtifacts: %s", result.artifactDir)
	}
	observation := recordedGateLiveObservation(t, fixture, before, commandLog, recordedGateReviewFromClaudeStream(result.stream))
	if err := assertRecordedGateLifecycle(observation); err != nil {
		t.Fatalf("recorded gate lifecycle graded FAIL: %v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model())
}

func newClaudeLiveRunner(t *testing.T) claudeLiveRunner {
	t.Helper()
	binary := buildRecordedGateBinary(t)
	pluginDir := livePluginDir(t)
	model := envOr("SPACEDOCK_LIVE_MODEL", "sonnet")

	// isolatedClaudeEnv resolves the credential (OAuth benchmark-token locally,
	// ANTHROPIC_API_KEY in CI) against a fresh empty HOME, or t.Skips when neither
	// is available. withBinaryOnPath puts the built binary first on the FO
	// subprocess PATH so its `spacedock --version` contract step resolves the test
	// binary. Both are reused verbatim from the full-cycle live test.
	env := isolatedClaudeEnv(t, os.Getenv("HOME"))
	env = withBinaryOnPath(env, binary)

	homeDir, _ := envValue(env, "HOME")
	return claudeLiveRunner{
		binary:       binary,
		pluginDir:    pluginDir,
		env:          env,
		modelName:    model,
		artifactRoot: claudeLiveArtifactDir(t, "claude-shared-scenarios"),
		homeDir:      homeDir,
	}
}

// claudeLiveRunner satisfies the liveDriver seam: the shared per-scenario
// orchestration drives it through the interface, oblivious to the `-p` transport.
var _ liveDriver = claudeLiveRunner{}

func (r claudeLiveRunner) model() string { return r.modelName }
func (r claudeLiveRunner) home() string  { return r.homeDir }

// withStubPATH returns a runner copy whose launched FO subprocess resolves a stub
// binary in dir first (the shallow-boot scenario's stub `gh` reporting MERGED). It
// never mutates the receiver's env, so parallel scenarios sharing the runner stay
// race-free.
func (r claudeLiveRunner) withStubPATH(dir string) liveDriver {
	r.env = withPATHPrefix(r.env, dir)
	return r
}

func (r claudeLiveRunner) withExtraPluginDir(dir string) claudeLiveRunner {
	r.extraPluginDirs = append(append([]string(nil), r.extraPluginDirs...), dir)
	return r
}

func runClaudeGateGuardrailScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	fixture := writeGateWorkflow(t, workflowRoot)
	if scenario.name == "default-headless-recorded-gate-stop" {
		writeFile(t, filepath.Join(fixture.root, "README.md"), strings.Replace(recordedGateReadme(), "### validation", "### implementation\n\nAppend an implementation stage report, then return completion.\n\n### validation", 1))
		writeFile(t, fixture.entity, strings.Replace(recordedGateEntity(), "status: validation", "status: implementation", 1))
		gitCommitPathScoped(t, fixture.stateRoot, "recorded-gate-task/index.md", "start before gate")
	}
	before := readFile(t, fixture.entity)
	commandLog := filepath.Join(fixture.root, "evidence", "command.log")
	shimDir := writeRecordedGateLoggingShim(t, buildRecordedGateBinary(t), commandLog)
	runner = runner.withStubPATH(shimDir)

	result := runner.run(t, scenario, workflowRoot, gatePrompt(workflowRoot))
	if _, err := os.Stat(filepath.Join(fixture.stateRoot, "_archive", "recorded-gate-task", "index.md")); !os.IsNotExist(err) {
		t.Fatalf("recorded-gate-task was archived while waiting at the gate; stat err=%v", err)
	}
	after := readFile(t, fixture.entity)
	if err := assertGateHeld(before, after, recordedGateReviewFromClaudeStream(result.stream)); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if err := assertRecordedGateHoldLog(readFile(t, commandLog)); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model())
}

func runClaudeRejectionFlowScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeRejectionWorkflow(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, rejectionPrompt(workflowRoot))
	after := readFile(t, entityPath)
	if err := assertRejectionFlow(after, result.finalMessage+"\n"+result.stream); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if err := assertRejectionRecordedRound(workflowRoot, entityPath, "validation", claudeRecordedRejectionRound(result.stream)); err != nil {
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
	emitClaudeScenarioMetrics(t, scenario, result, runner.model())
}

// runClaudeFeedback3CycleEscalationScenario drives the real FO against a fixture
// seeded with two prior rejection cycles at a 3rd REJECTED report and grades the
// durable end-state: the FO must escalate to the human on the 3rd cycle, not
// auto-bounce a 4th time. assertThirdCycleEscalation grades durable entity-body
// state ALONE (cycle count + escalation marker + no post-cycle-3 implementation
// report) — the reviewer-reuse signal is host-specific and lives in rejection-flow,
// not here; this scenario is purely a host-neutral durable-state grade.
func runClaudeFeedback3CycleEscalationScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeEscalationWorkflow(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, escalationPrompt(workflowRoot))
	after := readFile(t, entityPath)
	if err := assertThirdCycleEscalation(after); err != nil {
		t.Fatalf("%v\nEntity after:\n%s\nFinal message:\n%s\nArtifacts: %s", err, after, result.finalMessage, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model())
}

func runClaudeMergeHookGuardrailScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeMergeHookGuardWorkflow(t, workflowRoot)
	before := readFile(t, entityPath)

	result := runner.run(t, scenario, workflowRoot, mergeHookGuardPrompt(workflowRoot))
	after := readFile(t, entityPath)
	if err := assertMergeHookGuardHeld(before, after, result.finalMessage+"\n"+result.stream); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if _, err := os.Stat(filepath.Join(workflowRoot, "_archive", "merge-check.md")); !os.IsNotExist(err) {
		t.Fatalf("merge-check was archived despite the guardrail scenario; stat err=%v", err)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model())
}

// runClaudeSelfEvidenceMergeTriageScenario drives the real FO against the
// merge/triage fixture (a diff touching a live-lane-exercised path, the required lane
// unapproved, a prior-session handoff mislabelling this run's live-CI red) and grades
// the FO's OWN decision against the INCIDENT class only (Cycle-3): it must NOT wave the
// change through — terminalize-as-PASSED or merge-without-evidence — while the required
// lane is unapproved; a hold or a reject-without-merge both pass. Its final-message
// diagnosis must name THIS run's failing test, not the inherited "known flake" label.
// The this-run token is graded in the FINAL MESSAGE only — the fixture body carries it
// so the FO can read it, so grading the transcript would pass on the entity-read echo.
func runClaudeSelfEvidenceMergeTriageScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeMergeTriageWorkflow(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, mergeTriagePrompt(workflowRoot))
	after := readMergeTriageAfter(t, workflowRoot, entityPath)
	if err := assertSelfEvidenceMergeTriage(after, result.finalMessage); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model())
}

// runClaudeSmallestSufficientMechanismScenario drives the real FO against the
// smallest-sufficient-mechanism fixture (a commissioned workflow with two ready
// entities plus two plain deterministic-edit notes whose content the prompt hands the
// FO) and grades the FO's tool-call STREAM in both directions of the ladder: the
// deterministic edits are FO-authored with a direct commit and NO worker/PR climb, and
// the commissioned ready entities are engaged via the standing dispatch loop WITHOUT a
// per-entity justification. The trace is graded, not the durable end-state, which is
// identical whether the FO climbed or not.
func runClaudeSmallestSufficientMechanismScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	writeSmallestMechanismWorkflow(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, smallestMechanismPrompt(workflowRoot))
	if err := assertClaudeSmallestSufficientMechanism(result.stream, ssmEditFiles(), ssmCommissioned()); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model())
}

// runClaudeKeepMovingScenario drives the real FO against the keep-moving-posture fixture
// (four independent entities: a just-approved gate, two independent ready entities, and a
// questioned entity) and grades the FO's tool-call STREAM plus its FINAL MESSAGE against
// the four false-stop patterns: it advances + dispatches the approved entity with no
// permission question, dispatches both independent entities, re-shapes the questioned
// entity and pauses only its dispatch, and does not end its turn on an async wait. The
// durable end-state cannot distinguish keep-moving from a false stop, so the motion trace
// (actions) and the turn-ending postures (final message) are graded, not the entity files.
func runClaudeKeepMovingScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	writeKeepMovingWorkflow(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, keepMovingPrompt(workflowRoot))
	if err := assertClaudeKeepMoving(result.stream, result.finalMessage, kmIndependent()); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model())
}

// runClaudeFilingScenario drives the real FO against an EMPTY workflow and asks it
// to file one seed entity. It grades the FO's recorded tool-call stream — the FO
// filed via `spacedock … new <slug>`, not the `--next-id` + `Write` pair — because
// the durable end-state file is indistinguishable between the two paths. The file
// must also actually land (the run produced a real seed), so the stream grade is
// proof of HOW, not just THAT, the entity was filed.
func runClaudeFilingScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeFilingWorkflow(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, filingPrompt(workflowRoot))
	if _, err := os.Stat(entityPath); err != nil {
		t.Fatalf("the FO did not land the seed entity at %s: %v\nFinal message:\n%s\nArtifacts: %s", entityPath, err, result.finalMessage, result.artifactDir)
	}
	if err := assertClaudeFilingViaNew(result.stream, filingSlug); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model())
}

// runClaudeShallowBootScenario drives the real FO against the shallow-boot fixture
// (one gate-check entity at a human gate) with a per-run isolated team root, and
// grades the durable end-state: the FO greets with the accurate held-gate state,
// NO entity mutation occurs, NO team config lands on disk, and no durable dispatch
// fingerprint appears. It then asserts the AC-2 behavioral signal (no TeamCreate
// before the greet) and records the AC-6 measured signal over the captured stream.
func runClaudeShallowBootScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	fixture := writeShallowBootWorkflow(t, workflowRoot)
	gateBefore := readFile(t, fixture.gateEntityPath)

	result := runner.run(t, scenario, workflowRoot, shallowBootPrompt(workflowRoot))

	// The Claude team root is {home}/.claude/teams — the exact path the comm-officer
	// startup hook membership-checks and TeamCreate writes a team config.json under.
	teamRoot := filepath.Join(runner.home(), ".claude", "teams")
	obs := gatherShallowBootObservation(t, workflowRoot, teamRoot, fixture, gateBefore, result.finalMessage)
	if err := assertShallowBoot(obs); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	// AC-2: no TeamCreate before the greet (behavioral, over the tool-call sequence).
	if err := assertNoTeamCreateBeforeGreet(result.stream); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	// The boot-window oracle: a greet turn was produced (structural only — the
	// former ~60k ceiling/spike thresholds no longer gate CI, see
	// assertShallowBootMeasuredTurns).
	if err := assertShallowBootMeasured(result.stream); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	// Record (don't gate on) the greet turn's full token usage as a distinct
	// shallow-boot-window observation, riding the same journeymetrics ledger pipe
	// emitClaudeScenarioMetrics below already uses.
	emitShallowBootWindowMetrics(t, result.stream, runner.model())
	emitClaudeScenarioMetrics(t, scenario, result, runner.model())
}

// run launches the real `spacedock claude` front door for one shared scenario and
// returns the (finalMessage, full stream) the shared assertions consume. The
// launch shape is the spike WINNER: --plugin-dir + --skip-compat-check are the
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
func (r claudeLiveRunner) run(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, prompt string) liveResult {
	t.Helper()
	artifactDir := filepath.Join(r.artifactRoot, scenario.name)
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	streamPath := filepath.Join(artifactDir, "claude-stream.jsonl")
	finalPath := filepath.Join(artifactDir, "claude-final-message.txt")

	argv := []string{"claude",
		"--plugin-dir", r.pluginDir,
		"--skip-compat-check",
		"--"}
	for _, dir := range r.extraPluginDirs {
		argv = append(argv, "--plugin-dir", dir)
	}
	argv = append(argv,
		"-p", prompt+" "+antiShutdownOverride,
		"--permission-mode", "bypassPermissions",
		"--output-format", "stream-json",
		"--verbose",
		"--model", r.modelName,
	)
	cmd := exec.Command(r.binary, argv...)
	cmd.Dir = workflowRoot
	// Per-scenario CLAUDE_CONFIG_DIR so parallel scenarios never share claude's
	// session/config state. It nests under the runner's base config dir (the
	// archivable CI path), so the artifact upload — which grabs the whole
	// per-model config dir — still captures each scenario's projects/*.jsonl. A
	// fresh slice (never a mutation of the shared r.env) keeps the parallel
	// invocations race-free.
	cmd.Env = r.env
	configDir, _ := envValue(r.env, "CLAUDE_CONFIG_DIR")
	if base, ok := envValue(r.env, "CLAUDE_CONFIG_DIR"); ok {
		configDir = filepath.Join(base, scenario.name)
		cmd.Env = withClaudeConfigDir(r.env, configDir)
	}
	if err := seedStoredLoginCredential(configDir); err == nil {
		cmd.Env = withoutEnvKey(cmd.Env, "CLAUDE_CODE_OAUTH_TOKEN")
	}

	// The resolved cwd is what Claude Code encodes into its projects path; the FO
	// subprocess runs in workflowRoot, so resolve its symlinks (macOS t.TempDir is
	// under /var -> /private/var) to match the on-disk subagents dir.
	resolvedCwd := workflowRoot
	if resolved, err := filepath.EvalSymlinks(workflowRoot); err == nil {
		resolvedCwd = resolved
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

	// Wrong-root and broad-search observations come from the full transcript, so
	// they cannot reliably classify a boot phase or decide pass/fail. Retain the
	// evidence only as secondary context: one cleanup reports it after an existing
	// runner or scenario failure and remains silent when the shared scenario passes.
	registerClaudeLiveFailureDiagnostic(t, detectClaudeLiveFailureDiagnostic(stream, workflowRoot))
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

	return liveResult{
		finalMessage: finalMessage,
		stream:       stream,
		artifactDir:  artifactDir,
		duration:     duration,
		configDir:    configDir,
		cwd:          resolvedCwd,
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

// withoutEnvKey returns env with every KEY=... entry for key removed. The pty
// driver uses it to drop CLAUDE_CODE_OAUTH_TOKEN once it has seeded the stored-login
// credential, so the seeded login is the child's only (and authoritative) credential.
func withoutEnvKey(env []string, key string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env))
	for _, e := range env {
		if strings.HasPrefix(e, prefix) {
			continue
		}
		out = append(out, e)
	}
	return out
}
