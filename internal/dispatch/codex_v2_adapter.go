package dispatch

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var codexV2TaskNameUnsafe = regexp.MustCompile(`[^a-z0-9_]+`)

type CodexMultiAgentV2Identity struct {
	Name  string
	Slug  string
	Stage string
}

type CodexMultiAgentV2Spawn struct {
	TaskName string
	Message  string
	Identity CodexMultiAgentV2Identity
}

func (s CodexMultiAgentV2Spawn) ToolArgs() map[string]string {
	return map[string]string{
		"task_name":  s.TaskName,
		"message":    s.Message,
		"fork_turns": "none",
	}
}

func CodexMultiAgentV2SpawnInput(raw []byte) (CodexMultiAgentV2Spawn, error) {
	return CodexMultiAgentV2SpawnInputWithSeen(raw, nil)
}

func CodexMultiAgentV2SpawnInputWithSeen(raw []byte, seen map[string]string) (CodexMultiAgentV2Spawn, error) {
	var helper struct {
		Prompt string  `json:"prompt"`
		Name   *string `json:"name"`
	}
	if err := json.Unmarshal(raw, &helper); err != nil {
		return CodexMultiAgentV2Spawn{}, fmt.Errorf("decode dispatch build output: %w", err)
	}
	if helper.Name == nil || *helper.Name == "" {
		return CodexMultiAgentV2Spawn{}, fmt.Errorf("dispatch build output missing name")
	}
	if helper.Prompt == "" {
		return CodexMultiAgentV2Spawn{}, fmt.Errorf("dispatch build output missing prompt")
	}

	taskName := CodexMultiAgentV2TaskName(*helper.Name)
	if seen != nil {
		if previous, ok := seen[taskName]; ok && previous != *helper.Name {
			return CodexMultiAgentV2Spawn{}, fmt.Errorf("codex multi_agent_v2 task_name collision: %q maps both %q and %q", taskName, previous, *helper.Name)
		}
		seen[taskName] = *helper.Name
	}

	return CodexMultiAgentV2Spawn{
		TaskName: taskName,
		Message:  helper.Prompt,
		Identity: CodexMultiAgentV2Identity{
			Name:  *helper.Name,
			Slug:  codexMultiAgentV2Slug(*helper.Name),
			Stage: codexMultiAgentV2Stage(*helper.Name),
		},
	}, nil
}

func CodexMultiAgentV2TaskName(name string) string {
	taskName := strings.ToLower(name)
	taskName = codexV2TaskNameUnsafe.ReplaceAllString(taskName, "_")
	taskName = strings.Trim(taskName, "_")
	if taskName == "" {
		return "spacedock_worker"
	}
	return taskName
}

func codexMultiAgentV2Slug(name string) string {
	parts := strings.Split(name, "-")
	if len(parts) <= 4 {
		return ""
	}
	return strings.Join(parts[2:len(parts)-1], "-")
}

func codexMultiAgentV2Stage(name string) string {
	if i := strings.LastIndex(name, "-"); i >= 0 && i+1 < len(name) {
		return name[i+1:]
	}
	return ""
}
