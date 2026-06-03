//go:build live

package ensigncycle

import (
	"context"
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
		"gate-guardrail":       runCodexGateGuardrailScenario,
		"rejection-flow":       runCodexRejectionFlowScenario,
		"merge-hook-guardrail": runCodexMergeHookGuardrailScenario,
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

func (r codexLiveRunner) run(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, prompt string) codexScenarioResult {
	t.Helper()
	artifactDir := filepath.Join(r.artifactRoot, scenario.name)
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	finalPath := filepath.Join(artifactDir, "codex-final-message.txt")
	jsonlPath := filepath.Join(artifactDir, "codex-exec.jsonl")
	stderrPath := filepath.Join(artifactDir, "codex-exec.stderr.txt")

	ctx, cancel := context.WithTimeout(context.Background(), scenario.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, r.codexBin,
		"exec",
		"--json",
		"--dangerously-bypass-approvals-and-sandbox",
		"--cd", workflowRoot,
		"--output-last-message", finalPath,
		prompt,
	)
	cmd.Env = r.env
	stdout, err := os.Create(jsonlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer stdout.Close()
	stderr, err := os.Create(stderrPath)
	if err != nil {
		t.Fatal(err)
	}
	defer stderr.Close()
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	started := time.Now()
	err = cmd.Run()
	duration := time.Since(started)
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("codex exec did not finish within %s for %s; artifacts in %s", scenario.timeout, scenario.name, artifactDir)
	}
	if err != nil {
		t.Fatalf("codex exec failed for %s: %v; artifacts in %s", scenario.name, err, artifactDir)
	}

	return codexScenarioResult{
		finalMessage: readFile(t, finalPath),
		jsonl:        readFile(t, jsonlPath),
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
