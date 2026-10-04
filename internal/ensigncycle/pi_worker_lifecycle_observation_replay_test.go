package ensigncycle

import (
	"strings"
	"testing"
)

// Pi observes a delegated worker's completion through surfaces the assert did
// not credit: the captured root transcript from run 37101046846 (preserved under
// docs/dev/.spacedock-state/pi-native-completion-evidence/) spawns an async
// `subagent`, waits with `bg_wait`, and receives a native `subagent-notify` —
// never the `subagent_wait` result or `State: complete` status it keyed on, so it
// computed completed=-1 for a dispatch that succeeded. These lines reproduce that
// shape deterministically and carry ONLY the two new surfaces, so removing the new
// credit turns the positive case RED.
const (
	piValidationSpawn = `{"type":"message","message":{"role":"assistant","content":[{"type":"toolCall","id":"call-spawn","name":"subagent","arguments":{"agent":"worker","task":"Read /tmp/spacedock-dispatch/auto-continue-task-validation.md and treat its content as your assignment.","async":true}}]}}`
	piSpawnResult     = `{"type":"message","message":{"role":"toolResult","toolCallId":"call-spawn","toolName":"subagent","details":{"runId":"run-ac0f4b70"}}}`
	piBgWaitCall      = `{"type":"message","message":{"role":"assistant","content":[{"type":"toolCall","id":"call-wait","name":"bg_wait","arguments":{"id":"run-ac0f4b70","timeoutMs":1800000}}]}}`
	piBgWaitResult    = `{"type":"message","message":{"role":"toolResult","toolCallId":"call-wait","toolName":"bg_wait","content":[{"type":"text","text":"Waited 1m34s for run \"run-ac0f4b70\"; done. Outcome: 1 complete. Completion/control events have been observed; inspect status if a notification is not visible yet."}]}}`
	piNativeNotice    = `{"type":"custom_message","customType":"subagent-notify","content":"Background task completed: **worker**\n\nworker:\nValidated AC-1 and appended a PASSED validation report.\n\nSession file: /tmp/session.jsonl"}`
	piGatePrepare     = `{"type":"message","message":{"role":"assistant","content":[{"type":"toolCall","id":"call-gate","name":"bash","arguments":{"command":"spacedock gate prepare auto-continue-task"}}]}}`
)

func piWorkerLifecycleStream(parts ...string) string {
	return strings.Join(parts, "\n")
}

// TestPiWorkerLifecycleObservationReplay feeds the captured Pi transcript shape
// through the shared assert and proves both new completion surfaces are load
// bearing: either surface alone grades the run GREEN, removing both turns it RED,
// a missing dispatch still fails, and a completion after the validation
// transition still violates the ordering contract.
func TestPiWorkerLifecycleObservationReplay(t *testing.T) {
	entity := "---\nstatus: validation\n---\n# Task\n\n## Stage Report: validation\n\n- DONE: verify\n"
	captured := []string{piValidationSpawn, piSpawnResult, piBgWaitCall, piBgWaitResult, piNativeNotice, piGatePrepare}

	// Both surfaces together: the positive case. Before the new credit this is the
	// captured transcript that yielded completed=-1.
	if err := assertWorkerLifecycle(piWorkerLifecycleStream(captured...), entity, "validation", "gate prepare"); err != nil {
		t.Fatalf("captured Pi bg_wait + native notice lifecycle rejected: %v", err)
	}

	// Each new surface is independently sufficient.
	if err := assertWorkerLifecycle(piWorkerLifecycleStream(piValidationSpawn, piSpawnResult, piBgWaitCall, piBgWaitResult, piGatePrepare), entity, "validation", "gate prepare"); err != nil {
		t.Fatalf("bg_wait completion alone rejected: %v", err)
	}
	if err := assertWorkerLifecycle(piWorkerLifecycleStream(piValidationSpawn, piSpawnResult, piBgWaitCall, piNativeNotice, piGatePrepare), entity, "validation", "gate prepare"); err != nil {
		t.Fatalf("native completion notice alone rejected: %v", err)
	}

	// Negative control: remove both new surfaces and the run must RED. This is the
	// falsifying edit for AC-1 — drop the new credit and this case stops failing.
	noBgWait := strings.Replace(piBgWaitResult, "done. Outcome: 1 complete", "running. Outcome: 0 complete", 1)
	noNotice := strings.Replace(piNativeNotice, "Background task completed", "Background task failed", 1)
	if err := assertWorkerLifecycle(piWorkerLifecycleStream(piValidationSpawn, piSpawnResult, piBgWaitCall, noBgWait, noNotice, piGatePrepare), entity, "validation", "gate prepare"); err == nil {
		t.Fatal("a lifecycle with both new completion surfaces removed passed")
	}

	// AC-2 zero-spawn control: no dispatch for this stage, so the native notice and
	// the wait cannot credit it. Credit any transcript regardless of spawn count and
	// this case stops failing.
	noSpawn := strings.Replace(piValidationSpawn, "auto-continue-task-validation.md", "auto-continue-task-implementation.md", 1)
	if err := assertWorkerLifecycle(piWorkerLifecycleStream(noSpawn, piSpawnResult, piBgWaitCall, piBgWaitResult, piNativeNotice, piGatePrepare), entity, "validation", "gate prepare"); err == nil {
		t.Fatal("a Pi lifecycle with no validation-stage spawn passed")
	}

	// AC-3 ordering control: a completion that lands after the validation transition
	// must not grade GREEN. Drop the ordering check and this case stops failing.
	inverted := []string{piValidationSpawn, piSpawnResult, piGatePrepare, piBgWaitCall, piBgWaitResult, piNativeNotice}
	if err := assertWorkerLifecycle(piWorkerLifecycleStream(inverted...), entity, "validation", "gate prepare"); err == nil {
		t.Fatal("a Pi completion after the validation transition passed")
	}
}
