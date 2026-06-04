// ABOUTME: Pi runtime adapter helpers for pi-agent-teams action payloads.
// ABOUTME: Keeps Spacedock lifecycle intents away from Claude team-tool shapes.
package piruntime

import (
	"encoding/json"
	"strings"
)

// TeamsTask is one task entry for pi-agent-teams action=delegate.
type TeamsTask struct {
	Text     string `json:"text"`
	Assignee string `json:"assignee,omitempty"`
}

// TeamsAction is the small Spacedock-owned subset of the pi-agent-teams `teams`
// tool schema needed to express dispatch, follow-up, shutdown, and team done.
type TeamsAction struct {
	Action  string      `json:"action"`
	Tasks   []TeamsTask `json:"tasks,omitempty"`
	Name    string      `json:"name,omitempty"`
	Message string      `json:"message,omitempty"`
	Reason  string      `json:"reason,omitempty"`
	All     bool        `json:"all,omitempty"`
}

func TeamsDelegateAction(workerName, assignment string) TeamsAction {
	return TeamsAction{
		Action: "delegate",
		Tasks: []TeamsTask{{
			Text:     assignment,
			Assignee: workerName,
		}},
	}
}

func TeamsDirectMessageAction(workerName, message string) TeamsAction {
	return TeamsAction{Action: "message_dm", Name: workerName, Message: message}
}

func TeamsShutdownAction(workerName, reason string) TeamsAction {
	return TeamsAction{Action: "member_shutdown", Name: workerName, Reason: reason}
}

func TeamsDoneAction(force bool) TeamsAction {
	return TeamsAction{Action: "team_done", All: force}
}

// ContainsClaudeTeamToolName reports whether an adapter payload accidentally
// embeds Claude team-tool names. It is intentionally string-based over the JSON
// payload so it catches both field values and nested task text.
func ContainsClaudeTeamToolName(action TeamsAction) bool {
	b, _ := json.Marshal(action)
	s := string(b)
	for _, banned := range []string{"Agent", "SendMessage", "TeamCreate", "TeamDelete"} {
		if strings.Contains(s, banned) {
			return true
		}
	}
	return false
}
