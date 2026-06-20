package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCodexMultiAgentV2RuntimeReferencesUseLiveHostBindings(t *testing.T) {
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
			"send_input",
			"wait_agent(handle)",
		} {
			if strings.Contains(text, banned) {
				t.Errorf("%s contains unversioned stale Codex v2 wording %q", rel, banned)
			}
		}
	}

	for _, want := range []string{
		"Codex multi_agent_v2",
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
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("Codex v2 references missing %q", want)
		}
	}

	if _, err := os.Stat(filepath.Join(root, "skills", "first-officer", "references", "codex-v2-first-officer-runtime.md")); !os.IsNotExist(err) {
		t.Fatalf("Codex v2 must be a host-binding section in the existing Codex runtime reference, not a separate runtime file")
	}
}
