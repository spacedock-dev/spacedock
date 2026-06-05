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
// call, Codex via a `send_input` call. They live under the DEFAULT build tags
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

// codexCollabItem is one `codex exec --json` stream item. Codex surfaces its
// multi-agent calls as `collab_tool_call` items (tool = spawn_agent / send_input /
// wait / close_agent); the worker is addressed by opaque `receiver_thread_ids`,
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

// assertCodexReviewerReuse scans the `codex exec --json` transcript for the FO
// reusing the kept-alive cycle-1 validation reviewer for the cycle-2 re-review —
// the durable producer signal it did NOT fresh-dispatch the cycle-2 validator (the
// #141 keepalive contract). Codex reuses via a `send_input` collab_tool_call to the
// reviewer's thread; the thread is bound by the `spawn_agent` whose prompt
// dispatched the validation stage. It enforces BOTH halves of "genuine reuse",
// because binding ANY validation-prompt spawn alone false-passes a fresh-dispatch
// (the cycle-8 M2 hole):
//
//  1. NO fresh cycle-2 validation spawn. The rejection-flow drives ONE cycle-1
//     validation spawn_agent; a reuse run re-engages that same thread. A run that
//     fresh-spawns the cycle-2 validator emits a SECOND validation spawn_agent — the
//     forbidden behavior — so >1 validation spawn_agent is an immediate FAIL, even
//     though a send_input to the fresh thread otherwise looks like reuse.
//  2. A send_input to the cycle-1 reviewer's bound thread. With exactly one
//     validation spawn, a send_input to its thread is the re-review — not the
//     feedback-to-implementation send_input (which targets the impl thread), not
//     loose narration.
func assertCodexReviewerReuse(jsonl string) error {
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

	// Second pass: with exactly one validation reviewer, a send_input to its bound
	// thread is the cycle-2 re-review reuse signal.
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
		if it.Type != "collab_tool_call" || it.Tool != "send_input" {
			continue
		}
		for _, tid := range it.ReceiverThreadIDs {
			if validationThreads[tid] {
				return nil
			}
		}
	}
	return fmt.Errorf("the FO spawned exactly one validation reviewer but sent it no send_input for the cycle-2 re-review")
}

// codexDispatchesValidation reports whether a spawn_agent prompt dispatched the
// validation stage — either the validation dispatch file
// (spacedock-ensign-{slug}-validation.md) or a prompt naming the validation stage.
func codexDispatchesValidation(prompt string) bool {
	lower := strings.ToLower(prompt)
	return strings.Contains(lower, "-validation.md") || strings.Contains(lower, "validation")
}
