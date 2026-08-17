//go:build live

package ensigncycle

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

type piSharedLiveDriver struct {
	t                                                  *testing.T
	binary, pluginDir, modelName, artifactRoot, piHome string
	env                                                []string
}

func newPiSharedLiveDriver(t *testing.T) piSharedLiveDriver {
	t.Helper()
	repo := repoRoot(t)
	binary := piSpacedockBinary(t, repo)
	piHome := t.TempDir()
	decision := seedPiLiveAuth(t, piHome, os.Getenv("HOME"), os.Getenv("CODEX_AUTH_JSON"), os.Getenv("OPENAI_API_KEY"), os.Getenv("SPACEDOCK_PI_LIVE_REQUIRED"))
	writeFile(t, filepath.Join(piHome, "settings.json"), fmt.Sprintf("{\"packages\":[%q]}\n", "file:"+repo))
	return piSharedLiveDriver{
		t:      t,
		binary: binary, pluginDir: repo, modelName: piLiveChildModel(decision), piHome: piHome,
		artifactRoot: piLiveArtifactDir(t, "pi-common"),
		env:          piLiveEnvForAuth(piHome, t.TempDir(), t.TempDir(), filepath.Dir(binary), piSubagentsPackageRoot(t), os.Getenv("OPENAI_API_KEY"), decision.mode),
	}
}

func (d piSharedLiveDriver) model() string { return d.modelName }
func (d piSharedLiveDriver) home() string  { return d.piHome }

// Pi records subagent toolCalls and their completion in the session transcript that
// run() folds into result.stream, so the public stream is the lifecycle stream.
func (d piSharedLiveDriver) lifecycleStream(_ *testing.T, result liveResult) string {
	return result.stream
}
func (d piSharedLiveDriver) smallestMechanismTrace(result liveResult, edits, commissioned []string) mechanismTrace {
	return claudeMechanismTrace(result.stream, edits, commissioned)
}
func (d piSharedLiveDriver) withStubPATH(dir string) liveDriver {
	d.env = withPATHPrefix(d.env, dir)
	d.env = withSpacedockShimShellEnv(d.t, d.env, dir)
	return d
}
func (d piSharedLiveDriver) emitMetrics(t *testing.T, scenario sharedRuntimeScenario, result liveResult) {
	emitPiScenarioMetrics(t, scenario, result, d.modelName)
}
func (d piSharedLiveDriver) gradeShallowBootObservation(*testing.T, liveResult) {}
func (d piSharedLiveDriver) prepareRecordedGate(*testing.T) (liveDriver, func(liveResult)) {
	return d, noLiveGrade
}

func (d piSharedLiveDriver) run(t *testing.T, scenario sharedRuntimeScenario, root, prompt string) liveResult {
	t.Helper()
	artifactDir := filepath.Join(d.artifactRoot, scenario.name)
	sessionDir := filepath.Join(artifactDir, "sessions")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), piLiveRunTimeout(12*time.Minute))
	defer cancel()
	cmd := exec.CommandContext(ctx, d.binary, "pi", prompt, "--plugin-dir", d.pluginDir, "--", "--print", "--model", d.modelName, "--session-dir", sessionDir)
	cmd.Dir, cmd.Env = root, d.env
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	started := time.Now()
	err := cmd.Run()
	duration := time.Since(started)
	writeFile(t, filepath.Join(artifactDir, "pi-stdout.txt"), stdout.String())
	writeFile(t, filepath.Join(artifactDir, "pi-stderr.txt"), stderr.String())
	writeFile(t, filepath.Join(artifactDir, "model.txt"), d.modelName+"\n")
	writeFile(t, filepath.Join(artifactDir, "duration.txt"), duration.String()+"\n")
	writeFile(t, filepath.Join(artifactDir, "process-status.txt"), fmt.Sprintf("error=%v timeout=%t\n", err, ctx.Err() == context.DeadlineExceeded))
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("Pi journey %q exceeded the per-run cap (SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES, default 12m); artifacts: %s", scenario.name, artifactDir)
	}
	if err != nil {
		t.Fatalf("Pi journey %q failed: %v; artifacts: %s\n%s", scenario.name, err, artifactDir, tail(stderr.String(), 4000))
	}
	rootSession := onePiSession(t, filepath.Join(sessionDir, "*.jsonl"), "root")
	stderr.WriteString("\n" + readFile(t, rootSession))
	return liveResult{finalMessage: stdout.String(), stream: stdout.String() + "\n" + stderr.String(), commands: piObservedCommands(t, rootSession), artifactDir: artifactDir, duration: duration}
}

func piObservedCommands(t *testing.T, sessionPath string) []string {
	return piTranscriptToolValues(t, sessionPath, "bash", "shell")
}
