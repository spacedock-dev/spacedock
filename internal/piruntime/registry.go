// ABOUTME: Minimal Pi worker registry for epoch-scoped completion evidence.
// ABOUTME: Prevents stale completion text from satisfying a later routed turn.
package piruntime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// WorkerRecord is the minimal workflow metadata needed to reason about a Pi
// worker handle without building a second session system on top of Pi.
type WorkerRecord struct {
	WorkerLabel     string `json:"worker_label"`
	Substrate       string `json:"substrate"`
	RunID           string `json:"run_id,omitempty"`
	SessionFile     string `json:"session_file,omitempty"`
	EntitySlug      string `json:"entity_slug,omitempty"`
	Stage           string `json:"stage,omitempty"`
	State           string `json:"state"`
	CompletionEpoch int    `json:"completion_epoch"`
}

// CompletionEvidence is the completion tuple observed from a Pi substrate.
type CompletionEvidence struct {
	WorkerLabel     string
	RunID           string
	CompletionEpoch int
}

type Registry struct {
	path    string
	records map[string]WorkerRecord
}

func OpenRegistry(path string) (*Registry, error) {
	r := &Registry{path: path, records: map[string]WorkerRecord{}}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return r, nil
		}
		return nil, err
	}
	if len(b) == 0 {
		return r, nil
	}
	if err := json.Unmarshal(b, &r.records); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Registry) Get(workerLabel string) (WorkerRecord, bool) {
	rec, ok := r.records[workerLabel]
	return rec, ok
}

func (r *Registry) Upsert(rec WorkerRecord) error {
	if rec.WorkerLabel == "" {
		return fmt.Errorf("worker label is required")
	}
	r.records[rec.WorkerLabel] = rec
	return r.save()
}

func (r *Registry) MarkActiveAgain(workerLabel, stage string) (WorkerRecord, error) {
	rec, ok := r.records[workerLabel]
	if !ok {
		return WorkerRecord{}, fmt.Errorf("worker %q not found", workerLabel)
	}
	if rec.State == "active" {
		return WorkerRecord{}, fmt.Errorf("worker %q is already active", workerLabel)
	}
	rec.State = "active"
	rec.Stage = stage
	rec.CompletionEpoch++
	r.records[workerLabel] = rec
	return rec, r.save()
}

func (r *Registry) MarkCompleted(workerLabel, runID string) error {
	rec, ok := r.records[workerLabel]
	if !ok {
		return fmt.Errorf("worker %q not found", workerLabel)
	}
	rec.State = "completed"
	rec.RunID = runID
	r.records[workerLabel] = rec
	return r.save()
}

func (r *Registry) CompletionIsCurrent(ev CompletionEvidence) bool {
	rec, ok := r.records[ev.WorkerLabel]
	if !ok {
		return false
	}
	return rec.State == "completed" && rec.RunID == ev.RunID && rec.CompletionEpoch == ev.CompletionEpoch
}

func (r *Registry) save() error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(r.records, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(r.path, b, 0o644)
}
