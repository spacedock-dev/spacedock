//go:build live

package ensigncycle

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
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
		recordedGatePrompt(fixture.root)+"\n\nPi host requirement: make the successor subagent call with `async: true` and explicit model `"+envOr("SPACEDOCK_PI_LIVE_CHILD_MODEL", "openrouter/openai/gpt-4.1-mini")+"`; omit the `output` argument entirely, and do not substitute a synchronous retry.",
	)

	handle, childSession := waitForRecordedGatePiChild(t, readRecordedGatePiSessions(t, artifactDir))
	after := resolveRecordedGateEntity(fixture)
	events := recordedGateEventsFromCommandLog(readFile(t, commandLog))
	session := readRecordedGatePiSessions(t, artifactDir)
	if err := assertRecordedGateRuntimeLoadOrder("pi", session); err != nil {
		t.Fatalf("Pi recorded gate lifecycle load order graded FAIL: %v; artifacts in %s", err, artifactDir)
	}
	_, review := recordedGatePiObservation(session, after)
	dispatch := recordedGateDispatchProof{spawned: handle != "", handle: handle, workerOutput: strings.Contains(childSession, recordedGateDispatchMarker) && strings.Contains(after, recordedGateDispatchMarker)}
	if err := assertRecordedGateLifecycle(recordedGateObservation{
		events: events, before: before, after: after, dispatch: dispatch,
		gateReview: review, expectedNext: "handoff",
	}); err != nil {
		t.Fatalf("Pi recorded gate lifecycle graded FAIL: %v; artifacts in %s\n--- entity after ---\n%s", err, artifactDir, after)
	}
}

func waitForRecordedGatePiChild(t *testing.T, session string) (string, string) {
	t.Helper()
	matches := regexp.MustCompile(`"asyncDir":"([^"]+)"`).FindAllStringSubmatch(session, -1)
	if len(matches) == 0 {
		t.Fatal("Pi successor spawn returned no async directory")
	}
	dir := matches[len(matches)-1][1]
	for deadline := time.Now().Add(8 * time.Minute); time.Now().Before(deadline); time.Sleep(time.Second) {
		body, err := os.ReadFile(filepath.Join(dir, "status.json"))
		if err != nil || strings.Contains(string(body), `"state": "running"`) || strings.Contains(string(body), `"state": "queued"`) {
			continue
		}
		sessions := regexp.MustCompile(`"sessionFile":\s*"([^"]+)"`).FindAllStringSubmatch(string(body), -1)
		if !strings.Contains(string(body), `"state": "complete"`) || filepath.Base(dir) == "" || len(sessions) == 0 {
			t.Fatalf("Pi async successor ended without a correlated completion: %s", body)
		}
		return filepath.Base(dir), readFile(t, sessions[len(sessions)-1][1])
	}
	t.Fatal("Pi async successor did not complete within 8 minutes")
	return "", ""
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
