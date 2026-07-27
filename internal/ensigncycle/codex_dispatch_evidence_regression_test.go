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
	return codexCommandOutput(command, output, 0, "completed")
}

func codexCommandOutput(command, output string, exit int, status string) string {
	line := map[string]any{
		"type": "item.completed",
		"item": map[string]any{
			"type":              "command_execution",
			"command":           command,
			"aggregated_output": output,
			"exit_code":         exit,
			"status":            status,
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
	results = append(results, fmt.Sprintf("error: dispatch build failed for spacedock-ensign-%s-%s", failed, stage))
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

func TestCodexDurableStageReportCountIgnoresProseMentions(t *testing.T) {
	command := "sed -n '1,260p' alpha.md beta.md"
	output := "## Stage Report: done\n\nThe Stage Report is durable.\n"
	reported := codexDurableStageReportTargets(command, output, []string{"alpha", "beta"})
	if reported["alpha"] || reported["beta"] {
		t.Fatalf("one anchored report heading must not prove two named entities: %#v", reported)
	}
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

func TestCodexDispatchEvidenceDoesNotCrossAttributeFailedBatchSuccess(t *testing.T) {
	succeeded := kmReadyOne
	failed := kmReadyTwo
	stream := strings.Join([]string{
		codexMixedDispatchBuildEvidence([]string{succeeded}, failed, kmNextStage),
		codexWaitCompleted(),
		codexEntityReadEvidence(succeeded, kmNextStage, kmNextStage),
		codexEntityReadEvidence(failed, kmNextStage, kmNextStage),
	}, "\n")

	evidence := codexDispatchCompletionEvidenceFromJSONL(stream, []string{succeeded, failed})
	if !evidence.stageReport[succeeded] {
		t.Fatalf("successful build %q was not credited: %+v", succeeded, evidence)
	}
	if evidence.stageReport[failed] || evidence.doneReport[failed] {
		t.Fatalf("failed build %q borrowed another target's success record: %+v", failed, evidence)
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

func TestCodexDispatchEvidenceDoesNotMultiplyAnonymousBatchedReport(t *testing.T) {
	entities := []string{kmReadyOne, kmReadyTwo}
	stream := strings.Join([]string{
		codexDispatchBuildEvidence(kmReadyOne, kmNextStage),
		codexDispatchBuildEvidence(kmReadyTwo, kmNextStage),
		codexWaitCompleted(),
		codexCompletedCommandOutput(
			"sed -n '14,40p' ready-one.md; sed -n '14,40p' ready-two.md",
			"## Stage Report: implementation\n\n- DONE: completed.\n",
		),
	}, "\n")

	evidence := codexDispatchCompletionEvidenceFromJSONL(stream, entities)
	for _, entity := range entities {
		if evidence.stageReport[entity] || evidence.doneReport[entity] {
			t.Errorf("one anonymous report was cross-attributed to %q: %+v", entity, evidence)
		}
	}
}

func TestCodexDispatchEvidencePhaseBindsBatchedFinalization(t *testing.T) {
	entity := kmReadyOne
	build := codexDispatchBuildEvidence(entity, kmNextStage)
	wait := codexWaitCompleted()
	report := codexEntityReadEvidence(entity, kmNextStage, kmNextStage)
	mergeCommand := `for slug in ready-one; do spacedock merge guard "$slug" --workflow-dir .; done`
	finalized := "finalized: " + entity + " -> done (verdict passed), archived. State durability: local-only (no origin remote).\n"
	merge := codexCompletedCommandOutput(mergeCommand, finalized)

	positive := strings.Join([]string{build, wait, report, merge}, "\n")
	evidence := codexDispatchCompletionEvidenceFromJSONL(positive, []string{entity})
	if !evidence.stageReport[entity] || !evidence.doneReport[entity] {
		t.Fatalf("retained build/wait/report/finalization shape was not credited: %+v", evidence)
	}

	controls := map[string]string{
		"missing-build":        strings.Join([]string{wait, report, merge}, "\n"),
		"missing-wait":         strings.Join([]string{build, report, merge}, "\n"),
		"missing-report":       strings.Join([]string{build, wait, merge}, "\n"),
		"failed-command":       strings.Join([]string{build, wait, report, codexCommandOutput(mergeCommand, finalized, 1, "failed")}, "\n"),
		"missing-entity":       strings.Join([]string{build, wait, report, codexCompletedCommandOutput(mergeCommand, "finalized: ready-two -> done\n")}, "\n"),
		"literal-command-only": strings.Join([]string{build, wait, report, codexCompletedCommandOutput("spacedock merge guard "+entity, "")}, "\n"),
		"planted-output":       strings.Join([]string{build, wait, report, codexCompletedCommandOutput("printf planted", finalized)}, "\n"),
		"planted-command-output": strings.Join([]string{build, wait, report,
			codexCompletedCommandOutput("printf 'spacedock merge guard "+entity+"'", finalized)}, "\n"),
	}
	for name, stream := range controls {
		t.Run(name, func(t *testing.T) {
			got := codexDispatchCompletionEvidenceFromJSONL(stream, []string{entity})
			if got.doneReport[entity] {
				t.Fatalf("unqualified finalization was credited: %+v", got)
			}
		})
	}
}
