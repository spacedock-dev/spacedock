//go:build live

package ensigncycle

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
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
	binary       string
	repoRoot     string
	env          []string
	model        string
	artifactRoot string
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

	for _, scenario := range claudeLiveScenarios(t) {
		t.Run(scenario.name, func(t *testing.T) {
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
		"gate-guardrail":       runClaudeGateGuardrailScenario,
		"rejection-flow":       runClaudeRejectionFlowScenario,
		"merge-hook-guardrail": runClaudeMergeHookGuardrailScenario,
	}
}

func newClaudeLiveRunner(t *testing.T) claudeLiveRunner {
	t.Helper()
	binary := spacedockBinary(t)
	repo := repoRoot(t)
	model := envOr("SPACEDOCK_LIVE_MODEL", "sonnet")

	// isolatedClaudeEnv resolves the credential (OAuth benchmark-token locally,
	// ANTHROPIC_API_KEY in CI) against a fresh empty HOME, or t.Skips when neither
	// is available. withBinaryOnPath puts the built binary first on the FO
	// subprocess PATH so its `spacedock --version` contract step resolves the test
	// binary. Both are reused verbatim from the full-cycle live test.
	env := isolatedClaudeEnv(t, os.Getenv("HOME"))
	env = withBinaryOnPath(env, binary)

	return claudeLiveRunner{
		binary:       binary,
		repoRoot:     repo,
		env:          env,
		model:        model,
		artifactRoot: claudeLiveArtifactDir(t, "claude-shared-scenarios"),
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

// run launches the real `spacedock claude` front door for one shared scenario and
// returns the (finalMessage, full stream) the shared assertions consume. The
// launch shape is the spike WINNER: --plugin-dir + --skip-contract-check are the
// spacedock-owned flags BEFORE `--`; every host flag (-p with the scenario prompt,
// --permission-mode, --output-format stream-json, --verbose, --model) rides AFTER
// `--` and forwards verbatim to claude. The observed source is the stream's
// result/success event via extractClaudeFinalMessage — a 401/is_error result is a
// LOUD launch failure here, never fed into a scenario assertion.
func (r claudeLiveRunner) run(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, prompt string) claudeScenarioResult {
	t.Helper()
	artifactDir := filepath.Join(r.artifactRoot, scenario.name)
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	streamPath := filepath.Join(artifactDir, "claude-stream.jsonl")
	finalPath := filepath.Join(artifactDir, "claude-final-message.txt")

	ctx, cancel := context.WithTimeout(context.Background(), scenario.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, r.binary, "claude",
		"--plugin-dir", r.repoRoot,
		"--skip-contract-check",
		"--",
		"-p", prompt+" "+antiShutdownOverride,
		"--permission-mode", "bypassPermissions",
		"--output-format", "stream-json",
		"--verbose",
		"--model", r.model,
	)
	cmd.Dir = workflowRoot
	cmd.Env = r.env

	// stdout carries the stream-json transcript; stderr is folded in so a launch
	// error (e.g. a stale-token 401 printed to stderr) is captured alongside it.
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	started := time.Now()
	runErr := cmd.Run()
	duration := time.Since(started)
	stream := buf.String()
	if writeErr := os.WriteFile(streamPath, []byte(stream), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("spacedock claude did not finish within %s for %s; artifacts in %s", scenario.timeout, scenario.name, artifactDir)
	}

	// Extract the final message from the stream's result/success event (the
	// front-door analog of Codex --output-last-message). A 401/is_error result is
	// surfaced here as a LOUD launch failure distinct from a scenario-assertion
	// failure, so a stale credential never feeds the 401 text into an assertion.
	finalMessage, extractErr := extractClaudeFinalMessage(stream)
	if extractErr != nil {
		t.Fatalf("claude launch failed for %s (run err=%v): %v; artifacts in %s\nStream tail:\n%s",
			scenario.name, runErr, extractErr, artifactDir, tail(stream, 4000))
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
