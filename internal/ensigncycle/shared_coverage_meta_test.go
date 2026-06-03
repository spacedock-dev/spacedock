//go:build live

package ensigncycle

import "testing"

// TestSharedScenarioRunnerCoverage is the AC-2/AC-3 parity guard against drift: it
// iterates the host-neutral sharedRuntimeScenarios() table and fails if EITHER the
// Codex runner map or the Claude runner map lacks a runner for any shared scenario
// ID. A scenario that exists for one host only — the exact split this task closes
// — turns this red. It is live-tagged because the runner maps reference the
// live-only runner adapter types, but it spends no model: it only inspects the
// maps' key coverage.
func TestSharedScenarioRunnerCoverage(t *testing.T) {
	codexRunners := codexScenarioRunners()
	claudeRunners := claudeScenarioRunners()

	scenarios := sharedRuntimeScenarios()
	if len(scenarios) == 0 {
		t.Fatal("sharedRuntimeScenarios() is empty; the parity guard has nothing to check")
	}

	for _, scenario := range scenarios {
		if codexRunners[scenario.name] == nil {
			t.Errorf("shared scenario %q has no Codex runner", scenario.name)
		}
		if claudeRunners[scenario.name] == nil {
			t.Errorf("shared scenario %q has no Claude runner", scenario.name)
		}
	}

	// A runner with no matching shared scenario is also drift: a host scenario the
	// shared table forgot. Guard both directions so the maps and the table stay in
	// lockstep.
	defined := map[string]bool{}
	for _, scenario := range scenarios {
		defined[scenario.name] = true
	}
	for name := range codexRunners {
		if !defined[name] {
			t.Errorf("Codex runner %q has no shared scenario definition", name)
		}
	}
	for name := range claudeRunners {
		if !defined[name] {
			t.Errorf("Claude runner %q has no shared scenario definition", name)
		}
	}
}
