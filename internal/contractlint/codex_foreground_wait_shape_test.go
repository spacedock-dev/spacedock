// ABOUTME: Quarantined structural checks for Codex foreground-wait operator cues.
// ABOUTME: Keeps instruction-text reads out of behavior/integration tests.
package contractlint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCodexForegroundWaitSectionCarriesOperatorInterruptionShape(t *testing.T) {
	path := filepath.Join(repoRoot(t), "skills", "first-officer", "references", "codex-first-officer-runtime.md")
	section := markdownSubsection(t, path, "### Foreground wait")

	requireForegroundWaitLifecycleClaim(t, path, section)
}

func TestForegroundWaitLifecycleClaimRejectsTerminalMutations(t *testing.T) {
	tests := []struct {
		name    string
		section string
	}{
		{
			name: "worker is terminalized",
			section: "Before calling `wait_agent`, tell the captain that pressing Esc or otherwise " +
				"causing an operator interruption only returns control; the worker is failed, " +
				"closed, or redispatched, and this failure still retries the same handle.",
		},
		{
			name: "interruption marks worker terminal",
			section: "Before calling `wait_agent`, tell the captain that Esc or operator interruption " +
				"returns control and marks the worker failed, closes it, or redispatches it before " +
				"the next foreground wait retries the same handle after the failure.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := foregroundWaitLifecycleClaimError(tt.section); err == nil {
				t.Fatal("expected a terminal-lifecycle mutation to fail")
			}
		})
	}
}

func TestForegroundWaitLifecycleClaimAcceptsNegatedClaims(t *testing.T) {
	tests := []struct {
		name    string
		section string
	}{
		{
			name: "worker is not terminalized",
			section: "Before calling `wait_agent`, tell the captain that pressing Esc or otherwise " +
				"causing an operator interruption only returns control; the worker is not failed, " +
				"closed, or redispatched, and the next foreground wait retries the same handle.",
		},
		{
			name: "does not mark worker terminal",
			section: "Before calling `wait_agent`, tell the captain that Esc or operator interruption " +
				"returns control and does not mark the worker failed, close it, or redispatch it; " +
				"the next foreground wait retries the same handle.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := foregroundWaitLifecycleClaimError(tt.section); err != nil {
				t.Fatalf("expected negated lifecycle claim to pass: %v", err)
			}
		})
	}
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

func requireForegroundWaitLifecycleClaim(t *testing.T, path, section string) {
	t.Helper()
	if err := foregroundWaitLifecycleClaimError(section); err != nil {
		t.Errorf("%s section has unsafe foreground-wait lifecycle wording: %v", path, err)
	}
}

func foregroundWaitLifecycleClaimError(section string) error {
	normalized := strings.ToLower(strings.Join(strings.Fields(section), " "))
	for _, phrase := range []string{
		"worker is failed, closed, or redispatched",
		"worker is failed, closed, or re-dispatched",
		"marks the worker failed",
		"marks the worker as failed",
		"closes it",
		"redispatches it",
		"re-dispatches it",
	} {
		if strings.Contains(normalized, phrase) {
			return fmt.Errorf("contains affirmative terminal lifecycle claim %q", phrase)
		}
	}

	for _, phrase := range []string{
		"worker is not failed, closed, or redispatched",
		"worker is not failed, closed, or re-dispatched",
		"does not mark the worker failed, close it, or redispatch it",
		"does not mark the worker failed, close it, or re-dispatch it",
		"does not mark the worker as failed, close it, or redispatch it",
		"does not mark the worker as failed, close it, or re-dispatch it",
		"does not fail, close, or redispatch the worker",
		"does not fail, close, or re-dispatch the worker",
		"does not fail, close, or redispatch it",
		"does not fail, close, or re-dispatch it",
	} {
		if strings.Contains(normalized, phrase) {
			return nil
		}
	}

	return fmt.Errorf("missing explicit negated failed/closed/redispatched claim")
}
