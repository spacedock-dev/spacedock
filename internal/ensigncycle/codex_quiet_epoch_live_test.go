//go:build live

package ensigncycle

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	codexQuietEpochScenarioName  = "foreground-wait-quiet-epoch"
	codexQuietEpochEntity        = "foreground-wait-quiet-epoch"
	codexQuietEpochHoldDuration  = 6 * time.Minute
	codexQuietEpochSilenceBudget = 8 * time.Minute
)

func TestLiveCodexForegroundWaitRetriesQuietlyWithinEpoch(t *testing.T) {
	result, err := runCodexQuietEpochScenario(t)
	if err != nil {
		t.Fatalf("quiet-epoch scenario: %v\nArtifacts: %s", err, result.artifactDir)
	}
	if err := assertCodexQuietEpochRollout(result.rollout, result.holdStarted, result.holdFinished); err != nil {
		t.Fatalf("quiet-epoch rollout: %v\nArtifacts: %s", err, result.artifactDir)
	}
	if err := assertCodexQuietEpochTrace(result.trace); err != nil {
		t.Fatalf("quiet-epoch narration: %v\nArtifacts: %s", err, result.artifactDir)
	}
}

type codexQuietEpochResult struct {
	codexForegroundWaitResult
	holdStarted  time.Time
	holdFinished time.Time
}

func runCodexQuietEpochScenario(t *testing.T) (codexQuietEpochResult, error) {
	t.Helper()
	runner := newCodexLiveRunner(t)
	workflowRoot := t.TempDir()
	stateRoot, entityPath := writeCodexQuietEpochWorkflow(t, workflowRoot)

	foregroundResult, err := runner.runForegroundWaitScenario(t, workflowRoot, codexQuietEpochPrompt(workflowRoot), codexQuietEpochScenarioName, codexQuietEpochSilenceBudget)
	result := codexQuietEpochResult{codexForegroundWaitResult: foregroundResult}
	if err != nil {
		return result, err
	}

	entity := readFile(t, entityPath)
	holdStarted, holdFinished, err := codexForegroundWaitHoldTimes(entity)
	if err != nil {
		return result, fmt.Errorf("read quiet-epoch hold markers: %w", err)
	}
	if holdFinished.Sub(holdStarted) < codexQuietEpochHoldDuration {
		return result, fmt.Errorf("quiet-epoch hold = %s, want at least %s", holdFinished.Sub(holdStarted), codexQuietEpochHoldDuration)
	}
	result.holdStarted = holdStarted
	result.holdFinished = holdFinished
	entityRelativePath := filepath.Join(codexQuietEpochEntity, "index.md")
	if !codexForegroundWaitReportCommitted(t, stateRoot, entityRelativePath) {
		return result, fmt.Errorf("quiet-epoch report was not committed path-scoped to %s", entityRelativePath)
	}
	if status := strings.TrimSpace(git(t, stateRoot, "status", "--short", "--", entityRelativePath)); status != "" {
		return result, fmt.Errorf("quiet-epoch worker entity has uncommitted changes: %s", status)
	}
	return result, nil
}

func writeCodexQuietEpochWorkflow(t *testing.T, root string) (stateRoot, entityPath string) {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), codexQuietEpochReadme())
	stateRoot = filepath.Join(root, ".spacedock-state")
	entityPath = filepath.Join(stateRoot, codexQuietEpochEntity, "index.md")
	writeFile(t, entityPath, "---\n"+
		"id: "+codexQuietEpochEntity+"\n"+
		"title: Codex Foreground Wait Quiet Epoch\n"+
		"status: backlog\n"+
		"completed:\n"+
		"verdict:\n"+
		"worktree:\n"+
		"---\n"+
		"# Codex Foreground Wait Quiet Epoch\n")
	gitInit(t, root)
	gitInit(t, stateRoot)
	return stateRoot, entityPath
}

func codexQuietEpochReadme() string {
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
		"# Codex Foreground Wait Quiet Epoch Fixture\n\n" +
		"### implementation\n\n" +
		"Before writing any repository file, run exactly `started=$(date +%s); sleep 360; finished=$(date +%s); printf '%s %s\\n' \"$started\" \"$finished\"`. Do not shorten or background the delay. Then append an implementation stage report to `.spacedock-state/" + codexQuietEpochEntity + "/index.md` containing exactly these machine-readable lines, with the captured values substituted:\n\n" +
		"```text\n" +
		"CODEX_FOREGROUND_WAIT_STARTED_EPOCH=<started>\n" +
		"CODEX_FOREGROUND_WAIT_FINISHED_EPOCH=<finished>\n" +
		"CODEX_FOREGROUND_WAIT_PROBE_COMPLETE\n" +
		"```\n\n" +
		"The report must include a `- DONE:` item and `### Summary`. Commit only `" + codexQuietEpochEntity + "/index.md` in the state checkout before reporting completion.\n\n" +
		"### done\n\nTerminal state.\n"
}

func codexQuietEpochPrompt(workflowRoot string) string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: "+workflowRoot,
		"The split-root state checkout is `"+filepath.Join(workflowRoot, ".spacedock-state")+"`, and the entity is `"+filepath.Join(workflowRoot, ".spacedock-state", codexQuietEpochEntity, "index.md")+"`.",
		"Process only the entity `"+codexQuietEpochEntity+"`. Dispatch its sole implementation worker exactly once; do not do the worker's timed hold yourself and do not dispatch additional workers.",
		"After the worker completes, verify both its implementation stage report and its path-scoped git commit before presenting the implementation gate. Do not approve, advance, archive, or redispatch the entity.",
		"This is a foreground-wait quiet-epoch control: follow the Codex first-officer runtime adapter while the owned worker is unresolved, then stop after the gate presentation.",
	)
}

func assertCodexQuietEpochRollout(rollout codexForegroundWaitRollout, holdStarted, holdFinished time.Time) error {
	if rollout.spawnCalls != 1 {
		return fmt.Errorf("isolated Codex rollout spawned %d workers, want one quiet-epoch worker (%s)", rollout.spawnCalls, rollout.artifactPath)
	}
	if len(rollout.waits) < 2 {
		return fmt.Errorf("isolated Codex rollout has %d waits, want a timeout and reinstalled wait (%s)", len(rollout.waits), rollout.artifactPath)
	}
	first, second := rollout.waits[0], rollout.waits[1]
	if first.timeoutMS != 300000 || second.timeoutMS != 300000 {
		return fmt.Errorf("quiet-epoch waits used %d then %d ms, want explicit 300000 ms", first.timeoutMS, second.timeoutMS)
	}
	if first.returnedAt.IsZero() || first.timedOut == nil || !*first.timedOut {
		return fmt.Errorf("first quiet-epoch wait must return timed_out:true")
	}
	if second.calledAt.IsZero() || !second.calledAt.After(first.returnedAt) {
		return fmt.Errorf("quiet-epoch re-wait did not begin after the ordinary timeout return")
	}
	return assertCodexQuietEpochTiming(rollout, holdStarted, holdFinished)
}

func TestCodexQuietEpochTimingAcceptsTimeoutAndRewaitDuringHold(t *testing.T) {
	holdStarted := time.Date(2026, time.July, 12, 3, 0, 0, 0, time.UTC)
	holdFinished := holdStarted.Add(codexQuietEpochHoldDuration)
	rollout := quietEpochTimingRollout(
		holdStarted.Add(-5*time.Second),
		holdFinished.Add(-2*time.Second),
		holdFinished.Add(-time.Second),
	)

	if err := assertCodexQuietEpochRollout(rollout, holdStarted, holdFinished); err != nil {
		t.Fatalf("timeout and re-wait inside the unresolved hold must pass: %v", err)
	}
}

func TestCodexQuietEpochTimingRejectsBoundaryRegressions(t *testing.T) {
	holdStarted := time.Date(2026, time.July, 12, 3, 0, 0, 0, time.UTC)
	holdFinished := holdStarted.Add(codexQuietEpochHoldDuration)
	tests := []struct {
		name    string
		rollout codexForegroundWaitRollout
	}{
		{
			name: "first wait misses hold",
			rollout: quietEpochTimingRollout(
				holdFinished.Add(time.Second),
				holdFinished.Add(5*time.Minute),
				holdFinished.Add(5*time.Minute+time.Second),
			),
		},
		{
			name: "timeout reaches hold finish",
			rollout: quietEpochTimingRollout(
				holdStarted,
				holdFinished,
				holdFinished.Add(time.Second),
			),
		},
		{
			name: "rewait begins after hold",
			rollout: quietEpochTimingRollout(
				holdStarted,
				holdFinished.Add(-time.Second),
				holdFinished,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := assertCodexQuietEpochRollout(tt.rollout, holdStarted, holdFinished); err == nil {
				t.Fatal("late or non-overlapping wait epoch must fail")
			}
		})
	}
}

func quietEpochTimingRollout(firstStarted, firstReturned, rewaitStarted time.Time) codexForegroundWaitRollout {
	timedOut := true
	return codexForegroundWaitRollout{
		spawnCalls: 1,
		waits: []codexForegroundWaitCall{
			{
				calledAt:   firstStarted,
				returnedAt: firstReturned,
				timeoutMS:  300000,
				timedOut:   &timedOut,
			},
			{
				calledAt:  rewaitStarted,
				timeoutMS: 300000,
			},
		},
	}
}

// assertCodexQuietEpochTiming binds the timeout/re-wait trace to the period in
// which the worker is still unresolved. Epoch silence is meaningful only while
// the first wait overlaps the worker hold and both the timeout and re-wait occur
// before that hold ends.
func assertCodexQuietEpochTiming(rollout codexForegroundWaitRollout, holdStarted, holdFinished time.Time) error {
	if holdStarted.IsZero() || holdFinished.IsZero() || !holdFinished.After(holdStarted) {
		return fmt.Errorf("quiet-epoch hold markers are invalid: %s through %s", holdStarted, holdFinished)
	}
	if len(rollout.waits) < 2 {
		return fmt.Errorf("quiet-epoch timing needs a timeout and reinstalled wait")
	}
	first, second := rollout.waits[0], rollout.waits[1]
	if first.calledAt.IsZero() || first.returnedAt.IsZero() || second.calledAt.IsZero() {
		return fmt.Errorf("quiet-epoch timing is missing a wait timestamp")
	}
	if !first.calledAt.Before(holdFinished) || !first.returnedAt.After(holdStarted) {
		return fmt.Errorf("first wait does not overlap the unresolved hold")
	}
	if !first.returnedAt.Before(holdFinished) {
		return fmt.Errorf("ordinary timeout returned at or after the hold finished")
	}
	if !second.calledAt.Before(holdFinished) {
		return fmt.Errorf("silent re-wait began at or after the hold finished")
	}
	return nil
}

func TestCodexQuietEpochTraceRejectsRepeatedNarration(t *testing.T) {
	trace := []timedCodexLine{
		{line: `{"type":"item.completed","item":{"type":"collab_tool_call","tool":"wait_agent","status":"completed"}}`},
		{line: `{"type":"item.completed","item":{"type":"agent_message","text":"continuing with the same unresolved worker"}}`},
		{line: `{"type":"item.started","item":{"type":"collab_tool_call","tool":"wait_agent"}}`},
	}

	if err := assertCodexQuietEpochTrace(trace); err == nil {
		t.Fatal("repeated narration between an ordinary timeout and its re-wait must fail")
	}
}

func TestCodexQuietEpochTraceAcceptsSilentRewait(t *testing.T) {
	trace := []timedCodexLine{
		{line: `{"type":"item.completed","item":{"type":"collab_tool_call","tool":"wait_agent","status":"completed"}}`},
		{line: `{"type":"item.started","item":{"type":"collab_tool_call","tool":"wait_agent"}}`},
	}

	if err := assertCodexQuietEpochTrace(trace); err != nil {
		t.Fatalf("silent re-wait must remain in the monitoring epoch: %v", err)
	}
}

// assertCodexQuietEpochTrace inspects the live Codex event stream, rather than
// the adapter prose. Its caller separately proves that the first completed wait
// is an ordinary timeout; this helper then requires the next wait to begin with
// no intervening FO narration.
func assertCodexQuietEpochTrace(trace []timedCodexLine) error {
	seenWaitReturn := false
	for _, entry := range trace {
		var event struct {
			Type string `json:"type"`
			Item struct {
				Type string `json:"type"`
				Tool string `json:"tool"`
				Text string `json:"text"`
			} `json:"item"`
		}
		if err := json.Unmarshal([]byte(entry.line), &event); err != nil {
			continue
		}
		isWait := event.Item.Type == "collab_tool_call" && isCodexWaitTool(strings.ToLower(event.Item.Tool))
		if !seenWaitReturn && isWait && event.Type == "item.completed" {
			seenWaitReturn = true
			continue
		}
		if seenWaitReturn && isWait && event.Type == "item.started" {
			return nil
		}
		if seenWaitReturn && event.Item.Type == "agent_message" && strings.TrimSpace(event.Item.Text) != "" {
			return fmt.Errorf("FO narrated between an ordinary wait return and its re-wait")
		}
	}
	if !seenWaitReturn {
		return fmt.Errorf("trace has no completed foreground wait")
	}
	return fmt.Errorf("trace has no reinstalled foreground wait")
}
