//go:build live

// ABOUTME: Live check of /spacedock:review's choice flow: runs the real skill in claude -p
// ABOUTME: against the built binary and stub programs, and grades the tool calls and the state file.
package integration

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// PATH holds no subspace-tui, so a local hand-off stops at the launcher instead of opening a
// terminal window on the machine running the test.
// Needs SPACEDOCK_BIN (the built spacedock) and SPACEDOCK_LIVE_CLAUDE_CONFIG_DIR (a logged-in
// Claude Code config). SPACEDOCK_LIVE_OUT, when set, receives each turn's raw stream.

type reviewLive struct {
	t       *testing.T
	root    string
	proj    string
	env     []string
	relay   string
	session string
	turns   int
	tools   []liveTool
	text    string
}

type liveTool struct {
	Name  string
	Input map[string]any
}

func (l *reviewLive) flat() string {
	b, _ := json.Marshal(l.tools)
	return string(b)
}

func newReviewLive(t *testing.T) *reviewLive {
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
	return &reviewLive{t: t, root: root, proj: proj, relay: relayLog, env: []string{
		"HOME=" + os.Getenv("HOME"), "USER=" + os.Getenv("USER"), "LOGNAME=" + os.Getenv("USER"),
		"CLAUDE_CONFIG_DIR=" + cfg, "XDG_STATE_HOME=" + filepath.Join(root, "state"),
		"XDG_DATA_HOME=" + filepath.Join(root, "data"), "TERM=dumb", "PATH=" + progs + ":/usr/bin:/bin",
	}}
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

// turn sends one prompt (resuming the previous turn's session after the first) and records the
// tool calls and the assistant text of that turn only.
func (l *reviewLive) turn(prompt string) {
	l.t.Helper()
	args := []string{"-p", prompt, "--plugin-dir", repoRoot(l.t), "--allowedTools", "Bash,Read,Write,Skill",
		"--output-format", "stream-json", "--verbose", "--model", "sonnet"}
	if l.session != "" {
		args = append(args, "--resume", l.session)
	}
	cmd := exec.Command(filepath.Join(l.root, "bin", "claude"), args...)
	cmd.Dir, cmd.Env = l.proj, l.env
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		l.t.Fatalf("claude turn %q: %v\n%s", prompt, err, out.String())
	}
	l.turns++
	if dir := os.Getenv("SPACEDOCK_LIVE_OUT"); dir != "" {
		name := filepath.Join(dir, strings.ReplaceAll(l.t.Name(), "/", "_")+"-turn"+string(rune('0'+l.turns))+".jsonl")
		_ = os.WriteFile(name, out.Bytes(), 0o644)
	}
	l.tools, l.text = nil, ""
	sc := bufio.NewScanner(&out)
	sc.Buffer(make([]byte, 1<<20), 1<<26)
	for sc.Scan() {
		var ev struct {
			Type      string `json:"type"`
			SessionID string `json:"session_id"`
			IsError   bool   `json:"is_error"`
			Message   struct {
				Content []struct {
					Type  string         `json:"type"`
					Text  string         `json:"text"`
					Name  string         `json:"name"`
					Input map[string]any `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal(sc.Bytes(), &ev) != nil {
			continue
		}
		if ev.SessionID != "" {
			l.session = ev.SessionID
		}
		if ev.Type == "result" && ev.IsError {
			l.t.Fatalf("claude turn %q ended in error:\n%s", prompt, sc.Text())
		}
		if ev.Type != "assistant" {
			continue
		}
		for _, c := range ev.Message.Content {
			switch c.Type {
			case "text":
				l.text += c.Text + "\n"
			case "tool_use":
				l.tools = append(l.tools, liveTool{c.Name, c.Input})
			}
		}
	}
}

func (l *reviewLive) said(want string) bool {
	return strings.Contains(strings.NewReplacer("*", "", "`", "").Replace(l.text), want)
}

func (l *reviewLive) noPublish() {
	l.t.Helper()
	calls, _ := os.ReadFile(l.relay)
	if strings.Contains(string(calls), "publish") {
		l.t.Errorf("subspace-relay-client was asked to publish: %q", calls)
	}
}

func TestLiveReviewNoArgumentShowsTheChoiceAndPublishesNothing(t *testing.T) {
	l := newReviewLive(t)
	l.spacedock("review-mode", "set", "local")
	l.turn("/spacedock:review")
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

func TestLiveReviewWebNeedsASecondAgreementBeforeRecording(t *testing.T) {
	l := newReviewLive(t)
	l.turn("/spacedock:review")
	if !l.said("Terminal (TUI)") || !l.said("Web") {
		t.Fatalf("first question does not offer both views:\n%s", l.text)
	}
	l.turn("Web")
	for _, want := range []string{"review.spacedock.md", "seven days"} {
		if !l.said(want) {
			t.Errorf("consent question lacks %q:\n%s", want, l.text)
		}
	}
	if strings.Contains(l.flat(), "review-mode set") || l.mode() != "" {
		t.Fatalf("Web was recorded before the user agreed: %s", l.flat())
	}
	l.turn("No")
	if strings.Contains(l.flat(), "review-mode set remote") || l.mode() != "" {
		t.Errorf("declining the upload still recorded %q: %s", l.mode(), l.flat())
	}
	if !l.said("Terminal (TUI)") {
		t.Errorf("a decline did not offer the Terminal (TUI):\n%s", l.text)
	}
	if !l.said("/spacedock:review") {
		t.Errorf("the answer did not say the setting can be changed through /spacedock:review:\n%s", l.text)
	}
	l.noPublish()
}

func TestLiveReviewSwitchingToWebAsksForAgreementAndDecliningKeepsTheChoice(t *testing.T) {
	l := newReviewLive(t)
	l.spacedock("review-mode", "set", "local")
	l.turn("/spacedock:review")
	l.turn("Switch to Web")
	if !l.said("seven days") {
		t.Errorf("switching to Web did not ask the upload question:\n%s", l.text)
	}
	if strings.Contains(l.flat(), "review-mode set") || l.mode() != "local" {
		t.Fatalf("the choice changed before the user agreed: %s", l.flat())
	}
	l.turn("No")
	if l.mode() != "local" {
		t.Errorf("declining changed the choice to %q", l.mode())
	}
	l.noPublish()
}

func TestLiveReviewNonFileArgumentWritesTheArtifactThenReviewsIt(t *testing.T) {
	l := newReviewLive(t)
	l.spacedock("review-mode", "set", "local")
	l.turn("/spacedock:review write a three-line plan for a hello-world command")
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

func TestLiveReviewFileArgumentHandsOffTheFileUnchanged(t *testing.T) {
	l := newReviewLive(t)
	l.spacedock("review-mode", "set", "local")
	l.turn("/spacedock:review plan.md")
	handed := false
	for _, tool := range l.tools {
		if skill, _ := tool.Input["skill"].(string); tool.Name == "Skill" && skill == "subspace:r" {
			handed = strings.Contains(l.flat(), "plan.md")
		}
		if tool.Name == "Write" {
			t.Errorf("a file argument wrote a file: %v", tool.Input)
		}
	}
	if !handed {
		t.Errorf("plan.md was not handed to subspace:r: %s", l.flat())
	}
	l.noPublish()
}
