//go:build live

package ensigncycle

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/cli"
)

// The Codex runner adapter: it turns a host-neutral sharedRuntimeScenario into a
// real `codex exec --json` launch and returns the (before, after, observed) state
// the shared assertions consume. Auth/HOME isolation (isolated CODEX_HOME +
// copied auth.json / OPENAI_API_KEY), the local Codex marketplace+plugin install,
// and the `--output-last-message` observed-extract are the ONLY Codex-specific
// surface; the scenario table, fixtures, prompts, and assertions are shared with
// the Claude runner.

type codexLiveRunner struct {
	codexBin         string
	env              []string
	artifactRoot     string
	firstOfficerBase string
}

type codexLiveScenario struct {
	sharedRuntimeScenario
	run func(*testing.T, codexLiveRunner, sharedRuntimeScenario)
}

func TestLiveCodexSharedScenarios(t *testing.T) {
	runner := newCodexLiveRunner(t)

	for _, scenario := range codexLiveScenarios(t) {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.run(t, runner, scenario.sharedRuntimeScenario)
		})
	}
}

func codexLiveScenarios(t *testing.T) []codexLiveScenario {
	t.Helper()
	runners := codexScenarioRunners()

	var scenarios []codexLiveScenario
	for _, scenario := range sharedRuntimeScenarios() {
		run := runners[scenario.name]
		if run == nil {
			t.Fatalf("shared scenario %q has no Codex live runner", scenario.name)
		}
		scenarios = append(scenarios, codexLiveScenario{
			sharedRuntimeScenario: scenario,
			run:                   run,
		})
	}
	return scenarios
}

// codexScenarioRunners maps each shared scenario ID to its Codex runner. It is the
// Codex side of the parity guard: the shared coverage meta-test fails if this map
// lacks a runner for any sharedRuntimeScenarios() ID.
func codexScenarioRunners() map[string]func(*testing.T, codexLiveRunner, sharedRuntimeScenario) {
	return map[string]func(*testing.T, codexLiveRunner, sharedRuntimeScenario){
		"gate-guardrail":                runCodexGateGuardrailScenario,
		"rejection-flow":                runCodexRejectionFlowScenario,
		"feedback-3-cycle-escalation":   runCodexFeedback3CycleEscalationScenario,
		"merge-hook-guardrail":          runCodexMergeHookGuardrailScenario,
		"filing":                        runCodexFilingScenario,
		"shallow-boot":                  runCodexShallowBootScenario,
		"self-evidence-merge-triage":    runCodexSelfEvidenceMergeTriageScenario,
		"smallest-sufficient-mechanism": runCodexSmallestSufficientMechanismScenario,
		"keep-moving-posture":           runCodexKeepMovingScenario,
	}
}

func newCodexLiveRunner(t *testing.T) codexLiveRunner {
	t.Helper()
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	realHome := os.Getenv("HOME")
	decision := decideCodexLiveAuth(openAIAPIKey, codexLocalAuthAvailable(realHome), os.Getenv("SPACEDOCK_CODEX_LIVE_REQUIRED"))
	switch decision.mode {
	case codexAuthSkip:
		t.Skip(decision.message)
	case codexAuthFatal:
		t.Fatal(decision.message)
	}
	codexBin, err := exec.LookPath("codex")
	if err != nil {
		t.Fatal("codex not on PATH; install Codex CLI before running the live Codex suite")
	}

	binary := spacedockBinary(t)
	repo := repoRoot(t)
	artifactRoot := codexLiveArtifactDir(t, "codex-shared-scenarios")
	codexHome := newCodexLiveIsolatedHome(t, repo, artifactRoot)
	cleanHome := t.TempDir()
	if decision.mode == codexAuthLocal {
		if err := seedCodexLocalAuth(codexHome, realHome); err != nil {
			t.Fatalf("seed local Codex auth: %v", err)
		}
	}
	env := codexLiveEnv(codexHome, cleanHome, filepath.Dir(binary), openAIAPIKey)

	setupDir := filepath.Join(artifactRoot, "_setup")
	if err := os.MkdirAll(setupDir, 0o755); err != nil {
		t.Fatal(err)
	}
	install, err := cli.WriteCodexLocalMarketplace(t.TempDir(), repo, "spacedock")
	if err != nil {
		t.Fatalf("write local Codex marketplace: %v", err)
	}
	switch decision.mode {
	case codexAuthAPIKey:
		runCodexLiveCommand(t, setupDir, "codex-login.txt", openAIAPIKey+"\n", env, codexBin, "login", "--with-api-key")
	case codexAuthLocal:
		runCodexLiveCommand(t, setupDir, "codex-login-status.txt", "", env, codexBin, "login", "status")
	}
	runCodexLiveCommand(t, setupDir, "codex-marketplace-add.txt", "", env, codexBin, "plugin", "marketplace", "add", install.MarketplaceRoot)
	runCodexLiveCommand(t, setupDir, "codex-plugin-add.txt", "", env, codexBin, "plugin", "add", "spacedock@spacedock")
	listing := runCodexLiveCommand(t, setupDir, "codex-plugin-list.txt", "", env, codexBin, "plugin", "list")
	if !strings.Contains(listing, install.PluginPath) {
		t.Fatalf("codex plugin list did not point at the local checkout path %q:\n%s", install.PluginPath, listing)
	}
	if strings.Contains(listing, "github.com") || strings.Contains(listing, "ref `next`") {
		t.Fatalf("codex plugin list points at remote next, not the local checkout:\n%s", listing)
	}

	adapterPath := filepath.Join(install.PluginPath, "skills", "first-officer", "references", "codex-first-officer-runtime.md")
	if _, err := os.Stat(adapterPath); err != nil {
		t.Fatalf("current-checkout plugin cache is missing r0 Codex adapter %s: %v", adapterPath, err)
	}
	if err := os.WriteFile(filepath.Join(setupDir, "codex-runtime-adapter-present.txt"), []byte(adapterPath+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	firstOfficerBase := codexInstalledFirstOfficerBase(t, codexHome)

	return codexLiveRunner{
		codexBin:         codexBin,
		env:              env,
		artifactRoot:     artifactRoot,
		firstOfficerBase: firstOfficerBase,
	}
}

// codexInstalledFirstOfficerBase resolves the one installed cache entry in this
// runner's isolated CODEX_HOME. This is the loader-visible base emitted in Codex
// command events; the marketplace source symlink is not the installed path.
func codexInstalledFirstOfficerBase(t *testing.T, codexHome string) string {
	t.Helper()
	cacheRoot := filepath.Join(codexHome, "plugins", "cache", "spacedock", "spacedock")
	entries, err := os.ReadDir(cacheRoot)
	if err != nil {
		t.Fatalf("read isolated Codex plugin cache %s: %v", cacheRoot, err)
	}
	var candidates []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		base := filepath.Join(cacheRoot, entry.Name(), "skills", "first-officer")
		if _, err := os.Stat(filepath.Join(base, "SKILL.md")); err == nil {
			candidates = append(candidates, base)
		}
	}
	if len(candidates) != 1 {
		t.Fatalf("isolated Codex plugin cache has %d first-officer installs under %s, want exactly one: %v", len(candidates), cacheRoot, candidates)
	}
	return candidates[0]
}

func newCodexLiveIsolatedHome(t *testing.T, repo, artifactRoot string) string {
	t.Helper()
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Fatalf("resolve user cache dir for isolated CODEX_HOME: %v", err)
	}
	var failures []string
	for _, parent := range codexLiveIsolatedHomeParentCandidates(cacheDir, repo, artifactRoot) {
		if err := os.MkdirAll(parent, 0o700); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", parent, err))
			continue
		}
		dir, err := os.MkdirTemp(parent, "codex-home-")
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", parent, err))
			continue
		}
		t.Cleanup(func() {
			_ = os.RemoveAll(dir)
		})
		return dir
	}
	t.Fatalf("create isolated CODEX_HOME outside system temp: %s", strings.Join(failures, "; "))
	return ""
}

func runCodexGateGuardrailScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeGateWorkflow(t, workflowRoot)
	before := readFile(t, entityPath)

	result, err := runner.run(t, scenario, workflowRoot, gatePrompt(workflowRoot), 0)
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	after := readFile(t, entityPath)
	if err := assertGateHeld(before, after, result.finalMessage); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if err := assertFOGateLoadBoundary(codexFOLoadTrace(result.jsonl, mustFOLoadSpec(t, runner.firstOfficerBase, "codex"), nil)); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	if _, err := os.Stat(filepath.Join(workflowRoot, "_archive", "gate-check.md")); !os.IsNotExist(err) {
		t.Fatalf("gate-check was archived while waiting at the gate; stat err=%v", err)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

func runCodexRejectionFlowScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	attempt, err := runCodexRejectionFlowWithRetry(func(attempt int) (codexRejectionFlowAttempt, error) {
		workflowRoot := t.TempDir()
		entityPath := writeRejectionWorkflow(t, workflowRoot)
		result, err := runner.run(t, scenario, workflowRoot, rejectionPrompt(workflowRoot), attempt)
		if err != nil {
			return codexRejectionFlowAttempt{result: result}, err
		}
		return codexRejectionFlowAttempt{
			entityAfter: readFile(t, entityPath),
			result:      result,
		}, nil
	})
	if err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, attempt.result.finalMessage, attempt.result.artifactDir)
	}
	if err := assertFOWriteBeforeFirstMutation(codexFOLoadTrace(attempt.result.jsonl, mustFOLoadSpec(t, runner.firstOfficerBase, "codex"), nil)); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, attempt.result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, attempt.result)
}

// runCodexFeedback3CycleEscalationScenario drives the real FO against a fixture
// seeded with two prior rejection cycles at a 3rd REJECTED report and grades the
// durable end-state: the FO must escalate to the human on the 3rd cycle, not
// auto-bounce a 4th time. assertThirdCycleEscalation grades durable entity-body
// state ALONE (cycle count + escalation marker + no post-cycle-3 implementation
// report) — host-neutral, the same assertion the Claude runner feeds.
func runCodexFeedback3CycleEscalationScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeEscalationWorkflow(t, workflowRoot)

	result, err := runner.run(t, scenario, workflowRoot, escalationPrompt(workflowRoot), 0)
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	after := readFile(t, entityPath)
	if err := assertThirdCycleEscalation(after); err != nil {
		t.Fatalf("%v\nEntity after:\n%s\nFinal message:\n%s\nArtifacts: %s", err, after, result.finalMessage, result.artifactDir)
	}
	if err := assertFOWriteBeforeFirstMutation(codexFOLoadTrace(result.jsonl, mustFOLoadSpec(t, runner.firstOfficerBase, "codex"), nil)); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

func runCodexMergeHookGuardrailScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeMergeHookGuardWorkflow(t, workflowRoot)
	before := readFile(t, entityPath)

	result, err := runner.run(t, scenario, workflowRoot, mergeHookGuardPrompt(workflowRoot), 0)
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	after := readFile(t, entityPath)
	if err := assertMergeHookGuardHeld(before, after, result.finalMessage+"\n"+result.jsonl); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if err := assertFOTerminalLoadBoundary(codexFOLoadTrace(result.jsonl, mustFOLoadSpec(t, runner.firstOfficerBase, "codex"), codexTerminalAction("merge-check"))); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	if _, err := os.Stat(filepath.Join(workflowRoot, "_archive", "merge-check.md")); !os.IsNotExist(err) {
		t.Fatalf("merge-check was archived despite the guardrail scenario; stat err=%v", err)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

// runCodexSelfEvidenceMergeTriageScenario drives the real FO against the merge/triage
// fixture and grades the SAME host-neutral decision assertSelfEvidenceMergeTriage the
// Claude runner feeds, against the INCIDENT class only (Cycle-3): the FO must NOT wave
// the change through — terminalize-as-PASSED or merge-without-evidence — while the
// required lane is unapproved; a hold or a reject-without-merge both pass. Its
// final-message diagnosis must name THIS run's failing test, not the inherited "known
// flake" label. The this-run token is graded in the FINAL MESSAGE only, so the
// entity-read echo in the transcript cannot satisfy it.
func runCodexSelfEvidenceMergeTriageScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeMergeTriageWorkflow(t, workflowRoot)

	result, err := runner.run(t, scenario, workflowRoot, mergeTriagePrompt(workflowRoot), 0)
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	after := readMergeTriageAfter(t, workflowRoot, entityPath)
	if err := assertSelfEvidenceMergeTriage(after, result.finalMessage); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if err := assertFOWriteBeforeObservedMutation(codexFOLoadTrace(result.jsonl, mustFOLoadSpec(t, runner.firstOfficerBase, "codex"), nil)); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

// runCodexSmallestSufficientMechanismScenario drives the real FO against the
// smallest-sufficient-mechanism fixture and grades the SAME host-neutral ladder
// assertCodexSmallestSufficientMechanism the Claude runner feeds, over the Codex
// command/collab transcript: the deterministic edits are FO-authored (in-house
// apply_patch) with a direct commit and NO worker/PR climb, and the commissioned ready
// entities are engaged via the standing dispatch loop without a per-entity justification.
func runCodexSmallestSufficientMechanismScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	writeSmallestMechanismWorkflow(t, workflowRoot)

	result, err := runner.run(t, scenario, workflowRoot, smallestMechanismPrompt(workflowRoot), 0)
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	if err := assertCodexSmallestSufficientMechanism(result.jsonl, ssmEditFiles(), ssmCommissioned()); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if err := assertFOWriteBeforeFirstMutation(codexFOLoadTrace(result.jsonl, mustFOLoadSpec(t, runner.firstOfficerBase, "codex"), nil)); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

// runCodexKeepMovingScenario drives the real FO against the keep-moving-posture fixture
// and grades the SAME host-neutral patterns assertCodexKeepMoving the Claude runner feeds,
// over the Codex command/collab/file_change transcript plus the final message: advance +
// dispatch the approved entity with no permission question, dispatch both independent
// entities, re-shape the questioned entity and pause its dispatch, and no turn-end on an
// async wait.
func runCodexKeepMovingScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	writeKeepMovingWorkflow(t, workflowRoot)

	result, err := runner.run(t, scenario, workflowRoot, keepMovingPrompt(workflowRoot), 0)
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	if err := assertCodexKeepMoving(result.jsonl, result.finalMessage, kmIndependent()); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if err := assertFOWriteBeforeFirstMutation(codexFOLoadTrace(result.jsonl, mustFOLoadSpec(t, runner.firstOfficerBase, "codex"), nil)); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

// runCodexFilingScenario drives the real FO against an EMPTY workflow and asks it
// to file one seed entity. Like the Claude runner it grades the FO's recorded
// command stream — the FO filed via `spacedock … new <slug>`, not a `--next-id`
// preview-then-write — because the durable end-state file is indistinguishable
// between the two paths. The file must also actually land, so the stream grade is
// proof of HOW, not just THAT, the entity was filed.
func runCodexFilingScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeFilingWorkflow(t, workflowRoot)

	result, err := runner.run(t, scenario, workflowRoot, filingPrompt(workflowRoot), 0)
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	if _, err := os.Stat(entityPath); err != nil {
		t.Fatalf("the FO did not land the seed entity at %s: %v\nFinal message:\n%s\nArtifacts: %s", entityPath, err, result.finalMessage, result.artifactDir)
	}
	if err := assertCodexFilingViaNew(result.jsonl, filingSlug); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if err := assertFOFilingLoadBoundary(codexFOLoadTrace(result.jsonl, mustFOLoadSpec(t, runner.firstOfficerBase, "codex"), codexFilingAction(filingSlug))); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

// runCodexShallowBootScenario drives the real FO against the shallow-boot fixture
// and grades the SAME host-neutral durable end-state assertShallowBoot the Claude
// runner feeds: the FO greets and presents the gate, S7b advances+archives the
// merged PR before-greet, and NO worker is dispatched (the gate entity is
// unchanged, not archived, no worktree). Codex has no Claude team root, so the
// no-team-config check is host-neutral-vacuous (empty teamRoot); the no-dispatch
// proof rides the durable gate-unchanged + no-worktree facts. The AC-2/AC-6 Claude
// token-stream measurements are Claude-specific and live in the Claude runner.
func runCodexShallowBootScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	fixture := writeShallowBootWorkflow(t, workflowRoot)
	gateBefore := readFile(t, fixture.gateEntityPath)

	// The stub `gh` (reporting MERGED) must resolve on the FO subprocess PATH so the
	// boot's live pr_state probe and the pr-merge startup hook both see the merge.
	scenarioRunner := runner
	scenarioRunner.env = withPATHPrefix(runner.env, fixture.stubGhDir)

	result, err := scenarioRunner.run(t, scenario, workflowRoot, shallowBootPrompt(workflowRoot), 0)
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}

	obs := gatherShallowBootObservation(t, workflowRoot, "", fixture, gateBefore, result.finalMessage)
	if err := assertShallowBoot(obs); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if err := assertFOWriteBeforeFirstMutation(codexFOLoadTrace(result.jsonl, mustFOLoadSpec(t, runner.firstOfficerBase, "codex"), nil)); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

func codexExecArgv(workflowRoot, finalPath, prompt string) []string {
	return []string{
		"exec",
		"--json",
		"--enable", "multi_agent_v2",
		"--dangerously-bypass-approvals-and-sandbox",
		"--cd", workflowRoot,
		"--output-last-message", finalPath,
		prompt,
	}
}

// run launches `codex exec --json` for one shared scenario. Liveness still uses
// the shared streamWatcher for the 60s stream-silence guard, with a Codex-specific
// foreground-wait watchdog layered beside it so repeated wait-loop JSONL does not
// count as progress forever. The watchdog returns a typed stall; it does not make
// a scenario pass, because the pass assertions still grade durable state plus the
// Codex producer signal.
func (r codexLiveRunner) run(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, prompt string, attempt int) (codexScenarioResult, error) {
	t.Helper()
	artifactDir := codexAttemptArtifactDir(r.artifactRoot, scenario.name, attempt)
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	finalPath := filepath.Join(artifactDir, "codex-final-message.txt")
	jsonlPath := filepath.Join(artifactDir, "codex-exec.jsonl")
	stderrPath := filepath.Join(artifactDir, "codex-exec.stderr.txt")

	cmd := exec.Command(r.codexBin, codexExecArgv(workflowRoot, finalPath, prompt)...)
	cmd.Env = r.env
	// stdout (the --json event stream) flows through the watcher's pipe for the
	// no-progress liveness guard; stderr goes to its own artifact file. The
	// cmdPoller closes the pipe write-end on exit so the scanner EOFs.
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	stderr, err := os.Create(stderrPath)
	if err != nil {
		t.Fatal(err)
	}
	defer stderr.Close()
	cmd.Stderr = stderr

	started := time.Now()
	if startErr := cmd.Start(); startErr != nil {
		return codexScenarioResult{artifactDir: artifactDir}, fmt.Errorf("codex exec failed to start for %s: %w", scenario.name, startErr)
	}
	poller := newCmdPoller(cmd, pw)
	defer poller.kill()
	watcher := newStreamWatcher(newPipeLineSource(pr), poller, discardStreamLine)
	probe, err := newWorkflowStateProbe(workflowRoot)
	if err != nil {
		return codexScenarioResult{artifactDir: artifactDir}, fmt.Errorf("snapshot workflow state for %s: %w", scenario.name, err)
	}
	watchdog := newCodexCollabWaitWatchdog(scenario.name, artifactDir, probe)

	// drainToExit runs the process to exit accumulating the full transcript, OR
	// kills it on a 60s no-progress stall; the deferred poller.kill() reaps it.
	jsonl, stallErr := drainCodexToExitWithWaitWatchdog(watcher, quietBudgetDefault, "codex shared scenario "+scenario.name, watchdog)
	duration := time.Since(started)

	if writeErr := os.WriteFile(jsonlPath, []byte(jsonl), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if stallErr != nil {
		return codexScenarioResult{
			jsonl:       jsonl,
			artifactDir: artifactDir,
			duration:    duration,
		}, stallErr
	}

	return codexScenarioResult{
		finalMessage: readFile(t, finalPath),
		jsonl:        jsonl,
		artifactDir:  artifactDir,
		duration:     duration,
	}, nil
}

func codexLiveArtifactDir(t *testing.T, name string) string {
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

func runCodexLiveCommand(t *testing.T, artifactDir, artifactName, stdin string, env []string, argv ...string) string {
	t.Helper()
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Env = env
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	out, err := cmd.CombinedOutput()
	if writeErr := os.WriteFile(filepath.Join(artifactDir, artifactName), out, 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if err != nil {
		t.Fatalf("%s failed: %v\n%s", strings.Join(argv, " "), err, out)
	}
	return string(out)
}

func argvHasAdjacent(args []string, left, right string) bool {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == left && args[i+1] == right {
			return true
		}
	}
	return false
}

func TestCodexLiveRunnerExecArgvEnablesMultiAgentV2(t *testing.T) {
	args := codexExecArgv("/tmp/workflow", "/tmp/final-message.txt", "run the scenario")
	if !argvHasAdjacent(args, "--enable", "multi_agent_v2") {
		t.Fatalf("codex live exec argv must explicitly enable multi_agent_v2 because CODEX_HOME is isolated; args=%v", args)
	}
}
