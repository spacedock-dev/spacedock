package ensigncycle

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// codexCompletedCommandOutput builds a completed codex exec command item with durable
// stdout. The PR #486 multi_agent_v2 failures surfaced the dispatch evidence here:
// dispatch-build stdout, wait records, then entity reads showing stage reports.
func codexCompletedCommandOutput(command, output string) string {
	line := map[string]any{
		"type": "item.completed",
		"item": map[string]any{
			"type":              "command_execution",
			"command":           command,
			"aggregated_output": output,
			"exit_code":         0,
			"status":            "completed",
		},
	}
	b, err := json.Marshal(line)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func codexWaitCompleted() string {
	line := map[string]any{
		"type": "item.completed",
		"item": map[string]any{
			"type":   "collab_tool_call",
			"tool":   "wait",
			"status": "completed",
		},
	}
	b, err := json.Marshal(line)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func codexDispatchBuildEvidence(entity, stage string) string {
	command := fmt.Sprintf("spacedock dispatch build --workflow-dir . --entity-path %s.md --stage %s", entity, stage)
	output := fmt.Sprintf(`{
  "schema_version": 2,
  "subagent_type": "spacedock:ensign",
  "description": "%s: %s",
  "dispatch_file_path": "/tmp/spacedock-dispatch/spacedock-ensign-%s-%s.md",
  "prompt": "Read /tmp/spacedock-dispatch/spacedock-ensign-%s-%s.md and treat its content as your assignment.",
  "name": "spacedock-ensign-%s-%s"
}`, entity, stage, entity, stage, entity, stage, entity, stage)
	return codexCompletedCommandOutput(command, output)
}

func codexEntityReadEvidence(entity, status, reportStage string) string {
	command := fmt.Sprintf("sed -n '1,260p' %s.md", entity)
	output := fmt.Sprintf(`---
id: %s
title: %s
status: %s
completed: 2026-07-08
verdict: passed
worktree:
---
# %s

## Stage Report: %s

- DONE: Record the stage result.

### Summary

%s completed %s.
`, entity, entity, status, entity, reportStage, entity, reportStage)
	return codexCompletedCommandOutput(command, output)
}

func TestCodexSmallestSufficientCreditsDispatchBuildWaitAndDurableRead(t *testing.T) {
	stream := strings.Join([]string{
		codexFileChange("/tmp/wf/"+ssmEditFileA, "/tmp/wf/"+ssmEditFileB),
		codexCommand("git commit -m note -- " + ssmStrategyDoc),
		codexDispatchBuildEvidence(ssmCommissionedA, "ready"),
		codexWaitCompleted(),
		codexEntityReadEvidence(ssmCommissionedA, "done", "ready"),
		codexDispatchBuildEvidence(ssmCommissionedB, "ready"),
		codexWaitCompleted(),
		codexEntityReadEvidence(ssmCommissionedB, "done", "ready"),
	}, "\n")

	if err := assertCodexSmallestSufficientMechanism(stream, ssmEditFiles(), ssmCommissioned()); err != nil {
		t.Fatalf("multi_agent_v2 dispatch-build/wait/durable-read evidence must credit commissioned dispatches: %v", err)
	}
}

func TestCodexKeepMovingCreditsDispatchBuildWaitAndDurableRead(t *testing.T) {
	stream := strings.Join([]string{
		codexCommand("spacedock status --workflow-dir . --set " + kmApprovedGate + " status=" + kmNextStage + " verdict=approved"),
		codexDispatchBuildEvidence(kmReadyOne, "done"),
		codexDispatchBuildEvidence(kmApprovedGate, "done"),
		codexDispatchBuildEvidence(kmReadyTwo, "done"),
		codexWaitCompleted(),
		codexEntityReadEvidence(kmReadyOne, "done", "done"),
		codexWaitCompleted(),
		codexEntityReadEvidence(kmReadyTwo, "done", "done"),
		codexWaitCompleted(),
		codexEntityReadEvidence(kmApprovedGate, "done", "done"),
		codexDispatchBuildEvidence(kmQuestioned, "review"),
		codexWaitCompleted(),
		codexEntityReadEvidence(kmQuestioned, "review", "review"),
	}, "\n")

	if err := assertCodexKeepMoving(stream, kmCorrectFinal(), kmIndependent()); err != nil {
		t.Fatalf("multi_agent_v2 dispatch-build/wait/durable-read evidence must credit keep-moving dispatches and re-shape: %v", err)
	}
}
