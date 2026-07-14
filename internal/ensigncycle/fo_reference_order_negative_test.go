package ensigncycle

import "testing"

const (
	claudeFirstOfficerBaseEvent = `{"type":"user","message":{"content":[{"type":"text","text":"Base directory for this skill: /plugin/skills/first-officer\n\n# First Officer"}]},"isSynthetic":true}`
	codexFirstOfficerBaseEvent  = `{"type":"item.completed","item":{"type":"command_execution","command":"sed -n '1,80p' /plugin/skills/first-officer/SKILL.md","status":"completed","exit_code":0,"aggregated_output":"name: first-officer"}}`
)

func TestFOReferenceOrderNormalizersPreserveHostEventOrder(t *testing.T) {
	claude := claudeFirstOfficerBaseEvent + "\n" + `{"type":"assistant","message":{"id":"a","content":[{"type":"tool_use","id":"1","name":"Read","input":{"file_path":"/plugin/skills/first-officer/references/first-officer-shared-core.md"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"1","content":"# First Officer Shared Core"}]}}
{"type":"assistant","message":{"id":"b","content":[{"type":"tool_use","id":"2","name":"Read","input":{"file_path":"/plugin/skills/first-officer/references/claude-first-officer-runtime.md"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"2","content":"# Claude Code First Officer Runtime"}]}}
{"type":"assistant","message":{"id":"engage","content":[{"type":"tool_use","id":"ready","name":"Bash","input":{"command":"$B state ready --workflow-dir /tmp/workflow --json"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"ready","content":"{\"ready\":true,\"mod-block\":\"merge:pr-merge\"}"}]}}
{"type":"assistant","message":{"id":"c","content":[{"type":"text","text":"write.classify says allowed"}]}}
{"type":"assistant","message":{"id":"d","content":[{"type":"tool_use","id":"3","name":"Read","input":{"file_path":"/plugin/skills/first-officer/references/fo-write-core.md"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"3","content":"# First Officer Write Core"}]}}
{"type":"assistant","message":{"id":"e","content":[{"type":"tool_use","id":"4","name":"Read","input":{"file_path":"/plugin/skills/first-officer/references/fo-merge-core.md"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"4","content":"# First Officer Merge Core"}]}}
{"type":"assistant","message":{"id":"f","content":[{"type":"tool_use","id":"5","name":"Bash","input":{"command":"$B merge guard merged-pr --verdict passed --json"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"5","content":"{\"signal\":\"finalized\"}"}]}}`
	claudeEvents := normalizeClaudeFOReferenceEvents(claude)
	if err := assertFOReferenceJourney(claudeEvents, "recovery"); err != nil {
		t.Fatalf("Claude positive recovery stream: %v", err)
	}

	codex := codexFirstOfficerBaseEvent + "\n" + `{"type":"item.completed","item":{"type":"command_execution","command":"sed -n 1,220p /plugin/skills/first-officer/references/first-officer-shared-core.md && sed -n 1,220p /plugin/skills/first-officer/references/codex-first-officer-runtime.md","aggregated_output":"# First Officer Shared Core\n# Codex First Officer Runtime"}}
{"type":"item.completed","item":{"type":"command_execution","command":"$B state ready --workflow-dir /tmp/workflow --json","status":"completed","exit_code":0,"aggregated_output":"{\"ready\":true,\"mod-block\":\"merge:pr-merge\"}"}}
{"type":"item.completed","item":{"type":"agent_message","text":"write.classify says allowed"}}
{"type":"item.completed","item":{"type":"command_execution","command":"sed -n 1,220p /plugin/skills/first-officer/references/fo-write-core.md","aggregated_output":"# First Officer Write Core"}}
{"type":"item.completed","item":{"type":"command_execution","command":"sed -n 1,220p /plugin/skills/first-officer/references/fo-merge-core.md","aggregated_output":"# First Officer Merge Core"}}
{"type":"item.completed","item":{"type":"command_execution","command":"$B merge guard merged-pr --verdict passed --json"}}`
	codexEvents := normalizeCodexFOReferenceEvents(codex)
	if err := assertFOReferenceJourney(codexEvents, "recovery"); err != nil {
		t.Fatalf("Codex positive recovery stream: %v", err)
	}
}

func TestFOReferenceOrderOracleRejectsAdversarialControls(t *testing.T) {
	base := []foReferenceEvent{foSharedRead, foRuntimeRead, foWriteRead, foMutation}
	if err := assertFOReferenceJourney(base, "filing"); err != nil {
		t.Fatalf("positive filing events: %v", err)
	}
	controls := []struct {
		name    string
		journey string
		events  []foReferenceEvent
	}{
		{"eager gate", "gate", []foReferenceEvent{foSharedRead, foWriteRead, foRuntimeRead}},
		{"gate mutation", "gate", []foReferenceEvent{foSharedRead, foRuntimeRead, foMutation}},
		{"gate terminal", "gate", []foReferenceEvent{foSharedRead, foRuntimeRead, foTerminal}},
		{"missing shared precondition", "filing", []foReferenceEvent{foRuntimeRead, foWriteRead, foMutation}},
		{"failed shared precondition", "filing", []foReferenceEvent{foFailedRead, foRuntimeRead, foWriteRead, foMutation}},
		{"missing runtime precondition", "filing", []foReferenceEvent{foSharedRead, foWriteRead, foMutation}},
		{"failed runtime precondition", "filing", []foReferenceEvent{foSharedRead, foFailedRead, foWriteRead, foMutation}},
		{"runtime before shared", "filing", []foReferenceEvent{foRuntimeRead, foSharedRead, foWriteRead, foMutation}},
		{"write before runtime", "filing", []foReferenceEvent{foSharedRead, foWriteRead, foRuntimeRead, foMutation}},
		{"missing write", "filing", []foReferenceEvent{foSharedRead, foRuntimeRead, foMutation}},
		{"write after mutation", "filing", []foReferenceEvent{foSharedRead, foRuntimeRead, foMutation, foWriteRead}},
		{"wrong path retry", "filing", []foReferenceEvent{foSharedRead, foRuntimeRead, foWrongPath, foWriteRead, foMutation}},
		{"wrapper invocation", "filing", []foReferenceEvent{foSharedRead, foRuntimeRead, foWrapperSkill, foWriteRead, foMutation}},
		{"broad search", "filing", []foReferenceEvent{foSharedRead, foRuntimeRead, foBroadSearch, foWriteRead, foMutation}},
		{"merge after guard", "terminal", []foReferenceEvent{foSharedRead, foRuntimeRead, foWriteRead, foMutation, foMergeGuard, foMergeRead}},
		{"mod block before engage", "recovery", []foReferenceEvent{foSharedRead, foRuntimeRead, foModBlockSeen, foEngage, foWriteRead, foMergeRead, foMergeGuard, foMutation}},
		{"write before mod block", "recovery", []foReferenceEvent{foSharedRead, foRuntimeRead, foEngage, foWriteRead, foModBlockSeen, foMergeRead, foMergeGuard, foMutation}},
		{"reversed recovery", "recovery", []foReferenceEvent{foSharedRead, foRuntimeRead, foEngage, foModBlockSeen, foWriteRead, foMutation, foMergeGuard, foMergeRead}},
		{"repeated merge work", "recovery", []foReferenceEvent{foSharedRead, foRuntimeRead, foEngage, foModBlockSeen, foWriteRead, foMergeRead, foRepeatedMerge, foMergeAction, foMutation}},
	}
	for _, control := range controls {
		t.Run(control.name, func(t *testing.T) {
			if err := assertFOReferenceJourney(control.events, control.journey); err == nil {
				t.Fatalf("planted %s regression passed: %v", control.name, control.events)
			}
		})
	}
}

func TestFOReferenceOrderOracleIgnoresClassificationNarration(t *testing.T) {
	claudeNarration := `{"type":"assistant","message":{"id":"a","content":[{"type":"text","text":"I invoked write.classify before editing."}]}}`
	if events := normalizeClaudeFOReferenceEvents(claudeNarration); len(events) != 0 {
		t.Fatalf("Claude narration produced execution events: %v", events)
	}
	codexNarration := `{"type":"item.completed","item":{"type":"agent_message","text":"write.classify approved the mutation"}}`
	if events := normalizeCodexFOReferenceEvents(codexNarration); len(events) != 0 {
		t.Fatalf("Codex narration produced execution events: %v", events)
	}
}

func TestFOReferenceOrderRequiresSuccessfulExactReads(t *testing.T) {
	pathMention := codexFirstOfficerBaseEvent + "\n" + `{"type":"item.completed","item":{"type":"command_execution","command":"echo references/fo-write-core.md","status":"completed","exit_code":0}}
{"type":"item.completed","item":{"type":"command_execution","command":"spacedock new task","status":"completed","exit_code":0}}`
	if err := assertFOReferenceJourney(normalizeCodexFOReferenceEvents(pathMention), "filing"); err == nil {
		t.Fatal("a non-read path mention satisfied the write-core boundary")
	}
	sedMention := codexFirstOfficerBaseEvent + "\n" + `{"type":"item.completed","item":{"type":"command_execution","command":"sed -n '/references/fo-write-core.md/p' README.md","status":"completed","exit_code":0,"aggregated_output":"references/fo-write-core.md"}}
{"type":"item.completed","item":{"type":"command_execution","command":"spacedock new task","status":"completed","exit_code":0}}`
	if err := assertFOReferenceJourney(normalizeCodexFOReferenceEvents(sedMention), "filing"); err == nil {
		t.Fatal("a sed output that only echoed the path mention satisfied the canonical-body boundary")
	}

	failedClaudeRead := claudeFirstOfficerBaseEvent + "\n" + `{"type":"assistant","message":{"id":"a","content":[{"type":"tool_use","id":"read","name":"Read","input":{"file_path":"/plugin/skills/first-officer/references/fo-write-core.md"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"read","is_error":true,"content":"not found"}]}}
{"type":"assistant","message":{"id":"b","content":[{"type":"tool_use","id":"write","name":"Bash","input":{"command":"spacedock new task"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"write","content":"created"}]}}`
	if err := assertFOReferenceJourney(normalizeClaudeFOReferenceEvents(failedClaudeRead), "filing"); err == nil {
		t.Fatal("a failed Claude Read satisfied the write-core boundary")
	}

	failedThenRetried := codexFirstOfficerBaseEvent + "\n" + `{"type":"item.completed","item":{"type":"command_execution","command":"cat /plugin/skills/first-officer/references/fo-write-core.md","status":"failed","exit_code":1}}
{"type":"item.completed","item":{"type":"command_execution","command":"cat /plugin/skills/first-officer/references/fo-write-core.md","status":"completed","exit_code":0,"aggregated_output":"# First Officer Write Core"}}
{"type":"item.completed","item":{"type":"command_execution","command":"spacedock new task","status":"completed","exit_code":0}}`
	if err := assertFOReferenceJourney(normalizeCodexFOReferenceEvents(failedThenRetried), "filing"); err == nil {
		t.Fatal("a failed-then-retried Codex read satisfied the no-retry boundary")
	}

	emptyRead := codexFirstOfficerBaseEvent + "\n" + `{"type":"item.completed","item":{"type":"command_execution","command":"cat /plugin/skills/first-officer/references/fo-write-core.md","status":"completed","exit_code":0,"aggregated_output":""}}
{"type":"item.completed","item":{"type":"command_execution","command":"spacedock new task","status":"completed","exit_code":0}}`
	if err := assertFOReferenceJourney(normalizeCodexFOReferenceEvents(emptyRead), "filing"); err == nil {
		t.Fatal("an empty successful-status read satisfied the canonical body boundary")
	}

	wrongThenRetried := codexFirstOfficerBaseEvent + "\n" + `{"type":"item.completed","item":{"type":"command_execution","command":"cat /plugin/fo-write-core.md","status":"failed","exit_code":1}}
{"type":"item.completed","item":{"type":"command_execution","command":"cat /plugin/skills/first-officer/references/fo-write-core.md","status":"completed","exit_code":0,"aggregated_output":"# First Officer Write Core"}}
{"type":"item.completed","item":{"type":"command_execution","command":"spacedock new task","status":"completed","exit_code":0}}`
	if err := assertFOReferenceJourney(normalizeCodexFOReferenceEvents(wrongThenRetried), "filing"); err == nil {
		t.Fatal("a wrong-path retry satisfied the exact no-retry boundary")
	}

	wrongBaseSameSuffix := codexFirstOfficerBaseEvent + "\n" + `{"type":"item.completed","item":{"type":"command_execution","command":"cat /other/skills/first-officer/references/fo-write-core.md","status":"completed","exit_code":0,"aggregated_output":"# First Officer Write Core"}}
{"type":"item.completed","item":{"type":"command_execution","command":"spacedock new task","status":"completed","exit_code":0}}`
	if err := assertFOReferenceJourney(normalizeCodexFOReferenceEvents(wrongBaseSameSuffix), "filing"); err == nil {
		t.Fatal("a same-suffix read under a different skill base satisfied the canonical boundary")
	}
	wrongBaseClaudeRead := claudeFirstOfficerBaseEvent + "\n" + `{"type":"assistant","message":{"id":"a","content":[{"type":"tool_use","id":"read","name":"Read","input":{"file_path":"/other/skills/first-officer/references/fo-write-core.md"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"read","content":"# First Officer Write Core"}]}}
{"type":"assistant","message":{"id":"b","content":[{"type":"tool_use","id":"write","name":"Bash","input":{"command":"spacedock new task"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"write","content":"created"}]}}`
	if err := assertFOReferenceJourney(normalizeClaudeFOReferenceEvents(wrongBaseClaudeRead), "filing"); err == nil {
		t.Fatal("a Claude Read with the right suffix under the wrong skill base satisfied the canonical boundary")
	}

	compoundWrongThenCanonical := codexFirstOfficerBaseEvent + "\n" + `{"type":"item.completed","item":{"type":"command_execution","command":"cat /other/skills/first-officer/references/fo-write-core.md /plugin/skills/first-officer/references/fo-write-core.md","status":"completed","exit_code":0,"aggregated_output":"# First Officer Write Core"}}
{"type":"item.completed","item":{"type":"command_execution","command":"spacedock new task","status":"completed","exit_code":0}}`
	if err := assertFOReferenceJourney(normalizeCodexFOReferenceEvents(compoundWrongThenCanonical), "filing"); err == nil {
		t.Fatal("a compound wrong-base then canonical read suppressed the wrong-path hazard")
	}
}

func TestFOReferenceOrderPreservesLoadedSkillBaseCase(t *testing.T) {
	const skillBase = "/PluginCache/skills/first-officer"
	events := classifyFOCommand("cat "+skillBase+"/references/fo-write-core.md", skillBase)
	if eventIndex(events, foWriteRead) < 0 || eventIndex(events, foWrongPath) >= 0 {
		t.Fatalf("uppercase-containing loaded skill base classified as %v, want write-read", events)
	}
}

func TestFOReferenceOrderDetectsShellMutationBoundaries(t *testing.T) {
	commands := []string{
		`printf '%s' value > task.md`,
		`printf '%s' value >> task.md`,
		`sed -i.bak 's/a/b/' task.md`,
		`mv task.tmp task.md`,
		`git -C /repo add task.md`,
		`git -C /repo commit -m update`,
	}
	for _, command := range commands {
		t.Run(command, func(t *testing.T) {
			if eventIndex(classifyFOCommand(command, ""), foMutation) < 0 {
				t.Fatalf("shell mutation was not detected: %s", command)
			}
			stream := `{"type":"item.completed","item":{"type":"command_execution","command":` + mustJSONString(command) + `,"status":"completed","exit_code":0}}`
			if err := assertFOReferenceJourney(normalizeCodexFOReferenceEvents(stream), "filing"); err == nil {
				t.Fatalf("mutation %q passed without a successful write-core read", command)
			}
		})
	}
}
