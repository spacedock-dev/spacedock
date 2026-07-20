// ABOUTME: Binds the Codex spawn signature in the FO adapter to the arg shape the Go
// ABOUTME: emitter produces, and guards the shared feedback skill's host-neutrality.
package contractlint

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/dispatch"
)

// codexSpawnCallRe captures the arg list of a `spawn_agent(...)` shape in adapter
// prose; codexSpawnArgRe splits it into `name` / `name="value"` entries.
var (
	codexSpawnCallRe = regexp.MustCompile(`spawn_agent\(([^)]*)\)`)
	codexSpawnArgRe  = regexp.MustCompile(`([a-z_]+)(?:="([^"]*)")?`)
)

// codexSpawnSignatureViolations holds every `spawn_agent(...)` signature the
// adapter spells out to the arg shape the Go emitter produces: the arg-NAME set
// must equal the ToolArgs keys, and a default the signature spells out
// (`fork_turns="none"`) must equal the emitted value. The sides are independent —
// renaming a ToolArgs key or changing the emitted fork_turns value reds without
// the doc moving, and a doc edit reds without the Go moving.
func codexSpawnSignatureViolations(text string, toolArgs map[string]string) []string {
	calls := codexSpawnCallRe.FindAllStringSubmatch(text, -1)
	if len(calls) == 0 {
		return []string{"no `spawn_agent(...)` signature found — extractor bug; the binding would pass vacuously"}
	}
	wanted := map[string]bool{}
	for name := range toolArgs {
		wanted[name] = true
	}
	var out []string
	for _, call := range calls {
		names := map[string]bool{}
		for _, arg := range codexSpawnArgRe.FindAllStringSubmatch(call[1], -1) {
			names[arg[1]] = true
			if emitted, ok := toolArgs[arg[1]]; arg[2] != "" && ok && emitted != arg[2] {
				out = append(out, fmt.Sprintf("signature %q spells default %s=%q but the Go emitter produces %q", call[0], arg[1], arg[2], emitted))
			}
		}
		if !setEqual(names, wanted) {
			out = append(out, fmt.Sprintf("signature %q arg set %v does not match the Go emitter's ToolArgs keys %v; neither side may rename, add, or drop an arg without the other", call[0], sortedSet(names), sortedSet(wanted)))
		}
	}
	return out
}

// TestCodexSpawnSignatureBindsToolArgs binds the Codex FO adapter's spawn
// signature to `dispatch.CodexMultiAgentV2Spawn.ToolArgs()`, the Go surface
// Spacedock emits for a Codex spawn. The other Codex tools the adapter names
// (`send_message`, `followup_task`, `wait_agent`, `list_agents`,
// `interrupt_agent`) have no Spacedock-emitted Go source; their runtime meaning is
// owned by internal/ensigncycle/codex_live_runner_test.go.
func TestCodexSpawnSignatureBindsToolArgs(t *testing.T) {
	emitted := dispatch.CodexMultiAgentV2Spawn{TaskName: "spacedock_worker", Message: "assignment"}.ToolArgs()
	if len(emitted) == 0 {
		t.Fatal("ToolArgs() returned no args — the binding would pass vacuously")
	}
	for _, msg := range codexSpawnSignatureViolations(readRepoFile(t, codexFORuntimeRel), emitted) {
		t.Errorf("%s: %s", codexFORuntimeRel, msg)
	}
}

// hostNeutralBannedTokens are the host names and host-specific tool names a
// SHARED skill must never carry. It is a structural-absence vocabulary: it asserts
// nothing about meaning, only that the skill speaks in `«capability»` terms every
// host adapter can bind.
var hostNeutralBannedTokens = []string{
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
}

func hostNeutralViolations(text string) []string {
	var out []string
	for _, banned := range hostNeutralBannedTokens {
		if strings.Contains(text, banned) {
			out = append(out, fmt.Sprintf("contains runtime-specific host or tool name %q", banned))
		}
	}
	return out
}

// TestFeedbackRejectionFlowStaysHostNeutral keeps the shared feedback-rejection
// skill free of host and host-tool names. What its `«capability»` vocabulary MEANS
// is bound by TestRuntimeCapabilitySetAgreesAcrossCoreAndAdapters, which requires
// every capability the skill references to exist in the dispatch core.
func TestFeedbackRejectionFlowStaysHostNeutral(t *testing.T) {
	for _, msg := range hostNeutralViolations(readRepoFile(t, feedbackFlowRel)) {
		t.Errorf("%s %s", feedbackFlowRel, msg)
	}
}

// TestHostNeutralGuardDiscriminates is the non-vacuity control: capability-only
// prose PASSES, a host name and each shape of host tool name RED. An emptied
// vocabulary or a predicate that stopped comparing fails here rather than waving
// through a shared skill that has drifted host-specific.
func TestHostNeutralGuardDiscriminates(t *testing.T) {
	pass := []struct{ why, text string }{
		{"capability-only prose", "Re-run the reviewer through `«addressable-worker»` and wait for `«completion-signal»`.\n"},
		{"empty document", ""},
	}
	for _, c := range pass {
		if v := hostNeutralViolations(c.text); len(v) != 0 {
			t.Fatalf("control: the %s was wrongly flagged: %v", c.why, v)
		}
	}
	red := []struct{ why, text string }{
		{"host name (Claude)", "On Claude, keep the reviewer alive.\n"},
		{"host name (Codex)", "On Codex, keep the reviewer alive.\n"},
		{"turn-starting tool name", "Re-run the reviewer with followup_task(target,message).\n"},
		{"non-triggering tool name", "Send context with send_message(target,message).\n"},
		{"Claude team tool name", "Re-run the reviewer with SendMessage.\n"},
	}
	for _, c := range red {
		if v := hostNeutralViolations(c.text); len(v) == 0 {
			t.Fatalf("control: the %s was not flagged — the guard stopped biting", c.why)
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
