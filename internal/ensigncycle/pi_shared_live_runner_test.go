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

type piSharedLiveAdapter struct{}

func (piSharedLiveAdapter) runtimeName() string { return "pi" }

func (piSharedLiveAdapter) runSharedScenario(t *testing.T, scenario sharedRuntimeScenario) {
	driver := newPiSharedLiveDriver(t)
	switch scenario.name {
	case "full-ensign-cycle":
		runFullEnsignCycleJourney(t, driver, scenario)
	case "gate-guardrail", "default-headless-gate-stop":
		runClaudeGateGuardrailScenario(t, driver, scenario)
	case "withdrawn-gate-recovery":
		runClaudeWithdrawnGateRecoveryScenario(t, driver, scenario)
	case "recorded-gate-lifecycle":
		runClaudeRecordedGateLifecycleScenario(t, driver, scenario)
	case "rejection-flow":
		runClaudeRejectionFlowScenario(t, driver, scenario)
	case "feedback-3-cycle-escalation":
		runClaudeFeedback3CycleEscalationScenario(t, driver, scenario)
	case "merge-hook-guardrail":
		runClaudeMergeHookGuardrailScenario(t, driver, scenario)
	case "filing":
		runClaudeFilingScenario(t, driver, scenario)
	case "shallow-boot":
		runClaudeShallowBootScenario(t, driver, scenario)
	case "zero-discovery":
		runZeroDiscoveryJourney(t, driver, scenario)
	case "auto-continue-after-implementation":
		runAutoContinueJourney(t, driver, scenario)
	case "self-evidence-merge-triage":
		runClaudeSelfEvidenceMergeTriageScenario(t, driver, scenario)
	case "smallest-sufficient-mechanism":
		runClaudeSmallestSufficientMechanismScenario(t, driver, scenario)
	case "keep-moving-posture":
		runClaudeKeepMovingScenario(t, driver, scenario)
	case "ac-value-reanchor":
		runACValueReanchorJourney(t, driver, scenario)
	default:
		t.Fatalf("unknown shared journey %q", scenario.name)
	}
}

type piSharedLiveDriver struct {
	binary       string
	pluginDir    string
	env          []string
	modelName    string
	artifactRoot string
	piHome       string
}

func newPiSharedLiveDriver(t *testing.T) piSharedLiveDriver {
	t.Helper()
	repo := repoRoot(t)
	packageRoot := piSubagentsPackageRoot(t)
	binary := piSpacedockBinary(t, repo)
	piHome := t.TempDir()
	seedPiLiveAuth(t, piHome, os.Getenv("HOME"), os.Getenv("OPENAI_API_KEY"), os.Getenv("SPACEDOCK_PI_LIVE_REQUIRED"))
	writeFile(t, filepath.Join(piHome, "settings.json"), fmt.Sprintf("{\"packages\":[%q]}\n", "file:"+repo))
	return piSharedLiveDriver{
		binary:       binary,
		pluginDir:    repo,
		env:          piLiveEnv(piHome, t.TempDir(), t.TempDir(), filepath.Dir(binary), packageRoot),
		modelName:    piLiveModelName(),
		artifactRoot: piLiveArtifactDir(t, "pi-shared-scenarios"),
		piHome:       piHome,
	}
}

func (d piSharedLiveDriver) runtimeName() string { return "pi" }
func (d piSharedLiveDriver) model() string       { return d.modelName }
func (d piSharedLiveDriver) home() string        { return d.piHome }
func (d piSharedLiveDriver) withStubPATH(dir string) liveDriver {
	d.env = withPATHPrefix(d.env, dir)
	return d
}

func (d piSharedLiveDriver) run(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, prompt string) liveResult {
	t.Helper()
	artifactDir := filepath.Join(d.artifactRoot, scenario.name)
	sessionDir := filepath.Join(artifactDir, "sessions")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, d.binary,
		"pi", prompt,
		"--plugin-dir", d.pluginDir,
		"--",
		"--print",
		"--model", d.modelName,
		"--session-dir", sessionDir,
	)
	cmd.Dir = workflowRoot
	cmd.Env = d.env
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	started := time.Now()
	err := cmd.Run()
	duration := time.Since(started)
	writeFile(t, filepath.Join(artifactDir, "pi-stdout.txt"), stdout.String())
	writeFile(t, filepath.Join(artifactDir, "pi-stderr.txt"), stderr.String())
	writeFile(t, filepath.Join(artifactDir, "model.txt"), d.modelName+"\n")
	writeFile(t, filepath.Join(artifactDir, "duration.txt"), duration.String()+"\n")
	writeFile(t, filepath.Join(artifactDir, "process-status.txt"), fmt.Sprintf("error=%v timeout=%t\n", err, ctx.Err() == context.DeadlineExceeded))
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("Pi journey %q exceeded its 12-minute deadline; artifacts: %s", scenario.name, artifactDir)
	}
	if err != nil {
		t.Fatalf("Pi journey %q failed: %v; artifacts: %s\nstderr tail:\n%s", scenario.name, err, artifactDir, tail(stderr.String(), 4000))
	}
	return liveResult{
		finalMessage: stdout.String(),
		stream:       stdout.String() + "\n" + stderr.String(),
		artifactDir:  artifactDir,
		duration:     duration,
		configDir:    sessionDir,
		cwd:          workflowRoot,
	}
}
