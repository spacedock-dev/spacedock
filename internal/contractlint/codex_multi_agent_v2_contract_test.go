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
		"`followup_task` starts a worker turn on the addressed worker",
		"`send_input` is both the compatibility address and turn-starting route",
		"completed-but-still-addressable worker",
		"`«completion-signal»`",
		"`wait_agent(timeout_ms)`",
		"queued/activity-driven",
		"`«roster-reconcile»`",
		"`list_agents(path_prefix?)`",
		"`«worker.shutdown»`",
		"Do not bless `interrupt_agent`",
		"Do not infer capabilities from a Codex version name",
		"re-run the kept-alive reviewer through the same `«addressable-worker»` capability",
		"reviewer already sent a completion signal",
		"Fresh-dispatch the reviewer only when the existing reviewer is no longer addressable or reuse conditions fail",
		"Do not infer that the turn-starting addressable-worker route is absent from the absence of newer follow-up tools",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("Codex current-runtime references missing %q", want)
		}
	}

	if _, err := os.Stat(filepath.Join(root, "skills", "first-officer", "references", "codex-v2-first-officer-runtime.md")); !os.IsNotExist(err) {
		t.Fatalf("Codex v2 must be a host-binding section in the existing Codex runtime reference, not a separate runtime file")
	}
}

func TestFeedbackRejectionFlowStaysHostNeutral(t *testing.T) {
	root := repoRoot(t)
	rel := filepath.Join("skills", "feedback-rejection-flow", "SKILL.md")
	data, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	text := string(data)

	for _, banned := range []string{
		"Codex",
		"Claude",
		"followup_task",
		"send_message",
		"send_input",
		"SendMessage",
		"list_agents",
		"wait_agent",
		"spawn_agent",
		"interrupt_agent",
	} {
		if strings.Contains(text, banned) {
			t.Errorf("%s contains runtime-specific host or tool name %q", rel, banned)
		}
	}

	for _, want := range []string{
		"`«addressable-worker»`",
		"`«completion-signal»`",
		"Fresh-dispatch the reviewer only when the existing reviewer is no longer addressable or reuse conditions fail",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("%s missing host-neutral capability wording %q", rel, want)
		}
	}
}

func TestCodexToolNamesStayInRuntimeBindingSection(t *testing.T) {
	root := repoRoot(t)
	codexRuntimeRel := filepath.Join("skills", "first-officer", "references", "codex-first-officer-runtime.md")
	data, err := os.ReadFile(filepath.Join(root, codexRuntimeRel))
	if err != nil {
		t.Fatalf("read %s: %v", codexRuntimeRel, err)
	}
	codexRuntime := string(data)
	probeSection, outsideProbe := extractMarkdownSection(t, codexRuntime, "## Live Tool Surface Probe")
	runtimeSection, outsideRuntime := extractMarkdownSection(t, outsideProbe, "## Runtime implementation")
	allowedSections := probeSection + "\n" + runtimeSection
	waitSection, outsideWait := extractMarkdownSection(t, outsideRuntime, "## Codex wait notes")

	for _, want := range []string{
		"`spawn_agent(task_name,message,fork_turns)`",
		"`send_message(target,message)`",
		"`followup_task(target,message)`",
		"`wait_agent(timeout_ms)`",
		"`list_agents(path_prefix?)`",
		"`interrupt_agent`",
	} {
		if !strings.Contains(allowedSections, want) {
			t.Errorf("%s Codex runtime binding/probe sections missing %q", codexRuntimeRel, want)
		}
	}

	paths := map[string]string{
		codexRuntimeRel: outsideWait,
	}
	for _, rel := range []string{
		filepath.Join("skills", "ensign", "references", "codex-ensign-runtime.md"),
		filepath.Join("skills", "feedback-rejection-flow", "SKILL.md"),
		filepath.Join("docs", "dev", "codex-idle-notification-probe.md"),
	} {
		data, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		paths[rel] = string(data)
	}

	for rel, text := range paths {
		for _, banned := range []string{
			"spawn_agent(",
			"send_message(",
			"followup_task(",
			"wait_agent(",
			"list_agents(",
			"send_input",
			"SendMessage",
			"interrupt_agent",
		} {
			if strings.Contains(text, banned) {
				t.Errorf("%s contains runtime tool outside Codex live binding section: %q", rel, banned)
			}
		}
	}

	if strings.Contains(waitSection, "spawn_agent(") ||
		strings.Contains(waitSection, "send_message(") ||
		strings.Contains(waitSection, "followup_task(") ||
		strings.Contains(waitSection, "list_agents(") ||
		strings.Contains(waitSection, "send_input") ||
		strings.Contains(waitSection, "SendMessage") ||
		strings.Contains(waitSection, "interrupt_agent") {
		t.Errorf("%s Codex wait notes may mention only wait_agent as a concrete runtime tool", codexRuntimeRel)
	}
}

func TestCodexFirstOfficerRuntimeAvoidsNegativeHostContrast(t *testing.T) {
	root := repoRoot(t)
	rel := filepath.Join("skills", "first-officer", "references", "codex-first-officer-runtime.md")
	data, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	text := string(data)

	for _, banned := range []string{
		"## Terminal Teardown",
		"Merge-and-Cleanup step",
		"10. **Teardown agents at terminal.**",
		"## Team Creation",
		"## Backstop",
		"TeamDelete",
		"TeamCreate",
		"team registry",
		"team-lead",
		"Claude specifics",
		"Do not create teams",
		"Codex declares none",
		"no-op",
	} {
		if strings.Contains(text, banned) {
			t.Errorf("%s contains negative host-contrast wording %q", rel, banned)
		}
	}

	for _, want := range []string{
		"`«worker.shutdown»` remains unresolved until probed.",
		"`«roster-reconcile»` may provide active/completed task-path reads when bound.",
		"`«addressable-worker»` may carry cooperative preservation text when present.",
		"Durable workflow state remains authoritative.",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("%s missing positive Codex capability wording %q", rel, want)
		}
	}
}

func TestCodexEnsignRuntimeAvoidsNegativeHostContrast(t *testing.T) {
	root := repoRoot(t)
	rel := filepath.Join("skills", "ensign", "references", "codex-ensign-runtime.md")
	data, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	text := string(data)

	for _, banned := range []string{
		"Claude",
		"another host",
		"Codex declares none",
		"Do not reconstruct",
		"Do not send follow-up",
	} {
		if strings.Contains(text, banned) {
			t.Errorf("%s contains negative host-contrast wording %q", rel, banned)
		}
	}

	for _, want := range []string{
		"`«context-budget»` is unavailable unless a future probe binds it.",
		"Use the generated file as the source of truth.",
		"After sending the completion signal, stop unless the FO routes more work through `«addressable-worker»`.",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("%s missing positive Codex ensign wording %q", rel, want)
		}
	}
}

func extractMarkdownSection(t *testing.T, text, heading string) (section string, remainder string) {
	t.Helper()
	start := strings.Index(text, heading)
	if start == -1 {
		t.Fatalf("missing heading %q", heading)
	}
	afterStart := start + len(heading)
	nextRel := strings.Index(text[afterStart:], "\n## ")
	if nextRel == -1 {
		return text[start:], text[:start]
	}
	end := afterStart + nextRel
	return text[start:end], text[:start] + text[end:]
}
