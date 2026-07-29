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
		"recorded-gate-lifecycle": {
			mode:   "gap",
			reason: "TODO(9w59t6m1qc46hccd54p04z2j): the runner remains registered and selected but is quarantined because delegated gate presentation-to-application/dispatch is unreliable; it provides no live capability evidence.",
		},
		"rejection-flow": {
			mode:   "gap",
			reason: "Pi currently has durable live coverage for subagent dispatch/front-door setup, but not a live-safe shared first-officer rejection-flow runner.",
		},
		"feedback-3-cycle-escalation": {
			mode:   "gap",
			reason: "Pi currently has durable live coverage for subagent dispatch/front-door setup, but not a live-safe shared first-officer 3-cycle-escalation runner.",
		},
		"merge-hook-guardrail": {
			mode:   "gap",
			reason: "Pi currently has durable live coverage for subagent dispatch/front-door setup, but not a live-safe shared first-officer merge-hook runner.",
		},
		"filing": {
			mode:   "gap",
			reason: "Pi currently has durable live coverage for subagent dispatch/front-door setup, but not a live-safe shared first-officer filing runner.",
		},
		"shallow-boot": {
			mode:   "gap",
			reason: "Pi currently has durable live coverage for subagent dispatch/front-door setup, but not a live-safe shared first-officer shallow-boot runner.",
		},
		"self-evidence-merge-triage": {
			mode:   "gap",
			reason: "Pi currently has durable live coverage for subagent dispatch/front-door setup, but not a live-safe shared first-officer merge/triage self-evidence runner.",
		},
		"smallest-sufficient-mechanism": {
			mode:   "gap",
			reason: "Pi currently has durable live coverage for subagent dispatch/front-door setup, but not a live-safe shared first-officer smallest-sufficient-mechanism runner.",
		},
		"keep-moving-posture": {
			mode:   "gap",
			reason: "Pi currently has durable live coverage for subagent dispatch/front-door setup, but not a live-safe shared first-officer keep-moving-posture runner.",
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
