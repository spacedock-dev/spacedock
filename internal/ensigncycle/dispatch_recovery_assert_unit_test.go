// ABOUTME: Offline RED-control unit tests for the fo-dispatch-recovery scenario
// ABOUTME: oracles — synthetic stream-json, no model spend, no live credential.
package ensigncycle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDispatchRecoveryPromptsSelectIntendedDispatchModes(t *testing.T) {
	if !strings.Contains(bareReachablePrompt(), "`widget-task`") {
		t.Fatal("bare prompt must name its single entity to select blocking bare dispatch")
	}
	if strings.Contains(breakGlassShimPrompt(), "`widget-task`") {
		t.Fatal("break-glass prompt must not name one entity and accidentally select blocking bare dispatch")
	}
}

// boundedRetryGoodStream is a hand-authored representative stream shaped like the
// real multi-delta runner output: a text delta acknowledging an Agent() error, the
// failed initial Agent() dispatch, then exactly one bounded re-dispatch carrying the
// `-retry` suffix on the same `{worker_key}-{slug}-{stage}` stem, and NO third
// attempt.
const boundedRetryGoodStream = `{"type":"assistant","message":{"id":"m1","content":[{"type":"text","text":"Agent() returned an error dispatching widget-task implementation; recording retry 1 and re-dispatching once."}]}}
{"type":"assistant","message":{"id":"m1","content":[{"type":"tool_use","id":"t1","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation","run_in_background":true,"description":"initial dispatch"}}]}}
{"type":"assistant","message":{"id":"m2","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation-retry","run_in_background":true,"description":"bounded re-dispatch"}}]}}`

func TestAssertBoundedRetryObservablesOffline(t *testing.T) {
	if err := assertBoundedRetryObservables(boundedRetryGoodStream); err != nil {
		t.Fatalf("the positive fixture (one dispatch, one -retry re-dispatch, no third attempt) must pass: %v", err)
	}
}

// TestAssertBoundedRetryObservablesCatchesNoRetry is the RED control for the bound's
// lower edge: a dispatch that failed but was NEVER retried (no `-retry` re-dispatch)
// must fail — a dispatch failure must be retried once before anything holds.
func TestAssertBoundedRetryObservablesCatchesNoRetry(t *testing.T) {
	stream := `{"type":"assistant","message":{"id":"m1","content":[{"type":"tool_use","id":"t1","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation","run_in_background":true,"description":"initial dispatch"}}]}}`
	if err := assertBoundedRetryObservables(stream); err == nil {
		t.Fatal("a stream with a failed dispatch but NO -retry re-dispatch must fail — a dispatch failure must be retried once")
	}
}

// TestAssertBoundedRetryObservablesCatchesThirdAttempt is the RED control for the
// bound's upper edge: a THIRD Agent() dispatch for one `(entity, stage)` (past the
// initial + one `-retry`) must fail — the retry is bounded to a single re-attempt.
func TestAssertBoundedRetryObservablesCatchesThirdAttempt(t *testing.T) {
	stream := `{"type":"assistant","message":{"id":"m1","content":[{"type":"tool_use","id":"t1","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation","run_in_background":true}}]}}
{"type":"assistant","message":{"id":"m2","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation-retry","run_in_background":true}}]}}
{"type":"assistant","message":{"id":"m3","content":[{"type":"tool_use","id":"t3","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation-retry-again","run_in_background":true}}]}}`
	if err := assertBoundedRetryObservables(stream); err == nil {
		t.Fatal("a third Agent() dispatch for one (entity, stage) must fail — the retry is bounded to one re-attempt")
	}
}

// bareReachableGoodStream is a hand-authored representative post-retirement bare
// drive: a bare-shaped Agent() call (no name and explicit `run_in_background:false`), with NO
// retired Degraded Mode captain report and NO recovery-skill load.
const bareReachableGoodStream = `{"type":"assistant","message":{"id":"m1","content":[{"type":"text","text":"The captain asked for bare dispatch; dispatching one worker at a time, blocking on each."}]}}
{"type":"assistant","message":{"id":"m2","content":[{"type":"tool_use","id":"t1","name":"Agent","input":{"subagent_type":"spacedock:ensign","description":"bare dispatch","prompt":"...","run_in_background":false}}]}}`

func TestAssertBareReachableObservablesOffline(t *testing.T) {
	if err := assertBareReachableObservables(bareReachableGoodStream); err != nil {
		t.Fatalf("the positive fixture (bare Agent, no retired report, no recovery-skill load) must pass: %v", err)
	}
}

// TestAssertBareReachableObservablesCatchesRetiredReport is the wrong-way RED
// control (ii): a bare drive that STILL emits the retired Degraded Mode captain
// report must fail — the report was retired with Degraded Mode.
func TestAssertBareReachableObservablesCatchesRetiredReport(t *testing.T) {
	stream := `{"type":"assistant","message":{"id":"m1","content":[{"type":"text","text":"Falling back to bare mode for the remainder of this session due to infrastructure failure."}]}}
{"type":"assistant","message":{"id":"m2","content":[{"type":"tool_use","id":"t1","name":"Agent","input":{"subagent_type":"spacedock:ensign","description":"bare dispatch"}}]}}`
	if err := assertBareReachableObservables(stream); err == nil {
		t.Fatal("a stream still emitting the retired Degraded Mode captain report must fail — the report was retired")
	}
}

// TestAssertBareReachableObservablesCatchesRecoverySkillLoad is the wrong-way RED
// control (iii): a post-retirement bare drive that loads
// Skill(skill="spacedock:fo-dispatch-recovery") must fail.
func TestAssertBareReachableObservablesCatchesRecoverySkillLoad(t *testing.T) {
	stream := `{"type":"assistant","message":{"id":"m1","content":[{"type":"tool_use","id":"t1","name":"Skill","input":{"skill":"spacedock:fo-dispatch-recovery"}}]}}
{"type":"assistant","message":{"id":"m2","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","description":"bare dispatch"}}]}}`
	if err := assertBareReachableObservables(stream); err == nil {
		t.Fatal("a post-retirement bare drive that loads spacedock:fo-dispatch-recovery must fail")
	}
}

// TestAssertBareReachableObservablesCatchesNoBareAgent is the RED control for
// observable (i): a stream whose only Agent() call carries `name`/`run_in_background`
// (a team-shaped call, not bare) must fail — bare dispatch was not reached.
func TestAssertBareReachableObservablesCatchesNoBareAgent(t *testing.T) {
	stream := `{"type":"assistant","message":{"id":"m1","content":[{"type":"tool_use","id":"t1","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation","run_in_background":true,"description":"team-shaped"}}]}}`
	if err := assertBareReachableObservables(stream); err == nil {
		t.Fatal("a stream with only a name/run_in_background Agent() (no bare-shaped call) must fail — bare dispatch was not reached")
	}
}

// breakGlassGoodStream is a hand-authored representative stream: a text block
// naming the failed `dispatch build` helper BEFORE any Agent() call, a
// Skill(spacedock:fo-dispatch-recovery) load, then a break-glass-shaped Agent()
// call (run_in_background=true, a {worker_key}-{slug}-{stage} name, and a prompt
// carrying the ensign skill invocation plus an inline stage definition).
const breakGlassBareGoodStream = `{"type":"assistant","message":{"id":"msg1","content":[{"type":"text","text":"spacedock dispatch build exited non-zero (exit 1); reporting the helper failure before proceeding."}]}}
{"type":"assistant","message":{"id":"msg2","content":[{"type":"tool_use","id":"t1","name":"Skill","input":{"skill":"spacedock:fo-dispatch-recovery"}}]}}
{"type":"assistant","message":{"id":"msg3","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","description":"Widget Task: implementation","prompt":"## First action\n\nSkill(skill=\"spacedock:ensign\")\n\n### Stage definition:\n\nAppend a one-line marker to ` + "`widget-task.md`" + ` proving the worker ran.\n\n- **Outputs:** An implementation stage report.\n\n### Stage report"}}]}}`

const breakGlassTeamGoodStream = `{"type":"assistant","message":{"id":"msg1","content":[{"type":"text","text":"spacedock dispatch build exited non-zero (exit 1); reporting the helper failure before proceeding."}]}}
{"type":"assistant","message":{"id":"msg2","content":[{"type":"tool_use","id":"t1","name":"Skill","input":{"skill":"spacedock:fo-dispatch-recovery"}}]}}
{"type":"assistant","message":{"id":"msg3","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation","run_in_background":true,"description":"Widget Task: implementation","prompt":"## First action\n\nSkill(skill=\"spacedock:ensign\")\n\n### Stage definition:\n\nAppend a one-line marker to ` + "`widget-task.md`" + ` proving the worker ran.\n\n- **Outputs:** An implementation stage report.\n\n### Stage report\n\nAppend the report.\n\n### Completion Signal\n\nSendMessage(to=\"team-lead\", message=\"Done\")"}}]}}`

func TestAssertBreakGlassObservablesOffline(t *testing.T) {
	explicitFalseBareStream := strings.Replace(
		breakGlassBareGoodStream,
		`"description":"Widget Task: implementation"`,
		`"description":"Widget Task: implementation","run_in_background":false`,
		1,
	)
	for _, tc := range []struct {
		name   string
		mode   dispatchMode
		stream string
	}{
		{"selected bare omits run_in_background", dispatchModeBare, breakGlassBareGoodStream},
		{"selected bare serializes run_in_background false", dispatchModeBare, explicitFalseBareStream},
		{"selected team uses team call", dispatchModeTeam, breakGlassTeamGoodStream},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := assertBreakGlassObservables(tc.stream, tc.mode); err != nil {
				t.Fatalf("mode-preserving fixture must pass: %v", err)
			}
		})
	}
	if err := assertBreakGlassObservables(breakGlassBareGoodStream, dispatchModeTeam); err == nil {
		t.Fatal("selected team with a bare call must fail")
	}
	if err := assertBreakGlassObservables(breakGlassTeamGoodStream, dispatchModeBare); err == nil {
		t.Fatal("selected bare with a team call must fail")
	}
	if err := assertBreakGlassObservables(breakGlassBareGoodStream+"\n"+breakGlassBareGoodStream, dispatchModeBare); err == nil {
		t.Fatal("multiple workers must fail")
	}
	for _, mutation := range []struct {
		name string
		old  string
		new  string
	}{
		{"background true", `"description":"Widget Task: implementation"`, `"description":"Widget Task: implementation","run_in_background":true`},
		{"name present", `"description":"Widget Task: implementation"`, `"description":"Widget Task: implementation","name":"ensign-widget-implementation"`},
		{"team_name present", `"description":"Widget Task: implementation"`, `"description":"Widget Task: implementation","team_name":"legacy"`},
	} {
		t.Run("selected bare rejects "+mutation.name, func(t *testing.T) {
			stream := strings.Replace(breakGlassBareGoodStream, mutation.old, mutation.new, 1)
			if err := assertBreakGlassObservables(stream, dispatchModeBare); err == nil {
				t.Fatalf("selected bare must reject %s", mutation.name)
			}
		})
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
	if err := assertBreakGlassObservables(stream, dispatchModeTeam); err == nil {
		t.Fatal("a helper-failure report observed only AFTER the Agent() call must fail — the first action is reporting, not dispatching")
	}
}

// TestAssertBreakGlassObservablesCatchesMissingSkillLoad is the RED control for
// observable (ii).
func TestAssertBreakGlassObservablesCatchesMissingSkillLoad(t *testing.T) {
	stream := `{"type":"assistant","message":{"id":"msg1","content":[{"type":"text","text":"spacedock dispatch build exited non-zero; reporting now."}]}}
{"type":"assistant","message":{"id":"msg2","content":[{"type":"tool_use","id":"t2","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"ensign-widget-implementation","run_in_background":true,"prompt":"Skill(skill=\"spacedock:ensign\")\n### Stage definition:\nbody"}}]}}`
	if err := assertBreakGlassObservables(stream, dispatchModeTeam); err == nil {
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
	if err := assertBreakGlassObservables(stream, dispatchModeTeam); err == nil {
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
	if err := assertBreakGlassObservables(stream, dispatchModeTeam); err == nil {
		t.Fatal("an Agent() prompt missing the inline ### Stage definition must fail the break-glass shape check")
	}
}

func TestAssertBreakGlassDurableResultMutations(t *testing.T) {
	complete := dispatchRecoveryEntity() + "\n" + dispatchRecoveryMarker + "\n\n## Stage Report: implementation\n\n- DONE: worker ran\n  Marker and report committed.\n\n### Summary\n\nRecovery completed.\n"
	for _, tc := range []struct {
		name   string
		mutate func(string) string
		dirty  bool
		wantOK bool
	}{
		{"complete", func(s string) string { return s }, false, true},
		{"missing marker", func(s string) string { return strings.ReplaceAll(s, dispatchRecoveryMarker, "marker-removed") }, false, false},
		{"missing heading", func(s string) string { return strings.Replace(s, "## Stage Report: implementation", "## Notes", 1) }, false, false},
		{"missing done", func(s string) string { return strings.Replace(s, "- DONE:", "- SKIPPED:", 1) }, false, false},
		{"missing summary", func(s string) string { return strings.Replace(s, "### Summary", "### Notes", 1) }, false, false},
		{"uncommitted result", func(s string) string { return dispatchRecoveryEntity() }, true, false},
		{"dirty entity", func(s string) string { return s }, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			entityPath := filepath.Join(root, "widget-task.md")
			writeFile(t, filepath.Join(root, "README.md"), dispatchRecoveryReadme())
			writeFile(t, entityPath, dispatchRecoveryEntity())
			gitInit(t, root)
			if tc.name != "uncommitted result" {
				writeFile(t, entityPath, tc.mutate(complete))
				git(t, root, "add", "widget-task.md")
				git(t, root, "commit", "-m", "worker result", "--", "widget-task.md")
			} else {
				writeFile(t, entityPath, complete)
			}
			if tc.dirty && tc.name == "dirty entity" {
				if err := os.WriteFile(entityPath, []byte(readFile(t, entityPath)+"dirty\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			err := assertBreakGlassDurableResult(root, entityPath)
			if (err == nil) != tc.wantOK {
				t.Fatalf("assertBreakGlassDurableResult error = %v, wantOK %v", err, tc.wantOK)
			}
		})
	}
}

func TestAssertCompleteRecoveryReportRejectsScatteredTokens(t *testing.T) {
	content := dispatchRecoveryMarker + "\n\n- DONE: unrelated prior work\n\n### Summary\n\nPrior summary.\n\n## Stage Report: implementation\n\nNo completion entries here.\n"
	if err := assertCompleteRecoveryReport(content); err == nil {
		t.Fatal("DONE and Summary outside the exact implementation Stage Report must not satisfy durable recovery grading")
	}
}
