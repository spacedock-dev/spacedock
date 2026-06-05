package ensigncycle

import "testing"

// Offline table tests for the host-specific reviewer-reuse assertions. They prove
// each assertion requires a REAL reuse tool call targeting the validation reviewer
// — not loose narration, not an unrelated tool, not a call to a different
// recipient. Without these the assertions are exercised only by the model-spending
// live runners; here they cost milliseconds and pin the discriminating behavior.

func TestAssertClaudeReviewerReuse(t *testing.T) {
	// A real SendMessage tool_use whose `to` names the validation reviewer — the
	// name-addressed reuse shape.
	realReuse := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"spacedock-ensign-validation","message":"re-review the fix"}}]}}`

	// The agentId-addressed reuse shape — the actual shape the FO emits in a Claude
	// team. A validation-stage Agent spawn returns an agentId on completion; the FO
	// resumes that completed reviewer by agentId (not name) for the cycle-2 re-review.
	// Both events must be present and correlated for this to pass.
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

	cases := []struct {
		name    string
		stream  string
		wantErr bool
	}{
		{"real SendMessage to validation by name", realReuse, false},
		{"reuse by correlated validation agentId", agentIDReuse, false},
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

	cases := []struct {
		name    string
		jsonl   string
		wantErr bool
	}{
		{"real send_input to the validation reviewer thread", realReuse, false},
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
