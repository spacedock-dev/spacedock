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
// A reused completed agent is re-engaged by this agentId, not by its spawn name.
var claudeAgentIDResult = regexp.MustCompile(`agentId:\s*(a[0-9a-f]+)`)

// assertClaudeReviewerReuse scans the stream-json transcript for the FO reusing
// the kept-alive validation reviewer for the cycle-2 re-review — the durable
// producer signal it did NOT fresh-dispatch (the #141 keepalive contract the Go
// port dropped). A reused reviewer is reachable two ways, and BOTH are genuine
// reuse: by its spawn name (a SendMessage `to` naming the validation stage) OR by
// the opaque `agentId` the validation-reviewer subagent returned on completion (a
// completed background agent in a Claude team is resumed by agentId, not name —
// the actual shape the FO emits). The agentId path correlates the SendMessage
// target back to the Agent/Task spawn whose description named the validation
// stage, so it cannot be satisfied by a SendMessage to an unrelated agent or by
// the word "validation" in narration: only a real tool_use to the validation
// reviewer's handle passes.
func assertClaudeReviewerReuse(stream string) error {
	// First pass: collect the agentIds returned by Agent/Task spawns whose
	// description named the validation stage. These are the kept-alive validation
	// reviewer's resumable handles.
	validationSpawnIDs := map[string]bool{}  // tool_use_id of a validation-stage spawn
	validationAgentIDs := map[string]bool{}  // agentId those spawns returned
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry struct {
			Message *struct {
				Content []struct {
					Type      string          `json:"type"`
					Name      string          `json:"name"`
					ID        string          `json:"id"`
					ToolUseID string          `json:"tool_use_id"`
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
				validationSpawnIDs[block.ID] = true
			}
			if block.Type == "tool_result" && validationSpawnIDs[block.ToolUseID] {
				if m := claudeAgentIDResult.FindStringSubmatch(string(block.Content)); m != nil {
					validationAgentIDs[m[1]] = true
				}
			}
		}
	}

	// Second pass: a SendMessage `to` that names the validation stage OR equals a
	// validation-reviewer agentId is the reuse signal.
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
	return fmt.Errorf("no SendMessage tool_use targeting the validation reviewer (by name or by the kept-alive reviewer's agentId) found in the stream — the FO did not reuse the kept-alive reviewer for the cycle-2 re-review")
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
// reusing the kept-alive validation reviewer for the cycle-2 re-review — the
// durable producer signal it did NOT fresh-dispatch (the #141 keepalive contract
// the Go port dropped). Codex reuses via a `send_input` collab_tool_call to the
// reviewer's thread. The validation reviewer's thread cannot be matched by a
// "validation" name (Codex addresses threads by opaque id), so the assertion
// correlates: the spawn_agent whose prompt dispatched the validation stage binds
// the reviewer's thread id, and a later send_input to THAT thread is the reuse.
// This cannot be satisfied by the OTHER send_input the flow emits — the
// feedback-to-implementation routing — because that targets the implementation
// worker's thread, not the validation reviewer's; nor by loose narration, since it
// matches the collab_tool_call shape, not prose.
func assertCodexReviewerReuse(jsonl string) error {
	// First pass: bind the validation reviewer's thread id(s) from the spawn_agent
	// whose prompt dispatched the validation stage.
	validationThreads := map[string]bool{}
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
		if it.Type != "collab_tool_call" || it.Tool != "spawn_agent" {
			continue
		}
		if codexDispatchesValidation(it.Prompt) {
			for _, tid := range it.ReceiverThreadIDs {
				validationThreads[tid] = true
			}
		}
	}

	// Second pass: a send_input to a validation reviewer thread is the reuse signal.
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
	return fmt.Errorf("no send_input collab_tool_call to the kept-alive validation reviewer's thread found in the transcript — the FO did not reuse the reviewer for the cycle-2 re-review")
}

// codexDispatchesValidation reports whether a spawn_agent prompt dispatched the
// validation stage — either the validation dispatch file
// (spacedock-ensign-{slug}-validation.md) or a prompt naming the validation stage.
func codexDispatchesValidation(prompt string) bool {
	lower := strings.ToLower(prompt)
	return strings.Contains(lower, "-validation.md") || strings.Contains(lower, "validation")
}
