package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
		"`spawn_agent(task_name,message,fork_turns=\"none\")`",
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
