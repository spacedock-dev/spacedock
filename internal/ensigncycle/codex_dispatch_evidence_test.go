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
	pendingDispatch := map[string]bool{}
	waitedAfterDispatch := map[string]bool{}
	durableDoneReport := map[string]bool{}
	durableStageReport := map[string]bool{}

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
			} `json:"item"`
		}
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		if ev.Item.Type == "command_execution" && ev.Type == "item.completed" {
			for _, entity := range codexDispatchBuildTargets(ev.Item.Command, ev.Item.AggregatedOutput, entities) {
				pendingDispatch[entity] = true
			}
			for _, entity := range entities {
				if codexDurableStageReportForEntity(ev.Item.AggregatedOutput, entity) {
					durableStageReport[entity] = true
					if codexDurableStatusForEntity(ev.Item.AggregatedOutput, entity, "done") {
						durableDoneReport[entity] = true
					}
				}
				if codexMergeGuardCompletedEntity(ev.Item.Command, ev.Item.AggregatedOutput, entity) {
					durableDoneReport[entity] = true
					durableStageReport[entity] = true
				}
			}
		}
		if ev.Item.Type == "collab_tool_call" && ev.Type == "item.completed" && codexWaitTool(ev.Item.Tool) {
			for entity := range pendingDispatch {
				waitedAfterDispatch[entity] = true
			}
		}
	}

	result := newCodexDispatchCompletionEvidence()
	for _, entity := range entities {
		if waitedAfterDispatch[entity] && durableDoneReport[entity] {
			result.doneReport[entity] = true
		}
		if waitedAfterDispatch[entity] && durableStageReport[entity] {
			result.stageReport[entity] = true
		}
	}
	return result
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

func codexDurableStageReportForEntity(output, entity string) bool {
	if !strings.Contains(output, "Stage Report") {
		return false
	}
	return strings.Contains(output, "id: "+entity) ||
		strings.Contains(output, `"id":"`+entity+`"`) ||
		strings.Contains(output, entity+".md")
}

func codexDurableStatusForEntity(output, entity, status string) bool {
	if !codexDurableStageReportForEntity(output, entity) {
		return false
	}
	return strings.Contains(output, "status: "+status) ||
		strings.Contains(output, `"status":"`+status+`"`)
}

func codexMergeGuardCompletedEntity(command, output, entity string) bool {
	return kmMergeGuardTerminalizes(command, entity) &&
		strings.Contains(output, "finalized: "+entity+" -> done")
}
