//go:build live

package ensigncycle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	seedPiLiveAuth(t, piHome, os.Getenv("HOME"), os.Getenv("OPENAI_API_KEY"), os.Getenv("SPACEDOCK_PI_LIVE_REQUIRED"))
	writeFile(t, filepath.Join(piHome, "settings.json"), fmt.Sprintf("{\"packages\":[%q]}\n", "file:"+repo))
	return piSharedLiveDriver{
		t:      t,
		binary: binary, pluginDir: repo, modelName: piLiveModelName(), piHome: piHome,
		artifactRoot: piLiveArtifactDir(t, "pi-common"),
		env:          piLiveEnv(piHome, t.TempDir(), t.TempDir(), filepath.Dir(binary), piSubagentsPackageRoot(t)),
	}
}

func (d piSharedLiveDriver) model() string { return d.modelName }
func (d piSharedLiveDriver) home() string  { return d.piHome }
func (d piSharedLiveDriver) withStubPATH(dir string) liveDriver {
	d.env = withPATHPrefix(d.env, dir)
	d.env = withSpacedockShimShellEnv(d.t, d.env, dir)
	return d
}
func (d piSharedLiveDriver) emitMetrics(*testing.T, sharedRuntimeScenario, liveResult) {}
func (d piSharedLiveDriver) gradeShallowBootObservation(*testing.T, liveResult)        {}

func (d piSharedLiveDriver) run(t *testing.T, scenario sharedRuntimeScenario, root, prompt string) liveResult {
	t.Helper()
	artifactDir := filepath.Join(d.artifactRoot, scenario.name)
	sessionDir := filepath.Join(artifactDir, "sessions")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Minute)
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
		t.Fatalf("Pi journey %q exceeded 12 minutes; artifacts: %s", scenario.name, artifactDir)
	}
	if err != nil {
		t.Fatalf("Pi journey %q failed: %v; artifacts: %s\n%s", scenario.name, err, artifactDir, tail(stderr.String(), 4000))
	}
	rootSession := onePiSession(t, filepath.Join(sessionDir, "*.jsonl"), "root")
	return liveResult{finalMessage: stdout.String(), stream: stdout.String() + "\n" + stderr.String(), commands: piObservedCommands(t, rootSession), artifactDir: artifactDir, duration: duration}
}

func piObservedCommands(t *testing.T, sessionPath string) []string {
	t.Helper()
	var commands []string
	for lineNo, line := range strings.Split(readFile(t, sessionPath), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var record struct {
			Message struct {
				Content json.RawMessage `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("Pi root session %s line %d is not JSON: %v", sessionPath, lineNo+1, err)
		}
		var blocks []struct {
			Type      string `json:"type"`
			Name      string `json:"name"`
			Arguments struct {
				Command string `json:"command"`
			} `json:"arguments"`
		}
		if json.Unmarshal(record.Message.Content, &blocks) != nil {
			continue
		}
		for _, block := range blocks {
			if block.Type == "toolCall" && (block.Name == "bash" || block.Name == "shell") && block.Arguments.Command != "" {
				commands = append(commands, block.Arguments.Command)
			}
		}
	}
	return commands
}
