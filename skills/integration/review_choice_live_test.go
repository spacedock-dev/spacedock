//go:build live

// ABOUTME: Live check of /spacedock:review's argument handling and choice flow: runs the real skill
// ABOUTME: in claude -p against the built binary and stub programs, and grades tool calls and state.
package integration

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// PATH holds no subspace-tui, so a local hand-off stops at the launcher instead of opening a
// terminal window on the machine running the test.
// Needs SPACEDOCK_BIN (the built spacedock) and SPACEDOCK_LIVE_CLAUDE_CONFIG_DIR (a logged-in
// Claude Code config). SPACEDOCK_LIVE_OUT, when set, receives each test's raw stream.
// claude -p offers AskUserQuestion only when a host answers permission prompts, so each test
// runs one claude process with --permission-prompt-tool stdio and answers the menus itself.

type reviewLive struct {
	t       *testing.T
	root    string
	proj    string
	env     []string
	relay   string
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  *bufio.Reader
	raw     bytes.Buffer
	answers map[string]string
	tools   []liveTool
	asked   []liveQuestion
	text    string
}

type liveTool struct {
	Name  string
	Input map[string]any
}

type liveQuestion struct {
	Question string
	Options  []string
}

func (q liveQuestion) kind() string {
	switch {
	case q.has("Keep") && q.has("Leave"):
		return "keep"
	case strings.Contains(q.Question, "seven days"):
		return "agree"
	case q.has("Web"):
		return "choice"
	}
	return "unknown"
}

func (q liveQuestion) has(token string) bool { return q.label(token) != "" }

func (q liveQuestion) label(token string) string {
	for _, o := range q.Options {
		if strings.Contains(strings.ToLower(o), strings.ToLower(token)) {
			return o
		}
	}
	return ""
}

func (l *reviewLive) flat() string {
	b, _ := json.Marshal(l.tools)
	return string(b)
}

func (l *reviewLive) kinds() []string {
	var out []string
	for _, q := range l.asked {
		out = append(out, q.kind())
	}
	return out
}

func newReviewLive(t *testing.T, answers map[string]string) *reviewLive {
	t.Helper()
	bin, cfg := os.Getenv("SPACEDOCK_BIN"), os.Getenv("SPACEDOCK_LIVE_CLAUDE_CONFIG_DIR")
	if bin == "" || cfg == "" {
		t.Skip("set SPACEDOCK_BIN and SPACEDOCK_LIVE_CLAUDE_CONFIG_DIR")
	}
	claude, err := exec.LookPath("claude")
	if err != nil {
		t.Skip("claude is not on PATH")
	}
	root := t.TempDir()
	progs := filepath.Join(root, "bin")
	if err := os.MkdirAll(progs, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, target := range map[string]string{"spacedock": bin, "claude": claude} {
		if err := os.Symlink(target, filepath.Join(progs, name)); err != nil {
			t.Fatal(err)
		}
	}
	relayLog := filepath.Join(root, "relay-calls")
	stubProgram(t, progs, "subspace-relay-client", "0.16.0-beta.4", "relay-publish-v1", relayLog)
	stubProgram(t, progs, "spacedock-review", "0.16.0-beta.4", "review-invocation-prepare-v1", filepath.Join(root, "review-calls"))
	proj := filepath.Join(root, "proj")
	if err := os.MkdirAll(proj, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(proj, "plan.md"), []byte("# Plan\n\nShip the thing.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("git", "-C", proj, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	l := &reviewLive{t: t, root: root, proj: proj, relay: relayLog, answers: answers, env: []string{
		"HOME=" + os.Getenv("HOME"), "USER=" + os.Getenv("USER"), "LOGNAME=" + os.Getenv("USER"),
		"CLAUDE_CONFIG_DIR=" + cfg, "XDG_STATE_HOME=" + filepath.Join(root, "state"),
		"XDG_DATA_HOME=" + filepath.Join(root, "data"), "TERM=dumb", "PATH=" + progs + ":/usr/bin:/bin",
	}}
	l.cmd = exec.Command(filepath.Join(progs, "claude"), "-p", "--input-format", "stream-json",
		"--output-format", "stream-json", "--verbose", "--permission-prompt-tool", "stdio",
		"--plugin-dir", repoRoot(t), "--allowedTools", "Bash,Read,Write,Skill", "--model", "sonnet")
	l.cmd.Dir, l.cmd.Env = proj, l.env
	if l.stdin, err = l.cmd.StdinPipe(); err != nil {
		t.Fatal(err)
	}
	out, err := l.cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	l.stdout = bufio.NewReaderSize(out, 1<<20)
	if err := l.cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = l.stdin.Close()
		_ = l.cmd.Wait()
		if dir := os.Getenv("SPACEDOCK_LIVE_OUT"); dir != "" {
			_ = os.WriteFile(filepath.Join(dir, strings.ReplaceAll(t.Name(), "/", "_")+".jsonl"), l.raw.Bytes(), 0o644)
		}
	})
	return l
}

func (l *reviewLive) spacedock(args ...string) string {
	l.t.Helper()
	cmd := exec.Command(filepath.Join(l.root, "bin", "spacedock"), args...)
	cmd.Dir, cmd.Env = l.proj, l.env
	out, err := cmd.CombinedOutput()
	if err != nil {
		l.t.Fatalf("spacedock %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func (l *reviewLive) mode() string {
	l.t.Helper()
	var s struct{ Mode string }
	if err := json.Unmarshal([]byte(l.spacedock("review-mode", "get", "--json")), &s); err != nil {
		l.t.Fatal(err)
	}
	return s.Mode
}

func (l *reviewLive) send(v any) {
	l.t.Helper()
	b, _ := json.Marshal(v)
	if _, err := l.stdin.Write(append(b, '\n')); err != nil {
		l.t.Errorf("write to claude: %v", err)
	}
}

func (l *reviewLive) answer(requestID string, input map[string]any) {
	l.t.Helper()
	raw, _ := input["questions"].([]any)
	reply := map[string]any{}
	denied := ""
	for _, r := range raw {
		m, _ := r.(map[string]any)
		q := liveQuestion{Question: m["question"].(string)}
		for _, o := range m["options"].([]any) {
			q.Options = append(q.Options, o.(map[string]any)["label"].(string))
		}
		l.asked = append(l.asked, q)
		want, ok := l.answers[q.kind()]
		if !ok || q.label(want) == "" {
			l.t.Errorf("unexpected question %q with options %v (kind %s)", q.Question, q.Options, q.kind())
			denied = "no answer is available for this question"
			continue
		}
		reply[q.Question] = q.label(want)
	}
	resp := map[string]any{"behavior": "allow", "updatedInput": map[string]any{"questions": raw, "answers": reply}}
	if denied != "" {
		resp = map[string]any{"behavior": "deny", "message": denied}
	}
	l.send(map[string]any{"type": "control_response", "response": map[string]any{
		"subtype": "success", "request_id": requestID, "response": resp}})
}

func (l *reviewLive) turn(prompt string) {
	l.t.Helper()
	l.tools, l.asked, l.text = nil, nil, ""
	l.send(map[string]any{"type": "user", "message": map[string]any{"role": "user", "content": prompt}})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			line, err := l.stdout.ReadBytes('\n')
			if len(line) > 0 {
				l.raw.Write(line)
				if l.handle(prompt, line) {
					return
				}
			}
			if err != nil {
				l.t.Errorf("claude ended before the turn %q finished: %v", prompt, err)
				return
			}
		}
	}()
	select {
	case <-done:
	case <-time.After(6 * time.Minute):
		_ = l.stdin.Close()
		l.t.Fatalf("claude turn %q did not finish in six minutes", prompt)
	}
}

func (l *reviewLive) handle(prompt string, line []byte) bool {
	var ev struct {
		Type      string `json:"type"`
		RequestID string `json:"request_id"`
		IsError   bool   `json:"is_error"`
		Request   struct {
			Subtype  string         `json:"subtype"`
			ToolName string         `json:"tool_name"`
			Input    map[string]any `json:"input"`
		} `json:"request"`
		Message struct {
			Content []struct {
				Type  string         `json:"type"`
				Text  string         `json:"text"`
				Name  string         `json:"name"`
				Input map[string]any `json:"input"`
			} `json:"content"`
		} `json:"message"`
	}
	if json.Unmarshal(line, &ev) != nil {
		return false
	}
	switch ev.Type {
	case "control_request":
		if ev.Request.Subtype != "can_use_tool" {
			return false
		}
		if ev.Request.ToolName == "AskUserQuestion" {
			l.answer(ev.RequestID, ev.Request.Input)
			return false
		}
		l.t.Errorf("claude asked permission for %s, which the test does not allow: %v", ev.Request.ToolName, ev.Request.Input)
		l.send(map[string]any{"type": "control_response", "response": map[string]any{
			"subtype": "success", "request_id": ev.RequestID,
			"response": map[string]any{"behavior": "deny", "message": "not allowed in this test"}}})
	case "result":
		if ev.IsError {
			l.t.Errorf("claude turn %q ended in error:\n%s", prompt, line)
		}
		return true
	case "assistant":
		for _, c := range ev.Message.Content {
			switch c.Type {
			case "text":
				l.text += c.Text + "\n"
			case "tool_use":
				l.tools = append(l.tools, liveTool{c.Name, c.Input})
			}
		}
	}
	return false
}

func (l *reviewLive) said(want string) bool {
	return strings.Contains(strings.NewReplacer("*", "", "`", "").Replace(l.text), want)
}

func (l *reviewLive) noPublish() {
	l.t.Helper()
	calls, _ := os.ReadFile(l.relay)
	if reviewPublishCount(calls) != 0 {
		l.t.Errorf("subspace-relay-client was asked to publish: %q", calls)
	}
}

func (l *reviewLive) wantMenus(want ...string) {
	l.t.Helper()
	got := l.kinds()
	if strings.Join(got, ",") != strings.Join(want, ",") {
		l.t.Fatalf("AskUserQuestion menus = %v, want %v\ntools: %s\ntext:\n%s", got, want, l.flat(), l.text)
	}
}

func (l *reviewLive) handedTo(skill, file string) bool {
	for _, tool := range l.tools {
		if s, _ := tool.Input["skill"].(string); tool.Name == "Skill" && s == skill {
			if args, _ := tool.Input["args"].(string); strings.Contains(args, file) {
				return true
			}
		}
	}
	return false
}

func TestLiveReviewNoArgumentShowsTheChoiceAndPublishesNothing(t *testing.T) {
	l := newReviewLive(t, map[string]string{"keep": "Leave"})
	l.spacedock("review-mode", "set", "local")
	l.turn("/spacedock:review")
	l.wantMenus("keep")
	if q := l.asked[0]; !q.has("Keep") || !q.has("Switch") || !q.has("Clear") || !q.has("Leave") {
		t.Errorf("menu does not offer keep, switch, clear and leave: %v", q.Options)
	}
	if !strings.Contains(l.flat(), "review-mode get") {
		t.Errorf("did not read the recorded choice: %s", l.flat())
	}
	for _, banned := range []string{"review-mode set", "review-mode clear", "subspace:", "plan.md"} {
		if strings.Contains(l.flat(), banned) {
			t.Errorf("no-argument run acted on %q: %s", banned, l.flat())
		}
	}
	if !l.said("Terminal (TUI): the file stays on this machine") {
		t.Errorf("choice not shown in the agreed words:\n%s", l.text)
	}
	if got := l.mode(); got != "local" {
		t.Errorf("recorded choice changed to %q", got)
	}
	l.noPublish()
}

func TestLiveReviewNoArgumentWithNothingRecordedAsksTerminalOrWeb(t *testing.T) {
	l := newReviewLive(t, map[string]string{"choice": "Web", "agree": "No"})
	l.turn("/spacedock:review")
	l.wantMenus("choice", "agree")
	if q := l.asked[1]; !strings.Contains(q.Question, "review.spacedock.md") || !q.has("Yes") || !q.has("No") {
		t.Errorf("agreement question lacks the host or its Yes/No options: %q %v", q.Question, q.Options)
	}
	if strings.Contains(l.flat(), "review-mode set") || l.mode() != "" {
		t.Errorf("declining the upload still recorded %q: %s", l.mode(), l.flat())
	}
	if !l.said("Terminal") {
		t.Errorf("a decline did not offer the Terminal (TUI):\n%s", l.text)
	}
	if !l.said("/spacedock:review") {
		t.Errorf("the answer did not say the setting can be changed through /spacedock:review:\n%s", l.text)
	}
	l.noPublish()
}

func TestLiveReviewSwitchingToWebAsksForAgreementAndDecliningKeepsTheChoice(t *testing.T) {
	l := newReviewLive(t, map[string]string{"keep": "Switch", "agree": "No"})
	l.spacedock("review-mode", "set", "local")
	l.turn("/spacedock:review")
	l.wantMenus("keep", "agree")
	if strings.Contains(l.flat(), "review-mode set") || l.mode() != "local" {
		t.Errorf("the choice changed although the user declined the upload: %s", l.flat())
	}
	l.noPublish()
}

func TestLiveReviewTerminalWordRecordsTerminalAndSaysSo(t *testing.T) {
	l := newReviewLive(t, map[string]string{})
	l.spacedock("review-mode", "set", "remote")
	l.turn("/spacedock:review TUI")
	l.wantMenus()
	if got := l.mode(); got != "local" {
		t.Errorf("recorded choice = %q, want local", got)
	}
	if !l.said("Terminal (TUI)") || !l.said("/spacedock:review") {
		t.Errorf("did not say Terminal (TUI) is recorded and how to change it:\n%s", l.text)
	}
	if strings.Contains(l.flat(), "subspace:") || strings.Contains(l.flat(), "plan.md") {
		t.Errorf("a reserved word opened or reviewed a file: %s", l.flat())
	}
	l.noPublish()
}

func TestLiveReviewWebWordDeclinedRecordsNothing(t *testing.T) {
	l := newReviewLive(t, map[string]string{"agree": "No"})
	l.turn("/spacedock:review WebUI")
	l.wantMenus("agree")
	if q := l.asked[0]; !strings.Contains(q.Question, "review.spacedock.md") || !q.has("Yes") || !q.has("No") {
		t.Errorf("agreement question lacks the host or its Yes/No options: %q %v", q.Question, q.Options)
	}
	if strings.Contains(l.flat(), "review-mode set") || l.mode() != "" {
		t.Errorf("declining the upload still recorded %q: %s", l.mode(), l.flat())
	}
	if strings.Contains(l.flat(), "subspace:") || strings.Contains(l.flat(), "plan.md") {
		t.Errorf("a reserved word opened or reviewed a file: %s", l.flat())
	}
	if !l.said("/spacedock:review") {
		t.Errorf("the answer did not say the setting can be changed through /spacedock:review:\n%s", l.text)
	}
	l.noPublish()
}

func TestLiveReviewFileArgumentWithNothingRecordedAsksAndDecliningUploadsNothing(t *testing.T) {
	l := newReviewLive(t, map[string]string{"choice": "Web", "agree": "No"})
	l.turn("/spacedock:review plan.md")
	l.wantMenus("choice", "agree")
	if strings.Contains(l.flat(), "review-mode set") || l.mode() != "" {
		t.Errorf("declining the upload still recorded %q: %s", l.mode(), l.flat())
	}
	if strings.Contains(l.flat(), "subspace:") {
		t.Errorf("the file was handed off although the user declined the upload: %s", l.flat())
	}
	l.noPublish()
}

func TestLiveReviewFileArgumentUnderARecordedChoiceAsksNothing(t *testing.T) {
	l := newReviewLive(t, map[string]string{})
	l.spacedock("review-mode", "set", "local")
	l.turn("/spacedock:review plan.md")
	l.wantMenus()
	if !l.handedTo("subspace:r", "plan.md") {
		t.Errorf("plan.md was not handed to subspace:r: %s", l.flat())
	}
	for _, tool := range l.tools {
		if tool.Name == "Write" {
			t.Errorf("a file argument wrote a file: %v", tool.Input)
		}
	}
	l.noPublish()
}

func TestLiveReviewFreeTextIsCarriedOutBeforeAskingAndDecliningUploadsNothing(t *testing.T) {
	l := newReviewLive(t, map[string]string{"choice": "Web", "agree": "No"})
	l.turn("/spacedock:review read plan.md and write a one-sentence summary of it into summary.md")
	l.wantMenus("choice", "agree")
	read, wrote, firstMenu := -1, -1, -1
	for i, tool := range l.tools {
		path, _ := tool.Input["file_path"].(string)
		cmd, _ := tool.Input["command"].(string)
		if tool.Name == "AskUserQuestion" && firstMenu < 0 {
			firstMenu = i
		}
		if read < 0 && strings.Contains(path+cmd, "plan.md") {
			read = i
		}
		if wrote < 0 && (tool.Name == "Write" && strings.HasSuffix(path, "summary.md") ||
			tool.Name == "Bash" && strings.Contains(cmd, "summary.md") && strings.Contains(cmd, ">")) {
			wrote = i
		}
	}
	if read < 0 || wrote < read || firstMenu < wrote {
		t.Errorf("want plan.md read, summary.md written, then the first question; got read=%d write=%d question=%d: %s", read, wrote, firstMenu, l.flat())
	}
	if body, _ := os.ReadFile(filepath.Join(l.proj, "summary.md")); len(strings.TrimSpace(string(body))) == 0 {
		t.Errorf("summary.md is missing or empty")
	}
	if strings.Contains(l.flat(), "review-mode set") || l.mode() != "" {
		t.Errorf("declining the upload still recorded %q: %s", l.mode(), l.flat())
	}
	if strings.Contains(l.flat(), "subspace:") {
		t.Errorf("the file was handed off although the user declined the upload: %s", l.flat())
	}
	l.noPublish()
}

func TestLiveReviewFreeTextWithoutAPathLeavesNoStrayFileInTheProject(t *testing.T) {
	l := newReviewLive(t, map[string]string{})
	l.spacedock("review-mode", "set", "local")
	l.turn("/spacedock:review write a three-line plan for a hello-world command")
	l.wantMenus()
	wrote, reviewed := -1, -1
	for i, tool := range l.tools {
		path, _ := tool.Input["file_path"].(string)
		cmd, _ := tool.Input["command"].(string)
		writes := tool.Name == "Write" && strings.HasSuffix(path, ".md") ||
			tool.Name == "Bash" && strings.Contains(cmd, ".md") && strings.Contains(cmd, ">")
		if writes && wrote < 0 {
			wrote = i
		}
		if skill, _ := tool.Input["skill"].(string); tool.Name == "Skill" && skill == "subspace:r" && reviewed < 0 {
			reviewed = i
		}
	}
	if wrote < 0 || reviewed < 0 || wrote > reviewed {
		t.Errorf("want a Markdown write before the subspace:r hand-off, got write=%d review=%d: %s", wrote, reviewed, l.flat())
	}
	if strings.Contains(l.flat(), "review-mode set") || strings.Contains(l.flat(), "review-mode clear") {
		t.Errorf("changed the recorded choice although one is recorded: %s", l.flat())
	}
	if out, _ := exec.Command("git", "-C", l.proj, "status", "--porcelain").Output(); strings.TrimSpace(string(out)) != "?? plan.md" {
		t.Errorf("the project gained a file nobody asked for:\n%s", out)
	}
	l.noPublish()
}

func (l *reviewLive) choiceBytes() []byte {
	l.t.Helper()
	paths, err := filepath.Glob(filepath.Join(l.root, "state", "subspace", "review-mode", "*.json"))
	if err != nil {
		l.t.Fatal(err)
	}
	var snapshot []byte
	for _, path := range paths {
		body, err := os.ReadFile(path)
		if err != nil {
			l.t.Fatal(err)
		}
		snapshot = append(snapshot, body...)
	}
	return snapshot
}

func (l *reviewLive) wantUnchangedChoice(before []byte) {
	l.t.Helper()
	if !bytes.Equal(before, l.choiceBytes()) {
		l.t.Error("one-off review changed the saved preference bytes")
	}
	for _, banned := range []string{"review-mode set", "review-mode clear"} {
		if strings.Contains(l.flat(), banned) {
			l.t.Errorf("one-off review attempted %s: %s", banned, l.flat())
		}
	}
}

func (l *reviewLive) wantHandoff(skill, file string, flags ...string) {
	l.t.Helper()
	var calls []string
	for _, tool := range l.tools {
		if s, _ := tool.Input["skill"].(string); tool.Name == "Skill" && strings.HasPrefix(s, "subspace:") {
			if s != skill {
				l.t.Errorf("handed to %s, want %s", s, skill)
			}
			args, _ := tool.Input["args"].(string)
			calls = append(calls, args)
		}
		if tool.Name == "Write" || tool.Name == "Edit" {
			l.t.Errorf("file review generated or edited an artifact: %v", tool.Input)
		}
	}
	if len(calls) != 1 {
		l.t.Fatalf("want exactly one companion handoff, got %v: %s", calls, l.flat())
	}
	args := calls[0]
	if !strings.Contains(args, file) {
		l.t.Errorf("handoff args %q do not contain the original file %q", args, file)
	}
	if strings.Contains(file, " ") && !strings.Contains(args, `"`+file+`"`) && !strings.Contains(args, "'"+file+"'") && !strings.Contains(args, strings.ReplaceAll(file, " ", `\ `)) {
		l.t.Errorf("path with spaces was not kept as one locator: %q", args)
	}
	for _, flag := range flags {
		if !strings.Contains(args, flag) {
			l.t.Errorf("handoff args %q omit %q", args, flag)
		}
	}
	if strings.Contains(args, "--allow-question") && strings.Contains(args, "--mode binding") {
		l.t.Errorf("unsupported binding questions were launched: %q", args)
	}
	body, err := os.ReadFile(filepath.Join(l.proj, file))
	if err != nil || string(body) != "# Plan\n\nShip the thing.\n" {
		l.t.Errorf("original artifact changed: %v", err)
	}
}

func TestLiveReviewProductNameUsesSavedChoice(t *testing.T) {
	cases := []struct {
		name, saved, skill, flag string
		answers                  map[string]string
		menus                    []string
	}{
		{"saved-terminal", "local", "subspace:r", "--mode binding", nil, nil},
		{"nothing-saved", "", "subspace:r", "--mode binding", map[string]string{"choice": "Terminal"}, []string{"choice"}},
		{"saved-web", "remote", "subspace:remote-review", "decision", nil, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, run := range []string{"run-1", "run-2", "run-3"} {
				t.Run(run, func(t *testing.T) {
					l := newReviewLive(t, tc.answers)
					if tc.saved != "" {
						l.spacedock("review-mode", "set", tc.saved)
					}
					before := l.choiceBytes()
					l.turn("Send plan.md to me in Spacedock Review for a decision.")
					l.wantMenus(tc.menus...)
					l.wantHandoff(tc.skill, "plan.md", tc.flag)
					if tc.saved == "" {
						if l.mode() != "local" {
							t.Error("natural-language review did not save the selected Terminal choice")
						}
					} else {
						l.wantUnchangedChoice(before)
					}
					if tc.skill == "subspace:r" {
						l.noPublish()
					}
				})
			}
		})
	}
}

func TestLiveReviewKeywordForms(t *testing.T) {
	cases := []struct {
		name, request, saved, skill, flag string
		answers                           map[string]string
		menus                             []string
	}{
		{"terminal-unset", "plan.md terminal decision", "", "subspace:r", "--mode binding", nil, nil},
		{"terminal-matching", "plan.md terminal decision", "local", "subspace:r", "--mode binding", nil, nil},
		{"terminal-conflicting", "plan.md terminal decision", "remote", "subspace:r", "--mode binding", nil, nil},
		{"web-unset-yes", "plan.md web decision", "", "subspace:remote-review", "decision", map[string]string{"agree": "Yes"}, []string{"agree"}},
		{"web-conflicting-yes", "plan.md web decision", "local", "subspace:remote-review", "decision", map[string]string{"agree": "Yes"}, []string{"agree"}},
		{"web-matching", "plan.md web decision", "remote", "subspace:remote-review", "decision", nil, nil},
		{"decision-local", "plan.md decision", "local", "subspace:r", "--mode binding", nil, nil},
		{"decision-remote", "plan.md decision", "remote", "subspace:remote-review", "decision", nil, nil},
		{"decision-unset", "plan.md decision", "", "subspace:r", "--mode binding", map[string]string{"choice": "Terminal"}, []string{"choice"}},
		{"questions-unset", "plan.md terminal questions", "", "subspace:r", "--allow-question", nil, nil},
		{"questions-matching", "plan.md terminal questions", "local", "subspace:r", "--allow-question", nil, nil},
		{"questions-conflicting", "plan.md terminal questions", "remote", "subspace:r", "--allow-question", nil, nil},
		{"quoted-spaces", `"launch plan.md" terminal decision`, "remote", "subspace:r", "--mode binding", nil, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := newReviewLive(t, tc.answers)
			if tc.saved != "" {
				l.spacedock("review-mode", "set", tc.saved)
			}
			file := "plan.md"
			if tc.name == "quoted-spaces" {
				file = "launch plan.md"
				if err := os.WriteFile(filepath.Join(l.proj, file), []byte("# Plan\n\nShip the thing.\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			before := l.choiceBytes()
			l.turn("/spacedock:review " + tc.request)
			l.wantMenus(tc.menus...)
			if tc.flag == "" {
				l.wantHandoff(tc.skill, file)
			} else {
				l.wantHandoff(tc.skill, file, tc.flag)
			}
			if tc.name == "decision-unset" {
				if l.mode() != "local" {
					t.Error("omitted location did not save the selected Terminal choice")
				}
			} else {
				l.wantUnchangedChoice(before)
			}
			if tc.skill == "subspace:r" {
				l.noPublish()
			}
			if len(tc.menus) == 1 && tc.menus[0] == "agree" {
				q := l.asked[0]
				for _, term := range []string{"review.spacedock.md", "Anyone with the link", "seven days"} {
					if !strings.Contains(q.Question, term) {
						t.Errorf("upload agreement omitted %q: %s", term, q.Question)
					}
				}
				menu, handoff := -1, -1
				for i, tool := range l.tools {
					if tool.Name == "AskUserQuestion" && menu < 0 {
						menu = i
					}
					if skill, _ := tool.Input["skill"].(string); tool.Name == "Skill" && skill == "subspace:remote-review" {
						handoff = i
					}
				}
				if menu < 0 || handoff <= menu {
					t.Errorf("Web handoff preceded upload agreement: %s", l.flat())
				}
			}
		})
	}
}

func TestLiveReviewOneOffWebDeclined(t *testing.T) {
	for _, saved := range []string{"", "local"} {
		t.Run("saved-"+saved, func(t *testing.T) {
			l := newReviewLive(t, map[string]string{"agree": "No"})
			if saved != "" {
				l.spacedock("review-mode", "set", saved)
			}
			before := l.choiceBytes()
			l.turn("/spacedock:review plan.md web decision")
			l.wantMenus("agree")
			q := l.asked[0]
			for _, term := range []string{"review.spacedock.md", "Anyone with the link", "seven days"} {
				if !strings.Contains(q.Question, term) {
					t.Errorf("upload agreement omitted %q: %s", term, q.Question)
				}
			}
			l.wantUnchangedChoice(before)
			if strings.Contains(l.flat(), "subspace:") {
				t.Errorf("declined upload still handed off the file: %s", l.flat())
			}
			l.noPublish()
		})
	}
}

func TestLiveReviewUnsupportedKeywordCombinations(t *testing.T) {
	for _, suffix := range []string{"terminal decision questions", "terminal questions decision", "web questions", "web decision questions", "terminal maybe", "web", "terminal", "tui", "local", "webui", "remote", "questions"} {
		t.Run(strings.ReplaceAll(suffix, " ", "-"), func(t *testing.T) {
			l := newReviewLive(t, nil)
			l.spacedock("review-mode", "set", "remote")
			before := l.choiceBytes()
			l.turn("/spacedock:review plan.md " + suffix)
			l.wantMenus()
			l.wantUnchangedChoice(before)
			if strings.Contains(l.flat(), "subspace:") {
				t.Errorf("unsupported request was handed off: %s", l.flat())
			}
			for _, tool := range l.tools {
				if tool.Name == "Write" || tool.Name == "Edit" {
					t.Errorf("unsupported request generated an artifact: %v", tool.Input)
				}
			}
			if strings.Contains(suffix, "questions") && !strings.Contains(l.text, "questions") {
				t.Errorf("refusal dropped the questions intent: %s", l.text)
			}
			if !reviewCorrectionsSupported(l.text) {
				t.Errorf("refusal suggested an unsupported request or no correction: %s", l.text)
			}
			l.noPublish()
		})
	}
}

func TestLiveReviewUnrelatedSuffixRemainsAPrompt(t *testing.T) {
	l := newReviewLive(t, nil)
	l.spacedock("review-mode", "set", "local")
	l.turn("/spacedock:review plan.md summarize this file into summary.md")
	l.wantMenus()
	if body, err := os.ReadFile(filepath.Join(l.proj, "summary.md")); err != nil || len(strings.TrimSpace(string(body))) == 0 {
		t.Errorf("unrelated suffix was not carried out as a prompt: %v", err)
	}
	if !l.handedTo("subspace:r", "summary.md") {
		t.Errorf("prompt output was not reviewed: %s", l.flat())
	}
	l.noPublish()
}

func reviewCorrectionsSupported(text string) bool {
	text = strings.ReplaceAll(text, "*", "")
	requests := regexp.MustCompile("/spacedock:review[ \t]+(?:plan\\.md|\"plan\\.md\")[ \t]+([^`\r\n]+)").FindAllStringSubmatch(text, -1)
	if len(requests) == 0 {
		return false
	}
	for _, request := range requests {
		switch strings.Trim(request[1], " \t\"'.;:") {
		case "terminal decision", "web decision", "decision", "terminal questions":
		default:
			return false
		}
	}
	return true
}

func TestReviewRefusalCorrection(t *testing.T) {
	for _, form := range []string{"terminal decision", "web decision", "decision", "terminal questions"} {
		if !reviewCorrectionsSupported("Try `/spacedock:review plan.md " + form + "`.") {
			t.Errorf("supported correction rejected: %s", form)
		}
	}
	if !reviewCorrectionsSupported("Try `/spacedock:review plan.md terminal decision` or `/spacedock:review plan.md terminal questions`.") {
		t.Error("supported halves of decision with questions were rejected")
	}
	for _, text := range []string{
		"Try `/spacedock:review plan.md web`.",
		"Try `/spacedock:review plan.md terminal decision` or `/spacedock:review plan.md web`.",
		"Try `/spacedock:review plan.md terminal decision questions`.",
		"No corrected request.",
	} {
		if reviewCorrectionsSupported(text) {
			t.Errorf("unsupported or missing correction accepted: %s", text)
		}
	}
}
