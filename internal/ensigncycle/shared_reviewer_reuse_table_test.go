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

// TestClaudeSingleEntityRejectionFlow is the CONTRACT-correct single-entity (`-p`)
// reviewer producer-signal table (the option-(b) correction of the AC-3 finding,
// backlog seed e3z). It pins the deterministic bare-mode end-state — two distinct
// fresh validation spawns, fix-agent and reviewer separate — over committed,
// model-free transcripts, including the TWO observed non-deterministic live shapes
// (2-fresh-spawns and impl-reused-through-validation), so the validator's live
// re-run has an offline oracle.
//
// Root cause (verified in the contract, not assumed):
//   - The Claude runner launches `spacedock claude -- -p {prompt}` and the prompt
//     names one entity, so the run is single-entity → bare (first-officer-shared-
//     core.md `## Single-Entity Mode`; claude-fo-dispatch.md "In single-entity mode,
//     skip team creation. Use bare-mode dispatch for all agent spawning"). That
//     clause predates P2 (since the original vendoring `83c73494`), so the single-
//     entity reviewer is bare both before and after lazy-TeamCreate — the contradiction
//     is pre-existing; the old assertClaudeReviewerReuse encoded a team-mode
//     assumption the `-p` run can never satisfy.
//   - In bare mode the contract makes the flow DETERMINISTIC and SEQUENTIAL
//     (claude-fo-dispatch.md `## Feedback Rejection Flow (bare mode)`: "dispatch fix
//     agent (wait for completion), then dispatch reviewer (wait for completion)").
//     So the contract-correct end-state is two distinct fresh validation spawns with
//     the fix agent and reviewer as SEPARATE dispatches — which is exactly what
//     assertClaudeSingleEntityRejectionFlow asserts.
//
// NOTE the Codex contrast: Codex has no team registry, so its reviewer reuse via
// `send_input` to a persistent thread is contract-valid even in single-entity
// (context-dependent reuse — codex-first-officer-runtime.md `## Reuse And Feedback
// Routing`); assertCodexReviewerReuse stays correct for the Codex `-p` run. The
// Claude/Codex single-entity behaviors legitimately differ.
func TestClaudeSingleEntityRejectionFlow(t *testing.T) {
	// CONTRACT-CORRECT (bare mode): two distinct fresh validation spawns (cycle-1 +
	// cycle-2), the bare-mode-sequential end-state — PASS.
	bareCycle1Validation := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_BV1","input":{"description":"Rejection Task: validation","subagent_type":"spacedock:ensign"}}]}}`
	bareCycle2Validation := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_BV2","input":{"description":"Rejection Task: validation (cycle 2 fresh)","subagent_type":"spacedock:ensign"}}]}}`
	twoFreshSpawns := bareCycle1Validation + "\n" + bareCycle2Validation

	// VIOLATION: the cycle-2 re-review collapsed onto the implementation worker via a
	// SendMessage instructing it to validate its own rework — the impl-as-validator
	// shape. Must FAIL.
	implRework := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_BI","input":{"description":"Rejection Task: implementation rework","subagent_type":"spacedock:ensign"}}]}}`
	implAsValidator := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"spacedock-ensign-rejection-task-implementation","message":"now validate your own rework"}}]}}`
	implReusedThroughValidation := bareCycle1Validation + "\n" + implRework + "\n" + implAsValidator

	// VIOLATION: only one validation spawn, no second re-review at all. Must FAIL on
	// the spawn-count check.
	onlyCycle1 := bareCycle1Validation

	// CONTRACT-CORRECT (team mode supersede-teardown): two fresh validation spawns
	// PLUS a shutdown_request to the superseded cycle-1 implementation worker. The
	// shutdown is the contract-mandated supersede teardown, not a re-review routing —
	// it must NOT trip the impl-as-validator check. This is the original opus CI
	// false-positive shape (run 27861449587 att1). Must PASS.
	implShutdownObject := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"spacedock-ensign-rejection-task-implementation","message":{"type":"shutdown_request","reason":"superseded by fresh implementation-rework dispatch"}}}]}}`
	teamSupersedeTeardown := bareCycle1Validation + "\n" + implRework + "\n" + implShutdownObject + "\n" + bareCycle2Validation

	// CONTRACT-CORRECT (team mode reviewer reuse): one validation spawn, then the
	// cycle-2 re-review reuses the kept-alive `…-validation` reviewer via SendMessage,
	// with a shutdown_request reaping the superseded cycle-1 impl worker. This is the
	// local opus team-mode reuse shape. Must PASS.
	reviewerReuseMessage := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"spacedock-ensign-rejection-task-validation","message":"Advancing to next stage: validation (cycle 2 re-review). Re-review the updated entity state."}}]}}`
	teamReviewerReuse := bareCycle1Validation + "\n" + implRework + "\n" + implShutdownObject + "\n" + reviewerReuseMessage

	cases := []struct {
		name    string
		stream  string
		wantErr bool
	}{
		{"contract-correct two fresh validation spawns (bare-mode sequential)", twoFreshSpawns, false},
		{"impl reused through validation (impl-as-validator) must RED", implReusedThroughValidation, true},
		{"only the cycle-1 validation spawn (no re-review) must RED", onlyCycle1, true},
		{"empty stream must RED", "", true},
		{"team-mode supersede shutdown_request to impl handle must PASS", teamSupersedeTeardown, false},
		{"team-mode reviewer reuse (1 spawn + reuse message) must PASS", teamReviewerReuse, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := assertClaudeSingleEntityRejectionFlow(tc.stream)
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
	// The Codex collab_tool_call shape. A spawn_agent dispatching the validation
	// stage binds the reviewer's thread id (vThread); a later followup_task
	// (or legacy send_input fixture) to vThread is the cycle-2 reviewer reuse.
	// Both events, correlated, are the only passing shape.
	const vThread = "019e9695-1a3a-7481-be31-a1e81d53ec6d"
	const iThread = "019e9693-ad11-7701-b615-90d6a86b2dc8"
	spawnValidation := `{"type":"item.completed","item":{"type":"collab_tool_call","tool":"spawn_agent","receiver_thread_ids":["` + vThread + `"],"prompt":"Read /tmp/spacedock-dispatch/spacedock-ensign-rejection-task-validation.md and treat its content as your assignment."}}`
	spawnImpl := `{"type":"item.completed","item":{"type":"collab_tool_call","tool":"spawn_agent","receiver_thread_ids":["` + iThread + `"],"prompt":"Read /tmp/spacedock-dispatch/spacedock-ensign-rejection-task-implementation.md and treat its content as your assignment."}}`
	reuseValidation := `{"type":"item.started","item":{"type":"collab_tool_call","tool":"send_input","receiver_thread_ids":["` + vThread + `"],"prompt":"Re-run validation for rejection-task as cycle 2 using your existing validation reviewer context."}}`
	reuseValidationV2 := `{"type":"item.started","item":{"type":"collab_tool_call","tool":"followup_task","receiver_thread_ids":["` + vThread + `"],"prompt":"Re-run validation for rejection-task as cycle 2 using your existing validation reviewer context."}}`
	feedbackToImpl := `{"type":"item.started","item":{"type":"collab_tool_call","tool":"send_input","receiver_thread_ids":["` + iThread + `"],"prompt":"Feedback routed from validation to implementation for rejection-task. The fix marker is absent."}}`

	realReuse := spawnValidation + "\n" + spawnImpl + "\n" + feedbackToImpl + "\n" + reuseValidation
	realReuseV2 := spawnValidation + "\n" + spawnImpl + "\n" + feedbackToImpl + "\n" + reuseValidationV2

	// FRESH-DISPATCH false-pass (the cycle-8 M2 hole): the FO drives cycle-1
	// validation (spawn #1 → vThread), then FRESH-spawns a SECOND cycle-2 validation
	// reviewer (vThread2 — its prompt also names validation) and routes to THAT
	// new thread. The old assertion bound ANY validation-prompt spawn as "the
	// validation thread", so a message to the fresh cycle-2 thread passed as
	// "reuse". The strengthened assertion must RED it: two validation spawn_agents
	// means the FO fresh-dispatched, it did not reuse the kept-alive cycle-1 reviewer.
	const vThread2 = "019e9696-2b4b-8592-cf42-b2f92e64fd7e"
	freshCycle2Spawn := `{"type":"item.completed","item":{"type":"collab_tool_call","tool":"spawn_agent","receiver_thread_ids":["` + vThread2 + `"],"prompt":"Read /tmp/spacedock-dispatch/spacedock-ensign-rejection-task-validation.md and treat its content as your assignment."}}`
	sendInputToFresh := `{"type":"item.started","item":{"type":"collab_tool_call","tool":"send_input","receiver_thread_ids":["` + vThread2 + `"],"prompt":"Re-run validation for rejection-task as cycle 2."}}`
	freshDispatch := spawnValidation + "\n" + spawnImpl + "\n" + feedbackToImpl + "\n" + freshCycle2Spawn + "\n" + sendInputToFresh
	noAddressableSurface := `{"type":"item.completed","item":{"type":"agent_message","text":"The live Codex tool surface has spawn_agent and wait_agent, but no followup_task/send_message reuse route exposed here."}}`
	currentNoAddressableSurface := `{"type":"item.completed","item":{"type":"agent_message","text":"Codex has no turn-starting follow-up route for a completed worker, so there is no completed-worker follow-up route for validation reuse."}}`
	absentTwoFresh := noAddressableSurface + "\n" + spawnValidation + "\n" + spawnImpl + "\n" + freshCycle2Spawn
	currentAbsentTwoFresh := currentNoAddressableSurface + "\n" + spawnValidation + "\n" + spawnImpl + "\n" + freshCycle2Spawn
	absentOnlyOneValidation := noAddressableSurface + "\n" + spawnValidation + "\n" + spawnImpl
	absentButReuseTool := noAddressableSurface + "\n" + realReuseV2

	// The absence declarations a live Codex FO actually emitted when the host
	// surface exposed only spawn_agent/wait_agent and it contract-correctly
	// fresh-dispatched the cycle-2 reviewer. These are verbatim from the six
	// identical failing rejection-flow transcripts the #141 keepalive lane
	// produced. The brittle fixed-phrase absent detector matched NONE of them, so
	// the assertion took its PRESENT branch and RED-flagged a contract-correct
	// fresh-dispatch (the M7-shape false-positive). The detector must recognize a
	// genuine reuse-route-absent declaration however the FO words it.
	liveAbsentWordings := []string{
		"The available multi-agent surface exposes spawn_agent and wait_agent, but no follow-up/send-message binding.",
		"This host has no follow-up/send binding for worker reuse.",
		"Since the cycle-1 reviewer is not addressable on this host, I'm dispatching a fresh validation reviewer for cycle 2.",
		"This surface has spawn_agent and wait_agent, but no followup_task or equivalent turn-starting reuse route.",
		"Reviewer reuse is not supported by this tool surface.",
		"This host cannot address the kept-alive reviewer for follow-up.",
		"This host exposes the older multi-agent surface: I can spawn and wait for workers, but I do not have a followup_task/live reuse call.",
		"This Codex surface exposes spawn/wait but no addressable follow-up tool.",
		"Since this host does not expose a turn-starting addressable-worker binding, I'm fresh-dispatching validation instead of reusing the cycle-1 reviewer.",
		"This Codex host did not expose an addressable reviewer reuse/follow-up tool.",
	}
	liveAbsentTwoFresh := make([]string, len(liveAbsentWordings))
	for i, w := range liveAbsentWordings {
		decl := `{"type":"item.completed","item":{"type":"agent_message","text":` + mustJSONString(w) + `}}`
		liveAbsentTwoFresh[i] = decl + "\n" + spawnValidation + "\n" + spawnImpl + "\n" + freshCycle2Spawn
	}
	falseAbsentNarrations := []string{
		"The validation reviewer stays alive for the rejection-flow re-review; this does not close or redispatch the worker.",
		"The host has no dedicated shutdown tool, so I will keep the two reusable workers and route validation back to the kept-alive reviewer.",
		"The validation reviewer is being reused for re-review only; the implementation worker that applied the fix is not doing validation.",
	}
	falseAbsentReuse := make([]string, len(falseAbsentNarrations))
	for i, w := range falseAbsentNarrations {
		decl := `{"type":"item.completed","item":{"type":"agent_message","text":` + mustJSONString(w) + `}}`
		falseAbsentReuse[i] = decl + "\n" + realReuseV2
	}

	cases := []struct {
		name    string
		jsonl   string
		wantErr bool
	}{
		{"real send_input to the validation reviewer thread", realReuse, false},
		{"real followup_task to the validation reviewer thread", realReuseV2, false},
		{"fresh cycle-2 spawn_agent + send_input must RED", freshDispatch, true},
		{"addressable-worker absent: two fresh validation spawns", absentTwoFresh, false},
		{"addressable-worker absent: current live wording plus two fresh validation spawns", currentAbsentTwoFresh, false},
		{"addressable-worker absent: missing second fresh validation spawn must RED", absentOnlyOneValidation, true},
		{"addressable-worker absent: reuse tool contradicts absent classification", absentButReuseTool, true},
		{"live FO absence wording 0 + two fresh validation spawns", liveAbsentTwoFresh[0], false},
		{"live FO absence wording 1 + two fresh validation spawns", liveAbsentTwoFresh[1], false},
		{"live FO absence wording 2 + two fresh validation spawns", liveAbsentTwoFresh[2], false},
		{"live FO absence wording 3 + two fresh validation spawns", liveAbsentTwoFresh[3], false},
		{"live FO absence wording 4 + two fresh validation spawns", liveAbsentTwoFresh[4], false},
		{"live FO absence wording 5 + two fresh validation spawns", liveAbsentTwoFresh[5], false},
		{"live FO absence wording 6 + two fresh validation spawns", liveAbsentTwoFresh[6], false},
		{"live FO absence wording 7 + two fresh validation spawns", liveAbsentTwoFresh[7], false},
		{"live FO absence wording 8 + two fresh validation spawns", liveAbsentTwoFresh[8], false},
		{"live FO absence wording 9 + two fresh validation spawns", liveAbsentTwoFresh[9], false},
		{"PR #464 affirmative reuse with redispatch negation must not classify absent", falseAbsentReuse[0], false},
		{"PR #464 reusable-workers narration with shutdown negation must not classify absent", falseAbsentReuse[1], false},
		{"PR #465 validation reviewer reused while implementation worker is not validating", falseAbsentReuse[2], false},
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

func TestAssertCodexReviewerReuseAcceptsAdvanceModeValidationReroute(t *testing.T) {
	jsonl := strings.Join([]string{
		codexAgentMessageLine("The Codex runtime has a reusable worker route (`followup_task`), so I will keep the cycle-1 validation reviewer addressable."),
		codexCommandLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage implementation --checklist-file impl.checklist`),
		codexCommandLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage validation --checklist-file validation-cycle1.checklist`),
		codexWaitLine(),
		codexCommandLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage implementation --checklist-file rework.checklist --feedback-context-file rework.feedback --feedback-reflow --advance`),
		codexWaitLine(),
		codexCommandLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage validation --checklist-file validation-cycle2.checklist --advance`),
		codexAgentMessageLine("The validation re-review assignment is ready. I am routing it to the kept-alive cycle-1 validation reviewer now."),
		codexWaitLine(),
	}, "\n")

	if err := assertCodexReviewerReuse(jsonl); err != nil {
		t.Fatalf("Codex 0.142 live transcript with --advance validation re-review should pass: %v", err)
	}

	noValidationAdvance := strings.Replace(jsonl, " --stage validation --checklist-file validation-cycle2.checklist --advance", " --stage validation --checklist-file validation-cycle2.checklist", 1)
	if err := assertCodexReviewerReuse(noValidationAdvance); err == nil {
		t.Fatal("expected transcript without validation --advance re-review to fail")
	}
}

func TestAssertCodexReviewerReuseAcceptsLiveAdvanceModeReuseNarration(t *testing.T) {
	jsonl := strings.Join([]string{
		codexAgentMessageLine("The Codex runtime has `followup_task`, so the cycle-1 validation reviewer is reusable if it remains addressable."),
		codexCommandStartedLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage implementation --checklist-file impl-cycle1.checklist`),
		codexCommandLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage implementation --checklist-file impl-cycle1.checklist`),
		codexCommandStartedLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage validation --checklist-file validation-cycle1.checklist`),
		codexCommandLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage validation --checklist-file validation-cycle1.checklist`),
		codexWaitLine(),
		codexCommandStartedLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage implementation --checklist-file impl-cycle2.checklist --feedback-context-file cycle1-feedback.txt --feedback-reflow --advance`),
		codexCommandLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage implementation --checklist-file impl-cycle2.checklist --feedback-context-file cycle1-feedback.txt --feedback-reflow --advance`),
		codexWaitLine(),
		codexAgentMessageLine("The rework satisfies its checklist: 1 done, 0 skipped, 0 failed, and the standalone fix marker is present. I'm advancing back to validation and reusing the kept-alive cycle-1 validation reviewer for the second-cycle re-review."),
		codexCommandStartedLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage validation --checklist-file validation-cycle2.checklist --feedback-context-file cycle2-review.txt --advance`),
		codexCommandLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage validation --checklist-file validation-cycle2.checklist --feedback-context-file cycle2-review.txt --advance`),
		codexAgentMessageLine("The second-cycle validation transition is committed. I'm sending the re-review to the kept-alive validation reviewer, which keeps the fix worker separate from the reviewer."),
		codexWaitLine(),
	}, "\n")

	if err := assertCodexReviewerReuse(jsonl); err != nil {
		t.Fatalf("Codex live advance-mode reuse narration should pass: %v", err)
	}
}

func codexAgentMessageLine(text string) string {
	return `{"type":"item.completed","item":{"type":"agent_message","text":` + mustJSONString(text) + `}}`
}

func codexCommandLine(command string) string {
	return `{"type":"item.completed","item":{"type":"command_execution","command":` + mustJSONString(command) + `,"status":"completed","exit_code":0}}`
}

func codexCommandStartedLine(command string) string {
	return `{"type":"item.started","item":{"type":"command_execution","command":` + mustJSONString(command) + `,"status":"in_progress"}}`
}

func codexWaitLine() string {
	return `{"type":"item.completed","item":{"type":"collab_tool_call","tool":"wait","receiver_thread_ids":[],"status":"completed"}}`
}
