//go:build live

package ensigncycle

import (
	"reflect"
	"strings"
	"testing"
)

func TestSharedScenarioRunnerCoverageFinal(t *testing.T) {
	want := []string{
		"full-ensign-cycle",
		"gate-guardrail",
		"default-headless-gate-stop",
		"withdrawn-gate-recovery",
		"recorded-gate-lifecycle",
		"rejection-flow",
		"feedback-3-cycle-escalation",
		"merge-hook-guardrail",
		"filing",
		"shallow-boot",
		"zero-discovery",
		"auto-continue-after-implementation",
		"self-evidence-merge-triage",
		"smallest-sufficient-mechanism",
		"keep-moving-posture",
		"ac-value-reanchor",
	}

	var scenarios []string
	for _, scenario := range sharedRuntimeScenarios() {
		scenarios = append(scenarios, scenario.name)
	}
	if !reflect.DeepEqual(scenarios, want) {
		t.Fatalf("shared scenarios = %v, want %v", scenarios, want)
	}

	runners := sharedScenarioRunners()
	if len(runners) != len(want) {
		t.Fatalf("shared runner count = %d, want %d", len(runners), len(want))
	}
	for _, name := range want {
		if runners[name] == nil {
			t.Errorf("shared journey %q has no runner", name)
		}
	}
	for name := range runners {
		if !containsString(want, name) {
			t.Errorf("runner %q has no registry journey", name)
		}
	}
}

func TestSharedLiveRuntimeSelection(t *testing.T) {
	for _, runtime := range []string{"claude", "codex", "pi"} {
		t.Run(runtime, func(t *testing.T) {
			adapter, err := selectSharedLiveRuntime(runtime)
			if err != nil {
				t.Fatal(err)
			}
			if got := adapter.runtimeName(); got != runtime {
				t.Fatalf("adapter runtime = %q, want %q", got, runtime)
			}
		})
	}
	if _, err := selectSharedLiveRuntime(""); err == nil {
		t.Fatal("empty SPACEDOCK_LIVE_RUNTIME accepted")
	}
	if _, err := selectSharedLiveRuntime("other"); err == nil {
		t.Fatal("unknown SPACEDOCK_LIVE_RUNTIME accepted")
	}
}

func TestPromotedCommonJourneyEntrypoints(t *testing.T) {
	for _, name := range []string{
		"full-ensign-cycle",
		"default-headless-gate-stop",
		"withdrawn-gate-recovery",
		"zero-discovery",
		"auto-continue-after-implementation",
		"ac-value-reanchor",
	} {
		if sharedScenarioRunners()[name] == nil {
			t.Errorf("promoted journey %q is not selected by the common runner", name)
		}
	}
}

func TestSharedScenarioSequenceStopsAfterFirstFailure(t *testing.T) {
	scenarios := []sharedRuntimeScenario{{name: "first"}, {name: "second"}}
	var ran []string
	runSharedScenarioSequence(scenarios, func(scenario sharedRuntimeScenario) bool {
		ran = append(ran, scenario.name)
		return false
	})
	if want := []string{"first"}; !reflect.DeepEqual(ran, want) {
		t.Fatalf("ran scenarios = %v, want %v", ran, want)
	}
}

func TestSharedLiveTODOEvidenceSet(t *testing.T) {
	want := map[string]string{
		"auto-continue-after-implementation": liveDurableJourneyDefectID,
		"smallest-sufficient-mechanism":      liveDurableJourneyDefectID,
		"keep-moving-posture":                liveDurableJourneyDefectID,
	}
	for _, scenario := range sharedRuntimeScenarios() {
		reason := liveDurableJourneyTODO(scenario.name)
		owner, missing := want[scenario.name]
		if missing {
			if reason == "" || !strings.HasPrefix(reason, "TODO("+owner+"):") {
				t.Errorf("TODO journey %q reason = %q, want exact owner TODO(%s)", scenario.name, reason, owner)
			}
			delete(want, scenario.name)
			continue
		}
		if reason != "" {
			t.Errorf("implemented journey %q unexpectedly has missing-evidence TODO %q", scenario.name, reason)
		}
	}
	if len(want) != 0 {
		t.Fatalf("TODO journeys are hidden from the shared scenario table: %v", want)
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
