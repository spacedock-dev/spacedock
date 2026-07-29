package ensigncycle

import (
	"strings"
	"testing"
)

// The mandatory offline negative for the keep-moving-posture scenario, plus the cycle-1
// replay fixtures. It proves the grader holds the invariant triple for a captain correction
// (corrected entity does not drive FORWARD until folded; everything else keeps moving; the
// folded re-shape SURFACES or its rework is honestly named in-flight) from the REAL fixture
// tokens: the correct trace PASSES, and every incident reds in isolation — a serialized
// independent dispatch, a dropped re-shape, a FORWARD drive of the corrected entity, and a
// silent wait/park that never surfaces the re-shape.
// The re-shape MECHANISM is free: a routed rework dispatch of the corrected entity is NOT a
// violation (only a forward status is). The replay fixtures distil the three recorded
// cycle-1 drives (run 28843487385): codex green (surfaced a gate review), sonnet green
// (honestly named the in-flight rework), opus red (silent wait — the mini false-stop this
// entity targets). Offline (default tag): the grader and extractors are pure functions.

// kmCorrectFinal is a keep-moving final message: it reports the advance, parallel dispatch,
// and re-shape, and SURFACES the corrected entity's re-shape for gate review — no permission
// question, no silent park.
func kmCorrectFinal() string {
	return "Advanced " + kmApprovedGate + " to " + kmNextStage + " and dispatched it. " +
		"Dispatched " + kmReadyOne + " and " + kmReadyTwo + " in parallel in one turn. " +
		"Re-shaped " + kmQuestioned + " to fold the captain's symlink correction and presented it for gate review, " +
		"holding it from advancing until the correction is folded."
}

// kmSilentWaitFinal is the opus cycle-1 specimen verbatim: everything dispatched, then a
// silent wait that never names the corrected entity's re-shape — the invariant's silent-park
// failure (a mini false-stop, this scenario's own target class).
func kmSilentWaitFinal() string {
	return "All four workers dispatched as named background agents. No completion signals yet and no independent FO work remains, so I'll wait for their completion notifications before advancing."
}

// kmClaudeCorrectStream builds the tool-call stream of a correct Claude run: a `status --set`
// advance of the approved entity, an Agent dispatch of its next stage, an Agent dispatch of
// each independent entity, and an Edit re-shaping the corrected entity in-house.
func kmClaudeCorrectStream() string {
	lines := []string{
		claudeToolUse("Bash", `{"command":"spacedock status --workflow-dir . --set `+kmApprovedGate+` status=`+kmNextStage+`"}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch the `+kmNextStage+` stage for `+kmApprovedGate+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyOne+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyTwo+`."}`),
		claudeToolUse("Edit", `{"file_path":"`+kmQuestioned+`.md"}`),
	}
	return strings.Join(lines, "\n")
}

func TestGradeKeepMoving(t *testing.T) {
	independent := kmIndependent()

	kmCorrectTrace := func() keepMovingTrace {
		tr := newKeepMovingTrace()
		tr.approvedAdvanced = true
		tr.approvedDispatched = true
		for _, e := range independent {
			tr.independentDispatched[e] = true
		}
		tr.correctedReshaped = true
		tr.correctedAddressed = true
		return tr
	}

	// Positive: the correct trace passes.
	if err := gradeKeepMoving(kmCorrectTrace(), independent); err != nil {
		t.Fatalf("the correct keep-moving trace must pass: %v", err)
	}

	// Negative (S1): the approved entity was never advanced.
	noAdvance := kmCorrectTrace()
	noAdvance.approvedAdvanced = false
	if err := gradeKeepMoving(noAdvance, independent); err == nil {
		t.Fatal("expected a missing advance of the approved entity to fail")
	}

	// Negative (S1 isolating): advanced but never dispatched.
	noDispatch := kmCorrectTrace()
	noDispatch.approvedDispatched = false
	if err := gradeKeepMoving(noDispatch, independent); err == nil {
		t.Fatal("expected an advance-without-dispatch to fail")
	}

	// Negative (S2 serialize): an independent ready entity never dispatched.
	serialized := kmCorrectTrace()
	serialized.independentDispatched[kmReadyTwo] = false
	if err := gradeKeepMoving(serialized, independent); err == nil {
		t.Fatal("expected an undispatched independent entity (serialized) to fail")
	}

	// Negative (S4): the corrected entity was never re-shaped.
	noReshape := kmCorrectTrace()
	noReshape.correctedReshaped = false
	if err := gradeKeepMoving(noReshape, independent); err == nil {
		t.Fatal("expected a dropped re-shape of the corrected entity to fail")
	}

	// Negative (S4 forward-drive): the corrected entity was driven forward before folding.
	driven := kmCorrectTrace()
	driven.correctedDriven = true
	if err := gradeKeepMoving(driven, independent); err == nil {
		t.Fatal("expected driving the corrected entity forward to fail")
	}

	// Negative (S4 silent park): the re-shape was neither surfaced nor honestly named in-flight.
	unaddressed := kmCorrectTrace()
	unaddressed.correctedAddressed = false
	if err := gradeKeepMoving(unaddressed, independent); err == nil {
		t.Fatal("expected a silently-parked re-shape to fail")
	}

	// Negative (S4 halt — pattern 4): the correction halted the independent work.
	halted := kmCorrectTrace()
	for _, e := range independent {
		halted.independentDispatched[e] = false
	}
	if err := gradeKeepMoving(halted, independent); err == nil {
		t.Fatal("expected a correction that halted the independent work to fail")
	}
}

func TestAssertClaudeKeepMoving(t *testing.T) {
	independent := kmIndependent()

	// Positive: the correct stream + a surfacing final message passes.
	if err := assertClaudeKeepMoving(kmClaudeCorrectStream(), kmCorrectFinal(), independent); err != nil {
		t.Fatalf("the correct Claude keep-moving trace must pass: %v", err)
	}

	// Negative (S4 silent park): correct actions, a silent-wait final that never names the
	// corrected entity — the re-shape is silently parked.
	if err := assertClaudeKeepMoving(kmClaudeCorrectStream(), kmSilentWaitFinal(), independent); err == nil {
		t.Fatal("expected a silent-wait final message to fail the Claude assertion")
	}

	// Negative (S4 forward-drive): the FO advanced the corrected entity forward (to
	// implementation) before the re-shape folded.
	driven := kmClaudeCorrectStream() + "\n" +
		claudeToolUse("Bash", `{"command":"spacedock status --workflow-dir . --set `+kmQuestioned+` status=`+kmNextStage+`"}`)
	if err := assertClaudeKeepMoving(driven, kmCorrectFinal(), independent); err == nil {
		t.Fatal("expected driving the corrected entity forward to fail the Claude assertion")
	}

	// Positive (mechanism-free re-shape): the corrected entity re-shaped via a ROUTED rework
	// dispatch (Agent naming it) instead of an in-house edit — legitimate, must not be read as
	// a forward drive.
	routedReshape := strings.Join([]string{
		claudeToolUse("Bash", `{"command":"spacedock status --workflow-dir . --set `+kmApprovedGate+` status=`+kmNextStage+`"}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch the `+kmNextStage+` stage for `+kmApprovedGate+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyOne+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyTwo+`."}`),
		claudeToolUse("Agent", `{"prompt":"Rework `+kmQuestioned+` at `+kmReopenStage+` to fold the correction.","description":"`+kmQuestioned+`: `+kmReopenStage+` re-shape"}`),
	}, "\n")
	if err := assertClaudeKeepMoving(routedReshape, kmCorrectFinal(), independent); err != nil {
		t.Fatalf("a routed rework re-shape of the corrected entity must pass: %v", err)
	}

	// Positive (bulk-command scoping regression): a correct run that sets the approved advance
	// AND the corrected entity's re-open in ONE bulk `--set … --set …` Bash command (the opus
	// cycle-1 shape). The segment-scoped extractor must attribute status=implementation to
	// approved-gate and status=ideation to questioned — never cross-attributing the forward
	// status to the corrected entity as a false forward-drive.
	bulk := strings.Join([]string{
		claudeToolUse("Bash", `{"command":"spacedock status --set `+kmApprovedGate+` status=`+kmNextStage+` ; spacedock status --set `+kmQuestioned+` status=`+kmReopenStage+`"}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch the `+kmNextStage+` stage for `+kmApprovedGate+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyOne+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyTwo+`."}`),
		claudeToolUse("Edit", `{"file_path":"`+kmQuestioned+`.md"}`),
	}, "\n")
	if err := assertClaudeKeepMoving(bulk, kmCorrectFinal(), independent); err != nil {
		t.Fatalf("a correct run using a bulk multi---set command must pass (segment-scoping regression): %v", err)
	}

	// Negative (S1): no advance of the approved entity.
	noAdvance := strings.Join([]string{
		claudeToolUse("Agent", `{"prompt":"Dispatch the `+kmNextStage+` stage for `+kmApprovedGate+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyOne+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyTwo+`."}`),
		claudeToolUse("Edit", `{"file_path":"`+kmQuestioned+`.md"}`),
	}, "\n")
	if err := assertClaudeKeepMoving(noAdvance, kmCorrectFinal(), independent); err == nil {
		t.Fatal("expected a stream with no advance to fail the Claude assertion")
	}

	// Negative (S2 serialize): ready-two never dispatched.
	serialized := strings.Join([]string{
		claudeToolUse("Bash", `{"command":"spacedock status --workflow-dir . --set `+kmApprovedGate+` status=`+kmNextStage+`"}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch the `+kmNextStage+` stage for `+kmApprovedGate+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyOne+`."}`),
		claudeToolUse("Edit", `{"file_path":"`+kmQuestioned+`.md"}`),
	}, "\n")
	if err := assertClaudeKeepMoving(serialized, kmCorrectFinal(), independent); err == nil {
		t.Fatal("expected a serialized dispatch (ready-two missing) to fail the Claude assertion")
	}

	// Negative (S4): the corrected entity was never re-shaped.
	noReshape := strings.Join([]string{
		claudeToolUse("Bash", `{"command":"spacedock status --workflow-dir . --set `+kmApprovedGate+` status=`+kmNextStage+`"}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch the `+kmNextStage+` stage for `+kmApprovedGate+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyOne+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyTwo+`."}`),
	}, "\n")
	if err := assertClaudeKeepMoving(noReshape, kmCorrectFinal(), independent); err == nil {
		t.Fatal("expected a stream with no re-shape to fail the Claude assertion")
	}
}

// kmCodexCorrectStream builds the transcript of a correct Codex run: a `status --set`
// advance, a spawn_agent dispatch of the approved next stage and of each independent entity,
// and a file_change re-shaping the corrected entity (the codex 0.142.5 surface).
func kmCodexCorrectStream() string {
	lines := []string{
		codexCommand("spacedock status --workflow-dir . --set " + kmApprovedGate + " status=" + kmNextStage),
		codexSpawn("Dispatch the " + kmNextStage + " stage for " + kmApprovedGate + "."),
		codexSpawn("Dispatch " + kmNextStage + " for " + kmReadyOne + "."),
		codexSpawn("Dispatch " + kmNextStage + " for " + kmReadyTwo + "."),
		codexFileChange(kmQuestioned + ".md"),
	}
	return strings.Join(lines, "\n")
}

func TestAssertCodexKeepMoving(t *testing.T) {
	independent := kmIndependent()

	// Positive: the correct Codex transcript + a surfacing final message passes.
	if err := assertCodexKeepMoving(kmCodexCorrectStream(), kmCorrectFinal(), independent); err != nil {
		t.Fatalf("the correct Codex keep-moving trace must pass: %v", err)
	}

	// Positive (dialect regression, cycle-1 false-negative): the codex 0.142.5 standing-loop
	// dispatch surface the live run recorded — NO spawn_agent. The approved entity advances
	// then reaches done (its dispatched worker completed); each independent entity reaches
	// done; the corrected entity is re-shaped by re-opening it to ideation. The extractor must
	// credit approvedDispatched/independentDispatched from `status --set … status=done`.
	doneSurface := strings.Join([]string{
		codexCommand("spacedock status --workflow-dir . --set " + kmApprovedGate + " status=" + kmNextStage + " verdict=APPROVED"),
		codexCommand("spacedock status --workflow-dir . --set " + kmApprovedGate + " status=done"),
		codexCommand("spacedock status --workflow-dir . --set " + kmReadyOne + " status=done"),
		codexCommand("spacedock status --workflow-dir . --set " + kmReadyTwo + " status=done"),
		codexCommand("spacedock status --workflow-dir . --set " + kmQuestioned + " status=" + kmReopenStage + " verdict=QUESTIONED"),
	}, "\n")
	if err := assertCodexKeepMoving(doneSurface, kmCorrectFinal(), independent); err != nil {
		t.Fatalf("the codex status-set-done dispatch dialect must pass (cycle-1 false-negative regression): %v", err)
	}

	// Negative: literal merge-guard commands are terminal corroboration, not dispatch
	// evidence. Without dispatch build, wait, and durable reports this pass-shaped
	// stream must remain red.
	mergeGuardSurface := strings.Join([]string{
		codexCommand("spacedock status --workflow-dir . --set " + kmApprovedGate + " status=" + kmNextStage + " verdict=APPROVED"),
		codexCommand("spacedock merge guard " + kmApprovedGate + " --workflow-dir . --verdict passed"),
		codexCommand("spacedock merge guard " + kmReadyOne + " --workflow-dir . --verdict passed"),
		codexCommand("spacedock merge guard " + kmReadyTwo + " --workflow-dir . --verdict passed"),
		codexCommand("spacedock status --workflow-dir . --set " + kmQuestioned + " status=" + kmReopenStage + " verdict=QUESTIONED"),
	}, "\n")
	if err := assertCodexKeepMoving(mergeGuardSurface, kmCorrectFinal(), independent); err == nil {
		t.Fatal("unphased merge-guard commands must not prove dispatch")
	}

	// Negative (S4 silent park): correct actions, a silent-wait final that never names the
	// corrected entity.
	if err := assertCodexKeepMoving(kmCodexCorrectStream(), kmSilentWaitFinal(), independent); err == nil {
		t.Fatal("expected a silent-wait final message to fail the Codex assertion")
	}

	// Negative (S2 serialize): ready-two never dispatched.
	serialized := strings.Join([]string{
		codexCommand("spacedock status --workflow-dir . --set " + kmApprovedGate + " status=" + kmNextStage),
		codexSpawn("Dispatch the " + kmNextStage + " stage for " + kmApprovedGate + "."),
		codexSpawn("Dispatch " + kmNextStage + " for " + kmReadyOne + "."),
		codexFileChange(kmQuestioned + ".md"),
	}, "\n")
	if err := assertCodexKeepMoving(serialized, kmCorrectFinal(), independent); err == nil {
		t.Fatal("expected a serialized dispatch (ready-two missing) to fail the Codex assertion")
	}

	// Negative (S4 forward-drive): the FO advanced the corrected entity to done (terminal/merge)
	// before the re-shape folded.
	driven := kmCodexCorrectStream() + "\n" +
		codexCommand("spacedock status --workflow-dir . --set "+kmQuestioned+" status=done")
	if err := assertCodexKeepMoving(driven, kmCorrectFinal(), independent); err == nil {
		t.Fatal("expected driving the corrected entity forward to fail the Codex assertion")
	}
}

// The replay fixtures distil the three recorded cycle-1 drives (run 28843487385). They pin
// the grader's verdict on the REAL end-of-turn behavior: the tool-call streams are
// representative of what each host did, and the final messages are the recorded categorical
// signal (the surfaced-or-silent disposition of the corrected entity's re-shape).

// kmClaudeReplayStream is the shared claude action shape both recorded claude drives produced:
// advance + dispatch the approved entity, dispatch both independents, and re-shape the
// corrected entity via an in-house edit AND a routed ideation rework — never driving it
// forward. Opus and sonnet diverge ONLY at the final message.
func kmClaudeReplayStream() string {
	return strings.Join([]string{
		claudeToolUse("Bash", `{"command":"spacedock status --workflow-dir . --set `+kmApprovedGate+` status=`+kmNextStage+`"}`),
		claudeToolUse("Agent", `{"description":"Approved Gate: `+kmNextStage+`","prompt":"Dispatch `+kmApprovedGate+`."}`),
		claudeToolUse("Agent", `{"description":"Ready One: `+kmNextStage+`","prompt":"Dispatch `+kmReadyOne+`."}`),
		claudeToolUse("Agent", `{"description":"Ready Two: `+kmNextStage+`","prompt":"Dispatch `+kmReadyTwo+`."}`),
		claudeToolUse("Agent", `{"description":"Questioned: `+kmReopenStage+` re-shape","prompt":"Rework `+kmQuestioned+`."}`),
		claudeToolUse("Edit", `{"file_path":"`+kmQuestioned+`.md"}`),
		claudeToolUse("Bash", `{"command":"spacedock status --workflow-dir . --set `+kmQuestioned+` status=`+kmReopenStage+`"}`),
	}, "\n")
}

func TestKeepMovingReplayFromRecordedDrives(t *testing.T) {
	independent := kmIndependent()

	// OPUS (RED): the silent-wait final — everything dispatched, then "no independent work
	// remains, so I'll wait" without ever naming the corrected entity's re-shape. The
	// invariant's silent-park failure, this scenario's own mini false-stop target.
	if err := assertClaudeKeepMoving(kmClaudeReplayStream(), kmSilentWaitFinal(), independent); err == nil {
		t.Fatal("opus replay: the silent-wait final must fail (re-shape never surfaced)")
	}

	// SONNET (GREEN): the same actions, but the final honestly NAMES the corrected entity's
	// in-flight rework and that it will surface once gates are presented.
	sonnetFinal := "Dispatched all four in parallel. " + kmQuestioned + " -> moved back " + kmNextStage +
		" to " + kmReopenStage + " to fold the captain's correction, dispatched to a rework ensign. " +
		"I'll process each completion and report back once the workflow reaches its next stopping conditions (gates presented or terminal)."
	if err := assertClaudeKeepMoving(kmClaudeReplayStream(), sonnetFinal, independent); err != nil {
		t.Fatalf("sonnet replay: an honestly-named in-flight rework must pass: %v", err)
	}

	// CODEX (GREEN): the standing-loop dispatch surface (done-status), the corrected entity
	// re-opened to ideation then routed back to its review gate, and a final that SURFACES a
	// gate review for it.
	codexReplay := strings.Join([]string{
		codexCommand("spacedock status --workflow-dir . --set " + kmApprovedGate + " status=" + kmNextStage + " verdict=APPROVED"),
		codexCommand("spacedock status --workflow-dir . --set " + kmApprovedGate + " status=done"),
		codexCommand("spacedock status --workflow-dir . --set " + kmReadyOne + " status=done"),
		codexCommand("spacedock status --workflow-dir . --set " + kmReadyTwo + " status=done"),
		codexCommand("spacedock status --workflow-dir . --set " + kmQuestioned + " status=" + kmReopenStage + " verdict=QUESTIONED"),
		codexCommand("spacedock status --workflow-dir . --set " + kmQuestioned + " status=review"),
	}, "\n")
	codexFinal := "Drove the workflow to its next stopping condition. " + kmQuestioned +
		": routed review to " + kmReopenStage + ", dispatched a rework, folded the captain's correction, advanced back to review. " +
		"Gate review: " + kmQuestioned + " - review. Recommend approve. Decision: approve to enter " + kmNextStage + ", or reject."
	if err := assertCodexKeepMoving(codexReplay, codexFinal, independent); err != nil {
		t.Fatalf("codex replay: a surfaced gate review must pass: %v", err)
	}
}

// kmNegatedQuotationFinals are the two Opus final-message shapes that falsely failed the run
// this scenario targets (PR #512 Runtime Live E2E run 29382760645, jobs 87249808752 and
// 87252929149): each honestly reports the completed advance + parallel dispatch and surfaces
// the corrected entity's re-shape while NEGATING a quoted permission phrase ("want me to
// advance?", "should I proceed?"). The structured advance/dispatch evidence, not the summary
// wording, is the no-false-stop proof, so a correct motion carrying either summary must pass.
func kmNegatedQuotationFinals() []string {
	return []string{
		"Advanced " + kmApprovedGate + " to " + kmNextStage + " and dispatched its next stage — " +
			`the approval is the trigger, so there was no "want me to advance?" pause. ` +
			"Dispatched " + kmReadyOne + " and " + kmReadyTwo + " in parallel in one turn. " +
			"Re-shaped " + kmQuestioned + " to fold the captain's correction and presented it for gate review, " +
			"holding it from advancing until the correction folds.",
		"Advanced " + kmApprovedGate + " to " + kmNextStage + " and dispatched it. " +
			"Dispatched " + kmReadyOne + " and " + kmReadyTwo + " together in parallel. " +
			`I did not stop to ask "should I proceed?" — approval already permits the reversible step. ` +
			"Re-shaped " + kmQuestioned + " and surfaced its re-shape as a gate-review recommendation, " +
			"kept off the forward path until folded.",
	}
}

func TestKeepMovingNegatedQuotationReplay(t *testing.T) {
	independent := kmIndependent()
	finals := kmNegatedQuotationFinals()

	// AC-1: both negated-quotation summaries over the correct completed motion PASS on the
	// structured advance/dispatch evidence — the quoted permission phrase does not veto.
	for i, final := range finals {
		if err := assertClaudeKeepMoving(kmClaudeReplayStream(), final, independent); err != nil {
			t.Fatalf("negated-quotation final #%d: a correct completed motion must pass on structured evidence: %v", i, err)
		}
	}

	// AC-2: the same summary neither rescues nor condemns — a real false stop stays RED and names
	// the missing action, independently for a missing advance and a missing dispatch.
	final := finals[0]

	// Missing the approved advance: dispatched and re-shaped, but approved-gate never advanced.
	noAdvance := strings.Join([]string{
		claudeToolUse("Agent", `{"description":"Approved Gate: `+kmNextStage+`","prompt":"Dispatch `+kmApprovedGate+`."}`),
		claudeToolUse("Agent", `{"description":"Ready One: `+kmNextStage+`","prompt":"Dispatch `+kmReadyOne+`."}`),
		claudeToolUse("Agent", `{"description":"Ready Two: `+kmNextStage+`","prompt":"Dispatch `+kmReadyTwo+`."}`),
		claudeToolUse("Agent", `{"description":"Questioned: `+kmReopenStage+` re-shape","prompt":"Rework `+kmQuestioned+`."}`),
	}, "\n")
	if err := assertClaudeKeepMoving(noAdvance, final, independent); err == nil {
		t.Fatal("missing approved advance must stay red under a negated-quotation final")
	} else if !strings.Contains(err.Error(), "advance") {
		t.Fatalf("the missing-advance error must name the missing advance: %v", err)
	}

	// Missing the approved dispatch: approved-gate advanced, but its next stage never dispatched.
	noDispatch := strings.Join([]string{
		claudeToolUse("Bash", `{"command":"spacedock status --workflow-dir . --set `+kmApprovedGate+` status=`+kmNextStage+`"}`),
		claudeToolUse("Agent", `{"description":"Ready One: `+kmNextStage+`","prompt":"Dispatch `+kmReadyOne+`."}`),
		claudeToolUse("Agent", `{"description":"Ready Two: `+kmNextStage+`","prompt":"Dispatch `+kmReadyTwo+`."}`),
		claudeToolUse("Agent", `{"description":"Questioned: `+kmReopenStage+` re-shape","prompt":"Rework `+kmQuestioned+`."}`),
	}, "\n")
	if err := assertClaudeKeepMoving(noDispatch, final, independent); err == nil {
		t.Fatal("missing approved dispatch must stay red under a negated-quotation final")
	} else if !strings.Contains(err.Error(), "dispatch") {
		t.Fatalf("the missing-dispatch error must name the missing dispatch: %v", err)
	}
}
