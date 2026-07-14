//go:build live

package ensigncycle

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

const codexInteractiveBootBudget = 8 * time.Minute

// runInteractiveGreet launches the shipped `spacedock codex` front door inside a
// real tmux TTY. It sends no follow-up prompt: the launcher supplies the normal
// first-officer bootstrap, the durable rollout supplies the TUI source marker and
// task-complete turn end, and the still-live tmux session proves greet-and-stop.
func (r codexLiveRunner) runInteractiveGreet(t *testing.T, scenario sharedRuntimeScenario, workflowRoot string) (codexScenarioResult, error) {
	return r.runInteractiveSession(t, scenario, workflowRoot, "")
}

// runInteractiveEngage exercises the real two-turn boundary: the launcher-owned
// greet reaches a committed idle first, then the captain sends only the literal
// interaction verb. The recovery state is discovered by engage, not disclosed in
// a headless scenario prompt that can cue eager owner reads before convergence.
func (r codexLiveRunner) runInteractiveEngage(t *testing.T, scenario sharedRuntimeScenario, workflowRoot string) (codexScenarioResult, error) {
	return r.runInteractiveSession(t, scenario, workflowRoot, "engage")
}

func (r codexLiveRunner) runInteractiveSession(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, followup string) (codexScenarioResult, error) {
	t.Helper()
	artifactDir := codexAttemptArtifactDir(r.artifactRoot, scenario.name+"-interactive", 0)
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	resolvedCwd := workflowRoot
	if resolved, err := filepath.EvalSymlinks(workflowRoot); err == nil {
		resolvedCwd = resolved
	}
	if err := seedInteractiveCodexTrust(r.codexHome, resolvedCwd); err != nil {
		return codexScenarioResult{artifactDir: artifactDir}, fmt.Errorf("seed interactive Codex trust: %w", err)
	}

	session := fmt.Sprintf("sdcodex-%s-%d", scenario.name, time.Now().UnixNano())
	proc := newTmuxLiveProc(session)
	defer proc.kill()

	// CODEX_THREAD_ID/CODEX_CI belong to a parent Codex process when the harness
	// itself runs under Codex. Removing them gives the child an independent TUI
	// session and durable rollout. The actual Spacedock front door remains the
	// launched program; --no-alt-screen only makes pane diagnostics legible.
	// Every isolated interactive live session runs with the same explicit
	// no-prompt/full-write posture as the headless runner. Use the compatible
	// component flags rather than the aggregate bypass flag: on an unsandboxed
	// front door the aggregate flag otherwise coexists with the launcher's
	// injected `on-request` mode and Codex rejects the conflicting pair.
	hostArgs := []string{"--ask-for-approval", "never", "--sandbox", "danger-full-access", "--no-alt-screen"}
	launch := shellJoin(append([]string{
		"env", "-u", "CODEX_THREAD_ID", "-u", "CODEX_CI",
		r.spacedockBin, "codex", "--skip-compat-check", "--",
	}, hostArgs...))
	args := []string{"new-session", "-d", "-s", session, "-x", "220", "-y", "50", "-c", workflowRoot}
	for _, kv := range r.env {
		args = append(args, "-e", kv)
	}
	args = append(args, launch)
	started := time.Now()
	if out, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		return codexScenarioResult{artifactDir: artifactDir}, fmt.Errorf("tmux interactive Codex launch failed: %w (%s)", err, out)
	}

	rolloutPath, err := waitForCodexInteractiveRollout(r.codexHome, resolvedCwd, started, proc, codexInteractiveBootBudget)
	if err != nil {
		pane := captureTmuxPane(session)
		_ = os.WriteFile(filepath.Join(artifactDir, "codex-pane-no-rollout.txt"), []byte(pane), 0o644)
		return codexScenarioResult{artifactDir: artifactDir}, fmt.Errorf("%w\nCodex pane:\n%s", err, pane)
	}
	if err := waitForCodexInteractiveIdle(rolloutPath, proc, codexInteractiveBootBudget); err != nil {
		pane := captureTmuxPane(session)
		_ = os.WriteFile(filepath.Join(artifactDir, "codex-pane-at-stall.txt"), []byte(pane), 0o644)
		return codexScenarioResult{artifactDir: artifactDir}, fmt.Errorf("%w\nCodex pane:\n%s", err, pane)
	}
	if followup != "" {
		stream, readErr := os.ReadFile(rolloutPath)
		if readErr != nil {
			return codexScenarioResult{artifactDir: artifactDir}, fmt.Errorf("read interactive Codex rollout before follow-up: %w", readErr)
		}
		completedBefore := codexRolloutTaskCompleteCount(string(stream))
		if err := sendCodexInteractiveInput(session, followup); err != nil {
			return codexScenarioResult{artifactDir: artifactDir}, fmt.Errorf("send interactive Codex follow-up: %w", err)
		}
		if err := waitForCodexInteractiveNextIdle(rolloutPath, completedBefore, proc, codexInteractiveBootBudget); err != nil {
			pane := captureTmuxPane(session)
			_ = os.WriteFile(filepath.Join(artifactDir, "codex-pane-at-followup-stall.txt"), []byte(pane), 0o644)
			return codexScenarioResult{artifactDir: artifactDir}, fmt.Errorf("%w\nCodex pane:\n%s", err, pane)
		}
	}

	streamBytes, err := os.ReadFile(rolloutPath)
	if err != nil {
		return codexScenarioResult{artifactDir: artifactDir}, fmt.Errorf("read interactive Codex rollout: %w", err)
	}
	stream := string(streamBytes)
	finalMessage, err := extractCodexInteractiveFinalMessage(stream)
	if err != nil {
		return codexScenarioResult{jsonl: stream, artifactDir: artifactDir}, err
	}
	_, exited := proc.poll()
	pane := captureTmuxPane(session)
	if err := os.WriteFile(filepath.Join(artifactDir, "codex-rollout.jsonl"), streamBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(artifactDir, "codex-pane-at-greet.txt"), []byte(pane), 0o644)
	_ = os.WriteFile(filepath.Join(artifactDir, "codex-final-message.txt"), []byte(finalMessage), 0o644)

	return codexScenarioResult{
		finalMessage: finalMessage,
		jsonl:        stream,
		artifactDir:  artifactDir,
		duration:     time.Since(started),
		interactive:  codexRolloutIsInteractive(stream),
		resident:     !exited,
	}, nil
}

// seedInteractiveCodexTrust pre-accepts the TUI's per-worktree trust picker in
// the isolated live-test config. The fixture is generated by the test itself and
// its path is unique, so this does not weaken or mutate operator trust state.
func seedInteractiveCodexTrust(codexHome, resolvedCwd string) error {
	path := filepath.Join(codexHome, "config.toml")
	body, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	header := "[projects." + strconv.Quote(resolvedCwd) + "]"
	if strings.Contains(string(body), header) {
		return nil
	}
	addition := "\n\n" + header + "\ntrust_level = \"trusted\"\n"
	return os.WriteFile(path, append(body, addition...), 0o600)
}

func waitForCodexInteractiveRollout(codexHome, cwd string, notBefore time.Time, proc procPoller, budget time.Duration) (string, error) {
	deadline := time.Now().Add(budget)
	for {
		if path := codexInteractiveRolloutForCWD(filepath.Join(codexHome, "sessions"), cwd, notBefore); path != "" {
			return path, nil
		}
		if _, exited := proc.poll(); exited {
			return "", fmt.Errorf("interactive Codex process exited before a TUI rollout appeared for %s", cwd)
		}
		if time.Now().After(deadline) {
			return "", fmt.Errorf("timed out after %s waiting for an interactive Codex rollout for %s", budget, cwd)
		}
		time.Sleep(ptyPollInterval)
	}
}

func codexInteractiveRolloutForCWD(root, cwd string, notBefore time.Time) string {
	var newest string
	var newestTime time.Time
	_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || filepath.Ext(path) != ".jsonl" {
			return nil
		}
		info, err := entry.Info()
		if err != nil || info.ModTime().Before(notBefore.Add(-time.Second)) || info.ModTime().Before(newestTime) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		for _, line := range strings.Split(string(data), "\n") {
			var record codexRolloutRecord
			if json.Unmarshal([]byte(line), &record) != nil || record.Type != "session_meta" {
				continue
			}
			if filepath.Clean(record.Payload.Cwd) == filepath.Clean(cwd) && record.Payload.Source == "cli" && record.Payload.Originator == "codex-tui" {
				newest = path
				newestTime = info.ModTime()
			}
			break
		}
		return nil
	})
	return newest
}

func waitForCodexInteractiveIdle(path string, proc procPoller, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	consecutive := 0
	for {
		if data, err := os.ReadFile(path); err == nil && codexRolloutReachedIdle(string(data)) {
			consecutive++
			if consecutive >= ptyIdleStablePolls {
				return nil
			}
		} else {
			consecutive = 0
		}
		if _, exited := proc.poll(); exited {
			return fmt.Errorf("interactive Codex process exited before reaching a committed greet turn")
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("no completed interactive Codex greet within %s", budget)
		}
		time.Sleep(ptyPollInterval)
	}
}

func codexRolloutTaskCompleteCount(stream string) int {
	count := 0
	for _, record := range parseCodexRollout(stream) {
		if record.Type == "event_msg" && record.Payload.Type == "task_complete" {
			count++
		}
	}
	return count
}

func waitForCodexInteractiveNextIdle(path string, completedBefore int, proc procPoller, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	consecutive := 0
	for {
		if data, err := os.ReadFile(path); err == nil && codexRolloutTaskCompleteCount(string(data)) > completedBefore && codexRolloutReachedIdle(string(data)) {
			consecutive++
			if consecutive >= ptyIdleStablePolls {
				return nil
			}
		} else {
			consecutive = 0
		}
		if _, exited := proc.poll(); exited {
			return fmt.Errorf("interactive Codex process exited before completing its follow-up turn")
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("no completed interactive Codex follow-up within %s", budget)
		}
		time.Sleep(ptyPollInterval)
	}
}

func sendCodexInteractiveInput(session, input string) error {
	if out, err := exec.Command("tmux", "send-keys", "-t", session, "-l", input).CombinedOutput(); err != nil {
		return fmt.Errorf("tmux send-keys (literal): %w (%s)", err, out)
	}
	time.Sleep(300 * time.Millisecond)
	if out, err := exec.Command("tmux", "send-keys", "-t", session, "Enter").CombinedOutput(); err != nil {
		return fmt.Errorf("tmux send-keys (enter): %w (%s)", err, out)
	}
	return nil
}

func captureTmuxPane(session string) string {
	out, err := exec.Command("tmux", "capture-pane", "-p", "-t", session).CombinedOutput()
	if err != nil {
		return ""
	}
	return string(out)
}
