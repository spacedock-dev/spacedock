package journeymetrics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestTrackingOptInDoesNotOwnBehaviorOutcome(t *testing.T) {
	dir := t.TempDir()

	unitResult, err := Track(dir, JourneySpec{
		ID:     "unit-fake",
		Source: "unit-test",
		Host:   "go-test",
		Model:  "fake",
	}, func() BehaviorResult {
		return BehaviorResult{Passed: true}
	}, func() Observation {
		return Observation{
			MetricsState:    StateMeasured,
			Duration:        25 * time.Millisecond,
			Turns:           1,
			ToolCalls:       1,
			ToolCallsByName: map[string]int{"Bash": 1},
			Tokens:          TokenTotals{Input: 10, Output: 2, Total: 12},
		}
	})
	if err != nil {
		t.Fatalf("track unit journey: %v", err)
	}
	if !unitResult.Passed {
		t.Fatalf("unit behavior was changed by metrics tracking: %+v", unitResult)
	}

	liveResult, err := Track(dir, JourneySpec{
		ID:     "live-fake",
		Source: "live-harness",
		Host:   "claude",
		Model:  "sonnet",
	}, func() BehaviorResult {
		return BehaviorResult{Passed: false, Failure: "behavior assertion failed"}
	}, func() Observation {
		return Observation{
			MetricsState:    StateMeasured,
			Duration:        150 * time.Millisecond,
			Turns:           2,
			ToolCalls:       2,
			ToolCallsByName: map[string]int{"Read": 1, "Bash": 1},
			Tokens:          TokenTotals{Input: 20, Output: 5, Total: 25},
		}
	})
	if err != nil {
		t.Fatalf("track live journey: %v", err)
	}
	if liveResult.Passed || liveResult.Failure == "" {
		t.Fatalf("live behavior outcome was hidden by metrics tracking: %+v", liveResult)
	}

	records := readRecordFiles(t, dir)
	if len(records) != 2 {
		t.Fatalf("emitted records = %d, want 2", len(records))
	}
	byID := map[string]Record{}
	for _, record := range records {
		byID[record.JourneyID] = record
	}
	if byID["unit-fake"].Outcome.Status != "passed" {
		t.Errorf("unit outcome = %q, want passed", byID["unit-fake"].Outcome.Status)
	}
	if byID["live-fake"].Outcome.Status != "failed" {
		t.Errorf("live outcome = %q, want failed", byID["live-fake"].Outcome.Status)
	}
	if byID["unit-fake"].ToolCalls != 1 || byID["live-fake"].ToolCalls != 2 {
		t.Errorf("metrics not emitted independently: %+v", byID)
	}
}

func TestBudgetPolicyIsExplicitAndCoarse(t *testing.T) {
	record := Record{
		JourneyID:       "measured",
		MetricsState:    StateMeasured,
		Tokens:          TokenTotals{Total: 101},
		ToolCalls:       3,
		ToolCallsByName: map[string]int{"Bash": 3},
	}

	if got := EvaluateBudget(record, Budget{}); got.Blocking {
		t.Fatalf("no-budget journey blocked from cost drift: %+v", got)
	}

	maxTokens := 100
	got := EvaluateBudget(record, Budget{MaxTotalTokens: &maxTokens})
	if !got.Blocking || !strings.Contains(strings.Join(got.Violations, "\n"), "max_total_tokens") {
		t.Fatalf("token ceiling did not report a configured violation: %+v", got)
	}

	record.Tokens.Total = 100
	got = EvaluateBudget(record, Budget{MaxTotalTokens: &maxTokens})
	if got.Blocking {
		t.Fatalf("exact token equality should not be required; equal-to ceiling must pass: %+v", got)
	}

	maxToolCalls := 2
	record.ToolCalls = 3
	got = EvaluateBudget(record, Budget{MaxToolCalls: &maxToolCalls})
	if !got.Blocking || !strings.Contains(strings.Join(got.Violations, "\n"), "max_tool_calls") {
		t.Fatalf("tool-call ceiling did not report a configured violation: %+v", got)
	}
}

func readRecordFiles(t *testing.T, dir string) []Record {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var records []Record
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		var record Record
		if err := json.Unmarshal(data, &record); err != nil {
			t.Fatalf("parse %s: %v\n%s", entry.Name(), err, data)
		}
		records = append(records, record)
	}
	return records
}
