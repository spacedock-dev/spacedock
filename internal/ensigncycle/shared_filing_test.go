package ensigncycle

import "fmt"

// codexCommandItem is the shared shape used by scenario assertions that still
// grade Codex tool selection. Filing deliberately does not use it: its producer
// proof comes from the execution-grounded invocation ledger.
type codexCommandItem struct {
	Type string `json:"type"`
	Item struct {
		Type    string `json:"type"`
		Command string `json:"command"`
	} `json:"item"`
}

func claudeToolUse(name, inputJSON string) string {
	return `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"` + name + `","input":` + inputJSON + `}]}}`
}

func codexCommand(command string) string {
	return `{"type":"item.completed","item":{"type":"command_execution","command":"` + command + `"}}`
}

func assertFilingViaNew(invocations []testInvocation, slug string) error {
	filedViaNew := false
	sawNextID := false
	for _, invocation := range invocations {
		if invocation.tool != "spacedock" || len(invocation.args) == 0 {
			continue
		}
		if invocation.args[0] == "new" && len(invocation.args) > 1 && invocation.args[1] == slug {
			filedViaNew = true
		}
		if invocation.args[0] == "status" && invocationHasAdjacentArgs(invocation, "--new", slug) {
			filedViaNew = true
		}
		if invocation.args[0] == "status" && invocationHasArg(invocation, "--next-id") {
			sawNextID = true
		}
	}
	if !filedViaNew {
		return fmt.Errorf("the launcher ledger contains no atomic `spacedock new %s` execution", slug)
	}
	if sawNextID {
		return fmt.Errorf("the launcher ledger contains a `status --next-id` execution; atomic filing must mint the id itself")
	}
	return nil
}
