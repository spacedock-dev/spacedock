//go:build live

package ensigncycle

import (
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
	piBin := piBinaryOrSkip(t)
	repo := repoRoot(t)
	piSubagentsRoot := piSubagentsPackageRoot(t)
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
		"--skill", filepath.Join(repo, "skills", "fo-gate-lifecycle"), "--provider", envOr("SPACEDOCK_PI_LIVE_PROVIDER", "openai-codex"), "--model", envOr("SPACEDOCK_PI_LIVE_MODEL", "gpt-5.3-codex"),
		"--skill", filepath.Join(repo, "skills", "ensign"),
		recordedGatePrompt(fixture.root),
	)

	after := resolveRecordedGateEntity(fixture)
	events := recordedGateEventsFromCommandLog(readFile(t, commandLog))
	session := readRecordedGatePiSessions(t, artifactDir)
	if err := assertRecordedGateRuntimeLoadOrder("pi", session); err != nil {
		t.Fatalf("Pi recorded gate lifecycle load order graded FAIL: %v; artifacts in %s", err, artifactDir)
	}
	dispatch, review := recordedGatePiObservation(session, after)
	if err := assertRecordedGateLifecycle(recordedGateObservation{
		events: events, before: before, after: after, dispatch: dispatch,
		gateReview: review, expectedNext: "handoff",
	}); err != nil {
		t.Fatalf("Pi recorded gate lifecycle graded FAIL: %v; artifacts in %s\n--- entity after ---\n%s", err, artifactDir, after)
	}
}

func readRecordedGatePiSessions(t *testing.T, artifactDir string) string {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(artifactDir, "sessions", "*.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var session strings.Builder
	for _, path := range paths {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		session.Write(body)
		session.WriteByte('\n')
	}
	return session.String()
}
