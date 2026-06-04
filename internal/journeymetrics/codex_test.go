package journeymetrics

import (
	"strings"
	"testing"
)

func TestCodexFixtureIsCharacterizedNotMeasured(t *testing.T) {
	data := readTestdata(t, "codex-exec.jsonl")

	characterization, err := CharacterizeCodexExecJSONL(data)
	if err != nil {
		t.Fatalf("CharacterizeCodexExecJSONL: %v", err)
	}

	wantKinds := []string{
		"message",
		"session.completed",
		"session.started",
		"tool_call.completed",
		"tool_call.started",
		"turn.completed",
		"turn.started",
	}
	if strings.Join(characterization.EventKinds, "\n") != strings.Join(wantKinds, "\n") {
		t.Fatalf("event kinds = %#v, want %#v", characterization.EventKinds, wantKinds)
	}
	if !hasFields(characterization.FieldsByEvent["tool_call.completed"], "call_id", "exit_code", "name", "type") {
		t.Fatalf("tool_call.completed fields = %#v", characterization.FieldsByEvent["tool_call.completed"])
	}

	record := CodexCharacterizedRecord(JourneySpec{
		ScenarioID: "runtime-dispatch",
		Source:     "live-harness",
		Mode:       ModeLLMLive,
		Runtime:    "codex",
		Executor:   "llm",
		Host:       "codex",
		Model:      characterization.Model,
	}, characterization, BehaviorResult{Passed: true})
	if record.MetricsState != StateCharacterized {
		t.Fatalf("metrics_state = %q, want characterized", record.MetricsState)
	}
	if record.Tokens.Total != 0 {
		t.Fatalf("characterized Codex record made an unverified token claim: %+v", record.Tokens)
	}
	if record.Mode != ModeLLMLive {
		t.Fatalf("mode = %q, want llm-live", record.Mode)
	}
	if record.Model != characterization.Model {
		t.Fatalf("model = %q, want characterization model %q", record.Model, characterization.Model)
	}

	recordWithoutMode := CodexCharacterizedRecord(JourneySpec{
		ScenarioID: "runtime-dispatch",
		Source:     "live-harness",
		Runtime:    "codex",
		Executor:   "llm",
		Host:       "codex",
	}, characterization, BehaviorResult{Passed: true})
	if recordWithoutMode.Mode != "" {
		t.Fatalf("mode without explicit evidence mode = %q, want empty rather than model fallback", recordWithoutMode.Mode)
	}
	if recordWithoutMode.Model != characterization.Model {
		t.Fatalf("model without explicit model = %q, want characterization model %q", recordWithoutMode.Model, characterization.Model)
	}

	maxTokens := 1
	budget := EvaluateBudget(record, Budget{MaxTotalTokens: &maxTokens})
	if !budget.Blocking || !strings.Contains(strings.Join(budget.Violations, "\n"), "measured") {
		t.Fatalf("Codex characterized token ceiling was not rejected: %+v", budget)
	}
}

func hasFields(got []string, want ...string) bool {
	seen := map[string]bool{}
	for _, field := range got {
		seen[field] = true
	}
	for _, field := range want {
		if !seen[field] {
			return false
		}
	}
	return true
}
