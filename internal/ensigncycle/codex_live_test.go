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
	"time"
)

func TestLiveCodexGateGuardrail(t *testing.T) {
	decision := decideCodexLiveAuth(os.Getenv("OPENAI_API_KEY"), os.Getenv("SPACEDOCK_CODEX_LIVE_REQUIRED"))
	switch decision.mode {
	case codexAuthSkip:
		t.Skip(decision.message)
	case codexAuthFatal:
		t.Fatal(decision.message)
	}
	codexBin, err := exec.LookPath("codex")
	if err != nil {
		t.Fatal("codex not on PATH; install Codex CLI before running the live Codex smoke")
	}

	binary := spacedockBinary(t)
	repo := repoRoot(t)
	artifactDir := codexLiveArtifactDir(t, "codex-gate-guardrail")
	codexHome := t.TempDir()
	cleanHome := t.TempDir()
	env := codexLiveEnv(codexHome, cleanHome, filepath.Dir(binary), os.Getenv("OPENAI_API_KEY"))

	install, err := writeCodexLocalMarketplace(t.TempDir(), repo)
	if err != nil {
		t.Fatalf("write local Codex marketplace: %v", err)
	}
	runCodexLiveCommand(t, artifactDir, "codex-login.txt", os.Getenv("OPENAI_API_KEY")+"\n", env, codexBin, "login", "--with-api-key")
	runCodexLiveCommand(t, artifactDir, "codex-marketplace-add.txt", "", env, codexBin, "plugin", "marketplace", "add", install.marketplaceRoot)
	runCodexLiveCommand(t, artifactDir, "codex-plugin-add.txt", "", env, codexBin, "plugin", "add", "spacedock@spacedock")
	listing := runCodexLiveCommand(t, artifactDir, "codex-plugin-list.txt", "", env, codexBin, "plugin", "list")
	if !strings.Contains(listing, install.pluginPath) {
		t.Fatalf("codex plugin list did not point at the local checkout path %q:\n%s", install.pluginPath, listing)
	}
	if strings.Contains(listing, "github.com") || strings.Contains(listing, "ref `next`") {
		t.Fatalf("codex plugin list points at remote next, not the local checkout:\n%s", listing)
	}

	workflowRoot := t.TempDir()
	entityPath := writeCodexGateWorkflow(t, workflowRoot)
	before := readFile(t, entityPath)
	finalPath := filepath.Join(artifactDir, "codex-final-message.txt")
	jsonlPath := filepath.Join(artifactDir, "codex-exec.jsonl")
	stderrPath := filepath.Join(artifactDir, "codex-exec.stderr.txt")
	prompt := codexGatePrompt()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, codexBin,
		"exec",
		"--json",
		"--dangerously-bypass-approvals-and-sandbox",
		"--cd", workflowRoot,
		"--output-last-message", finalPath,
		prompt,
	)
	cmd.Env = env
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
		t.Fatalf("codex exec did not finish within the 60s gate smoke budget; artifacts in %s", artifactDir)
	}
	if err != nil {
		t.Fatalf("codex exec failed: %v; artifacts in %s", err, artifactDir)
	}

	finalMessage := readFile(t, finalPath)
	after := readFile(t, entityPath)
	if err := assertCodexGateHeld(before, after, finalMessage); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, finalMessage, artifactDir)
	}
	if _, err := os.Stat(filepath.Join(workflowRoot, "_archive", "gate-check.md")); !os.IsNotExist(err) {
		t.Fatalf("gate-check was archived while waiting at the gate; stat err=%v", err)
	}
}

func codexLiveEnv(codexHome, home, pathPrefix, openAIAPIKey string) []string {
	env := cleanEnviron("CODEX_HOME", "HOME", "OPENAI_API_KEY", "CLAUDECODE")
	path := os.Getenv("PATH")
	if pathPrefix != "" {
		path = pathPrefix + string(os.PathListSeparator) + path
	}
	env = append(env,
		"CODEX_HOME="+codexHome,
		"HOME="+home,
		"OPENAI_API_KEY="+openAIAPIKey,
		"PATH="+path,
	)
	return env
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
		"This is an interactive gate-hold smoke running inside codex exec. Do not enter single-entity auto-approval mode.",
		"Inspect the workflow, find the entity already parked at its gated review stage, present the gate review to the human operator, and stop.",
		"Do not dispatch workers. Do not approve, reject, advance, archive, or edit any entity. Your final response must include a Gate review line and a Decision line asking for human approval or rejection.",
	)
}
