package ensigncycle

// sharedRuntimeScenario is the host-neutral definition of one runtime regression
// scenario. Every field is a runtime-neutral fact about the user journey — its
// ID, its old Python provenance, and the behavior intent it guards. It carries NO
// launch, auth, plugin, artifact, transcript, OR timeout field: those host-specific
// concerns live behind the per-host runner adapters (the Codex runner in
// codex_live_runner_test.go, the Claude runner in claude_live_runner_test.go), each
// of which implements the same scenario IDs. Liveness is guarded uniformly by the
// per-stage no-progress quiet budget in those runners (the shared streamWatcher's
// quietBudgetDefault, 60s) — a per-scenario basket timeout is banned, so the shared
// table carries no timeout. The shared coverage meta-tests fail if either host
// lacks a runner for a shared scenario, so a scenario cannot drift to one host only.
type sharedRuntimeScenario struct {
	name          string
	oldPythonTest string
	intent        string
}

func sharedRuntimeScenarios() []sharedRuntimeScenario {
	return []sharedRuntimeScenario{
		{
			name:          "gate-guardrail",
			oldPythonTest: "tests/test_gate_guardrail.py",
			intent:        "FO halts at a human gate and presents the review without self-approval, mutation, or archival.",
		},
		{
			name:          "rejection-flow",
			oldPythonTest: "tests/test_rejection_flow.py",
			intent:        "FO drives a two-cycle rejection trajectory — route back, re-implement, re-validate via reviewer reuse — restoring the dropped second cycle.",
		},
		{
			name:          "feedback-3-cycle-escalation",
			oldPythonTest: "tests/test_rejection_flow.py",
			intent:        "On the 3rd consecutive REJECTED validation the FO escalates to the human instead of auto-bouncing a 4th time.",
		},
		{
			name:          "merge-hook-guardrail",
			oldPythonTest: "tests/test_merge_hook_guardrail.py",
			intent:        "FO cannot bypass a registered merge hook by terminalizing without pr, mod-block, or force.",
		},
		{
			name:          "filing",
			oldPythonTest: "n/a (new behavior — `spacedock new` adopted post-Python port)",
			intent:        "FO files a new seed entity via the atomic `spacedock new <slug>` path, not the drift-prone `--next-id` + hand-write pair.",
		},
	}
}
