package ensigncycle

import (
	"errors"
	"strings"
	"testing"
)

// TestExtractClaudeFinalMessage is the AC-7 offline guard for the Claude runner's
// `observed` extractor — the front-door analog of Codex `--output-last-message`.
// It is pure string parsing over a synthetic stream-json transcript, so it spends
// no live credential and gives the riskiest sharing mechanism a fast guard.
func TestExtractClaudeFinalMessage(t *testing.T) {
	t.Run("prefers_result_success_event", func(t *testing.T) {
		// A normal stream: assistant text blocks, then the terminal
		// result/success event carrying the FO's final message in `result`. The
		// extractor must prefer the result event over the assistant text.
		stream := strings.Join([]string{
			`{"type":"assistant","message":{"content":[{"type":"text","text":"Looking at the workflow."}]}}`,
			`{"type":"assistant","message":{"content":[{"type":"text","text":"This is an intermediate block, not the final message."}]}}`,
			`{"type":"result","subtype":"success","is_error":false,"result":"Gate review: Gate Check - review\nDecision: approve or reject?"}`,
		}, "\n")

		got, err := extractClaudeFinalMessage(stream)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := "Gate review: Gate Check - review\nDecision: approve or reject?"
		if got != want {
			t.Fatalf("final message = %q, want the result event's result text %q", got, want)
		}
	})

	t.Run("falls_back_to_last_assistant_text_block", func(t *testing.T) {
		// No result event at all (e.g. a transcript truncated before the terminal
		// event). The extractor falls back to the LAST assistant text block, which
		// the spike pinned as the secondary source.
		stream := strings.Join([]string{
			`{"type":"assistant","message":{"content":[{"type":"text","text":"First assistant block."}]}}`,
			`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Bash","input":{}}]}}`,
			`{"type":"assistant","message":{"content":[{"type":"text","text":"Gate review: held.\nDecision: pending."}]}}`,
		}, "\n")

		got, err := extractClaudeFinalMessage(stream)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := "Gate review: held.\nDecision: pending."
		if got != want {
			t.Fatalf("fallback final message = %q, want last assistant text block %q", got, want)
		}
	})

	t.Run("surfaces_401_result_as_launch_failure", func(t *testing.T) {
		// The machine gotcha from the spike: a stale OAuth benchmark-token yields a
		// result event with is_error:true and api_error_status:401 BEFORE any FO
		// work. This MUST be surfaced as a launch failure (a distinct error), NOT
		// returned as `observed` text that a scenario assertion would then judge.
		stream := strings.Join([]string{
			`{"type":"system","subtype":"init"}`,
			`{"type":"result","subtype":"error_during_execution","is_error":true,"api_error_status":401,"result":"API Error: 401 OAuth token has expired"}`,
		}, "\n")

		got, err := extractClaudeFinalMessage(stream)
		if err == nil {
			t.Fatalf("expected a launch failure for a 401 result event, got message %q", got)
		}
		if !errors.Is(err, errClaudeLaunchFailed) {
			t.Fatalf("error = %v, want it to wrap errClaudeLaunchFailed so the runner fails loudly distinct from an assertion failure", err)
		}
		// The launch-failure error must carry the 401 signal so the loud failure
		// names the real cause, not a generic miss.
		if !strings.Contains(err.Error(), "401") {
			t.Fatalf("launch-failure error = %q, want it to name the 401 status", err.Error())
		}
	})

	t.Run("error_result_without_401_is_still_launch_failure", func(t *testing.T) {
		// Any is_error:true result event means the launch itself failed (no FO
		// work happened), so it is a launch failure rather than `observed` text.
		stream := `{"type":"result","subtype":"error_during_execution","is_error":true,"result":"some other launch error"}`

		_, err := extractClaudeFinalMessage(stream)
		if !errors.Is(err, errClaudeLaunchFailed) {
			t.Fatalf("error = %v, want errClaudeLaunchFailed for any is_error result", err)
		}
	})

	t.Run("no_result_and_no_assistant_text_is_an_error", func(t *testing.T) {
		// A transcript with neither a result event nor any assistant text block
		// cannot yield an `observed` string; the extractor errors rather than
		// returning empty text a scenario assertion would silently judge.
		stream := `{"type":"system","subtype":"init"}`

		if _, err := extractClaudeFinalMessage(stream); err == nil {
			t.Fatal("expected an error when no result event and no assistant text exist")
		}
	})
}
