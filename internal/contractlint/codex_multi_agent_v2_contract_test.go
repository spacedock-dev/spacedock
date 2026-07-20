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

// codexSpawnCallRe captures a `spawn_agent(...)` arg list; codexSpawnArgRe matches
// ONE whole `name` / `name="value"` entry. The anchors are load-bearing: an
// unanchored scan skips unparseable text, so `spawn_agent(task_name,,message)` passed.
var (
	codexSpawnCallRe = regexp.MustCompile(`spawn_agent\(([^)]*)\)`)
	codexSpawnArgRe  = regexp.MustCompile(`^([a-z_]+)(?:="([^"]*)")?$`)
)

type codexSpawnArg struct {
	name, value string
	hasDefault  bool
}

func codexSpawnArgs(argList string) ([]codexSpawnArg, error) {
	var out []codexSpawnArg
	for _, entry := range strings.Split(argList, ",") {
		trimmed := strings.TrimSpace(entry)
		m := codexSpawnArgRe.FindStringSubmatch(trimmed)
		if m == nil {
			return nil, fmt.Errorf("unparseable argument %q", trimmed)
		}
		out = append(out, codexSpawnArg{name: m[1], value: m[2], hasDefault: strings.Contains(trimmed, "=")})
	}
	return out, nil
}

// codexSpawnSignatureViolations holds every `spawn_agent(...)` signature the
// adapter spells out to the arg shape `ToolArgs` DECLARES: the arg-NAME set must
// equal its keys, and a spelled-out default must equal its value. Either side
// moving alone reds.
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
		args, err := codexSpawnArgs(call[1])
		if err != nil {
			out = append(out, fmt.Sprintf("signature %q is malformed: %v", call[0], err))
			continue
		}
		names := map[string]bool{}
		for _, arg := range args {
			if names[arg.name] {
				out = append(out, fmt.Sprintf("signature %q declares argument %q twice", call[0], arg.name))
			}
			names[arg.name] = true
			if emitted, ok := toolArgs[arg.name]; arg.hasDefault && ok && emitted != arg.value {
				out = append(out, fmt.Sprintf("signature %q spells default %s=%q but the Go emitter produces %q", call[0], arg.name, arg.value, emitted))
			}
		}
		if !setEqual(names, wanted) {
			out = append(out, fmt.Sprintf("signature %q arg set %v does not match the Go emitter's ToolArgs keys %v; neither side may rename, add, or drop an arg without the other", call[0], sortedSet(names), sortedSet(wanted)))
		}
	}
	return out
}

// TestCodexSpawnSignatureBindsToolArgs binds the Codex FO adapter's spawn
// signature to the arg shape `dispatch.CodexMultiAgentV2Spawn.ToolArgs()`
// declares. Scope of the claim: `ToolArgs` has NO production caller, so this is
// doc-to-Go-declaration agreement, NOT evidence of what a Codex spawn puts on the
// wire. It still reds on real divergence — the two sides can move independently —
// but nothing here observes runtime behavior. The adapter's other tool names are
// annotated on codexToolTokens in runtime_binding_block_test.go.
func TestCodexSpawnSignatureBindsToolArgs(t *testing.T) {
	emitted := dispatch.CodexMultiAgentV2Spawn{TaskName: "spacedock_worker", Message: "assignment"}.ToolArgs()
	if len(emitted) == 0 {
		t.Fatal("ToolArgs() returned no args — the binding would pass vacuously")
	}
	for _, msg := range codexSpawnSignatureViolations(readRepoFile(t, codexFORuntimeRel), emitted) {
		t.Errorf("%s: %s", codexFORuntimeRel, msg)
	}
}

// hostNeutralBannedTokens are the host and host-tool names a SHARED skill must
// never carry — a structural-absence vocabulary, asserting nothing about meaning.
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
// is bound by TestRuntimeCapabilitySetAgreesAcrossCoreAndAdapters.
func TestFeedbackRejectionFlowStaysHostNeutral(t *testing.T) {
	for _, msg := range hostNeutralViolations(readRepoFile(t, feedbackFlowRel)) {
		t.Errorf("%s %s", feedbackFlowRel, msg)
	}
}

// TestHostNeutralGuardDiscriminates is the non-vacuity control: capability-only
// prose PASSES, a host name and each shape of host tool name RED.
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
