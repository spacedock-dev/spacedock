//go:build live

package ensigncycle

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLivePiRecordedGateLifecycle is the Pi binding for the same host-neutral
// fixture, prompt, command trace, durable state, and dispatch oracle used by the
// Claude and Codex runners. The outer launcher is the fresh binary; the nested FO
// resolves the logging shim first on PATH and the shim delegates every call back
// to that exact binary.
func TestLivePiRecordedGateLifecycle(t *testing.T) {
	piBin := piBinaryOrSkip(t)
	repo := repoRoot(t)
	piSubagentsRoot := piSubagentsPackageRoot(t)
	binary := buildRecordedGateBinary(t)
	fixture := writePreparedRecordedGateFixture(t)
	before := readFile(t, fixture.entity)

	piHome := t.TempDir()
	sessionDir := t.TempDir()
	cleanHome := t.TempDir()
	seedPiLocalAuth(t, piHome, os.Getenv("HOME"))
	artifactDir := filepath.Join(piLiveArtifactDir(t, "pi-recorded-gate-lifecycle"), "run")
	if err := os.MkdirAll(filepath.Join(artifactDir, "sessions"), 0o755); err != nil {
		t.Fatal(err)
	}
	commandLog := filepath.Join(fixture.root, "evidence", "command.log")
	shimDir := writeRecordedGateLoggingShim(t, binary, commandLog)
	env := piLiveEnv(piHome, sessionDir, cleanHome, filepath.Dir(binary), piSubagentsRoot)
	env = withPATHPrefix(env, shimDir)
	env = withRecordedGateEnv(env, "SPACEDOCK_BIN", filepath.Join(shimDir, "spacedock"))
	extension := filepath.Join(piSubagentsRoot, "src", "extension", "index.ts")
	if _, err := os.Stat(extension); os.IsNotExist(err) {
		extension = filepath.Join(piSubagentsRoot, "index.ts")
	}

	runPiLiveCommand(t, artifactDir, fixture.root, env, piBin,
		"--print",
		"--session-dir", filepath.Join(artifactDir, "sessions"),
		"--extension", extension,
		"--skill", filepath.Join(piSubagentsRoot, "skills", "pi-subagents"),
		"--skill", filepath.Join(repo, "skills", "first-officer"),
		"--skill", filepath.Join(repo, "skills", "fo-gate-lifecycle"),
		"--skill", filepath.Join(repo, "skills", "present-gate"),
		"--skill", filepath.Join(repo, "skills", "ensign"),
		recordedGatePrompt(fixture.root)+"\n\nPi harness requirement: stamp the successor dispatch with explicit model `"+envOr("SPACEDOCK_PI_LIVE_CHILD_MODEL", "openrouter/openai/gpt-4.1-mini")+"`.",
	)

	session := readRecordedGatePiRootSession(t, artifactDir)
	observation := recordedGateLiveObservation(t, fixture, before, commandLog, recordedGateReviewFromPiSession(session))
	if err := assertRecordedGateLifecycle(observation); err != nil {
		t.Fatalf("Pi recorded gate lifecycle graded FAIL: %v; artifacts in %s\n--- entity after ---\n%s", err, artifactDir, observation.after)
	}
}

func readRecordedGatePiRootSession(t *testing.T, artifactDir string) string {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(artifactDir, "sessions", "*.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 1 {
		t.Fatalf("Pi root sessions=%d, want exactly one flat JSONL: %v", len(paths), paths)
	}
	body, err := os.ReadFile(paths[0])
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}
