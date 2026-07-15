package ensigncycle

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const codexScenarioTimeout = 15 * time.Minute

type codexScenarioResult struct {
	finalMessage string
	jsonl        string
	artifactDir  string
	duration     time.Duration
	exitCode     int
	timedOut     bool
}

type codexProcessSpec struct {
	bin         string
	argv        []string
	env         []string
	artifactDir string
	finalPath   string
	timeout     time.Duration
}

// runCodexProcess is the complete Codex scenario process boundary. It starts one
// command under one fixed wall-clock deadline and writes stdout and stderr
// directly to the scenario artifacts. Runtime events and workflow writes are not
// observed until after Run returns, so they cannot alter liveness or cause a retry.
func runCodexProcess(spec codexProcessSpec) (codexScenarioResult, error) {
	result := codexScenarioResult{artifactDir: spec.artifactDir, exitCode: -1}
	if err := os.MkdirAll(spec.artifactDir, 0o755); err != nil {
		return result, fmt.Errorf("create Codex artifact directory: %w", err)
	}

	jsonlPath := filepath.Join(spec.artifactDir, "codex-exec.jsonl")
	stderrPath := filepath.Join(spec.artifactDir, "codex-exec.stderr.txt")
	stdout, err := os.Create(jsonlPath)
	if err != nil {
		return result, fmt.Errorf("create Codex JSONL artifact: %w", err)
	}
	stderr, err := os.Create(stderrPath)
	if err != nil {
		_ = stdout.Close()
		return result, fmt.Errorf("create Codex stderr artifact: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), spec.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, spec.bin, spec.argv...)
	cmd.Env = spec.env
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	started := time.Now()
	runErr := cmd.Run()
	result.duration = time.Since(started)
	result.timedOut = errors.Is(ctx.Err(), context.DeadlineExceeded)
	if runErr == nil {
		result.exitCode = 0
	} else {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			result.exitCode = exitErr.ExitCode()
		}
	}
	closeErr := errors.Join(stdout.Close(), stderr.Close())

	processResult := fmt.Sprintf(
		"started: %s\nexit_code: %d\ntimed_out: %t\nduration: %s\n",
		started.UTC().Format(time.RFC3339Nano), result.exitCode, result.timedOut, result.duration,
	)
	resultErr := os.WriteFile(filepath.Join(spec.artifactDir, "codex-process-result.txt"), []byte(processResult), 0o644)
	jsonl, readErr := os.ReadFile(jsonlPath)
	result.jsonl = string(jsonl)
	if final, err := os.ReadFile(spec.finalPath); err == nil {
		result.finalMessage = string(final)
	} else if !errors.Is(err, os.ErrNotExist) {
		readErr = errors.Join(readErr, err)
	}

	if artifactErr := errors.Join(closeErr, resultErr, readErr); artifactErr != nil {
		return result, fmt.Errorf("finalize Codex process artifacts: %w", artifactErr)
	}
	if runErr != nil {
		if result.timedOut {
			return result, fmt.Errorf("codex exec exceeded fixed %s deadline", spec.timeout)
		}
		return result, fmt.Errorf("codex exec exited %d: %w", result.exitCode, runErr)
	}
	return result, nil
}

// captureCodexRejectionEvidence snapshots the durable rejection-flow outcome
// after the sole Codex process exits. These artifacts are diagnostics only; none
// of them participate in process liveness.
func captureCodexRejectionEvidence(workflowRoot, entityPath, artifactDir string) (string, error) {
	entity, entityErr := os.ReadFile(entityPath)
	var errs []error
	if entityErr != nil {
		errs = append(errs, fmt.Errorf("read rejection entity: %w", entityErr))
	} else if err := os.WriteFile(filepath.Join(artifactDir, "rejection-task.after.md"), entity, 0o644); err != nil {
		errs = append(errs, fmt.Errorf("write rejection entity artifact: %w", err))
	}

	commands := []struct {
		name string
		args []string
	}{
		{"git-head.txt", []string{"rev-parse", "HEAD"}},
		{"git-log.txt", []string{"log", "--oneline", "--decorate", "-20"}},
		{"git-status.txt", []string{"status", "--short"}},
	}
	for _, artifact := range commands {
		argv := append([]string{"-C", workflowRoot}, artifact.args...)
		out, err := exec.Command("git", argv...).CombinedOutput()
		if writeErr := os.WriteFile(filepath.Join(artifactDir, artifact.name), out, 0o644); writeErr != nil {
			errs = append(errs, fmt.Errorf("write %s: %w", artifact.name, writeErr))
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("capture %s: %w", artifact.name, err))
		}
	}
	return string(entity), errors.Join(errs...)
}

func TestCodexSingleRunPreservesFaultEvidence(t *testing.T) {
	cases := []struct {
		name         string
		mode         string
		timeout      time.Duration
		wantExit     int
		wantTimedOut bool
	}{
		{name: "nonzero exit", mode: "exit-23", timeout: 2 * time.Second, wantExit: 23},
		{name: "hard deadline", mode: "stall", timeout: 2 * time.Second, wantExit: -1, wantTimedOut: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			workflowRoot := t.TempDir()
			entityPath := filepath.Join(workflowRoot, "rejection-task.md")
			writeFile(t, entityPath, "# Rejection Task\n")
			gitInit(t, workflowRoot)
			artifactRoot := t.TempDir()
			artifactDir := filepath.Join(artifactRoot, "rejection-flow")
			invocations := filepath.Join(t.TempDir(), "invocations.txt")
			finalPath := filepath.Join(artifactDir, "codex-final-message.txt")

			result, runErr := runCodexProcess(codexProcessSpec{
				bin: os.Args[0],
				argv: []string{
					"-test.run=^TestCodexSingleRunHelperProcess$", "--",
					tc.mode, workflowRoot, invocations,
				},
				env:         append(os.Environ(), "GO_WANT_CODEX_SINGLE_RUN_HELPER=1"),
				artifactDir: artifactDir,
				finalPath:   finalPath,
				timeout:     tc.timeout,
			})
			entityAfter, captureErr := captureCodexRejectionEvidence(workflowRoot, entityPath, artifactDir)
			if runErr == nil {
				t.Fatal("fault-injected scenario unexpectedly passed")
			}
			if captureErr != nil {
				t.Fatalf("capture post-run evidence: %v", captureErr)
			}
			if got := strings.Fields(readFile(t, invocations)); len(got) != 1 {
				t.Fatalf("codex invocation records = %v, want exactly one", got)
			}
			if result.exitCode != tc.wantExit || result.timedOut != tc.wantTimedOut {
				t.Fatalf("process classification = exit %d timeout %t, want exit %d timeout %t", result.exitCode, result.timedOut, tc.wantExit, tc.wantTimedOut)
			}
			if !strings.Contains(result.jsonl, `"tool":"wait_agent"`) {
				t.Fatalf("JSONL artifact did not retain wait-shaped output: %q", result.jsonl)
			}
			if got := readFile(t, filepath.Join(artifactDir, "codex-exec.stderr.txt")); !strings.Contains(got, "single-run stderr marker") {
				t.Fatalf("stderr artifact missing marker: %q", got)
			}
			process := readFile(t, filepath.Join(artifactDir, "codex-process-result.txt"))
			for _, want := range []string{fmt.Sprintf("exit_code: %d", tc.wantExit), fmt.Sprintf("timed_out: %t", tc.wantTimedOut)} {
				if !strings.Contains(process, want) {
					t.Fatalf("process result missing %q:\n%s", want, process)
				}
			}
			if !strings.Contains(entityAfter, "fault-injected durable state") {
				t.Fatalf("entity evidence missing durable state: %q", entityAfter)
			}
			for _, name := range []string{"rejection-task.after.md", "git-head.txt", "git-log.txt", "git-status.txt"} {
				if _, err := os.Stat(filepath.Join(artifactDir, name)); err != nil {
					t.Fatalf("missing %s: %v", name, err)
				}
			}
			if log := readFile(t, filepath.Join(artifactDir, "git-log.txt")); !strings.Contains(log, "fault: durable first-run state") {
				t.Fatalf("git log did not retain the helper commit: %q", log)
			}
			attempts, err := filepath.Glob(filepath.Join(artifactRoot, "**", "attempt-*"))
			if err != nil || len(attempts) != 0 {
				t.Fatalf("unexpected retry artifact directories: %v (glob error %v)", attempts, err)
			}
		})
	}
}

func TestCodexSingleRunHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_CODEX_SINGLE_RUN_HELPER") != "1" {
		return
	}
	args := helperArgs(os.Args)
	if len(args) != 3 {
		fmt.Fprintf(os.Stderr, "helper args = %v\n", args)
		os.Exit(2)
	}
	mode, workflowRoot, invocationPath := args[0], args[1], args[2]
	invocations, err := os.OpenFile(invocationPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	_, writeErr := invocations.WriteString("1\n")
	closeInvocationsErr := invocations.Close()
	if err := errors.Join(writeErr, closeInvocationsErr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Fprintln(os.Stderr, "single-run stderr marker")
	fmt.Println(`{"type":"item.started","item":{"type":"collab_tool_call","tool":"wait_agent"}}`)
	entityPath := filepath.Join(workflowRoot, "rejection-task.md")
	f, err := os.OpenFile(entityPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	_, err = fmt.Fprintln(f, "\nfault-injected durable state")
	closeErr := f.Close()
	if err != nil || closeErr != nil {
		fmt.Fprintln(os.Stderr, errors.Join(err, closeErr))
		os.Exit(2)
	}
	gitArgs := func(args ...string) error {
		full := append([]string{"-C", workflowRoot, "-c", "user.email=t@t", "-c", "user.name=t"}, args...)
		return exec.Command("git", full...).Run()
	}
	if err := gitArgs("add", "--", "rejection-task.md"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if err := gitArgs("commit", "-q", "-m", "fault: durable first-run state", "--", "rejection-task.md"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if mode == "exit-23" {
		os.Exit(23)
	}
	for {
		fmt.Println(`{"type":"item.started","item":{"type":"collab_tool_call","tool":"wait_agent"}}`)
		time.Sleep(20 * time.Millisecond)
	}
}

func helperArgs(argv []string) []string {
	for i, arg := range argv {
		if arg == "--" {
			return argv[i+1:]
		}
	}
	return nil
}

func passingRejectionEntity() string {
	return "---\nstatus: validation\n---\n" +
		rejectionFixMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: Initial implementation\n\n" +
		"## Stage Report: implementation\n\n- DONE: Applied rejection fix\n\n" +
		"### Feedback Cycles\n\n- Cycle 1: REJECTED\n- Cycle 2: PASSED\n"
}

func passingRejectionEntityWithValidationReports() string {
	return "---\nstatus: validation\n---\n" +
		rejectionFixMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: Initial implementation\n\n" +
		"## Stage Report: validation\n\nRecommendation: REJECTED.\n\n" +
		"## Stage Report: implementation\n\n- DONE: Applied rejection fix\n\n" +
		"## Stage Report: validation\n\nRecommendation: PASSED.\n\n" +
		"### Feedback Cycles\n\n- Cycle 1: REJECTED\n- Cycle 2: PASSED\n"
}
