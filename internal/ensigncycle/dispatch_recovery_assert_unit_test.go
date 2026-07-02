// ABOUTME: Offline RED-control unit tests for the fo-dispatch-recovery scenario
// ABOUTME: oracles — synthetic stream-json, no model spend, no live credential.
package ensigncycle

import "testing"

// degradedBareGoodStream is a hand-authored representative stream (NOT a captured
// live run — the live baseline capture this scenario needs is a separate,
// credentialed step) shaped like the real multi-delta runner output: a thinking
// delta, a Skill(spacedock:fo-dispatch-recovery) tool_use delta, a text delta
// carrying the verbatim captain report, then two bare-mode Agent() calls (neither
// carries `name` nor `run_in_background`).
const degradedBareGoodStream = `{"type":"assistant","message":{"id":"msg1","content":[{"type":"thinking","thinking":"second dispatch failure — tripping Degraded Mode"}]}}
{"type":"assistant","message":{"id":"msg1","content":[{"type":"tool_use","id":"t1","name":"Skill","input":{"skill":"spacedock:fo-dispatch-recovery"}}]}}
{"type":"assistant","message":{"id":"msg2","content":[{"type":"text","text":"Falling back to bare mode for the remainder of this session due to infrastructure failure. Prior background agents are presumed-zombified; I will not route work to them or through the team registry."}]}}
{"type":"assistant","message":{"id":"msg3","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","description":"bare dispatch","prompt":"..."}}]}}
{"type":"assistant","message":{"id":"msg4","content":[{"type":"tool_use","id":"t3","name":"Agent","input":{"subagent_type":"spacedock:ensign","description":"bare dispatch 2","prompt":"..."}}]}}`

func TestAssertDegradedBareObservablesOffline(t *testing.T) {
	if err := assertDegradedBareObservables(degradedBareGoodStream); err != nil {
		t.Fatalf("the positive fixture (recovery skill loaded, verbatim report, bare Agent calls) must pass: %v", err)
	}
}

// TestAssertDegradedBareObservablesCatchesMissingSkillLoad is the RED control for
// observable (i): a stream with the captain report and bare Agent calls but NO
// Skill(skill="spacedock:fo-dispatch-recovery") tool_use must fail.
func TestAssertDegradedBareObservablesCatchesMissingSkillLoad(t *testing.T) {
	stream := `{"type":"assistant","message":{"id":"msg2","content":[{"type":"text","text":"Falling back to bare mode for the remainder of this session due to infrastructure failure."}]}}
{"type":"assistant","message":{"id":"msg3","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","description":"bare dispatch"}}]}}`
	if err := assertDegradedBareObservables(stream); err == nil {
		t.Fatal("a stream with no Skill(skill=\"spacedock:fo-dispatch-recovery\") tool_use must fail — the trigger did not load the recovery skill")
	}
}

// TestAssertDegradedBareObservablesCatchesMissingCaptainReport is the RED control
// for observable (ii): the skill loads and Agent calls are bare, but no text block
// carries the verbatim captain report sentence.
func TestAssertDegradedBareObservablesCatchesMissingCaptainReport(t *testing.T) {
	stream := `{"type":"assistant","message":{"id":"msg1","content":[{"type":"tool_use","id":"t1","name":"Skill","input":{"skill":"spacedock:fo-dispatch-recovery"}}]}}
{"type":"assistant","message":{"id":"msg3","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","description":"bare dispatch"}}]}}`
	if err := assertDegradedBareObservables(stream); err == nil {
		t.Fatal("a stream missing the verbatim captain report sentence must fail")
	}
}

// TestAssertDegradedBareObservablesCatchesNamedAgentAfterTrip is the RED control for
// observable (iii): an Agent() call that STILL carries `name` (a reuse-shaped call,
// the exact bug Degraded Mode's Effects forbid) after the trip must fail — even
// though the skill loaded and the report fired.
func TestAssertDegradedBareObservablesCatchesNamedAgentAfterTrip(t *testing.T) {
	stream := `{"type":"assistant","message":{"id":"msg1","content":[{"type":"tool_use","id":"t1","name":"Skill","input":{"skill":"spacedock:fo-dispatch-recovery"}}]}}
{"type":"assistant","message":{"id":"msg2","content":[{"type":"text","text":"Falling back to bare mode for the remainder of this session due to infrastructure failure."}]}}
{"type":"assistant","message":{"id":"msg3","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"stale-worker-name","run_in_background":true,"description":"still team-shaped"}}]}}`
	if err := assertDegradedBareObservables(stream); err == nil {
		t.Fatal("an Agent() call carrying `name`/`run_in_background` after Degraded Mode tripped must fail — bare mode requires both omitted")
	}
}

// breakGlassGoodStream is a hand-authored representative stream: a text block
// naming the failed `dispatch build` helper BEFORE any Agent() call, a
// Skill(spacedock:fo-dispatch-recovery) load, then a break-glass-shaped Agent()
// call (run_in_background=true, a {worker_key}-{slug}-{stage} name, and a prompt
// carrying the ensign skill invocation plus an inline stage definition).
const breakGlassGoodStream = `{"type":"assistant","message":{"id":"msg1","content":[{"type":"text","text":"spacedock dispatch build exited non-zero (exit 1); reporting the helper failure before proceeding."}]}}
{"type":"assistant","message":{"id":"msg2","content":[{"type":"tool_use","id":"t1","name":"Skill","input":{"skill":"spacedock:fo-dispatch-recovery"}}]}}
{"type":"assistant","message":{"id":"msg3","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation","run_in_background":true,"description":"break-glass dispatch","prompt":"## First action\n\nSkill(skill=\"spacedock:ensign\")\n\n### Stage definition:\n\ncopy the stage body here"}}]}}`

func TestAssertBreakGlassObservablesOffline(t *testing.T) {
	if err := assertBreakGlassObservables(breakGlassGoodStream); err != nil {
		t.Fatalf("the positive fixture (report before Agent, recovery skill loaded, break-glass-shaped Agent call) must pass: %v", err)
	}
}

// TestAssertBreakGlassObservablesCatchesReportAfterAgent is the RED control for
// observable (i): the helper-failure report text appears only AFTER an Agent() call
// already fired — the "never hand-assemble a dispatch while the helper works, report
// first" invariant is violated.
func TestAssertBreakGlassObservablesCatchesReportAfterAgent(t *testing.T) {
	stream := `{"type":"assistant","message":{"id":"msg1","content":[{"type":"tool_use","id":"t1","name":"Skill","input":{"skill":"spacedock:fo-dispatch-recovery"}}]}}
{"type":"assistant","message":{"id":"msg2","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation","run_in_background":true,"prompt":"Skill(skill=\"spacedock:ensign\")\n### Stage definition:\nbody"}}]}}
{"type":"assistant","message":{"id":"msg3","content":[{"type":"text","text":"spacedock dispatch build exited non-zero; reporting now."}]}}`
	if err := assertBreakGlassObservables(stream); err == nil {
		t.Fatal("a helper-failure report observed only AFTER the Agent() call must fail — the first action is reporting, not dispatching")
	}
}

// TestAssertBreakGlassObservablesCatchesMissingSkillLoad is the RED control for
// observable (ii).
func TestAssertBreakGlassObservablesCatchesMissingSkillLoad(t *testing.T) {
	stream := `{"type":"assistant","message":{"id":"msg1","content":[{"type":"text","text":"spacedock dispatch build exited non-zero; reporting now."}]}}
{"type":"assistant","message":{"id":"msg2","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation","run_in_background":true,"prompt":"Skill(skill=\"spacedock:ensign\")\n### Stage definition:\nbody"}}]}}`
	if err := assertBreakGlassObservables(stream); err == nil {
		t.Fatal("a stream with no Skill(skill=\"spacedock:fo-dispatch-recovery\") tool_use must fail")
	}
}

// TestAssertBreakGlassObservablesCatchesWrongAgentShape is the RED control for
// observable (iii): the Agent() call is missing run_in_background=true, so it is
// NOT break-glass-shaped even though the report and skill load are both present.
func TestAssertBreakGlassObservablesCatchesWrongAgentShape(t *testing.T) {
	stream := `{"type":"assistant","message":{"id":"msg1","content":[{"type":"text","text":"spacedock dispatch build exited non-zero; reporting now."}]}}
{"type":"assistant","message":{"id":"msg2","content":[{"type":"tool_use","id":"t1","name":"Skill","input":{"skill":"spacedock:fo-dispatch-recovery"}}]}}
{"type":"assistant","message":{"id":"msg3","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation","prompt":"Skill(skill=\"spacedock:ensign\")\n### Stage definition:\nbody"}}]}}`
	if err := assertBreakGlassObservables(stream); err == nil {
		t.Fatal("an Agent() call missing run_in_background=true must fail the break-glass shape check")
	}
}

// TestAssertBreakGlassObservablesCatchesMissingStageDefInPrompt is the RED control
// for the prompt-shape half of observable (iii): run_in_background=true and a
// shaped name are present, but the prompt lacks the inline ### Stage definition —
// the FO hand-waved the dispatch instead of inlining the stage body verbatim.
func TestAssertBreakGlassObservablesCatchesMissingStageDefInPrompt(t *testing.T) {
	stream := `{"type":"assistant","message":{"id":"msg1","content":[{"type":"text","text":"spacedock dispatch build exited non-zero; reporting now."}]}}
{"type":"assistant","message":{"id":"msg2","content":[{"type":"tool_use","id":"t1","name":"Skill","input":{"skill":"spacedock:fo-dispatch-recovery"}}]}}
{"type":"assistant","message":{"id":"msg3","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation","run_in_background":true,"prompt":"Skill(skill=\"spacedock:ensign\")\ngo do the task"}}]}}`
	if err := assertBreakGlassObservables(stream); err == nil {
		t.Fatal("an Agent() prompt missing the inline ### Stage definition must fail the break-glass shape check")
	}
}
