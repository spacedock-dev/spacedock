package ensigncycle

import (
	"errors"
	"strings"
	"testing"
)

// Offline table tests for the host-specific reviewer-reuse assertions. They prove
// each assertion reads reviewer identity ONLY from structured native handles — the
// Claude validation spawn's agentId/input.name correlated to a SendMessage.to, the
// Codex validation spawn_agent's receiver_thread_ids correlated to a followup_task/
// send_input — and reports an explicit errReviewerIdentityUnsupported when no such
// handle exists. Command strings, wait counts, durable reports, and free-form
// narration never manufacture a reuse claim. Without these the assertions are
// exercised only by the model-spending live runners; here they cost milliseconds and
// pin the discriminating behavior.

func TestAssertClaudeReviewerReuse(t *testing.T) {
	// The agentId-addressed reuse shape — one shape the FO emits in a Claude team
	// (a teams-mode kept-alive completed agent can be resumed by the agentId its
	// spawn returned). A validation-stage Agent spawn (exactly one) returns an
	// agentId on completion; the FO resumes THAT reviewer by agentId for the cycle-2
	// re-review. Both events must be present and correlated for this to pass.
	validationSpawn := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_V","input":{"description":"Rejection Task: validation","subagent_type":"spacedock:ensign","name":"spacedock-ensign-rejection-task-validation"}}]}}`
	validationResult := `{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_V","content":[{"type":"text","text":"Verdict: REJECTED."},{"type":"text","text":"agentId: a94abe89c85f9f4cc (use SendMessage with to: 'a94abe89c85f9f4cc' to continue this agent)"}]}]}}`
	reuseByAgentID := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"a94abe89c85f9f4cc","message":"re-review: the fix marker is now present"}}]}}`
	agentIDReuse := validationSpawn + "\n" + validationResult + "\n" + reuseByAgentID

	// An agentId from an IMPLEMENTATION spawn, not the validation reviewer: a
	// SendMessage to it is NOT reviewer reuse and must NOT pass. There is no
	// validation spawn to correlate, so the identity surface is unsupported.
	implSpawn := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_I","input":{"description":"Rejection Task: implementation","subagent_type":"spacedock:ensign","name":"spacedock-ensign-rejection-task-implementation"}}]}}`
	implResult := `{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_I","content":[{"type":"text","text":"agentId: a000impl0000 (use SendMessage with to: 'a000impl0000' to continue this agent)"}]}]}}`
	reuseImplAgentID := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"a000impl0000","message":"apply the fix"}}]}}`
	wrongAgentIDReuse := implSpawn + "\n" + implResult + "\n" + reuseImplAgentID

	// A SendMessage to a validation-reviewer agentId that was NEVER spawned/returned
	// in this transcript: an opaque handle with no correlating validation spawn — the
	// identity surface is unsupported, not a reuse pass.
	uncorrelatedAgentID := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"adeadbeef0000000","message":"re-review"}}]}}`

	// Genuine NAME-addressed reuse: the FO spawns the cycle-1 validation reviewer
	// ONCE declaring its teammate handle in input.name, then re-engages it for cycle 2
	// by that exact handle (the production shape — every recorded spawn carries
	// input.name and SendMessage.to matches it). Exactly one validation spawn + a
	// name-correlated message to it = genuine reuse → PASS.
	nameSpawnValidation := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_NV","input":{"description":"Rejection Task: validation","subagent_type":"spacedock:ensign","name":"spacedock-ensign-rejection-task-validation"}}]}}`
	reuseByName := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"spacedock-ensign-rejection-task-validation","message":"re-review: the fix marker is now present"}}]}}`
	nameReuse := nameSpawnValidation + "\n" + reuseByName

	// FRESH-DISPATCH (two validation spawns): the FO drives cycle-1 validation
	// (spawn #1), then FRESH-dispatches the cycle-2 validator (a SECOND validation
	// Agent spawn) and kicks it off by its validation NAME. Two validation spawns is
	// structurally FRESH, not reuse — the team-mode #141 keepalive assertion RED-flags
	// it (the reviewer was not kept alive and reused).
	freshCycle2Spawn := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_FV","input":{"description":"Rejection Task: validation (cycle 2)","subagent_type":"spacedock:ensign","name":"spacedock-ensign-rejection-task-validation-cycle2"}}]}}`
	freshDispatchNameMessage := nameSpawnValidation + "\n" + freshCycle2Spawn + "\n" + reuseByName

	cases := []struct {
		name            string
		stream          string
		wantErr         bool
		wantUnsupported bool
	}{
		{"genuine name-addressed reuse (one validation spawn)", nameReuse, false, false},
		{"reuse by correlated validation agentId", agentIDReuse, false, false},
		{"fresh cycle-2 dispatch (two validation spawns) must RED as fresh", freshDispatchNameMessage, true, false},
		{
			"loose narration only",
			`{"type":"assistant","message":{"content":[{"type":"text","text":"I will reuse the validation reviewer via SendMessage."}]}}`,
			true, true,
		},
		{
			"validation spawn exposing no agentId or input.name handle",
			`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","input":{"to":"spacedock-ensign-validation","description":"fresh validation dispatch"}}]}}`,
			true, true,
		},
		{
			"SendMessage to a non-validation recipient with no validation spawn",
			`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"spacedock-ensign-implementation","message":"apply the fix"}}]}}`,
			true, true,
		},
		{"reuse of an implementation agentId, no validation spawn", wrongAgentIDReuse, true, true},
		{"SendMessage to an uncorrelated agentId, no validation spawn", uncorrelatedAgentID, true, true},
		{"empty stream", "", true, true},
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
			if got := errors.Is(err, errReviewerIdentityUnsupported); got != tc.wantUnsupported {
				t.Fatalf("errReviewerIdentityUnsupported = %v for %q, want %v (err: %v)", got, tc.name, tc.wantUnsupported, err)
			}
		})
	}
}

// TestClaudeSingleEntityRejectionFlow is the single-entity (`-p`) reviewer
// producer-signal table. It pins the two contract-valid end-states — bare-mode two
// distinct fresh validation spawns, and team-mode one spawn reused via SendMessage —
// over committed, model-free transcripts, so the validator's live re-run has an
// offline oracle.
//
// Root cause (verified in the contract, not assumed):
//   - The Claude runner launches `spacedock claude -- -p {prompt}` and the prompt
//     names one entity, so the run is single-entity → bare (first-officer-shared-
//     core.md `## Single-Entity Mode`; claude-fo-dispatch.md "In single-entity mode,
//     skip team creation. Use bare-mode dispatch for all agent spawning"), which
//     fresh-dispatches a distinct reviewer per cycle — two validation spawns. When
//     the `-p` FO opts into team mode it keeps the cycle-1 reviewer alive and reuses
//     it. Both reach a validation worker; neither collapses onto the fix worker.
//
// NOTE the Codex contrast: Codex has no team registry, so its reviewer reuse via
// `send_input` to a persistent thread is contract-valid even in single-entity
// (context-dependent reuse — codex-first-officer-runtime.md `## Reuse And Feedback
// Routing`); assertCodexReviewerReuse stays correct for the Codex `-p` run. The
// Claude/Codex single-entity behaviors legitimately differ.
func TestAssertCodexReviewerReuse(t *testing.T) {
	// The Codex collab_tool_call shape. A spawn_agent dispatching the validation
	// stage binds the reviewer's thread id (vThread); a later followup_task
	// (or legacy send_input fixture) to vThread is the cycle-2 reviewer reuse.
	// Both events, correlated, are the only reuse-passing shape.
	const vThread = "019e9695-1a3a-7481-be31-a1e81d53ec6d"
	const iThread = "019e9693-ad11-7701-b615-90d6a86b2dc8"
	const uThread = "019e9698-3c5c-96a3-d053-c3fa3f75fe8f"
	spawnValidation := `{"type":"item.completed","item":{"type":"collab_tool_call","tool":"spawn_agent","receiver_thread_ids":["` + vThread + `"],"prompt":"Read /tmp/spacedock-dispatch/spacedock-ensign-rejection-task-validation.md and treat its content as your assignment."}}`
	spawnImpl := `{"type":"item.completed","item":{"type":"collab_tool_call","tool":"spawn_agent","receiver_thread_ids":["` + iThread + `"],"prompt":"Read /tmp/spacedock-dispatch/spacedock-ensign-rejection-task-implementation.md and treat its content as your assignment."}}`
	reuseValidation := `{"type":"item.started","item":{"type":"collab_tool_call","tool":"send_input","receiver_thread_ids":["` + vThread + `"],"prompt":"Re-run validation for rejection-task as cycle 2 using your existing validation reviewer context."}}`
	reuseValidationV2 := `{"type":"item.started","item":{"type":"collab_tool_call","tool":"followup_task","receiver_thread_ids":["` + vThread + `"],"prompt":"Re-run validation for rejection-task as cycle 2 using your existing validation reviewer context."}}`
	feedbackToImpl := `{"type":"item.started","item":{"type":"collab_tool_call","tool":"send_input","receiver_thread_ids":["` + iThread + `"],"prompt":"Feedback routed from validation to implementation for rejection-task. The fix marker is absent."}}`

	realReuse := spawnValidation + "\n" + spawnImpl + "\n" + feedbackToImpl + "\n" + reuseValidation
	realReuseV2 := spawnValidation + "\n" + spawnImpl + "\n" + feedbackToImpl + "\n" + reuseValidationV2

	// Two distinct validation spawn_agents with NO narration line: the FO
	// fresh-dispatched a cycle-2 reviewer distinct from cycle-1. This is structurally
	// FRESH — the contract-correct choice when the host exposes no addressable-worker
	// reuse route — and passes without any narration proving it (narration is no
	// longer load-bearing).
	const vThread2 = "019e9696-2b4b-8592-cf42-b2f92e64fd7e"
	freshCycle2Spawn := `{"type":"item.completed","item":{"type":"collab_tool_call","tool":"spawn_agent","receiver_thread_ids":["` + vThread2 + `"],"prompt":"Read /tmp/spacedock-dispatch/spacedock-ensign-rejection-task-validation.md and treat its content as your assignment."}}`
	twoFreshValidation := spawnValidation + "\n" + spawnImpl + "\n" + freshCycle2Spawn

	// A follow-up to an uncorrelated thread while a validation reviewer exists: the
	// re-review was routed to a thread that is not the bound validation reviewer's —
	// a wrong-reviewer route, not reuse.
	reuseUncorrelated := `{"type":"item.started","item":{"type":"collab_tool_call","tool":"send_input","receiver_thread_ids":["` + uThread + `"],"prompt":"Re-run validation for rejection-task as cycle 2."}}`
	followupToUncorrelated := spawnValidation + "\n" + reuseUncorrelated

	cases := []struct {
		name            string
		jsonl           string
		wantErr         bool
		wantUnsupported bool
	}{
		{"real send_input to the validation reviewer thread", realReuse, false, false},
		{"real followup_task to the validation reviewer thread", realReuseV2, false, false},
		{"two distinct validation spawn_agents (fresh reviewers, no narration)", twoFreshValidation, false, false},
		{
			// The feedback-to-implementation send_input alone — its prompt mentions
			// "validation" but it targets the IMPLEMENTATION thread, not the reviewer.
			// A validation reviewer exists, so this is a wrong-reviewer route (no_reuse),
			// not reuse.
			"send_input to the implementation worker, not the reviewer",
			spawnValidation + "\n" + spawnImpl + "\n" + feedbackToImpl,
			true, false,
		},
		{"followup_task to an uncorrelated thread while a validation reviewer exists", followupToUncorrelated, true, false},
		{
			// A validation spawn with no reuse tool-call at all — the reviewer was
			// created but never reused.
			"validation spawn with no reuse follow-up",
			spawnValidation + "\n" + `{"type":"item.completed","item":{"type":"command_execution","command":"echo validation"}}`,
			true, false,
		},
		{
			// A send_input to the validation thread with NO correlating validation
			// spawn in the transcript: an uncorrelated thread id with no structured
			// spawn to bind it is unsupported.
			"send_input to a thread with no correlating validation spawn",
			reuseValidation,
			true, true,
		},
		{
			"loose narration only",
			`{"type":"item.completed","item":{"type":"agent_message","text":"I will send_input to the validation worker."}}`,
			true, true,
		},
		{"empty transcript", "", true, true},
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
			if got := errors.Is(err, errReviewerIdentityUnsupported); got != tc.wantUnsupported {
				t.Fatalf("errReviewerIdentityUnsupported = %v for %q, want %v (err: %v)", got, tc.name, tc.wantUnsupported, err)
			}
		})
	}
}

// The command-string/wait and narration transcripts the previous oracle graded as a
// reuse pass by inference. They carry no spawn_agent/thread handle, so structured
// correlation cannot prove WHO re-reviewed — the honest outcome is unsupported.

func codexAdvanceValidationRerouteJSONL() string {
	return strings.Join([]string{
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
}

func codexCIAdvanceJSONL() string {
	return strings.Join([]string{
		codexAgentMessageLine("Because Codex has an addressable worker route, I am routing the rework back to the existing implementation worker rather than spawning the reviewer to do fix work."),
		codexCommandLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage implementation --checklist-file impl.checklist`),
		codexCommandLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage validation --checklist-file validation-cycle1.checklist`),
		codexWaitLine(),
		codexCommandLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage implementation --checklist-file rework.checklist --advance`),
		codexWaitLine(),
		codexAgentMessageLine("I am advancing back to validation and reusing the kept-alive cycle-1 validation reviewer for the re-review, keeping it separate from the worker that applied the fix."),
		codexCommandLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage validation --checklist-file validation-cycle2.checklist --advance`),
		codexAgentMessageLine("The cycle-2 validation assignment is built. I am routing it to the original validation reviewer through the addressable worker handle, not to the implementation rework worker."),
		codexWaitLine(),
	}, "\n")
}

func codexLiveAdvanceNarrationJSONL() string {
	return strings.Join([]string{
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
}

func codexV2AssignmentJSONL() string {
	return strings.Join([]string{
		codexAgentMessageLine("Codex exec does not expose spawn metadata here, so I am grading the durable assignments and keeping the validation reviewer separate."),
		codexCommandLine(`printf '{"entity_path":"rejection-task.md","stage":"implementation"}' | ${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --json`),
		codexCommandLine(`printf '{"entity_path":"rejection-task.md","stage":"validation"}' | ${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --json`),
		codexWaitLine(),
		codexCommandLine(`printf '{"entity_path":"rejection-task.md","stage":"implementation","advance":true}' | ${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --json`),
		codexWaitLine(),
		codexAgentMessageLine("The cycle-1 reviewer stayed available; I kept validation separate for the re-review."),
		codexCommandLine(`printf '{"entity_path":"rejection-task.md","stage":"validation","advance":true}' | ${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --json`),
		codexWaitLine(),
	}, "\n")
}

// codexDurableCommandWaitJSONL is the command+wait transcript the previous oracle
// paired with a durable validation-report entity to claim reuse. Identity was never
// in the transcript — the durable report proves a re-review OCCURRED, not who did it.
func codexDurableCommandWaitJSONL() string {
	return strings.Join([]string{
		codexCommandLine(`${SPACEDOCK_BIN:-spacedock} dispatch build --workflow-dir . --entity-path rejection-task.md --stage implementation`),
		codexWaitLine(),
	}, "\n")
}

func codexNarrationOnlyJSONL() string {
	return strings.Join([]string{
		codexAgentMessageLine("I am using a separate validation reviewer and will keep it alive for the re-review."),
		codexAgentMessageLine("The separate validation reviewer passed the re-review."),
	}, "\n")
}

// TestAssertCodexReviewerReuseRetiresCommandStringWaitInference converts the former
// advance-mode "accepts" tests: command-string + wait transcripts carry no structured
// spawn/thread handle, so identity is unsupported — never a reuse pass (AC-2).
func TestAssertCodexReviewerReuseRetiresCommandStringWaitInference(t *testing.T) {
	for name, jsonl := range map[string]string{
		"advance-mode validation reroute":         codexAdvanceValidationRerouteJSONL(),
		"CI advance-mode without feedback reflow": codexCIAdvanceJSONL(),
		"live advance-mode reuse narration":       codexLiveAdvanceNarrationJSONL(),
		"current-v2 assignment surfaces":          codexV2AssignmentJSONL(),
	} {
		err := assertCodexReviewerReuse(jsonl)
		if !errors.Is(err, errReviewerIdentityUnsupported) {
			t.Fatalf("%s: command-string/wait inference must yield errReviewerIdentityUnsupported, got: %v", name, err)
		}
	}
}

// TestAssertCodexReviewerReuseRetiresNarrationAndDurableInference converts the former
// durable-state and narration "accepts" tests: neither a durable validation report
// nor free-form narration carries a structured reviewer handle, so identity is
// unsupported — never a reuse pass (AC-2, AC-3).
func TestAssertCodexReviewerReuseRetiresNarrationAndDurableInference(t *testing.T) {
	for name, jsonl := range map[string]string{
		"durable-report command+wait shape": codexDurableCommandWaitJSONL(),
		"narration only":                    codexNarrationOnlyJSONL(),
	} {
		err := assertCodexReviewerReuse(jsonl)
		if !errors.Is(err, errReviewerIdentityUnsupported) {
			t.Fatalf("%s: inference must yield errReviewerIdentityUnsupported, got: %v", name, err)
		}
	}
}

// TestReviewerReuseAC1InferenceRetirement is the AC-1 value measurement. It
// enumerates the committed inference-only fixtures the previous oracle graded as a
// reuse pass and asserts the count that now yields a reuse pass is 0, while every
// structured fixture (Codex thread-id follow-up; Claude agentId/input.name
// SendMessage) still passes. Baseline under the previous oracle: 5 inference-only
// fixtures yielded a reuse pass (the 3 advance-mode command/wait fixtures, the
// current-v2 assignment fixture, and the durable-report fixture); the structured
// fixtures passed. New count: 0.
func TestReviewerReuseAC1InferenceRetirement(t *testing.T) {
	inferenceOnly := []struct {
		name  string
		jsonl string
	}{
		{"advance-mode validation reroute (command-string/wait)", codexAdvanceValidationRerouteJSONL()},
		{"CI advance-mode (command-string/wait)", codexCIAdvanceJSONL()},
		{"live advance-mode reuse narration (command-string/wait)", codexLiveAdvanceNarrationJSONL()},
		{"current-v2 assignment surfaces (command-string/wait)", codexV2AssignmentJSONL()},
		{"durable-report command+wait shape", codexDurableCommandWaitJSONL()},
		{"narration only", codexNarrationOnlyJSONL()},
	}
	inferenceReusePasses := 0
	for _, tc := range inferenceOnly {
		result, _ := codexReviewerIdentity(tc.jsonl)
		if result == reviewerReuse {
			inferenceReusePasses++
			t.Errorf("inference-only fixture %q yielded a reuse pass — identity must not be reconstructed from command strings, wait counts, durable reports, or narration", tc.name)
		}
		if !errors.Is(assertCodexReviewerReuse(tc.jsonl), errReviewerIdentityUnsupported) {
			t.Errorf("inference-only fixture %q must yield errReviewerIdentityUnsupported", tc.name)
		}
	}
	if inferenceReusePasses != 0 {
		t.Fatalf("AC-1: inference-only reuse-pass count = %d, want 0 (baseline was 5 under the previous oracle)", inferenceReusePasses)
	}

	// Structured fixtures still pass. The same-handle reuse fixtures grade as
	// reviewerReuse specifically; the two-distinct-spawn fixture is a fresh pass.
	claudeAgentIDReuse := strings.Join([]string{
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_V","input":{"description":"Rejection Task: validation","name":"spacedock-ensign-rejection-task-validation"}}]}}`,
		`{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_V","content":[{"type":"text","text":"agentId: a94abe89c85f9f4cc (use SendMessage with to: 'a94abe89c85f9f4cc')"}]}]}}`,
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"a94abe89c85f9f4cc","message":"re-review"}}]}}`,
	}, "\n")
	claudeNameReuse := strings.Join([]string{
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":"toolu_NV","input":{"description":"Rejection Task: validation","name":"spacedock-ensign-rejection-task-validation"}}]}}`,
		`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":"spacedock-ensign-rejection-task-validation","message":"re-review"}}]}}`,
	}, "\n")

	structured := []struct {
		name      string
		host      string
		stream    string
		wantReuse bool
	}{
		{"Codex thread-id send_input/followup reuse", "codex", codexReviewerReuseJSONL("thread-validation-1", "thread-implementation-1"), true},
		{"Claude agentId-correlated SendMessage reuse", "claude", claudeAgentIDReuse, true},
		{"Claude input.name-correlated SendMessage reuse", "claude", claudeNameReuse, true},
	}
	structuredPasses := 0
	structuredReuseGrades := 0
	for _, tc := range structured {
		var assertErr error
		var result reviewerIdentityResult
		switch tc.host {
		case "codex":
			assertErr = assertCodexReviewerReuse(tc.stream)
			result, _ = codexReviewerIdentity(tc.stream)
		default:
			assertErr = assertClaudeReviewerReuse(tc.stream)
			result, _ = claudeReviewerIdentity(tc.stream)
		}
		if assertErr == nil {
			structuredPasses++
		} else {
			t.Errorf("structured fixture %q must still pass, got: %v", tc.name, assertErr)
		}
		if tc.wantReuse && result == reviewerReuse {
			structuredReuseGrades++
		} else if tc.wantReuse {
			t.Errorf("structured fixture %q must grade as reviewerReuse, got result %d", tc.name, result)
		}
	}
	if structuredPasses != len(structured) {
		t.Fatalf("AC-1: structured-fixture pass count = %d, want %d (unchanged)", structuredPasses, len(structured))
	}
	if structuredReuseGrades != len(structured) {
		t.Fatalf("AC-1: structured same-handle reuse-grade count = %d, want %d", structuredReuseGrades, len(structured))
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

func codexCollabToolLine(eventType, tool, threadID, prompt string) string {
	return `{"type":` + mustJSONString(eventType) +
		`,"item":{"type":"collab_tool_call","tool":` + mustJSONString(tool) +
		`,"receiver_thread_ids":[` + mustJSONString(threadID) +
		`],"prompt":` + mustJSONString(prompt) + `}}`
}

func codexReviewerReuseJSONL(validationThread, implementationThread string) string {
	return strings.Join([]string{
		codexCollabToolLine("item.completed", "spawn_agent", validationThread, "Read /tmp/spacedock-dispatch/spacedock-ensign-rejection-task-validation.md and treat its content as your assignment."),
		codexCollabToolLine("item.completed", "spawn_agent", implementationThread, "Read /tmp/spacedock-dispatch/spacedock-ensign-rejection-task-implementation.md and treat its content as your assignment."),
		codexCollabToolLine("item.started", "send_input", implementationThread, "Feedback routed from validation to implementation for rejection-task. The fix marker is absent."),
		codexCollabToolLine("item.started", "send_input", validationThread, "Re-run validation for rejection-task as cycle 2 using your existing validation reviewer context."),
	}, "\n")
}
