package ensigncycle

import (
	"encoding/json"
	"regexp"
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
			reported := codexDurableStageReportTargets(ev.Item.Command, ev.Item.AggregatedOutput, entities)
			for _, entity := range entities {
				if phase[entity] == dispatchWaited && reported[entity] {
					result.stageReport[entity] = true
					if codexDurableStatusForEntity(ev.Item.AggregatedOutput, "done") {
						result.doneReport[entity] = true
					}
				}
				if result.stageReport[entity] && codexMergeGuardCompletedEntity(ev.Item.Command, ev.Item.AggregatedOutput, entity) {
					result.doneReport[entity] = true
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

	type dispatchResult struct {
		DispatchFilePath string `json:"dispatch_file_path"`
	}
	matched := make([]string, 0, len(entities))
	dec := json.NewDecoder(strings.NewReader(output))
	for {
		var result dispatchResult
		if dec.Decode(&result) != nil {
			break
		}
		if result.DispatchFilePath == "" {
			continue
		}
		for _, entity := range entities {
			if strings.Contains(result.DispatchFilePath, "spacedock-ensign-"+entity+"-") {
				matched = append(matched, entity)
				break
			}
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

func codexDurableStageReportTargets(command, output string, entities []string) map[string]bool {
	reported := map[string]bool{}
	reportCount := len(regexp.MustCompile(`(?m)^##[ \t]+Stage Report:`).FindAllStringIndex(output, -1))
	if reportCount == 0 {
		return reported
	}

	var named []string
	for _, entity := range entities {
		if strings.Contains(command, entity+".md") || strings.Contains(command, "status --read "+entity) {
			named = append(named, entity)
		}
	}
	// A command that names several files may prove anonymous report blocks only
	// when it returns a distinct block for every named target.
	if reportCount >= len(named) {
		for _, entity := range named {
			reported[entity] = true
		}
	}
	return reported
}

func codexDurableStatusForEntity(output, status string) bool {
	return strings.Contains(output, "status: "+status) ||
		strings.Contains(output, `"status":"`+status+`"`)
}

func codexMergeGuardCompletedEntity(command, output, entity string) bool {
	if !codexMergeGuardTargetsEntity(command, entity) {
		return false
	}
	prefix := "finalized: " + entity + " -> done"
	full := regexp.MustCompile(`^` + regexp.QuoteMeta(prefix) + ` \(verdict [^)\r\n]+\), archived\.(?: State durability: [^\r\n]+)?$`)
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == prefix || full.MatchString(line) {
			return true
		}
	}
	return false
}

func codexMergeGuardTargetsEntity(command, entity string) bool {
	loopVars := map[string]bool{}
	for _, match := range regexp.MustCompile(`\bfor\s+([A-Za-z_][A-Za-z0-9_]*)\s+in\s+([^;]+)`).FindAllStringSubmatch(command, -1) {
		for _, target := range strings.Fields(match[2]) {
			if strings.Trim(target, `"'`) == entity {
				loopVars[match[1]] = true
			}
		}
	}
	for _, segment := range regexp.MustCompile(`;|\r?\n|&&|\|\|`).Split(command, -1) {
		segment = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(segment), "do "))
		fields := strings.Fields(segment)
		if len(fields) < 4 {
			continue
		}
		launcher := strings.Trim(fields[0], `"'`)
		if launcher != "spacedock" && launcher != "spacedock_launcher" &&
			!strings.HasSuffix(launcher, "/spacedock") && !strings.Contains(launcher, "SPACEDOCK_BIN") {
			continue
		}
		if fields[1] != "merge" || fields[2] != "guard" {
			continue
		}
		target := strings.Trim(fields[3], `"'`)
		if target == entity || (strings.HasPrefix(target, "$") && loopVars[strings.Trim(strings.TrimPrefix(target, "$"), "{}")]) {
			return true
		}
	}
	return false
}
