package ensigncycle

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type toggledProgressProbe struct {
	mu         sync.Mutex
	changedNow bool
}

func (p *toggledProgressProbe) markChanged() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.changedNow = true
}

func (p *toggledProgressProbe) changed() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.changedNow {
		return false
	}
	p.changedNow = false
	return true
}

func TestCodexCollabWaitWatchdogRepeatedWaitStall(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newTestWatcher(src, proc)
	probe := &toggledProgressProbe{}
	watchdog := newCodexCollabWaitWatchdog("rejection-flow", "/tmp/rejection-flow/attempt-1", probe)

	src.push(codexCollabToolLine("item.started", "wait_agent", "thread-validation-1", ""))
	done := make(chan struct{})
	go func() {
		defer close(done)
		deadline := time.Now().Add(2 * w.quietBudget)
		for time.Now().Before(deadline) {
			src.push(codexCollabToolLine("item.started", "wait_agent", "thread-validation-1", ""))
			time.Sleep(w.pollInterval)
		}
	}()

	jsonl, err := drainCodexToExitWithWaitWatchdog(w, w.quietBudget, "codex shared scenario rejection-flow", watchdog)
	<-done

	stall := requireCodexWaitStall(t, err)
	if stall.scenario != "rejection-flow" {
		t.Fatalf("stall.scenario = %q, want rejection-flow", stall.scenario)
	}
	if stall.handle != "thread-validation-1" {
		t.Fatalf("stall.handle = %q, want thread-validation-1", stall.handle)
	}
	if stall.arm != codexWaitStallRepeatedWait {
		t.Fatalf("stall.arm = %q, want %q", stall.arm, codexWaitStallRepeatedWait)
	}
	if !proc.wasKilled() {
		t.Fatal("watchdog stall must kill the fake Codex proc")
	}
	for _, want := range []string{"rejection-flow", "thread-validation-1", string(codexWaitStallRepeatedWait), "/tmp/rejection-flow/attempt-1"} {
		if !strings.Contains(stall.Error(), want) {
			t.Fatalf("stall diagnostic %q is missing %q", stall.Error(), want)
		}
	}
	if !strings.Contains(jsonl, "wait_agent") {
		t.Fatalf("returned JSONL should preserve the wait transcript, got %q", jsonl)
	}
}

func TestCodexCollabWaitWatchdogSilentAfterWaitStall(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newTestWatcher(src, proc)
	watchdog := newCodexCollabWaitWatchdog("rejection-flow", "/tmp/rejection-flow/attempt-1", &toggledProgressProbe{})

	src.push(codexCollabToolLine("item.started", "collab:wait", "thread-validation-2", ""))

	_, err := drainCodexToExitWithWaitWatchdog(w, w.quietBudget, "codex shared scenario rejection-flow", watchdog)

	stall := requireCodexWaitStall(t, err)
	if stall.arm != codexWaitStallSilentAfterWait {
		t.Fatalf("stall.arm = %q, want %q", stall.arm, codexWaitStallSilentAfterWait)
	}
	if stall.handle != "thread-validation-2" {
		t.Fatalf("stall.handle = %q, want thread-validation-2", stall.handle)
	}
	var timeout *stepTimeout
	if errors.As(err, &timeout) {
		t.Fatalf("silent-after-wait must be typed as codexCollabWaitStallError, not generic stepTimeout: %v", err)
	}
	if !proc.wasKilled() {
		t.Fatal("silent-after-wait stall must kill the fake Codex proc")
	}
}

func TestCodexCollabWaitWatchdogSilentAfterWaitWithDurableProgressStillTyped(t *testing.T) {
	src := &fakeLineSource{}
	proc := &fakeProc{}
	w := newTestWatcher(src, proc)
	probe := &toggledProgressProbe{}
	watchdog := newCodexCollabWaitWatchdog("rejection-flow", "/tmp/rejection-flow/attempt-1", probe)

	src.push(codexCollabToolLine("item.started", "wait", "thread-validation-6", ""))
	go func() {
		time.Sleep(2 * w.pollInterval)
		probe.markChanged()
	}()

	_, err := drainCodexToExitWithWaitWatchdog(w, w.quietBudget, "codex shared scenario rejection-flow", watchdog)

	stall := requireCodexWaitStall(t, err)
	if stall.arm != codexWaitStallSilentAfterWait {
		t.Fatalf("stall.arm = %q, want %q", stall.arm, codexWaitStallSilentAfterWait)
	}
	if !stall.durableProgress {
		t.Fatal("stall should record that durable progress was observed during the silent wait")
	}
	if !strings.Contains(stall.Error(), "durable_progress=true") {
		t.Fatalf("stall diagnostic should distinguish durable progress during the wait: %q", stall.Error())
	}
	if !proc.wasKilled() {
		t.Fatal("silent wait with no parent progress must still kill the fake Codex proc")
	}
}

func TestCodexCollabWaitWatchdogPositiveControls(t *testing.T) {
	tests := []struct {
		name  string
		drive func(src *fakeLineSource, proc *fakeProc, probe *toggledProgressProbe, w *streamWatcher)
	}{
		{
			name: "durable entity progress clears the active wait",
			drive: func(src *fakeLineSource, proc *fakeProc, probe *toggledProgressProbe, w *streamWatcher) {
				src.push(codexCollabToolLine("item.started", "wait_agent", "thread-validation-3", ""))
				time.Sleep(2 * w.pollInterval)
				probe.markChanged()
				src.push(codexCollabToolLine("item.started", "wait_agent", "thread-validation-3", ""))
				time.Sleep(2 * w.pollInterval)
				proc.setExited(0)
			},
		},
		{
			name: "completed foreground wait clears the active wait",
			drive: func(src *fakeLineSource, proc *fakeProc, probe *toggledProgressProbe, w *streamWatcher) {
				src.push(codexCollabToolLine("item.started", "wait", "thread-validation-4", ""))
				time.Sleep(w.quietBudget/2 + 2*w.pollInterval)
				src.push(codexCollabWaitCompletedLine("thread-validation-4", "Done: validation completed."))
				time.Sleep(w.quietBudget/2 + 3*w.pollInterval)
				src.push(codexCollabToolLine("item.started", "wait", "thread-validation-4", ""))
				time.Sleep(2 * w.pollInterval)
				proc.setExited(0)
			},
		},
		{
			name: "non-wait scenario progress clears the active wait",
			drive: func(src *fakeLineSource, proc *fakeProc, probe *toggledProgressProbe, w *streamWatcher) {
				src.push(codexCollabToolLine("item.started", "wait_agent", "thread-validation-5", ""))
				time.Sleep(2 * w.pollInterval)
				src.push(codexCollabToolLine("item.started", "send_input", "thread-validation-5", "Re-run validation for cycle 2."))
				time.Sleep(2 * w.pollInterval)
				proc.setExited(0)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			src := &fakeLineSource{}
			proc := &fakeProc{}
			w := newTestWatcher(src, proc)
			probe := &toggledProgressProbe{}
			watchdog := newCodexCollabWaitWatchdog("rejection-flow", "/tmp/rejection-flow/attempt-1", probe)

			done := make(chan struct{})
			go func() {
				defer close(done)
				tc.drive(src, proc, probe, w)
			}()

			_, err := drainCodexToExitWithWaitWatchdog(w, w.quietBudget, "codex shared scenario rejection-flow", watchdog)
			<-done
			if err != nil {
				t.Fatalf("watchdog must not trip when external progress clears the wait: %v", err)
			}
			if proc.wasKilled() {
				t.Fatal("positive control must not kill the fake Codex proc")
			}
		})
	}
}

func TestRunCodexRejectionFlowRetryAfterTypedWatchdog(t *testing.T) {
	artifactRoot := t.TempDir()
	var attempts []int

	result, err := runCodexRejectionFlowWithRetry(func(attempt int) (codexRejectionFlowAttempt, error) {
		attempts = append(attempts, attempt)
		artifactDir := codexAttemptArtifactDir(artifactRoot, "rejection-flow", attempt)
		if attempt == 1 {
			return codexRejectionFlowAttempt{}, &codexCollabWaitStallError{
				scenario:    "rejection-flow",
				handle:      "thread-validation-1",
				arm:         codexWaitStallSilentAfterWait,
				artifactDir: artifactDir,
			}
		}
		return codexRejectionFlowAttempt{
			entityAfter: passingRejectionEntity(),
			result: codexScenarioResult{
				finalMessage: "first-cycle rejection routed back to implementation; second-cycle re-validation PASSED",
				jsonl:        codexReviewerReuseJSONL("thread-validation-1", "thread-implementation-1"),
				artifactDir:  artifactDir,
			},
		}, nil
	})
	if err != nil {
		t.Fatalf("typed watchdog stall should retry once and return the clean attempt: %v", err)
	}
	if got := fmt.Sprint(attempts); got != "[1 2]" {
		t.Fatalf("attempts = %s, want [1 2]", got)
	}
	if !strings.HasSuffix(result.result.artifactDir, filepath.Join("rejection-flow", "attempt-2")) {
		t.Fatalf("result artifact dir = %q, want attempt-2", result.result.artifactDir)
	}
}

func TestRunCodexRejectionFlowRetryIsNarrow(t *testing.T) {
	t.Run("non-watchdog errors stop after one attempt", func(t *testing.T) {
		attempts := 0
		_, err := runCodexRejectionFlowWithRetry(func(attempt int) (codexRejectionFlowAttempt, error) {
			attempts++
			return codexRejectionFlowAttempt{}, errors.New("codex auth setup failed")
		})
		if err == nil {
			t.Fatal("expected non-watchdog error")
		}
		if attempts != 1 {
			t.Fatalf("attempts = %d, want 1", attempts)
		}
	})

	t.Run("assertion failures stop after one attempt", func(t *testing.T) {
		attempts := 0
		_, err := runCodexRejectionFlowWithRetry(func(attempt int) (codexRejectionFlowAttempt, error) {
			attempts++
			return codexRejectionFlowAttempt{
				entityAfter: rejectionEntity(),
				result: codexScenarioResult{
					finalMessage: "the transcript mentions a rejection and implementation follow-up",
					jsonl:        codexReviewerReuseJSONL("thread-validation-1", "thread-implementation-1"),
					artifactDir:  codexAttemptArtifactDir(t.TempDir(), "rejection-flow", attempt),
				},
			}, nil
		})
		if err == nil {
			t.Fatal("expected durable-state assertion failure")
		}
		if attempts != 1 {
			t.Fatalf("attempts = %d, want 1", attempts)
		}
	})
}

func TestRejectionFlowRejectsWaitReuseTranscriptWithoutDurableSecondCycle(t *testing.T) {
	entity := rejectionFixMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: initial\n\n" +
		"## Stage Report: implementation\n\n- DONE: applied fix\n\n" +
		"### Feedback Cycles\n\n- Cycle 1: REJECTED\n"
	jsonl := codexCollabToolLine("item.started", "wait_agent", "thread-validation-1", "") + "\n" +
		codexReviewerReuseJSONL("thread-validation-1", "thread-implementation-1")

	if err := assertCodexReviewerReuse(jsonl); err != nil {
		t.Fatalf("test fixture should contain a real Codex reuse signal: %v", err)
	}
	if err := assertRejectionFlow(entity, "validation rejected; implementation follow-up\n"+jsonl); err == nil {
		t.Fatal("durable rejection-flow assertion must fail when the entity lacks the second feedback cycle")
	}
}

func TestCodexForegroundWaitWatchdogDocsContract(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	readmePath := filepath.Join(wd, "..", "..", "docs", "dev", "README.md")
	b, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("read docs/dev/README.md: %v", err)
	}
	doc := string(b)

	mustContain := []string{
		"Codex foreground-wait watchdog",
		"`collab:wait` / `wait_agent`",
		"durable workflow-state progress",
		"typed stall",
	}
	for _, clause := range mustContain {
		if !strings.Contains(doc, clause) {
			t.Errorf("docs/dev/README.md is missing the required Codex foreground-wait watchdog clause: %q", clause)
		}
	}
}

func requireCodexWaitStall(t *testing.T, err error) *codexCollabWaitStallError {
	t.Helper()
	var stall *codexCollabWaitStallError
	if !errors.As(err, &stall) {
		t.Fatalf("want codexCollabWaitStallError, got %T: %v", err, err)
	}
	return stall
}

func codexCollabToolLine(eventType, tool, threadID, prompt string) string {
	return fmt.Sprintf(`{"type":%q,"item":{"type":"collab_tool_call","tool":%q,"receiver_thread_ids":[%q],"prompt":%q}}`,
		eventType, tool, threadID, prompt)
}

func codexCollabWaitCompletedLine(threadID, message string) string {
	return fmt.Sprintf(`{"type":"item.completed","item":{"type":"collab_tool_call","tool":"wait","receiver_thread_ids":[%q],"agents_states":{%q:{"status":"completed","message":%q}},"status":"completed"}}`,
		threadID, threadID, message)
}

func codexReviewerReuseJSONL(validationThread, implementationThread string) string {
	return strings.Join([]string{
		codexCollabToolLine("item.completed", "spawn_agent", validationThread, "Read /tmp/spacedock-dispatch/spacedock-ensign-rejection-task-validation.md and treat its content as your assignment."),
		codexCollabToolLine("item.completed", "spawn_agent", implementationThread, "Read /tmp/spacedock-dispatch/spacedock-ensign-rejection-task-implementation.md and treat its content as your assignment."),
		codexCollabToolLine("item.started", "send_input", implementationThread, "Feedback routed from validation to implementation for rejection-task. The fix marker is absent."),
		codexCollabToolLine("item.started", "send_input", validationThread, "Re-run validation for rejection-task as cycle 2 using your existing validation reviewer context."),
	}, "\n")
}

func passingRejectionEntity() string {
	return "---\nstatus: validation\n---\n" +
		rejectionFixMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: Initial implementation\n\n" +
		"## Stage Report: implementation\n\n- DONE: Applied rejection fix\n\n" +
		"### Feedback Cycles\n\n- Cycle 1: REJECTED\n- Cycle 2: PASSED\n"
}
