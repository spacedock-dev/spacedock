package ensigncycle

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

type conflictOwnerHandoffExpectation struct {
	DispatchFile string
	WorkerName   string
	Entity       string
	Stage        string
	Branch       string
	Worktree     string
	Marker       string
}

func assertCodexConflictOwnerHandoff(jsonl string, want conflictOwnerHandoffExpectation) error {
	var spawns, startedFollowups, completedFollowups []codexCollabItem
	for _, line := range strings.Split(jsonl, "\n") {
		var event codexCollabItem
		if json.Unmarshal([]byte(line), &event) != nil || event.Item.Type != "collab_tool_call" {
			continue
		}
		switch event.Item.Tool {
		case "spawn_agent":
			if event.Type == "item.completed" {
				spawns = append(spawns, event)
			}
		case "followup_task":
			switch event.Type {
			case "item.started":
				startedFollowups = append(startedFollowups, event)
			case "item.completed":
				completedFollowups = append(completedFollowups, event)
			}
		}
	}
	if len(spawns) != 1 {
		return fmt.Errorf("completed spawn_agent calls = %d, want exactly 1", len(spawns))
	}
	spawn := spawns[0].Item
	if len(spawn.ReceiverThreadIDs) != 1 {
		return fmt.Errorf("owner spawn receiver handles = %v, want exactly 1", spawn.ReceiverThreadIDs)
	}
	for label, token := range map[string]string{"dispatch file": want.DispatchFile, "worker name": want.WorkerName} {
		if token == "" || !strings.Contains(spawn.Prompt, token) {
			return fmt.Errorf("owner spawn prompt does not bind stamped %s %q", label, token)
		}
	}

	followups := startedFollowups
	if len(followups) == 0 {
		followups = completedFollowups
	}
	if len(followups) != 1 {
		return fmt.Errorf("followup_task calls = %d, want exactly 1", len(followups))
	}
	followup := followups[0].Item
	if len(followup.ReceiverThreadIDs) != 1 || followup.ReceiverThreadIDs[0] != spawn.ReceiverThreadIDs[0] {
		return fmt.Errorf("followup receiver handles = %v, want stamped owner handle %q", followup.ReceiverThreadIDs, spawn.ReceiverThreadIDs[0])
	}
	for label, token := range map[string]string{
		"entity": want.Entity, "stage": want.Stage, "branch": want.Branch,
		"worktree": want.Worktree, "marker": want.Marker,
	} {
		if token == "" || !strings.Contains(followup.Prompt, token) {
			return fmt.Errorf("owner followup prompt does not carry %s %q", label, token)
		}
	}
	return nil
}

func TestCodexConflictOwnerHandoffCorrelatesTypedEvents(t *testing.T) {
	want := conflictOwnerHandoffExpectation{
		DispatchFile: "/tmp/spacedock-dispatch/spacedock-ensign-conflict-owner-implementation.md",
		WorkerName:   "spacedock-ensign-conflict-owner-implementation",
		Entity:       "conflict-owner",
		Stage:        "implementation",
		Branch:       "spacedock-ensign/conflict-owner",
		Worktree:     ".worktrees/spacedock-ensign-conflict-owner",
		Marker:       "runtime-worker-owner",
	}
	spawn := func(thread string) string {
		return codexCollabToolLine("item.completed", "spawn_agent", thread, "worker "+want.WorkerName+" reads "+want.DispatchFile)
	}
	followup := func(thread, marker, branch string) string {
		return codexCollabToolLine("item.started", "followup_task", thread, strings.Join([]string{want.Entity, want.Stage, branch, want.Worktree, marker}, " "))
	}
	valid := strings.Join([]string{spawn("owner-thread"), followup("owner-thread", want.Marker, want.Branch)}, "\n")

	tests := []struct {
		name  string
		jsonl string
		pass  bool
	}{
		{name: "one stamped owner handoff", jsonl: valid, pass: true},
		{name: "contract text only", jsonl: `{"type":"item.completed","item":{"type":"agent_message","text":"spawn_agent then followup_task"}}`},
		{name: "zero events", jsonl: ""},
		{name: "extra spawn", jsonl: valid + "\n" + spawn("other-thread")},
		{name: "zero followup", jsonl: spawn("owner-thread")},
		{name: "extra followup", jsonl: valid + "\n" + followup("owner-thread", want.Marker, want.Branch)},
		{name: "wrong handle", jsonl: strings.Join([]string{spawn("owner-thread"), followup("other-thread", want.Marker, want.Branch)}, "\n")},
		{name: "wrong marker", jsonl: strings.Join([]string{spawn("owner-thread"), followup("owner-thread", "wrong-marker", want.Branch)}, "\n")},
		{name: "wrong tuple", jsonl: strings.Join([]string{spawn("owner-thread"), followup("owner-thread", want.Marker, "other/branch")}, "\n")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := assertCodexConflictOwnerHandoff(test.jsonl, want)
			if test.pass && err != nil {
				t.Fatalf("valid typed handoff rejected: %v", err)
			}
			if !test.pass && err == nil {
				t.Fatal("invalid handoff evidence passed")
			}
		})
	}
}
