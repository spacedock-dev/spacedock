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
		ID:     "codex-runtime-dispatch",
		Source: "live-harness",
		Host:   "codex",
		Model:  characterization.Model,
	}, characterization, BehaviorResult{Passed: true})
	if record.MetricsState != StateCharacterized {
		t.Fatalf("metrics_state = %q, want characterized", record.MetricsState)
	}
	if record.Tokens.Total != 0 {
		t.Fatalf("characterized Codex record made an unverified token claim: %+v", record.Tokens)
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
