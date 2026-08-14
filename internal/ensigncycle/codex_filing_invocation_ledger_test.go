package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type codexCommandExecution struct {
	Command  []string `json:"command"`
	Status   string   `json:"status"`
	Stdout   string   `json:"stdout"`
	ExitCode *int     `json:"exit_code"`
}

func codexCommandExecutions(rollout string) ([]codexCommandExecution, error) {
	var executions []codexCommandExecution
	for _, line := range strings.Split(rollout, "\n") {
		var event struct {
			Type    string `json:"type"`
			Payload struct {
				Type string `json:"type"`
				Item struct {
					Type string `json:"type"`
					codexCommandExecution
				} `json:"item"`
			} `json:"payload"`
		}
		if json.Unmarshal([]byte(line), &event) != nil || event.Type != "event_msg" ||
			event.Payload.Type != "item_completed" || event.Payload.Item.Type != "CommandExecution" {
			continue
		}
		executions = append(executions, event.Payload.Item.codexCommandExecution)
	}
	if len(executions) == 0 {
		return nil, fmt.Errorf("correlated Codex rollout has no completed CommandExecution items")
	}
	return executions, nil
}

func assertCodexFilingTransaction(executions []codexCommandExecution, entityPath, slug string) (string, error) {
	createPattern := regexp.MustCompile(`(?:\bnew\b|--new)[ \t]+` + regexp.QuoteMeta(slug) + `(?:[ \t;|&]|$)`)
	var filing *codexCommandExecution
	creates := 0
	for i := range executions {
		if nextIDInvocation.MatchString(strings.Join(executions[i].Command, " ")) {
			return "", fmt.Errorf("filing previewed --next-id instead of using the atomic new path")
		}
		if len(executions[i].Command) != 3 || executions[i].Command[1] != "-lc" && executions[i].Command[1] != "-c" {
			continue
		}
		command := executions[i].Command[2]
		if !commandFilesViaNew(command, slug) {
			continue
		}
		count := len(createPattern.FindAllStringIndex(command, -1))
		creates += count
		if count == 1 && executions[i].Status == "completed" && executions[i].ExitCode != nil && *executions[i].ExitCode == 0 {
			filing = &executions[i]
		}
	}
	if creates != 1 || filing == nil {
		return "", fmt.Errorf("correlated Codex rollout has %d atomic creates for %s, want one completed exit-0 transaction", creates, slug)
	}

	path := filepath.Clean(entityPath)
	receiptPattern := regexp.MustCompile(`(?m)^created: ` + regexp.QuoteMeta(path) + ` id=([^[:space:]]+)$`)
	receipts := receiptPattern.FindAllStringSubmatch(filing.Stdout, -1)
	if len(receipts) != 1 {
		return "", fmt.Errorf("atomic create stdout has %d exact receipts for %s, want one", len(receipts), path)
	}
	if err := assertCodexFiledEntity(path, receipts[0][1]); err != nil {
		return "", err
	}
	return filing.Command[2], nil
}

func assertCodexFiledEntity(path, id string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read landed filing entity: %w", err)
	}
	parts := strings.SplitN(string(data), "---", 3)
	if len(parts) != 3 || strings.TrimSpace(parts[0]) != "" {
		return fmt.Errorf("landed filing entity has malformed frontmatter")
	}
	for _, field := range []string{"\ntitle: Wire The Thing\n", "\nstatus: backlog\n", "\nid: " + id + "\n"} {
		if strings.Count("\n"+parts[1]+"\n", field) != 1 {
			return fmt.Errorf("landed filing entity has invalid %s", strings.TrimSpace(field))
		}
	}
	if body := strings.TrimSpace(parts[2]); body == "" || strings.ContainsAny(body, "\r\n") {
		return fmt.Errorf("landed filing entity body is not one nonblank line")
	}
	return nil
}
