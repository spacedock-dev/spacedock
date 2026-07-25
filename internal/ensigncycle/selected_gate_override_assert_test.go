// ABOUTME: Observed root-stream oracle for the selected Subspace gate override.
// ABOUTME: Proves one post-bind room-only handoff without probes or fallback.
package ensigncycle

import (
	"encoding/json"
	"fmt"
	"strings"
)

type selectedOverrideAction struct {
	kind  string
	name  string
	value string
	text  string
}

func selectedOverrideRootActions(stream string) []selectedOverrideAction {
	var actions []selectedOverrideAction
	for _, line := range strings.Split(stream, "\n") {
		var row struct {
			Type    string          `json:"type"`
			Parent  json.RawMessage `json:"parent_tool_use_id"`
			Message *struct {
				Content []streamContentBlock `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal([]byte(line), &row) != nil || row.Type != "assistant" || row.Message == nil ||
			(len(row.Parent) > 0 && string(row.Parent) != "null") {
			continue
		}
		for _, block := range row.Message.Content {
			switch block.Type {
			case "tool_use":
				actions = append(actions, selectedOverrideAction{
					kind:  "tool",
					name:  block.Name,
					value: inputStringField(block.Input, "command"),
					text:  inputStringField(block.Input, "skill") + "\x00" + inputStringField(block.Input, "args"),
				})
			case "text":
				actions = append(actions, selectedOverrideAction{kind: "text", text: block.Text})
			}
		}
	}
	return actions
}

func assertSelectedGateOverride(commandLog, stream, room string) error {
	countCommands := func(command string, successOnly bool) int {
		count := 0
		command = strings.TrimSpace(command)
		for _, line := range strings.Split(commandLog, "\n") {
			fields := strings.SplitN(line, "\t", 2)
			if len(fields) != 2 || !strings.HasPrefix(fields[0], "exit=") ||
				(successOnly && fields[0] != "exit=0") {
				continue
			}
			if fields[1] == command || strings.HasPrefix(fields[1], command+" ") {
				count++
			}
		}
		return count
	}
	successful := func(command string) int { return countCommands(command, true) }
	if got := countCommands("gate --help", false); got != 1 {
		return fmt.Errorf("fresh gate help attempts = %d, want 1", got)
	}
	if got := successful("gate --help"); got != 1 {
		return fmt.Errorf("successful fresh gate help calls = %d, want 1", got)
	}
	if got := successful("gate prepare recorded-gate-task "); got != 1 {
		return fmt.Errorf("successful gate prepare calls = %d, want 1", got)
	}
	if got := successful("state commit recorded-gate-task "); got != 1 {
		return fmt.Errorf("successful prepared-binding commits = %d, want 1", got)
	}
	for _, forbidden := range []string{"gate record recorded-gate-task ", "gate consume recorded-gate-task "} {
		if got := successful(forbidden); got != 0 {
			return fmt.Errorf("successful post-handoff %q calls = %d, want 0", forbidden, got)
		}
	}

	prepareAt, commitAt, overrideAt := -1, -1, -1
	overrideCount := 0
	actions := selectedOverrideRootActions(stream)
	for i, action := range actions {
		if action.kind == "tool" && action.name == "Bash" {
			lower := strings.ToLower(action.value)
			switch {
			case strings.Contains(action.value, "gate prepare recorded-gate-task"):
				if prepareAt < 0 {
					prepareAt = i
				}
			case strings.Contains(action.value, "state commit recorded-gate-task") && prepareAt >= 0:
				if commitAt < 0 {
					commitAt = i
				}
			}
			if strings.Contains(lower, "subspace") {
				return fmt.Errorf("selected Subspace override was probed or invoked through Bash: %q", action.value)
			}
			if strings.Contains(action.value, "gate record recorded-gate-task") ||
				strings.Contains(action.value, "gate consume recorded-gate-task") {
				return fmt.Errorf("gate was closed after selected override handoff: %q", action.value)
			}
		}
		if action.kind == "tool" && action.name == "Agent" {
			return fmt.Errorf("Agent detour observed in selected override trajectory")
		}
		if action.kind == "tool" && action.name == "Skill" {
			skill, args, _ := strings.Cut(action.text, "\x00")
			if skill != "subspace:r" {
				continue
			}
			overrideCount++
			if overrideAt < 0 {
				overrideAt = i
			}
			if args != "gate "+room {
				return fmt.Errorf("selected override args = %q, want exact room-only %q", args, "gate "+room)
			}
		}
		if action.kind == "text" && assertConciseRecordedGateReview(action.text) == nil {
			return fmt.Errorf("chat gate review observed in selected override trajectory")
		}
	}
	if prepareAt < 0 || commitAt < 0 || overrideAt < 0 || !(prepareAt < commitAt && commitAt < overrideAt) {
		return fmt.Errorf("selected override order prepare=%d commit=%d override=%d", prepareAt, commitAt, overrideAt)
	}
	if overrideCount != 1 {
		return fmt.Errorf("selected override calls = %d, want 1", overrideCount)
	}
	return nil
}
