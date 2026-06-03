package ensigncycle

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

// TestSharedRuntimeScenarioDefinitions is the AC-1 guard: the shared runtime
// scenarios are defined once, in host-neutral code. It pins the scenario ID set,
// requires every scenario to carry its runtime-neutral facts (provenance, intent,
// positive live timeout), and reflects over the table type to prove it encodes NO
// Claude-only or Codex-only field — the structural guard against a runner concern
// (auth, plugin, launch) leaking back into the shared definition.
func TestSharedRuntimeScenarioDefinitions(t *testing.T) {
	scenarios := sharedRuntimeScenarios()

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
		t.Fatalf("shared runtime scenarios = %v, want %v", got, want)
	}

	timeoutByName := map[string]time.Duration{}
	for _, scenario := range scenarios {
		timeoutByName[scenario.name] = scenario.timeout
	}
	if timeoutByName["rejection-flow"] <= timeoutByName["gate-guardrail"] {
		t.Fatalf("rejection-flow timeout = %s, want more budget than gate-guardrail timeout %s",
			timeoutByName["rejection-flow"], timeoutByName["gate-guardrail"])
	}

	// AC-1: the host-neutral table type encodes ONLY runtime-neutral facts. Any
	// field naming a single host (codex/claude) would mean a runner concern leaked
	// into the shared definition, the exact parity drift this table exists to
	// prevent. Pin the exact field set and reject any host-named field.
	typ := reflect.TypeOf(sharedRuntimeScenario{})
	wantFields := map[string]bool{
		"name":          true,
		"oldPythonTest": true,
		"intent":        true,
		"timeout":       true,
	}
	if typ.NumField() != len(wantFields) {
		t.Fatalf("sharedRuntimeScenario has %d fields, want %d host-neutral fields", typ.NumField(), len(wantFields))
	}
	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Name
		lower := strings.ToLower(name)
		if strings.Contains(lower, "codex") || strings.Contains(lower, "claude") {
			t.Fatalf("sharedRuntimeScenario field %q names a single host; runner concerns must live in the host adapters, not the shared table", name)
		}
		if !wantFields[name] {
			t.Fatalf("sharedRuntimeScenario has unexpected field %q; the host-neutral table carries only %v", name, keysOf(wantFields))
		}
	}
}

func keysOf(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
