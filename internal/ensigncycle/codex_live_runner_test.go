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

// The Codex runner adapter: it turns a host-neutral sharedRuntimeScenario into a
// real `codex exec --json` launch and returns the (before, after, observed) state
// the shared assertions consume. Auth/HOME isolation (isolated CODEX_HOME +
// copied auth.json / OPENAI_API_KEY), the local Codex marketplace+plugin install,
// and the `--output-last-message` observed-extract are the ONLY Codex-specific
// surface; the scenario table, fixtures, prompts, and assertions are shared with
// the Claude runner.

type codexLiveRunner struct {
	codexBin     string
	env          []string
	artifactRoot string
}

type codexScenarioResult struct {
	finalMessage string
	jsonl        string
	artifactDir  string
	duration     time.Duration
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
		"gate-guardrail":              runCodexGateGuardrailScenario,
		"rejection-flow":              runCodexRejectionFlowScenario,
		"feedback-3-cycle-escalation": runCodexFeedback3CycleEscalationScenario,
		"merge-hook-guardrail":        runCodexMergeHookGuardrailScenario,
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
	codexHome := t.TempDir()
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
	install, err := writeCodexLocalMarketplace(t.TempDir(), repo)
	if err != nil {
		t.Fatalf("write local Codex marketplace: %v", err)
	}
	switch decision.mode {
	case codexAuthAPIKey:
		runCodexLiveCommand(t, setupDir, "codex-login.txt", openAIAPIKey+"\n", env, codexBin, "login", "--with-api-key")
	case codexAuthLocal:
		runCodexLiveCommand(t, setupDir, "codex-login-status.txt", "", env, codexBin, "login", "status")
	}
	runCodexLiveCommand(t, setupDir, "codex-marketplace-add.txt", "", env, codexBin, "plugin", "marketplace", "add", install.marketplaceRoot)
	runCodexLiveCommand(t, setupDir, "codex-plugin-add.txt", "", env, codexBin, "plugin", "add", "spacedock@spacedock")
	listing := runCodexLiveCommand(t, setupDir, "codex-plugin-list.txt", "", env, codexBin, "plugin", "list")
	if !strings.Contains(listing, install.pluginPath) {
		t.Fatalf("codex plugin list did not point at the local checkout path %q:\n%s", install.pluginPath, listing)
	}
	if strings.Contains(listing, "github.com") || strings.Contains(listing, "ref `next`") {
		t.Fatalf("codex plugin list points at remote next, not the local checkout:\n%s", listing)
	}

	adapterPath := filepath.Join(install.pluginPath, "skills", "first-officer", "references", "codex-first-officer-runtime.md")
	if _, err := os.Stat(adapterPath); err != nil {
		t.Fatalf("current-checkout plugin cache is missing r0 Codex adapter %s: %v", adapterPath, err)
	}
	if err := os.WriteFile(filepath.Join(setupDir, "codex-runtime-adapter-present.txt"), []byte(adapterPath+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	return codexLiveRunner{codexBin: codexBin, env: env, artifactRoot: artifactRoot}
}

func runCodexGateGuardrailScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
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
	emitCodexScenarioMetrics(t, scenario, result)
}

func runCodexRejectionFlowScenario(t *testing.T, runner codexLiveRunner, scenario sharedRuntimeScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeRejectionWorkflow(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, rejectionPrompt())
	after := readFile(t, entityPath)
	if err := assertRejectionFlow(after, result.finalMessage+"\n"+result.jsonl); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	// AC-4 reviewer-reuse: on Codex the FO routes the cycle-2 re-review to the
	// kept-alive validation worker with a `send_input` call (the Codex analog of
	// Claude's SendMessage reuse), not a fresh dispatch. Host-specific producer
	// signal, graded by the runner — not the shared host-neutral assertion.
	if err := assertCodexReviewerReuse(result.jsonl); err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
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

	result := runner.run(t, scenario, workflowRoot, escalationPrompt())
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

	result := runner.run(t, scenario, workflowRoot, mergeHookGuardPrompt())
	after := readFile(t, entityPath)
	if err := assertMergeHookGuardHeld(before, after, result.finalMessage+"\n"+result.jsonl); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if _, err := os.Stat(filepath.Join(workflowRoot, "_archive", "merge-check.md")); !os.IsNotExist(err) {
		t.Fatalf("merge-check was archived despite the guardrail scenario; stat err=%v", err)
	}
	emitCodexScenarioMetrics(t, scenario, result)
}

// run launches `codex exec --json` for one shared scenario. Liveness is the SAME
// streamWatcher the Claude runner and the live cycle use — one mechanism, no second
// impl. drainToExit runs the process to exit accumulating the full --json
// transcript, bounded by the per-step no-progress quietBudgetDefault (60s): the
// deadline resets on every event line, so a genuine multi-minute run never trips as
// long as Codex keeps emitting item.* events, and only silence past the budget
// kills the process — the same ≤60s AC-1-guarded discipline.
func (r codexLiveRunner) run(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, prompt string) codexScenarioResult {
	t.Helper()
	artifactDir := filepath.Join(r.artifactRoot, scenario.name)
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	finalPath := filepath.Join(artifactDir, "codex-final-message.txt")
	jsonlPath := filepath.Join(artifactDir, "codex-exec.jsonl")
	stderrPath := filepath.Join(artifactDir, "codex-exec.stderr.txt")

	cmd := exec.Command(r.codexBin,
		"exec",
		"--json",
		"--dangerously-bypass-approvals-and-sandbox",
		"--cd", workflowRoot,
		"--output-last-message", finalPath,
		prompt,
	)
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
		t.Fatalf("codex exec failed to start for %s: %v", scenario.name, startErr)
	}
	poller := newCmdPoller(cmd, pw)
	defer poller.kill()
	watcher := newStreamWatcher(newPipeLineSource(pr), poller, func(line string) { t.Log(line) })

	// drainToExit runs the process to exit accumulating the full transcript, OR
	// kills it on a 60s no-progress stall; the deferred poller.kill() reaps it.
	jsonl, stallErr := watcher.drainToExit(quietBudgetDefault, "codex shared scenario "+scenario.name)
	duration := time.Since(started)

	if writeErr := os.WriteFile(jsonlPath, []byte(jsonl), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if stallErr != nil {
		t.Fatalf("%v\nArtifacts: %s", stallErr, artifactDir)
	}

	return codexScenarioResult{
		finalMessage: readFile(t, finalPath),
		jsonl:        jsonl,
		artifactDir:  artifactDir,
		duration:     duration,
	}
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
