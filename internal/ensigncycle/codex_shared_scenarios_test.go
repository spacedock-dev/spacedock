package ensigncycle

import (
	"reflect"
	"testing"
	"time"
)

func TestCodexSharedScenarioDefinitions(t *testing.T) {
	scenarios := codexSharedScenarios()

	var got []string
	for _, scenario := range scenarios {
		got = append(got, scenario.name)
		if scenario.oldPythonTest == "" {
			t.Fatalf("scenario %q is missing its old Python scenario source", scenario.name)
		}
		if scenario.intent == "" {
			t.Fatalf("scenario %q is missing its shared behavior intent", scenario.name)
		}
		if scenario.timeout <= 0 {
			t.Fatalf("scenario %q timeout = %s, want positive live timeout", scenario.name, scenario.timeout)
		}
	}

	want := []string{
		"gate-guardrail",
		"rejection-flow",
		"merge-hook-guardrail",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("shared Codex scenarios = %v, want %v", got, want)
	}

	timeoutByName := map[string]time.Duration{}
	for _, scenario := range scenarios {
		timeoutByName[scenario.name] = scenario.timeout
	}
	if timeoutByName["rejection-flow"] <= timeoutByName["gate-guardrail"] {
		t.Fatalf("rejection-flow timeout = %s, want more budget than gate-guardrail timeout %s",
			timeoutByName["rejection-flow"], timeoutByName["gate-guardrail"])
	}
}
