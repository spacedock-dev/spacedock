// ABOUTME: Quarantined structural checks for Codex foreground-wait operator cues.
// ABOUTME: Keeps instruction-text reads out of behavior/integration tests.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCodexForegroundWaitSectionCarriesOperatorInterruptionShape(t *testing.T) {
	path := filepath.Join(repoRoot(t), "skills", "first-officer", "references", "codex-first-officer-runtime.md")
	section := markdownSubsection(t, path, "### Foreground wait")

	requireSectionContains(t, path, section, "`wait_agent(handle)`")
	requireSectionContains(t, path, section, "Before calling `wait_agent`")
	requireSectionContains(t, path, section, "Esc")
	requireSectionContains(t, path, section, "operator interruption")
	requireSectionContains(t, path, section, "returns control")
	requireSectionContains(t, path, section, "same handle")

	for _, terminalWord := range []string{"fail", "failed", "failure", "close", "redispatch"} {
		requireSectionContains(t, path, section, terminalWord)
	}
}

func TestCodexIdleProbeForegroundWaitInterruptionIsNonTerminal(t *testing.T) {
	path := filepath.Join(repoRoot(t), "docs", "dev", "codex-idle-notification-probe.md")
	section := markdownSubsection(t, path, "## Foreground wait comparison")

	requireSectionContains(t, path, section, "`wait_agent(handle)`")
	requireSectionContains(t, path, section, "operator interruption")
	requireSectionContains(t, path, section, "returns control")
	requireSectionContains(t, path, section, "non-terminal")
	requireSectionContains(t, path, section, "same handle")
	requireSectionContains(t, path, section, "final status")
}

func markdownSubsection(t *testing.T, path, heading string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	text := string(data)
	start := strings.Index(text, heading+"\n")
	if start < 0 {
		t.Fatalf("%s: missing heading %q", path, heading)
	}
	section := text[start+len(heading)+1:]
	prefix := "#"
	if strings.HasPrefix(heading, "### ") {
		prefix = "### "
	} else if strings.HasPrefix(heading, "## ") {
		prefix = "## "
	}
	if end := strings.Index(section, "\n"+prefix); end >= 0 {
		section = section[:end]
	}
	return section
}

func requireSectionContains(t *testing.T, path, section, want string) {
	t.Helper()
	if !strings.Contains(section, want) {
		t.Errorf("%s section is missing %q", path, want)
	}
}
