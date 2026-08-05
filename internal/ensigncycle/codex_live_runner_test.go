//go:build live

package ensigncycle

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The Codex runner adapter: it turns a host-neutral sharedRuntimeScenario into a
// real `spacedock codex` launch and returns the (before, after, observed) state
// the shared assertions consume. Auth/HOME isolation (isolated CODEX_HOME +
// minimal config plus copied auth.json / OPENAI_API_KEY), Spacedock-owned local
// plugin setup, and the `--output-last-message` observed-extract are the ONLY
// Codex-specific surface; the scenario table, fixtures, prompts, and assertions
// are shared with the Claude runner.

type codexLiveRunner struct {
	binary       string
	pluginDir    string
	codexBin     string
	env          []string
	artifactRoot string
}

type codexLiveScenario struct {
	sharedRuntimeScenario
	run func(*testing.T, codexLiveRunner, sharedRuntimeScenario)
}

func TestLiveCodexSharedScenarios(t *testing.T) {
	runner := newCodexLiveRunner(t)

	for _, scenario := range codexLiveScenarios(t) {
		t.Run(scenario.name, func(t *testing.T) {
			if reason := liveDurableJourneyTODO(scenario.name); reason != "" {
				t.Skip(reason)
			}
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
		"recorded-gate-lifecycle":       runCodexRecordedGateLifecycleScenario,
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

func (r codexLiveRunner) withStubPATH(t *testing.T, dir string) codexLiveRunner {
	t.Helper()
	r.env = withPATHPrefix(r.env, dir)
	r.env = withRecordedGateEnv(r.env, "SPACEDOCK_BIN", filepath.Join(dir, "spacedock"))
	// The Spacedock front door re-pins SPACEDOCK_BIN to its own executable before
	// launching Codex. Put a Codex shim in front of the host only for recorded-gate
	// fixtures so the child FO still resolves the fixture's command logger.
	shim := filepath.Join(dir, "codex")
	script := "#!/bin/sh\n" +
		"export SPACEDOCK_BIN=" + shellQuote(filepath.Join(dir, "spacedock")) + "\n" +
		"exec " + shellQuote(r.codexBin) + " \"$@\"\n"
	if err := os.WriteFile(shim, []byte(script), 0o755); err != nil {
		t.Fatalf("write recorded-gate Codex shim: %v", err)
	}
	return r
}

func runCodexRecordedGateLifecycleScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	fixture := writePreparedRecordedGateFixture(t)
	before := readFile(t, fixture.entity)
	commandLog := filepath.Join(fixture.root, "evidence", "command.log")
	shimDir := writeRecordedGateLoggingShim(t, buildRecordedGateBinary(t), commandLog)
	result, err := runner.withStubPATH(t, shimDir).run(t, scenario, fixture.root, recordedGatePrompt(fixture.root))
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	observation := recordedGateLiveObservation(t, fixture, before, commandLog)
	if err := assertRecordedGateLifecycle(observation); err != nil {
		t.Fatalf("recorded gate lifecycle graded FAIL: %v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, result)
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
	if err := seedCodexLiveConfig(codexHome); err != nil {
		t.Fatalf("seed live Codex config: %v", err)
	}
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
	switch decision.mode {
	case codexAuthAPIKey:
		runCodexLiveCommand(t, setupDir, "codex-login.txt", openAIAPIKey+"\n", env, codexBin, "login", "--with-api-key")
	case codexAuthLocal:
		runCodexLiveCommand(t, setupDir, "codex-login-status.txt", "", env, codexBin, "login", "status")
	}

	adapterPath := filepath.Join(repo, "skills", "first-officer", "references", "codex-first-officer-runtime.md")
	if _, err := os.Stat(adapterPath); err != nil {
		t.Fatalf("current-checkout Codex adapter is missing %s: %v", adapterPath, err)
	}
	if err := os.WriteFile(filepath.Join(setupDir, "codex-runtime-adapter-present.txt"), []byte(adapterPath+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	sourceHead := runCodexLiveCommand(t, setupDir, "source-head.txt", "", os.Environ(), "git", "-C", repo, "rev-parse", "HEAD")
	if strings.TrimSpace(sourceHead) == "" {
		t.Fatal("current-checkout source HEAD is empty")
	}

	return codexLiveRunner{binary: binary, pluginDir: repo, codexBin: codexBin, env: env, artifactRoot: artifactRoot}
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
	fixture := writeGateWorkflow(t, workflowRoot)
	before := readFile(t, fixture.entity)
	commandLog := filepath.Join(fixture.root, "evidence", "command.log")
	shimDir := writeRecordedGateLoggingShim(t, buildRecordedGateBinary(t), commandLog)

	result, err := runner.withStubPATH(t, shimDir).run(t, scenario, workflowRoot, gatePrompt(workflowRoot))
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	if _, err := os.Stat(filepath.Join(fixture.stateRoot, "_archive", "recorded-gate-task", "index.md")); !os.IsNotExist(err) {
		t.Fatalf("recorded-gate-task was archived while waiting at the gate; stat err=%v", err)
	}
	after := readFile(t, fixture.entity)
	expected, err := recordedGateHeldExpectation(fixture)
	if err != nil {
		t.Fatalf("read prepared gate expectation: %v\nArtifacts: %s", err, result.artifactDir)
	}
	if err := assertGateHeld(before, after, expected); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if err := assertRecordedGateHoldLog(readFile(t, commandLog)); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

func runCodexRejectionFlowScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	t.Skip("TODO(zbcj98qfwtax61vxdzrf615e): Codex must reliably bind a distinct post-rework Briefing before re-enabling this journey")

	workflowRoot := t.TempDir()
	entityPath := writeRejectionWorkflow(t, workflowRoot)
	result, runErr := runner.run(t, scenario, workflowRoot, rejectionPrompt(workflowRoot))
	entityAfter, captureErr := captureCodexRejectionEvidence(workflowRoot, entityPath, result.artifactDir)
	if captureErr != nil {
		t.Fatalf("capture rejection-flow evidence: %v\nArtifacts: %s", captureErr, result.artifactDir)
	}
	if runErr != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", runErr, result.finalMessage, result.artifactDir)
	}
	if err := assertRejectionFlow(entityAfter, result.finalMessage+"\n"+result.jsonl); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if err := assertRejectionRecordedRound(workflowRoot, entityPath, "validation", codexRecordedRejectionRound(result.jsonl)); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	// assertRejectionFlow (above) proves the two-cycle re-review OCCURRED from the
	// durable entity body. assertCodexReviewerReuse grades WHO performed it from
	// structured handles: a re-review routed to a non-validation thread fails; reuse
	// and structurally-distinct fresh both pass; an absent identity handle is
	// identity-not-provable — log it, do not fail, since the durable two-cycle proof
	// already stands.
	if err := assertCodexReviewerReuse(result.jsonl); err != nil {
		if errors.Is(err, errReviewerIdentityUnsupported) {
			t.Logf("rejection-flow reviewer identity abstained (re-review OCCURRED per durable state, but WHO is not structurally provable): %s", err)
		} else {
			t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
		}
	}
	emitCodexScenarioMetrics(t, scenario, result)
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

	result, err := runner.run(t, scenario, workflowRoot, escalationPrompt(workflowRoot))
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	after := readFile(t, entityPath)
	if err := assertThirdCycleEscalation(after); err != nil {
		t.Fatalf("%v\nEntity after:\n%s\nFinal message:\n%s\nArtifacts: %s", err, after, result.finalMessage, result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

func runCodexMergeHookGuardrailScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeMergeHookGuardWorkflow(t, workflowRoot)
	before := readFile(t, entityPath)

	result, err := runner.run(t, scenario, workflowRoot, mergeHookGuardPrompt(workflowRoot))
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	after := readFile(t, entityPath)
	if err := assertMergeHookGuardHeld(before, after, result.finalMessage+"\n"+result.jsonl); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
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

	result, err := runner.run(t, scenario, workflowRoot, mergeTriagePrompt(workflowRoot))
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	after := readMergeTriageAfter(t, workflowRoot, entityPath)
	if err := assertSelfEvidenceMergeTriage(after, result.finalMessage); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

// runCodexSmallestSufficientMechanismScenario drives the real FO against the
// durable fixture while preserving transcript checks unrelated to completion.
func runCodexSmallestSufficientMechanismScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	writeSmallestMechanismWorkflow(t, workflowRoot)

	result, err := runner.run(t, scenario, workflowRoot, smallestMechanismPrompt(workflowRoot))
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	trace := codexMechanismTrace(result.jsonl, ssmEditFiles(), ssmCommissioned())
	if err := assertDurableSmallestMechanism(t, workflowRoot, trace, ssmEditFiles(), ssmCommissioned()); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

// runCodexKeepMovingScenario grades each completed task from its own ordered
// Git history and keeps the questioned task active after a durable re-shape.
func runCodexKeepMovingScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot, err := os.MkdirTemp("", "spacedock-keep-moving-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cleanupKeepMovingRoot(t, workflowRoot, t.Failed()) })
	writeKeepMovingWorkflow(t, workflowRoot)

	result, err := runner.run(t, scenario, workflowRoot, keepMovingPrompt(workflowRoot))
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	if err := assertDurableKeepMoving(t, workflowRoot); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
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

	result, err := runner.run(t, scenario, workflowRoot, filingPrompt(workflowRoot))
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	if _, err := os.Stat(entityPath); err != nil {
		t.Fatalf("the FO did not land the seed entity at %s: %v\nFinal message:\n%s\nArtifacts: %s", entityPath, err, result.finalMessage, result.artifactDir)
	}
	if err := assertCodexFilingViaNew(result.jsonl, filingSlug); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

// runCodexShallowBootScenario drives the real FO against the shallow-boot fixture
// and grades the SAME host-neutral durable end-state assertShallowBoot the Claude
// runner feeds: the FO greets with the accurate held-gate state, performs no
// persisted entity mutation, and leaves no durable dispatch fingerprint (the gate
// entity is unchanged, not archived, no worktree). Codex has no Claude team root,
// so the no-team-config check is host-neutral-vacuous (empty teamRoot). Absence of
// transient dispatch commands is outside this durable oracle. The AC-2/AC-6 Claude
// token-stream measurements are Claude-specific and live in the Claude runner.
func runCodexShallowBootScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	fixture := writeShallowBootWorkflow(t, workflowRoot)
	gateBefore := readFile(t, fixture.gateEntityPath)

	result, err := runner.run(t, scenario, workflowRoot, shallowBootPrompt(workflowRoot))
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}

	obs := gatherShallowBootObservation(t, workflowRoot, "", fixture, gateBefore, result.finalMessage)
	if err := assertShallowBoot(obs); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

// run launches one `codex exec --json` for one shared scenario. Each complete
// JSONL line resets the shared quiet budget. Stream silence kills that sole
// process; activity and durable writes never trigger another launch.
func (r codexLiveRunner) run(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, prompt string) (codexScenarioResult, error) {
	t.Helper()
	artifactDir := filepath.Join(r.artifactRoot, scenario.name)
	finalPath := filepath.Join(artifactDir, "codex-final-message.txt")
	return runCodexProcess(codexProcessSpec{
		bin:         r.binary,
		argv:        codexLiveFrontDoorArgv(r.pluginDir, workflowRoot, finalPath, prompt),
		env:         r.env,
		artifactDir: artifactDir,
		finalPath:   finalPath,
		quietBudget: quietBudgetDefault,
	})
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
