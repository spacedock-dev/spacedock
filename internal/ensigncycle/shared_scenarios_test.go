package ensigncycle

import "time"

// sharedRuntimeScenario is the host-neutral definition of one runtime regression
// scenario. Every field is a runtime-neutral fact about the user journey — its
// ID, its old Python provenance, the behavior intent it guards, and its live
// timeout/cost class. It carries NO launch, auth, plugin, artifact, or transcript
// field: those host-specific concerns live behind the per-host runner adapters
// (the Codex runner in codex_live_runner_test.go, the Claude runner in claude_live_runner_test.go),
// each of which implements the same scenario IDs. The shared coverage meta-tests
// fail if either host lacks a runner for a shared scenario, so a scenario cannot
// drift to one host only.
type sharedRuntimeScenario struct {
	name          string
	oldPythonTest string
	intent        string
	timeout       time.Duration
}

func sharedRuntimeScenarios() []sharedRuntimeScenario {
	return []sharedRuntimeScenario{
		{
			name:          "gate-guardrail",
			oldPythonTest: "tests/test_gate_guardrail.py",
			intent:        "FO halts at a human gate and presents the review without self-approval, mutation, or archival.",
			timeout:       2 * time.Minute,
		},
		{
			name:          "rejection-flow",
			oldPythonTest: "tests/test_rejection_flow.py",
			intent:        "FO observes a rejected validation report and routes the concrete finding back through implementation.",
			timeout:       4 * time.Minute,
		},
		{
			name:          "merge-hook-guardrail",
			oldPythonTest: "tests/test_merge_hook_guardrail.py",
			intent:        "FO cannot bypass a registered merge hook by terminalizing without pr, mod-block, or force.",
			timeout:       90 * time.Second,
		},
	}
}
