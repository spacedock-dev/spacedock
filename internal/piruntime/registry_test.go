// ABOUTME: Pi worker registry contract tests — follow-up completion evidence is
// ABOUTME: epoch-scoped so stale results cannot satisfy a reused worker turn.
package piruntime

import (
	"path/filepath"
	"testing"
)

func TestRegistryRejectsStaleCompletionAfterFollowup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "workers.json")
	reg, err := OpenRegistry(path)
	if err != nil {
		t.Fatal(err)
	}

	initial := WorkerRecord{
		WorkerLabel:     "spacedock-ensign-pi-runtime-support-implementation",
		Substrate:       "subagents",
		RunID:           "run-1",
		EntitySlug:      "pi-runtime-support",
		Stage:           "implementation",
		State:           "completed",
		CompletionEpoch: 0,
	}
	if err := reg.Upsert(initial); err != nil {
		t.Fatal(err)
	}
	stale := CompletionEvidence{WorkerLabel: initial.WorkerLabel, RunID: "run-1", CompletionEpoch: 0}
	if !reg.CompletionIsCurrent(stale) {
		t.Fatal("initial completion should be current before follow-up")
	}

	followup, err := reg.MarkActiveAgain(initial.WorkerLabel, "validation")
	if err != nil {
		t.Fatal(err)
	}
	if followup.CompletionEpoch != 1 || followup.State != "active" || followup.Stage != "validation" {
		t.Fatalf("bad follow-up record: %#v", followup)
	}
	if reg.CompletionIsCurrent(stale) {
		t.Fatal("stale completion from epoch 0 must not satisfy epoch 1 follow-up")
	}

	current := CompletionEvidence{WorkerLabel: initial.WorkerLabel, RunID: "run-2", CompletionEpoch: 1}
	if err := reg.MarkCompleted(initial.WorkerLabel, current.RunID); err != nil {
		t.Fatal(err)
	}
	if !reg.CompletionIsCurrent(current) {
		t.Fatal("epoch 1 completion should be current after MarkCompleted")
	}
}

func TestRegistryPersistsRecords(t *testing.T) {
	path := filepath.Join(t.TempDir(), "workers.json")
	reg, err := OpenRegistry(path)
	if err != nil {
		t.Fatal(err)
	}
	want := WorkerRecord{WorkerLabel: "worker-a", Substrate: "subagents", RunID: "run-1", State: "completed"}
	if err := reg.Upsert(want); err != nil {
		t.Fatal(err)
	}

	reopened, err := OpenRegistry(path)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := reopened.Get("worker-a")
	if !ok {
		t.Fatal("persisted worker missing after reopen")
	}
	if got.RunID != want.RunID || got.Substrate != want.Substrate || got.State != want.State {
		t.Fatalf("persisted record mismatch: got %#v want %#v", got, want)
	}
}
