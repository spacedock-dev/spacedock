// ABOUTME: Pi teams adapter contract tests — Spacedock lifecycle intents map
// ABOUTME: to pi-agent-teams `teams` action payloads, not Claude team tools.
package piruntime

import "testing"

func TestTeamsAdapterMapsLifecycleActions(t *testing.T) {
	spawn := TeamsDelegateAction("spacedock-ensign-pi-runtime-support-implementation", "Read /tmp/dispatch.md and implement.")
	if spawn.Action != "delegate" {
		t.Fatalf("spawn action=%q, want delegate", spawn.Action)
	}
	if len(spawn.Tasks) != 1 || spawn.Tasks[0].Assignee != "spacedock-ensign-pi-runtime-support-implementation" {
		t.Fatalf("delegate task assignment not preserved: %#v", spawn.Tasks)
	}
	if spawn.Tasks[0].Text != "Read /tmp/dispatch.md and implement." {
		t.Fatalf("delegate task text drifted: %q", spawn.Tasks[0].Text)
	}

	followup := TeamsDirectMessageAction("worker-a", "Fix the rejected AC-2 evidence.")
	if followup.Action != "message_dm" || followup.Name != "worker-a" || followup.Message != "Fix the rejected AC-2 evidence." {
		t.Fatalf("bad follow-up mapping: %#v", followup)
	}

	shutdown := TeamsShutdownAction("worker-a", "stage complete")
	if shutdown.Action != "member_shutdown" || shutdown.Name != "worker-a" || shutdown.Reason != "stage complete" {
		t.Fatalf("bad shutdown mapping: %#v", shutdown)
	}

	done := TeamsDoneAction(true)
	if done.Action != "team_done" || !done.All {
		t.Fatalf("bad team done mapping: %#v", done)
	}
}

func TestTeamsAdapterPayloadsContainNoClaudeToolNames(t *testing.T) {
	payloads := []TeamsAction{
		TeamsDelegateAction("worker-a", "Do work"),
		TeamsDirectMessageAction("worker-a", "Continue"),
		TeamsShutdownAction("worker-a", "done"),
		TeamsDoneAction(true),
	}
	for _, payload := range payloads {
		if ContainsClaudeTeamToolName(payload) {
			t.Fatalf("payload contains Claude team-tool name: %#v", payload)
		}
	}
}
