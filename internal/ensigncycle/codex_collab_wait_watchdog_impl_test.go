package ensigncycle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type codexScenarioResult struct {
	finalMessage string
	jsonl        string
	artifactDir  string
	duration     time.Duration
	interactive  bool
	resident     bool
}

type codexWaitStallArm string

const (
	codexWaitStallRepeatedWait    codexWaitStallArm = "repeated-wait"
	codexWaitStallSilentAfterWait codexWaitStallArm = "silent-after-wait"
)

type codexCollabWaitStallError struct {
	scenario        string
	handle          string
	arm             codexWaitStallArm
	artifactDir     string
	durableProgress bool
}

func (e *codexCollabWaitStallError) Error() string {
	handle := e.handle
	if handle == "" {
		handle = "unknown"
	}
	return fmt.Sprintf("Codex foreground-wait watchdog typed stall: scenario=%s wait_handle=%s arm=%s durable_progress=%t artifacts=%s",
		e.scenario, handle, e.arm, e.durableProgress, e.artifactDir)
}

type codexDurableProgressProbe interface {
	changed() bool
}

type codexCollabWaitWatchdog struct {
	scenario    string
	artifactDir string
	probe       codexDurableProgressProbe

	active              bool
	handle              string
	waitStarted         time.Time
	lastMeaningfulEvent codexWaitStallArm
}

func newCodexCollabWaitWatchdog(scenario, artifactDir string, probe codexDurableProgressProbe) *codexCollabWaitWatchdog {
	return &codexCollabWaitWatchdog{
		scenario:    scenario,
		artifactDir: artifactDir,
		probe:       probe,
	}
}

func (w *codexCollabWaitWatchdog) observeLine(line string, now time.Time, waitBudget time.Duration) error {
	w.observeDurableProgress()

	if clears, ok := lineIsCodexWorkerCompletion(line); ok && clears {
		w.clearWait()
		return nil
	}

	var ev codexCollabItem
	if err := json.Unmarshal([]byte(line), &ev); err != nil {
		return nil
	}
	it := ev.Item
	if it.Type != "collab_tool_call" {
		return nil
	}

	tool := strings.ToLower(it.Tool)
	handle := codexWaitHandle(it.ReceiverThreadIDs)
	if isCodexWaitTool(tool) {
		if w.active && w.handle == handle && !w.waitStarted.IsZero() && now.Sub(w.waitStarted) >= waitBudget {
			return w.stall(codexWaitStallRepeatedWait, false)
		}
		if !w.active || w.handle != handle {
			w.waitStarted = now
		}
		w.active = true
		w.handle = handle
		w.lastMeaningfulEvent = codexWaitStallSilentAfterWait
		return nil
	}

	if isCodexScenarioProgressTool(tool) {
		w.clearWait()
	}
	return nil
}

func (w *codexCollabWaitWatchdog) silenceStall() error {
	if !w.active && w.lastMeaningfulEvent != codexWaitStallSilentAfterWait {
		return nil
	}
	return w.stall(codexWaitStallSilentAfterWait, false)
}

func (w *codexCollabWaitWatchdog) clearWait() {
	w.active = false
	w.handle = ""
	w.waitStarted = time.Time{}
	w.lastMeaningfulEvent = ""
}

func (w *codexCollabWaitWatchdog) probeChanged() bool {
	return w.probe != nil && w.probe.changed()
}

func (w *codexCollabWaitWatchdog) observeDurableProgress() bool {
	if !w.probeChanged() {
		return false
	}
	w.clearWait()
	return true
}

func (w *codexCollabWaitWatchdog) stall(arm codexWaitStallArm, durableProgress bool) error {
	return &codexCollabWaitStallError{
		scenario:        w.scenario,
		handle:          w.handle,
		arm:             arm,
		artifactDir:     w.artifactDir,
		durableProgress: durableProgress,
	}
}

func drainCodexToExitWithWaitWatchdog(watcher *streamWatcher, budget time.Duration, label string, watchdog *codexCollabWaitWatchdog) (string, error) {
	deadline := time.Now().Add(budget)
	for {
		before := len(watcher.transcript)
		_, drained := watcher.drainEntries()
		if drained > 0 {
			deadline = time.Now().Add(budget)
			for _, line := range watcher.transcript[before:] {
				if err := watchdog.observeLine(line, time.Now(), budget); err != nil {
					watcher.proc.kill()
					return watcher.fullTranscript(), err
				}
			}
		}
		if watchdog.observeDurableProgress() {
			deadline = time.Now().Add(budget)
			continue
		}
		if _, exited := watcher.proc.poll(); exited {
			watcher.drainEntries()
			return watcher.fullTranscript(), nil
		}
		if time.Now().After(deadline) {
			if err := watchdog.silenceStall(); err != nil {
				watcher.proc.kill()
				return watcher.fullTranscript(), err
			}
			watcher.proc.kill()
			return watcher.fullTranscript(), &stepTimeout{
				label: label,
				msg: fmt.Sprintf("%s made no stream progress within %s (no-progress quiet budget) — a hung stage; killed the subprocess.\nTranscript tail:\n%s",
					label, budget, watcher.transcriptTail()),
			}
		}
		time.Sleep(watcher.pollInterval)
	}
}

func isCodexWaitTool(tool string) bool {
	switch tool {
	case "wait", "wait_agent", "collab:wait":
		return true
	default:
		return false
	}
}

func isCodexScenarioProgressTool(tool string) bool {
	return tool != "" && !isCodexWaitTool(tool)
}

func codexWaitHandle(threadIDs []string) string {
	if len(threadIDs) == 0 {
		return "unknown"
	}
	ids := append([]string(nil), threadIDs...)
	sort.Strings(ids)
	return strings.Join(ids, ",")
}

func lineIsCodexWorkerCompletion(line string) (bool, bool) {
	var sys struct {
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		Status  string `json:"status"`
	}
	if err := json.Unmarshal([]byte(line), &sys); err == nil &&
		sys.Type == "system" && sys.Subtype == "task_notification" && sys.Status == "completed" {
		return true, true
	}
	var ev struct {
		Type string `json:"type"`
		Item struct {
			Type         string `json:"type"`
			Tool         string `json:"tool"`
			Status       string `json:"status"`
			AgentsStates map[string]struct {
				Status string `json:"status"`
			} `json:"agents_states"`
		} `json:"item"`
	}
	if err := json.Unmarshal([]byte(line), &ev); err == nil &&
		ev.Type == "item.completed" &&
		ev.Item.Type == "collab_tool_call" &&
		isCodexWaitTool(strings.ToLower(ev.Item.Tool)) &&
		ev.Item.Status == "completed" {
		for _, state := range ev.Item.AgentsStates {
			if state.Status == "completed" {
				return true, true
			}
		}
	}
	return false, false
}

type workflowStateProbe struct {
	root string
	last string
}

func newWorkflowStateProbe(root string) (*workflowStateProbe, error) {
	fp, err := workflowStateFingerprint(root)
	if err != nil {
		return nil, err
	}
	return &workflowStateProbe{root: root, last: fp}, nil
}

func (p *workflowStateProbe) changed() bool {
	fp, err := workflowStateFingerprint(p.root)
	if err != nil || fp == p.last {
		return false
	}
	p.last = fp
	return true
}

func workflowStateFingerprint(root string) (string, error) {
	var paths []string
	if err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		paths = append(paths, rel)
		return nil
	}); err != nil {
		return "", err
	}
	sort.Strings(paths)

	h := sha256.New()
	for _, rel := range paths {
		b, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			return "", err
		}
		h.Write([]byte(rel))
		h.Write([]byte{0})
		h.Write(b)
		h.Write([]byte{0})
	}
	for _, args := range [][]string{{"rev-parse", "HEAD"}, {"status", "--short"}} {
		if out, ok := gitOutput(root, args...); ok {
			h.Write([]byte(strings.Join(args, " ")))
			h.Write([]byte{0})
			h.Write([]byte(out))
			h.Write([]byte{0})
		}
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func gitOutput(root string, args ...string) (string, bool) {
	full := append([]string{"-C", root}, args...)
	out, err := exec.Command("git", full...).CombinedOutput()
	if err != nil {
		return "", false
	}
	return string(out), true
}

func codexAttemptArtifactDir(root, scenario string, attempt int) string {
	if attempt > 0 {
		return filepath.Join(root, scenario, fmt.Sprintf("attempt-%d", attempt))
	}
	return filepath.Join(root, scenario)
}

type codexRejectionFlowAttempt struct {
	entityAfter string
	result      codexScenarioResult
}

type codexRejectionFlowAttemptFunc func(attempt int) (codexRejectionFlowAttempt, error)

func runCodexRejectionFlowWithRetry(runAttempt codexRejectionFlowAttemptFunc) (codexRejectionFlowAttempt, error) {
	var last codexRejectionFlowAttempt
	for attempt := 1; attempt <= 2; attempt++ {
		got, err := runAttempt(attempt)
		last = got
		if err != nil {
			if isCodexCollabWaitStall(err) && attempt == 1 {
				continue
			}
			return got, err
		}
		if err := assertRejectionFlow(got.entityAfter, got.result.finalMessage+"\n"+got.result.jsonl); err != nil {
			return got, err
		}
		if err := assertCodexReviewerReuseWithDurableState(got.result.jsonl, got.entityAfter); err != nil {
			return got, err
		}
		return got, nil
	}
	return last, fmt.Errorf("codex rejection-flow retry loop exhausted without a result")
}

func isCodexCollabWaitStall(err error) bool {
	var stall *codexCollabWaitStallError
	return errors.As(err, &stall)
}
