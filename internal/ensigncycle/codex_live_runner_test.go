//go:build live

package ensigncycle

import (
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
// checkout-candidate plugin setup, and the `--output-last-message` observed-extract are the ONLY
// Codex-specific surface; the common declarations, fixtures, prompts, and assertions
// are shared with the Claude runner.
type codexLiveRunner struct {
	binary       string
	pluginDir    string
	codexBin     string
	codexHome    string
	env          []string
	artifactRoot string
}

type codexAsLiveDriver struct {
	t      *testing.T
	runner codexLiveRunner
}

func (d codexAsLiveDriver) run(t *testing.T, scenario sharedRuntimeScenario, root, prompt string) liveResult {
	result, err := d.runner.run(t, scenario, root, prompt)
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}
	commands := successfulCodexCommands(result.jsonl)
	if scenario.name == "filing" {
		command, gradeErr := assertCodexPublicFilingTransaction(result.jsonl, filepath.Join(root, filingSlug+".md"), filingSlug)
		if gradeErr != nil {
			t.Logf("Codex public filing transaction not established: %v", gradeErr)
			commands = nil
		} else {
			commands = []string{command}
		}
	}
	return liveResult{finalMessage: result.finalMessage, stream: result.jsonl, commands: commands, artifactDir: result.artifactDir, duration: result.duration}
}
func (d codexAsLiveDriver) emitMetrics(t *testing.T, scenario sharedRuntimeScenario, result liveResult) {
	emitCodexScenarioMetrics(t, scenario, codexScenarioResult{finalMessage: result.finalMessage, jsonl: result.stream, artifactDir: result.artifactDir, duration: result.duration})
}
func (d codexAsLiveDriver) gradeShallowBootObservation(*testing.T, liveResult) {}

// Codex's public stream carries no sub-agent lifecycle; the spawn_agent call and
// its completion live in the CODEX_HOME rollout, so the lifecycle stream is the
// native one resolved from the public stream's thread id.
func (d codexAsLiveDriver) lifecycleStream(t *testing.T, result liveResult) string {
	return nativeLifecycleStream(t, d, result)
}
func (d codexAsLiveDriver) prepareRecordedGate(*testing.T) (liveDriver, func(liveResult)) {
	return d, noLiveGrade
}
func (d codexAsLiveDriver) model() string { return envOr("SPACEDOCK_CODEX_LIVE_MODEL", "codex") }
func (d codexAsLiveDriver) home() string  { return d.runner.codexHome }
func (d codexAsLiveDriver) smallestMechanismTrace(result liveResult, edits, commissioned []string) mechanismTrace {
	return smallestMechanismTraceForDialect("codex", result.stream, edits, commissioned)
}
func (d codexAsLiveDriver) withStubPATH(dir string) liveDriver {
	d.runner = d.runner.withStubPATH(d.t, dir)
	return d
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

func newCodexLiveRunner(t *testing.T, setupIDs ...string) codexLiveRunner {
	t.Helper()
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	oauthJSON := os.Getenv("CODEX_AUTH_JSON")
	realHome := os.Getenv("HOME")
	decision := decideCodexLiveAuthForCI(oauthJSON, openAIAPIKey, codexLocalAuthAvailable(realHome), os.Getenv("SPACEDOCK_CODEX_LIVE_REQUIRED"))
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
	switch decision.mode {
	case codexAuthOAuth:
		if err := seedCodexOAuthAuth(codexHome, oauthJSON); err != nil {
			t.Fatalf("seed Codex OAuth auth: %v", err)
		}
	case codexAuthLocal:
		if err := seedCodexLocalAuth(codexHome, realHome); err != nil {
			t.Fatalf("seed local Codex auth: %v", err)
		}
	}
	env := codexLiveEnv(codexHome, cleanHome, filepath.Dir(binary), openAIAPIKey, decision.mode)

	setupID := ""
	if len(setupIDs) > 0 {
		setupID = setupIDs[0]
	}
	setupDir := codexLiveSetupArtifactDir(artifactRoot, setupID)
	if err := os.MkdirAll(setupDir, 0o755); err != nil {
		t.Fatal(err)
	}
	switch decision.mode {
	case codexAuthAPIKey:
		runCodexLiveCommand(t, setupDir, "codex-login.txt", openAIAPIKey+"\n", env, codexBin, "login", "--with-api-key")
	case codexAuthOAuth, codexAuthLocal:
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

	return codexLiveRunner{binary: binary, pluginDir: livePluginDir(t), codexBin: codexBin, codexHome: codexHome, env: env, artifactRoot: artifactRoot}
}

func codexLiveSetupArtifactDir(artifactRoot, setupID string) string {
	if setupID == "" {
		return filepath.Join(artifactRoot, "_setup")
	}
	return filepath.Join(artifactRoot, "_setup", setupID)
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

// run launches one `codex exec --json` for one shared scenario. Each complete
// JSONL line resets the shared quiet budget. Stream silence kills that sole
// process; activity and durable writes never trigger another launch.
func (r codexLiveRunner) run(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, prompt string) (codexScenarioResult, error) {
	t.Helper()
	artifactDir := filepath.Join(r.artifactRoot, scenario.name)
	finalPath := filepath.Join(artifactDir, "codex-final-message.txt")
	return runCodexProcess(codexProcessSpec{
		bin:         r.binary,
		argv:        codexLiveFrontDoorArgvForScenario(r.pluginDir, workflowRoot, finalPath, prompt, scenario.name),
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
