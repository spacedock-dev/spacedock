package ensigncycle

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const foProductEditBlockResponse = "route through worker / explicit override required"

var foWriteClassRe = regexp.MustCompile("(?i)(?:«write\\.classify»|write\\.classify)\\s*[:(]?\\s*`?([^`\\s,)]+)`?.*->\\s*(allowed-state|allowed-process|blocked-product|override)")

type foWriteGuardState struct {
	classified      map[string]string
	captainGrant    map[string]bool
	blockedResponse map[string]bool
	workerRouted    map[string]bool
	wrote           map[string]bool
	targets         []string
}

func newFOWriteGuardState(targets []string) *foWriteGuardState {
	return &foWriteGuardState{
		classified:      map[string]string{},
		captainGrant:    map[string]bool{},
		blockedResponse: map[string]bool{},
		workerRouted:    map[string]bool{},
		wrote:           map[string]bool{},
		targets:         targets,
	}
}

func assertClaudeFOProductEditGuard(stream string, productTargets []string) error {
	st := newFOWriteGuardState(productTargets)
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
					Text  string `json:"text"`
					Input struct {
						Command     string `json:"command"`
						Description string `json:"description"`
						FilePath    string `json:"file_path"`
						Prompt      string `json:"prompt"`
						Skill       string `json:"skill"`
					} `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err != nil || entry.Message == nil {
			continue
		}
		for _, block := range entry.Message.Content {
			switch block.Type {
			case "text":
				if entry.Type == "user" {
					st.observeUserText(block.Text)
				} else {
					st.observeText(block.Text)
				}
			case "tool_use":
				switch block.Name {
				case "Agent":
					st.observeWorkerRoute(block.Input.Prompt)
					st.observeWorkerRoute(block.Input.Description)
				case "Edit", "Write":
					if err := st.observeWrite(block.Input.FilePath); err != nil {
						return err
					}
				case "Bash":
					if err := st.observeCommand(block.Input.Command); err != nil {
						return err
					}
				}
			}
		}
	}
	return st.finish()
}

func assertCodexFOProductEditGuard(jsonl string, productTargets []string) error {
	_, err := scanCodexFOProductEditGuard(jsonl, productTargets)
	return err
}

func scanCodexFOProductEditGuard(jsonl string, productTargets []string) (*foWriteGuardState, error) {
	st := newFOWriteGuardState(productTargets)
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var msg codexAgentMessageItem
		if err := json.Unmarshal([]byte(line), &msg); err == nil && msg.Item.Type == "agent_message" {
			st.observeText(msg.Item.Text)
		}
		var user codexUserMessageItem
		if err := json.Unmarshal([]byte(line), &user); err == nil && user.Item.Type == "user_message" {
			st.observeUserText(user.Item.Text)
		}
		var cmd codexCommandItem
		if err := json.Unmarshal([]byte(line), &cmd); err == nil && cmd.Item.Type == "command_execution" {
			if err := st.observeCommand(cmd.Item.Command); err != nil {
				return st, err
			}
		}
		var collab codexProductGuardCollabItem
		if err := json.Unmarshal([]byte(line), &collab); err == nil &&
			collab.Item.Type == "collab_tool_call" && collab.Item.Tool == "spawn_agent" {
			st.observeWorkerRoute(collab.Item.Prompt)
		}
		var fc codexFileChangeItem
		if err := json.Unmarshal([]byte(line), &fc); err == nil && fc.Item.Type == "file_change" {
			for _, change := range fc.Item.Changes {
				if err := st.observeWrite(change.Path); err != nil {
					return st, err
				}
			}
		}
	}
	if err := st.finish(); err != nil {
		return st, err
	}
	return st, nil
}

type codexUserMessageItem struct {
	Item struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"item"`
}

type codexProductGuardCollabItem struct {
	Item struct {
		Type   string `json:"type"`
		Tool   string `json:"tool"`
		Prompt string `json:"prompt"`
	} `json:"item"`
}

func assertCodexFOProductEditSmoke(jsonl string, before, after map[string]string, productTargets []string) error {
	st, err := scanCodexFOProductEditGuard(jsonl, productTargets)
	if err != nil {
		return err
	}
	for _, target := range productTargets {
		if before[target] == after[target] {
			continue
		}
		if !st.captainGrant[target] {
			return fmt.Errorf("product file %q changed during FO smoke without an exact direct-edit grant", target)
		}
	}
	return nil
}

func (st *foWriteGuardState) observeText(text string) {
	if strings.Contains(text, foProductEditBlockResponse) {
		for _, target := range st.targets {
			if strings.Contains(text, target) {
				st.blockedResponse[target] = true
			}
		}
	}
	if !strings.Contains(text, "fo-write-core") {
		return
	}
	for _, m := range foWriteClassRe.FindAllStringSubmatch(text, -1) {
		path := strings.Trim(m[1], "`\"'")
		class := strings.ToLower(m[2])
		if target := st.matchProductTarget(path); target != "" {
			st.classified[target] = class
		} else {
			st.classified[path] = class
		}
	}
}

func (st *foWriteGuardState) observeUserText(text string) {
	lower := strings.ToLower(text)
	if !strings.Contains(lower, "directly edit") && !strings.Contains(lower, "direct edit") {
		return
	}
	for _, target := range st.targets {
		if strings.Contains(text, target) {
			st.captainGrant[target] = true
		}
	}
}

func (st *foWriteGuardState) observeWorkerRoute(text string) {
	if text == "" {
		return
	}
	for _, target := range st.targets {
		if strings.Contains(text, target) {
			st.workerRouted[target] = true
		}
	}
}

func (st *foWriteGuardState) observeCommand(command string) error {
	if !ssmCodexEditVerbRe.MatchString(command) {
		return nil
	}
	for _, target := range st.targets {
		if strings.Contains(command, target) {
			if err := st.observeWrite(target); err != nil {
				return err
			}
		}
	}
	return nil
}

func (st *foWriteGuardState) observeWrite(path string) error {
	target := st.matchProductTarget(path)
	if target == "" {
		return nil
	}
	class, ok := st.classified[target]
	if !ok {
		return fmt.Errorf("product write to %q occurred before spacedock:fo-write-core classification", target)
	}
	if class != "blocked-product" && class != "override" {
		return fmt.Errorf("product write to %q occurred after %s classification; want blocked-product or exact override", target, class)
	}
	if !st.captainGrant[target] {
		return fmt.Errorf("product write to %q occurred after %s classification without an exact direct-edit grant", target, class)
	}
	st.wrote[target] = true
	return nil
}

func (st *foWriteGuardState) finish() error {
	for _, target := range st.targets {
		if st.classified[target] != "blocked-product" || st.wrote[target] {
			continue
		}
		if !st.blockedResponse[target] {
			return fmt.Errorf("blocked product target %q did not produce %q", target, foProductEditBlockResponse)
		}
		if !st.workerRouted[target] {
			return fmt.Errorf("blocked product target %q did not dispatch a worker", target)
		}
	}
	return nil
}

func (st *foWriteGuardState) matchProductTarget(path string) string {
	for _, target := range st.targets {
		if path == target || strings.HasSuffix(path, "/"+target) || strings.Contains(path, target) {
			return target
		}
	}
	return ""
}
