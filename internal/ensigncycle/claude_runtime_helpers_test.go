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
	prepare := strings.Index(log, prepareToken)
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
	case len(requireImplementation) > 0 && requireImplementation[0] && (len(dispatches) != 2 || !strings.Contains(" "+strings.SplitN(dispatches[1], "\n", 2)[0]+" ", " --stage implementation ") || successfulStatusSet(log[:strings.Index(log, "exit=0\tdispatch build ")], "status=implementation started") || successfulStatusSet(strings.SplitN(dispatches[1], "\n", 2)[1], "status=validation started")):
		return &gradedErr{code: "implementation-worker-not-dispatched", msg: boundary + "implementation was not dispatched before validation"}
	}
	return nil
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
		"exit=1\tgate prepare recorded-gate-task validation\n" +
		"exit=0\tgate prepare recorded-gate-task validation\n" +
		"exit=0\tstate commit recorded-gate-task\n" +
		"state-head\tabc123\n"
	if err := assertRecordedGateHoldLog(prepared, true); err != nil {
		t.Fatalf("prepare-first hold log rejected: %v", err)
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

func gradeLive(expected string, errs ...error) liveGrade {
	seen := map[string]bool{}
	grade := liveGrade{}
	for _, err := range errs {
		if graded, ok := err.(*gradedErr); ok {
			seen[graded.code] = true
		} else if err != nil {
			seen["untyped-semantic-failure"] = true
		}
	}
	for code := range seen {
		grade.codes = append(grade.codes, code)
	}
	sort.Strings(grade.codes)
	switch {
	case expected == "" && len(grade.codes) == 0:
		grade.status = "pass"
	case expected != "" && len(grade.codes) == 0:
		grade.status = "xpass"
	case expected != "" && len(grade.codes) == 1 && grade.codes[0] == expected:
		grade.status = "xfail"
	default:
		grade.status = "fail"
	}
	return grade
}
