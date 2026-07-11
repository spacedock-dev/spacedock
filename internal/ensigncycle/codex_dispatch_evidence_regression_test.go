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

func codexFailedDispatchBuildEvidence(entity, stage string) string {
	line := map[string]any{
		"type": "item.completed",
		"item": map[string]any{
			"type":              "command_execution",
			"command":           fmt.Sprintf("spacedock dispatch build --workflow-dir . --entity-path %s.md --stage %s", entity, stage),
			"aggregated_output": "error: dispatch build failed",
			"exit_code":         1,
			"status":            "failed",
		},
	}
	b, err := json.Marshal(line)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func codexMixedDispatchBuildEvidence(succeeded []string, failed, stage string) string {
	commands := make([]string, 0, len(succeeded)+1)
	results := make([]string, 0, len(succeeded)+1)
	for _, entity := range succeeded {
		commands = append(commands, fmt.Sprintf("spacedock dispatch build --entity-path %s.md --stage %s", entity, stage))
		results = append(results, fmt.Sprintf(`{"dispatch_file_path":"/tmp/spacedock-dispatch/spacedock-ensign-%s-%s.md"}`, entity, stage))
	}
	commands = append(commands, fmt.Sprintf("spacedock dispatch build --entity-path %s.md --stage %s", failed, stage))
	results = append(results, "error: dispatch build failed")
	line := map[string]any{
		"type": "item.completed",
		"item": map[string]any{
			"type":              "command_execution",
			"command":           strings.Join(commands, "; "),
			"aggregated_output": strings.Join(results, "\n"),
			"exit_code":         1,
			"status":            "failed",
		},
	}
	b, err := json.Marshal(line)
	if err != nil {
		panic(err)
	}
	return string(b)
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

func TestCodexKeepMovingCreditsPostWaitWorkingStageReports(t *testing.T) {
	stream := strings.Join([]string{
		codexCommand("spacedock status --workflow-dir . --set " + kmApprovedGate + " status=" + kmNextStage + " verdict=approved"),
		codexDispatchBuildEvidence(kmReadyOne, kmNextStage),
		codexDispatchBuildEvidence(kmApprovedGate, kmNextStage),
		codexDispatchBuildEvidence(kmReadyTwo, kmNextStage),
		codexWaitCompleted(),
		codexEntityReadEvidence(kmReadyOne, kmNextStage, kmNextStage),
		codexEntityReadEvidence(kmApprovedGate, kmNextStage, kmNextStage),
		codexEntityReadEvidence(kmReadyTwo, kmNextStage, kmNextStage),
		codexCommand("spacedock status --workflow-dir . --set " + kmQuestioned + " status=" + kmReopenStage + " verdict=questioned"),
	}, "\n")

	if err := assertCodexKeepMoving(stream, kmCorrectFinal(), kmIndependent()); err != nil {
		t.Fatalf("build + wait + durable working-stage report must prove dispatch before FO terminalization: %v", err)
	}
}

func TestCodexDispatchEvidenceRejectsOutOfOrderAndFailedBuilds(t *testing.T) {
	entity := kmReadyOne
	cases := map[string]string{
		"stale report before build": strings.Join([]string{
			codexEntityReadEvidence(entity, kmNextStage, kmNextStage),
			codexDispatchBuildEvidence(entity, kmNextStage),
			codexWaitCompleted(),
		}, "\n"),
		"report before wait": strings.Join([]string{
			codexDispatchBuildEvidence(entity, kmNextStage),
			codexEntityReadEvidence(entity, kmNextStage, kmNextStage),
			codexWaitCompleted(),
		}, "\n"),
		"failed build then wait and report": strings.Join([]string{
			codexFailedDispatchBuildEvidence(entity, kmNextStage),
			codexWaitCompleted(),
			codexEntityReadEvidence(entity, kmNextStage, kmNextStage),
		}, "\n"),
	}

	for name, stream := range cases {
		t.Run(name, func(t *testing.T) {
			evidence := codexDispatchCompletionEvidenceFromJSONL(stream, []string{entity})
			if evidence.stageReport[entity] || evidence.doneReport[entity] {
				t.Fatalf("invalid temporal stream credited dispatch evidence: %+v", evidence)
			}
		})
	}
}

func TestCodexDispatchEvidenceCreditsOnlySuccessfulBuildsInFailedBatch(t *testing.T) {
	succeeded := []string{kmApprovedGate, kmReadyOne, kmReadyTwo}
	failed := kmQuestioned
	stream := []string{codexMixedDispatchBuildEvidence(succeeded, failed, kmNextStage), codexWaitCompleted()}
	for _, entity := range succeeded {
		stream = append(stream, codexEntityReadEvidence(entity, kmNextStage, kmNextStage))
	}
	stream = append(stream, codexEntityReadEvidence(failed, kmNextStage, kmNextStage))

	evidence := codexDispatchCompletionEvidenceFromJSONL(strings.Join(stream, "\n"), append(succeeded, failed))
	for _, entity := range succeeded {
		if !evidence.stageReport[entity] {
			t.Errorf("successful build %q in failed batch was not credited", entity)
		}
	}
	if evidence.stageReport[failed] || evidence.doneReport[failed] {
		t.Fatalf("failed build %q in mixed batch was credited: %+v", failed, evidence)
	}
}

func TestCodexDispatchEvidenceCreditsNamedBatchedDurableReads(t *testing.T) {
	entities := []string{kmApprovedGate, kmReadyOne, kmReadyTwo}
	stream := []string{codexDispatchBuildEvidence(kmApprovedGate, kmNextStage)}
	stream = append(stream, codexDispatchBuildEvidence(kmReadyOne, kmNextStage))
	stream = append(stream, codexDispatchBuildEvidence(kmReadyTwo, kmNextStage), codexWaitCompleted())
	stream = append(stream, codexCompletedCommandOutput(
		"sed -n '14,40p' approved-gate.md; sed -n '14,40p' ready-one.md; sed -n '14,40p' ready-two.md",
		strings.Repeat("## Stage Report: implementation\n\n- DONE: completed.\n", len(entities)),
	))

	evidence := codexDispatchCompletionEvidenceFromJSONL(strings.Join(stream, "\n"), entities)
	for _, entity := range entities {
		if !evidence.stageReport[entity] {
			t.Errorf("named entity %q in successful batched durable read was not credited", entity)
		}
	}
}
