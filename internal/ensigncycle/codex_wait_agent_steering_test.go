package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

type codexSteeringEvidence struct {
	SchemaVersion   int                  `json:"schema_version"`
	RunID           string               `json:"run_id"`
	Host            string               `json:"host"`
	CodexCLIVersion string               `json:"codex_cli_version"`
	Source          string               `json:"source"`
	Worker          codexSteeringWorker  `json:"worker"`
	Events          []codexSteeringEvent `json:"events"`
}

type codexSteeringWorker struct {
	TaskPath        string `json:"task_path"`
	CompletionEpoch int    `json:"completion_epoch"`
}

type codexSteeringEvent struct {
	Sequence        int     `json:"sequence"`
	AtUTC           *string `json:"at_utc"`
	Type            string  `json:"type"`
	TaskPath        string  `json:"task_path"`
	CompletionEpoch int     `json:"completion_epoch"`
	Status          string  `json:"status"`
	Message         string  `json:"message"`
	Count           int     `json:"count"`
}

func TestCodexWaitAgentSteeringEvidence(t *testing.T) {
	evidence := loadCodexSteeringEvidence(t)
	if err := assertCodexWaitAgentSteering(evidence); err != nil {
		t.Fatalf("reduced Codex steering trace must prove active-loop resumption: %v", err)
	}
}

func TestCodexWaitAgentSteeringRejectsIndependentMutants(t *testing.T) {
	base := loadCodexSteeringEvidence(t)
	target := base.Worker
	cases := []struct {
		name   string
		want   string
		mutate func(*codexSteeringEvidence)
	}{
		{
			name: "target cancellation",
			want: "lifecycle mutation",
			mutate: func(e *codexSteeringEvidence) {
				insertCodexSteeringEvent(e, 6, codexSteeringEvent{
					Type: "interrupt_agent_called", TaskPath: target.TaskPath, CompletionEpoch: target.CompletionEpoch,
				})
			},
		},
		{
			name: "target redispatch",
			want: "spawn count",
			mutate: func(e *codexSteeringEvidence) {
				insertCodexSteeringEvent(e, 6, codexSteeringEvent{
					Type: "spawn_agent_called", TaskPath: target.TaskPath, CompletionEpoch: target.CompletionEpoch,
				})
			},
		},
		{
			name: "monitoring before active scope is empty",
			want: "before active scope became empty",
			mutate: func(e *codexSteeringEvidence) {
				insertCodexSteeringEvent(e, 7, codexSteeringEvent{
					Type: "wait_agent_called", TaskPath: target.TaskPath, CompletionEpoch: target.CompletionEpoch,
				})
			},
		},
		{
			name: "wrong completion identity",
			want: "matching final status",
			mutate: func(e *codexSteeringEvidence) {
				eventOfType(t, e, "final_status").TaskPath = "/root/unrelated-worker"
			},
		},
		{
			name: "stale completion epoch",
			want: "matching final status",
			mutate: func(e *codexSteeringEvidence) {
				eventOfType(t, e, "final_status").CompletionEpoch--
			},
		},
		{
			name: "wrong durable-report identity",
			want: "durable report",
			mutate: func(e *codexSteeringEvidence) {
				eventOfType(t, e, "durable_report_read").TaskPath = "/root/unrelated-worker"
			},
		},
		{
			name: "stale durable-report epoch",
			want: "durable report",
			mutate: func(e *codexSteeringEvidence) {
				eventOfType(t, e, "durable_report_read").CompletionEpoch--
			},
		},
		{
			name: "wait return alone",
			want: "matching final status",
			mutate: func(e *codexSteeringEvidence) {
				removeCodexSteeringEvents(e, "final_status", "durable_report_read")
			},
		},
		{
			name: "repeated captain-facing interruption disclaimer",
			want: "captain-facing harness label",
			mutate: func(e *codexSteeringEvidence) {
				disclaimer := codexSteeringEvent{
					Type:    "agent_message",
					Message: "Wait interrupted by new input; an interruption returns control and the worker is not failed, closed, or redispatched.",
				}
				insertCodexSteeringEvent(e, 2, disclaimer)
				insertCodexSteeringEvent(e, 10, disclaimer)
			},
		},
		{
			name: "final status without durable report",
			want: "durable report",
			mutate: func(e *codexSteeringEvidence) {
				removeCodexSteeringEvents(e, "durable_report_read")
			},
		},
		{
			name: "durable report without final status",
			want: "matching final status",
			mutate: func(e *codexSteeringEvidence) {
				removeCodexSteeringEvents(e, "final_status")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			evidence := cloneCodexSteeringEvidence(t, base)
			tc.mutate(&evidence)
			err := assertCodexWaitAgentSteering(evidence)
			if err == nil {
				t.Fatal("mutated steering trace unexpectedly passed")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("mutation failed for the wrong invariant: got %v, want error containing %q", err, tc.want)
			}
		})
	}
}

func assertCodexWaitAgentSteering(evidence codexSteeringEvidence) error {
	if evidence.SchemaVersion != 1 || evidence.Host != "codex" ||
		evidence.RunID == "" || evidence.CodexCLIVersion == "" || evidence.Source == "" {
		return fmt.Errorf("invalid Codex steering evidence metadata")
	}
	target := evidence.Worker
	if target.TaskPath == "" || target.CompletionEpoch < 1 {
		return fmt.Errorf("invalid correlated worker identity")
	}

	spawnCount := 0
	firstWait := -1
	captainInput := -1
	runningAfterInput := -1
	usefulWork := 0
	activeScopeEmpty := -1
	secondWait := -1
	finalStatus := -1
	durableReport := -1
	var lastTimestamp time.Time

	for i, event := range evidence.Events {
		if event.Sequence != i+1 {
			return fmt.Errorf("event sequence %d has ordinal %d", i+1, event.Sequence)
		}
		if event.AtUTC != nil {
			timestamp, err := time.Parse(time.RFC3339Nano, *event.AtUTC)
			if err != nil {
				return fmt.Errorf("event %d has invalid timestamp: %w", event.Sequence, err)
			}
			if !lastTimestamp.IsZero() && timestamp.Before(lastTimestamp) {
				return fmt.Errorf("event %d timestamp precedes the prior observed timestamp", event.Sequence)
			}
			lastTimestamp = timestamp
		}
		if event.Type == "agent_message" && isOldCodexWaitDisclaimer(event.Message) {
			return fmt.Errorf("captain-facing harness label or interruption disclaimer at event %d", event.Sequence)
		}
		if isTargetLifecycleMutation(event, target) {
			return fmt.Errorf("target lifecycle mutation %q at event %d", event.Type, event.Sequence)
		}

		switch event.Type {
		case "spawn_agent_called":
			if sameCodexSteeringWorker(event, target) {
				spawnCount++
			}
		case "wait_agent_called":
			if !sameCodexSteeringWorker(event, target) {
				continue
			}
			if firstWait < 0 {
				firstWait = i
			} else if captainInput >= 0 {
				if activeScopeEmpty < 0 {
					return fmt.Errorf("monitoring resumed before active scope became empty at event %d", event.Sequence)
				}
				if secondWait < 0 {
					secondWait = i
				}
			}
		case "captain_input":
			if captainInput >= 0 {
				return fmt.Errorf("multiple captain-input events in reduced trace")
			}
			captainInput = i
		case "list_agents_result":
			if sameCodexSteeringWorker(event, target) && event.Status == "running" && captainInput >= 0 && i > captainInput {
				runningAfterInput = i
			}
		case "active_work":
			if captainInput >= 0 && i > captainInput && activeScopeEmpty < 0 && event.Count > 0 {
				usefulWork += event.Count
			}
		case "active_scope_empty":
			if captainInput >= 0 && usefulWork > 0 && activeScopeEmpty < 0 {
				activeScopeEmpty = i
			}
		case "final_status":
			if sameCodexSteeringWorker(event, target) && event.Status == "completed" {
				finalStatus = i
			}
		case "durable_report_read":
			if sameCodexSteeringWorker(event, target) && event.Status == "completed" {
				durableReport = i
			}
		case "agent_message", "harness_output", "interrupt_agent_called", "close_agent_called":
		default:
			return fmt.Errorf("event %d has unknown type %q", event.Sequence, event.Type)
		}
	}

	if spawnCount != 1 {
		return fmt.Errorf("target spawn count = %d, want 1", spawnCount)
	}
	if firstWait < 0 || captainInput <= firstWait {
		return fmt.Errorf("captain input did not follow the first monitoring call")
	}
	if runningAfterInput <= captainInput {
		return fmt.Errorf("same correlated worker was not still running after captain input")
	}
	if usefulWork < 1 || activeScopeEmpty < 0 {
		return fmt.Errorf("captain-authorized active work was not exhausted before monitoring resumed")
	}
	if secondWait <= activeScopeEmpty || secondWait <= runningAfterInput {
		return fmt.Errorf("same correlated worker was not monitored again after active work")
	}
	if finalStatus <= secondWait {
		return fmt.Errorf("matching final status was not observed after resumed monitoring")
	}
	if durableReport <= finalStatus {
		return fmt.Errorf("durable report was not read after matching final status")
	}
	return nil
}

func sameCodexSteeringWorker(event codexSteeringEvent, worker codexSteeringWorker) bool {
	return event.TaskPath == worker.TaskPath && event.CompletionEpoch == worker.CompletionEpoch
}

func isTargetLifecycleMutation(event codexSteeringEvent, worker codexSteeringWorker) bool {
	if !sameCodexSteeringWorker(event, worker) {
		return false
	}
	switch event.Type {
	case "interrupt_agent_called", "close_agent_called":
		return true
	default:
		return false
	}
}

func isOldCodexWaitDisclaimer(message string) bool {
	message = strings.ToLower(message)
	return strings.Contains(message, "wait interrupted by new input") ||
		strings.Contains(message, "interruption returns control") ||
		strings.Contains(message, "interruption only returns control")
}

func loadCodexSteeringEvidence(t *testing.T) codexSteeringEvidence {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve steering test source path")
	}
	path := filepath.Join(filepath.Dir(source), "..", "..", "docs", "dev", "_evidence",
		"codex-wait-agent-steering-semantics", "2026-07-23-dogfood.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read Codex steering evidence: %v", err)
	}
	var evidence codexSteeringEvidence
	if err := json.Unmarshal(data, &evidence); err != nil {
		t.Fatalf("parse Codex steering evidence: %v", err)
	}
	return evidence
}

func cloneCodexSteeringEvidence(t *testing.T, evidence codexSteeringEvidence) codexSteeringEvidence {
	t.Helper()
	data, err := json.Marshal(evidence)
	if err != nil {
		t.Fatal(err)
	}
	var cloned codexSteeringEvidence
	if err := json.Unmarshal(data, &cloned); err != nil {
		t.Fatal(err)
	}
	return cloned
}

func eventOfType(t *testing.T, evidence *codexSteeringEvidence, eventType string) *codexSteeringEvent {
	t.Helper()
	for i := range evidence.Events {
		if evidence.Events[i].Type == eventType {
			return &evidence.Events[i]
		}
	}
	t.Fatalf("fixture has no %s event", eventType)
	return nil
}

func insertCodexSteeringEvent(evidence *codexSteeringEvidence, at int, event codexSteeringEvent) {
	evidence.Events = append(evidence.Events, codexSteeringEvent{})
	copy(evidence.Events[at+1:], evidence.Events[at:])
	evidence.Events[at] = event
	resequenceCodexSteeringEvents(evidence)
}

func removeCodexSteeringEvents(evidence *codexSteeringEvidence, eventTypes ...string) {
	remove := map[string]bool{}
	for _, eventType := range eventTypes {
		remove[eventType] = true
	}
	kept := evidence.Events[:0]
	for _, event := range evidence.Events {
		if !remove[event.Type] {
			kept = append(kept, event)
		}
	}
	evidence.Events = kept
	resequenceCodexSteeringEvents(evidence)
}

func resequenceCodexSteeringEvents(evidence *codexSteeringEvidence) {
	for i := range evidence.Events {
		evidence.Events[i].Sequence = i + 1
	}
}
