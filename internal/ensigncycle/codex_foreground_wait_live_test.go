//go:build live

package ensigncycle

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	codexForegroundWaitScenarioName = "foreground-wait-timeout"
	codexForegroundWaitEntity       = "foreground-wait-timeout"
	codexForegroundWaitReportMarker = "CODEX_FOREGROUND_WAIT_PROBE_COMPLETE"

	codexForegroundWaitHoldDuration = 45 * time.Second
	codexForegroundWaitCrossingMin  = 35 * time.Second
)

var (
	codexForegroundWaitHoldStartedEpoch  = regexp.MustCompile(`(?m)^CODEX_FOREGROUND_WAIT_STARTED_EPOCH=(\d+)$`)
	codexForegroundWaitHoldFinishedEpoch = regexp.MustCompile(`(?m)^CODEX_FOREGROUND_WAIT_FINISHED_EPOCH=(\d+)$`)
	codexForegroundWaitTimedOut          = regexp.MustCompile(`(?i)"timed_out"\s*:\s*true`)
)

// TestLiveCodexForegroundWaitUsesFiveMinutePerCallPolicy is the runtime proof for
// the Codex adapter's explicit foreground wait. The worker does no repository write
// before a 45-second hold, then commits its implementation stage report. Its recorded
// timestamps mark the delayed hold only, not report, final-status, or completion time.
// The test records the live collaboration stream and requires one non-timeout wait
// across the old default horizon, no pre-return re-wait churn, and separate durable
// report/path-scoped-commit verification.
func TestLiveCodexForegroundWaitUsesFiveMinutePerCallPolicy(t *testing.T) {
	runner := newCodexLiveRunner(t)
	workflowRoot := t.TempDir()
	stateRoot, entityPath := writeCodexForegroundWaitWorkflow(t, workflowRoot)

	result, err := runner.runForegroundWaitScenario(t, workflowRoot, codexForegroundWaitPrompt(workflowRoot))
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}

	entity := readFile(t, entityPath)
	holdStarted, holdFinished, err := codexForegroundWaitHoldTimes(entity)
	if err != nil {
		t.Fatalf("read delayed-hold timing report: %v\nEntity:\n%s\nArtifacts: %s", err, entity, result.artifactDir)
	}
	if holdFinished.Sub(holdStarted) < codexForegroundWaitHoldDuration {
		t.Fatalf("delayed hold = %s, want at least %s\nEntity:\n%s\nArtifacts: %s", holdFinished.Sub(holdStarted), codexForegroundWaitHoldDuration, entity, result.artifactDir)
	}
	entityRelativePath := filepath.Join(codexForegroundWaitEntity, "index.md")
	if !codexForegroundWaitReportCommitted(t, stateRoot, entityRelativePath) {
		t.Fatalf("delayed-hold timing report was not committed path-scoped to state checkout %s\nArtifacts: %s", entityRelativePath, result.artifactDir)
	}
	if status := strings.TrimSpace(git(t, stateRoot, "status", "--short", "--", entityRelativePath)); status != "" {
		t.Fatalf("state-checkout worker entity has uncommitted changes after completion: %s\nArtifacts: %s", status, result.artifactDir)
	}

	if err := assertCodexForegroundWaitTrace(result.trace, result.rollout, holdFinished); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	if err := writeCodexForegroundWaitEvidence(t, result, holdStarted, holdFinished, stateRoot, entityRelativePath); err != nil {
		t.Fatalf("write foreground-wait evidence: %v", err)
	}

	t.Logf("Codex foreground wait evidence: hold_duration=%s artifacts=%s", holdFinished.Sub(holdStarted), result.artifactDir)
}

type timedCodexLine struct {
	at   time.Time
	line string
}

type codexForegroundWaitResult struct {
	codexScenarioResult
	trace   []timedCodexLine
	rollout codexForegroundWaitRollout
}

type codexForegroundWaitRollout struct {
	sourcePath   string
	artifactPath string
	spawnCalls   int
	waits        []codexForegroundWaitCall
}

type codexForegroundWaitCall struct {
	callID     string
	calledAt   time.Time
	returnedAt time.Time
	timeoutMS  int
	timedOut   *bool
}

func captureCodexForegroundWaitRollout(codexHome, artifactDir string) (codexForegroundWaitRollout, error) {
	var selected codexForegroundWaitRollout
	var selectedData []byte
	sessionRoot := filepath.Join(codexHome, "sessions")
	err := filepath.WalkDir(sessionRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".jsonl" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		candidate := parseCodexForegroundWaitRollout(path, data)
		if candidate.spawnCalls == 0 || len(candidate.waits) == 0 {
			return nil
		}
		if selected.sourcePath == "" || candidate.waits[0].calledAt.After(selected.waits[0].calledAt) {
			selected = candidate
			selectedData = data
		}
		return nil
	})
	if err != nil {
		return codexForegroundWaitRollout{}, fmt.Errorf("scan isolated Codex rollouts under %s: %w", sessionRoot, err)
	}
	if selected.sourcePath == "" {
		return codexForegroundWaitRollout{}, fmt.Errorf("no isolated Codex rollout contained both spawn_agent and wait_agent")
	}
	selected.artifactPath = filepath.Join(artifactDir, "codex-foreground-wait-rollout.jsonl")
	if err := os.WriteFile(selected.artifactPath, selectedData, 0o644); err != nil {
		return codexForegroundWaitRollout{}, fmt.Errorf("copy selected Codex rollout: %w", err)
	}
	return selected, nil
}

func parseCodexForegroundWaitRollout(path string, data []byte) codexForegroundWaitRollout {
	rollout := codexForegroundWaitRollout{sourcePath: path}
	waitByCallID := map[string]int{}
	for _, line := range strings.Split(string(data), "\n") {
		var record struct {
			Timestamp string `json:"timestamp"`
			Type      string `json:"type"`
			Payload   struct {
				Type      string          `json:"type"`
				CallID    string          `json:"call_id"`
				Name      string          `json:"name"`
				Arguments string          `json:"arguments"`
				Output    json.RawMessage `json:"output"`
			} `json:"payload"`
		}
		if err := json.Unmarshal([]byte(line), &record); err != nil || record.Type != "response_item" {
			continue
		}
		at, err := time.Parse(time.RFC3339Nano, record.Timestamp)
		if err != nil {
			continue
		}
		switch record.Payload.Type {
		case "function_call":
			switch record.Payload.Name {
			case "spawn_agent":
				rollout.spawnCalls++
			case "wait_agent":
				var arguments struct {
					TimeoutMS int `json:"timeout_ms"`
				}
				if err := json.Unmarshal([]byte(record.Payload.Arguments), &arguments); err != nil {
					continue
				}
				rollout.waits = append(rollout.waits, codexForegroundWaitCall{
					callID:    record.Payload.CallID,
					calledAt:  at,
					timeoutMS: arguments.TimeoutMS,
				})
				waitByCallID[record.Payload.CallID] = len(rollout.waits) - 1
			}
		case "function_call_output":
			waitIndex, ok := waitByCallID[record.Payload.CallID]
			if !ok {
				continue
			}
			rollout.waits[waitIndex].returnedAt = at
			if timedOut, found := codexRolloutTimedOut(record.Payload.Output); found {
				rollout.waits[waitIndex].timedOut = &timedOut
			}
		}
	}
	return rollout
}

func codexRolloutTimedOut(raw json.RawMessage) (bool, bool) {
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return false, false
	}
	return codexRolloutTimedOutValue(value)
}

func codexRolloutTimedOutValue(value any) (bool, bool) {
	switch value := value.(type) {
	case map[string]any:
		for key, child := range value {
			if strings.EqualFold(key, "timed_out") {
				if timedOut, ok := child.(bool); ok {
					return timedOut, true
				}
			}
			if timedOut, found := codexRolloutTimedOutValue(child); found {
				return timedOut, true
			}
		}
	case []any:
		for _, child := range value {
			if timedOut, found := codexRolloutTimedOutValue(child); found {
				return timedOut, true
			}
		}
	case string:
		var nested any
		if json.Unmarshal([]byte(value), &nested) == nil {
			if timedOut, found := codexRolloutTimedOutValue(nested); found {
				return timedOut, true
			}
		}
		match := regexp.MustCompile(`(?i)"?timed_out"?\s*:\s*(true|false)`).FindStringSubmatch(value)
		if len(match) == 2 {
			return strings.EqualFold(match[1], "true"), true
		}
	}
	return false, false
}

func (r codexLiveRunner) runForegroundWaitScenario(t *testing.T, workflowRoot, prompt string) (codexForegroundWaitResult, error) {
	t.Helper()
	artifactDir := codexAttemptArtifactDir(r.artifactRoot, codexForegroundWaitScenarioName, 0)
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	finalPath := filepath.Join(artifactDir, "codex-final-message.txt")
	jsonlPath := filepath.Join(artifactDir, "codex-exec.jsonl")
	stderrPath := filepath.Join(artifactDir, "codex-exec.stderr.txt")

	cmd := exec.Command(r.codexBin, codexExecArgv(workflowRoot, finalPath, prompt)...)
	cmd.Env = r.env
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	stderr, err := os.Create(stderrPath)
	if err != nil {
		t.Fatal(err)
	}
	defer stderr.Close()
	cmd.Stderr = stderr

	started := time.Now()
	if err := cmd.Start(); err != nil {
		return codexForegroundWaitResult{codexScenarioResult: codexScenarioResult{artifactDir: artifactDir}}, fmt.Errorf("codex exec failed to start foreground-wait scenario: %w", err)
	}
	poller := newCmdPoller(cmd, pw)
	defer poller.kill()
	trace := make([]timedCodexLine, 0, 64)
	watcher := newStreamWatcher(newPipeLineSource(pr), poller, func(line string) {
		trace = append(trace, timedCodexLine{at: time.Now(), line: line})
	})
	probe, err := newWorkflowStateProbe(workflowRoot)
	if err != nil {
		return codexForegroundWaitResult{codexScenarioResult: codexScenarioResult{artifactDir: artifactDir}}, fmt.Errorf("snapshot workflow state for foreground-wait scenario: %w", err)
	}
	watchdog := newCodexCollabWaitWatchdog(codexForegroundWaitScenarioName, artifactDir, probe)
	// The parent stream is intentionally quiet while the worker boots its assigned
	// skill and performs the 45-second hold. Use the existing dispatch-close
	// no-progress budget so that setup overhead cannot be misclassified as a Codex
	// wait timeout; this is still a stream-silence watchdog, not a per-call timeout.
	jsonl, runErr := drainCodexToExitWithWaitWatchdog(watcher, quietBudgetDispatchClose, "codex foreground-wait scenario", watchdog)
	result := codexForegroundWaitResult{
		codexScenarioResult: codexScenarioResult{
			jsonl:       jsonl,
			artifactDir: artifactDir,
			duration:    time.Since(started),
		},
		trace: trace,
	}
	if err := os.WriteFile(jsonlPath, []byte(jsonl), 0o644); err != nil {
		t.Fatal(err)
	}
	codexHome, ok := envValue(r.env, "CODEX_HOME")
	if !ok || codexHome == "" {
		return result, fmt.Errorf("isolated Codex runner did not expose CODEX_HOME")
	}
	rollout, rolloutErr := captureCodexForegroundWaitRollout(codexHome, artifactDir)
	result.rollout = rollout
	if runErr != nil {
		return result, runErr
	}
	if rolloutErr != nil {
		return result, rolloutErr
	}
	result.finalMessage = readFile(t, finalPath)
	return result, nil
}

func writeCodexForegroundWaitWorkflow(t *testing.T, root string) (stateRoot, entityPath string) {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), codexForegroundWaitReadme())
	stateRoot = filepath.Join(root, ".spacedock-state")
	entityPath = filepath.Join(stateRoot, codexForegroundWaitEntity, "index.md")
	writeFile(t, entityPath, codexForegroundWaitEntityBody())
	gitInit(t, root)
	gitInit(t, stateRoot)
	return stateRoot, entityPath
}

func codexForegroundWaitReadme() string {
	return "---\n" +
		"commissioned-by: spacedock@1\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"state: .spacedock-state\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: backlog\n" +
		"      initial: true\n" +
		"    - name: implementation\n" +
		"      gate: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Codex Foreground Wait Fixture\n\n" +
		"This fixture has one owned implementation worker. Its delay is an intentional live timing control, not background work.\n\n" +
		"### implementation\n\n" +
		"Before writing any repository file, run exactly `started=$(date +%s); sleep 45; finished=$(date +%s); printf '%s %s\\n' \"$started\" \"$finished\"`. Do not shorten or background the delay. These are delayed-hold endpoints, not report or final-status timestamps. Then append an implementation stage report to `.spacedock-state/foreground-wait-timeout/index.md` containing exactly these machine-readable lines, with the captured values substituted:\n\n" +
		"```text\n" +
		"CODEX_FOREGROUND_WAIT_STARTED_EPOCH=<started>\n" +
		"CODEX_FOREGROUND_WAIT_FINISHED_EPOCH=<finished>\n" +
		"CODEX_FOREGROUND_WAIT_PROBE_COMPLETE\n" +
		"```\n\n" +
		"The report must include a `- DONE:` item and `### Summary`. Commit only `foreground-wait-timeout/index.md` in the state checkout before reporting completion.\n\n" +
		"- **Outputs:** A committed implementation stage report carrying the exact timing markers after the 45-second hold.\n\n" +
		"### done\n\nTerminal state.\n"
}

func codexForegroundWaitEntityBody() string {
	return "---\n" +
		"id: " + codexForegroundWaitEntity + "\n" +
		"title: Codex Foreground Wait Timeout\n" +
		"status: backlog\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Codex Foreground Wait Timeout\n\n" +
		"This entity exists only to prove a delayed owned worker wakes the first officer's foreground wait without timeout churn.\n"
}

func codexForegroundWaitPrompt(workflowRoot string) string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: "+workflowRoot,
		"The split-root state checkout is `"+filepath.Join(workflowRoot, ".spacedock-state")+"`, and the entity is `"+filepath.Join(workflowRoot, ".spacedock-state", codexForegroundWaitEntity, "index.md")+"`.",
		"Process only the entity `"+codexForegroundWaitEntity+"`. Dispatch its sole implementation worker exactly once; do not do the worker's timed hold yourself and do not dispatch additional workers.",
		"After the worker completes, verify both its implementation stage report and its path-scoped git commit before presenting the implementation gate. Do not approve, advance, archive, or redispatch the entity.",
		"This is a foreground-wait timing control: follow the Codex first-officer runtime adapter while the owned worker is unresolved, then stop after the gate presentation.",
	)
}

func codexForegroundWaitHoldTimes(entity string) (time.Time, time.Time, error) {
	start := codexForegroundWaitHoldStartedEpoch.FindStringSubmatch(entity)
	finish := codexForegroundWaitHoldFinishedEpoch.FindStringSubmatch(entity)
	if len(start) != 2 || len(finish) != 2 || !strings.Contains(entity, codexForegroundWaitReportMarker) {
		return time.Time{}, time.Time{}, fmt.Errorf("stage report is missing required delayed-hold timing markers")
	}
	startSeconds, err := strconv.ParseInt(start[1], 10, 64)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse delayed-hold start epoch %q: %w", start[1], err)
	}
	finishSeconds, err := strconv.ParseInt(finish[1], 10, 64)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse delayed-hold finish epoch %q: %w", finish[1], err)
	}
	holdStarted := time.Unix(startSeconds, 0)
	holdFinished := time.Unix(finishSeconds, 0)
	if !holdFinished.After(holdStarted) {
		return time.Time{}, time.Time{}, fmt.Errorf("delayed-hold finish %s is not after start %s", holdFinished, holdStarted)
	}
	return holdStarted, holdFinished, nil
}

func codexForegroundWaitReportCommitted(t *testing.T, root, entityName string) bool {
	t.Helper()
	for _, sha := range strings.Fields(git(t, root, "log", "--format=%H", "--", entityName)) {
		body := git(t, root, "show", sha+":"+entityName)
		if !strings.Contains(body, codexForegroundWaitReportMarker) {
			continue
		}
		names := strings.Fields(git(t, root, "show", "--format=", "--name-only", sha))
		if len(names) == 1 && names[0] == entityName {
			return true
		}
	}
	return false
}

func assertCodexForegroundWaitTrace(trace []timedCodexLine, rollout codexForegroundWaitRollout, holdFinished time.Time) error {
	if rollout.spawnCalls != 1 {
		return fmt.Errorf("isolated Codex rollout spawned %d workers, want exactly the one owned timing worker (%s)", rollout.spawnCalls, rollout.artifactPath)
	}
	if len(rollout.waits) == 0 {
		return fmt.Errorf("isolated Codex rollout contains no wait_agent call (%s)", rollout.artifactPath)
	}
	wait := rollout.waits[0]
	for _, observedWait := range rollout.waits {
		if observedWait.timeoutMS != 300000 {
			return fmt.Errorf("wait_agent timeout_ms = %d, want explicit per-call 300000 (%s)", observedWait.timeoutMS, rollout.artifactPath)
		}
		if observedWait.returnedAt.IsZero() {
			return fmt.Errorf("wait_agent call %s has no function_call_output timestamp (%s)", observedWait.callID, rollout.artifactPath)
		}
		if observedWait.timedOut == nil {
			return fmt.Errorf("wait_agent call %s has no timed_out result (%s)", observedWait.callID, rollout.artifactPath)
		}
		if *observedWait.timedOut {
			return fmt.Errorf("wait_agent call %s timed out (%s)", observedWait.callID, rollout.artifactPath)
		}
	}
	for _, laterWait := range rollout.waits[1:] {
		if !laterWait.calledAt.After(wait.returnedAt) {
			return fmt.Errorf("wait_agent call %s began before the owned worker's foreground wait returned, showing timeout/re-wait churn (%s)", laterWait.callID, rollout.artifactPath)
		}
	}
	if wait.returnedAt.Sub(wait.calledAt) < codexForegroundWaitCrossingMin {
		return fmt.Errorf("wait_agent lasted %s, want at least %s to cross the old short default (%s)", wait.returnedAt.Sub(wait.calledAt), codexForegroundWaitCrossingMin, rollout.artifactPath)
	}
	if wait.returnedAt.Before(holdFinished) {
		return fmt.Errorf("wait_agent returned at %s before the delayed-hold end marker at %s", wait.returnedAt.Format(time.RFC3339Nano), holdFinished.Format(time.RFC3339Nano))
	}
	var waitStarts, waitReturns, completions, timedOut []timedCodexLine
	for _, entry := range trace {
		if codexForegroundWaitTimedOut.MatchString(entry.line) {
			timedOut = append(timedOut, entry)
		}
		var event struct {
			Type string `json:"type"`
			Item struct {
				Type string `json:"type"`
				Tool string `json:"tool"`
			} `json:"item"`
		}
		if err := json.Unmarshal([]byte(entry.line), &event); err != nil {
			continue
		}
		if event.Item.Type == "collab_tool_call" && isCodexWaitTool(strings.ToLower(event.Item.Tool)) {
			switch event.Type {
			case "item.started":
				waitStarts = append(waitStarts, entry)
			case "item.completed":
				waitReturns = append(waitReturns, entry)
			}
		}
		if completed, _ := lineIsCodexWorkerCompletion(entry.line); completed {
			completions = append(completions, entry)
		}
	}
	if len(waitStarts) == 0 && len(waitReturns) == 0 {
		// Codex's --json surface may omit collaboration metadata; the timestamped
		// rollout above remains the definitive function-call record in that case.
		return nil
	}
	if len(waitStarts) == 0 || len(waitReturns) == 0 {
		return fmt.Errorf("collaboration stream has incomplete foreground-wait evidence (%d starts, %d returns)\n%s", len(waitStarts), len(waitReturns), codexForegroundWaitTraceTail(trace))
	}
	firstReturn := waitReturns[0]
	startsBeforeFirstReturn := 0
	for _, started := range waitStarts {
		if !started.at.After(firstReturn.at) {
			startsBeforeFirstReturn++
		}
	}
	if startsBeforeFirstReturn != 1 {
		return fmt.Errorf("collaboration stream has %d wait starts before the owned worker's first wait returned, want 1\n%s", startsBeforeFirstReturn, codexForegroundWaitTraceTail(trace))
	}
	if len(completions) == 0 {
		// Current multi-agent traces can return an empty agents_states object on a
		// non-timeout wait. The paired isolated rollout's timed_out:false output is
		// the authoritative non-timeout classification for this headless run.
		return nil
	}
	completion := completions[0]
	waitsBeforeCompletion := 0
	for _, started := range waitStarts {
		if !started.at.After(completion.at) {
			waitsBeforeCompletion++
		}
	}
	if waitsBeforeCompletion != 1 {
		return fmt.Errorf("foreground wait starts before owned-worker completion = %d, want 1 (a timeout/re-wait churned)\n%s", waitsBeforeCompletion, codexForegroundWaitTraceTail(trace))
	}
	for _, timeout := range timedOut {
		if !timeout.at.After(completion.at) {
			return fmt.Errorf("collaboration stream reports a timeout before the owned worker completed\n%s", codexForegroundWaitTraceTail(trace))
		}
	}
	if len(waitReturns) == 0 {
		return fmt.Errorf("collaboration stream contains no completed foreground-wait event\n%s", codexForegroundWaitTraceTail(trace))
	}
	return nil
}

func codexForegroundWaitTraceTail(trace []timedCodexLine) string {
	start := len(trace) - transcriptTailLines
	if start < 0 {
		start = 0
	}
	var out strings.Builder
	for _, entry := range trace[start:] {
		fmt.Fprintf(&out, "%s\t%s\n", entry.at.Format(time.RFC3339Nano), entry.line)
	}
	return out.String()
}

func writeCodexForegroundWaitEvidence(t *testing.T, result codexForegroundWaitResult, holdStarted, holdFinished time.Time, stateRoot, entityName string) error {
	t.Helper()
	var trace strings.Builder
	for _, entry := range result.trace {
		fmt.Fprintf(&trace, "%s\t%s\n", entry.at.Format(time.RFC3339Nano), entry.line)
	}
	if err := os.WriteFile(filepath.Join(result.artifactDir, "codex-foreground-wait-timed-trace.txt"), []byte(trace.String()), 0o644); err != nil {
		return err
	}
	gitLog := git(t, stateRoot, "log", "--oneline", "--", entityName)
	wait := result.rollout.waits[0]
	return os.WriteFile(filepath.Join(result.artifactDir, "codex-foreground-wait-evidence.txt"), []byte(fmt.Sprintf("hold_started=%s\nhold_finished=%s\nhold_duration=%s\nprocess_duration=%s\nwait_call=%s\nwait_started=%s\nwait_returned=%s\nwait_duration=%s\ntimeout_ms=%d\ntimed_out=%t\nrollout_source=%s\nrollout_artifact=%s\n\nstate checkout git log --oneline -- %s\n%s", holdStarted.Format(time.RFC3339Nano), holdFinished.Format(time.RFC3339Nano), holdFinished.Sub(holdStarted), result.duration, wait.callID, wait.calledAt.Format(time.RFC3339Nano), wait.returnedAt.Format(time.RFC3339Nano), wait.returnedAt.Sub(wait.calledAt), wait.timeoutMS, *wait.timedOut, result.rollout.sourcePath, result.rollout.artifactPath, entityName, gitLog)), 0o644)
}
