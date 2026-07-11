package ensigncycle

import (
	"encoding/json"
	"strings"
)

type codexDispatchCompletionEvidence struct {
	doneReport  map[string]bool
	stageReport map[string]bool
}

func newCodexDispatchCompletionEvidence() codexDispatchCompletionEvidence {
	return codexDispatchCompletionEvidence{
		doneReport:  map[string]bool{},
		stageReport: map[string]bool{},
	}
}

// codexDispatchCompletionEvidenceFromJSONL credits the multi_agent_v2 shape where
// the stream omits spawn_agent records but still shows dispatch build, a foreground
// wait, and durable entity/report state after completion.
func codexDispatchCompletionEvidenceFromJSONL(jsonl string, entities []string) codexDispatchCompletionEvidence {
	type dispatchPhase uint8
	const (
		dispatchNone dispatchPhase = iota
		dispatchBuilt
		dispatchWaited
	)
	phase := map[string]dispatchPhase{}
	result := newCodexDispatchCompletionEvidence()

	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ev struct {
			Type string `json:"type"`
			Item struct {
				Type             string `json:"type"`
				Tool             string `json:"tool"`
				Command          string `json:"command"`
				AggregatedOutput string `json:"aggregated_output"`
				ExitCode         *int   `json:"exit_code"`
				Status           string `json:"status"`
			} `json:"item"`
		}
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		if ev.Item.Type == "command_execution" && ev.Type == "item.completed" {
			for _, entity := range codexSuccessfulDispatchBuildTargets(
				ev.Item.Command,
				ev.Item.AggregatedOutput,
				ev.Item.ExitCode,
				ev.Item.Status,
				entities,
			) {
				phase[entity] = dispatchBuilt
			}
			if ev.Item.ExitCode == nil || *ev.Item.ExitCode != 0 || ev.Item.Status == "failed" {
				continue
			}
			for _, entity := range entities {
				if phase[entity] == dispatchWaited && codexDurableStageReportForEntity(ev.Item.Command, ev.Item.AggregatedOutput, entity) {
					result.stageReport[entity] = true
					if codexDurableStatusForEntity(ev.Item.Command, ev.Item.AggregatedOutput, entity, "done") {
						result.doneReport[entity] = true
					}
				}
				if codexMergeGuardCompletedEntity(ev.Item.Command, ev.Item.AggregatedOutput, entity) {
					result.doneReport[entity] = true
					result.stageReport[entity] = true
				}
			}
		}
		if ev.Item.Type == "collab_tool_call" && ev.Type == "item.completed" &&
			codexWaitTool(ev.Item.Tool) && ev.Item.Status == "completed" {
			for entity, current := range phase {
				if current == dispatchBuilt {
					phase[entity] = dispatchWaited
				}
			}
		}
	}
	return result
}

// codexSuccessfulDispatchBuildTargets accounts for Codex batching several
// dispatch builds into one shell item. A zero-exit item proves every addressed
// build succeeded. For a non-zero batch, only the targets with a complete
// dispatch-build JSON result are proven; a later command in the batch may have
// failed after those results were emitted.
func codexSuccessfulDispatchBuildTargets(command, output string, exitCode *int, status string, entities []string) []string {
	if !strings.Contains(command, "dispatch build") {
		return nil
	}
	if exitCode != nil && *exitCode == 0 && status != "failed" {
		return codexDispatchBuildTargets(command, output, entities)
	}

	var matched []string
	for _, entity := range entities {
		if strings.Contains(output, `"dispatch_file_path":`) &&
			strings.Contains(output, "spacedock-ensign-"+entity+"-") {
			matched = append(matched, entity)
		}
	}
	return matched
}

func codexDispatchBuildTargets(command, output string, entities []string) []string {
	if !strings.Contains(command, "dispatch build") {
		return nil
	}
	haystack := command + "\n" + output
	var matched []string
	for _, entity := range entities {
		if strings.Contains(haystack, entity+".md") ||
			strings.Contains(haystack, "spacedock-ensign-"+entity+"-") ||
			strings.Contains(haystack, `"name": "spacedock-ensign-`+entity) ||
			strings.Contains(haystack, `"name":"spacedock-ensign-`+entity) {
			matched = append(matched, entity)
		}
	}
	return matched
}

func codexWaitTool(tool string) bool {
	return tool == "wait" || tool == "wait_agent" || tool == "collab:wait"
}

func codexDurableStageReportForEntity(command, output, entity string) bool {
	if !strings.Contains(output, "Stage Report") {
		return false
	}
	return strings.Contains(output, "id: "+entity) ||
		strings.Contains(output, `"id":"`+entity+`"`) ||
		strings.Contains(output, entity+".md") ||
		strings.Contains(command, entity+".md") ||
		strings.Contains(command, "status --read "+entity)
}

func codexDurableStatusForEntity(command, output, entity, status string) bool {
	if !codexDurableStageReportForEntity(command, output, entity) {
		return false
	}
	return strings.Contains(output, "status: "+status) ||
		strings.Contains(output, `"status":"`+status+`"`)
}

func codexMergeGuardCompletedEntity(command, output, entity string) bool {
	return kmMergeGuardTerminalizes(command, entity) &&
		strings.Contains(output, "finalized: "+entity+" -> done")
}
