//go:build live

package ensigncycle

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type codexLiveRunner struct {
	codexBin     string
	env          []string
	artifactRoot string
}

type codexScenarioResult struct {
	finalMessage string
	jsonl        string
	artifactDir  string
}

type codexLiveScenario struct {
	codexSharedScenario
	run func(*testing.T, codexLiveRunner, codexSharedScenario)
}

func TestLiveCodexSharedScenarios(t *testing.T) {
	runner := newCodexLiveRunner(t)

	for _, scenario := range codexLiveScenarios(t) {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.run(t, runner, scenario.codexSharedScenario)
		})
	}
}

func codexLiveScenarios(t *testing.T) []codexLiveScenario {
	t.Helper()
	runners := map[string]func(*testing.T, codexLiveRunner, codexSharedScenario){
		"gate-guardrail":       runCodexGateGuardrailScenario,
		"rejection-flow":       runCodexRejectionFlowScenario,
		"merge-hook-guardrail": runCodexMergeHookGuardrailScenario,
	}

	var scenarios []codexLiveScenario
	for _, scenario := range codexSharedScenarios() {
		run := runners[scenario.name]
		if run == nil {
			t.Fatalf("shared scenario %q has no Codex live runner", scenario.name)
		}
		scenarios = append(scenarios, codexLiveScenario{
			codexSharedScenario: scenario,
			run:                 run,
		})
	}
	return scenarios
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

func runCodexGateGuardrailScenario(t *testing.T, runner codexLiveRunner, scenario codexSharedScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeCodexGateWorkflow(t, workflowRoot)
	before := readFile(t, entityPath)

	result := runner.run(t, scenario, workflowRoot, codexGatePrompt())
	after := readFile(t, entityPath)
	if err := assertCodexGateHeld(before, after, result.finalMessage); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if _, err := os.Stat(filepath.Join(workflowRoot, "_archive", "gate-check.md")); !os.IsNotExist(err) {
		t.Fatalf("gate-check was archived while waiting at the gate; stat err=%v", err)
	}
}

func runCodexRejectionFlowScenario(t *testing.T, runner codexLiveRunner, scenario codexSharedScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeCodexRejectionWorkflow(t, workflowRoot)

	result := runner.run(t, scenario, workflowRoot, codexRejectionPrompt())
	after := readFile(t, entityPath)
	if err := assertCodexRejectionFlow(after, result.finalMessage+"\n"+result.jsonl); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
}

func runCodexMergeHookGuardrailScenario(t *testing.T, runner codexLiveRunner, scenario codexSharedScenario) {
	t.Helper()
	workflowRoot := t.TempDir()
	entityPath := writeCodexMergeHookGuardWorkflow(t, workflowRoot)
	before := readFile(t, entityPath)

	result := runner.run(t, scenario, workflowRoot, codexMergeHookGuardPrompt())
	after := readFile(t, entityPath)
	if err := assertCodexMergeHookGuardHeld(before, after, result.finalMessage+"\n"+result.jsonl); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if _, err := os.Stat(filepath.Join(workflowRoot, "_archive", "merge-check.md")); !os.IsNotExist(err) {
		t.Fatalf("merge-check was archived despite the guardrail scenario; stat err=%v", err)
	}
}

func (r codexLiveRunner) run(t *testing.T, scenario codexSharedScenario, workflowRoot, prompt string) codexScenarioResult {
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

	err = cmd.Run()
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

func writeCodexGateWorkflow(t *testing.T, root string) string {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), codexGateReadme())
	entityPath := filepath.Join(root, "gate-check.md")
	writeFile(t, entityPath, codexGateEntity())
	gitInit(t, root)
	return entityPath
}

func codexGateReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: draft\n" +
		"      initial: true\n" +
		"    - name: review\n" +
		"      gate: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Codex Gate Fixture\n\n" +
		"### draft\n\nCreate the draft.\n\n- **Outputs:** A draft stage report.\n\n" +
		"### review\n\nHuman approval gate.\n\n- **Outputs:** A gate review for the human operator.\n\n" +
		"### done\n\nTerminal state.\n"
}

func codexGateEntity() string {
	return "---\n" +
		"id: gate-check\n" +
		"title: Gate Check\n" +
		"status: review\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Gate Check\n\n" +
		"This fixture starts at the human review gate.\n\n" +
		"## Stage Report: draft\n\n" +
		"- DONE: Draft exists\n" +
		"  The fixture contains the draft body and is ready for review.\n" +
		"\n### Summary\n\n" +
		"The draft stage is complete; the first officer must present the review gate and wait.\n"
}

func codexGatePrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"This is an interactive gate-hold scenario running inside codex exec. Do not enter single-entity auto-approval mode.",
		"Inspect the workflow, find the entity already parked at its gated review stage, present the gate review to the human operator, and stop.",
		"Do not dispatch workers. Do not approve, reject, advance, archive, or edit any entity. Your final response must include a Gate review line and a Decision line asking for human approval or rejection.",
	)
}

func writeCodexRejectionWorkflow(t *testing.T, root string) string {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), codexRejectionReadme())
	entityPath := filepath.Join(root, "rejection-task.md")
	writeFile(t, entityPath, codexRejectionEntity())
	gitInit(t, root)
	return entityPath
}

func codexRejectionReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: backlog\n" +
		"      initial: true\n" +
		"    - name: implementation\n" +
		"    - name: validation\n" +
		"      gate: true\n" +
		"      feedback-to: implementation\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Codex Rejection Fixture\n\n" +
		"### implementation\n\n" +
		"Apply the validation rejection by appending this exact standalone line to `rejection-task.md`:\n\n" +
		"`codex-rejection-fix: applied`\n\n" +
		"Then append a `## Stage Report: implementation` section with one `- DONE:` item naming the fix.\n\n" +
		"- **Outputs:** The exact fix marker and an implementation stage report.\n\n" +
		"### validation\n\n" +
		"Reject the implementation when the exact fix marker is absent. If it is present, report PASSED.\n\n" +
		"- **Outputs:** A PASSED or REJECTED validation stage report.\n\n" +
		"### done\n\nTerminal state.\n"
}

func codexRejectionEntity() string {
	return "---\n" +
		"id: rejection-task\n" +
		"title: Rejection Task\n" +
		"status: validation\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Rejection Task\n\n" +
		"The implementation is intentionally missing the exact fix marker.\n\n" +
		"## Stage Report: implementation\n\n" +
		"- DONE: Initial implementation exists\n" +
		"  The initial implementation deliberately omits the required fix marker.\n" +
		"\n### Summary\n\n" +
		"Ready for validation.\n\n" +
		"## Stage Report: validation\n\n" +
		"- FAILED: Fix marker is absent\n" +
		"  REJECTED: expected exact line `codex-rejection-fix: applied`, but it is missing. Route this back to implementation.\n" +
		"\n### Summary\n\n" +
		"Recommendation: REJECTED. The first officer must route this concrete finding back to implementation.\n"
}

func codexRejectionPrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"Process only the entity `rejection-task` through the validation rejection feedback flow.",
		"The latest validation report already recommends REJECTED. Route the concrete finding back to the implementation target, dispatch a Codex worker if needed, wait for the follow-up implementation completion, and then stop.",
		"Do not advance the entity to validation again or to done. Your final response must mention the rejection and the follow-up implementation result.",
	)
}

func writeCodexMergeHookGuardWorkflow(t *testing.T, root string) string {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), codexMergeHookGuardReadme())
	modsDir := filepath.Join(root, "_mods")
	if err := os.MkdirAll(modsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(modsDir, "local-merge.md"), codexMergeHookMod())
	entityPath := filepath.Join(root, "merge-check.md")
	writeFile(t, entityPath, codexMergeHookGuardEntity())
	gitInit(t, root)
	return entityPath
}

func codexMergeHookGuardReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: implementation\n" +
		"      initial: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Codex Merge Hook Fixture\n\n" +
		"### implementation\n\nWork is ready.\n\n- **Outputs:** Existing implementation report.\n\n" +
		"### done\n\nTerminal state.\n"
}

func codexMergeHookMod() string {
	return "---\n" +
		"name: local-merge\n" +
		"description: Registers a merge hook for the live Codex guardrail scenario.\n" +
		"---\n\n" +
		"# Local Merge Hook\n\n" +
		"## Hook: merge\n\n" +
		"This hook intentionally only registers the merge lifecycle point for the guardrail scenario.\n"
}

func codexMergeHookGuardEntity() string {
	return "---\n" +
		"id: merge-check\n" +
		"title: Merge Check\n" +
		"status: implementation\n" +
		"completed:\n" +
		"verdict:\n" +
		"pr:\n" +
		"mod-block:\n" +
		"worktree:\n" +
		"---\n" +
		"# Merge Check\n\n" +
		"This entity must not advance to done while a merge hook is registered and no `mod-block` or `pr` signal exists.\n\n" +
		"## Stage Report: implementation\n\n" +
		"- DONE: Work exists\n" +
		"  The fixture is ready to test the terminal merge-hook guard.\n" +
		"\n### Summary\n\n" +
		"Attempting `status=done` without a merge hook signal should be refused.\n"
}

func codexMergeHookGuardPrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"This is the merge-hook guardrail scenario. First inspect startup/status so the registered merge hook is visible.",
		"Then intentionally run `spacedock status --workflow-dir . --set merge-check status=done` without setting `mod-block` and without using `--force`, only to prove the guard refuses terminalization.",
		"Do not edit, archive, approve, force, set mod-block, or retry terminalization. Your final response must include the guard error mentioning merge hooks.",
	)
}
