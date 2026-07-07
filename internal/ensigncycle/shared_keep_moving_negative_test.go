package ensigncycle

import (
	"strings"
	"testing"
)

// The mandatory offline negative for the keep-moving-posture scenario. It proves the
// grader catches each of the four 0223 false-stop patterns from the REAL fixture tokens:
// the correct keep-moving trace (advance+dispatch the approved entity, dispatch both
// independent ready entities, re-shape the questioned entity and pause its dispatch, no
// permission question, no async wait) PASSES, and every INCIDENT reds in isolation —
// a post-approval permission question, a serialized independent dispatch, a dropped
// re-shape, driving the corrected entity, a turn-end-on-async wait, and a correction that
// halted the independent work. A tautological assertion that only checked "an advance
// happened" would stay green on the permission-question and async-wait incidents; those
// cases fail it loudly. Offline (default tag): the grader and extractors are pure
// functions over trace / transcript strings, so they spend no model.

// kmCorrectTrace is the trace of a run that kept moving on every axis: it advanced and
// dispatched the approved entity, dispatched both independent ready entities, re-shaped
// the questioned entity, drove it no further, and ended with neither a permission
// question nor an async wait.
func kmCorrectTrace() keepMovingTrace {
	tr := newKeepMovingTrace()
	tr.approvedAdvanced = true
	tr.approvedDispatched = true
	for _, e := range kmIndependent() {
		tr.independentDispatched[e] = true
	}
	tr.correctedReshaped = true
	return tr
}

func TestGradeKeepMoving(t *testing.T) {
	independent := kmIndependent()

	// Positive: the correct keep-moving trace passes.
	if err := gradeKeepMoving(kmCorrectTrace(), independent); err != nil {
		t.Fatalf("the correct keep-moving trace must pass: %v", err)
	}

	// Negative (S1 pause): the FO asked permission to advance/dispatch after the approval.
	asked := kmCorrectTrace()
	asked.askedPermission = true
	if err := gradeKeepMoving(asked, independent); err == nil {
		t.Fatal("expected a post-approval permission question to fail")
	}

	// Negative (S1): the approved entity was never advanced past its gate.
	noAdvance := kmCorrectTrace()
	noAdvance.approvedAdvanced = false
	if err := gradeKeepMoving(noAdvance, independent); err == nil {
		t.Fatal("expected a missing advance of the approved entity to fail")
	}

	// Negative (S1 isolating): advanced but never dispatched the next stage. Fails on the
	// dispatch half even though the advance is present, so the advance cannot mask it.
	noDispatch := kmCorrectTrace()
	noDispatch.approvedDispatched = false
	if err := gradeKeepMoving(noDispatch, independent); err == nil {
		t.Fatal("expected an advance-without-dispatch to fail")
	}

	// Negative (S2 serialize): an independent ready entity the FO never dispatched — the
	// serialized-behind-a-pause pattern.
	serialized := kmCorrectTrace()
	serialized.independentDispatched[kmReadyTwo] = false
	if err := gradeKeepMoving(serialized, independent); err == nil {
		t.Fatal("expected an undispatched independent entity (serialized) to fail")
	}

	// Negative (S4): the questioned entity was never re-shaped to fold the correction.
	noReshape := kmCorrectTrace()
	noReshape.correctedReshaped = false
	if err := gradeKeepMoving(noReshape, independent); err == nil {
		t.Fatal("expected a dropped re-shape of the questioned entity to fail")
	}

	// Negative (S4 isolating): the FO drove the questioned entity forward while its
	// mechanism was under correction — its dispatch should have paused.
	drove := kmCorrectTrace()
	drove.correctedDriven = true
	if err := gradeKeepMoving(drove, independent); err == nil {
		t.Fatal("expected driving the corrected entity to fail")
	}

	// Negative (S3): the FO ended its turn on an async wait while independent work remained.
	waited := kmCorrectTrace()
	waited.endedWaitingOnAsync = true
	if err := gradeKeepMoving(waited, independent); err == nil {
		t.Fatal("expected a turn-end-on-async wait to fail")
	}

	// Negative (S4 halt — pattern 4): the correction halted the independent work. The
	// questioned entity was re-shaped, but neither independent entity kept moving — the
	// captain's question was conflated with "stop the session." Reds on the parallel check.
	halted := kmCorrectTrace()
	for _, e := range independent {
		halted.independentDispatched[e] = false
	}
	if err := gradeKeepMoving(halted, independent); err == nil {
		t.Fatal("expected a correction that halted the independent work to fail")
	}
}

// kmCorrectFinal is a keep-moving final message: it reports the completed advance,
// parallel dispatch, and re-shape with neither a permission question nor an async wait.
func kmCorrectFinal() string {
	return "Advanced " + kmApprovedGate + " to " + kmNextStage + " and dispatched it. " +
		"Dispatched " + kmReadyOne + " and " + kmReadyTwo + " in parallel in one turn. " +
		"Re-shaped " + kmQuestioned + " to fold the captain's symlink correction and paused its dispatch. " +
		"All independent work is moving; nothing is blocked."
}

func kmPermissionFinal() string {
	return "I reviewed the captain's approval of " + kmApprovedGate + ". Want me to advance it to " + kmNextStage + " and dispatch the stage?"
}

func kmWaitFinal() string {
	return "Dispatched the " + kmNextStage + " ensign for " + kmApprovedGate + ". I'll wait for it to finish before doing more."
}

// kmClaudeCorrectStream builds the tool-call stream of a correct Claude run: a `status
// --set` advance of the approved entity, an Agent dispatch of its next stage, an Agent
// dispatch of each independent entity, and an Edit re-shaping the questioned entity.
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

func TestAssertClaudeKeepMoving(t *testing.T) {
	independent := kmIndependent()

	// Positive: the correct stream + a keep-moving final message passes.
	if err := assertClaudeKeepMoving(kmClaudeCorrectStream(), kmCorrectFinal(), independent); err != nil {
		t.Fatalf("the correct Claude keep-moving trace must pass: %v", err)
	}

	// Negative (S1 pause): the correct actions but a permission-question final message.
	// The durable state would look identical; the final-message posture red-flags the stop.
	if err := assertClaudeKeepMoving(kmClaudeCorrectStream(), kmPermissionFinal(), independent); err == nil {
		t.Fatal("expected a permission-question final message to fail the Claude assertion")
	}

	// Negative (S3): the correct actions but an async-wait final message.
	if err := assertClaudeKeepMoving(kmClaudeCorrectStream(), kmWaitFinal(), independent); err == nil {
		t.Fatal("expected an async-wait final message to fail the Claude assertion")
	}

	// Negative (S3 resume-wait): a subtler wait — suspending further work until the async
	// returns. Proves the wait regex catches the resume posture, not only "I'll wait".
	resumeWait := "Dispatched the ensigns. I will proceed once they return."
	if err := assertClaudeKeepMoving(kmClaudeCorrectStream(), resumeWait, independent); err == nil {
		t.Fatal("expected a resume-once-they-return final message to fail the Claude assertion")
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

	// Negative (S4): the questioned entity was never re-shaped.
	noReshape := strings.Join([]string{
		claudeToolUse("Bash", `{"command":"spacedock status --workflow-dir . --set `+kmApprovedGate+` status=`+kmNextStage+`"}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch the `+kmNextStage+` stage for `+kmApprovedGate+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyOne+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyTwo+`."}`),
	}, "\n")
	if err := assertClaudeKeepMoving(noReshape, kmCorrectFinal(), independent); err == nil {
		t.Fatal("expected a stream with no re-shape to fail the Claude assertion")
	}

	// Negative (S4 isolating): the FO dispatched the questioned entity despite the correction.
	drove := kmClaudeCorrectStream() + "\n" +
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmQuestioned+`."}`)
	if err := assertClaudeKeepMoving(drove, kmCorrectFinal(), independent); err == nil {
		t.Fatal("expected driving the corrected entity to fail the Claude assertion")
	}
}

// kmCodexCorrectStream builds the transcript of a correct Codex run: a `status --set`
// advance, a spawn_agent dispatch of the approved next stage and of each independent
// entity, and a file_change re-shaping the questioned entity (the codex 0.142.5 surface).
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

	// Positive: the correct Codex transcript + a keep-moving final message passes.
	if err := assertCodexKeepMoving(kmCodexCorrectStream(), kmCorrectFinal(), independent); err != nil {
		t.Fatalf("the correct Codex keep-moving trace must pass: %v", err)
	}

	// Positive (dialect regression): the same run with the re-shape as an apply_patch
	// command instead of a file_change item. The extractor must read both edit surfaces.
	applyPatchReshape := strings.Join([]string{
		codexCommand("spacedock status --workflow-dir . --set " + kmApprovedGate + " status=" + kmNextStage),
		codexSpawn("Dispatch the " + kmNextStage + " stage for " + kmApprovedGate + "."),
		codexSpawn("Dispatch " + kmNextStage + " for " + kmReadyOne + "."),
		codexSpawn("Dispatch " + kmNextStage + " for " + kmReadyTwo + "."),
		codexCommand("apply_patch " + kmQuestioned + ".md"),
	}, "\n")
	if err := assertCodexKeepMoving(applyPatchReshape, kmCorrectFinal(), independent); err != nil {
		t.Fatalf("the apply_patch re-shape dialect must pass: %v", err)
	}

	// Negative (S1 pause): correct actions, a permission-question final message.
	if err := assertCodexKeepMoving(kmCodexCorrectStream(), kmPermissionFinal(), independent); err == nil {
		t.Fatal("expected a permission-question final message to fail the Codex assertion")
	}

	// Negative (S3): correct actions, an async-wait final message.
	if err := assertCodexKeepMoving(kmCodexCorrectStream(), kmWaitFinal(), independent); err == nil {
		t.Fatal("expected an async-wait final message to fail the Codex assertion")
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

	// Negative (S1): no advance of the approved entity.
	noAdvance := strings.Join([]string{
		codexSpawn("Dispatch the " + kmNextStage + " stage for " + kmApprovedGate + "."),
		codexSpawn("Dispatch " + kmNextStage + " for " + kmReadyOne + "."),
		codexSpawn("Dispatch " + kmNextStage + " for " + kmReadyTwo + "."),
		codexFileChange(kmQuestioned + ".md"),
	}, "\n")
	if err := assertCodexKeepMoving(noAdvance, kmCorrectFinal(), independent); err == nil {
		t.Fatal("expected a transcript with no advance to fail the Codex assertion")
	}

	// Negative (S4 isolating): the FO advanced the questioned entity despite the correction
	// (the `status --set questioned` surface — driving the entity its dispatch should pause).
	drove := kmCodexCorrectStream() + "\n" +
		codexCommand("spacedock status --workflow-dir . --set " + kmQuestioned + " status=" + kmNextStage)
	if err := assertCodexKeepMoving(drove, kmCorrectFinal(), independent); err == nil {
		t.Fatal("expected driving the corrected entity to fail the Codex assertion")
	}
}
