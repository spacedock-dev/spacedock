package ensigncycle

import (
	"strings"
	"testing"
)

// Offline table tests for the host-specific reviewer-reuse assertions. They prove
// each assertion requires a REAL reuse tool call targeting the validation reviewer
// — not loose narration, not an unrelated tool, not a call to a different
// recipient. Without these the assertions are exercised only by the model-spending
// live runners; here they cost milliseconds and pin the discriminating behavior.

func TestAssertClaudeReviewerReuse(t *testing.T) {
	// The agentId-addressed reuse shape — one shape the FO emits in a Claude team
	// (a teams-mode kept-alive completed agent can be resumed by the agentId its
	// spawn returned). A validation-stage Agent spawn (exactly one) returns an
	// agentId on completion; the FO resumes THAT reviewer by agentId for the cycle-2
	// re-review. Both events must be present and correlated for this to pass.
	validationSpawn := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_V","input":{"description":"Rejection Task: validation","subagent_type":"spacedock:ensign"}}]}}`
	validationResult := `{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_V","content":[{"type":"text","text":"Verdict: REJECTED."},{"type":"text","text":"agentId: a94abe89c85f9f4cc (use SendMessage with to: 'a94abe89c85f9f4cc' to continue this agent)"}]}]}}`
	reuseByAgentID := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"a94abe89c85f9f4cc","message":"re-review: the fix marker is now present"}}]}}`
	agentIDReuse := validationSpawn + "\n" + validationResult + "\n" + reuseByAgentID

	// An agentId from an IMPLEMENTATION spawn, not the validation reviewer: a
	// SendMessage to it is NOT reviewer reuse and must NOT pass. This pins the
	// correlation — the assertion must tie the reused handle back to the validation
	// spawn, not accept any agentId.
	implSpawn := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_I","input":{"description":"Rejection Task: implementation","subagent_type":"spacedock:ensign"}}]}}`
	implResult := `{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_I","content":[{"type":"text","text":"agentId: a000impl0000 (use SendMessage with to: 'a000impl0000' to continue this agent)"}]}]}}`
	reuseImplAgentID := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"a000impl0000","message":"apply the fix"}}]}}`
	wrongAgentIDReuse := implSpawn + "\n" + implResult + "\n" + reuseImplAgentID

	// A SendMessage to a validation-reviewer agentId that was NEVER spawned/returned
	// in this transcript: an opaque handle with no correlating validation spawn must
	// NOT pass.
	uncorrelatedAgentID := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"adeadbeef0000000","message":"re-review"}}]}}`

	// Genuine NAME-addressed reuse: the FO spawns the cycle-1 validation reviewer
	// ONCE, then re-engages it for cycle 2 by its spawn NAME (the production shape —
	// testdata/*.jsonl address members by name, not agentId). Exactly one validation
	// spawn + a name-message to it = genuine reuse → PASS.
	nameSpawnValidation := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_NV","input":{"description":"Rejection Task: validation","subagent_type":"spacedock:ensign"}}]}}`
	reuseByName := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"spacedock-ensign-rejection-task-validation","message":"re-review: the fix marker is now present"}}]}}`
	nameReuse := nameSpawnValidation + "\n" + reuseByName

	// FRESH-DISPATCH false-pass (the #302/cycle-8 hole): the FO drives cycle-1
	// validation (spawn #1), then FRESH-dispatches the cycle-2 validator (a SECOND
	// validation Agent spawn — the FORBIDDEN behavior the #141 keepalive contract
	// prohibits) and kicks it off by its validation NAME. The old assertion greened
	// this on the bare name substring; the strengthened one must RED it — two
	// validation spawns means the FO did NOT reuse the kept-alive cycle-1 reviewer.
	freshCycle2Spawn := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_FV","input":{"description":"Rejection Task: validation (cycle 2)","subagent_type":"spacedock:ensign"}}]}}`
	freshDispatchNameMessage := nameSpawnValidation + "\n" + freshCycle2Spawn + "\n" + reuseByName

	cases := []struct {
		name    string
		stream  string
		wantErr bool
	}{
		{"genuine name-addressed reuse (one validation spawn)", nameReuse, false},
		{"reuse by correlated validation agentId", agentIDReuse, false},
		{"fresh cycle-2 dispatch + name message must RED", freshDispatchNameMessage, true},
		{
			"loose narration only",
			`{"type":"assistant","message":{"content":[{"type":"text","text":"I will reuse the validation reviewer via SendMessage."}]}}`,
			true,
		},
		{
			"unrelated tool targeting validation",
			`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","input":{"to":"spacedock-ensign-validation","description":"fresh validation dispatch"}}]}}`,
			true,
		},
		{
			"SendMessage to a non-validation recipient by name",
			`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"spacedock-ensign-implementation","message":"apply the fix"}}]}}`,
			true,
		},
		{"reuse of an implementation agentId, not the reviewer", wrongAgentIDReuse, true},
		{"SendMessage to an uncorrelated agentId", uncorrelatedAgentID, true},
		{"empty stream", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := assertClaudeReviewerReuse(tc.stream)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %q, got nil", tc.name)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected pass for %q, got: %v", tc.name, err)
			}
		})
	}
}

// TestSingleEntityBareReviewerNonReuseRepro is the DETERMINISTIC, offline,
// model-free repro of the AC-3 `rejection-flow` finding (backlog seed e3z,
// bare-mode-coverage-baseline). It does NOT run the live scenario — it pins the
// END-STATE shape the live single-entity (`-p`) Claude run produces, so the
// causality can be reasoned about without spending a model or touching the
// validator's live run.
//
// Root cause (verified in the contract, not assumed):
//   - The `rejection-flow` scenario launches `spacedock claude -- -p {prompt}` and
//     the prompt says "Process only the entity `rejection-task`". Both single-entity
//     activation conditions hold (non-interactive `claude -p` AND the prompt names a
//     specific entity — first-officer-shared-core.md `## Single-Entity Mode`).
//   - In single-entity mode the contract mandates bare dispatch: "In single-entity
//     mode, skip team creation. Use bare-mode dispatch for all agent spawning"
//     (claude-fo-dispatch.md). This clause predates P2 (it has lived since the
//     original vendoring `83c73494`), so single-entity Claude reviewers are bare
//     both before and after lazy-TeamCreate.
//   - A bare reviewer fails Claude reuse-condition-1, "Not in bare mode (teams
//     available)" (claude-fo-dispatch.md). With no team, `dispatch build --bare-mode`
//     emits an Agent call with `name`/`team_name` ABSENT (build.go: Name/TeamName are
//     *string omitempty; the bare-mode parity case pins it), so there is no kept-alive
//     handle to SendMessage at all. The cycle-2 re-review therefore fresh-dispatches a
//     SECOND validation Agent → exactly the >1-validation-spawn the assertion reds on.
//
// The two bare validation spawns below carry NO `id`/agentId-returning tool_result
// (a bare worker is not a team member, so it returns no `agentId:` resume handle).
// assertClaudeReviewerReuse reds on the spawn COUNT (>1) — the #141 keepalive
// violation — independent of any handle, which is exactly the live failure's shape.
//
// NOTE the contrast with Codex: the Codex `rejection-flow` reuses via `send_input`
// to a persistent thread (codex-first-officer-runtime.md: "no team_name lifecycle …
// use Codex task names and mailbox notifications as the worker handle"), which does
// NOT depend on a Claude team, so single-entity Codex reuse is unaffected. The
// finding is Claude-specific — it is the team-gated reuse-condition-1 meeting
// single-entity bare mode.
func TestSingleEntityBareReviewerNonReuseRepro(t *testing.T) {
	// Two BARE validation Agent spawns (no team_name, no agentId resume handle),
	// the deterministic end-state of a single-entity `-p` rejection-flow run: the FO
	// fresh-dispatches the cycle-2 validator because the bare cycle-1 reviewer fails
	// reuse-condition-1.
	bareCycle1Validation := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_BV1","input":{"description":"Rejection Task: validation","subagent_type":"spacedock:ensign"}}]}}`
	bareCycle2Validation := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_BV2","input":{"description":"Rejection Task: validation (cycle 2 fresh)","subagent_type":"spacedock:ensign"}}]}}`
	bareTwoSpawnStream := bareCycle1Validation + "\n" + bareCycle2Validation

	err := assertClaudeReviewerReuse(bareTwoSpawnStream)
	if err == nil {
		t.Fatal("single-entity bare two-validation-spawn end-state must RED assertClaudeReviewerReuse — this is the AC-3 finding's deterministic shape")
	}
	if !strings.Contains(err.Error(), "FRESH-dispatched the cycle-2 validator") {
		t.Fatalf("the repro must red on the #141 keepalive violation (>1 validation spawn), got a different failure: %v", err)
	}

	// Falsifiability control: the SAME run in TEAM mode (a single kept-alive cycle-1
	// reviewer reused by agentId) PASSES — proving the repro's red is caused by the
	// extra bare spawn, not by the assertion being unsatisfiable. This is the
	// before/after the fix would restore (option (a): give the -p feedback flow a
	// team so the reviewer is reusable).
	teamCycle1Spawn := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_TV","input":{"description":"Rejection Task: validation","subagent_type":"spacedock:ensign"}}]}}`
	teamCycle1Result := `{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_TV","content":[{"type":"text","text":"agentId: a1111deadbeef0 (use SendMessage with to: 'a1111deadbeef0')"}]}]}}`
	teamReuse := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"a1111deadbeef0","message":"re-review: fix marker present"}}]}}`
	teamModeStream := teamCycle1Spawn + "\n" + teamCycle1Result + "\n" + teamReuse
	if err := assertClaudeReviewerReuse(teamModeStream); err != nil {
		t.Fatalf("the team-mode control (one reviewer reused by agentId) must PASS — else the repro's red is not attributable to the extra bare spawn: %v", err)
	}
}

func TestAssertCodexReviewerReuse(t *testing.T) {
	// The real Codex collab_tool_call shape. A spawn_agent dispatching the validation
	// stage binds the reviewer's thread id (vThread); a later send_input to vThread is
	// the cycle-2 reviewer reuse. Both events, correlated, are the only passing shape.
	const vThread = "019e9695-1a3a-7481-be31-a1e81d53ec6d"
	const iThread = "019e9693-ad11-7701-b615-90d6a86b2dc8"
	spawnValidation := `{"type":"item.completed","item":{"type":"collab_tool_call","tool":"spawn_agent","receiver_thread_ids":["` + vThread + `"],"prompt":"Read /tmp/spacedock-dispatch/spacedock-ensign-rejection-task-validation.md and treat its content as your assignment."}}`
	spawnImpl := `{"type":"item.completed","item":{"type":"collab_tool_call","tool":"spawn_agent","receiver_thread_ids":["` + iThread + `"],"prompt":"Read /tmp/spacedock-dispatch/spacedock-ensign-rejection-task-implementation.md and treat its content as your assignment."}}`
	reuseValidation := `{"type":"item.started","item":{"type":"collab_tool_call","tool":"send_input","receiver_thread_ids":["` + vThread + `"],"prompt":"Re-run validation for rejection-task as cycle 2 using your existing validation reviewer context."}}`
	feedbackToImpl := `{"type":"item.started","item":{"type":"collab_tool_call","tool":"send_input","receiver_thread_ids":["` + iThread + `"],"prompt":"Feedback routed from validation to implementation for rejection-task. The fix marker is absent."}}`

	realReuse := spawnValidation + "\n" + spawnImpl + "\n" + feedbackToImpl + "\n" + reuseValidation

	// FRESH-DISPATCH false-pass (the cycle-8 M2 hole): the FO drives cycle-1
	// validation (spawn #1 → vThread), then FRESH-spawns a SECOND cycle-2 validation
	// reviewer (vThread2 — its prompt also names validation) and send_inputs to THAT
	// new thread. The old assertion bound ANY validation-prompt spawn as "the
	// validation thread", so a send_input to the fresh cycle-2 thread passed as
	// "reuse". The strengthened assertion must RED it: two validation spawn_agents
	// means the FO fresh-dispatched, it did not reuse the kept-alive cycle-1 reviewer.
	const vThread2 = "019e9696-2b4b-8592-cf42-b2f92e64fd7e"
	freshCycle2Spawn := `{"type":"item.completed","item":{"type":"collab_tool_call","tool":"spawn_agent","receiver_thread_ids":["` + vThread2 + `"],"prompt":"Read /tmp/spacedock-dispatch/spacedock-ensign-rejection-task-validation.md and treat its content as your assignment."}}`
	sendInputToFresh := `{"type":"item.started","item":{"type":"collab_tool_call","tool":"send_input","receiver_thread_ids":["` + vThread2 + `"],"prompt":"Re-run validation for rejection-task as cycle 2."}}`
	freshDispatch := spawnValidation + "\n" + spawnImpl + "\n" + feedbackToImpl + "\n" + freshCycle2Spawn + "\n" + sendInputToFresh

	cases := []struct {
		name    string
		jsonl   string
		wantErr bool
	}{
		{"real send_input to the validation reviewer thread", realReuse, false},
		{"fresh cycle-2 spawn_agent + send_input must RED", freshDispatch, true},
		{
			"loose narration only",
			`{"type":"item.completed","item":{"type":"agent_message","text":"I will send_input to the validation worker."}}`,
			true,
		},
		{
			"unrelated tool referencing validation",
			spawnValidation + "\n" + `{"type":"item.completed","item":{"type":"command_execution","command":"echo validation"}}`,
			true,
		},
		{
			// The feedback-to-implementation send_input alone — its prompt mentions
			// "validation" but it targets the IMPLEMENTATION thread, not the reviewer.
			// Must NOT pass: it is feedback routing, not reviewer reuse.
			"send_input to the implementation worker, not the reviewer",
			spawnValidation + "\n" + spawnImpl + "\n" + feedbackToImpl,
			true,
		},
		{
			// A send_input to the validation thread with NO correlating validation
			// spawn in the transcript: an uncorrelated thread id must not pass.
			"send_input to an uncorrelated thread",
			reuseValidation,
			true,
		},
		{"empty transcript", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := assertCodexReviewerReuse(tc.jsonl)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %q, got nil", tc.name)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected pass for %q, got: %v", tc.name, err)
			}
		})
	}
}
