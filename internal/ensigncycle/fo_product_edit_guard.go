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
	classified map[string]string
	routed     map[string]bool
	wrote      map[string]bool
	targets    []string
}

func newFOWriteGuardState(targets []string) *foWriteGuardState {
	return &foWriteGuardState{
		classified: map[string]string{},
		routed:     map[string]bool{},
		wrote:      map[string]bool{},
		targets:    targets,
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
			Message *struct {
				Content []struct {
					Type  string `json:"type"`
					Name  string `json:"name"`
					Text  string `json:"text"`
					Input struct {
						Command  string `json:"command"`
						FilePath string `json:"file_path"`
						Skill    string `json:"skill"`
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
				st.observeText(block.Text)
			case "tool_use":
				switch block.Name {
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
		var cmd codexCommandItem
		if err := json.Unmarshal([]byte(line), &cmd); err == nil && cmd.Item.Type == "command_execution" {
			if err := st.observeCommand(cmd.Item.Command); err != nil {
				return err
			}
		}
		var fc codexFileChangeItem
		if err := json.Unmarshal([]byte(line), &fc); err == nil && fc.Item.Type == "file_change" {
			for _, change := range fc.Item.Changes {
				if err := st.observeWrite(change.Path); err != nil {
					return err
				}
			}
		}
	}
	return st.finish()
}

func (st *foWriteGuardState) observeText(text string) {
	if strings.Contains(text, foProductEditBlockResponse) {
		for _, target := range st.targets {
			if strings.Contains(text, target) {
				st.routed[target] = true
			}
		}
	}
	if !strings.Contains(text, "fo-write-core") {
		return
	}
	for _, m := range foWriteClassRe.FindAllStringSubmatch(text, -1) {
		path := strings.Trim(m[1], "`\"'")
		class := strings.ToLower(m[2])
		st.classified[path] = class
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
	if class != "override" {
		return fmt.Errorf("product write to %q occurred after %s classification; want exact override", target, class)
	}
	st.wrote[target] = true
	return nil
}

func (st *foWriteGuardState) finish() error {
	for _, target := range st.targets {
		if st.classified[target] == "blocked-product" && !st.wrote[target] && !st.routed[target] {
			return fmt.Errorf("blocked product target %q did not produce %q", target, foProductEditBlockResponse)
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
