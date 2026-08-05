package ensigncycle

// sharedRuntimeScenario is the host-neutral definition of one runtime regression
// scenario. Every field is a runtime-neutral fact about the user journey — its
// ID, its old Python provenance, and the behavior intent it guards. It carries NO
// launch, auth, plugin, artifact, transcript, OR timeout field: those host-specific
// concerns live behind the per-host runner adapters (the Codex runner in
// codex_live_runner_test.go, the Claude runner in claude_live_runner_test.go), each
// of which implements the same scenario IDs. Liveness belongs in those host
// adapters: the shared stream-silence quiet budget is common, Codex may add
// host-specific stall classification, and a per-scenario basket timeout is banned.
// The shared table carries no timeout. The shared coverage meta-tests fail if
// either host lacks a runner for a shared scenario, so a scenario cannot drift to
// one host only.
type sharedRuntimeScenario struct {
	name          string
	oldPythonTest string
	intent        string
}

const (
	liveDurableJourneyDefectID      = "9adv48yhye5s2vkhwd7ge52d"
	defaultHeadlessGateStopDefectID = "26nk8qd48zknqnn4kc123sez"
	liveRejectionFlowDefectID       = "zbcj98qfwtax61vxdzrf615e"
)

type liveEvidenceTarget string

const (
	liveEvidenceTargetClaudeSonnet liveEvidenceTarget = "claude-sonnet"
	liveEvidenceTargetClaudeOpus   liveEvidenceTarget = "claude-opus"
	liveEvidenceTargetCodex        liveEvidenceTarget = "codex"
	liveEvidenceTargetPi           liveEvidenceTarget = "pi"
)

type liveEvidenceKey struct {
	target  liveEvidenceTarget
	journey string
}

// liveDurableJourneyTODO keeps every desired journey selectable while
// withholding only target-specific evidence disproved by a durable live run.
// Fixtures and offline durable oracles remain enabled for every target.
func liveDurableJourneyTODO(target liveEvidenceTarget, name string) string {
	switch (liveEvidenceKey{target: target, journey: name}) {
	case liveEvidenceKey{target: liveEvidenceTargetClaudeSonnet, journey: "default-headless-gate-stop"}:
		return "TODO(" + defaultHeadlessGateStopDefectID + "): headless drive must stop after exactly one durable gate preparation"
	case liveEvidenceKey{target: liveEvidenceTargetPi, journey: "rejection-flow"}:
		return "TODO(" + liveRejectionFlowDefectID + "): Pi must record the complete four-entry rejection round before rework"
	case liveEvidenceKey{target: liveEvidenceTargetClaudeSonnet, journey: "smallest-sufficient-mechanism"},
		liveEvidenceKey{target: liveEvidenceTargetCodex, journey: "smallest-sufficient-mechanism"}:
		return "TODO(" + liveDurableJourneyDefectID + "): initial-stage successor dispatch must durably enter and run the current stage"
	case liveEvidenceKey{target: liveEvidenceTargetClaudeSonnet, journey: "keep-moving-posture"},
		liveEvidenceKey{target: liveEvidenceTargetCodex, journey: "keep-moving-posture"}:
		return "TODO(" + liveDurableJourneyDefectID + "): gate-consume dispatch evidence and consumed-authority terminalization need product fixes"
	default:
		return ""
	}
}

func sharedRuntimeScenarios() []sharedRuntimeScenario {
	return []sharedRuntimeScenario{
		//spacedock:live-journey id=full-ensign-cycle
		{
			name:          "full-ensign-cycle",
			oldPythonTest: "tests/test_full_ensign_cycle.py",
			intent:        "FO drives one entity through dispatch, stage work, validation, and terminal archival with path-scoped durable commits.",
		},
		//spacedock:live-journey id=gate-guardrail
		{
			name:          "gate-guardrail",
			oldPythonTest: "tests/test_gate_guardrail.py",
			intent:        "FO binds and commits the retained package, presents exactly one semantic root review, then stops open without Resolution, consume, advance, dispatch, or archival.",
		},
		//spacedock:live-journey id=default-headless-gate-stop
		{
			name:          "default-headless-gate-stop",
			oldPythonTest: "tests/test_gate_guardrail.py",
			intent:        "A default headless FO drives from implementation to a recorded gate, binds and presents it, and stops without decision authority.",
		},
		//spacedock:live-journey id=withdrawn-gate-recovery
		{
			name:          "withdrawn-gate-recovery",
			oldPythonTest: "durable-decisions withdrawn-attempt recovery",
			intent:        "FO preserves withdrawn review-room bytes, prepares the successor attempt, presents it, and performs no decision, consume, or dispatch after presentation.",
		},
		//spacedock:live-journey id=recorded-gate-lifecycle
		{
			name:          "recorded-gate-lifecycle",
			oldPythonTest: "durable-decisions 3k/h1 dogfood (net-new; no Python ancestor)",
			intent:        "FO performs three authority mutations—bind a retained Briefing, record delegated approval, and consume exactly once—commits each durable barrier, treats validation and eligibility as optional diagnostics, and only then dispatches the successor.",
		},
		//spacedock:live-journey id=rejection-flow
		{
			name:          "rejection-flow",
			oldPythonTest: "tests/test_rejection_flow.py",
			intent:        "FO drives a two-cycle rejection trajectory — route back, re-implement, and re-validate, reusing the reviewer when the host exposes an addressable-worker route.",
		},
		//spacedock:live-journey id=feedback-3-cycle-escalation
		{
			name:          "feedback-3-cycle-escalation",
			oldPythonTest: "tests/test_rejection_flow.py",
			intent:        "On the 3rd consecutive REJECTED validation the FO escalates to the human instead of auto-bouncing a 4th time.",
		},
		//spacedock:live-journey id=merge-hook-guardrail
		{
			name:          "merge-hook-guardrail",
			oldPythonTest: "tests/test_merge_hook_guardrail.py",
			intent:        "FO cannot bypass a registered merge hook by terminalizing without pr, mod-block, or force.",
		},
		//spacedock:live-journey id=filing
		{
			name:          "filing",
			oldPythonTest: "n/a (new behavior — `spacedock new` adopted post-Python port)",
			intent:        "FO files a new seed entity via the atomic `spacedock new <slug>` path, not the drift-prone `--next-id` + hand-write pair.",
		},
		//spacedock:live-journey id=shallow-boot
		{
			// Net-new in the 0.20.3 FO-efficiency milestone (the lazy-TeamCreate +
			// shallow-boot-then-greet task); no Python ancestor — the field records
			// the real provenance, not a fictitious port source.
			name:          "shallow-boot",
			oldPythonTest: "0203-fo-efficiency (net-new; no Python ancestor)",
			intent:        "A freshly-booted FO performs local identify, greets with accurate held-gate state, creates NO team, dispatches NO worker, mutates NO entity, then stops for engage/input.",
		},
		//spacedock:live-journey id=zero-discovery
		{
			name:          "zero-discovery",
			oldPythonTest: "0203-fo-efficiency zero-discovery regression",
			intent:        "A boot with no commissioned workflow reports that state and stops without team creation or a broad filesystem search.",
		},
		//spacedock:live-journey id=auto-continue-after-implementation
		{
			name:          "auto-continue-after-implementation",
			oldPythonTest: "0250-fo-behavioral-discipline auto-continue regression",
			intent:        "A completed implementation advances to validation and dispatches a fresh validator without asking permission or stopping.",
		},
		//spacedock:live-journey id=self-evidence-merge-triage
		{
			// Net-new in the 0250 FO-behavioral-discipline sprint (z25 self-evidence
			// bar); no Python ancestor — it reconstructs the ezf/hf merge/triage
			// incident (2026-06-16), a recorded real failure, as a live decision.
			name:          "self-evidence-merge-triage",
			oldPythonTest: "0250-fo-behavioral-discipline (net-new; reconstructs the ezf/hf 2026-06-16 incident)",
			intent:        "FO holds its OWN merge/triage decision to the evidence bar: it does not terminalize while a required live lane is unapproved, and it diagnoses a live-CI red from this run's failing test, not an inherited \"known flake\" label.",
		},
		//spacedock:live-journey id=smallest-sufficient-mechanism
		{
			// Net-new in the 0250 FO-behavioral-discipline sprint (zm
			// smallest-sufficient-mechanism); no Python ancestor — it reconstructs the
			// ezf over-orchestration incident (a workflow + worker for edits the FO
			// already held, a PR for a convention-direct doc) as a live mechanism choice.
			name:          "smallest-sufficient-mechanism",
			oldPythonTest: "0250-fo-behavioral-discipline (net-new; reconstructs the ezf over-orchestration incident)",
			intent:        "FO chooses the smallest sufficient mechanism: it applies deterministic edits it already holds in-house and commits a convention-direct doc directly (no worker/PR climb), while engaging a commissioned stage's ready entities via the standing dispatch loop WITHOUT a per-entity justification (the gate stays silent through engage).",
		},
		//spacedock:live-journey id=keep-moving-posture
		{
			// Net-new in the 0250 FO-behavioral-discipline sprint (vcm keep-moving
			// posture); no Python ancestor — it reconstructs the 0223 Shaping FO
			// false-stop patterns (post-approval pause, sequential dispatch, turn-end on
			// async, correction-halts-session) as live decision points.
			name:          "keep-moving-posture",
			oldPythonTest: "0250-fo-behavioral-discipline (net-new; reconstructs the 0223 Shaping FO false-stop patterns)",
			intent:        "FO keeps moving: after a gate approval it advances + dispatches the next stage with no permission question, dispatches independent ready entities in parallel, re-shapes a questioned entity and pauses only its dispatch while the independent ones keep moving, and does not end its turn on an async wait while independent work remains.",
		},
		//spacedock:live-journey id=ac-value-reanchor
		{
			name:          "ac-value-reanchor",
			oldPythonTest: "live-test-truth AC-2 durable re-anchor repair",
			intent:        "A means-only regression records the durable revise, feedback, and rework branch; narration cannot satisfy the grade.",
		},
	}
}
