package ensigncycle

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// The host-specific reviewer-reuse assertions for the rejection-flow scenario.
// They are the producer signal for who performed the cycle-2 re-review after a
// cycle-1 rejection: the FO must reuse the kept-alive cycle-1 validation reviewer
// (the #141 keepalive contract) or, on a host with no addressable-worker route,
// fresh-dispatch a distinct cycle-2 reviewer — either way the reviewer and the
// fix worker stay separate. Reviewer identity is read ONLY from native structured
// spawn/follow-up handles: Claude correlates the validation Agent/Task spawn's
// returned agentId and declared input.name to a later SendMessage.to; Codex binds
// the validation spawn_agent's receiver_thread_ids and correlates a followup_task/
// send_input to that thread. When a transcript carries no such handle, identity is
// reported as unsupported (errReviewerIdentityUnsupported) — never a reuse pass;
// a durable report can prove a re-review OCCURRED but not WHO performed it. These
// live under the DEFAULT build tags (parsing only stdlib JSON) so the offline table
// tests exercise them without spending a model, alongside the //go:build live
// runners that feed them real transcripts.

// reviewerIdentityResult is the outcome of a per-host structured identity
// correlation over one rejection-flow transcript.
type reviewerIdentityResult int

const (
	// reviewerReuse: exactly one validation reviewer, re-review routed back to its
	// own structured handle (Claude agentId/name, Codex bound thread).
	reviewerReuse reviewerIdentityResult = iota
	// reviewerFresh: two or more distinct validation reviewers — the re-review used
	// a fresh reviewer, structurally proven distinct (the contract-correct choice on
	// a host with no addressable-worker reuse route).
	reviewerFresh
	// reviewerNoReuse: a validation reviewer exists but the re-review was routed to a
	// non-validation handle, or no re-review reached the reviewer at all.
	reviewerNoReuse
	// reviewerUnsupported: the transcript exposes no structured spawn/follow-up handle
	// to correlate identity from. Surfaced as errReviewerIdentityUnsupported; never a
	// reuse pass.
	reviewerUnsupported
)

// errReviewerIdentityUnsupported is the sentinel an assertion returns (errors.Is
// checkable) when the transcript carries no structured handle that could prove which
// process performed the re-review. Identity is never derived from command strings,
// wait counts, durable reports, or free-form narration.
var errReviewerIdentityUnsupported = errors.New("reviewer identity unsupported: no structured spawn/follow-up handle to correlate")

// claudeAgentIDResult extracts the `agentId: <id>` a completed subagent's
// tool_result returns ("agentId: a94abe89c85f9f4cc (use SendMessage with to: …)").
// The teams-mode reuse handle for a kept-alive completed agent.
var claudeAgentIDResult = regexp.MustCompile(`agentId:\s*(a[0-9a-f]+)`)

// claudeReviewerIdentity correlates the Claude stream-json transcript's validation
// reviewer identity structurally. It binds each validation-stage Agent/Task spawn to
// the handles that spawn exposes — the agentId its tool_result returns AND the
// teammate handle declared in its input.name — then reads a later SendMessage.to
// against those handles:
//
//   - exactly one validation spawn + a SendMessage.to equal to one of its correlated
//     handles → reviewerReuse (a shutdown_request is teardown, not a re-review, so it
//     is skipped);
//   - two or more validation spawns → reviewerFresh (distinct reviewers by count);
//   - exactly one validation spawn but the re-review reached a non-validation handle,
//     or no SendMessage reached it → reviewerNoReuse;
//   - no validation spawn, or a single validation spawn exposing neither an agentId
//     nor an input.name handle → reviewerUnsupported.
func claudeReviewerIdentity(stream string) (reviewerIdentityResult, string) {
	validationSpawnIDs := map[string]bool{} // tool_use id of each validation-stage spawn
	validationSpawnCount := 0
	validationHandles := map[string]bool{} // input.name + returned agentId of the validation spawns
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry struct {
			Message *struct {
				Content []struct {
					Type      string `json:"type"`
					Name      string `json:"name"`
					ID        string `json:"id"`
					ToolUseID string `json:"tool_use_id"`
					Input     struct {
						Description string `json:"description"`
						Name        string `json:"name"`
					} `json:"input"`
					Content json.RawMessage `json:"content"`
				} `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err != nil || entry.Message == nil {
			continue
		}
		for _, block := range entry.Message.Content {
			if block.Type == "tool_use" && (block.Name == "Agent" || block.Name == "Task") &&
				strings.Contains(strings.ToLower(block.Input.Description), "validation") {
				validationSpawnCount++
				validationSpawnIDs[block.ID] = true
				if block.Input.Name != "" {
					validationHandles[block.Input.Name] = true
				}
			}
			if block.Type == "tool_result" && validationSpawnIDs[block.ToolUseID] {
				if m := claudeAgentIDResult.FindStringSubmatch(string(block.Content)); m != nil {
					validationHandles[m[1]] = true
				}
			}
		}
	}

	if validationSpawnCount == 0 {
		return reviewerUnsupported, "no validation-stage Agent/Task spawn — no reviewer to correlate"
	}
	if validationSpawnCount >= 2 {
		return reviewerFresh, fmt.Sprintf("%d distinct validation-stage spawns — fresh reviewers", validationSpawnCount)
	}
	if len(validationHandles) == 0 {
		return reviewerUnsupported, "the validation spawn exposed neither an agentId nor an input.name handle"
	}

	reuseHit := false
	wrongRoute := false
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry struct {
			Message *struct {
				Content []struct {
					Type  string `json:"type"`
					Name  string `json:"name"`
					Input struct {
						To      string          `json:"to"`
						Message json.RawMessage `json:"message"`
					} `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err != nil || entry.Message == nil {
			continue
		}
		for _, block := range entry.Message.Content {
			if block.Type != "tool_use" || block.Name != "SendMessage" {
				continue
			}
			// A shutdown_request is the contract-mandated supersede teardown of a
			// prior worker, not a re-review routing — skip it, whichever handle it
			// targets.
			if isShutdownRequest(block.Input.Message) {
				continue
			}
			if validationHandles[block.Input.To] {
				reuseHit = true
			} else if block.Input.To != "" {
				wrongRoute = true
			}
		}
	}
	if reuseHit {
		return reviewerReuse, "SendMessage to the correlated validation reviewer handle"
	}
	if wrongRoute {
		return reviewerNoReuse, "the cycle-2 re-review was routed to a non-validation handle while a validation reviewer existed"
	}
	return reviewerNoReuse, "one validation reviewer spawned but no reuse SendMessage reached its handle"
}

// assertClaudeReviewerReuse enforces the team-mode #141 keepalive contract: the FO
// reused the kept-alive cycle-1 validation reviewer for the cycle-2 re-review. Only a
// structured reuse to the correlated reviewer handle passes; a second validation
// spawn (fresh dispatch), a re-review routed elsewhere, or an absent identity handle
// all fail.
func assertClaudeReviewerReuse(stream string) error {
	result, detail := claudeReviewerIdentity(stream)
	switch result {
	case reviewerReuse:
		return nil
	case reviewerFresh:
		return fmt.Errorf("the FO fresh-dispatched the cycle-2 validator instead of reusing the kept-alive cycle-1 reviewer (the #141 keepalive contract violation): %s", detail)
	case reviewerNoReuse:
		return fmt.Errorf("the FO did not reuse the kept-alive validation reviewer for the cycle-2 re-review: %s", detail)
	default:
		return fmt.Errorf("%w: %s", errReviewerIdentityUnsupported, detail)
	}
}

// isShutdownRequest reports whether a SendMessage `message` payload is a
// cooperative teardown (`{"type":"shutdown_request",...}`), sent either as a JSON
// object or as a JSON string carrying that type. The supersede-shutdown contract
// requires the FO to reap a superseded worker by this message, so a shutdown_request
// to a worker handle is mandatory teardown, not a re-review routing.
func isShutdownRequest(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var obj struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(raw, &obj); err == nil && obj.Type == "shutdown_request" {
		return true
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		var inner struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal([]byte(s), &inner); err == nil && inner.Type == "shutdown_request" {
			return true
		}
		return strings.Contains(s, "shutdown_request")
	}
	return false
}

// codexCollabItem is one `codex exec --json` stream item. Codex surfaces its
// multi-agent calls as `collab_tool_call` items (tool = spawn_agent / followup_task /
// send_input for legacy pre-v2 fixtures / wait / close_agent); the worker is addressed by opaque `receiver_thread_ids`,
// not by a name. The validation reviewer's thread is bound by the spawn_agent whose
// prompt dispatches the validation stage (its completed item carries the spawned
// thread id in receiver_thread_ids).
type codexCollabItem struct {
	Type string `json:"type"`
	Item struct {
		Type              string   `json:"type"`
		Tool              string   `json:"tool"`
		ReceiverThreadIDs []string `json:"receiver_thread_ids"`
		Prompt            string   `json:"prompt"`
	} `json:"item"`
}

// codexReviewerIdentity correlates the `codex exec --json` transcript's validation
// reviewer identity structurally. It binds each validation spawn_agent
// (`item.completed`, prompt dispatches validation) to its receiver_thread_ids, then
// reads the follow-up routing:
//
//   - exactly one validation spawn + a followup_task/send_input to its bound thread
//     → reviewerReuse (a feedback send_input to the implementation thread is normal
//     routing and does not itself count as reuse);
//   - two or more distinct validation spawn_agents → reviewerFresh (distinct
//     reviewers by count — the contract-correct behavior when the host exposes no
//     addressable-worker reuse route);
//   - exactly one validation spawn but the follow-up reached a non-validation thread,
//     or no follow-up reached the reviewer → reviewerNoReuse;
//   - no validation spawn_agent, or a validation spawn exposing no thread id → no
//     structured handle to correlate → reviewerUnsupported.
func codexReviewerIdentity(jsonl string) (reviewerIdentityResult, string) {
	validationThreads := map[string]bool{}
	validationSpawnCount := 0
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ev codexCollabItem
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		it := ev.Item
		// spawn_agent surfaces as item.started (no threads yet) AND item.completed
		// (with the spawned receiver_thread_ids); count the completed one so the
		// thread binding and the count stay in lockstep (one spawn = one count).
		if it.Type != "collab_tool_call" || it.Tool != "spawn_agent" || ev.Type != "item.completed" {
			continue
		}
		if codexDispatchesValidation(it.Prompt) {
			validationSpawnCount++
			for _, tid := range it.ReceiverThreadIDs {
				validationThreads[tid] = true
			}
		}
	}

	if validationSpawnCount == 0 {
		return reviewerUnsupported, "no validation spawn_agent — no structured thread handle to correlate"
	}
	if validationSpawnCount >= 2 {
		return reviewerFresh, fmt.Sprintf("%d distinct validation spawn_agents — fresh reviewers", validationSpawnCount)
	}
	if len(validationThreads) == 0 {
		return reviewerUnsupported, "the validation spawn_agent exposed no receiver_thread_ids handle"
	}

	reuseHit := false
	wrongRoute := false
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ev codexCollabItem
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		it := ev.Item
		if it.Type != "collab_tool_call" || !codexReviewerReuseTool(it.Tool) {
			continue
		}
		hitValidation := false
		for _, tid := range it.ReceiverThreadIDs {
			if validationThreads[tid] {
				hitValidation = true
			}
		}
		if hitValidation {
			reuseHit = true
		} else if len(it.ReceiverThreadIDs) > 0 {
			wrongRoute = true
		}
	}
	if reuseHit {
		return reviewerReuse, "followup_task/send_input to the bound validation reviewer thread"
	}
	if wrongRoute {
		return reviewerNoReuse, "the cycle-2 re-review was routed to a non-validation thread while a validation reviewer existed"
	}
	return reviewerNoReuse, "one validation reviewer bound but no followup_task/send_input reused its thread"
}

// assertCodexReviewerReuse enforces the rejection-flow reviewer routing over the
// `codex exec --json` transcript. Structured reuse (followup_task/send_input to the
// bound validation thread) and structurally-distinct fresh reviewers both pass — the
// latter is the contract-correct choice when the host exposes no addressable-worker
// route. A re-review routed to a non-validation thread fails; an absent structured
// handle yields errReviewerIdentityUnsupported (never a reuse pass).
func assertCodexReviewerReuse(jsonl string) error {
	result, detail := codexReviewerIdentity(jsonl)
	switch result {
	case reviewerReuse, reviewerFresh:
		return nil
	case reviewerNoReuse:
		return fmt.Errorf("the FO did not route the cycle-2 re-review to a validation reviewer: %s", detail)
	default:
		return fmt.Errorf("%w: %s", errReviewerIdentityUnsupported, detail)
	}
}

func codexReviewerReuseTool(tool string) bool {
	return tool == "followup_task" || tool == "send_input"
}

// ---------------------------------------------------------------------------
// Rejection-flow worker topology (native transcripts)
// ---------------------------------------------------------------------------

// rejectionRoute is one ordered routing observation in a rejection-flow run: the FO
// OPENED a worker for a stage, ROUTED follow-up work to a worker it had already
// opened, or a dispatched worker reported DONE. Identity is the host's own
// structured handle — a Codex task path, a Claude teammate name — and never prompt
// content: a Codex `spawn_agent`'s arguments are an encrypted blob EXCEPT the
// plaintext task path, so the path is the only identity the rollout exposes at all.
type rejectionRoute struct {
	index  int
	event  string
	stage  string
	target string
}

const (
	routeSpawn = "spawn"
	routeReuse = "reuse"
	routeDone  = "done"
)

// rejectionBranch is the contract-observable that decides which ordered chain a run
// owes. It is NOT a choice between two acceptable behaviors: on any given run exactly
// one branch is conforming, fixed by whether the reuse route survives reuse
// condition 0 (`fo-dispatch-core.md:49`, fail-safe — "if it reports the worker over
// budget, or the probe is unavailable, dispatch fresh").
type rejectionBranch string

const (
	rejectionBranchReuse rejectionBranch = "reuse"
	rejectionBranchFresh rejectionBranch = "fresh"
)

// rejectionStageOfHandle reads the stage out of a worker handle. Both hosts derive
// the handle from (slug, stage) — `spacedock_ensign_rejection_task_validation`,
// `spacedock-ensign-rejection-task-validation` — so the trailing stage token is a
// structural property of the handle, not prose the model chose.
func rejectionStageOfHandle(handle string) string {
	switch {
	case strings.HasSuffix(handle, "validation"):
		return "validation"
	case strings.HasSuffix(handle, "implementation"):
		return "implementation"
	}
	return ""
}

// codexRejectionRoutes extracts the ordered worker topology from a Codex NATIVE
// rollout. The public `codex exec --json` stream cannot serve: its only
// `collab_tool_call` items are `wait` (verified across the preserved streams and
// both spike runs), so it carries no topology at all. The rollout does, as three
// correlated payloads: a `spawn_agent` `function_call`, the `function_call_output`
// that returns the spawned worker's task path, and later `followup_task` calls that
// name a task path directly. Worker completion is the `agent_message` whose author
// is that task path and whose content carries the ensign's `Done:` signal.
func codexRejectionRoutes(rollout string) []rejectionRoute {
	var routes []rejectionRoute
	pendingSpawn := map[string]bool{} // call_id of a spawn_agent awaiting its output
	// awaiting holds the task paths currently owed a completion. A dispatch (spawn or
	// follow-up) opens the debt and the round's FIRST `Done:` closes it, so a worker
	// that narrates "Done:" more than once in one round still contributes one event —
	// the same first-transition ordering semantic the shared lifecycle helper uses.
	awaiting := map[string]bool{}
	for i, line := range strings.Split(rollout, "\n") {
		var event struct {
			Payload struct {
				Type      string          `json:"type"`
				Name      string          `json:"name"`
				CallID    string          `json:"call_id"`
				Arguments string          `json:"arguments"`
				Output    string          `json:"output"`
				Author    string          `json:"author"`
				Content   json.RawMessage `json:"content"`
			} `json:"payload"`
		}
		if json.Unmarshal([]byte(line), &event) != nil {
			continue
		}
		p := event.Payload
		switch {
		case p.Type == "function_call" && p.Name == "spawn_agent":
			pendingSpawn[p.CallID] = true
		case p.Type == "function_call_output" && pendingSpawn[p.CallID]:
			delete(pendingSpawn, p.CallID)
			if handle := codexHandleFromJSON(p.Output); handle != "" {
				awaiting[handle] = true
				routes = append(routes, rejectionRoute{index: i, event: routeSpawn, stage: rejectionStageOfHandle(handle), target: handle})
			}
		case p.Type == "function_call" && p.Name == "followup_task":
			if handle := codexHandleFromJSON(p.Arguments); handle != "" {
				awaiting[handle] = true
				routes = append(routes, rejectionRoute{index: i, event: routeReuse, stage: rejectionStageOfHandle(handle), target: handle})
			}
		// A worker's completion is its FINAL_ANSWER carrying the ensign `Done:`
		// signal. Requiring both markers skips the intermediate `Message Type:
		// MESSAGE` traffic a live worker also emits, which a bare `Done:` search
		// would miscount as a second completion for the same round.
		case p.Type == "agent_message" && awaiting[codexBareHandle(p.Author)] &&
			strings.Contains(string(p.Content), "FINAL_ANSWER") && strings.Contains(string(p.Content), "Done:"):
			handle := codexBareHandle(p.Author)
			delete(awaiting, handle)
			routes = append(routes, rejectionRoute{index: i, event: routeDone, stage: rejectionStageOfHandle(handle), target: handle})
		}
	}
	return routes
}

// codexHandleFromJSON reads a worker identity out of a rollout JSON payload, and
// returns "" when the payload names no worker. Both facts below are read off real
// multi-agent rollout bytes rather than assumed, and both matter:
//
//   - a spawn's `function_call_output` returns `{"task_name":"/root/NAME",…}`, while a
//     `followup_task`'s arguments carry the identity under a DIFFERENT key, `target`,
//     alongside an encrypted `message`. Reading only `task_name` silently drops every
//     follow-up and makes a conforming reuse chain look two events short.
//   - a REFUSED spawn returns a plain error string ("agent path `/root/NAME` already
//     exists"), not JSON. That spawn opened no worker, so it must contribute no
//     routing event; treating the error text as a handle invents one.
//
// Reading `task_name`/`target` and never `message` is also what keeps identity on
// task paths rather than prompt content — not a stylistic choice here, since the
// prompt is an encrypted blob.
func codexHandleFromJSON(raw string) string {
	var doc struct {
		TaskName string `json:"task_name"`
		Target   string `json:"target"`
	}
	if json.Unmarshal([]byte(raw), &doc) != nil {
		return ""
	}
	if doc.TaskName != "" {
		return codexBareHandle(doc.TaskName)
	}
	return codexBareHandle(doc.Target)
}

// codexBareHandle strips the `/root/` rooting the rollout applies inconsistently —
// spawn outputs and `agent_message.author` are rooted, a `followup_task` target may
// be either — so correlating a reuse back to its spawn compares like with like.
func codexBareHandle(handle string) string {
	if cut := strings.LastIndex(handle, "/"); cut >= 0 {
		return handle[cut+1:]
	}
	return handle
}

// claudeRejectionRoutes extracts the ordered worker topology AND the run's branch
// key from a Claude stream-json transcript. The branch key is the `dispatch
// context-budget` probe result read off the probe's own tool_result: the command
// prints `reuse_ok` on success and, on the fail-safe path, exits non-zero with no
// `reuse_ok` on stdout at all (internal/dispatch contextbudget parity contract). A
// run where no probe ever reported `reuse_ok` owes the fail-safe FRESH chain; one
// where a probe did owes the REUSE chain.
func claudeRejectionRoutes(stream string) ([]rejectionRoute, rejectionBranch) {
	var routes []rejectionRoute
	// A reused worker is the SAME background task, so its completion notification
	// carries the tool_use id of the Agent call that opened it. Completions are
	// therefore tracked as a debt against that id, which a dispatch (spawn or reuse
	// advance) arms and the round's notification clears — one completion per round,
	// and a reused worker can report twice without a second spawn.
	openedBy := map[string]string{}       // teammate name -> the tool_use id that opened it
	awaitingStage := map[string]string{}  // armed tool_use id -> stage awaited
	awaitingTarget := map[string]string{} // armed tool_use id -> teammate name
	probeIDs := map[string]bool{}         // Bash tool_use id of a context-budget probe
	live := map[string]bool{}             // teammate names this run has opened
	branch := rejectionBranchFresh
	for i, line := range strings.Split(stream, "\n") {
		var event struct {
			Type      string `json:"type"`
			Subtype   string `json:"subtype"`
			Status    string `json:"status"`
			ToolUseID string `json:"tool_use_id"`
			Message   *struct {
				Content []struct {
					Type      string `json:"type"`
					Name      string `json:"name"`
					ID        string `json:"id"`
					ToolUseID string `json:"tool_use_id"`
					Input     struct {
						Description string          `json:"description"`
						Name        string          `json:"name"`
						Command     string          `json:"command"`
						To          string          `json:"to"`
						Message     json.RawMessage `json:"message"`
					} `json:"input"`
					Content json.RawMessage `json:"content"`
				} `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal([]byte(line), &event) != nil {
			continue
		}
		if event.Type == "system" && event.Subtype == "task_notification" && event.Status == "completed" {
			if stage, ok := awaitingStage[event.ToolUseID]; ok {
				routes = append(routes, rejectionRoute{index: i, event: routeDone, stage: stage, target: awaitingTarget[event.ToolUseID]})
				delete(awaitingStage, event.ToolUseID)
			}
		}
		if event.Message == nil {
			continue
		}
		for _, block := range event.Message.Content {
			switch {
			case block.Type == "tool_use" && (block.Name == "Agent" || block.Name == "Task"):
				stage := rejectionStageOfHandle(block.Input.Name)
				if stage == "" {
					stage = rejectionStageOfDescription(block.Input.Description)
				}
				openedBy[block.Input.Name] = block.ID
				awaitingStage[block.ID], awaitingTarget[block.ID] = stage, block.Input.Name
				live[block.Input.Name] = true
				routes = append(routes, rejectionRoute{index: i, event: routeSpawn, stage: stage, target: block.Input.Name})
			case block.Type == "tool_use" && block.Name == "Bash" && strings.Contains(block.Input.Command, "dispatch context-budget"):
				probeIDs[block.ID] = true
			// A shutdown_request is the contract's supersede teardown, not a reuse
			// advance, so it never counts as routing follow-up work.
			case block.Type == "tool_use" && block.Name == "SendMessage" && live[block.Input.To] && !isShutdownRequest(block.Input.Message):
				stage := rejectionStageOfHandle(block.Input.To)
				if id := openedBy[block.Input.To]; id != "" {
					awaitingStage[id], awaitingTarget[id] = stage, block.Input.To
				}
				routes = append(routes, rejectionRoute{index: i, event: routeReuse, stage: stage, target: block.Input.To})
			case block.Type == "tool_result" && probeIDs[block.ToolUseID] && strings.Contains(string(block.Content), "reuse_ok"):
				branch = rejectionBranchReuse
			}
		}
	}
	return routes, branch
}

// rejectionStageOfDescription is the fallback stage read for a Claude spawn whose
// `name` is absent (a bare `Task`); the dispatch description names the stage.
func rejectionStageOfDescription(description string) string {
	lower := strings.ToLower(description)
	switch {
	case strings.Contains(lower, "validation"):
		return "validation"
	case strings.Contains(lower, "implementation"):
		return "implementation"
	}
	return ""
}

// codexDispatchesValidation reports whether a spawn_agent prompt dispatched the
// validation stage — either the validation dispatch file
// (spacedock-ensign-{slug}-validation.md) or a prompt naming the validation stage.
func codexDispatchesValidation(prompt string) bool {
	lower := strings.ToLower(prompt)
	return strings.Contains(lower, "-validation.md") || strings.Contains(lower, "validation")
}
