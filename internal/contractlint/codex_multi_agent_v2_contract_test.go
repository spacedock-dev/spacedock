package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCodexCurrentMultiAgentRuntimeReferencesUseLiveToolSurfaceProbe(t *testing.T) {
	root := repoRoot(t)
	paths := []string{
		filepath.Join("skills", "first-officer", "references", "codex-first-officer-runtime.md"),
		filepath.Join("skills", "ensign", "references", "codex-ensign-runtime.md"),
		filepath.Join("skills", "feedback-rejection-flow", "SKILL.md"),
		filepath.Join("docs", "dev", "codex-idle-notification-probe.md"),
	}

	joined := ""
	for _, rel := range paths {
		data, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		text := string(data)
		joined += "\n--- " + rel + " ---\n" + text
		for _, banned := range []string{
			"Codex multi_agent_v2",
			"multi_agent_v2",
			"Codex v2",
			"pre-v2",
			"send_input",
			"wait_agent(handle)",
		} {
			if strings.Contains(text, banned) {
				t.Errorf("%s contains stale or version-assuming Codex wording %q", rel, banned)
			}
		}
	}

	for _, want := range []string{
		"live tool surface",
		"not from a runtime-version label",
		"`«worker.spawn»`",
		"`spawn_agent(task_name,message,fork_turns)`",
		"`«worker-identity»`",
		"`«addressable-worker»`",
		"`send_message(target,message)`",
		"`followup_task(target,message)`",
		"`«completion-signal»`",
		"`wait_agent(timeout_ms)`",
		"queued/activity-driven",
		"`«roster-reconcile»`",
		"`list_agents(path_prefix?)`",
		"`«worker.shutdown»`",
		"Do not bless `interrupt_agent`",
		"Do not infer capabilities from a Codex version name",
		"re-run the kept-alive reviewer through the same `«addressable-worker»` binding",
		"Fresh-spawn the reviewer only when the existing reviewer is no longer addressable or reuse conditions fail",
		"Do not infer that `followup_task` is absent from the absence of earlier follow-up events",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("Codex current-runtime references missing %q", want)
		}
	}

	if _, err := os.Stat(filepath.Join(root, "skills", "first-officer", "references", "codex-v2-first-officer-runtime.md")); !os.IsNotExist(err) {
		t.Fatalf("Codex v2 must be a host-binding section in the existing Codex runtime reference, not a separate runtime file")
	}
}
