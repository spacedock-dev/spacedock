package ensigncycle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type codexScenarioResult struct {
	finalMessage string
	jsonl        string
	artifactDir  string
	duration     time.Duration
	exitCode     int
	terminal     bool
	timedOut     bool
	lastEvent    string
}

type codexProcessSpec struct {
	bin         string
	argv        []string
	env         []string
	artifactDir string
	finalPath   string
	quietBudget time.Duration
}

// runCodexProcess is the complete Codex scenario process boundary. It starts one
// command, writes stdout and stderr to scenario artifacts, and also streams each
// complete stdout JSONL line through streamWatcher. Stream progress resets the
// quiet budget. Silence kills the sole process; no retry or second launch occurs.
func runCodexProcess(spec codexProcessSpec) (codexScenarioResult, error) {
	result := codexScenarioResult{artifactDir: spec.artifactDir, exitCode: -1}
	if spec.quietBudget <= 0 {
		return result, fmt.Errorf("Codex quiet budget must be positive")
	}
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

	streamReader, streamWriter := io.Pipe()
	cmd := exec.Command(spec.bin, spec.argv...)
	cmd.Env = spec.env
	cmd.Stdout = io.MultiWriter(stdout, streamWriter)
	cmd.Stderr = stderr

	started := time.Now()
	startErr := cmd.Start()
	var runErr, stallErr error
	var terminalObserved bool
	if startErr != nil {
		runErr = startErr
		_ = streamWriter.CloseWithError(startErr)
	} else {
		poller := newCmdPoller(cmd, streamWriter)
		watcher := newStreamWatcher(newPipeLineSource(streamReader), poller, func(string) {})
		_, terminalObserved, stallErr = drainCodexToTerminal(watcher, spec.quietBudget)
		result.exitCode, runErr = poller.wait()
		var timeout *stepTimeout
		result.timedOut = errors.As(stallErr, &timeout)
	}
	result.duration = time.Since(started)
	_ = streamReader.Close()
	closeErr := errors.Join(stdout.Close(), stderr.Close())

	jsonl, readErr := os.ReadFile(jsonlPath)
	result.jsonl = string(jsonl)
	result.lastEvent = lastCodexEvent(result.jsonl)
	// A terminal line can be written to the durable stdout artifact in the same
	// instant the quiet deadline fires, after the pipe watcher has made its last
	// drain. Reconcile that Codex-specific boundary after Wait; generic watcher
	// timeout semantics remain unchanged when no exact terminal event was emitted.
	terminalObserved = terminalObserved || codexTranscriptHasTerminal(result.jsonl)
	if final, err := os.ReadFile(spec.finalPath); err == nil {
		result.finalMessage = string(final)
	} else if !errors.Is(err, os.ErrNotExist) {
		readErr = errors.Join(readErr, err)
	}
	var terminalEvidenceErr error
	if terminalObserved {
		stallErr = nil
		result.timedOut = false
		if strings.TrimSpace(result.finalMessage) == "" {
			terminalEvidenceErr = fmt.Errorf("codex turn.completed had no non-empty final message; terminal evidence is incomplete\nLast event: %s\nArtifacts: %s", result.lastEvent, result.artifactDir)
		} else {
			result.terminal = true
		}
	}
	processResult := fmt.Sprintf(
		"started: %s\nexit_code: %d\nterminal: %t\ntimed_out: %t\nduration: %s\nlast_event: %s\n",
		started.UTC().Format(time.RFC3339Nano), result.exitCode, result.terminal, result.timedOut, result.duration, result.lastEvent,
	)
	resultErr := os.WriteFile(filepath.Join(spec.artifactDir, "codex-process-result.txt"), []byte(processResult), 0o644)

	if artifactErr := errors.Join(closeErr, resultErr, readErr); artifactErr != nil {
		return result, fmt.Errorf("finalize Codex process artifacts: %w", artifactErr)
	}
	if stallErr != nil {
		return result, fmt.Errorf("%w\nLast event: %s\nArtifacts: %s", stallErr, result.lastEvent, result.artifactDir)
	}
	if terminalEvidenceErr != nil {
		return result, terminalEvidenceErr
	}
	if result.terminal {
		// Codex can emit turn.completed while its front-door process remains alive;
		// the cleanup wait also gives --output-last-message time to flush. The
		// Codex-specific drain owns that semantic terminal boundary, so an exit error
		// from the intentional cleanup kill is not a scenario failure.
		return result, nil
	}
	if runErr != nil {
		return result, fmt.Errorf("codex exec exited %d: %w", result.exitCode, runErr)
	}
	return result, nil
}

func codexTranscriptHasTerminal(jsonl string) bool {
	for _, line := range strings.Split(jsonl, "\n") {
		var event struct {
			Type string `json:"type"`
		}
		if json.Unmarshal([]byte(line), &event) == nil && event.Type == "turn.completed" {
			return true
		}
	}
	return false
}

// drainCodexToTerminal is the Codex-specific process boundary. Codex's
// turn.completed event is the semantic end of a run, even when the front-door
// process remains alive. The final-message file is checked after the intentional
// cleanup wait in runCodexProcess. The shared streamWatcher keeps its generic
// OS-exit contract; only this Codex adapter recognizes the host-specific terminal
// event and reaps the process after killing that idle tail.
func drainCodexToTerminal(w *streamWatcher, budget time.Duration) (string, bool, error) {
	deadline := time.Now().Add(budget)
	terminalSeen := false
	for {
		entries, drained := w.drainEntries()
		if drained > 0 {
			deadline = time.Now().Add(budget)
		}
		for _, entry := range entries {
			if entry.Type == "turn.completed" {
				terminalSeen = true
			}
		}

		if terminalSeen {
			w.proc.kill()
			return w.fullTranscript(), true, nil
		}

		if _, exited := w.proc.poll(); exited {
			entries, _ = w.drainEntries()
			for _, entry := range entries {
				if entry.Type == "turn.completed" {
					terminalSeen = true
				}
			}
			if terminalSeen {
				return w.fullTranscript(), true, nil
			}
			return w.fullTranscript(), false, nil
		}

		if time.Now().After(deadline) {
			w.proc.kill()
			return w.fullTranscript(), false, &stepTimeout{
				label: "codex exec",
				msg: fmt.Sprintf("codex exec made no stream progress within %s (no-progress quiet budget) — a hung stage; killed the subprocess.\nTranscript tail:\n%s",
					budget, w.transcriptTail()),
			}
		}
		time.Sleep(w.pollInterval)
	}
}

func lastCodexEvent(jsonl string) string {
	lines := strings.Split(strings.TrimSpace(jsonl), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if event := strings.TrimSpace(lines[i]); event != "" {
			return event
		}
	}
	return "<none>"
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
		quietBudget  time.Duration
		wantExit     int
		wantTimedOut bool
	}{
		{name: "nonzero exit", mode: "exit-23", quietBudget: 2 * time.Second, wantExit: 23},
		{name: "quiet timeout", mode: "stall", quietBudget: 2 * time.Second, wantExit: -1, wantTimedOut: true},
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
				quietBudget: tc.quietBudget,
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

func codexProcessFixture(t *testing.T, name, mode string, quietBudget time.Duration) (codexProcessSpec, string) {
	t.Helper()
	workflowRoot := t.TempDir()
	writeFile(t, filepath.Join(workflowRoot, "rejection-task.md"), "# Rejection Task\n")
	gitInit(t, workflowRoot)
	artifactDir := filepath.Join(t.TempDir(), name)
	invocations := filepath.Join(t.TempDir(), "invocations.txt")
	return codexProcessSpec{
		bin: os.Args[0], argv: []string{
			"-test.run=^TestCodexSingleRunHelperProcess$", "--", mode, workflowRoot, invocations,
		},
		env: append(os.Environ(), "GO_WANT_CODEX_SINGLE_RUN_HELPER=1", "GORACE=atexit_sleep_ms=0"), artifactDir: artifactDir,
		finalPath: filepath.Join(artifactDir, "codex-final-message.txt"), quietBudget: quietBudget,
	}, invocations
}

func TestCodexProcessActivityResetsQuietBudget(t *testing.T) {
	const quietBudget = 250 * time.Millisecond
	spec, invocations := codexProcessFixture(t, "activity-reset", "progress-then-exit", quietBudget)
	result, err := runCodexProcess(spec)
	if err != nil {
		t.Fatalf("progressing process should stay alive beyond its quiet budget: %v", err)
	}
	if result.duration <= 4*quietBudget {
		t.Fatalf("helper duration = %s, want more than four quiet budgets (%s)", result.duration, 4*quietBudget)
	}
	if result.exitCode != 0 || result.timedOut {
		t.Fatalf("process classification = exit %d timeout %t, want exit 0 timeout false", result.exitCode, result.timedOut)
	}
	if got := strings.Count(strings.TrimSpace(result.jsonl), "\n") + 1; got < 5 {
		t.Fatalf("complete JSONL events = %d, want at least 5", got)
	}
	if got := strings.Fields(readFile(t, invocations)); len(got) != 1 {
		t.Fatalf("codex invocation records = %v, want exactly one", got)
	}
}

func TestCodexProcessQuietTimeoutPreservesFaultEvidence(t *testing.T) {
	const quietBudget = 250 * time.Millisecond
	spec, invocations := codexProcessFixture(t, "quiet-timeout", "stall", quietBudget)
	result, runErr := runCodexProcess(spec)
	if runErr == nil {
		t.Fatal("silent process unexpectedly passed")
	}
	lastEvent := strings.TrimSpace(result.jsonl)
	if !strings.Contains(lastEvent, `"sequence":1`) || strings.Contains(lastEvent, "\n") {
		t.Fatalf("partial JSONL must preserve exactly the last complete event: %q", result.jsonl)
	}
	if !strings.Contains(runErr.Error(), lastEvent) || !strings.Contains(runErr.Error(), spec.artifactDir) {
		t.Fatalf("quiet-timeout error must name last event and artifact directory: %v", runErr)
	}
	if result.exitCode != -1 || !result.timedOut {
		t.Fatalf("stalled process was not killed: exit %d timeout %t", result.exitCode, result.timedOut)
	}
	if got := readFile(t, filepath.Join(spec.artifactDir, "codex-exec.stderr.txt")); !strings.Contains(got, "single-run stderr marker") {
		t.Fatalf("stderr artifact missing marker: %q", got)
	}
	process := readFile(t, filepath.Join(spec.artifactDir, "codex-process-result.txt"))
	for _, want := range []string{"exit_code: -1", "timed_out: true", lastEvent} {
		if !strings.Contains(process, want) {
			t.Fatalf("process result missing %q:\n%s", want, process)
		}
	}
	if got := strings.Fields(readFile(t, invocations)); len(got) != 1 {
		t.Fatalf("codex invocation records = %v, want exactly one", got)
	}
}

func TestCodexProcessRecognizesTerminalTurnBeforeOSExit(t *testing.T) {
	const quietBudget = 250 * time.Millisecond
	spec, invocations := codexProcessFixture(t, "terminal-before-exit", "terminal-then-hang", quietBudget)
	spec.env = append(spec.env, "GO_WANT_CODEX_FINAL_MESSAGE_PATH="+spec.finalPath)

	result, runErr := runCodexProcess(spec)
	if runErr != nil {
		t.Fatalf("terminal Codex process should finish after turn.completed and final message: %v", runErr)
	}
	if result.timedOut {
		t.Fatal("terminal Codex process was classified as a quiet timeout")
	}
	if !strings.Contains(result.jsonl, `{"type":"turn.completed"`) {
		t.Fatalf("terminal Codex event was not preserved in the transcript: %q", result.jsonl)
	}
	if result.finalMessage != "Gate review: terminal evidence\n" {
		t.Fatalf("final message = %q, want the complete terminal message", result.finalMessage)
	}
	if got := strings.Fields(readFile(t, invocations)); len(got) != 1 {
		t.Fatalf("codex invocation records = %v, want exactly one", got)
	}
	if !result.terminal {
		t.Fatal("terminal Codex process did not record the semantic terminal boundary")
	}
}

func TestCodexProcessRequiresFinalMessageForTerminalTurn(t *testing.T) {
	const quietBudget = 250 * time.Millisecond
	spec, _ := codexProcessFixture(t, "terminal-without-final-message", "terminal-without-final", quietBudget)
	result, runErr := runCodexProcess(spec)
	if runErr == nil {
		t.Fatal("turn.completed without a final message unexpectedly passed")
	}
	if result.timedOut {
		t.Fatalf("terminal process without a final message should fail on missing terminal evidence, not quiet timeout: %+v", result)
	}
	if result.terminal {
		t.Fatal("turn.completed without a final message was incorrectly accepted as terminal")
	}
	if !strings.Contains(result.jsonl, `{"type":"turn.completed"`) {
		t.Fatalf("terminal event was not preserved for the failed case: %q", result.jsonl)
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
	fmt.Println(`{"type":"item.started","sequence":1,"item":{"type":"collab_tool_call","tool":"wait_agent"}}`)
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
	if mode == "progress-then-exit" {
		for sequence := 2; sequence <= 31; sequence++ {
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("{\"type\":\"item.started\",\"sequence\":%d,\"item\":{\"type\":\"collab_tool_call\",\"tool\":\"wait_agent\"}}\n", sequence)
		}
		os.Exit(0)
	}
	if mode == "terminal-then-hang" {
		finalPath := os.Getenv("GO_WANT_CODEX_FINAL_MESSAGE_PATH")
		if finalPath == "" {
			fmt.Fprintln(os.Stderr, "terminal fixture missing final-message path")
			os.Exit(2)
		}
		if err := os.WriteFile(finalPath, []byte("Gate review: terminal evidence\n"), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		fmt.Println(`{"type":"turn.completed","usage":{"input_tokens":1,"output_tokens":1}}`)
	}
	if mode == "terminal-without-final" {
		fmt.Println(`{"type":"turn.completed","usage":{"input_tokens":1,"output_tokens":1}}`)
	}
	for {
		time.Sleep(10 * time.Millisecond)
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
