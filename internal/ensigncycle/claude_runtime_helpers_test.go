package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

func TestClaudeTODOModelScope(t *testing.T) {
	for _, tc := range []struct {
		model, family string
		want          bool
	}{
		{"sonnet", "sonnet", true},
		{"claude-sonnet-5", "sonnet", true},
		{"claude-opus-4-8", "sonnet", false},
		{"opus", "opus", true},
		{"claude-opus-4-8", "opus", true},
		{"openrouter/opossum", "opus", false},
	} {
		if got := claudeModelFamily(tc.model, tc.family); got != tc.want {
			t.Errorf("claudeModelFamily(%q, %q) = %t, want %t", tc.model, tc.family, got, tc.want)
		}
	}
}

func TestClaudeRejectionFlowTODOModelScope(t *testing.T) {
	for model, want := range map[string]bool{
		"sonnet": true, "claude-sonnet-5": true,
		"opus": true, "claude-opus-4-8": true,
		"haiku": false, "openrouter/opossum": false,
	} {
		if got := claudeRejectionFlowTODOModel(model); got != want {
			t.Errorf("claudeRejectionFlowTODOModel(%q) = %t, want %t", model, got, want)
		}
	}
}

func claudeModelFamily(model, family string) bool {
	model = strings.ToLower(strings.TrimSpace(model))
	return model == family || strings.Contains(model, "claude-"+family+"-")
}

func claudeRejectionFlowTODOModel(model string) bool {
	return claudeModelFamily(model, "opus") || claudeModelFamily(model, "sonnet")
}

func claudeSonnetGateGuardrailTODO(model string) string {
	if !claudeModelFamily(model, "sonnet") {
		return ""
	}
	return "TODO(3zzpdw704df1g8pg1x9thzmw): Claude Sonnet gate-guardrail is temporarily non-evidence for local task sonnet-gate-guardrail-no-authority pending a fresh passing promotion run"
}

func TestClaudeSonnetGateGuardrailTODOModelScope(t *testing.T) {
	for model, wantSkip := range map[string]bool{
		"sonnet":          true,
		"claude-sonnet-5": true,
		"opus":            false,
		"claude-opus-4-8": false,
		"haiku":           false,
	} {
		got := claudeSonnetGateGuardrailTODO(model)
		if (got != "") != wantSkip {
			t.Errorf("claudeSonnetGateGuardrailTODO(%q) = %q, wantSkip %t", model, got, wantSkip)
		}
	}
}

func assertRecordedGateHoldLog(log string) error {
	const prepareToken = "exit=0\tgate prepare recorded-gate-task "
	prepare, commit, head := strings.Index(log, prepareToken), strings.LastIndex(log, "exit=0\tstate commit recorded-gate-task"), strings.LastIndex(log, "state-head\t")
	if prepare < 0 || commit < prepare || head < commit || strings.Count(log, prepareToken) != 1 ||
		strings.Contains(log[prepare:], " --decision ") || strings.Contains(log[prepare:], "gate consume recorded-gate-task") || strings.Contains(log[prepare:], "dispatch build ") {
		return errGraded("gate hold crossed its committed no-authority boundary")
	}
	return nil
}

func TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle(t *testing.T) {
	const prepared = "exit=1\tgate prepare recorded-gate-task validation\n" +
		"exit=0\tgate prepare recorded-gate-task validation\n" +
		"exit=0\tstate commit recorded-gate-task\n" +
		"state-head\tabc123\n"
	if err := assertRecordedGateHoldLog(prepared); err != nil {
		t.Fatalf("prepare-first hold log rejected: %v", err)
	}
	for name, mutation := range map[string]string{
		"retired bind":      strings.Replace(prepared, "exit=0\tgate prepare recorded-gate-task validation", "exit=0\tgate record recorded-gate-task --briefing briefing.md", 1),
		"missing commit":    strings.Replace(prepared, "exit=0\tstate commit recorded-gate-task\n", "", 1),
		"decision":          prepared + "exit=0\tgate record recorded-gate-task --decision approve\n",
		"consume":           prepared + "exit=0\tgate consume recorded-gate-task\n",
		"successor build":   prepared + "exit=0\tdispatch build successor\n",
		"duplicate prepare": prepared + "exit=0\tgate prepare recorded-gate-task validation\n",
	} {
		t.Run(name, func(t *testing.T) {
			if err := assertRecordedGateHoldLog(mutation); err == nil {
				t.Fatal("mutated hold log unexpectedly accepted")
			}
		})
	}
}

func errGraded(msg string) error { return &gradedErr{msg} }

type gradedErr struct{ msg string }

func (e *gradedErr) Error() string { return e.msg }
