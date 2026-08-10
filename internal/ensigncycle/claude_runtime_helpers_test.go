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

func assertRecordedGateHoldLog(log string, requireImplementation ...bool) error {
	const prepareToken = "exit=0\tgate prepare recorded-gate-task "
	const cleanupToken = " --set recorded-gate-task completed= verdict="
	prepare := strings.Index(log, prepareToken)
	cleanup := strings.Index(log, cleanupToken)
	validation := strings.Index(log, " --set recorded-gate-task status=validation started")
	commit := strings.LastIndex(log, "exit=0\tstate commit recorded-gate-task")
	head := strings.LastIndex(log, "state-head\t")
	dispatches := strings.Split(log[:max(prepare, 0)], "exit=0\tdispatch build ")
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
	case len(requireImplementation) > 0 && requireImplementation[0] && strings.Count(log[:prepare], cleanupToken) > 1:
		return errGraded(boundary + "more than one terminal-field cleanup recorded")
	case len(requireImplementation) > 0 && requireImplementation[0] && (cleanup < 0 || cleanup > validation):
		return errGraded(boundary + "terminal-field cleanup missing or after validation")
	case len(requireImplementation) > 0 && requireImplementation[0] && strings.Count(log[:prepare], " --stage implementation") == 1 && strings.Count(log[:prepare], " --stage validation") != 1:
		return &gradedErr{code: "dispatch-envelope-not-acknowledged", msg: boundary + "validation dispatch envelope count was not one"}
	case len(requireImplementation) > 0 && requireImplementation[0] && (len(dispatches) != 3 || !strings.Contains(" "+strings.SplitN(dispatches[1], "\n", 2)[0]+" ", " --stage implementation ") || successfulStatusSet(log[:strings.Index(log, "exit=0\tdispatch build ")], "status=implementation started") || successfulStatusSet(strings.SplitN(dispatches[1], "\n", 2)[1], "completed= verdict=", "status=validation started")):
		return &gradedErr{code: "implementation-worker-not-dispatched", msg: boundary + "implementation was not dispatched before validation"}
	}
	return nil
}

func assertImplementationWorkerLifecycle(stream, entity string) error {
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
			Content                       json.RawMessage
		}
	}
	spawnID, codexWorker, spawns, completed, validation := "", "", 0, -1, -1
	for i, line := range strings.Split(stream, "\n") {
		var event row
		if json.Unmarshal([]byte(line), &event) != nil {
			continue
		}
		if event.Payload.Type == "function_call" && event.Payload.Name == "spawn_agent" && strings.Contains(event.Payload.Arguments, "implementation") {
			spawns++
			var args struct {
				TaskName string `json:"task_name"`
			}
			_ = json.Unmarshal([]byte(event.Payload.Arguments), &args)
			codexWorker = "/root/" + args.TaskName
		}
		if event.Payload.Type == "agent_message" && event.Payload.Author == codexWorker && strings.Contains(string(event.Payload.Content), "Done:") {
			completed = i
		}
		if event.Payload.Type == "custom_tool_call" && strings.Contains(line, "status=validation") {
			validation = i
		}
		if event.Message != nil {
			for _, item := range event.Message.Content {
				if item.Type == "tool_use" && item.Name == "Agent" && strings.Contains(strings.ToLower(item.Input.Description), "implementation") {
					spawns++
					spawnID = item.ID
				}
				if item.Type == "tool_use" && item.Name == "Bash" && strings.Contains(item.Input.Command, "status=validation") {
					validation = i
				}
			}
		}
		if event.Type == "system" && event.Subtype == "task_notification" && event.Status == "completed" && event.ToolUseID == spawnID {
			completed = i
		}
	}
	spans, reportErr := statuspkg.FindSectionSpans([]byte(entity), []string{"Stage Report: implementation"})
	if spawns != 1 || completed < 0 || validation < 0 || completed >= validation || reportErr != nil || len(spans) != 1 || !strings.Contains(entity[spans[0].Start:spans[0].End], "- DONE:") {
		return &gradedErr{code: "implementation-worker-not-dispatched", msg: fmt.Sprintf("implementation lifecycle incomplete: spawns=%d completed=%d validation=%d report=%v", spawns, completed, validation, reportErr)}
	}
	return nil
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
	root := t.TempDir()
	if assertObserverOutsideWorkflow(root, filepath.Join(root, "evidence", "command.log")) == nil {
		t.Fatal("in-workflow observer path passed")
	}
	if err := assertObserverOutsideWorkflow(root, filepath.Join(t.TempDir(), "command.log")); err != nil {
		t.Fatal(err)
	}
}

func successfulStatusSet(log string, allowed ...string) bool {
	var found []string
	for _, line := range strings.Split(log, "\n") {
		if strings.HasPrefix(line, "exit=0\tstatus ") && strings.Contains(line, " --set ") {
			found = append(found, line)
		}
	}
	if len(found) != len(allowed) {
		return true
	}
	for i := range found {
		if !strings.HasSuffix(found[i], " "+allowed[i]) {
			return true
		}
	}
	return false
}

func TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle(t *testing.T) {
	const prepared = "exit=0\tstatus --workflow-dir /tmp/workflow --set recorded-gate-task status=implementation started\n" +
		"exit=0\tstate commit recorded-gate-task\n" +
		"state-head\timplementation\n" +
		"exit=0\tdispatch build --stage implementation\n" +
		"exit=0\tstatus --workflow-dir /tmp/workflow --set recorded-gate-task completed= verdict=\n" +
		"exit=0\tstatus --workflow-dir /tmp/workflow --set recorded-gate-task status=validation started\n" +
		"exit=0\tstate commit recorded-gate-task\n" +
		"state-head\tvalidation\n" +
		"exit=0\tdispatch build --stage validation\n" +
		"exit=1\tgate prepare recorded-gate-task validation\n" +
		"exit=0\tgate prepare recorded-gate-task validation\n" +
		"exit=0\tstate commit recorded-gate-task\n" +
		"state-head\tabc123\n"
	if err := assertRecordedGateHoldLog(prepared, true); err != nil {
		t.Fatalf("prepare-first hold log rejected: %v", err)
	}
	if err := assertImplementationWorkerLifecycle(prepared, "---\nstatus: validation\n---\n# Task\n"); err == nil {
		t.Fatal("command-only baseline passed the native lifecycle oracle")
	}
	for name, tc := range map[string]struct {
		mutation string
		want     string
	}{
		"retired bind":         {strings.Replace(prepared, "exit=0\tgate prepare recorded-gate-task validation", "exit=0\tgate record recorded-gate-task --briefing briefing.md", 1), "no successful gate prepare recorded"},
		"missing commit":       {strings.Replace(prepared, "exit=0\tgate prepare recorded-gate-task validation\nexit=0\tstate commit recorded-gate-task\n", "exit=0\tgate prepare recorded-gate-task validation\n", 1), "state commit missing or before the successful gate prepare"},
		"decision":             {prepared + "exit=0\tgate record recorded-gate-task --decision approve\n", "a decision was recorded after prepare"},
		"consume":              {prepared + "exit=0\tgate consume recorded-gate-task\n", "the gate was consumed after prepare"},
		"withdraw":             {prepared + "exit=0\tgate withdraw recorded-gate-task\n", "the gate was withdrawn after prepare"},
		"duplicate cleanup":    {strings.Replace(prepared, "completed= verdict=\n", "completed= verdict=\nexit=0\tstatus --set recorded-gate-task completed= verdict=\n", 1), "more than one terminal-field cleanup"},
		"reordered cleanup":    {strings.Replace(prepared, "completed= verdict=\nexit=0\tstatus --workflow-dir /tmp/workflow --set recorded-gate-task status=validation started", "status=validation started\nexit=0\tstatus --set recorded-gate-task completed= verdict=", 1), "cleanup missing or after validation"},
		"missing cleanup":      {strings.Replace(prepared, "exit=0\tstatus --workflow-dir /tmp/workflow --set recorded-gate-task completed= verdict=\n", "", 1), "cleanup missing or after validation"},
		"post-prepare cleanup": {strings.Replace(prepared, "exit=0\tstatus --workflow-dir /tmp/workflow --set recorded-gate-task completed= verdict=\n", "", 1) + "exit=0\tstatus --set recorded-gate-task completed= verdict=\n", "status changed after prepare"},
		"successor build":      {prepared + "exit=0\tdispatch build successor\n", "a successor was dispatched after prepare"},
		"duplicate prepare":    {prepared + "exit=0\tgate prepare recorded-gate-task validation\n", "more than one successful gate prepare recorded"},
	} {
		t.Run(name, func(t *testing.T) {
			err := assertRecordedGateHoldLog(tc.mutation, true)
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
