package ensigncycle

import "testing"

func TestFOReferenceOrderNormalizersPreserveHostEventOrder(t *testing.T) {
	claude := `{"type":"assistant","message":{"id":"a","content":[{"type":"tool_use","id":"1","name":"Read","input":{"file_path":"/plugin/skills/first-officer/references/first-officer-shared-core.md"}}]}}
{"type":"assistant","message":{"id":"b","content":[{"type":"tool_use","id":"2","name":"Read","input":{"file_path":"/plugin/skills/first-officer/references/claude-first-officer-runtime.md"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","content":"{\"mod-block\":\"merge:pr-merge\"}"}]}}
{"type":"assistant","message":{"id":"c","content":[{"type":"text","text":"write.classify says allowed"}]}}
{"type":"assistant","message":{"id":"d","content":[{"type":"tool_use","id":"3","name":"Read","input":{"file_path":"/plugin/skills/first-officer/references/fo-write-core.md"}}]}}
{"type":"assistant","message":{"id":"e","content":[{"type":"tool_use","id":"4","name":"Read","input":{"file_path":"/plugin/skills/first-officer/references/fo-merge-core.md"}}]}}
{"type":"assistant","message":{"id":"f","content":[{"type":"tool_use","id":"5","name":"Bash","input":{"command":"$B merge guard merged-pr --verdict passed --json"}}]}}`
	claudeEvents := normalizeClaudeFOReferenceEvents(claude)
	if err := assertFOReferenceJourney(claudeEvents, "recovery"); err != nil {
		t.Fatalf("Claude positive recovery stream: %v", err)
	}

	codex := `{"type":"item.completed","item":{"type":"command_execution","command":"sed -n 1,220p /plugin/skills/first-officer/references/first-officer-shared-core.md && sed -n 1,220p /plugin/skills/first-officer/references/codex-first-officer-runtime.md"}}
{"type":"item.completed","item":{"type":"command_execution","command":"$B status --boot --json","aggregated_output":"{\"mod-block\":\"merge:pr-merge\"}"}}
{"type":"item.completed","item":{"type":"agent_message","text":"write.classify says allowed"}}
{"type":"item.completed","item":{"type":"command_execution","command":"sed -n 1,220p /plugin/skills/first-officer/references/fo-write-core.md"}}
{"type":"item.completed","item":{"type":"command_execution","command":"sed -n 1,220p /plugin/skills/first-officer/references/fo-merge-core.md"}}
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
		{"missing write", "filing", []foReferenceEvent{foSharedRead, foRuntimeRead, foMutation}},
		{"write after mutation", "filing", []foReferenceEvent{foSharedRead, foRuntimeRead, foMutation, foWriteRead}},
		{"wrong path retry", "filing", []foReferenceEvent{foSharedRead, foRuntimeRead, foWrongPath, foWriteRead, foMutation}},
		{"wrapper invocation", "filing", []foReferenceEvent{foSharedRead, foRuntimeRead, foWrapperSkill, foWriteRead, foMutation}},
		{"broad search", "filing", []foReferenceEvent{foSharedRead, foRuntimeRead, foBroadSearch, foWriteRead, foMutation}},
		{"merge after guard", "terminal", []foReferenceEvent{foSharedRead, foRuntimeRead, foWriteRead, foMutation, foMergeGuard, foMergeRead}},
		{"reversed recovery", "recovery", []foReferenceEvent{foSharedRead, foRuntimeRead, foModBlockSeen, foWriteRead, foMutation, foMergeGuard, foMergeRead}},
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
