//go:build live

package ensigncycle

import "testing"

type piSharedScenarioCoverage struct {
	mode   string
	reason string
}

func piSharedScenarioCoverageMap() map[string]piSharedScenarioCoverage {
	return map[string]piSharedScenarioCoverage{
		"gate-guardrail": {
			mode:   "gap",
			reason: "Pi currently has durable live coverage for subagent dispatch/front-door setup, but not a live-safe shared first-officer gate runner.",
		},
		"rejection-flow": {
			mode:   "gap",
			reason: "Pi currently has durable live coverage for subagent dispatch/front-door setup, but not a live-safe shared first-officer rejection-flow runner.",
		},
		"merge-hook-guardrail": {
			mode:   "gap",
			reason: "Pi currently has durable live coverage for subagent dispatch/front-door setup, but not a live-safe shared first-officer merge-hook runner.",
		},
	}
}

func TestPiSharedScenarioCoverage(t *testing.T) {
	coverage := piSharedScenarioCoverageMap()
	defined := map[string]bool{}
	for _, scenario := range sharedRuntimeScenarios() {
		defined[scenario.name] = true
		entry, ok := coverage[scenario.name]
		if !ok {
			t.Errorf("shared scenario %q has no Pi coverage entry", scenario.name)
			continue
		}
		switch entry.mode {
		case "live", "codified", "gap":
		default:
			t.Errorf("shared scenario %q has invalid Pi coverage mode %q", scenario.name, entry.mode)
		}
		if entry.reason == "" {
			t.Errorf("shared scenario %q Pi coverage entry needs an honest reason", scenario.name)
		}
	}
	for name := range coverage {
		if !defined[name] {
			t.Errorf("Pi coverage entry %q has no shared scenario definition", name)
		}
	}
}
