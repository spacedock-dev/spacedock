package ensigncycle

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// The host-specific reviewer-reuse assertions for the rejection-flow scenario.
// They are the producer signal that the FO reused the kept-alive validation
// reviewer for the cycle-2 re-review rather than fresh-dispatching (the #141
// keepalive contract the Go port dropped): Claude reuses via a SendMessage tool
// call, Codex multi_agent_v2 via `followup_task`, and legacy Codex fixtures via
// `send_input`. They live under the DEFAULT build tags
// (parsing only stdlib JSON) so the offline table tests exercise them without
// spending a model, alongside the //go:build live runners that feed them real
// transcripts.

// claudeAgentIDResult extracts the `agentId: <id>` a completed subagent's
// tool_result returns ("agentId: a94abe89c85f9f4cc (use SendMessage with to: …)").
// The teams-mode reuse handle for a kept-alive completed agent.
var claudeAgentIDResult = regexp.MustCompile(`agentId:\s*(a[0-9a-f]+)`)

// assertClaudeReviewerReuse scans the stream-json transcript for the FO reusing
// the kept-alive cycle-1 validation reviewer for the cycle-2 re-review — the
// durable producer signal it did NOT fresh-dispatch the cycle-2 validator (the
// #141 keepalive contract). It enforces BOTH halves of "genuine reuse", because
// either alone false-passes a fresh-dispatch (the cycle-8 audit hole):
//
//  1. NO fresh cycle-2 validation spawn. The rejection-flow drives ONE cycle-1
//     validation; a reuse run re-engages that same reviewer, so there is EXACTLY
//     ONE validation-stage Agent/Task spawn. A run that fresh-dispatches the
//     cycle-2 validator emits a SECOND validation spawn — the forbidden behavior —
//     so >1 validation spawn is an immediate FAIL regardless of any reuse-shaped
//     message. This is the discriminator the bare name/agentId match lacked: a
//     fresh cycle-2 reviewer messaged by its validation NAME looks identical to a
//     reused one EXCEPT for the extra spawn.
//  2. A message to the cycle-1 reviewer's handle. With exactly one validation
//     reviewer in play, a SendMessage to it — by the agentId its spawn returned
//     (teams-mode kept-alive resume) OR by its spawn NAME (`…-validation`, the
//     production shape the recorded testdata uses) — is the re-review. Narration
//     or a message to a non-validation recipient does not count.
func assertClaudeReviewerReuse(stream string) error {
	// First pass: count validation-stage Agent/Task spawns and collect the agentIds
	// those spawns returned. A reused reviewer drives one validation spawn; a fresh
	// cycle-2 dispatch drives a second.
	validationSpawnIDs := map[string]bool{} // tool_use_id of each validation-stage spawn
	validationSpawnCount := 0
	validationAgentIDs := map[string]bool{} // agentIds those spawns returned (reuse handles)
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
			}
			if block.Type == "tool_result" && validationSpawnIDs[block.ToolUseID] {
				if m := claudeAgentIDResult.FindStringSubmatch(string(block.Content)); m != nil {
					validationAgentIDs[m[1]] = true
				}
			}
		}
	}

	if validationSpawnCount == 0 {
		return fmt.Errorf("no validation-stage Agent/Task spawn found — the FO never created a cycle-1 reviewer to reuse")
	}
	if validationSpawnCount > 1 {
		return fmt.Errorf("the FO emitted %d validation-stage Agent/Task spawns — it FRESH-dispatched the cycle-2 validator instead of reusing the kept-alive cycle-1 reviewer (the #141 keepalive contract violation)", validationSpawnCount)
	}

	// Second pass: with exactly one validation reviewer, a SendMessage to it (by its
	// returned agentId, or by a name naming the validation stage) is the cycle-2
	// re-review reuse signal.
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
						To string `json:"to"`
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
			to := block.Input.To
			if strings.Contains(strings.ToLower(to), "validation") || validationAgentIDs[to] {
				return nil
			}
		}
	}
	return fmt.Errorf("the FO spawned exactly one validation reviewer but sent it no reuse SendMessage (by name or agentId) for the cycle-2 re-review")
}

// assertClaudeSingleEntityRejectionFlow is the single-entity (`-p`) Claude
// producer-signal assertion for the rejection-flow scenario. The Claude runner
// launches `spacedock claude -- -p {prompt}` with a prompt naming one entity, so
// the run is single-entity → bare (claude-fo-dispatch.md: "In single-entity mode,
// skip team creation. Use bare-mode dispatch for all agent spawning"). In bare mode
// the contract makes the feedback flow DETERMINISTIC and SEQUENTIAL
// (claude-fo-dispatch.md `## Feedback Rejection Flow (bare mode)`: "dispatch fix
// agent (wait for completion), then dispatch reviewer (wait for completion)"). So
// the contract-correct end-state is: the cycle-2 re-review is a DISTINCT, FRESHLY
// DISPATCHED validation worker — NOT a reuse of the bare cycle-1 reviewer (reuse-
// condition-1 hard-fails in bare mode), and NOT the implementation worker serving as
// its own validator (the fix agent and the reviewer are separate sequential
// dispatches). It enforces BOTH halves, because either alone false-passes a wrong
// run:
//
//  1. AT LEAST TWO distinct validation-stage Agent/Task spawns. The bare flow
//     fresh-dispatches a validation reviewer for cycle-1 AND a fresh one for cycle-2
//     (no reuse handle exists). A run with fewer than two validation spawns either
//     never re-reviewed or collapsed the cycle-2 re-review onto a non-validation
//     worker — both forbidden. This is the discriminator that catches the observed
//     non-deterministic "reused the impl ensign through validation" run (which left
//     only the cycle-1 validation spawn).
//  2. The cycle-2 re-review is NOT routed to an implementation worker. A SendMessage
//     to an `…-implementation` handle telling it to validate is the impl-as-validator
//     violation; the re-review must be a validation-stage spawn, not a message to the
//     fix worker.
//
// This is the CONTRACT-correct single-entity assertion. It replaces the team-mode
// assertClaudeReviewerReuse for this scenario's `-p` run, which wrongly assumed a
// kept-alive reviewer the bare contract cannot produce (the AC-3 finding). The real
// reviewer-continuity question — whether single-entity SHOULD create a team so the
// reviewer is reusable — is the spun-off option-(a) task, not this correction.
func assertClaudeSingleEntityRejectionFlow(stream string) error {
	validationSpawnCount := 0
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
						Description string `json:"description"`
						To          string `json:"to"`
					} `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err != nil || entry.Message == nil {
			continue
		}
		for _, block := range entry.Message.Content {
			if block.Type != "tool_use" {
				continue
			}
			desc := strings.ToLower(block.Input.Description)
			if (block.Name == "Agent" || block.Name == "Task") && strings.Contains(desc, "validation") {
				validationSpawnCount++
			}
			// The impl-as-validator violation: a SendMessage telling an
			// implementation-named worker to validate / re-review.
			if block.Name == "SendMessage" && strings.Contains(strings.ToLower(block.Input.To), "implementation") {
				return fmt.Errorf("the cycle-2 re-review was routed to an implementation worker (%q) — the fix agent and the reviewer must be SEPARATE sequential dispatches in bare mode; the impl worker must never serve as its own validator", block.Input.To)
			}
		}
	}
	if validationSpawnCount < 2 {
		return fmt.Errorf("single-entity bare rejection-flow produced %d validation-stage spawns, want >= 2 (a fresh cycle-1 reviewer AND a fresh cycle-2 reviewer — bare mode cannot reuse a kept-alive reviewer, so each cycle fresh-dispatches a distinct validation worker)", validationSpawnCount)
	}
	return nil
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

// assertCodexReviewerReuse scans the `codex exec --json` transcript for the FO's
// rejection-flow reviewer routing. When the transcript proves a turn-starting
// `«addressable-worker»` route exists, the FO must reuse the kept-alive cycle-1
// validation reviewer for the cycle-2 re-review — the durable producer signal it
// did NOT fresh-dispatch the cycle-2 validator (the #141 keepalive contract).
// Current public Codex live surfaces may expose only spawn/wait; when the FO
// explicitly observes that no turn-starting reuse route is exposed, the contract is
// characterized as addressable-worker ABSENT and the cycle-2 reviewer is fresh.
// Codex multi_agent_v2 reuses via a `followup_task` collab_tool_call to the
// reviewer's thread; legacy pre-v2 fixtures use `send_input`. The thread is bound
// by the `spawn_agent` whose prompt dispatched the validation stage. In the PRESENT
// branch it enforces BOTH halves of "genuine reuse", because binding ANY
// validation-prompt spawn alone false-passes a fresh-dispatch (the cycle-8 M2 hole):
//
//  1. NO fresh cycle-2 validation spawn. The rejection-flow drives ONE cycle-1
//     validation spawn_agent; a reuse run re-engages that same thread. A run that
//     fresh-spawns the cycle-2 validator emits a SECOND validation spawn_agent — the
//     forbidden behavior — so >1 validation spawn_agent is an immediate FAIL, even
//     though a follow-up to the fresh thread otherwise looks like reuse.
//  2. A turn-triggering reuse call to the cycle-1 reviewer's bound thread. With
//     exactly one validation spawn, a follow-up to its thread is the re-review —
//     not feedback-to-implementation routing (which targets the impl thread), not
//     loose narration.
func assertCodexReviewerReuse(jsonl string) error {
	if codexAddressableWorkerAbsent(jsonl) {
		return assertCodexFreshValidationWhenAddressableAbsent(jsonl)
	}

	// First pass: count validation spawn_agents and bind their thread id(s). A
	// reused reviewer drives one validation spawn; a fresh cycle-2 dispatch drives a
	// second.
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
		return fmt.Errorf("no validation spawn_agent found — the FO never created a cycle-1 reviewer to reuse")
	}
	if validationSpawnCount > 1 {
		return fmt.Errorf("the FO emitted %d validation spawn_agents — it FRESH-dispatched the cycle-2 validator instead of reusing the kept-alive cycle-1 reviewer (the #141 keepalive contract violation)", validationSpawnCount)
	}

	// Second pass: with exactly one validation reviewer, a turn-triggering reuse
	// call to its bound thread is the cycle-2 re-review reuse signal.
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
		for _, tid := range it.ReceiverThreadIDs {
			if validationThreads[tid] {
				return nil
			}
		}
	}
	return fmt.Errorf("the FO spawned exactly one validation reviewer but sent it no followup_task/send_input for the cycle-2 re-review")
}

func codexReviewerReuseTool(tool string) bool {
	return tool == "followup_task" || tool == "send_input"
}

func codexAddressableWorkerAbsent(jsonl string) bool {
	lower := strings.ToLower(jsonl)
	return strings.Contains(lower, "no followup_task/send_message reuse route exposed") ||
		strings.Contains(lower, "followup_task/message reuse is unavailable") ||
		strings.Contains(lower, "no turn-starting follow-up route for a completed worker") ||
		strings.Contains(lower, "no completed-worker follow-up route") ||
		strings.Contains(lower, "reviewer reuse is not available in this host")
}

func assertCodexFreshValidationWhenAddressableAbsent(jsonl string) error {
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
		if it.Type != "collab_tool_call" {
			continue
		}
		if codexReviewerReuseTool(it.Tool) {
			return fmt.Errorf("Codex addressable-worker was characterized ABSENT, but transcript contains turn-starting reuse tool %q", it.Tool)
		}
		if it.Tool == "spawn_agent" && ev.Type == "item.completed" && codexDispatchesValidation(it.Prompt) {
			validationSpawnCount++
		}
	}
	if validationSpawnCount < 2 {
		return fmt.Errorf("Codex addressable-worker ABSENT rejection-flow produced %d validation spawn_agents, want >= 2 fresh validation reviewers for cycle 1 and cycle 2", validationSpawnCount)
	}
	return nil
}

// codexDispatchesValidation reports whether a spawn_agent prompt dispatched the
// validation stage — either the validation dispatch file
// (spacedock-ensign-{slug}-validation.md) or a prompt naming the validation stage.
func codexDispatchesValidation(prompt string) bool {
	lower := strings.ToLower(prompt)
	return strings.Contains(lower, "-validation.md") || strings.Contains(lower, "validation")
}
