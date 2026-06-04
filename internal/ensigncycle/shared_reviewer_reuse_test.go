package ensigncycle

import (
	"encoding/json"
	"fmt"
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

// assertClaudeReviewerReuse scans the stream-json transcript for a SendMessage
// tool_use whose `to` targets the validation reviewer — the durable producer
// signal that the FO reused the kept-alive reviewer for the cycle-2 re-review
// rather than fresh-dispatching. It parses the tool_use JSON shape (not loose
// prose), so it cannot be satisfied by the word "validation" appearing in
// narration.
func assertClaudeReviewerReuse(stream string) error {
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry struct {
			Type    string `json:"type"`
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
			if block.Type == "tool_use" && block.Name == "SendMessage" &&
				strings.Contains(strings.ToLower(block.Input.To), "validation") {
				return nil
			}
		}
	}
	return fmt.Errorf("no SendMessage tool_use targeting the validation reviewer found in the stream — the FO did not reuse the kept-alive reviewer for the cycle-2 re-review")
}

// assertCodexReviewerReuse scans the `codex exec --json` transcript for a
// `send_input` tool call whose arguments reference the validation worker — the
// durable producer signal that the FO reused the kept-alive reviewer for the
// cycle-2 re-review rather than fresh-dispatching. It matches the tool-call event's
// `name` field (the adapter's `send_input` reuse call, per
// codex-first-officer-runtime.md), not loose narration, so the FO merely writing
// the words does not satisfy it.
func assertCodexReviewerReuse(jsonl string) error {
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var event struct {
			Type      string          `json:"type"`
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		if event.Name == "send_input" && strings.Contains(strings.ToLower(string(event.Arguments)), "validation") {
			return nil
		}
	}
	return fmt.Errorf("no send_input tool call targeting the validation worker found in the transcript — the FO did not reuse the kept-alive reviewer for the cycle-2 re-review")
}
