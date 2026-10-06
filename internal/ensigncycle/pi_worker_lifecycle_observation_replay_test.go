package ensigncycle

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// Pi observes a delegated worker's completion through surfaces the assert did
// not credit: the captured root transcript spawned an async `subagent`, waited
// with `bg_wait`, and received a native `subagent-notify` — never the
// `subagent_wait` result or `State: complete` status it keyed on — so it
// computed completed=-1 for a dispatch that succeeded.
//
// The captured baseline is testdata/pi_worker_lifecycle/auto-continue-validation-parent.jsonl,
// selected from the auto-continue-after-implementation validation parent session
// (2026-10-04); see the sibling provenance.json for the exact source path,
// SHA-256, selected line numbers and projection. The native notice names its own
// run in a `Retention-managed async directory: …/async-subagent-runs/<id>` field,
// which is what the correlation guard below matches against the spawn's run id.
func TestPiWorkerLifecycleObservationReplay(t *testing.T) {
	entity := "---\nstatus: validation\n---\n# Task\n\n## Stage Report: validation\n\n- DONE: verify\n"
	captured := strings.Split(strings.TrimRight(readFile(t, filepath.Join("testdata", "pi_worker_lifecycle", "auto-continue-validation-parent.jsonl")), "\n"), "\n")
	// Original record order: spawn, spawn result, bg_wait call, bg_wait result,
	// native notice, gate prepare.
	spawn, spawnResult, bgWaitCall, bgWaitResult, nativeNotice, gatePrepare := captured[0], captured[1], captured[2], captured[3], captured[4], captured[5]

	var spawned struct {
		Message struct{ Details struct{ RunID string } }
	}
	if err := json.Unmarshal([]byte(spawnResult), &spawned); err != nil || spawned.Message.Details.RunID == "" {
		t.Fatalf("captured spawn result carries no run id: %v (%s)", err, spawnResult)
	}
	runID := spawned.Message.Details.RunID

	stream := func(parts ...string) string { return strings.Join(parts, "\n") }
	// Both surfaces together: the captured baseline. Before the new credit this is
	// the captured transcript that yielded completed=-1.
	if err := assertWorkerLifecycle(stream(spawn, spawnResult, bgWaitCall, bgWaitResult, nativeNotice, gatePrepare), entity, "validation", "gate prepare"); err != nil {
		t.Fatalf("captured Pi bg_wait + native notice lifecycle rejected: %v", err)
	}

	// Each new surface is independently sufficient.
	if err := assertWorkerLifecycle(stream(spawn, spawnResult, bgWaitCall, bgWaitResult, gatePrepare), entity, "validation", "gate prepare"); err != nil {
		t.Fatalf("bg_wait completion alone rejected: %v", err)
	}
	if err := assertWorkerLifecycle(stream(spawn, spawnResult, bgWaitCall, nativeNotice, gatePrepare), entity, "validation", "gate prepare"); err != nil {
		t.Fatalf("native completion notice alone rejected: %v", err)
	}

	// AC-1 correlation control: a native notice that names a DIFFERENT run must not
	// credit the spawned run. Drop the run-id match and this case stops failing.
	unrelatedNotice := strings.Replace(nativeNotice, "async-subagent-runs/"+runID, "async-subagent-runs/0badc0de-0000-0000-0000-000000000000", 1)
	if unrelatedNotice == nativeNotice {
		t.Fatal("captured native notice does not name its async run id")
	}
	if err := assertWorkerLifecycle(stream(spawn, spawnResult, bgWaitCall, unrelatedNotice, gatePrepare), entity, "validation", "gate prepare"); err == nil {
		t.Fatal("a native completion notice for another run credited the spawned run")
	}

	// Negative control: remove both new surfaces and the run must RED. This is the
	// falsifying edit for AC-1 — drop either new credit and this case stops failing.
	noBgWait := strings.Replace(bgWaitResult, "done. Outcome: 1 complete", "running. Outcome: 0 complete", 1)
	noNotice := strings.Replace(nativeNotice, "Background task completed", "Background task failed", 1)
	if err := assertWorkerLifecycle(stream(spawn, spawnResult, bgWaitCall, noBgWait, noNotice, gatePrepare), entity, "validation", "gate prepare"); err == nil {
		t.Fatal("a lifecycle with both new completion surfaces removed passed")
	}

	// AC-2 zero-spawn control: no dispatch for this stage, so the native notice and
	// the wait cannot credit it. Credit any transcript regardless of spawn count and
	// this case stops failing.
	noSpawn := strings.Replace(spawn, "auto-continue-task-validation.md", "auto-continue-task-implementation.md", 1)
	if err := assertWorkerLifecycle(stream(noSpawn, spawnResult, bgWaitCall, bgWaitResult, nativeNotice, gatePrepare), entity, "validation", "gate prepare"); err == nil {
		t.Fatal("a Pi lifecycle with no validation-stage spawn passed")
	}

	// AC-3 ordering control: a completion that lands after the validation transition
	// must not grade GREEN. Drop the ordering check and this case stops failing.
	inverted := []string{spawn, spawnResult, gatePrepare, bgWaitCall, bgWaitResult, nativeNotice}
	if err := assertWorkerLifecycle(stream(inverted...), entity, "validation", "gate prepare"); err == nil {
		t.Fatal("a Pi completion after the validation transition passed")
	}
}
