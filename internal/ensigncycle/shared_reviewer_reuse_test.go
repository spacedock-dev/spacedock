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

// assertClaudeSingleEntityRejectionFlow is the single-entity (`-p`) Claude
// producer-signal assertion for the rejection-flow scenario. The Claude runner
// launches `spacedock claude -- -p {prompt}` with a prompt naming one entity. The
// `-p` FO drives EITHER bare OR team mode (it opts into the background inter-agent
// communication when SendMessage is exposed), so the contract admits two valid end-states: bare
// fresh-dispatches a distinct reviewer per cycle (reviewerFresh); team keeps the
// cycle-1 reviewer alive and reuses it (reviewerReuse). The invariant across both is
// that the cycle-2 re-review reaches a VALIDATION worker and NEVER collapses onto the
// implementation worker (reviewerNoReuse) — a shutdown_request to the superseded
// fix worker is exempt teardown, not a re-review routing.
func assertClaudeSingleEntityRejectionFlow(stream string) error {
	result, detail := claudeReviewerIdentity(stream)
	switch result {
	case reviewerReuse, reviewerFresh:
		return nil
	case reviewerNoReuse:
		return fmt.Errorf("the single-entity cycle-2 re-review did not reach a validation worker (fresh or reused): %s", detail)
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

// codexDispatchesValidation reports whether a spawn_agent prompt dispatched the
// validation stage — either the validation dispatch file
// (spacedock-ensign-{slug}-validation.md) or a prompt naming the validation stage.
func codexDispatchesValidation(prompt string) bool {
	lower := strings.ToLower(prompt)
	return strings.Contains(lower, "-validation.md") || strings.Contains(lower, "validation")
}
