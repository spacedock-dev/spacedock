package ensigncycle

import "time"

type codexSharedScenario struct {
	name          string
	oldPythonTest string
	intent        string
	timeout       time.Duration
}

func codexSharedScenarios() []codexSharedScenario {
	return []codexSharedScenario{
		{
			name:          "gate-guardrail",
			oldPythonTest: "tests/test_gate_guardrail.py --runtime codex",
			intent:        "FO halts at a human gate and presents the review without self-approval, mutation, or archival.",
			timeout:       60 * time.Second,
		},
		{
			name:          "rejection-flow",
			oldPythonTest: "tests/test_rejection_flow.py --runtime codex",
			intent:        "FO observes a rejected validation report and routes the concrete finding back through implementation.",
			timeout:       4 * time.Minute,
		},
		{
			name:          "merge-hook-guardrail",
			oldPythonTest: "tests/test_merge_hook_guardrail.py --runtime codex",
			intent:        "FO cannot bypass a registered merge hook by terminalizing without pr, mod-block, or force.",
			timeout:       90 * time.Second,
		},
	}
}
