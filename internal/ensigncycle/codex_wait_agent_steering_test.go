package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	EvidenceDir     string               `json:"-"`
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
	ArtifactPath    string  `json:"artifact_path"`
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
		mutate func(*testing.T, *codexSteeringEvidence)
	}{
		{
			name: "target cancellation",
			want: "lifecycle mutation",
			mutate: func(_ *testing.T, e *codexSteeringEvidence) {
				insertCodexSteeringEvent(e, 6, codexSteeringEvent{
					Type: "interrupt_agent_called", TaskPath: target.TaskPath, CompletionEpoch: target.CompletionEpoch,
				})
			},
		},
		{
			name: "target cancellation with omitted epoch",
			want: "lifecycle mutation",
			mutate: func(_ *testing.T, e *codexSteeringEvidence) {
				insertCodexSteeringEvent(e, 6, codexSteeringEvent{
					Type: "interrupt_agent_called", TaskPath: target.TaskPath,
				})
			},
		},
		{
			name: "target cancellation with invalid epoch",
			want: "lifecycle mutation",
			mutate: func(_ *testing.T, e *codexSteeringEvidence) {
				insertCodexSteeringEvent(e, 6, codexSteeringEvent{
					Type: "interrupt_agent_called", TaskPath: target.TaskPath, CompletionEpoch: -1,
				})
			},
		},
		{
			name: "target cancellation with stale epoch",
			want: "lifecycle mutation",
			mutate: func(_ *testing.T, e *codexSteeringEvidence) {
				insertCodexSteeringEvent(e, 6, codexSteeringEvent{
					Type: "interrupt_agent_called", TaskPath: target.TaskPath, CompletionEpoch: target.CompletionEpoch - 1,
				})
			},
		},
		{
			name: "target redispatch",
			want: "spawn count",
			mutate: func(_ *testing.T, e *codexSteeringEvidence) {
				insertCodexSteeringEvent(e, 6, codexSteeringEvent{
					Type: "spawn_agent_called", TaskPath: target.TaskPath, CompletionEpoch: target.CompletionEpoch,
				})
			},
		},
		{
			name: "monitoring before active scope is empty",
			want: "before active scope became empty",
			mutate: func(_ *testing.T, e *codexSteeringEvidence) {
				insertCodexSteeringEvent(e, 7, codexSteeringEvent{
					Type: "wait_agent_called", TaskPath: target.TaskPath, CompletionEpoch: target.CompletionEpoch,
				})
			},
		},
		{
			name: "missing harness wait return",
			want: "exactly one harness wait-return",
			mutate: func(_ *testing.T, e *codexSteeringEvidence) {
				removeCodexSteeringEvents(e, "harness_output")
			},
		},
		{
			name: "harness wait return after captain input",
			want: "did not occur between",
			mutate: func(t *testing.T, e *codexSteeringEvidence) {
				swapCodexSteeringEvents(t, e, "harness_output", "captain_input")
			},
		},
		{
			name: "duplicate harness wait return",
			want: "exactly one harness wait-return",
			mutate: func(t *testing.T, e *codexSteeringEvidence) {
				event := *eventOfType(t, e, "harness_output")
				insertCodexSteeringEvent(e, 4, event)
			},
		},
		{
			name: "wrong completion identity",
			want: "matching final status",
			mutate: func(t *testing.T, e *codexSteeringEvidence) {
				eventOfType(t, e, "final_status").TaskPath = "/root/unrelated-worker"
			},
		},
		{
			name: "stale completion epoch",
			want: "matching final status",
			mutate: func(t *testing.T, e *codexSteeringEvidence) {
				eventOfType(t, e, "final_status").CompletionEpoch--
			},
		},
		{
			name: "wrong durable-report identity",
			want: "durable report",
			mutate: func(t *testing.T, e *codexSteeringEvidence) {
				eventOfType(t, e, "durable_report_read").TaskPath = "/root/unrelated-worker"
			},
		},
		{
			name: "stale durable-report epoch",
			want: "durable report",
			mutate: func(t *testing.T, e *codexSteeringEvidence) {
				eventOfType(t, e, "durable_report_read").CompletionEpoch--
			},
		},
		{
			name: "durable artifact has wrong identity",
			want: "worker identity",
			mutate: func(t *testing.T, e *codexSteeringEvidence) {
				mutateCodexSteeringReport(t, e, target.TaskPath, "/root/unrelated-worker")
			},
		},
		{
			name: "durable artifact has stale epoch",
			want: "completion epoch",
			mutate: func(t *testing.T, e *codexSteeringEvidence) {
				mutateCodexSteeringReport(t, e, "Completion epoch: `2`", "Completion epoch: `1`")
			},
		},
		{
			name: "durable artifact lacks stage report",
			want: "implementation stage report",
			mutate: func(t *testing.T, e *codexSteeringEvidence) {
				mutateCodexSteeringReport(t, e, "## Stage Report: implementation", "## Notes")
			},
		},
		{
			name: "durable artifact lacks completed worker status",
			want: "completed worker status",
			mutate: func(t *testing.T, e *codexSteeringEvidence) {
				mutateCodexSteeringReport(t, e, "Worker status: `completed`", "Worker status: `running`")
			},
		},
		{
			name: "wait return alone",
			want: "matching final status",
			mutate: func(_ *testing.T, e *codexSteeringEvidence) {
				removeCodexSteeringEvents(e, "final_status", "durable_report_read")
			},
		},
		{
			name: "repeated captain-facing interruption disclaimer",
			want: "captain-facing harness label",
			mutate: func(_ *testing.T, e *codexSteeringEvidence) {
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
			mutate: func(_ *testing.T, e *codexSteeringEvidence) {
				removeCodexSteeringEvents(e, "durable_report_read")
			},
		},
		{
			name: "durable report without final status",
			want: "matching final status",
			mutate: func(_ *testing.T, e *codexSteeringEvidence) {
				removeCodexSteeringEvents(e, "final_status")
			},
		},
		{
			name: "invalid event sequence",
			want: "declares sequence",
			mutate: func(_ *testing.T, e *codexSteeringEvidence) {
				e.Events[3].Sequence = 99
			},
		},
		{
			name: "timestamp moves backward",
			want: "precedes the prior",
			mutate: func(t *testing.T, e *codexSteeringEvidence) {
				timestamp := "2026-07-23T14:30:00Z"
				eventOfType(t, e, "list_agents_result").AtUTC = &timestamp
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			evidence := cloneCodexSteeringEvidence(t, base)
			tc.mutate(t, &evidence)
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
	harnessReturn := -1
	harnessReturnCount := 0
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
			return fmt.Errorf("event at position %d declares sequence %d", i+1, event.Sequence)
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
		case "harness_output":
			if strings.TrimSpace(event.Message) != "Wait interrupted by new input" {
				return fmt.Errorf("event %d has unrecognized harness wait-return label", event.Sequence)
			}
			harnessReturn = i
			harnessReturnCount++
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
				if err := validateCodexSteeringReport(evidence.EvidenceDir, event.ArtifactPath, target); err != nil {
					return fmt.Errorf("durable report is invalid: %w", err)
				}
				durableReport = i
			}
		case "agent_message", "interrupt_agent_called", "close_agent_called":
		default:
			return fmt.Errorf("event %d has unknown type %q", event.Sequence, event.Type)
		}
	}

	if spawnCount != 1 {
		return fmt.Errorf("target spawn count = %d, want 1", spawnCount)
	}
	if harnessReturnCount != 1 {
		return fmt.Errorf("expected exactly one harness wait-return event, got %d", harnessReturnCount)
	}
	if firstWait < 0 || harnessReturn <= firstWait || captainInput <= harnessReturn {
		return fmt.Errorf("harness wait return did not occur between the first monitoring call and captain input")
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
	switch event.Type {
	case "interrupt_agent_called", "close_agent_called":
		return event.TaskPath == "" || event.TaskPath == worker.TaskPath
	default:
		return false
	}
}

var codexImplementationStageReport = regexp.MustCompile(`(?m)^## Stage Report: implementation[ \t]*$`)

func validateCodexSteeringReport(evidenceDir, artifactPath string, worker codexSteeringWorker) error {
	if evidenceDir == "" || artifactPath == "" || filepath.Base(artifactPath) != artifactPath {
		return fmt.Errorf("artifact path must name one file in the evidence directory")
	}
	data, err := os.ReadFile(filepath.Join(evidenceDir, artifactPath))
	if err != nil {
		return fmt.Errorf("read artifact: %w", err)
	}
	report := string(data)
	if !codexDurableStatusForEntity(report, "implementation") {
		return fmt.Errorf("artifact does not retain implementation entity status")
	}
	if !strings.Contains(report, "Worker status: `completed`") {
		return fmt.Errorf("artifact does not retain completed worker status")
	}
	if !strings.Contains(report, "Worker task path: `"+worker.TaskPath+"`") {
		return fmt.Errorf("artifact does not retain the correlated worker identity")
	}
	if !strings.Contains(report, fmt.Sprintf("Completion epoch: `%d`", worker.CompletionEpoch)) {
		return fmt.Errorf("artifact does not retain the correlated completion epoch")
	}
	if !codexImplementationStageReport.MatchString(report) {
		return fmt.Errorf("artifact does not contain an implementation stage report")
	}
	return nil
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
	evidence.EvidenceDir = filepath.Dir(path)
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
	cloned.EvidenceDir = evidence.EvidenceDir
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

func swapCodexSteeringEvents(t *testing.T, evidence *codexSteeringEvidence, first, second string) {
	t.Helper()
	firstEvent := eventOfType(t, evidence, first)
	secondEvent := eventOfType(t, evidence, second)
	*firstEvent, *secondEvent = *secondEvent, *firstEvent
	resequenceCodexSteeringEvents(evidence)
}

func mutateCodexSteeringReport(t *testing.T, evidence *codexSteeringEvidence, old, replacement string) {
	t.Helper()
	event := eventOfType(t, evidence, "durable_report_read")
	data, err := os.ReadFile(filepath.Join(evidence.EvidenceDir, event.ArtifactPath))
	if err != nil {
		t.Fatal(err)
	}
	report := string(data)
	if !strings.Contains(report, old) {
		t.Fatalf("durable report lacks mutation target %q", old)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, event.ArtifactPath), []byte(strings.Replace(report, old, replacement, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	evidence.EvidenceDir = dir
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
