//go:build live

package ensigncycle

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/gates"
)

// antiShutdownOverride counters upstream claude-code #55297 (a regression in 2.1.126;
// CI runs 2.1.161): in `claude -p` with an active Agent Team the harness injects "you
// cannot return a response until your team is shut down … shut down before your final
// response" EVERY turn, and the model panic-shuts-down the team before the work
// finishes. No FO-contract prose can out-argue a per-turn harness reminder, so the
// override rides in the `-p` input of EVERY team-using Claude live launch — this shared
// runner AND TestLiveCommonFullEnsignCycle's drivePrompt. It is GENERIC: it governs shutdown
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
// --output-last-message) via extractClaudeFinalMessage. The common declarations,
// fixtures, prompts, and assertions are shared with the Codex runner.

type claudeLiveRunner struct {
	t            *testing.T
	binary       string
	pluginDir    string
	env          []string
	modelName    string
	artifactRoot string
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

// liveDriver turns a prompt and workflow root into the observed final message
// and stream that the shared assertions consume. model and home expose the two
// facts the per-scenario orchestration needs beyond the headless launch.
// withStubPATH returns a driver copy whose launched FO subprocess resolves a stub
// binary in dir first (the shallow-boot scenario's stub `gh` reporting MERGED).
//
// lifecycleStream returns the stream in which THIS host records worker spawns and
// completions, which is not always the public one: Codex reports its sub-agent
// lifecycle only in the CODEX_HOME rollout. It is a required method rather than an
// optional interface a driver may assert into because omission was the defect — an
// optional interface is satisfiable by forgetting it, so a host that skipped it
// graded green with the dispatch-evidence check silently not running. A missing
// method is now a build error.
type liveDriver interface {
	run(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, prompt string) liveResult
	emitMetrics(t *testing.T, scenario sharedRuntimeScenario, result liveResult)
	gradeShallowBootObservation(t *testing.T, result liveResult)
	prepareRecordedGate(t *testing.T) (liveDriver, func(liveResult))
	smallestMechanismTrace(result liveResult, edits, commissioned []string) mechanismTrace
	lifecycleStream(t *testing.T, result liveResult) string
	model() string
	home() string
	withStubPATH(dir string) liveDriver
}

// liveResult is the host-neutral observed state the shared assertions consume.
// finalMessage is the headless stream's result/success event and stream is its
// stream-json transcript.
type liveResult struct {
	finalMessage string
	stream       string
	commands     []string
	artifactDir  string
	duration     time.Duration
	// configDir and cwd locate the dispatched-ensign sub-agent transcripts on disk
	// (under {configDir}/projects/{encode(cwd)}/{FO-session-id}/subagents), so the
	// journey-metrics fold can observe the ensign's --read adoption. cwd is the
	// EvalSymlinks-resolved FO working dir — the form Claude Code encodes into the
	// projects path.
	configDir string
	cwd       string
}

func runClaudeRecordedGateLifecycleScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario, build func(*testing.T) recordedGateFixture, assert func(recordedGateObservation) error) {
	t.Helper()
	runner, gradePreparation := runner.prepareRecordedGate(t)
	fixture := build(t)
	before := readFile(t, fixture.entity)
	commandLog := filepath.Join(fixture.root, "evidence", "command.log")
	shimDir := writeRecordedGateLoggingShim(t, buildRecordedGateBinary(t), commandLog)
	runner = runner.withStubPATH(shimDir)
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
	gradePreparation(result)
	observation := recordedGateLiveObservation(t, fixture, before, commandLog)
	finishLiveScenario(t, runner, scenario, result,
		durableSemantic("recorded-gate-lifecycle-violation", assert(observation)))
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
		t:            t,
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
func (r claudeLiveRunner) smallestMechanismTrace(result liveResult, edits, commissioned []string) mechanismTrace {
	return smallestMechanismTraceForDialect("claude", result.stream, edits, commissioned)
}
func (r claudeLiveRunner) home() string { return r.homeDir }

// Claude records the Agent spawn and its task_notification completion in the public
// stream-json transcript, so the public stream IS the lifecycle stream.
func (r claudeLiveRunner) lifecycleStream(_ *testing.T, result liveResult) string {
	return result.stream
}
func (r claudeLiveRunner) emitMetrics(t *testing.T, scenario sharedRuntimeScenario, result liveResult) {
	emitClaudeScenarioMetrics(t, scenario, result, r.modelName)
}
func (r claudeLiveRunner) gradeShallowBootObservation(t *testing.T, result liveResult) {
	t.Helper()
	if err := assertNoTeamCreateBeforeGreet(result.stream); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	if err := assertShallowBootMeasured(result.stream); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	emitShallowBootWindowMetrics(t, result.stream, r.modelName)
}
func (r claudeLiveRunner) prepareRecordedGate(t *testing.T) (liveDriver, func(liveResult)) {
	source := r.pluginDir
	r.pluginDir = t.TempDir()
	if err := copyTree(source, r.pluginDir); err != nil {
		t.Fatal(err)
	}
	return r, func(result liveResult) {
		if !strings.Contains(result.stream, r.pluginDir) || !strings.Contains(result.stream, "# First Officer Gate Lifecycle") {
			t.Fatalf("recorded gate lifecycle did not load the copied skill body\nArtifacts: %s", result.artifactDir)
		}
	}
}

// withStubPATH returns a runner copy whose launched FO subprocess resolves a stub
// binary in dir first (the shallow-boot scenario's stub `gh` reporting MERGED). It
// never mutates the receiver's env, so parallel scenarios sharing the runner stay
// race-free.
func (r claudeLiveRunner) withStubPATH(dir string) liveDriver {
	r.env = withPATHPrefix(r.env, dir)
	r.env = withSpacedockShimShellEnv(r.t, r.env, dir)
	return r
}

func finishLiveScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario, result liveResult, semantic ...error) {
	t.Helper()
	scenario.grade = gradeLive(scenario.gap.kind == "xfail", semantic...)
	runner.emitMetrics(t, scenario, result)
	if scenario.grade.status == "xfail" {
		t.Logf("XFAIL %s/%s owner=%s observed=%v", scenario.gap.target, scenario.name, scenario.gap.owner, scenario.grade.codes)
	}
	if scenario.grade.status == "xpass" {
		t.Logf("XPASS ALERT %s/%s owner=%s observed=%v", scenario.gap.target, scenario.name, scenario.gap.owner, scenario.grade.codes)
	}
	if liveGradeFailsLane(scenario.grade.status) {
		// The durable end state each finding graded lives in a t.TempDir that is
		// gone by the time anyone reads CI, so the messages are the only surviving
		// account of what the codes mean.
		details := ""
		if len(scenario.grade.details) > 0 {
			details = "\nFindings:\n  " + strings.Join(scenario.grade.details, "\n  ")
		}
		t.Fatalf("%s %s/%s owner=%s observed=%v%s\nFinal message:\n%s\nArtifacts: %s", strings.ToUpper(scenario.grade.status), scenario.gap.target, scenario.name, scenario.gap.owner, scenario.grade.codes, details, result.finalMessage, result.artifactDir)
	}
}

func runGateStopScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario, build func(*testing.T, string) recordedGateFixture, assert func(string, string, gateHeldExpectation) error) {
	t.Helper()
	workflowRoot := t.TempDir()
	fixture := build(t, workflowRoot)
	before := readFile(t, fixture.entity)
	observerRoot := t.TempDir()
	if err := os.Symlink(fixture.stateRoot, filepath.Join(observerRoot, ".spacedock-state")); err != nil {
		t.Fatal(err)
	}
	commandLog := filepath.Join(observerRoot, "command.log")
	if err := assertObserverOutsideWorkflow(workflowRoot, commandLog); err != nil {
		t.Fatal(err)
	}
	shimDir := writeRecordedGateLoggingShim(t, buildRecordedGateBinary(t), commandLog)
	runner = runner.withStubPATH(shimDir)

	result := runner.run(t, scenario, workflowRoot, gatePrompt(workflowRoot))
	writeFile(t, filepath.Join(result.artifactDir, "command.log"), readFile(t, commandLog))
	if _, err := os.Stat(filepath.Join(fixture.stateRoot, "_archive", "recorded-gate-task", "index.md")); !os.IsNotExist(err) {
		t.Fatalf("recorded-gate-task was archived while waiting at the gate; stat err=%v", err)
	}
	after := readFile(t, fixture.entity)
	var semantic []error
	expected, err := semanticGateHeldExpectation(fixture)
	if err != nil {
		semantic = append(semantic, err)
	}
	if err := assert(before, after, expected); err != nil {
		semantic = append(semantic, &gradedErr{code: "gate-not-held", msg: err.Error()})
	}
	semantic = append(semantic, assertRecordedGateHoldLog(readFile(t, commandLog)))
	if scenario.name == "default-headless-gate-stop" {
		semantic = append(semantic, assertImplementationWorkerLifecycle(nativeLifecycleStream(t, runner, result), after))
	}
	finishLiveScenario(t, runner, scenario, result, semantic...)
}

func nativeLifecycleStream(t *testing.T, runner liveDriver, result liveResult) string {
	t.Helper()
	stream, err := codexNativeLifecycleStream(runner.home(), result.stream)
	if err != nil {
		t.Fatal(err)
	}
	return stream
}
func runClaudeWithdrawnGateRecoveryScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario, build func(*testing.T, string) recordedGateFixture, assert func(*gates.Document) error) {
	t.Helper()
	workflowRoot := t.TempDir()
	fixture := build(t, workflowRoot)
	binary := buildRecordedGateBinary(t)
	commitRecordedGateState(t, binary, fixture, "commit selected gate inputs")
	prepared := mustRecordedGate(t, binary, fixture.root,
		"gate", "prepare", "recorded-gate-task",
		"--question", "Should the stale candidate advance?",
		"--artifact", fixture.gateReview,
		"--summary", "Stale candidate.",
		"--reference", fixture.references[0],
		"--reference", fixture.references[1],
		"--workflow-dir", fixture.root)
	firstRoom := outputValue(prepared.stdout, "room")
	firstBriefing := readFile(t, filepath.Join(firstRoom, "gate-briefing.json"))
	firstRequest := readFile(t, filepath.Join(firstRoom, "request.json"))
	commitRecordedGateState(t, binary, fixture, "prepare stale attempt")
	mustRecordedGate(t, binary, fixture.root,
		"gate", "withdraw", "recorded-gate-task",
		"--reason", "Sprint re-scope replaced the reviewed candidate.",
		"--workflow-dir", fixture.root)
	commitRecordedGateState(t, binary, fixture, "withdraw stale attempt")

	commandLog := filepath.Join(fixture.root, "evidence", "command.log")
	shimDir := writeRecordedGateLoggingShim(t, binary, commandLog)
	runner = runner.withStubPATH(shimDir)
	result := runner.run(t, scenario, workflowRoot, gatePrompt(workflowRoot))

	doc, _, err := gates.Read(fixture.entity)
	if err != nil {
		t.Fatalf("read recovered entity: %v\nArtifacts: %s", err, result.artifactDir)
	}
	if err := assert(doc); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	current := doc.Records[0].Attempts[1]
	if readFile(t, filepath.Join(firstRoom, "gate-briefing.json")) != firstBriefing ||
		readFile(t, filepath.Join(firstRoom, "request.json")) != firstRequest {
		t.Fatalf("recovery rewrote withdrawn room bytes\nArtifacts: %s", result.artifactDir)
	}
	secondRoom := filepath.Join(filepath.Dir(fixture.entity), filepath.FromSlash(current.Briefing.RoomRef))
	entries, err := os.ReadDir(secondRoom)
	if err != nil || len(entries) != 2 {
		t.Fatalf("successor room is not the emitted two-file room: entries=%v err=%v\nArtifacts: %s", entries, err, result.artifactDir)
	}
	log := readFile(t, commandLog)
	const prepare = "exit=0\tgate prepare recorded-gate-task "
	if strings.Count(log, prepare) != 1 ||
		strings.Count(log, "exit=0\tstate commit recorded-gate-task") != 1 ||
		strings.Contains(log, "gate record recorded-gate-task") ||
		strings.Contains(log, "gate consume recorded-gate-task") ||
		strings.Contains(log, "dispatch build ") {
		t.Fatalf("withdrawn recovery crossed its prepare/commit gate-stop boundary\n%s\nArtifacts: %s", log, result.artifactDir)
	}
	if !validatingStatus.MatchString(readFile(t, fixture.entity)) {
		t.Fatalf("withdrawn recovery changed workflow status\nArtifacts: %s", result.artifactDir)
	}
	runner.emitMetrics(t, scenario, result)
}

func assertWithdrawnGateRecovery(doc *gates.Document) error {
	if len(doc.Records) != 1 || len(doc.Records[0].Attempts) != 2 {
		return fmt.Errorf("recovery attempts = %#v", doc.Records)
	}
	withdrawn, current := doc.Records[0].Attempts[0], doc.Records[0].Attempts[1]
	if withdrawn.Withdrawal == nil || withdrawn.Resolution != nil || withdrawn.Application != nil {
		return fmt.Errorf("withdrawn attempt lost its clean terminal state: %#v", withdrawn)
	}
	if current.Withdrawal != nil || current.Resolution != nil || current.Application != nil || !strings.HasSuffix(current.ID, "-2") {
		return fmt.Errorf("recovery did not stop on open successor N+1: %#v", current)
	}
	return nil
}

// codexRejectionBranch is the Codex branch key, and it is a CONSTANT by derivation
// rather than for convenience: `«context-budget»` is ABSENT on Codex
// (`codex-first-officer-runtime.md:28`), so reuse condition 0 is satisfied by
// definition, and `followup_task` is the host's live reuse-advance handle, so
// condition 1 holds for any worker the run already opened. The reuse route is
// therefore always available and the FO always owes the reuse chain. Deriving the
// branch from whether a run HAPPENED to call followup_task would bless FO residual
// mode 2 — fresh-dispatching the fix target while the followup route is live — by
// re-grading that very deviation as the fresh branch.
const codexRejectionBranch = rejectionBranchReuse

// rejectionHostRealization is the host half of the team-mode invocation. The shared
// prompt carries the host-neutral mode requirement; this names the concrete calls
// THIS host opens and re-routes workers with. Same split as
// merged_team_mode_live_test.go's forceMergedTeamCue, appended by the runner the way
// antiShutdownOverride already is.
func rejectionHostRealization(runner liveDriver) string {
	if _, ok := runner.(codexAsLiveDriver); ok {
		return "On this host that means `spawn_agent` to open each worker, and `followup_task` to route follow-up work to a worker that is still live."
	}
	return "On this host that means an Agent with a name set and run_in_background true to open each worker, and SendMessage to route follow-up work to a worker that is still live."
}

// writeRejectionTopologyDigest persists the run's branch chain beside its other
// artifacts. AC-1 counts a conforming green only when the focused test exits 0 AND
// this file shows that run's chain in order, so it is written on every run, pass or
// fail — the in-process extraction dies with the isolated CODEX_HOME that t.Cleanup
// removes.
func writeRejectionTopologyDigest(t *testing.T, artifactDir string, branch rejectionBranch, routes []rejectionRoute) {
	t.Helper()
	if artifactDir == "" {
		t.Fatal("rejection-flow run exposed no artifact dir to persist the topology digest into")
	}
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatalf("create rejection topology artifact dir: %v", err)
	}
	digest := rejectionTopologyDigest(branch, routes)
	path := filepath.Join(artifactDir, "rejection-topology.tsv")
	if err := os.WriteFile(path, []byte(digest), 0o644); err != nil {
		t.Fatalf("persist rejection topology digest: %v", err)
	}
	t.Logf("rejection topology digest %s:\n%s", path, digest)
}

func runClaudeRejectionFlowScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario, build func(*testing.T, string) string, assert func(string, string) error) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := build(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, rejectionPrompt(workflowRoot)+"\n"+rejectionHostRealization(runner))
	after := readFile(t, entityPath)
	recordedRound := claudeRecordedRejectionRound(result.stream)
	// The publication counter reads the same run stream, so it grades the FO
	// behavior that was observed, not the wording of the skill that produced it.
	publications := claudeRejectionRoundPublications(result.stream)
	// Worker topology comes from each host's NATIVE transcript: the Claude
	// stream-json spawns/notifications, and for Codex the parent rollout, because the
	// public `codex exec --json` stream carries only `wait` collab items and no
	// topology at all.
	routes, branch := claudeRejectionRoutes(result.stream)
	if _, ok := runner.(codexAsLiveDriver); ok {
		recordedRound = codexRecordedRejectionRound(result.stream)
		publications = codexRejectionRoundPublications(result.stream)
		routes, branch = codexRejectionRoutes(nativeLifecycleStream(t, runner, result)), codexRejectionBranch
	}
	writeRejectionTopologyDigest(t, result.artifactDir, branch, routes)
	// Every check below is host-neutral. The gate-prepared check in particular was
	// wired Codex-only, which made FO residual mode 1 (ends without `gate prepare`)
	// invisible on Claude and Pi even though it grades durable on-disk state.
	finishLiveScenario(t, runner, scenario, result,
		durableSemantic("rejection-flow-state", assert(after, result.finalMessage)),
		durableSemantic("rejection-round-missing", assertRejectionRecordedRound(workflowRoot, entityPath, "validation", recordedRound)),
		durableSemantic("rejection-round-publication-count", assertSingleRejectionRoundPublication(publications)),
		durableSemantic("rejection-gate-not-prepared", assertRejectionGatePrepared(entityPath)),
		durableSemantic("rejection-cycle-line", assertRejectionCycleLine(entityPath)),
		durableSemantic("rejection-worker-topology", assertRejectionWorkerTopology(branch, routes)))
}

// runClaudeFeedback3CycleEscalationScenario drives the real FO against a fixture
// seeded with two prior rejection cycles at a 3rd REJECTED report and grades the
// durable end-state: the FO must escalate to the human on the 3rd cycle, not
// auto-bounce a 4th time. assertThirdCycleEscalation grades durable entity-body
// state ALONE (cycle count + escalation marker + no post-cycle-3 implementation
// report) — the reviewer-reuse signal is host-specific and lives in rejection-flow,
// not here; this scenario is purely a host-neutral durable-state grade.
func runClaudeFeedback3CycleEscalationScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario, build func(*testing.T, string) string, assert func(string) error) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := build(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, escalationPrompt(workflowRoot))
	after := readFile(t, entityPath)
	if err := assert(after); err != nil {
		t.Fatalf("%v\nEntity after:\n%s\nFinal message:\n%s\nArtifacts: %s", err, after, result.finalMessage, result.artifactDir)
	}
	runner.emitMetrics(t, scenario, result)
}

func runClaudeMergeHookGuardrailScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario, build func(*testing.T, string) string, assert func(string, string, string) error) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := build(t, workflowRoot)
	before := readFile(t, entityPath)

	result := runner.run(t, scenario, workflowRoot, mergeHookGuardPrompt(workflowRoot))
	after := readFile(t, entityPath)
	if err := assert(before, after, result.finalMessage+"\n"+result.stream); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if _, err := os.Stat(filepath.Join(workflowRoot, "_archive", "merge-check.md")); !os.IsNotExist(err) {
		t.Fatalf("merge-check was archived despite the guardrail scenario; stat err=%v", err)
	}
	runner.emitMetrics(t, scenario, result)
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
func runClaudeSelfEvidenceMergeTriageScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario, build func(*testing.T, string) string, assert func(string, string) error) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := build(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, mergeTriagePrompt(workflowRoot))
	after := readMergeTriageAfter(t, workflowRoot, entityPath)
	if err := assert(after, result.finalMessage); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	runner.emitMetrics(t, scenario, result)
}

// runClaudeSmallestSufficientMechanismScenario drives the real FO against the
// smallest-sufficient-mechanism fixture (a commissioned workflow with two ready
// entities plus two plain deterministic-edit notes whose content the prompt hands the
// FO) and grades the FO's tool-call STREAM in both directions of the ladder: the
// deterministic edits are FO-authored with a direct commit and NO worker/PR climb, and
// the commissioned ready entities are engaged via the standing dispatch loop WITHOUT a
// per-entity justification. The trace is graded, not the durable end-state, which is
// identical whether the FO climbed or not.
func runClaudeSmallestSufficientMechanismScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario, build func(*testing.T, string) string, assert func(*testing.T, string, mechanismTrace, []string, []string) error) {
	t.Helper()
	workflowRoot := build(t, t.TempDir())

	result := runner.run(t, scenario, workflowRoot, smallestMechanismPrompt(workflowRoot))
	trace := runner.smallestMechanismTrace(result, ssmEditFiles(), ssmCommissioned())
	finishLiveScenario(t, runner, scenario, result,
		durableSemantic("smallest-mechanism-violation", assert(t, workflowRoot, trace, ssmEditFiles(), ssmCommissioned())))
}

// runClaudeKeepMovingScenario grades each completed task from its own ordered,
// path-scoped Git history and keeps the questioned task active after a durable re-shape.
func runClaudeKeepMovingScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario, build func(*testing.T, string) string, assert func(*testing.T, string) error) {
	t.Helper()
	workflowRoot := build(t, t.TempDir())

	result := runner.run(t, scenario, workflowRoot, keepMovingPrompt(workflowRoot))
	finishLiveScenario(t, runner, scenario, result,
		durableSemantic("keep-moving-violation", assert(t, workflowRoot)))
}

// runClaudeFilingScenario drives the real FO against an EMPTY workflow and asks it
// to file one seed entity. It grades the FO's recorded tool-call stream — the FO
// filed via `spacedock … new <slug>`, not the `--next-id` + `Write` pair — because
// the durable end-state file is indistinguishable between the two paths. The file
// must also actually land (the run produced a real seed), so the stream grade is
// proof of HOW, not just THAT, the entity was filed.
func runClaudeFilingScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario, build func(*testing.T, string) string, assert func([]string, string) error) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := build(t, workflowRoot)
	result := runner.run(t, scenario, workflowRoot, filingPrompt(workflowRoot))
	if _, err := os.Stat(entityPath); err != nil {
		t.Fatalf("the FO did not land the seed entity at %s: %v\nFinal message:\n%s\nArtifacts: %s", entityPath, err, result.finalMessage, result.artifactDir)
	}
	if err := assert(result.commands, filingSlug); err != nil {
		finishLiveScenario(t, runner, scenario, result, durableSemantic("filing-command-not-observed", err))
		return
	}
	finishLiveScenario(t, runner, scenario, result)
}

func assertFilingCommands(commands []string, slug string) error {
	filed := false
	for _, command := range commands {
		if nextIDInvocation.MatchString(command) {
			return fmt.Errorf("filing previewed --next-id instead of using the atomic new path")
		}
		if commandFilesViaNew(command, slug) {
			filed = true
		}
	}
	if !filed {
		return fmt.Errorf("filing command log has no spacedock new %s invocation", slug)
	}
	return nil
}

// runClaudeShallowBootScenario drives the real FO against the shallow-boot fixture
// (one gate-check entity at a human gate) with a per-run isolated team root, and
// grades the durable end-state: the FO greets with the accurate held-gate state,
// NO entity mutation occurs, NO team config lands on disk, and no durable dispatch
// fingerprint appears. It then asserts the AC-2 behavioral signal (no TeamCreate
// before the greet) and records the AC-6 measured signal over the captured stream.
func runClaudeShallowBootScenario(t *testing.T, runner liveDriver, scenario sharedRuntimeScenario, build func(*testing.T, string) shallowBootFixture, assert func(shallowBootObservation) error) {
	t.Helper()
	workflowRoot := t.TempDir()
	fixture := build(t, workflowRoot)
	gateBefore := readFile(t, fixture.gateEntityPath)

	result := runner.run(t, scenario, workflowRoot, shallowBootPrompt(workflowRoot))

	// The Claude team root is {home}/.claude/teams — the exact path the comm-officer
	// startup hook membership-checks and TeamCreate writes a team config.json under.
	teamRoot := filepath.Join(runner.home(), ".claude", "teams")
	obs := gatherShallowBootObservation(t, workflowRoot, teamRoot, fixture, gateBefore, result.finalMessage)
	if err := assert(obs); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	runner.gradeShallowBootObservation(t, result)
	runner.emitMetrics(t, scenario, result)
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
// FOStreamWatcher, shared with TestLiveCommonFullEnsignCycle) — one mechanism, no second
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

	cmd := exec.Command(r.binary, "claude",
		"--plugin-dir", r.pluginDir,
		"--skip-compat-check",
		"--",
		"-p", prompt+" "+antiShutdownOverride,
		"--permission-mode", "bypassPermissions",
		"--output-format", "stream-json",
		"--verbose",
		"--model", r.modelName,
	)
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
		commands:     claudeObservedCommands(stream),
		artifactDir:  artifactDir,
		duration:     duration,
		configDir:    configDir,
		cwd:          resolvedCwd,
	}
}

func claudeObservedCommands(stream string) []string {
	var commands []string
	for _, line := range strings.Split(stream, "\n") {
		var entry streamEntry
		if json.Unmarshal([]byte(line), &entry) != nil {
			continue
		}
		for _, block := range entry.toolUseBlocks() {
			if block.Name == "Bash" {
				commands = append(commands, block.Input.Command)
			}
		}
	}
	return commands
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
