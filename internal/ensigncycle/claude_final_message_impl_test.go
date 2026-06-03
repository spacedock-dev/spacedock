package ensigncycle

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// errClaudeLaunchFailed marks a Claude launch that never reached FO work — an
// is_error result event (the canonical case is a stale OAuth benchmark-token's
// api_error_status:401). The Claude runner must surface this as a LOUD launch
// failure, distinct from a scenario-assertion failure, so a credential problem is
// never fed into a behavior assertion and misread as a runtime regression.
var errClaudeLaunchFailed = errors.New("claude launch failed before any FO work")

// claudeResultEvent is the terminal stream-json event of a `spacedock claude`
// run: the front-door analog of Codex `--output-last-message`. Only the fields
// the extractor reads are modeled. It is parsed independently of streamEntry
// (which the watcher owns and does not parse the result event) so the extractor
// stays self-contained.
type claudeResultEvent struct {
	Type           string `json:"type"`
	Subtype        string `json:"subtype"`
	IsError        bool   `json:"is_error"`
	APIErrorStatus int    `json:"api_error_status"`
	Result         string `json:"result"`
}

// claudeAssistantEvent is an assistant stream entry; the extractor reads its text
// blocks for the fallback path when no result event is present.
type claudeAssistantEvent struct {
	Type    string `json:"type"`
	Message struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"message"`
}

// extractClaudeFinalMessage returns the Claude run's final-message string from a
// stream-json transcript — the `observed` source the shared scenario assertions
// consume. The spike pinned the precedence:
//
//	(1) a terminal result/success event's `result` field is the front-door analog
//	    of Codex --output-last-message and is preferred;
//	(2) absent any result event, the LAST assistant text block is the fallback;
//	(3) an is_error result event (e.g. a stale-token api_error_status:401) is a
//	    LAUNCH failure (errClaudeLaunchFailed), never `observed` text, so a
//	    credential problem fails loudly rather than being judged as behavior.
//
// A transcript with no result event and no assistant text yields an error rather
// than an empty string a scenario assertion would silently accept.
func extractClaudeFinalMessage(stream string) (string, error) {
	var lastAssistantText string
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var result claudeResultEvent
		if err := json.Unmarshal([]byte(line), &result); err == nil && result.Type == "result" {
			if result.IsError {
				if result.APIErrorStatus != 0 {
					return "", fmt.Errorf("%w: result event reported is_error with api_error_status %d: %s",
						errClaudeLaunchFailed, result.APIErrorStatus, result.Result)
				}
				return "", fmt.Errorf("%w: result event reported is_error: %s", errClaudeLaunchFailed, result.Result)
			}
			return result.Result, nil
		}

		var assistant claudeAssistantEvent
		if err := json.Unmarshal([]byte(line), &assistant); err == nil && assistant.Type == "assistant" {
			for _, block := range assistant.Message.Content {
				if block.Type == "text" && block.Text != "" {
					lastAssistantText = block.Text
				}
			}
		}
	}

	if lastAssistantText != "" {
		return lastAssistantText, nil
	}
	return "", fmt.Errorf("no result/success event and no assistant text block in the Claude stream transcript")
}
