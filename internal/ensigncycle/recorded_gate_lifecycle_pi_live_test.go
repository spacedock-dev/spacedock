//go:build live

package ensigncycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLivePiRecordedGateLifecycle is the Pi binding for the same host-neutral
// fixture, prompt, command trace, durable state, and dispatch oracle used by the
// Claude and Codex runners. The outer launcher is the fresh binary; the nested FO
// resolves the logging shim first on PATH and the shim delegates every call back
// to that exact binary.
func TestLivePiRecordedGateLifecycle(t *testing.T) {
	t.Skip("TODO(9w59t6m1qc46hccd54p04z2j): delegated gate presentation-to-application/dispatch is unreliable")
	repo := repoRoot(t)
	piSubagentsRoot := piSubagentsPackageRoot(t)
	expectedChildModel := piLiveModelName()
	binary := buildRecordedGateBinary(t)
	fixture := writeRecordedGateFixture(t)
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
	bashEnv := filepath.Join(shimDir, "recorded-gate-env.sh")
	writeFile(t, bashEnv, "export SPACEDOCK_BIN="+filepath.Join(shimDir, "spacedock")+"\n")
	env := piLiveEnv(piHome, sessionDir, cleanHome, filepath.Dir(binary), piSubagentsRoot)
	env = withPATHPrefix(env, shimDir)
	env = withRecordedGateEnv(env, "BASH_ENV", bashEnv)

	runPiLiveCommand(t, artifactDir, fixture.root, env, binary,
		"pi",
		recordedGatePrompt(fixture.root)+"\n\nPi harness requirement: stamp the successor dispatch with explicit model `"+expectedChildModel+"`.",
		"--plugin-dir", repo,
		"--",
		"--print",
		"--model", expectedChildModel,
		"--session-dir", filepath.Join(artifactDir, "sessions"),
	)

	session := readRecordedGatePiRootSession(t, artifactDir)
	assertRecordedGatePiChildModel(t, artifactDir, expectedChildModel)
	observation := recordedGateLiveObservation(t, fixture, before, commandLog, recordedGateReviewFromPiSession(session))
	if err := assertRecordedGateLifecycle(observation); err != nil {
		t.Fatalf("Pi recorded gate lifecycle graded FAIL: %v; artifacts in %s\n--- entity after ---\n%s", err, artifactDir, observation.after)
	}
}

func assertRecordedGatePiChildModel(t *testing.T, artifactDir, want string) {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(artifactDir, "sessions", "*", "*", "run-*", "session.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 1 {
		t.Fatalf("Pi child sessions=%d, want exactly one: %v", len(paths), paths)
	}
	var models []string
	for _, line := range strings.Split(readFile(t, paths[0]), "\n") {
		var event struct {
			Type     string `json:"type"`
			Provider string `json:"provider"`
			ModelID  string `json:"modelId"`
		}
		if json.Unmarshal([]byte(line), &event) == nil && event.Type == "model_change" {
			models = append(models, event.Provider+"/"+event.ModelID)
		}
	}
	if len(models) != 1 || models[0] != want {
		t.Fatalf("Pi child model changes=%v, want exactly [%s]; child session %s", models, want, paths[0])
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
