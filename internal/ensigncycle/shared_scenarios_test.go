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
			intent:        "FO drives a two-cycle rejection trajectory — route back, re-implement, re-validate via reviewer reuse — restoring the dropped second cycle.",
			// Sized against a MEASURED local sonnet run of the now-heavier
			// before-cycle-1 fixture: 4 live ensign phases (impl-c1 → val-c1 REJECT →
			// impl-rework → val-c2 reuse) totalled 13.65 min end-to-end (~3.1 min of
			// ensign work + ~10.5 min headless teams-mode inbox-polling/orchestration
			// overhead, which the fixture cannot trim). 22m leaves margin for the
			// slower opus lane; the old 8m predates the live cycle-1 validation phase.
			timeout:       22 * time.Minute,
		},
		{
			name:          "feedback-3-cycle-escalation",
			oldPythonTest: "tests/test_rejection_flow.py",
			intent:        "On the 3rd consecutive REJECTED validation the FO escalates to the human instead of auto-bouncing a 4th time.",
			timeout:       8 * time.Minute,
		},
		{
			name:          "merge-hook-guardrail",
			oldPythonTest: "tests/test_merge_hook_guardrail.py",
			intent:        "FO cannot bypass a registered merge hook by terminalizing without pr, mod-block, or force.",
			timeout:       90 * time.Second,
		},
	}
}
