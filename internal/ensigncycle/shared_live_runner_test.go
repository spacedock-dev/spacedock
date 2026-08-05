//go:build live

package ensigncycle

import (
	"fmt"
	"os"
	"testing"
)

//spacedock:live-suite lanes=claude-live,codex-live,pi-live
func TestLiveSharedScenarios(t *testing.T) {
	adapter, err := selectSharedLiveRuntime(os.Getenv("SPACEDOCK_LIVE_RUNTIME"))
	if err != nil {
		t.Fatal(err)
	}
	runners := sharedScenarioRunners()
	for _, scenario := range sharedRuntimeScenarios() {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			run := runners[scenario.name]
			if run == nil {
				t.Fatalf("shared scenario %q has no common runner", scenario.name)
			}
			if reason := liveDurableJourneyTODO(scenario.name); reason != "" {
				t.Skip(reason)
			}
			run(t, adapter, scenario)
		})
	}
}

type sharedLiveRuntimeAdapter interface {
	runtimeName() string
	runSharedScenario(*testing.T, sharedRuntimeScenario)
}

type sharedScenarioRunner func(*testing.T, sharedLiveRuntimeAdapter, sharedRuntimeScenario)

func runSharedScenario(t *testing.T, adapter sharedLiveRuntimeAdapter, scenario sharedRuntimeScenario) {
	t.Helper()
	adapter.runSharedScenario(t, scenario)
}

// sharedScenarioRunners is the sole executable common-journey map. Runtime
// adapters select transport; they do not own a second journey registry.
func sharedScenarioRunners() map[string]sharedScenarioRunner {
	return map[string]sharedScenarioRunner{
		"full-ensign-cycle":                  runSharedScenario,
		"gate-guardrail":                     runSharedScenario,
		"default-headless-gate-stop":         runSharedScenario,
		"withdrawn-gate-recovery":            runSharedScenario,
		"recorded-gate-lifecycle":            runSharedScenario,
		"rejection-flow":                     runSharedScenario,
		"feedback-3-cycle-escalation":        runSharedScenario,
		"merge-hook-guardrail":               runSharedScenario,
		"filing":                             runSharedScenario,
		"shallow-boot":                       runSharedScenario,
		"zero-discovery":                     runSharedScenario,
		"auto-continue-after-implementation": runSharedScenario,
		"self-evidence-merge-triage":         runSharedScenario,
		"smallest-sufficient-mechanism":      runSharedScenario,
		"keep-moving-posture":                runSharedScenario,
		"ac-value-reanchor":                  runSharedScenario,
	}
}

func selectSharedLiveRuntime(name string) (sharedLiveRuntimeAdapter, error) {
	switch name {
	case "claude":
		return claudeSharedLiveAdapter{}, nil
	case "codex":
		return codexSharedLiveAdapter{}, nil
	case "pi":
		return piSharedLiveAdapter{}, nil
	case "":
		return nil, fmt.Errorf("SPACEDOCK_LIVE_RUNTIME is required (claude, codex, or pi)")
	default:
		return nil, fmt.Errorf("unsupported SPACEDOCK_LIVE_RUNTIME %q (want claude, codex, or pi)", name)
	}
}
