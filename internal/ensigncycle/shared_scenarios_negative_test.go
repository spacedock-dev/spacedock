package ensigncycle

import (
	"strings"
	"testing"
)

// Negative-case discipline: each shared scenario's assertion is behavior/state
// oriented, not a transcript-shape tautology. For every shared scenario these cases build the
// SPECIFIC broken end-state the scenario guards against — from the real shared
// fixture, not an arbitrary string — and prove the assertion goes red. A
// tautological assertion (one that only checks the transcript echoed a phrase)
// would stay green on these broken states; these cases fail it loudly.
//
// They are offline (default tag): the assertions are pure functions over
// entity-state + observed strings, so the negative cases spend no model.

func TestGateGuardrailNegativeBrokenStateTransition(t *testing.T) {
	before := gateEntity()
	// The held-gate baseline passes: the entity is unmutated, still at review, and
	// the final message presents the gate review + decision.
	heldFinal := "Gate review: Gate Check at review.\nDecision: approve or reject?"
	if err := assertGateHeld(before, before, heldFinal); err != nil {
		t.Fatalf("held-gate baseline must pass: %v", err)
	}

	// Broken state transition: the FO advanced the gated entity to done instead of
	// presenting the gate. Even if it ALSO emitted a gate-review-shaped final
	// message, the state assertion must catch the advanced status.
	advanced := strings.Replace(before, "status: review", "status: done", 1)
	if advanced == before {
		t.Fatal("fixture must contain `status: review` to advance")
	}
	if err := assertGateHeld(before, advanced, heldFinal); err == nil {
		t.Fatal("expected a gate advanced to status: done to fail assertGateHeld even with a gate-review final message")
	}

	// Broken state transition: the FO self-approved by setting a verdict on the
	// gated entity.
	withVerdict := strings.Replace(before, "verdict:\n", "verdict: passed\n", 1)
	if withVerdict == before {
		t.Fatal("fixture must contain an empty `verdict:` line")
	}
	if err := assertGateHeld(before, withVerdict, heldFinal); err == nil {
		t.Fatal("expected a self-approved (verdict set) gate to fail assertGateHeld")
	}
}

func TestRejectionFlowNegativeSingleCycle(t *testing.T) {
	rejectedObserved := "validation was REJECTED; routing the finding back to implementation"

	// Un-driven fixture: the rejection scenario now starts BEFORE the first
	// validation, at status: implementation with NO stage reports and NO seeded
	// rejection. The seeded fixture must NOT pre-satisfy assertRejectionFlow — a live
	// pass requires the real producer to drive BOTH cycles (omit the fix, get
	// rejected, rework, re-validate). If this seeded state passed, a live run that did
	// nothing would falsely pass.
	seeded := rejectionEntity()
	if !strings.Contains(seeded, "status: implementation") {
		t.Fatal("rejection fixture must now start at status: implementation, before the first validation")
	}
	if got := len(implementationReport.FindAllString(seeded, -1)); got != 0 {
		t.Fatalf("rejection fixture must start with no implementation reports (live producer writes them), got %d", got)
	}
	if err := assertRejectionFlow(seeded, rejectedObserved); err == nil {
		t.Fatal("expected the un-driven rejection fixture (no reports, no cycles) to fail assertRejectionFlow")
	}

	// No-reviewer-created shape — the exact flaw that shipped on PR #302. The OLD
	// fixture pre-wrote a `## Stage Report: validation` REJECTED, so cycle-1
	// validation never ran live, no reviewer was ever spawned, and the cycle-2
	// reviewer-reuse signal was unreachable: the FO correctly fresh-dispatched. The
	// redesigned fixture must NOT pre-contain a validation report, so the FO drives a
	// real cycle-1 validation that spawns a reviewer to keep alive and reuse.
	if strings.Contains(seeded, "## Stage Report: validation") {
		t.Fatal("rejection fixture must NOT pre-contain a validation report — a pre-written cycle-1 rejection means no live reviewer is ever spawned, making the cycle-2 reviewer-reuse signal unreachable (the PR #302 regression)")
	}

	// AC-4 single-cycle end-state — the Go-port regression the evolved scenario
	// restores: the FO applied the fix and left a SECOND implementation report, but
	// stopped after one cycle, never driving the second validation round (only one
	// recorded cycle). The two-implementation-report check passes, so this MUST fail
	// on the second-cycle check — proving the evolved assertion catches the
	// single-route-back simplification the Python test never had.
	singleCycle := "---\nstatus: implementation\n---\n# Rejection Task\n\n" +
		rejectionFixMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: initial (no marker)\n\n" +
		"## Stage Report: implementation\n\n- DONE: applied fix\n\n" +
		"### Feedback Cycles\n\n- Cycle 1: REJECTED\n"
	if len(implementationReport.FindAllString(singleCycle, -1)) < 2 {
		t.Fatal("single-cycle body must carry a second implementation report")
	}
	if len(feedbackCycleEntry.FindAllString(singleCycle, -1)) != 1 {
		t.Fatal("single-cycle body must carry exactly one recorded cycle (only Cycle 1)")
	}
	if err := assertRejectionFlow(singleCycle, rejectedObserved); err == nil {
		t.Fatal("expected a single-cycle end-state (fix applied, second implementation report, but only one recorded cycle) to fail assertRejectionFlow on the second-cycle check")
	}

	// No-reuse run shape — the producer-signal half of the shipped flaw. A run whose
	// transcript never carries a reuse tool-call (because no reviewer was kept alive
	// to reuse) must RED on the host reuse assertions. This is the offline pin for
	// "a run that never creates-or-reuses a reviewer"; the live legs grade the real
	// transcript.
	noReuseClaude := `{"type":"assistant","message":{"content":[{"type":"text","text":"fresh-dispatching a new validation reviewer; no prior reviewer to reuse"}]}}`
	if err := assertClaudeReviewerReuse(noReuseClaude); err == nil {
		t.Fatal("expected a transcript that never reuses a reviewer to fail assertClaudeReviewerReuse")
	}
	noReuseCodex := `{"type":"message","role":"assistant","content":"fresh-dispatching a new validation worker; no prior worker to reuse"}`
	if err := assertCodexReviewerReuse(noReuseCodex); err == nil {
		t.Fatal("expected a transcript that never reuses a reviewer to fail assertCodexReviewerReuse")
	}
}

func TestThirdCycleEscalationNegativeAutoBounce(t *testing.T) {
	// The escalated end-state the live run must reach passes: the real fixture plus
	// the third cycle entry and the escalation marker, with NO new implementation
	// report — the FO parked for the human instead of bouncing a fourth time.
	escalated := escalationEntity() +
		"- Cycle 3: REJECTED — third consecutive rejection.\n" +
		escalationMarker + "\n"
	if err := assertThirdCycleEscalation(escalated); err != nil {
		t.Fatalf("escalated baseline must pass: %v", err)
	}

	// Broken end-state — 4th auto-bounce: built from the REAL fixture, the FO
	// recorded a third cycle but routed back to implementation a fourth time (a new
	// implementation report) instead of escalating, and recorded no marker. The
	// state assertion must catch the extra implementation report even though the
	// body still mentions three rejection rounds.
	autoBounced := escalationEntity() +
		"- Cycle 3: REJECTED — routed back to implementation again.\n\n" +
		"## Stage Report: implementation\n\n- DONE: reworked a fourth time\n"
	if implementationReport.MatchString(escalationEntity()) {
		if len(implementationReport.FindAllString(autoBounced, -1)) != 2 {
			t.Fatal("4th-auto-bounce body must carry two implementation reports built from the real fixture")
		}
	}
	if err := assertThirdCycleEscalation(autoBounced); err == nil {
		t.Fatal("expected a 4th auto-bounce (third cycle routed back + a new implementation report, no marker) to fail assertThirdCycleEscalation")
	}

	// Broken end-state — stalled at cycle 2: the real fixture as-staged carries only
	// the two seeded cycle entries and no escalation marker — the FO never reached
	// the third-cycle decision. Must fail on the cycle-count check, not pass on any
	// transcript shape.
	stalled := escalationEntity()
	if got := len(feedbackCycleEntry.FindAllString(stalled, -1)); got != 2 {
		t.Fatalf("escalation fixture must start with exactly two seeded `### Feedback Cycles` entries, got %d", got)
	}
	if err := assertThirdCycleEscalation(stalled); err == nil {
		t.Fatal("expected a stalled-at-cycle-2 end-state (only two cycle entries, no marker) to fail assertThirdCycleEscalation")
	}

	// Isolating case for the no-post-cycle-3-report check: marker present AND three
	// recorded cycles (so the cycle-count and marker checks BOTH pass), but a stray
	// post-cycle-3 `## Stage Report: implementation` — the ONLY defect. This must
	// still fail, and it fails ONLY on the report-count check. Without this case,
	// deleting that check leaves the suite green: every OTHER escalation negative
	// that carries a stray report also lacks the marker, so they red on the marker
	// check first and never exercise the report-count clause. This is the one
	// assertion proving the FO did not auto-bounce a fourth time, so it must be
	// independently covered.
	markerWithStrayReport := escalationEntity() +
		"- Cycle 3: REJECTED — third consecutive rejection.\n" +
		escalationMarker + "\n\n" +
		"## Stage Report: implementation\n\n- DONE: stray fourth-round report\n"
	if !strings.Contains(markerWithStrayReport, escalationMarker) {
		t.Fatal("report-count isolating case must carry the escalation marker so it passes the marker check")
	}
	if got := len(feedbackCycleEntry.FindAllString(markerWithStrayReport, -1)); got < 3 {
		t.Fatalf("report-count isolating case must carry at least three cycle entries, got %d", got)
	}
	if err := assertThirdCycleEscalation(markerWithStrayReport); err == nil {
		t.Fatal("expected a marker-present, three-cycle end-state with a stray post-cycle-3 implementation report to fail assertThirdCycleEscalation on the report-count check")
	}

	// Isolating case for the park-not-advance (non-terminal) check: marker present,
	// three recorded cycles, exactly one implementation report (so the cycle-count,
	// marker, and report-count checks ALL pass), but the FO terminalized the entity
	// to status: done — escalate-then-terminalize, auto-resolving instead of parking
	// for the human (the escalation prompt says "do not advance to done"). This must
	// still fail, and it fails ONLY on the non-terminal check.
	markerButTerminalized := strings.Replace(escalationEntity(), "status: validation", "status: done", 1) +
		"- Cycle 3: REJECTED — third consecutive rejection.\n" +
		escalationMarker + "\n"
	if !strings.Contains(markerButTerminalized, "status: done") {
		t.Fatal("non-terminal isolating case must carry status: done")
	}
	if got := len(implementationReport.FindAllString(markerButTerminalized, -1)); got != 1 {
		t.Fatalf("non-terminal isolating case must carry exactly one implementation report so only the non-terminal check rejects it, got %d", got)
	}
	if err := assertThirdCycleEscalation(markerButTerminalized); err == nil {
		t.Fatal("expected a marker-present but terminalized (status: done) end-state to fail assertThirdCycleEscalation on the park-not-advance check")
	}
}

func TestMergeHookGuardrailNegativeBypass(t *testing.T) {
	before := mergeHookGuardEntity()
	guardObserved := "Error: entity merge-check cannot advance to terminal - workflow has merge hook(s) [local-merge]"
	// The held-guard baseline passes: entity unmutated, still implementation, and
	// the observed output named the merge hook + terminal guard refusal.
	if err := assertMergeHookGuardHeld(before, before, guardObserved); err != nil {
		t.Fatalf("held merge-hook guard baseline must pass: %v", err)
	}

	// Merge-hook bypass: the FO terminalized the entity to done despite the
	// registered hook. The state assertion must catch the advanced status even if
	// the observed transcript still mentions a merge hook.
	bypassed := strings.Replace(before, "status: implementation", "status: done", 1)
	if bypassed == before {
		t.Fatal("merge-hook fixture must contain `status: implementation`")
	}
	if err := assertMergeHookGuardHeld(before, bypassed, guardObserved); err == nil {
		t.Fatal("expected a terminalized (status: done) entity to fail assertMergeHookGuardHeld even with merge-hook mention in the transcript")
	}

	// Bypass with no guard error in the observed output: the FO advanced and the
	// run never surfaced the guard refusal. Must fail on the missing guard signal.
	if err := assertMergeHookGuardHeld(before, before, "terminalized merge-check to done"); err == nil {
		t.Fatal("expected a run with no terminal-guard refusal in observed to fail assertMergeHookGuardHeld")
	}

	// Isolating case for the `cannot advance to terminal` guard-error check: the
	// observed output mentions a merge hook (so the mention check passes) but never
	// reports the terminal-guard refusal — the FO touched the hook yet the guard
	// never FIRED. This must still fail, and it fails ONLY on the guard-error
	// check. Without this case, deleting that check leaves the suite green: every
	// other merge-hook negative observed string also lacks the merge-hook mention,
	// so they trip the earlier mention check and never exercise the guard-error
	// clause. This is the one assertion proving the merge guard actually fired, so
	// it must be independently covered.
	hookMentionedNoGuard := "Inspected startup: workflow registers a merge hook [local-merge]. Proceeding without terminalization."
	if err := assertMergeHookGuardHeld(before, before, hookMentionedNoGuard); err == nil {
		t.Fatal("expected observed that mentions a merge hook but omits the terminal-guard refusal to fail assertMergeHookGuardHeld on the guard-error check")
	}
}

func TestShallowBootNegativeBrokenEndStates(t *testing.T) {
	// The realized shallow-boot end-state passes: no team config, the gate entity
	// unchanged, the merged entity advanced+archived, the greet present.
	gate := shallowBootGateEntity()
	mergedArchived := "---\nid: merged-pr\nstatus: done\ncompleted: 2026-06-13T00:00:00Z\nverdict: PASSED\npr: \"#42\"\nmod-block:\nworktree:\n---\n"
	greet := "Workflow overview: 1 task at the review gate; merged-pr advanced (PR #42 merged).\nGate review: Gate Check at review.\nDecision: approve or reject?"
	good := shallowBootObservation{
		finalMessage: greet, gateBefore: gate, gateAfter: gate,
		mergedAfter: mergedArchived, mergedArchived: true,
	}
	if err := assertShallowBoot(good); err != nil {
		t.Fatalf("the realized shallow-boot end-state must pass: %v", err)
	}

	// Broken: a team config landed on disk — lazy-TeamCreate was not honored. The
	// eager team is the exact regression P2 prevents.
	eagerTeam := good
	eagerTeam.teamConfigOnDisk = true
	if err := assertShallowBoot(eagerTeam); err == nil {
		t.Fatal("expected a team config.json on disk (eager TeamCreate) to fail assertShallowBoot")
	}

	// Broken: a worker was dispatched — the gate entity was advanced past its gate.
	dispatched := good
	dispatched.gateAfter = strings.Replace(gate, "status: review", "status: done", 1)
	if dispatched.gateAfter == gate {
		t.Fatal("gate fixture must contain `status: review`")
	}
	if err := assertShallowBoot(dispatched); err == nil {
		t.Fatal("expected a dispatched (gate advanced) end-state to fail assertShallowBoot")
	}

	// Broken: the FO self-approved the gate (verdict set) instead of presenting it.
	selfApproved := good
	selfApproved.gateAfter = strings.Replace(gate, "verdict:\n", "verdict: passed\n", 1)
	if selfApproved.gateAfter == gate {
		t.Fatal("gate fixture must contain an empty `verdict:` line")
	}
	if err := assertShallowBoot(selfApproved); err == nil {
		t.Fatal("expected a self-approved gate (verdict set) end-state to fail assertShallowBoot")
	}

	// Broken: a worktree was created for the gate entity — a dispatch happened.
	worktreeCreated := good
	worktreeCreated.gateWorktreeCreated = true
	if err := assertShallowBoot(worktreeCreated); err == nil {
		t.Fatal("expected a worktree created for the gated entity to fail assertShallowBoot")
	}

	// Broken: the merged-PR entity was NOT advanced (S7b skipped) — still active,
	// not archived. This is the M3 failure a greet-and-stop boot would have without
	// the before-greet sweep.
	s7bSkipped := good
	s7bSkipped.mergedArchived = false
	s7bSkipped.mergedAfter = shallowBootMergedEntity() // still at implementation, no verdict
	if err := assertShallowBoot(s7bSkipped); err == nil {
		t.Fatal("expected an un-advanced merged-PR entity (S7b skipped) to fail assertShallowBoot")
	}

	// Isolating: archived but verdict not set — advancement incomplete. Isolates the
	// verdict check from the archived check so neither can be silently dropped.
	noVerdict := good
	noVerdict.mergedAfter = strings.Replace(mergedArchived, "verdict: PASSED", "verdict:", 1)
	if err := assertShallowBoot(noVerdict); err == nil {
		t.Fatal("expected an archived-but-no-verdict merged entity to fail assertShallowBoot on the verdict check")
	}

	// Isolating: advanced+archived but a mod-block still set — the clear was skipped.
	modBlockLeft := good
	modBlockLeft.mergedAfter = strings.Replace(mergedArchived, "mod-block:\n", "mod-block: merge:pr-merge\n", 1)
	if err := assertShallowBoot(modBlockLeft); err == nil {
		t.Fatal("expected an advanced entity with a lingering mod-block to fail assertShallowBoot on the mod-block-clear check")
	}

	// Broken: no greet — the final message lacks the gate review / decision prompt.
	noGreet := good
	noGreet.finalMessage = "Advanced merged-pr; nothing else to do."
	if err := assertShallowBoot(noGreet); err == nil {
		t.Fatal("expected a final message with no gate review/decision prompt to fail assertShallowBoot")
	}
}
