package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	statuspkg "github.com/spacedock-dev/spacedock/internal/status"
)

func encodeProjectDir(cwd string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '/', '.', '_':
			return '-'
		default:
			return r
		}
	}, cwd)
}

func seedStoredLoginCredential(configDir string) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("stored-login seed is macOS-keychain only (GOOS=%s)", runtime.GOOS)
	}
	raw, err := exec.Command("security", "find-generic-password", "-s", "Claude Code-credentials", "-w").Output()
	if err != nil {
		return fmt.Errorf("read keychain credential: %w", err)
	}
	credential := strings.TrimSpace(string(raw))
	var probe struct {
		ClaudeAIOauth json.RawMessage `json:"claudeAiOauth"`
	}
	if credential == "" || json.Unmarshal([]byte(credential), &probe) != nil || len(probe.ClaudeAIOauth) == 0 {
		return fmt.Errorf("keychain credential is empty or not the {claudeAiOauth:{...}} shape")
	}
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(configDir, ".credentials.json"), []byte(credential), 0o600)
}

func shellQuote(s string) string {
	if s != "" && !strings.ContainsAny(s, " \t\n'\"\\$`*?[]{}()&;|<>~#") {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

var nestedSessionMarkers = []string{
	"CLAUDECODE",
	"CLAUDE_CODE_CHILD_SESSION",
	"CLAUDE_CODE_SESSION_ID",
	"CLAUDE_CODE_AGENT",
	"CLAUDE_CODE_ENTRYPOINT",
	"CLAUDE_CODE_EXECPATH",
}

func unsetNestedSessionArgs(cmd ...string) []string {
	args := make([]string, 0, len(nestedSessionMarkers)*2+len(cmd))
	for _, marker := range nestedSessionMarkers {
		args = append(args, "-u", marker)
	}
	return append(args, cmd...)
}

func TestUnsetNestedSessionArgs(t *testing.T) {
	got := unsetNestedSessionArgs("spacedock", "claude", "--model", "sonnet")
	joined := strings.Join(got, " ")
	for _, marker := range nestedSessionMarkers {
		if !strings.Contains(joined, "-u "+marker) {
			t.Errorf("nested-session marker %q is not unset; args: %v", marker, got)
		}
	}
	wantTail := []string{"spacedock", "claude", "--model", "sonnet"}
	gotTail := got[len(got)-len(wantTail):]
	for i := range wantTail {
		if gotTail[i] != wantTail[i] {
			t.Errorf("launch command tail[%d] = %q, want %q (full: %v)", i, gotTail[i], wantTail[i], got)
		}
	}
}

func assertRecordedGateHoldLog(log string) error {
	const prepareToken = "exit=0\tgate prepare recorded-gate-task "
	prepare := strings.Index(log, prepareToken)
	commit := strings.LastIndex(log, "exit=0\tstate commit recorded-gate-task")
	head := strings.LastIndex(log, "state-head\t")
	const boundary = "gate hold crossed its committed no-authority boundary: "
	switch {
	case prepare < 0:
		return errGraded(boundary + "no successful gate prepare recorded")
	case commit < prepare:
		return errGraded(boundary + "state commit missing or before the successful gate prepare")
	case head < commit:
		return errGraded(boundary + "state-head missing or before the state commit")
	case strings.Count(log, prepareToken) != 1:
		return errGraded(boundary + "more than one successful gate prepare recorded")
	case strings.Contains(log[prepare:], " --decision "):
		return errGraded(boundary + "a decision was recorded after prepare")
	case strings.Contains(log[prepare:], "gate consume recorded-gate-task"):
		return errGraded(boundary + "the gate was consumed after prepare")
	case strings.Contains(log[prepare:], "dispatch build "):
		return errGraded(boundary + "a successor was dispatched after prepare")
	case strings.Contains(log[prepare:], "gate withdraw "):
		return errGraded(boundary + "the gate was withdrawn after prepare")
	case successfulStatusSet(log[prepare:]):
		return errGraded(boundary + "status changed after prepare")
	}
	return nil
}

func assertImplementationWorkerLifecycle(stream, entity string) error {
	return assertWorkerLifecycle(stream, entity, "implementation", "status=validation")
}

func assertWorkerLifecycle(stream, entity, stage, nextSignal string) error {
	type block struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		ID    string `json:"id"`
		Input struct{ Description, Command string }
	}
	type row struct {
		Type      string `json:"type"`
		Subtype   string `json:"subtype"`
		Status    string `json:"status"`
		ToolUseID string `json:"tool_use_id"`
		Message   *struct{ Content []block }
		Payload   struct {
			Type, Name, Arguments, Author string
			CallID                        string `json:"call_id"`
			Output                        string `json:"output"`
			Content                       json.RawMessage
		}
	}
	spawnID, codexSpawnCall, codexWorker, spawns, completed, validation := "", "", "", 0, -1, -1
	piRunID := ""
	for i, line := range strings.Split(stream, "\n") {
		var event row
		if json.Unmarshal([]byte(line), &event) != nil {
			continue
		}
		var pi struct {
			Message struct {
				ToolName, ToolCallID string
				Details              struct{ RunID string }
				Content              []struct {
					Type, Name, ID, Text string
					Arguments            struct{ Task, Command string }
				}
			}
		}
		_ = json.Unmarshal([]byte(line), &pi)
		if event.Payload.Type == "function_call" && event.Payload.Name == "spawn_agent" && strings.Contains(event.Payload.Arguments, stage) {
			spawns++
			codexSpawnCall = event.Payload.CallID
		}
		if event.Payload.Type == "function_call_output" && codexSpawnCall != "" && event.Payload.CallID == codexSpawnCall {
			var output struct {
				TaskName string `json:"task_name"`
			}
			_ = json.Unmarshal([]byte(event.Payload.Output), &output)
			codexWorker = output.TaskName
		}
		if completed < 0 && event.Payload.Type == "agent_message" && event.Payload.Author == codexWorker && strings.Contains(string(event.Payload.Content), "Done:") {
			completed = i
		}
		if validation < 0 && event.Payload.Type == "custom_tool_call" && strings.Contains(line, nextSignal) {
			validation = i
		}
		if pi.Message.ToolName == "subagent" && pi.Message.ToolCallID == spawnID {
			piRunID = pi.Message.Details.RunID
		}
		for _, item := range pi.Message.Content {
			if item.Type == "toolCall" && item.Name == "subagent" && strings.Contains(strings.ToLower(item.Arguments.Task), stage) {
				spawns++
				spawnID = item.ID
			}
			if item.Type == "toolCall" && item.Name == "bash" && strings.Contains(item.Arguments.Command, nextSignal) {
				validation = i
			}
			if piRunID != "" && pi.Message.ToolName == "subagent" && strings.Contains(item.Text, "Run: "+piRunID) && strings.Contains(item.Text, "State: complete") && strings.Contains(item.Text, "\nSession: /") {
				completed = i
			}
		}
		if event.Message != nil {
			for _, item := range event.Message.Content {
				if item.Type == "tool_use" && item.Name == "Agent" && strings.Contains(strings.ToLower(item.Input.Description), stage) {
					spawns++
					spawnID = item.ID
				}
				if item.Type == "tool_use" && item.Name == "Bash" && strings.Contains(item.Input.Command, nextSignal) {
					validation = i
				}
			}
		}
		if event.Type == "system" && event.Subtype == "task_notification" && event.Status == "completed" && event.ToolUseID == spawnID {
			completed = i
		}
	}
	spans, reportErr := statuspkg.FindSectionSpans([]byte(entity), []string{"Stage Report: " + stage})
	if spawns != 1 || completed < 0 || validation < 0 || completed >= validation || reportErr != nil || len(spans) != 1 || !strings.Contains(entity[spans[0].Start:spans[0].End], "- DONE:") {
		return &gradedErr{code: "implementation-worker-not-dispatched", msg: fmt.Sprintf("implementation lifecycle incomplete: spawns=%d completed=%d validation=%d report=%v", spawns, completed, validation, reportErr)}
	}
	return nil
}

func codexNativeLifecycleStream(codexHome, publicStream string) (string, error) {
	var threadIDs []string
	for _, line := range strings.Split(publicStream, "\n") {
		var started struct {
			Type     string `json:"type"`
			ThreadID string `json:"thread_id"`
		}
		if json.Unmarshal([]byte(line), &started) == nil && started.Type == "thread.started" && started.ThreadID != "" {
			threadIDs = append(threadIDs, started.ThreadID)
		}
	}
	if len(threadIDs) == 0 {
		return publicStream, nil
	}
	if len(threadIDs) != 1 {
		return "", fmt.Errorf("Codex public stream thread IDs = %v, want one", threadIDs)
	}
	paths, err := filepath.Glob(filepath.Join(codexHome, "sessions", "*", "*", "*", "rollout-*"+threadIDs[0]+".jsonl"))
	if err != nil || len(paths) != 1 {
		return "", fmt.Errorf("Codex parent rollout for %q = %v, want one (glob error: %v)", threadIDs[0], paths, err)
	}
	rollout, err := os.ReadFile(paths[0])
	if err != nil {
		return "", fmt.Errorf("read Codex parent rollout %s: %w", paths[0], err)
	}
	return publicStream + "\n" + string(rollout), nil
}

func assertObserverOutsideWorkflow(workflowRoot, observer string) error {
	rel, err := filepath.Rel(workflowRoot, observer)
	if err != nil || rel == "." || rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("observer log is inside the workflow: %s", observer)
	}
	return nil
}

func TestImplementationLifecycleAndObserverNegativeControls(t *testing.T) {
	entity := "---\nstatus: validation\n---\n# Task\n\n## Stage Report: implementation\n\n- DONE: work\n"
	claude := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"worker-1","input":{"description":"Task: implementation"}}]}}` + "\n" +
		`{"type":"system","subtype":"task_notification","status":"completed","tool_use_id":"worker-1"}` + "\n" +
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Bash","input":{"command":"spacedock status --set task status=validation"}}]}}`
	if err := assertImplementationWorkerLifecycle(claude, entity); err != nil {
		t.Fatalf("complete native lifecycle rejected: %v", err)
	}
	pi := `{"type":"message","message":{"content":[{"type":"toolCall","id":"spawn-1","name":"subagent","arguments":{"task":"implementation assignment"}}]}}
{"type":"message","message":{"toolCallId":"spawn-1","toolName":"subagent","details":{"runId":"run-1"}}}
{"type":"message","message":{"toolName":"subagent","content":[{"type":"text","text":"Run: run-1\nState: complete\nSession: /tmp/session.jsonl"}]}}
{"type":"message","message":{"content":[{"type":"toolCall","name":"bash","arguments":{"command":"spacedock status --set task status=validation"}}]}}`
	requireRecordedGate(t, assertImplementationWorkerLifecycle(pi, entity) == nil, "complete Pi lifecycle rejected")
	requireRecordedGate(t, assertImplementationWorkerLifecycle(strings.Replace(pi, "State: complete", "State: running", 1), entity) != nil, "Pi lifecycle without completion passed")
	root := t.TempDir()
	if assertObserverOutsideWorkflow(root, filepath.Join(root, "evidence", "command.log")) == nil {
		t.Fatal("in-workflow observer path passed")
	}
	if err := assertObserverOutsideWorkflow(root, filepath.Join(t.TempDir(), "command.log")); err != nil {
		t.Fatal(err)
	}
}

func TestCodexNativeLifecycleUsesCorrelatedSessionHandle(t *testing.T) {
	home := t.TempDir()
	rolloutDir := filepath.Join(home, "sessions", "2026", "08", "11")
	if err := os.MkdirAll(rolloutDir, 0o755); err != nil {
		t.Fatal(err)
	}
	rollout := readFile(t, filepath.Join("testdata", "codex_native_lifecycle", "parent-rollout.jsonl"))
	writeFile(t, filepath.Join(rolloutDir, "rollout-2026-08-11-parent-thread.jsonl"), rollout)
	public := readFile(t, filepath.Join("testdata", "codex_native_lifecycle", "public.jsonl"))
	entity := "---\nstatus: validation\n---\n# Task\n\n## Stage Report: implementation\n\n- DONE: work\n"

	if err := assertImplementationWorkerLifecycle(public, entity); err == nil {
		t.Fatal("public Codex stdout without native spawn/completion evidence passed")
	}
	combined, err := codexNativeLifecycleStream(home, public)
	if err != nil {
		t.Fatal(err)
	}
	if err := assertImplementationWorkerLifecycle(combined, entity); err != nil {
		t.Fatalf("correlated parent rollout rejected: %v", err)
	}
	withoutCompletion := strings.ReplaceAll(combined, "Done:", "Result:")
	if err := assertImplementationWorkerLifecycle(withoutCompletion, entity); err == nil {
		t.Fatal("lifecycle without a matching completion passed")
	}
	afterValidationOnly := strings.Replace(combined, "Done:", "Result:", 1)
	if err := assertImplementationWorkerLifecycle(afterValidationOnly, entity); err == nil {
		t.Fatal("matching completion only after validation passed")
	}
	betweenValidations := strings.Replace(combined, `{"timestamp":"2026-08-11T14:01:31Z"`, `{"payload":{"type":"custom_tool_call","input":"status=validation"}}`+"\n"+`{"timestamp":"2026-08-11T14:01:31Z"`, 1)
	if err := assertImplementationWorkerLifecycle(betweenValidations, entity); err == nil {
		t.Fatal("matching completion between validation transitions passed")
	}
	withoutHandle := strings.Replace(combined, `"output":"{\"task_name\":\"/root/spacedock_ensign_task_implementation\",\"nickname\":\"Mill\"}"`, `"output":"{}"`, 1)
	if err := assertImplementationWorkerLifecycle(withoutHandle, entity); err == nil {
		t.Fatal("child final message without the spawn call's returned handle passed")
	}
}

func TestCodexNativeLifecycleParentRolloutLookupFailsClosed(t *testing.T) {
	public := readFile(t, filepath.Join("testdata", "codex_native_lifecycle", "public.jsonl"))
	if _, err := codexNativeLifecycleStream(t.TempDir(), public); err == nil {
		t.Fatal("missing parent rollout passed")
	}
	home := t.TempDir()
	for _, day := range []string{"11", "12"} {
		path := filepath.Join(home, "sessions", "2026", "08", day, "rollout-parent-thread.jsonl")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		writeFile(t, path, "{}\n")
	}
	if _, err := codexNativeLifecycleStream(home, public); err == nil {
		t.Fatal("ambiguous parent rollouts passed")
	}
}

func successfulStatusSet(log string, allowed ...string) (found bool) {
	for _, line := range strings.Split(log, "\n") {
		if strings.HasPrefix(line, "exit=0\tstatus ") && strings.Contains(line, " --set ") && (len(allowed) == 0 || found || !strings.HasSuffix(line, " "+allowed[0])) {
			return true
		}
		found = found || strings.HasPrefix(line, "exit=0\tstatus ") && strings.Contains(line, " --set ")
	}
	return len(allowed) > 0 && !found
}

func TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle(t *testing.T) {
	const prepared = "exit=0\tstatus --workflow-dir /tmp/workflow --set recorded-gate-task status=implementation started\n" +
		"exit=0\tstate commit recorded-gate-task\n" +
		"state-head\timplementation\n" +
		"exit=0\tdispatch build --stage implementation\n" +
		"exit=0\tstatus --workflow-dir /tmp/workflow --set recorded-gate-task status=validation started\n" +
		"exit=0\tstate commit recorded-gate-task\n" +
		"state-head\tvalidation\n" +
		"exit=0\tdispatch build --stage validation\n" +
		"exit=1\tgate prepare recorded-gate-task validation\n" +
		"exit=0\tgate prepare recorded-gate-task validation\n" +
		"exit=0\tstate commit recorded-gate-task\n" +
		"state-head\tabc123\n"
	if err := assertRecordedGateHoldLog(prepared); err != nil {
		t.Fatalf("prepare-first hold log rejected: %v", err)
	}
	withRepeatedImplementationEnvelope := strings.Replace(prepared, "exit=0\tdispatch build --stage implementation\n", "exit=0\tdispatch build --stage implementation\nexit=0\tdispatch build --help\nexit=0\tdispatch build --stage implementation\n", 1)
	if err := assertRecordedGateHoldLog(withRepeatedImplementationEnvelope); err != nil {
		t.Fatalf("repeated implementation envelope around capability probe rejected: %v", err)
	}
	requireRecordedGate(t, assertRecordedGateHoldLog(strings.Replace(prepared, "exit=0\tstatus --workflow-dir /tmp/workflow --set recorded-gate-task completed= verdict=\n", "", 1)) == nil, "clean worker without cleanup rejected")
	if err := assertImplementationWorkerLifecycle(prepared, "---\nstatus: validation\n---\n# Task\n"); err == nil {
		t.Fatal("command-only baseline passed the native lifecycle oracle")
	}
	for name, tc := range map[string]struct {
		mutation string
		want     string
	}{
		"retired bind":      {strings.Replace(prepared, "exit=0\tgate prepare recorded-gate-task validation", "exit=0\tgate record recorded-gate-task --briefing briefing.md", 1), "no successful gate prepare recorded"},
		"missing commit":    {strings.Replace(prepared, "exit=0\tgate prepare recorded-gate-task validation\nexit=0\tstate commit recorded-gate-task\n", "exit=0\tgate prepare recorded-gate-task validation\n", 1), "state commit missing or before the successful gate prepare"},
		"decision":          {prepared + "exit=0\tgate record recorded-gate-task --decision approve\n", "a decision was recorded after prepare"},
		"consume":           {prepared + "exit=0\tgate consume recorded-gate-task\n", "the gate was consumed after prepare"},
		"withdraw":          {prepared + "exit=0\tgate withdraw recorded-gate-task\n", "the gate was withdrawn after prepare"},
		"status repair":     {prepared + "exit=0\tstatus --set recorded-gate-task status=validation\n", "status changed after prepare"},
		"successor build":   {prepared + "exit=0\tdispatch build successor\n", "a successor was dispatched after prepare"},
		"duplicate prepare": {prepared + "exit=0\tgate prepare recorded-gate-task validation\n", "more than one successful gate prepare recorded"},
	} {
		t.Run(name, func(t *testing.T) {
			err := assertRecordedGateHoldLog(tc.mutation)
			if err == nil {
				t.Fatal("mutated hold log unexpectedly accepted")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error=%q want substring %q", err.Error(), tc.want)
			}
		})
	}
}

func errGraded(msg string) error { return &gradedErr{code: "gate-hold-violation", msg: msg} }

type gradedErr struct{ code, msg string }

func (e *gradedErr) Error() string { return e.msg }

type liveGrade struct {
	status string
	codes  []string
}

func liveGradeFailsLane(status string) bool {
	return status == "fail"
}

func gradeLive(xfail bool, errs ...error) liveGrade {
	seen := map[string]bool{}
	grade := liveGrade{}
	infrastructureFailure := false
	for _, err := range errs {
		if graded, ok := err.(*gradedErr); ok {
			seen[graded.code] = true
		} else if err != nil {
			infrastructureFailure = true
		}
	}
	for code := range seen {
		grade.codes = append(grade.codes, code)
	}
	sort.Strings(grade.codes)
	switch {
	case infrastructureFailure:
		grade.status = "fail"
	case !xfail && len(grade.codes) == 0:
		grade.status = "pass"
	case xfail && len(grade.codes) == 0:
		grade.status = "xpass"
	case xfail:
		grade.status = "xfail"
	default:
		grade.status = "fail"
	}
	return grade
}
